package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// LearningLogService は学習ログのビジネスロジックを提供する。
// 学習記録のCRUD操作と、ストリーク・カレンダーデータの取得を担当する。
// goalRepoが設定されている場合、ログ作成時にゴール進捗を自動更新する。
type LearningLogService struct {
	repo     repository.LearningLogRepositoryInterface
	goalRepo repository.LearningGoalRepositoryInterface
}

// NewLearningLogService は新しいLearningLogServiceインスタンスを生成する。
// goalRepoはnil可（ゴール連携なし）。
func NewLearningLogService(repo repository.LearningLogRepositoryInterface, goalRepo repository.LearningGoalRepositoryInterface) *LearningLogService {
	return &LearningLogService{repo: repo, goalRepo: goalRepo}
}

// Create は新しい学習ログを作成する。
// Duration、Category、Sourceのバリデーションを行う。
// GoalIDが指定されている場合、ゴールの進捗を自動更新する。
func (s *LearningLogService) Create(log *model.LearningLog) error {
	// Title/Content バリデーション
	if err := domain.ValidateStringLength(log.Title, 1, 200, "タイトル"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(log.Content, 0, 10000, "内容"); err != nil {
		return err
	}
	// Duration: 0以上1440以下（24時間）
	if err := validateDuration(log.Duration); err != nil {
		return err
	}
	// Category: 空文字の場合はデフォルト値を設定、それ以外は有効な値のみ許可
	if log.Category == "" {
		log.Category = model.LogCategoryOther
	} else if !model.ValidCategories[log.Category] {
		return domain.NewError(domain.ErrCodeBadRequest, msgInvalidCategory, nil)
	}
	// Source: 空文字（デフォルト"manual"）または有効な値のみ許可
	if log.Source != "" && !model.ValidSources[log.Source] {
		return domain.NewError(domain.ErrCodeBadRequest, msgInvalidSource, nil)
	}

	// ゴール紐付けバリデーション
	var goal *model.LearningGoal
	if log.GoalID != nil && s.goalRepo != nil {
		var err error
		goal, err = s.goalRepo.FindByID(*log.GoalID)
		if err != nil {
			return domain.NewError(domain.ErrCodeNotFound, "指定されたゴールが見つかりません", err)
		}
		if goal.UserID != log.UserID {
			return domain.NewError(domain.ErrCodeForbidden, "他のユーザーのゴールには紐付けできません", nil)
		}
	}

	if err := s.repo.Create(log); err != nil {
		return domain.NewError(domain.ErrCodeInternal, "学習ログの作成に失敗しました", err)
	}

	// ゴール進捗自動更新（TargetHoursが設定されている場合のみ）
	if goal != nil && goal.TargetHours > 0 {
		totalMinutes, err := s.repo.SumDurationByGoalID(*log.GoalID)
		if err != nil {
			return nil // ログ作成は成功、進捗更新失敗はサイレント
		}
		goal.Progress = CalculateGoalProgressPercentage(totalMinutes, goal.TargetHours)
		if goal.Progress >= 100 && goal.Status == model.GoalStatusActive {
			goal.Status = model.GoalStatusCompleted
			now := time.Now()
			goal.CompletedAt = &now
		}
		_ = s.goalRepo.Update(goal)
	}

	return nil
}

// GetLinkedLogs は指定ゴールに紐付いた学習ログを取得する。所有権を検証する。
func (s *LearningLogService) GetLinkedLogs(goalID, userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	if s.goalRepo == nil {
		return nil, 0, domain.NewError(domain.ErrCodeBadRequest, "ゴール連携が有効ではありません", nil)
	}
	goal, err := s.goalRepo.FindByID(goalID)
	if err != nil {
		return nil, 0, domain.NewError(domain.ErrCodeNotFound, "指定されたゴールが見つかりません", err)
	}
	if goal.UserID != userID {
		return nil, 0, domain.NewError(domain.ErrCodeForbidden, "他のユーザーのゴールのログは参照できません", nil)
	}
	return s.repo.GetByGoalID(goalID, limit, offset)
}

// GetGoalProgress は指定ゴールの実績時間 vs 目標時間の進捗情報を返す。所有権を検証する。
func (s *LearningLogService) GetGoalProgress(goalID, userID uint) (*model.GoalProgress, error) {
	if s.goalRepo == nil {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "ゴール連携が有効ではありません", nil)
	}
	goal, err := s.goalRepo.FindByID(goalID)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "指定されたゴールが見つかりません", err)
	}
	if goal.UserID != userID {
		return nil, domain.NewError(domain.ErrCodeForbidden, "他のユーザーのゴール進捗は参照できません", nil)
	}

	actualMinutes, err := s.repo.SumDurationByGoalID(goalID)
	if err != nil {
		return nil, err
	}

	return &model.GoalProgress{
		GoalID:        goalID,
		TargetHours:   goal.TargetHours,
		ActualMinutes: actualMinutes,
		Percentage:    CalculateGoalProgressPercentage(actualMinutes, goal.TargetHours),
	}, nil
}

// BatchCreate は複数の学習ログを一括作成する。
// 各ログのバリデーションを行い、全てパスした場合のみ一括保存する。
// 最大50件まで。
func (s *LearningLogService) BatchCreate(userID uint, logs []model.LearningLog) ([]model.LearningLog, error) {
	if len(logs) == 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "学習ログは1件以上指定してください", nil)
	}
	if len(logs) > 50 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "学習ログは50件以下で指定してください", nil)
	}

	for i := range logs {
		logs[i].UserID = userID

		if err := validateDuration(logs[i].Duration); err != nil {
			return nil, err
		}
		if logs[i].Category == "" {
			logs[i].Category = model.LogCategoryOther
		} else if !model.ValidCategories[logs[i].Category] {
			return nil, domain.NewError(domain.ErrCodeBadRequest, msgInvalidCategory, nil)
		}
		if logs[i].Source != "" && !model.ValidSources[logs[i].Source] {
			return nil, domain.NewError(domain.ErrCodeBadRequest, msgInvalidSource, nil)
		}
	}

	if err := s.repo.CreateBatch(logs); err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "学習ログの一括作成に失敗しました", err)
	}

	return logs, nil
}

// GetByID は指定IDの学習ログを取得する。所有権を検証する。
func (s *LearningLogService) GetByID(id, userID uint) (*model.LearningLog, error) {
	return s.findAndCheckOwnership(id, userID)
}

// GetByUserID は指定ユーザーの学習ログをページネーション付きで取得する。
func (s *LearningLogService) GetByUserID(userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	return s.repo.GetByUserID(userID, limit, offset)
}

// GetByCategory は指定ユーザーの学習ログをカテゴリでフィルタリングして取得する。
func (s *LearningLogService) GetByCategory(userID uint, category string) ([]model.LearningLog, error) {
	if !model.ValidCategories[model.LogCategory(category)] {
		return nil, domain.NewError(domain.ErrCodeBadRequest, msgInvalidCategory, nil)
	}
	return s.repo.GetByCategory(userID, category)
}

// GetBySource は指定ユーザーの学習ログをソース（manual/pomodoro）でフィルタリングして取得する。
func (s *LearningLogService) GetBySource(userID uint, source string) ([]model.LearningLog, error) {
	if !model.ValidSources[model.LogSource(source)] {
		return nil, domain.NewError(domain.ErrCodeBadRequest, msgInvalidSource, nil)
	}
	return s.repo.GetBySource(userID, source)
}

// CalculateGoalProgressPercentage はゴールの進捗率（0〜100）を算出する純粋関数。
// totalMinutes: 実績の学習時間（分）、targetHours: 目標時間（時間）。
// targetHoursが0以下の場合は0を返す。
func CalculateGoalProgressPercentage(totalMinutes, targetHours int) int {
	if targetHours <= 0 {
		return 0
	}
	progress := totalMinutes * 100 / (targetHours * 60)
	if progress > 100 {
		progress = 100
	}
	return progress
}

// validateDuration は学習時間（分）が有効な範囲（0〜1440）かを検証する。
func validateDuration(duration int) error {
	if duration < 0 || duration > 1440 {
		return domain.NewError(domain.ErrCodeBadRequest, "学習時間は0〜1440分の範囲で指定してください", nil)
	}
	return nil
}

// findAndCheckOwnership は学習ログを取得し、指定ユーザーが所有者かを検証する。
func (s *LearningLogService) findAndCheckOwnership(id, userID uint) (*model.LearningLog, error) {
	return checkOwnership(s.repo.FindByID, id, userID, func(l *model.LearningLog) uint { return l.UserID })
}

// Update は所有権を検証した後、学習ログを更新する。
func (s *LearningLogService) Update(id, userID uint, updates *model.LearningLog) (*model.LearningLog, error) {
	log, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if title := strings.TrimSpace(updates.Title); title != "" {
		if err := domain.ValidateStringLength(title, 1, 200, "タイトル"); err != nil {
			return nil, err
		}
		log.Title = title
	}
	if content := strings.TrimSpace(updates.Content); content != "" {
		if err := domain.ValidateStringLength(content, 1, 10000, "内容"); err != nil {
			return nil, err
		}
		log.Content = content
	}
	if cat := strings.TrimSpace(string(updates.Category)); cat != "" {
		log.Category = updates.Category
	}
	if updates.Duration != 0 {
		if err := validateDuration(updates.Duration); err != nil {
			return nil, err
		}
		log.Duration = updates.Duration
	}

	if err := s.repo.Update(log); err != nil {
		return nil, err
	}
	return log, nil
}

// FavoriteLog は所有権を検証した後、学習ログをお気に入りに設定する。
func (s *LearningLogService) FavoriteLog(id, userID uint) error {
	log, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return err
	}
	log.IsFavorite = true
	return s.repo.Update(log)
}

// UnfavoriteLog は所有権を検証した後、学習ログのお気に入りを解除する。
func (s *LearningLogService) UnfavoriteLog(id, userID uint) error {
	log, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return err
	}
	log.IsFavorite = false
	return s.repo.Update(log)
}

// Delete は所有権を検証した後、学習ログを削除する。
func (s *LearningLogService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id, userID)
}

// GetWeeklyDuration は指定ユーザーの過去7日間の学習時間合計（分）を返す。
func (s *LearningLogService) GetWeeklyDuration(userID uint) (int, error) {
	return s.repo.SumDurationByPeriod(userID, 7)
}

// GetStreakInfo は指定ユーザーの学習ストリーク情報を取得する。
func (s *LearningLogService) GetStreakInfo(userID uint) (*model.StreakInfo, error) {
	return s.repo.GetStreakInfo(userID)
}

// GetCalendarData はカレンダー表示用の日別学習ログ集計データを取得する。
func (s *LearningLogService) GetCalendarData(userID uint) ([]model.CalendarEntry, error) {
	return s.repo.GetCalendarData(userID)
}

// GetRecentCategories はユーザーの最近よく使うカテゴリを頻度順で返す。
func (s *LearningLogService) GetRecentCategories(userID uint) ([]string, error) {
	return s.repo.GetRecentCategories(userID, 5)
}

// GetFavorites はユーザーのお気に入り学習ログ一覧を取得する。
func (s *LearningLogService) GetFavorites(userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	return s.repo.GetFavorites(userID, limit, offset)
}

// GetMonthlySummary はユーザーの月別学習サマリー（直近Nヶ月）を取得する。
func (s *LearningLogService) GetMonthlySummary(userID uint, months int) ([]model.MonthlySummary, error) {
	if months < 1 || months > 24 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "monthsは1〜24の範囲で指定してください", nil)
	}
	return s.repo.GetMonthlySummary(userID, months)
}

// ExportCSV は指定ユーザーの学習ログをCSV形式でエクスポートする。
// days: 取得する過去の日数（0は全期間）。負の値はバリデーションエラー。
func (s *LearningLogService) ExportCSV(userID uint, days int) ([]byte, error) {
	if days < 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "期間は0以上の値を指定してください", nil)
	}

	logs, err := s.repo.GetByPeriod(userID, days)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// ヘッダー行（BOM付きUTF-8でExcel互換）
	buf.WriteString("\xef\xbb\xbf")
	if err := w.Write([]string{"日付", "カテゴリ", "タイトル", "学習時間(分)", "メモ"}); err != nil {
		return nil, err
	}

	for _, log := range logs {
		row := []string{
			log.CreatedAt.Format("2006-01-02"),
			string(log.Category),
			log.Title,
			fmt.Sprintf("%d", log.Duration),
			log.Content,
		}
		if err := w.Write(row); err != nil {
			return nil, err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// jsonLogEntry はJSONエクスポート用の構造体。
type jsonLogEntry struct {
	Date     string `json:"date"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Duration int    `json:"duration"`
	Content  string `json:"content"`
}

// ExportJSON は指定ユーザーの学習ログをJSON形式でエクスポートする。
// days: 取得する過去の日数（0は全期間）。負の値はバリデーションエラー。
func (s *LearningLogService) ExportJSON(userID uint, days int) ([]byte, error) {
	if days < 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "期間は0以上の値を指定してください", nil)
	}

	logs, err := s.repo.GetByPeriod(userID, days)
	if err != nil {
		return nil, err
	}

	entries := make([]jsonLogEntry, len(logs))
	for i, log := range logs {
		entries[i] = jsonLogEntry{
			Date:     log.CreatedAt.Format("2006-01-02"),
			Title:    log.Title,
			Category: string(log.Category),
			Duration: log.Duration,
			Content:  log.Content,
		}
	}

	return json.Marshal(entries)
}

// ImportCSV はCSVデータから学習ログを一括インポートする。
// エクスポート形式（BOM付きUTF-8、ヘッダー: 日付,カテゴリ,タイトル,学習時間(分),メモ）に対応する。
// 最大50件まで。各行のバリデーションを行い、全てパスした場合のみ保存する。
func (s *LearningLogService) ImportCSV(userID uint, data []byte) ([]model.LearningLog, error) {
	// BOMを除去
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

	reader := csv.NewReader(io.NopCloser(bytes.NewReader(data)))
	reader.TrimLeadingSpace = true

	// ヘッダー行を読み飛ばす
	header, err := reader.Read()
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "CSVファイルの読み取りに失敗しました", err)
	}
	if len(header) < 5 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "CSVのカラム数が不足しています（5列必要）", nil)
	}

	var logs []model.LearningLog
	lineNum := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, domain.NewError(domain.ErrCodeBadRequest, fmt.Sprintf("CSV %d行目の読み取りに失敗しました", lineNum+1), err)
		}
		lineNum++

		if len(record) < 5 {
			return nil, domain.NewError(domain.ErrCodeBadRequest, fmt.Sprintf("CSV %d行目のカラム数が不足しています", lineNum), nil)
		}

		// 日付パース
		dateStr := strings.TrimSpace(record[0])
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, domain.NewError(domain.ErrCodeBadRequest, fmt.Sprintf("CSV %d行目の日付形式が不正です: %s", lineNum, dateStr), err)
		}

		// カテゴリ
		category := model.LogCategory(strings.TrimSpace(record[1]))
		if category != "" && !model.ValidCategories[category] {
			return nil, domain.NewError(domain.ErrCodeBadRequest, fmt.Sprintf("CSV %d行目の無効なカテゴリです: %s", lineNum, category), nil)
		}
		if category == "" {
			category = model.LogCategoryOther
		}

		// タイトル
		title := strings.TrimSpace(record[2])
		if title == "" {
			return nil, domain.NewError(domain.ErrCodeBadRequest, fmt.Sprintf("CSV %d行目のタイトルが空です", lineNum), nil)
		}

		// 学習時間
		durationStr := strings.TrimSpace(record[3])
		duration, err := strconv.Atoi(durationStr)
		if err != nil {
			return nil, domain.NewError(domain.ErrCodeBadRequest, fmt.Sprintf("CSV %d行目の学習時間が不正です: %s", lineNum, durationStr), err)
		}
		if err := validateDuration(duration); err != nil {
			return nil, domain.NewError(domain.ErrCodeBadRequest, fmt.Sprintf("CSV %d行目: 学習時間は0〜1440分の範囲で指定してください", lineNum), nil)
		}

		// メモ
		content := strings.TrimSpace(record[4])
		if content == "" {
			content = title // メモが空の場合はタイトルをデフォルトに
		}

		logs = append(logs, model.LearningLog{
			UserID:    userID,
			Title:     title,
			Content:   content,
			Category:  category,
			Duration:  duration,
			Source:    model.LogSourceManual,
			CreatedAt: parsedDate,
		})
	}

	if len(logs) == 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "インポートするデータがありません", nil)
	}
	if len(logs) > 50 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "一度にインポートできるのは50件までです", nil)
	}

	if err := s.repo.CreateBatch(logs); err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "学習ログのインポートに失敗しました", err)
	}

	return logs, nil
}

// CountByUserID は指定ユーザーの学習ログ総数を返す。
func (s *LearningLogService) CountByUserID(userID uint) (int64, error) {
	return s.repo.CountByUserID(userID)
}

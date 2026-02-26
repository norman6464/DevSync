package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// LearningLogService は学習ログのビジネスロジックを提供する。
// 学習記録のCRUD操作と、ストリーク・カレンダーデータの取得を担当する。
type LearningLogService struct {
	repo repository.LearningLogRepositoryInterface
}

// NewLearningLogService は新しいLearningLogServiceインスタンスを生成する。
func NewLearningLogService(repo repository.LearningLogRepositoryInterface) *LearningLogService {
	return &LearningLogService{repo: repo}
}

// Create は新しい学習ログを作成する。
// Duration、Category、Sourceのバリデーションを行う。
func (s *LearningLogService) Create(log *model.LearningLog) error {
	// Duration: 0以上1440以下（24時間）
	if err := validateDuration(log.Duration); err != nil {
		return err
	}
	// Category: 空文字（デフォルト適用）または有効な値のみ許可
	if log.Category != "" && !model.ValidCategories[log.Category] {
		return ErrBadRequest
	}
	// Source: 空文字（デフォルト"manual"）または有効な値のみ許可
	if log.Source != "" && !model.ValidSources[log.Source] {
		return ErrBadRequest
	}
	return s.repo.Create(log)
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
		if logs[i].Category != "" && !model.ValidCategories[logs[i].Category] {
			return nil, ErrBadRequest
		}
		if logs[i].Source != "" && !model.ValidSources[logs[i].Source] {
			return nil, ErrBadRequest
		}
	}

	if err := s.repo.CreateBatch(logs); err != nil {
		return nil, err
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
		return nil, domain.NewError(domain.ErrCodeBadRequest, "無効なカテゴリです", nil)
	}
	return s.repo.GetByCategory(userID, category)
}

// GetBySource は指定ユーザーの学習ログをソース（manual/pomodoro）でフィルタリングして取得する。
func (s *LearningLogService) GetBySource(userID uint, source string) ([]model.LearningLog, error) {
	if !model.ValidSources[model.LogSource(source)] {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "無効なソースです", nil)
	}
	return s.repo.GetBySource(userID, source)
}

// validateDuration は学習時間（分）が有効な範囲（0〜1440）かを検証する。
func validateDuration(duration int) error {
	if duration < 0 || duration > 1440 {
		return ErrBadRequest
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

	if strings.TrimSpace(updates.Title) != "" {
		log.Title = updates.Title
	}
	if strings.TrimSpace(updates.Content) != "" {
		log.Content = updates.Content
	}
	if strings.TrimSpace(string(updates.Category)) != "" {
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

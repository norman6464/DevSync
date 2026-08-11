package usecase

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

const (
	// msgInvalidLogCategory は無効なカテゴリのメッセージ。
	msgInvalidLogCategory = "無効なカテゴリです"
	// msgInvalidLogSource は無効なソースのメッセージ。
	msgInvalidLogSource = "無効なソースです"
	// msgGoalNotFound はゴールを取得できなかったときのメッセージ。
	msgGoalNotFound = "指定されたゴールが見つかりません"
	// msgGoalLinkDisabled はゴール連携が注入されていないときのメッセージ。
	msgGoalLinkDisabled = "ゴール連携が有効ではありません"

	// maxLearningLogBatchSize は一括作成・インポートの上限件数。
	maxLearningLogBatchSize = 50
	// recentLearningCategoryLimit は「最近よく使うカテゴリ」の件数。
	recentLearningCategoryLimit = 5
	// weeklyLearningDurationDays は週間学習時間の集計日数。
	weeklyLearningDurationDays = 7
)

// learningLogOwnerOf は所有権チェック用に学習ログの所有者 ID を取り出す。
func learningLogOwnerOf(l *model.LearningLog) uint { return l.UserID }

// validateLearningLogDuration は学習時間（分）が 0〜1440 の範囲かを検証する。
func validateLearningLogDuration(duration int) error {
	if duration < 0 || duration > 1440 {
		return domain.NewError(domain.ErrCodeBadRequest, "学習時間は0〜1440分の範囲で指定してください", nil)
	}
	return nil
}

// normalizeLearningLogCategory は空のカテゴリを既定値で埋め、無効な値を弾く。
func normalizeLearningLogCategory(log *model.LearningLog) error {
	if log.Category == "" {
		log.Category = model.LogCategoryOther
		return nil
	}
	if !model.ValidCategories[log.Category] {
		return domain.NewError(domain.ErrCodeBadRequest, msgInvalidLogCategory, nil)
	}
	return nil
}

// validateLearningLogSource はソースを検証する。空文字は既定値扱いで許容する。
func validateLearningLogSource(source model.LogSource) error {
	if source != "" && !model.ValidSources[source] {
		return domain.NewError(domain.ErrCodeBadRequest, msgInvalidLogSource, nil)
	}
	return nil
}

// CalculateGoalProgressPercentage はゴールの進捗率（0〜100）を算出する純粋関数。
// totalMinutes は実績の学習時間（分）、targetHours は目標時間（時間）。
// targetHours が 0 以下の場合は 0 を返す。
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

// findOwnedGoal は紐付け対象のゴールを取得する。
// 取得できない場合は 404、他人のゴールの場合は forbiddenMsg で 403 を返す
// （不在と DB 障害をどちらも 404 に潰す移行前の挙動を維持している）。
func findOwnedGoal(ctx context.Context, goals repository.LearningGoalLinker, goalID, userID uint, forbiddenMsg string) (*model.LearningGoal, error) {
	goal, err := goals.FindByID(ctx, goalID)
	if err != nil || goal == nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, msgGoalNotFound, err)
	}
	if goal.UserID != userID {
		return nil, domain.NewError(domain.ErrCodeForbidden, forbiddenMsg, nil)
	}
	return goal, nil
}

// CreateLearningLogUseCase は学習ログを作成する。
// ゴールが紐付いている場合は、そのゴールの進捗も更新する。
type CreateLearningLogUseCase struct {
	logs  repository.LearningLogRepository
	goals repository.LearningGoalLinker
}

// NewCreateLearningLogUseCase は CreateLearningLogUseCase を生成する。
// goals は nil 可（ゴール連携なし）。
func NewCreateLearningLogUseCase(logs repository.LearningLogRepository, goals repository.LearningGoalLinker) *CreateLearningLogUseCase {
	return &CreateLearningLogUseCase{logs: logs, goals: goals}
}

// Execute は各項目を検証し、ゴール紐付けを確認したうえで学習ログを作成する。
func (uc *CreateLearningLogUseCase) Execute(ctx context.Context, log *model.LearningLog) error {
	if err := domain.ValidateStringLength(log.Title, 1, 200, "タイトル"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(log.Content, 0, 10000, "内容"); err != nil {
		return err
	}
	if err := validateLearningLogDuration(log.Duration); err != nil {
		return err
	}
	if err := normalizeLearningLogCategory(log); err != nil {
		return err
	}
	if err := validateLearningLogSource(log.Source); err != nil {
		return err
	}

	var goal *model.LearningGoal
	if log.GoalID != nil && uc.goals != nil {
		var err error
		goal, err = findOwnedGoal(ctx, uc.goals, *log.GoalID, log.UserID, "他のユーザーのゴールには紐付けできません")
		if err != nil {
			return err
		}
	}

	if err := uc.logs.Create(ctx, log); err != nil {
		return domain.NewError(domain.ErrCodeInternal, "学習ログの作成に失敗しました", err)
	}

	uc.refreshGoalProgress(ctx, goal, log)
	return nil
}

// refreshGoalProgress は紐付いたゴールの進捗を更新する。
// ログ作成自体は成功しているため、集計や更新に失敗しても呼び出し元へはエラーを返さない
// （移行前の挙動を維持している）。
func (uc *CreateLearningLogUseCase) refreshGoalProgress(ctx context.Context, goal *model.LearningGoal, log *model.LearningLog) {
	if goal == nil || goal.TargetHours <= 0 {
		return
	}

	totalMinutes, err := uc.logs.SumDurationByGoalID(ctx, *log.GoalID)
	if err != nil {
		return
	}

	goal.Progress = CalculateGoalProgressPercentage(totalMinutes, goal.TargetHours)
	if goal.Progress >= 100 && goal.Status == model.GoalStatusActive {
		goal.Status = model.GoalStatusCompleted
		now := time.Now()
		goal.CompletedAt = &now
	}
	_ = uc.goals.Update(ctx, goal)
}

// BatchCreateLearningLogsUseCase は学習ログを一括作成する。
type BatchCreateLearningLogsUseCase struct {
	logs repository.LearningLogRepository
}

// NewBatchCreateLearningLogsUseCase は BatchCreateLearningLogsUseCase を生成する。
func NewBatchCreateLearningLogsUseCase(logs repository.LearningLogRepository) *BatchCreateLearningLogsUseCase {
	return &BatchCreateLearningLogsUseCase{logs: logs}
}

// Execute は全件の検証が通った場合のみ一括保存する。
func (uc *BatchCreateLearningLogsUseCase) Execute(ctx context.Context, userID uint, logs []model.LearningLog) ([]model.LearningLog, error) {
	if len(logs) == 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "学習ログは1件以上指定してください", nil)
	}
	if len(logs) > maxLearningLogBatchSize {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "学習ログは50件以下で指定してください", nil)
	}

	for i := range logs {
		logs[i].UserID = userID

		if err := validateLearningLogDuration(logs[i].Duration); err != nil {
			return nil, err
		}
		if err := normalizeLearningLogCategory(&logs[i]); err != nil {
			return nil, err
		}
		if err := validateLearningLogSource(logs[i].Source); err != nil {
			return nil, err
		}
	}

	if err := uc.logs.CreateBatch(ctx, logs); err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "学習ログの一括作成に失敗しました", err)
	}
	return logs, nil
}

// GetLearningLogUseCase は学習ログを 1 件取得する。
type GetLearningLogUseCase struct {
	logs repository.LearningLogRepository
}

// NewGetLearningLogUseCase は GetLearningLogUseCase を生成する。
func NewGetLearningLogUseCase(logs repository.LearningLogRepository) *GetLearningLogUseCase {
	return &GetLearningLogUseCase{logs: logs}
}

// Execute は所有権を検証したうえで学習ログを返す。
func (uc *GetLearningLogUseCase) Execute(ctx context.Context, id, userID uint) (*model.LearningLog, error) {
	return ensureOwner(ctx, uc.logs.FindByID, id, userID, learningLogOwnerOf)
}

// ListLearningLogsUseCase はユーザーの学習ログ一覧を取得する。
type ListLearningLogsUseCase struct {
	logs repository.LearningLogRepository
}

// NewListLearningLogsUseCase は ListLearningLogsUseCase を生成する。
func NewListLearningLogsUseCase(logs repository.LearningLogRepository) *ListLearningLogsUseCase {
	return &ListLearningLogsUseCase{logs: logs}
}

// Execute は学習ログを作成日の新しい順で返す。
func (uc *ListLearningLogsUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	return uc.logs.GetByUserID(ctx, userID, limit, offset)
}

// UpdateLearningLogUseCase は学習ログを更新する。
type UpdateLearningLogUseCase struct {
	logs repository.LearningLogRepository
}

// NewUpdateLearningLogUseCase は UpdateLearningLogUseCase を生成する。
func NewUpdateLearningLogUseCase(logs repository.LearningLogRepository) *UpdateLearningLogUseCase {
	return &UpdateLearningLogUseCase{logs: logs}
}

// Execute は所有権を検証し、トリム後に空でないフィールドだけを更新する。
// 学習時間は 0 のとき変更しない。
func (uc *UpdateLearningLogUseCase) Execute(ctx context.Context, id, userID uint, updates *model.LearningLog) (*model.LearningLog, error) {
	log, err := ensureOwner(ctx, uc.logs.FindByID, id, userID, learningLogOwnerOf)
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
	if category := strings.TrimSpace(string(updates.Category)); category != "" {
		log.Category = updates.Category
	}
	if updates.Duration != 0 {
		if err := validateLearningLogDuration(updates.Duration); err != nil {
			return nil, err
		}
		log.Duration = updates.Duration
	}

	if err := uc.logs.Update(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

// DeleteLearningLogUseCase は学習ログを削除する。
type DeleteLearningLogUseCase struct {
	logs repository.LearningLogRepository
}

// NewDeleteLearningLogUseCase は DeleteLearningLogUseCase を生成する。
func NewDeleteLearningLogUseCase(logs repository.LearningLogRepository) *DeleteLearningLogUseCase {
	return &DeleteLearningLogUseCase{logs: logs}
}

// Execute は所有権を検証したうえで学習ログを削除する。
func (uc *DeleteLearningLogUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensureOwner(ctx, uc.logs.FindByID, id, userID, learningLogOwnerOf); err != nil {
		return err
	}
	return uc.logs.Delete(ctx, id, userID)
}

// FavoriteLearningLogUseCase は学習ログをお気に入りに設定する。
type FavoriteLearningLogUseCase struct {
	logs repository.LearningLogRepository
}

// NewFavoriteLearningLogUseCase は FavoriteLearningLogUseCase を生成する。
func NewFavoriteLearningLogUseCase(logs repository.LearningLogRepository) *FavoriteLearningLogUseCase {
	return &FavoriteLearningLogUseCase{logs: logs}
}

// Execute は所有権を検証したうえでお気に入りに設定する。
func (uc *FavoriteLearningLogUseCase) Execute(ctx context.Context, id, userID uint) error {
	return setLearningLogFavorite(ctx, uc.logs, id, userID, true)
}

// UnfavoriteLearningLogUseCase は学習ログのお気に入りを解除する。
type UnfavoriteLearningLogUseCase struct {
	logs repository.LearningLogRepository
}

// NewUnfavoriteLearningLogUseCase は UnfavoriteLearningLogUseCase を生成する。
func NewUnfavoriteLearningLogUseCase(logs repository.LearningLogRepository) *UnfavoriteLearningLogUseCase {
	return &UnfavoriteLearningLogUseCase{logs: logs}
}

// Execute は所有権を検証したうえでお気に入りを解除する。
func (uc *UnfavoriteLearningLogUseCase) Execute(ctx context.Context, id, userID uint) error {
	return setLearningLogFavorite(ctx, uc.logs, id, userID, false)
}

// setLearningLogFavorite はお気に入り状態を書き換える共通処理。
func setLearningLogFavorite(ctx context.Context, logs repository.LearningLogRepository, id, userID uint, favorite bool) error {
	log, err := ensureOwner(ctx, logs.FindByID, id, userID, learningLogOwnerOf)
	if err != nil {
		return err
	}
	log.IsFavorite = favorite
	return logs.Update(ctx, log)
}

// ListFavoriteLearningLogsUseCase はお気に入りの学習ログ一覧を取得する。
type ListFavoriteLearningLogsUseCase struct {
	logs repository.LearningLogRepository
}

// NewListFavoriteLearningLogsUseCase は ListFavoriteLearningLogsUseCase を生成する。
func NewListFavoriteLearningLogsUseCase(logs repository.LearningLogRepository) *ListFavoriteLearningLogsUseCase {
	return &ListFavoriteLearningLogsUseCase{logs: logs}
}

// Execute はお気に入りの学習ログを作成日の新しい順で返す。
func (uc *ListFavoriteLearningLogsUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	return uc.logs.GetFavorites(ctx, userID, limit, offset)
}

// ListLearningLogsByCategoryUseCase はカテゴリで絞り込んだ学習ログを取得する。
type ListLearningLogsByCategoryUseCase struct {
	logs repository.LearningLogRepository
}

// NewListLearningLogsByCategoryUseCase は ListLearningLogsByCategoryUseCase を生成する。
func NewListLearningLogsByCategoryUseCase(logs repository.LearningLogRepository) *ListLearningLogsByCategoryUseCase {
	return &ListLearningLogsByCategoryUseCase{logs: logs}
}

// Execute はカテゴリを検証したうえで該当する学習ログを返す。
func (uc *ListLearningLogsByCategoryUseCase) Execute(ctx context.Context, userID uint, category string) ([]model.LearningLog, error) {
	if !model.ValidCategories[model.LogCategory(category)] {
		return nil, domain.NewError(domain.ErrCodeBadRequest, msgInvalidLogCategory, nil)
	}
	return uc.logs.GetByCategory(ctx, userID, category)
}

// ListLearningLogsBySourceUseCase はソースで絞り込んだ学習ログを取得する。
type ListLearningLogsBySourceUseCase struct {
	logs repository.LearningLogRepository
}

// NewListLearningLogsBySourceUseCase は ListLearningLogsBySourceUseCase を生成する。
func NewListLearningLogsBySourceUseCase(logs repository.LearningLogRepository) *ListLearningLogsBySourceUseCase {
	return &ListLearningLogsBySourceUseCase{logs: logs}
}

// Execute はソースを検証したうえで該当する学習ログを返す。
func (uc *ListLearningLogsBySourceUseCase) Execute(ctx context.Context, userID uint, source string) ([]model.LearningLog, error) {
	if !model.ValidSources[model.LogSource(source)] {
		return nil, domain.NewError(domain.ErrCodeBadRequest, msgInvalidLogSource, nil)
	}
	return uc.logs.GetBySource(ctx, userID, source)
}

// GetWeeklyLearningDurationUseCase は過去 7 日間の学習時間合計を取得する。
type GetWeeklyLearningDurationUseCase struct {
	logs repository.LearningLogRepository
}

// NewGetWeeklyLearningDurationUseCase は GetWeeklyLearningDurationUseCase を生成する。
func NewGetWeeklyLearningDurationUseCase(logs repository.LearningLogRepository) *GetWeeklyLearningDurationUseCase {
	return &GetWeeklyLearningDurationUseCase{logs: logs}
}

// Execute は過去 7 日間の学習時間合計（分）を返す。
func (uc *GetWeeklyLearningDurationUseCase) Execute(ctx context.Context, userID uint) (int, error) {
	return uc.logs.SumDurationByPeriod(ctx, userID, weeklyLearningDurationDays)
}

// GetLearningStreakUseCase は学習ストリーク情報を取得する。
type GetLearningStreakUseCase struct {
	logs repository.LearningLogRepository
}

// NewGetLearningStreakUseCase は GetLearningStreakUseCase を生成する。
func NewGetLearningStreakUseCase(logs repository.LearningLogRepository) *GetLearningStreakUseCase {
	return &GetLearningStreakUseCase{logs: logs}
}

// Execute は連続学習日数などのストリーク情報を返す。
func (uc *GetLearningStreakUseCase) Execute(ctx context.Context, userID uint) (*model.StreakInfo, error) {
	return uc.logs.GetStreakInfo(ctx, userID)
}

// GetLearningCalendarUseCase はカレンダー表示用の日別集計を取得する。
type GetLearningCalendarUseCase struct {
	logs repository.LearningLogRepository
}

// NewGetLearningCalendarUseCase は GetLearningCalendarUseCase を生成する。
func NewGetLearningCalendarUseCase(logs repository.LearningLogRepository) *GetLearningCalendarUseCase {
	return &GetLearningCalendarUseCase{logs: logs}
}

// Execute は日別の学習ログ件数を返す。
func (uc *GetLearningCalendarUseCase) Execute(ctx context.Context, userID uint) ([]model.CalendarEntry, error) {
	return uc.logs.GetCalendarData(ctx, userID)
}

// ListRecentLearningCategoriesUseCase は最近よく使うカテゴリを取得する。
type ListRecentLearningCategoriesUseCase struct {
	logs repository.LearningLogRepository
}

// NewListRecentLearningCategoriesUseCase は ListRecentLearningCategoriesUseCase を生成する。
func NewListRecentLearningCategoriesUseCase(logs repository.LearningLogRepository) *ListRecentLearningCategoriesUseCase {
	return &ListRecentLearningCategoriesUseCase{logs: logs}
}

// Execute は使用回数の多い順にカテゴリを 5 件返す。
func (uc *ListRecentLearningCategoriesUseCase) Execute(ctx context.Context, userID uint) ([]string, error) {
	return uc.logs.GetRecentCategories(ctx, userID, recentLearningCategoryLimit)
}

// GetLearningLogMonthlySummaryUseCase は月別の学習サマリーを取得する。
type GetLearningLogMonthlySummaryUseCase struct {
	logs repository.LearningLogRepository
}

// NewGetLearningLogMonthlySummaryUseCase は GetLearningLogMonthlySummaryUseCase を生成する。
func NewGetLearningLogMonthlySummaryUseCase(logs repository.LearningLogRepository) *GetLearningLogMonthlySummaryUseCase {
	return &GetLearningLogMonthlySummaryUseCase{logs: logs}
}

// Execute は直近 months ヶ月の月別サマリーを返す。months は 1〜24 の範囲。
func (uc *GetLearningLogMonthlySummaryUseCase) Execute(ctx context.Context, userID uint, months int) ([]model.MonthlySummary, error) {
	if months < 1 || months > 24 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "monthsは1〜24の範囲で指定してください", nil)
	}
	return uc.logs.GetMonthlySummary(ctx, userID, months)
}

// CountLearningLogsUseCase は学習ログの総数を取得する。
type CountLearningLogsUseCase struct {
	logs repository.LearningLogRepository
}

// NewCountLearningLogsUseCase は CountLearningLogsUseCase を生成する。
func NewCountLearningLogsUseCase(logs repository.LearningLogRepository) *CountLearningLogsUseCase {
	return &CountLearningLogsUseCase{logs: logs}
}

// Execute は学習ログの総数を返す。
func (uc *CountLearningLogsUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.logs.CountByUserID(ctx, userID)
}

// ListGoalLinkedLogsUseCase は指定ゴールに紐付いた学習ログ一覧を取得する。
type ListGoalLinkedLogsUseCase struct {
	logs  repository.LearningLogRepository
	goals repository.LearningGoalLinker
}

// NewListGoalLinkedLogsUseCase は ListGoalLinkedLogsUseCase を生成する。
func NewListGoalLinkedLogsUseCase(logs repository.LearningLogRepository, goals repository.LearningGoalLinker) *ListGoalLinkedLogsUseCase {
	return &ListGoalLinkedLogsUseCase{logs: logs, goals: goals}
}

// Execute はゴールの所有権を検証したうえで紐付いた学習ログを返す。
func (uc *ListGoalLinkedLogsUseCase) Execute(ctx context.Context, goalID, userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	if uc.goals == nil {
		return nil, 0, domain.NewError(domain.ErrCodeBadRequest, msgGoalLinkDisabled, nil)
	}
	if _, err := findOwnedGoal(ctx, uc.goals, goalID, userID, "他のユーザーのゴールのログは参照できません"); err != nil {
		return nil, 0, err
	}
	return uc.logs.GetByGoalID(ctx, goalID, limit, offset)
}

// GetGoalProgressUseCase は指定ゴールの実績時間と目標時間の進捗を取得する。
type GetGoalProgressUseCase struct {
	logs  repository.LearningLogRepository
	goals repository.LearningGoalLinker
}

// NewGetGoalProgressUseCase は GetGoalProgressUseCase を生成する。
func NewGetGoalProgressUseCase(logs repository.LearningLogRepository, goals repository.LearningGoalLinker) *GetGoalProgressUseCase {
	return &GetGoalProgressUseCase{logs: logs, goals: goals}
}

// Execute はゴールの所有権を検証したうえで進捗情報を返す。
func (uc *GetGoalProgressUseCase) Execute(ctx context.Context, goalID, userID uint) (*model.GoalProgress, error) {
	if uc.goals == nil {
		return nil, domain.NewError(domain.ErrCodeBadRequest, msgGoalLinkDisabled, nil)
	}
	goal, err := findOwnedGoal(ctx, uc.goals, goalID, userID, "他のユーザーのゴール進捗は参照できません")
	if err != nil {
		return nil, err
	}

	actualMinutes, err := uc.logs.SumDurationByGoalID(ctx, goalID)
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

// ExportLearningLogsCSVUseCase は学習ログを CSV 形式で書き出す。
type ExportLearningLogsCSVUseCase struct {
	logs repository.LearningLogRepository
}

// NewExportLearningLogsCSVUseCase は ExportLearningLogsCSVUseCase を生成する。
func NewExportLearningLogsCSVUseCase(logs repository.LearningLogRepository) *ExportLearningLogsCSVUseCase {
	return &ExportLearningLogsCSVUseCase{logs: logs}
}

// Execute は直近 days 日の学習ログを CSV（BOM 付き UTF-8）で返す。days が 0 なら全期間。
func (uc *ExportLearningLogsCSVUseCase) Execute(ctx context.Context, userID uint, days int) ([]byte, error) {
	logs, err := fetchLogsForExport(ctx, uc.logs, userID, days)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// Excel で開いたときに文字化けしないよう BOM を先頭に付ける。
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

// jsonLearningLogEntry は JSON エクスポートの 1 件分。
type jsonLearningLogEntry struct {
	Date     string `json:"date"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Duration int    `json:"duration"`
	Content  string `json:"content"`
}

// ExportLearningLogsJSONUseCase は学習ログを JSON 形式で書き出す。
type ExportLearningLogsJSONUseCase struct {
	logs repository.LearningLogRepository
}

// NewExportLearningLogsJSONUseCase は ExportLearningLogsJSONUseCase を生成する。
func NewExportLearningLogsJSONUseCase(logs repository.LearningLogRepository) *ExportLearningLogsJSONUseCase {
	return &ExportLearningLogsJSONUseCase{logs: logs}
}

// Execute は直近 days 日の学習ログを JSON で返す。days が 0 なら全期間。
func (uc *ExportLearningLogsJSONUseCase) Execute(ctx context.Context, userID uint, days int) ([]byte, error) {
	logs, err := fetchLogsForExport(ctx, uc.logs, userID, days)
	if err != nil {
		return nil, err
	}

	entries := make([]jsonLearningLogEntry, len(logs))
	for i, log := range logs {
		entries[i] = jsonLearningLogEntry{
			Date:     log.CreatedAt.Format("2006-01-02"),
			Title:    log.Title,
			Category: string(log.Category),
			Duration: log.Duration,
			Content:  log.Content,
		}
	}
	return json.Marshal(entries)
}

// fetchLogsForExport はエクスポート対象の学習ログを取得する共通処理。
func fetchLogsForExport(ctx context.Context, logs repository.LearningLogRepository, userID uint, days int) ([]model.LearningLog, error) {
	if days < 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "期間は0以上の値を指定してください", nil)
	}
	return logs.GetByPeriod(ctx, userID, days)
}

// ImportLearningLogsCSVUseCase は CSV から学習ログを一括インポートする。
type ImportLearningLogsCSVUseCase struct {
	logs repository.LearningLogRepository
}

// NewImportLearningLogsCSVUseCase は ImportLearningLogsCSVUseCase を生成する。
func NewImportLearningLogsCSVUseCase(logs repository.LearningLogRepository) *ImportLearningLogsCSVUseCase {
	return &ImportLearningLogsCSVUseCase{logs: logs}
}

// Execute はエクスポートと同じ書式の CSV を読み取り、全行の検証が通った場合のみ保存する。
func (uc *ImportLearningLogsCSVUseCase) Execute(ctx context.Context, userID uint, data []byte) ([]model.LearningLog, error) {
	logs, err := parseLearningLogsCSV(userID, data)
	if err != nil {
		return nil, err
	}

	if err := uc.logs.CreateBatch(ctx, logs); err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "学習ログのインポートに失敗しました", err)
	}
	return logs, nil
}

// parseLearningLogsCSV は CSV を学習ログへ変換する。行番号つきのエラーを返す。
func parseLearningLogsCSV(userID uint, data []byte) ([]model.LearningLog, error) {
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))

	reader := csv.NewReader(io.NopCloser(bytes.NewReader(data)))
	reader.TrimLeadingSpace = true

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

		dateStr := strings.TrimSpace(record[0])
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, domain.NewError(domain.ErrCodeBadRequest, fmt.Sprintf("CSV %d行目の日付形式が不正です: %s", lineNum, dateStr), err)
		}

		category := model.LogCategory(strings.TrimSpace(record[1]))
		if category != "" && !model.ValidCategories[category] {
			return nil, domain.NewError(domain.ErrCodeBadRequest, fmt.Sprintf("CSV %d行目の無効なカテゴリです: %s", lineNum, category), nil)
		}
		if category == "" {
			category = model.LogCategoryOther
		}

		title := strings.TrimSpace(record[2])
		if title == "" {
			return nil, domain.NewError(domain.ErrCodeBadRequest, fmt.Sprintf("CSV %d行目のタイトルが空です", lineNum), nil)
		}

		durationStr := strings.TrimSpace(record[3])
		duration, err := strconv.Atoi(durationStr)
		if err != nil {
			return nil, domain.NewError(domain.ErrCodeBadRequest, fmt.Sprintf("CSV %d行目の学習時間が不正です: %s", lineNum, durationStr), err)
		}
		if err := validateLearningLogDuration(duration); err != nil {
			return nil, domain.NewError(domain.ErrCodeBadRequest, fmt.Sprintf("CSV %d行目: 学習時間は0〜1440分の範囲で指定してください", lineNum), nil)
		}

		// メモが空のときはタイトルで補う。
		content := strings.TrimSpace(record[4])
		if content == "" {
			content = title
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
	if len(logs) > maxLearningLogBatchSize {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "一度にインポートできるのは50件までです", nil)
	}
	return logs, nil
}

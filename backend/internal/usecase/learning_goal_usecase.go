package usecase

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// msgInvalidGoalCategory は無効なカテゴリのメッセージ。
// 旧 service/errors.go の msgInvalidCategory と同一文言（あちらは未移行スライスが使うため残している）。
const msgInvalidGoalCategory = "無効なカテゴリです"

// validGoalCategories は有効な目標カテゴリの集合。
var validGoalCategories = map[string]bool{
	string(model.GoalCategoryLanguage):  true,
	string(model.GoalCategoryFramework): true,
	string(model.GoalCategorySkill):     true,
	string(model.GoalCategoryProject):   true,
	string(model.GoalCategoryOther):     true,
}

// validGoalStatuses は有効な目標ステータスの集合。
var validGoalStatuses = map[string]bool{
	string(model.GoalStatusActive):    true,
	string(model.GoalStatusCompleted): true,
	string(model.GoalStatusPaused):    true,
}

// learningGoalOwnerOf は学習目標の所有者 ID を返す。
func learningGoalOwnerOf(g *model.LearningGoal) uint { return g.UserID }

// DeadlineStatus はアクティブな目標のデッドライン状態を判定する純粋関数。
// "overdue"（期限超過）、"approaching"（3日以内）、""（安全/対象外）を返す。
func DeadlineStatus(goal *model.LearningGoal, now time.Time) string {
	if goal.TargetDate == nil || goal.Status != model.GoalStatusActive {
		return ""
	}
	days := int(goal.TargetDate.Sub(now).Hours() / 24)
	if goal.TargetDate.Before(now) && days < 0 {
		return "overdue"
	}
	if days <= 3 {
		return "approaching"
	}
	return ""
}

// DaysUntilDeadline は目標の期限までの残り日数を返す。
// TargetDate が nil または期限超過の場合は -1 を返す。
func DaysUntilDeadline(goal *model.LearningGoal, now time.Time) int {
	if goal.TargetDate == nil {
		return -1
	}
	days := int(goal.TargetDate.Sub(now).Hours() / 24)
	if days < 0 {
		return -1
	}
	return days
}

// CalculateGoalForecast は目標達成予測を算出する純粋関数。
// dailyAverageMinutes は過去14日間の日平均学習時間（分）。
func CalculateGoalForecast(goal *model.LearningGoal, actualMinutes, dailyAverageMinutes int, now time.Time) *model.GoalForecast {
	forecast := &model.GoalForecast{
		GoalID:              goal.ID,
		Title:               goal.Title,
		CurrentProgress:     goal.Progress,
		TargetHours:         goal.TargetHours,
		ActualMinutes:       actualMinutes,
		DailyAverageMinutes: dailyAverageMinutes,
		EstimatedDaysLeft:   -1,
		DaysUntilDeadline:   -1,
		Difficulty:          "unknown",
	}

	if goal.TargetDate != nil {
		forecast.DaysUntilDeadline = int(math.Ceil(goal.TargetDate.Sub(now).Hours() / 24))
		if forecast.DaysUntilDeadline < 0 {
			forecast.DaysUntilDeadline = 0
		}
	}

	// 目標時間が設定されている場合、残り時間から予測日数を算出
	if goal.TargetHours > 0 && dailyAverageMinutes > 0 {
		remainingMinutes := goal.TargetHours*60 - actualMinutes
		if remainingMinutes <= 0 {
			forecast.EstimatedDaysLeft = 0
			forecast.OnTrack = true
			forecast.Difficulty = "easy"
		} else {
			forecast.EstimatedDaysLeft = int(math.Ceil(float64(remainingMinutes) / float64(dailyAverageMinutes)))
			if forecast.DaysUntilDeadline >= 0 {
				forecast.OnTrack = forecast.EstimatedDaysLeft <= forecast.DaysUntilDeadline
				ratio := float64(forecast.EstimatedDaysLeft) / float64(max(forecast.DaysUntilDeadline, 1))
				switch {
				case ratio <= 0.5:
					forecast.Difficulty = "easy"
				case ratio <= 1.0:
					forecast.Difficulty = "medium"
				default:
					forecast.Difficulty = "hard"
				}
			}
		}
	} else if goal.TargetHours > 0 && dailyAverageMinutes == 0 {
		forecast.Difficulty = "hard"
	}

	return forecast
}

// CreateLearningGoalUseCase は学習目標を作成する。
type CreateLearningGoalUseCase struct {
	goals repository.LearningGoalRepository
}

// NewCreateLearningGoalUseCase は CreateLearningGoalUseCase を生成する。
func NewCreateLearningGoalUseCase(goals repository.LearningGoalRepository) *CreateLearningGoalUseCase {
	return &CreateLearningGoalUseCase{goals: goals}
}

// Execute はタイトルを検証し、カテゴリ未指定なら「その他」を入れて作成する。
func (uc *CreateLearningGoalUseCase) Execute(ctx context.Context, goal *model.LearningGoal) error {
	if err := domain.ValidateStringLength(goal.Title, 1, 200, "タイトル"); err != nil {
		return err
	}
	goal.Title = strings.TrimSpace(goal.Title)
	if goal.Category == "" {
		goal.Category = model.GoalCategoryOther
	}
	return uc.goals.Create(ctx, goal)
}

// GetLearningGoalUseCase は所有者の学習目標を取得する。
type GetLearningGoalUseCase struct {
	goals repository.LearningGoalRepository
}

// NewGetLearningGoalUseCase は GetLearningGoalUseCase を生成する。
func NewGetLearningGoalUseCase(goals repository.LearningGoalRepository) *GetLearningGoalUseCase {
	return &GetLearningGoalUseCase{goals: goals}
}

// Execute は所有権を検証したうえで目標を返す。
func (uc *GetLearningGoalUseCase) Execute(ctx context.Context, id, userID uint) (*model.LearningGoal, error) {
	return ensureOwner(ctx, uc.goals.FindByID, id, userID, learningGoalOwnerOf)
}

// ListLearningGoalsUseCase は指定ユーザーの学習目標をページネーション付きで取得する。
type ListLearningGoalsUseCase struct {
	goals repository.LearningGoalRepository
}

// NewListLearningGoalsUseCase は ListLearningGoalsUseCase を生成する。
func NewListLearningGoalsUseCase(goals repository.LearningGoalRepository) *ListLearningGoalsUseCase {
	return &ListLearningGoalsUseCase{goals: goals}
}

// Execute は目標一覧と総件数を返す。
func (uc *ListLearningGoalsUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.LearningGoal, int64, error) {
	return uc.goals.GetByUserID(ctx, userID, limit, offset)
}

// ListActiveLearningGoalsUseCase は進行中の学習目標を取得する。
type ListActiveLearningGoalsUseCase struct {
	goals repository.LearningGoalRepository
}

// NewListActiveLearningGoalsUseCase は ListActiveLearningGoalsUseCase を生成する。
func NewListActiveLearningGoalsUseCase(goals repository.LearningGoalRepository) *ListActiveLearningGoalsUseCase {
	return &ListActiveLearningGoalsUseCase{goals: goals}
}

// Execute は進行中の目標一覧を返す。
func (uc *ListActiveLearningGoalsUseCase) Execute(ctx context.Context, userID uint) ([]model.LearningGoal, error) {
	return uc.goals.GetActiveByUserID(ctx, userID)
}

// ListLearningGoalsByCategoryUseCase はカテゴリで学習目標を絞り込む。
type ListLearningGoalsByCategoryUseCase struct {
	goals repository.LearningGoalRepository
}

// NewListLearningGoalsByCategoryUseCase は ListLearningGoalsByCategoryUseCase を生成する。
func NewListLearningGoalsByCategoryUseCase(goals repository.LearningGoalRepository) *ListLearningGoalsByCategoryUseCase {
	return &ListLearningGoalsByCategoryUseCase{goals: goals}
}

// Execute はカテゴリを検証したうえで目標一覧を返す。
func (uc *ListLearningGoalsByCategoryUseCase) Execute(ctx context.Context, userID uint, category string) ([]model.LearningGoal, error) {
	if !validGoalCategories[category] {
		return nil, domain.NewError(domain.ErrCodeBadRequest, msgInvalidGoalCategory, nil)
	}
	return uc.goals.GetByCategory(ctx, userID, category)
}

// ListLearningGoalsByStatusUseCase はステータスで学習目標を絞り込む。
type ListLearningGoalsByStatusUseCase struct {
	goals repository.LearningGoalRepository
}

// NewListLearningGoalsByStatusUseCase は ListLearningGoalsByStatusUseCase を生成する。
func NewListLearningGoalsByStatusUseCase(goals repository.LearningGoalRepository) *ListLearningGoalsByStatusUseCase {
	return &ListLearningGoalsByStatusUseCase{goals: goals}
}

// Execute はステータスを検証したうえで目標一覧を返す。
func (uc *ListLearningGoalsByStatusUseCase) Execute(ctx context.Context, userID uint, status string) ([]model.LearningGoal, error) {
	if !validGoalStatuses[status] {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "無効なステータスです", nil)
	}
	return uc.goals.GetByStatus(ctx, userID, status)
}

// GetLearningGoalStatsUseCase は学習目標の統計を取得する。
type GetLearningGoalStatsUseCase struct {
	goals repository.LearningGoalRepository
}

// NewGetLearningGoalStatsUseCase は GetLearningGoalStatsUseCase を生成する。
func NewGetLearningGoalStatsUseCase(goals repository.LearningGoalRepository) *GetLearningGoalStatsUseCase {
	return &GetLearningGoalStatsUseCase{goals: goals}
}

// Execute は統計を返す。
func (uc *GetLearningGoalStatsUseCase) Execute(ctx context.Context, userID uint) (*model.LearningGoalStats, error) {
	return uc.goals.GetStats(ctx, userID)
}

// UpdateLearningGoalUseCase は所有者の学習目標を更新する。
type UpdateLearningGoalUseCase struct {
	goals repository.LearningGoalRepository
}

// NewUpdateLearningGoalUseCase は UpdateLearningGoalUseCase を生成する。
func NewUpdateLearningGoalUseCase(goals repository.LearningGoalRepository) *UpdateLearningGoalUseCase {
	return &UpdateLearningGoalUseCase{goals: goals}
}

// Execute は所有権を検証したうえで目標を部分更新する。
// 進捗が 100 に達した場合はステータスを完了へ自動遷移させる。
func (uc *UpdateLearningGoalUseCase) Execute(ctx context.Context, id, userID uint, updates *model.LearningGoal) (*model.LearningGoal, error) {
	goal, err := ensureOwner(ctx, uc.goals.FindByID, id, userID, learningGoalOwnerOf)
	if err != nil {
		return nil, err
	}

	if title := strings.TrimSpace(updates.Title); title != "" {
		if err := domain.ValidateStringLength(title, 1, 200, "タイトル"); err != nil {
			return nil, err
		}
		goal.Title = title
	}
	if desc := strings.TrimSpace(updates.Description); desc != "" {
		if err := domain.ValidateStringLength(desc, 1, 1000, "説明"); err != nil {
			return nil, err
		}
		goal.Description = desc
	}
	if cat := strings.TrimSpace(string(updates.Category)); cat != "" {
		goal.Category = updates.Category
	}
	if updates.TargetDate != nil {
		goal.TargetDate = updates.TargetDate
	}
	// 0 以上なら常に反映する（未指定の 0 も書き込まれる。移行前の挙動を維持している）。
	if updates.Progress >= 0 {
		progress := updates.Progress
		if progress > 100 {
			progress = 100
		}
		goal.Progress = progress

		if progress == 100 && goal.Status == model.GoalStatusActive {
			goal.Status = model.GoalStatusCompleted
			now := time.Now()
			goal.CompletedAt = &now
		}
	}
	if st := strings.TrimSpace(string(updates.Status)); st != "" {
		goal.Status = updates.Status
		if goal.Status == model.GoalStatusCompleted && goal.CompletedAt == nil {
			now := time.Now()
			goal.CompletedAt = &now
		}
	}

	if err := uc.goals.Update(ctx, goal); err != nil {
		return nil, err
	}
	return goal, nil
}

// GetGoalDeadlineAlertsUseCase は期限が近い/超過した目標のアラートを取得する。
type GetGoalDeadlineAlertsUseCase struct {
	goals repository.LearningGoalRepository
}

// NewGetGoalDeadlineAlertsUseCase は GetGoalDeadlineAlertsUseCase を生成する。
func NewGetGoalDeadlineAlertsUseCase(goals repository.LearningGoalRepository) *GetGoalDeadlineAlertsUseCase {
	return &GetGoalDeadlineAlertsUseCase{goals: goals}
}

// Execute は進行中の目標からアラート対象を抽出して返す。
func (uc *GetGoalDeadlineAlertsUseCase) Execute(ctx context.Context, userID uint) ([]model.GoalDeadlineAlert, error) {
	goals, err := uc.goals.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var alerts []model.GoalDeadlineAlert
	for _, g := range goals {
		status := DeadlineStatus(&g, now)
		if status != "" {
			alerts = append(alerts, model.GoalDeadlineAlert{
				Goal:     g,
				Status:   status,
				DaysLeft: DaysUntilDeadline(&g, now),
			})
		}
	}
	return alerts, nil
}

// DuplicateLearningGoalUseCase は学習目標を複製する。
type DuplicateLearningGoalUseCase struct {
	goals repository.LearningGoalRepository
}

// NewDuplicateLearningGoalUseCase は DuplicateLearningGoalUseCase を生成する。
func NewDuplicateLearningGoalUseCase(goals repository.LearningGoalRepository) *DuplicateLearningGoalUseCase {
	return &DuplicateLearningGoalUseCase{goals: goals}
}

// Execute は所有権を検証したうえで、進捗 0・進行中の複製を作る。
func (uc *DuplicateLearningGoalUseCase) Execute(ctx context.Context, id, userID uint) (*model.LearningGoal, error) {
	goal, err := ensureOwner(ctx, uc.goals.FindByID, id, userID, learningGoalOwnerOf)
	if err != nil {
		return nil, err
	}

	newGoal := &model.LearningGoal{
		UserID:      goal.UserID,
		Title:       goal.Title + " (コピー)",
		Description: goal.Description,
		Category:    goal.Category,
		TargetDate:  goal.TargetDate,
		Progress:    0,
		Status:      model.GoalStatusActive,
	}
	if err := uc.goals.Create(ctx, newGoal); err != nil {
		return nil, err
	}
	return newGoal, nil
}

// ToggleLearningGoalShareUseCase は学習目標の公開/非公開を切り替える。
type ToggleLearningGoalShareUseCase struct {
	goals repository.LearningGoalRepository
}

// NewToggleLearningGoalShareUseCase は ToggleLearningGoalShareUseCase を生成する。
func NewToggleLearningGoalShareUseCase(goals repository.LearningGoalRepository) *ToggleLearningGoalShareUseCase {
	return &ToggleLearningGoalShareUseCase{goals: goals}
}

// Execute は所有権を検証したうえで公開状態を反転する。
func (uc *ToggleLearningGoalShareUseCase) Execute(ctx context.Context, id, userID uint) (*model.LearningGoal, error) {
	goal, err := ensureOwner(ctx, uc.goals.FindByID, id, userID, learningGoalOwnerOf)
	if err != nil {
		return nil, err
	}

	goal.IsPublic = !goal.IsPublic
	if err := uc.goals.Update(ctx, goal); err != nil {
		return nil, err
	}
	return goal, nil
}

// ListPublicLearningGoalsUseCase は全ユーザーの公開目標を取得する。
type ListPublicLearningGoalsUseCase struct {
	goals repository.LearningGoalRepository
}

// NewListPublicLearningGoalsUseCase は ListPublicLearningGoalsUseCase を生成する。
func NewListPublicLearningGoalsUseCase(goals repository.LearningGoalRepository) *ListPublicLearningGoalsUseCase {
	return &ListPublicLearningGoalsUseCase{goals: goals}
}

// Execute は公開目標一覧と総件数を返す。
func (uc *ListPublicLearningGoalsUseCase) Execute(ctx context.Context, limit, offset int) ([]model.LearningGoal, int64, error) {
	return uc.goals.GetPublicGoals(ctx, limit, offset)
}

// ListPublicLearningGoalsByUserUseCase は指定ユーザーの公開目標を取得する。
type ListPublicLearningGoalsByUserUseCase struct {
	goals repository.LearningGoalRepository
}

// NewListPublicLearningGoalsByUserUseCase は ListPublicLearningGoalsByUserUseCase を生成する。
func NewListPublicLearningGoalsByUserUseCase(goals repository.LearningGoalRepository) *ListPublicLearningGoalsByUserUseCase {
	return &ListPublicLearningGoalsByUserUseCase{goals: goals}
}

// Execute は指定ユーザーの公開目標一覧と総件数を返す。
func (uc *ListPublicLearningGoalsByUserUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.LearningGoal, int64, error) {
	return uc.goals.GetPublicByUserID(ctx, userID, limit, offset)
}

// CountLearningGoalsUseCase は指定ユーザーの学習目標総数を返す。
type CountLearningGoalsUseCase struct {
	goals repository.LearningGoalRepository
}

// NewCountLearningGoalsUseCase は CountLearningGoalsUseCase を生成する。
func NewCountLearningGoalsUseCase(goals repository.LearningGoalRepository) *CountLearningGoalsUseCase {
	return &CountLearningGoalsUseCase{goals: goals}
}

// Execute は目標総数を返す。
func (uc *CountLearningGoalsUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.goals.CountByUserID(ctx, userID)
}

// DeleteLearningGoalUseCase は所有者の学習目標を削除する。
type DeleteLearningGoalUseCase struct {
	goals repository.LearningGoalRepository
}

// NewDeleteLearningGoalUseCase は DeleteLearningGoalUseCase を生成する。
func NewDeleteLearningGoalUseCase(goals repository.LearningGoalRepository) *DeleteLearningGoalUseCase {
	return &DeleteLearningGoalUseCase{goals: goals}
}

// Execute は所有権を検証したうえで目標を削除する。
func (uc *DeleteLearningGoalUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensureOwner(ctx, uc.goals.FindByID, id, userID, learningGoalOwnerOf); err != nil {
		return err
	}
	return uc.goals.Delete(ctx, id)
}

// GoalProgressUpdate は進捗一括更新の 1 件分。
type GoalProgressUpdate struct {
	GoalID   uint
	Progress int
}

// BatchUpdateGoalProgressUseCase は複数の学習目標の進捗をまとめて更新する。
type BatchUpdateGoalProgressUseCase struct {
	update *UpdateLearningGoalUseCase
}

// NewBatchUpdateGoalProgressUseCase は BatchUpdateGoalProgressUseCase を生成する。
func NewBatchUpdateGoalProgressUseCase(update *UpdateLearningGoalUseCase) *BatchUpdateGoalProgressUseCase {
	return &BatchUpdateGoalProgressUseCase{update: update}
}

// Execute は順に進捗を更新する。
// 途中で失敗した場合はそこで中断し、それまでの更新は取り消さない（移行前の挙動を維持している）。
func (uc *BatchUpdateGoalProgressUseCase) Execute(ctx context.Context, userID uint, updates []GoalProgressUpdate) ([]model.LearningGoal, error) {
	var results []model.LearningGoal
	for _, u := range updates {
		goal, err := uc.update.Execute(ctx, u.GoalID, userID, &model.LearningGoal{Progress: u.Progress})
		if err != nil {
			return nil, err
		}
		results = append(results, *goal)
	}
	return results, nil
}

// GetGoalForecastUseCase は進行中の目標に対する達成予測を返す。
type GetGoalForecastUseCase struct {
	goals repository.LearningGoalRepository
}

// NewGetGoalForecastUseCase は GetGoalForecastUseCase を生成する。
func NewGetGoalForecastUseCase(goals repository.LearningGoalRepository) *GetGoalForecastUseCase {
	return &GetGoalForecastUseCase{goals: goals}
}

// Execute は進行中の目標ごとに達成予測を算出して返す。
func (uc *GetGoalForecastUseCase) Execute(ctx context.Context, userID uint) ([]model.GoalForecast, error) {
	goals, err := uc.goals.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var forecasts []model.GoalForecast
	for _, goal := range goals {
		forecasts = append(forecasts, *CalculateGoalForecast(&goal, 0, 0, now))
	}
	return forecasts, nil
}

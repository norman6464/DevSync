package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetWeeklyGoalProgressUseCase は指定ユーザーの全カテゴリ週間目標の達成状況を返す。
type GetWeeklyGoalProgressUseCase struct {
	goals repository.WeeklyGoalRepository
}

// NewGetWeeklyGoalProgressUseCase は GetWeeklyGoalProgressUseCase を生成する。
func NewGetWeeklyGoalProgressUseCase(goals repository.WeeklyGoalRepository) *GetWeeklyGoalProgressUseCase {
	return &GetWeeklyGoalProgressUseCase{goals: goals}
}

// Execute は各カテゴリの目標時間と今週の実績から達成率を算出して返す。
func (uc *GetWeeklyGoalProgressUseCase) Execute(ctx context.Context, userID uint) ([]model.WeeklyGoalProgress, error) {
	goals, err := uc.goals.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var progress []model.WeeklyGoalProgress
	for _, g := range goals {
		actual, err := uc.goals.SumDurationByUserCategoryThisWeek(ctx, userID, string(g.Category))
		if err != nil {
			return nil, err
		}
		pct := 0
		if g.TargetMinutes > 0 {
			pct = actual * 100 / g.TargetMinutes
		}
		progress = append(progress, model.WeeklyGoalProgress{
			Category:        g.Category,
			TargetMinutes:   g.TargetMinutes,
			ActualMinutes:   actual,
			ProgressPercent: pct,
		})
	}
	return progress, nil
}

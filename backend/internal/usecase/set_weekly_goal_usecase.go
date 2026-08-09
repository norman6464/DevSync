package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// SetWeeklyGoalUseCase はカテゴリ別の週間学習目標を設定する（既存があれば更新）。
type SetWeeklyGoalUseCase struct {
	goals repository.WeeklyGoalRepository
}

// NewSetWeeklyGoalUseCase は SetWeeklyGoalUseCase を生成する。
func NewSetWeeklyGoalUseCase(goals repository.WeeklyGoalRepository) *SetWeeklyGoalUseCase {
	return &SetWeeklyGoalUseCase{goals: goals}
}

// Execute はカテゴリと目標時間を検証し、週間目標を upsert する。
func (uc *SetWeeklyGoalUseCase) Execute(ctx context.Context, userID uint, category string, targetMinutes int) (*model.WeeklyGoal, error) {
	if !model.ValidCategories[model.LogCategory(category)] {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "無効なカテゴリです", nil)
	}
	if targetMinutes < 0 || targetMinutes > 10080 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "目標時間は0〜10080分（1週間）で指定してください", nil)
	}

	goal := &model.WeeklyGoal{
		UserID:        userID,
		Category:      model.LogCategory(category),
		TargetMinutes: targetMinutes,
	}
	if err := uc.goals.Upsert(ctx, goal); err != nil {
		return nil, err
	}
	return goal, nil
}

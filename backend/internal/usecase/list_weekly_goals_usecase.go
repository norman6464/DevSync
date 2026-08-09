package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// ListWeeklyGoalsUseCase は指定ユーザーの全カテゴリ週間目標を取得する。
type ListWeeklyGoalsUseCase struct {
	goals repository.WeeklyGoalRepository
}

// NewListWeeklyGoalsUseCase は ListWeeklyGoalsUseCase を生成する。
func NewListWeeklyGoalsUseCase(goals repository.WeeklyGoalRepository) *ListWeeklyGoalsUseCase {
	return &ListWeeklyGoalsUseCase{goals: goals}
}

// Execute はユーザーの全カテゴリ週間目標を返す。
func (uc *ListWeeklyGoalsUseCase) Execute(ctx context.Context, userID uint) ([]model.WeeklyGoal, error) {
	return uc.goals.GetByUserID(ctx, userID)
}

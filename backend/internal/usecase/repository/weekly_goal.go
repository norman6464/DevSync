package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// WeeklyGoalRepository はカテゴリ別週間学習目標の永続化に対する、usecase 側が要求する契約。
type WeeklyGoalRepository interface {
	Upsert(ctx context.Context, goal *model.WeeklyGoal) error
	GetByUserID(ctx context.Context, userID uint) ([]model.WeeklyGoal, error)
	SumDurationByUserCategoryThisWeek(ctx context.Context, userID uint, category string) (int, error)
}

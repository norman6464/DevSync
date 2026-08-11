package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// LearningGoalRepository は学習目標の永続化に対する、usecase 側が要求する契約。
type LearningGoalRepository interface {
	Create(ctx context.Context, goal *model.LearningGoal) error
	Update(ctx context.Context, goal *model.LearningGoal) error
	Delete(ctx context.Context, id uint) error
	// FindByID は指定 ID の目標を返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindByID(ctx context.Context, id uint) (*model.LearningGoal, error)
	GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.LearningGoal, int64, error)
	GetActiveByUserID(ctx context.Context, userID uint) ([]model.LearningGoal, error)
	GetByCategory(ctx context.Context, userID uint, category string) ([]model.LearningGoal, error)
	GetByStatus(ctx context.Context, userID uint, status string) ([]model.LearningGoal, error)
	GetPublicByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.LearningGoal, int64, error)
	GetPublicGoals(ctx context.Context, limit, offset int) ([]model.LearningGoal, int64, error)
	GetStats(ctx context.Context, userID uint) (*model.LearningGoalStats, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
}

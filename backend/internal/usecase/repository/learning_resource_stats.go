package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// LearningResourceStatsRepository はユーザー学習リソース活動集計統計の取得に対する、usecase 側が要求する契約。
type LearningResourceStatsRepository interface {
	GetLearningResourceStats(ctx context.Context, userID uint) (*model.LearningResourceStats, error)
}

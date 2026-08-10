package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// LearningLogStatsRepository はユーザー学習ログ集計統計の取得に対する、usecase 側が要求する契約。
type LearningLogStatsRepository interface {
	GetLearningLogStats(ctx context.Context, userID uint) (*model.LearningLogStats, error)
}

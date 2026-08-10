package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// FollowStatsRepository はユーザーフォロー関係集計統計の取得に対する、usecase 側が要求する契約。
type FollowStatsRepository interface {
	GetFollowStats(ctx context.Context, userID uint) (*model.FollowStats, error)
}

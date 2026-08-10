package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// PostStatsRepository はユーザー投稿集計統計の取得に対する、usecase 側が要求する契約。
type PostStatsRepository interface {
	GetPostStats(ctx context.Context, userID uint) (*model.PostStats, error)
}

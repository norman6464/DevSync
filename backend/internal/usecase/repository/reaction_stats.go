package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// ReactionStatsRepository はユーザーリアクション集計統計の取得に対する、usecase 側が要求する契約。
type ReactionStatsRepository interface {
	GetReactionStats(ctx context.Context, userID uint) (*model.ReactionStats, error)
	GetEmojiBreakdown(ctx context.Context, userID uint) ([]model.ReactionCount, error)
	GetTopReactedPosts(ctx context.Context, userID uint, limit int) ([]model.TopReactedPost, error)
}

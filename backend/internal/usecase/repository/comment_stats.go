package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// CommentStatsRepository はユーザーコメント活動集計統計の取得に対する、usecase 側が要求する契約。
type CommentStatsRepository interface {
	GetCommentStats(ctx context.Context, userID uint) (*model.CommentStats, error)
}

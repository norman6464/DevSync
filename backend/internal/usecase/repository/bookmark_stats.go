package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// BookmarkStatsRepository はユーザーブックマーク集計統計の取得に対する、usecase 側が要求する契約。
type BookmarkStatsRepository interface {
	GetBookmarkStats(ctx context.Context, userID uint) (*model.BookmarkStats, error)
}

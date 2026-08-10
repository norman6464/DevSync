package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// BookReviewStatsRepository はユーザー書籍レビュー集計統計の取得に対する、usecase 側が要求する契約。
type BookReviewStatsRepository interface {
	GetBookReviewStats(ctx context.Context, userID uint) (*model.BookReviewStats, error)
}

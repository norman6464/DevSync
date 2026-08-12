package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// PostBookmarkRepository は投稿のブックマークの永続化に対する、usecase 側が要求する契約。
type PostBookmarkRepository interface {
	Bookmark(ctx context.Context, userID, postID uint) error
	Unbookmark(ctx context.Context, userID, postID uint) error
	// HasBookmarked は指定ユーザーが投稿をブックマーク済みかを返す。
	HasBookmarked(ctx context.Context, userID, postID uint) (bool, error)
	// FindBookmarkedByUserID は指定ユーザーのブックマーク済み投稿を新しい順に返す（総件数も返す）。
	FindBookmarkedByUserID(ctx context.Context, userID uint, page, limit int) ([]model.Post, int64, error)
	// CountBookmarkedByUserID は指定ユーザーのブックマーク件数を返す。
	CountBookmarkedByUserID(ctx context.Context, userID uint) (int64, error)
}

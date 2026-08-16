package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// BookmarkCollectionRepository はブックマークコレクションの永続化に対する、usecase 側が要求する契約。
type BookmarkCollectionRepository interface {
	Create(ctx context.Context, collection *model.BookmarkCollection) error
	FindByID(ctx context.Context, id uint) (*model.BookmarkCollection, error)
	FindByUserID(ctx context.Context, userID uint) ([]model.BookmarkCollection, error)
	Update(ctx context.Context, collection *model.BookmarkCollection) error
	Delete(ctx context.Context, id uint) error
	// AddPost はコレクションへ投稿を原子的に追加する。
	// 既に同じ投稿が入っている場合は追加せず (false, nil) を返す（同時実行でも重複しない）。
	AddPost(ctx context.Context, item *model.BookmarkCollectionItem) (bool, error)
	RemovePost(ctx context.Context, collectionID, postID uint) error
	GetPosts(ctx context.Context, collectionID uint, limit, offset int) ([]model.Post, int64, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
}

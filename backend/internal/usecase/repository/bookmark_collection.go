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
	AddPost(ctx context.Context, item *model.BookmarkCollectionItem) error
	RemovePost(ctx context.Context, collectionID, postID uint) error
	GetPosts(ctx context.Context, collectionID uint, limit, offset int) ([]model.Post, int64, error)
	HasPost(ctx context.Context, collectionID, postID uint) (bool, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
}

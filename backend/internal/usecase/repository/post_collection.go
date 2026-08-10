package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// PostCollectionRepository は投稿コレクションの永続化に対する、usecase 側が要求する契約。
type PostCollectionRepository interface {
	Create(ctx context.Context, collection *model.PostCollection) error
	FindByID(ctx context.Context, id uint) (*model.PostCollection, error)
	FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.PostCollection, int64, error)
	FindPublicByUserID(ctx context.Context, userID uint) ([]model.PostCollection, error)
	Update(ctx context.Context, collection *model.PostCollection) error
	Delete(ctx context.Context, id uint) error
	AddPost(ctx context.Context, item *model.PostCollectionItem) error
	RemovePost(ctx context.Context, collectionID, postID uint) error
	HasPost(ctx context.Context, collectionID, postID uint) (bool, error)
	GetPostsByCollectionID(ctx context.Context, collectionID uint) ([]model.PostCollectionItem, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
}

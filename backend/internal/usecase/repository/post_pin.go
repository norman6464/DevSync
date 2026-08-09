package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// PostPinRepository は投稿ピン留めの永続化に対する、usecase 側が要求する契約。
type PostPinRepository interface {
	Pin(ctx context.Context, pin *model.PostPin) error
	Unpin(ctx context.Context, userID, postID uint) error
	GetByUserID(ctx context.Context, userID uint) ([]model.PostPin, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
	IsPinned(ctx context.Context, userID, postID uint) (bool, error)
	UpdateOrder(ctx context.Context, userID uint, postIDs []uint) error
}

// PostReader は所有権チェックに必要な投稿読み取りだけを切り出した最小 port（-er）。
type PostReader interface {
	FindByID(ctx context.Context, id uint) (*model.Post, error)
}

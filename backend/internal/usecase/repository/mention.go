package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// MentionRepository はメンションの永続化に対する、usecase 側が要求する契約。
type MentionRepository interface {
	Create(ctx context.Context, mention *model.Mention) error
	// FindByUserID は指定ユーザー宛のメンションを新しい順にページネーションして返す。
	FindByUserID(ctx context.Context, userID uint, page, limit int) ([]model.Mention, error)
	// FindByPostID は指定投稿に紐づくメンションを返す。
	FindByPostID(ctx context.Context, postID uint) ([]model.Mention, error)
	// DeleteByPostID は指定投稿に紐づくメンションをすべて削除する。
	DeleteByPostID(ctx context.Context, postID uint) error
	// DeleteByCommentID は指定コメントに紐づくメンションをすべて削除する。
	DeleteByCommentID(ctx context.Context, commentID uint) error
}

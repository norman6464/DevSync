package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// PostCommentRepository は投稿コメントの永続化に対する、usecase 側が要求する契約。
// 単体取得は所有権チェック用の最小 port である CommentReader を再利用する
// （不在は error で表現する。CommentReader の契約に合わせる）。
type PostCommentRepository interface {
	CommentReader

	Create(ctx context.Context, comment *model.Comment) error
	Update(ctx context.Context, comment *model.Comment) error
	Delete(ctx context.Context, id uint) error
	ListByPostID(ctx context.Context, postID uint) ([]model.Comment, error)
	ListReplies(ctx context.Context, parentID uint) ([]model.Comment, error)
}

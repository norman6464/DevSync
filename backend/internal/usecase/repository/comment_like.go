package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// CommentLikeRepository はコメントいいねの永続化に対する、usecase 側が要求する契約。
type CommentLikeRepository interface {
	Like(ctx context.Context, userID, commentID uint) error
	Unlike(ctx context.Context, userID, commentID uint) error
	HasLiked(ctx context.Context, userID, commentID uint) (bool, error)
	CountByCommentID(ctx context.Context, commentID uint) (int64, error)
}

// CommentReader は所有権・存在チェックに必要なコメント読み取りだけを切り出した最小 port。
// 投稿(post)スライスの読み取りへの依存を、必要なメソッドだけに絞る（-er port）。
type CommentReader interface {
	FindCommentByID(ctx context.Context, id uint) (*model.Comment, error)
}

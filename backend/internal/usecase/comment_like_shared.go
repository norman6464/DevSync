package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// ensureNotSelfComment はコメントの存在を確認し、自分のコメントへの操作を禁止する。
// Like / Unlike で共有する前提チェック。存在しなければ ErrNotFound、自分のものなら ErrForbidden。
func ensureNotSelfComment(ctx context.Context, reader repository.CommentReader, userID, commentID uint) error {
	comment, err := reader.FindCommentByID(ctx, commentID)
	if err != nil {
		return domain.ErrNotFound
	}
	if comment.UserID == userID {
		return domain.ErrForbidden
	}
	return nil
}

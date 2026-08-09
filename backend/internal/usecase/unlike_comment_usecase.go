package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// UnlikeCommentUseCase はコメントのいいねを取り消す。
type UnlikeCommentUseCase struct {
	likes    repository.CommentLikeRepository
	comments repository.CommentReader
}

// NewUnlikeCommentUseCase は UnlikeCommentUseCase を生成する。
func NewUnlikeCommentUseCase(likes repository.CommentLikeRepository, comments repository.CommentReader) *UnlikeCommentUseCase {
	return &UnlikeCommentUseCase{likes: likes, comments: comments}
}

// Execute はコメントのいいねを取り消す。存在しない・自分のコメント・未いいねはエラー。
func (uc *UnlikeCommentUseCase) Execute(ctx context.Context, userID, commentID uint) error {
	if err := ensureNotSelfComment(ctx, uc.comments, userID, commentID); err != nil {
		return err
	}
	liked, err := uc.likes.HasLiked(ctx, userID, commentID)
	if err != nil {
		return err
	}
	if !liked {
		return domain.NewError(domain.ErrCodeNotFound, "いいねしていません", nil)
	}
	return uc.likes.Unlike(ctx, userID, commentID)
}

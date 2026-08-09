package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// LikeCommentUseCase はコメントにいいねする。
type LikeCommentUseCase struct {
	likes    repository.CommentLikeRepository
	comments repository.CommentReader
}

// NewLikeCommentUseCase は LikeCommentUseCase を生成する。
func NewLikeCommentUseCase(likes repository.CommentLikeRepository, comments repository.CommentReader) *LikeCommentUseCase {
	return &LikeCommentUseCase{likes: likes, comments: comments}
}

// Execute はコメントにいいねする。存在しない・自分のコメント・いいね済みはエラー。
func (uc *LikeCommentUseCase) Execute(ctx context.Context, userID, commentID uint) error {
	if err := ensureNotSelfComment(ctx, uc.comments, userID, commentID); err != nil {
		return err
	}
	liked, err := uc.likes.HasLiked(ctx, userID, commentID)
	if err != nil {
		return err
	}
	if liked {
		return domain.NewError(domain.ErrCodeBadRequest, "すでにいいね済みです", nil)
	}
	return uc.likes.Like(ctx, userID, commentID)
}

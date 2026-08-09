package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetCommentLikeStatusUseCase はコメントのいいね状態（いいね済みか・いいね数）を取得する。
type GetCommentLikeStatusUseCase struct {
	likes    repository.CommentLikeRepository
	comments repository.CommentReader
}

// NewGetCommentLikeStatusUseCase は GetCommentLikeStatusUseCase を生成する。
func NewGetCommentLikeStatusUseCase(likes repository.CommentLikeRepository, comments repository.CommentReader) *GetCommentLikeStatusUseCase {
	return &GetCommentLikeStatusUseCase{likes: likes, comments: comments}
}

// Execute はコメントの (いいね済みか, いいね数) を返す。コメントが存在しなければ ErrNotFound。
func (uc *GetCommentLikeStatusUseCase) Execute(ctx context.Context, userID, commentID uint) (bool, int64, error) {
	if _, err := uc.comments.FindCommentByID(ctx, commentID); err != nil {
		return false, 0, domain.ErrNotFound
	}
	liked, err := uc.likes.HasLiked(ctx, userID, commentID)
	if err != nil {
		return false, 0, err
	}
	count, err := uc.likes.CountByCommentID(ctx, commentID)
	if err != nil {
		return false, 0, err
	}
	return liked, count, nil
}

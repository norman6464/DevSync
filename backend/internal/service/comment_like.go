package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// CommentLikeService はコメントへのいいね機能のビジネスロジックを提供する。
type CommentLikeService struct {
	likeRepo repository.CommentLikeRepositoryInterface
	postRepo repository.PostRepositoryInterface
}

// NewCommentLikeService は新しいCommentLikeServiceインスタンスを生成する。
func NewCommentLikeService(
	likeRepo repository.CommentLikeRepositoryInterface,
	postRepo repository.PostRepositoryInterface,
) *CommentLikeService {
	return &CommentLikeService{likeRepo: likeRepo, postRepo: postRepo}
}

// Like はコメントにいいねする。
// コメントが存在しない場合、自分のコメントの場合、すでにいいね済みの場合はエラーを返す。
func (s *CommentLikeService) Like(userID, commentID uint) error {
	comment, err := s.postRepo.FindCommentByID(commentID)
	if err != nil {
		return ErrNotFound
	}
	if comment.UserID == userID {
		return ErrForbidden
	}

	liked, err := s.likeRepo.HasLiked(userID, commentID)
	if err != nil {
		return err
	}
	if liked {
		return domain.NewError(domain.ErrCodeBadRequest, "すでにいいね済みです", nil)
	}

	return s.likeRepo.Like(userID, commentID)
}

// Unlike はコメントのいいねを取り消す。
// コメントが存在しない場合、自分のコメントの場合、いいねしていない場合はエラーを返す。
func (s *CommentLikeService) Unlike(userID, commentID uint) error {
	comment, err := s.postRepo.FindCommentByID(commentID)
	if err != nil {
		return ErrNotFound
	}
	if comment.UserID == userID {
		return ErrForbidden
	}

	liked, err := s.likeRepo.HasLiked(userID, commentID)
	if err != nil {
		return err
	}
	if !liked {
		return domain.NewError(domain.ErrCodeNotFound, "いいねしていません", nil)
	}

	return s.likeRepo.Unlike(userID, commentID)
}

// GetStatus はコメントのいいね状態（いいねしているか・いいね数）を返す。
func (s *CommentLikeService) GetStatus(userID, commentID uint) (bool, int64, error) {
	if _, err := s.postRepo.FindCommentByID(commentID); err != nil {
		return false, 0, ErrNotFound
	}

	liked, err := s.likeRepo.HasLiked(userID, commentID)
	if err != nil {
		return false, 0, err
	}

	count, err := s.likeRepo.CountByCommentID(commentID)
	if err != nil {
		return false, 0, err
	}

	return liked, count, nil
}

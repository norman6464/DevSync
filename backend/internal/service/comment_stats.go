package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// CommentStatsService はユーザーコメント活動集計統計のビジネスロジックを提供する。
type CommentStatsService struct {
	repo repository.CommentStatsRepositoryInterface
}

// NewCommentStatsService は新しいCommentStatsServiceインスタンスを生成する。
func NewCommentStatsService(repo repository.CommentStatsRepositoryInterface) *CommentStatsService {
	return &CommentStatsService{repo: repo}
}

// GetCommentStats は指定ユーザーのコメント活動集計統計を取得する。
func (s *CommentStatsService) GetCommentStats(userID uint) (*model.CommentStats, error) {
	if userID == 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "userIDは必須です", nil)
	}
	return s.repo.GetCommentStats(userID)
}

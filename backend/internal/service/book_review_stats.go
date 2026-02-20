package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// BookReviewStatsService はユーザー書籍レビュー集計統計のビジネスロジックを提供する。
type BookReviewStatsService struct {
	repo repository.BookReviewStatsRepositoryInterface
}

// NewBookReviewStatsService は新しいBookReviewStatsServiceインスタンスを生成する。
func NewBookReviewStatsService(repo repository.BookReviewStatsRepositoryInterface) *BookReviewStatsService {
	return &BookReviewStatsService{repo: repo}
}

// GetBookReviewStats は指定ユーザーの書籍レビュー集計統計を取得する。
func (s *BookReviewStatsService) GetBookReviewStats(userID uint) (*model.BookReviewStats, error) {
	if err := validateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return s.repo.GetBookReviewStats(userID)
}

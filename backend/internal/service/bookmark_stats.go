package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// BookmarkStatsService はユーザーブックマーク集計統計のビジネスロジックを提供する。
type BookmarkStatsService struct {
	repo repository.BookmarkStatsRepositoryInterface
}

// NewBookmarkStatsService は新しいBookmarkStatsServiceインスタンスを生成する。
func NewBookmarkStatsService(repo repository.BookmarkStatsRepositoryInterface) *BookmarkStatsService {
	return &BookmarkStatsService{repo: repo}
}

// GetBookmarkStats は指定ユーザーのブックマーク集計統計を取得する。
func (s *BookmarkStatsService) GetBookmarkStats(userID uint) (*model.BookmarkStats, error) {
	if userID == 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "userIDは必須です", nil)
	}
	return s.repo.GetBookmarkStats(userID)
}

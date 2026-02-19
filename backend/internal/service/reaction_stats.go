package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// ReactionStatsService はユーザーリアクション集計統計のビジネスロジックを提供する。
type ReactionStatsService struct {
	repo repository.ReactionStatsRepositoryInterface
}

// NewReactionStatsService は新しいReactionStatsServiceインスタンスを生成する。
func NewReactionStatsService(repo repository.ReactionStatsRepositoryInterface) *ReactionStatsService {
	return &ReactionStatsService{repo: repo}
}

// GetReactionStats は指定ユーザーのリアクション集計統計を取得する。
func (s *ReactionStatsService) GetReactionStats(userID uint) (*model.ReactionStats, error) {
	if userID == 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "userIDは必須です", nil)
	}
	return s.repo.GetReactionStats(userID)
}

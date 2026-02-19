package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// StudyCircleStatsService はスタディサークルの集計統計のビジネスロジックを提供する。
type StudyCircleStatsService struct {
	repo repository.StudyCircleStatsRepositoryInterface
}

// NewStudyCircleStatsService は新しいStudyCircleStatsServiceインスタンスを生成する。
func NewStudyCircleStatsService(repo repository.StudyCircleStatsRepositoryInterface) *StudyCircleStatsService {
	return &StudyCircleStatsService{repo: repo}
}

// GetCircleStats は指定サークルの集計統計を返す。
func (s *StudyCircleStatsService) GetCircleStats(circleID uint) (*model.StudyCircleStats, error) {
	if circleID == 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "circleIDは必須です", nil)
	}
	return s.repo.GetCircleStats(circleID)
}

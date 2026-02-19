package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// NoteStatsService はユーザーのノート集計統計のビジネスロジックを提供する。
type NoteStatsService struct {
	repo repository.NoteStatsRepositoryInterface
}

// NewNoteStatsService は新しいNoteStatsServiceインスタンスを生成する。
func NewNoteStatsService(repo repository.NoteStatsRepositoryInterface) *NoteStatsService {
	return &NoteStatsService{repo: repo}
}

// GetNoteStats は指定ユーザーのノート集計統計を返す。
func (s *NoteStatsService) GetNoteStats(userID uint) (*model.NoteStats, error) {
	if userID == 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "userIDは必須です", nil)
	}
	return s.repo.GetNoteStats(userID)
}

package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetNoteStatsUseCase は指定ユーザーのノート集計統計を取得する。
type GetNoteStatsUseCase struct {
	stats repository.NoteStatsRepository
}

// NewGetNoteStatsUseCase は GetNoteStatsUseCase を生成する。
func NewGetNoteStatsUseCase(stats repository.NoteStatsRepository) *GetNoteStatsUseCase {
	return &GetNoteStatsUseCase{stats: stats}
}

// Execute はユーザー ID を検証し、ノート集計統計を返す。
func (uc *GetNoteStatsUseCase) Execute(ctx context.Context, userID uint) (*model.NoteStats, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return uc.stats.GetNoteStats(ctx, userID)
}

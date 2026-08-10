package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// NoteStatsRepository はノート集計統計の取得に対する、usecase 側が要求する契約。
type NoteStatsRepository interface {
	GetNoteStats(ctx context.Context, userID uint) (*model.NoteStats, error)
}

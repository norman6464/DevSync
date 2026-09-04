package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// noteStatsRepository は [repository.NoteStatsRepository] の sqlc(pgx) 実装。
type noteStatsRepository struct {
	q *sqlcgen.Queries
}

// NewNoteStatsRepository は NoteStatsRepository の sqlc(pgx) 実装を返す。
func NewNoteStatsRepository(q *sqlcgen.Queries) repository.NoteStatsRepository {
	return &noteStatsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NoteStatsRepository = (*noteStatsRepository)(nil)

// GetNoteStats は指定ユーザーのノート集計統計を返す。
func (r *noteStatsRepository) GetNoteStats(ctx context.Context, userID uint) (*model.NoteStats, error) {
	total, err := r.q.CountNotesByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	archived, err := r.q.CountArchivedNotesByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	favorite, err := r.q.CountFavoriteNotesByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	folders, err := r.q.CountNoteFoldersByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	weekAgo := domain.DaysAgo(time.Now(), 7)
	thisWeek, err := r.q.CountNotesByUserSince(ctx, sqlcgen.CountNotesByUserSinceParams{
		UserID:    int64(userID),
		CreatedAt: toTimestamptzNotNull(weekAgo),
	})
	if err != nil {
		return nil, err
	}

	return &model.NoteStats{
		TotalNotes:    total,
		ArchivedNotes: archived,
		FavoriteNotes: favorite,
		TotalFolders:  folders,
		NotesThisWeek: thisWeek,
	}, nil
}

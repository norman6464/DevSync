package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// noteUpdater は [repository.NoteUpdater] の sqlc(pgx) 実装。
// バージョン復元でノート本体を書き戻すためだけに使う。
type noteUpdater struct {
	q *sqlcgen.Queries
}

// NewNoteUpdater は NoteUpdater の sqlc(pgx) 実装を返す。
func NewNoteUpdater(q *sqlcgen.Queries) repository.NoteUpdater {
	return &noteUpdater{q: q}
}

var _ repository.NoteUpdater = (*noteUpdater)(nil)

// Update はノート本体を書き戻す（既存 NoteRepository.Update と同じ）。
func (r *noteUpdater) Update(ctx context.Context, note *model.Note) error {
	row, err := r.q.UpdateNote(ctx, sqlcgen.UpdateNoteParams{
		ID:         int64(note.ID),
		Title:      note.Title,
		Content:    &note.Content,
		Tags:       &note.Tags,
		FolderID:   toInt64PtrFromUintPtr(note.FolderID),
		IsFavorite: &note.IsFavorite,
		IsArchived: &note.IsArchived,
	})
	if err != nil {
		return err
	}
	*note = toModelNote(row)
	return nil
}

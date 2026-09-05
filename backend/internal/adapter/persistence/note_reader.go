package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// noteReader は [repository.NoteReader] の sqlc(pgx) 実装。
// 所有権チェックに必要なノート読み取りだけを提供する。
type noteReader struct {
	q *sqlcgen.Queries
}

// NewNoteReader は NoteReader の sqlc(pgx) 実装を返す。
func NewNoteReader(q *sqlcgen.Queries) repository.NoteReader {
	return &noteReader{q: q}
}

var _ repository.NoteReader = (*noteReader)(nil)

// FindByID は ID でノートを取得する（既存 NoteRepository.FindByID と同じ preload）。
// 不在の場合は (nil, nil) を返す。
func (r *noteReader) FindByID(ctx context.Context, id uint) (*model.Note, error) {
	row, err := r.q.GetNoteByID(ctx, int64(id))
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	note := toModelNote(row.Note)
	note.User = toModelUser(row.User)
	attachFolder(&note, row.FolderID2, row.FolderUserID, row.FolderParentID, row.FolderName, row.FolderCreatedAt, row.FolderUpdatedAt)
	return &note, nil
}

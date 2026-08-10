package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// NoteLinkRepository はノート間リンクの永続化に対する、usecase 側が要求する契約。
type NoteLinkRepository interface {
	Create(ctx context.Context, link *model.NoteLink) error
	FindBySourceNoteID(ctx context.Context, sourceNoteID uint) ([]model.NoteLink, error)
	FindByTargetNoteID(ctx context.Context, targetNoteID uint) ([]model.NoteLink, error)
	Delete(ctx context.Context, sourceNoteID, targetNoteID uint) error
	Exists(ctx context.Context, sourceNoteID, targetNoteID uint) (bool, error)
	CountBySourceNoteID(ctx context.Context, noteID uint) (int64, error)
	CountByTargetNoteID(ctx context.Context, noteID uint) (int64, error)
}

// NoteReader は所有権チェックに必要なノート読み取りだけを切り出した最小 port（-er）。
type NoteReader interface {
	FindByID(ctx context.Context, id uint) (*model.Note, error)
}

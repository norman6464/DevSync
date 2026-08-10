package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// noteReader は [repository.NoteReader] の GORM 実装。
// 所有権チェックに必要なノート読み取りだけを提供する。
type noteReader struct {
	db *gorm.DB
}

// NewNoteReader は NoteReader の GORM 実装を返す。
func NewNoteReader(db *gorm.DB) repository.NoteReader {
	return &noteReader{db: db}
}

var _ repository.NoteReader = (*noteReader)(nil)

// FindByID は ID でノートを取得する（既存 NoteRepository.FindByID と同じ preload）。
// 不在の場合は (nil, nil) を返す。
func (r *noteReader) FindByID(ctx context.Context, id uint) (*model.Note, error) {
	var note model.Note
	err := r.db.WithContext(ctx).Preload("User").Preload("Folder").First(&note, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &note, nil
}

package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// noteUpdater は [repository.NoteUpdater] の GORM 実装。
// バージョン復元でノート本体を書き戻すためだけに使う。
type noteUpdater struct {
	db *gorm.DB
}

// NewNoteUpdater は NoteUpdater の GORM 実装を返す。
func NewNoteUpdater(db *gorm.DB) repository.NoteUpdater {
	return &noteUpdater{db: db}
}

var _ repository.NoteUpdater = (*noteUpdater)(nil)

// Update はノート本体を書き戻す（既存 NoteRepository.Update と同じ）。
func (r *noteUpdater) Update(ctx context.Context, note *model.Note) error {
	return r.db.WithContext(ctx).Save(note).Error
}

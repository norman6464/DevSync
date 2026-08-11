package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// noteVersionRepository は [repository.NoteVersionRepository] の GORM 実装。
type noteVersionRepository struct {
	db *gorm.DB
}

// NewNoteVersionRepository は NoteVersionRepository の GORM 実装を返す。
func NewNoteVersionRepository(db *gorm.DB) repository.NoteVersionRepository {
	return &noteVersionRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NoteVersionRepository = (*noteVersionRepository)(nil)

// Create は新しいノートバージョンを保存する。
func (r *noteVersionRepository) Create(ctx context.Context, version *model.NoteVersion) error {
	return r.db.WithContext(ctx).Create(version).Error
}

// FindByNoteID は指定ノートのバージョン履歴を新しい順に取得する。
func (r *noteVersionRepository) FindByNoteID(ctx context.Context, noteID uint, limit, offset int) ([]model.NoteVersion, int64, error) {
	var versions []model.NoteVersion
	var total int64
	if err := r.db.WithContext(ctx).Model(&model.NoteVersion{}).
		Where("note_id = ?", noteID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.WithContext(ctx).
		Where("note_id = ?", noteID).
		Order("version_number DESC").Limit(limit).Offset(offset).
		Find(&versions).Error
	return versions, total, err
}

// FindByID は指定 ID のバージョンを取得する。不在の場合は (nil, nil) を返す。
func (r *noteVersionRepository) FindByID(ctx context.Context, id uint) (*model.NoteVersion, error) {
	var version model.NoteVersion
	err := r.db.WithContext(ctx).First(&version, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &version, nil
}

// GetLatestVersionNumber は指定ノートの最新バージョン番号を返す。バージョンがない場合は 0 を返す。
func (r *noteVersionRepository) GetLatestVersionNumber(ctx context.Context, noteID uint) (int, error) {
	var maxVersion int
	err := r.db.WithContext(ctx).Model(&model.NoteVersion{}).
		Where("note_id = ?", noteID).
		Select("COALESCE(MAX(version_number), 0)").
		Scan(&maxVersion).Error
	return maxVersion, err
}

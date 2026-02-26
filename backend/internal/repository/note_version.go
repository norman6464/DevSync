package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// NoteVersionRepository はノートバージョン履歴のデータアクセスを提供する。
type NoteVersionRepository struct {
	db *gorm.DB
}

// NewNoteVersionRepository は新しいNoteVersionRepositoryインスタンスを生成する。
func NewNoteVersionRepository(db *gorm.DB) *NoteVersionRepository {
	return &NoteVersionRepository{db: db}
}

// Create は新しいノートバージョンを保存する。
func (r *NoteVersionRepository) Create(version *model.NoteVersion) error {
	return r.db.Create(version).Error
}

// FindByNoteID は指定ノートのバージョン履歴を新しい順に取得する。
func (r *NoteVersionRepository) FindByNoteID(noteID uint, limit, offset int) ([]model.NoteVersion, int64, error) {
	var versions []model.NoteVersion
	var total int64
	r.db.Model(&model.NoteVersion{}).Where("note_id = ?", noteID).Count(&total)
	err := r.db.Where("note_id = ?", noteID).Order("version_number DESC").Limit(limit).Offset(offset).Find(&versions).Error
	return versions, total, err
}

// FindByID は指定IDのバージョンを取得する。
func (r *NoteVersionRepository) FindByID(id uint) (*model.NoteVersion, error) {
	var version model.NoteVersion
	if err := r.db.First(&version, id).Error; err != nil {
		return nil, err
	}
	return &version, nil
}

// GetLatestVersionNumber は指定ノートの最新バージョン番号を返す。バージョンがない場合は0を返す。
func (r *NoteVersionRepository) GetLatestVersionNumber(noteID uint) (int, error) {
	var maxVersion int
	err := r.db.Model(&model.NoteVersion{}).Where("note_id = ?", noteID).Select("COALESCE(MAX(version_number), 0)").Scan(&maxVersion).Error
	return maxVersion, err
}

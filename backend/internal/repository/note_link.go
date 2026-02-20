package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// NoteLinkRepository はNoteLinkのデータアクセス層。
type NoteLinkRepository struct {
	db *gorm.DB
}

// NewNoteLinkRepository は新しいNoteLinkRepositoryインスタンスを生成する。
func NewNoteLinkRepository(db *gorm.DB) *NoteLinkRepository {
	return &NoteLinkRepository{db: db}
}

// Create は新しいリンクを作成する。
func (r *NoteLinkRepository) Create(link *model.NoteLink) error {
	return r.db.Create(link).Error
}

// FindBySourceNoteID は指定ノートからのリンク一覧を取得する。
// TargetNoteをPreloadして返す。
func (r *NoteLinkRepository) FindBySourceNoteID(sourceNoteID uint) ([]model.NoteLink, error) {
	var links []model.NoteLink
	err := r.db.Preload("TargetNote").
		Where("source_note_id = ?", sourceNoteID).
		Order("created_at DESC").
		Find(&links).Error
	return links, err
}

// FindByTargetNoteID は指定ノートへのリンク一覧（バックリンク）を取得する。
// SourceNoteをPreloadして返す。
func (r *NoteLinkRepository) FindByTargetNoteID(targetNoteID uint) ([]model.NoteLink, error) {
	var links []model.NoteLink
	err := r.db.Preload("SourceNote").
		Where("target_note_id = ?", targetNoteID).
		Order("created_at DESC").
		Find(&links).Error
	return links, err
}

// Delete は指定のリンクを削除する。
func (r *NoteLinkRepository) Delete(sourceNoteID, targetNoteID uint) error {
	return r.db.Where("source_note_id = ? AND target_note_id = ?", sourceNoteID, targetNoteID).
		Delete(&model.NoteLink{}).Error
}

// Exists は指定のリンクが既に存在するかチェックする。
func (r *NoteLinkRepository) Exists(sourceNoteID, targetNoteID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.NoteLink{}).
		Where("source_note_id = ? AND target_note_id = ?", sourceNoteID, targetNoteID).
		Count(&count).Error
	return count > 0, err
}

// CountBySourceNoteID は指定ノートからのリンク数を返す。
func (r *NoteLinkRepository) CountBySourceNoteID(noteID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.NoteLink{}).Where("source_note_id = ?", noteID).Count(&count).Error
	return count, err
}

// CountByTargetNoteID は指定ノートへのリンク数（バックリンク数）を返す。
func (r *NoteLinkRepository) CountByTargetNoteID(noteID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.NoteLink{}).Where("target_note_id = ?", noteID).Count(&count).Error
	return count, err
}

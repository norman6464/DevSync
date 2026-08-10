package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// noteLinkRepository は [repository.NoteLinkRepository] の GORM 実装。
type noteLinkRepository struct {
	db *gorm.DB
}

// NewNoteLinkRepository は NoteLinkRepository の GORM 実装を返す。
func NewNoteLinkRepository(db *gorm.DB) repository.NoteLinkRepository {
	return &noteLinkRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NoteLinkRepository = (*noteLinkRepository)(nil)

// Create は新しいリンクを作成する。
func (r *noteLinkRepository) Create(ctx context.Context, link *model.NoteLink) error {
	return r.db.WithContext(ctx).Create(link).Error
}

// FindBySourceNoteID は指定ノートからのリンク一覧を取得する（リンク先を Preload する）。
func (r *noteLinkRepository) FindBySourceNoteID(ctx context.Context, sourceNoteID uint) ([]model.NoteLink, error) {
	var links []model.NoteLink
	err := r.db.WithContext(ctx).Preload("TargetNote").
		Where("source_note_id = ?", sourceNoteID).
		Order("created_at DESC").
		Find(&links).Error
	return links, err
}

// FindByTargetNoteID は指定ノートへのリンク一覧（バックリンク）を取得する（リンク元を Preload する）。
func (r *noteLinkRepository) FindByTargetNoteID(ctx context.Context, targetNoteID uint) ([]model.NoteLink, error) {
	var links []model.NoteLink
	err := r.db.WithContext(ctx).Preload("SourceNote").
		Where("target_note_id = ?", targetNoteID).
		Order("created_at DESC").
		Find(&links).Error
	return links, err
}

// Delete は指定のリンクを削除する。
func (r *noteLinkRepository) Delete(ctx context.Context, sourceNoteID, targetNoteID uint) error {
	return r.db.WithContext(ctx).
		Where("source_note_id = ? AND target_note_id = ?", sourceNoteID, targetNoteID).
		Delete(&model.NoteLink{}).Error
}

// Exists は指定のリンクが既に存在するかを返す。
func (r *noteLinkRepository) Exists(ctx context.Context, sourceNoteID, targetNoteID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.NoteLink{}).
		Where("source_note_id = ? AND target_note_id = ?", sourceNoteID, targetNoteID).
		Count(&count).Error
	return count > 0, err
}

// CountBySourceNoteID は指定ノートからのリンク数を返す。
func (r *noteLinkRepository) CountBySourceNoteID(ctx context.Context, noteID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.NoteLink{}).
		Where("source_note_id = ?", noteID).Count(&count).Error
	return count, err
}

// CountByTargetNoteID は指定ノートへのリンク数（バックリンク数）を返す。
func (r *noteLinkRepository) CountByTargetNoteID(ctx context.Context, noteID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.NoteLink{}).
		Where("target_note_id = ?", noteID).Count(&count).Error
	return count, err
}

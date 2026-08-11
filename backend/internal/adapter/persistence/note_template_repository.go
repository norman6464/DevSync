package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// noteTemplateRepository は [repository.NoteTemplateRepository] の GORM 実装。
type noteTemplateRepository struct {
	db *gorm.DB
}

// NewNoteTemplateRepository は NoteTemplateRepository の GORM 実装を返す。
func NewNoteTemplateRepository(db *gorm.DB) repository.NoteTemplateRepository {
	return &noteTemplateRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NoteTemplateRepository = (*noteTemplateRepository)(nil)

// Create は新しいテンプレートを作成する。
func (r *noteTemplateRepository) Create(ctx context.Context, template *model.NoteTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

// Update は既存のテンプレートを更新する。
func (r *noteTemplateRepository) Update(ctx context.Context, template *model.NoteTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

// Delete はテンプレートを削除する。
func (r *noteTemplateRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.NoteTemplate{}, id).Error
}

// FindByID は指定 ID のテンプレートを取得する。不在の場合は (nil, nil) を返す。
func (r *noteTemplateRepository) FindByID(ctx context.Context, id uint) (*model.NoteTemplate, error) {
	var template model.NoteTemplate
	err := r.db.WithContext(ctx).First(&template, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// FindByUserID は指定ユーザーの全テンプレートを作成日の新しい順で取得する。
func (r *noteTemplateRepository) FindByUserID(ctx context.Context, userID uint) ([]model.NoteTemplate, error) {
	var templates []model.NoteTemplate
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&templates).Error
	return templates, err
}

// FindDefaultByUserID はデフォルトに設定されたテンプレートを取得する。未設定の場合は (nil, nil) を返す。
func (r *noteTemplateRepository) FindDefaultByUserID(ctx context.Context, userID uint) (*model.NoteTemplate, error) {
	var template model.NoteTemplate
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_default = ?", userID, true).
		First(&template).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// ClearDefaultFlag は指定ユーザーの全テンプレートのデフォルト指定を外す。
func (r *noteTemplateRepository) ClearDefaultFlag(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&model.NoteTemplate{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error
}

// CountByUserID は指定ユーザーのテンプレート総数を返す。
func (r *noteTemplateRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.NoteTemplate{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

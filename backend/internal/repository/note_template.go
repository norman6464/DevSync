package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// NoteTemplateRepository はNoteTemplateのデータアクセス層。
type NoteTemplateRepository struct {
	db *gorm.DB
}

// NewNoteTemplateRepository は新しいNoteTemplateRepositoryインスタンスを生成する。
func NewNoteTemplateRepository(db *gorm.DB) *NoteTemplateRepository {
	return &NoteTemplateRepository{db: db}
}

// Create は新しいテンプレートを作成する。
func (r *NoteTemplateRepository) Create(template *model.NoteTemplate) error {
	return r.db.Create(template).Error
}

// FindByID は指定IDのテンプレートを取得する。
func (r *NoteTemplateRepository) FindByID(id uint) (*model.NoteTemplate, error) {
	var template model.NoteTemplate
	err := r.db.First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// FindByUserID は指定ユーザーの全テンプレートを取得する。
func (r *NoteTemplateRepository) FindByUserID(userID uint) ([]model.NoteTemplate, error) {
	var templates []model.NoteTemplate
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&templates).Error
	return templates, err
}

// FindDefaultByUserID は指定ユーザーのデフォルトテンプレートを取得する。
func (r *NoteTemplateRepository) FindDefaultByUserID(userID uint) (*model.NoteTemplate, error) {
	var template model.NoteTemplate
	err := r.db.Where("user_id = ? AND is_default = ?", userID, true).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// Update はテンプレートを更新する。
func (r *NoteTemplateRepository) Update(template *model.NoteTemplate) error {
	return r.db.Save(template).Error
}

// Delete はテンプレートを削除する。
func (r *NoteTemplateRepository) Delete(id uint) error {
	return r.db.Delete(&model.NoteTemplate{}, id).Error
}

// ClearDefaultFlag は指定ユーザーの全テンプレートのis_defaultをfalseにする。
// 新しいデフォルトテンプレートを設定する前に使用する。
func (r *NoteTemplateRepository) ClearDefaultFlag(userID uint) error {
	return r.db.Model(&model.NoteTemplate{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error
}

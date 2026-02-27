package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// LearningLogTemplateRepository は学習ログテンプレートのデータアクセス層。
type LearningLogTemplateRepository struct {
	db *gorm.DB
}

// NewLearningLogTemplateRepository は新しいLearningLogTemplateRepositoryインスタンスを生成する。
func NewLearningLogTemplateRepository(db *gorm.DB) *LearningLogTemplateRepository {
	return &LearningLogTemplateRepository{db: db}
}

// Create は新しいテンプレートを作成する。
func (r *LearningLogTemplateRepository) Create(template *model.LearningLogTemplate) error {
	return r.db.Create(template).Error
}

// FindByID は指定IDのテンプレートを取得する。
func (r *LearningLogTemplateRepository) FindByID(id uint) (*model.LearningLogTemplate, error) {
	var template model.LearningLogTemplate
	err := r.db.First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// FindByUserID は指定ユーザーの全テンプレートを取得する。
func (r *LearningLogTemplateRepository) FindByUserID(userID uint) ([]model.LearningLogTemplate, error) {
	var templates []model.LearningLogTemplate
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&templates).Error
	return templates, err
}

// FindDefaultByUserID は指定ユーザーのデフォルトテンプレートを取得する。
func (r *LearningLogTemplateRepository) FindDefaultByUserID(userID uint) (*model.LearningLogTemplate, error) {
	var template model.LearningLogTemplate
	err := r.db.Where("user_id = ? AND is_default = ?", userID, true).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// Update はテンプレートを更新する。
func (r *LearningLogTemplateRepository) Update(template *model.LearningLogTemplate) error {
	return r.db.Save(template).Error
}

// Delete はテンプレートを削除する。
func (r *LearningLogTemplateRepository) Delete(id uint) error {
	return r.db.Delete(&model.LearningLogTemplate{}, id).Error
}

// CountByUserID は指定ユーザーのテンプレート総数を返す。
func (r *LearningLogTemplateRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.LearningLogTemplate{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// ClearDefaultFlag は指定ユーザーの全テンプレートのis_defaultをfalseにする。
func (r *LearningLogTemplateRepository) ClearDefaultFlag(userID uint) error {
	return r.db.Model(&model.LearningLogTemplate{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error
}

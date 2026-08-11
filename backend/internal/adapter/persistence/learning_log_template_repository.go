package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// learningLogTemplateRepository は [repository.LearningLogTemplateRepository] の GORM 実装。
type learningLogTemplateRepository struct {
	db *gorm.DB
}

// NewLearningLogTemplateRepository は LearningLogTemplateRepository の GORM 実装を返す。
func NewLearningLogTemplateRepository(db *gorm.DB) repository.LearningLogTemplateRepository {
	return &learningLogTemplateRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.LearningLogTemplateRepository = (*learningLogTemplateRepository)(nil)

// Create は新しいテンプレートを作成する。
func (r *learningLogTemplateRepository) Create(ctx context.Context, template *model.LearningLogTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

// Update は既存のテンプレートを更新する。
func (r *learningLogTemplateRepository) Update(ctx context.Context, template *model.LearningLogTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

// Delete はテンプレートを削除する。
func (r *learningLogTemplateRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.LearningLogTemplate{}, id).Error
}

// FindByID は指定 ID のテンプレートを取得する。不在の場合は (nil, nil) を返す。
func (r *learningLogTemplateRepository) FindByID(ctx context.Context, id uint) (*model.LearningLogTemplate, error) {
	var template model.LearningLogTemplate
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
func (r *learningLogTemplateRepository) FindByUserID(ctx context.Context, userID uint) ([]model.LearningLogTemplate, error) {
	var templates []model.LearningLogTemplate
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&templates).Error
	return templates, err
}

// FindDefaultByUserID はデフォルトに設定されたテンプレートを取得する。未設定の場合は (nil, nil) を返す。
func (r *learningLogTemplateRepository) FindDefaultByUserID(ctx context.Context, userID uint) (*model.LearningLogTemplate, error) {
	var template model.LearningLogTemplate
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
func (r *learningLogTemplateRepository) ClearDefaultFlag(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&model.LearningLogTemplate{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error
}

// CountByUserID は指定ユーザーのテンプレート総数を返す。
func (r *learningLogTemplateRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.LearningLogTemplate{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

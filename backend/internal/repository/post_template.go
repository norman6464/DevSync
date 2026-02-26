package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// PostTemplateRepository は投稿テンプレートデータへのアクセスを提供するリポジトリ実装。
type PostTemplateRepository struct {
	db *gorm.DB
}

// NewPostTemplateRepository は新しいPostTemplateRepositoryインスタンスを生成する。
func NewPostTemplateRepository(db *gorm.DB) *PostTemplateRepository {
	return &PostTemplateRepository{db: db}
}

// Create は新しい投稿テンプレートをデータベースに作成する。
func (r *PostTemplateRepository) Create(template *model.PostTemplate) error {
	return r.db.Create(template).Error
}

// FindByID は指定IDの投稿テンプレートを取得する。
func (r *PostTemplateRepository) FindByID(id uint) (*model.PostTemplate, error) {
	var template model.PostTemplate
	err := r.db.First(&template, id).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// FindByUserID は指定ユーザーの投稿テンプレートをページネーション付きで取得する（新しい順）。
func (r *PostTemplateRepository) FindByUserID(userID uint, limit, offset int) ([]model.PostTemplate, int64, error) {
	var templates []model.PostTemplate
	var total int64
	query := r.db.Where("user_id = ?", userID)
	query.Model(&model.PostTemplate{}).Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&templates).Error
	return templates, total, err
}

// Update は既存の投稿テンプレートを更新する。
func (r *PostTemplateRepository) Update(template *model.PostTemplate) error {
	return r.db.Save(template).Error
}

// Delete は指定IDの投稿テンプレートを削除する。
func (r *PostTemplateRepository) Delete(id uint) error {
	return r.db.Delete(&model.PostTemplate{}, id).Error
}

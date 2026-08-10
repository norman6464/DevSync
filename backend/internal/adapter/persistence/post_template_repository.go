package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// postTemplateRepository は [repository.PostTemplateRepository] の GORM 実装。
type postTemplateRepository struct {
	db *gorm.DB
}

// NewPostTemplateRepository は PostTemplateRepository の GORM 実装を返す。
func NewPostTemplateRepository(db *gorm.DB) repository.PostTemplateRepository {
	return &postTemplateRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.PostTemplateRepository = (*postTemplateRepository)(nil)

// Create は新しい投稿テンプレートをデータベースに作成する。
func (r *postTemplateRepository) Create(ctx context.Context, template *model.PostTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

// FindByID は指定IDの投稿テンプレートを取得する。
func (r *postTemplateRepository) FindByID(ctx context.Context, id uint) (*model.PostTemplate, error) {
	var template model.PostTemplate
	if err := r.db.WithContext(ctx).First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

// FindByUserID は指定ユーザーの投稿テンプレートをページネーション付きで取得する（新しい順）。
func (r *postTemplateRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.PostTemplate, int64, error) {
	var templates []model.PostTemplate
	var total int64
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	query.Model(&model.PostTemplate{}).Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&templates).Error
	return templates, total, err
}

// Update は既存の投稿テンプレートを更新する。
func (r *postTemplateRepository) Update(ctx context.Context, template *model.PostTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

// Delete は指定IDの投稿テンプレートを削除する。
func (r *postTemplateRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.PostTemplate{}, id).Error
}

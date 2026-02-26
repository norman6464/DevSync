package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// ResourceReviewRepository は学習リソースレビューのデータアクセス層。
type ResourceReviewRepository struct {
	db *gorm.DB
}

// NewResourceReviewRepository は新しいResourceReviewRepositoryインスタンスを生成する。
func NewResourceReviewRepository(db *gorm.DB) *ResourceReviewRepository {
	return &ResourceReviewRepository{db: db}
}

// Create は新しいレビューを作成する。
func (r *ResourceReviewRepository) Create(review *model.ResourceReview) error {
	return r.db.Create(review).Error
}

// FindByID は指定IDのレビューをユーザー情報付きで取得する。
func (r *ResourceReviewRepository) FindByID(id uint) (*model.ResourceReview, error) {
	var review model.ResourceReview
	err := r.db.Preload("User").First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// FindByResourceID は指定リソースのレビュー一覧をページネーション付きで取得する。
func (r *ResourceReviewRepository) FindByResourceID(resourceID uint, limit, offset int) ([]model.ResourceReview, int64, error) {
	var reviews []model.ResourceReview
	var total int64

	query := r.db.Model(&model.ResourceReview{}).Where("resource_id = ?", resourceID)
	query.Count(&total)

	err := query.Preload("User").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&reviews).Error

	return reviews, total, err
}

// FindByUserAndResource は指定ユーザーの指定リソースへのレビューを取得する。
func (r *ResourceReviewRepository) FindByUserAndResource(userID, resourceID uint) (*model.ResourceReview, error) {
	var review model.ResourceReview
	err := r.db.Where("user_id = ? AND resource_id = ?", userID, resourceID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// Update はレビューを更新する。
func (r *ResourceReviewRepository) Update(review *model.ResourceReview) error {
	return r.db.Save(review).Error
}

// Delete はレビューを削除する。
func (r *ResourceReviewRepository) Delete(id uint) error {
	return r.db.Delete(&model.ResourceReview{}, id).Error
}

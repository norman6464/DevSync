package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// resourceReviewRepository は [repository.ResourceReviewRepository] の GORM 実装。
type resourceReviewRepository struct {
	db *gorm.DB
}

// NewResourceReviewRepository は ResourceReviewRepository の GORM 実装を返す。
func NewResourceReviewRepository(db *gorm.DB) repository.ResourceReviewRepository {
	return &resourceReviewRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ResourceReviewRepository = (*resourceReviewRepository)(nil)

func (r *resourceReviewRepository) Create(ctx context.Context, review *model.ResourceReview) error {
	return r.db.WithContext(ctx).Create(review).Error
}

// FindByID は指定 ID のレビューをユーザー情報付きで取得する。
func (r *resourceReviewRepository) FindByID(ctx context.Context, id uint) (*model.ResourceReview, error) {
	var review model.ResourceReview
	if err := r.db.WithContext(ctx).Preload("User").First(&review, id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

// FindByResourceID は指定リソースのレビュー一覧をページネーション付きで取得する。
func (r *resourceReviewRepository) FindByResourceID(ctx context.Context, resourceID uint, limit, offset int) ([]model.ResourceReview, int64, error) {
	var reviews []model.ResourceReview
	var total int64

	query := r.db.WithContext(ctx).Model(&model.ResourceReview{}).Where("resource_id = ?", resourceID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("User").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&reviews).Error

	return reviews, total, err
}

// FindByUserAndResource は指定ユーザーの指定リソースへのレビューを取得する。
func (r *resourceReviewRepository) FindByUserAndResource(ctx context.Context, userID, resourceID uint) (*model.ResourceReview, error) {
	var review model.ResourceReview
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND resource_id = ?", userID, resourceID).
		First(&review).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *resourceReviewRepository) Update(ctx context.Context, review *model.ResourceReview) error {
	return r.db.WithContext(ctx).Save(review).Error
}

func (r *resourceReviewRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.ResourceReview{}, id).Error
}

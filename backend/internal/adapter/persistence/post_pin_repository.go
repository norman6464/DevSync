package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// postPinRepository は [repository.PostPinRepository] の GORM 実装。
type postPinRepository struct {
	db *gorm.DB
}

// NewPostPinRepository は PostPinRepository の GORM 実装を返す。
func NewPostPinRepository(db *gorm.DB) repository.PostPinRepository {
	return &postPinRepository{db: db}
}

var _ repository.PostPinRepository = (*postPinRepository)(nil)

func (r *postPinRepository) Pin(ctx context.Context, pin *model.PostPin) error {
	return r.db.WithContext(ctx).Create(pin).Error
}

func (r *postPinRepository) Unpin(ctx context.Context, userID, postID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Delete(&model.PostPin{}).Error
}

func (r *postPinRepository) GetByUserID(ctx context.Context, userID uint) ([]model.PostPin, error) {
	var pins []model.PostPin
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("pin_order ASC").
		Preload("Post").Preload("Post.User").
		Find(&pins).Error
	return pins, err
}

func (r *postPinRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.PostPin{}).
		Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *postPinRepository) IsPinned(ctx context.Context, userID, postID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.PostPin{}).
		Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	return count > 0, err
}

func (r *postPinRepository) UpdateOrder(ctx context.Context, userID uint, postIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, postID := range postIDs {
			if err := tx.Model(&model.PostPin{}).
				Where("user_id = ? AND post_id = ?", userID, postID).
				Update("pin_order", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

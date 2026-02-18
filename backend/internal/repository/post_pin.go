package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

type PostPinRepository struct {
	db *gorm.DB
}

func NewPostPinRepository(db *gorm.DB) *PostPinRepository {
	return &PostPinRepository{db: db}
}

func (r *PostPinRepository) Pin(pin *model.PostPin) error {
	return r.db.Create(pin).Error
}

func (r *PostPinRepository) Unpin(userID, postID uint) error {
	return r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&model.PostPin{}).Error
}

func (r *PostPinRepository) GetByUserID(userID uint) ([]model.PostPin, error) {
	var pins []model.PostPin
	err := r.db.Where("user_id = ?", userID).Order("pin_order ASC").Preload("Post").Preload("Post.User").Find(&pins).Error
	return pins, err
}

func (r *PostPinRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.PostPin{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *PostPinRepository) IsPinned(userID, postID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.PostPin{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	return count > 0, err
}

func (r *PostPinRepository) UpdateOrder(userID uint, postIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, postID := range postIDs {
			if err := tx.Model(&model.PostPin{}).Where("user_id = ? AND post_id = ?", userID, postID).Update("pin_order", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

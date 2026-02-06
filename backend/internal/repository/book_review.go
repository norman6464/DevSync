package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

type BookReviewRepository struct {
	db *gorm.DB
}

func NewBookReviewRepository(db *gorm.DB) *BookReviewRepository {
	return &BookReviewRepository{db: db}
}

func (r *BookReviewRepository) Create(review *model.BookReview) error {
	return r.db.Create(review).Error
}

func (r *BookReviewRepository) FindByID(id uint) (*model.BookReview, error) {
	var review model.BookReview
	err := r.db.Preload("User").First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *BookReviewRepository) FindByUserID(userID uint) ([]model.BookReview, error) {
	var reviews []model.BookReview
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

func (r *BookReviewRepository) FindAll(limit, offset int) ([]model.BookReview, int64, error) {
	var reviews []model.BookReview
	var total int64

	r.db.Model(&model.BookReview{}).Count(&total)

	err := r.db.Preload("User").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&reviews).Error

	return reviews, total, err
}

func (r *BookReviewRepository) Update(review *model.BookReview) error {
	return r.db.Save(review).Error
}

func (r *BookReviewRepository) Delete(id uint) error {
	return r.db.Delete(&model.BookReview{}, id).Error
}

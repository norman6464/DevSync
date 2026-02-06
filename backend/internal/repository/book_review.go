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
	return &review, err
}

func (r *BookReviewRepository) FindByUserID(userID uint) ([]model.BookReview, error) {
	var reviews []model.BookReview
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&reviews).Error
	return reviews, err
}

func (r *BookReviewRepository) FindAll(limit, offset int) ([]model.BookReview, error) {
	var reviews []model.BookReview
	err := r.db.Preload("User").Order("created_at DESC").Limit(limit).Offset(offset).Find(&reviews).Error
	return reviews, err
}

func (r *BookReviewRepository) Update(review *model.BookReview) error {
	return r.db.Save(review).Error
}

func (r *BookReviewRepository) Delete(id uint) error {
	return r.db.Delete(&model.BookReview{}, id).Error
}

func (r *BookReviewRepository) IncrementLikeCount(id uint) error {
	return r.db.Model(&model.BookReview{}).Where("id = ?", id).Update("like_count", gorm.Expr("like_count + 1")).Error
}

func (r *BookReviewRepository) DecrementLikeCount(id uint) error {
	return r.db.Model(&model.BookReview{}).Where("id = ?", id).Update("like_count", gorm.Expr("like_count - 1")).Error
}

// Like operations
func (r *BookReviewRepository) CreateLike(like *model.BookReviewLike) error {
	return r.db.Create(like).Error
}

func (r *BookReviewRepository) DeleteLike(userID, reviewID uint) error {
	return r.db.Where("user_id = ? AND book_review_id = ?", userID, reviewID).Delete(&model.BookReviewLike{}).Error
}

func (r *BookReviewRepository) HasLiked(userID, reviewID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.BookReviewLike{}).Where("user_id = ? AND book_review_id = ?", userID, reviewID).Count(&count).Error
	return count > 0, err
}

func (r *BookReviewRepository) GetUserLikedReviewIDs(userID uint, reviewIDs []uint) ([]uint, error) {
	var likedIDs []uint
	err := r.db.Model(&model.BookReviewLike{}).
		Where("user_id = ? AND book_review_id IN ?", userID, reviewIDs).
		Pluck("book_review_id", &likedIDs).Error
	return likedIDs, err
}

func (r *BookReviewRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.BookReview{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

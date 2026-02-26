package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// BookReviewRepository は書籍レビューデータへのアクセスを提供するリポジトリ実装。
type BookReviewRepository struct {
	db *gorm.DB
}

// NewBookReviewRepository は新しいBookReviewRepositoryインスタンスを生成する。
func NewBookReviewRepository(db *gorm.DB) *BookReviewRepository {
	return &BookReviewRepository{db: db}
}

// Create は新しい書籍レビューをデータベースに作成する。
func (r *BookReviewRepository) Create(review *model.BookReview) error {
	return r.db.Create(review).Error
}

// FindByID は指定IDの書籍レビューをユーザー情報付きで取得する。
func (r *BookReviewRepository) FindByID(id uint) (*model.BookReview, error) {
	var review model.BookReview
	err := r.db.Preload("User").First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// FindByUserID は指定ユーザーの書籍レビューをページネーション付きで取得する（新しい順）。
func (r *BookReviewRepository) FindByUserID(userID uint, limit, offset int) ([]model.BookReview, int64, error) {
	var reviews []model.BookReview
	var total int64
	query := r.db.Where("user_id = ?", userID)
	query.Model(&model.BookReview{}).Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&reviews).Error
	return reviews, total, err
}

// FindAll は全書籍レビューをページネーション付きで取得する。
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

// FindByRating は指定ユーザーの書籍レビューを評価範囲でフィルタリングして取得する（新しい順）。
func (r *BookReviewRepository) FindByRating(userID uint, minRating, maxRating int) ([]model.BookReview, error) {
	var reviews []model.BookReview
	err := r.db.Where("user_id = ? AND rating >= ? AND rating <= ?", userID, minRating, maxRating).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

// Search は書籍レビューをタイトル・著者名・ISBNからキーワード検索する（新しい順）。
func (r *BookReviewRepository) Search(query string, limit, offset int) ([]model.BookReview, int64, error) {
	var reviews []model.BookReview
	var total int64
	like := EscapeLikePattern(query)
	q := r.db.Where("title ILIKE ? OR author ILIKE ? OR isbn ILIKE ?", like, like, like)
	q.Model(&model.BookReview{}).Count(&total)
	err := q.Preload("User").Order("created_at DESC").Limit(limit).Offset(offset).Find(&reviews).Error
	return reviews, total, err
}

// Update は既存の書籍レビューを更新する。
func (r *BookReviewRepository) Update(review *model.BookReview) error {
	return r.db.Save(review).Error
}

// Delete は指定IDの書籍レビューを削除する。
func (r *BookReviewRepository) Delete(id uint) error {
	return r.db.Delete(&model.BookReview{}, id).Error
}

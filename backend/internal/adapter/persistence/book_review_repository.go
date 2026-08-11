package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// bookReviewRepository は [repository.BookReviewRepository] の GORM 実装。
// BookReview は論理削除（DeletedAt）を使うため、GORM が削除済み行を自動的に除外する。
type bookReviewRepository struct {
	db *gorm.DB
}

// NewBookReviewRepository は BookReviewRepository の GORM 実装を返す。
func NewBookReviewRepository(db *gorm.DB) repository.BookReviewRepository {
	return &bookReviewRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.BookReviewRepository = (*bookReviewRepository)(nil)

// Create は新しい書籍レビューを作成する。
func (r *bookReviewRepository) Create(ctx context.Context, review *model.BookReview) error {
	return r.db.WithContext(ctx).Create(review).Error
}

// FindByID は指定 ID のレビューをユーザー情報付きで取得する。不在の場合は (nil, nil) を返す。
func (r *bookReviewRepository) FindByID(ctx context.Context, id uint) (*model.BookReview, error) {
	var review model.BookReview
	err := r.db.WithContext(ctx).Preload("User").First(&review, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// FindByUserID は指定ユーザーのレビューをページネーション付きで取得する（新しい順）。
func (r *bookReviewRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.BookReview, int64, error) {
	db := r.db.WithContext(ctx)

	var total int64
	if err := db.Model(&model.BookReview{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var reviews []model.BookReview
	err := db.Where("user_id = ?", userID).
		Order("created_at DESC").Limit(limit).Offset(offset).
		Find(&reviews).Error
	return reviews, total, err
}

// FindAll は全レビューをページネーション付きで取得する（新しい順）。
func (r *bookReviewRepository) FindAll(ctx context.Context, limit, offset int) ([]model.BookReview, int64, error) {
	db := r.db.WithContext(ctx)

	var total int64
	if err := db.Model(&model.BookReview{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var reviews []model.BookReview
	err := db.Preload("User").
		Order("created_at DESC").Limit(limit).Offset(offset).
		Find(&reviews).Error
	return reviews, total, err
}

// FindByRating は指定ユーザーのレビューを評価範囲で絞り込んで取得する（新しい順）。
func (r *bookReviewRepository) FindByRating(ctx context.Context, userID uint, minRating, maxRating int) ([]model.BookReview, error) {
	var reviews []model.BookReview
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND rating >= ? AND rating <= ?", userID, minRating, maxRating).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}

// Search はタイトル・著者名・ISBN からレビューをキーワード検索する（新しい順）。
func (r *bookReviewRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.BookReview, int64, error) {
	db := r.db.WithContext(ctx)
	like := escapeLikePattern(query)
	const cond = "title ILIKE ? OR author ILIKE ? OR isbn ILIKE ?"

	var total int64
	if err := db.Model(&model.BookReview{}).Where(cond, like, like, like).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var reviews []model.BookReview
	err := db.Where(cond, like, like, like).
		Preload("User").Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&reviews).Error
	return reviews, total, err
}

// Update は既存のレビューを更新する。
func (r *bookReviewRepository) Update(ctx context.Context, review *model.BookReview) error {
	return r.db.WithContext(ctx).Save(review).Error
}

// Delete はレビューを削除する（論理削除）。
func (r *bookReviewRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.BookReview{}, id).Error
}

// CountByUserID は指定ユーザーのレビュー総数を返す。
func (r *bookReviewRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.BookReview{}).
		Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

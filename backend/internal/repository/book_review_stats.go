package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// BookReviewStatsRepository はユーザー書籍レビュー集計統計の取得を担当するリポジトリ実装。
type BookReviewStatsRepository struct {
	db *gorm.DB
}

// NewBookReviewStatsRepository は新しいBookReviewStatsRepositoryインスタンスを生成する。
func NewBookReviewStatsRepository(db *gorm.DB) *BookReviewStatsRepository {
	return &BookReviewStatsRepository{db: db}
}

// GetBookReviewStats は指定ユーザーの書籍レビュー集計統計を返す。
func (r *BookReviewStatsRepository) GetBookReviewStats(userID uint) (*model.BookReviewStats, error) {
	var stats model.BookReviewStats

	// 総レビュー数
	if err := r.db.Model(&model.BookReview{}).Where("user_id = ?", userID).Count(&stats.TotalReviews).Error; err != nil {
		return nil, err
	}

	if stats.TotalReviews == 0 {
		return &stats, nil
	}

	// 平均評価
	var avgRating *float64
	if err := r.db.Model(&model.BookReview{}).Where("user_id = ?", userID).Select("AVG(rating)").Scan(&avgRating).Error; err != nil {
		return nil, err
	}
	if avgRating != nil {
		stats.AverageRating = *avgRating
	}

	// 最高評価
	var maxRating *int
	if err := r.db.Model(&model.BookReview{}).Where("user_id = ?", userID).Select("MAX(rating)").Scan(&maxRating).Error; err != nil {
		return nil, err
	}
	if maxRating != nil {
		stats.MaxRating = *maxRating
	}

	// 最低評価
	var minRating *int
	if err := r.db.Model(&model.BookReview{}).Where("user_id = ?", userID).Select("MIN(rating)").Scan(&minRating).Error; err != nil {
		return nil, err
	}
	if minRating != nil {
		stats.MinRating = *minRating
	}

	// 5つ星レビュー数
	if err := r.db.Model(&model.BookReview{}).Where("user_id = ? AND rating = ?", userID, 5).Count(&stats.FiveStarCount).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

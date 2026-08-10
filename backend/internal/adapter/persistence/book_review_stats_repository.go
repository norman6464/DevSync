package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// bookReviewStatsRepository は [repository.BookReviewStatsRepository] の GORM 実装。
type bookReviewStatsRepository struct {
	db *gorm.DB
}

// NewBookReviewStatsRepository は BookReviewStatsRepository の GORM 実装を返す。
func NewBookReviewStatsRepository(db *gorm.DB) repository.BookReviewStatsRepository {
	return &bookReviewStatsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.BookReviewStatsRepository = (*bookReviewStatsRepository)(nil)

// GetBookReviewStats は指定ユーザーの書籍レビュー集計統計を返す。
// レビューが 0 件の場合は集計クエリを実行せずゼロ値を返す。
func (r *bookReviewStatsRepository) GetBookReviewStats(ctx context.Context, userID uint) (*model.BookReviewStats, error) {
	db := r.db.WithContext(ctx)
	var stats model.BookReviewStats

	// 総レビュー数
	if err := db.Model(&model.BookReview{}).Where("user_id = ?", userID).Count(&stats.TotalReviews).Error; err != nil {
		return nil, err
	}

	if stats.TotalReviews == 0 {
		return &stats, nil
	}

	// 平均評価
	var avgRating *float64
	if err := db.Model(&model.BookReview{}).Where("user_id = ?", userID).Select("AVG(rating)").Scan(&avgRating).Error; err != nil {
		return nil, err
	}
	if avgRating != nil {
		stats.AverageRating = *avgRating
	}

	// 最高評価
	var maxRating *int
	if err := db.Model(&model.BookReview{}).Where("user_id = ?", userID).Select("MAX(rating)").Scan(&maxRating).Error; err != nil {
		return nil, err
	}
	if maxRating != nil {
		stats.MaxRating = *maxRating
	}

	// 最低評価
	var minRating *int
	if err := db.Model(&model.BookReview{}).Where("user_id = ?", userID).Select("MIN(rating)").Scan(&minRating).Error; err != nil {
		return nil, err
	}
	if minRating != nil {
		stats.MinRating = *minRating
	}

	// 5つ星レビュー数
	if err := db.Model(&model.BookReview{}).Where("user_id = ? AND rating = ?", userID, 5).Count(&stats.FiveStarCount).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

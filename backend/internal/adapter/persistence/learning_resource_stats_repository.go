package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// learningResourceStatsRepository は [repository.LearningResourceStatsRepository] の GORM 実装。
type learningResourceStatsRepository struct {
	db *gorm.DB
}

// NewLearningResourceStatsRepository は LearningResourceStatsRepository の GORM 実装を返す。
func NewLearningResourceStatsRepository(db *gorm.DB) repository.LearningResourceStatsRepository {
	return &learningResourceStatsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.LearningResourceStatsRepository = (*learningResourceStatsRepository)(nil)

// GetLearningResourceStats は指定ユーザーの学習リソース活動集計統計を返す。
func (r *learningResourceStatsRepository) GetLearningResourceStats(ctx context.Context, userID uint) (*model.LearningResourceStats, error) {
	db := r.db.WithContext(ctx)
	var stats model.LearningResourceStats

	// 総リソース数
	if err := db.Model(&model.LearningResource{}).Where("user_id = ?", userID).Count(&stats.TotalResources).Error; err != nil {
		return nil, err
	}

	// いいね総数（SUM of like_count）
	var totalLikes *int64
	if err := db.Model(&model.LearningResource{}).Where("user_id = ?", userID).Select("SUM(like_count)").Scan(&totalLikes).Error; err != nil {
		return nil, err
	}
	if totalLikes != nil {
		stats.TotalLikes = *totalLikes
	}

	// 保存総数（SUM of save_count）
	var totalSaves *int64
	if err := db.Model(&model.LearningResource{}).Where("user_id = ?", userID).Select("SUM(save_count)").Scan(&totalSaves).Error; err != nil {
		return nil, err
	}
	if totalSaves != nil {
		stats.TotalSaves = *totalSaves
	}

	// 使用カテゴリ数（DISTINCT category）
	if err := db.Model(&model.LearningResource{}).Where("user_id = ?", userID).Distinct("category").Count(&stats.CategoryCount).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

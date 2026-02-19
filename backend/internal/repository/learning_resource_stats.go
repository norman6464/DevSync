package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// LearningResourceStatsRepository はユーザー学習リソース活動集計統計の取得を担当するリポジトリ実装。
type LearningResourceStatsRepository struct {
	db *gorm.DB
}

// NewLearningResourceStatsRepository は新しいLearningResourceStatsRepositoryインスタンスを生成する。
func NewLearningResourceStatsRepository(db *gorm.DB) *LearningResourceStatsRepository {
	return &LearningResourceStatsRepository{db: db}
}

// GetLearningResourceStats は指定ユーザーの学習リソース活動集計統計を返す。
func (r *LearningResourceStatsRepository) GetLearningResourceStats(userID uint) (*model.LearningResourceStats, error) {
	var stats model.LearningResourceStats

	// 総リソース数
	if err := r.db.Model(&model.LearningResource{}).Where("user_id = ?", userID).Count(&stats.TotalResources).Error; err != nil {
		return nil, err
	}

	// いいね総数（SUM of like_count）
	var totalLikes *int64
	if err := r.db.Model(&model.LearningResource{}).Where("user_id = ?", userID).Select("SUM(like_count)").Scan(&totalLikes).Error; err != nil {
		return nil, err
	}
	if totalLikes != nil {
		stats.TotalLikes = *totalLikes
	}

	// 保存総数（SUM of save_count）
	var totalSaves *int64
	if err := r.db.Model(&model.LearningResource{}).Where("user_id = ?", userID).Select("SUM(save_count)").Scan(&totalSaves).Error; err != nil {
		return nil, err
	}
	if totalSaves != nil {
		stats.TotalSaves = *totalSaves
	}

	// 使用カテゴリ数（DISTINCT category）
	if err := r.db.Model(&model.LearningResource{}).Where("user_id = ?", userID).Distinct("category").Count(&stats.CategoryCount).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

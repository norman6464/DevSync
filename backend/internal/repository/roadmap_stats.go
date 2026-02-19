package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// RoadmapStatsRepository はユーザーロードマップ統計の取得を担当するリポジトリ実装。
type RoadmapStatsRepository struct {
	db *gorm.DB
}

// NewRoadmapStatsRepository は新しいRoadmapStatsRepositoryインスタンスを生成する。
func NewRoadmapStatsRepository(db *gorm.DB) *RoadmapStatsRepository {
	return &RoadmapStatsRepository{db: db}
}

// GetRoadmapStats は指定ユーザーのロードマップ統計を返す。
func (r *RoadmapStatsRepository) GetRoadmapStats(userID uint) (*model.RoadmapStats, error) {
	var stats model.RoadmapStats

	var totalCount int64
	if err := r.db.Model(&model.Roadmap{}).Where("user_id = ?", userID).Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats.TotalRoadmaps = int(totalCount)

	var activeCount int64
	if err := r.db.Model(&model.Roadmap{}).Where("user_id = ? AND status = ?", userID, model.RoadmapStatusActive).Count(&activeCount).Error; err != nil {
		return nil, err
	}
	stats.ActiveRoadmaps = int(activeCount)

	var completedCount int64
	if err := r.db.Model(&model.Roadmap{}).Where("user_id = ? AND status = ?", userID, model.RoadmapStatusCompleted).Count(&completedCount).Error; err != nil {
		return nil, err
	}
	stats.CompletedRoadmaps = int(completedCount)

	var totalSteps int64
	if err := r.db.Model(&model.Roadmap{}).Where("user_id = ?", userID).Select("COALESCE(SUM(step_count), 0)").Scan(&totalSteps).Error; err != nil {
		return nil, err
	}
	stats.TotalSteps = int(totalSteps)

	var completedSteps int64
	if err := r.db.Model(&model.Roadmap{}).Where("user_id = ?", userID).Select("COALESCE(SUM(completed_step_count), 0)").Scan(&completedSteps).Error; err != nil {
		return nil, err
	}
	stats.CompletedSteps = int(completedSteps)

	return &stats, nil
}

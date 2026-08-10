package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// roadmapStatsRepository は [repository.RoadmapStatsRepository] の GORM 実装。
type roadmapStatsRepository struct {
	db *gorm.DB
}

// NewRoadmapStatsRepository は RoadmapStatsRepository の GORM 実装を返す。
func NewRoadmapStatsRepository(db *gorm.DB) repository.RoadmapStatsRepository {
	return &roadmapStatsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.RoadmapStatsRepository = (*roadmapStatsRepository)(nil)

// GetRoadmapStats は指定ユーザーのロードマップ統計を返す。
func (r *roadmapStatsRepository) GetRoadmapStats(ctx context.Context, userID uint) (*model.RoadmapStats, error) {
	db := r.db.WithContext(ctx)
	var stats model.RoadmapStats

	var totalCount int64
	if err := db.Model(&model.Roadmap{}).Where("user_id = ?", userID).Count(&totalCount).Error; err != nil {
		return nil, err
	}
	stats.TotalRoadmaps = int(totalCount)

	var activeCount int64
	if err := db.Model(&model.Roadmap{}).Where("user_id = ? AND status = ?", userID, model.RoadmapStatusActive).Count(&activeCount).Error; err != nil {
		return nil, err
	}
	stats.ActiveRoadmaps = int(activeCount)

	var completedCount int64
	if err := db.Model(&model.Roadmap{}).Where("user_id = ? AND status = ?", userID, model.RoadmapStatusCompleted).Count(&completedCount).Error; err != nil {
		return nil, err
	}
	stats.CompletedRoadmaps = int(completedCount)

	var totalSteps int64
	if err := db.Model(&model.Roadmap{}).Where("user_id = ?", userID).Select("COALESCE(SUM(step_count), 0)").Scan(&totalSteps).Error; err != nil {
		return nil, err
	}
	stats.TotalSteps = int(totalSteps)

	var completedSteps int64
	if err := db.Model(&model.Roadmap{}).Where("user_id = ?", userID).Select("COALESCE(SUM(completed_step_count), 0)").Scan(&completedSteps).Error; err != nil {
		return nil, err
	}
	stats.CompletedSteps = int(completedSteps)

	return &stats, nil
}

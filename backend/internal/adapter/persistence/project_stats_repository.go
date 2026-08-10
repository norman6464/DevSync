package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// projectStatsRepository は [repository.ProjectStatsRepository] の GORM 実装。
type projectStatsRepository struct {
	db *gorm.DB
}

// NewProjectStatsRepository は ProjectStatsRepository の GORM 実装を返す。
func NewProjectStatsRepository(db *gorm.DB) repository.ProjectStatsRepository {
	return &projectStatsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ProjectStatsRepository = (*projectStatsRepository)(nil)

// GetProjectStats は指定ユーザーのプロジェクト活動集計統計を返す。
func (r *projectStatsRepository) GetProjectStats(ctx context.Context, userID uint) (*model.ProjectStats, error) {
	db := r.db.WithContext(ctx)
	var stats model.ProjectStats

	// 総プロジェクト数
	if err := db.Model(&model.Project{}).Where("user_id = ?", userID).Count(&stats.TotalProjects).Error; err != nil {
		return nil, err
	}

	// 注目プロジェクト数
	if err := db.Model(&model.Project{}).Where("user_id = ? AND featured = ?", userID, true).Count(&stats.FeaturedProjects).Error; err != nil {
		return nil, err
	}

	// 進行中プロジェクト数（end_date が NULL）
	if err := db.Model(&model.Project{}).Where("user_id = ? AND end_date IS NULL", userID).Count(&stats.OngoingProjects).Error; err != nil {
		return nil, err
	}

	// 完了プロジェクト数（end_date が NOT NULL）
	if err := db.Model(&model.Project{}).Where("user_id = ? AND end_date IS NOT NULL", userID).Count(&stats.CompletedProjects).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

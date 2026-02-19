package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// ProjectStatsRepository はユーザープロジェクト活動集計統計の取得を担当するリポジトリ実装。
type ProjectStatsRepository struct {
	db *gorm.DB
}

// NewProjectStatsRepository は新しいProjectStatsRepositoryインスタンスを生成する。
func NewProjectStatsRepository(db *gorm.DB) *ProjectStatsRepository {
	return &ProjectStatsRepository{db: db}
}

// GetProjectStats は指定ユーザーのプロジェクト活動集計統計を返す。
func (r *ProjectStatsRepository) GetProjectStats(userID uint) (*model.ProjectStats, error) {
	var stats model.ProjectStats

	// 総プロジェクト数
	if err := r.db.Model(&model.Project{}).Where("user_id = ?", userID).Count(&stats.TotalProjects).Error; err != nil {
		return nil, err
	}

	// 注目プロジェクト数
	if err := r.db.Model(&model.Project{}).Where("user_id = ? AND featured = ?", userID, true).Count(&stats.FeaturedProjects).Error; err != nil {
		return nil, err
	}

	// 進行中プロジェクト数（end_dateがNULL）
	if err := r.db.Model(&model.Project{}).Where("user_id = ? AND end_date IS NULL", userID).Count(&stats.OngoingProjects).Error; err != nil {
		return nil, err
	}

	// 完了プロジェクト数（end_dateがNOT NULL）
	if err := r.db.Model(&model.Project{}).Where("user_id = ? AND end_date IS NOT NULL", userID).Count(&stats.CompletedProjects).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// LearningLogStatsRepository はユーザー学習ログ集計統計の取得を担当するリポジトリ実装。
type LearningLogStatsRepository struct {
	db *gorm.DB
}

// NewLearningLogStatsRepository は新しいLearningLogStatsRepositoryインスタンスを生成する。
func NewLearningLogStatsRepository(db *gorm.DB) *LearningLogStatsRepository {
	return &LearningLogStatsRepository{db: db}
}

// GetLearningLogStats は指定ユーザーの学習ログ集計統計を返す。
func (r *LearningLogStatsRepository) GetLearningLogStats(userID uint) (*model.LearningLogStats, error) {
	var stats model.LearningLogStats

	// 総ログ数
	if err := r.db.Model(&model.LearningLog{}).Where("user_id = ?", userID).Count(&stats.TotalLogs).Error; err != nil {
		return nil, err
	}

	// 総学習時間（分単位）
	if err := r.db.Model(&model.LearningLog{}).Where("user_id = ?", userID).Select("COALESCE(SUM(duration), 0)").Scan(&stats.TotalDuration).Error; err != nil {
		return nil, err
	}

	// カテゴリ数
	if err := r.db.Model(&model.LearningLog{}).Where("user_id = ?", userID).Select("COUNT(DISTINCT category)").Scan(&stats.CategoryCount).Error; err != nil {
		return nil, err
	}

	// 今月のログ数
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	if err := r.db.Model(&model.LearningLog{}).Where("user_id = ? AND created_at >= ?", userID, monthStart).Count(&stats.LogsThisMonth).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

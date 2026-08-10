package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// learningLogStatsRepository は [repository.LearningLogStatsRepository] の GORM 実装。
type learningLogStatsRepository struct {
	db *gorm.DB
}

// NewLearningLogStatsRepository は LearningLogStatsRepository の GORM 実装を返す。
func NewLearningLogStatsRepository(db *gorm.DB) repository.LearningLogStatsRepository {
	return &learningLogStatsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.LearningLogStatsRepository = (*learningLogStatsRepository)(nil)

// GetLearningLogStats は指定ユーザーの学習ログ集計統計を返す。
func (r *learningLogStatsRepository) GetLearningLogStats(ctx context.Context, userID uint) (*model.LearningLogStats, error) {
	db := r.db.WithContext(ctx)
	var stats model.LearningLogStats

	// 総ログ数
	if err := db.Model(&model.LearningLog{}).Where("user_id = ?", userID).Count(&stats.TotalLogs).Error; err != nil {
		return nil, err
	}

	// 総学習時間（分単位）
	if err := db.Model(&model.LearningLog{}).Where("user_id = ?", userID).Select("COALESCE(SUM(duration), 0)").Scan(&stats.TotalDuration).Error; err != nil {
		return nil, err
	}

	// カテゴリ数
	if err := db.Model(&model.LearningLog{}).Where("user_id = ?", userID).Select("COUNT(DISTINCT category)").Scan(&stats.CategoryCount).Error; err != nil {
		return nil, err
	}

	// 今月のログ数
	monthStart := domain.StartOfMonth(time.Now())
	if err := db.Model(&model.LearningLog{}).Where("user_id = ? AND created_at >= ?", userID, monthStart).Count(&stats.LogsThisMonth).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// FollowStatsRepository はユーザーフォロー関係集計統計の取得を担当するリポジトリ実装。
type FollowStatsRepository struct {
	db *gorm.DB
}

// NewFollowStatsRepository は新しいFollowStatsRepositoryインスタンスを生成する。
func NewFollowStatsRepository(db *gorm.DB) *FollowStatsRepository {
	return &FollowStatsRepository{db: db}
}

// GetFollowStats は指定ユーザーのフォロー関係集計統計を返す。
func (r *FollowStatsRepository) GetFollowStats(userID uint) (*model.FollowStats, error) {
	var stats model.FollowStats

	// フォロワー数
	if err := r.db.Model(&model.Follow{}).Where("followee_id = ?", userID).Count(&stats.FollowerCount).Error; err != nil {
		return nil, err
	}

	// フォロー数
	if err := r.db.Model(&model.Follow{}).Where("follower_id = ?", userID).Count(&stats.FollowingCount).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

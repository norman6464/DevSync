package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// followStatsRepository は [repository.FollowStatsRepository] の GORM 実装。
type followStatsRepository struct {
	db *gorm.DB
}

// NewFollowStatsRepository は FollowStatsRepository の GORM 実装を返す。
func NewFollowStatsRepository(db *gorm.DB) repository.FollowStatsRepository {
	return &followStatsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.FollowStatsRepository = (*followStatsRepository)(nil)

// GetFollowStats は指定ユーザーのフォロー関係集計統計を返す。
func (r *followStatsRepository) GetFollowStats(ctx context.Context, userID uint) (*model.FollowStats, error) {
	db := r.db.WithContext(ctx)
	var stats model.FollowStats

	// フォロワー数
	if err := db.Model(&model.Follow{}).Where("followee_id = ?", userID).Count(&stats.FollowerCount).Error; err != nil {
		return nil, err
	}

	// フォロー数
	if err := db.Model(&model.Follow{}).Where("follower_id = ?", userID).Count(&stats.FollowingCount).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

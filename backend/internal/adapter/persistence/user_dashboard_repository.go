package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// userDashboardRepository は [repository.UserDashboardRepository] の GORM 実装。
type userDashboardRepository struct {
	db *gorm.DB
}

// NewUserDashboardRepository は UserDashboardRepository の GORM 実装を返す。
func NewUserDashboardRepository(db *gorm.DB) repository.UserDashboardRepository {
	return &userDashboardRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.UserDashboardRepository = (*userDashboardRepository)(nil)

// GetDashboardStats は指定ユーザーのダッシュボード統計情報を集計して返す。
// 投稿数・受信いいね数・受信コメント数・受信閲覧数・フォロワー数・フォロー数を返す。
func (r *userDashboardRepository) GetDashboardStats(ctx context.Context, userID uint) (*model.UserDashboardStats, error) {
	db := r.db.WithContext(ctx)
	stats := &model.UserDashboardStats{}

	// 投稿数（下書き除外）
	if err := db.Model(&model.Post{}).
		Where("user_id = ? AND is_draft = ?", userID, false).
		Count(&stats.PostCount).Error; err != nil {
		return nil, err
	}

	// 受信いいね数
	if err := db.Model(&model.Post{}).
		Select("COALESCE(SUM(like_count), 0)").
		Where("user_id = ?", userID).
		Scan(&stats.LikesReceived).Error; err != nil {
		return nil, err
	}

	// 受信コメント数
	if err := db.Model(&model.Post{}).
		Select("COALESCE(SUM(comment_count), 0)").
		Where("user_id = ?", userID).
		Scan(&stats.CommentsReceived).Error; err != nil {
		return nil, err
	}

	// 受信閲覧数
	if err := db.Model(&model.PostView{}).
		Joins("JOIN posts ON posts.id = post_views.post_id").
		Where("posts.user_id = ?", userID).
		Count(&stats.ViewsReceived).Error; err != nil {
		return nil, err
	}

	// フォロワー数
	if err := db.Model(&model.Follow{}).
		Where("followee_id = ?", userID).
		Count(&stats.FollowerCount).Error; err != nil {
		return nil, err
	}

	// フォロー数
	if err := db.Model(&model.Follow{}).
		Where("follower_id = ?", userID).
		Count(&stats.FollowingCount).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

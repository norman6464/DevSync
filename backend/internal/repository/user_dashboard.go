package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// UserDashboardRepository はユーザーダッシュボード統計データへのアクセスを提供する。
type UserDashboardRepository struct {
	db *gorm.DB
}

// NewUserDashboardRepository は新しいUserDashboardRepositoryインスタンスを生成する。
func NewUserDashboardRepository(db *gorm.DB) *UserDashboardRepository {
	return &UserDashboardRepository{db: db}
}

// GetDashboardStats は指定ユーザーのダッシュボード統計情報を集計して返す。
// 投稿数・受信いいね数・受信コメント数・受信閲覧数・フォロワー数・フォロー数を返す。
func (r *UserDashboardRepository) GetDashboardStats(userID uint) (*model.UserDashboardStats, error) {
	stats := &model.UserDashboardStats{}

	// 投稿数（下書き除外）
	if err := r.db.Model(&model.Post{}).
		Where("user_id = ? AND is_draft = ?", userID, false).
		Count(&stats.PostCount).Error; err != nil {
		return nil, err
	}

	// 受信いいね数
	if err := r.db.Model(&model.Post{}).
		Select("COALESCE(SUM(like_count), 0)").
		Where("user_id = ?", userID).
		Scan(&stats.LikesReceived).Error; err != nil {
		return nil, err
	}

	// 受信コメント数
	if err := r.db.Model(&model.Post{}).
		Select("COALESCE(SUM(comment_count), 0)").
		Where("user_id = ?", userID).
		Scan(&stats.CommentsReceived).Error; err != nil {
		return nil, err
	}

	// 受信閲覧数
	if err := r.db.Model(&model.PostView{}).
		Joins("JOIN posts ON posts.id = post_views.post_id").
		Where("posts.user_id = ?", userID).
		Count(&stats.ViewsReceived).Error; err != nil {
		return nil, err
	}

	// フォロワー数
	if err := r.db.Model(&model.Follow{}).
		Where("followee_id = ?", userID).
		Count(&stats.FollowerCount).Error; err != nil {
		return nil, err
	}

	// フォロー数
	if err := r.db.Model(&model.Follow{}).
		Where("follower_id = ?", userID).
		Count(&stats.FollowingCount).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

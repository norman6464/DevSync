package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// FollowRepository はフォロー関係データへのアクセスを提供するリポジトリ実装。
type FollowRepository struct {
	db *gorm.DB
}

// NewFollowRepository は新しいFollowRepositoryインスタンスを生成する。
func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

// Follow はフォロー関係を作成する。
func (r *FollowRepository) Follow(followerID, followeeID uint) error {
	follow := &model.Follow{FollowerID: followerID, FolloweeID: followeeID}
	return r.db.Create(follow).Error
}

// Unfollow はフォロー関係を解除する。
func (r *FollowRepository) Unfollow(followerID, followeeID uint) error {
	return r.db.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Delete(&model.Follow{}).Error
}

// IsFollowing はフォロー関係が存在するかどうかを判定する。
func (r *FollowRepository) IsFollowing(followerID, followeeID uint) bool {
	var count int64
	r.db.Model(&model.Follow{}).Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Count(&count)
	return count > 0
}

// GetFollowers は指定ユーザーのフォロワー一覧をページネーション付きで取得する。
func (r *FollowRepository) GetFollowers(userID uint, limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	r.db.Raw(`SELECT COUNT(*) FROM follows WHERE followee_id = ?`, userID).Scan(&total)
	err := r.db.Raw(`SELECT u.* FROM users u JOIN follows f ON f.follower_id = u.id WHERE f.followee_id = ? ORDER BY f.created_at DESC LIMIT ? OFFSET ?`, userID, limit, offset).Scan(&users).Error
	return users, total, err
}

// GetFollowing は指定ユーザーがフォロー中のユーザー一覧をページネーション付きで取得する。
func (r *FollowRepository) GetFollowing(userID uint, limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	r.db.Raw(`SELECT COUNT(*) FROM follows WHERE follower_id = ?`, userID).Scan(&total)
	err := r.db.Raw(`SELECT u.* FROM users u JOIN follows f ON f.followee_id = u.id WHERE f.follower_id = ? ORDER BY f.created_at DESC LIMIT ? OFFSET ?`, userID, limit, offset).Scan(&users).Error
	return users, total, err
}

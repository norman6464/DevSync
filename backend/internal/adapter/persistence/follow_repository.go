// Package persistence は usecase/repository の port を GORM で実装する adapter 層。
package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// followRepository は [repository.FollowRepository] の GORM 実装。
type followRepository struct {
	db *gorm.DB
}

// NewFollowRepository は FollowRepository の GORM 実装を返す。
func NewFollowRepository(db *gorm.DB) repository.FollowRepository {
	return &followRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.FollowRepository = (*followRepository)(nil)

func (r *followRepository) Follow(ctx context.Context, followerID, followeeID uint) error {
	follow := &model.Follow{FollowerID: followerID, FolloweeID: followeeID}
	return r.db.WithContext(ctx).Create(follow).Error
}

func (r *followRepository) Unfollow(ctx context.Context, followerID, followeeID uint) error {
	return r.db.WithContext(ctx).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Delete(&model.Follow{}).Error
}

func (r *followRepository) GetFollowers(ctx context.Context, userID uint, limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	db := r.db.WithContext(ctx)
	db.Raw(`SELECT COUNT(*) FROM follows WHERE followee_id = ?`, userID).Scan(&total)
	err := db.Raw(
		`SELECT u.* FROM users u JOIN follows f ON f.follower_id = u.id WHERE f.followee_id = ? ORDER BY f.created_at DESC LIMIT ? OFFSET ?`,
		userID, limit, offset,
	).Scan(&users).Error
	return users, total, err
}

func (r *followRepository) GetFollowing(ctx context.Context, userID uint, limit, offset int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	db := r.db.WithContext(ctx)
	db.Raw(`SELECT COUNT(*) FROM follows WHERE follower_id = ?`, userID).Scan(&total)
	err := db.Raw(
		`SELECT u.* FROM users u JOIN follows f ON f.followee_id = u.id WHERE f.follower_id = ? ORDER BY f.created_at DESC LIMIT ? OFFSET ?`,
		userID, limit, offset,
	).Scan(&users).Error
	return users, total, err
}

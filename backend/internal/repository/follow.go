package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

type FollowRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

func (r *FollowRepository) Follow(followerID, followeeID uint) error {
	follow := &model.Follow{FollowerID: followerID, FolloweeID: followeeID}
	return r.db.Create(follow).Error
}

func (r *FollowRepository) Unfollow(followerID, followeeID uint) error {
	return r.db.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Delete(&model.Follow{}).Error
}

func (r *FollowRepository) IsFollowing(followerID, followeeID uint) bool {
	var count int64
	r.db.Model(&model.Follow{}).Where("follower_id = ? AND followee_id = ?", followerID, followeeID).Count(&count)
	return count > 0
}

func (r *FollowRepository) GetFollowers(userID uint) ([]model.User, error) {
	var users []model.User
	err := r.db.Raw(`SELECT u.* FROM users u JOIN follows f ON f.follower_id = u.id WHERE f.followee_id = ?`, userID).Scan(&users).Error
	return users, err
}

func (r *FollowRepository) GetFollowing(userID uint) ([]model.User, error) {
	var users []model.User
	err := r.db.Raw(`SELECT u.* FROM users u JOIN follows f ON f.followee_id = u.id WHERE f.follower_id = ?`, userID).Scan(&users).Error
	return users, err
}

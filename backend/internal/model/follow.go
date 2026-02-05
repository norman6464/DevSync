package model

import "time"

type Follow struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	FollowerID uint      `json:"follower_id" gorm:"not null;uniqueIndex:idx_follower_following"`
	Follower   User      `json:"follower" gorm:"foreignKey:FollowerID"`
	FolloweeID uint      `json:"followee_id" gorm:"not null;uniqueIndex:idx_follower_following"`
	Followee   User      `json:"followee" gorm:"foreignKey:FolloweeID"`
	CreatedAt  time.Time `json:"created_at"`
}

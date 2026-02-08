package model

import "time"

// Follow はユーザー間のフォロー関係を表す。
// uniqueIndex制約で同一ペアの重複フォローを防止する。
type Follow struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	FollowerID uint      `json:"follower_id" gorm:"not null;uniqueIndex:idx_follower_following"` // フォローする側
	Follower   User      `json:"follower" gorm:"foreignKey:FollowerID"`
	FolloweeID uint      `json:"followee_id" gorm:"not null;uniqueIndex:idx_follower_following"` // フォローされる側
	Followee   User      `json:"followee" gorm:"foreignKey:FolloweeID"`
	CreatedAt  time.Time `json:"created_at"`
}

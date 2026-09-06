package model

import "time"

// Follow はユーザー間のフォロー関係を表す。
// uniqueIndex制約で同一ペアの重複フォローを防止する。
type Follow struct {
	ID         uint      `json:"id"`
	FollowerID uint      `json:"follower_id"` // フォローする側
	Follower   User      `json:"follower"`
	FolloweeID uint      `json:"followee_id"` // フォローされる側
	Followee   User      `json:"followee"`
	CreatedAt  time.Time `json:"created_at"`
}

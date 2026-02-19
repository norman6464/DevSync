package model

// FollowStats はユーザーのフォロー関係集計統計を表す。
type FollowStats struct {
	FollowerCount  int64 `json:"follower_count"`
	FollowingCount int64 `json:"following_count"`
}

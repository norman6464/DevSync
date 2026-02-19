package model

// UserDashboardStats はユーザーダッシュボードの統計情報を表す。
// プロフィールページや管理画面で表示する集計データ。
type UserDashboardStats struct {
	PostCount        int64 `json:"post_count"`
	LikesReceived    int64 `json:"likes_received"`
	CommentsReceived int64 `json:"comments_received"`
	ViewsReceived    int64 `json:"views_received"`
	FollowerCount    int64 `json:"follower_count"`
	FollowingCount   int64 `json:"following_count"`
}

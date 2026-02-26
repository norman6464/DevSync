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

// LearningDashboardSummary は学習ダッシュボード統合サマリーを表す。
// 複数APIの結果を一括で取得するためのレスポンス構造体。
type LearningDashboardSummary struct {
	StreakInfo       *StreakInfo          `json:"streak_info"`
	WeeklyMinutes   int                 `json:"weekly_minutes"`
	ActiveGoalCount int                 `json:"active_goal_count"`
	TodayMinutes    int                 `json:"today_minutes"`
	ProductivityScore *ProductivityScore `json:"productivity_score"`
}

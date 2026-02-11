package model

// BadgeStats はバッジ判定に必要な集計統計を保持する構造体。
// BadgeRepository が複数テーブルから集計して返す。
type BadgeStats struct {
	TotalContributions int // GitHub総コントリビューション数
	CurrentStreak      int // GitHub連続コントリビューション日数
	LearningLogStreak  int // 学習ログ連続記録日数
	TotalPosts         int // 投稿総数
	TotalLikesReceived int // 受け取ったいいね総数
	FollowerCount      int // フォロワー数
	FollowingCount     int // フォロー中の数
	QAAnswerCount      int // Q&A回答数
	CompletedGoals     int // 完了した学習目標数
}

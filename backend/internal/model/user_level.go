package model

// XPStats はXP計算に必要な集計統計を保持する構造体。
// LevelRepositoryがDBから集計した結果を格納する。
type XPStats struct {
	LearningLogCount         int // 学習ログ数
	LearningLogTotalDuration int // 学習ログ合計時間（分）
	PostCount                int // 投稿数
	GitHubContributionDays   int // GitHubコントリビューション日数
	CompletedGoals           int // 完了した学習目標数
	CommentCount             int // コメント数
	LikesReceived            int // 受け取ったいいね数
	CurrentStreak            int // 現在の学習ログ連続日数
}

// LevelInfo はユーザーのレベル情報をAPIレスポンス用に表す。
type LevelInfo struct {
	Level           int     `json:"level"`            // 現在のレベル
	TotalXP         int     `json:"total_xp"`         // 累計XP
	CurrentLevelXP  int     `json:"current_level_xp"` // 現在レベルの開始XP
	NextLevelXP     int     `json:"next_level_xp"`    // 次レベルに必要な累計XP
	ProgressXP      int     `json:"progress_xp"`      // 現在レベル内の進捗XP
	ProgressPercent float64 `json:"progress_percent"` // 現在レベル内の進捗率(0-100)
}

// XPBreakdown はXPの内訳をAPIレスポンス用に表す。
type XPBreakdown struct {
	LearningLogs int `json:"learning_logs"` // 学習ログから獲得したXP
	Posts        int `json:"posts"`         // 投稿から獲得したXP
	GitHub       int `json:"github"`        // GitHubコントリビューションから獲得したXP
	Goals        int `json:"goals"`         // 目標完了から獲得したXP
	Comments     int `json:"comments"`      // コメントから獲得したXP
	Likes        int `json:"likes"`         // いいね受取から獲得したXP
	StreakBonus  int `json:"streak_bonus"`  // ストリークボーナスから獲得したXP
	Total        int `json:"total"`         // 合計XP
}

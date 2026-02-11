package model

// HeatmapEntry は曜日×時間帯ごとの学習時間を表す（ヒートマップ用）。
type HeatmapEntry struct {
	DayOfWeek    int `json:"day_of_week"`    // 0=日曜 〜 6=土曜
	Hour         int `json:"hour"`           // 0〜23
	TotalMinutes int `json:"total_minutes"`  // その時間帯の合計学習時間（分）
}

// CategoryBreakdown はカテゴリ別学習時間の内訳を表す。
type CategoryBreakdown struct {
	Category     string  `json:"category"`      // カテゴリ名（coding, reading, course, meetup, other）
	TotalMinutes int     `json:"total_minutes"`  // 合計学習時間（分）
	LogCount     int     `json:"log_count"`      // ログ件数
	Percentage   float64 `json:"percentage"`     // 全体に占める割合（0-100）
}

// WeeklyTrend は週ごとの学習時間推移を表す。
type WeeklyTrend struct {
	WeekStart    string `json:"week_start"`    // 週の開始日（"YYYY-MM-DD"形式）
	TotalMinutes int    `json:"total_minutes"` // その週の合計学習時間（分）
	LogCount     int    `json:"log_count"`     // その週のログ件数
}

// ProductivityStats は生産性スコア計算に必要な統計データを保持する。
// リポジトリがDBから集計した結果を格納する。
type ProductivityStats struct {
	PomodoroSessions int `json:"pomodoro_sessions"` // ポモドーロ完了セッション数
	ManualSessions   int `json:"manual_sessions"`   // 手動記録セッション数
	CompletedGoals   int `json:"completed_goals"`   // 完了した目標数
	TotalGoals       int `json:"total_goals"`       // 全目標数
	CurrentStreak    int `json:"current_streak"`     // 現在のストリーク日数
	LongestStreak    int `json:"longest_streak"`     // 最長ストリーク日数
	TotalLogDays     int `json:"total_log_days"`     // 学習記録がある日数（過去12週）
	TotalDaysInRange int `json:"total_days_in_range"` // 期間内の総日数（過去12週=84日）
}

// ProductivityScore は生産性スコアのAPIレスポンス用構造体。
type ProductivityScore struct {
	PomodoroRate      float64 `json:"pomodoro_rate"`      // ポモドーロ活用率（0-100）
	GoalRate          float64 `json:"goal_rate"`           // 目標達成率（0-100）
	StreakConsistency float64 `json:"streak_consistency"`  // ストリーク継続率（0-100）
	OverallScore      float64 `json:"overall_score"`       // 総合スコア（0-100）
}

// AIInsight はAIが生成した学習インサイトを表す。
// フロントエンド側でi18nキー analytics.insight_{Type} を使ってメッセージを生成する。
type AIInsight struct {
	Type   string                 `json:"type"`   // インサイトの種類（peak_time, category_focus, weekly_growth_up, weekly_growth_down, streak_active, streak_record）
	Params map[string]interface{} `json:"params"` // i18n補間用パラメータ
}

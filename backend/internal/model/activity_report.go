package model

import "time"

// ReportPeriod はアクティビティレポートの集計期間を表す型。
type ReportPeriod string

// アクティビティレポートの集計期間定数群。
const (
	ReportPeriodWeekly  ReportPeriod = "weekly"  // 週次
	ReportPeriodMonthly ReportPeriod = "monthly" // 月次
)

// ActivityReport はユーザーの一定期間のアクティビティレポートを表す。
// DBテーブルには対応せず、各種データを集計して生成されるDTO。
// DailyContributions に日別の内訳、TopLanguages に使用言語のランキングを格納する。
type ActivityReport struct {
	Period             ReportPeriod       `json:"period"`     // 集計期間（weekly/monthly）
	StartDate          time.Time          `json:"start_date"` // 集計開始日
	EndDate            time.Time          `json:"end_date"`   // 集計終了日
	UserID             uint               `json:"user_id"`
	TotalContributions int                `json:"total_contributions"` // 総コントリビューション数
	PostsCreated       int                `json:"posts_created"`       // 投稿数
	CommentsCreated    int                `json:"comments_created"`    // コメント数
	LikesReceived      int                `json:"likes_received"`      // 受け取ったいいね数
	GoalsCompleted     int                `json:"goals_completed"`     // 完了した目標数
	GoalsProgress      int                `json:"goals_progress"`      // アクティブ目標の平均進捗率
	NewFollowers       int                `json:"new_followers"`       // 新規フォロワー数
	MessagesExchanged  int                `json:"messages_exchanged"`  // メッセージ送受信数
	DailyContributions []DailyActivity    `json:"daily_contributions"` // 日別アクティビティの内訳（チャート用）
	TopLanguages       []LanguageActivity `json:"top_languages"`       // 使用言語ランキング
}

// DailyActivity は1日分のアクティビティ内訳を表す。
// チャート表示のために日別の投稿数・コメント数を格納する。
type DailyActivity struct {
	Date          string `json:"date"`          // 日付（"YYYY-MM-DD" 形式）
	Contributions int    `json:"contributions"` // その日の総コントリビューション数
	Posts         int    `json:"posts"`         // その日の投稿数
	Comments      int    `json:"comments"`      // その日のコメント数
}

// LanguageActivity は集計期間内のプログラミング言語使用状況を表す。
// GitHub連携データから算出される。
type LanguageActivity struct {
	Language string `json:"language"` // 言語名
	Bytes    int64  `json:"bytes"`    // コード量（バイト）
	Repos    int    `json:"repos"`    // 使用リポジトリ数
}

// ReportComparison は前期間との比較データを表す。
// 各Diffフィールドは前期間との差分値、TrendPercentage は全体的な活動傾向を示す。
type ReportComparison struct {
	ContributionsDiff int     `json:"contributions_diff"` // コントリビューション数の差分
	PostsDiff         int     `json:"posts_diff"`         // 投稿数の差分
	FollowersDiff     int     `json:"followers_diff"`     // フォロワー数の差分
	GoalsDiff         int     `json:"goals_diff"`         // 目標完了数の差分
	TrendPercentage   float64 `json:"trend_percentage"`   // 全体的な活動傾向（%）
}

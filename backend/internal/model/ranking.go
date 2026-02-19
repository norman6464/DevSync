package model

// RankingEntry はランキング表示用の1エントリを表す。
type RankingEntry struct {
	UserID    uint   `json:"user_id"`    // ユーザーID
	Username  string `json:"username"`   // ユーザー名（URLスラッグ用）
	Name      string `json:"name"`       // 表示名
	AvatarURL string `json:"avatar_url"` // アバター画像URL
	Score     int64  `json:"score"`      // ランキングスコア
}

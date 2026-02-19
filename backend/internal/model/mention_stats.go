package model

// MentionStats はユーザーのメンション活動に関する集計統計を表す。
type MentionStats struct {
	MentionsReceived  int64 `json:"mentions_received"`
	MentionsMade      int64 `json:"mentions_made"`
	MentionsThisMonth int64 `json:"mentions_this_month"`
}

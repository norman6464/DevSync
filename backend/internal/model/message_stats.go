package model

// MessageStats はユーザーのダイレクトメッセージに関する集計統計を表す。
type MessageStats struct {
	TotalSent          int64 `json:"total_sent"`
	TotalReceived      int64 `json:"total_received"`
	ConversationsCount int64 `json:"conversations_count"`
	MessagesThisMonth  int64 `json:"messages_this_month"`
}

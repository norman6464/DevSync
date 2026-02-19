package model

// ReactionStats はユーザーが受け取ったリアクション集計統計を表す。
type ReactionStats struct {
	TotalReactionsReceived int64 `json:"total_reactions_received"`
	UniqueReactors         int64 `json:"unique_reactors"`
	ReactionsThisMonth     int64 `json:"reactions_this_month"`
}

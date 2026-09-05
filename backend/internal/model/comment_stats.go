package model

// CommentStats はユーザーのコメント活動に関する集計統計を表す。
type CommentStats struct {
	TotalComments     int64 `json:"total_comments"`
	TotalReplies      int64 `json:"total_replies"`
	CommentsReceived  int64 `json:"comments_received"`
	CommentsThisMonth int64 `json:"comments_this_month"`
}

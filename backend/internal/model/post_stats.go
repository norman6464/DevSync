package model

// PostStats はユーザーの投稿集計統計を表す。
type PostStats struct {
	TotalPosts         int64 `json:"total_posts"`
	PublishedPosts     int64 `json:"published_posts"`
	DraftPosts         int64 `json:"draft_posts"`
	TotalLikesReceived int64 `json:"total_likes_received"`
	TotalComments      int64 `json:"total_comments"`
	PostsThisMonth     int64 `json:"posts_this_month"`
}

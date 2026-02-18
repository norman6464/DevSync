package dto

// MentionResponse はメンション情報のレスポンス。
type MentionResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	ActorID   uint   `json:"actor_id"`
	ActorName string `json:"actor_name"`
	PostID    *uint  `json:"post_id,omitempty"`
	CommentID *uint  `json:"comment_id,omitempty"`
	CreatedAt string `json:"created_at"`
}

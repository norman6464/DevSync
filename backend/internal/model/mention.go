package model

import "time"

// Mention はテキスト内の @ユーザー名 によるメンションを表す。
// PostID または CommentID のいずれかが設定される。
type Mention struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"` // メンションされたユーザー
	User      User      `json:"user"`
	ActorID   uint      `json:"actor_id"` // メンションしたユーザー
	Actor     User      `json:"actor"`
	PostID    *uint     `json:"post_id,omitempty"` // 関連投稿
	Post      *Post     `json:"post,omitempty"`
	CommentID *uint     `json:"comment_id,omitempty"` // 関連コメント
	CreatedAt time.Time `json:"created_at"`
}

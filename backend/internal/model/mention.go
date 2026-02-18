package model

import "time"

// Mention はテキスト内の @ユーザー名 によるメンションを表す。
// PostID または CommentID のいずれかが設定される。
type Mention struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`            // メンションされたユーザー
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	ActorID   uint      `json:"actor_id" gorm:"not null"`                 // メンションしたユーザー
	Actor     User      `json:"actor" gorm:"foreignKey:ActorID"`
	PostID    *uint     `json:"post_id,omitempty" gorm:"index"`           // 関連投稿
	Post      *Post     `json:"post,omitempty" gorm:"foreignKey:PostID"`
	CommentID *uint     `json:"comment_id,omitempty" gorm:"index"`        // 関連コメント
	CreatedAt time.Time `json:"created_at"`
}

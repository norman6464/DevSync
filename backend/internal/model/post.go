package model

import "time"

// Post はユーザーの投稿（学習報告など）を表す。
type Post struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	User         User      `json:"user" gorm:"foreignKey:UserID"`
	Title        string    `json:"title" gorm:"not null"`
	Content      string    `json:"content" gorm:"type:text;not null"`
	ImageURLs    string    `json:"image_urls" gorm:"type:text"` // カンマ区切りの画像URL
	IsDraft      bool      `json:"is_draft" gorm:"default:false;index"` // 下書きフラグ
	LikeCount    int       `json:"like_count" gorm:"default:0"`
	CommentCount int            `json:"comment_count" gorm:"default:0"`
	CodeSnippets []CodeSnippet  `json:"code_snippets,omitempty" gorm:"foreignKey:PostID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// Like は投稿への「いいね」を記録する。
// uniqueIndex制約でユーザーごとに1投稿1いいねを保証する。
type Like struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_post_like"`
	PostID    uint      `json:"post_id" gorm:"not null;uniqueIndex:idx_user_post_like;index"`
	CreatedAt time.Time `json:"created_at"`
}

// Bookmark は投稿のブックマーク（保存）を記録する。
// uniqueIndex制約でユーザーごとに1投稿1ブックマークを保証する。
type Bookmark struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_post_bookmark"`
	PostID    uint      `json:"post_id" gorm:"not null;uniqueIndex:idx_user_post_bookmark;index"`
	Post      Post      `json:"post,omitempty" gorm:"foreignKey:PostID"`
	CreatedAt time.Time `json:"created_at"`
}

// Reaction は投稿への絵文字リアクションを記録する。
// uniqueIndex制約でユーザーごとに1投稿1絵文字1リアクションを保証する。
type Reaction struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_post_emoji"`
	PostID    uint      `json:"post_id" gorm:"not null;uniqueIndex:idx_user_post_emoji;index"`
	Emoji     string    `json:"emoji" gorm:"not null;uniqueIndex:idx_user_post_emoji;size:10"`
	CreatedAt time.Time `json:"created_at"`
}

// ReactionCount はリアクション種別ごとの集計。
type ReactionCount struct {
	Emoji string `json:"emoji"`
	Count int    `json:"count"`
}

// Comment は投稿へのコメントを表す。
type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	PostID    uint      `json:"post_id" gorm:"not null;index"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

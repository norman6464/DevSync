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

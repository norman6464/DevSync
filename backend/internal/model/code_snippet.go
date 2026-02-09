// Package model はDevSyncアプリケーションのデータモデルを定義する。
package model

import "time"

// CodeSnippet は投稿に添付されるコードスニペットを表す。
// 1投稿に複数のスニペットを添付可能。
type CodeSnippet struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	PostID       uint      `json:"post_id" gorm:"not null;index"`
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	Language     string    `json:"language" gorm:"not null"`
	FileName     string    `json:"file_name"`
	Code         string    `json:"code" gorm:"type:text;not null"`
	CommentCount int       `json:"comment_count" gorm:"default:0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SnippetComment はコードスニペットの特定行へのインラインコメントを表す。
// GitHub PRレビュー風の行単位フィードバックを提供する。
type SnippetComment struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	SnippetID  uint      `json:"snippet_id" gorm:"not null;index"`
	UserID     uint      `json:"user_id" gorm:"not null;index"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	LineNumber int       `json:"line_number" gorm:"not null"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

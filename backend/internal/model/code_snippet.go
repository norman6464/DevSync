// Package model はDevSyncアプリケーションのデータモデルを定義する。
package model

import "time"

// CodeSnippet は投稿に添付されるコードスニペットを表す。
// 1投稿に複数のスニペットを添付可能。
type CodeSnippet struct {
	ID           uint      `json:"id"`
	PostID       uint      `json:"post_id"`
	UserID       uint      `json:"user_id"`
	Language     string    `json:"language"`
	FileName     string    `json:"file_name"`
	Code         string    `json:"code"`
	CommentCount int       `json:"comment_count"`
	ForkedFromID *uint     `json:"forked_from_id"`
	ForkCount    int       `json:"fork_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CodeSnippetFavorite はコードスニペットのお気に入りを記録する。
// ユーザーとスニペットの組み合わせでユニークインデックスを持つ。
type CodeSnippetFavorite struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	SnippetID uint      `json:"snippet_id"`
	CreatedAt time.Time `json:"created_at"`
}

// SnippetComment はコードスニペットの特定行へのインラインコメントを表す。
// GitHub PRレビュー風の行単位フィードバックを提供する。
type SnippetComment struct {
	ID         uint      `json:"id"`
	SnippetID  uint      `json:"snippet_id"`
	UserID     uint      `json:"user_id"`
	User       User      `json:"user"`
	LineNumber int       `json:"line_number"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

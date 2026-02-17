package dto

import "github.com/norman6464/devsync/backend/internal/model"

// CodeSnippetInput はコードスニペットの入力データ。
type CodeSnippetInput struct {
	Language string `json:"language"`
	FileName string `json:"file_name"`
	Code     string `json:"code"`
}

// CreatePostRequest は投稿作成リクエスト。
type CreatePostRequest struct {
	Title        string             `json:"title" binding:"required"`
	Content      string             `json:"content" binding:"required"`
	ImageURLs    string             `json:"image_urls"`
	IsDraft      bool               `json:"is_draft"`
	CodeSnippets []CodeSnippetInput `json:"code_snippets"`
}

// UpdatePostRequest は投稿更新リクエスト。
type UpdatePostRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	ImageURLs string `json:"image_urls"`
}

// CreateCommentRequest はコメント作成リクエスト。
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// PostDetailResponse は投稿詳細レスポンス（いいね済みフラグ付き）。
type PostDetailResponse struct {
	model.Post
	Liked bool `json:"liked"`
}

package dto

import "github.com/norman6464/devsync/backend/internal/model"

// CodeSnippetFavoritesResponse はお気に入りスニペット一覧レスポンス。
type CodeSnippetFavoritesResponse struct {
	Snippets []model.CodeSnippet `json:"snippets"`
	Total    int64               `json:"total"`
}

// CreateCodeSnippetRequest はコードスニペット作成リクエスト。
type CreateCodeSnippetRequest struct {
	Language string `json:"language" binding:"required,max=100" validate:"required,max=100"`
	FileName string `json:"file_name" binding:"omitempty,max=255"`
	Code     string `json:"code" binding:"required,max=50000" validate:"required,max=50000"`
}

// UpdateCodeSnippetRequest はコードスニペット更新リクエスト。
type UpdateCodeSnippetRequest struct {
	Language string `json:"language" binding:"omitempty,max=100"`
	FileName string `json:"file_name" binding:"omitempty,max=255"`
	Code     string `json:"code" binding:"omitempty,max=50000"`
}

// ForkSnippetRequest はスニペットフォークリクエスト。
type ForkSnippetRequest struct {
	TargetPostID uint `json:"target_post_id" binding:"required"`
}

// CreateSnippetCommentRequest はスニペットコメント作成リクエスト。
type CreateSnippetCommentRequest struct {
	LineNumber int    `json:"line_number" binding:"required" validate:"required"`
	Content    string `json:"content" binding:"required,max=5000" validate:"required,max=5000"`
}

package dto

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

// CreateSnippetCommentRequest はスニペットコメント作成リクエスト。
type CreateSnippetCommentRequest struct {
	LineNumber int    `json:"line_number" binding:"required" validate:"required"`
	Content    string `json:"content" binding:"required,max=5000" validate:"required,max=5000"`
}

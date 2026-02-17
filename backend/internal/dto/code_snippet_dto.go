package dto

// CreateCodeSnippetRequest はコードスニペット作成リクエスト。
type CreateCodeSnippetRequest struct {
	Language string `json:"language" binding:"required"`
	FileName string `json:"file_name"`
	Code     string `json:"code" binding:"required"`
}

// UpdateCodeSnippetRequest はコードスニペット更新リクエスト。
type UpdateCodeSnippetRequest struct {
	Language string `json:"language"`
	FileName string `json:"file_name"`
	Code     string `json:"code"`
}

// CreateSnippetCommentRequest はスニペットコメント作成リクエスト。
type CreateSnippetCommentRequest struct {
	LineNumber int    `json:"line_number" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

package model

// CodeSnippetStats はユーザーのコードスニペット活動集計統計を表す。
type CodeSnippetStats struct {
	TotalSnippets int64 `json:"total_snippets"`
	TotalComments int64 `json:"total_comments"`
	LanguageCount int64 `json:"language_count"`
}

package dto

// CreateNoteTemplateRequest はテンプレート作成リクエスト。
type CreateNoteTemplateRequest struct {
	Name            string `json:"name" binding:"required"`
	Description     string `json:"description"`
	DefaultTitle    string `json:"default_title"`
	ContentTemplate string `json:"content_template" binding:"required"`
	DefaultTags     string `json:"default_tags"`
	IsDefault       bool   `json:"is_default"`
}

// UpdateNoteTemplateRequest はテンプレート更新リクエスト。
type UpdateNoteTemplateRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	DefaultTitle    string `json:"default_title"`
	ContentTemplate string `json:"content_template"`
	DefaultTags     string `json:"default_tags"`
	IsDefault       *bool  `json:"is_default"`
}

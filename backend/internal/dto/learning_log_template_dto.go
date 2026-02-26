package dto

// CreateLearningLogTemplateRequest は学習ログテンプレート作成リクエスト。
type CreateLearningLogTemplateRequest struct {
	Name            string `json:"name" binding:"required"`
	DefaultTitle    string `json:"default_title"`
	DefaultContent  string `json:"default_content"`
	DefaultCategory string `json:"default_category"`
	DefaultDuration int    `json:"default_duration"`
	IsDefault       bool   `json:"is_default"`
}

// UpdateLearningLogTemplateRequest は学習ログテンプレート更新リクエスト。
type UpdateLearningLogTemplateRequest struct {
	Name            string `json:"name"`
	DefaultTitle    string `json:"default_title"`
	DefaultContent  string `json:"default_content"`
	DefaultCategory string `json:"default_category"`
	DefaultDuration *int   `json:"default_duration"`
	IsDefault       *bool  `json:"is_default"`
}

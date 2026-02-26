package dto

import "github.com/norman6464/devsync/backend/internal/model"

// CreatePostTemplateRequest は投稿テンプレート作成リクエスト。
type CreatePostTemplateRequest struct {
	Name            string `json:"name" binding:"required,min=1,max=100"`
	TitleTemplate   string `json:"title_template" binding:"omitempty,max=200"`
	ContentTemplate string `json:"content_template" binding:"required,min=1,max=50000"`
}

// UpdatePostTemplateRequest は投稿テンプレート更新リクエスト。
type UpdatePostTemplateRequest struct {
	Name            *string `json:"name" binding:"omitempty,max=100"`
	TitleTemplate   *string `json:"title_template" binding:"omitempty,max=200"`
	ContentTemplate *string `json:"content_template" binding:"omitempty,max=50000"`
}

// PostTemplateListResponse は投稿テンプレート一覧レスポンス（ページネーション付き）。
type PostTemplateListResponse struct {
	Templates []model.PostTemplate `json:"templates"`
	Total     int64                `json:"total"`
	Limit     int                  `json:"limit"`
	Offset    int                  `json:"offset"`
}

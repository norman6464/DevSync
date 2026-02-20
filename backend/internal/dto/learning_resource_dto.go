package dto

import "github.com/norman6464/devsync/backend/internal/model"

// ResourceDetailResponse は学習リソース詳細レスポンス（いいね・保存状態付き）。
type ResourceDetailResponse struct {
	Resource model.LearningResource `json:"resource"`
	HasLiked bool                   `json:"has_liked"`
	HasSaved bool                   `json:"has_saved"`
}

// ResourceListResponse は学習リソース一覧レスポンス。
type ResourceListResponse struct {
	Resources []model.LearningResource `json:"resources"`
	Total     int64                    `json:"total"`
	Limit     int                      `json:"limit"`
	Offset    int                      `json:"offset"`
}

// CreateResourceRequest は学習リソース作成のリクエストボディ。
type CreateResourceRequest struct {
	Title       string `json:"title" binding:"required,max=300"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Category    string `json:"category" binding:"required"`
	Difficulty  string `json:"difficulty"`
	Tags        string `json:"tags"`
	ImageURL    string `json:"image_url"`
	IsPublic    *bool  `json:"is_public"`
}

// UpdateResourceRequest は学習リソース更新のリクエストボディ。
type UpdateResourceRequest struct {
	Title       string `json:"title" binding:"max=300"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Category    string `json:"category"`
	Difficulty  string `json:"difficulty"`
	Tags        string `json:"tags"`
	ImageURL    string `json:"image_url"`
	IsPublic    *bool  `json:"is_public"`
}

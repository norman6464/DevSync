package dto

import "github.com/norman6464/devsync/backend/internal/model"

// UpsertResourceProgressRequest はリソース進捗のUPSERTリクエスト。
type UpsertResourceProgressRequest struct {
	ResourceID        uint   `json:"resource_id" binding:"required"`
	Status            string `json:"status" binding:"required"`
	CompletionPercent int    `json:"completion_percent"`
	Note              string `json:"note"`
}

// ResourceProgressResponse はリソース進捗レスポンス。
type ResourceProgressResponse struct {
	Progress model.ResourceProgress `json:"progress"`
}

// ResourceProgressListResponse はリソース進捗一覧レスポンス。
type ResourceProgressListResponse struct {
	Progresses []model.ResourceProgress `json:"progresses"`
	Total      int64                    `json:"total"`
	Limit      int                      `json:"limit"`
	Offset     int                      `json:"offset"`
}

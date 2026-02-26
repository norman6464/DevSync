package dto

import "github.com/norman6464/devsync/backend/internal/model"

// ResourceReviewListResponse はリソースレビュー一覧レスポンス。
type ResourceReviewListResponse struct {
	Reviews []model.ResourceReview `json:"reviews"`
	Total   int64                  `json:"total"`
}

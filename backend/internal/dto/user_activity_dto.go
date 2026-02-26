package dto

import "github.com/norman6464/devsync/backend/internal/model"

// UserActivityListResponse はアクティビティタイムラインレスポンス。
type UserActivityListResponse struct {
	Activities []model.UserActivity `json:"activities"`
	Total      int64                `json:"total"`
	Limit      int                  `json:"limit"`
	Offset     int                  `json:"offset"`
}

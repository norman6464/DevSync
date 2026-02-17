package dto

import "github.com/norman6464/devsync/backend/internal/model"

// ProjectListResponse はプロジェクト一覧レスポンス。
type ProjectListResponse struct {
	Projects []model.Project `json:"projects"`
	Total    int64           `json:"total"`
	Limit    int             `json:"limit"`
	Offset   int             `json:"offset"`
}

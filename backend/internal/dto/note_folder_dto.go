package dto

import "github.com/norman6464/devsync/backend/internal/model"

// CreateNoteFolderRequest はフォルダ作成リクエスト。
type CreateNoteFolderRequest struct {
	Name     string `json:"name" binding:"required"`
	ParentID *uint  `json:"parent_id"`
}

// UpdateNoteFolderRequest はフォルダ更新リクエスト。
type UpdateNoteFolderRequest struct {
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

// NoteFolderListResponse はフォルダ一覧レスポンス（ページネーション付き）。
type NoteFolderListResponse struct {
	Folders []model.NoteFolder `json:"folders"`
	Total   int64              `json:"total"`
	Limit   int                `json:"limit"`
	Offset  int                `json:"offset"`
}

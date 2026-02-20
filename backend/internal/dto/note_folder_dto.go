package dto

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

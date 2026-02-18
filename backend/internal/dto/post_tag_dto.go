package dto

// SetTagsRequest はタグ設定のリクエスト。
type SetTagsRequest struct {
	Tags []string `json:"tags" binding:"required"`
}

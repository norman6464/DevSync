package dto

// SetTagsRequest はタグ設定のリクエスト。
type SetTagsRequest struct {
	Tags []string `json:"tags" binding:"required,max=10,dive,max=50"`
}

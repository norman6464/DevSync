package dto

// SetTagsRequest はタグ設定のリクエストボディ。
type SetTagsRequest struct {
	Tags []string `json:"tags"`
}

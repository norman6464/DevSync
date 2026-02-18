package dto

// ReorderPinsRequest はピン留め順序変更のリクエスト。
type ReorderPinsRequest struct {
	PostIDs []uint `json:"post_ids" binding:"required,max=100"`
}

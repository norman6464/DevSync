package dto

// CreateNoteLinkRequest はノートリンク作成リクエスト。
type CreateNoteLinkRequest struct {
	TargetNoteID uint `json:"target_note_id" binding:"required"`
}

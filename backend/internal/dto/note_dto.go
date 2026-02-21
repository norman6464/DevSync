package dto

// CreateNoteRequest はノート作成リクエスト。
type CreateNoteRequest struct {
	Title    string `json:"title" binding:"required,min=1,max=200"`
	Content  string `json:"content" binding:"required,min=1,max=50000"`
	Tags     string `json:"tags"`
	FolderID *uint  `json:"folder_id"`
}

// UpdateNoteRequest はノート更新リクエスト。
type UpdateNoteRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Tags     string `json:"tags"`
	FolderID *uint  `json:"folder_id"`
}

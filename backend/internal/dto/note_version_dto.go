package dto

import "github.com/norman6464/devsync/backend/internal/model"

// NoteVersionListResponse はノートバージョン一覧レスポンス（ページネーション付き）。
type NoteVersionListResponse struct {
	Versions []model.NoteVersion `json:"versions"`
	Total    int64               `json:"total"`
	Limit    int                 `json:"limit"`
	Offset   int                 `json:"offset"`
}

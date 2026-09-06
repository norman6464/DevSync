package model

import "time"

// NoteVersion はノートの過去バージョンを保存する。
// ノート更新時に自動的に旧バージョンが保存される。
type NoteVersion struct {
	ID            uint      `json:"id"`
	NoteID        uint      `json:"note_id"`
	VersionNumber int       `json:"version_number"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Tags          string    `json:"tags"`
	CreatedAt     time.Time `json:"created_at"`
}

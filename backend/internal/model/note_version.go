package model

import "time"

// NoteVersion はノートの過去バージョンを保存する。
// ノート更新時に自動的に旧バージョンが保存される。
type NoteVersion struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	NoteID        uint      `gorm:"not null;index" json:"note_id"`
	VersionNumber int       `gorm:"not null" json:"version_number"`
	Title         string    `gorm:"not null" json:"title"`
	Content       string    `gorm:"type:text" json:"content"`
	Tags          string    `gorm:"type:text" json:"tags"`
	CreatedAt     time.Time `json:"created_at"`
}

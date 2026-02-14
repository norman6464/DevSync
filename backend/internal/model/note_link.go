package model

import "time"

// NoteLink はノート間のリンクを表すモデル。
// あるノート（SourceNote）から別のノート（TargetNote）へのリンクを保存する。
type NoteLink struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SourceNoteID uint      `gorm:"not null;index;uniqueIndex:idx_note_link_unique" json:"source_note_id"`
	TargetNoteID uint      `gorm:"not null;index;uniqueIndex:idx_note_link_unique" json:"target_note_id"`
	SourceNote   *Note     `gorm:"foreignKey:SourceNoteID" json:"source_note,omitempty"`
	TargetNote   *Note     `gorm:"foreignKey:TargetNoteID" json:"target_note,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

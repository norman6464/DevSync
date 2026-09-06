package model

import "time"

// NoteLink はノート間のリンクを表すモデル。
// あるノート（SourceNote）から別のノート（TargetNote）へのリンクを保存する。
type NoteLink struct {
	ID           uint      `json:"id"`
	SourceNoteID uint      `json:"source_note_id"`
	TargetNoteID uint      `json:"target_note_id"`
	SourceNote   *Note     `json:"source_note,omitempty"`
	TargetNote   *Note     `json:"target_note,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// NoteLinkStats はノートのリンク統計情報を表す。
type NoteLinkStats struct {
	NoteID           uint  `json:"note_id"`
	ForwardLinkCount int64 `json:"forward_link_count"`
	BacklinkCount    int64 `json:"backlink_count"`
}

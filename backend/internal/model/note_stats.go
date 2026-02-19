package model

// NoteStats はユーザーのノートに関する集計統計を表す。
type NoteStats struct {
	TotalNotes     int64 `json:"total_notes"`
	ArchivedNotes  int64 `json:"archived_notes"`
	FavoriteNotes  int64 `json:"favorite_notes"`
	TotalFolders   int64 `json:"total_folders"`
	NotesThisWeek  int64 `json:"notes_this_week"`
}

package model

import "time"

// Note はユーザーの学習ノートを表す。
// Notion風の学習ノート機能で、マークダウンをサポートする。
type Note struct {
	ID       uint        `json:"id"`
	UserID   uint        `json:"user_id"`
	User     User        `json:"user,omitempty"`
	FolderID *uint       `json:"folder_id,omitempty"`
	Folder   *NoteFolder `json:"folder,omitempty"`

	Title      string `json:"title"`
	Content    string `json:"content"` // マークダウン形式
	Tags       string `json:"tags"`    // カンマ区切りのタグ
	IsFavorite bool   `json:"is_favorite"`
	IsArchived bool   `json:"is_archived"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NoteFolder はノートを整理するフォルダを表す。
// 階層構造をサポートするため、ParentIDを持つ。
type NoteFolder struct {
	ID       uint        `json:"id"`
	UserID   uint        `json:"user_id"`
	User     User        `json:"user,omitempty"`
	ParentID *uint       `json:"parent_id,omitempty"`
	Parent   *NoteFolder `json:"parent,omitempty"`

	Name string `json:"name"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

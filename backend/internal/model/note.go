package model

import "time"

// Note はユーザーの学習ノートを表す。
// Notion風の学習ノート機能で、マークダウンをサポートする。
type Note struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"not null;index" json:"user_id"`
	User       User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	FolderID   *uint      `gorm:"index" json:"folder_id,omitempty"`
	Folder     *NoteFolder `gorm:"foreignKey:FolderID" json:"folder,omitempty"`

	Title      string     `gorm:"not null" json:"title"`
	Content    string     `gorm:"type:text" json:"content"` // マークダウン形式
	Tags       string     `gorm:"type:text" json:"tags"`    // カンマ区切りのタグ
	IsFavorite bool       `gorm:"default:false" json:"is_favorite"`
	IsArchived bool       `gorm:"default:false" json:"is_archived"`

	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// NoteFolder はノートを整理するフォルダを表す。
// 階層構造をサポートするため、ParentIDを持つ。
type NoteFolder struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ParentID  *uint      `gorm:"index" json:"parent_id,omitempty"`
	Parent    *NoteFolder `gorm:"foreignKey:ParentID" json:"parent,omitempty"`

	Name      string     `gorm:"not null" json:"name"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

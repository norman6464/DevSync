package model

import "time"

// NoteTemplate はノートテンプレートのモデル。
// よく使うノートの雛形を保存し、再利用できる。
type NoteTemplate struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"not null;index" json:"user_id"`
	Name            string    `gorm:"type:varchar(100);not null" json:"name"`
	Description     string    `gorm:"type:varchar(500)" json:"description"`
	DefaultTitle    string    `gorm:"type:varchar(200)" json:"default_title"`
	ContentTemplate string    `gorm:"type:text;not null" json:"content_template"`
	DefaultTags     string    `gorm:"type:varchar(255)" json:"default_tags"`
	IsDefault       bool      `gorm:"default:false;index:idx_note_templates_is_default" json:"is_default"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

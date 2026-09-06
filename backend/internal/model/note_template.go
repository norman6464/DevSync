package model

import "time"

// NoteTemplate はノートテンプレートのモデル。
// よく使うノートの雛形を保存し、再利用できる。
type NoteTemplate struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	DefaultTitle    string    `json:"default_title"`
	ContentTemplate string    `json:"content_template"`
	DefaultTags     string    `json:"default_tags"`
	IsDefault       bool      `json:"is_default"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

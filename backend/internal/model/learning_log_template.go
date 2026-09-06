package model

import "time"

// LearningLogTemplate は学習ログテンプレートのモデル。
// よく使う学習ログのパターンを保存し、再利用できる。
type LearningLogTemplate struct {
	ID              uint        `json:"id"`
	UserID          uint        `json:"user_id"`
	Name            string      `json:"name"`
	DefaultTitle    string      `json:"default_title"`
	DefaultContent  string      `json:"default_content"`
	DefaultCategory LogCategory `json:"default_category"`
	DefaultDuration int         `json:"default_duration"`
	IsDefault       bool        `json:"is_default"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

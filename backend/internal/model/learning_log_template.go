package model

import "time"

// LearningLogTemplate は学習ログテンプレートのモデル。
// よく使う学習ログのパターンを保存し、再利用できる。
type LearningLogTemplate struct {
	ID              uint        `gorm:"primaryKey" json:"id"`
	UserID          uint        `gorm:"not null;index" json:"user_id"`
	Name            string      `gorm:"type:varchar(100);not null" json:"name"`
	DefaultTitle    string      `gorm:"type:varchar(200)" json:"default_title"`
	DefaultContent  string      `gorm:"type:text" json:"default_content"`
	DefaultCategory LogCategory `gorm:"type:varchar(50);default:'other'" json:"default_category"`
	DefaultDuration int         `gorm:"default:0" json:"default_duration"`
	IsDefault       bool        `gorm:"default:false;index:idx_log_templates_is_default" json:"is_default"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

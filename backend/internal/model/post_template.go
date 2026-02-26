package model

import "time"

// PostTemplate は投稿テンプレートのモデル。
// よく使う投稿フォーマットを保存し、投稿作成時に再利用できる。
type PostTemplate struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"not null;index" json:"user_id"`
	Name            string    `gorm:"type:varchar(100);not null" json:"name"`
	TitleTemplate   string    `gorm:"type:varchar(200)" json:"title_template"`
	ContentTemplate string    `gorm:"type:text;not null" json:"content_template"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

package model

import "time"

// PostTemplate は投稿テンプレートのモデル。
// よく使う投稿フォーマットを保存し、投稿作成時に再利用できる。
type PostTemplate struct {
	ID              uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	Name            string    `json:"name"`
	TitleTemplate   string    `json:"title_template"`
	ContentTemplate string    `json:"content_template"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

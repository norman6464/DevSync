package model

import "time"

// WidgetSettings はダッシュボードのウィジェット表示設定。
// ユーザーごとにウィジェットの表示/非表示と並び順をJSON形式で保存する。
type WidgetSettings struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Settings  string    `json:"settings"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

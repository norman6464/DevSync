package model

import "time"

// WidgetSettings はダッシュボードのウィジェット表示設定。
// ユーザーごとにウィジェットの表示/非表示と並び順をJSON形式で保存する。
type WidgetSettings struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	Settings  string    `gorm:"type:text;not null" json:"settings"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

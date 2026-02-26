package dto

import "encoding/json"

// UpdateWidgetSettingsRequest はウィジェット設定更新リクエスト。
type UpdateWidgetSettingsRequest struct {
	Settings json.RawMessage `json:"settings" binding:"required"`
}

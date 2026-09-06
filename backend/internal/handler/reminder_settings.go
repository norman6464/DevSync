package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// ReminderSettingsHandler は学習リマインダー設定関連のHTTPハンドラ。
type ReminderSettingsHandler struct {
	getSettings    *usecase.GetReminderSettingsUseCase
	updateSettings *usecase.UpdateReminderSettingsUseCase
}

// NewReminderSettingsHandler は新しいReminderSettingsHandlerインスタンスを生成する。
func NewReminderSettingsHandler(
	getSettings *usecase.GetReminderSettingsUseCase,
	updateSettings *usecase.UpdateReminderSettingsUseCase,
) *ReminderSettingsHandler {
	return &ReminderSettingsHandler{getSettings: getSettings, updateSettings: updateSettings}
}

// GetSettings は認証ユーザーのリマインダー設定を取得する。
func (h *ReminderSettingsHandler) GetSettings(c *gin.Context) {
	settings, err := h.getSettings.Execute(c.Request.Context(), c.GetUint("userID"))
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, settings)
}

// updateReminderSettingsRequest はリマインダー設定更新のリクエストボディ。
type updateReminderSettingsRequest struct {
	Enabled          bool                    `json:"enabled"`
	Frequency        model.ReminderFrequency `json:"frequency" binding:"omitempty,max=50"`
	NotificationTime string                  `json:"notification_time" binding:"omitempty,max=10"`
	InactiveDays     int                     `json:"inactive_days" binding:"omitempty,min=0,max=365"`
	EnableWeb        bool                    `json:"enable_web"`
	EnableEmail      bool                    `json:"enable_email"`
}

// UpdateSettings は認証ユーザーのリマインダー設定を更新する。
func (h *ReminderSettingsHandler) UpdateSettings(c *gin.Context) {
	req := bindJSON[updateReminderSettingsRequest](c)
	if req == nil {
		return
	}

	settings, err := h.updateSettings.Execute(c.Request.Context(), usecase.UpdateReminderSettingsInput{
		UserID:           c.GetUint("userID"),
		Enabled:          req.Enabled,
		Frequency:        req.Frequency,
		NotificationTime: req.NotificationTime,
		InactiveDays:     req.InactiveDays,
		EnableWeb:        req.EnableWeb,
		EnableEmail:      req.EnableEmail,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, settings)
}

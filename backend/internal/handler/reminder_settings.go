package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
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

// UpdateSettings は認証ユーザーのリマインダー設定を更新する。
func (h *ReminderSettingsHandler) UpdateSettings(c *gin.Context) {
	req := bindJSON[dto.UpdateReminderSettingsRequest](c)
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

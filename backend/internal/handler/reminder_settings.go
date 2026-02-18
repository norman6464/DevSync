package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// ReminderSettingsServiceInterface はReminderSettingsHandlerが依存するサービスのインターフェース。
type ReminderSettingsServiceInterface interface {
	GetSettings(userID uint) (*model.ReminderSettings, error)
	UpdateSettings(userID uint, updates *model.ReminderSettings) (*model.ReminderSettings, error)
}

// ReminderSettingsHandler は学習リマインダー設定関連のHTTPハンドラ。
type ReminderSettingsHandler struct {
	service ReminderSettingsServiceInterface
}

// NewReminderSettingsHandler は新しいReminderSettingsHandlerインスタンスを生成する。
func NewReminderSettingsHandler(s ReminderSettingsServiceInterface) *ReminderSettingsHandler {
	return &ReminderSettingsHandler{service: s}
}

// GetSettings は認証ユーザーのリマインダー設定を取得する。
func (h *ReminderSettingsHandler) GetSettings(c *gin.Context) {
	userID := c.GetUint("userID")

	settings, err := h.service.GetSettings(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, settings)
}

// UpdateSettings は認証ユーザーのリマインダー設定を更新する。
func (h *ReminderSettingsHandler) UpdateSettings(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[dto.UpdateReminderSettingsRequest](c)
	if req == nil {
		return
	}

	updates := &model.ReminderSettings{
		Enabled:          req.Enabled,
		Frequency:        req.Frequency,
		NotificationTime: req.NotificationTime,
		InactiveDays:     req.InactiveDays,
		EnableWeb:        req.EnableWeb,
		EnableEmail:      req.EnableEmail,
	}

	settings, err := h.service.UpdateSettings(userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, settings)
}

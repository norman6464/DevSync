package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

// UpdateReminderSettingsInput はリマインダー設定更新リクエストの入力構造体。
type UpdateReminderSettingsInput struct {
	Enabled          bool                      `json:"enabled"`
	Frequency        model.ReminderFrequency   `json:"frequency"`
	NotificationTime string                    `json:"notification_time"`
	InactiveDays     int                       `json:"inactive_days"`
	EnableWeb        bool                      `json:"enable_web"`
	EnableEmail      bool                      `json:"enable_email"`
}

// GetSettings は認証ユーザーのリマインダー設定を取得する。
func (h *ReminderSettingsHandler) GetSettings(c *gin.Context) {
	userID := c.GetUint("userID")

	settings, err := h.service.GetSettings(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateSettings は認証ユーザーのリマインダー設定を更新する。
func (h *ReminderSettingsHandler) UpdateSettings(c *gin.Context) {
	userID := c.GetUint("userID")

	var req UpdateReminderSettingsInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	c.JSON(http.StatusOK, settings)
}

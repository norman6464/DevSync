package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// NotificationSettingsServiceInterface は通知設定サービスの抽象インターフェース。
type NotificationSettingsServiceInterface interface {
	GetSettings(userID uint) (*model.NotificationSettings, error)
	UpdateSettings(userID uint, updates *model.NotificationSettings) (*model.NotificationSettings, error)
}

// NotificationSettingsHandler は通知設定関連のHTTPハンドラ。
type NotificationSettingsHandler struct {
	service NotificationSettingsServiceInterface
}

// NewNotificationSettingsHandler は新しいNotificationSettingsHandlerインスタンスを生成する。
func NewNotificationSettingsHandler(s NotificationSettingsServiceInterface) *NotificationSettingsHandler {
	return &NotificationSettingsHandler{service: s}
}

// GetSettings は認証ユーザーの通知設定を取得する。
func (h *NotificationSettingsHandler) GetSettings(c *gin.Context) {
	userID := c.GetUint("userID")

	settings, err := h.service.GetSettings(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, settings)
}

// UpdateSettings は認証ユーザーの通知設定を更新する。
func (h *NotificationSettingsHandler) UpdateSettings(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[dto.UpdateNotificationSettingsRequest](c)
	if req == nil {
		return
	}

	// リクエストから更新オブジェクトを作成
	updates := &model.NotificationSettings{
		EnableLikes:    req.EnableLikes,
		EnableComments: req.EnableComments,
		EnableFollows:  req.EnableFollows,
		EnableMessages: req.EnableMessages,
		EnableMentions: req.EnableMentions,
		EnableWebPush:  req.EnableWebPush,
		EnableEmail:    req.EnableEmail,
		EnableSound:    req.EnableSound,
	}

	settings, err := h.service.UpdateSettings(userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, settings)
}

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// NotificationSettingsHandler は通知設定関連のHTTPハンドラ。
type NotificationSettingsHandler struct {
	service *service.NotificationSettingsService
}

// NewNotificationSettingsHandler は新しいNotificationSettingsHandlerインスタンスを生成する。
func NewNotificationSettingsHandler(s *service.NotificationSettingsService) *NotificationSettingsHandler {
	return &NotificationSettingsHandler{service: s}
}

// UpdateNotificationSettingsInput は通知設定更新のリクエストボディ。
type UpdateNotificationSettingsInput struct {
	EnableLikes    bool `json:"enable_likes"`
	EnableComments bool `json:"enable_comments"`
	EnableFollows  bool `json:"enable_follows"`
	EnableMessages bool `json:"enable_messages"`
	EnableMentions bool `json:"enable_mentions"`
	EnableWebPush  bool `json:"enable_web_push"`
	EnableEmail    bool `json:"enable_email"`
	EnableSound    bool `json:"enable_sound"`
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

	req := bindJSON[UpdateNotificationSettingsInput](c)
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

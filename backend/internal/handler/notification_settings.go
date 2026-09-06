package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// NotificationSettingsHandler は通知設定関連のHTTPハンドラ。
type NotificationSettingsHandler struct {
	getSettings    *usecase.GetNotificationSettingsUseCase
	updateSettings *usecase.UpdateNotificationSettingsUseCase
}

// NewNotificationSettingsHandler は新しいNotificationSettingsHandlerインスタンスを生成する。
func NewNotificationSettingsHandler(
	getSettings *usecase.GetNotificationSettingsUseCase,
	updateSettings *usecase.UpdateNotificationSettingsUseCase,
) *NotificationSettingsHandler {
	return &NotificationSettingsHandler{getSettings: getSettings, updateSettings: updateSettings}
}

// GetSettings は認証ユーザーの通知設定を取得する。
func (h *NotificationSettingsHandler) GetSettings(c *gin.Context) {
	settings, err := h.getSettings.Execute(c.Request.Context(), c.GetUint("userID"))
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, settings)
}

// updateNotificationSettingsRequest は通知設定更新のリクエストボディ。
type updateNotificationSettingsRequest struct {
	EnableLikes    bool `json:"enable_likes"`
	EnableComments bool `json:"enable_comments"`
	EnableFollows  bool `json:"enable_follows"`
	EnableMessages bool `json:"enable_messages"`
	EnableMentions bool `json:"enable_mentions"`
	EnableWebPush  bool `json:"enable_web_push"`
	EnableEmail    bool `json:"enable_email"`
	EnableSound    bool `json:"enable_sound"`
}

// UpdateSettings は認証ユーザーの通知設定を更新する。
func (h *NotificationSettingsHandler) UpdateSettings(c *gin.Context) {
	req := bindJSON[updateNotificationSettingsRequest](c)
	if req == nil {
		return
	}

	settings, err := h.updateSettings.Execute(c.Request.Context(), usecase.UpdateNotificationSettingsInput{
		UserID:         c.GetUint("userID"),
		EnableLikes:    req.EnableLikes,
		EnableComments: req.EnableComments,
		EnableFollows:  req.EnableFollows,
		EnableMessages: req.EnableMessages,
		EnableMentions: req.EnableMentions,
		EnableWebPush:  req.EnableWebPush,
		EnableEmail:    req.EnableEmail,
		EnableSound:    req.EnableSound,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, settings)
}

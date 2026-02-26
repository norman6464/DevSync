package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// WidgetSettingsServiceInterface はWidgetSettingsHandlerが依存するサービスのインターフェース。
type WidgetSettingsServiceInterface interface {
	GetSettings(userID uint) (*model.WidgetSettings, error)
	UpdateSettings(userID uint, settings string) error
}

// WidgetSettingsHandler はダッシュボードウィジェット設定のHTTPハンドラ。
type WidgetSettingsHandler struct {
	service WidgetSettingsServiceInterface
}

// NewWidgetSettingsHandler は新しいWidgetSettingsHandlerインスタンスを生成する。
func NewWidgetSettingsHandler(s WidgetSettingsServiceInterface) *WidgetSettingsHandler {
	return &WidgetSettingsHandler{service: s}
}

// GetSettings は認証ユーザーのウィジェット設定を取得する。
func (h *WidgetSettingsHandler) GetSettings(c *gin.Context) {
	userID := c.GetUint("userID")

	settings, err := h.service.GetSettings(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, settings)
}

// UpdateSettings は認証ユーザーのウィジェット設定を更新する。
func (h *WidgetSettingsHandler) UpdateSettings(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[dto.UpdateWidgetSettingsRequest](c)
	if input == nil {
		return
	}

	if err := h.service.UpdateSettings(userID, string(input.Settings)); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"message": "設定を更新しました"})
}

package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// WidgetSettingsHandler はダッシュボードウィジェット設定の HTTP ハンドラ。
type WidgetSettingsHandler struct {
	getSettings    *usecase.GetWidgetSettingsUseCase
	updateSettings *usecase.UpdateWidgetSettingsUseCase
}

// NewWidgetSettingsHandler は WidgetSettingsHandler を生成する。
func NewWidgetSettingsHandler(
	getSettings *usecase.GetWidgetSettingsUseCase,
	updateSettings *usecase.UpdateWidgetSettingsUseCase,
) *WidgetSettingsHandler {
	return &WidgetSettingsHandler{getSettings: getSettings, updateSettings: updateSettings}
}

// GetSettings は認証ユーザーのウィジェット設定を取得する。
func (h *WidgetSettingsHandler) GetSettings(c *gin.Context) {
	userID := c.GetUint("userID")

	settings, err := h.getSettings.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, settings)
}

// updateWidgetSettingsRequest はウィジェット設定更新リクエスト。
type updateWidgetSettingsRequest struct {
	Settings json.RawMessage `json:"settings" binding:"required"`
}

// UpdateSettings は認証ユーザーのウィジェット設定を更新する。
func (h *WidgetSettingsHandler) UpdateSettings(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[updateWidgetSettingsRequest](c)
	if input == nil {
		return
	}

	if err := h.updateSettings.Execute(c.Request.Context(), userID, string(input.Settings)); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("設定を更新しました"))
}

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// EmailPreferencesHandler はメール配信設定関連のHTTPハンドラ。
// ユーザーのメール配信設定の取得・更新を処理する。
type EmailPreferencesHandler struct {
	get    *usecase.GetEmailPreferencesUseCase
	update *usecase.UpdateEmailPreferencesUseCase
}

// NewEmailPreferencesHandler は新しいEmailPreferencesHandlerインスタンスを生成する。
func NewEmailPreferencesHandler(
	get *usecase.GetEmailPreferencesUseCase,
	update *usecase.UpdateEmailPreferencesUseCase,
) *EmailPreferencesHandler {
	return &EmailPreferencesHandler{get: get, update: update}
}

// GetPreferences はユーザーのメール配信設定を取得する。
func (h *EmailPreferencesHandler) GetPreferences(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := h.get.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.EmailPreferencesResponse{
		EmailWeeklyReport: user.EmailWeeklyReport,
		EmailLanguage:     user.EmailLanguage,
	})
}

// UpdatePreferences はユーザーのメール配信設定を更新する。
func (h *EmailPreferencesHandler) UpdatePreferences(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[dto.UpdateEmailPreferencesRequest](c)
	if input == nil {
		return
	}

	user, err := h.update.Execute(c.Request.Context(), userID, input.EmailWeeklyReport, input.EmailLanguage)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.EmailPreferencesResponse{
		EmailWeeklyReport: user.EmailWeeklyReport,
		EmailLanguage:     user.EmailLanguage,
	})
}

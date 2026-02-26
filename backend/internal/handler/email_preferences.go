package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// EmailPreferencesServiceInterface はメール配信設定に必要なサービスの抽象インターフェース。
type EmailPreferencesServiceInterface interface {
	GetByID(id uint) (*model.User, error)
	UpdateEmailPreferences(userID uint, weeklyReport *bool, language *string) (*model.User, error)
}

// EmailPreferencesHandler はメール配信設定関連のHTTPハンドラ。
// ユーザーのメール配信設定の取得・更新を処理する。
type EmailPreferencesHandler struct {
	userService EmailPreferencesServiceInterface
}

// NewEmailPreferencesHandler は新しいEmailPreferencesHandlerインスタンスを生成する。
func NewEmailPreferencesHandler(userService EmailPreferencesServiceInterface) *EmailPreferencesHandler {
	return &EmailPreferencesHandler{userService: userService}
}

// GetPreferences はユーザーのメール配信設定を取得する。
func (h *EmailPreferencesHandler) GetPreferences(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := h.userService.GetByID(userID)
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

	user, err := h.userService.UpdateEmailPreferences(userID, input.EmailWeeklyReport, input.EmailLanguage)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.EmailPreferencesResponse{
		EmailWeeklyReport: user.EmailWeeklyReport,
		EmailLanguage:     user.EmailLanguage,
	})
}

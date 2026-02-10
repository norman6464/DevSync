package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/service"
)

// EmailPreferencesHandler はメール配信設定関連のHTTPハンドラ。
// ユーザーのメール配信設定の取得・更新を処理する。
type EmailPreferencesHandler struct {
	userService *service.UserService
}

// NewEmailPreferencesHandler は新しいEmailPreferencesHandlerインスタンスを生成する。
func NewEmailPreferencesHandler(userService *service.UserService) *EmailPreferencesHandler {
	return &EmailPreferencesHandler{userService: userService}
}

// GetPreferences はユーザーのメール配信設定を取得する。
func (h *EmailPreferencesHandler) GetPreferences(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := h.userService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email_weekly_report": user.EmailWeeklyReport,
		"email_language":      user.EmailLanguage,
	})
}

// UpdatePreferences はユーザーのメール配信設定を更新する。
func (h *EmailPreferencesHandler) UpdatePreferences(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		EmailWeeklyReport *bool   `json:"email_weekly_report"`
		EmailLanguage     *string `json:"email_language"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if req.EmailWeeklyReport != nil {
		user.EmailWeeklyReport = *req.EmailWeeklyReport
	}
	if req.EmailLanguage != nil {
		// 言語バリデーション
		validLangs := map[string]bool{
			"ja": true, "en": true, "ko": true, "zh-CN": true, "zh-TW": true,
			"es": true, "fr": true, "de": true, "pt": true, "ru": true,
		}
		if !validLangs[*req.EmailLanguage] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email language"})
			return
		}
		user.EmailLanguage = *req.EmailLanguage
	}

	if err := h.userService.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update preferences"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email_weekly_report": user.EmailWeeklyReport,
		"email_language":      user.EmailLanguage,
	})
}

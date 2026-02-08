package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/service"
)

// BadgeHandler はバッジ関連のHTTPハンドラ。
// ユーザーバッジの取得・バッジ獲得通知の作成を処理する。
type BadgeHandler struct {
	service *service.BadgeService
}

// NewBadgeHandler は新しいBadgeHandlerインスタンスを生成する。
func NewBadgeHandler(s *service.BadgeService) *BadgeHandler {
	return &BadgeHandler{service: s}
}

// GetUserBadges は指定ユーザーの全バッジを獲得状況付きで返す。
func (h *BadgeHandler) GetUserBadges(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	badges, err := h.service.GetUserBadges(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"badges": badges})
}

// NotifyBadgeEarned は新しく獲得したバッジの通知を作成する。
func (h *BadgeHandler) NotifyBadgeEarned(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		BadgeID string `json:"badge_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "badge_id is required"})
		return
	}

	if err := h.service.NotifyBadgeEarned(userID, req.BadgeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "badge notification created"})
}

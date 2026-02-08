package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/service"
)

type BadgeHandler struct {
	service *service.BadgeService
}

func NewBadgeHandler(s *service.BadgeService) *BadgeHandler {
	return &BadgeHandler{service: s}
}

// GetUserBadges returns all badges with earned status for the given user.
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

// NotifyBadgeEarned creates a notification for a newly earned badge.
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

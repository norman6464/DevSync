package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
	"github.com/norman6464/devsync/backend/internal/service"
	"gorm.io/gorm"
)

type BadgeHandler struct {
	db               *gorm.DB
	notificationRepo *repository.NotificationRepository
}

func NewBadgeHandler(db *gorm.DB, notificationRepo *repository.NotificationRepository) *BadgeHandler {
	return &BadgeHandler{db: db, notificationRepo: notificationRepo}
}

// GetUserBadges returns all badges with earned status for the given user.
func (h *BadgeHandler) GetUserBadges(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	stats, err := service.GetBadgeStats(h.db, uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	badges := service.EvaluateBadges(stats)
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

	notification := &model.Notification{
		UserID:  userID,
		Type:    model.NotificationTypeBadge,
		ActorID: userID,
		BadgeID: &req.BadgeID,
	}
	if err := h.notificationRepo.Create(notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "badge notification created"})
}

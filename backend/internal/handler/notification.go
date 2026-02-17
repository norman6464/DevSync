package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// NotificationServiceInterface はNotificationServiceが実装すべきインターフェース。
type NotificationServiceInterface interface {
	GetByUserID(userID uint, page, limit int, notificationType string) ([]model.Notification, error)
	CountByUserID(userID uint, notificationType string) (int64, error)
	CountUnread(userID uint) (int64, error)
	MarkAsRead(id, userID uint) error
	MarkAllAsRead(userID uint) error
	Delete(id, userID uint) error
}

// NotificationHandler は通知関連のHTTPハンドラ。
// 通知の取得・既読処理・削除を処理する。
type NotificationHandler struct {
	service NotificationServiceInterface
}

// NewNotificationHandler は新しいNotificationHandlerインスタンスを生成する。
func NewNotificationHandler(s NotificationServiceInterface) *NotificationHandler {
	return &NotificationHandler{service: s}
}

// GetAll は認証ユーザーの通知一覧をページネーション付きで取得する。
func (h *NotificationHandler) GetAll(c *gin.Context) {
	userID := c.GetUint("userID")
	page, limit := parsePagination(c)
	notificationType := c.DefaultQuery("type", "")

	notifications, err := h.service.GetByUserID(userID, page, limit, notificationType)
	if err != nil {
		respondError(c, err)
		return
	}

	total, err := h.service.CountByUserID(userID, notificationType)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{
		"notifications": notifications,
		"total":         total,
	})
}

// GetUnreadCount は認証ユーザーの未読通知数を取得する。
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.service.CountUnread(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"count": count})
}

// MarkAsRead は指定された通知を既読にする。
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.MarkAsRead(id, userID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"message": "marked as read"})
}

// MarkAllAsRead は認証ユーザーの全通知を既読にする。
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.service.MarkAllAsRead(userID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"message": "all marked as read"})
}

// Delete は指定された通知を削除する。
func (h *NotificationHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(id, userID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"message": "deleted"})
}

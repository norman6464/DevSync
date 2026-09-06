package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// NotificationHandler は通知関連のHTTPハンドラ。
// 通知の取得・既読処理・削除を処理する。
type NotificationHandler struct {
	list        *usecase.ListNotificationsUseCase
	countUnread *usecase.CountUnreadNotificationsUseCase
	markAsRead  *usecase.MarkNotificationAsReadUseCase
	markAllRead *usecase.MarkAllNotificationsAsReadUseCase
	remove      *usecase.DeleteNotificationUseCase
}

// NewNotificationHandler は新しいNotificationHandlerインスタンスを生成する。
func NewNotificationHandler(
	list *usecase.ListNotificationsUseCase,
	countUnread *usecase.CountUnreadNotificationsUseCase,
	markAsRead *usecase.MarkNotificationAsReadUseCase,
	markAllRead *usecase.MarkAllNotificationsAsReadUseCase,
	remove *usecase.DeleteNotificationUseCase,
) *NotificationHandler {
	return &NotificationHandler{
		list:        list,
		countUnread: countUnread,
		markAsRead:  markAsRead,
		markAllRead: markAllRead,
		remove:      remove,
	}
}

// notificationListResponse は通知一覧レスポンス
type notificationListResponse struct {
	Notifications []model.Notification `json:"notifications"`
	Total         int64                `json:"total"`
	Page          int                  `json:"page"`
	Limit         int                  `json:"limit"`
}

// GetAll は認証ユーザーの通知一覧をページネーション付きで取得する。
func (h *NotificationHandler) GetAll(c *gin.Context) {
	userID := c.GetUint("userID")
	page, limit := parsePagination(c)
	notificationType := c.DefaultQuery("type", "")

	notifications, total, err := h.list.Execute(c.Request.Context(), userID, page, limit, notificationType)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, notificationListResponse{
		Notifications: notifications,
		Total:         total,
		Page:          page,
		Limit:         limit,
	})
}

// countResponse はカウントレスポンス
type countResponse struct {
	Count int64 `json:"count"`
}

// GetUnreadCount は認証ユーザーの未読通知数を取得する。
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.countUnread.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, countResponse{Count: count})
}

// MarkAsRead は指定された通知を既読にする。
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.markAsRead.Execute(c.Request.Context(), id, userID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("marked as read"))
}

// MarkAllAsRead は認証ユーザーの全通知を既読にする。
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.markAllRead.Execute(c.Request.Context(), userID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("all marked as read"))
}

// Delete は指定された通知を削除する。
func (h *NotificationHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.remove.Execute(c.Request.Context(), id, userID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("deleted"))
}

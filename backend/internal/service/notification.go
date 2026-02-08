package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// NotificationService handles notification business logic.
type NotificationService struct {
	repo repository.NotificationRepositoryInterface
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(repo repository.NotificationRepositoryInterface) *NotificationService {
	return &NotificationService{repo: repo}
}

// GetByUserID returns paginated notifications for a user.
func (s *NotificationService) GetByUserID(userID uint, page, limit int, notificationType string) ([]model.Notification, error) {
	return s.repo.FindByUserID(userID, page, limit, notificationType)
}

// CountByUserID returns the total count of notifications for a user.
func (s *NotificationService) CountByUserID(userID uint, notificationType string) (int64, error) {
	return s.repo.CountByUserID(userID, notificationType)
}

// CountUnread returns the count of unread notifications.
func (s *NotificationService) CountUnread(userID uint) (int64, error) {
	return s.repo.CountUnread(userID)
}

// MarkAsRead marks a notification as read.
func (s *NotificationService) MarkAsRead(id, userID uint) error {
	return s.repo.MarkAsRead(id, userID)
}

// MarkAllAsRead marks all notifications as read for a user.
func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return s.repo.MarkAllAsRead(userID)
}

// Delete removes a notification.
func (s *NotificationService) Delete(id, userID uint) error {
	return s.repo.Delete(id, userID)
}

// NotifyFollowers creates notifications for all followers of a user.
func (s *NotificationService) NotifyFollowers(actorID uint, postID uint, notificationType model.NotificationType) {
	followerIDs, err := s.repo.GetFollowerIDs(actorID)
	if err != nil || len(followerIDs) == 0 {
		return
	}
	var notifications []*model.Notification
	for _, followerID := range followerIDs {
		notifications = append(notifications, &model.Notification{
			UserID:  followerID,
			Type:    notificationType,
			ActorID: actorID,
			PostID:  &postID,
		})
	}
	s.repo.CreateBatch(notifications)
}

// CreateNotification creates a single notification.
func (s *NotificationService) CreateNotification(notification *model.Notification) error {
	return s.repo.Create(notification)
}

// CreateBatch creates multiple notifications at once.
func (s *NotificationService) CreateBatch(notifications []*model.Notification) error {
	return s.repo.CreateBatch(notifications)
}

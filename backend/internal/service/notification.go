package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// NotificationService は通知のビジネスロジックを提供する。
// 個別通知の作成・既読管理に加え、フォロワーへの一括通知を担当する。
type NotificationService struct {
	repo repository.NotificationRepositoryInterface
}

// NewNotificationService は新しいNotificationServiceインスタンスを生成する。
func NewNotificationService(repo repository.NotificationRepositoryInterface) *NotificationService {
	return &NotificationService{repo: repo}
}

// GetByUserID は指定ユーザーの通知をページネーション付きで取得する。
// notificationTypeで通知種別のフィルタリングが可能。
func (s *NotificationService) GetByUserID(userID uint, page, limit int, notificationType string) ([]model.Notification, error) {
	return s.repo.FindByUserID(userID, page, limit, notificationType)
}

// CountByUserID は指定ユーザーの通知総数を取得する。
func (s *NotificationService) CountByUserID(userID uint, notificationType string) (int64, error) {
	return s.repo.CountByUserID(userID, notificationType)
}

// CountUnread は未読通知の数を取得する。
func (s *NotificationService) CountUnread(userID uint) (int64, error) {
	return s.repo.CountUnread(userID)
}

// MarkAsRead は指定通知を既読にマークする。
func (s *NotificationService) MarkAsRead(id, userID uint) error {
	return s.repo.MarkAsRead(id, userID)
}

// MarkAllAsRead は指定ユーザーの全通知を既読にマークする。
func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return s.repo.MarkAllAsRead(userID)
}

// Delete は通知を削除する。
func (s *NotificationService) Delete(id, userID uint) error {
	return s.repo.Delete(id, userID)
}

// NotifyFollowers は指定ユーザーの全フォロワーに通知を一括作成する。
// 投稿作成時などに非同期で呼び出される。
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

// CreateNotification は単一の通知を作成する。
func (s *NotificationService) CreateNotification(notification *model.Notification) error {
	return s.repo.Create(notification)
}

// CreateBatch は複数の通知を一括作成する。
func (s *NotificationService) CreateBatch(notifications []*model.Notification) error {
	return s.repo.CreateBatch(notifications)
}

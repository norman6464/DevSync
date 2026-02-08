package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// MessageService handles message business logic.
type MessageService struct {
	repo                repository.MessageRepositoryInterface
	notificationService *NotificationService
}

// NewMessageService creates a new MessageService.
func NewMessageService(repo repository.MessageRepositoryInterface, notificationService *NotificationService) *MessageService {
	return &MessageService{repo: repo, notificationService: notificationService}
}

// GetConversations returns all conversations for a user.
func (s *MessageService) GetConversations(userID uint) ([]repository.ConversationSummary, error) {
	return s.repo.GetConversations(userID)
}

// GetConversation returns paginated messages between two users.
func (s *MessageService) GetConversation(userID, otherUserID uint, page, limit int) ([]model.Message, error) {
	// Mark messages as read
	s.repo.MarkAsRead(otherUserID, userID)

	return s.repo.GetConversation(userID, otherUserID, page, limit)
}

// SendMessage sends a message and notifies the receiver.
func (s *MessageService) SendMessage(msg *model.Message) error {
	if err := s.repo.Create(msg); err != nil {
		return err
	}

	// Create notification for message receiver
	go func(senderID, receiverID uint) {
		notification := &model.Notification{
			UserID:  receiverID,
			Type:    model.NotificationTypeMessage,
			ActorID: senderID,
		}
		s.notificationService.CreateNotification(notification)
	}(msg.SenderID, msg.ReceiverID)

	return nil
}

// MarkAsRead marks messages from a sender as read for a receiver.
func (s *MessageService) MarkAsRead(senderID, receiverID uint) error {
	return s.repo.MarkAsRead(senderID, receiverID)
}

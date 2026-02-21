package service

import (
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// MessageService はダイレクトメッセージのビジネスロジックを提供する。
// メッセージ送信時に受信者への通知も自動作成する。
type MessageService struct {
	repo                repository.MessageRepositoryInterface
	notificationService NotificationServiceInterface
}

// NewMessageService は新しいMessageServiceインスタンスを生成する。
func NewMessageService(repo repository.MessageRepositoryInterface, notificationService NotificationServiceInterface) *MessageService {
	return &MessageService{repo: repo, notificationService: notificationService}
}

// GetConversations は指定ユーザーの全会話一覧を取得する。
func (s *MessageService) GetConversations(userID uint) ([]model.ConversationSummary, error) {
	return s.repo.GetConversations(userID)
}

// GetConversation は2ユーザー間のメッセージをページネーション付きで取得する。
// 取得と同時に、相手からのメッセージを既読にマークする。
func (s *MessageService) GetConversation(userID, otherUserID uint, page, limit int) ([]model.Message, error) {
	// 相手からのメッセージを既読にマーク
	s.repo.MarkAsRead(otherUserID, userID)

	return s.repo.GetConversation(userID, otherUserID, page, limit)
}

// SendMessage はメッセージを送信し、受信者に非同期で通知を作成する。
// 自分自身へのメッセージ送信は許可しない。
func (s *MessageService) SendMessage(msg *model.Message) error {
	if err := domain.ValidateStringLength(msg.Content, 1, 5000, "メッセージ内容"); err != nil {
		return err
	}
	msg.Content = strings.TrimSpace(msg.Content)
	if msg.SenderID == msg.ReceiverID {
		return ErrBadRequest
	}
	if err := s.repo.Create(msg); err != nil {
		return err
	}

	// 受信者への通知を非同期で作成
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

// MarkAsRead は指定送信者から受信者へのメッセージを既読にマークする。
func (s *MessageService) MarkAsRead(senderID, receiverID uint) error {
	return s.repo.MarkAsRead(senderID, receiverID)
}

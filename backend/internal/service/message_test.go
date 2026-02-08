package service

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestMessageService() (*MessageService, *MockMessageRepository, *MockNotificationRepository) {
	msgRepo := new(MockMessageRepository)
	notifRepo := new(MockNotificationRepository)
	notifService := NewNotificationService(notifRepo)
	svc := NewMessageService(msgRepo, notifService)
	return svc, msgRepo, notifRepo
}

// ============================================================
// GetConversations
// ============================================================

func TestMessageGetConversations_Success(t *testing.T) {
	svc, msgRepo, _ := newTestMessageService()

	summaries := []repository.ConversationSummary{
		{UserID: 2, Name: "Alice", LastMessage: "Hello"},
	}
	msgRepo.On("GetConversations", uint(1)).Return(summaries, nil)

	result, err := svc.GetConversations(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Alice", result[0].Name)
	msgRepo.AssertExpectations(t)
}

// ============================================================
// GetConversation（既読マーク付き）
// ============================================================

func TestMessageGetConversation_Success(t *testing.T) {
	svc, msgRepo, _ := newTestMessageService()

	// MarkAsRead が先に呼ばれる
	msgRepo.On("MarkAsRead", uint(2), uint(1)).Return(nil)

	messages := []model.Message{{Content: "Hello"}, {Content: "Hi"}}
	msgRepo.On("GetConversation", uint(1), uint(2), 1, 20).Return(messages, nil)

	result, err := svc.GetConversation(1, 2, 1, 20)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	msgRepo.AssertCalled(t, "MarkAsRead", uint(2), uint(1))
	msgRepo.AssertExpectations(t)
}

// ============================================================
// SendMessage
// ============================================================

func TestMessageSendMessage_Success(t *testing.T) {
	svc, msgRepo, notifRepo := newTestMessageService()

	msg := &model.Message{SenderID: 1, ReceiverID: 2, Content: "Hello!"}
	msgRepo.On("Create", msg).Return(nil)

	// 通知はgoroutineで呼ばれる
	notifRepo.On("Create", mock.AnythingOfType("*model.Notification")).Return(nil).Maybe()

	err := svc.SendMessage(msg)
	assert.NoError(t, err)
	msgRepo.AssertExpectations(t)
}

// ============================================================
// MarkAsRead
// ============================================================

func TestMessageMarkAsRead_Success(t *testing.T) {
	svc, msgRepo, _ := newTestMessageService()

	msgRepo.On("MarkAsRead", uint(2), uint(1)).Return(nil)

	err := svc.MarkAsRead(2, 1)
	assert.NoError(t, err)
	msgRepo.AssertExpectations(t)
}

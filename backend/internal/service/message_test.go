package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newTestMessageService はMessageServiceのテスト用インスタンスを生成するヘルパー。
func newTestMessageService() (*MessageService, *MockMessageRepository, *MockNotificationRepository) {
	msgRepo := new(MockMessageRepository)
	notifRepo := new(MockNotificationRepository)
	notifService := NewNotificationService(notifRepo)
	svc := NewMessageService(msgRepo, notifService)
	return svc, msgRepo, notifRepo
}

// ============================================================
// 会話一覧取得テスト
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
// メッセージ送信テスト
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
// 既読マークテスト
// ============================================================

func TestMessageMarkAsRead_Success(t *testing.T) {
	svc, msgRepo, _ := newTestMessageService()

	msgRepo.On("MarkAsRead", uint(2), uint(1)).Return(nil)

	err := svc.MarkAsRead(2, 1)
	assert.NoError(t, err)
	msgRepo.AssertExpectations(t)
}

// ============================================================
// 会話一覧取得エラーテスト
// ============================================================

func TestMessageGetConversations_RepoError(t *testing.T) {
	svc, msgRepo, _ := newTestMessageService()

	msgRepo.On("GetConversations", uint(1)).Return([]repository.ConversationSummary(nil), errors.New("db error"))

	result, err := svc.GetConversations(1)
	assert.Error(t, err)
	assert.Nil(t, result)
	msgRepo.AssertExpectations(t)
}

func TestMessageGetConversations_Empty(t *testing.T) {
	svc, msgRepo, _ := newTestMessageService()

	msgRepo.On("GetConversations", uint(1)).Return([]repository.ConversationSummary{}, nil)

	result, err := svc.GetConversations(1)
	assert.NoError(t, err)
	assert.Empty(t, result)
	msgRepo.AssertExpectations(t)
}

// ============================================================
// GetConversation エラー・エッジケーステスト
// ============================================================

func TestMessageGetConversation_MarkAsReadFails_StillReturnsMessages(t *testing.T) {
	svc, msgRepo, _ := newTestMessageService()

	// MarkAsReadが失敗してもメッセージ取得は成功する
	msgRepo.On("MarkAsRead", uint(2), uint(1)).Return(errors.New("mark read failed"))

	messages := []model.Message{{Content: "Hello"}}
	msgRepo.On("GetConversation", uint(1), uint(2), 1, 20).Return(messages, nil)

	result, err := svc.GetConversation(1, 2, 1, 20)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	msgRepo.AssertExpectations(t)
}

func TestMessageGetConversation_RepoError(t *testing.T) {
	svc, msgRepo, _ := newTestMessageService()

	msgRepo.On("MarkAsRead", uint(2), uint(1)).Return(nil)
	msgRepo.On("GetConversation", uint(1), uint(2), 1, 20).Return([]model.Message(nil), errors.New("db error"))

	result, err := svc.GetConversation(1, 2, 1, 20)
	assert.Error(t, err)
	assert.Nil(t, result)
	msgRepo.AssertExpectations(t)
}

// ============================================================
// メッセージ送信エラーテスト
// ============================================================

func TestMessageSendMessage_CreateError(t *testing.T) {
	svc, msgRepo, notifRepo := newTestMessageService()

	msg := &model.Message{SenderID: 1, ReceiverID: 2, Content: "Hello!"}
	msgRepo.On("Create", msg).Return(errors.New("db error"))

	err := svc.SendMessage(msg)
	assert.Error(t, err)
	msgRepo.AssertExpectations(t)
	// Create失敗時は通知が送られないことを確認
	notifRepo.AssertNotCalled(t, "Create")
}

// ============================================================
// 既読マークエラーテスト
// ============================================================

func TestMessageMarkAsRead_RepoError(t *testing.T) {
	svc, msgRepo, _ := newTestMessageService()

	msgRepo.On("MarkAsRead", uint(2), uint(1)).Return(errors.New("db error"))

	err := svc.MarkAsRead(2, 1)
	assert.Error(t, err)
	msgRepo.AssertExpectations(t)
}

package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockMessagePort は usecase/repository.MessageRepository のモック。
type mockMessagePort struct{ mock.Mock }

func (m *mockMessagePort) Create(ctx context.Context, msg *model.Message) error {
	return m.Called(ctx, msg).Error(0)
}

func (m *mockMessagePort) GetConversation(ctx context.Context, userID, otherUserID uint, page, limit int) ([]model.Message, error) {
	args := m.Called(ctx, userID, otherUserID, page, limit)
	msgs, _ := args.Get(0).([]model.Message)
	return msgs, args.Error(1)
}

func (m *mockMessagePort) GetConversations(ctx context.Context, userID uint) ([]model.ConversationSummary, error) {
	args := m.Called(ctx, userID)
	c, _ := args.Get(0).([]model.ConversationSummary)
	return c, args.Error(1)
}

func (m *mockMessagePort) MarkAsRead(ctx context.Context, senderID, receiverID uint) error {
	return m.Called(ctx, senderID, receiverID).Error(0)
}

// recordingNotificationCreator は非同期に作られる通知を記録する。
type recordingNotificationCreator struct {
	mu            sync.Mutex
	notifications []*model.Notification
	err           error
	created       chan struct{}
}

func newRecordingNotificationCreator() *recordingNotificationCreator {
	return &recordingNotificationCreator{created: make(chan struct{}, 4)}
}

func (r *recordingNotificationCreator) Create(ctx context.Context, notification *model.Notification) error {
	r.mu.Lock()
	r.notifications = append(r.notifications, notification)
	r.mu.Unlock()
	r.created <- struct{}{}
	return r.err
}

// wait は通知が作られるまで待ち、作られなければ失敗させる。
func (r *recordingNotificationCreator) wait(t *testing.T) *model.Notification {
	t.Helper()
	select {
	case <-r.created:
		r.mu.Lock()
		defer r.mu.Unlock()
		return r.notifications[len(r.notifications)-1]
	case <-time.After(2 * time.Second):
		t.Fatal("通知が作成されなかった")
		return nil
	}
}

// newTestMessageHandler は本物の usecase に port モックを注入したハンドラーを生成する。
func newTestMessageHandler() (*MessageHandler, *mockMessagePort, *recordingNotificationCreator) {
	messages := new(mockMessagePort)
	notifications := newRecordingNotificationCreator()
	h := NewMessageHandler(
		usecase.NewListConversationsUseCase(messages),
		usecase.NewGetConversationUseCase(messages),
		usecase.NewSendMessageUseCase(messages, notifications),
		usecase.NewMarkMessagesAsReadUseCase(messages),
	)
	return h, messages, notifications
}

// ---------- GetConversations ----------

func TestMessageGetConversations_Success(t *testing.T) {
	h, messages, _ := newTestMessageHandler()
	r := newRouter(1)
	r.GET("/messages", h.GetConversations)

	messages.On("GetConversations", mock.Anything, uint(1)).
		Return([]model.ConversationSummary{{UserID: 2, Name: "相手", UnreadCount: 3}}, nil)

	w := doRequest(r, http.MethodGet, "/messages", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"unread_count":3`)
	messages.AssertExpectations(t)
}

// 会話が無ければ空配列を返す。
func TestMessageGetConversations_Empty(t *testing.T) {
	h, messages, _ := newTestMessageHandler()
	r := newRouter(1)
	r.GET("/messages", h.GetConversations)

	messages.On("GetConversations", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/messages", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

func TestMessageGetConversations_RepositoryError(t *testing.T) {
	h, messages, _ := newTestMessageHandler()
	r := newRouter(1)
	r.GET("/messages", h.GetConversations)

	messages.On("GetConversations", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/messages", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- GetMessages ----------

func TestMessageGetMessages_Success(t *testing.T) {
	h, messages, _ := newTestMessageHandler()
	r := newRouter(1)
	r.GET("/messages/:userId", h.GetMessages)

	// 取得の前に相手からのメッセージを既読にする。
	messages.On("MarkAsRead", mock.Anything, uint(2), uint(1)).Return(nil)
	messages.On("GetConversation", mock.Anything, uint(1), uint(2), 1, 20).
		Return([]model.Message{{Content: "やあ"}}, nil)

	w := doRequest(r, http.MethodGet, "/messages/2", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"content":"やあ"`)
	messages.AssertExpectations(t)
}

// 既読処理に失敗しても会話は返す（移行前と同じ）。
func TestMessageGetMessages_MarkAsReadFailureIsIgnored(t *testing.T) {
	h, messages, _ := newTestMessageHandler()
	r := newRouter(1)
	r.GET("/messages/:userId", h.GetMessages)

	messages.On("MarkAsRead", mock.Anything, uint(2), uint(1)).Return(errors.New("db error"))
	messages.On("GetConversation", mock.Anything, uint(1), uint(2), 1, 20).Return([]model.Message{}, nil)

	w := doRequest(r, http.MethodGet, "/messages/2", nil)
	assertStatus(t, w, http.StatusOK)
	messages.AssertExpectations(t)
}

func TestMessageGetMessages_WithPagination(t *testing.T) {
	h, messages, _ := newTestMessageHandler()
	r := newRouter(1)
	r.GET("/messages/:userId", h.GetMessages)

	messages.On("MarkAsRead", mock.Anything, uint(2), uint(1)).Return(nil)
	messages.On("GetConversation", mock.Anything, uint(1), uint(2), 3, 5).Return([]model.Message{}, nil)

	w := doRequest(r, http.MethodGet, "/messages/2?page=3&limit=5", nil)
	assertStatus(t, w, http.StatusOK)
	messages.AssertExpectations(t)
}

func TestMessageGetMessages_InvalidID(t *testing.T) {
	h, messages, _ := newTestMessageHandler()
	r := newRouter(1)
	r.GET("/messages/:userId", h.GetMessages)

	w := doRequest(r, http.MethodGet, "/messages/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
	messages.AssertNotCalled(t, "MarkAsRead", mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageGetMessages_RepositoryError(t *testing.T) {
	h, messages, _ := newTestMessageHandler()
	r := newRouter(1)
	r.GET("/messages/:userId", h.GetMessages)

	messages.On("MarkAsRead", mock.Anything, uint(2), uint(1)).Return(nil)
	messages.On("GetConversation", mock.Anything, uint(1), uint(2), 1, 20).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/messages/2", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- MarkAsRead ----------

func TestMessageMarkAsRead_Success(t *testing.T) {
	h, messages, _ := newTestMessageHandler()
	r := newRouter(1)
	r.PUT("/messages/:userId/read", h.MarkAsRead)

	messages.On("MarkAsRead", mock.Anything, uint(2), uint(1)).Return(nil)

	w := doRequest(r, http.MethodPut, "/messages/2/read", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "既読にしました")
	messages.AssertExpectations(t)
}

func TestMessageMarkAsRead_RepositoryError(t *testing.T) {
	h, messages, _ := newTestMessageHandler()
	r := newRouter(1)
	r.PUT("/messages/:userId/read", h.MarkAsRead)

	messages.On("MarkAsRead", mock.Anything, uint(2), uint(1)).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPut, "/messages/2/read", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- SendMessage ----------

func TestMessageSendMessage_Success(t *testing.T) {
	h, messages, notifications := newTestMessageHandler()
	r := newRouter(1)
	r.POST("/messages/:userId", h.SendMessage)

	messages.On("Create", mock.Anything, mock.MatchedBy(func(msg *model.Message) bool {
		// 前後の空白は落として保存する。
		return msg.SenderID == 1 && msg.ReceiverID == 2 && msg.Content == "こんにちは"
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/messages/2", map[string]string{"content": "  こんにちは  "})
	assertStatus(t, w, http.StatusCreated)
	assert.Contains(t, w.Body.String(), `"content":"こんにちは"`)

	// 受信者へメッセージ受信の通知が非同期で作られる。
	n := notifications.wait(t)
	assert.Equal(t, uint(2), n.UserID)
	assert.Equal(t, uint(1), n.ActorID)
	assert.Equal(t, model.NotificationTypeMessage, n.Type)
	messages.AssertExpectations(t)
}

// 自分自身へは送信できない。
func TestMessageSendMessage_ToSelf(t *testing.T) {
	h, messages, _ := newTestMessageHandler()
	r := newRouter(1)
	r.POST("/messages/:userId", h.SendMessage)

	w := doRequest(r, http.MethodPost, "/messages/1", map[string]string{"content": "自分宛"})
	assertStatus(t, w, http.StatusBadRequest)
	messages.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestMessageSendMessage_EmptyContent(t *testing.T) {
	h, messages, _ := newTestMessageHandler()
	r := newRouter(1)
	r.POST("/messages/:userId", h.SendMessage)

	w := doRequest(r, http.MethodPost, "/messages/2", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
	messages.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// 空白だけの本文も弾く。
func TestMessageSendMessage_BlankContent(t *testing.T) {
	h, messages, _ := newTestMessageHandler()
	r := newRouter(1)
	r.POST("/messages/:userId", h.SendMessage)

	w := doRequest(r, http.MethodPost, "/messages/2", map[string]string{"content": "   "})
	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "メッセージ内容を入力してください")
	messages.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// 5000 文字を超える本文は弾く。
func TestMessageSendMessage_TooLong(t *testing.T) {
	h, messages, _ := newTestMessageHandler()
	r := newRouter(1)
	r.POST("/messages/:userId", h.SendMessage)

	w := doRequest(r, http.MethodPost, "/messages/2", map[string]string{"content": strings.Repeat("あ", 5001)})
	assertStatus(t, w, http.StatusBadRequest)
	messages.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestMessageSendMessage_RepositoryError(t *testing.T) {
	h, messages, notifications := newTestMessageHandler()
	r := newRouter(1)
	r.POST("/messages/:userId", h.SendMessage)

	messages.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/messages/2", map[string]string{"content": "こんにちは"})
	assertStatus(t, w, http.StatusInternalServerError)
	// 保存に失敗したら通知も作らない。
	assert.Empty(t, notifications.created)
}

func TestMessageSendMessage_InvalidID(t *testing.T) {
	h, messages, _ := newTestMessageHandler()
	r := newRouter(1)
	r.POST("/messages/:userId", h.SendMessage)

	w := doRequest(r, http.MethodPost, "/messages/abc", map[string]string{"content": "こんにちは"})
	assertStatus(t, w, http.StatusBadRequest)
	messages.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

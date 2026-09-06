package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// stubAdviceContext はルールエンジンが参照するデータの取得をすべて成功させる。
// 生成されるアドバイスの内容は usecase 側のテストで検証する。
func stubAdviceContext(ports *aiPorts) {
	ports.Logs.On("GetStreakInfo", mock.Anything, mock.Anything).
		Return(&model.StreakInfo{CurrentStreak: 3, LongestStreak: 5, TotalDays: 10}, nil).Maybe()
	ports.Goals.On("GetByUserID", mock.Anything, mock.Anything, 100, 0).
		Return([]model.LearningGoal{}, int64(0), nil).Maybe()
	ports.Goals.On("GetStats", mock.Anything, mock.Anything).Return(&model.LearningGoalStats{}, nil).Maybe()
	ports.Roadmaps.On("GetByUserID", mock.Anything, mock.Anything, 100, 0).
		Return([]model.Roadmap{}, int64(0), nil).Maybe()
	ports.GitHub.On("GetLanguageStats", mock.Anything, mock.Anything).
		Return([]model.GitHubLanguageStat{}, nil).Maybe()
	ports.Logs.On("GetByUserID", mock.Anything, mock.Anything, 100, 0).
		Return([]model.LearningLog{}, int64(0), nil).Maybe()
	ports.Resources.On("FindByUserID", mock.Anything, mock.Anything, true, 100, 0).
		Return([]model.LearningResource{}, int64(0), nil).Maybe()
	ports.Users.On("FindByID", mock.Anything, mock.Anything).Return(&model.User{ID: 1}, nil).Maybe()
}

// ---------- GetAdvice ----------

func TestAIAdvice_GetAdvice_Success(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	stubAdviceContext(ports)
	ports.Conversations.On("CountTodayMessages", mock.Anything, uint(1)).Return(int64(2), nil)

	r := newRouter(1)
	r.GET("/advice", h.GetAdvice)
	w := doRequest(r, http.MethodGet, "/advice", nil)

	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"llm_available":true`)
	// 上限 5 回のうち 2 回使用しているので残り 3 回
	assert.Contains(t, w.Body.String(), `"daily_chat_remaining":3`)
	ports.Conversations.AssertExpectations(t)
}

// LLM 未設定のときは残り回数を数えない。
func TestAIAdvice_GetAdvice_LLMUnavailable(t *testing.T) {
	h, ports := setupAIAdviceHandlerWithoutLLM()
	stubAdviceContext(ports)

	r := newRouter(1)
	r.GET("/advice", h.GetAdvice)
	w := doRequest(r, http.MethodGet, "/advice", nil)

	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"llm_available":false`)
	assert.Contains(t, w.Body.String(), `"daily_chat_remaining":0`)
	ports.Conversations.AssertNotCalled(t, "CountTodayMessages", mock.Anything, mock.Anything)
}

// 残り回数の取得に失敗しても 200 で 0 を返す。
func TestAIAdvice_GetAdvice_RemainingError(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	stubAdviceContext(ports)
	ports.Conversations.On("CountTodayMessages", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/advice", h.GetAdvice)
	w := doRequest(r, http.MethodGet, "/advice", nil)

	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"daily_chat_remaining":0`)
}

// 学習状況の取得に失敗してもアドバイスは空配列で返す。
func TestAIAdvice_GetAdvice_ContextError(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	ports.Logs.On("GetStreakInfo", mock.Anything, uint(1)).Return(nil, errors.New("db error"))
	ports.Conversations.On("CountTodayMessages", mock.Anything, uint(1)).Return(int64(0), nil)

	r := newRouter(1)
	r.GET("/advice", h.GetAdvice)
	w := doRequest(r, http.MethodGet, "/advice", nil)

	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"advices":[]`)
}

// ---------- MarkAsRead ----------

func TestAIAdvice_MarkAsRead_Success(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	ports.Advices.On("MarkAsRead", mock.Anything, uint(10), uint(1)).Return(nil)

	r := newRouter(1)
	r.PUT("/advice/:id/read", h.MarkAsRead)
	w := doRequest(r, http.MethodPut, "/advice/10/read", nil)

	assertStatus(t, w, http.StatusOK)
	ports.Advices.AssertExpectations(t)
}

func TestAIAdvice_MarkAsRead_NotFound(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	ports.Advices.On("MarkAsRead", mock.Anything, uint(10), uint(1)).Return(errors.New("not found"))

	r := newRouter(1)
	r.PUT("/advice/:id/read", h.MarkAsRead)
	w := doRequest(r, http.MethodPut, "/advice/10/read", nil)

	assertStatus(t, w, http.StatusNotFound)
	ports.Advices.AssertExpectations(t)
}

func TestAIAdvice_MarkAsRead_InvalidID(t *testing.T) {
	h, ports := setupAIAdviceHandler()

	r := newRouter(1)
	r.PUT("/advice/:id/read", h.MarkAsRead)
	w := doRequest(r, http.MethodPut, "/advice/abc/read", nil)

	assertStatus(t, w, http.StatusBadRequest)
	ports.Advices.AssertNotCalled(t, "MarkAsRead", mock.Anything, mock.Anything, mock.Anything)
}

// ---------- Chat ----------

func TestAIAdvice_Chat_Success(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	stubAdviceContext(ports)
	ports.Conversations.On("CountTodayMessages", mock.Anything, uint(1)).Return(int64(0), nil)
	ports.Conversations.On("CreateConversation", mock.Anything, mock.MatchedBy(func(c *model.AIConversation) bool {
		return c.UserID == 1 && c.Title == "Goの学習方法を教えて"
	})).Return(nil).Run(func(args mock.Arguments) {
		args.Get(1).(*model.AIConversation).ID = 7
	})
	ports.Conversations.On("AddMessage", mock.Anything, mock.MatchedBy(func(m *model.AIMessage) bool {
		return m.Role == model.AIMessageRoleUser && m.Content == "Goの学習方法を教えて"
	})).Return(nil)
	ports.LLM.On("Complete", mock.Anything, mock.MatchedBy(func(msgs []model.ChatMessage) bool {
		// 先頭はシステムプロンプト、末尾がユーザーの発言
		return len(msgs) == 2 && msgs[0].Role == "system" && msgs[1].Content == "Goの学習方法を教えて"
	})).Return(&model.ChatResponse{Content: "まずは公式ツアーから", TokensUsed: 120}, nil)
	ports.Conversations.On("AddMessage", mock.Anything, mock.MatchedBy(func(m *model.AIMessage) bool {
		return m.Role == model.AIMessageRoleAssistant && m.TokensUsed == 120
	})).Return(nil)
	ports.Conversations.On("FindConversationByID", mock.Anything, uint(7)).
		Return(&model.AIConversation{ID: 7, UserID: 1}, nil)

	r := newRouter(1)
	r.POST("/chat", h.Chat)
	w := doRequest(r, http.MethodPost, "/chat", map[string]interface{}{"message": "Goの学習方法を教えて"})

	assertStatus(t, w, http.StatusOK)
	ports.LLM.AssertExpectations(t)
	ports.Conversations.AssertExpectations(t)
}

func TestAIAdvice_Chat_ValidationError(t *testing.T) {
	h, ports := setupAIAdviceHandler()

	r := newRouter(1)
	r.POST("/chat", h.Chat)
	w := doRequest(r, http.MethodPost, "/chat", map[string]interface{}{"message": ""})

	assertStatus(t, w, http.StatusBadRequest)
	ports.LLM.AssertNotCalled(t, "Complete", mock.Anything, mock.Anything)
}

// LLM 未設定なら 503 を返す。
func TestAIAdvice_Chat_LLMNotConfigured(t *testing.T) {
	h, ports := setupAIAdviceHandlerWithoutLLM()

	r := newRouter(1)
	r.POST("/chat", h.Chat)
	w := doRequest(r, http.MethodPost, "/chat", map[string]interface{}{"message": "こんにちは"})

	assertStatus(t, w, http.StatusServiceUnavailable)
	ports.Conversations.AssertNotCalled(t, "CountTodayMessages", mock.Anything, mock.Anything)
}

// 1 日の上限に達していたら 429 を返す。
func TestAIAdvice_Chat_RateLimited(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	ports.Conversations.On("CountTodayMessages", mock.Anything, uint(1)).Return(int64(5), nil)

	r := newRouter(1)
	r.POST("/chat", h.Chat)
	w := doRequest(r, http.MethodPost, "/chat", map[string]interface{}{"message": "こんにちは"})

	assertStatus(t, w, http.StatusTooManyRequests)
	ports.LLM.AssertNotCalled(t, "Complete", mock.Anything, mock.Anything)
}

// 他人の会話には投稿できない。
func TestAIAdvice_Chat_ForbiddenConversation(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	stubAdviceContext(ports)
	ports.Conversations.On("CountTodayMessages", mock.Anything, uint(1)).Return(int64(0), nil)
	ports.Conversations.On("FindConversationByID", mock.Anything, uint(9)).
		Return(&model.AIConversation{ID: 9, UserID: 999}, nil)

	r := newRouter(1)
	r.POST("/chat", h.Chat)
	w := doRequest(r, http.MethodPost, "/chat", map[string]interface{}{"message": "こんにちは", "conversation_id": 9})

	assertStatus(t, w, http.StatusForbidden)
	ports.LLM.AssertNotCalled(t, "Complete", mock.Anything, mock.Anything)
}

func TestAIAdvice_Chat_LLMError(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	stubAdviceContext(ports)
	ports.Conversations.On("CountTodayMessages", mock.Anything, uint(1)).Return(int64(0), nil)
	ports.Conversations.On("CreateConversation", mock.Anything, mock.Anything).Return(nil)
	ports.Conversations.On("AddMessage", mock.Anything, mock.Anything).Return(nil)
	ports.LLM.On("Complete", mock.Anything, mock.Anything).Return(nil, errors.New("api error"))

	r := newRouter(1)
	r.POST("/chat", h.Chat)
	w := doRequest(r, http.MethodPost, "/chat", map[string]interface{}{"message": "こんにちは"})

	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- 会話 ----------

func TestAIAdvice_GetConversations_Success(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	ports.Conversations.On("FindConversationsByUserID", mock.Anything, uint(1), 20, 0).
		Return([]model.AIConversation{{ID: 1, UserID: 1, Title: "会話"}}, nil)

	r := newRouter(1)
	r.GET("/conversations", h.GetConversations)
	w := doRequest(r, http.MethodGet, "/conversations", nil)

	assertStatus(t, w, http.StatusOK)
	ports.Conversations.AssertExpectations(t)
}

func TestAIAdvice_GetConversations_RepositoryError(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	ports.Conversations.On("FindConversationsByUserID", mock.Anything, uint(1), 20, 0).
		Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/conversations", h.GetConversations)
	w := doRequest(r, http.MethodGet, "/conversations", nil)

	assertStatus(t, w, http.StatusInternalServerError)
}

func TestAIAdvice_GetConversation_Success(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	ports.Conversations.On("FindConversationByID", mock.Anything, uint(3)).
		Return(&model.AIConversation{ID: 3, UserID: 1}, nil)

	r := newRouter(1)
	r.GET("/conversations/:id", h.GetConversation)
	w := doRequest(r, http.MethodGet, "/conversations/3", nil)

	assertStatus(t, w, http.StatusOK)
	ports.Conversations.AssertExpectations(t)
}

func TestAIAdvice_GetConversation_NotFound(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	ports.Conversations.On("FindConversationByID", mock.Anything, uint(3)).Return(nil, nil)

	r := newRouter(1)
	r.GET("/conversations/:id", h.GetConversation)
	w := doRequest(r, http.MethodGet, "/conversations/3", nil)

	assertStatus(t, w, http.StatusNotFound)
}

func TestAIAdvice_GetConversation_Forbidden(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	ports.Conversations.On("FindConversationByID", mock.Anything, uint(3)).
		Return(&model.AIConversation{ID: 3, UserID: 999}, nil)

	r := newRouter(1)
	r.GET("/conversations/:id", h.GetConversation)
	w := doRequest(r, http.MethodGet, "/conversations/3", nil)

	assertStatus(t, w, http.StatusForbidden)
}

func TestAIAdvice_DeleteConversation_Success(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	ports.Conversations.On("FindConversationByID", mock.Anything, uint(3)).
		Return(&model.AIConversation{ID: 3, UserID: 1}, nil)
	ports.Conversations.On("DeleteConversation", mock.Anything, uint(3), uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/conversations/:id", h.DeleteConversation)
	w := doRequest(r, http.MethodDelete, "/conversations/3", nil)

	assertStatus(t, w, http.StatusOK)
	ports.Conversations.AssertExpectations(t)
}

func TestAIAdvice_DeleteConversation_Forbidden(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	ports.Conversations.On("FindConversationByID", mock.Anything, uint(3)).
		Return(&model.AIConversation{ID: 3, UserID: 999}, nil)

	r := newRouter(1)
	r.DELETE("/conversations/:id", h.DeleteConversation)
	w := doRequest(r, http.MethodDelete, "/conversations/3", nil)

	assertStatus(t, w, http.StatusForbidden)
	ports.Conversations.AssertNotCalled(t, "DeleteConversation", mock.Anything, mock.Anything, mock.Anything)
}

func TestAIAdvice_DeleteConversation_InvalidID(t *testing.T) {
	h, ports := setupAIAdviceHandler()

	r := newRouter(1)
	r.DELETE("/conversations/:id", h.DeleteConversation)
	w := doRequest(r, http.MethodDelete, "/conversations/abc", nil)

	assertStatus(t, w, http.StatusBadRequest)
	ports.Conversations.AssertNotCalled(t, "FindConversationByID", mock.Anything, mock.Anything)
}

// ---------- 未読アドバイス ----------

func TestAIAdviceGetUnreadAdvice_Success(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	ports.Advices.On("FindUnreadByUserID", mock.Anything, uint(1)).
		Return([]model.AIAdvice{{ID: 1, TitleKey: "advice.streakBroken"}}, nil)

	r := newRouter(1)
	r.GET("/advice/unread", h.GetUnreadAdvice)
	w := doRequest(r, http.MethodGet, "/advice/unread", nil)

	assertStatus(t, w, http.StatusOK)
	ports.Advices.AssertExpectations(t)
}

func TestAIAdviceGetUnreadAdvice_Empty(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	ports.Advices.On("FindUnreadByUserID", mock.Anything, uint(1)).Return(nil, nil)

	r := newRouter(1)
	r.GET("/advice/unread", h.GetUnreadAdvice)
	w := doRequest(r, http.MethodGet, "/advice/unread", nil)

	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

func TestAIAdviceGetUnreadAdvice_RepositoryError(t *testing.T) {
	h, ports := setupAIAdviceHandler()
	ports.Advices.On("FindUnreadByUserID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/advice/unread", h.GetUnreadAdvice)
	w := doRequest(r, http.MethodGet, "/advice/unread", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	ports.Advices.AssertExpectations(t)
}

func TestAIChatRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request aiChatRequest
		wantErr bool
	}{
		{
			name: "有効なリクエスト（新規会話）",
			request: aiChatRequest{
				Message: "Go言語について教えてください",
			},
			wantErr: false,
		},
		{
			name: "有効なリクエスト（既存会話）",
			request: aiChatRequest{
				Message:        "続きを教えてください",
				ConversationID: 42,
			},
			wantErr: false,
		},
		{
			name: "メッセージが空",
			request: aiChatRequest{
				Message: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

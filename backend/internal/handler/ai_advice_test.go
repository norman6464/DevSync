package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestAIAdvice_GetAdvice_Success(t *testing.T) {
	h, svc := setupAIAdviceHandler()
	advices := []model.AIAdvice{{ID: 1, TitleKey: "test"}}
	svc.On("GenerateAdvice", uint(1)).Return(advices)
	svc.On("IsLLMAvailable").Return(true)
	svc.On("GetDailyChatRemaining", uint(1)).Return(5, nil)

	r := gin.New()
	r.GET("/advice", authMiddleware(1), h.GetAdvice)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/advice", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestAIAdvice_GetAdvice_LLMUnavailable(t *testing.T) {
	h, svc := setupAIAdviceHandler()
	svc.On("GenerateAdvice", uint(1)).Return([]model.AIAdvice{})
	svc.On("IsLLMAvailable").Return(false)

	r := gin.New()
	r.GET("/advice", authMiddleware(1), h.GetAdvice)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/advice", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestAIAdvice_MarkAsRead_Success(t *testing.T) {
	h, svc := setupAIAdviceHandler()
	svc.On("MarkAsRead", uint(10), uint(1)).Return(nil)

	r := gin.New()
	r.PUT("/advice/:id/read", authMiddleware(1), h.MarkAsRead)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/advice/10/read", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestAIAdvice_MarkAsRead_NotFound(t *testing.T) {
	h, svc := setupAIAdviceHandler()
	svc.On("MarkAsRead", uint(99), uint(1)).Return(errors.New("not found"))

	r := gin.New()
	r.PUT("/advice/:id/read", authMiddleware(1), h.MarkAsRead)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/advice/99/read", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

func TestAIAdvice_Chat_Success(t *testing.T) {
	h, svc := setupAIAdviceHandler()
	conv := &model.AIConversation{ID: 1, Title: "Test"}
	svc.On("Chat", uint(1), "hello", uint(0)).Return(conv, nil)

	r := gin.New()
	r.POST("/advice/chat", authMiddleware(1), h.Chat)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/advice/chat", jsonBody(map[string]interface{}{
		"message": "hello",
	}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestAIAdvice_Chat_ValidationError(t *testing.T) {
	h, _ := setupAIAdviceHandler()

	r := gin.New()
	r.POST("/advice/chat", authMiddleware(1), h.Chat)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/advice/chat", jsonBody(map[string]interface{}{}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestAIAdvice_DeleteConversation_Success(t *testing.T) {
	h, svc := setupAIAdviceHandler()
	svc.On("DeleteConversation", uint(5), uint(1)).Return(nil)

	r := gin.New()
	r.DELETE("/conversations/:id", authMiddleware(1), h.DeleteConversation)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/conversations/5", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestAIAdvice_GetConversations_Success(t *testing.T) {
	h, svc := setupAIAdviceHandler()
	convs := []model.AIConversation{{ID: 1}, {ID: 2}}
	svc.On("GetConversations", uint(1), 20, 0).Return(convs, nil)

	r := gin.New()
	r.GET("/conversations", authMiddleware(1), h.GetConversations)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/conversations", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestAIAdvice_GetConversation_Success(t *testing.T) {
	h, svc := setupAIAdviceHandler()
	conv := &model.AIConversation{ID: 3, Title: "Test Conv"}
	svc.On("GetConversation", uint(3), uint(1)).Return(conv, nil)

	r := gin.New()
	r.GET("/conversations/:id", authMiddleware(1), h.GetConversation)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/conversations/3", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestAIAdvice_GetConversation_ServiceError(t *testing.T) {
	h, svc := setupAIAdviceHandler()
	svc.On("GetConversation", uint(99), uint(1)).Return(nil, errors.New("not found"))

	r := gin.New()
	r.GET("/conversations/:id", authMiddleware(1), h.GetConversation)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/conversations/99", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetUnreadAdvice テスト
// ============================================================

func TestAIAdviceGetUnreadAdvice_Success(t *testing.T) {
	h, svc := setupAIAdviceHandler()
	r := newRouter(1)
	r.GET("/ai-advice/unread", h.GetUnreadAdvice)

	advices := []model.AIAdvice{
		{TitleKey: "advice.study_daily"},
	}
	svc.On("GetUnreadAdvice", uint(1)).Return(advices, nil)

	w := doRequest(r, http.MethodGet, "/ai-advice/unread", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestAIAdviceGetUnreadAdvice_Empty(t *testing.T) {
	h, svc := setupAIAdviceHandler()
	r := newRouter(1)
	r.GET("/ai-advice/unread", h.GetUnreadAdvice)

	svc.On("GetUnreadAdvice", uint(1)).Return([]model.AIAdvice(nil), nil)

	w := doRequest(r, http.MethodGet, "/ai-advice/unread", nil)
	assertStatus(t, w, http.StatusOK)
	// nilの場合は空配列[]に変換されるべき
	assert.Equal(t, "[]", w.Body.String())
	svc.AssertExpectations(t)
}

// ============================================================
// DeleteConversation テスト
// ============================================================

func TestAIAdviceDeleteConversation_Success(t *testing.T) {
	h, svc := setupAIAdviceHandler()
	r := newRouter(1)
	r.DELETE("/ai-advice/conversations/:id", h.DeleteConversation)

	svc.On("DeleteConversation", uint(5), uint(1)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/ai-advice/conversations/5", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestAIAdviceDeleteConversation_ServiceError(t *testing.T) {
	h, svc := setupAIAdviceHandler()
	r := newRouter(1)
	r.DELETE("/ai-advice/conversations/:id", h.DeleteConversation)

	svc.On("DeleteConversation", uint(99), uint(1)).Return(errors.New("not found"))

	w := doRequest(r, http.MethodDelete, "/ai-advice/conversations/99", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestAIAdviceDeleteConversation_InvalidID(t *testing.T) {
	h, _ := setupAIAdviceHandler()
	r := newRouter(1)
	r.DELETE("/ai-advice/conversations/:id", h.DeleteConversation)

	w := doRequest(r, http.MethodDelete, "/ai-advice/conversations/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestAIAdviceGetUnreadAdvice_ServiceError(t *testing.T) {
	h, svc := setupAIAdviceHandler()
	r := newRouter(1)
	r.GET("/ai-advice/unread", h.GetUnreadAdvice)

	svc.On("GetUnreadAdvice", uint(1)).Return([]model.AIAdvice(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/ai-advice/unread", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// カバレッジ向上テスト
// ============================================================

func TestAIAdvice_GetConversations_ServiceError(t *testing.T) {
	h, svc := setupAIAdviceHandler()
	r := newRouter(1)
	r.GET("/conversations", h.GetConversations)

	svc.On("GetConversations", uint(1), 20, 0).Return([]model.AIConversation(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/conversations", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestAIAdvice_Chat_ServiceError(t *testing.T) {
	h, svc := setupAIAdviceHandler()
	r := newRouter(1)
	r.POST("/advice/chat", h.Chat)

	svc.On("Chat", uint(1), "hello", uint(0)).Return(nil, errors.New("llm error"))

	w := doRequest(r, http.MethodPost, "/advice/chat", map[string]interface{}{
		"message": "hello",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestAIAdvice_GetAdvice_RemainingError(t *testing.T) {
	h, svc := setupAIAdviceHandler()
	r := newRouter(1)
	r.GET("/advice", h.GetAdvice)

	svc.On("GenerateAdvice", uint(1)).Return([]model.AIAdvice{})
	svc.On("IsLLMAvailable").Return(true)
	svc.On("GetDailyChatRemaining", uint(1)).Return(0, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/advice", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestAIAdvice_MarkAsRead_InvalidID(t *testing.T) {
	h, _ := setupAIAdviceHandler()
	r := newRouter(1)
	r.PUT("/advice/:id/read", h.MarkAsRead)

	w := doRequest(r, http.MethodPut, "/advice/abc/read", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestAIAdvice_GetConversation_InvalidID(t *testing.T) {
	h, _ := setupAIAdviceHandler()
	r := newRouter(1)
	r.GET("/conversations/:id", h.GetConversation)

	w := doRequest(r, http.MethodGet, "/conversations/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

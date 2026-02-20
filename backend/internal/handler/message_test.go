package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

func TestMessage_GetConversations_Success(t *testing.T) {
	h, svc := setupMessageHandler()
	svc.On("GetConversations", uint(1)).Return([]model.ConversationSummary{
		{UserID: 2, Name: "Alice", LastMessage: "Hello"},
		{UserID: 3, Name: "Bob", LastMessage: "Hi"},
	}, nil)

	r := newRouter(1)
	r.GET("/conversations", h.GetConversations)

	w := doRequest(r, http.MethodGet, "/conversations", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestMessage_GetConversations_ServiceError(t *testing.T) {
	h, svc := setupMessageHandler()
	svc.On("GetConversations", uint(1)).Return([]model.ConversationSummary(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/conversations", h.GetConversations)

	w := doRequest(r, http.MethodGet, "/conversations", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestMessage_GetMessages_Success(t *testing.T) {
	h, svc := setupMessageHandler()
	svc.On("GetConversation", uint(1), uint(5), 1, 20).Return([]model.Message{
		{Content: "Hello"},
	}, nil)

	r := newRouter(1)
	r.GET("/messages/:userId", h.GetMessages)

	w := doRequest(r, http.MethodGet, "/messages/5", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestMessage_GetMessages_WithPagination(t *testing.T) {
	h, svc := setupMessageHandler()
	svc.On("GetConversation", uint(1), uint(5), 2, 30).Return([]model.Message{}, nil)

	r := newRouter(1)
	r.GET("/messages/:userId", h.GetMessages)

	w := doRequest(r, http.MethodGet, "/messages/5?page=2&limit=30", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestMessage_SendMessage_Success(t *testing.T) {
	h, svc := setupMessageHandler()
	svc.On("SendMessage", mock.AnythingOfType("*model.Message")).Return(nil)

	r := newRouter(1)
	r.POST("/messages/:userId", h.SendMessage)

	w := doRequest(r, http.MethodPost, "/messages/5", map[string]string{
		"content": "Hello!",
	})
	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
}

func TestMessage_SendMessage_ValidationError(t *testing.T) {
	h, _ := setupMessageHandler()

	r := newRouter(1)
	r.POST("/messages/:userId", h.SendMessage)

	w := doRequest(r, http.MethodPost, "/messages/5", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestMessage_SendMessage_ServiceError(t *testing.T) {
	h, svc := setupMessageHandler()
	svc.On("SendMessage", mock.AnythingOfType("*model.Message")).Return(errors.New("send failed"))

	r := newRouter(1)
	r.POST("/messages/:userId", h.SendMessage)

	w := doRequest(r, http.MethodPost, "/messages/5", map[string]string{
		"content": "Hello!",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestMessage_GetMessages_InvalidID(t *testing.T) {
	h, _ := setupMessageHandler()

	r := newRouter(1)
	r.GET("/messages/:userId", h.GetMessages)

	w := doRequest(r, http.MethodGet, "/messages/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestMessage_GetMessages_ServiceError(t *testing.T) {
	h, svc := setupMessageHandler()
	svc.On("GetConversation", uint(1), uint(5), 1, 20).Return([]model.Message(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/messages/:userId", h.GetMessages)

	w := doRequest(r, http.MethodGet, "/messages/5", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

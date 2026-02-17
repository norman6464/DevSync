package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/mock"
)

func TestAtCoder_GetRating_Success(t *testing.T) {
	h, atcoderSvc, _ := setupAtCoderHandler()
	atcoderSvc.On("GetRating", "tourist").Return(&service.AtCoderRatingInfo{
		Rating: 3000,
	}, nil)

	r := newRouter(1)
	r.GET("/atcoder/:username/rating", h.GetRating)

	w := doRequest(r, http.MethodGet, "/atcoder/tourist/rating", nil)
	assertStatus(t, w, http.StatusOK)
	atcoderSvc.AssertExpectations(t)
}

func TestAtCoder_GetRating_Error(t *testing.T) {
	h, atcoderSvc, _ := setupAtCoderHandler()
	atcoderSvc.On("GetRating", "invalid").Return(nil, errors.New("not found"))

	r := newRouter(1)
	r.GET("/atcoder/:username/rating", h.GetRating)

	w := doRequest(r, http.MethodGet, "/atcoder/invalid/rating", nil)
	assertStatus(t, w, http.StatusBadRequest)
	atcoderSvc.AssertExpectations(t)
}

func TestAtCoder_Connect_Success(t *testing.T) {
	h, atcoderSvc, userSvc := setupAtCoderHandler()
	atcoderSvc.On("ValidateUsername", "myuser").Return(true)
	userSvc.On("GetByID", uint(1)).Return(&model.User{}, nil)
	userSvc.On("Update", mock.AnythingOfType("*model.User")).Return(nil)

	r := newRouter(1)
	r.POST("/atcoder/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/atcoder/connect", map[string]string{
		"username": "myuser",
	})
	assertStatus(t, w, http.StatusOK)
	atcoderSvc.AssertExpectations(t)
	userSvc.AssertExpectations(t)
}

func TestAtCoder_Connect_MissingUsername(t *testing.T) {
	h, _, _ := setupAtCoderHandler()

	r := newRouter(1)
	r.POST("/atcoder/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/atcoder/connect", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestAtCoder_Connect_InvalidUsername(t *testing.T) {
	h, atcoderSvc, _ := setupAtCoderHandler()
	atcoderSvc.On("ValidateUsername", "bad!user").Return(false)

	r := newRouter(1)
	r.POST("/atcoder/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/atcoder/connect", map[string]string{
		"username": "bad!user",
	})
	assertStatus(t, w, http.StatusBadRequest)
	atcoderSvc.AssertExpectations(t)
}

func TestAtCoder_Disconnect_Success(t *testing.T) {
	h, _, userSvc := setupAtCoderHandler()
	userSvc.On("GetByID", uint(1)).Return(&model.User{AtCoderUsername: "myuser"}, nil)
	userSvc.On("Update", mock.AnythingOfType("*model.User")).Return(nil)

	r := newRouter(1)
	r.DELETE("/atcoder/disconnect", h.Disconnect)

	w := doRequest(r, http.MethodDelete, "/atcoder/disconnect", nil)
	assertStatus(t, w, http.StatusOK)
	userSvc.AssertExpectations(t)
}

func TestAtCoder_Disconnect_UserNotFound(t *testing.T) {
	h, _, userSvc := setupAtCoderHandler()
	userSvc.On("GetByID", uint(1)).Return(nil, errors.New("not found"))

	r := newRouter(1)
	r.DELETE("/atcoder/disconnect", h.Disconnect)

	w := doRequest(r, http.MethodDelete, "/atcoder/disconnect", nil)
	assertStatus(t, w, http.StatusNotFound)
	userSvc.AssertExpectations(t)
}

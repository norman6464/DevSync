package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

func TestAtCoder_GetRating_Success(t *testing.T) {
	h, atcoderSvc := setupAtCoderHandler()
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
	h, atcoderSvc := setupAtCoderHandler()
	atcoderSvc.On("GetRating", "invalid").Return(nil, errors.New("not found"))

	r := newRouter(1)
	r.GET("/atcoder/:username/rating", h.GetRating)

	w := doRequest(r, http.MethodGet, "/atcoder/invalid/rating", nil)
	assertStatus(t, w, http.StatusBadRequest)
	atcoderSvc.AssertExpectations(t)
}

func TestAtCoder_Connect_Success(t *testing.T) {
	h, atcoderSvc := setupAtCoderHandler()
	atcoderSvc.On("ConnectAtCoder", uint(1), "myuser").Return(&model.User{AtCoderUsername: "myuser"}, nil)

	r := newRouter(1)
	r.POST("/atcoder/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/atcoder/connect", map[string]string{
		"username": "myuser",
	})
	assertStatus(t, w, http.StatusOK)
	atcoderSvc.AssertExpectations(t)
}

func TestAtCoder_Connect_MissingUsername(t *testing.T) {
	h, _ := setupAtCoderHandler()

	r := newRouter(1)
	r.POST("/atcoder/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/atcoder/connect", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestAtCoder_Connect_Error(t *testing.T) {
	h, atcoderSvc := setupAtCoderHandler()
	atcoderSvc.On("ConnectAtCoder", uint(1), "bad!user").Return(nil, errors.New("invalid username"))

	r := newRouter(1)
	r.POST("/atcoder/connect", h.Connect)

	w := doRequest(r, http.MethodPost, "/atcoder/connect", map[string]string{
		"username": "bad!user",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	atcoderSvc.AssertExpectations(t)
}

func TestAtCoder_Disconnect_Success(t *testing.T) {
	h, atcoderSvc := setupAtCoderHandler()
	atcoderSvc.On("DisconnectAtCoder", uint(1)).Return(&model.User{}, nil)

	r := newRouter(1)
	r.DELETE("/atcoder/disconnect", h.Disconnect)

	w := doRequest(r, http.MethodDelete, "/atcoder/disconnect", nil)
	assertStatus(t, w, http.StatusOK)
	atcoderSvc.AssertExpectations(t)
}

func TestAtCoder_Disconnect_Error(t *testing.T) {
	h, atcoderSvc := setupAtCoderHandler()
	atcoderSvc.On("DisconnectAtCoder", uint(1)).Return(nil, errors.New("not found"))

	r := newRouter(1)
	r.DELETE("/atcoder/disconnect", h.Disconnect)

	w := doRequest(r, http.MethodDelete, "/atcoder/disconnect", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	atcoderSvc.AssertExpectations(t)
}

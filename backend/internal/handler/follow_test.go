package handler

import (
	"errors"
	"fmt"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/assert"
)

// ---------- Follow ----------

func TestFollow_Success(t *testing.T) {
	h, svc := setupFollowHandler()
	svc.On("Follow", uint(1), uint(2)).Return(nil)

	r := newRouter(1)
	r.POST("/users/:id/follow", h.Follow)
	w := doRequest(r, "POST", "/users/2/follow", nil)
	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestFollow_InvalidID(t *testing.T) {
	h, _ := setupFollowHandler()

	r := newRouter(1)
	r.POST("/users/:id/follow", h.Follow)
	w := doRequest(r, "POST", "/users/abc/follow", nil)
	assertStatus(t, w, 400)
}

func TestFollow_ServiceError(t *testing.T) {
	h, svc := setupFollowHandler()
	svc.On("Follow", uint(1), uint(2)).Return(service.ErrBadRequest)

	r := newRouter(1)
	r.POST("/users/:id/follow", h.Follow)
	w := doRequest(r, "POST", "/users/2/follow", nil)
	assertStatus(t, w, 400)
}

// ---------- Unfollow ----------

func TestUnfollow_Success(t *testing.T) {
	h, svc := setupFollowHandler()
	svc.On("Unfollow", uint(1), uint(2)).Return(nil)

	r := newRouter(1)
	r.DELETE("/users/:id/follow", h.Unfollow)
	w := doRequest(r, "DELETE", "/users/2/follow", nil)
	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestUnfollow_ServiceError(t *testing.T) {
	h, svc := setupFollowHandler()
	svc.On("Unfollow", uint(1), uint(2)).Return(service.ErrNotFound)

	r := newRouter(1)
	r.DELETE("/users/:id/follow", h.Unfollow)
	w := doRequest(r, "DELETE", "/users/2/follow", nil)
	assertStatus(t, w, 404)
}

// ---------- GetFollowers ----------

func TestGetFollowers_Success(t *testing.T) {
	h, svc := setupFollowHandler()
	users := []model.User{{Name: "alice"}, {Name: "bob"}}
	svc.On("GetFollowers", uint(2)).Return(users, nil)

	r := newRouter(1)
	r.GET("/users/:id/followers", h.GetFollowers)
	w := doRequest(r, "GET", "/users/2/followers", nil)
	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestGetFollowers_Empty(t *testing.T) {
	h, svc := setupFollowHandler()
	var empty []model.User
	svc.On("GetFollowers", uint(2)).Return(empty, nil)

	r := newRouter(1)
	r.GET("/users/:id/followers", h.GetFollowers)
	w := doRequest(r, "GET", "/users/2/followers", nil)
	assertStatus(t, w, 200)
	// nilの場合は空配列を返すことを確認
	body := w.Body.String()
	assert.Contains(t, body, "[]")
}

func TestGetFollowers_ServiceError(t *testing.T) {
	h, svc := setupFollowHandler()
	var empty []model.User
	svc.On("GetFollowers", uint(2)).Return(empty, fmt.Errorf("internal error"))

	r := newRouter(1)
	r.GET("/users/:id/followers", h.GetFollowers)
	w := doRequest(r, "GET", "/users/2/followers", nil)
	assertStatus(t, w, 500)
}

// ---------- GetFollowing ----------

func TestGetFollowing_Success(t *testing.T) {
	h, svc := setupFollowHandler()
	users := []model.User{{Name: "charlie"}}
	svc.On("GetFollowing", uint(2)).Return(users, nil)

	r := newRouter(1)
	r.GET("/users/:id/following", h.GetFollowing)
	w := doRequest(r, "GET", "/users/2/following", nil)
	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestGetFollowing_Empty(t *testing.T) {
	h, svc := setupFollowHandler()
	var empty []model.User
	svc.On("GetFollowing", uint(2)).Return(empty, nil)

	r := newRouter(1)
	r.GET("/users/:id/following", h.GetFollowing)
	w := doRequest(r, "GET", "/users/2/following", nil)
	assertStatus(t, w, 200)
	body := w.Body.String()
	assert.Contains(t, body, "[]")
}

func TestGetFollowing_ServiceError(t *testing.T) {
	h, svc := setupFollowHandler()
	svc.On("GetFollowing", uint(2)).Return([]model.User(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:id/following", h.GetFollowing)
	w := doRequest(r, "GET", "/users/2/following", nil)
	assertStatus(t, w, 500)
	svc.AssertExpectations(t)
}

func TestGetFollowing_InvalidID(t *testing.T) {
	h, _ := setupFollowHandler()
	r := newRouter(1)
	r.GET("/users/:id/following", h.GetFollowing)
	w := doRequest(r, "GET", "/users/abc/following", nil)
	assertStatus(t, w, 400)
}

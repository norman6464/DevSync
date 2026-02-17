package handler

import (
	"fmt"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
)

// ---------- GetMyLevelInfo ----------

func TestLevelGetMyLevelInfo_Success(t *testing.T) {
	h, svc := setupLevelHandler()
	info := &model.LevelInfo{Level: 5, TotalXP: 1200}
	svc.On("GetLevelInfo", uint(1)).Return(info, nil)

	r := newRouter(1)
	r.GET("/levels/me", h.GetMyLevelInfo)
	w := doRequest(r, "GET", "/levels/me", nil)

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestLevelGetMyLevelInfo_ServiceError(t *testing.T) {
	h, svc := setupLevelHandler()
	svc.On("GetLevelInfo", uint(1)).Return(nil, fmt.Errorf("internal error"))

	r := newRouter(1)
	r.GET("/levels/me", h.GetMyLevelInfo)
	w := doRequest(r, "GET", "/levels/me", nil)

	assertStatus(t, w, 500)
	svc.AssertExpectations(t)
}

// ---------- GetLevelInfo ----------

func TestLevelGetLevelInfo_Success(t *testing.T) {
	h, svc := setupLevelHandler()
	info := &model.LevelInfo{Level: 3, TotalXP: 600}
	svc.On("GetLevelInfo", uint(5)).Return(info, nil)

	r := newRouter(1)
	r.GET("/users/:userId/level", h.GetLevelInfo)
	w := doRequest(r, "GET", "/users/5/level", nil)

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestLevelGetLevelInfo_InvalidID(t *testing.T) {
	h, _ := setupLevelHandler()

	r := newRouter(1)
	r.GET("/users/:userId/level", h.GetLevelInfo)
	w := doRequest(r, "GET", "/users/abc/level", nil)

	assertStatus(t, w, 400)
}

// ---------- GetXPBreakdown ----------

func TestLevelGetXPBreakdown_Success(t *testing.T) {
	h, svc := setupLevelHandler()
	breakdown := &model.XPBreakdown{Total: 1200}
	svc.On("GetXPBreakdown", uint(5)).Return(breakdown, nil)

	r := newRouter(1)
	r.GET("/users/:userId/xp", h.GetXPBreakdown)
	w := doRequest(r, "GET", "/users/5/xp", nil)

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestLevelGetXPBreakdown_InvalidID(t *testing.T) {
	h, _ := setupLevelHandler()

	r := newRouter(1)
	r.GET("/users/:userId/xp", h.GetXPBreakdown)
	w := doRequest(r, "GET", "/users/abc/xp", nil)

	assertStatus(t, w, 400)
}

func TestLevelGetXPBreakdown_ServiceError(t *testing.T) {
	h, svc := setupLevelHandler()
	svc.On("GetXPBreakdown", uint(5)).Return(nil, fmt.Errorf("internal error"))

	r := newRouter(1)
	r.GET("/users/:userId/xp", h.GetXPBreakdown)
	w := doRequest(r, "GET", "/users/5/xp", nil)

	assertStatus(t, w, 500)
	svc.AssertExpectations(t)
}

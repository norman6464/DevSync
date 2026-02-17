package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/mock"
)

func TestLearningLog_Create_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()
	svc.On("Create", mock.AnythingOfType("*model.LearningLog")).Return(nil)

	r := gin.New()
	r.POST("/logs", authMiddleware(1), h.Create)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logs", jsonBody(map[string]interface{}{
		"title":   "Go学習",
		"content": "インターフェースを学んだ",
	}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
}

func TestLearningLog_Create_ValidationError(t *testing.T) {
	h, _ := setupLearningLogHandler()

	r := gin.New()
	r.POST("/logs", authMiddleware(1), h.Create)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logs", jsonBody(map[string]interface{}{}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestLearningLog_GetByID_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()
	log := &model.LearningLog{Title: "Test", Content: "Content"}
	svc.On("GetByID", uint(1)).Return(log, nil)

	r := gin.New()
	r.GET("/logs/:id", h.GetByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/logs/1", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningLog_GetByID_ServiceError(t *testing.T) {
	h, svc := setupLearningLogHandler()
	svc.On("GetByID", uint(99)).Return(nil, errors.New("not found"))

	r := gin.New()
	r.GET("/logs/:id", h.GetByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/logs/99", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestLearningLog_GetMyLogs_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()
	logs := []model.LearningLog{{Title: "Log1"}, {Title: "Log2"}}
	svc.On("GetByUserID", uint(1)).Return(logs, nil)

	r := gin.New()
	r.GET("/me/logs", authMiddleware(1), h.GetMyLogs)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me/logs", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningLog_GetByUserID_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()
	logs := []model.LearningLog{{Title: "Log1"}}
	svc.On("GetByUserID", uint(5)).Return(logs, nil)

	r := gin.New()
	r.GET("/users/:userId/logs", h.GetByUserID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/5/logs", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningLog_Update_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()
	updated := &model.LearningLog{Title: "Updated"}
	svc.On("Update", uint(1), uint(3), mock.AnythingOfType("*model.LearningLog")).Return(updated, nil)

	r := gin.New()
	r.PUT("/logs/:id", authMiddleware(3), h.Update)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/logs/1", jsonBody(map[string]interface{}{
		"title": "Updated",
	}))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningLog_Delete_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()
	svc.On("Delete", uint(1), uint(3)).Return(nil)

	r := gin.New()
	r.DELETE("/logs/:id", authMiddleware(3), h.Delete)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/logs/1", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningLog_Delete_ServiceError(t *testing.T) {
	h, svc := setupLearningLogHandler()
	svc.On("Delete", uint(99), uint(3)).Return(errors.New("not found"))

	r := gin.New()
	r.DELETE("/logs/:id", authMiddleware(3), h.Delete)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/logs/99", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestLearningLog_GetStreakInfo_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()
	info := &model.StreakInfo{CurrentStreak: 5, LongestStreak: 10}
	svc.On("GetStreakInfo", uint(1)).Return(info, nil)

	r := gin.New()
	r.GET("/users/:userId/streak", h.GetStreakInfo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1/streak", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningLog_GetCalendarData_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()
	entries := []model.CalendarEntry{{Date: "2026-02-17", Count: 3}}
	svc.On("GetCalendarData", uint(1)).Return(entries, nil)

	r := gin.New()
	r.GET("/users/:userId/calendar", h.GetCalendarData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1/calendar", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningLog_GetCalendarData_ServiceError(t *testing.T) {
	h, svc := setupLearningLogHandler()
	svc.On("GetCalendarData", uint(99)).Return([]model.CalendarEntry(nil), errors.New("db error"))

	r := gin.New()
	r.GET("/users/:userId/calendar", h.GetCalendarData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/99/calendar", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

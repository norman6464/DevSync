package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/assert"
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
	svc.On("GetByID", uint(1), uint(1)).Return(log, nil)

	r := newRouter(1)
	r.GET("/logs/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/logs/1", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningLog_GetByID_ServiceError(t *testing.T) {
	h, svc := setupLearningLogHandler()
	svc.On("GetByID", uint(99), uint(1)).Return(nil, errors.New("not found"))

	r := newRouter(1)
	r.GET("/logs/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/logs/99", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestLearningLog_GetMyLogs_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()
	logs := []model.LearningLog{{Title: "Log1"}, {Title: "Log2"}}
	svc.On("GetByUserID", uint(1), 20, 0).Return(logs, int64(2), nil)

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
	svc.On("GetByUserID", uint(5), 20, 0).Return(logs, int64(1), nil)

	r := gin.New()
	r.GET("/users/:userId/logs", h.GetByUserID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/5/logs", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningLog_GetMyLogs_ServiceError(t *testing.T) {
	h, svc := setupLearningLogHandler()
	svc.On("GetByUserID", uint(1), 20, 0).Return([]model.LearningLog(nil), int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/me/logs", h.GetMyLogs)

	w := doRequest(r, http.MethodGet, "/me/logs", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestLearningLog_GetByUserID_ServiceError(t *testing.T) {
	h, svc := setupLearningLogHandler()
	svc.On("GetByUserID", uint(5), 20, 0).Return([]model.LearningLog(nil), int64(0), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:userId/logs", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/5/logs", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestLearningLog_GetByUserID_InvalidID(t *testing.T) {
	h, _ := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/users/:userId/logs", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/abc/logs", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestLearningLog_Update_InvalidJSON(t *testing.T) {
	h, _ := setupLearningLogHandler()
	r := newRouter(1)
	r.PUT("/logs/:id", h.Update)

	w := doRequestRaw(r, http.MethodPut, "/logs/1", "not json")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestLearningLog_Update_InvalidID(t *testing.T) {
	h, _ := setupLearningLogHandler()
	r := newRouter(1)
	r.PUT("/logs/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/logs/abc", map[string]interface{}{"title": "test"})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestLearningLog_Update_ServiceError(t *testing.T) {
	h, svc := setupLearningLogHandler()
	svc.On("Update", uint(1), uint(1), mock.AnythingOfType("*model.LearningLog")).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.PUT("/logs/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/logs/1", map[string]interface{}{"title": "Updated"})
	assertStatus(t, w, http.StatusInternalServerError)
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

func TestLearningLog_GetStreakInfo_ServiceError(t *testing.T) {
	h, svc := setupLearningLogHandler()
	svc.On("GetStreakInfo", uint(1)).Return(nil, errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:userId/streak", h.GetStreakInfo)
	w := doRequest(r, http.MethodGet, "/users/1/streak", nil)

	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestLearningLog_GetStreakInfo_InvalidID(t *testing.T) {
	h, _ := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/users/:userId/streak", h.GetStreakInfo)
	w := doRequest(r, http.MethodGet, "/users/abc/streak", nil)
	assertStatus(t, w, http.StatusBadRequest)
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

// ============================================================
// GetByCategory テスト
// ============================================================

func TestLearningLog_GetByCategory_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/category/:category", h.GetByCategory)

	logs := []model.LearningLog{{Category: "programming"}}
	svc.On("GetByCategory", uint(1), "programming").Return(logs, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/category/programming", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningLog_GetByCategory_NilResult(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/category/:category", h.GetByCategory)

	svc.On("GetByCategory", uint(1), "reading").Return([]model.LearningLog(nil), nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/category/reading", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
	svc.AssertExpectations(t)
}

func TestLearningLog_GetByCategory_ServiceError(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/category/:category", h.GetByCategory)

	svc.On("GetByCategory", uint(1), "invalid").Return([]model.LearningLog(nil), errors.New("bad category"))

	w := doRequest(r, http.MethodGet, "/learning-logs/category/invalid", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// GetBySource テスト
// ============================================================

func TestLearningLog_GetBySource_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/source/:source", h.GetBySource)

	logs := []model.LearningLog{{Source: "manual"}}
	svc.On("GetBySource", uint(1), "manual").Return(logs, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/source/manual", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningLog_GetBySource_ServiceError(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/source/:source", h.GetBySource)

	svc.On("GetBySource", uint(1), "invalid").Return([]model.LearningLog(nil), errors.New("bad source"))

	w := doRequest(r, http.MethodGet, "/learning-logs/source/invalid", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// ExportLogs テスト
// ============================================================

func TestLearningLog_ExportLogs_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/export", h.ExportLogs)

	csvData := []byte("date,category,minutes\n2026-01-01,programming,60")
	svc.On("ExportCSV", uint(1), 30).Return(csvData, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/export", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
	assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")
	svc.AssertExpectations(t)
}

func TestLearningLog_ExportLogs_AllPeriod(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/export", h.ExportLogs)

	csvData := []byte("date,category,minutes\n")
	svc.On("ExportCSV", uint(1), 0).Return(csvData, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/export?period=all", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Header().Get("Content-Disposition"), "learning-logs-all-")
	svc.AssertExpectations(t)
}

func TestLearningLog_ExportLogs_InvalidPeriod(t *testing.T) {
	h, _ := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/export", h.ExportLogs)

	w := doRequest(r, http.MethodGet, "/learning-logs/export?period=abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestLearningLog_ExportLogs_ServiceError(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/export", h.ExportLogs)

	svc.On("ExportCSV", uint(1), 30).Return(nil, errors.New("export error"))

	w := doRequest(r, http.MethodGet, "/learning-logs/export", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// ExportLogs JSON形式テスト
// ============================================================

func TestLearningLog_ExportLogs_JSONSuccess(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/export", h.ExportLogs)

	jsonData := []byte(`[{"date":"2026-02-19","title":"Go基礎","category":"coding","duration":60,"content":"変数を学んだ"}]`)
	svc.On("ExportJSON", uint(1), 30).Return(jsonData, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/export?format=json", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".json")
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	svc.AssertExpectations(t)
}

func TestLearningLog_ExportLogs_JSONAllPeriod(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/export", h.ExportLogs)

	jsonData := []byte(`[]`)
	svc.On("ExportJSON", uint(1), 0).Return(jsonData, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/export?format=json&period=all", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Header().Get("Content-Disposition"), "learning-logs-all-")
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".json")
	svc.AssertExpectations(t)
}

func TestLearningLog_ExportLogs_JSONServiceError(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/export", h.ExportLogs)

	svc.On("ExportJSON", uint(1), 30).Return(nil, errors.New("export error"))

	w := doRequest(r, http.MethodGet, "/learning-logs/export?format=json", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestLearningLog_ExportLogs_InvalidFormat(t *testing.T) {
	h, _ := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/export", h.ExportLogs)

	w := doRequest(r, http.MethodGet, "/learning-logs/export?format=xml", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// GetWeeklyDuration テスト
// ============================================================

func TestLearningLog_GetWeeklyDuration_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/weekly-duration/:userId", h.GetWeeklyDuration)

	svc.On("GetWeeklyDuration", uint(1)).Return(120, nil)

	w := doRequest(r, http.MethodGet, "/learning-logs/weekly-duration/1", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "120")
	svc.AssertExpectations(t)
}

func TestLearningLog_GetWeeklyDuration_ServiceError(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/weekly-duration/:userId", h.GetWeeklyDuration)

	svc.On("GetWeeklyDuration", uint(1)).Return(0, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/learning-logs/weekly-duration/1", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestLearningLog_GetWeeklyDuration_InvalidID(t *testing.T) {
	h, _ := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/learning-logs/weekly-duration/:userId", h.GetWeeklyDuration)

	w := doRequest(r, http.MethodGet, "/learning-logs/weekly-duration/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestLearningLog_Favorite_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.PUT("/logs/:id/favorite", h.Favorite)

	svc.On("FavoriteLog", uint(10), uint(1)).Return(nil)

	w := doRequest(r, http.MethodPut, "/logs/10/favorite", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningLog_Favorite_Forbidden(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.PUT("/logs/:id/favorite", h.Favorite)

	svc.On("FavoriteLog", uint(10), uint(1)).Return(service.ErrForbidden)

	w := doRequest(r, http.MethodPut, "/logs/10/favorite", nil)
	assertStatus(t, w, http.StatusForbidden)
	svc.AssertExpectations(t)
}

func TestLearningLog_Favorite_InvalidID(t *testing.T) {
	h, _ := setupLearningLogHandler()
	r := newRouter(1)
	r.PUT("/logs/:id/favorite", h.Favorite)

	w := doRequest(r, http.MethodPut, "/logs/abc/favorite", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestLearningLog_Unfavorite_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.PUT("/logs/:id/unfavorite", h.Unfavorite)

	svc.On("UnfavoriteLog", uint(10), uint(1)).Return(nil)

	w := doRequest(r, http.MethodPut, "/logs/10/unfavorite", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningLog_Unfavorite_InvalidID(t *testing.T) {
	h, _ := setupLearningLogHandler()
	r := newRouter(1)
	r.PUT("/logs/:id/unfavorite", h.Unfavorite)

	w := doRequest(r, http.MethodPut, "/logs/abc/unfavorite", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// BatchCreate テスト
// ============================================================

func TestLearningLog_BatchCreate_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.POST("/logs/batch", h.BatchCreate)

	resultLogs := []model.LearningLog{
		{Title: "Go基礎", Content: "変数を学んだ", UserID: 1, Duration: 60},
		{Title: "React入門", Content: "コンポーネント作成", UserID: 1, Duration: 45},
	}
	svc.On("BatchCreate", uint(1), mock.AnythingOfType("[]model.LearningLog")).Return(resultLogs, nil)

	w := doRequest(r, http.MethodPost, "/logs/batch", map[string]interface{}{
		"logs": []map[string]interface{}{
			{"title": "Go基礎", "content": "変数を学んだ", "duration": 60},
			{"title": "React入門", "content": "コンポーネント作成", "duration": 45},
		},
	})
	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
}

func TestLearningLog_BatchCreate_EmptyLogs(t *testing.T) {
	h, _ := setupLearningLogHandler()
	r := newRouter(1)
	r.POST("/logs/batch", h.BatchCreate)

	w := doRequest(r, http.MethodPost, "/logs/batch", map[string]interface{}{
		"logs": []map[string]interface{}{},
	})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestLearningLog_BatchCreate_InvalidJSON(t *testing.T) {
	h, _ := setupLearningLogHandler()
	r := newRouter(1)
	r.POST("/logs/batch", h.BatchCreate)

	w := doRequestRaw(r, http.MethodPost, "/logs/batch", "not json")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestLearningLog_BatchCreate_ServiceError(t *testing.T) {
	h, svc := setupLearningLogHandler()
	r := newRouter(1)
	r.POST("/logs/batch", h.BatchCreate)

	svc.On("BatchCreate", uint(1), mock.AnythingOfType("[]model.LearningLog")).Return(nil, errors.New("validation error"))

	w := doRequest(r, http.MethodPost, "/logs/batch", map[string]interface{}{
		"logs": []map[string]interface{}{
			{"title": "テスト", "content": "内容", "duration": 30},
		},
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestLearningLog_BatchCreate_MissingTitle(t *testing.T) {
	h, _ := setupLearningLogHandler()
	r := newRouter(1)
	r.POST("/logs/batch", h.BatchCreate)

	w := doRequest(r, http.MethodPost, "/logs/batch", map[string]interface{}{
		"logs": []map[string]interface{}{
			{"content": "内容のみ"},
		},
	})
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// GetLinkedLogs ハンドラーテスト
// ============================================================

func TestLearningLogGetLinkedLogs_Handler_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()
	logs := []model.LearningLog{{ID: 1, Title: "テスト"}}
	svc.On("GetLinkedLogs", uint(5), uint(1), 20, 0).Return(logs, int64(1), nil)
	r := newRouter(1)
	r.GET("/goals/:id/linked-logs", h.GetLinkedLogs)
	w := doRequest(r, http.MethodGet, "/goals/5/linked-logs", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningLogGetLinkedLogs_Handler_InvalidID(t *testing.T) {
	h, _ := setupLearningLogHandler()
	r := newRouter(1)
	r.GET("/goals/:id/linked-logs", h.GetLinkedLogs)
	w := doRequest(r, http.MethodGet, "/goals/abc/linked-logs", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestLearningLogGetLinkedLogs_Handler_ServiceError(t *testing.T) {
	h, svc := setupLearningLogHandler()
	svc.On("GetLinkedLogs", uint(5), uint(1), 20, 0).Return([]model.LearningLog(nil), int64(0), errors.New("forbidden"))
	r := newRouter(1)
	r.GET("/goals/:id/linked-logs", h.GetLinkedLogs)
	w := doRequest(r, http.MethodGet, "/goals/5/linked-logs", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

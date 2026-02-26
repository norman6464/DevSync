package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ============================================================
// クイックエントリー: 最近のカテゴリ取得ハンドラーテスト
// ============================================================

func TestGetRecentCategories_Handler_Success(t *testing.T) {
	h, svc := setupLearningLogHandler()

	categories := []string{"coding", "reading"}
	svc.On("GetRecentCategories", uint(1)).Return(categories, nil)

	r := gin.New()
	r.GET("/learning-logs/recent-categories", authMiddleware(1), h.GetRecentCategories)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/learning-logs/recent-categories", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestGetRecentCategories_Handler_Empty(t *testing.T) {
	h, svc := setupLearningLogHandler()

	svc.On("GetRecentCategories", uint(1)).Return([]string{}, nil)

	r := gin.New()
	r.GET("/learning-logs/recent-categories", authMiddleware(1), h.GetRecentCategories)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/learning-logs/recent-categories", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestGetRecentCategories_Handler_ServiceError(t *testing.T) {
	h, svc := setupLearningLogHandler()

	svc.On("GetRecentCategories", uint(1)).Return([]string(nil), errors.New("db error"))

	r := gin.New()
	r.GET("/learning-logs/recent-categories", authMiddleware(1), h.GetRecentCategories)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/learning-logs/recent-categories", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

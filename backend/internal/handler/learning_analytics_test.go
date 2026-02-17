package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

func TestLearningAnalytics_GetHeatmap_Success(t *testing.T) {
	h, svc := setupLearningAnalyticsHandler()
	data := []model.HeatmapEntry{{DayOfWeek: 1, Hour: 10, TotalMinutes: 60}}
	svc.On("GetHeatmap", uint(1)).Return(data, nil)

	r := gin.New()
	r.GET("/users/:userId/analytics/heatmap", h.GetHeatmap)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1/analytics/heatmap", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningAnalytics_GetHeatmap_InvalidID(t *testing.T) {
	h, _ := setupLearningAnalyticsHandler()

	r := gin.New()
	r.GET("/users/:userId/analytics/heatmap", h.GetHeatmap)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/abc/analytics/heatmap", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestLearningAnalytics_GetHeatmap_ServiceError(t *testing.T) {
	h, svc := setupLearningAnalyticsHandler()
	svc.On("GetHeatmap", uint(1)).Return([]model.HeatmapEntry(nil), errors.New("db error"))

	r := gin.New()
	r.GET("/users/:userId/analytics/heatmap", h.GetHeatmap)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1/analytics/heatmap", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestLearningAnalytics_GetCategoryBreakdown_Success(t *testing.T) {
	h, svc := setupLearningAnalyticsHandler()
	data := []model.CategoryBreakdown{{Category: "coding", TotalMinutes: 120}}
	svc.On("GetCategoryBreakdown", uint(2)).Return(data, nil)

	r := gin.New()
	r.GET("/users/:userId/analytics/categories", h.GetCategoryBreakdown)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/2/analytics/categories", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningAnalytics_GetProductivityScore_Success(t *testing.T) {
	h, svc := setupLearningAnalyticsHandler()
	score := &model.ProductivityScore{OverallScore: 85.5}
	svc.On("GetProductivityScore", uint(1)).Return(score, nil)

	r := gin.New()
	r.GET("/users/:userId/analytics/productivity", h.GetProductivityScore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1/analytics/productivity", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningAnalytics_GetWeeklyTrends_Success(t *testing.T) {
	h, svc := setupLearningAnalyticsHandler()
	data := []model.WeeklyTrend{{WeekStart: "2026-02-10", TotalMinutes: 300}}
	svc.On("GetWeeklyTrends", uint(1), 12).Return(data, nil)

	r := gin.New()
	r.GET("/users/:userId/analytics/trends", h.GetWeeklyTrends)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1/analytics/trends", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningAnalytics_GetWeeklyTrends_CustomWeeks(t *testing.T) {
	h, svc := setupLearningAnalyticsHandler()
	data := []model.WeeklyTrend{{WeekStart: "2026-02-10", TotalMinutes: 200}}
	svc.On("GetWeeklyTrends", uint(1), 4).Return(data, nil)

	r := gin.New()
	r.GET("/users/:userId/analytics/trends", h.GetWeeklyTrends)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1/analytics/trends?weeks=4", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningAnalytics_GetInsights_Success(t *testing.T) {
	h, svc := setupLearningAnalyticsHandler()
	insights := []model.AIInsight{{Type: "peak_time", Params: map[string]interface{}{"hour": float64(14)}}}
	svc.On("GetInsights", uint(3)).Return(insights, nil)

	r := gin.New()
	r.GET("/me/analytics/insights", authMiddleware(3), h.GetInsights)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me/analytics/insights", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestLearningAnalytics_GetInsights_ServiceError(t *testing.T) {
	h, svc := setupLearningAnalyticsHandler()
	svc.On("GetInsights", uint(3)).Return([]model.AIInsight(nil), errors.New("db error"))

	r := gin.New()
	r.GET("/me/analytics/insights", authMiddleware(3), h.GetInsights)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me/analytics/insights", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

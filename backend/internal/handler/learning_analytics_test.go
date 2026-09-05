package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockLearningAnalyticsRepo は usecase/repository.LearningAnalyticsRepository のモック（ctx 付き）。
type mockLearningAnalyticsRepo struct{ mock.Mock }

func (m *mockLearningAnalyticsRepo) GetHeatmapData(ctx context.Context, userID uint) ([]model.HeatmapEntry, error) {
	args := m.Called(ctx, userID)
	e, _ := args.Get(0).([]model.HeatmapEntry)
	return e, args.Error(1)
}

func (m *mockLearningAnalyticsRepo) GetCategoryBreakdown(ctx context.Context, userID uint) ([]model.CategoryBreakdown, error) {
	args := m.Called(ctx, userID)
	c, _ := args.Get(0).([]model.CategoryBreakdown)
	return c, args.Error(1)
}

func (m *mockLearningAnalyticsRepo) GetWeeklyTrends(ctx context.Context, userID uint, weeks int) ([]model.WeeklyTrend, error) {
	args := m.Called(ctx, userID, weeks)
	t, _ := args.Get(0).([]model.WeeklyTrend)
	return t, args.Error(1)
}

func (m *mockLearningAnalyticsRepo) GetProductivityStats(ctx context.Context, userID uint) (*model.ProductivityStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.ProductivityStats)
	return s, args.Error(1)
}

// newTestLearningAnalyticsHandler は本物の usecase に port モックを注入したハンドラーを生成する。
func newTestLearningAnalyticsHandler() (*LearningAnalyticsHandler, *mockLearningAnalyticsRepo) {
	repo := new(mockLearningAnalyticsRepo)
	h := NewLearningAnalyticsHandler(
		usecase.NewGetLearningHeatmapUseCase(repo),
		usecase.NewGetCategoryBreakdownUseCase(repo),
		usecase.NewGetWeeklyTrendsUseCase(repo),
		usecase.NewGetDayOfWeekSummaryUseCase(repo),
		usecase.NewGetProductivityScoreUseCase(repo),
		usecase.NewGetLearningInsightsUseCase(repo),
	)
	return h, repo
}

// ============================================================
// ヒートマップ
// ============================================================

func TestLearningAnalyticsHandler_GetHeatmap(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/users/:userId/analytics/heatmap", h.GetHeatmap)

	repo.On("GetHeatmapData", mock.Anything, uint(7)).
		Return([]model.HeatmapEntry{{DayOfWeek: 1, Hour: 9, TotalMinutes: 60}}, nil)

	w := doRequest(r, http.MethodGet, "/users/7/analytics/heatmap", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"total_minutes":60`)
	repo.AssertExpectations(t)
}

func TestLearningAnalyticsHandler_GetHeatmap_InvalidID(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/users/:userId/analytics/heatmap", h.GetHeatmap)

	w := doRequest(r, http.MethodGet, "/users/abc/analytics/heatmap", nil)
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "GetHeatmapData", mock.Anything, mock.Anything)
}

func TestLearningAnalyticsHandler_GetHeatmap_RepositoryError(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/users/:userId/analytics/heatmap", h.GetHeatmap)

	repo.On("GetHeatmapData", mock.Anything, uint(7)).
		Return([]model.HeatmapEntry(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/7/analytics/heatmap", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// カテゴリ別
// ============================================================

func TestLearningAnalyticsHandler_GetCategoryBreakdown(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/users/:userId/analytics/categories", h.GetCategoryBreakdown)

	repo.On("GetCategoryBreakdown", mock.Anything, uint(7)).Return([]model.CategoryBreakdown{
		{Category: "coding", TotalMinutes: 75, LogCount: 3},
		{Category: "reading", TotalMinutes: 25, LogCount: 1},
	}, nil)

	w := doRequest(r, http.MethodGet, "/users/7/analytics/categories", nil)
	assertStatus(t, w, http.StatusOK)
	// 割合は usecase で計算されてレスポンスに乗る。
	assert.Contains(t, w.Body.String(), `"percentage":75`)
	repo.AssertExpectations(t)
}

func TestLearningAnalyticsHandler_GetCategoryBreakdown_InvalidID(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/users/:userId/analytics/categories", h.GetCategoryBreakdown)

	w := doRequest(r, http.MethodGet, "/users/abc/analytics/categories", nil)
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "GetCategoryBreakdown", mock.Anything, mock.Anything)
}

func TestLearningAnalyticsHandler_GetCategoryBreakdown_RepositoryError(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/users/:userId/analytics/categories", h.GetCategoryBreakdown)

	repo.On("GetCategoryBreakdown", mock.Anything, uint(7)).
		Return([]model.CategoryBreakdown(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/7/analytics/categories", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// 生産性スコア
// ============================================================

func TestLearningAnalyticsHandler_GetProductivityScore(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/users/:userId/analytics/productivity", h.GetProductivityScore)

	repo.On("GetProductivityStats", mock.Anything, uint(7)).Return(&model.ProductivityStats{
		PomodoroSessions: 10, ManualSessions: 10,
		CompletedGoals: 5, TotalGoals: 10,
		TotalLogDays: 42, TotalDaysInRange: 84,
	}, nil)

	w := doRequest(r, http.MethodGet, "/users/7/analytics/productivity", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"overall_score":50`)
	repo.AssertExpectations(t)
}

func TestLearningAnalyticsHandler_GetProductivityScore_InvalidID(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/users/:userId/analytics/productivity", h.GetProductivityScore)

	w := doRequest(r, http.MethodGet, "/users/abc/analytics/productivity", nil)
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "GetProductivityStats", mock.Anything, mock.Anything)
}

func TestLearningAnalyticsHandler_GetProductivityScore_RepositoryError(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/users/:userId/analytics/productivity", h.GetProductivityScore)

	repo.On("GetProductivityStats", mock.Anything, uint(7)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/7/analytics/productivity", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// 週間トレンド
// ============================================================

func TestLearningAnalyticsHandler_GetWeeklyTrends(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/users/:userId/analytics/trends", h.GetWeeklyTrends)

	repo.On("GetWeeklyTrends", mock.Anything, uint(7), 12).
		Return([]model.WeeklyTrend{{WeekStart: "2026-01-05", TotalMinutes: 120, LogCount: 3}}, nil)

	w := doRequest(r, http.MethodGet, "/users/7/analytics/trends", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "2026-01-05")
	repo.AssertExpectations(t)
}

func TestLearningAnalyticsHandler_GetWeeklyTrends_CustomWeeks(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/users/:userId/analytics/trends", h.GetWeeklyTrends)

	repo.On("GetWeeklyTrends", mock.Anything, uint(7), 4).Return([]model.WeeklyTrend{}, nil)

	w := doRequest(r, http.MethodGet, "/users/7/analytics/trends?weeks=4", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// 不正・0・負の weeks はデフォルトの 12 週にフォールバックする。
func TestLearningAnalyticsHandler_GetWeeklyTrends_FallbackToDefault(t *testing.T) {
	for _, query := range []string{"?weeks=abc", "?weeks=0", "?weeks=-5"} {
		h, repo := newTestLearningAnalyticsHandler()
		r := newRouter(1)
		r.GET("/users/:userId/analytics/trends", h.GetWeeklyTrends)

		repo.On("GetWeeklyTrends", mock.Anything, uint(7), 12).Return([]model.WeeklyTrend{}, nil)

		w := doRequest(r, http.MethodGet, "/users/7/analytics/trends"+query, nil)
		assertStatus(t, w, http.StatusOK)
		repo.AssertExpectations(t)
	}
}

func TestLearningAnalyticsHandler_GetWeeklyTrends_InvalidID(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/users/:userId/analytics/trends", h.GetWeeklyTrends)

	w := doRequest(r, http.MethodGet, "/users/abc/analytics/trends", nil)
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "GetWeeklyTrends", mock.Anything, mock.Anything, mock.Anything)
}

func TestLearningAnalyticsHandler_GetWeeklyTrends_RepositoryError(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/users/:userId/analytics/trends", h.GetWeeklyTrends)

	repo.On("GetWeeklyTrends", mock.Anything, uint(7), 12).
		Return([]model.WeeklyTrend(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/7/analytics/trends", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// 曜日別サマリー
// ============================================================

func TestLearningAnalyticsHandler_GetDayOfWeekSummary(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/users/:userId/analytics/day-of-week", h.GetDayOfWeekSummary)

	repo.On("GetHeatmapData", mock.Anything, uint(7)).Return([]model.HeatmapEntry{
		{DayOfWeek: 1, Hour: 9, TotalMinutes: 60},
		{DayOfWeek: 1, Hour: 21, TotalMinutes: 30},
	}, nil)

	w := doRequest(r, http.MethodGet, "/users/7/analytics/day-of-week", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"average_minutes":45`)
	repo.AssertExpectations(t)
}

func TestLearningAnalyticsHandler_GetDayOfWeekSummary_InvalidID(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/users/:userId/analytics/day-of-week", h.GetDayOfWeekSummary)

	w := doRequest(r, http.MethodGet, "/users/abc/analytics/day-of-week", nil)
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "GetHeatmapData", mock.Anything, mock.Anything)
}

func TestLearningAnalyticsHandler_GetDayOfWeekSummary_RepositoryError(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/users/:userId/analytics/day-of-week", h.GetDayOfWeekSummary)

	repo.On("GetHeatmapData", mock.Anything, uint(7)).
		Return([]model.HeatmapEntry(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/7/analytics/day-of-week", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// AI インサイト
// ============================================================

func TestLearningAnalyticsHandler_GetInsights(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/analytics/insights", h.GetInsights)

	repo.On("GetHeatmapData", mock.Anything, uint(1)).
		Return([]model.HeatmapEntry{{DayOfWeek: 2, Hour: 21, TotalMinutes: 120}}, nil)
	repo.On("GetCategoryBreakdown", mock.Anything, uint(1)).
		Return([]model.CategoryBreakdown{{Category: "coding", TotalMinutes: 100}}, nil)
	repo.On("GetWeeklyTrends", mock.Anything, uint(1), 12).Return([]model.WeeklyTrend{}, nil)
	repo.On("GetProductivityStats", mock.Anything, uint(1)).
		Return(&model.ProductivityStats{CurrentStreak: 8}, nil)

	w := doRequest(r, http.MethodGet, "/analytics/insights", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "peak_time")
	assert.Contains(t, w.Body.String(), "streak_active")
	repo.AssertExpectations(t)
}

// 該当するインサイトが無くても null ではなく空配列を返す。
func TestLearningAnalyticsHandler_GetInsights_Empty(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/analytics/insights", h.GetInsights)

	repo.On("GetHeatmapData", mock.Anything, uint(1)).Return([]model.HeatmapEntry{}, nil)
	repo.On("GetCategoryBreakdown", mock.Anything, uint(1)).Return([]model.CategoryBreakdown{}, nil)
	repo.On("GetWeeklyTrends", mock.Anything, uint(1), 12).Return([]model.WeeklyTrend{}, nil)
	repo.On("GetProductivityStats", mock.Anything, uint(1)).Return(&model.ProductivityStats{}, nil)

	w := doRequest(r, http.MethodGet, "/analytics/insights", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

func TestLearningAnalyticsHandler_GetInsights_RepositoryError(t *testing.T) {
	h, repo := newTestLearningAnalyticsHandler()
	r := newRouter(1)
	r.GET("/analytics/insights", h.GetInsights)

	repo.On("GetHeatmapData", mock.Anything, uint(1)).
		Return([]model.HeatmapEntry(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/analytics/insights", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

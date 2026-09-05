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

// mockLearningLogSummaryReader は usecase/repository.LearningLogSummaryReader のモック（ctx 付き）。
type mockLearningLogSummaryReader struct{ mock.Mock }

func (m *mockLearningLogSummaryReader) GetStreakInfo(ctx context.Context, userID uint) (*model.StreakInfo, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.StreakInfo)
	return s, args.Error(1)
}

func (m *mockLearningLogSummaryReader) SumDurationByPeriod(ctx context.Context, userID uint, days int) (int, error) {
	args := m.Called(ctx, userID, days)
	return args.Int(0), args.Error(1)
}

// mockActiveLearningGoalReader は usecase/repository.ActiveLearningGoalReader のモック（ctx 付き）。
type mockActiveLearningGoalReader struct{ mock.Mock }

func (m *mockActiveLearningGoalReader) GetActiveByUserID(ctx context.Context, userID uint) ([]model.LearningGoal, error) {
	args := m.Called(ctx, userID)
	g, _ := args.Get(0).([]model.LearningGoal)
	return g, args.Error(1)
}

// mockProductivityStatsReader は usecase/repository.ProductivityStatsReader のモック（ctx 付き）。
type mockProductivityStatsReader struct{ mock.Mock }

func (m *mockProductivityStatsReader) GetProductivityStats(ctx context.Context, userID uint) (*model.ProductivityStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.ProductivityStats)
	return s, args.Error(1)
}

// learningDashboardPorts は LearningDashboardHandler が使う port モックの束。
type learningDashboardPorts struct {
	Logs      *mockLearningLogSummaryReader
	Goals     *mockActiveLearningGoalReader
	Analytics *mockProductivityStatsReader
}

// newTestLearningDashboardHandler は本物の usecase に port モックを注入したハンドラーを生成する。
func newTestLearningDashboardHandler() (*LearningDashboardHandler, learningDashboardPorts) {
	p := learningDashboardPorts{
		Logs:      new(mockLearningLogSummaryReader),
		Goals:     new(mockActiveLearningGoalReader),
		Analytics: new(mockProductivityStatsReader),
	}
	h := NewLearningDashboardHandler(
		usecase.NewGetLearningDashboardSummaryUseCase(p.Logs, p.Goals, p.Analytics),
	)
	return h, p
}

func TestLearningDashboardHandler_GetSummary(t *testing.T) {
	h, p := newTestLearningDashboardHandler()
	r := newRouter(1)
	r.GET("/dashboard/summary", h.GetSummary)

	p.Logs.On("GetStreakInfo", mock.Anything, uint(1)).
		Return(&model.StreakInfo{CurrentStreak: 3, LongestStreak: 7}, nil)
	p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 7).Return(300, nil)
	p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 1).Return(45, nil)
	p.Goals.On("GetActiveByUserID", mock.Anything, uint(1)).
		Return([]model.LearningGoal{{Title: "Go学習"}}, nil)
	p.Analytics.On("GetProductivityStats", mock.Anything, uint(1)).Return(&model.ProductivityStats{
		PomodoroSessions: 10, ManualSessions: 10,
		CompletedGoals: 5, TotalGoals: 10,
		TotalLogDays: 42, TotalDaysInRange: 84,
	}, nil)

	w := doRequest(r, http.MethodGet, "/dashboard/summary", nil)
	assertStatus(t, w, http.StatusOK)
	body := w.Body.String()
	assert.Contains(t, body, `"weekly_minutes":300`)
	assert.Contains(t, body, `"today_minutes":45`)
	assert.Contains(t, body, `"active_goal_count":1`)
	assert.Contains(t, body, `"overall_score":50`)
	p.Logs.AssertExpectations(t)
	p.Goals.AssertExpectations(t)
	p.Analytics.AssertExpectations(t)
}

func TestLearningDashboardHandler_GetSummary_Empty(t *testing.T) {
	h, p := newTestLearningDashboardHandler()
	r := newRouter(1)
	r.GET("/dashboard/summary", h.GetSummary)

	p.Logs.On("GetStreakInfo", mock.Anything, uint(1)).Return(&model.StreakInfo{}, nil)
	p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 7).Return(0, nil)
	p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 1).Return(0, nil)
	p.Goals.On("GetActiveByUserID", mock.Anything, uint(1)).Return([]model.LearningGoal{}, nil)
	p.Analytics.On("GetProductivityStats", mock.Anything, uint(1)).Return(&model.ProductivityStats{}, nil)

	w := doRequest(r, http.MethodGet, "/dashboard/summary", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"active_goal_count":0`)
}

func TestLearningDashboardHandler_GetSummary_RepositoryError(t *testing.T) {
	h, p := newTestLearningDashboardHandler()
	r := newRouter(1)
	r.GET("/dashboard/summary", h.GetSummary)

	p.Logs.On("GetStreakInfo", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/dashboard/summary", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	p.Goals.AssertNotCalled(t, "GetActiveByUserID", mock.Anything, mock.Anything)
}

// 集計の途中で失敗した場合も 500 になり、以降の取得は行わない。
func TestLearningDashboardHandler_GetSummary_StatsError(t *testing.T) {
	h, p := newTestLearningDashboardHandler()
	r := newRouter(1)
	r.GET("/dashboard/summary", h.GetSummary)

	p.Logs.On("GetStreakInfo", mock.Anything, uint(1)).Return(&model.StreakInfo{}, nil)
	p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 7).Return(100, nil)
	p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 1).Return(20, nil)
	p.Goals.On("GetActiveByUserID", mock.Anything, uint(1)).Return([]model.LearningGoal{}, nil)
	p.Analytics.On("GetProductivityStats", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/dashboard/summary", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

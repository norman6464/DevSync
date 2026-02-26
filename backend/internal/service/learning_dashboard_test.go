package service

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// newTestDashboardService はLearningDashboardServiceのテスト用インスタンスを生成するヘルパー。
func newTestDashboardService() (*LearningDashboardService, *MockLearningLogRepository, *MockLearningGoalRepository, *MockLearningAnalyticsRepository) {
	logRepo := new(MockLearningLogRepository)
	goalRepo := new(MockLearningGoalRepository)
	analyticsRepo := new(MockLearningAnalyticsRepository)
	svc := NewLearningDashboardService(logRepo, goalRepo, analyticsRepo)
	return svc, logRepo, goalRepo, analyticsRepo
}

// ============================================================
// GetSummary 正常系テスト
// ============================================================

func TestGetSummary_Success(t *testing.T) {
	svc, logRepo, goalRepo, analyticsRepo := newTestDashboardService()

	streakInfo := &model.StreakInfo{CurrentStreak: 5, LongestStreak: 10}
	logRepo.On("GetStreakInfo", uint(1)).Return(streakInfo, nil)
	logRepo.On("SumDurationByPeriod", uint(1), 7).Return(420, nil)
	logRepo.On("SumDurationByPeriod", uint(1), 1).Return(60, nil)
	goalRepo.On("GetActiveByUserID", uint(1)).Return([]model.LearningGoal{
		{Title: "Go学習"},
		{Title: "React学習"},
	}, nil)
	analyticsRepo.On("GetProductivityStats", uint(1)).Return(&model.ProductivityStats{
		PomodoroSessions: 10,
		ManualSessions:   5,
		CompletedGoals:   3,
		TotalGoals:       5,
		TotalLogDays:     20,
		TotalDaysInRange: 30,
	}, nil)

	summary, err := svc.GetSummary(1)
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, streakInfo, summary.StreakInfo)
	assert.Equal(t, 420, summary.WeeklyMinutes)
	assert.Equal(t, 60, summary.TodayMinutes)
	assert.Equal(t, 2, summary.ActiveGoalCount)
	assert.NotNil(t, summary.ProductivityScore)
	logRepo.AssertExpectations(t)
	goalRepo.AssertExpectations(t)
	analyticsRepo.AssertExpectations(t)
}

func TestGetSummary_EmptyData(t *testing.T) {
	svc, logRepo, goalRepo, analyticsRepo := newTestDashboardService()

	logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
	logRepo.On("SumDurationByPeriod", uint(1), 7).Return(0, nil)
	logRepo.On("SumDurationByPeriod", uint(1), 1).Return(0, nil)
	goalRepo.On("GetActiveByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
	analyticsRepo.On("GetProductivityStats", uint(1)).Return(&model.ProductivityStats{}, nil)

	summary, err := svc.GetSummary(1)
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, 0, summary.WeeklyMinutes)
	assert.Equal(t, 0, summary.TodayMinutes)
	assert.Equal(t, 0, summary.ActiveGoalCount)
	assert.Equal(t, 0.0, summary.ProductivityScore.OverallScore)
	logRepo.AssertExpectations(t)
}

// ============================================================
// GetSummary エラー系テスト
// ============================================================

func TestGetSummary_StreakInfoError(t *testing.T) {
	svc, logRepo, _, _ := newTestDashboardService()

	logRepo.On("GetStreakInfo", uint(1)).Return((*model.StreakInfo)(nil), assert.AnError)

	summary, err := svc.GetSummary(1)
	assert.Error(t, err)
	assert.Nil(t, summary)
}

func TestGetSummary_WeeklyMinutesError(t *testing.T) {
	svc, logRepo, _, _ := newTestDashboardService()

	logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
	logRepo.On("SumDurationByPeriod", uint(1), 7).Return(0, assert.AnError)

	summary, err := svc.GetSummary(1)
	assert.Error(t, err)
	assert.Nil(t, summary)
}

func TestGetSummary_TodayMinutesError(t *testing.T) {
	svc, logRepo, _, _ := newTestDashboardService()

	logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
	logRepo.On("SumDurationByPeriod", uint(1), 7).Return(100, nil)
	logRepo.On("SumDurationByPeriod", uint(1), 1).Return(0, assert.AnError)

	summary, err := svc.GetSummary(1)
	assert.Error(t, err)
	assert.Nil(t, summary)
}

func TestGetSummary_ActiveGoalsError(t *testing.T) {
	svc, logRepo, goalRepo, _ := newTestDashboardService()

	logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
	logRepo.On("SumDurationByPeriod", uint(1), 7).Return(100, nil)
	logRepo.On("SumDurationByPeriod", uint(1), 1).Return(30, nil)
	goalRepo.On("GetActiveByUserID", uint(1)).Return([]model.LearningGoal(nil), assert.AnError)

	summary, err := svc.GetSummary(1)
	assert.Error(t, err)
	assert.Nil(t, summary)
}

func TestGetSummary_ProductivityStatsError(t *testing.T) {
	svc, logRepo, goalRepo, analyticsRepo := newTestDashboardService()

	logRepo.On("GetStreakInfo", uint(1)).Return(&model.StreakInfo{}, nil)
	logRepo.On("SumDurationByPeriod", uint(1), 7).Return(100, nil)
	logRepo.On("SumDurationByPeriod", uint(1), 1).Return(30, nil)
	goalRepo.On("GetActiveByUserID", uint(1)).Return([]model.LearningGoal{}, nil)
	analyticsRepo.On("GetProductivityStats", uint(1)).Return((*model.ProductivityStats)(nil), assert.AnError)

	summary, err := svc.GetSummary(1)
	assert.Error(t, err)
	assert.Nil(t, summary)
}

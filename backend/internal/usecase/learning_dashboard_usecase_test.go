package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockLearningLogSummaryReader は usecase/repository.LearningLogSummaryReader のモック。
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

// mockActiveLearningGoalReader は usecase/repository.ActiveLearningGoalReader のモック。
type mockActiveLearningGoalReader struct{ mock.Mock }

func (m *mockActiveLearningGoalReader) GetActiveByUserID(ctx context.Context, userID uint) ([]model.LearningGoal, error) {
	args := m.Called(ctx, userID)
	g, _ := args.Get(0).([]model.LearningGoal)
	return g, args.Error(1)
}

// mockProductivityStatsReader は usecase/repository.ProductivityStatsReader のモック。
type mockProductivityStatsReader struct{ mock.Mock }

func (m *mockProductivityStatsReader) GetProductivityStats(ctx context.Context, userID uint) (*model.ProductivityStats, error) {
	args := m.Called(ctx, userID)
	s, _ := args.Get(0).(*model.ProductivityStats)
	return s, args.Error(1)
}

// dashboardPorts はダッシュボード usecase が使う port モックの束。
type dashboardPorts struct {
	Logs      *mockLearningLogSummaryReader
	Goals     *mockActiveLearningGoalReader
	Analytics *mockProductivityStatsReader
}

// newDashboardUseCase は port モックを注入した usecase を生成する。
func newDashboardUseCase() (*usecase.GetLearningDashboardSummaryUseCase, dashboardPorts) {
	p := dashboardPorts{
		Logs:      new(mockLearningLogSummaryReader),
		Goals:     new(mockActiveLearningGoalReader),
		Analytics: new(mockProductivityStatsReader),
	}
	return usecase.NewGetLearningDashboardSummaryUseCase(p.Logs, p.Goals, p.Analytics), p
}

func TestGetLearningDashboardSummaryUseCase_Execute(t *testing.T) {
	t.Run("各値をまとめて返す", func(t *testing.T) {
		uc, p := newDashboardUseCase()
		streak := &model.StreakInfo{CurrentStreak: 5, LongestStreak: 10}
		p.Logs.On("GetStreakInfo", mock.Anything, uint(1)).Return(streak, nil)
		p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 7).Return(420, nil)
		p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 1).Return(60, nil)
		p.Goals.On("GetActiveByUserID", mock.Anything, uint(1)).
			Return([]model.LearningGoal{{Title: "Go学習"}, {Title: "React学習"}}, nil)
		p.Analytics.On("GetProductivityStats", mock.Anything, uint(1)).Return(&model.ProductivityStats{
			PomodoroSessions: 10, ManualSessions: 10,
			CompletedGoals: 5, TotalGoals: 10,
			TotalLogDays: 42, TotalDaysInRange: 84,
		}, nil)

		summary, err := uc.Execute(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, streak, summary.StreakInfo)
		assert.Equal(t, 420, summary.WeeklyMinutes)
		assert.Equal(t, 60, summary.TodayMinutes)
		// 進行中の目標は件数だけを返す。
		assert.Equal(t, 2, summary.ActiveGoalCount)
		require.NotNil(t, summary.ProductivityScore)
		assert.Equal(t, 50.0, summary.ProductivityScore.OverallScore)
		p.Logs.AssertExpectations(t)
		p.Goals.AssertExpectations(t)
		p.Analytics.AssertExpectations(t)
	})

	t.Run("データが無くても 0 のサマリーを返す", func(t *testing.T) {
		uc, p := newDashboardUseCase()
		p.Logs.On("GetStreakInfo", mock.Anything, uint(1)).Return(&model.StreakInfo{}, nil)
		p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 7).Return(0, nil)
		p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 1).Return(0, nil)
		p.Goals.On("GetActiveByUserID", mock.Anything, uint(1)).Return([]model.LearningGoal{}, nil)
		p.Analytics.On("GetProductivityStats", mock.Anything, uint(1)).Return(&model.ProductivityStats{}, nil)

		summary, err := uc.Execute(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, 0, summary.WeeklyMinutes)
		assert.Equal(t, 0, summary.TodayMinutes)
		assert.Equal(t, 0, summary.ActiveGoalCount)
		assert.Equal(t, 0.0, summary.ProductivityScore.OverallScore)
	})

	t.Run("週間と当日で異なる日数を渡す", func(t *testing.T) {
		uc, p := newDashboardUseCase()
		p.Logs.On("GetStreakInfo", mock.Anything, uint(1)).Return(&model.StreakInfo{}, nil)
		p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 7).Return(100, nil)
		p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 1).Return(20, nil)
		p.Goals.On("GetActiveByUserID", mock.Anything, uint(1)).Return([]model.LearningGoal{}, nil)
		p.Analytics.On("GetProductivityStats", mock.Anything, uint(1)).Return(&model.ProductivityStats{}, nil)

		summary, err := uc.Execute(context.Background(), 1)

		require.NoError(t, err)
		assert.Equal(t, 100, summary.WeeklyMinutes)
		assert.Equal(t, 20, summary.TodayMinutes)
		p.Logs.AssertNumberOfCalls(t, "SumDurationByPeriod", 2)
	})

	// 途中で失敗したらその時点で中断し、以降の取得は行わない。
	t.Run("各段階のエラーで中断する", func(t *testing.T) {
		t.Run("ストリークの取得に失敗", func(t *testing.T) {
			uc, p := newDashboardUseCase()
			p.Logs.On("GetStreakInfo", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

			summary, err := uc.Execute(context.Background(), 1)

			assert.Nil(t, summary)
			assert.EqualError(t, err, "db error")
			p.Logs.AssertNotCalled(t, "SumDurationByPeriod", mock.Anything, mock.Anything, mock.Anything)
			p.Goals.AssertNotCalled(t, "GetActiveByUserID", mock.Anything, mock.Anything)
		})

		t.Run("週間学習時間の取得に失敗", func(t *testing.T) {
			uc, p := newDashboardUseCase()
			p.Logs.On("GetStreakInfo", mock.Anything, uint(1)).Return(&model.StreakInfo{}, nil)
			p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 7).Return(0, errors.New("db error"))

			summary, err := uc.Execute(context.Background(), 1)

			assert.Nil(t, summary)
			assert.EqualError(t, err, "db error")
			p.Goals.AssertNotCalled(t, "GetActiveByUserID", mock.Anything, mock.Anything)
		})

		t.Run("当日の学習時間の取得に失敗", func(t *testing.T) {
			uc, p := newDashboardUseCase()
			p.Logs.On("GetStreakInfo", mock.Anything, uint(1)).Return(&model.StreakInfo{}, nil)
			p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 7).Return(100, nil)
			p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 1).Return(0, errors.New("db error"))

			summary, err := uc.Execute(context.Background(), 1)

			assert.Nil(t, summary)
			assert.EqualError(t, err, "db error")
			p.Goals.AssertNotCalled(t, "GetActiveByUserID", mock.Anything, mock.Anything)
		})

		t.Run("進行中の目標の取得に失敗", func(t *testing.T) {
			uc, p := newDashboardUseCase()
			p.Logs.On("GetStreakInfo", mock.Anything, uint(1)).Return(&model.StreakInfo{}, nil)
			p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 7).Return(100, nil)
			p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 1).Return(20, nil)
			p.Goals.On("GetActiveByUserID", mock.Anything, uint(1)).
				Return([]model.LearningGoal(nil), errors.New("db error"))

			summary, err := uc.Execute(context.Background(), 1)

			assert.Nil(t, summary)
			assert.EqualError(t, err, "db error")
			p.Analytics.AssertNotCalled(t, "GetProductivityStats", mock.Anything, mock.Anything)
		})

		t.Run("生産性統計の取得に失敗", func(t *testing.T) {
			uc, p := newDashboardUseCase()
			p.Logs.On("GetStreakInfo", mock.Anything, uint(1)).Return(&model.StreakInfo{}, nil)
			p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 7).Return(100, nil)
			p.Logs.On("SumDurationByPeriod", mock.Anything, uint(1), 1).Return(20, nil)
			p.Goals.On("GetActiveByUserID", mock.Anything, uint(1)).Return([]model.LearningGoal{}, nil)
			p.Analytics.On("GetProductivityStats", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

			summary, err := uc.Execute(context.Background(), 1)

			assert.Nil(t, summary)
			assert.EqualError(t, err, "db error")
		})
	})
}

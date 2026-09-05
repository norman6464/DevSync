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

// mockLearningAnalyticsRepo は usecase/repository.LearningAnalyticsRepository のモック。
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

// ============================================================
// ヒートマップ / カテゴリ別
// ============================================================

func TestGetLearningHeatmapUseCase_Execute(t *testing.T) {
	t.Run("集計結果をそのまま返す", func(t *testing.T) {
		repo := new(mockLearningAnalyticsRepo)
		repo.On("GetHeatmapData", mock.Anything, uint(5)).
			Return([]model.HeatmapEntry{{DayOfWeek: 1, Hour: 9, TotalMinutes: 60}}, nil)
		uc := usecase.NewGetLearningHeatmapUseCase(repo)

		entries, err := uc.Execute(context.Background(), 5)

		require.NoError(t, err)
		assert.Len(t, entries, 1)
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockLearningAnalyticsRepo)
		repo.On("GetHeatmapData", mock.Anything, uint(5)).
			Return([]model.HeatmapEntry(nil), errors.New("db error"))
		uc := usecase.NewGetLearningHeatmapUseCase(repo)

		_, err := uc.Execute(context.Background(), 5)

		assert.EqualError(t, err, "db error")
	})
}

func TestGetCategoryBreakdownUseCase_Execute(t *testing.T) {
	t.Run("全体に占める割合を計算する", func(t *testing.T) {
		repo := new(mockLearningAnalyticsRepo)
		repo.On("GetCategoryBreakdown", mock.Anything, uint(5)).Return([]model.CategoryBreakdown{
			{Category: "coding", TotalMinutes: 90, LogCount: 3},
			{Category: "reading", TotalMinutes: 30, LogCount: 1},
		}, nil)
		uc := usecase.NewGetCategoryBreakdownUseCase(repo)

		categories, err := uc.Execute(context.Background(), 5)

		require.NoError(t, err)
		assert.Equal(t, 75.0, categories[0].Percentage)
		assert.Equal(t, 25.0, categories[1].Percentage)
	})

	t.Run("割り切れない割合は小数第 2 位まで丸める", func(t *testing.T) {
		repo := new(mockLearningAnalyticsRepo)
		repo.On("GetCategoryBreakdown", mock.Anything, uint(5)).Return([]model.CategoryBreakdown{
			{Category: "coding", TotalMinutes: 1},
			{Category: "reading", TotalMinutes: 2},
		}, nil)
		uc := usecase.NewGetCategoryBreakdownUseCase(repo)

		categories, err := uc.Execute(context.Background(), 5)

		require.NoError(t, err)
		assert.Equal(t, 33.33, categories[0].Percentage)
		assert.Equal(t, 66.67, categories[1].Percentage)
	})

	t.Run("合計が 0 なら割合も 0", func(t *testing.T) {
		repo := new(mockLearningAnalyticsRepo)
		repo.On("GetCategoryBreakdown", mock.Anything, uint(5)).Return([]model.CategoryBreakdown{
			{Category: "coding", TotalMinutes: 0, LogCount: 2},
		}, nil)
		uc := usecase.NewGetCategoryBreakdownUseCase(repo)

		categories, err := uc.Execute(context.Background(), 5)

		require.NoError(t, err)
		assert.Equal(t, 0.0, categories[0].Percentage)
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockLearningAnalyticsRepo)
		repo.On("GetCategoryBreakdown", mock.Anything, uint(5)).
			Return([]model.CategoryBreakdown(nil), errors.New("db error"))
		uc := usecase.NewGetCategoryBreakdownUseCase(repo)

		_, err := uc.Execute(context.Background(), 5)

		assert.EqualError(t, err, "db error")
	})
}

// ============================================================
// 週間トレンド / 曜日別
// ============================================================

func TestGetWeeklyTrendsUseCase_Execute(t *testing.T) {
	t.Run("指定した週数で取得する", func(t *testing.T) {
		repo := new(mockLearningAnalyticsRepo)
		repo.On("GetWeeklyTrends", mock.Anything, uint(5), 4).
			Return([]model.WeeklyTrend{{WeekStart: "2026-01-05", TotalMinutes: 120}}, nil)
		uc := usecase.NewGetWeeklyTrendsUseCase(repo)

		trends, err := uc.Execute(context.Background(), 5, 4)

		require.NoError(t, err)
		assert.Len(t, trends, 1)
		repo.AssertExpectations(t)
	})

	t.Run("0 以下の週数はデフォルトの 12 週になる", func(t *testing.T) {
		for _, weeks := range []int{0, -1} {
			repo := new(mockLearningAnalyticsRepo)
			repo.On("GetWeeklyTrends", mock.Anything, uint(5), 12).Return([]model.WeeklyTrend{}, nil)
			uc := usecase.NewGetWeeklyTrendsUseCase(repo)

			_, err := uc.Execute(context.Background(), 5, weeks)

			require.NoError(t, err)
			repo.AssertExpectations(t)
		}
	})
}

func TestGetDayOfWeekSummaryUseCase_Execute(t *testing.T) {
	t.Run("ヒートマップを曜日別に集計する", func(t *testing.T) {
		repo := new(mockLearningAnalyticsRepo)
		repo.On("GetHeatmapData", mock.Anything, uint(5)).Return([]model.HeatmapEntry{
			{DayOfWeek: 1, Hour: 9, TotalMinutes: 60},
			{DayOfWeek: 1, Hour: 20, TotalMinutes: 30},
			{DayOfWeek: 3, Hour: 10, TotalMinutes: 45},
		}, nil)
		uc := usecase.NewGetDayOfWeekSummaryUseCase(repo)

		summary, err := uc.Execute(context.Background(), 5)

		require.NoError(t, err)
		require.Len(t, summary, 7)
		assert.Equal(t, 90, summary[1].TotalMinutes)
		assert.Equal(t, 2, summary[1].LogCount)
		assert.Equal(t, 45, summary[1].AverageMinutes)
		assert.Equal(t, 45, summary[3].TotalMinutes)
		assert.Equal(t, 0, summary[0].LogCount)
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockLearningAnalyticsRepo)
		repo.On("GetHeatmapData", mock.Anything, uint(5)).
			Return([]model.HeatmapEntry(nil), errors.New("db error"))
		uc := usecase.NewGetDayOfWeekSummaryUseCase(repo)

		_, err := uc.Execute(context.Background(), 5)

		assert.EqualError(t, err, "db error")
	})
}

func TestAggregateDayOfWeek(t *testing.T) {
	t.Run("ログが無くても 7 件返す", func(t *testing.T) {
		summary := usecase.AggregateDayOfWeek(nil)

		require.Len(t, summary, 7)
		for i, s := range summary {
			assert.Equal(t, i, s.DayOfWeek)
			assert.Zero(t, s.TotalMinutes)
			assert.Zero(t, s.LogCount)
			assert.Zero(t, s.AverageMinutes)
		}
	})

	t.Run("平均は整数除算で切り捨てる", func(t *testing.T) {
		summary := usecase.AggregateDayOfWeek([]model.HeatmapEntry{
			{DayOfWeek: 2, TotalMinutes: 10},
			{DayOfWeek: 2, TotalMinutes: 11},
		})

		assert.Equal(t, 21, summary[2].TotalMinutes)
		assert.Equal(t, 10, summary[2].AverageMinutes)
	})

	t.Run("範囲外の曜日は無視する", func(t *testing.T) {
		summary := usecase.AggregateDayOfWeek([]model.HeatmapEntry{
			{DayOfWeek: 7, TotalMinutes: 100},
			{DayOfWeek: -1, TotalMinutes: 100},
			{DayOfWeek: 0, TotalMinutes: 5},
		})

		require.Len(t, summary, 7)
		assert.Equal(t, 5, summary[0].TotalMinutes)
		total := 0
		for _, s := range summary {
			total += s.TotalMinutes
		}
		assert.Equal(t, 5, total)
	})
}

// ============================================================
// 生産性スコア
// ============================================================

func TestGetProductivityScoreUseCase_Execute(t *testing.T) {
	t.Run("統計からスコアを算出する", func(t *testing.T) {
		repo := new(mockLearningAnalyticsRepo)
		repo.On("GetProductivityStats", mock.Anything, uint(5)).Return(&model.ProductivityStats{
			PomodoroSessions: 10, ManualSessions: 10,
			CompletedGoals: 5, TotalGoals: 10,
			TotalLogDays: 42, TotalDaysInRange: 84,
		}, nil)
		uc := usecase.NewGetProductivityScoreUseCase(repo)

		score, err := uc.Execute(context.Background(), 5)

		require.NoError(t, err)
		assert.Equal(t, 50.0, score.PomodoroRate)
		assert.Equal(t, 50.0, score.GoalRate)
		assert.Equal(t, 50.0, score.StreakConsistency)
		assert.Equal(t, 50.0, score.OverallScore)
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockLearningAnalyticsRepo)
		repo.On("GetProductivityStats", mock.Anything, uint(5)).Return(nil, errors.New("db error"))
		uc := usecase.NewGetProductivityScoreUseCase(repo)

		_, err := uc.Execute(context.Background(), 5)

		assert.EqualError(t, err, "db error")
	})
}

func TestCalculateProductivityScore(t *testing.T) {
	t.Run("すべて 0 ならスコアも 0", func(t *testing.T) {
		score := usecase.CalculateProductivityScore(&model.ProductivityStats{})

		assert.Equal(t, 0.0, score.PomodoroRate)
		assert.Equal(t, 0.0, score.GoalRate)
		assert.Equal(t, 0.0, score.StreakConsistency)
		assert.Equal(t, 0.0, score.OverallScore)
	})

	t.Run("すべて満点なら 100", func(t *testing.T) {
		score := usecase.CalculateProductivityScore(&model.ProductivityStats{
			PomodoroSessions: 10, ManualSessions: 0,
			CompletedGoals: 5, TotalGoals: 5,
			TotalLogDays: 84, TotalDaysInRange: 84,
		})

		assert.Equal(t, 100.0, score.PomodoroRate)
		assert.Equal(t, 100.0, score.GoalRate)
		assert.Equal(t, 100.0, score.StreakConsistency)
		assert.Equal(t, 100.0, score.OverallScore)
	})

	t.Run("目標が無ければ目標達成率は 0", func(t *testing.T) {
		score := usecase.CalculateProductivityScore(&model.ProductivityStats{
			PomodoroSessions: 10, ManualSessions: 10,
		})

		assert.Equal(t, 50.0, score.PomodoroRate)
		assert.Equal(t, 0.0, score.GoalRate)
		assert.Equal(t, 15.0, score.OverallScore)
	})

	t.Run("重みはポモドーロ 30% 目標 40% ストリーク 30%", func(t *testing.T) {
		score := usecase.CalculateProductivityScore(&model.ProductivityStats{
			PomodoroSessions: 1, ManualSessions: 0, // 100%
			CompletedGoals: 0, TotalGoals: 4, // 0%
			TotalLogDays: 42, TotalDaysInRange: 84, // 50%
		})

		assert.Equal(t, 45.0, score.OverallScore)
	})

	t.Run("割り切れない率は小数第 2 位まで丸める", func(t *testing.T) {
		score := usecase.CalculateProductivityScore(&model.ProductivityStats{
			PomodoroSessions: 1, ManualSessions: 2,
		})

		assert.Equal(t, 33.33, score.PomodoroRate)
	})
}

// ============================================================
// AI インサイト
// ============================================================

func TestGetLearningInsightsUseCase_Execute(t *testing.T) {
	newRepo := func() *mockLearningAnalyticsRepo {
		repo := new(mockLearningAnalyticsRepo)
		repo.On("GetHeatmapData", mock.Anything, uint(5)).Return([]model.HeatmapEntry{
			{DayOfWeek: 2, Hour: 21, TotalMinutes: 120},
		}, nil)
		repo.On("GetCategoryBreakdown", mock.Anything, uint(5)).Return([]model.CategoryBreakdown{
			{Category: "coding", TotalMinutes: 90},
			{Category: "reading", TotalMinutes: 10},
		}, nil)
		repo.On("GetWeeklyTrends", mock.Anything, uint(5), 12).Return([]model.WeeklyTrend{
			{WeekStart: "2026-01-05", TotalMinutes: 100},
			{WeekStart: "2026-01-12", TotalMinutes: 150},
		}, nil)
		repo.On("GetProductivityStats", mock.Anything, uint(5)).
			Return(&model.ProductivityStats{CurrentStreak: 9, LongestStreak: 9}, nil)
		return repo
	}

	t.Run("条件を満たしたインサイトを順に返す", func(t *testing.T) {
		uc := usecase.NewGetLearningInsightsUseCase(newRepo())

		insights, err := uc.Execute(context.Background(), 5)

		require.NoError(t, err)
		require.Len(t, insights, 4)
		assert.Equal(t, "peak_time", insights[0].Type)
		assert.Equal(t, "category_focus", insights[1].Type)
		assert.Equal(t, "weekly_growth_up", insights[2].Type)
		assert.Equal(t, "streak_active", insights[3].Type)
	})

	t.Run("データが無ければ空配列を返す", func(t *testing.T) {
		repo := new(mockLearningAnalyticsRepo)
		repo.On("GetHeatmapData", mock.Anything, uint(5)).Return([]model.HeatmapEntry{}, nil)
		repo.On("GetCategoryBreakdown", mock.Anything, uint(5)).Return([]model.CategoryBreakdown{}, nil)
		repo.On("GetWeeklyTrends", mock.Anything, uint(5), 12).Return([]model.WeeklyTrend{}, nil)
		repo.On("GetProductivityStats", mock.Anything, uint(5)).Return(&model.ProductivityStats{}, nil)
		uc := usecase.NewGetLearningInsightsUseCase(repo)

		insights, err := uc.Execute(context.Background(), 5)

		require.NoError(t, err)
		assert.NotNil(t, insights)
		assert.Empty(t, insights)
	})

	t.Run("DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockLearningAnalyticsRepo)
		repo.On("GetHeatmapData", mock.Anything, uint(5)).
			Return([]model.HeatmapEntry(nil), errors.New("db error"))
		uc := usecase.NewGetLearningInsightsUseCase(repo)

		_, err := uc.Execute(context.Background(), 5)

		assert.EqualError(t, err, "db error")
	})
}

func TestBuildLearningInsights(t *testing.T) {
	t.Run("カテゴリ集中は 70% 以上で出る", func(t *testing.T) {
		atThreshold := usecase.BuildLearningInsights(nil, []model.CategoryBreakdown{
			{Category: "coding", TotalMinutes: 70},
			{Category: "reading", TotalMinutes: 30},
		}, nil, &model.ProductivityStats{})
		require.Len(t, atThreshold, 1)
		assert.Equal(t, "category_focus", atThreshold[0].Type)

		below := usecase.BuildLearningInsights(nil, []model.CategoryBreakdown{
			{Category: "coding", TotalMinutes: 69},
			{Category: "reading", TotalMinutes: 31},
		}, nil, &model.ProductivityStats{})
		assert.Empty(t, below)
	})

	t.Run("週間成長は増加なら up、20% を超える減少なら down", func(t *testing.T) {
		up := usecase.BuildLearningInsights(nil, nil, []model.WeeklyTrend{
			{TotalMinutes: 100}, {TotalMinutes: 110},
		}, &model.ProductivityStats{})
		require.Len(t, up, 1)
		assert.Equal(t, "weekly_growth_up", up[0].Type)
		assert.Equal(t, 10.0, up[0].Params["rate"])

		down := usecase.BuildLearningInsights(nil, nil, []model.WeeklyTrend{
			{TotalMinutes: 100}, {TotalMinutes: 70},
		}, &model.ProductivityStats{})
		require.Len(t, down, 1)
		assert.Equal(t, "weekly_growth_down", down[0].Type)
		assert.Equal(t, 30.0, down[0].Params["rate"])

		// ちょうど 20% の減少は閾値を超えないため何も出ない。
		flat := usecase.BuildLearningInsights(nil, nil, []model.WeeklyTrend{
			{TotalMinutes: 100}, {TotalMinutes: 80},
		}, &model.ProductivityStats{})
		assert.Empty(t, flat)
	})

	t.Run("直近の週が 1 件以下なら成長判定はしない", func(t *testing.T) {
		insights := usecase.BuildLearningInsights(nil, nil, []model.WeeklyTrend{{TotalMinutes: 100}},
			&model.ProductivityStats{})
		assert.Empty(t, insights)
	})

	t.Run("ストリークは 7 日以上で active、自己ベスト未満なら record", func(t *testing.T) {
		active := usecase.BuildLearningInsights(nil, nil, nil,
			&model.ProductivityStats{CurrentStreak: 7, LongestStreak: 10})
		require.Len(t, active, 1)
		assert.Equal(t, "streak_active", active[0].Type)

		record := usecase.BuildLearningInsights(nil, nil, nil,
			&model.ProductivityStats{CurrentStreak: 2, LongestStreak: 10})
		require.Len(t, record, 1)
		assert.Equal(t, "streak_record", record[0].Type)

		none := usecase.BuildLearningInsights(nil, nil, nil,
			&model.ProductivityStats{CurrentStreak: 3, LongestStreak: 3})
		assert.Empty(t, none)
	})

	t.Run("ピークタイムは最大の曜日×時間帯を返す", func(t *testing.T) {
		insights := usecase.BuildLearningInsights([]model.HeatmapEntry{
			{DayOfWeek: 1, Hour: 9, TotalMinutes: 30},
			{DayOfWeek: 5, Hour: 22, TotalMinutes: 120},
		}, nil, nil, &model.ProductivityStats{})

		require.Len(t, insights, 1)
		assert.Equal(t, "peak_time", insights[0].Type)
		assert.Equal(t, 5, insights[0].Params["day_of_week"])
		assert.Equal(t, 22, insights[0].Params["hour"])
		assert.Equal(t, 120, insights[0].Params["minutes"])
	})

	t.Run("学習時間が全部 0 ならピークタイムは出ない", func(t *testing.T) {
		insights := usecase.BuildLearningInsights([]model.HeatmapEntry{
			{DayOfWeek: 1, Hour: 9, TotalMinutes: 0},
		}, nil, nil, &model.ProductivityStats{})
		assert.Empty(t, insights)
	})
}

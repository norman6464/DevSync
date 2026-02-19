package service

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// === ヘルパー関数 ===

func newTestAnalyticsService() (*LearningAnalyticsService, *MockLearningAnalyticsRepository) {
	repo := new(MockLearningAnalyticsRepository)
	svc := NewLearningAnalyticsService(repo)
	return svc, repo
}

// ============================================================
// ヒートマップ取得テスト
// ============================================================

func TestGetHeatmap_Success(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	heatmap := []model.HeatmapEntry{
		{DayOfWeek: 1, Hour: 21, TotalMinutes: 120},
		{DayOfWeek: 3, Hour: 20, TotalMinutes: 90},
	}
	repo.On("GetHeatmapData", uint(1)).Return(heatmap, nil)

	result, err := svc.GetHeatmap(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 120, result[0].TotalMinutes)
	assert.Equal(t, 1, result[0].DayOfWeek)
	repo.AssertExpectations(t)
}

func TestGetHeatmap_Empty(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	repo.On("GetHeatmapData", uint(1)).Return([]model.HeatmapEntry{}, nil)

	result, err := svc.GetHeatmap(1)
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestGetHeatmap_RepoError(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	repo.On("GetHeatmapData", uint(1)).Return([]model.HeatmapEntry(nil), assert.AnError)

	result, err := svc.GetHeatmap(1)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// カテゴリ別学習時間テスト
// ============================================================

func TestGetCategoryBreakdown_Success(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	// リポジトリはPercentage未設定で返す（Serviceが計算する）
	categories := []model.CategoryBreakdown{
		{Category: "coding", TotalMinutes: 300, LogCount: 10},
		{Category: "reading", TotalMinutes: 100, LogCount: 5},
		{Category: "course", TotalMinutes: 100, LogCount: 3},
	}
	repo.On("GetCategoryBreakdown", uint(1)).Return(categories, nil)

	result, err := svc.GetCategoryBreakdown(1)
	assert.NoError(t, err)
	assert.Len(t, result, 3)

	// Percentageがサービスで計算されている
	assert.Equal(t, 60.0, result[0].Percentage) // 300/500 * 100
	assert.Equal(t, 20.0, result[1].Percentage) // 100/500 * 100
	assert.Equal(t, 20.0, result[2].Percentage) // 100/500 * 100
	repo.AssertExpectations(t)
}

func TestGetCategoryBreakdown_Empty(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	repo.On("GetCategoryBreakdown", uint(1)).Return([]model.CategoryBreakdown{}, nil)

	result, err := svc.GetCategoryBreakdown(1)
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestGetCategoryBreakdown_SingleCategory(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	categories := []model.CategoryBreakdown{
		{Category: "coding", TotalMinutes: 200, LogCount: 8},
	}
	repo.On("GetCategoryBreakdown", uint(1)).Return(categories, nil)

	result, err := svc.GetCategoryBreakdown(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 100.0, result[0].Percentage) // 唯一のカテゴリは100%
	repo.AssertExpectations(t)
}

func TestGetCategoryBreakdown_AllZeroMinutes(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	categories := []model.CategoryBreakdown{
		{Category: "coding", TotalMinutes: 0, LogCount: 2},
		{Category: "reading", TotalMinutes: 0, LogCount: 1},
	}
	repo.On("GetCategoryBreakdown", uint(1)).Return(categories, nil)

	result, err := svc.GetCategoryBreakdown(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	// 合計0分の場合、Percentageは0
	assert.Equal(t, 0.0, result[0].Percentage)
	assert.Equal(t, 0.0, result[1].Percentage)
	repo.AssertExpectations(t)
}

func TestGetCategoryBreakdown_RepoError(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	repo.On("GetCategoryBreakdown", uint(1)).Return([]model.CategoryBreakdown(nil), assert.AnError)

	result, err := svc.GetCategoryBreakdown(1)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// 週間トレンドテスト
// ============================================================

func TestGetWeeklyTrends_Success(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	trends := []model.WeeklyTrend{
		{WeekStart: "2025-01-06", TotalMinutes: 300, LogCount: 10},
		{WeekStart: "2025-01-13", TotalMinutes: 450, LogCount: 15},
		{WeekStart: "2025-01-20", TotalMinutes: 200, LogCount: 7},
	}
	repo.On("GetWeeklyTrends", uint(1), 12).Return(trends, nil)

	result, err := svc.GetWeeklyTrends(1, 12)
	assert.NoError(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "2025-01-06", result[0].WeekStart)
	assert.Equal(t, 450, result[1].TotalMinutes)
	repo.AssertExpectations(t)
}

func TestGetWeeklyTrends_DefaultWeeks(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	repo.On("GetWeeklyTrends", uint(1), 12).Return([]model.WeeklyTrend{}, nil)

	// weeks=0の場合、デフォルトで12週
	result, err := svc.GetWeeklyTrends(1, 0)
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestGetWeeklyTrends_RepoError(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	repo.On("GetWeeklyTrends", uint(1), 12).Return([]model.WeeklyTrend(nil), assert.AnError)

	result, err := svc.GetWeeklyTrends(1, 12)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// 生産性スコア計算テスト
// ============================================================

func TestCalculateProductivityScore_AllActive(t *testing.T) {
	stats := &model.ProductivityStats{
		PomodoroSessions: 20,
		ManualSessions:   10,
		CompletedGoals:   8,
		TotalGoals:       10,
		CurrentStreak:    14,
		LongestStreak:    14,
		TotalLogDays:     60,
		TotalDaysInRange: 84,
	}

	score := CalculateProductivityScore(stats)

	// ポモドーロ率: 20 / (20+10) * 100 = 66.67
	assert.InDelta(t, 66.67, score.PomodoroRate, 0.01)
	// 目標達成率: 8/10 * 100 = 80
	assert.Equal(t, 80.0, score.GoalRate)
	// ストリーク継続率: 60/84 * 100 = 71.43
	assert.InDelta(t, 71.43, score.StreakConsistency, 0.01)
	// 総合スコア: (66.67 * 0.3 + 80 * 0.4 + 71.43 * 0.3) = 20.0 + 32.0 + 21.43 = 73.43
	assert.InDelta(t, 73.43, score.OverallScore, 0.01)
}

func TestCalculateProductivityScore_AllZero(t *testing.T) {
	stats := &model.ProductivityStats{}

	score := CalculateProductivityScore(stats)

	assert.Equal(t, 0.0, score.PomodoroRate)
	assert.Equal(t, 0.0, score.GoalRate)
	assert.Equal(t, 0.0, score.StreakConsistency)
	assert.Equal(t, 0.0, score.OverallScore)
}

func TestCalculateProductivityScore_NoGoals(t *testing.T) {
	stats := &model.ProductivityStats{
		PomodoroSessions: 5,
		ManualSessions:   5,
		TotalLogDays:     30,
		TotalDaysInRange: 84,
	}

	score := CalculateProductivityScore(stats)

	assert.Equal(t, 50.0, score.PomodoroRate) // 5/10 * 100
	assert.Equal(t, 0.0, score.GoalRate)       // 目標なしは0
	assert.InDelta(t, 35.71, score.StreakConsistency, 0.01) // 30/84 * 100
	repo_expected := 50.0*0.3 + 0.0*0.4 + 35.71*0.3
	assert.InDelta(t, repo_expected, score.OverallScore, 0.01)
}

func TestCalculateProductivityScore_PerfectScore(t *testing.T) {
	stats := &model.ProductivityStats{
		PomodoroSessions: 100,
		ManualSessions:   0,
		CompletedGoals:   10,
		TotalGoals:       10,
		TotalLogDays:     84,
		TotalDaysInRange: 84,
	}

	score := CalculateProductivityScore(stats)

	assert.Equal(t, 100.0, score.PomodoroRate)
	assert.Equal(t, 100.0, score.GoalRate)
	assert.Equal(t, 100.0, score.StreakConsistency)
	assert.Equal(t, 100.0, score.OverallScore)
}

func TestGetProductivityScore_Success(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	stats := &model.ProductivityStats{
		PomodoroSessions: 10,
		ManualSessions:   5,
		CompletedGoals:   3,
		TotalGoals:       5,
		TotalLogDays:     40,
		TotalDaysInRange: 84,
	}
	repo.On("GetProductivityStats", uint(1)).Return(stats, nil)

	score, err := svc.GetProductivityScore(1)
	assert.NoError(t, err)
	assert.NotNil(t, score)
	assert.InDelta(t, 66.67, score.PomodoroRate, 0.01)
	assert.Equal(t, 60.0, score.GoalRate)
	repo.AssertExpectations(t)
}

func TestGetProductivityScore_RepoError(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	repo.On("GetProductivityStats", uint(1)).Return((*model.ProductivityStats)(nil), assert.AnError)

	score, err := svc.GetProductivityScore(1)
	assert.Error(t, err)
	assert.Nil(t, score)
	repo.AssertExpectations(t)
}

// ============================================================
// AIインサイト生成テスト
// ============================================================

func TestGetInsights_PeakTime(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	// 火曜21時が最大学習時間
	heatmap := []model.HeatmapEntry{
		{DayOfWeek: 2, Hour: 21, TotalMinutes: 200},
		{DayOfWeek: 1, Hour: 20, TotalMinutes: 50},
		{DayOfWeek: 3, Hour: 19, TotalMinutes: 80},
	}
	categories := []model.CategoryBreakdown{
		{Category: "coding", TotalMinutes: 300, LogCount: 10},
	}
	trends := []model.WeeklyTrend{
		{WeekStart: "2025-01-13", TotalMinutes: 300, LogCount: 10},
		{WeekStart: "2025-01-20", TotalMinutes: 400, LogCount: 12},
	}
	stats := &model.ProductivityStats{
		CurrentStreak: 7, LongestStreak: 7,
		TotalLogDays: 30, TotalDaysInRange: 84,
	}

	repo.On("GetHeatmapData", uint(1)).Return(heatmap, nil)
	repo.On("GetCategoryBreakdown", uint(1)).Return(categories, nil)
	repo.On("GetWeeklyTrends", uint(1), 12).Return(trends, nil)
	repo.On("GetProductivityStats", uint(1)).Return(stats, nil)

	insights, err := svc.GetInsights(1)
	assert.NoError(t, err)
	assert.NotEmpty(t, insights)

	// peak_timeインサイトが含まれることを検証
	hasPeakTime := false
	for _, ins := range insights {
		if ins.Type == "peak_time" {
			hasPeakTime = true
		}
	}
	assert.True(t, hasPeakTime)
	repo.AssertExpectations(t)
}

func TestGetInsights_StreakTrend(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	heatmap := []model.HeatmapEntry{}
	categories := []model.CategoryBreakdown{}
	trends := []model.WeeklyTrend{
		{WeekStart: "2025-01-13", TotalMinutes: 100, LogCount: 3},
		{WeekStart: "2025-01-20", TotalMinutes: 200, LogCount: 7},
	}
	stats := &model.ProductivityStats{
		CurrentStreak: 14, LongestStreak: 14,
		TotalLogDays: 50, TotalDaysInRange: 84,
	}

	repo.On("GetHeatmapData", uint(1)).Return(heatmap, nil)
	repo.On("GetCategoryBreakdown", uint(1)).Return(categories, nil)
	repo.On("GetWeeklyTrends", uint(1), 12).Return(trends, nil)
	repo.On("GetProductivityStats", uint(1)).Return(stats, nil)

	insights, err := svc.GetInsights(1)
	assert.NoError(t, err)

	// streak_activeインサイトが含まれることを検証
	hasStreakActive := false
	for _, ins := range insights {
		if ins.Type == "streak_active" {
			hasStreakActive = true
		}
	}
	assert.True(t, hasStreakActive)
	repo.AssertExpectations(t)
}

func TestGetInsights_WeeklyGrowth(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	heatmap := []model.HeatmapEntry{}
	categories := []model.CategoryBreakdown{}
	// 前週比で増加
	trends := []model.WeeklyTrend{
		{WeekStart: "2025-01-13", TotalMinutes: 100, LogCount: 3},
		{WeekStart: "2025-01-20", TotalMinutes: 150, LogCount: 5},
	}
	stats := &model.ProductivityStats{
		TotalLogDays: 10, TotalDaysInRange: 84,
	}

	repo.On("GetHeatmapData", uint(1)).Return(heatmap, nil)
	repo.On("GetCategoryBreakdown", uint(1)).Return(categories, nil)
	repo.On("GetWeeklyTrends", uint(1), 12).Return(trends, nil)
	repo.On("GetProductivityStats", uint(1)).Return(stats, nil)

	insights, err := svc.GetInsights(1)
	assert.NoError(t, err)

	hasWeeklyGrowthUp := false
	for _, ins := range insights {
		if ins.Type == "weekly_growth_up" {
			hasWeeklyGrowthUp = true
		}
	}
	assert.True(t, hasWeeklyGrowthUp)
	repo.AssertExpectations(t)
}

func TestGetInsights_Empty(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	repo.On("GetHeatmapData", uint(1)).Return([]model.HeatmapEntry{}, nil)
	repo.On("GetCategoryBreakdown", uint(1)).Return([]model.CategoryBreakdown{}, nil)
	repo.On("GetWeeklyTrends", uint(1), 12).Return([]model.WeeklyTrend{}, nil)
	repo.On("GetProductivityStats", uint(1)).Return(&model.ProductivityStats{}, nil)

	insights, err := svc.GetInsights(1)
	assert.NoError(t, err)
	assert.NotNil(t, insights) // 空でもnilにはならない
	repo.AssertExpectations(t)
}

func TestGetInsights_RepoError(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	repo.On("GetHeatmapData", uint(1)).Return([]model.HeatmapEntry(nil), assert.AnError)

	insights, err := svc.GetInsights(1)
	assert.Error(t, err)
	assert.Nil(t, insights)
	repo.AssertExpectations(t)
}

// ============================================================
// analyzeWeeklyGrowth テスト（GetInsights経由）
// ============================================================

func TestGetInsights_WeeklyGrowthUp(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	// 直近2週間で成長あり（prev=100, current=150 → +50%成長）
	trends := []model.WeeklyTrend{
		{TotalMinutes: 100},
		{TotalMinutes: 150},
	}

	repo.On("GetHeatmapData", uint(1)).Return([]model.HeatmapEntry{}, nil)
	repo.On("GetCategoryBreakdown", uint(1)).Return([]model.CategoryBreakdown{}, nil)
	repo.On("GetWeeklyTrends", uint(1), 12).Return(trends, nil)
	repo.On("GetProductivityStats", uint(1)).Return(&model.ProductivityStats{}, nil)

	insights, err := svc.GetInsights(1)
	assert.NoError(t, err)
	// weekly_growth_upのインサイトが含まれているはず
	found := false
	for _, ins := range insights {
		if ins.Type == "weekly_growth_up" {
			found = true
			break
		}
	}
	assert.True(t, found, "weekly_growth_up インサイトが含まれていること")
	repo.AssertExpectations(t)
}

func TestGetInsights_WeeklyGrowthDown(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	// 直近2週間で大きく下落（prev=200, current=100 → -50%下落）
	trends := []model.WeeklyTrend{
		{TotalMinutes: 200},
		{TotalMinutes: 100},
	}

	repo.On("GetHeatmapData", uint(1)).Return([]model.HeatmapEntry{}, nil)
	repo.On("GetCategoryBreakdown", uint(1)).Return([]model.CategoryBreakdown{}, nil)
	repo.On("GetWeeklyTrends", uint(1), 12).Return(trends, nil)
	repo.On("GetProductivityStats", uint(1)).Return(&model.ProductivityStats{}, nil)

	insights, err := svc.GetInsights(1)
	assert.NoError(t, err)
	// weekly_growth_downのインサイトが含まれているはず
	found := false
	for _, ins := range insights {
		if ins.Type == "weekly_growth_down" {
			found = true
			break
		}
	}
	assert.True(t, found, "weekly_growth_down インサイトが含まれていること")
	repo.AssertExpectations(t)
}

func TestGetInsights_WeeklyGrowthNoChange(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	// 変化なし（-20%以内の微小下落はnilを返す）
	trends := []model.WeeklyTrend{
		{TotalMinutes: 100},
		{TotalMinutes: 90}, // -10%（-20%以内なのでweekly_growth_downなし）
	}

	repo.On("GetHeatmapData", uint(1)).Return([]model.HeatmapEntry{}, nil)
	repo.On("GetCategoryBreakdown", uint(1)).Return([]model.CategoryBreakdown{}, nil)
	repo.On("GetWeeklyTrends", uint(1), 12).Return(trends, nil)
	repo.On("GetProductivityStats", uint(1)).Return(&model.ProductivityStats{}, nil)

	insights, err := svc.GetInsights(1)
	assert.NoError(t, err)
	// weekly_growth系インサイトは含まれないはず
	for _, ins := range insights {
		assert.NotEqual(t, "weekly_growth_up", ins.Type)
		assert.NotEqual(t, "weekly_growth_down", ins.Type)
	}
	repo.AssertExpectations(t)
}

func TestGetInsights_WeeklyGrowthPrevZero(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	// prev=0の場合は除算を避けてnilを返す
	trends := []model.WeeklyTrend{
		{TotalMinutes: 0},
		{TotalMinutes: 100},
	}

	repo.On("GetHeatmapData", uint(1)).Return([]model.HeatmapEntry{}, nil)
	repo.On("GetCategoryBreakdown", uint(1)).Return([]model.CategoryBreakdown{}, nil)
	repo.On("GetWeeklyTrends", uint(1), 12).Return(trends, nil)
	repo.On("GetProductivityStats", uint(1)).Return(&model.ProductivityStats{}, nil)

	insights, err := svc.GetInsights(1)
	assert.NoError(t, err)
	// weekly_growth系インサイトは含まれないはず（prev=0で除算スキップ）
	for _, ins := range insights {
		assert.NotEqual(t, "weekly_growth_up", ins.Type)
		assert.NotEqual(t, "weekly_growth_down", ins.Type)
	}
	repo.AssertExpectations(t)
}

// ============================================================
// analyzeStreak 追加テスト
// ============================================================

func TestGetInsights_StreakRecord(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	// CurrentStreak < 7 かつ LongestStreak > CurrentStreak → streak_record
	stats := &model.ProductivityStats{
		CurrentStreak: 3,
		LongestStreak: 10,
	}

	repo.On("GetHeatmapData", uint(1)).Return([]model.HeatmapEntry{}, nil)
	repo.On("GetCategoryBreakdown", uint(1)).Return([]model.CategoryBreakdown{}, nil)
	repo.On("GetWeeklyTrends", uint(1), 12).Return([]model.WeeklyTrend{}, nil)
	repo.On("GetProductivityStats", uint(1)).Return(stats, nil)

	insights, err := svc.GetInsights(1)
	assert.NoError(t, err)

	hasStreakRecord := false
	for _, ins := range insights {
		if ins.Type == "streak_record" {
			hasStreakRecord = true
		}
	}
	assert.True(t, hasStreakRecord)
	repo.AssertExpectations(t)
}

func TestGetInsights_StreakNil(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	// CurrentStreak < 7 かつ LongestStreak == 0 → analyzeStreak nil
	stats := &model.ProductivityStats{
		CurrentStreak: 2,
		LongestStreak: 0,
	}

	repo.On("GetHeatmapData", uint(1)).Return([]model.HeatmapEntry{}, nil)
	repo.On("GetCategoryBreakdown", uint(1)).Return([]model.CategoryBreakdown{}, nil)
	repo.On("GetWeeklyTrends", uint(1), 12).Return([]model.WeeklyTrend{}, nil)
	repo.On("GetProductivityStats", uint(1)).Return(stats, nil)

	insights, err := svc.GetInsights(1)
	assert.NoError(t, err)
	for _, ins := range insights {
		assert.NotEqual(t, "streak_active", ins.Type)
		assert.NotEqual(t, "streak_record", ins.Type)
	}
	repo.AssertExpectations(t)
}

// ============================================================
// analyzeCategoryFocus 追加テスト
// ============================================================

func TestGetInsights_CategoryFocusBelow70(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	// 最大カテゴリが全体の50%（70%未満）→ nil
	categories := []model.CategoryBreakdown{
		{Category: "coding", TotalMinutes: 100},
		{Category: "reading", TotalMinutes: 100},
	}

	repo.On("GetHeatmapData", uint(1)).Return([]model.HeatmapEntry{}, nil)
	repo.On("GetCategoryBreakdown", uint(1)).Return(categories, nil)
	repo.On("GetWeeklyTrends", uint(1), 12).Return([]model.WeeklyTrend{}, nil)
	repo.On("GetProductivityStats", uint(1)).Return(&model.ProductivityStats{}, nil)

	insights, err := svc.GetInsights(1)
	assert.NoError(t, err)
	for _, ins := range insights {
		assert.NotEqual(t, "category_focus", ins.Type)
	}
	repo.AssertExpectations(t)
}

func TestGetInsights_CategoryFocusZeroMinutes(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	// 全カテゴリのTotalMinutes=0 → totalMinutes==0でnil
	categories := []model.CategoryBreakdown{
		{Category: "coding", TotalMinutes: 0},
	}

	repo.On("GetHeatmapData", uint(1)).Return([]model.HeatmapEntry{}, nil)
	repo.On("GetCategoryBreakdown", uint(1)).Return(categories, nil)
	repo.On("GetWeeklyTrends", uint(1), 12).Return([]model.WeeklyTrend{}, nil)
	repo.On("GetProductivityStats", uint(1)).Return(&model.ProductivityStats{}, nil)

	insights, err := svc.GetInsights(1)
	assert.NoError(t, err)
	for _, ins := range insights {
		assert.NotEqual(t, "category_focus", ins.Type)
	}
	repo.AssertExpectations(t)
}

// ============================================================
// GetInsights エラーパス 追加テスト
// ============================================================

func TestGetInsights_CategoryBreakdownError(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	repo.On("GetHeatmapData", uint(1)).Return([]model.HeatmapEntry{}, nil)
	repo.On("GetCategoryBreakdown", uint(1)).Return([]model.CategoryBreakdown(nil), assert.AnError)

	insights, err := svc.GetInsights(1)
	assert.Error(t, err)
	assert.Nil(t, insights)
	repo.AssertExpectations(t)
}

func TestGetInsights_WeeklyTrendsError(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	repo.On("GetHeatmapData", uint(1)).Return([]model.HeatmapEntry{}, nil)
	repo.On("GetCategoryBreakdown", uint(1)).Return([]model.CategoryBreakdown{}, nil)
	repo.On("GetWeeklyTrends", uint(1), 12).Return([]model.WeeklyTrend(nil), assert.AnError)

	insights, err := svc.GetInsights(1)
	assert.Error(t, err)
	assert.Nil(t, insights)
	repo.AssertExpectations(t)
}

func TestGetInsights_ProductivityStatsError(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	repo.On("GetHeatmapData", uint(1)).Return([]model.HeatmapEntry{}, nil)
	repo.On("GetCategoryBreakdown", uint(1)).Return([]model.CategoryBreakdown{}, nil)
	repo.On("GetWeeklyTrends", uint(1), 12).Return([]model.WeeklyTrend{}, nil)
	repo.On("GetProductivityStats", uint(1)).Return((*model.ProductivityStats)(nil), assert.AnError)

	insights, err := svc.GetInsights(1)
	assert.Error(t, err)
	assert.Nil(t, insights)
	repo.AssertExpectations(t)
}

// ============================================================
// analyzePeakTime 追加テスト
// ============================================================

func TestGetInsights_PeakTimeAllZeroMinutes(t *testing.T) {
	svc, repo := newTestAnalyticsService()

	// 全エントリのTotalMinutes=0 → maxEntry.TotalMinutes==0でnil
	heatmap := []model.HeatmapEntry{
		{DayOfWeek: 1, Hour: 10, TotalMinutes: 0},
		{DayOfWeek: 2, Hour: 20, TotalMinutes: 0},
	}

	repo.On("GetHeatmapData", uint(1)).Return(heatmap, nil)
	repo.On("GetCategoryBreakdown", uint(1)).Return([]model.CategoryBreakdown{}, nil)
	repo.On("GetWeeklyTrends", uint(1), 12).Return([]model.WeeklyTrend{}, nil)
	repo.On("GetProductivityStats", uint(1)).Return(&model.ProductivityStats{}, nil)

	insights, err := svc.GetInsights(1)
	assert.NoError(t, err)
	for _, ins := range insights {
		assert.NotEqual(t, "peak_time", ins.Type)
	}
	repo.AssertExpectations(t)
}

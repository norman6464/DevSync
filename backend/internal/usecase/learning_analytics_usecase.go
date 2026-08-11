package usecase

import (
	"context"
	"math"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// defaultWeeklyTrendsWeeks は週間トレンドのデフォルト週数。
const defaultWeeklyTrendsWeeks = 12

// daysInWeek は曜日別サマリーの要素数。
const daysInWeek = 7

// roundToTwoDecimals は割合を小数第 2 位で四捨五入する。
func roundToTwoDecimals(v float64) float64 {
	return math.Round(v*100) / 100
}

// percentageOf は part / total を百分率（小数第 2 位まで）で返す。total が 0 以下なら 0。
func percentageOf(part, total int) float64 {
	if total <= 0 {
		return 0
	}
	return roundToTwoDecimals(float64(part) / float64(total) * 100)
}

// GetLearningHeatmapUseCase は学習時間のヒートマップを取得する。
type GetLearningHeatmapUseCase struct {
	analytics repository.LearningAnalyticsRepository
}

// NewGetLearningHeatmapUseCase は GetLearningHeatmapUseCase を生成する。
func NewGetLearningHeatmapUseCase(analytics repository.LearningAnalyticsRepository) *GetLearningHeatmapUseCase {
	return &GetLearningHeatmapUseCase{analytics: analytics}
}

// Execute は曜日×時間帯ごとの学習時間を返す。
func (uc *GetLearningHeatmapUseCase) Execute(ctx context.Context, userID uint) ([]model.HeatmapEntry, error) {
	return uc.analytics.GetHeatmapData(ctx, userID)
}

// GetCategoryBreakdownUseCase はカテゴリ別の学習時間と割合を取得する。
type GetCategoryBreakdownUseCase struct {
	analytics repository.LearningAnalyticsRepository
}

// NewGetCategoryBreakdownUseCase は GetCategoryBreakdownUseCase を生成する。
func NewGetCategoryBreakdownUseCase(analytics repository.LearningAnalyticsRepository) *GetCategoryBreakdownUseCase {
	return &GetCategoryBreakdownUseCase{analytics: analytics}
}

// Execute はカテゴリ別の学習時間に、全体に占める割合を付けて返す。
func (uc *GetCategoryBreakdownUseCase) Execute(ctx context.Context, userID uint) ([]model.CategoryBreakdown, error) {
	categories, err := uc.analytics.GetCategoryBreakdown(ctx, userID)
	if err != nil {
		return nil, err
	}
	return WithCategoryPercentages(categories), nil
}

// WithCategoryPercentages はカテゴリ別集計に全体比の割合を埋めて返す純粋関数。
func WithCategoryPercentages(categories []model.CategoryBreakdown) []model.CategoryBreakdown {
	totalMinutes := 0
	for _, c := range categories {
		totalMinutes += c.TotalMinutes
	}
	for i := range categories {
		categories[i].Percentage = percentageOf(categories[i].TotalMinutes, totalMinutes)
	}
	return categories
}

// GetWeeklyTrendsUseCase は週間の学習トレンドを取得する。
type GetWeeklyTrendsUseCase struct {
	analytics repository.LearningAnalyticsRepository
}

// NewGetWeeklyTrendsUseCase は GetWeeklyTrendsUseCase を生成する。
func NewGetWeeklyTrendsUseCase(analytics repository.LearningAnalyticsRepository) *GetWeeklyTrendsUseCase {
	return &GetWeeklyTrendsUseCase{analytics: analytics}
}

// Execute は直近 weeks 週のトレンドを返す。weeks が 0 以下ならデフォルト週数を使う。
func (uc *GetWeeklyTrendsUseCase) Execute(ctx context.Context, userID uint, weeks int) ([]model.WeeklyTrend, error) {
	if weeks <= 0 {
		weeks = defaultWeeklyTrendsWeeks
	}
	return uc.analytics.GetWeeklyTrends(ctx, userID, weeks)
}

// GetDayOfWeekSummaryUseCase は曜日別の学習サマリーを取得する。
type GetDayOfWeekSummaryUseCase struct {
	analytics repository.LearningAnalyticsRepository
}

// NewGetDayOfWeekSummaryUseCase は GetDayOfWeekSummaryUseCase を生成する。
func NewGetDayOfWeekSummaryUseCase(analytics repository.LearningAnalyticsRepository) *GetDayOfWeekSummaryUseCase {
	return &GetDayOfWeekSummaryUseCase{analytics: analytics}
}

// Execute はヒートマップを曜日ごとに集計して返す。
func (uc *GetDayOfWeekSummaryUseCase) Execute(ctx context.Context, userID uint) ([]model.DayOfWeekSummary, error) {
	heatmap, err := uc.analytics.GetHeatmapData(ctx, userID)
	if err != nil {
		return nil, err
	}
	return AggregateDayOfWeek(heatmap), nil
}

// AggregateDayOfWeek はヒートマップを曜日別に集計する純粋関数。
// ログの有無にかかわらず日曜〜土曜の 7 件を返す。
func AggregateDayOfWeek(heatmap []model.HeatmapEntry) []model.DayOfWeekSummary {
	result := make([]model.DayOfWeekSummary, daysInWeek)
	for i := range result {
		result[i].DayOfWeek = i
	}

	for _, entry := range heatmap {
		if entry.DayOfWeek < 0 || entry.DayOfWeek >= daysInWeek {
			continue
		}
		result[entry.DayOfWeek].TotalMinutes += entry.TotalMinutes
		result[entry.DayOfWeek].LogCount++
	}

	for i := range result {
		if result[i].LogCount > 0 {
			result[i].AverageMinutes = result[i].TotalMinutes / result[i].LogCount
		}
	}
	return result
}

// GetProductivityScoreUseCase は生産性スコアを取得する。
type GetProductivityScoreUseCase struct {
	analytics repository.LearningAnalyticsRepository
}

// NewGetProductivityScoreUseCase は GetProductivityScoreUseCase を生成する。
func NewGetProductivityScoreUseCase(analytics repository.LearningAnalyticsRepository) *GetProductivityScoreUseCase {
	return &GetProductivityScoreUseCase{analytics: analytics}
}

// Execute は統計を集計し、生産性スコアを算出して返す。
func (uc *GetProductivityScoreUseCase) Execute(ctx context.Context, userID uint) (*model.ProductivityScore, error) {
	stats, err := uc.analytics.GetProductivityStats(ctx, userID)
	if err != nil {
		return nil, err
	}
	return CalculateProductivityScore(stats), nil
}

// CalculateProductivityScore は統計から生産性スコアを算出する純粋関数。
// 重みはポモドーロ活用率 30%・目標達成率 40%・ストリーク継続率 30%。
func CalculateProductivityScore(stats *model.ProductivityStats) *model.ProductivityScore {
	score := &model.ProductivityScore{}

	score.PomodoroRate = percentageOf(stats.PomodoroSessions, stats.PomodoroSessions+stats.ManualSessions)
	score.GoalRate = percentageOf(stats.CompletedGoals, stats.TotalGoals)
	score.StreakConsistency = percentageOf(stats.TotalLogDays, stats.TotalDaysInRange)
	score.OverallScore = roundToTwoDecimals(
		score.PomodoroRate*0.3 + score.GoalRate*0.4 + score.StreakConsistency*0.3)

	return score
}

// GetLearningInsightsUseCase は学習データから AI インサイトを生成する。
type GetLearningInsightsUseCase struct {
	analytics repository.LearningAnalyticsRepository
}

// NewGetLearningInsightsUseCase は GetLearningInsightsUseCase を生成する。
func NewGetLearningInsightsUseCase(analytics repository.LearningAnalyticsRepository) *GetLearningInsightsUseCase {
	return &GetLearningInsightsUseCase{analytics: analytics}
}

// Execute はヒートマップ・カテゴリ・トレンド・生産性統計をまとめて取得し、インサイトを組み立てる。
func (uc *GetLearningInsightsUseCase) Execute(ctx context.Context, userID uint) ([]model.AIInsight, error) {
	heatmap, err := uc.analytics.GetHeatmapData(ctx, userID)
	if err != nil {
		return nil, err
	}
	categories, err := uc.analytics.GetCategoryBreakdown(ctx, userID)
	if err != nil {
		return nil, err
	}
	trends, err := uc.analytics.GetWeeklyTrends(ctx, userID, defaultWeeklyTrendsWeeks)
	if err != nil {
		return nil, err
	}
	stats, err := uc.analytics.GetProductivityStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	return BuildLearningInsights(heatmap, categories, trends, stats), nil
}

// BuildLearningInsights は各集計から該当するインサイトだけを順に並べて返す純粋関数。
func BuildLearningInsights(
	heatmap []model.HeatmapEntry,
	categories []model.CategoryBreakdown,
	trends []model.WeeklyTrend,
	stats *model.ProductivityStats,
) []model.AIInsight {
	insights := []model.AIInsight{}
	for _, insight := range []*model.AIInsight{
		analyzePeakTime(heatmap),
		analyzeCategoryFocus(categories),
		analyzeWeeklyGrowth(trends),
		analyzeStreak(stats),
	} {
		if insight != nil {
			insights = append(insights, *insight)
		}
	}
	return insights
}

// analyzePeakTime は最も学習時間が長い曜日×時間帯を特定する。
func analyzePeakTime(heatmap []model.HeatmapEntry) *model.AIInsight {
	if len(heatmap) == 0 {
		return nil
	}

	var maxEntry model.HeatmapEntry
	for _, entry := range heatmap {
		if entry.TotalMinutes > maxEntry.TotalMinutes {
			maxEntry = entry
		}
	}
	if maxEntry.TotalMinutes == 0 {
		return nil
	}

	return &model.AIInsight{
		Type: "peak_time",
		Params: map[string]interface{}{
			"day_of_week": maxEntry.DayOfWeek,
			"hour":        maxEntry.Hour,
			"minutes":     maxEntry.TotalMinutes,
		},
	}
}

// analyzeCategoryFocus は特定カテゴリへの偏り（70% 以上）を検出する。
func analyzeCategoryFocus(categories []model.CategoryBreakdown) *model.AIInsight {
	if len(categories) == 0 {
		return nil
	}

	totalMinutes := 0
	var topCategory model.CategoryBreakdown
	for _, c := range categories {
		totalMinutes += c.TotalMinutes
		if c.TotalMinutes > topCategory.TotalMinutes {
			topCategory = c
		}
	}
	if totalMinutes == 0 {
		return nil
	}

	percentage := math.Round(float64(topCategory.TotalMinutes) / float64(totalMinutes) * 100)
	if percentage < 70 {
		return nil
	}

	return &model.AIInsight{
		Type: "category_focus",
		Params: map[string]interface{}{
			"percentage": percentage,
			"category":   topCategory.Category,
		},
	}
}

// analyzeWeeklyGrowth は直近 2 週の増減から成長傾向を判定する。
func analyzeWeeklyGrowth(trends []model.WeeklyTrend) *model.AIInsight {
	if len(trends) < 2 {
		return nil
	}

	prev := trends[len(trends)-2]
	current := trends[len(trends)-1]
	if prev.TotalMinutes == 0 {
		return nil
	}

	growthRate := math.Round(float64(current.TotalMinutes-prev.TotalMinutes) / float64(prev.TotalMinutes) * 100)
	switch {
	case growthRate > 0:
		return &model.AIInsight{
			Type:   "weekly_growth_up",
			Params: map[string]interface{}{"rate": growthRate},
		}
	case growthRate < -20:
		return &model.AIInsight{
			Type:   "weekly_growth_down",
			Params: map[string]interface{}{"rate": math.Abs(growthRate)},
		}
	}
	return nil
}

// analyzeStreak は継続中のストリークか、自己ベストとの差からインサイトを作る。
func analyzeStreak(stats *model.ProductivityStats) *model.AIInsight {
	if stats.CurrentStreak >= 7 {
		return &model.AIInsight{
			Type:   "streak_active",
			Params: map[string]interface{}{"days": stats.CurrentStreak},
		}
	}
	if stats.LongestStreak > 0 && stats.CurrentStreak < stats.LongestStreak {
		return &model.AIInsight{
			Type: "streak_record",
			Params: map[string]interface{}{
				"longest": stats.LongestStreak,
				"current": stats.CurrentStreak,
			},
		}
	}
	return nil
}

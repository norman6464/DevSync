package service

import (
	"math"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// LearningAnalyticsService は学習分析のビジネスロジックを提供する。
// ヒートマップ、カテゴリ別集計、生産性スコア計算、AIインサイト生成を担当する。
type LearningAnalyticsService struct {
	repo repository.LearningAnalyticsRepositoryInterface
}

// NewLearningAnalyticsService は新しいLearningAnalyticsServiceインスタンスを生成する。
func NewLearningAnalyticsService(repo repository.LearningAnalyticsRepositoryInterface) *LearningAnalyticsService {
	return &LearningAnalyticsService{repo: repo}
}

// GetHeatmap は指定ユーザーの学習時間ヒートマップデータを取得する。
func (s *LearningAnalyticsService) GetHeatmap(userID uint) ([]model.HeatmapEntry, error) {
	return s.repo.GetHeatmapData(userID)
}

// GetCategoryBreakdown は指定ユーザーのカテゴリ別学習時間を取得し、割合を計算する。
func (s *LearningAnalyticsService) GetCategoryBreakdown(userID uint) ([]model.CategoryBreakdown, error) {
	categories, err := s.repo.GetCategoryBreakdown(userID)
	if err != nil {
		return nil, err
	}

	// 合計学習時間を計算
	totalMinutes := 0
	for _, c := range categories {
		totalMinutes += c.TotalMinutes
	}

	// 割合を計算
	for i := range categories {
		if totalMinutes > 0 {
			categories[i].Percentage = math.Round(float64(categories[i].TotalMinutes)/float64(totalMinutes)*10000) / 100
		} else {
			categories[i].Percentage = 0
		}
	}

	return categories, nil
}

// GetWeeklyTrends は指定ユーザーの週間学習トレンドを取得する。
// weeksが0以下の場合はデフォルトの12週に設定される。
func (s *LearningAnalyticsService) GetWeeklyTrends(userID uint, weeks int) ([]model.WeeklyTrend, error) {
	if weeks <= 0 {
		weeks = 12
	}
	return s.repo.GetWeeklyTrends(userID, weeks)
}

// GetProductivityScore は指定ユーザーの生産性スコアを計算して返す。
func (s *LearningAnalyticsService) GetProductivityScore(userID uint) (*model.ProductivityScore, error) {
	stats, err := s.repo.GetProductivityStats(userID)
	if err != nil {
		return nil, err
	}
	return CalculateProductivityScore(stats), nil
}

// CalculateProductivityScore はProductivityStatsから生産性スコアを算出する純粋関数。
// 重み: ポモドーロ活用率30%、目標達成率40%、ストリーク継続率30%
func CalculateProductivityScore(stats *model.ProductivityStats) *model.ProductivityScore {
	score := &model.ProductivityScore{}

	// ポモドーロ活用率: ポモドーロセッション数 / 全セッション数 * 100
	totalSessions := stats.PomodoroSessions + stats.ManualSessions
	if totalSessions > 0 {
		score.PomodoroRate = math.Round(float64(stats.PomodoroSessions)/float64(totalSessions)*10000) / 100
	}

	// 目標達成率: 完了目標数 / 全目標数 * 100
	if stats.TotalGoals > 0 {
		score.GoalRate = math.Round(float64(stats.CompletedGoals)/float64(stats.TotalGoals)*10000) / 100
	}

	// ストリーク継続率: 学習記録日数 / 期間内総日数 * 100
	if stats.TotalDaysInRange > 0 {
		score.StreakConsistency = math.Round(float64(stats.TotalLogDays)/float64(stats.TotalDaysInRange)*10000) / 100
	}

	// 総合スコア: ポモドーロ30% + 目標40% + ストリーク30%
	score.OverallScore = math.Round((score.PomodoroRate*0.3+score.GoalRate*0.4+score.StreakConsistency*0.3)*100) / 100

	return score
}

// GetInsights は指定ユーザーの学習データからAIインサイトを生成する。
// ヒートマップ、カテゴリ、トレンド、生産性データを総合的に分析する。
func (s *LearningAnalyticsService) GetInsights(userID uint) ([]model.AIInsight, error) {
	heatmap, err := s.repo.GetHeatmapData(userID)
	if err != nil {
		return nil, err
	}

	categories, err := s.repo.GetCategoryBreakdown(userID)
	if err != nil {
		return nil, err
	}

	trends, err := s.repo.GetWeeklyTrends(userID, 12)
	if err != nil {
		return nil, err
	}

	stats, err := s.repo.GetProductivityStats(userID)
	if err != nil {
		return nil, err
	}

	insights := []model.AIInsight{}

	// ピークタイム分析
	if insight := analyzePeakTime(heatmap); insight != nil {
		insights = append(insights, *insight)
	}

	// カテゴリ集中度分析
	if insight := analyzeCategoryFocus(categories); insight != nil {
		insights = append(insights, *insight)
	}

	// 週間成長分析
	if insight := analyzeWeeklyGrowth(trends); insight != nil {
		insights = append(insights, *insight)
	}

	// ストリーク分析
	if insight := analyzeStreak(stats); insight != nil {
		insights = append(insights, *insight)
	}

	return insights, nil
}

// analyzePeakTime はヒートマップから最も学習時間が長い曜日×時間帯を特定する。
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

// analyzeCategoryFocus はカテゴリ別データから学習の偏りを分析する。
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
	if percentage >= 70 {
		return &model.AIInsight{
			Type: "category_focus",
			Params: map[string]interface{}{
				"percentage": percentage,
				"category":   topCategory.Category,
			},
		}
	}

	return nil
}

// analyzeWeeklyGrowth は週間トレンドから成長傾向を分析する。
func analyzeWeeklyGrowth(trends []model.WeeklyTrend) *model.AIInsight {
	if len(trends) < 2 {
		return nil
	}

	// 直近2週間を比較
	prev := trends[len(trends)-2]
	current := trends[len(trends)-1]

	if prev.TotalMinutes == 0 {
		return nil
	}

	growthRate := math.Round(float64(current.TotalMinutes-prev.TotalMinutes) / float64(prev.TotalMinutes) * 100)

	if growthRate > 0 {
		return &model.AIInsight{
			Type: "weekly_growth_up",
			Params: map[string]interface{}{
				"rate": growthRate,
			},
		}
	} else if growthRate < -20 {
		return &model.AIInsight{
			Type: "weekly_growth_down",
			Params: map[string]interface{}{
				"rate": math.Abs(growthRate),
			},
		}
	}

	return nil
}

// analyzeStreak はストリーク情報からインサイトを生成する。
func analyzeStreak(stats *model.ProductivityStats) *model.AIInsight {
	if stats.CurrentStreak >= 7 {
		return &model.AIInsight{
			Type: "streak_active",
			Params: map[string]interface{}{
				"days": stats.CurrentStreak,
			},
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

package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// productivityRangeDays は生産性スコアの集計期間（過去 12 週）。
const productivityRangeDays = 84

// learningAnalyticsRepository は [repository.LearningAnalyticsRepository] の GORM 実装。
// learning_logs / learning_goals から Raw SQL で集計する。
type learningAnalyticsRepository struct {
	db *gorm.DB
}

// NewLearningAnalyticsRepository は LearningAnalyticsRepository の GORM 実装を返す。
func NewLearningAnalyticsRepository(db *gorm.DB) repository.LearningAnalyticsRepository {
	return &learningAnalyticsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.LearningAnalyticsRepository = (*learningAnalyticsRepository)(nil)

// GetHeatmapData は曜日×時間帯ごとの学習時間を返す。
func (r *learningAnalyticsRepository) GetHeatmapData(ctx context.Context, userID uint) ([]model.HeatmapEntry, error) {
	var entries []model.HeatmapEntry
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			EXTRACT(DOW FROM created_at)::int AS day_of_week,
			EXTRACT(HOUR FROM created_at)::int AS hour,
			COALESCE(SUM(duration), 0) AS total_minutes
		FROM learning_logs
		WHERE user_id = ?
		GROUP BY day_of_week, hour
		ORDER BY day_of_week, hour
	`, userID).Scan(&entries).Error
	if err != nil {
		return nil, err
	}
	return entries, nil
}

// GetCategoryBreakdown はカテゴリ別の学習時間とログ件数を返す。割合は usecase 側で計算する。
func (r *learningAnalyticsRepository) GetCategoryBreakdown(ctx context.Context, userID uint) ([]model.CategoryBreakdown, error) {
	var categories []model.CategoryBreakdown
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			category,
			COALESCE(SUM(duration), 0) AS total_minutes,
			COUNT(*) AS log_count
		FROM learning_logs
		WHERE user_id = ?
		GROUP BY category
		ORDER BY total_minutes DESC
	`, userID).Scan(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetWeeklyTrends は直近 weeks 週の週別集計を返す。
func (r *learningAnalyticsRepository) GetWeeklyTrends(ctx context.Context, userID uint, weeks int) ([]model.WeeklyTrend, error) {
	startDate := time.Now().AddDate(0, 0, -7*weeks).Format("2006-01-02")

	var trends []model.WeeklyTrend
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			TO_CHAR(DATE_TRUNC('week', created_at), 'YYYY-MM-DD') AS week_start,
			COALESCE(SUM(duration), 0) AS total_minutes,
			COUNT(*) AS log_count
		FROM learning_logs
		WHERE user_id = ? AND created_at >= ?
		GROUP BY week_start
		ORDER BY week_start
	`, userID, startDate).Scan(&trends).Error
	if err != nil {
		return nil, err
	}
	return trends, nil
}

// GetProductivityStats は生産性スコアの算出に必要な統計を集計する。
func (r *learningAnalyticsRepository) GetProductivityStats(ctx context.Context, userID uint) (*model.ProductivityStats, error) {
	db := r.db.WithContext(ctx)

	stats := &model.ProductivityStats{TotalDaysInRange: productivityRangeDays}
	startDate := time.Now().AddDate(0, 0, -productivityRangeDays).Format("2006-01-02")

	counts := []struct {
		query string
		args  []interface{}
		dest  *int
	}{
		{`SELECT COUNT(*) FROM learning_logs WHERE user_id = ? AND source = 'pomodoro' AND created_at >= ?`,
			[]interface{}{userID, startDate}, &stats.PomodoroSessions},
		{`SELECT COUNT(*) FROM learning_logs WHERE user_id = ? AND source = 'manual' AND created_at >= ?`,
			[]interface{}{userID, startDate}, &stats.ManualSessions},
		{`SELECT COUNT(*) FROM learning_goals WHERE user_id = ? AND status = 'completed'`,
			[]interface{}{userID}, &stats.CompletedGoals},
		{`SELECT COUNT(*) FROM learning_goals WHERE user_id = ?`,
			[]interface{}{userID}, &stats.TotalGoals},
		{`SELECT COUNT(DISTINCT DATE(created_at)) FROM learning_logs WHERE user_id = ? AND created_at >= ?`,
			[]interface{}{userID, startDate}, &stats.TotalLogDays},
	}
	for _, c := range counts {
		if err := db.Raw(c.query, c.args...).Scan(c.dest).Error; err != nil {
			return nil, err
		}
	}

	var rows []struct {
		LogDate time.Time
	}
	if err := db.Raw(`
		SELECT DISTINCT DATE(created_at) AS log_date
		FROM learning_logs
		WHERE user_id = ?
		ORDER BY log_date DESC
	`, userID).Scan(&rows).Error; err != nil {
		return nil, err
	}

	dates := make([]time.Time, 0, len(rows))
	for _, row := range rows {
		dates = append(dates, row.LogDate)
	}
	stats.CurrentStreak, stats.LongestStreak = calcAnalyticsStreaks(dates, time.Now())

	return stats, nil
}

// calcAnalyticsStreaks は降順の学習日一覧から、現在の連続日数と最長連続日数を求める。
// 直近の学習日が今日または昨日のときだけ現在の連続日数を数える。
func calcAnalyticsStreaks(dates []time.Time, now time.Time) (current, longest int) {
	if len(dates) == 0 {
		return 0, 0
	}

	today := normalizeToCalendarDay(now)
	first := normalizeToCalendarDay(dates[0])

	if isTodayOrYesterday(first, today) {
		current = 1
		prev := first
		for i := 1; i < len(dates); i++ {
			curr := normalizeToCalendarDay(dates[i])
			if !isNextCalendarDay(prev, curr) {
				break
			}
			current++
			prev = curr
		}
	}

	longest = 1
	streak := 1
	prev := normalizeToCalendarDay(dates[0])
	for i := 1; i < len(dates); i++ {
		curr := normalizeToCalendarDay(dates[i])
		if isNextCalendarDay(prev, curr) {
			streak++
		} else {
			if streak > longest {
				longest = streak
			}
			streak = 1
		}
		prev = curr
	}
	if streak > longest {
		longest = streak
	}

	return current, longest
}

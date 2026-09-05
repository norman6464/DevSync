package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// productivityRangeDays は生産性スコアの集計期間（過去 12 週）。
const productivityRangeDays = 84

// learningAnalyticsRepository は [repository.LearningAnalyticsRepository] の sqlc(pgx) 実装。
// learning_logs / learning_goals から集計する読み取り専用の操作。
type learningAnalyticsRepository struct {
	q *sqlcgen.Queries
}

// NewLearningAnalyticsRepository は LearningAnalyticsRepository の sqlc(pgx) 実装を返す。
func NewLearningAnalyticsRepository(q *sqlcgen.Queries) repository.LearningAnalyticsRepository {
	return &learningAnalyticsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.LearningAnalyticsRepository = (*learningAnalyticsRepository)(nil)

// startOfDay は指定日時の属するカレンダー日の午前0時を返す。
// 移行前のGo実装（`time.Format("2006-01-02")`をSQLへ文字列として渡す）と同じく、
// 時刻を切り捨てて「その日の始まり」だけを比較条件に使う。
func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// GetHeatmapData は曜日×時間帯ごとの学習時間を返す。
func (r *learningAnalyticsRepository) GetHeatmapData(ctx context.Context, userID uint) ([]model.HeatmapEntry, error) {
	rows, err := r.q.GetLearningHeatmapData(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	entries := make([]model.HeatmapEntry, len(rows))
	for i, row := range rows {
		entries[i] = model.HeatmapEntry{
			DayOfWeek:    int(row.DayOfWeek),
			Hour:         int(row.Hour),
			TotalMinutes: int(row.TotalMinutes),
		}
	}
	return entries, nil
}

// GetCategoryBreakdown はカテゴリ別の学習時間とログ件数を返す。割合は usecase 側で計算する。
func (r *learningAnalyticsRepository) GetCategoryBreakdown(ctx context.Context, userID uint) ([]model.CategoryBreakdown, error) {
	rows, err := r.q.GetLearningCategoryBreakdown(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	categories := make([]model.CategoryBreakdown, len(rows))
	for i, row := range rows {
		categories[i] = model.CategoryBreakdown{
			Category:     fromStringPtr(row.Category),
			TotalMinutes: int(row.TotalMinutes),
			LogCount:     int(row.LogCount),
		}
	}
	return categories, nil
}

// GetWeeklyTrends は直近 weeks 週の週別集計を返す。
func (r *learningAnalyticsRepository) GetWeeklyTrends(ctx context.Context, userID uint, weeks int) ([]model.WeeklyTrend, error) {
	startDate := startOfDay(time.Now().AddDate(0, 0, -7*weeks))

	rows, err := r.q.GetLearningWeeklyTrends(ctx, sqlcgen.GetLearningWeeklyTrendsParams{
		UserID:    int64(userID),
		CreatedAt: toTimestamptzNotNull(startDate),
	})
	if err != nil {
		return nil, err
	}
	trends := make([]model.WeeklyTrend, len(rows))
	for i, row := range rows {
		trends[i] = model.WeeklyTrend{
			WeekStart:    row.WeekStart,
			TotalMinutes: int(row.TotalMinutes),
			LogCount:     int(row.LogCount),
		}
	}
	return trends, nil
}

// GetProductivityStats は生産性スコアの算出に必要な統計を集計する。
func (r *learningAnalyticsRepository) GetProductivityStats(ctx context.Context, userID uint) (*model.ProductivityStats, error) {
	stats := &model.ProductivityStats{TotalDaysInRange: productivityRangeDays}
	startDate := toTimestamptzNotNull(startOfDay(time.Now().AddDate(0, 0, -productivityRangeDays)))

	pomodoro, err := r.q.CountLearningLogsBySourceSince(ctx, sqlcgen.CountLearningLogsBySourceSinceParams{
		UserID: int64(userID), Source: strPtr("pomodoro"), CreatedAt: startDate,
	})
	if err != nil {
		return nil, err
	}
	stats.PomodoroSessions = int(pomodoro)

	manual, err := r.q.CountLearningLogsBySourceSince(ctx, sqlcgen.CountLearningLogsBySourceSinceParams{
		UserID: int64(userID), Source: strPtr("manual"), CreatedAt: startDate,
	})
	if err != nil {
		return nil, err
	}
	stats.ManualSessions = int(manual)

	completedGoals, err := r.q.CountCompletedLearningGoalsByUser(ctx, sqlcgen.CountCompletedLearningGoalsByUserParams{
		UserID: int64(userID), Status: strPtr("completed"),
	})
	if err != nil {
		return nil, err
	}
	stats.CompletedGoals = int(completedGoals)

	totalGoals, err := r.q.CountLearningGoalsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	stats.TotalGoals = int(totalGoals)

	totalLogDays, err := r.q.CountLearningLogDaysSince(ctx, sqlcgen.CountLearningLogDaysSinceParams{
		UserID: int64(userID), CreatedAt: startDate,
	})
	if err != nil {
		return nil, err
	}
	stats.TotalLogDays = int(totalLogDays)

	dateRows, err := r.q.ListDistinctLearningLogDates(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	dates := make([]time.Time, len(dateRows))
	for i, d := range dateRows {
		dates[i] = d.Time
	}
	stats.CurrentStreak, stats.LongestStreak = calcAnalyticsStreaks(dates, time.Now())

	return stats, nil
}

// strPtr は文字列リテラルを固定値フィルタとして渡すための小さなヘルパー。
func strPtr(s string) *string {
	return &s
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

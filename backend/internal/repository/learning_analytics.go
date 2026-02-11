package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// LearningAnalyticsRepository は学習分析データの集計を提供するリポジトリ実装。
// learning_logsテーブルと関連テーブルからRaw SQLで統計を集計する。
type LearningAnalyticsRepository struct {
	db *gorm.DB
}

// NewLearningAnalyticsRepository は新しいLearningAnalyticsRepositoryインスタンスを生成する。
func NewLearningAnalyticsRepository(db *gorm.DB) *LearningAnalyticsRepository {
	return &LearningAnalyticsRepository{db: db}
}

// GetHeatmapData は曜日×時間帯ごとの学習時間集計を返す。
func (r *LearningAnalyticsRepository) GetHeatmapData(userID uint) ([]model.HeatmapEntry, error) {
	var entries []model.HeatmapEntry
	err := r.db.Raw(`
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

// GetCategoryBreakdown はカテゴリ別の学習時間・ログ件数を集計する。
// Percentageはサービス層で計算するため、ここでは0を返す。
func (r *LearningAnalyticsRepository) GetCategoryBreakdown(userID uint) ([]model.CategoryBreakdown, error) {
	var categories []model.CategoryBreakdown
	err := r.db.Raw(`
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

// GetWeeklyTrends は過去N週間の週ごとの学習時間・ログ件数を集計する。
func (r *LearningAnalyticsRepository) GetWeeklyTrends(userID uint, weeks int) ([]model.WeeklyTrend, error) {
	startDate := time.Now().AddDate(0, 0, -7*weeks).Format("2006-01-02")

	var trends []model.WeeklyTrend
	err := r.db.Raw(`
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

// GetProductivityStats は生産性スコア計算に必要な統計データを集計する。
func (r *LearningAnalyticsRepository) GetProductivityStats(userID uint) (*model.ProductivityStats, error) {
	stats := &model.ProductivityStats{}

	// 過去12週間の期間
	daysInRange := 84
	stats.TotalDaysInRange = daysInRange
	startDate := time.Now().AddDate(0, 0, -daysInRange).Format("2006-01-02")

	// ポモドーロセッション数
	if err := r.db.Raw(
		`SELECT COUNT(*) FROM learning_logs WHERE user_id = ? AND source = 'pomodoro' AND created_at >= ?`,
		userID, startDate,
	).Scan(&stats.PomodoroSessions).Error; err != nil {
		return nil, err
	}

	// 手動記録セッション数
	if err := r.db.Raw(
		`SELECT COUNT(*) FROM learning_logs WHERE user_id = ? AND source = 'manual' AND created_at >= ?`,
		userID, startDate,
	).Scan(&stats.ManualSessions).Error; err != nil {
		return nil, err
	}

	// 完了した目標数
	if err := r.db.Raw(
		`SELECT COUNT(*) FROM learning_goals WHERE user_id = ? AND status = 'completed'`,
		userID,
	).Scan(&stats.CompletedGoals).Error; err != nil {
		return nil, err
	}

	// 全目標数
	if err := r.db.Raw(
		`SELECT COUNT(*) FROM learning_goals WHERE user_id = ?`,
		userID,
	).Scan(&stats.TotalGoals).Error; err != nil {
		return nil, err
	}

	// 学習記録がある日数（過去12週）
	if err := r.db.Raw(
		`SELECT COUNT(DISTINCT DATE(created_at)) FROM learning_logs WHERE user_id = ? AND created_at >= ?`,
		userID, startDate,
	).Scan(&stats.TotalLogDays).Error; err != nil {
		return nil, err
	}

	// ストリーク情報：学習記録日を取得して連続日数を算出
	type logDate struct {
		LogDate time.Time
	}
	var logDates []logDate
	if err := r.db.Raw(`
		SELECT DISTINCT DATE(created_at) AS log_date
		FROM learning_logs
		WHERE user_id = ?
		ORDER BY log_date DESC
	`, userID).Scan(&logDates).Error; err != nil {
		return nil, err
	}

	if len(logDates) > 0 {
		normalize := func(t time.Time) time.Time {
			return t.Truncate(24 * time.Hour)
		}

		today := normalize(time.Now())
		first := normalize(logDates[0].LogDate)

		// current streak: 今日から遡って連続している日数
		if first.Equal(today) || first.Equal(today.AddDate(0, 0, -1)) {
			stats.CurrentStreak = 1
			prev := first
			for i := 1; i < len(logDates); i++ {
				curr := normalize(logDates[i].LogDate)
				if int(prev.Sub(curr).Hours()/24) == 1 {
					stats.CurrentStreak++
					prev = curr
				} else {
					break
				}
			}
		}

		// longest streak: 全期間での最大連続日数
		longestStreak := 1
		streakLen := 1
		prev := normalize(logDates[0].LogDate)
		for i := 1; i < len(logDates); i++ {
			curr := normalize(logDates[i].LogDate)
			if int(prev.Sub(curr).Hours()/24) == 1 {
				streakLen++
			} else {
				if streakLen > longestStreak {
					longestStreak = streakLen
				}
				streakLen = 1
			}
			prev = curr
		}
		if streakLen > longestStreak {
			longestStreak = streakLen
		}
		stats.LongestStreak = longestStreak
	}

	return stats, nil
}

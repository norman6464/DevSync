package persistence

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// learningLogRepository は [repository.LearningLogRepository] の GORM 実装。
type learningLogRepository struct {
	db *gorm.DB
}

// NewLearningLogRepository は LearningLogRepository の GORM 実装を返す。
func NewLearningLogRepository(db *gorm.DB) repository.LearningLogRepository {
	return &learningLogRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.LearningLogRepository = (*learningLogRepository)(nil)

// Create は新しい学習ログを作成する。
func (r *learningLogRepository) Create(ctx context.Context, log *model.LearningLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// CreateBatch は複数の学習ログを一括作成する。
func (r *learningLogRepository) CreateBatch(ctx context.Context, logs []model.LearningLog) error {
	return r.db.WithContext(ctx).Create(&logs).Error
}

// Update は既存の学習ログを更新する。
func (r *learningLogRepository) Update(ctx context.Context, log *model.LearningLog) error {
	return r.db.WithContext(ctx).Save(log).Error
}

// Delete は所有者本人の学習ログを削除する。
func (r *learningLogRepository) Delete(ctx context.Context, id, userID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.LearningLog{}).Error
}

// FindByID は指定 ID の学習ログを取得する。不在の場合は (nil, nil) を返す。
func (r *learningLogRepository) FindByID(ctx context.Context, id uint) (*model.LearningLog, error) {
	var log model.LearningLog
	err := r.db.WithContext(ctx).First(&log, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetByUserID はユーザーの学習ログを作成日の新しい順でページ取得し、総数も返す。
func (r *learningLogRepository) GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	return r.pagedLogs(ctx, "user_id = ?", []interface{}{userID}, limit, offset)
}

// GetFavorites はお気に入りの学習ログを作成日の新しい順でページ取得し、総数も返す。
func (r *learningLogRepository) GetFavorites(ctx context.Context, userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	return r.pagedLogs(ctx, "user_id = ? AND is_favorite = ?", []interface{}{userID, true}, limit, offset)
}

// GetByGoalID は指定ゴールに紐付いた学習ログを作成日の新しい順でページ取得し、総数も返す。
func (r *learningLogRepository) GetByGoalID(ctx context.Context, goalID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	return r.pagedLogs(ctx, "goal_id = ?", []interface{}{goalID}, limit, offset)
}

// pagedLogs は絞り込み条件を受け取り、総数とページを返す共通処理。
func (r *learningLogRepository) pagedLogs(ctx context.Context, where string, args []interface{}, limit, offset int) ([]model.LearningLog, int64, error) {
	scope := r.db.WithContext(ctx).Model(&model.LearningLog{}).Where(where, args...)

	var total int64
	if err := scope.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var logs []model.LearningLog
	err := scope.Session(&gorm.Session{}).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&logs).Error
	return logs, total, err
}

// GetByCategory はカテゴリで絞り込んだ学習ログを作成日の新しい順で取得する。
func (r *learningLogRepository) GetByCategory(ctx context.Context, userID uint, category string) ([]model.LearningLog, error) {
	var logs []model.LearningLog
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND category = ?", userID, category).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// GetBySource はソースで絞り込んだ学習ログを作成日の新しい順で取得する。
func (r *learningLogRepository) GetBySource(ctx context.Context, userID uint, source string) ([]model.LearningLog, error) {
	var logs []model.LearningLog
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND source = ?", userID, source).
		Order("created_at DESC").
		Find(&logs).Error
	return logs, err
}

// GetByPeriod は直近 days 日の学習ログを作成日の新しい順で取得する。days が 0 以下なら全期間。
func (r *learningLogRepository) GetByPeriod(ctx context.Context, userID uint, days int) ([]model.LearningLog, error) {
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if days > 0 {
		query = query.Where("created_at >= ?", time.Now().AddDate(0, 0, -days))
	}

	var logs []model.LearningLog
	err := query.Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// SumDurationByPeriod は直近 days 日の学習時間合計（分）を返す。
func (r *learningLogRepository) SumDurationByPeriod(ctx context.Context, userID uint, days int) (int, error) {
	var total int
	err := r.db.WithContext(ctx).Model(&model.LearningLog{}).
		Where("user_id = ? AND created_at >= ?", userID, time.Now().AddDate(0, 0, -days)).
		Select("COALESCE(SUM(duration), 0)").
		Scan(&total).Error
	return total, err
}

// SumDurationByGoalID は指定ゴールに紐付いた学習ログの学習時間合計（分）を返す。
func (r *learningLogRepository) SumDurationByGoalID(ctx context.Context, goalID uint) (int, error) {
	var total int
	err := r.db.WithContext(ctx).Model(&model.LearningLog{}).
		Where("goal_id = ?", goalID).
		Select("COALESCE(SUM(duration), 0)").
		Scan(&total).Error
	return total, err
}

// CountByUserID はユーザーの学習ログ総数を返す。
func (r *learningLogRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.LearningLog{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

// GetStreakInfo は学習ログの日付から連続学習情報を算出する。
func (r *learningLogRepository) GetStreakInfo(ctx context.Context, userID uint) (*model.StreakInfo, error) {
	var rows []struct {
		Date time.Time
	}
	err := r.db.WithContext(ctx).
		Raw("SELECT DISTINCT DATE(created_at) as date FROM learning_logs WHERE user_id = ? ORDER BY date DESC", userID).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	dates := make([]time.Time, 0, len(rows))
	for _, row := range rows {
		dates = append(dates, row.Date)
	}
	return calcStreakInfo(dates, time.Now()), nil
}

// calcStreakInfo は学習日の一覧から連続学習情報を組み立てる。
// 直近のログが今日または昨日のときだけ現在の連続日数を数える。
func calcStreakInfo(dates []time.Time, now time.Time) *model.StreakInfo {
	info := &model.StreakInfo{TotalDays: len(dates)}
	if len(dates) == 0 {
		return info
	}

	sort.Slice(dates, func(i, j int) bool { return dates[i].After(dates[j]) })
	info.LastLogDate = dates[0].Format("2006-01-02")

	if isTodayOrYesterday(normalizeToCalendarDay(dates[0]), normalizeToCalendarDay(now)) {
		info.CurrentStreak = 1
		for i := 1; i < len(dates); i++ {
			// 同日の重複は「翌暦日」にならないため連続とみなさない。
			if !isNextCalendarDay(normalizeToCalendarDay(dates[i-1]), normalizeToCalendarDay(dates[i])) {
				break
			}
			info.CurrentStreak++
		}
	}

	longest, streak := 1, 1
	for i := 1; i < len(dates); i++ {
		if isNextCalendarDay(normalizeToCalendarDay(dates[i-1]), normalizeToCalendarDay(dates[i])) {
			streak++
			if streak > longest {
				longest = streak
			}
			continue
		}
		streak = 1
	}
	info.LongestStreak = longest

	return info
}

// GetRecentCategories は使用回数の多い順にカテゴリを limit 件返す。
func (r *learningLogRepository) GetRecentCategories(ctx context.Context, userID uint, limit int) ([]string, error) {
	var categories []string
	err := r.db.WithContext(ctx).Model(&model.LearningLog{}).
		Select("category").
		Where("user_id = ?", userID).
		Group("category").
		Order("COUNT(*) DESC").
		Limit(limit).
		Pluck("category", &categories).Error
	return categories, err
}

// GetCalendarData はカレンダー表示用の日別ログ件数を取得する。
func (r *learningLogRepository) GetCalendarData(ctx context.Context, userID uint) ([]model.CalendarEntry, error) {
	var entries []model.CalendarEntry
	err := r.db.WithContext(ctx).Model(&model.LearningLog{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("DATE(created_at)").
		Order("date ASC").
		Find(&entries).Error
	return entries, err
}

// GetMonthlySummary は直近 months ヶ月の月別サマリー（合計時間・ログ件数）を取得する。
func (r *learningLogRepository) GetMonthlySummary(ctx context.Context, userID uint, months int) ([]model.MonthlySummary, error) {
	startDate := time.Now().AddDate(0, -months, 0).Format("2006-01-02")

	var summaries []model.MonthlySummary
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			TO_CHAR(DATE_TRUNC('month', created_at), 'YYYY-MM-DD') AS month,
			COALESCE(SUM(duration), 0) AS total_minutes,
			COUNT(*) AS log_count
		FROM learning_logs
		WHERE user_id = ? AND created_at >= ?
		GROUP BY month
		ORDER BY month
	`, userID, startDate).Scan(&summaries).Error
	return summaries, err
}

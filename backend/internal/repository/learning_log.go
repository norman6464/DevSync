package repository

import (
	"sort"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// LearningLogRepository は学習ログデータへのアクセスを提供するリポジトリ実装。
type LearningLogRepository struct {
	db *gorm.DB
}

// NewLearningLogRepository は新しいLearningLogRepositoryインスタンスを生成する。
func NewLearningLogRepository(db *gorm.DB) *LearningLogRepository {
	return &LearningLogRepository{db: db}
}

// Create は新しい学習ログをデータベースに作成する。
func (r *LearningLogRepository) Create(log *model.LearningLog) error {
	return r.db.Create(log).Error
}

// CreateBatch は複数の学習ログをトランザクション内で一括作成する。
func (r *LearningLogRepository) CreateBatch(logs []model.LearningLog) error {
	return r.db.Create(&logs).Error
}

// Update は既存の学習ログを更新する。
func (r *LearningLogRepository) Update(log *model.LearningLog) error {
	return r.db.Save(log).Error
}

// Delete は指定IDかつ指定ユーザーの学習ログを削除する（所有権チェック付き）。
func (r *LearningLogRepository) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.LearningLog{}).Error
}

// FindByID は指定IDの学習ログを取得する。
func (r *LearningLogRepository) FindByID(id uint) (*model.LearningLog, error) {
	var log model.LearningLog
	err := r.db.First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetByUserID は指定ユーザーの学習ログをページネーション付きで取得する（新しい順）。
func (r *LearningLogRepository) GetByUserID(userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	var logs []model.LearningLog
	var total int64
	query := r.db.Where("user_id = ?", userID)
	query.Model(&model.LearningLog{}).Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

// GetByCategory は指定ユーザーの学習ログをカテゴリでフィルタリングして取得する（新しい順）。
func (r *LearningLogRepository) GetByCategory(userID uint, category string) ([]model.LearningLog, error) {
	var logs []model.LearningLog
	err := r.db.Where("user_id = ? AND category = ?", userID, category).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// GetBySource は指定ユーザーの学習ログをソースでフィルタリングして取得する（新しい順）。
func (r *LearningLogRepository) GetBySource(userID uint, source string) ([]model.LearningLog, error) {
	var logs []model.LearningLog
	err := r.db.Where("user_id = ? AND source = ?", userID, source).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// GetByPeriod は指定ユーザーの指定期間分の学習ログを取得する。
// days=0 の場合は全期間を取得する。
func (r *LearningLogRepository) GetByPeriod(userID uint, days int) ([]model.LearningLog, error) {
	var logs []model.LearningLog
	query := r.db.Where("user_id = ?", userID)
	if days > 0 {
		since := time.Now().AddDate(0, 0, -days)
		query = query.Where("created_at >= ?", since)
	}
	err := query.Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// SumDurationByPeriod は指定ユーザーの指定期間内の学習時間合計（分）を返す。
func (r *LearningLogRepository) SumDurationByPeriod(userID uint, days int) (int, error) {
	var total int
	since := time.Now().AddDate(0, 0, -days)
	err := r.db.Model(&model.LearningLog{}).
		Where("user_id = ? AND created_at >= ?", userID, since).
		Select("COALESCE(SUM(duration), 0)").
		Scan(&total).Error
	return total, err
}

// GetStreakInfo は学習ログから連続学習情報を算出する。
// 現在の連続日数、最長連続日数、合計学習日数、最終ログ日を返す。
func (r *LearningLogRepository) GetStreakInfo(userID uint) (*model.StreakInfo, error) {
	var dates []struct {
		Date time.Time
	}
	err := r.db.Raw("SELECT DISTINCT DATE(created_at) as date FROM learning_logs WHERE user_id = ? ORDER BY date DESC", userID).Scan(&dates).Error
	if err != nil {
		return nil, err
	}

	info := &model.StreakInfo{
		TotalDays: len(dates),
	}

	if len(dates) == 0 {
		return info, nil
	}

	info.LastLogDate = dates[0].Date.Format("2006-01-02")

	// 日付を降順にソート
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Date.After(dates[j].Date)
	})

	// 現在の連続日数を計算
	today := time.Now().UTC().Truncate(24 * time.Hour)
	currentStreak := 0
	startIdx := 0

	firstDate := dates[0].Date.UTC().Truncate(24 * time.Hour)
	diffToToday := today.Sub(firstDate)

	// 今日または昨日にログがある場合のみ連続日数をカウント
	if diffToToday < 48*time.Hour {
		currentStreak = 1
		startIdx = 1
		for i := startIdx; i < len(dates); i++ {
			prev := dates[i-1].Date.UTC().Truncate(24 * time.Hour)
			curr := dates[i].Date.UTC().Truncate(24 * time.Hour)
			diff := prev.Sub(curr)
			if diff >= 24*time.Hour && diff < 48*time.Hour {
				currentStreak++
			} else {
				break
			}
		}
	}

	info.CurrentStreak = currentStreak

	// 最長連続日数を計算
	longest := 1
	streak := 1
	for i := 1; i < len(dates); i++ {
		prev := dates[i-1].Date.UTC().Truncate(24 * time.Hour)
		curr := dates[i].Date.UTC().Truncate(24 * time.Hour)
		diff := prev.Sub(curr)
		if diff >= 24*time.Hour && diff < 48*time.Hour {
			streak++
			if streak > longest {
				longest = streak
			}
		} else {
			streak = 1
		}
	}
	info.LongestStreak = longest

	return info, nil
}

// GetRecentCategories はユーザーの最近よく使うカテゴリを頻度順で返す。
func (r *LearningLogRepository) GetRecentCategories(userID uint, limit int) ([]string, error) {
	var categories []string
	err := r.db.Model(&model.LearningLog{}).
		Select("category").
		Where("user_id = ?", userID).
		Group("category").
		Order("COUNT(*) DESC").
		Limit(limit).
		Pluck("category", &categories).Error
	return categories, err
}

// GetByGoalID は指定ゴールに紐付いた学習ログをページネーション付きで取得する。
func (r *LearningLogRepository) GetByGoalID(goalID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	var logs []model.LearningLog
	var total int64
	query := r.db.Where("goal_id = ?", goalID)
	query.Model(&model.LearningLog{}).Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

// SumDurationByGoalID は指定ゴールに紐付いた学習ログの合計学習時間（分）を返す。
func (r *LearningLogRepository) SumDurationByGoalID(goalID uint) (int, error) {
	var total int
	err := r.db.Model(&model.LearningLog{}).
		Where("goal_id = ?", goalID).
		Select("COALESCE(SUM(duration), 0)").
		Scan(&total).Error
	return total, err
}

// GetFavorites はお気に入り学習ログをページネーション付きで取得する（新しい順）。
func (r *LearningLogRepository) GetFavorites(userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	var logs []model.LearningLog
	var total int64
	query := r.db.Where("user_id = ? AND is_favorite = ?", userID, true)
	query.Model(&model.LearningLog{}).Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

// GetCalendarData はカレンダー表示用の日別ログ件数を取得する。
func (r *LearningLogRepository) GetCalendarData(userID uint) ([]model.CalendarEntry, error) {
	var entries []model.CalendarEntry
	err := r.db.Model(&model.LearningLog{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("DATE(created_at)").
		Order("date ASC").
		Find(&entries).Error
	return entries, err
}

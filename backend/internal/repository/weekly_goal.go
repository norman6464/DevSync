package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// WeeklyGoalRepository はカテゴリ別週間学習目標のリポジトリ実装。
type WeeklyGoalRepository struct {
	db *gorm.DB
}

// NewWeeklyGoalRepository は新しいWeeklyGoalRepositoryインスタンスを生成する。
func NewWeeklyGoalRepository(db *gorm.DB) *WeeklyGoalRepository {
	return &WeeklyGoalRepository{db: db}
}

// Upsert はWeeklyGoalを作成または更新する（user_id+categoryでユニーク）。
func (r *WeeklyGoalRepository) Upsert(goal *model.WeeklyGoal) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "category"}},
		DoUpdates: clause.AssignmentColumns([]string{"target_minutes", "updated_at"}),
	}).Create(goal).Error
}

// GetByUserID は指定ユーザーの全カテゴリ週間目標を取得する。
func (r *WeeklyGoalRepository) GetByUserID(userID uint) ([]model.WeeklyGoal, error) {
	var goals []model.WeeklyGoal
	err := r.db.Where("user_id = ?", userID).Order("category ASC").Find(&goals).Error
	return goals, err
}

// SumDurationByUserCategoryThisWeek は指定ユーザー・カテゴリの今週の学習時間合計（分）を返す。
func (r *WeeklyGoalRepository) SumDurationByUserCategoryThisWeek(userID uint, category string) (int, error) {
	var total int
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // 日曜を7にして月曜始まり
	}
	startOfWeek := now.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)

	err := r.db.Model(&model.LearningLog{}).
		Where("user_id = ? AND category = ? AND created_at >= ?", userID, category, startOfWeek).
		Select("COALESCE(SUM(duration), 0)").
		Scan(&total).Error
	return total, err
}

package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// weeklyGoalRepository は [repository.WeeklyGoalRepository] の GORM 実装。
type weeklyGoalRepository struct {
	db *gorm.DB
}

// NewWeeklyGoalRepository は WeeklyGoalRepository の GORM 実装を返す。
func NewWeeklyGoalRepository(db *gorm.DB) repository.WeeklyGoalRepository {
	return &weeklyGoalRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.WeeklyGoalRepository = (*weeklyGoalRepository)(nil)

// Upsert は WeeklyGoal を作成または更新する（user_id + category でユニーク）。
func (r *weeklyGoalRepository) Upsert(ctx context.Context, goal *model.WeeklyGoal) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "category"}},
		DoUpdates: clause.AssignmentColumns([]string{"target_minutes", "updated_at"}),
	}).Create(goal).Error
}

// GetByUserID は指定ユーザーの全カテゴリ週間目標を取得する。
func (r *weeklyGoalRepository) GetByUserID(ctx context.Context, userID uint) ([]model.WeeklyGoal, error) {
	var goals []model.WeeklyGoal
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("category ASC").Find(&goals).Error
	return goals, err
}

// SumDurationByUserCategoryThisWeek は指定ユーザー・カテゴリの今週の学習時間合計（分）を返す。
func (r *weeklyGoalRepository) SumDurationByUserCategoryThisWeek(ctx context.Context, userID uint, category string) (int, error) {
	var total int
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // 日曜を 7 にして月曜始まり
	}
	startOfWeek := now.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)

	err := r.db.WithContext(ctx).Model(&model.LearningLog{}).
		Where("user_id = ? AND category = ? AND created_at >= ?", userID, category, startOfWeek).
		Select("COALESCE(SUM(duration), 0)").
		Scan(&total).Error
	return total, err
}

package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// weeklyGoalRepository は [repository.WeeklyGoalRepository] の sqlc(pgx) 実装。
type weeklyGoalRepository struct {
	q *sqlcgen.Queries
}

// NewWeeklyGoalRepository は WeeklyGoalRepository の sqlc(pgx) 実装を返す。
func NewWeeklyGoalRepository(q *sqlcgen.Queries) repository.WeeklyGoalRepository {
	return &weeklyGoalRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.WeeklyGoalRepository = (*weeklyGoalRepository)(nil)

// toModelWeeklyGoal は sqlc の生成行を model.WeeklyGoal へ変換する。
func toModelWeeklyGoal(row sqlcgen.WeeklyGoal) model.WeeklyGoal {
	return model.WeeklyGoal{
		ID:            uint(row.ID),
		UserID:        uint(row.UserID),
		Category:      model.LogCategory(row.Category),
		TargetMinutes: int(row.TargetMinutes),
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}
}

// Upsert は WeeklyGoal を作成または更新する（user_id + category でユニーク）。
func (r *weeklyGoalRepository) Upsert(ctx context.Context, goal *model.WeeklyGoal) error {
	row, err := r.q.UpsertWeeklyGoal(ctx, sqlcgen.UpsertWeeklyGoalParams{
		UserID:        int64(goal.UserID),
		Category:      string(goal.Category),
		TargetMinutes: int64(goal.TargetMinutes),
	})
	if err != nil {
		return err
	}
	*goal = toModelWeeklyGoal(row)
	return nil
}

// GetByUserID は指定ユーザーの全カテゴリ週間目標を取得する。
func (r *weeklyGoalRepository) GetByUserID(ctx context.Context, userID uint) ([]model.WeeklyGoal, error) {
	rows, err := r.q.ListWeeklyGoalsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	goals := make([]model.WeeklyGoal, len(rows))
	for i, row := range rows {
		goals[i] = toModelWeeklyGoal(row)
	}
	return goals, nil
}

// SumDurationByUserCategoryThisWeek は指定ユーザー・カテゴリの今週の学習時間合計（分）を返す。
func (r *weeklyGoalRepository) SumDurationByUserCategoryThisWeek(ctx context.Context, userID uint, category string) (int, error) {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // 日曜を 7 にして月曜始まり
	}
	startOfWeek := now.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)

	total, err := r.q.SumLearningLogDurationByUserCategorySince(ctx, sqlcgen.SumLearningLogDurationByUserCategorySinceParams{
		UserID:    int64(userID),
		Category:  &category,
		CreatedAt: toTimestamptzNotNull(startOfWeek),
	})
	if err != nil {
		return 0, err
	}
	return int(total), nil
}

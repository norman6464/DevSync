package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// learningGoalRepository は [repository.LearningGoalRepository] の sqlc(pgx) 実装。
type learningGoalRepository struct {
	q *sqlcgen.Queries
}

// NewLearningGoalRepository は LearningGoalRepository の sqlc(pgx) 実装を返す。
func NewLearningGoalRepository(q *sqlcgen.Queries) repository.LearningGoalRepository {
	return &learningGoalRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.LearningGoalRepository = (*learningGoalRepository)(nil)

func toModelLearningGoal(row sqlcgen.LearningGoal) model.LearningGoal {
	return model.LearningGoal{
		ID:          uint(row.ID),
		UserID:      uint(row.UserID),
		Title:       row.Title,
		Description: fromStringPtr(row.Description),
		Category:    model.GoalCategory(fromStringPtr(row.Category)),
		TargetDate:  fromTimestamptz(row.TargetDate),
		Progress:    int(fromInt64PtrValue(row.Progress)),
		TargetHours: int(fromInt64PtrValue(row.TargetHours)),
		Status:      model.GoalStatus(fromStringPtr(row.Status)),
		IsPublic:    fromBoolPtr(row.IsPublic),
		CreatedAt:   timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:   timeValue(fromTimestamptz(row.UpdatedAt)),
		CompletedAt: fromTimestamptz(row.CompletedAt),
	}
}

// Create は新しい学習目標を作成する。
// Status はGORMの `gorm:"default:'active'"` に相当し、未指定（ゼロ値）なら active を補う。
func (r *learningGoalRepository) Create(ctx context.Context, goal *model.LearningGoal) error {
	status := goal.Status
	if status == "" {
		status = model.GoalStatusActive
	}

	row, err := r.q.CreateLearningGoal(ctx, sqlcgen.CreateLearningGoalParams{
		UserID:      int64(goal.UserID),
		Title:       goal.Title,
		Description: &goal.Description,
		Category:    (*string)(&goal.Category),
		TargetDate:  toTimestamptz(goal.TargetDate),
		Progress:    toInt64Ptr(goal.Progress),
		TargetHours: toInt64Ptr(goal.TargetHours),
		Status:      (*string)(&status),
		IsPublic:    &goal.IsPublic,
	})
	if err != nil {
		return err
	}
	*goal = toModelLearningGoal(row)
	return nil
}

// Update は既存の学習目標を更新する（GORMのSave＝全カラム上書きに相当）。
func (r *learningGoalRepository) Update(ctx context.Context, goal *model.LearningGoal) error {
	row, err := r.q.UpdateLearningGoal(ctx, sqlcgen.UpdateLearningGoalParams{
		ID:          int64(goal.ID),
		Title:       goal.Title,
		Description: &goal.Description,
		Category:    (*string)(&goal.Category),
		TargetDate:  toTimestamptz(goal.TargetDate),
		Progress:    toInt64Ptr(goal.Progress),
		TargetHours: toInt64Ptr(goal.TargetHours),
		Status:      (*string)(&goal.Status),
		IsPublic:    &goal.IsPublic,
		CompletedAt: toTimestamptz(goal.CompletedAt),
	})
	if err != nil {
		return err
	}
	*goal = toModelLearningGoal(row)
	return nil
}

// Delete は学習目標を削除する。
func (r *learningGoalRepository) Delete(ctx context.Context, id uint) error {
	return r.q.DeleteLearningGoal(ctx, int64(id))
}

// FindByID は指定 ID の学習目標を取得する。不在の場合は (nil, nil) を返す。
func (r *learningGoalRepository) FindByID(ctx context.Context, id uint) (*model.LearningGoal, error) {
	row, err := r.q.GetLearningGoalByID(ctx, int64(id))
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	goal := toModelLearningGoal(row)
	return &goal, nil
}

// GetByUserID は指定ユーザーの学習目標をページネーション付きで取得する（新しい順）。
func (r *learningGoalRepository) GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.LearningGoal, int64, error) {
	total, err := r.q.CountLearningGoalsByUser(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListLearningGoalsByUser(ctx, sqlcgen.ListLearningGoalsByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	return toModelLearningGoals(rows), total, nil
}

func toModelLearningGoals(rows []sqlcgen.LearningGoal) []model.LearningGoal {
	goals := make([]model.LearningGoal, len(rows))
	for i, row := range rows {
		goals[i] = toModelLearningGoal(row)
	}
	return goals
}

// GetActiveByUserID は指定ユーザーの進行中の学習目標を取得する（新しい順）。
func (r *learningGoalRepository) GetActiveByUserID(ctx context.Context, userID uint) ([]model.LearningGoal, error) {
	rows, err := r.q.ListActiveLearningGoalsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	return toModelLearningGoals(rows), nil
}

// GetByCategory は指定ユーザーの学習目標をカテゴリで絞り込んで取得する（新しい順）。
func (r *learningGoalRepository) GetByCategory(ctx context.Context, userID uint, category string) ([]model.LearningGoal, error) {
	rows, err := r.q.ListLearningGoalsByCategory(ctx, sqlcgen.ListLearningGoalsByCategoryParams{
		UserID:   int64(userID),
		Category: &category,
	})
	if err != nil {
		return nil, err
	}
	return toModelLearningGoals(rows), nil
}

// GetByStatus は指定ユーザーの学習目標をステータスで絞り込んで取得する（新しい順）。
func (r *learningGoalRepository) GetByStatus(ctx context.Context, userID uint, status string) ([]model.LearningGoal, error) {
	rows, err := r.q.ListLearningGoalsByStatus(ctx, sqlcgen.ListLearningGoalsByStatusParams{
		UserID: int64(userID),
		Status: &status,
	})
	if err != nil {
		return nil, err
	}
	return toModelLearningGoals(rows), nil
}

// GetPublicByUserID は指定ユーザーの公開済み学習目標をページネーション付きで取得する（新しい順）。
func (r *learningGoalRepository) GetPublicByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.LearningGoal, int64, error) {
	total, err := r.q.CountPublicLearningGoalsByUser(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListPublicLearningGoalsByUser(ctx, sqlcgen.ListPublicLearningGoalsByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	return toModelLearningGoals(rows), total, nil
}

// GetPublicGoals は全ユーザーの公開済み学習目標をページネーション付きで取得する（新しい順）。
func (r *learningGoalRepository) GetPublicGoals(ctx context.Context, limit, offset int) ([]model.LearningGoal, int64, error) {
	total, err := r.q.CountPublicLearningGoals(ctx)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListPublicLearningGoals(ctx, sqlcgen.ListPublicLearningGoalsParams{
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	return toModelLearningGoals(rows), total, nil
}

// CountByUserID は指定ユーザーの学習目標総数を返す。
func (r *learningGoalRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountLearningGoalsByUser(ctx, int64(userID))
}

// GetStats は指定ユーザーの学習目標統計（総数・アクティブ数・完了数・平均進捗）を返す。
func (r *learningGoalRepository) GetStats(ctx context.Context, userID uint) (*model.LearningGoalStats, error) {
	total, err := r.q.CountLearningGoalsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	active, err := r.q.CountActiveLearningGoalsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	completedStatus := string(model.GoalStatusCompleted)
	completed, err := r.q.CountCompletedLearningGoalsByUser(ctx, sqlcgen.CountCompletedLearningGoalsByUserParams{
		UserID: int64(userID),
		Status: &completedStatus,
	})
	if err != nil {
		return nil, err
	}

	avgProgress, err := r.q.GetAverageActiveProgressByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	return &model.LearningGoalStats{
		TotalGoals:      int(total),
		ActiveGoals:     int(active),
		CompletedGoals:  int(completed),
		AverageProgress: int(avgProgress),
	}, nil
}

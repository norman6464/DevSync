package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// learningGoalRepository は [repository.LearningGoalRepository] の GORM 実装。
//
// 旧 repository パッケージにも同じテーブルを扱う実装が残っている。learning_log /
// recommendation / learning_dashboard がまだそちらを使っているため、移行が一巡するまで
// 新旧のアダプタが並存する。
type learningGoalRepository struct {
	db *gorm.DB
}

// NewLearningGoalRepository は LearningGoalRepository の GORM 実装を返す。
func NewLearningGoalRepository(db *gorm.DB) repository.LearningGoalRepository {
	return &learningGoalRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.LearningGoalRepository = (*learningGoalRepository)(nil)

// Create は新しい学習目標を作成する。
func (r *learningGoalRepository) Create(ctx context.Context, goal *model.LearningGoal) error {
	return r.db.WithContext(ctx).Create(goal).Error
}

// Update は既存の学習目標を更新する。
func (r *learningGoalRepository) Update(ctx context.Context, goal *model.LearningGoal) error {
	return r.db.WithContext(ctx).Save(goal).Error
}

// Delete は学習目標を削除する。
func (r *learningGoalRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.LearningGoal{}, id).Error
}

// FindByID は指定 ID の学習目標を取得する。不在の場合は (nil, nil) を返す。
func (r *learningGoalRepository) FindByID(ctx context.Context, id uint) (*model.LearningGoal, error) {
	var goal model.LearningGoal
	err := r.db.WithContext(ctx).First(&goal, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

// GetByUserID は指定ユーザーの学習目標をページネーション付きで取得する（新しい順）。
func (r *learningGoalRepository) GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.LearningGoal, int64, error) {
	return r.paginated(ctx, r.db.WithContext(ctx).Where("user_id = ?", userID), limit, offset)
}

// GetActiveByUserID は指定ユーザーの進行中の学習目標を取得する（新しい順）。
func (r *learningGoalRepository) GetActiveByUserID(ctx context.Context, userID uint) ([]model.LearningGoal, error) {
	var goals []model.LearningGoal
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, model.GoalStatusActive).
		Order("created_at DESC").Find(&goals).Error
	return goals, err
}

// GetByCategory は指定ユーザーの学習目標をカテゴリで絞り込んで取得する（新しい順）。
func (r *learningGoalRepository) GetByCategory(ctx context.Context, userID uint, category string) ([]model.LearningGoal, error) {
	var goals []model.LearningGoal
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND category = ?", userID, category).
		Order("created_at DESC").Find(&goals).Error
	return goals, err
}

// GetByStatus は指定ユーザーの学習目標をステータスで絞り込んで取得する（新しい順）。
func (r *learningGoalRepository) GetByStatus(ctx context.Context, userID uint, status string) ([]model.LearningGoal, error) {
	var goals []model.LearningGoal
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, status).
		Order("created_at DESC").Find(&goals).Error
	return goals, err
}

// GetPublicByUserID は指定ユーザーの公開済み学習目標をページネーション付きで取得する（新しい順）。
func (r *learningGoalRepository) GetPublicByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.LearningGoal, int64, error) {
	return r.paginated(ctx, r.db.WithContext(ctx).Where("user_id = ? AND is_public = ?", userID, true), limit, offset)
}

// GetPublicGoals は全ユーザーの公開済み学習目標をページネーション付きで取得する（新しい順）。
func (r *learningGoalRepository) GetPublicGoals(ctx context.Context, limit, offset int) ([]model.LearningGoal, int64, error) {
	return r.paginated(ctx, r.db.WithContext(ctx).Where("is_public = ?", true), limit, offset)
}

// paginated は絞り込み済みクエリに対して総件数とページを取得する共通処理。
func (r *learningGoalRepository) paginated(ctx context.Context, scope *gorm.DB, limit, offset int) ([]model.LearningGoal, int64, error) {
	var total int64
	if err := scope.Session(&gorm.Session{}).Model(&model.LearningGoal{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var goals []model.LearningGoal
	err := scope.Session(&gorm.Session{}).
		Order("created_at DESC").Limit(limit).Offset(offset).
		Find(&goals).Error
	return goals, total, err
}

// CountByUserID は指定ユーザーの学習目標総数を返す。
func (r *learningGoalRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.LearningGoal{}).
		Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// GetStats は指定ユーザーの学習目標統計（総数・アクティブ数・完了数・平均進捗）を返す。
func (r *learningGoalRepository) GetStats(ctx context.Context, userID uint) (*model.LearningGoalStats, error) {
	db := r.db.WithContext(ctx)
	var stats model.LearningGoalStats

	count := func(scope *gorm.DB) (int64, error) {
		var n int64
		err := scope.Count(&n).Error
		return n, err
	}

	total, err := count(db.Model(&model.LearningGoal{}).Where("user_id = ?", userID))
	if err != nil {
		return nil, err
	}
	stats.TotalGoals = int(total)

	active, err := count(db.Model(&model.LearningGoal{}).Where("user_id = ? AND status = ?", userID, model.GoalStatusActive))
	if err != nil {
		return nil, err
	}
	stats.ActiveGoals = int(active)

	completed, err := count(db.Model(&model.LearningGoal{}).Where("user_id = ? AND status = ?", userID, model.GoalStatusCompleted))
	if err != nil {
		return nil, err
	}
	stats.CompletedGoals = int(completed)

	var avgProgress float64
	if err := db.Model(&model.LearningGoal{}).
		Where("user_id = ? AND status = ?", userID, model.GoalStatusActive).
		Select("COALESCE(AVG(progress), 0)").Scan(&avgProgress).Error; err != nil {
		return nil, err
	}
	stats.AverageProgress = int(avgProgress)

	return &stats, nil
}

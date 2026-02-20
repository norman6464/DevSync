package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// LearningGoalRepository は学習目標データへのアクセスを提供するリポジトリ実装。
type LearningGoalRepository struct {
	db *gorm.DB
}

// NewLearningGoalRepository は新しいLearningGoalRepositoryインスタンスを生成する。
func NewLearningGoalRepository(db *gorm.DB) *LearningGoalRepository {
	return &LearningGoalRepository{db: db}
}

// Create は新しい学習目標をデータベースに作成する。
func (r *LearningGoalRepository) Create(goal *model.LearningGoal) error {
	return r.db.Create(goal).Error
}

// Update は既存の学習目標を更新する。
func (r *LearningGoalRepository) Update(goal *model.LearningGoal) error {
	return r.db.Save(goal).Error
}

// Delete は指定IDの学習目標を削除する。
func (r *LearningGoalRepository) Delete(id uint) error {
	return r.db.Delete(&model.LearningGoal{}, id).Error
}

// FindByID は指定IDの学習目標を取得する。
func (r *LearningGoalRepository) FindByID(id uint) (*model.LearningGoal, error) {
	var goal model.LearningGoal
	err := r.db.First(&goal, id).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

// GetByUserID は指定ユーザーの学習目標をページネーション付きで取得する（新しい順）。
func (r *LearningGoalRepository) GetByUserID(userID uint, limit, offset int) ([]model.LearningGoal, int64, error) {
	var goals []model.LearningGoal
	var total int64
	query := r.db.Where("user_id = ?", userID)
	query.Model(&model.LearningGoal{}).Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&goals).Error
	return goals, total, err
}

// GetActiveByUserID は指定ユーザーの進行中の学習目標を取得する（新しい順）。
func (r *LearningGoalRepository) GetActiveByUserID(userID uint) ([]model.LearningGoal, error) {
	var goals []model.LearningGoal
	err := r.db.Where("user_id = ? AND status = ?", userID, model.GoalStatusActive).Order("created_at DESC").Find(&goals).Error
	return goals, err
}

// GetByCategory は指定ユーザーの学習目標をカテゴリでフィルタリングして取得する（新しい順）。
func (r *LearningGoalRepository) GetByCategory(userID uint, category string) ([]model.LearningGoal, error) {
	var goals []model.LearningGoal
	err := r.db.Where("user_id = ? AND category = ?", userID, category).Order("created_at DESC").Find(&goals).Error
	return goals, err
}

// GetByStatus は指定ユーザーの学習目標をステータスでフィルタリングして取得する（新しい順）。
func (r *LearningGoalRepository) GetByStatus(userID uint, status string) ([]model.LearningGoal, error) {
	var goals []model.LearningGoal
	err := r.db.Where("user_id = ? AND status = ?", userID, status).Order("created_at DESC").Find(&goals).Error
	return goals, err
}

// GetStats は指定ユーザーの学習目標統計情報を算出する。
// 目標総数、アクティブ数、完了数、平均進捗率を返す。
func (r *LearningGoalRepository) GetStats(userID uint) (*model.LearningGoalStats, error) {
	var stats model.LearningGoalStats

	// 目標総数を取得
	var totalCount int64
	r.db.Model(&model.LearningGoal{}).Where("user_id = ?", userID).Count(&totalCount)
	stats.TotalGoals = int(totalCount)

	// アクティブ目標数を取得
	var activeCount int64
	r.db.Model(&model.LearningGoal{}).Where("user_id = ? AND status = ?", userID, model.GoalStatusActive).Count(&activeCount)
	stats.ActiveGoals = int(activeCount)

	// 完了済み目標数を取得
	var completedCount int64
	r.db.Model(&model.LearningGoal{}).Where("user_id = ? AND status = ?", userID, model.GoalStatusCompleted).Count(&completedCount)
	stats.CompletedGoals = int(completedCount)

	// アクティブ目標の平均進捗率を算出
	var avgProgress float64
	r.db.Model(&model.LearningGoal{}).Where("user_id = ? AND status = ?", userID, model.GoalStatusActive).Select("COALESCE(AVG(progress), 0)").Scan(&avgProgress)
	stats.AverageProgress = int(avgProgress)

	return &stats, nil
}

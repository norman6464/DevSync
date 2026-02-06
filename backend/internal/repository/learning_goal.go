package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

type LearningGoalRepository struct {
	db *gorm.DB
}

func NewLearningGoalRepository(db *gorm.DB) *LearningGoalRepository {
	return &LearningGoalRepository{db: db}
}

// Create creates a new learning goal
func (r *LearningGoalRepository) Create(goal *model.LearningGoal) error {
	return r.db.Create(goal).Error
}

// Update updates an existing learning goal
func (r *LearningGoalRepository) Update(goal *model.LearningGoal) error {
	return r.db.Save(goal).Error
}

// Delete deletes a learning goal
func (r *LearningGoalRepository) Delete(id uint) error {
	return r.db.Delete(&model.LearningGoal{}, id).Error
}

// FindByID finds a learning goal by ID
func (r *LearningGoalRepository) FindByID(id uint) (*model.LearningGoal, error) {
	var goal model.LearningGoal
	err := r.db.First(&goal, id).Error
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

// GetByUserID gets all learning goals for a user
func (r *LearningGoalRepository) GetByUserID(userID uint) ([]model.LearningGoal, error) {
	var goals []model.LearningGoal
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&goals).Error
	return goals, err
}

// GetActiveByUserID gets all active learning goals for a user
func (r *LearningGoalRepository) GetActiveByUserID(userID uint) ([]model.LearningGoal, error) {
	var goals []model.LearningGoal
	err := r.db.Where("user_id = ? AND status = ?", userID, model.GoalStatusActive).Order("created_at DESC").Find(&goals).Error
	return goals, err
}

// GetStats gets learning goal statistics for a user
func (r *LearningGoalRepository) GetStats(userID uint) (*model.LearningGoalStats, error) {
	var stats model.LearningGoalStats

	// Get total goals
	var totalCount int64
	r.db.Model(&model.LearningGoal{}).Where("user_id = ?", userID).Count(&totalCount)
	stats.TotalGoals = int(totalCount)

	// Get active goals
	var activeCount int64
	r.db.Model(&model.LearningGoal{}).Where("user_id = ? AND status = ?", userID, model.GoalStatusActive).Count(&activeCount)
	stats.ActiveGoals = int(activeCount)

	// Get completed goals
	var completedCount int64
	r.db.Model(&model.LearningGoal{}).Where("user_id = ? AND status = ?", userID, model.GoalStatusCompleted).Count(&completedCount)
	stats.CompletedGoals = int(completedCount)

	// Get average progress of active goals
	var avgProgress float64
	r.db.Model(&model.LearningGoal{}).Where("user_id = ? AND status = ?", userID, model.GoalStatusActive).Select("COALESCE(AVG(progress), 0)").Scan(&avgProgress)
	stats.AverageProgress = int(avgProgress)

	return &stats, nil
}

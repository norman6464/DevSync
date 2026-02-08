package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

type LearningLogRepository struct {
	db *gorm.DB
}

func NewLearningLogRepository(db *gorm.DB) *LearningLogRepository {
	return &LearningLogRepository{db: db}
}

// Create creates a new learning log
func (r *LearningLogRepository) Create(log *model.LearningLog) error {
	return r.db.Create(log).Error
}

// Update updates an existing learning log
func (r *LearningLogRepository) Update(log *model.LearningLog) error {
	return r.db.Save(log).Error
}

// Delete deletes a learning log by ID and user ID (ownership check)
func (r *LearningLogRepository) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.LearningLog{}).Error
}

// FindByID finds a learning log by ID
func (r *LearningLogRepository) FindByID(id uint) (*model.LearningLog, error) {
	var log model.LearningLog
	err := r.db.First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetByUserID gets all learning logs for a user
func (r *LearningLogRepository) GetByUserID(userID uint) ([]model.LearningLog, error) {
	var logs []model.LearningLog
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// GetCalendarData returns daily log counts for calendar visualization
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

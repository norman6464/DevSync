package repository

import (
	"sort"
	"time"

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

// GetStreakInfo calculates streak data from learning logs
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

	// Sort descending
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Date.After(dates[j].Date)
	})

	// Calculate current streak
	today := time.Now().UTC().Truncate(24 * time.Hour)
	currentStreak := 0
	startIdx := 0

	firstDate := dates[0].Date.UTC().Truncate(24 * time.Hour)
	diffToToday := today.Sub(firstDate)

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

	// Calculate longest streak
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

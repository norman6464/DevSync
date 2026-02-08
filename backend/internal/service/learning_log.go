package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// LearningLogService handles learning log business logic.
type LearningLogService struct {
	repo repository.LearningLogRepositoryInterface
}

// NewLearningLogService creates a new LearningLogService.
func NewLearningLogService(repo repository.LearningLogRepositoryInterface) *LearningLogService {
	return &LearningLogService{repo: repo}
}

// Create creates a new learning log.
func (s *LearningLogService) Create(log *model.LearningLog) error {
	return s.repo.Create(log)
}

// GetByID returns a learning log by ID.
func (s *LearningLogService) GetByID(id uint) (*model.LearningLog, error) {
	return s.repo.FindByID(id)
}

// GetByUserID returns all learning logs for a user.
func (s *LearningLogService) GetByUserID(userID uint) ([]model.LearningLog, error) {
	return s.repo.GetByUserID(userID)
}

// Update updates a learning log after verifying ownership.
func (s *LearningLogService) Update(id, userID uint, updates *model.LearningLog) (*model.LearningLog, error) {
	log, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if log.UserID != userID {
		return nil, ErrForbidden
	}

	if updates.Title != "" {
		log.Title = updates.Title
	}
	if updates.Content != "" {
		log.Content = updates.Content
	}
	if updates.Category != "" {
		log.Category = updates.Category
	}
	if updates.Duration != 0 {
		log.Duration = updates.Duration
	}

	if err := s.repo.Update(log); err != nil {
		return nil, err
	}
	return log, nil
}

// Delete deletes a learning log after verifying ownership.
func (s *LearningLogService) Delete(id, userID uint) error {
	log, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if log.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id, userID)
}

// GetStreakInfo returns streak data for a user.
func (s *LearningLogService) GetStreakInfo(userID uint) (*model.StreakInfo, error) {
	return s.repo.GetStreakInfo(userID)
}

// GetCalendarData returns daily log counts for calendar visualization.
func (s *LearningLogService) GetCalendarData(userID uint) ([]model.CalendarEntry, error) {
	return s.repo.GetCalendarData(userID)
}

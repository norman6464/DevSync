package service

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// LearningGoalService handles learning goal business logic.
type LearningGoalService struct {
	repo repository.LearningGoalRepositoryInterface
}

// NewLearningGoalService creates a new LearningGoalService.
func NewLearningGoalService(repo repository.LearningGoalRepositoryInterface) *LearningGoalService {
	return &LearningGoalService{repo: repo}
}

// Create creates a new learning goal.
func (s *LearningGoalService) Create(goal *model.LearningGoal) error {
	return s.repo.Create(goal)
}

// GetByID returns a learning goal by ID.
func (s *LearningGoalService) GetByID(id uint) (*model.LearningGoal, error) {
	return s.repo.FindByID(id)
}

// GetByUserID returns all learning goals for a user.
func (s *LearningGoalService) GetByUserID(userID uint) ([]model.LearningGoal, error) {
	return s.repo.GetByUserID(userID)
}

// GetActiveByUserID returns active learning goals for a user.
func (s *LearningGoalService) GetActiveByUserID(userID uint) ([]model.LearningGoal, error) {
	return s.repo.GetActiveByUserID(userID)
}

// GetStats returns learning goal statistics for a user.
func (s *LearningGoalService) GetStats(userID uint) (*model.LearningGoalStats, error) {
	return s.repo.GetStats(userID)
}

// Update updates a learning goal after verifying ownership.
func (s *LearningGoalService) Update(id, userID uint, updates *model.LearningGoal) (*model.LearningGoal, error) {
	goal, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if goal.UserID != userID {
		return nil, ErrForbidden
	}

	if updates.Title != "" {
		goal.Title = updates.Title
	}
	if updates.Description != "" {
		goal.Description = updates.Description
	}
	if updates.Category != "" {
		goal.Category = updates.Category
	}
	if updates.TargetDate != nil {
		goal.TargetDate = updates.TargetDate
	}
	if updates.Progress >= 0 {
		progress := updates.Progress
		if progress > 100 {
			progress = 100
		}
		goal.Progress = progress

		// Auto-complete if progress reaches 100
		if progress == 100 && goal.Status == model.GoalStatusActive {
			goal.Status = model.GoalStatusCompleted
			now := time.Now()
			goal.CompletedAt = &now
		}
	}
	if updates.Status != "" {
		goal.Status = updates.Status
		if goal.Status == model.GoalStatusCompleted && goal.CompletedAt == nil {
			now := time.Now()
			goal.CompletedAt = &now
		}
	}

	if err := s.repo.Update(goal); err != nil {
		return nil, err
	}
	return goal, nil
}

// Delete deletes a learning goal after verifying ownership.
func (s *LearningGoalService) Delete(id, userID uint) error {
	goal, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if goal.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}

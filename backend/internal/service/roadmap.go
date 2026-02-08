package service

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// StepOrder represents the order of a step (used in reordering).
type StepOrder = repository.StepOrder

// RoadmapService handles roadmap business logic.
type RoadmapService struct {
	repo repository.RoadmapRepositoryInterface
}

// NewRoadmapService creates a new RoadmapService.
func NewRoadmapService(repo repository.RoadmapRepositoryInterface) *RoadmapService {
	return &RoadmapService{repo: repo}
}

// Create creates a new roadmap.
func (s *RoadmapService) Create(roadmap *model.Roadmap) error {
	return s.repo.Create(roadmap)
}

// GetByID returns a roadmap by ID, checking visibility.
func (s *RoadmapService) GetByID(id, userID uint) (*model.Roadmap, error) {
	roadmap, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if roadmap.UserID != userID && !roadmap.IsPublic {
		return nil, ErrForbidden
	}
	return roadmap, nil
}

// GetByUserID returns all roadmaps for a user.
func (s *RoadmapService) GetByUserID(userID uint) ([]model.Roadmap, error) {
	return s.repo.GetByUserID(userID)
}

// GetPublicRoadmaps returns paginated public roadmaps.
func (s *RoadmapService) GetPublicRoadmaps(limit, offset int) ([]model.Roadmap, int64, error) {
	return s.repo.GetPublicRoadmaps(limit, offset)
}

// GetStats returns roadmap statistics for a user.
func (s *RoadmapService) GetStats(userID uint) (*model.RoadmapStats, error) {
	return s.repo.GetStats(userID)
}

// Update updates a roadmap after verifying ownership.
func (s *RoadmapService) Update(id, userID uint, updates *model.Roadmap) (*model.Roadmap, error) {
	roadmap, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if roadmap.UserID != userID {
		return nil, ErrForbidden
	}

	if updates.Title != "" {
		roadmap.Title = updates.Title
	}
	if updates.Description != "" {
		roadmap.Description = updates.Description
	}
	if updates.Category != "" {
		roadmap.Category = updates.Category
	}
	if updates.Status != "" {
		roadmap.Status = updates.Status
		if roadmap.Status == model.RoadmapStatusCompleted && roadmap.CompletedAt == nil {
			now := time.Now()
			roadmap.CompletedAt = &now
		}
	}

	if err := s.repo.Update(roadmap); err != nil {
		return nil, err
	}
	return roadmap, nil
}

// UpdateVisibility updates the public/private status of a roadmap after verifying ownership.
func (s *RoadmapService) UpdateVisibility(id, userID uint, isPublic bool) (*model.Roadmap, error) {
	roadmap, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if roadmap.UserID != userID {
		return nil, ErrForbidden
	}

	roadmap.IsPublic = isPublic

	if err := s.repo.Update(roadmap); err != nil {
		return nil, err
	}
	return roadmap, nil
}

// Delete deletes a roadmap after verifying ownership.
func (s *RoadmapService) Delete(id, userID uint) error {
	roadmap, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if roadmap.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}

// CopyRoadmap copies a public roadmap as a template.
func (s *RoadmapService) CopyRoadmap(roadmapID, userID uint) (*model.Roadmap, error) {
	original, err := s.repo.FindByID(roadmapID)
	if err != nil {
		return nil, err
	}
	if !original.IsPublic && original.UserID != userID {
		return nil, ErrForbidden
	}
	return s.repo.CopyRoadmap(roadmapID, userID)
}

// CreateStep creates a new step in a roadmap after verifying ownership.
func (s *RoadmapService) CreateStep(roadmapID, userID uint, step *model.RoadmapStep) error {
	roadmap, err := s.repo.FindByID(roadmapID)
	if err != nil {
		return err
	}
	if roadmap.UserID != userID {
		return ErrForbidden
	}

	step.RoadmapID = roadmapID
	return s.repo.CreateStep(step)
}

// UpdateStep updates a step after verifying roadmap ownership.
func (s *RoadmapService) UpdateStep(roadmapID, stepID, userID uint, updates *model.RoadmapStep) (*model.RoadmapStep, error) {
	roadmap, err := s.repo.FindByID(roadmapID)
	if err != nil {
		return nil, err
	}
	if roadmap.UserID != userID {
		return nil, ErrForbidden
	}

	step, err := s.repo.FindStepByID(stepID)
	if err != nil {
		return nil, err
	}
	if step.RoadmapID != roadmapID {
		return nil, ErrBadRequest
	}

	if updates.Title != "" {
		step.Title = updates.Title
	}
	if updates.Description != "" {
		step.Description = updates.Description
	}
	if updates.ResourceURL != "" {
		step.ResourceURL = updates.ResourceURL
	}

	if err := s.repo.UpdateStep(step); err != nil {
		return nil, err
	}
	return step, nil
}

// UpdateStepCompletion updates the completion status of a step.
func (s *RoadmapService) UpdateStepCompletion(roadmapID, stepID, userID uint, isCompleted bool) (*model.RoadmapStep, error) {
	roadmap, err := s.repo.FindByID(roadmapID)
	if err != nil {
		return nil, err
	}
	if roadmap.UserID != userID {
		return nil, ErrForbidden
	}

	step, err := s.repo.FindStepByID(stepID)
	if err != nil {
		return nil, err
	}
	if step.RoadmapID != roadmapID {
		return nil, ErrBadRequest
	}

	step.IsCompleted = isCompleted
	if isCompleted && step.CompletedAt == nil {
		now := time.Now()
		step.CompletedAt = &now
	} else if !isCompleted {
		step.CompletedAt = nil
	}

	if err := s.repo.UpdateStep(step); err != nil {
		return nil, err
	}
	return step, nil
}

// DeleteStep deletes a step after verifying roadmap ownership.
func (s *RoadmapService) DeleteStep(roadmapID, stepID, userID uint) error {
	roadmap, err := s.repo.FindByID(roadmapID)
	if err != nil {
		return err
	}
	if roadmap.UserID != userID {
		return ErrForbidden
	}

	step, err := s.repo.FindStepByID(stepID)
	if err != nil {
		return err
	}
	if step.RoadmapID != roadmapID {
		return ErrBadRequest
	}

	return s.repo.DeleteStep(stepID)
}

// ReorderSteps reorders steps within a roadmap after verifying ownership.
func (s *RoadmapService) ReorderSteps(roadmapID, userID uint, orders []repository.StepOrder) error {
	roadmap, err := s.repo.FindByID(roadmapID)
	if err != nil {
		return err
	}
	if roadmap.UserID != userID {
		return ErrForbidden
	}
	return s.repo.ReorderSteps(roadmapID, orders)
}

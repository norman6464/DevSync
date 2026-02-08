package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// LearningResourceService handles learning resource business logic.
type LearningResourceService struct {
	repo repository.LearningResourceRepositoryInterface
}

// NewLearningResourceService creates a new LearningResourceService.
func NewLearningResourceService(repo repository.LearningResourceRepositoryInterface) *LearningResourceService {
	return &LearningResourceService{repo: repo}
}

// Create creates a new learning resource.
func (s *LearningResourceService) Create(resource *model.LearningResource) error {
	return s.repo.Create(resource)
}

// GetByID returns a learning resource by ID, checking visibility.
func (s *LearningResourceService) GetByID(id, userID uint) (*model.LearningResource, error) {
	resource, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Check if resource is private and user is not owner
	if !resource.IsPublic && resource.UserID != userID {
		return nil, ErrForbidden
	}

	return resource, nil
}

// HasLiked checks if a user has liked a resource.
func (s *LearningResourceService) HasLiked(userID, resourceID uint) (bool, error) {
	return s.repo.HasLiked(userID, resourceID)
}

// HasSaved checks if a user has saved a resource.
func (s *LearningResourceService) HasSaved(userID, resourceID uint) (bool, error) {
	return s.repo.HasSaved(userID, resourceID)
}

// GetByUserID returns learning resources for a user.
func (s *LearningResourceService) GetByUserID(targetUserID, currentUserID uint) ([]model.LearningResource, error) {
	includePrivate := currentUserID == targetUserID
	return s.repo.FindByUserID(targetUserID, includePrivate)
}

// GetPublic returns paginated public learning resources.
func (s *LearningResourceService) GetPublic(limit, offset int, category, difficulty string) ([]model.LearningResource, int64, error) {
	return s.repo.FindPublic(limit, offset, category, difficulty)
}

// Search searches learning resources.
func (s *LearningResourceService) Search(query string, limit, offset int) ([]model.LearningResource, int64, error) {
	return s.repo.Search(query, limit, offset)
}

// Update updates a learning resource after verifying ownership.
func (s *LearningResourceService) Update(id, userID uint, updates *model.LearningResource) (*model.LearningResource, error) {
	resource, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if resource.UserID != userID {
		return nil, ErrForbidden
	}

	if updates.Title != "" {
		resource.Title = updates.Title
	}
	if updates.Description != "" {
		resource.Description = updates.Description
	}
	if updates.URL != "" {
		resource.URL = updates.URL
	}
	if updates.Category != "" {
		resource.Category = updates.Category
	}
	if updates.Difficulty != "" {
		resource.Difficulty = updates.Difficulty
	}
	if updates.Tags != "" {
		resource.Tags = updates.Tags
	}
	if updates.ImageURL != "" {
		resource.ImageURL = updates.ImageURL
	}

	if err := s.repo.Update(resource); err != nil {
		return nil, err
	}
	return resource, nil
}

// UpdateVisibility updates the public/private status of a resource after verifying ownership.
func (s *LearningResourceService) UpdateVisibility(id, userID uint, isPublic bool) (*model.LearningResource, error) {
	resource, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if resource.UserID != userID {
		return nil, ErrForbidden
	}

	resource.IsPublic = isPublic

	if err := s.repo.Update(resource); err != nil {
		return nil, err
	}
	return resource, nil
}

// Delete deletes a learning resource after verifying ownership.
func (s *LearningResourceService) Delete(id, userID uint) error {
	resource, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if resource.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}

// Like likes a learning resource.
func (s *LearningResourceService) Like(userID, resourceID uint) error {
	return s.repo.Like(userID, resourceID)
}

// Unlike unlikes a learning resource.
func (s *LearningResourceService) Unlike(userID, resourceID uint) error {
	return s.repo.Unlike(userID, resourceID)
}

// Save saves a learning resource.
func (s *LearningResourceService) Save(userID, resourceID uint) error {
	return s.repo.Save(userID, resourceID)
}

// Unsave unsaves a learning resource.
func (s *LearningResourceService) Unsave(userID, resourceID uint) error {
	return s.repo.Unsave(userID, resourceID)
}

// GetSavedByUserID returns paginated saved learning resources for a user.
func (s *LearningResourceService) GetSavedByUserID(userID uint, limit, offset int) ([]model.LearningResource, int64, error) {
	return s.repo.FindSavedByUserID(userID, limit, offset)
}

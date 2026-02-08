package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// ProjectService handles project business logic.
type ProjectService struct {
	repo repository.ProjectRepositoryInterface
}

// NewProjectService creates a new ProjectService.
func NewProjectService(repo repository.ProjectRepositoryInterface) *ProjectService {
	return &ProjectService{repo: repo}
}

// Create creates a new project.
func (s *ProjectService) Create(project *model.Project) error {
	return s.repo.Create(project)
}

// GetByID returns a project by ID.
func (s *ProjectService) GetByID(id uint) (*model.Project, error) {
	return s.repo.FindByID(id)
}

// GetByUserID returns all projects for a user.
func (s *ProjectService) GetByUserID(userID uint) ([]model.Project, error) {
	return s.repo.FindByUserID(userID)
}

// GetFeaturedByUserID returns featured projects for a user.
func (s *ProjectService) GetFeaturedByUserID(userID uint) ([]model.Project, error) {
	return s.repo.FindFeaturedByUserID(userID)
}

// GetAll returns paginated projects.
func (s *ProjectService) GetAll(limit, offset int) ([]model.Project, int64, error) {
	return s.repo.FindAll(limit, offset)
}

// Update updates a project after verifying ownership.
func (s *ProjectService) Update(id, userID uint, updates *model.Project) (*model.Project, error) {
	project, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrForbidden
	}

	if updates.Title != "" {
		project.Title = updates.Title
	}
	if updates.Description != "" {
		project.Description = updates.Description
	}
	if updates.TechStack != "" {
		project.TechStack = updates.TechStack
	}
	if updates.DemoURL != "" {
		project.DemoURL = updates.DemoURL
	}
	if updates.GithubURL != "" {
		project.GithubURL = updates.GithubURL
	}
	if updates.ImageURL != "" {
		project.ImageURL = updates.ImageURL
	}
	if updates.Role != "" {
		project.Role = updates.Role
	}
	if updates.StartDate != nil {
		project.StartDate = updates.StartDate
	}
	if updates.EndDate != nil {
		project.EndDate = updates.EndDate
	}
	if updates.GithubRepoID != nil {
		project.GithubRepoID = updates.GithubRepoID
	}

	if err := s.repo.Update(project); err != nil {
		return nil, err
	}
	return project, nil
}

// UpdateFeatured updates the featured status of a project after verifying ownership.
func (s *ProjectService) UpdateFeatured(id, userID uint, featured bool) (*model.Project, error) {
	project, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, ErrForbidden
	}

	project.Featured = featured

	if err := s.repo.Update(project); err != nil {
		return nil, err
	}
	return project, nil
}

// Delete deletes a project after verifying ownership.
func (s *ProjectService) Delete(id, userID uint) error {
	project, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if project.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}

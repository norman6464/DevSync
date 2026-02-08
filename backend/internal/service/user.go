package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// UserService handles user business logic.
type UserService struct {
	repo repository.UserRepositoryInterface
}

// NewUserService creates a new UserService.
func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

// GetAll returns all users, optionally filtered by search query.
func (s *UserService) GetAll(query string) ([]model.User, error) {
	if query != "" {
		return s.repo.Search(query)
	}
	return s.repo.FindAll()
}

// GetByID returns a user by ID.
func (s *UserService) GetByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

// FindByID returns a user by ID (alias for repository compatibility).
func (s *UserService) FindByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

// Update updates a user's information.
func (s *UserService) Update(user *model.User) error {
	return s.repo.Update(user)
}

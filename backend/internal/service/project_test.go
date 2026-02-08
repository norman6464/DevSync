package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newTestProjectService() (*ProjectService, *MockProjectRepository) {
	repo := new(MockProjectRepository)
	svc := NewProjectService(repo)
	return svc, repo
}

// ============================================================
// Update
// ============================================================

func TestProjectUpdate_Success(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{Title: "Old", Description: "Old Desc", UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.Project{Title: "New", TechStack: "Go, React"}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New", result.Title)
	assert.Equal(t, "Go, React", result.TechStack)
	assert.Equal(t, "Old Desc", result.Description)
	repo.AssertExpectations(t)
}

func TestProjectUpdate_Forbidden(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.Project{Title: "New"}
	result, err := svc.Update(1, 999, updates)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestProjectUpdate_NotFound(t *testing.T) {
	svc, repo := newTestProjectService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	updates := &model.Project{Title: "New"}
	result, err := svc.Update(999, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// UpdateFeatured
// ============================================================

func TestProjectUpdateFeatured_Success(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 1, Featured: false}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.UpdateFeatured(1, 1, true)
	assert.NoError(t, err)
	assert.True(t, result.Featured)
	repo.AssertExpectations(t)
}

func TestProjectUpdateFeatured_Forbidden(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	result, err := svc.UpdateFeatured(1, 999, true)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// Delete
// ============================================================

func TestProjectDelete_Success(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestProjectDelete_Forbidden(t *testing.T) {
	svc, repo := newTestProjectService()

	existing := &model.Project{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

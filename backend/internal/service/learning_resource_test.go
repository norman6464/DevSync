package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newTestLearningResourceService() (*LearningResourceService, *MockLearningResourceRepository) {
	repo := new(MockLearningResourceRepository)
	svc := NewLearningResourceService(repo)
	return svc, repo
}

// ============================================================
// GetByID（可視性チェック）
// ============================================================

func TestLearningResourceGetByID_PublicResource(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	resource := &model.LearningResource{Title: "Public Resource", UserID: 1, IsPublic: true}
	resource.ID = 1

	repo.On("FindByID", uint(1)).Return(resource, nil)

	// 他人がアクセスしても公開なのでOK
	result, err := svc.GetByID(1, 999)
	assert.NoError(t, err)
	assert.Equal(t, "Public Resource", result.Title)
	repo.AssertExpectations(t)
}

func TestLearningResourceGetByID_PrivateResourceOwner(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	resource := &model.LearningResource{Title: "Private Resource", UserID: 1, IsPublic: false}
	resource.ID = 1

	repo.On("FindByID", uint(1)).Return(resource, nil)

	// 所有者がアクセス → OK
	result, err := svc.GetByID(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Private Resource", result.Title)
	repo.AssertExpectations(t)
}

func TestLearningResourceGetByID_PrivateResourceForbidden(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	resource := &model.LearningResource{Title: "Private Resource", UserID: 1, IsPublic: false}
	resource.ID = 1

	repo.On("FindByID", uint(1)).Return(resource, nil)

	// 他人がアクセス → Forbidden
	result, err := svc.GetByID(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestLearningResourceGetByID_NotFound(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.GetByID(999, 1)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByUserID（自分 vs 他人の可視性）
// ============================================================

func TestLearningResourceGetByUserID_Self(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	resources := []model.LearningResource{
		{Title: "Public", IsPublic: true},
		{Title: "Private", IsPublic: false},
	}
	// 自分 → includePrivate=true
	repo.On("FindByUserID", uint(1), true).Return(resources, nil)

	result, err := svc.GetByUserID(1, 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestLearningResourceGetByUserID_Other(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	resources := []model.LearningResource{
		{Title: "Public Only", IsPublic: true},
	}
	// 他人 → includePrivate=false
	repo.On("FindByUserID", uint(1), false).Return(resources, nil)

	result, err := svc.GetByUserID(1, 999)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	repo.AssertExpectations(t)
}

// ============================================================
// Update
// ============================================================

func TestLearningResourceUpdate_Success(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	existing := &model.LearningResource{Title: "Old", Description: "Desc", UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.LearningResource{Title: "New Title", Category: "video"}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	assert.Equal(t, model.ResourceCategory("video"), result.Category)
	assert.Equal(t, "Desc", result.Description) // 変更なし
	repo.AssertExpectations(t)
}

func TestLearningResourceUpdate_Forbidden(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	existing := &model.LearningResource{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.LearningResource{Title: "New"}
	result, err := svc.Update(1, 999, updates)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// UpdateVisibility
// ============================================================

func TestLearningResourceUpdateVisibility_Success(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	existing := &model.LearningResource{UserID: 1, IsPublic: false}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	result, err := svc.UpdateVisibility(1, 1, true)
	assert.NoError(t, err)
	assert.True(t, result.IsPublic)
	repo.AssertExpectations(t)
}

func TestLearningResourceUpdateVisibility_Forbidden(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	existing := &model.LearningResource{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	result, err := svc.UpdateVisibility(1, 999, true)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// Delete
// ============================================================

func TestLearningResourceDelete_Success(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	existing := &model.LearningResource{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLearningResourceDelete_Forbidden(t *testing.T) {
	svc, repo := newTestLearningResourceService()

	existing := &model.LearningResource{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

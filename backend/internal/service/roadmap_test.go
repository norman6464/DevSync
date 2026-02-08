package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newTestRoadmapService() (*RoadmapService, *MockRoadmapRepository) {
	repo := new(MockRoadmapRepository)
	svc := NewRoadmapService(repo)
	return svc, repo
}

// ============================================================
// GetByID（可視性チェック）
// ============================================================

func TestRoadmapGetByID_PublicRoadmap(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{Title: "Public Roadmap", UserID: 1, IsPublic: true}
	roadmap.ID = 1

	repo.On("FindByID", uint(1)).Return(roadmap, nil)

	result, err := svc.GetByID(1, 999) // 他人がアクセス
	assert.NoError(t, err)
	assert.Equal(t, "Public Roadmap", result.Title)
	repo.AssertExpectations(t)
}

func TestRoadmapGetByID_PrivateOwner(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{Title: "Private", UserID: 1, IsPublic: false}
	roadmap.ID = 1

	repo.On("FindByID", uint(1)).Return(roadmap, nil)

	result, err := svc.GetByID(1, 1) // 所有者がアクセス
	assert.NoError(t, err)
	assert.Equal(t, "Private", result.Title)
	repo.AssertExpectations(t)
}

func TestRoadmapGetByID_PrivateForbidden(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{Title: "Private", UserID: 1, IsPublic: false}
	roadmap.ID = 1

	repo.On("FindByID", uint(1)).Return(roadmap, nil)

	result, err := svc.GetByID(1, 999) // 他人がアクセス
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// CopyRoadmap
// ============================================================

func TestCopyRoadmap_PublicSuccess(t *testing.T) {
	svc, repo := newTestRoadmapService()

	original := &model.Roadmap{Title: "Original", UserID: 1, IsPublic: true}
	original.ID = 1

	copied := &model.Roadmap{Title: "Original", UserID: 2}
	copied.ID = 2

	repo.On("FindByID", uint(1)).Return(original, nil)
	repo.On("CopyRoadmap", uint(1), uint(2)).Return(copied, nil)

	result, err := svc.CopyRoadmap(1, 2)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(2), result.UserID)
	repo.AssertExpectations(t)
}

func TestCopyRoadmap_PrivateForbidden(t *testing.T) {
	svc, repo := newTestRoadmapService()

	original := &model.Roadmap{Title: "Private", UserID: 1, IsPublic: false}
	original.ID = 1

	repo.On("FindByID", uint(1)).Return(original, nil)

	result, err := svc.CopyRoadmap(1, 999) // 他人が非公開をコピー
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestCopyRoadmap_PrivateOwnerCanCopy(t *testing.T) {
	svc, repo := newTestRoadmapService()

	original := &model.Roadmap{Title: "Private", UserID: 1, IsPublic: false}
	original.ID = 1

	copied := &model.Roadmap{Title: "Private", UserID: 1}
	copied.ID = 2

	repo.On("FindByID", uint(1)).Return(original, nil)
	repo.On("CopyRoadmap", uint(1), uint(1)).Return(copied, nil)

	result, err := svc.CopyRoadmap(1, 1) // 所有者は自分のをコピー可能
	assert.NoError(t, err)
	assert.NotNil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// Update
// ============================================================

func TestRoadmapUpdate_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	existing := &model.Roadmap{Title: "Old", UserID: 1, Status: model.RoadmapStatusActive}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.Roadmap{Title: "New Title"}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	repo.AssertExpectations(t)
}

func TestRoadmapUpdate_Forbidden(t *testing.T) {
	svc, repo := newTestRoadmapService()

	existing := &model.Roadmap{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.Roadmap{Title: "New"}
	result, err := svc.Update(1, 999, updates)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestRoadmapUpdate_CompletedStatus(t *testing.T) {
	svc, repo := newTestRoadmapService()

	existing := &model.Roadmap{Title: "Roadmap", UserID: 1, Status: model.RoadmapStatusActive}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.Roadmap{Status: model.RoadmapStatusCompleted}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, model.RoadmapStatusCompleted, result.Status)
	assert.NotNil(t, result.CompletedAt)
	repo.AssertExpectations(t)
}

// ============================================================
// UpdateStep（ステップ所属チェック）
// ============================================================

func TestRoadmapUpdateStep_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	step := &model.RoadmapStep{RoadmapID: 10, Title: "Old Step"}
	step.ID = 5

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("UpdateStep", step).Return(nil)

	updates := &model.RoadmapStep{Title: "New Step"}
	result, err := svc.UpdateStep(10, 5, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New Step", result.Title)
	repo.AssertExpectations(t)
}

func TestRoadmapUpdateStep_StepBelongsToDifferentRoadmap(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	// ステップは別のロードマップに所属
	step := &model.RoadmapStep{RoadmapID: 20, Title: "Step"}
	step.ID = 5

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)

	updates := &model.RoadmapStep{Title: "New"}
	result, err := svc.UpdateStep(10, 5, 1, updates)
	assert.ErrorIs(t, err, ErrBadRequest)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestRoadmapUpdateStep_NotOwner(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	repo.On("FindByID", uint(10)).Return(roadmap, nil)

	updates := &model.RoadmapStep{Title: "New"}
	result, err := svc.UpdateStep(10, 5, 999, updates)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// DeleteStep
// ============================================================

func TestRoadmapDeleteStep_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	step := &model.RoadmapStep{RoadmapID: 10}
	step.ID = 5

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("DeleteStep", uint(5)).Return(nil)

	err := svc.DeleteStep(10, 5, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRoadmapDeleteStep_StepBelongsToDifferentRoadmap(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	step := &model.RoadmapStep{RoadmapID: 20}
	step.ID = 5

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)

	err := svc.DeleteStep(10, 5, 1)
	assert.ErrorIs(t, err, ErrBadRequest)
	repo.AssertExpectations(t)
}

// ============================================================
// Delete
// ============================================================

func TestRoadmapDelete_Success(t *testing.T) {
	svc, repo := newTestRoadmapService()

	existing := &model.Roadmap{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRoadmapDelete_Forbidden(t *testing.T) {
	svc, repo := newTestRoadmapService()

	existing := &model.Roadmap{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

// ============================================================
// UpdateStepCompletion
// ============================================================

func TestRoadmapUpdateStepCompletion_Complete(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	step := &model.RoadmapStep{RoadmapID: 10, IsCompleted: false}
	step.ID = 5

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("UpdateStep", step).Return(nil)

	result, err := svc.UpdateStepCompletion(10, 5, 1, true)
	assert.NoError(t, err)
	assert.True(t, result.IsCompleted)
	assert.NotNil(t, result.CompletedAt)
	repo.AssertExpectations(t)
}

func TestRoadmapUpdateStepCompletion_Uncomplete(t *testing.T) {
	svc, repo := newTestRoadmapService()

	roadmap := &model.Roadmap{UserID: 1}
	roadmap.ID = 10

	step := &model.RoadmapStep{RoadmapID: 10, IsCompleted: true}
	step.ID = 5

	repo.On("FindByID", uint(10)).Return(roadmap, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("UpdateStep", step).Return(nil)

	result, err := svc.UpdateStepCompletion(10, 5, 1, false)
	assert.NoError(t, err)
	assert.False(t, result.IsCompleted)
	assert.Nil(t, result.CompletedAt)
	repo.AssertExpectations(t)
}

// ============================================================
// Not Found
// ============================================================

func TestRoadmapGetByID_NotFound(t *testing.T) {
	svc, repo := newTestRoadmapService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	result, err := svc.GetByID(999, 1)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

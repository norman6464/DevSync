package service

import (
	"errors"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newTestResourceProgressService はテスト用のResourceProgressServiceとモックを返す。
func newTestResourceProgressService() (*ResourceProgressService, *MockResourceProgressRepository, *MockLearningResourceRepository) {
	progressRepo := new(MockResourceProgressRepository)
	resourceRepo := new(MockLearningResourceRepository)
	svc := NewResourceProgressService(progressRepo, resourceRepo)
	return svc, progressRepo, resourceRepo
}

// --- UpsertProgress ---

func TestResourceProgressService_UpsertProgress_NewProgress(t *testing.T) {
	svc, progressRepo, resourceRepo := newTestResourceProgressService()

	resourceRepo.On("FindByID", uint(1)).Return(&model.LearningResource{ID: 1, UserID: 2}, nil)
	progressRepo.On("Upsert", mock.MatchedBy(func(p *model.ResourceProgress) bool {
		return p.UserID == 10 && p.ResourceID == 1 && p.Status == model.ResourceProgressInProgress && p.CompletionPercent == 30
	})).Return(nil)
	progressRepo.On("FindByUserAndResource", uint(10), uint(1)).Return(&model.ResourceProgress{
		ID: 1, UserID: 10, ResourceID: 1, Status: model.ResourceProgressInProgress, CompletionPercent: 30,
	}, nil)

	result, err := svc.UpsertProgress(10, 1, string(model.ResourceProgressInProgress), 30, "学習中")
	assert.NoError(t, err)
	assert.Equal(t, uint(10), result.UserID)
	assert.Equal(t, model.ResourceProgressInProgress, result.Status)
	assert.Equal(t, 30, result.CompletionPercent)
}

func TestResourceProgressService_UpsertProgress_Completed(t *testing.T) {
	svc, progressRepo, resourceRepo := newTestResourceProgressService()

	resourceRepo.On("FindByID", uint(1)).Return(&model.LearningResource{ID: 1, UserID: 2}, nil)
	progressRepo.On("Upsert", mock.MatchedBy(func(p *model.ResourceProgress) bool {
		return p.Status == model.ResourceProgressCompleted && p.CompletionPercent == 100 && p.CompletedAt != nil
	})).Return(nil)
	progressRepo.On("FindByUserAndResource", uint(10), uint(1)).Return(&model.ResourceProgress{
		ID: 1, UserID: 10, ResourceID: 1, Status: model.ResourceProgressCompleted, CompletionPercent: 100,
	}, nil)

	result, err := svc.UpsertProgress(10, 1, string(model.ResourceProgressCompleted), 100, "完了")
	assert.NoError(t, err)
	assert.Equal(t, model.ResourceProgressCompleted, result.Status)
	assert.Equal(t, 100, result.CompletionPercent)
}

func TestResourceProgressService_UpsertProgress_InProgressSetsStartedAt(t *testing.T) {
	svc, progressRepo, resourceRepo := newTestResourceProgressService()

	resourceRepo.On("FindByID", uint(1)).Return(&model.LearningResource{ID: 1, UserID: 2}, nil)
	progressRepo.On("Upsert", mock.MatchedBy(func(p *model.ResourceProgress) bool {
		return p.StartedAt != nil
	})).Return(nil)
	progressRepo.On("FindByUserAndResource", uint(10), uint(1)).Return(&model.ResourceProgress{
		ID: 1, UserID: 10, ResourceID: 1, Status: model.ResourceProgressInProgress,
	}, nil)

	_, err := svc.UpsertProgress(10, 1, string(model.ResourceProgressInProgress), 50, "")
	assert.NoError(t, err)
}

func TestResourceProgressService_UpsertProgress_InvalidStatus(t *testing.T) {
	svc, _, resourceRepo := newTestResourceProgressService()

	resourceRepo.On("FindByID", uint(1)).Return(&model.LearningResource{ID: 1, UserID: 2}, nil)

	_, err := svc.UpsertProgress(10, 1, "invalid_status", 50, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ステータス")
}

func TestResourceProgressService_UpsertProgress_InvalidPercent_Negative(t *testing.T) {
	svc, _, resourceRepo := newTestResourceProgressService()

	resourceRepo.On("FindByID", uint(1)).Return(&model.LearningResource{ID: 1, UserID: 2}, nil)

	_, err := svc.UpsertProgress(10, 1, string(model.ResourceProgressInProgress), -1, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "進捗率")
}

func TestResourceProgressService_UpsertProgress_InvalidPercent_Over100(t *testing.T) {
	svc, _, resourceRepo := newTestResourceProgressService()

	resourceRepo.On("FindByID", uint(1)).Return(&model.LearningResource{ID: 1, UserID: 2}, nil)

	_, err := svc.UpsertProgress(10, 1, string(model.ResourceProgressInProgress), 101, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "進捗率")
}

func TestResourceProgressService_UpsertProgress_ResourceNotFound(t *testing.T) {
	svc, _, resourceRepo := newTestResourceProgressService()

	resourceRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	_, err := svc.UpsertProgress(10, 999, string(model.ResourceProgressInProgress), 50, "")
	assert.Error(t, err)
}

func TestResourceProgressService_UpsertProgress_UpsertRepoError(t *testing.T) {
	svc, progressRepo, resourceRepo := newTestResourceProgressService()

	resourceRepo.On("FindByID", uint(1)).Return(&model.LearningResource{ID: 1, UserID: 2}, nil)
	progressRepo.On("Upsert", mock.Anything).Return(errors.New("db error"))

	_, err := svc.UpsertProgress(10, 1, string(model.ResourceProgressInProgress), 50, "")
	assert.Error(t, err)
}

// --- GetProgress ---

func TestResourceProgressService_GetProgress_Found(t *testing.T) {
	svc, progressRepo, _ := newTestResourceProgressService()

	now := time.Now()
	progressRepo.On("FindByUserAndResource", uint(10), uint(1)).Return(&model.ResourceProgress{
		ID: 1, UserID: 10, ResourceID: 1, Status: model.ResourceProgressInProgress,
		CompletionPercent: 50, StartedAt: &now,
	}, nil)

	result, err := svc.GetProgress(10, 1)
	assert.NoError(t, err)
	assert.Equal(t, 50, result.CompletionPercent)
}

func TestResourceProgressService_GetProgress_NotFound(t *testing.T) {
	svc, progressRepo, _ := newTestResourceProgressService()

	progressRepo.On("FindByUserAndResource", uint(10), uint(999)).Return(nil, errors.New("not found"))

	_, err := svc.GetProgress(10, 999)
	assert.Error(t, err)
}

// --- GetProgressList ---

func TestResourceProgressService_GetProgressList_WithStatusFilter(t *testing.T) {
	svc, progressRepo, _ := newTestResourceProgressService()

	progressRepo.On("FindByUserID", uint(10), "in_progress", 20, 0).Return(
		[]model.ResourceProgress{{ID: 1, UserID: 10, Status: model.ResourceProgressInProgress}},
		int64(1), nil,
	)

	list, total, err := svc.GetProgressList(10, "in_progress", 20, 0)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, int64(1), total)
}

func TestResourceProgressService_GetProgressList_NoFilter(t *testing.T) {
	svc, progressRepo, _ := newTestResourceProgressService()

	progressRepo.On("FindByUserID", uint(10), "", 20, 0).Return(
		[]model.ResourceProgress{
			{ID: 1, UserID: 10, Status: model.ResourceProgressInProgress},
			{ID: 2, UserID: 10, Status: model.ResourceProgressCompleted},
		},
		int64(2), nil,
	)

	list, total, err := svc.GetProgressList(10, "", 20, 0)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, int64(2), total)
}

func TestResourceProgressService_GetProgressList_RepoError(t *testing.T) {
	svc, progressRepo, _ := newTestResourceProgressService()

	progressRepo.On("FindByUserID", uint(10), "", 20, 0).Return(
		[]model.ResourceProgress{}, int64(0), errors.New("db error"),
	)

	_, _, err := svc.GetProgressList(10, "", 20, 0)
	assert.Error(t, err)
}

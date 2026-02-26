package service

import (
	"errors"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestProjectMilestoneService() (*ProjectMilestoneService, *MockProjectMilestoneRepository, *MockProjectRepository) {
	milestoneRepo := new(MockProjectMilestoneRepository)
	projectRepo := new(MockProjectRepository)
	svc := NewProjectMilestoneService(milestoneRepo, projectRepo)
	return svc, milestoneRepo, projectRepo
}

// --- Create ---

func TestProjectMilestoneService_Create_Success(t *testing.T) {
	svc, milestoneRepo, projectRepo := newTestProjectMilestoneService()

	projectRepo.On("FindByID", uint(1)).Return(&model.Project{ID: 1, UserID: 10}, nil)
	milestoneRepo.On("Create", mock.MatchedBy(func(m *model.ProjectMilestone) bool {
		return m.ProjectID == 1 && m.Title == "v1.0リリース"
	})).Return(nil)

	err := svc.Create(10, 1, "v1.0リリース", "初回リリース", nil)
	assert.NoError(t, err)
}

func TestProjectMilestoneService_Create_Forbidden(t *testing.T) {
	svc, _, projectRepo := newTestProjectMilestoneService()

	projectRepo.On("FindByID", uint(1)).Return(&model.Project{ID: 1, UserID: 99}, nil)

	err := svc.Create(10, 1, "v1.0リリース", "", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "権限")
}

func TestProjectMilestoneService_Create_ProjectNotFound(t *testing.T) {
	svc, _, projectRepo := newTestProjectMilestoneService()

	projectRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.Create(10, 999, "v1.0", "", nil)
	assert.Error(t, err)
}

func TestProjectMilestoneService_Create_TitleTooLong(t *testing.T) {
	svc, _, projectRepo := newTestProjectMilestoneService()

	projectRepo.On("FindByID", uint(1)).Return(&model.Project{ID: 1, UserID: 10}, nil)

	longTitle := make([]rune, 201)
	for i := range longTitle {
		longTitle[i] = 'あ'
	}
	err := svc.Create(10, 1, string(longTitle), "", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "タイトル")
}

func TestProjectMilestoneService_Create_EmptyTitle(t *testing.T) {
	svc, _, projectRepo := newTestProjectMilestoneService()

	projectRepo.On("FindByID", uint(1)).Return(&model.Project{ID: 1, UserID: 10}, nil)

	err := svc.Create(10, 1, "", "", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "タイトル")
}

func TestProjectMilestoneService_Create_DescriptionTooLong(t *testing.T) {
	svc, _, projectRepo := newTestProjectMilestoneService()

	projectRepo.On("FindByID", uint(1)).Return(&model.Project{ID: 1, UserID: 10}, nil)

	longDesc := string(make([]rune, 1001))
	err := svc.Create(10, 1, "v1.0リリース", longDesc, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "説明は1000文字以下")
}

// --- GetByProjectID ---

func TestProjectMilestoneService_GetByProjectID_Success(t *testing.T) {
	svc, milestoneRepo, _ := newTestProjectMilestoneService()

	milestoneRepo.On("FindByProjectID", uint(1)).Return([]model.ProjectMilestone{
		{ID: 1, ProjectID: 1, Title: "v1.0"},
		{ID: 2, ProjectID: 1, Title: "v2.0"},
	}, nil)

	list, err := svc.GetByProjectID(1)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

// --- Update ---

func TestProjectMilestoneService_Update_Success(t *testing.T) {
	svc, milestoneRepo, projectRepo := newTestProjectMilestoneService()

	milestoneRepo.On("FindByID", uint(1)).Return(&model.ProjectMilestone{ID: 1, ProjectID: 5, Title: "old"}, nil)
	projectRepo.On("FindByID", uint(5)).Return(&model.Project{ID: 5, UserID: 10}, nil)
	milestoneRepo.On("Update", mock.MatchedBy(func(m *model.ProjectMilestone) bool {
		return m.Title == "new title"
	})).Return(nil)

	result, err := svc.Update(10, 1, "new title", "desc", nil, "")
	assert.NoError(t, err)
	assert.Equal(t, "new title", result.Title)
}

func TestProjectMilestoneService_Update_Forbidden(t *testing.T) {
	svc, milestoneRepo, projectRepo := newTestProjectMilestoneService()

	milestoneRepo.On("FindByID", uint(1)).Return(&model.ProjectMilestone{ID: 1, ProjectID: 5}, nil)
	projectRepo.On("FindByID", uint(5)).Return(&model.Project{ID: 5, UserID: 99}, nil)

	_, err := svc.Update(10, 1, "new", "", nil, "")
	assert.Error(t, err)
}

func TestProjectMilestoneService_Update_StatusCompleted(t *testing.T) {
	svc, milestoneRepo, projectRepo := newTestProjectMilestoneService()

	milestoneRepo.On("FindByID", uint(1)).Return(&model.ProjectMilestone{ID: 1, ProjectID: 5, Title: "v1.0"}, nil)
	projectRepo.On("FindByID", uint(5)).Return(&model.Project{ID: 5, UserID: 10}, nil)
	milestoneRepo.On("Update", mock.MatchedBy(func(m *model.ProjectMilestone) bool {
		return m.Status == model.MilestoneCompleted && m.CompletedAt != nil
	})).Return(nil)

	result, err := svc.Update(10, 1, "", "", nil, string(model.MilestoneCompleted))
	assert.NoError(t, err)
	assert.Equal(t, model.MilestoneCompleted, result.Status)
	assert.NotNil(t, result.CompletedAt)
}

func TestProjectMilestoneService_Update_InvalidStatus(t *testing.T) {
	svc, milestoneRepo, projectRepo := newTestProjectMilestoneService()

	milestoneRepo.On("FindByID", uint(1)).Return(&model.ProjectMilestone{ID: 1, ProjectID: 5}, nil)
	projectRepo.On("FindByID", uint(5)).Return(&model.Project{ID: 5, UserID: 10}, nil)

	_, err := svc.Update(10, 1, "", "", nil, "invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ステータス")
}

func TestProjectMilestoneService_Update_NotFound(t *testing.T) {
	svc, milestoneRepo, _ := newTestProjectMilestoneService()

	milestoneRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	_, err := svc.Update(10, 999, "new", "", nil, "")
	assert.Error(t, err)
}

// --- Delete ---

func TestProjectMilestoneService_Delete_Success(t *testing.T) {
	svc, milestoneRepo, projectRepo := newTestProjectMilestoneService()

	milestoneRepo.On("FindByID", uint(1)).Return(&model.ProjectMilestone{ID: 1, ProjectID: 5}, nil)
	projectRepo.On("FindByID", uint(5)).Return(&model.Project{ID: 5, UserID: 10}, nil)
	milestoneRepo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(10, 1)
	assert.NoError(t, err)
}

func TestProjectMilestoneService_Delete_Forbidden(t *testing.T) {
	svc, milestoneRepo, projectRepo := newTestProjectMilestoneService()

	milestoneRepo.On("FindByID", uint(1)).Return(&model.ProjectMilestone{ID: 1, ProjectID: 5}, nil)
	projectRepo.On("FindByID", uint(5)).Return(&model.Project{ID: 5, UserID: 99}, nil)

	err := svc.Delete(10, 1)
	assert.Error(t, err)
}

func TestProjectMilestoneService_Delete_NotFound(t *testing.T) {
	svc, milestoneRepo, _ := newTestProjectMilestoneService()

	milestoneRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.Delete(10, 999)
	assert.Error(t, err)
}

// --- Create with DueDate ---

func TestProjectMilestoneService_Create_WithDueDate(t *testing.T) {
	svc, milestoneRepo, projectRepo := newTestProjectMilestoneService()

	dueDate := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	projectRepo.On("FindByID", uint(1)).Return(&model.Project{ID: 1, UserID: 10}, nil)
	milestoneRepo.On("Create", mock.MatchedBy(func(m *model.ProjectMilestone) bool {
		return m.DueDate != nil && m.DueDate.Equal(dueDate)
	})).Return(nil)

	err := svc.Create(10, 1, "v1.0", "", &dueDate)
	assert.NoError(t, err)
}

func TestProjectMilestoneService_Update_TitleTooLong(t *testing.T) {
	svc, milestoneRepo, projectRepo := newTestProjectMilestoneService()

	milestone := &model.ProjectMilestone{ID: 1, ProjectID: 5, Title: "既存"}
	milestoneRepo.On("FindByID", uint(1)).Return(milestone, nil)
	projectRepo.On("FindByID", uint(5)).Return(&model.Project{ID: 5, UserID: 10}, nil)

	longTitle := string(make([]rune, 201))
	_, err := svc.Update(10, 1, longTitle, "", nil, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "200文字")
}

func TestProjectMilestoneService_Update_DescriptionTooLong(t *testing.T) {
	svc, milestoneRepo, projectRepo := newTestProjectMilestoneService()

	milestone := &model.ProjectMilestone{ID: 1, ProjectID: 5, Title: "既存"}
	milestoneRepo.On("FindByID", uint(1)).Return(milestone, nil)
	projectRepo.On("FindByID", uint(5)).Return(&model.Project{ID: 5, UserID: 10}, nil)

	longDesc := string(make([]rune, 1001))
	_, err := svc.Update(10, 1, "", longDesc, nil, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "説明は1000文字以下")
}

func TestProjectMilestoneService_Update_StatusResetCompletedAt(t *testing.T) {
	svc, milestoneRepo, projectRepo := newTestProjectMilestoneService()

	now := time.Now()
	milestone := &model.ProjectMilestone{ID: 1, ProjectID: 5, Title: "既存", Status: model.MilestoneCompleted, CompletedAt: &now}
	milestoneRepo.On("FindByID", uint(1)).Return(milestone, nil)
	projectRepo.On("FindByID", uint(5)).Return(&model.Project{ID: 5, UserID: 10}, nil)
	milestoneRepo.On("Update", mock.MatchedBy(func(m *model.ProjectMilestone) bool {
		return m.Status == model.MilestoneInProgress && m.CompletedAt == nil
	})).Return(nil)

	result, err := svc.Update(10, 1, "", "", nil, "in_progress")
	assert.NoError(t, err)
	assert.Nil(t, result.CompletedAt)
	assert.Equal(t, model.MilestoneInProgress, result.Status)
}

func TestProjectMilestoneService_Update_RepoError(t *testing.T) {
	svc, milestoneRepo, projectRepo := newTestProjectMilestoneService()

	milestone := &model.ProjectMilestone{ID: 1, ProjectID: 5, Title: "既存"}
	milestoneRepo.On("FindByID", uint(1)).Return(milestone, nil)
	projectRepo.On("FindByID", uint(5)).Return(&model.Project{ID: 5, UserID: 10}, nil)
	milestoneRepo.On("Update", mock.Anything).Return(errors.New("db error"))

	_, err := svc.Update(10, 1, "新しいタイトル", "", nil, "")
	assert.Error(t, err)
}

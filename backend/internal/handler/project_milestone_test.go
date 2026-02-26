package handler

import (
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/mock"
)

// --- Mock ---

type MockProjectMilestoneService struct{ mock.Mock }

func (m *MockProjectMilestoneService) Create(userID, projectID uint, title, description string, dueDate *time.Time) error {
	return m.Called(userID, projectID, title, description, dueDate).Error(0)
}

func (m *MockProjectMilestoneService) GetByProjectID(projectID uint) ([]model.ProjectMilestone, error) {
	args := m.Called(projectID)
	return args.Get(0).([]model.ProjectMilestone), args.Error(1)
}

func (m *MockProjectMilestoneService) Update(userID, milestoneID uint, title, description string, dueDate *time.Time, status string) (*model.ProjectMilestone, error) {
	args := m.Called(userID, milestoneID, title, description, dueDate, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ProjectMilestone), args.Error(1)
}

func (m *MockProjectMilestoneService) Delete(userID, milestoneID uint) error {
	return m.Called(userID, milestoneID).Error(0)
}

func setupMilestoneHandler() (*ProjectMilestoneHandler, *MockProjectMilestoneService) {
	svc := new(MockProjectMilestoneService)
	h := NewProjectMilestoneHandler(svc)
	return h, svc
}

// --- Create ---

func TestProjectMilestoneHandler_Create_Success(t *testing.T) {
	h, svc := setupMilestoneHandler()

	svc.On("Create", uint(1), uint(10), "v1.0リリース", "初回リリース", (*time.Time)(nil)).Return(nil)

	r := newRouter(1)
	r.POST("/projects/:projectId/milestones", h.Create)

	w := doRequest(r, "POST", "/projects/10/milestones", map[string]interface{}{
		"title":       "v1.0リリース",
		"description": "初回リリース",
	})
	assertStatus(t, w, 201)
}

func TestProjectMilestoneHandler_Create_InvalidJSON(t *testing.T) {
	h, _ := setupMilestoneHandler()

	r := newRouter(1)
	r.POST("/projects/:projectId/milestones", h.Create)

	w := doRequestRaw(r, "POST", "/projects/10/milestones", "{invalid}")
	assertStatus(t, w, 400)
}

func TestProjectMilestoneHandler_Create_Forbidden(t *testing.T) {
	h, svc := setupMilestoneHandler()

	svc.On("Create", uint(1), uint(10), "v1.0", "", (*time.Time)(nil)).Return(service.ErrForbidden)

	r := newRouter(1)
	r.POST("/projects/:projectId/milestones", h.Create)

	w := doRequest(r, "POST", "/projects/10/milestones", map[string]interface{}{
		"title": "v1.0",
	})
	assertStatus(t, w, 403)
}

// --- GetByProjectID ---

func TestProjectMilestoneHandler_GetByProjectID_Success(t *testing.T) {
	h, svc := setupMilestoneHandler()

	svc.On("GetByProjectID", uint(10)).Return([]model.ProjectMilestone{
		{ID: 1, ProjectID: 10, Title: "v1.0"},
		{ID: 2, ProjectID: 10, Title: "v2.0"},
	}, nil)

	r := newRouter(1)
	r.GET("/projects/:projectId/milestones", h.GetByProjectID)

	w := doRequest(r, "GET", "/projects/10/milestones", nil)
	assertStatus(t, w, 200)
	data := parseJSON(t, w)
	milestones := data["milestones"].([]interface{})
	if len(milestones) != 2 {
		t.Errorf("expected 2 milestones, got %d", len(milestones))
	}
}

// --- Update ---

func TestProjectMilestoneHandler_Update_Success(t *testing.T) {
	h, svc := setupMilestoneHandler()

	svc.On("Update", uint(1), uint(5), "updated", "", (*time.Time)(nil), "in_progress").Return(&model.ProjectMilestone{
		ID: 5, Title: "updated", Status: model.MilestoneInProgress,
	}, nil)

	r := newRouter(1)
	r.PUT("/milestones/:milestoneId", h.Update)

	w := doRequest(r, "PUT", "/milestones/5", map[string]interface{}{
		"title":  "updated",
		"status": "in_progress",
	})
	assertStatus(t, w, 200)
}

func TestProjectMilestoneHandler_Update_NotFound(t *testing.T) {
	h, svc := setupMilestoneHandler()

	svc.On("Update", uint(1), uint(999), "new", "", (*time.Time)(nil), "").Return(nil, service.ErrNotFound)

	r := newRouter(1)
	r.PUT("/milestones/:milestoneId", h.Update)

	w := doRequest(r, "PUT", "/milestones/999", map[string]interface{}{
		"title": "new",
	})
	assertStatus(t, w, 404)
}

// --- Delete ---

func TestProjectMilestoneHandler_Delete_Success(t *testing.T) {
	h, svc := setupMilestoneHandler()

	svc.On("Delete", uint(1), uint(5)).Return(nil)

	r := newRouter(1)
	r.DELETE("/milestones/:milestoneId", h.Delete)

	w := doRequest(r, "DELETE", "/milestones/5", nil)
	assertStatus(t, w, 200)
}

func TestProjectMilestoneHandler_Delete_Forbidden(t *testing.T) {
	h, svc := setupMilestoneHandler()

	svc.On("Delete", uint(1), uint(5)).Return(service.ErrForbidden)

	r := newRouter(1)
	r.DELETE("/milestones/:milestoneId", h.Delete)

	w := doRequest(r, "DELETE", "/milestones/5", nil)
	assertStatus(t, w, 403)
}

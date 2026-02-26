package handler

import (
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/mock"
)

// --- Mock ---

type MockResourceProgressService struct{ mock.Mock }

func (m *MockResourceProgressService) UpsertProgress(userID, resourceID uint, status string, completionPercent int, note string) (*model.ResourceProgress, error) {
	args := m.Called(userID, resourceID, status, completionPercent, note)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ResourceProgress), args.Error(1)
}

func (m *MockResourceProgressService) GetProgress(userID, resourceID uint) (*model.ResourceProgress, error) {
	args := m.Called(userID, resourceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ResourceProgress), args.Error(1)
}

func (m *MockResourceProgressService) GetProgressList(userID uint, status string, limit, offset int) ([]model.ResourceProgress, int64, error) {
	args := m.Called(userID, status, limit, offset)
	return args.Get(0).([]model.ResourceProgress), args.Get(1).(int64), args.Error(2)
}

func setupResourceProgressHandler() (*ResourceProgressHandler, *MockResourceProgressService) {
	svc := new(MockResourceProgressService)
	h := NewResourceProgressHandler(svc)
	return h, svc
}

// --- Upsert ---

func TestResourceProgressHandler_Upsert_Success(t *testing.T) {
	h, svc := setupResourceProgressHandler()

	svc.On("UpsertProgress", uint(1), uint(10), "in_progress", 50, "学習中").Return(&model.ResourceProgress{
		ID: 1, UserID: 1, ResourceID: 10, Status: model.ResourceProgressInProgress, CompletionPercent: 50,
	}, nil)

	r := newRouter(1)
	r.PUT("/resources/progress", h.Upsert)

	w := doRequest(r, "PUT", "/resources/progress", map[string]interface{}{
		"resource_id":        10,
		"status":             "in_progress",
		"completion_percent": 50,
		"note":               "学習中",
	})

	assertStatus(t, w, 200)
	data := parseJSON(t, w)
	progress := data["progress"].(map[string]interface{})
	if progress["completion_percent"].(float64) != 50 {
		t.Errorf("expected completion_percent 50, got %v", progress["completion_percent"])
	}
}

func TestResourceProgressHandler_Upsert_InvalidJSON(t *testing.T) {
	h, _ := setupResourceProgressHandler()

	r := newRouter(1)
	r.PUT("/resources/progress", h.Upsert)

	w := doRequestRaw(r, "PUT", "/resources/progress", "{invalid}")
	assertStatus(t, w, 400)
}

func TestResourceProgressHandler_Upsert_ServiceError(t *testing.T) {
	h, svc := setupResourceProgressHandler()

	svc.On("UpsertProgress", uint(1), uint(10), "invalid", 50, "").Return(nil, service.ErrBadRequest)

	r := newRouter(1)
	r.PUT("/resources/progress", h.Upsert)

	w := doRequest(r, "PUT", "/resources/progress", map[string]interface{}{
		"resource_id":        10,
		"status":             "invalid",
		"completion_percent": 50,
	})

	assertStatus(t, w, 400)
}

// --- GetByResource ---

func TestResourceProgressHandler_GetByResource_Success(t *testing.T) {
	h, svc := setupResourceProgressHandler()

	svc.On("GetProgress", uint(1), uint(10)).Return(&model.ResourceProgress{
		ID: 1, UserID: 1, ResourceID: 10, Status: model.ResourceProgressInProgress, CompletionPercent: 75,
	}, nil)

	r := newRouter(1)
	r.GET("/resources/:resourceId/progress", h.GetByResource)

	w := doRequest(r, "GET", "/resources/10/progress", nil)
	assertStatus(t, w, 200)
}

func TestResourceProgressHandler_GetByResource_NotFound(t *testing.T) {
	h, svc := setupResourceProgressHandler()

	svc.On("GetProgress", uint(1), uint(999)).Return(nil, service.ErrNotFound)

	r := newRouter(1)
	r.GET("/resources/:resourceId/progress", h.GetByResource)

	w := doRequest(r, "GET", "/resources/999/progress", nil)
	assertStatus(t, w, 404)
}

// --- GetMyProgress ---

func TestResourceProgressHandler_GetMyProgress_Success(t *testing.T) {
	h, svc := setupResourceProgressHandler()

	svc.On("GetProgressList", uint(1), "", 20, 0).Return(
		[]model.ResourceProgress{
			{ID: 1, UserID: 1, Status: model.ResourceProgressInProgress},
			{ID: 2, UserID: 1, Status: model.ResourceProgressCompleted},
		},
		int64(2), nil,
	)

	r := newRouter(1)
	r.GET("/resources/progress", h.GetMyProgress)

	w := doRequest(r, "GET", "/resources/progress", nil)
	assertStatus(t, w, 200)
	data := parseJSON(t, w)
	progresses := data["progresses"].([]interface{})
	if len(progresses) != 2 {
		t.Errorf("expected 2 progresses, got %d", len(progresses))
	}
}

func TestResourceProgressHandler_GetMyProgress_WithStatusFilter(t *testing.T) {
	h, svc := setupResourceProgressHandler()

	svc.On("GetProgressList", uint(1), "completed", 20, 0).Return(
		[]model.ResourceProgress{{ID: 2, UserID: 1, Status: model.ResourceProgressCompleted}},
		int64(1), nil,
	)

	r := newRouter(1)
	r.GET("/resources/progress", h.GetMyProgress)

	w := doRequest(r, "GET", "/resources/progress?status=completed", nil)
	assertStatus(t, w, 200)
}

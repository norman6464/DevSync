package handler

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/mock"
)

// mockResourceProgressRepo は usecase/repository.ResourceProgressRepository のモック（ctx 付き）。
// （LearningResourceReader のモックは resource_review のテストで定義済みのものを再利用する）
type mockResourceProgressRepo struct{ mock.Mock }

func (m *mockResourceProgressRepo) Upsert(ctx context.Context, progress *model.ResourceProgress) error {
	return m.Called(ctx, progress).Error(0)
}
func (m *mockResourceProgressRepo) FindByUserAndResource(ctx context.Context, userID, resourceID uint) (*model.ResourceProgress, error) {
	args := m.Called(ctx, userID, resourceID)
	p, _ := args.Get(0).(*model.ResourceProgress)
	return p, args.Error(1)
}
func (m *mockResourceProgressRepo) FindByUserID(ctx context.Context, userID uint, status string, limit, offset int) ([]model.ResourceProgress, int64, error) {
	args := m.Called(ctx, userID, status, limit, offset)
	list, _ := args.Get(0).([]model.ResourceProgress)
	return list, args.Get(1).(int64), args.Error(2)
}

// setupResourceProgressHandler は本物の usecase + port モックで ResourceProgressHandler を組む。
func setupResourceProgressHandler() (*ResourceProgressHandler, *mockResourceProgressRepo, *mockLearningResourceReader) {
	progress := new(mockResourceProgressRepo)
	resources := new(mockLearningResourceReader)
	h := NewResourceProgressHandler(
		usecase.NewUpsertResourceProgressUseCase(progress, resources),
		usecase.NewGetResourceProgressUseCase(progress),
		usecase.NewListResourceProgressUseCase(progress),
	)
	return h, progress, resources
}

// --- Upsert ---

func TestResourceProgressHandler_Upsert_Success(t *testing.T) {
	h, progress, resources := setupResourceProgressHandler()
	resources.On("FindByID", mock.Anything, uint(10)).Return(&model.LearningResource{}, nil)
	progress.On("Upsert", mock.Anything, mock.AnythingOfType("*model.ResourceProgress")).Return(nil)
	progress.On("FindByUserAndResource", mock.Anything, uint(1), uint(10)).Return(&model.ResourceProgress{
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
	prog := data["progress"].(map[string]interface{})
	if prog["completion_percent"].(float64) != 50 {
		t.Errorf("expected completion_percent 50, got %v", prog["completion_percent"])
	}
	progress.AssertExpectations(t)
}

func TestResourceProgressHandler_Upsert_InvalidJSON(t *testing.T) {
	h, _, _ := setupResourceProgressHandler()

	r := newRouter(1)
	r.PUT("/resources/progress", h.Upsert)

	w := doRequestRaw(r, "PUT", "/resources/progress", "{invalid}")
	assertStatus(t, w, 400)
}

func TestResourceProgressHandler_Upsert_InvalidStatus(t *testing.T) {
	h, progress, resources := setupResourceProgressHandler()
	resources.On("FindByID", mock.Anything, uint(10)).Return(&model.LearningResource{}, nil)

	r := newRouter(1)
	r.PUT("/resources/progress", h.Upsert)

	w := doRequest(r, "PUT", "/resources/progress", map[string]interface{}{
		"resource_id":        10,
		"status":             "invalid",
		"completion_percent": 50,
	})

	assertStatus(t, w, 400)
	progress.AssertNotCalled(t, "Upsert")
}

// --- GetByResource ---

func TestResourceProgressHandler_GetByResource_Success(t *testing.T) {
	h, progress, _ := setupResourceProgressHandler()
	progress.On("FindByUserAndResource", mock.Anything, uint(1), uint(10)).Return(&model.ResourceProgress{
		ID: 1, UserID: 1, ResourceID: 10, Status: model.ResourceProgressInProgress, CompletionPercent: 75,
	}, nil)

	r := newRouter(1)
	r.GET("/resources/:resourceId/progress", h.GetByResource)

	w := doRequest(r, "GET", "/resources/10/progress", nil)
	assertStatus(t, w, 200)
	progress.AssertExpectations(t)
}

func TestResourceProgressHandler_GetByResource_NotFound(t *testing.T) {
	h, progress, _ := setupResourceProgressHandler()
	progress.On("FindByUserAndResource", mock.Anything, uint(1), uint(999)).Return((*model.ResourceProgress)(nil), domain.ErrNotFound)

	r := newRouter(1)
	r.GET("/resources/:resourceId/progress", h.GetByResource)

	w := doRequest(r, "GET", "/resources/999/progress", nil)
	assertStatus(t, w, 404)
	progress.AssertExpectations(t)
}

// --- GetMyProgress ---

func TestResourceProgressHandler_GetMyProgress_Success(t *testing.T) {
	h, progress, _ := setupResourceProgressHandler()
	progress.On("FindByUserID", mock.Anything, uint(1), "", 20, 0).Return(
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
	progress.AssertExpectations(t)
}

func TestResourceProgressHandler_GetMyProgress_WithStatusFilter(t *testing.T) {
	h, progress, _ := setupResourceProgressHandler()
	progress.On("FindByUserID", mock.Anything, uint(1), "completed", 20, 0).Return(
		[]model.ResourceProgress{{ID: 2, UserID: 1, Status: model.ResourceProgressCompleted}},
		int64(1), nil,
	)

	r := newRouter(1)
	r.GET("/resources/progress", h.GetMyProgress)

	w := doRequest(r, "GET", "/resources/progress?status=completed", nil)
	assertStatus(t, w, 200)
	progress.AssertExpectations(t)
}

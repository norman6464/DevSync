package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------- Create ----------

func TestRoadmapCreate_Success(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps", h.Create)

	repo.On("Create", mock.AnythingOfType("*model.Roadmap")).Return(nil)

	w := doRequest(r, http.MethodPost, "/roadmaps", map[string]interface{}{
		"title": "Learn Go", "category": "language", "is_public": true,
	})
	assertStatus(t, w, http.StatusCreated)
}

func TestRoadmapCreate_ValidationError(t *testing.T) {
	h, _ := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps", h.Create)

	// title は required
	w := doRequest(r, http.MethodPost, "/roadmaps", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestRoadmapCreate_DefaultCategory(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps", h.Create)

	repo.On("Create", mock.MatchedBy(func(rm *model.Roadmap) bool {
		return rm.Category == model.RoadmapCategoryOther
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/roadmaps", map[string]string{
		"title": "No Category",
	})
	assertStatus(t, w, http.StatusCreated)
}

// ---------- GetMyRoadmaps ----------

func TestRoadmapGetMy_Success(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/my", h.GetMyRoadmaps)

	repo.On("GetByUserID", uint(1), 20, 0).Return([]model.Roadmap{
		{Title: "My Roadmap"},
	}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/roadmaps/my", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// ---------- GetPublicRoadmaps ----------

func TestRoadmapGetPublic_Success(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/public", h.GetPublicRoadmaps)

	repo.On("GetPublicRoadmaps", 20, 0).Return(
		[]model.Roadmap{{Title: "Public"}}, int64(1), nil,
	)

	w := doRequest(r, http.MethodGet, "/roadmaps/public", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- GetByID ----------

func TestRoadmapGetByID_Success(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/:id", h.GetByID)

	rm := &model.Roadmap{Title: "Found", IsPublic: true}
	rm.ID = 10
	rm.UserID = 1
	repo.On("FindByID", uint(10)).Return(rm, nil)

	w := doRequest(r, http.MethodGet, "/roadmaps/10", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestRoadmapGetByID_NotFound(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/:id", h.GetByID)

	repo.On("FindByID", uint(999)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/roadmaps/999", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestRoadmapGetByID_ForbiddenPrivate(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/:id", h.GetByID)

	// 非公開で他ユーザーのロードマップ
	rm := &model.Roadmap{Title: "Private", IsPublic: false}
	rm.ID = 10
	rm.UserID = 999
	repo.On("FindByID", uint(10)).Return(rm, nil)

	w := doRequest(r, http.MethodGet, "/roadmaps/10", nil)
	assertStatus(t, w, http.StatusForbidden)
}

func TestRoadmapGetByID_InvalidID(t *testing.T) {
	h, _ := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/roadmaps/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- Update ----------

func TestRoadmapUpdate_Success(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id", h.Update)

	rm := &model.Roadmap{Title: "Old"}
	rm.ID = 10
	rm.UserID = 1
	repo.On("FindByID", uint(10)).Return(rm, nil)
	repo.On("Update", mock.AnythingOfType("*model.Roadmap")).Return(nil)

	title := "Updated"
	w := doRequest(r, http.MethodPut, "/roadmaps/10", map[string]interface{}{
		"title": &title,
	})
	assertStatus(t, w, http.StatusOK)
}

func TestRoadmapUpdate_Forbidden(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id", h.Update)

	rm := &model.Roadmap{Title: "Other"}
	rm.ID = 10
	rm.UserID = 999
	repo.On("FindByID", uint(10)).Return(rm, nil)

	w := doRequest(r, http.MethodPut, "/roadmaps/10", map[string]string{"title": "X"})
	assertStatus(t, w, http.StatusForbidden)
}

func TestRoadmapUpdate_NotFound(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id", h.Update)

	repo.On("FindByID", uint(10)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodPut, "/roadmaps/10", map[string]string{"title": "X"})
	assertStatus(t, w, http.StatusNotFound)
}

// ---------- Delete ----------

func TestRoadmapDelete_Success(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.DELETE("/roadmaps/:id", h.Delete)

	rm := &model.Roadmap{}
	rm.ID = 10
	rm.UserID = 1
	repo.On("FindByID", uint(10)).Return(rm, nil)
	repo.On("Delete", uint(10)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/roadmaps/10", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestRoadmapDelete_Forbidden(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.DELETE("/roadmaps/:id", h.Delete)

	rm := &model.Roadmap{}
	rm.ID = 10
	rm.UserID = 999
	repo.On("FindByID", uint(10)).Return(rm, nil)

	w := doRequest(r, http.MethodDelete, "/roadmaps/10", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- CopyRoadmap ----------

func TestRoadmapCopy_Success(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps/:id/copy", h.CopyRoadmap)

	original := &model.Roadmap{Title: "Template", IsPublic: true}
	original.ID = 10
	repo.On("FindByID", uint(10)).Return(original, nil)

	copied := &model.Roadmap{Title: "Template"}
	copied.ID = 20
	copied.UserID = 1
	repo.On("CopyRoadmap", uint(10), uint(1)).Return(copied, nil)

	w := doRequest(r, http.MethodPost, "/roadmaps/10/copy", nil)
	assertStatus(t, w, http.StatusCreated)
}

func TestRoadmapCopy_ForbiddenPrivate(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps/:id/copy", h.CopyRoadmap)

	original := &model.Roadmap{Title: "Private", IsPublic: false}
	original.ID = 10
	original.UserID = 999
	repo.On("FindByID", uint(10)).Return(original, nil)

	w := doRequest(r, http.MethodPost, "/roadmaps/10/copy", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- GetTemplates ----------

func TestRoadmapGetTemplates_Success(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/templates", h.GetTemplates)

	repo.On("GetTemplates").Return([]model.Roadmap{
		{Title: "Template 1", IsTemplate: true},
	}, nil)

	w := doRequest(r, http.MethodGet, "/roadmaps/templates", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- CreateStep ----------

func TestRoadmapCreateStep_Success(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps/:id/steps", h.CreateStep)

	rm := &model.Roadmap{Title: "My Roadmap"}
	rm.ID = 10
	rm.UserID = 1
	repo.On("FindByID", uint(10)).Return(rm, nil)
	repo.On("CreateStep", mock.AnythingOfType("*model.RoadmapStep")).Return(nil)

	w := doRequest(r, http.MethodPost, "/roadmaps/10/steps", map[string]string{
		"title": "Step 1",
	})
	assertStatus(t, w, http.StatusCreated)
}

func TestRoadmapCreateStep_Forbidden(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps/:id/steps", h.CreateStep)

	rm := &model.Roadmap{Title: "Other"}
	rm.ID = 10
	rm.UserID = 999
	repo.On("FindByID", uint(10)).Return(rm, nil)

	w := doRequest(r, http.MethodPost, "/roadmaps/10/steps", map[string]string{
		"title": "Step 1",
	})
	assertStatus(t, w, http.StatusForbidden)
}

func TestRoadmapCreateStep_ValidationError(t *testing.T) {
	h, _ := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps/:id/steps", h.CreateStep)

	// title は required
	w := doRequest(r, http.MethodPost, "/roadmaps/10/steps", map[string]string{})
	// Note: roadmap 10 の FindByID がモックされていないが、ShouldBindJSON が先に失敗する
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- UpdateStep ----------

func TestRoadmapUpdateStep_Success(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id/steps/:stepId", h.UpdateStep)

	rm := &model.Roadmap{Title: "My Roadmap"}
	rm.ID = 10
	rm.UserID = 1
	repo.On("FindByID", uint(10)).Return(rm, nil)

	step := &model.RoadmapStep{Title: "Old Step", RoadmapID: 10}
	step.ID = 5
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("UpdateStep", mock.AnythingOfType("*model.RoadmapStep")).Return(nil)

	title := "Updated Step"
	w := doRequest(r, http.MethodPut, "/roadmaps/10/steps/5", map[string]interface{}{
		"title": &title,
	})
	assertStatus(t, w, http.StatusOK)
}

func TestRoadmapUpdateStep_Forbidden(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id/steps/:stepId", h.UpdateStep)

	rm := &model.Roadmap{Title: "Other"}
	rm.ID = 10
	rm.UserID = 999
	repo.On("FindByID", uint(10)).Return(rm, nil)

	w := doRequest(r, http.MethodPut, "/roadmaps/10/steps/5", map[string]string{"title": "X"})
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- DeleteStep ----------

func TestRoadmapDeleteStep_Success(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.DELETE("/roadmaps/:id/steps/:stepId", h.DeleteStep)

	rm := &model.Roadmap{Title: "My Roadmap"}
	rm.ID = 10
	rm.UserID = 1
	repo.On("FindByID", uint(10)).Return(rm, nil)

	step := &model.RoadmapStep{Title: "Step", RoadmapID: 10}
	step.ID = 5
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("DeleteStep", uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/roadmaps/10/steps/5", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestRoadmapDeleteStep_Forbidden(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.DELETE("/roadmaps/:id/steps/:stepId", h.DeleteStep)

	rm := &model.Roadmap{Title: "Other"}
	rm.ID = 10
	rm.UserID = 999
	repo.On("FindByID", uint(10)).Return(rm, nil)

	w := doRequest(r, http.MethodDelete, "/roadmaps/10/steps/5", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- ReorderSteps ----------

func TestRoadmapReorderSteps_Success(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id/steps/reorder", h.ReorderSteps)

	rm := &model.Roadmap{Title: "My Roadmap"}
	rm.ID = 10
	rm.UserID = 1
	repo.On("FindByID", uint(10)).Return(rm, nil)
	repo.On("ReorderSteps", uint(10), mock.AnythingOfType("[]model.StepOrder")).Return(nil)

	w := doRequest(r, http.MethodPut, "/roadmaps/10/steps/reorder", map[string]interface{}{
		"orders": []map[string]interface{}{
			{"step_id": 1, "order_index": 0},
			{"step_id": 2, "order_index": 1},
		},
	})
	assertStatus(t, w, http.StatusOK)
}

func TestRoadmapReorderSteps_Forbidden(t *testing.T) {
	h, repo := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id/steps/reorder", h.ReorderSteps)

	rm := &model.Roadmap{Title: "Other"}
	rm.ID = 10
	rm.UserID = 999
	repo.On("FindByID", uint(10)).Return(rm, nil)

	w := doRequest(r, http.MethodPut, "/roadmaps/10/steps/reorder", map[string]interface{}{
		"orders": []map[string]interface{}{
			{"step_id": 1, "order_index": 0},
		},
	})
	assertStatus(t, w, http.StatusForbidden)
}

// ============================================================
// GetByStatus テスト
// ============================================================

func TestRoadmap_GetByStatus_Success(t *testing.T) {
	h, svc := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.GET("/roadmaps/status/:status", h.GetByStatus)

	roadmaps := []model.Roadmap{{Title: "Go入門"}}
	svc.On("GetByStatus", uint(1), "active").Return(roadmaps, nil)

	w := doRequest(r, http.MethodGet, "/roadmaps/status/active", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestRoadmap_GetByStatus_NilResult(t *testing.T) {
	h, svc := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.GET("/roadmaps/status/:status", h.GetByStatus)

	svc.On("GetByStatus", uint(1), "completed").Return([]model.Roadmap(nil), nil)

	w := doRequest(r, http.MethodGet, "/roadmaps/status/completed", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
	svc.AssertExpectations(t)
}

func TestRoadmap_GetByStatus_ServiceError(t *testing.T) {
	h, svc := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.GET("/roadmaps/status/:status", h.GetByStatus)

	svc.On("GetByStatus", uint(1), "invalid").Return([]model.Roadmap(nil), service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/roadmaps/status/invalid", nil)
	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

func TestRoadmapReorderSteps_ValidationError(t *testing.T) {
	h, _ := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id/steps/reorder", h.ReorderSteps)

	// orders は required
	w := doRequest(r, http.MethodPut, "/roadmaps/10/steps/reorder", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- CreateFromTemplate ----------

func TestRoadmapCreateFromTemplate_Success(t *testing.T) {
	h, svc := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.POST("/roadmaps/templates/:id/create", h.CreateFromTemplate)

	roadmap := &model.Roadmap{Title: "From Template"}
	roadmap.ID = 10
	svc.On("CreateFromTemplate", uint(5), uint(1)).Return(roadmap, nil)

	w := doRequest(r, http.MethodPost, "/roadmaps/templates/5/create", nil)
	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
}

func TestRoadmapCreateFromTemplate_InvalidID(t *testing.T) {
	h, _ := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.POST("/roadmaps/templates/:id/create", h.CreateFromTemplate)

	w := doRequest(r, http.MethodPost, "/roadmaps/templates/abc/create", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestRoadmapCreateFromTemplate_ServiceError(t *testing.T) {
	h, svc := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.POST("/roadmaps/templates/:id/create", h.CreateFromTemplate)

	svc.On("CreateFromTemplate", uint(5), uint(1)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodPost, "/roadmaps/templates/5/create", nil)
	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

// ---------- GetTemplates エラーパス ----------

func TestRoadmapGetTemplates_ServiceError(t *testing.T) {
	h, svc := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.GET("/roadmaps/templates", h.GetTemplates)

	svc.On("GetTemplates").Return([]model.Roadmap(nil), assert.AnError)

	w := doRequest(r, http.MethodGet, "/roadmaps/templates", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ---------- GetMyRoadmaps エラーパス ----------

func TestRoadmapGetMy_ServiceError(t *testing.T) {
	h, svc := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.GET("/roadmaps", h.GetMyRoadmaps)

	svc.On("GetByUserID", uint(1), 20, 0).Return([]model.Roadmap(nil), int64(0), assert.AnError)

	w := doRequest(r, http.MethodGet, "/roadmaps", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ---------- GetPublicRoadmaps エラーパス ----------

func TestRoadmapGetPublic_ServiceError(t *testing.T) {
	h, svc := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.GET("/roadmaps/public", h.GetPublicRoadmaps)

	svc.On("GetPublicRoadmaps", 20, 0).Return([]model.Roadmap(nil), int64(0), assert.AnError)

	w := doRequest(r, http.MethodGet, "/roadmaps/public", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ---------- UpdateStep (完了ステータスのみ) ----------

func TestRoadmapUpdateStep_CompletionOnly(t *testing.T) {
	h, svc := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.PUT("/roadmaps/:id/steps/:stepId", h.UpdateStep)

	completed := true
	step := &model.RoadmapStep{Title: "Step1"}
	step.ID = 2
	svc.On("UpdateStepCompletion", uint(1), uint(2), uint(1), completed).Return(step, nil)

	w := doRequest(r, http.MethodPut, "/roadmaps/1/steps/2", map[string]interface{}{
		"is_completed": true,
	})
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestRoadmapUpdateStep_CompletionError(t *testing.T) {
	h, svc := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.PUT("/roadmaps/:id/steps/:stepId", h.UpdateStep)

	svc.On("UpdateStepCompletion", uint(1), uint(2), uint(1), true).Return(nil, service.ErrForbidden)

	w := doRequest(r, http.MethodPut, "/roadmaps/1/steps/2", map[string]interface{}{
		"is_completed": true,
	})
	assertStatus(t, w, http.StatusForbidden)
	svc.AssertExpectations(t)
}

// ---------- UpdateStep InvalidStepID ----------

func TestRoadmapUpdateStep_InvalidStepID(t *testing.T) {
	h, _ := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.PUT("/roadmaps/:id/steps/:stepId", h.UpdateStep)

	w := doRequest(r, http.MethodPut, "/roadmaps/1/steps/abc", map[string]interface{}{
		"title": "test",
	})
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- DeleteStep InvalidStepID ----------

func TestRoadmapDeleteStep_InvalidStepID(t *testing.T) {
	h, _ := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.DELETE("/roadmaps/:id/steps/:stepId", h.DeleteStep)

	w := doRequest(r, http.MethodDelete, "/roadmaps/1/steps/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- CopyRoadmap InvalidID ----------

func TestRoadmapCopy_InvalidID(t *testing.T) {
	h, _ := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.POST("/roadmaps/:id/copy", h.CopyRoadmap)

	w := doRequest(r, http.MethodPost, "/roadmaps/abc/copy", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- CopyRoadmap ServiceError ----------

func TestRoadmapCopy_ServiceError(t *testing.T) {
	h, svc := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.POST("/roadmaps/:id/copy", h.CopyRoadmap)

	svc.On("CopyRoadmap", uint(5), uint(1)).Return(nil, assert.AnError)

	w := doRequest(r, http.MethodPost, "/roadmaps/5/copy", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ---------- ReorderSteps InvalidID ----------

func TestRoadmapReorderSteps_InvalidID(t *testing.T) {
	h, _ := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.PUT("/roadmaps/:id/steps/reorder", h.ReorderSteps)

	w := doRequest(r, http.MethodPut, "/roadmaps/abc/steps/reorder", map[string]interface{}{
		"orders": []map[string]interface{}{{"step_id": 1, "order_index": 0}},
	})
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// GetMyStats テスト
// ============================================================

func TestRoadmapGetMyStats_Success(t *testing.T) {
	h, svc := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.GET("/roadmaps/my/stats", h.GetMyStats)

	stats := &model.RoadmapStats{TotalRoadmaps: 3, ActiveRoadmaps: 2, CompletedRoadmaps: 1, TotalSteps: 10, CompletedSteps: 5}
	svc.On("GetStats", uint(1)).Return(stats, nil)

	w := doRequest(r, http.MethodGet, "/roadmaps/my/stats", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestRoadmapGetMyStats_ServiceError(t *testing.T) {
	h, svc := setupRoadmapHandlerMock()
	r := newRouter(1)
	r.GET("/roadmaps/my/stats", h.GetMyStats)

	svc.On("GetStats", uint(1)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/roadmaps/my/stats", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

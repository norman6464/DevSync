package handler

import (
	"encoding/json"
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

	repo.On("GetByUserID", uint(1)).Return([]model.Roadmap{
		{Title: "My Roadmap"},
	}, nil)

	w := doRequest(r, http.MethodGet, "/roadmaps/my", nil)
	assertStatus(t, w, http.StatusOK)

	var roadmaps []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &roadmaps)
	assert.Len(t, roadmaps, 1)
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

func TestRoadmapReorderSteps_ValidationError(t *testing.T) {
	h, _ := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id/steps/reorder", h.ReorderSteps)

	// orders は required
	w := doRequest(r, http.MethodPut, "/roadmaps/10/steps/reorder", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------- Create ----------

func TestRoadmapCreate_Success(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps", h.Create)

	ports.Roadmaps.On("Create", mock.Anything, mock.AnythingOfType("*model.Roadmap")).Return(nil)

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
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps", h.Create)

	ports.Roadmaps.On("Create", mock.Anything, mock.MatchedBy(func(rm *model.Roadmap) bool {
		return rm.Category == model.RoadmapCategoryOther
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/roadmaps", map[string]string{
		"title": "No Category",
	})
	assertStatus(t, w, http.StatusCreated)
}

// ---------- GetMyRoadmaps ----------

func TestRoadmapGetMy_Success(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/my", h.GetMyRoadmaps)

	ports.Roadmaps.On("GetByUserID", mock.Anything, uint(1), 20, 0).Return([]model.Roadmap{
		{Title: "My Roadmap"},
	}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/roadmaps/my", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Roadmaps.AssertExpectations(t)
}

// ---------- GetPublicRoadmaps ----------

func TestRoadmapGetPublic_Success(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/public", h.GetPublicRoadmaps)

	ports.Roadmaps.On("GetPublicRoadmaps", mock.Anything, 20, 0).Return(
		[]model.Roadmap{{Title: "Public"}}, int64(1), nil,
	)

	w := doRequest(r, http.MethodGet, "/roadmaps/public", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- GetByID ----------

func TestRoadmapGetByID_Success(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/:id", h.GetByID)

	rm := &model.Roadmap{Title: "Found", IsPublic: true}
	rm.ID = 10
	rm.UserID = 1
	ports.Roadmaps.On("FindByID", mock.Anything, uint(10)).Return(rm, nil)

	w := doRequest(r, http.MethodGet, "/roadmaps/10", nil)
	assertStatus(t, w, http.StatusOK)
}

// 不在のロードマップは 404 にならず 500 になる（移行前からの挙動）。
func TestRoadmapGetByID_MissingReturnsInternalError(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/:id", h.GetByID)

	// port は不在を (nil, nil) で表す。
	ports.Roadmaps.On("FindByID", mock.Anything, uint(999)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/roadmaps/999", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestRoadmapGetByID_ForbiddenPrivate(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/:id", h.GetByID)

	// 非公開で他ユーザーのロードマップ
	rm := &model.Roadmap{Title: "Private", IsPublic: false}
	rm.ID = 10
	rm.UserID = 999
	ports.Roadmaps.On("FindByID", mock.Anything, uint(10)).Return(rm, nil)

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
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id", h.Update)

	rm := &model.Roadmap{Title: "Old"}
	rm.ID = 10
	rm.UserID = 1
	ports.Roadmaps.On("FindByID", mock.Anything, uint(10)).Return(rm, nil)
	ports.Roadmaps.On("Update", mock.Anything, mock.AnythingOfType("*model.Roadmap")).Return(nil)

	title := "Updated"
	w := doRequest(r, http.MethodPut, "/roadmaps/10", map[string]interface{}{
		"title": &title,
	})
	assertStatus(t, w, http.StatusOK)
}

func TestRoadmapUpdate_Forbidden(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id", h.Update)

	rm := &model.Roadmap{Title: "Other"}
	rm.ID = 10
	rm.UserID = 999
	ports.Roadmaps.On("FindByID", mock.Anything, uint(10)).Return(rm, nil)

	w := doRequest(r, http.MethodPut, "/roadmaps/10", map[string]string{"title": "X"})
	assertStatus(t, w, http.StatusForbidden)
}

// 更新も不在は 404 にならず 500 になる（移行前からの挙動）。
func TestRoadmapUpdate_MissingReturnsInternalError(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id", h.Update)

	ports.Roadmaps.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)

	w := doRequest(r, http.MethodPut, "/roadmaps/10", map[string]string{"title": "X"})
	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- Delete ----------

func TestRoadmapDelete_Success(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.DELETE("/roadmaps/:id", h.Delete)

	rm := &model.Roadmap{}
	rm.ID = 10
	rm.UserID = 1
	ports.Roadmaps.On("FindByID", mock.Anything, uint(10)).Return(rm, nil)
	ports.Roadmaps.On("Delete", mock.Anything, uint(10)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/roadmaps/10", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestRoadmapDelete_Forbidden(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.DELETE("/roadmaps/:id", h.Delete)

	rm := &model.Roadmap{}
	rm.ID = 10
	rm.UserID = 999
	ports.Roadmaps.On("FindByID", mock.Anything, uint(10)).Return(rm, nil)

	w := doRequest(r, http.MethodDelete, "/roadmaps/10", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- CopyRoadmap ----------

func TestRoadmapCopy_Success(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps/:id/copy", h.CopyRoadmap)

	original := &model.Roadmap{Title: "Template", IsPublic: true}
	original.ID = 10
	ports.Roadmaps.On("FindByID", mock.Anything, uint(10)).Return(original, nil)

	copied := &model.Roadmap{Title: "Template"}
	copied.ID = 20
	copied.UserID = 1
	ports.Roadmaps.On("CopyRoadmap", mock.Anything, uint(10), uint(1)).Return(copied, nil)

	w := doRequest(r, http.MethodPost, "/roadmaps/10/copy", nil)
	assertStatus(t, w, http.StatusCreated)
}

func TestRoadmapCopy_ForbiddenPrivate(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps/:id/copy", h.CopyRoadmap)

	original := &model.Roadmap{Title: "Private", IsPublic: false}
	original.ID = 10
	original.UserID = 999
	ports.Roadmaps.On("FindByID", mock.Anything, uint(10)).Return(original, nil)

	w := doRequest(r, http.MethodPost, "/roadmaps/10/copy", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- GetTemplates ----------

func TestRoadmapGetTemplates_Success(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/templates", h.GetTemplates)

	ports.Roadmaps.On("GetTemplates", mock.Anything).Return([]model.Roadmap{
		{Title: "Template 1", IsTemplate: true},
	}, nil)

	w := doRequest(r, http.MethodGet, "/roadmaps/templates", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- CreateStep ----------

func TestRoadmapCreateStep_Success(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps/:id/steps", h.CreateStep)

	rm := &model.Roadmap{Title: "My Roadmap"}
	rm.ID = 10
	rm.UserID = 1
	ports.Roadmaps.On("FindByID", mock.Anything, uint(10)).Return(rm, nil)
	ports.Roadmaps.On("CreateStep", mock.Anything, mock.AnythingOfType("*model.RoadmapStep")).Return(nil)

	w := doRequest(r, http.MethodPost, "/roadmaps/10/steps", map[string]string{
		"title": "Step 1",
	})
	assertStatus(t, w, http.StatusCreated)
}

func TestRoadmapCreateStep_Forbidden(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps/:id/steps", h.CreateStep)

	rm := &model.Roadmap{Title: "Other"}
	rm.ID = 10
	rm.UserID = 999
	ports.Roadmaps.On("FindByID", mock.Anything, uint(10)).Return(rm, nil)

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
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id/steps/:stepId", h.UpdateStep)

	rm := &model.Roadmap{Title: "My Roadmap"}
	rm.ID = 10
	rm.UserID = 1
	ports.Roadmaps.On("FindByID", mock.Anything, uint(10)).Return(rm, nil)

	step := &model.RoadmapStep{Title: "Old Step", RoadmapID: 10}
	step.ID = 5
	ports.Roadmaps.On("FindStepByID", mock.Anything, uint(5)).Return(step, nil)
	ports.Roadmaps.On("UpdateStep", mock.Anything, mock.AnythingOfType("*model.RoadmapStep")).Return(nil)

	title := "Updated Step"
	w := doRequest(r, http.MethodPut, "/roadmaps/10/steps/5", map[string]interface{}{
		"title": &title,
	})
	assertStatus(t, w, http.StatusOK)
}

func TestRoadmapUpdateStep_Forbidden(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id/steps/:stepId", h.UpdateStep)

	rm := &model.Roadmap{Title: "Other"}
	rm.ID = 10
	rm.UserID = 999
	ports.Roadmaps.On("FindByID", mock.Anything, uint(10)).Return(rm, nil)

	w := doRequest(r, http.MethodPut, "/roadmaps/10/steps/5", map[string]string{"title": "X"})
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- DeleteStep ----------

func TestRoadmapDeleteStep_Success(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.DELETE("/roadmaps/:id/steps/:stepId", h.DeleteStep)

	rm := &model.Roadmap{Title: "My Roadmap"}
	rm.ID = 10
	rm.UserID = 1
	ports.Roadmaps.On("FindByID", mock.Anything, uint(10)).Return(rm, nil)

	step := &model.RoadmapStep{Title: "Step", RoadmapID: 10}
	step.ID = 5
	ports.Roadmaps.On("FindStepByID", mock.Anything, uint(5)).Return(step, nil)
	ports.Roadmaps.On("DeleteStep", mock.Anything, uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/roadmaps/10/steps/5", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestRoadmapDeleteStep_Forbidden(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.DELETE("/roadmaps/:id/steps/:stepId", h.DeleteStep)

	rm := &model.Roadmap{Title: "Other"}
	rm.ID = 10
	rm.UserID = 999
	ports.Roadmaps.On("FindByID", mock.Anything, uint(10)).Return(rm, nil)

	w := doRequest(r, http.MethodDelete, "/roadmaps/10/steps/5", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- ReorderSteps ----------

func TestRoadmapReorderSteps_Success(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id/steps/reorder", h.ReorderSteps)

	rm := &model.Roadmap{Title: "My Roadmap"}
	rm.ID = 10
	rm.UserID = 1
	ports.Roadmaps.On("FindByID", mock.Anything, uint(10)).Return(rm, nil)
	ports.Roadmaps.On("ReorderSteps", mock.Anything, uint(10), mock.AnythingOfType("[]model.StepOrder")).Return(nil)

	w := doRequest(r, http.MethodPut, "/roadmaps/10/steps/reorder", map[string]interface{}{
		"orders": []map[string]interface{}{
			{"step_id": 1, "order_index": 0},
			{"step_id": 2, "order_index": 1},
		},
	})
	assertStatus(t, w, http.StatusOK)
}

func TestRoadmapReorderSteps_Forbidden(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id/steps/reorder", h.ReorderSteps)

	rm := &model.Roadmap{Title: "Other"}
	rm.ID = 10
	rm.UserID = 999
	ports.Roadmaps.On("FindByID", mock.Anything, uint(10)).Return(rm, nil)

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
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/status/:status", h.GetByStatus)

	roadmaps := []model.Roadmap{{Title: "Go入門"}}
	ports.Roadmaps.On("GetByStatus", mock.Anything, uint(1), "active").Return(roadmaps, nil)

	w := doRequest(r, http.MethodGet, "/roadmaps/status/active", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Roadmaps.AssertExpectations(t)
}

func TestRoadmap_GetByStatus_NilResult(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/status/:status", h.GetByStatus)

	ports.Roadmaps.On("GetByStatus", mock.Anything, uint(1), "completed").Return([]model.Roadmap(nil), nil)

	w := doRequest(r, http.MethodGet, "/roadmaps/status/completed", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
	ports.Roadmaps.AssertExpectations(t)
}

// 未知のステータスはリポジトリを引かずに 400 を返す。
func TestRoadmap_GetByStatus_InvalidStatus(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/status/:status", h.GetByStatus)

	w := doRequest(r, http.MethodGet, "/roadmaps/status/invalid", nil)
	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "無効なステータスです")
	ports.Roadmaps.AssertNotCalled(t, "GetByStatus", mock.Anything, mock.Anything, mock.Anything)
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
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps/templates/:id/create", h.CreateFromTemplate)

	roadmap := &model.Roadmap{Title: "From Template"}
	roadmap.ID = 10
	ports.Roadmaps.On("FindByID", mock.Anything, uint(5)).Return(&model.Roadmap{ID: 5, IsTemplate: true}, nil)
	ports.Roadmaps.On("CopyRoadmap", mock.Anything, uint(5), uint(1)).Return(roadmap, nil)

	w := doRequest(r, http.MethodPost, "/roadmaps/templates/5/create", nil)
	assertStatus(t, w, http.StatusCreated)
	ports.Roadmaps.AssertExpectations(t)
}

func TestRoadmapCreateFromTemplate_InvalidID(t *testing.T) {
	h, _ := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps/templates/:id/create", h.CreateFromTemplate)

	w := doRequest(r, http.MethodPost, "/roadmaps/templates/abc/create", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// テンプレートが不在のときは 404 にならず 500 になる（移行前からの挙動）。
func TestRoadmapCreateFromTemplate_MissingReturnsInternalError(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps/templates/:id/create", h.CreateFromTemplate)

	// port は不在を (nil, nil) で表す。
	ports.Roadmaps.On("FindByID", mock.Anything, uint(5)).Return(nil, nil)

	w := doRequest(r, http.MethodPost, "/roadmaps/templates/5/create", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Roadmaps.AssertNotCalled(t, "CopyRoadmap", mock.Anything, mock.Anything, mock.Anything)
}

// テンプレートでないロードマップを指定すると 400。
func TestRoadmapCreateFromTemplate_NotATemplate(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps/templates/:id/create", h.CreateFromTemplate)

	ports.Roadmaps.On("FindByID", mock.Anything, uint(5)).
		Return(&model.Roadmap{ID: 5, IsTemplate: false}, nil)

	w := doRequest(r, http.MethodPost, "/roadmaps/templates/5/create", nil)
	assertStatus(t, w, http.StatusBadRequest)
	ports.Roadmaps.AssertNotCalled(t, "CopyRoadmap", mock.Anything, mock.Anything, mock.Anything)
}

// ---------- GetTemplates エラーパス ----------

func TestRoadmapGetTemplates_ServiceError(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/templates", h.GetTemplates)

	ports.Roadmaps.On("GetTemplates", mock.Anything).Return([]model.Roadmap(nil), assert.AnError)

	w := doRequest(r, http.MethodGet, "/roadmaps/templates", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Roadmaps.AssertExpectations(t)
}

// ---------- GetMyRoadmaps エラーパス ----------

func TestRoadmapGetMy_ServiceError(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps", h.GetMyRoadmaps)

	ports.Roadmaps.On("GetByUserID", mock.Anything, uint(1), 20, 0).Return([]model.Roadmap(nil), int64(0), assert.AnError)

	w := doRequest(r, http.MethodGet, "/roadmaps", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Roadmaps.AssertExpectations(t)
}

// ---------- GetPublicRoadmaps エラーパス ----------

func TestRoadmapGetPublic_ServiceError(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/public", h.GetPublicRoadmaps)

	ports.Roadmaps.On("GetPublicRoadmaps", mock.Anything, 20, 0).Return([]model.Roadmap(nil), int64(0), assert.AnError)

	w := doRequest(r, http.MethodGet, "/roadmaps/public", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Roadmaps.AssertExpectations(t)
}

// ---------- UpdateStep (完了ステータスのみ) ----------

func TestRoadmapUpdateStep_CompletionOnly(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id/steps/:stepId", h.UpdateStep)

	ports.Roadmaps.On("FindByID", mock.Anything, uint(1)).Return(&model.Roadmap{ID: 1, UserID: 1}, nil)
	ports.Roadmaps.On("FindStepByID", mock.Anything, uint(2)).Return(&model.RoadmapStep{ID: 2, RoadmapID: 1}, nil)
	ports.Roadmaps.On("UpdateStep", mock.Anything, mock.AnythingOfType("*model.RoadmapStep")).Return(nil)

	w := doRequest(r, http.MethodPut, "/roadmaps/1/steps/2", map[string]interface{}{
		"is_completed": true,
	})
	assertStatus(t, w, http.StatusOK)
	ports.Roadmaps.AssertExpectations(t)
}

func TestRoadmapUpdateStep_CompletionError(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id/steps/:stepId", h.UpdateStep)

	ports.Roadmaps.On("FindByID", mock.Anything, uint(1)).Return(&model.Roadmap{ID: 1, UserID: 999}, nil)

	w := doRequest(r, http.MethodPut, "/roadmaps/1/steps/2", map[string]interface{}{
		"is_completed": true,
	})
	assertStatus(t, w, http.StatusForbidden)
	ports.Roadmaps.AssertExpectations(t)
}

// ---------- UpdateStep InvalidStepID ----------

func TestRoadmapUpdateStep_InvalidStepID(t *testing.T) {
	h, _ := setupRoadmapHandler()
	r := newRouter(1)
	r.PUT("/roadmaps/:id/steps/:stepId", h.UpdateStep)

	w := doRequest(r, http.MethodPut, "/roadmaps/1/steps/abc", map[string]interface{}{
		"title": "test",
	})
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- DeleteStep InvalidStepID ----------

func TestRoadmapDeleteStep_InvalidStepID(t *testing.T) {
	h, _ := setupRoadmapHandler()
	r := newRouter(1)
	r.DELETE("/roadmaps/:id/steps/:stepId", h.DeleteStep)

	w := doRequest(r, http.MethodDelete, "/roadmaps/1/steps/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- CopyRoadmap InvalidID ----------

func TestRoadmapCopy_InvalidID(t *testing.T) {
	h, _ := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps/:id/copy", h.CopyRoadmap)

	w := doRequest(r, http.MethodPost, "/roadmaps/abc/copy", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- CopyRoadmap ServiceError ----------

func TestRoadmapCopy_ServiceError(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.POST("/roadmaps/:id/copy", h.CopyRoadmap)

	ports.Roadmaps.On("FindByID", mock.Anything, uint(5)).Return(&model.Roadmap{ID: 5, IsPublic: true}, nil)
	ports.Roadmaps.On("CopyRoadmap", mock.Anything, uint(5), uint(1)).Return(nil, assert.AnError)

	w := doRequest(r, http.MethodPost, "/roadmaps/5/copy", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Roadmaps.AssertExpectations(t)
}

// ---------- ReorderSteps InvalidID ----------

func TestRoadmapReorderSteps_InvalidID(t *testing.T) {
	h, _ := setupRoadmapHandler()
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
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/my/stats", h.GetMyStats)

	stats := &model.RoadmapStats{TotalRoadmaps: 3, ActiveRoadmaps: 2, CompletedRoadmaps: 1, TotalSteps: 10, CompletedSteps: 5}
	ports.Stats.On("GetRoadmapStats", mock.Anything, uint(1)).Return(stats, nil)

	w := doRequest(r, http.MethodGet, "/roadmaps/my/stats", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Roadmaps.AssertExpectations(t)
}

func TestRoadmapGetMyStats_ServiceError(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/my/stats", h.GetMyStats)

	ports.Stats.On("GetRoadmapStats", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/roadmaps/my/stats", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Roadmaps.AssertExpectations(t)
}

// ---------- GetMyCount ----------

func TestRoadmapGetMyCount_Success(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/my/count", h.GetMyCount)

	ports.Roadmaps.On("CountByUserID", mock.Anything, uint(1)).Return(int64(5), nil)

	w := doRequest(r, http.MethodGet, "/roadmaps/my/count", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"count":5`)
	ports.Roadmaps.AssertExpectations(t)
}

func TestRoadmapGetMyCount_ServiceError(t *testing.T) {
	h, ports := setupRoadmapHandler()
	r := newRouter(1)
	r.GET("/roadmaps/my/count", h.GetMyCount)

	ports.Roadmaps.On("CountByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/roadmaps/my/count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Roadmaps.AssertExpectations(t)
}

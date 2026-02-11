package handler

import (
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/mock"
)

// ---------- Create ----------

func TestStudyCircleCreate_Success(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.POST("/study-circles", h.Create)

	repo.On("Create", mock.AnythingOfType("*model.StudyCircle")).Return(nil)
	repo.On("AddMember", mock.AnythingOfType("uint"), mock.AnythingOfType("uint"), model.StudyCircleRoleOwner).Return(nil)

	w := doRequest(r, http.MethodPost, "/study-circles", map[string]interface{}{
		"name": "React学習会", "topic": "React入門",
	})
	assertStatus(t, w, http.StatusCreated)
}

func TestStudyCircleCreate_ValidationError(t *testing.T) {
	h, _ := setupStudyCircleHandler()
	r := newRouter(1)
	r.POST("/study-circles", h.Create)

	w := doRequest(r, http.MethodPost, "/study-circles", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- GetMyCircles ----------

func TestStudyCircleGetMyCircles_Success(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.GET("/study-circles", h.GetMyCircles)

	repo.On("FindByUserID", uint(1)).Return([]model.StudyCircle{
		{Name: "Circle A"}, {Name: "Circle B"},
	}, nil)

	w := doRequest(r, http.MethodGet, "/study-circles", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- GetByID ----------

func TestStudyCircleGetByID_Success(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.GET("/study-circles/:id", h.GetByID)

	repo.On("FindByID", uint(10)).Return(&model.StudyCircle{Name: "Test Circle", OwnerID: 1}, nil)
	repo.On("IsMember", uint(10), uint(1)).Return(true, nil)

	w := doRequest(r, http.MethodGet, "/study-circles/10", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestStudyCircleGetByID_NotMember(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.GET("/study-circles/:id", h.GetByID)

	repo.On("FindByID", uint(10)).Return(&model.StudyCircle{Name: "Test", OwnerID: 2}, nil)
	repo.On("IsMember", uint(10), uint(1)).Return(false, nil)

	w := doRequest(r, http.MethodGet, "/study-circles/10", nil)
	assertStatus(t, w, http.StatusForbidden)
}

func TestStudyCircleGetByID_InvalidID(t *testing.T) {
	h, _ := setupStudyCircleHandler()
	r := newRouter(1)
	r.GET("/study-circles/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/study-circles/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- Update ----------

func TestStudyCircleUpdate_Success(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.PUT("/study-circles/:id", h.Update)

	circle := &model.StudyCircle{Name: "Old", OwnerID: 1}
	circle.ID = 10
	repo.On("FindByID", uint(10)).Return(circle, nil)
	repo.On("Update", mock.AnythingOfType("*model.StudyCircle")).Return(nil)

	w := doRequest(r, http.MethodPut, "/study-circles/10", map[string]string{
		"name": "New Name",
	})
	assertStatus(t, w, http.StatusOK)
}

func TestStudyCircleUpdate_NotOwner(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.PUT("/study-circles/:id", h.Update)

	circle := &model.StudyCircle{Name: "Old", OwnerID: 999}
	circle.ID = 10
	repo.On("FindByID", uint(10)).Return(circle, nil)

	w := doRequest(r, http.MethodPut, "/study-circles/10", map[string]string{"name": "X"})
	assertStatus(t, w, http.StatusForbidden)
}

func TestStudyCircleUpdate_NotFound(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.PUT("/study-circles/:id", h.Update)

	repo.On("FindByID", uint(10)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodPut, "/study-circles/10", map[string]string{"name": "X"})
	assertStatus(t, w, http.StatusNotFound)
}

// ---------- Delete ----------

func TestStudyCircleDelete_Success(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.DELETE("/study-circles/:id", h.Delete)

	circle := &model.StudyCircle{OwnerID: 1}
	circle.ID = 10
	repo.On("FindByID", uint(10)).Return(circle, nil)
	repo.On("Delete", uint(10)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/study-circles/10", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestStudyCircleDelete_NotOwner(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.DELETE("/study-circles/:id", h.Delete)

	circle := &model.StudyCircle{OwnerID: 999}
	circle.ID = 10
	repo.On("FindByID", uint(10)).Return(circle, nil)

	w := doRequest(r, http.MethodDelete, "/study-circles/10", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- GetMembers ----------

func TestStudyCircleGetMembers_Success(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.GET("/study-circles/:id/members", h.GetMembers)

	repo.On("IsMember", uint(10), uint(1)).Return(true, nil)
	repo.On("GetMembers", uint(10)).Return([]model.StudyCircleMember{
		{UserID: 1}, {UserID: 2},
	}, nil)

	w := doRequest(r, http.MethodGet, "/study-circles/10/members", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestStudyCircleGetMembers_NotMember(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.GET("/study-circles/:id/members", h.GetMembers)

	repo.On("IsMember", uint(10), uint(1)).Return(false, nil)

	w := doRequest(r, http.MethodGet, "/study-circles/10/members", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- AddMember ----------

func TestStudyCircleAddMember_Success(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.POST("/study-circles/:id/members", h.AddMember)

	repo.On("IsMember", uint(10), uint(1)).Return(true, nil)
	repo.On("FindByID", uint(10)).Return(&model.StudyCircle{MaxMembers: 5, OwnerID: 1}, nil)
	repo.On("GetMemberCount", uint(10)).Return(3, nil)
	repo.On("AddMember", uint(10), uint(5), model.StudyCircleRoleMember).Return(nil)

	w := doRequest(r, http.MethodPost, "/study-circles/10/members", map[string]uint{
		"user_id": 5,
	})
	assertStatus(t, w, http.StatusOK)
}

func TestStudyCircleAddMember_LimitReached(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.POST("/study-circles/:id/members", h.AddMember)

	repo.On("IsMember", uint(10), uint(1)).Return(true, nil)
	repo.On("FindByID", uint(10)).Return(&model.StudyCircle{MaxMembers: 5, OwnerID: 1}, nil)
	repo.On("GetMemberCount", uint(10)).Return(5, nil)

	w := doRequest(r, http.MethodPost, "/study-circles/10/members", map[string]uint{
		"user_id": 5,
	})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestStudyCircleAddMember_NotMember(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.POST("/study-circles/:id/members", h.AddMember)

	repo.On("IsMember", uint(10), uint(1)).Return(false, nil)

	w := doRequest(r, http.MethodPost, "/study-circles/10/members", map[string]uint{
		"user_id": 5,
	})
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- RemoveMember ----------

func TestStudyCircleRemoveMember_OwnerRemoves(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.DELETE("/study-circles/:id/members/:userId", h.RemoveMember)

	circle := &model.StudyCircle{OwnerID: 1}
	circle.ID = 10
	repo.On("FindByID", uint(10)).Return(circle, nil)
	repo.On("RemoveMember", uint(10), uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/study-circles/10/members/5", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestStudyCircleRemoveMember_SelfLeave(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.DELETE("/study-circles/:id/members/:userId", h.RemoveMember)

	circle := &model.StudyCircle{OwnerID: 999}
	circle.ID = 10
	repo.On("FindByID", uint(10)).Return(circle, nil)
	repo.On("RemoveMember", uint(10), uint(1)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/study-circles/10/members/1", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestStudyCircleRemoveMember_Forbidden(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.DELETE("/study-circles/:id/members/:userId", h.RemoveMember)

	circle := &model.StudyCircle{OwnerID: 999}
	circle.ID = 10
	repo.On("FindByID", uint(10)).Return(circle, nil)

	w := doRequest(r, http.MethodDelete, "/study-circles/10/members/5", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- CreateStep ----------

func TestStudyCircleCreateStep_Success(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.POST("/study-circles/:id/steps", h.CreateStep)

	circle := &model.StudyCircle{OwnerID: 1}
	circle.ID = 10
	repo.On("FindByID", uint(10)).Return(circle, nil)
	repo.On("CreateStep", mock.AnythingOfType("*model.StudyCircleStep")).Return(nil)

	w := doRequest(r, http.MethodPost, "/study-circles/10/steps", map[string]interface{}{
		"title": "Step 1", "description": "First step",
	})
	assertStatus(t, w, http.StatusCreated)
}

func TestStudyCircleCreateStep_NotOwner(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.POST("/study-circles/:id/steps", h.CreateStep)

	circle := &model.StudyCircle{OwnerID: 999}
	circle.ID = 10
	repo.On("FindByID", uint(10)).Return(circle, nil)

	w := doRequest(r, http.MethodPost, "/study-circles/10/steps", map[string]interface{}{
		"title": "Step 1",
	})
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- UpdateProgress ----------

func TestStudyCircleUpdateProgress_Success(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.PUT("/study-circles/:id/steps/:stepId/progress", h.UpdateProgress)

	repo.On("IsMember", uint(10), uint(1)).Return(true, nil)
	repo.On("UpsertProgress", mock.AnythingOfType("*model.StudyCircleMemberProgress")).Return(nil)

	w := doRequest(r, http.MethodPut, "/study-circles/10/steps/5/progress", map[string]bool{
		"is_completed": true,
	})
	assertStatus(t, w, http.StatusOK)
}

func TestStudyCircleUpdateProgress_NotMember(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.PUT("/study-circles/:id/steps/:stepId/progress", h.UpdateProgress)

	repo.On("IsMember", uint(10), uint(1)).Return(false, nil)

	w := doRequest(r, http.MethodPut, "/study-circles/10/steps/5/progress", map[string]bool{
		"is_completed": true,
	})
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- CreateCheckin ----------

func TestStudyCircleCreateCheckin_Success(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.POST("/study-circles/:id/checkins", h.CreateCheckin)

	repo.On("IsMember", uint(10), uint(1)).Return(true, nil)
	repo.On("HasCheckedInToday", uint(10), uint(1)).Return(false, nil)
	repo.On("CreateCheckin", mock.AnythingOfType("*model.StudyCircleCheckin")).Return(nil)

	w := doRequest(r, http.MethodPost, "/study-circles/10/checkins", map[string]string{
		"content": "React hooks を学んだ",
	})
	assertStatus(t, w, http.StatusCreated)
}

func TestStudyCircleCreateCheckin_Duplicate(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.POST("/study-circles/:id/checkins", h.CreateCheckin)

	repo.On("IsMember", uint(10), uint(1)).Return(true, nil)
	repo.On("HasCheckedInToday", uint(10), uint(1)).Return(true, nil)

	w := doRequest(r, http.MethodPost, "/study-circles/10/checkins", map[string]string{
		"content": "duplicate",
	})
	assertStatus(t, w, http.StatusConflict)
}

func TestStudyCircleCreateCheckin_NotMember(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.POST("/study-circles/:id/checkins", h.CreateCheckin)

	repo.On("IsMember", uint(10), uint(1)).Return(false, nil)

	w := doRequest(r, http.MethodPost, "/study-circles/10/checkins", map[string]string{
		"content": "test",
	})
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- GetCheckins ----------

func TestStudyCircleGetCheckins_Success(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.GET("/study-circles/:id/checkins", h.GetCheckins)

	repo.On("IsMember", uint(10), uint(1)).Return(true, nil)
	repo.On("GetCheckins", uint(10)).Return([]model.StudyCircleCheckin{
		{Content: "React hooks"},
	}, nil)

	w := doRequest(r, http.MethodGet, "/study-circles/10/checkins", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- GetStreakRanking ----------

func TestStudyCircleGetStreakRanking_Success(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.GET("/study-circles/:id/streak-ranking", h.GetStreakRanking)

	repo.On("IsMember", uint(10), uint(1)).Return(true, nil)
	repo.On("GetStreakRanking", uint(10)).Return([]model.CircleMemberStreak{
		{UserID: 1, UserName: "Alice", CurrentStreak: 5, TotalCheckins: 10},
	}, nil)

	w := doRequest(r, http.MethodGet, "/study-circles/10/streak-ranking", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- ReorderSteps ----------

func TestStudyCircleReorderSteps_Success(t *testing.T) {
	h, repo := setupStudyCircleHandler()
	r := newRouter(1)
	r.PUT("/study-circles/:id/steps/reorder", h.ReorderSteps)

	circle := &model.StudyCircle{OwnerID: 1}
	circle.ID = 10
	repo.On("FindByID", uint(10)).Return(circle, nil)
	repo.On("ReorderSteps", uint(10), mock.AnythingOfType("[]repository.StepOrder")).Return(nil)

	w := doRequest(r, http.MethodPut, "/study-circles/10/steps/reorder", map[string]interface{}{
		"orders": []map[string]interface{}{
			{"step_id": 1, "order_index": 0},
			{"step_id": 2, "order_index": 1},
		},
	})
	assertStatus(t, w, http.StatusOK)
}

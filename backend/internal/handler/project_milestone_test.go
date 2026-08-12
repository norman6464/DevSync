package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockProjectMilestoneRepo は usecase/repository.ProjectMilestoneRepository のモック（ctx 付き）。
type mockProjectMilestoneRepo struct{ mock.Mock }

func (m *mockProjectMilestoneRepo) Create(ctx context.Context, milestone *model.ProjectMilestone) error {
	return m.Called(ctx, milestone).Error(0)
}

func (m *mockProjectMilestoneRepo) FindByID(ctx context.Context, id uint) (*model.ProjectMilestone, error) {
	args := m.Called(ctx, id)
	ms, _ := args.Get(0).(*model.ProjectMilestone)
	return ms, args.Error(1)
}

func (m *mockProjectMilestoneRepo) FindByProjectID(ctx context.Context, projectID uint) ([]model.ProjectMilestone, error) {
	args := m.Called(ctx, projectID)
	ms, _ := args.Get(0).([]model.ProjectMilestone)
	return ms, args.Error(1)
}

func (m *mockProjectMilestoneRepo) Update(ctx context.Context, milestone *model.ProjectMilestone) error {
	return m.Called(ctx, milestone).Error(0)
}

func (m *mockProjectMilestoneRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

// mockProjectReader は usecase/repository.ProjectReader のモック。
type mockProjectReader struct{ mock.Mock }

func (m *mockProjectReader) FindByID(ctx context.Context, id uint) (*model.Project, error) {
	args := m.Called(ctx, id)
	p, _ := args.Get(0).(*model.Project)
	return p, args.Error(1)
}

// setupMilestoneHandler は本物の usecase と port モックで ProjectMilestoneHandler を組む。
func setupMilestoneHandler() (*ProjectMilestoneHandler, *mockProjectMilestoneRepo, *mockProjectReader) {
	milestones := new(mockProjectMilestoneRepo)
	projects := new(mockProjectReader)
	h := NewProjectMilestoneHandler(
		usecase.NewCreateProjectMilestoneUseCase(milestones, projects),
		usecase.NewListProjectMilestonesUseCase(milestones),
		usecase.NewUpdateProjectMilestoneUseCase(milestones, projects),
		usecase.NewDeleteProjectMilestoneUseCase(milestones, projects),
	)
	return h, milestones, projects
}

// ownedProject は所有者が userID=1 のプロジェクトを返す。
func ownedProject(id uint) *model.Project {
	p := &model.Project{}
	p.ID = id
	p.UserID = 1
	return p
}

// --- Create ---

func TestProjectMilestoneHandler_Create_Success(t *testing.T) {
	h, milestones, projects := setupMilestoneHandler()
	projects.On("FindByID", mock.Anything, uint(10)).Return(ownedProject(10), nil)
	milestones.On("Create", mock.Anything, mock.MatchedBy(func(ms *model.ProjectMilestone) bool {
		return ms.ProjectID == 10 && ms.Title == "v1.0リリース" &&
			ms.Description == "初回リリース" && ms.Status == model.MilestoneNotStarted
	})).Return(nil)

	r := newRouter(1)
	r.POST("/projects/:id/milestones", h.Create)

	w := doRequest(r, "POST", "/projects/10/milestones", map[string]interface{}{
		"title":       "v1.0リリース",
		"description": "初回リリース",
	})
	assertStatus(t, w, http.StatusCreated)
	milestones.AssertExpectations(t)
}

// 他ユーザーのプロジェクトへの作成は 403 を返し、作成しない。
func TestProjectMilestoneHandler_Create_Forbidden(t *testing.T) {
	h, milestones, projects := setupMilestoneHandler()
	other := &model.Project{}
	other.ID = 10
	other.UserID = 999
	projects.On("FindByID", mock.Anything, uint(10)).Return(other, nil)

	r := newRouter(1)
	r.POST("/projects/:id/milestones", h.Create)

	w := doRequest(r, "POST", "/projects/10/milestones", map[string]interface{}{"title": "t"})
	assertStatus(t, w, http.StatusForbidden)
	milestones.AssertNotCalled(t, "Create")
}

// プロジェクトが存在しなければ 404。
func TestProjectMilestoneHandler_Create_ProjectNotFound(t *testing.T) {
	h, milestones, projects := setupMilestoneHandler()
	projects.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)

	r := newRouter(1)
	r.POST("/projects/:id/milestones", h.Create)

	w := doRequest(r, "POST", "/projects/10/milestones", map[string]interface{}{"title": "t"})
	assertStatus(t, w, http.StatusNotFound)
	milestones.AssertNotCalled(t, "Create")
}

// タイトルが空・長すぎる場合は 400 を返し、作成しない。
func TestProjectMilestoneHandler_Create_InvalidTitle(t *testing.T) {
	for name, title := range map[string]string{
		"空文字":    "",
		"201 文字": strings.Repeat("a", 201),
	} {
		t.Run(name, func(t *testing.T) {
			h, milestones, projects := setupMilestoneHandler()
			projects.On("FindByID", mock.Anything, uint(10)).Return(ownedProject(10), nil)

			r := newRouter(1)
			r.POST("/projects/:id/milestones", h.Create)

			w := doRequest(r, "POST", "/projects/10/milestones", map[string]interface{}{"title": title})
			assertStatus(t, w, http.StatusBadRequest)
			milestones.AssertNotCalled(t, "Create")
		})
	}
}

func TestProjectMilestoneHandler_Create_InvalidJSON(t *testing.T) {
	h, _, _ := setupMilestoneHandler()

	r := newRouter(1)
	r.POST("/projects/:id/milestones", h.Create)

	w := doRequestRaw(r, "POST", "/projects/10/milestones", "bad json")
	assertStatus(t, w, http.StatusBadRequest)
}

// 日付の形式が不正なら 400。
func TestProjectMilestoneHandler_Create_InvalidDueDate(t *testing.T) {
	h, milestones, _ := setupMilestoneHandler()

	r := newRouter(1)
	r.POST("/projects/:id/milestones", h.Create)

	w := doRequest(r, "POST", "/projects/10/milestones", map[string]interface{}{
		"title": "t", "due_date": "2026/01/01",
	})
	assertStatus(t, w, http.StatusBadRequest)
	milestones.AssertNotCalled(t, "Create")
}

// --- List ---

func TestProjectMilestoneHandler_GetByProjectID_Success(t *testing.T) {
	h, milestones, _ := setupMilestoneHandler()
	milestones.On("FindByProjectID", mock.Anything, uint(10)).
		Return([]model.ProjectMilestone{{ProjectID: 10, Title: "m1"}}, nil)

	r := newRouter(1)
	r.GET("/projects/:id/milestones", h.GetByProjectID)

	w := doRequest(r, "GET", "/projects/10/milestones", nil)
	assertStatus(t, w, http.StatusOK)
	milestones.AssertExpectations(t)
}

func TestProjectMilestoneHandler_GetByProjectID_RepoError(t *testing.T) {
	h, milestones, _ := setupMilestoneHandler()
	milestones.On("FindByProjectID", mock.Anything, uint(10)).
		Return([]model.ProjectMilestone(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/projects/:id/milestones", h.GetByProjectID)

	w := doRequest(r, "GET", "/projects/10/milestones", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// --- Update ---

func TestProjectMilestoneHandler_Update_Success(t *testing.T) {
	h, milestones, projects := setupMilestoneHandler()
	existing := &model.ProjectMilestone{ProjectID: 10, Title: "旧題", Description: "旧説明"}
	existing.ID = 7
	milestones.On("FindByID", mock.Anything, uint(7)).Return(existing, nil)
	projects.On("FindByID", mock.Anything, uint(10)).Return(ownedProject(10), nil)
	milestones.On("Update", mock.Anything, mock.MatchedBy(func(ms *model.ProjectMilestone) bool {
		return ms.Title == "新題" && ms.Description == "旧説明"
	})).Return(nil)

	r := newRouter(1)
	r.PUT("/projects/milestones/:milestoneId", h.Update)

	w := doRequest(r, "PUT", "/projects/milestones/7", map[string]interface{}{"title": "新題"})
	assertStatus(t, w, http.StatusOK)
	milestones.AssertExpectations(t)
}

// 完了へ遷移すると完了日時が入り、戻すと消える。
func TestProjectMilestoneHandler_Update_StatusTransitions(t *testing.T) {
	t.Run("完了で完了日時が入る", func(t *testing.T) {
		h, milestones, projects := setupMilestoneHandler()
		existing := &model.ProjectMilestone{ProjectID: 10, Title: "t"}
		existing.ID = 7
		milestones.On("FindByID", mock.Anything, uint(7)).Return(existing, nil)
		projects.On("FindByID", mock.Anything, uint(10)).Return(ownedProject(10), nil)
		milestones.On("Update", mock.Anything, mock.MatchedBy(func(ms *model.ProjectMilestone) bool {
			return ms.Status == model.MilestoneCompleted && ms.CompletedAt != nil
		})).Return(nil)

		r := newRouter(1)
		r.PUT("/projects/milestones/:milestoneId", h.Update)

		w := doRequest(r, "PUT", "/projects/milestones/7", map[string]interface{}{"status": "completed"})
		assertStatus(t, w, http.StatusOK)
		milestones.AssertExpectations(t)
	})

	t.Run("未着手へ戻すと完了日時が消える", func(t *testing.T) {
		h, milestones, projects := setupMilestoneHandler()
		now := time.Now()
		existing := &model.ProjectMilestone{ProjectID: 10, Title: "t", Status: model.MilestoneCompleted, CompletedAt: &now}
		existing.ID = 7
		milestones.On("FindByID", mock.Anything, uint(7)).Return(existing, nil)
		projects.On("FindByID", mock.Anything, uint(10)).Return(ownedProject(10), nil)
		milestones.On("Update", mock.Anything, mock.MatchedBy(func(ms *model.ProjectMilestone) bool {
			return ms.Status == model.MilestoneNotStarted && ms.CompletedAt == nil
		})).Return(nil)

		r := newRouter(1)
		r.PUT("/projects/milestones/:milestoneId", h.Update)

		w := doRequest(r, "PUT", "/projects/milestones/7", map[string]interface{}{"status": "not_started"})
		assertStatus(t, w, http.StatusOK)
		milestones.AssertExpectations(t)
	})
}

// 無効なステータスは 400 を返し、保存しない。
func TestProjectMilestoneHandler_Update_InvalidStatus(t *testing.T) {
	h, milestones, projects := setupMilestoneHandler()
	existing := &model.ProjectMilestone{ProjectID: 10, Title: "t"}
	existing.ID = 7
	milestones.On("FindByID", mock.Anything, uint(7)).Return(existing, nil)
	projects.On("FindByID", mock.Anything, uint(10)).Return(ownedProject(10), nil)

	r := newRouter(1)
	r.PUT("/projects/milestones/:milestoneId", h.Update)

	w := doRequest(r, "PUT", "/projects/milestones/7", map[string]interface{}{"status": "done"})
	assertStatus(t, w, http.StatusBadRequest)
	milestones.AssertNotCalled(t, "Update")
}

// 存在しないマイルストーンは 404。
func TestProjectMilestoneHandler_Update_NotFound(t *testing.T) {
	h, milestones, _ := setupMilestoneHandler()
	milestones.On("FindByID", mock.Anything, uint(7)).Return(nil, nil)

	r := newRouter(1)
	r.PUT("/projects/milestones/:milestoneId", h.Update)

	w := doRequest(r, "PUT", "/projects/milestones/7", map[string]interface{}{"title": "t"})
	assertStatus(t, w, http.StatusNotFound)
	milestones.AssertNotCalled(t, "Update")
}

// 他ユーザーのプロジェクトのマイルストーン更新は 403。
func TestProjectMilestoneHandler_Update_Forbidden(t *testing.T) {
	h, milestones, projects := setupMilestoneHandler()
	existing := &model.ProjectMilestone{ProjectID: 10, Title: "t"}
	existing.ID = 7
	milestones.On("FindByID", mock.Anything, uint(7)).Return(existing, nil)
	other := &model.Project{}
	other.ID = 10
	other.UserID = 999
	projects.On("FindByID", mock.Anything, uint(10)).Return(other, nil)

	r := newRouter(1)
	r.PUT("/projects/milestones/:milestoneId", h.Update)

	w := doRequest(r, "PUT", "/projects/milestones/7", map[string]interface{}{"title": "t"})
	assertStatus(t, w, http.StatusForbidden)
	milestones.AssertNotCalled(t, "Update")
}

// --- Delete ---

func TestProjectMilestoneHandler_Delete_Success(t *testing.T) {
	h, milestones, projects := setupMilestoneHandler()
	existing := &model.ProjectMilestone{ProjectID: 10}
	existing.ID = 7
	milestones.On("FindByID", mock.Anything, uint(7)).Return(existing, nil)
	projects.On("FindByID", mock.Anything, uint(10)).Return(ownedProject(10), nil)
	milestones.On("Delete", mock.Anything, uint(7)).Return(nil)

	r := newRouter(1)
	r.DELETE("/projects/milestones/:milestoneId", h.Delete)

	w := doRequest(r, "DELETE", "/projects/milestones/7", nil)
	assertStatus(t, w, http.StatusOK)
	milestones.AssertExpectations(t)
}

func TestProjectMilestoneHandler_Delete_Forbidden(t *testing.T) {
	h, milestones, projects := setupMilestoneHandler()
	existing := &model.ProjectMilestone{ProjectID: 10}
	existing.ID = 7
	milestones.On("FindByID", mock.Anything, uint(7)).Return(existing, nil)
	other := &model.Project{}
	other.ID = 10
	other.UserID = 999
	projects.On("FindByID", mock.Anything, uint(10)).Return(other, nil)

	r := newRouter(1)
	r.DELETE("/projects/milestones/:milestoneId", h.Delete)

	w := doRequest(r, "DELETE", "/projects/milestones/7", nil)
	assertStatus(t, w, http.StatusForbidden)
	milestones.AssertNotCalled(t, "Delete")
}

func TestProjectMilestoneHandler_Delete_InvalidID(t *testing.T) {
	h, _, _ := setupMilestoneHandler()

	r := newRouter(1)
	r.DELETE("/projects/milestones/:milestoneId", h.Delete)

	w := doRequest(r, "DELETE", "/projects/milestones/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// 403 のメッセージは共通 helper の汎用文言ではなく、このスライス固有のもの。
func TestProjectMilestoneHandler_Forbidden_Message(t *testing.T) {
	h, _, projects := setupMilestoneHandler()
	other := &model.Project{}
	other.ID = 10
	other.UserID = 999
	projects.On("FindByID", mock.Anything, uint(10)).Return(other, nil)

	r := newRouter(1)
	r.POST("/projects/:id/milestones", h.Create)

	w := doRequest(r, "POST", "/projects/10/milestones", map[string]interface{}{"title": "t"})
	body := parseJSON(t, w)
	assert.Equal(t, "このプロジェクトを編集する権限がありません", body["error"])
}

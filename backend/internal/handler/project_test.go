package handler

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockProjectRepo は usecase/repository.ProjectRepository のモック（ctx 付き）。
type mockProjectRepo struct{ mock.Mock }

func (m *mockProjectRepo) Create(ctx context.Context, project *model.Project) error {
	return m.Called(ctx, project).Error(0)
}

func (m *mockProjectRepo) Update(ctx context.Context, project *model.Project) error {
	return m.Called(ctx, project).Error(0)
}

func (m *mockProjectRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockProjectRepo) FindByID(ctx context.Context, id uint) (*model.Project, error) {
	args := m.Called(ctx, id)
	p, _ := args.Get(0).(*model.Project)
	return p, args.Error(1)
}

func (m *mockProjectRepo) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Project, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	p, _ := args.Get(0).([]model.Project)
	return p, args.Get(1).(int64), args.Error(2)
}

func (m *mockProjectRepo) FindFeaturedByUserID(ctx context.Context, userID uint) ([]model.Project, error) {
	args := m.Called(ctx, userID)
	p, _ := args.Get(0).([]model.Project)
	return p, args.Error(1)
}

func (m *mockProjectRepo) FindArchivedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Project, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	p, _ := args.Get(0).([]model.Project)
	return p, args.Get(1).(int64), args.Error(2)
}

func (m *mockProjectRepo) FindAll(ctx context.Context, limit, offset int) ([]model.Project, int64, error) {
	args := m.Called(ctx, limit, offset)
	p, _ := args.Get(0).([]model.Project)
	return p, args.Get(1).(int64), args.Error(2)
}

func (m *mockProjectRepo) Search(ctx context.Context, query string, limit, offset int) ([]model.Project, int64, error) {
	args := m.Called(ctx, query, limit, offset)
	p, _ := args.Get(0).([]model.Project)
	return p, args.Get(1).(int64), args.Error(2)
}

func (m *mockProjectRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockProjectRepo) Archive(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockProjectRepo) Unarchive(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

// newTestProjectHandler は本物の usecase に port モックを注入した ProjectHandler を生成する。
func newTestProjectHandler() (*ProjectHandler, *mockProjectRepo) {
	repo := new(mockProjectRepo)
	h := NewProjectHandler(
		usecase.NewCreateProjectUseCase(repo),
		usecase.NewGetProjectUseCase(repo),
		usecase.NewListProjectsByUserUseCase(repo),
		usecase.NewListFeaturedProjectsUseCase(repo),
		usecase.NewListAllProjectsUseCase(repo),
		usecase.NewListArchivedProjectsUseCase(repo),
		usecase.NewSearchProjectsUseCase(repo),
		usecase.NewUpdateProjectUseCase(repo),
		usecase.NewUpdateProjectFeaturedUseCase(repo),
		usecase.NewArchiveProjectUseCase(repo),
		usecase.NewUnarchiveProjectUseCase(repo),
		usecase.NewDeleteProjectUseCase(repo),
		usecase.NewCountProjectsUseCase(repo),
	)
	return h, repo
}

// projectOwnedBy は指定ユーザーが所有するプロジェクトを返すテスト用ヘルパー。
func projectOwnedBy(id, userID uint) *model.Project {
	return &model.Project{
		ID: id, UserID: userID,
		Title: "既存プロジェクト", Description: "説明", TechStack: "Go,React",
		Role: "Lead Developer",
	}
}

// ============================================================
// Create
// ============================================================

func TestProjectHandler_Create(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.POST("/projects", h.Create)

	repo.On("Create", mock.Anything, mock.MatchedBy(func(p *model.Project) bool {
		return p.UserID == 1 && p.Title == "新規プロジェクト"
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/projects", map[string]interface{}{
		"title": "新規プロジェクト", "description": "説明", "tech_stack": "Go",
	})
	assertStatus(t, w, http.StatusCreated)
	repo.AssertExpectations(t)
}

func TestProjectHandler_Create_WithDates(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.POST("/projects", h.Create)

	repo.On("Create", mock.Anything, mock.MatchedBy(func(p *model.Project) bool {
		return p.StartDate != nil && p.EndDate != nil &&
			p.StartDate.Format("2006-01-02") == "2026-01-01"
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/projects", map[string]interface{}{
		"title": "期間つき", "description": "説明",
		"start_date": "2026-01-01", "end_date": "2026-02-01",
	})
	assertStatus(t, w, http.StatusCreated)
	repo.AssertExpectations(t)
}

// タイトル欠落は DTO の binding で 400 になる。
func TestProjectHandler_Create_BadRequest(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.POST("/projects", h.Create)

	w := doRequest(r, http.MethodPost, "/projects", map[string]interface{}{"description": "説明"})
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// 説明の欠落は binding を通るが、usecase の検証で 400 になる。
func TestProjectHandler_Create_ValidationError(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.POST("/projects", h.Create)

	w := doRequest(r, http.MethodPost, "/projects", map[string]interface{}{"title": "題だけ"})
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// 不正なデモ URL は DTO の binding（http_url）で弾かれる。
func TestProjectHandler_Create_InvalidDemoURL(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.POST("/projects", h.Create)

	w := doRequest(r, http.MethodPost, "/projects", map[string]interface{}{
		"title": "題", "description": "説明", "demo_url": "not-a-url",
	})
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestProjectHandler_Create_RepositoryError(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.POST("/projects", h.Create)

	repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/projects", map[string]interface{}{
		"title": "題", "description": "説明",
	})
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// GetByID
// ============================================================

func TestProjectHandler_GetByID(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/projects/:id", h.GetByID)

	repo.On("FindByID", mock.Anything, uint(1)).Return(projectOwnedBy(1, 1), nil)

	w := doRequest(r, http.MethodGet, "/projects/1", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "既存プロジェクト")
}

func TestProjectHandler_GetByID_Forbidden(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/projects/:id", h.GetByID)

	repo.On("FindByID", mock.Anything, uint(1)).Return(projectOwnedBy(1, 2), nil)

	w := doRequest(r, http.MethodGet, "/projects/1", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// 不在のプロジェクトは 500 になる（移行前から変わらない挙動）。
func TestProjectHandler_GetByID_NotFoundIs500(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/projects/:id", h.GetByID)

	repo.On("FindByID", mock.Anything, uint(99)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/projects/99", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestProjectHandler_GetByID_InvalidID(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/projects/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/projects/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// ============================================================
// 一覧
// ============================================================

func TestProjectHandler_GetByUserID(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/users/:userId/projects", h.GetByUserID)

	repo.On("FindByUserID", mock.Anything, uint(7), 20, 0).
		Return([]model.Project{*projectOwnedBy(1, 7)}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/users/7/projects", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"total":1`)
	repo.AssertExpectations(t)
}

func TestProjectHandler_GetByUserID_Empty(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/users/:userId/projects", h.GetByUserID)

	repo.On("FindByUserID", mock.Anything, uint(7), 20, 0).
		Return([]model.Project(nil), int64(0), nil)

	w := doRequest(r, http.MethodGet, "/users/7/projects", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"total":0`)
}

func TestProjectHandler_GetByUserID_InvalidID(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/users/:userId/projects", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/abc/projects", nil)
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "FindByUserID", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestProjectHandler_GetByUserID_RepositoryError(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/users/:userId/projects", h.GetByUserID)

	repo.On("FindByUserID", mock.Anything, uint(7), 20, 0).
		Return([]model.Project(nil), int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/7/projects", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestProjectHandler_GetMyProjects(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/projects/my", h.GetMyProjects)

	repo.On("FindByUserID", mock.Anything, uint(1), 20, 0).
		Return([]model.Project{*projectOwnedBy(1, 1)}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/projects/my", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestProjectHandler_GetFeatured(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/users/:userId/projects/featured", h.GetFeatured)

	featured := projectOwnedBy(1, 7)
	featured.Featured = true
	repo.On("FindFeaturedByUserID", mock.Anything, uint(7)).Return([]model.Project{*featured}, nil)

	w := doRequest(r, http.MethodGet, "/users/7/projects/featured", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"featured":true`)
}

// 0 件でも null ではなく空配列を返す。
func TestProjectHandler_GetFeatured_Empty(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/users/:userId/projects/featured", h.GetFeatured)

	repo.On("FindFeaturedByUserID", mock.Anything, uint(7)).Return([]model.Project(nil), nil)

	w := doRequest(r, http.MethodGet, "/users/7/projects/featured", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

func TestProjectHandler_GetFeatured_InvalidID(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/users/:userId/projects/featured", h.GetFeatured)

	w := doRequest(r, http.MethodGet, "/users/abc/projects/featured", nil)
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "FindFeaturedByUserID", mock.Anything, mock.Anything)
}

func TestProjectHandler_GetFeatured_RepositoryError(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/users/:userId/projects/featured", h.GetFeatured)

	repo.On("FindFeaturedByUserID", mock.Anything, uint(7)).
		Return([]model.Project(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/7/projects/featured", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestProjectHandler_GetAll(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/projects", h.GetAll)

	repo.On("FindAll", mock.Anything, 20, 0).
		Return([]model.Project{*projectOwnedBy(1, 7)}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/projects", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestProjectHandler_GetAll_RepositoryError(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/projects", h.GetAll)

	repo.On("FindAll", mock.Anything, 20, 0).
		Return([]model.Project(nil), int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/projects", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestProjectHandler_GetArchived(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/projects/archived", h.GetArchived)

	archived := projectOwnedBy(1, 1)
	archived.IsArchived = true
	repo.On("FindArchivedByUserID", mock.Anything, uint(1), 20, 0).
		Return([]model.Project{*archived}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/projects/archived", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// ============================================================
// Search
// ============================================================

func TestProjectHandler_Search(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/projects/search", h.Search)

	repo.On("Search", mock.Anything, "go", 20, 0).
		Return([]model.Project{*projectOwnedBy(1, 7)}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/projects/search?q=go", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// q の欠落はハンドラーの検証で 400 になり、usecase まで届かない。
func TestProjectHandler_Search_MissingQuery(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/projects/search", h.Search)

	w := doRequest(r, http.MethodGet, "/projects/search", nil)
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

// 空白のみのキーワードはハンドラーを通るため、usecase の検証で 400 になる。
func TestProjectHandler_Search_BlankQuery(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/projects/search", h.Search)

	w := doRequest(r, http.MethodGet, "/projects/search?q=%20%20", nil)
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Search", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestProjectHandler_Search_RepositoryError(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/projects/search", h.Search)

	repo.On("Search", mock.Anything, "go", 20, 0).
		Return([]model.Project(nil), int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/projects/search?q=go", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// Update
// ============================================================

func TestProjectHandler_Update(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.PUT("/projects/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(1)).Return(projectOwnedBy(1, 1), nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Project) bool {
		// 指定しなかったフィールドは据え置かれる。
		return p.Title == "新タイトル" && p.Description == "説明" && p.TechStack == "Go,React"
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/projects/1", map[string]interface{}{"title": "新タイトル"})
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// featured が指定された場合は、更新のあとに注目切替がもう一度呼ばれる。
func TestProjectHandler_Update_WithFeatured(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.PUT("/projects/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(1)).Return(projectOwnedBy(1, 1), nil)
	repo.On("Update", mock.Anything, mock.Anything).Return(nil)

	w := doRequest(r, http.MethodPut, "/projects/1", map[string]interface{}{
		"title": "新タイトル", "featured": true,
	})
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"featured":true`)
	// 更新と注目切替で 2 回書き込む。
	repo.AssertNumberOfCalls(t, "Update", 2)
}

func TestProjectHandler_Update_WithDates(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.PUT("/projects/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(1)).Return(projectOwnedBy(1, 1), nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Project) bool {
		return p.StartDate != nil && p.StartDate.Format("2006-01-02") == "2026-03-01"
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/projects/1", map[string]interface{}{
		"start_date": "2026-03-01", "end_date": "2026-04-01",
	})
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestProjectHandler_Update_Forbidden(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.PUT("/projects/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(1)).Return(projectOwnedBy(1, 2), nil)

	w := doRequest(r, http.MethodPut, "/projects/1", map[string]interface{}{"title": "新タイトル"})
	assertStatus(t, w, http.StatusForbidden)
	repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// 注目切替の段でリポジトリが失敗したら 500 になる。
func TestProjectHandler_Update_FeaturedError(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.PUT("/projects/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(1)).Return(projectOwnedBy(1, 1), nil)
	repo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
	repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()

	w := doRequest(r, http.MethodPut, "/projects/1", map[string]interface{}{
		"title": "新タイトル", "featured": true,
	})
	assertStatus(t, w, http.StatusInternalServerError)
}

// 不正な GitHub URL は DTO の binding（http_url）で弾かれる。
func TestProjectHandler_Update_InvalidGithubURL(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.PUT("/projects/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/projects/1", map[string]interface{}{"github_url": "not-a-url"})
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

func TestProjectHandler_Update_InvalidID(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.PUT("/projects/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/projects/abc", map[string]interface{}{"title": "新タイトル"})
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// ============================================================
// アーカイブ / 削除 / 件数
// ============================================================

func TestProjectHandler_Archive(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.POST("/projects/:id/archive", h.Archive)

	repo.On("FindByID", mock.Anything, uint(1)).Return(projectOwnedBy(1, 1), nil)
	repo.On("Archive", mock.Anything, uint(1)).Return(nil)

	w := doRequest(r, http.MethodPost, "/projects/1/archive", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// アーカイブ済みのものを再度アーカイブすると 400。
func TestProjectHandler_Archive_AlreadyArchived(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.POST("/projects/:id/archive", h.Archive)

	archived := projectOwnedBy(1, 1)
	archived.IsArchived = true
	repo.On("FindByID", mock.Anything, uint(1)).Return(archived, nil)

	w := doRequest(r, http.MethodPost, "/projects/1/archive", nil)
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Archive", mock.Anything, mock.Anything)
}

func TestProjectHandler_Unarchive(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.POST("/projects/:id/unarchive", h.Unarchive)

	archived := projectOwnedBy(1, 1)
	archived.IsArchived = true
	repo.On("FindByID", mock.Anything, uint(1)).Return(archived, nil)
	repo.On("Unarchive", mock.Anything, uint(1)).Return(nil)

	w := doRequest(r, http.MethodPost, "/projects/1/unarchive", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// アーカイブされていないものを解除すると 400。
func TestProjectHandler_Unarchive_NotArchived(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.POST("/projects/:id/unarchive", h.Unarchive)

	repo.On("FindByID", mock.Anything, uint(1)).Return(projectOwnedBy(1, 1), nil)

	w := doRequest(r, http.MethodPost, "/projects/1/unarchive", nil)
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "Unarchive", mock.Anything, mock.Anything)
}

func TestProjectHandler_Delete(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.DELETE("/projects/:id", h.Delete)

	repo.On("FindByID", mock.Anything, uint(1)).Return(projectOwnedBy(1, 1), nil)
	repo.On("Delete", mock.Anything, uint(1)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/projects/1", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestProjectHandler_Delete_Forbidden(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.DELETE("/projects/:id", h.Delete)

	repo.On("FindByID", mock.Anything, uint(1)).Return(projectOwnedBy(1, 2), nil)

	w := doRequest(r, http.MethodDelete, "/projects/1", nil)
	assertStatus(t, w, http.StatusForbidden)
	repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

// 不在のプロジェクトの削除は 500 になる（移行前から変わらない挙動）。
func TestProjectHandler_Delete_NotFoundIs500(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.DELETE("/projects/:id", h.Delete)

	repo.On("FindByID", mock.Anything, uint(99)).Return(nil, nil)

	w := doRequest(r, http.MethodDelete, "/projects/99", nil)
	assertStatus(t, w, http.StatusNotFound)
	repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

func TestProjectHandler_Delete_InvalidID(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.DELETE("/projects/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/projects/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
	repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

func TestProjectHandler_GetMyCount(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/projects/my/count", h.GetMyCount)

	repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(4), nil)

	w := doRequest(r, http.MethodGet, "/projects/my/count", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"count":4`)
	repo.AssertExpectations(t)
}

func TestProjectHandler_GetMyCount_RepositoryError(t *testing.T) {
	h, repo := newTestProjectHandler()
	r := newRouter(1)
	r.GET("/projects/my/count", h.GetMyCount)

	repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/projects/my/count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

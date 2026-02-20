package handler

import (
	"fmt"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/mock"
)

// ---------- Create ----------

func TestProjectCreate_Success(t *testing.T) {
	h, svc := setupProjectHandler()
	svc.On("Create", mock.AnythingOfType("*model.Project")).Return(nil)

	r := newRouter(1)
	r.POST("/projects", h.Create)
	w := doRequest(r, "POST", "/projects", map[string]interface{}{
		"title":       "My Project",
		"description": "A test project",
		"tech_stack":  "Go, React",
	})

	assertStatus(t, w, 201)
	svc.AssertExpectations(t)
}

func TestProjectCreate_ValidationError(t *testing.T) {
	h, _ := setupProjectHandler()

	r := newRouter(1)
	r.POST("/projects", h.Create)
	// title is required
	w := doRequest(r, "POST", "/projects", map[string]interface{}{
		"description": "missing title",
	})

	assertStatus(t, w, 400)
}

func TestProjectCreate_ServiceError(t *testing.T) {
	h, svc := setupProjectHandler()
	svc.On("Create", mock.AnythingOfType("*model.Project")).Return(fmt.Errorf("internal error"))

	r := newRouter(1)
	r.POST("/projects", h.Create)
	w := doRequest(r, "POST", "/projects", map[string]interface{}{
		"title": "My Project",
	})

	assertStatus(t, w, 500)
	svc.AssertExpectations(t)
}

// ---------- GetByID ----------

func TestProjectGetByID_Success(t *testing.T) {
	h, svc := setupProjectHandler()
	project := &model.Project{Title: "Test", UserID: 1}
	project.ID = 1
	svc.On("GetByID", uint(1), uint(1)).Return(project, nil)

	r := newRouter(1)
	r.GET("/projects/:id", h.GetByID)
	w := doRequest(r, "GET", "/projects/1", nil)

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestProjectGetByID_InvalidID(t *testing.T) {
	h, _ := setupProjectHandler()

	r := newRouter(1)
	r.GET("/projects/:id", h.GetByID)
	w := doRequest(r, "GET", "/projects/abc", nil)

	assertStatus(t, w, 400)
}

func TestProjectGetByID_NotFound(t *testing.T) {
	h, svc := setupProjectHandler()
	svc.On("GetByID", uint(999), uint(1)).Return(nil, service.ErrNotFound)

	r := newRouter(1)
	r.GET("/projects/:id", h.GetByID)
	w := doRequest(r, "GET", "/projects/999", nil)

	assertStatus(t, w, 404)
	svc.AssertExpectations(t)
}

// ---------- GetByUserID ----------

func TestProjectGetByUserID_Success(t *testing.T) {
	h, svc := setupProjectHandler()
	projects := []model.Project{{Title: "P1"}, {Title: "P2"}}
	svc.On("GetByUserID", uint(1), 20, 0).Return(projects, int64(2), nil)

	r := newRouter(1)
	r.GET("/users/:userId/projects", h.GetByUserID)
	w := doRequest(r, "GET", "/users/1/projects", nil)

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestProjectGetByUserID_Empty(t *testing.T) {
	h, svc := setupProjectHandler()
	svc.On("GetByUserID", uint(1), 20, 0).Return([]model.Project{}, int64(0), nil)

	r := newRouter(1)
	r.GET("/users/:userId/projects", h.GetByUserID)
	w := doRequest(r, "GET", "/users/1/projects", nil)

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

// ---------- GetFeatured ----------

func TestProjectGetFeatured_Success(t *testing.T) {
	h, svc := setupProjectHandler()
	projects := []model.Project{{Title: "Featured"}}
	svc.On("GetFeaturedByUserID", uint(1)).Return(projects, nil)

	r := newRouter(1)
	r.GET("/users/:userId/projects/featured", h.GetFeatured)
	w := doRequest(r, "GET", "/users/1/projects/featured", nil)

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

// ---------- Update ----------

func TestProjectUpdate_Success(t *testing.T) {
	h, svc := setupProjectHandler()
	updated := &model.Project{Title: "Updated"}
	updated.ID = 1
	svc.On("Update", uint(1), uint(1), mock.AnythingOfType("*model.Project")).Return(updated, nil)

	r := newRouter(1)
	r.PUT("/projects/:id", h.Update)
	w := doRequest(r, "PUT", "/projects/1", map[string]interface{}{
		"title": "Updated",
	})

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestProjectUpdate_WithFeatured(t *testing.T) {
	h, svc := setupProjectHandler()
	updated := &model.Project{Title: "Updated"}
	updated.ID = 1
	featured := true
	svc.On("Update", uint(1), uint(1), mock.AnythingOfType("*model.Project")).Return(updated, nil)
	svc.On("UpdateFeatured", uint(1), uint(1), featured).Return(updated, nil)

	r := newRouter(1)
	r.PUT("/projects/:id", h.Update)
	w := doRequest(r, "PUT", "/projects/1", map[string]interface{}{
		"title":    "Updated",
		"featured": true,
	})

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestProjectUpdate_ServiceError(t *testing.T) {
	h, svc := setupProjectHandler()
	svc.On("Update", uint(1), uint(1), mock.AnythingOfType("*model.Project")).Return(nil, service.ErrNotFound)

	r := newRouter(1)
	r.PUT("/projects/:id", h.Update)
	w := doRequest(r, "PUT", "/projects/1", map[string]interface{}{
		"title": "Updated",
	})

	assertStatus(t, w, 404)
	svc.AssertExpectations(t)
}

// ---------- Delete ----------

func TestProjectDelete_Success(t *testing.T) {
	h, svc := setupProjectHandler()
	svc.On("Delete", uint(1), uint(1)).Return(nil)

	r := newRouter(1)
	r.DELETE("/projects/:id", h.Delete)
	w := doRequest(r, "DELETE", "/projects/1", nil)

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestProjectDelete_NotFound(t *testing.T) {
	h, svc := setupProjectHandler()
	svc.On("Delete", uint(999), uint(1)).Return(service.ErrNotFound)

	r := newRouter(1)
	r.DELETE("/projects/:id", h.Delete)
	w := doRequest(r, "DELETE", "/projects/999", nil)

	assertStatus(t, w, 404)
	svc.AssertExpectations(t)
}

func TestProjectDelete_Forbidden(t *testing.T) {
	h, svc := setupProjectHandler()
	svc.On("Delete", uint(1), uint(2)).Return(service.ErrForbidden)

	r := newRouter(2)
	r.DELETE("/projects/:id", h.Delete)
	w := doRequest(r, "DELETE", "/projects/1", nil)

	assertStatus(t, w, 403)
	svc.AssertExpectations(t)
}

// ---------- GetAll ----------

func TestProjectGetAll_Success(t *testing.T) {
	h, svc := setupProjectHandler()
	projects := []model.Project{{Title: "P1"}}
	svc.On("GetAll", 20, 0).Return(projects, int64(1), nil)

	r := newRouter(1)
	r.GET("/projects", h.GetAll)
	w := doRequest(r, "GET", "/projects", nil)

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

func TestProjectGetAll_ServiceError(t *testing.T) {
	h, svc := setupProjectHandler()
	svc.On("GetAll", 20, 0).Return([]model.Project{}, int64(0), fmt.Errorf("internal error"))

	r := newRouter(1)
	r.GET("/projects", h.GetAll)
	w := doRequest(r, "GET", "/projects", nil)

	assertStatus(t, w, 500)
	svc.AssertExpectations(t)
}

// ---------- Create (日付付き) ----------

func TestProjectCreate_WithDates(t *testing.T) {
	h, svc := setupProjectHandler()
	svc.On("Create", mock.AnythingOfType("*model.Project")).Return(nil)

	r := newRouter(1)
	r.POST("/projects", h.Create)
	w := doRequest(r, "POST", "/projects", map[string]interface{}{
		"title":      "日付付きプロジェクト",
		"start_date": "2024-01-15",
		"end_date":   "2024-06-30",
	})

	assertStatus(t, w, 201)
	svc.AssertExpectations(t)
}

// ---------- GetByUserID (エラーパス) ----------

func TestProjectGetByUserID_InvalidID(t *testing.T) {
	h, _ := setupProjectHandler()

	r := newRouter(1)
	r.GET("/users/:userId/projects", h.GetByUserID)
	w := doRequest(r, "GET", "/users/abc/projects", nil)

	assertStatus(t, w, 400)
}

func TestProjectGetByUserID_ServiceError(t *testing.T) {
	h, svc := setupProjectHandler()
	svc.On("GetByUserID", uint(1), 20, 0).Return([]model.Project(nil), int64(0), fmt.Errorf("db error"))

	r := newRouter(1)
	r.GET("/users/:userId/projects", h.GetByUserID)
	w := doRequest(r, "GET", "/users/1/projects", nil)

	assertStatus(t, w, 500)
	svc.AssertExpectations(t)
}

// ---------- GetFeatured (エラーパス) ----------

func TestProjectGetFeatured_InvalidID(t *testing.T) {
	h, _ := setupProjectHandler()

	r := newRouter(1)
	r.GET("/users/:userId/projects/featured", h.GetFeatured)
	w := doRequest(r, "GET", "/users/abc/projects/featured", nil)

	assertStatus(t, w, 400)
}

func TestProjectGetFeatured_ServiceError(t *testing.T) {
	h, svc := setupProjectHandler()
	svc.On("GetFeaturedByUserID", uint(1)).Return([]model.Project(nil), fmt.Errorf("db error"))

	r := newRouter(1)
	r.GET("/users/:userId/projects/featured", h.GetFeatured)
	w := doRequest(r, "GET", "/users/1/projects/featured", nil)

	assertStatus(t, w, 500)
	svc.AssertExpectations(t)
}

// ---------- Update (エラーパス) ----------

func TestProjectUpdate_InvalidID(t *testing.T) {
	h, _ := setupProjectHandler()

	r := newRouter(1)
	r.PUT("/projects/:id", h.Update)
	w := doRequest(r, "PUT", "/projects/abc", map[string]interface{}{"title": "t"})

	assertStatus(t, w, 400)
}

func TestProjectUpdate_FeaturedError(t *testing.T) {
	h, svc := setupProjectHandler()
	updated := &model.Project{Title: "Updated"}
	updated.ID = 1
	svc.On("Update", uint(1), uint(1), mock.AnythingOfType("*model.Project")).Return(updated, nil)
	svc.On("UpdateFeatured", uint(1), uint(1), true).Return(nil, service.ErrForbidden)

	r := newRouter(1)
	r.PUT("/projects/:id", h.Update)
	w := doRequest(r, "PUT", "/projects/1", map[string]interface{}{
		"title":    "Updated",
		"featured": true,
	})

	assertStatus(t, w, 403)
	svc.AssertExpectations(t)
}

func TestProjectUpdate_WithDates(t *testing.T) {
	h, svc := setupProjectHandler()
	updated := &model.Project{Title: "Updated"}
	updated.ID = 1
	svc.On("Update", uint(1), uint(1), mock.AnythingOfType("*model.Project")).Return(updated, nil)

	r := newRouter(1)
	r.PUT("/projects/:id", h.Update)
	w := doRequest(r, "PUT", "/projects/1", map[string]interface{}{
		"title":      "Updated",
		"start_date": "2024-01-01",
		"end_date":   "2024-12-31",
	})

	assertStatus(t, w, 200)
	svc.AssertExpectations(t)
}

// ---------- Delete (エラーパス) ----------

func TestProjectDelete_InvalidID(t *testing.T) {
	h, _ := setupProjectHandler()

	r := newRouter(1)
	r.DELETE("/projects/:id", h.Delete)
	w := doRequest(r, "DELETE", "/projects/abc", nil)

	assertStatus(t, w, 400)
}

func TestProjectCreate_InvalidDemoURL(t *testing.T) {
	h, _ := setupProjectHandler()

	r := newRouter(1)
	r.POST("/projects", h.Create)
	w := doRequest(r, "POST", "/projects", map[string]interface{}{
		"title":    "Test Project",
		"demo_url": "javascript:alert('xss')",
	})

	assertStatus(t, w, 400)
}

func TestProjectUpdate_InvalidGithubURL(t *testing.T) {
	h, _ := setupProjectHandler()

	r := newRouter(1)
	r.PUT("/projects/:id", h.Update)
	w := doRequest(r, "PUT", "/projects/1", map[string]interface{}{
		"github_url": "data:text/html,<script>alert('xss')</script>",
	})

	assertStatus(t, w, 400)
}

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

func TestQuestionCreate_Success(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.POST("/questions", h.Create)

	repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Question")).Return(nil)

	w := doRequest(r, http.MethodPost, "/questions", map[string]string{
		"title": "How to use Go?", "body": "I want to learn Go.",
	})
	assertStatus(t, w, http.StatusCreated)
}

func TestQuestionCreate_ValidationError(t *testing.T) {
	h, _ := setupQuestionHandler()
	r := newRouter(1)
	r.POST("/questions", h.Create)

	// title と body は required
	w := doRequest(r, http.MethodPost, "/questions", map[string]string{
		"title": "No body",
	})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestQuestionCreate_InvalidJSON(t *testing.T) {
	h, _ := setupQuestionHandler()
	r := newRouter(1)
	r.POST("/questions", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/questions", "not json")
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- GetAll ----------

func TestQuestionGetAll_Success(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.GET("/questions", h.GetAll)

	repo.On("FindAll", mock.Anything, 20, 0, "", "newest").Return([]model.Question{
		{Title: "Q1"}, {Title: "Q2"},
	}, int64(2), nil)

	w := doRequest(r, http.MethodGet, "/questions", nil)
	assertStatus(t, w, http.StatusOK)

	body := parseJSON(t, w)
	questions := body["questions"].([]interface{})
	assert.Len(t, questions, 2)
	assert.Equal(t, float64(2), body["total"])
}

func TestQuestionGetAll_WithFilters(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.GET("/questions", h.GetAll)

	repo.On("FindAll", mock.Anything, 10, 5, "go", "popular").Return([]model.Question{}, int64(0), nil)

	w := doRequest(r, http.MethodGet, "/questions?limit=10&offset=5&tag=go&sort=popular", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestQuestionGetAll_LimitCap(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.GET("/questions", h.GetAll)

	// limit=200 は 100 に制限される
	repo.On("FindAll", mock.Anything, 100, 0, "", "newest").Return([]model.Question{}, int64(0), nil)

	w := doRequest(r, http.MethodGet, "/questions?limit=200", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- Search ----------

func TestQuestionSearch_Success(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.GET("/questions/search", h.Search)

	repo.On("Search", mock.Anything, "golang", 20, 0).Return([]model.Question{
		{Title: "Go Question"},
	}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/questions/search?q=golang", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestQuestionSearch_EmptyQuery(t *testing.T) {
	h, _ := setupQuestionHandler()
	r := newRouter(1)
	r.GET("/questions/search", h.Search)

	w := doRequest(r, http.MethodGet, "/questions/search", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- GetByID ----------

func TestQuestionGetByID_Success(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.GET("/questions/:id", h.GetByID)

	repo.On("FindByID", mock.Anything, uint(5)).Return(&model.Question{Title: "Q"}, nil)
	repo.On("GetUserVote", mock.Anything, uint(1), uint(5)).Return(1, nil)

	w := doRequest(r, http.MethodGet, "/questions/5", nil)
	assertStatus(t, w, http.StatusOK)

	body := parseJSON(t, w)
	assert.Equal(t, float64(1), body["user_vote"])
}

// 単体取得は不在でも 404 にならない。リポジトリのエラーがそのまま返るため 500 になる（移行前からの挙動）。
func TestQuestionGetByID_MissingReturnsInternalError(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.GET("/questions/:id", h.GetByID)

	// port は不在を (nil, nil) で表す。
	repo.On("FindByID", mock.Anything, uint(999)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/questions/999", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestQuestionGetByID_InvalidID(t *testing.T) {
	h, _ := setupQuestionHandler()
	r := newRouter(1)
	r.GET("/questions/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/questions/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- Update ----------

func TestQuestionUpdate_Success(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.PUT("/questions/:id", h.Update)

	q := &model.Question{Title: "Old", Body: "Old Body"}
	q.ID = 5
	q.UserID = 1
	repo.On("FindByID", mock.Anything, uint(5)).Return(q, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*model.Question")).Return(nil)

	w := doRequest(r, http.MethodPut, "/questions/5", map[string]string{
		"title": "Updated",
	})
	assertStatus(t, w, http.StatusOK)
}

func TestQuestionUpdate_Forbidden(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.PUT("/questions/:id", h.Update)

	q := &model.Question{Title: "Other"}
	q.ID = 5
	q.UserID = 999
	repo.On("FindByID", mock.Anything, uint(5)).Return(q, nil)

	w := doRequest(r, http.MethodPut, "/questions/5", map[string]string{
		"title": "Hacked",
	})
	assertStatus(t, w, http.StatusForbidden)
}

// 更新も不在は 404 にならず 500 になる（移行前からの挙動）。
func TestQuestionUpdate_MissingReturnsInternalError(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.PUT("/questions/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(5)).Return(nil, nil)

	w := doRequest(r, http.MethodPut, "/questions/5", map[string]string{"title": "X"})
	assertStatus(t, w, http.StatusNotFound)
}

func TestQuestionUpdate_RepositoryError(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.PUT("/questions/:id", h.Update)

	repo.On("FindByID", mock.Anything, uint(5)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodPut, "/questions/5", map[string]string{"title": "X"})
	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- Delete ----------

func TestQuestionDelete_Success(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.DELETE("/questions/:id", h.Delete)

	q := &model.Question{}
	q.ID = 5
	q.UserID = 1
	repo.On("FindByID", mock.Anything, uint(5)).Return(q, nil)
	repo.On("Delete", mock.Anything, uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/questions/5", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestQuestionDelete_Forbidden(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.DELETE("/questions/:id", h.Delete)

	q := &model.Question{}
	q.ID = 5
	q.UserID = 999
	repo.On("FindByID", mock.Anything, uint(5)).Return(q, nil)

	w := doRequest(r, http.MethodDelete, "/questions/5", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- Vote / RemoveVote ----------

func TestQuestionVote_Success(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.POST("/questions/:id/vote", h.Vote)

	otherQuestion := &model.Question{UserID: 99}
	otherQuestion.ID = 5
	repo.On("FindByID", mock.Anything, uint(5)).Return(otherQuestion, nil)
	repo.On("Vote", mock.Anything, uint(1), uint(5), 1).Return(nil)

	w := doRequest(r, http.MethodPost, "/questions/5/vote", map[string]int{
		"value": 1,
	})
	assertStatus(t, w, http.StatusOK)
}

func TestQuestionVote_InvalidValue(t *testing.T) {
	h, _ := setupQuestionHandler()
	r := newRouter(1)
	r.POST("/questions/:id/vote", h.Vote)

	// value は 1 or -1 のみ
	w := doRequest(r, http.MethodPost, "/questions/5/vote", map[string]int{
		"value": 5,
	})
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// GetByUserID テスト
// ============================================================

func TestQuestion_GetByUserID_Success(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.GET("/users/:userId/questions", h.GetByUserID)

	questions := []model.Question{{Title: "Go質問"}}
	repo.On("FindByUserID", mock.Anything, uint(5), 20, 0).Return(questions, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/users/5/questions", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestQuestion_GetByUserID_Empty(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.GET("/users/:userId/questions", h.GetByUserID)

	repo.On("FindByUserID", mock.Anything, uint(99), 20, 0).Return([]model.Question{}, int64(0), nil)

	w := doRequest(r, http.MethodGet, "/users/99/questions", nil)
	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestQuestion_GetByUserID_InvalidID(t *testing.T) {
	h, _ := setupQuestionHandler()
	r := newRouter(1)
	r.GET("/users/:userId/questions", h.GetByUserID)

	w := doRequest(r, http.MethodGet, "/users/abc/questions", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- GetSolved ----------

func TestQuestionGetSolved_Success(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.GET("/questions/solved", h.GetSolved)

	repo.On("FindSolved", mock.Anything, 20, 0).Return(
		[]model.Question{{Title: "Solved Q"}},
		int64(1), nil,
	)

	w := doRequest(r, http.MethodGet, "/questions/solved", nil)
	assertStatus(t, w, http.StatusOK)

	body := parseJSON(t, w)
	assert.Equal(t, float64(1), body["total"])
}

func TestQuestionGetSolved_ServiceError(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.GET("/questions/solved", h.GetSolved)

	repo.On("FindSolved", mock.Anything, 20, 0).Return(
		[]model.Question{}, int64(0), errors.New("db error"),
	)

	w := doRequest(r, http.MethodGet, "/questions/solved", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestQuestionRemoveVote_Success(t *testing.T) {
	h, repo := setupQuestionHandler()
	r := newRouter(1)
	r.DELETE("/questions/:id/vote", h.RemoveVote)

	otherQuestion := &model.Question{UserID: 99}
	otherQuestion.ID = 5
	repo.On("FindByID", mock.Anything, uint(5)).Return(otherQuestion, nil)
	repo.On("RemoveVote", mock.Anything, uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/questions/5/vote", nil)
	assertStatus(t, w, http.StatusOK)
}

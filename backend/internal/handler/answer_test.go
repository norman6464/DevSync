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

// ============================================================
// GetByQuestionID テスト
// ============================================================

func TestAnswerGetByQuestionID_Success(t *testing.T) {
	h, svc := setupAnswerHandler()
	r := newRouter(1)
	r.GET("/questions/:id/answers", h.GetByQuestionID)

	answers := []model.Answer{
		{Body: "回答1"},
		{Body: "回答2"},
	}
	svc.On("GetByQuestionID", uint(1)).Return(answers, nil)

	w := doRequest(r, http.MethodGet, "/questions/1/answers", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestAnswerGetByQuestionID_InvalidID(t *testing.T) {
	h, _ := setupAnswerHandler()
	r := newRouter(1)
	r.GET("/questions/:id/answers", h.GetByQuestionID)

	w := doRequest(r, http.MethodGet, "/questions/abc/answers", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// Create テスト
// ============================================================

func TestAnswerCreate_Success(t *testing.T) {
	h, svc := setupAnswerHandler()
	r := newRouter(1)
	r.POST("/questions/:id/answers", h.Create)

	svc.On("Create", mock.AnythingOfType("*model.Answer")).Return(nil)

	w := doRequest(r, http.MethodPost, "/questions/1/answers", map[string]interface{}{
		"body": "これは回答です",
	})
	assertStatus(t, w, http.StatusCreated)
	svc.AssertExpectations(t)
}

func TestAnswerCreate_ValidationError(t *testing.T) {
	h, _ := setupAnswerHandler()
	r := newRouter(1)
	r.POST("/questions/:id/answers", h.Create)

	// body は required
	w := doRequest(r, http.MethodPost, "/questions/1/answers", map[string]interface{}{})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestAnswerCreate_InvalidJSON(t *testing.T) {
	h, _ := setupAnswerHandler()
	r := newRouter(1)
	r.POST("/questions/:id/answers", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/questions/1/answers", "not json")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestAnswerCreate_ServiceError(t *testing.T) {
	h, svc := setupAnswerHandler()
	r := newRouter(1)
	r.POST("/questions/:id/answers", h.Create)

	svc.On("Create", mock.AnythingOfType("*model.Answer")).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/questions/1/answers", map[string]interface{}{
		"body": "回答",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

// ============================================================
// Update テスト
// ============================================================

func TestAnswerUpdate_Success(t *testing.T) {
	h, svc := setupAnswerHandler()
	r := newRouter(1)
	r.PUT("/questions/:id/answers/:answerId", h.Update)

	updated := &model.Answer{Body: "更新後の回答"}
	updated.ID = 10
	svc.On("Update", uint(10), uint(1), "更新後の回答").Return(updated, nil)

	w := doRequest(r, http.MethodPut, "/questions/1/answers/10", map[string]interface{}{
		"body": "更新後の回答",
	})
	assertStatus(t, w, http.StatusOK)

	body := parseJSON(t, w)
	assert.Equal(t, "更新後の回答", body["body"])
	svc.AssertExpectations(t)
}

func TestAnswerUpdate_Forbidden(t *testing.T) {
	h, svc := setupAnswerHandler()
	r := newRouter(1)
	r.PUT("/questions/:id/answers/:answerId", h.Update)

	svc.On("Update", uint(10), uint(1), "不正な更新").Return(nil, service.ErrForbidden)

	w := doRequest(r, http.MethodPut, "/questions/1/answers/10", map[string]interface{}{
		"body": "不正な更新",
	})
	assertStatus(t, w, http.StatusForbidden)
	svc.AssertExpectations(t)
}

func TestAnswerUpdate_NotFound(t *testing.T) {
	h, svc := setupAnswerHandler()
	r := newRouter(1)
	r.PUT("/questions/:id/answers/:answerId", h.Update)

	svc.On("Update", uint(999), uint(1), "存在しない").Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodPut, "/questions/1/answers/999", map[string]interface{}{
		"body": "存在しない",
	})
	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

func TestAnswerUpdate_InvalidID(t *testing.T) {
	h, _ := setupAnswerHandler()
	r := newRouter(1)
	r.PUT("/questions/:id/answers/:answerId", h.Update)

	w := doRequest(r, http.MethodPut, "/questions/1/answers/abc", map[string]interface{}{
		"body": "テスト",
	})
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// Delete テスト
// ============================================================

func TestAnswerDelete_Success(t *testing.T) {
	h, svc := setupAnswerHandler()
	r := newRouter(1)
	r.DELETE("/questions/:id/answers/:answerId", h.Delete)

	svc.On("Delete", uint(10), uint(1)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/questions/1/answers/10", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestAnswerDelete_Forbidden(t *testing.T) {
	h, svc := setupAnswerHandler()
	r := newRouter(1)
	r.DELETE("/questions/:id/answers/:answerId", h.Delete)

	svc.On("Delete", uint(10), uint(1)).Return(service.ErrForbidden)

	w := doRequest(r, http.MethodDelete, "/questions/1/answers/10", nil)
	assertStatus(t, w, http.StatusForbidden)
	svc.AssertExpectations(t)
}

func TestAnswerDelete_NotFound(t *testing.T) {
	h, svc := setupAnswerHandler()
	r := newRouter(1)
	r.DELETE("/questions/:id/answers/:answerId", h.Delete)

	svc.On("Delete", uint(999), uint(1)).Return(service.ErrNotFound)

	w := doRequest(r, http.MethodDelete, "/questions/1/answers/999", nil)
	assertStatus(t, w, http.StatusNotFound)
	svc.AssertExpectations(t)
}

// ============================================================
// SetBestAnswer テスト
// ============================================================

func TestAnswerSetBestAnswer_Success(t *testing.T) {
	h, svc := setupAnswerHandler()
	r := newRouter(1)
	r.PUT("/questions/:id/answers/:answerId/best", h.SetBestAnswer)

	svc.On("SetBestAnswer", uint(1), uint(10), uint(1)).Return(nil)

	w := doRequest(r, http.MethodPut, "/questions/1/answers/10/best", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestAnswerSetBestAnswer_Forbidden(t *testing.T) {
	h, svc := setupAnswerHandler()
	r := newRouter(1)
	r.PUT("/questions/:id/answers/:answerId/best", h.SetBestAnswer)

	svc.On("SetBestAnswer", uint(1), uint(10), uint(1)).Return(service.ErrForbidden)

	w := doRequest(r, http.MethodPut, "/questions/1/answers/10/best", nil)
	assertStatus(t, w, http.StatusForbidden)
	svc.AssertExpectations(t)
}

// ============================================================
// Vote テスト
// ============================================================

func TestAnswerVote_Success(t *testing.T) {
	h, svc := setupAnswerHandler()
	r := newRouter(1)
	r.POST("/questions/:id/answers/:answerId/vote", h.Vote)

	svc.On("Vote", uint(1), uint(10), 1).Return(nil)

	w := doRequest(r, http.MethodPost, "/questions/1/answers/10/vote", map[string]interface{}{
		"value": 1,
	})
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestAnswerVote_ValidationError(t *testing.T) {
	h, _ := setupAnswerHandler()
	r := newRouter(1)
	r.POST("/questions/:id/answers/:answerId/vote", h.Vote)

	// value は required で 1 or -1 のみ許可
	w := doRequest(r, http.MethodPost, "/questions/1/answers/10/vote", map[string]interface{}{
		"value": 0,
	})
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// RemoveVote テスト
// ============================================================

func TestAnswerRemoveVote_Success(t *testing.T) {
	h, svc := setupAnswerHandler()
	r := newRouter(1)
	r.DELETE("/questions/:id/answers/:answerId/vote", h.RemoveVote)

	svc.On("RemoveVote", uint(1), uint(10)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/questions/1/answers/10/vote", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

// ============================================================
// GetByVoteRange テスト
// ============================================================

func TestAnswerGetByVoteRange_Success(t *testing.T) {
	h, svc := setupAnswerHandler()
	r := newRouter(1)
	r.GET("/questions/:id/answers/vote-range", h.GetByVoteRange)

	answers := []model.Answer{
		{Body: "回答A"},
	}
	svc.On("GetByVoteRange", uint(1), 5, 10).Return(answers, nil)

	w := doRequest(r, http.MethodGet, "/questions/1/answers/vote-range?min_vote=5&max_vote=10", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestAnswerGetByVoteRange_DefaultValues(t *testing.T) {
	h, svc := setupAnswerHandler()
	r := newRouter(1)
	r.GET("/questions/:id/answers/vote-range", h.GetByVoteRange)

	svc.On("GetByVoteRange", uint(1), 0, 100).Return([]model.Answer{}, nil)

	w := doRequest(r, http.MethodGet, "/questions/1/answers/vote-range", nil)
	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestAnswerGetByVoteRange_InvalidMinVote(t *testing.T) {
	h, _ := setupAnswerHandler()
	r := newRouter(1)
	r.GET("/questions/:id/answers/vote-range", h.GetByVoteRange)

	w := doRequest(r, http.MethodGet, "/questions/1/answers/vote-range?min_vote=abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestAnswerGetByVoteRange_InvalidID(t *testing.T) {
	h, _ := setupAnswerHandler()
	r := newRouter(1)
	r.GET("/questions/:id/answers/vote-range", h.GetByVoteRange)

	w := doRequest(r, http.MethodGet, "/questions/abc/answers/vote-range", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// RemoveVote テスト
// ============================================================

func TestAnswerRemoveVote_ServiceError(t *testing.T) {
	h, svc := setupAnswerHandler()
	r := newRouter(1)
	r.DELETE("/answers/:answerId/vote", h.RemoveVote)

	svc.On("RemoveVote", uint(1), uint(5)).Return(errors.New("not found"))

	w := doRequest(r, http.MethodDelete, "/answers/5/vote", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestAnswerRemoveVote_InvalidID(t *testing.T) {
	h, _ := setupAnswerHandler()
	r := newRouter(1)
	r.DELETE("/answers/:answerId/vote", h.RemoveVote)

	w := doRequest(r, http.MethodDelete, "/answers/abc/vote", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestAnswerGetByVoteRange_ServiceError(t *testing.T) {
	h, svc := setupAnswerHandler()
	r := newRouter(1)
	r.GET("/questions/:id/answers/vote-range", h.GetByVoteRange)

	svc.On("GetByVoteRange", uint(1), 0, 100).Return([]model.Answer(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/questions/1/answers/vote-range", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

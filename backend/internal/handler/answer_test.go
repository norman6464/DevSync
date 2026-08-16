package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ownedBy は指定ユーザーが投稿した回答を返すテスト用ヘルパー。
func ownedBy(id, userID, questionID uint) *model.Answer {
	return &model.Answer{ID: id, UserID: userID, QuestionID: questionID, Body: "既存の回答"}
}

// ============================================================
// GetByQuestionID テスト
// ============================================================

func TestAnswerGetByQuestionID_Success(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.GET("/questions/:id/answers", h.GetByQuestionID)

	ports.Answers.On("FindByQuestionID", mock.Anything, uint(1)).
		Return([]model.Answer{{Body: "回答1"}, {Body: "回答2"}}, nil)

	w := doRequest(r, http.MethodGet, "/questions/1/answers", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Answers.AssertExpectations(t)
}

func TestAnswerGetByQuestionID_InvalidID(t *testing.T) {
	h, _ := setupAnswerHandler()
	r := newRouter(1)
	r.GET("/questions/:id/answers", h.GetByQuestionID)

	w := doRequest(r, http.MethodGet, "/questions/abc/answers", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestAnswerGetByQuestionID_RepositoryError(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.GET("/questions/:id/answers", h.GetByQuestionID)

	ports.Answers.On("FindByQuestionID", mock.Anything, uint(1)).
		Return([]model.Answer(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/questions/1/answers", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// Create テスト
// ============================================================

func TestAnswerCreate_Success(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.POST("/questions/:id/answers", h.Create)

	ports.Questions.On("FindByID", mock.Anything, uint(1)).Return(&model.Question{ID: 1}, nil)
	ports.Answers.On("Create", mock.Anything, mock.AnythingOfType("*model.Answer")).Return(nil)

	w := doRequest(r, http.MethodPost, "/questions/1/answers", map[string]interface{}{
		"body": "これは回答です",
	})
	assertStatus(t, w, http.StatusCreated)
	ports.Answers.AssertExpectations(t)
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

// 質問が存在しなければ回答は作成されない。
func TestAnswerCreate_QuestionNotFound(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.POST("/questions/:id/answers", h.Create)

	// port は不在を (nil, nil) で表す。
	ports.Questions.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodPost, "/questions/1/answers", map[string]interface{}{"body": "回答"})
	assertStatus(t, w, http.StatusNotFound)
	ports.Answers.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestAnswerCreate_RepositoryError(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.POST("/questions/:id/answers", h.Create)

	ports.Questions.On("FindByID", mock.Anything, uint(1)).Return(&model.Question{ID: 1}, nil)
	ports.Answers.On("Create", mock.Anything, mock.AnythingOfType("*model.Answer")).
		Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/questions/1/answers", map[string]interface{}{"body": "回答"})
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// Update テスト
// ============================================================

func TestAnswerUpdate_Success(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.PUT("/questions/:id/answers/:answerId", h.Update)

	ports.Answers.On("FindByID", mock.Anything, uint(10)).Return(ownedBy(10, 1, 1), nil)
	ports.Answers.On("Update", mock.Anything, mock.AnythingOfType("*model.Answer")).Return(nil)

	w := doRequest(r, http.MethodPut, "/questions/1/answers/10", map[string]interface{}{
		"body": "更新後の回答",
	})
	assertStatus(t, w, http.StatusOK)

	body := parseJSON(t, w)
	assert.Equal(t, "更新後の回答", body["body"])
	ports.Answers.AssertExpectations(t)
}

func TestAnswerUpdate_Forbidden(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.PUT("/questions/:id/answers/:answerId", h.Update)

	ports.Answers.On("FindByID", mock.Anything, uint(10)).Return(ownedBy(10, 999, 1), nil)

	w := doRequest(r, http.MethodPut, "/questions/1/answers/10", map[string]interface{}{
		"body": "不正な更新",
	})
	assertStatus(t, w, http.StatusForbidden)
	ports.Answers.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// 不在の回答の更新は 404 にならず 500 になる（移行前からの挙動）。
func TestAnswerUpdate_MissingReturnsInternalError(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.PUT("/questions/:id/answers/:answerId", h.Update)

	ports.Answers.On("FindByID", mock.Anything, uint(999)).Return(nil, nil)

	w := doRequest(r, http.MethodPut, "/questions/1/answers/999", map[string]interface{}{
		"body": "存在しない",
	})
	assertStatus(t, w, http.StatusNotFound)
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
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.DELETE("/questions/:id/answers/:answerId", h.Delete)

	ports.Answers.On("FindByID", mock.Anything, uint(10)).Return(ownedBy(10, 1, 1), nil)
	ports.Answers.On("Delete", mock.Anything, mock.AnythingOfType("*model.Answer")).Return(nil)

	w := doRequest(r, http.MethodDelete, "/questions/1/answers/10", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Answers.AssertExpectations(t)
}

func TestAnswerDelete_Forbidden(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.DELETE("/questions/:id/answers/:answerId", h.Delete)

	ports.Answers.On("FindByID", mock.Anything, uint(10)).Return(ownedBy(10, 999, 1), nil)

	w := doRequest(r, http.MethodDelete, "/questions/1/answers/10", nil)
	assertStatus(t, w, http.StatusForbidden)
	ports.Answers.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

// 不在の回答の削除も 404 にならず 500 になる（移行前からの挙動）。
func TestAnswerDelete_MissingReturnsInternalError(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.DELETE("/questions/:id/answers/:answerId", h.Delete)

	ports.Answers.On("FindByID", mock.Anything, uint(999)).Return(nil, nil)

	w := doRequest(r, http.MethodDelete, "/questions/1/answers/999", nil)
	assertStatus(t, w, http.StatusNotFound)
}

// ============================================================
// SetBestAnswer テスト
// ============================================================

func TestAnswerSetBestAnswer_Success(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.PUT("/questions/:id/answers/:answerId/best", h.SetBestAnswer)

	ports.Questions.On("FindByID", mock.Anything, uint(1)).Return(&model.Question{ID: 1, UserID: 1}, nil)
	ports.Answers.On("FindByID", mock.Anything, uint(10)).Return(ownedBy(10, 2, 1), nil)
	ports.Answers.On("SetBestAnswer", mock.Anything, uint(1), uint(10)).Return(nil)

	w := doRequest(r, http.MethodPut, "/questions/1/answers/10/best", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Answers.AssertExpectations(t)
}

// 質問の投稿者でなければベストアンサーは設定できない。
func TestAnswerSetBestAnswer_Forbidden(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.PUT("/questions/:id/answers/:answerId/best", h.SetBestAnswer)

	ports.Questions.On("FindByID", mock.Anything, uint(1)).Return(&model.Question{ID: 1, UserID: 999}, nil)

	w := doRequest(r, http.MethodPut, "/questions/1/answers/10/best", nil)
	assertStatus(t, w, http.StatusForbidden)
	ports.Answers.AssertNotCalled(t, "SetBestAnswer", mock.Anything, mock.Anything, mock.Anything)
}

func TestAnswerSetBestAnswer_QuestionNotFound(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.PUT("/questions/:id/answers/:answerId/best", h.SetBestAnswer)

	ports.Questions.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

	w := doRequest(r, http.MethodPut, "/questions/1/answers/10/best", nil)
	assertStatus(t, w, http.StatusNotFound)
}

// 別の質問に属する回答を指定した場合は 400。
func TestAnswerSetBestAnswer_AnswerBelongsToOtherQuestion(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.PUT("/questions/:id/answers/:answerId/best", h.SetBestAnswer)

	ports.Questions.On("FindByID", mock.Anything, uint(1)).Return(&model.Question{ID: 1, UserID: 1}, nil)
	ports.Answers.On("FindByID", mock.Anything, uint(10)).Return(ownedBy(10, 2, 77), nil)

	w := doRequest(r, http.MethodPut, "/questions/1/answers/10/best", nil)
	assertStatus(t, w, http.StatusBadRequest)
	ports.Answers.AssertNotCalled(t, "SetBestAnswer", mock.Anything, mock.Anything, mock.Anything)
}

func TestAnswerSetBestAnswer_InvalidQuestionID(t *testing.T) {
	h, _ := setupAnswerHandler()
	r := newRouter(1)
	r.PUT("/questions/:id/best-answer/:answerId", h.SetBestAnswer)

	w := doRequest(r, http.MethodPut, "/questions/abc/best-answer/5", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestAnswerSetBestAnswer_InvalidAnswerID(t *testing.T) {
	h, _ := setupAnswerHandler()
	r := newRouter(1)
	r.PUT("/questions/:id/best-answer/:answerId", h.SetBestAnswer)

	w := doRequest(r, http.MethodPut, "/questions/1/best-answer/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestAnswerSetBestAnswer_RepositoryError(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.PUT("/questions/:id/best-answer/:answerId", h.SetBestAnswer)

	ports.Questions.On("FindByID", mock.Anything, uint(1)).Return(&model.Question{ID: 1, UserID: 1}, nil)
	ports.Answers.On("FindByID", mock.Anything, uint(5)).Return(ownedBy(5, 2, 1), nil)
	ports.Answers.On("SetBestAnswer", mock.Anything, uint(1), uint(5)).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPut, "/questions/1/best-answer/5", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ============================================================
// Vote テスト
// ============================================================

func TestAnswerVote_Success(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.POST("/questions/:id/answers/:answerId/vote", h.Vote)

	ports.Answers.On("FindByID", mock.Anything, uint(10)).Return(ownedBy(10, 2, 1), nil)
	ports.Answers.On("Vote", mock.Anything, uint(1), uint(10), 1).Return(nil)

	w := doRequest(r, http.MethodPost, "/questions/1/answers/10/vote", map[string]interface{}{
		"value": 1,
	})
	assertStatus(t, w, http.StatusOK)
	ports.Answers.AssertExpectations(t)
}

func TestAnswerVote_ValidationError(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.POST("/questions/:id/answers/:answerId/vote", h.Vote)

	// value は 1 か -1 のみ許可
	w := doRequest(r, http.MethodPost, "/questions/1/answers/10/vote", map[string]interface{}{
		"value": 0,
	})
	assertStatus(t, w, http.StatusBadRequest)
	ports.Answers.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// 自分の回答には投票できない。
func TestAnswerVote_SelfVoteForbidden(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.POST("/questions/:id/answers/:answerId/vote", h.Vote)

	ports.Answers.On("FindByID", mock.Anything, uint(10)).Return(ownedBy(10, 1, 1), nil)

	w := doRequest(r, http.MethodPost, "/questions/1/answers/10/vote", map[string]interface{}{"value": 1})
	assertStatus(t, w, http.StatusForbidden)
	ports.Answers.AssertNotCalled(t, "Vote", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestAnswerVote_AnswerNotFound(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.POST("/questions/:id/answers/:answerId/vote", h.Vote)

	ports.Answers.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)

	w := doRequest(r, http.MethodPost, "/questions/1/answers/10/vote", map[string]interface{}{"value": 1})
	assertStatus(t, w, http.StatusNotFound)
}

func TestAnswerVote_RepositoryError(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.POST("/answers/:answerId/vote", h.Vote)

	ports.Answers.On("FindByID", mock.Anything, uint(5)).Return(ownedBy(5, 2, 1), nil)
	ports.Answers.On("Vote", mock.Anything, uint(1), uint(5), 1).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/answers/5/vote", map[string]int{"value": 1})
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestAnswerVote_InvalidJSON(t *testing.T) {
	h, _ := setupAnswerHandler()
	r := newRouter(1)
	r.POST("/answers/:answerId/vote", h.Vote)

	w := doRequestRaw(r, http.MethodPost, "/answers/5/vote", "invalid json")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestAnswerVote_InvalidID(t *testing.T) {
	h, _ := setupAnswerHandler()
	r := newRouter(1)
	r.POST("/answers/:answerId/vote", h.Vote)

	w := doRequest(r, http.MethodPost, "/answers/abc/vote", map[string]int{"value": 1})
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// RemoveVote テスト
// ============================================================

func TestAnswerRemoveVote_Success(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.DELETE("/questions/:id/answers/:answerId/vote", h.RemoveVote)

	ports.Answers.On("FindByID", mock.Anything, uint(10)).Return(ownedBy(10, 2, 1), nil)
	ports.Answers.On("RemoveVote", mock.Anything, uint(1), uint(10)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/questions/1/answers/10/vote", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Answers.AssertExpectations(t)
}

// 自分の回答はそもそも投票できないため、取消も 403。
func TestAnswerRemoveVote_SelfVoteForbidden(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.DELETE("/answers/:answerId/vote", h.RemoveVote)

	ports.Answers.On("FindByID", mock.Anything, uint(5)).Return(ownedBy(5, 1, 1), nil)

	w := doRequest(r, http.MethodDelete, "/answers/5/vote", nil)
	assertStatus(t, w, http.StatusForbidden)
}

func TestAnswerRemoveVote_RepositoryError(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.DELETE("/answers/:answerId/vote", h.RemoveVote)

	ports.Answers.On("FindByID", mock.Anything, uint(5)).Return(ownedBy(5, 2, 1), nil)
	ports.Answers.On("RemoveVote", mock.Anything, uint(1), uint(5)).Return(errors.New("not found"))

	w := doRequest(r, http.MethodDelete, "/answers/5/vote", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestAnswerRemoveVote_InvalidID(t *testing.T) {
	h, _ := setupAnswerHandler()
	r := newRouter(1)
	r.DELETE("/answers/:answerId/vote", h.RemoveVote)

	w := doRequest(r, http.MethodDelete, "/answers/abc/vote", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// GetByVoteRange テスト
// ============================================================

func TestAnswerGetByVoteRange_Success(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.GET("/questions/:id/answers/vote-range", h.GetByVoteRange)

	ports.Answers.On("FindByVoteRange", mock.Anything, uint(1), 5, 10).
		Return([]model.Answer{{Body: "回答A"}}, nil)

	w := doRequest(r, http.MethodGet, "/questions/1/answers/vote-range?min_vote=5&max_vote=10", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Answers.AssertExpectations(t)
}

func TestAnswerGetByVoteRange_DefaultValues(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.GET("/questions/:id/answers/vote-range", h.GetByVoteRange)

	ports.Answers.On("FindByVoteRange", mock.Anything, uint(1), 0, 100).Return([]model.Answer{}, nil)

	w := doRequest(r, http.MethodGet, "/questions/1/answers/vote-range", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Answers.AssertExpectations(t)
}

// 下限が上限を上回る場合はリポジトリを引かずに 400。
func TestAnswerGetByVoteRange_InvertedRange(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.GET("/questions/:id/answers/vote-range", h.GetByVoteRange)

	w := doRequest(r, http.MethodGet, "/questions/1/answers/vote-range?min_vote=10&max_vote=5", nil)
	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "投票範囲が無効です")
	ports.Answers.AssertNotCalled(t, "FindByVoteRange", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
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

func TestAnswerGetByVoteRange_RepositoryError(t *testing.T) {
	h, ports := setupAnswerHandler()
	r := newRouter(1)
	r.GET("/questions/:id/answers/vote-range", h.GetByVoteRange)

	ports.Answers.On("FindByVoteRange", mock.Anything, uint(1), 0, 100).
		Return([]model.Answer(nil), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/questions/1/answers/vote-range", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

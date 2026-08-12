package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------- AddReaction ----------

func TestPostAddReaction_Success(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.POST("/posts/:id/reactions", h.AddReaction)

	// 他人の投稿にはリアクションできる。
	reactions.Authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(99), nil)
	reactions.Reactions.On("AddReaction", mock.Anything, uint(1), uint(5), "👍").Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/reactions", map[string]string{"emoji": "👍"})
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "reaction_added")
	reactions.Reactions.AssertExpectations(t)
	reactions.Authors.AssertExpectations(t)
}

// 自分の投稿にはリアクションできない。
func TestPostAddReaction_OwnPost(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.POST("/posts/:id/reactions", h.AddReaction)

	reactions.Authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(1), nil)

	w := doRequest(r, http.MethodPost, "/posts/5/reactions", map[string]string{"emoji": "👍"})
	assertStatus(t, w, http.StatusForbidden)
	reactions.Reactions.AssertNotCalled(t, "AddReaction", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

// 投稿が存在しなければ 404。
func TestPostAddReaction_PostNotFound(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.POST("/posts/:id/reactions", h.AddReaction)

	reactions.Authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(0), nil)

	w := doRequest(r, http.MethodPost, "/posts/5/reactions", map[string]string{"emoji": "👍"})
	assertStatus(t, w, http.StatusNotFound)
	reactions.Reactions.AssertNotCalled(t, "AddReaction", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

// 許可されていない絵文字は 400。投稿の取得も行わない。
func TestPostAddReaction_InvalidEmoji(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.POST("/posts/:id/reactions", h.AddReaction)

	w := doRequest(r, http.MethodPost, "/posts/5/reactions", map[string]string{"emoji": "🍣"})
	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "許可されていない絵文字です")
	reactions.Authors.AssertNotCalled(t, "FindAuthorID", mock.Anything, mock.Anything)
}

func TestPostAddReaction_InvalidID(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.POST("/posts/:id/reactions", h.AddReaction)

	w := doRequest(r, http.MethodPost, "/posts/abc/reactions", map[string]string{"emoji": "👍"})
	assertStatus(t, w, http.StatusBadRequest)
	reactions.Authors.AssertNotCalled(t, "FindAuthorID", mock.Anything, mock.Anything)
}

func TestPostAddReaction_InvalidJSON(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.POST("/posts/:id/reactions", h.AddReaction)

	w := doRequestRaw(r, http.MethodPost, "/posts/5/reactions", "{invalid json}")
	assertStatus(t, w, http.StatusBadRequest)
	reactions.Authors.AssertNotCalled(t, "FindAuthorID", mock.Anything, mock.Anything)
}

func TestPostAddReaction_RepositoryError(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.POST("/posts/:id/reactions", h.AddReaction)

	reactions.Authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(99), nil)
	reactions.Reactions.On("AddReaction", mock.Anything, uint(1), uint(5), "👍").Return(errors.New("duplicate"))

	w := doRequest(r, http.MethodPost, "/posts/5/reactions", map[string]string{"emoji": "👍"})
	assertStatus(t, w, http.StatusInternalServerError)
	reactions.Reactions.AssertExpectations(t)
}

// ---------- RemoveReaction ----------

func TestPostRemoveReaction_Success(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.DELETE("/posts/:id/reactions", h.RemoveReaction)

	reactions.Authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(99), nil)
	reactions.Reactions.On("RemoveReaction", mock.Anything, uint(1), uint(5), "👍").Return(nil)

	w := doRequest(r, http.MethodDelete, "/posts/5/reactions", map[string]string{"emoji": "👍"})
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "reaction_removed")
	reactions.Reactions.AssertExpectations(t)
}

func TestPostRemoveReaction_OwnPost(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.DELETE("/posts/:id/reactions", h.RemoveReaction)

	reactions.Authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(1), nil)

	w := doRequest(r, http.MethodDelete, "/posts/5/reactions", map[string]string{"emoji": "👍"})
	assertStatus(t, w, http.StatusForbidden)
	reactions.Reactions.AssertNotCalled(t, "RemoveReaction", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestPostRemoveReaction_InvalidEmoji(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.DELETE("/posts/:id/reactions", h.RemoveReaction)

	w := doRequest(r, http.MethodDelete, "/posts/5/reactions", map[string]string{"emoji": "🍣"})
	assertStatus(t, w, http.StatusBadRequest)
	reactions.Authors.AssertNotCalled(t, "FindAuthorID", mock.Anything, mock.Anything)
}

func TestPostRemoveReaction_InvalidID(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.DELETE("/posts/:id/reactions", h.RemoveReaction)

	w := doRequest(r, http.MethodDelete, "/posts/abc/reactions", map[string]string{"emoji": "👍"})
	assertStatus(t, w, http.StatusBadRequest)
	reactions.Authors.AssertNotCalled(t, "FindAuthorID", mock.Anything, mock.Anything)
}

func TestPostRemoveReaction_InvalidJSON(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.DELETE("/posts/:id/reactions", h.RemoveReaction)

	w := doRequestRaw(r, http.MethodDelete, "/posts/5/reactions", "{invalid json}")
	assertStatus(t, w, http.StatusBadRequest)
	reactions.Authors.AssertNotCalled(t, "FindAuthorID", mock.Anything, mock.Anything)
}

func TestPostRemoveReaction_RepositoryError(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.DELETE("/posts/:id/reactions", h.RemoveReaction)

	reactions.Authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(99), nil)
	reactions.Reactions.On("RemoveReaction", mock.Anything, uint(1), uint(5), "👍").Return(errors.New("db error"))

	w := doRequest(r, http.MethodDelete, "/posts/5/reactions", map[string]string{"emoji": "👍"})
	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- GetReactions ----------

func TestPostGetReactions_Success(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.GET("/posts/:id/reactions", h.GetReactions)

	reactions.Reactions.On("GetReactionsByPostID", mock.Anything, uint(5)).
		Return([]model.ReactionCount{{Emoji: "👍", Count: 3}}, nil)
	reactions.Reactions.On("GetUserReactions", mock.Anything, uint(1), uint(5)).Return([]string{"👍"}, nil)

	w := doRequest(r, http.MethodGet, "/posts/5/reactions", nil)
	assertStatus(t, w, http.StatusOK)
	body := w.Body.String()
	assert.Contains(t, body, `"emoji":"👍"`)
	assert.Contains(t, body, `"count":3`)
	reactions.Reactions.AssertExpectations(t)
}

// リアクションが無ければ空配列を返す（null にしない）。
func TestPostGetReactions_Empty(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.GET("/posts/:id/reactions", h.GetReactions)

	reactions.Reactions.On("GetReactionsByPostID", mock.Anything, uint(5)).Return(nil, nil)
	reactions.Reactions.On("GetUserReactions", mock.Anything, uint(1), uint(5)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/posts/5/reactions", nil)
	assertStatus(t, w, http.StatusOK)
	assert.JSONEq(t, `{"reactions":[],"user_reactions":[]}`, w.Body.String())
}

func TestPostGetReactions_InvalidID(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.GET("/posts/:id/reactions", h.GetReactions)

	w := doRequest(r, http.MethodGet, "/posts/abc/reactions", nil)
	assertStatus(t, w, http.StatusBadRequest)
	reactions.Reactions.AssertNotCalled(t, "GetReactionsByPostID", mock.Anything, mock.Anything)
}

func TestPostGetReactions_RepositoryError(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.GET("/posts/:id/reactions", h.GetReactions)

	reactions.Reactions.On("GetReactionsByPostID", mock.Anything, uint(5)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/posts/5/reactions", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	reactions.Reactions.AssertNotCalled(t, "GetUserReactions", mock.Anything, mock.Anything, mock.Anything)
}

func TestPostGetReactions_UserReactionsError(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.GET("/posts/:id/reactions", h.GetReactions)

	reactions.Reactions.On("GetReactionsByPostID", mock.Anything, uint(5)).Return([]model.ReactionCount{}, nil)
	reactions.Reactions.On("GetUserReactions", mock.Anything, uint(1), uint(5)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/posts/5/reactions", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- GetReactionsBatch ----------

func TestPostGetReactionsBatch_Success(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.POST("/posts/reactions/batch", h.GetReactionsBatch)

	postIDs := []uint{1, 2}
	reactions.Reactions.On("GetReactionsBatch", mock.Anything, postIDs).
		Return(map[uint][]model.ReactionCount{1: {{Emoji: "👍", Count: 2}}}, nil)
	reactions.Reactions.On("GetUserReactionsBatch", mock.Anything, uint(1), postIDs).
		Return(map[uint][]string{1: {"👍"}}, nil)

	w := doRequest(r, http.MethodPost, "/posts/reactions/batch", map[string]interface{}{"post_ids": postIDs})
	assertStatus(t, w, http.StatusOK)
	body := w.Body.String()
	assert.Contains(t, body, `"emoji":"👍"`)
	// リアクションが無い投稿にも空配列のエントリが入る。
	assert.Contains(t, body, `"2":{"reactions":[],"user_reactions":[]}`)
	reactions.Reactions.AssertExpectations(t)
}

// 51 件以上は 400。
func TestPostGetReactionsBatch_TooMany(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.POST("/posts/reactions/batch", h.GetReactionsBatch)

	postIDs := make([]uint, 51)
	for i := range postIDs {
		postIDs[i] = uint(i + 1)
	}

	w := doRequest(r, http.MethodPost, "/posts/reactions/batch", map[string]interface{}{"post_ids": postIDs})
	assertStatus(t, w, http.StatusBadRequest)
	reactions.Reactions.AssertNotCalled(t, "GetReactionsBatch", mock.Anything, mock.Anything)
}

func TestPostGetReactionsBatch_InvalidJSON(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.POST("/posts/reactions/batch", h.GetReactionsBatch)

	w := doRequestRaw(r, http.MethodPost, "/posts/reactions/batch", "{invalid json}")
	assertStatus(t, w, http.StatusBadRequest)
	reactions.Reactions.AssertNotCalled(t, "GetReactionsBatch", mock.Anything, mock.Anything)
}

func TestPostGetReactionsBatch_RepositoryError(t *testing.T) {
	h, reactions := setupPostHandlerWithReactionPorts()
	r := newRouter(1)
	r.POST("/posts/reactions/batch", h.GetReactionsBatch)

	postIDs := []uint{1}
	reactions.Reactions.On("GetReactionsBatch", mock.Anything, postIDs).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/posts/reactions/batch", map[string]interface{}{"post_ids": postIDs})
	assertStatus(t, w, http.StatusInternalServerError)
	reactions.Reactions.AssertNotCalled(t, "GetUserReactionsBatch", mock.Anything, mock.Anything, mock.Anything)
}

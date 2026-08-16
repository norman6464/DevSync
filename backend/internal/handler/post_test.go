package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// allowFollowerNotification はフォロワー通知（goroutine 実行）の呼び出しを許容する。
// 通知の内容は usecase 側のテストで検証する。
func allowFollowerNotification(ports *postHandlerPorts) {
	ports.Followers.On("FindFollowerIDs", mock.Anything, mock.Anything).Return(nil, nil).Maybe()
	ports.Followers.On("CreateBatch", mock.Anything, mock.Anything).Return(nil).Maybe()
}

// ---------- Create ----------

func TestPostCreate_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	allowFollowerNotification(ports)
	ports.Posts.On("Create", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
		return p.UserID == 1 && p.Title == "Test Post" && p.Content == "Hello"
	})).Return(nil)
	ports.Posts.On("FindByID", mock.Anything, mock.AnythingOfType("uint")).
		Return(&model.Post{Title: "Test Post", Content: "Hello"}, nil)

	w := doRequest(r, http.MethodPost, "/posts", map[string]string{
		"title": "Test Post", "content": "Hello",
	})
	assertStatus(t, w, http.StatusCreated)
	ports.Posts.AssertExpectations(t)
}

// 本文の @username からメンションを記録し、相手に通知する。
func TestPostCreate_ProcessesMentions(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	allowFollowerNotification(ports)
	ports.Posts.On("Create", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		args.Get(1).(*model.Post).ID = 42
	})
	ports.Posts.On("FindByID", mock.Anything, mock.AnythingOfType("uint")).
		Return(&model.Post{ID: 42}, nil).Maybe()
	ports.Usernames.On("FindByUsername", mock.Anything, "alice").Return(&model.User{ID: 2}, nil)
	ports.Mentions.On("FindByPostID", mock.Anything, uint(42)).Return([]model.Mention(nil), nil)
	ports.Mentions.On("Create", mock.Anything, mock.MatchedBy(func(m *model.Mention) bool {
		return m.UserID == 2 && m.ActorID == 1 && m.PostID != nil && *m.PostID == 42
	})).Return(true, nil)
	ports.Notifications.On("Create", mock.Anything, mock.MatchedBy(func(n *model.Notification) bool {
		return n.UserID == 2 && n.Type == model.NotificationTypeMention && n.PostID != nil && *n.PostID == 42
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts", map[string]string{
		"title": "Test Post", "content": "レビューお願いします @alice",
	})

	assertStatus(t, w, http.StatusCreated)
	ports.Mentions.AssertExpectations(t)
	ports.Notifications.AssertExpectations(t)
}

// メンションの記録に失敗しても投稿の作成は成功として返す。
func TestPostCreate_MentionFailureDoesNotFailPost(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	allowFollowerNotification(ports)
	ports.Posts.On("Create", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		args.Get(1).(*model.Post).ID = 42
	})
	ports.Posts.On("FindByID", mock.Anything, mock.AnythingOfType("uint")).
		Return(&model.Post{ID: 42}, nil).Maybe()
	ports.Usernames.On("FindByUsername", mock.Anything, "alice").Return(&model.User{ID: 2}, nil)
	ports.Mentions.On("FindByPostID", mock.Anything, uint(42)).Return([]model.Mention(nil), nil)
	ports.Mentions.On("Create", mock.Anything, mock.Anything).Return(false, errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/posts", map[string]string{
		"title": "Test Post", "content": "@alice",
	})

	assertStatus(t, w, http.StatusCreated)
}

// 作成後の再取得に失敗しても、作成した投稿をそのまま返す。
func TestPostCreate_RefetchFails(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	allowFollowerNotification(ports)
	ports.Posts.On("Create", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		args.Get(1).(*model.Post).ID = 7
	})
	ports.Posts.On("FindByID", mock.Anything, uint(7)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/posts", map[string]string{"title": "T", "content": "C"})
	assertStatus(t, w, http.StatusCreated)
	assert.Contains(t, w.Body.String(), `"id":7`)
}

// 読了時間は本文の文字数から算出して保存する。
func TestPostCreate_SetsEstimatedReadTime(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	content := ""
	for i := 0; i < 1500; i++ {
		content += "あ"
	}

	allowFollowerNotification(ports)
	ports.Posts.On("Create", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
		return p.EstimatedReadTime == 3
	})).Return(nil)
	ports.Posts.On("FindByID", mock.Anything, mock.AnythingOfType("uint")).Return(nil, nil)

	w := doRequest(r, http.MethodPost, "/posts", map[string]string{"title": "長文", "content": content})
	assertStatus(t, w, http.StatusCreated)
	ports.Posts.AssertExpectations(t)
}

func TestPostCreate_ValidationError(t *testing.T) {
	h, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	// title と content は required
	w := doRequest(r, http.MethodPost, "/posts", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCreate_InvalidJSON(t *testing.T) {
	h, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/posts", "{invalid json}")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCreate_RepositoryError(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	ports.Posts.On("Create", mock.Anything, mock.Anything).Return(domain.ErrBadRequest)

	w := doRequest(r, http.MethodPost, "/posts", map[string]string{"title": "Test", "content": "Content"})
	assertStatus(t, w, http.StatusBadRequest)
	ports.Posts.AssertExpectations(t)
}

// ---------- GetAll ----------

func TestPostGetAll_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts", h.GetAll)

	ports.Posts.On("FindAll", mock.Anything, 1, 20).Return([]model.Post{{Title: "A"}, {Title: "B"}}, nil)
	ports.Posts.On("CountAll", mock.Anything).Return(int64(2), nil)

	w := doRequest(r, http.MethodGet, "/posts", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"total":2`)
	ports.Posts.AssertExpectations(t)
}

func TestPostGetAll_WithPagination(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts", h.GetAll)

	ports.Posts.On("FindAll", mock.Anything, 2, 5).Return([]model.Post{}, nil)
	ports.Posts.On("CountAll", mock.Anything).Return(int64(10), nil)

	w := doRequest(r, http.MethodGet, "/posts?page=2&limit=5", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Posts.AssertExpectations(t)
}

func TestPostGetAll_RepositoryError(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts", h.GetAll)

	ports.Posts.On("FindAll", mock.Anything, 1, 20).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/posts", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Posts.AssertNotCalled(t, "CountAll", mock.Anything)
}

func TestPostGetAll_CountError(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts", h.GetAll)

	ports.Posts.On("FindAll", mock.Anything, 1, 20).Return([]model.Post{}, nil)
	ports.Posts.On("CountAll", mock.Anything).Return(int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/posts", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Posts.AssertExpectations(t)
}

// ---------- GetByID ----------

func TestPostGetByID_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id", h.GetByID)

	ports.Posts.On("FindByID", mock.Anything, uint(10)).Return(&model.Post{ID: 10, Title: "Found"}, nil)
	ports.Likes.On("HasLiked", mock.Anything, uint(1), uint(10)).Return(true, nil)
	ports.Bookmarks.On("HasBookmarked", mock.Anything, uint(1), uint(10)).Return(true, nil)

	w := doRequest(r, http.MethodGet, "/posts/10", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"liked":true`)
	assert.Contains(t, w.Body.String(), `"bookmarked":true`)
	ports.Posts.AssertExpectations(t)
	ports.Likes.AssertExpectations(t)
	ports.Bookmarks.AssertExpectations(t)
}

// いいね・ブックマークの判定に失敗しても投稿詳細は返す（移行前と同じく false になる）。
func TestPostGetByID_FlagQueriesFail(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id", h.GetByID)

	ports.Posts.On("FindByID", mock.Anything, uint(10)).Return(&model.Post{ID: 10}, nil)
	ports.Likes.On("HasLiked", mock.Anything, uint(1), uint(10)).Return(false, errors.New("db error"))
	ports.Bookmarks.On("HasBookmarked", mock.Anything, uint(1), uint(10)).Return(false, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/posts/10", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"liked":false`)
	assert.Contains(t, w.Body.String(), `"bookmarked":false`)
}

// 存在しない投稿は移行前と同じく内部エラーになる。
func TestPostGetByID_NotFound(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id", h.GetByID)

	ports.Posts.On("FindByID", mock.Anything, uint(999)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/posts/999", nil)
	assertStatus(t, w, http.StatusNotFound)
	ports.Posts.AssertExpectations(t)
}

func TestPostGetByID_InvalidID(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/posts/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
	ports.Posts.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// ---------- Update ----------

func TestPostUpdate_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id", h.Update)

	ports.Posts.On("FindByID", mock.Anything, uint(10)).
		Return(&model.Post{ID: 10, UserID: 1, Title: "Old", Content: "Old Content"}, nil)
	ports.Posts.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
		return p.Title == "New Title" && p.Content == "Old Content"
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/posts/10", map[string]string{"title": "New Title"})
	assertStatus(t, w, http.StatusOK)
	ports.Posts.AssertExpectations(t)
}

// 本文を更新すると読了時間も再計算される。
func TestPostUpdate_RecalculatesReadTime(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id", h.Update)

	content := ""
	for i := 0; i < 1000; i++ {
		content += "あ"
	}

	ports.Posts.On("FindByID", mock.Anything, uint(10)).
		Return(&model.Post{ID: 10, UserID: 1, EstimatedReadTime: 1}, nil)
	ports.Posts.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
		return p.EstimatedReadTime == 2
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/posts/10", map[string]string{"content": content})
	assertStatus(t, w, http.StatusOK)
	ports.Posts.AssertExpectations(t)
}

func TestPostUpdate_Forbidden(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1) // userID=1
	r.PUT("/posts/:id", h.Update)

	ports.Posts.On("FindByID", mock.Anything, uint(10)).
		Return(&model.Post{ID: 10, UserID: 999, Title: "Other's post"}, nil)

	w := doRequest(r, http.MethodPut, "/posts/10", map[string]string{"title": "Hacked"})
	assertStatus(t, w, http.StatusForbidden)
	ports.Posts.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// 存在しない投稿の更新は移行前と同じく内部エラーになる。
func TestPostUpdate_NotFound(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id", h.Update)

	ports.Posts.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)

	w := doRequest(r, http.MethodPut, "/posts/10", map[string]string{"title": "X"})
	assertStatus(t, w, http.StatusNotFound)
	ports.Posts.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestPostUpdate_WithImageUrls(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id", h.Update)

	imageURLs := `["https://example.com/image1.jpg"]`
	ports.Posts.On("FindByID", mock.Anything, uint(10)).
		Return(&model.Post{ID: 10, UserID: 1, Title: "Old", Content: "Old"}, nil)
	ports.Posts.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
		return p.ImageURLs == imageURLs
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/posts/10", map[string]interface{}{
		"title":      "Updated",
		"content":    "Updated content",
		"image_urls": imageURLs,
	})
	assertStatus(t, w, http.StatusOK)
	ports.Posts.AssertExpectations(t)
}

func TestPostUpdate_InvalidID(t *testing.T) {
	h, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/posts/abc", map[string]string{"title": "X"})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostUpdate_InvalidJSON(t *testing.T) {
	h, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id", h.Update)

	w := doRequestRaw(r, http.MethodPut, "/posts/10", "{invalid}")
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- Delete ----------

func TestPostDelete_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id", h.Delete)

	ports.Posts.On("FindByID", mock.Anything, uint(10)).Return(&model.Post{ID: 10, UserID: 1}, nil)
	ports.Posts.On("Delete", mock.Anything, uint(10)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/posts/10", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Posts.AssertExpectations(t)
}

func TestPostDelete_Forbidden(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id", h.Delete)

	ports.Posts.On("FindByID", mock.Anything, uint(10)).Return(&model.Post{ID: 10, UserID: 999}, nil)

	w := doRequest(r, http.MethodDelete, "/posts/10", nil)
	assertStatus(t, w, http.StatusForbidden)
	ports.Posts.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

func TestPostDelete_InvalidID(t *testing.T) {
	h, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/posts/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostDelete_NotFound(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id", h.Delete)

	ports.Posts.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)

	w := doRequest(r, http.MethodDelete, "/posts/10", nil)
	assertStatus(t, w, http.StatusNotFound)
	ports.Posts.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

// ---------- Timeline ----------

func TestPostTimeline_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/timeline", h.Timeline)

	ports.Posts.On("Timeline", mock.Anything, uint(1), 1, 20).
		Return([]model.Post{{Title: "Timeline Post"}}, nil)

	w := doRequest(r, http.MethodGet, "/posts/timeline", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Posts.AssertExpectations(t)
}

// タイムラインが空でも null ではなく空配列を返す。
func TestPostTimeline_Empty(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/timeline", h.Timeline)

	ports.Posts.On("Timeline", mock.Anything, uint(1), 1, 20).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/posts/timeline", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

func TestPostTimeline_RepositoryError(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/timeline", h.Timeline)

	ports.Posts.On("Timeline", mock.Anything, uint(1), 1, 20).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/posts/timeline", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Posts.AssertExpectations(t)
}

// ---------- Like / Unlike ----------

func TestPostLike_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/like", h.Like)

	ports.Authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(99), nil)
	ports.Likes.On("Like", mock.Anything, uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/like", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Likes.AssertExpectations(t)
	ports.Authors.AssertExpectations(t)
}

func TestPostLike_OwnPostForbidden(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/like", h.Like)

	ports.Authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(1), nil)

	w := doRequest(r, http.MethodPost, "/posts/5/like", nil)
	assertStatus(t, w, http.StatusForbidden)
	ports.Likes.AssertNotCalled(t, "Like", mock.Anything, mock.Anything, mock.Anything)
}

func TestPostLike_PostNotFound(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/like", h.Like)

	ports.Authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(0), nil)

	w := doRequest(r, http.MethodPost, "/posts/5/like", nil)
	assertStatus(t, w, http.StatusNotFound)
	ports.Likes.AssertNotCalled(t, "Like", mock.Anything, mock.Anything, mock.Anything)
}

func TestPostLike_InvalidID(t *testing.T) {
	h, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/like", h.Like)

	w := doRequest(r, http.MethodPost, "/posts/abc/like", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostUnlike_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/like", h.Unlike)

	ports.Authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(99), nil)
	ports.Likes.On("Unlike", mock.Anything, uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/posts/5/like", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Likes.AssertExpectations(t)
}

func TestPostUnlike_OwnPostForbidden(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/like", h.Unlike)

	ports.Authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(1), nil)

	w := doRequest(r, http.MethodDelete, "/posts/5/like", nil)
	assertStatus(t, w, http.StatusForbidden)
	ports.Likes.AssertNotCalled(t, "Unlike", mock.Anything, mock.Anything, mock.Anything)
}

func TestPostUnlike_InvalidID(t *testing.T) {
	h, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/like", h.Unlike)

	w := doRequest(r, http.MethodDelete, "/posts/abc/like", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostUnlike_RepositoryError(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/like", h.Unlike)

	ports.Authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(99), nil)
	ports.Likes.On("Unlike", mock.Anything, uint(1), uint(5)).Return(errors.New("db error"))

	w := doRequest(r, http.MethodDelete, "/posts/5/like", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Likes.AssertExpectations(t)
}

// ---------- Draft ----------

func TestPostCreate_Draft_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	ports.Posts.On("Create", mock.Anything, mock.Anything).Return(nil)
	ports.Posts.On("FindByID", mock.Anything, mock.AnythingOfType("uint")).
		Return(&model.Post{Title: "Draft Post", Content: "Draft content", IsDraft: true}, nil)

	w := doRequest(r, http.MethodPost, "/posts", map[string]interface{}{
		"title": "Draft Post", "content": "Draft content", "is_draft": true,
	})
	assertStatus(t, w, http.StatusCreated)
	// 下書きではフォロワー通知を行わない
	ports.Followers.AssertNotCalled(t, "FindFollowerIDs", mock.Anything, mock.Anything)
}

func TestPostGetDrafts_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/drafts", h.GetDrafts)

	ports.Posts.On("FindDraftsByUserID", mock.Anything, uint(1)).Return([]model.Post{
		{Title: "Draft 1", IsDraft: true},
		{Title: "Draft 2", IsDraft: true},
	}, nil)

	w := doRequest(r, http.MethodGet, "/posts/drafts", nil)
	assertStatus(t, w, http.StatusOK)

	var drafts []map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &drafts))
	assert.Len(t, drafts, 2)
}

func TestPostGetDrafts_RepositoryError(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/drafts", h.GetDrafts)

	ports.Posts.On("FindDraftsByUserID", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/posts/drafts", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Posts.AssertExpectations(t)
}

// ---------- Publish / Unpublish ----------

func TestPostPublish_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/publish", h.Publish)

	allowFollowerNotification(ports)
	ports.Posts.On("FindByID", mock.Anything, uint(5)).
		Return(&model.Post{ID: 5, UserID: 1, Title: "Draft", IsDraft: true}, nil)
	ports.Posts.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
		return !p.IsDraft
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/posts/5/publish", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Posts.AssertExpectations(t)
}

func TestPostPublish_Forbidden(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1) // userID=1
	r.PUT("/posts/:id/publish", h.Publish)

	ports.Posts.On("FindByID", mock.Anything, uint(5)).
		Return(&model.Post{ID: 5, UserID: 999, IsDraft: true}, nil)

	w := doRequest(r, http.MethodPut, "/posts/5/publish", nil)
	assertStatus(t, w, http.StatusForbidden)
	ports.Posts.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestPostPublish_NotDraft(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/publish", h.Publish)

	ports.Posts.On("FindByID", mock.Anything, uint(5)).
		Return(&model.Post{ID: 5, UserID: 1, IsDraft: false}, nil)

	w := doRequest(r, http.MethodPut, "/posts/5/publish", nil)
	assertStatus(t, w, http.StatusBadRequest)
	ports.Posts.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestPostPublish_InvalidID(t *testing.T) {
	h, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/publish", h.Publish)

	w := doRequest(r, http.MethodPut, "/posts/abc/publish", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostPublish_NotFound(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/publish", h.Publish)

	ports.Posts.On("FindByID", mock.Anything, uint(999)).Return(nil, nil)

	w := doRequest(r, http.MethodPut, "/posts/999/publish", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestPostUnpublish_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/unpublish", h.Unpublish)

	ports.Posts.On("FindByID", mock.Anything, uint(5)).
		Return(&model.Post{ID: 5, UserID: 1, IsDraft: false}, nil)
	ports.Posts.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
		return p.IsDraft
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/posts/5/unpublish", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Posts.AssertExpectations(t)
	// 下書きに戻すときはフォロワー通知を行わない
	ports.Followers.AssertNotCalled(t, "FindFollowerIDs", mock.Anything, mock.Anything)
}

func TestPostUnpublish_InvalidID(t *testing.T) {
	h, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/unpublish", h.Unpublish)

	w := doRequest(r, http.MethodPut, "/posts/abc/unpublish", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostUnpublish_NotFound(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/unpublish", h.Unpublish)

	ports.Posts.On("FindByID", mock.Anything, uint(999)).Return(nil, nil)

	w := doRequest(r, http.MethodPut, "/posts/999/unpublish", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestPostUnpublish_Forbidden(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/unpublish", h.Unpublish)

	ports.Posts.On("FindByID", mock.Anything, uint(5)).
		Return(&model.Post{ID: 5, UserID: 999, IsDraft: false}, nil)

	w := doRequest(r, http.MethodPut, "/posts/5/unpublish", nil)
	assertStatus(t, w, http.StatusForbidden)
}

func TestPostUnpublish_AlreadyDraft(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/unpublish", h.Unpublish)

	ports.Posts.On("FindByID", mock.Anything, uint(5)).
		Return(&model.Post{ID: 5, UserID: 1, IsDraft: true}, nil)

	w := doRequest(r, http.MethodPut, "/posts/5/unpublish", nil)
	assertStatus(t, w, http.StatusBadRequest)
	ports.Posts.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// ---------- GetUserPosts / GetMyPosts ----------

func TestPostGetUserPosts_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/users/:id/posts", h.GetUserPosts)

	ports.Posts.On("FindByUserID", mock.Anything, uint(2), 20, 0).
		Return([]model.Post{{Title: "User Post 1"}, {Title: "User Post 2"}}, int64(2), nil)

	w := doRequest(r, http.MethodGet, "/users/2/posts", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Posts.AssertExpectations(t)
}

func TestPostGetUserPosts_InvalidID(t *testing.T) {
	h, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/users/:id/posts", h.GetUserPosts)

	w := doRequest(r, http.MethodGet, "/users/abc/posts", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostGetUserPosts_RepositoryError(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/users/:id/posts", h.GetUserPosts)

	ports.Posts.On("FindByUserID", mock.Anything, uint(2), 20, 0).
		Return(nil, int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/users/2/posts", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Posts.AssertExpectations(t)
}

func TestPostGetMyPosts_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/my", h.GetMyPosts)

	ports.Posts.On("FindByUserID", mock.Anything, uint(1), 20, 0).
		Return([]model.Post{{Title: "My Post 1"}, {Title: "My Post 2"}}, int64(2), nil)

	w := doRequest(r, http.MethodGet, "/posts/my", nil)
	assertStatus(t, w, http.StatusOK)
	ports.Posts.AssertExpectations(t)
}

func TestPostGetMyPosts_RepositoryError(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/my", h.GetMyPosts)

	ports.Posts.On("FindByUserID", mock.Anything, uint(1), 20, 0).
		Return(nil, int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/posts/my", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Posts.AssertExpectations(t)
}

// ---------- 件数 ----------

func TestPostCounts_Success(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/my/count", h.GetMyCount)
	r.GET("/posts/drafts/count", h.GetDraftsCount)
	r.GET("/posts/scheduled/count", h.GetScheduledCount)

	ports.Posts.On("CountByUserID", mock.Anything, uint(1)).Return(int64(3), nil)
	ports.Posts.On("CountDraftsByUserID", mock.Anything, uint(1)).Return(int64(2), nil)
	ports.Posts.On("CountScheduledByUserID", mock.Anything, uint(1)).Return(int64(1), nil)

	w := doRequest(r, http.MethodGet, "/posts/my/count", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"count":3`)

	w = doRequest(r, http.MethodGet, "/posts/drafts/count", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"count":2`)

	w = doRequest(r, http.MethodGet, "/posts/scheduled/count", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"count":1`)

	ports.Posts.AssertExpectations(t)
}

func TestPostCounts_RepositoryError(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/my/count", h.GetMyCount)

	ports.Posts.On("CountByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/posts/my/count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Posts.AssertExpectations(t)
}

// ---------- Create with CodeSnippets ----------

func TestPostCreate_WithCodeSnippets(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	createdPost := &model.Post{ID: 10, Title: "Snippet Post", Content: "Hello #go"}
	allowFollowerNotification(ports)
	ports.Posts.On("Create", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		args.Get(1).(*model.Post).ID = 10
	})
	ports.Posts.On("FindByID", mock.Anything, mock.AnythingOfType("uint")).Return(createdPost, nil)
	// スニペット作成 usecase が投稿の存在確認 → 作成 → 再取得を行う
	ports.SnippetPosts.On("FindByID", mock.Anything, mock.AnythingOfType("uint")).Return(createdPost, nil)
	ports.Snippets.On("Create", mock.Anything, mock.AnythingOfType("*model.CodeSnippet")).Return(nil)
	ports.Snippets.On("FindByID", mock.Anything, mock.AnythingOfType("uint")).
		Return(&model.CodeSnippet{Language: "go", Code: "package main"}, nil)

	w := doRequest(r, http.MethodPost, "/posts", map[string]interface{}{
		"title":   "Snippet Post",
		"content": "Hello #go",
		"code_snippets": []map[string]string{
			{"language": "go", "file_name": "main.go", "code": "package main"},
		},
	})

	assertStatus(t, w, http.StatusCreated)
	ports.Snippets.AssertCalled(t, "Create", mock.Anything, mock.AnythingOfType("*model.CodeSnippet"))
}

func TestPostCreate_WithCodeSnippets_SkipsEmpty(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	createdPost := &model.Post{ID: 11, Title: "Post", Content: "Content"}
	allowFollowerNotification(ports)
	ports.Posts.On("Create", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		args.Get(1).(*model.Post).ID = 11
	})
	ports.Posts.On("FindByID", mock.Anything, uint(11)).Return(createdPost, nil)

	// language or code が空のスニペットはスキップされる
	w := doRequest(r, http.MethodPost, "/posts", map[string]interface{}{
		"title":   "Post",
		"content": "Content",
		"code_snippets": []map[string]string{
			{"language": "", "code": "some code"},
			{"language": "go", "code": ""},
		},
	})

	assertStatus(t, w, http.StatusCreated)
	ports.Snippets.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// ---------- Create with AutoTags ----------

// 自動タグ設定は本物の usecase を注入し、本文からの抽出と保存まで通す。
func TestPostCreate_WithAutoTags(t *testing.T) {
	h, ports := setupPostHandler()

	tagRepo := new(mockPostTagRepo)
	postReader := new(mockPostReader)
	h.SetAutoTagsUseCase(usecase.NewSetAutoPostTagsUseCase(
		usecase.NewSetPostTagsUseCase(tagRepo, postReader),
	))

	r := newRouter(1)
	r.POST("/posts", h.Create)

	createdPost := &model.Post{ID: 12, Title: "Tagged Post", Content: "Hello #golang"}
	allowFollowerNotification(ports)
	ports.Posts.On("Create", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		args.Get(1).(*model.Post).ID = 12
	})
	ports.Posts.On("FindByID", mock.Anything, uint(12)).Return(createdPost, nil)
	postReader.On("FindByID", mock.Anything, uint(12)).Return(ownedPost(12), nil)
	tagRepo.On("SetTags", mock.Anything, uint(12), []string{"golang"}).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts", map[string]string{
		"title": "Tagged Post", "content": "Hello #golang",
	})

	assertStatus(t, w, http.StatusCreated)
	// 本文の #golang が抽出され、正規化されたうえで保存される
	tagRepo.AssertCalled(t, "SetTags", mock.Anything, uint(12), []string{"golang"})
}

// ---------- AutoSaveDraft ----------

func TestPostAutoSaveDraft_CreateNew(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/drafts/auto-save", h.AutoSaveDraft)

	ports.Posts.On("Create", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
		return p.UserID == 1 && p.Title == "下書き" && p.IsDraft
	})).Run(func(args mock.Arguments) {
		args.Get(1).(*model.Post).ID = 100
	}).Return(nil)

	w := doRequest(r, http.MethodPut, "/posts/drafts/auto-save", map[string]interface{}{
		"title": "下書き", "content": "本文",
	})
	assertStatus(t, w, http.StatusOK)

	result := parseJSON(t, w)
	assert.Equal(t, float64(100), result["id"])
	assert.NotEmpty(t, result["updated_at"])
	ports.Posts.AssertExpectations(t)
}

func TestPostAutoSaveDraft_UpdateExisting(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/drafts/auto-save", h.AutoSaveDraft)

	existing := &model.Post{ID: 5, Title: "旧", Content: "旧本文", UserID: 1, IsDraft: true}
	ports.Posts.On("FindByID", mock.Anything, uint(5)).Return(existing, nil)
	ports.Posts.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
		return p.ID == 5 && p.Title == "新" && p.Content == "新本文"
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/posts/drafts/auto-save", map[string]interface{}{
		"id": 5, "title": "新", "content": "新本文",
	})
	assertStatus(t, w, http.StatusOK)

	result := parseJSON(t, w)
	assert.Equal(t, float64(5), result["id"])
	ports.Posts.AssertExpectations(t)
}

func TestPostAutoSaveDraft_Forbidden(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/drafts/auto-save", h.AutoSaveDraft)

	ports.Posts.On("FindByID", mock.Anything, uint(5)).
		Return(&model.Post{ID: 5, Title: "他人の下書き", UserID: 999, IsDraft: true}, nil)

	w := doRequest(r, http.MethodPut, "/posts/drafts/auto-save", map[string]interface{}{
		"id": 5, "title": "ハック",
	})
	assertStatus(t, w, http.StatusForbidden)
	ports.Posts.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestPostAutoSaveDraft_NotDraft(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/drafts/auto-save", h.AutoSaveDraft)

	ports.Posts.On("FindByID", mock.Anything, uint(5)).
		Return(&model.Post{ID: 5, Title: "公開済み", UserID: 1, IsDraft: false}, nil)

	w := doRequest(r, http.MethodPut, "/posts/drafts/auto-save", map[string]interface{}{
		"id": 5, "title": "更新",
	})
	assertStatus(t, w, http.StatusBadRequest)
	ports.Posts.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestPostAutoSaveDraft_InvalidJSON(t *testing.T) {
	h, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/drafts/auto-save", h.AutoSaveDraft)

	w := doRequestRaw(r, http.MethodPut, "/posts/drafts/auto-save", "{invalid}")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostAutoSaveDraft_RepositoryError(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/drafts/auto-save", h.AutoSaveDraft)

	ports.Posts.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPut, "/posts/drafts/auto-save", map[string]interface{}{
		"title": "下書き", "content": "本文",
	})
	assertStatus(t, w, http.StatusInternalServerError)
	ports.Posts.AssertExpectations(t)
}

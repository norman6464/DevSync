package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------- Create ----------

func TestPostCreate_Success(t *testing.T) {
	h, postRepo, notifRepo, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	postRepo.On("Create", mock.AnythingOfType("*model.Post")).Return(nil)
	notifRepo.On("GetFollowerIDs", uint(1)).Return([]uint{}, nil)
	postRepo.On("FindByID", mock.AnythingOfType("uint")).Return(&model.Post{
		Title: "Test Post", Content: "Hello",
	}, nil)

	w := doRequest(r, http.MethodPost, "/posts", map[string]string{
		"title": "Test Post", "content": "Hello",
	})

	assertStatus(t, w, http.StatusCreated)
}

func TestPostCreate_ValidationError(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	// title と content は required
	w := doRequest(r, http.MethodPost, "/posts", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCreate_InvalidJSON(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/posts", "{invalid json}")
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- GetAll ----------

func TestPostGetAll_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts", h.GetAll)

	postRepo.On("FindAll", 1, 20).Return([]model.Post{
		{Title: "A"}, {Title: "B"},
	}, nil)
	postRepo.On("CountAll").Return(int64(2), nil)

	w := doRequest(r, http.MethodGet, "/posts", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestPostGetAll_WithPagination(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts", h.GetAll)

	postRepo.On("FindAll", 2, 5).Return([]model.Post{}, nil)
	postRepo.On("CountAll").Return(int64(10), nil)

	w := doRequest(r, http.MethodGet, "/posts?page=2&limit=5", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- GetByID ----------

func TestPostGetByID_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id", h.GetByID)

	postRepo.On("FindByID", uint(10)).Return(&model.Post{Title: "Found"}, nil)
	postRepo.On("HasLiked", uint(1), uint(0)).Return(false)
	postRepo.On("HasBookmarked", uint(1), uint(0)).Return(false)

	w := doRequest(r, http.MethodGet, "/posts/10", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestPostGetByID_NotFound(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id", h.GetByID)

	postRepo.On("FindByID", uint(999)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/posts/999", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestPostGetByID_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/posts/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- Update ----------

func TestPostUpdate_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id", h.Update)

	post := &model.Post{Title: "Old", Content: "Old Content"}
	post.ID = 10
	post.UserID = 1
	postRepo.On("FindByID", uint(10)).Return(post, nil)
	postRepo.On("Update", mock.AnythingOfType("*model.Post")).Return(nil)

	w := doRequest(r, http.MethodPut, "/posts/10", map[string]string{
		"title": "New Title",
	})
	assertStatus(t, w, http.StatusOK)
}

func TestPostUpdate_Forbidden(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1) // userID=1
	r.PUT("/posts/:id", h.Update)

	post := &model.Post{Title: "Other's post"}
	post.ID = 10
	post.UserID = 999 // 別ユーザーの投稿
	postRepo.On("FindByID", uint(10)).Return(post, nil)

	w := doRequest(r, http.MethodPut, "/posts/10", map[string]string{
		"title": "Hacked",
	})
	assertStatus(t, w, http.StatusForbidden)
}

func TestPostUpdate_NotFound(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id", h.Update)

	postRepo.On("FindByID", uint(10)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodPut, "/posts/10", map[string]string{"title": "X"})
	assertStatus(t, w, http.StatusNotFound)
}

func TestPostUpdate_WithImageUrls(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id", h.Update)

	post := &model.Post{Title: "Old", Content: "Old", ImageURLs: ""}
	post.ID = 10
	post.UserID = 1
	postRepo.On("FindByID", uint(10)).Return(post, nil)
	postRepo.On("Update", mock.AnythingOfType("*model.Post")).Return(nil)

	imageUrls := `["https://example.com/image1.jpg"]`
	w := doRequest(r, http.MethodPut, "/posts/10", map[string]interface{}{
		"title":      "Updated",
		"content":    "Updated content",
		"image_urls": imageUrls,
	})
	assertStatus(t, w, http.StatusOK)
}

// ---------- Delete ----------

func TestPostDelete_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id", h.Delete)

	post := &model.Post{}
	post.ID = 10
	post.UserID = 1
	postRepo.On("FindByID", uint(10)).Return(post, nil)
	postRepo.On("Delete", uint(10)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/posts/10", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestPostDelete_Forbidden(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id", h.Delete)

	post := &model.Post{}
	post.ID = 10
	post.UserID = 999
	postRepo.On("FindByID", uint(10)).Return(post, nil)

	w := doRequest(r, http.MethodDelete, "/posts/10", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- Timeline ----------

func TestPostTimeline_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/timeline", h.Timeline)

	postRepo.On("Timeline", uint(1), 1, 20).Return([]model.Post{
		{Title: "Timeline Post"},
	}, nil)

	w := doRequest(r, http.MethodGet, "/posts/timeline", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- Like / Unlike ----------

func TestPostLike_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/like", h.Like)

	otherPost := &model.Post{UserID: 99}
	otherPost.ID = 5
	postRepo.On("FindByID", uint(5)).Return(otherPost, nil)
	postRepo.On("Like", uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/like", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestPostUnlike_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/like", h.Unlike)

	otherPost := &model.Post{UserID: 99}
	otherPost.ID = 5
	postRepo.On("FindByID", uint(5)).Return(otherPost, nil)
	postRepo.On("Unlike", uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/posts/5/like", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- Comments ----------

func TestPostGetComments_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id/comments", h.GetComments)

	postRepo.On("GetComments", uint(5)).Return([]model.Comment{
		{Content: "Nice!"},
	}, nil)

	w := doRequest(r, http.MethodGet, "/posts/5/comments", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestPostCreateComment_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	postRepo.On("CreateComment", mock.AnythingOfType("*model.Comment")).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]string{
		"content": "Great post!",
	})
	assertStatus(t, w, http.StatusCreated)
}

func TestPostCreateComment_ValidationError(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	// content は required
	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCreateComment_EmptyContent(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	// 空文字列は min=1 でエラー
	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]string{"content": ""})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCreateReply_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	parentComment := &model.Comment{PostID: 5, ParentID: nil}
	parentComment.ID = 10
	postRepo.On("FindCommentByID", uint(10)).Return(parentComment, nil)
	postRepo.On("CreateComment", mock.AnythingOfType("*model.Comment")).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]interface{}{
		"content":   "Reply to comment!",
		"parent_id": 10,
	})
	assertStatus(t, w, http.StatusCreated)
}

func TestPostDeleteComment_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/comments/:commentId", h.DeleteComment)

	comment := &model.Comment{UserID: 1}
	comment.ID = 3
	postRepo.On("FindCommentByID", uint(3)).Return(comment, nil)
	postRepo.On("DeleteComment", uint(3)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/posts/5/comments/3", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- Draft ----------

func TestPostCreate_Draft_Success(t *testing.T) {
	h, postRepo, notifRepo, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	postRepo.On("Create", mock.AnythingOfType("*model.Post")).Return(nil)
	// 下書きの場合、通知は送られない（GetFollowerIDsは呼ばれない）
	postRepo.On("FindByID", mock.AnythingOfType("uint")).Return(&model.Post{
		Title: "Draft Post", Content: "Draft content", IsDraft: true,
	}, nil)

	w := doRequest(r, http.MethodPost, "/posts", map[string]interface{}{
		"title": "Draft Post", "content": "Draft content", "is_draft": true,
	})

	assertStatus(t, w, http.StatusCreated)
	notifRepo.AssertNotCalled(t, "GetFollowerIDs")
}

func TestPostGetDrafts_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/drafts", h.GetDrafts)

	postRepo.On("FindDraftsByUserID", uint(1)).Return([]model.Post{
		{Title: "Draft 1", IsDraft: true},
		{Title: "Draft 2", IsDraft: true},
	}, nil)

	w := doRequest(r, http.MethodGet, "/posts/drafts", nil)
	assertStatus(t, w, http.StatusOK)

	var drafts []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &drafts)
	assert.Len(t, drafts, 2)
}

func TestPostPublish_Success(t *testing.T) {
	h, postRepo, notifRepo, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/publish", h.Publish)

	postRepo.On("FindByID", uint(5)).Return(&model.Post{
		ID: 5, UserID: 1, Title: "Draft", IsDraft: true,
	}, nil)
	postRepo.On("Update", mock.AnythingOfType("*model.Post")).Return(nil)
	notifRepo.On("GetFollowerIDs", uint(1)).Return([]uint{2, 3}, nil)
	notifRepo.On("CreateBatch", mock.AnythingOfType("[]*model.Notification")).Return(nil)

	w := doRequest(r, http.MethodPut, "/posts/5/publish", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestPostPublish_Forbidden(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1) // userID=1
	r.PUT("/posts/:id/publish", h.Publish)

	// 別のユーザーの下書き
	postRepo.On("FindByID", uint(5)).Return(&model.Post{
		ID: 5, UserID: 999, Title: "Draft", IsDraft: true,
	}, nil)

	w := doRequest(r, http.MethodPut, "/posts/5/publish", nil)
	assertStatus(t, w, http.StatusForbidden)
}

func TestPostPublish_NotDraft(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/publish", h.Publish)

	// すでに公開済み
	postRepo.On("FindByID", uint(5)).Return(&model.Post{
		ID: 5, UserID: 1, Title: "Published", IsDraft: false,
	}, nil)

	w := doRequest(r, http.MethodPut, "/posts/5/publish", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostPublish_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/publish", h.Publish)

	w := doRequest(r, http.MethodPut, "/posts/abc/publish", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostPublish_NotFound(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/publish", h.Publish)

	postRepo.On("FindByID", uint(999)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodPut, "/posts/999/publish", nil)
	assertStatus(t, w, http.StatusNotFound)
}

// ---------- Unpublish ----------

func TestPostUnpublish_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/unpublish", h.Unpublish)

	postRepo.On("FindByID", uint(5)).Return(&model.Post{
		ID: 5, UserID: 1, Title: "Published", IsDraft: false,
	}, nil)
	postRepo.On("Update", mock.AnythingOfType("*model.Post")).Return(nil)

	w := doRequest(r, http.MethodPut, "/posts/5/unpublish", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestPostUnpublish_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/unpublish", h.Unpublish)

	w := doRequest(r, http.MethodPut, "/posts/abc/unpublish", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostUnpublish_NotFound(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/unpublish", h.Unpublish)

	postRepo.On("FindByID", uint(999)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodPut, "/posts/999/unpublish", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestPostUnpublish_Forbidden(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/unpublish", h.Unpublish)

	postRepo.On("FindByID", uint(5)).Return(&model.Post{
		ID: 5, UserID: 999, Title: "Other's", IsDraft: false,
	}, nil)

	w := doRequest(r, http.MethodPut, "/posts/5/unpublish", nil)
	assertStatus(t, w, http.StatusForbidden)
}

func TestPostUnpublish_AlreadyDraft(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id/unpublish", h.Unpublish)

	postRepo.On("FindByID", uint(5)).Return(&model.Post{
		ID: 5, UserID: 1, Title: "Draft", IsDraft: true,
	}, nil)

	w := doRequest(r, http.MethodPut, "/posts/5/unpublish", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- GetUserPosts ----------

func TestPostGetUserPosts_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/users/:id/posts", h.GetUserPosts)

	postRepo.On("FindByUserID", uint(2), 20, 0).Return([]model.Post{
		{Title: "User Post 1"}, {Title: "User Post 2"},
	}, int64(2), nil)

	w := doRequest(r, http.MethodGet, "/users/2/posts", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestPostGetUserPosts_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/users/:id/posts", h.GetUserPosts)

	w := doRequest(r, http.MethodGet, "/users/abc/posts", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostGetUserPosts_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/users/:id/posts", h.GetUserPosts)

	postRepo.On("FindByUserID", uint(2), 20, 0).Return([]model.Post(nil), int64(0), service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/users/2/posts", nil)
	assertStatus(t, w, http.StatusNotFound)
}

// ---------- GetReplies ----------

func TestPostGetReplies_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/comments/:commentId/replies", h.GetReplies)

	postRepo.On("GetReplies", uint(10)).Return([]model.Comment{
		{Content: "Reply 1"}, {Content: "Reply 2"},
	}, nil)

	w := doRequest(r, http.MethodGet, "/comments/10/replies", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestPostGetReplies_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/comments/:commentId/replies", h.GetReplies)

	w := doRequest(r, http.MethodGet, "/comments/abc/replies", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostGetReplies_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/comments/:commentId/replies", h.GetReplies)

	postRepo.On("GetReplies", uint(10)).Return([]model.Comment(nil), service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/comments/10/replies", nil)
	assertStatus(t, w, http.StatusNotFound)
}

// ---------- Bookmark ----------

func TestPostBookmark_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/bookmark", h.Bookmark)

	otherPost := &model.Post{UserID: 99}
	otherPost.ID = 5
	postRepo.On("FindByID", uint(5)).Return(otherPost, nil)
	postRepo.On("Bookmark", uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/bookmark", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestPostBookmark_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/bookmark", h.Bookmark)

	w := doRequest(r, http.MethodPost, "/posts/abc/bookmark", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostBookmark_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/bookmark", h.Bookmark)

	postRepo.On("FindByID", uint(5)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodPost, "/posts/5/bookmark", nil)
	assertStatus(t, w, http.StatusNotFound)
}

// ---------- Unbookmark ----------

func TestPostUnbookmark_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/bookmark", h.Unbookmark)

	otherPost := &model.Post{UserID: 99}
	otherPost.ID = 5
	postRepo.On("FindByID", uint(5)).Return(otherPost, nil)
	postRepo.On("Unbookmark", uint(1), uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/posts/5/bookmark", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestPostUnbookmark_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/bookmark", h.Unbookmark)

	w := doRequest(r, http.MethodDelete, "/posts/abc/bookmark", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostUnbookmark_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/bookmark", h.Unbookmark)

	postRepo.On("FindByID", uint(5)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodDelete, "/posts/5/bookmark", nil)
	assertStatus(t, w, http.StatusNotFound)
}

// ---------- GetBookmarks ----------

func TestPostGetBookmarks_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/bookmarks", h.GetBookmarks)

	postRepo.On("FindBookmarkedByUserID", uint(1), 1, 20).Return([]model.Post{
		{Title: "Bookmarked 1"},
	}, int64(1), nil)

	w := doRequest(r, http.MethodGet, "/posts/bookmarks", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestPostGetBookmarks_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/bookmarks", h.GetBookmarks)

	postRepo.On("FindBookmarkedByUserID", uint(1), 1, 20).Return([]model.Post(nil), int64(0), service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/posts/bookmarks", nil)
	assertStatus(t, w, http.StatusNotFound)
}

// リアクションは DIP へ移行済み。テストは post_reaction_test.go に置く。

func TestPostCreate_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	postRepo.On("Create", mock.AnythingOfType("*model.Post")).Return(service.ErrBadRequest)

	w := doRequest(r, http.MethodPost, "/posts", map[string]string{
		"title": "Test", "content": "Content",
	})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostGetAll_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts", h.GetAll)

	postRepo.On("FindAll", 1, 20).Return([]model.Post(nil), service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/posts", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestPostGetAll_CountError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts", h.GetAll)

	postRepo.On("FindAll", 1, 20).Return([]model.Post{}, nil)
	postRepo.On("CountAll").Return(int64(0), service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/posts", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestPostTimeline_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/timeline", h.Timeline)

	postRepo.On("Timeline", uint(1), 1, 20).Return([]model.Post(nil), service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/posts/timeline", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestPostLike_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/like", h.Like)

	w := doRequest(r, http.MethodPost, "/posts/abc/like", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostLike_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/like", h.Like)

	postRepo.On("FindByID", uint(5)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodPost, "/posts/5/like", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestPostUnlike_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/like", h.Unlike)

	w := doRequest(r, http.MethodDelete, "/posts/abc/like", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostUnlike_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/like", h.Unlike)

	postRepo.On("FindByID", uint(5)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodDelete, "/posts/5/like", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestPostGetComments_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id/comments", h.GetComments)

	w := doRequest(r, http.MethodGet, "/posts/abc/comments", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostGetComments_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id/comments", h.GetComments)

	postRepo.On("GetComments", uint(5)).Return([]model.Comment(nil), service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/posts/5/comments", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestPostCreateComment_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	w := doRequest(r, http.MethodPost, "/posts/abc/comments", map[string]string{"content": "test"})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCreateComment_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	postRepo.On("CreateComment", mock.AnythingOfType("*model.Comment")).Return(service.ErrBadRequest)

	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]string{"content": "test"})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostDeleteComment_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/comments/:commentId", h.DeleteComment)

	w := doRequest(r, http.MethodDelete, "/posts/5/comments/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostDeleteComment_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/comments/:commentId", h.DeleteComment)

	postRepo.On("FindCommentByID", uint(3)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodDelete, "/posts/5/comments/3", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestPostGetDrafts_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/drafts", h.GetDrafts)

	postRepo.On("FindDraftsByUserID", uint(1)).Return([]model.Post(nil), service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/posts/drafts", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestPostUpdate_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id", h.Update)

	w := doRequest(r, http.MethodPut, "/posts/abc", map[string]string{"title": "X"})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostUpdate_InvalidJSON(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/:id", h.Update)

	w := doRequestRaw(r, http.MethodPut, "/posts/10", "{invalid}")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostDelete_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id", h.Delete)

	w := doRequest(r, http.MethodDelete, "/posts/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostDelete_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id", h.Delete)

	postRepo.On("FindByID", uint(10)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodDelete, "/posts/10", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestPostCreateComment_InvalidJSON(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	w := doRequestRaw(r, http.MethodPost, "/posts/5/comments", "{invalid}")
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- Create with CodeSnippets ----------

func TestPostCreate_WithCodeSnippets(t *testing.T) {
	h, postRepo, notifRepo, snippetPorts := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	createdPost := &model.Post{Title: "Snippet Post", Content: "Hello #go"}
	createdPost.ID = 10
	postRepo.On("Create", mock.AnythingOfType("*model.Post")).Return(nil).Run(func(args mock.Arguments) {
		p := args.Get(0).(*model.Post)
		p.ID = 10
	})
	notifRepo.On("GetFollowerIDs", uint(1)).Return([]uint{}, nil)
	// スニペット作成 usecase が投稿の存在確認 → 作成 → 再取得を行う
	postRepo.On("FindByID", mock.AnythingOfType("uint")).Return(createdPost, nil)
	snippetPorts.Posts.On("FindByID", mock.Anything, mock.AnythingOfType("uint")).Return(createdPost, nil)
	snippetPorts.Snippets.On("Create", mock.Anything, mock.AnythingOfType("*model.CodeSnippet")).Return(nil)
	snippetPorts.Snippets.On("FindByID", mock.Anything, mock.AnythingOfType("uint")).Return(&model.CodeSnippet{
		Language: "go", Code: "package main",
	}, nil)

	w := doRequest(r, http.MethodPost, "/posts", map[string]interface{}{
		"title":   "Snippet Post",
		"content": "Hello #go",
		"code_snippets": []map[string]string{
			{"language": "go", "file_name": "main.go", "code": "package main"},
		},
	})

	assertStatus(t, w, http.StatusCreated)
	snippetPorts.Snippets.AssertCalled(t, "Create", mock.Anything, mock.AnythingOfType("*model.CodeSnippet"))
}

func TestPostCreate_WithCodeSnippets_SkipsEmpty(t *testing.T) {
	h, postRepo, notifRepo, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	createdPost := &model.Post{Title: "Post", Content: "Content"}
	createdPost.ID = 11
	postRepo.On("Create", mock.AnythingOfType("*model.Post")).Return(nil).Run(func(args mock.Arguments) {
		p := args.Get(0).(*model.Post)
		p.ID = 11
	})
	notifRepo.On("GetFollowerIDs", uint(1)).Return([]uint{}, nil)
	postRepo.On("FindByID", uint(11)).Return(createdPost, nil)

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
}

// ---------- Create with AutoTags ----------

// 自動タグ設定は本物の usecase を注入し、本文からの抽出と保存まで通す。
func TestPostCreate_WithAutoTags(t *testing.T) {
	h, postRepo, notifRepo, _ := setupPostHandler()

	tagRepo := new(mockPostTagRepo)
	postReader := new(mockPostReader)
	h.SetAutoTagsUseCase(usecase.NewSetAutoPostTagsUseCase(
		usecase.NewSetPostTagsUseCase(tagRepo, postReader),
	))

	r := newRouter(1)
	r.POST("/posts", h.Create)

	createdPost := &model.Post{Title: "Tagged Post", Content: "Hello #golang"}
	createdPost.ID = 12
	postRepo.On("Create", mock.AnythingOfType("*model.Post")).Return(nil).Run(func(args mock.Arguments) {
		p := args.Get(0).(*model.Post)
		p.ID = 12
	})
	notifRepo.On("GetFollowerIDs", uint(1)).Return([]uint{}, nil)
	postRepo.On("FindByID", uint(12)).Return(createdPost, nil)
	postReader.On("FindByID", mock.Anything, uint(12)).Return(ownedPost(12), nil)
	tagRepo.On("SetTags", mock.Anything, uint(12), []string{"golang"}).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts", map[string]string{
		"title": "Tagged Post", "content": "Hello #golang",
	})

	assertStatus(t, w, http.StatusCreated)
	// 本文の #golang が抽出され、正規化されたうえで保存される
	tagRepo.AssertCalled(t, "SetTags", mock.Anything, uint(12), []string{"golang"})
}

// ---------- HideComment / UnhideComment ----------

func TestPostHideComment_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/comments/:commentId/hide", h.HideComment)

	comment := &model.Comment{PostID: 5, Content: "test"}
	comment.ID = 10
	comment.UserID = 1
	postRepo.On("FindCommentByID", uint(10)).Return(comment, nil)
	postRepo.On("UpdateComment", mock.AnythingOfType("*model.Comment")).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/comments/10/hide", nil)

	assertStatus(t, w, http.StatusOK)
	postRepo.AssertExpectations(t)
}

func TestPostHideComment_Forbidden(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/comments/:commentId/hide", h.HideComment)

	comment := &model.Comment{PostID: 5, Content: "test"}
	comment.ID = 10
	comment.UserID = 999 // 別ユーザー
	postRepo.On("FindCommentByID", uint(10)).Return(comment, nil)

	w := doRequest(r, http.MethodPost, "/posts/5/comments/10/hide", nil)

	assertStatus(t, w, http.StatusForbidden)
}

func TestPostHideComment_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/comments/:commentId/hide", h.HideComment)

	w := doRequest(r, http.MethodPost, "/posts/5/comments/abc/hide", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostUnhideComment_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/comments/:commentId/unhide", h.UnhideComment)

	comment := &model.Comment{PostID: 5, Content: "test", IsHidden: true}
	comment.ID = 10
	comment.UserID = 1
	postRepo.On("FindCommentByID", uint(10)).Return(comment, nil)
	postRepo.On("UpdateComment", mock.AnythingOfType("*model.Comment")).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/comments/10/unhide", nil)

	assertStatus(t, w, http.StatusOK)
	postRepo.AssertExpectations(t)
}

func TestPostUnhideComment_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/comments/:commentId/unhide", h.UnhideComment)

	w := doRequest(r, http.MethodPost, "/posts/5/comments/abc/unhide", nil)

	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- AutoSaveDraft ----------

func TestPostAutoSaveDraft_CreateNew(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/drafts/auto-save", h.AutoSaveDraft)

	postRepo.On("Create", mock.MatchedBy(func(p *model.Post) bool {
		return p.UserID == 1 && p.Title == "下書き" && p.IsDraft
	})).Run(func(args mock.Arguments) {
		p := args.Get(0).(*model.Post)
		p.ID = 100
	}).Return(nil)

	w := doRequest(r, http.MethodPut, "/posts/drafts/auto-save", map[string]interface{}{
		"title": "下書き", "content": "本文",
	})
	assertStatus(t, w, http.StatusOK)

	result := parseJSON(t, w)
	assert.Equal(t, float64(100), result["id"])
	assert.NotEmpty(t, result["updated_at"])
}

func TestPostAutoSaveDraft_UpdateExisting(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/drafts/auto-save", h.AutoSaveDraft)

	existing := &model.Post{Title: "旧", Content: "旧本文", UserID: 1, IsDraft: true}
	existing.ID = 5

	postRepo.On("FindByID", uint(5)).Return(existing, nil)
	postRepo.On("Update", existing).Return(nil)

	w := doRequest(r, http.MethodPut, "/posts/drafts/auto-save", map[string]interface{}{
		"id": 5, "title": "新", "content": "新本文",
	})
	assertStatus(t, w, http.StatusOK)

	result := parseJSON(t, w)
	assert.Equal(t, float64(5), result["id"])
}

func TestPostAutoSaveDraft_Forbidden(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/drafts/auto-save", h.AutoSaveDraft)

	existing := &model.Post{Title: "他人の下書き", UserID: 999, IsDraft: true}
	existing.ID = 5

	postRepo.On("FindByID", uint(5)).Return(existing, nil)

	w := doRequest(r, http.MethodPut, "/posts/drafts/auto-save", map[string]interface{}{
		"id": 5, "title": "ハック",
	})
	assertStatus(t, w, http.StatusForbidden)
}

func TestPostAutoSaveDraft_NotDraft(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/drafts/auto-save", h.AutoSaveDraft)

	published := &model.Post{Title: "公開済み", UserID: 1, IsDraft: false}
	published.ID = 5

	postRepo.On("FindByID", uint(5)).Return(published, nil)

	w := doRequest(r, http.MethodPut, "/posts/drafts/auto-save", map[string]interface{}{
		"id": 5, "title": "更新",
	})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostAutoSaveDraft_InvalidJSON(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/drafts/auto-save", h.AutoSaveDraft)

	w := doRequestRaw(r, http.MethodPut, "/posts/drafts/auto-save", "{invalid}")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostAutoSaveDraft_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.PUT("/posts/drafts/auto-save", h.AutoSaveDraft)

	postRepo.On("Create", mock.AnythingOfType("*model.Post")).Return(service.ErrNotFound)

	w := doRequest(r, http.MethodPut, "/posts/drafts/auto-save", map[string]interface{}{
		"title": "下書き", "content": "本文",
	})
	assertStatus(t, w, http.StatusNotFound)
}

func TestPostGetMyPosts_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/my", h.GetMyPosts)

	postRepo.On("FindByUserID", uint(1), 20, 0).Return([]model.Post{
		{Title: "My Post 1"}, {Title: "My Post 2"},
	}, int64(2), nil)

	w := doRequest(r, http.MethodGet, "/posts/my", nil)
	assertStatus(t, w, http.StatusOK)
	postRepo.AssertExpectations(t)
}

func TestPostGetMyPosts_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/my", h.GetMyPosts)

	postRepo.On("FindByUserID", uint(1), 20, 0).Return([]model.Post(nil), int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/posts/my", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	postRepo.AssertExpectations(t)
}

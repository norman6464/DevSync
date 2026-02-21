package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
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

// ---------- AddReaction ----------

func TestPostAddReaction_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/reactions", h.AddReaction)

	otherPost := &model.Post{UserID: 99}
	otherPost.ID = 5
	postRepo.On("FindByID", uint(5)).Return(otherPost, nil)
	postRepo.On("AddReaction", uint(1), uint(5), "👍").Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/reactions", map[string]string{"emoji": "👍"})
	assertStatus(t, w, http.StatusOK)
}

func TestPostAddReaction_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/reactions", h.AddReaction)

	w := doRequest(r, http.MethodPost, "/posts/abc/reactions", map[string]string{"emoji": "👍"})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostAddReaction_InvalidJSON(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/reactions", h.AddReaction)

	w := doRequestRaw(r, http.MethodPost, "/posts/5/reactions", "{invalid}")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostAddReaction_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/reactions", h.AddReaction)

	postRepo.On("FindByID", uint(5)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodPost, "/posts/5/reactions", map[string]string{"emoji": "👍"})
	assertStatus(t, w, http.StatusNotFound)
}

// ---------- RemoveReaction ----------

func TestPostRemoveReaction_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/reactions", h.RemoveReaction)

	otherPost := &model.Post{UserID: 99}
	otherPost.ID = 5
	postRepo.On("FindByID", uint(5)).Return(otherPost, nil)
	postRepo.On("RemoveReaction", uint(1), uint(5), "🔥").Return(nil)

	w := doRequest(r, http.MethodDelete, "/posts/5/reactions", map[string]string{"emoji": "🔥"})
	assertStatus(t, w, http.StatusOK)
}

func TestPostRemoveReaction_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/reactions", h.RemoveReaction)

	w := doRequest(r, http.MethodDelete, "/posts/abc/reactions", map[string]string{"emoji": "🔥"})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostRemoveReaction_InvalidJSON(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/reactions", h.RemoveReaction)

	w := doRequestRaw(r, http.MethodDelete, "/posts/5/reactions", "{invalid}")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostRemoveReaction_ServiceError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.DELETE("/posts/:id/reactions", h.RemoveReaction)

	postRepo.On("FindByID", uint(5)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodDelete, "/posts/5/reactions", map[string]string{"emoji": "🔥"})
	assertStatus(t, w, http.StatusNotFound)
}

// ---------- GetReactions ----------

func TestPostGetReactions_Success(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id/reactions", h.GetReactions)

	postRepo.On("GetReactionsByPostID", uint(5)).Return([]model.ReactionCount{
		{Emoji: "👍", Count: 3},
	}, nil)
	postRepo.On("GetUserReactions", uint(1), uint(5)).Return([]string{"👍"}, nil)

	w := doRequest(r, http.MethodGet, "/posts/5/reactions", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestPostGetReactions_InvalidID(t *testing.T) {
	h, _, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id/reactions", h.GetReactions)

	w := doRequest(r, http.MethodGet, "/posts/abc/reactions", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostGetReactions_GetReactionsError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id/reactions", h.GetReactions)

	postRepo.On("GetReactionsByPostID", uint(5)).Return([]model.ReactionCount(nil), service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/posts/5/reactions", nil)
	assertStatus(t, w, http.StatusNotFound)
}

func TestPostGetReactions_GetUserReactionsError(t *testing.T) {
	h, postRepo, _, _ := setupPostHandler()
	r := newRouter(1)
	r.GET("/posts/:id/reactions", h.GetReactions)

	postRepo.On("GetReactionsByPostID", uint(5)).Return([]model.ReactionCount{}, nil)
	postRepo.On("GetUserReactions", uint(1), uint(5)).Return([]string(nil), service.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/posts/5/reactions", nil)
	assertStatus(t, w, http.StatusNotFound)
}

// ---------- エラーパス追加テスト ----------

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
	h, postRepo, notifRepo, snippetRepo := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts", h.Create)

	createdPost := &model.Post{Title: "Snippet Post", Content: "Hello #go"}
	createdPost.ID = 10
	postRepo.On("Create", mock.AnythingOfType("*model.Post")).Return(nil).Run(func(args mock.Arguments) {
		p := args.Get(0).(*model.Post)
		p.ID = 10
	})
	notifRepo.On("GetFollowerIDs", uint(1)).Return([]uint{}, nil)
	// CodeSnippetService.Create 内部で postRepo.FindByID, snippetRepo.Create, snippetRepo.FindByID を呼ぶ
	postRepo.On("FindByID", mock.AnythingOfType("uint")).Return(createdPost, nil)
	snippetRepo.On("Create", mock.AnythingOfType("*model.CodeSnippet")).Return(nil)
	snippetRepo.On("FindByID", mock.AnythingOfType("uint")).Return(&model.CodeSnippet{
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
	snippetRepo.AssertCalled(t, "Create", mock.AnythingOfType("*model.CodeSnippet"))
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

// ---------- Create with TagService ----------

func TestPostCreate_WithTagService(t *testing.T) {
	h, postRepo, notifRepo, _ := setupPostHandler()

	tagSvc := new(MockPostTagService)
	h.SetTagService(tagSvc)

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
	tagSvc.On("SetAutoTags", uint(12), uint(1), "Hello #golang").Return(nil)

	w := doRequest(r, http.MethodPost, "/posts", map[string]string{
		"title": "Tagged Post", "content": "Hello #golang",
	})

	assertStatus(t, w, http.StatusCreated)
	tagSvc.AssertCalled(t, "SetAutoTags", uint(12), uint(1), "Hello #golang")
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

package handler

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// ---------- GetComments ----------

func TestPostGetComments_Success(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.GET("/posts/:id/comments", h.GetComments)

	comments.On("ListByPostID", mock.Anything, uint(5)).Return([]model.Comment{{ID: 1, Content: "Nice!"}}, nil)

	w := doRequest(r, http.MethodGet, "/posts/5/comments", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"content":"Nice!"`)
	comments.AssertExpectations(t)
}

// コメントが無ければ null ではなく空配列を返す。
func TestPostGetComments_Empty(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.GET("/posts/:id/comments", h.GetComments)

	comments.On("ListByPostID", mock.Anything, uint(5)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/posts/5/comments", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

func TestPostGetComments_InvalidID(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.GET("/posts/:id/comments", h.GetComments)

	w := doRequest(r, http.MethodGet, "/posts/abc/comments", nil)
	assertStatus(t, w, http.StatusBadRequest)
	comments.AssertNotCalled(t, "ListByPostID", mock.Anything, mock.Anything)
}

func TestPostGetComments_RepositoryError(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.GET("/posts/:id/comments", h.GetComments)

	comments.On("ListByPostID", mock.Anything, uint(5)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/posts/5/comments", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	comments.AssertExpectations(t)
}

// ---------- CreateComment ----------

func TestPostCreateComment_Success(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	comments.On("Create", mock.Anything, mock.MatchedBy(func(c *model.Comment) bool {
		return c.UserID == 1 && c.PostID == 5 && c.Content == "Great post!" && c.ParentID == nil
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]string{"content": "Great post!"})
	assertStatus(t, w, http.StatusCreated)
	comments.AssertExpectations(t)
}

// コメント本文の @username からメンションを記録し、通知から元の投稿へ辿れるようにする。
func TestPostCreateComment_ProcessesMentions(t *testing.T) {
	h, ports := setupPostHandler()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	ports.Comments.On("Create", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		args.Get(1).(*model.Comment).ID = 30
	})
	ports.Usernames.On("FindByUsername", mock.Anything, "alice").Return(&model.User{ID: 2}, nil)
	ports.Mentions.On("FindByCommentID", mock.Anything, uint(30)).Return([]model.Mention(nil), nil)
	ports.Mentions.On("Create", mock.Anything, mock.MatchedBy(func(m *model.Mention) bool {
		// 記録はコメントに紐づける（投稿本文のメンションと混ざらないようにする）
		return m.UserID == 2 && m.CommentID != nil && *m.CommentID == 30 && m.PostID == nil
	})).Return(nil)
	ports.Notifications.On("Create", mock.Anything, mock.MatchedBy(func(n *model.Notification) bool {
		// 通知からは元の投稿へ辿れる
		return n.UserID == 2 && n.Type == model.NotificationTypeMention && n.PostID != nil && *n.PostID == 5
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]string{"content": "@alice これどう思う"})

	assertStatus(t, w, http.StatusCreated)
	ports.Mentions.AssertExpectations(t)
	ports.Notifications.AssertExpectations(t)
}

// 前後の空白は保存前に取り除かれる。
func TestPostCreateComment_TrimsContent(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	comments.On("Create", mock.Anything, mock.MatchedBy(func(c *model.Comment) bool {
		return c.Content == "hello"
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]string{"content": "  hello  "})
	assertStatus(t, w, http.StatusCreated)
	comments.AssertExpectations(t)
}

// 空白だけの本文は DTO を通過するが、usecase のバリデーションで弾かれる。
func TestPostCreateComment_BlankContent(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]string{"content": "   "})
	assertStatus(t, w, http.StatusBadRequest)
	comments.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// 1000 文字を超える本文は usecase のバリデーションで弾かれる。
func TestPostCreateComment_TooLongContent(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]string{"content": strings.Repeat("a", 1001)})
	assertStatus(t, w, http.StatusBadRequest)
	comments.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestPostCreateComment_ValidationError(t *testing.T) {
	h, _ := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	// content は required
	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCreateComment_EmptyContent(t *testing.T) {
	h, _ := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	// 空文字列は min=1 でエラー
	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]string{"content": ""})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCreateComment_InvalidID(t *testing.T) {
	h, _ := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	w := doRequest(r, http.MethodPost, "/posts/abc/comments", map[string]string{"content": "test"})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCreateComment_InvalidJSON(t *testing.T) {
	h, _ := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	w := doRequestRaw(r, http.MethodPost, "/posts/5/comments", "{invalid}")
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostCreateComment_RepositoryError(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	comments.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]string{"content": "test"})
	assertStatus(t, w, http.StatusInternalServerError)
	comments.AssertExpectations(t)
}

// ---------- 返信（親コメント）の検証 ----------

func TestPostCreateReply_Success(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	comments.On("FindCommentByID", mock.Anything, uint(10)).Return(&model.Comment{ID: 10, PostID: 5}, nil)
	comments.On("Create", mock.Anything, mock.MatchedBy(func(c *model.Comment) bool {
		return c.ParentID != nil && *c.ParentID == 10
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]interface{}{
		"content":   "Reply to comment!",
		"parent_id": 10,
	})
	assertStatus(t, w, http.StatusCreated)
	comments.AssertExpectations(t)
}

func TestPostCreateReply_ParentNotFound(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	comments.On("FindCommentByID", mock.Anything, uint(99)).Return(nil, gorm.ErrRecordNotFound)

	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]interface{}{
		"content":   "Reply",
		"parent_id": 99,
	})
	assertStatus(t, w, http.StatusNotFound)
	assert.Contains(t, w.Body.String(), "親コメントが見つかりません")
	comments.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestPostCreateReply_ParentOnDifferentPost(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	comments.On("FindCommentByID", mock.Anything, uint(10)).Return(&model.Comment{ID: 10, PostID: 20}, nil)

	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]interface{}{
		"content":   "Reply",
		"parent_id": 10,
	})
	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "別の投稿に属しています")
	comments.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestPostCreateReply_NestedReplyNotAllowed(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments", h.CreateComment)

	grandParentID := uint(3)
	comments.On("FindCommentByID", mock.Anything, uint(10)).
		Return(&model.Comment{ID: 10, PostID: 5, ParentID: &grandParentID}, nil)

	w := doRequest(r, http.MethodPost, "/posts/5/comments", map[string]interface{}{
		"content":   "Nested reply",
		"parent_id": 10,
	})
	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "返信への返信はできません")
	comments.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// ---------- GetReplies ----------

func TestPostGetReplies_Success(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.GET("/comments/:commentId/replies", h.GetReplies)

	comments.On("ListReplies", mock.Anything, uint(10)).
		Return([]model.Comment{{ID: 1, Content: "Reply 1"}, {ID: 2, Content: "Reply 2"}}, nil)

	w := doRequest(r, http.MethodGet, "/comments/10/replies", nil)
	assertStatus(t, w, http.StatusOK)
	comments.AssertExpectations(t)
}

func TestPostGetReplies_Empty(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.GET("/comments/:commentId/replies", h.GetReplies)

	comments.On("ListReplies", mock.Anything, uint(10)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/comments/10/replies", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
}

func TestPostGetReplies_InvalidID(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.GET("/comments/:commentId/replies", h.GetReplies)

	w := doRequest(r, http.MethodGet, "/comments/abc/replies", nil)
	assertStatus(t, w, http.StatusBadRequest)
	comments.AssertNotCalled(t, "ListReplies", mock.Anything, mock.Anything)
}

func TestPostGetReplies_RepositoryError(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.GET("/comments/:commentId/replies", h.GetReplies)

	comments.On("ListReplies", mock.Anything, uint(10)).Return(nil, domain.ErrNotFound)

	w := doRequest(r, http.MethodGet, "/comments/10/replies", nil)
	assertStatus(t, w, http.StatusNotFound)
	comments.AssertExpectations(t)
}

// ---------- EditComment ----------

func TestPostEditComment_Success(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.PUT("/posts/:id/comments/:commentId", h.EditComment)

	comments.On("FindCommentByID", mock.Anything, uint(3)).
		Return(&model.Comment{ID: 3, UserID: 1, Content: "old"}, nil)
	comments.On("Update", mock.Anything, mock.MatchedBy(func(c *model.Comment) bool {
		return c.ID == 3 && c.Content == "new content"
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/posts/5/comments/3", map[string]string{"content": "new content"})
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"content":"new content"`)
	comments.AssertExpectations(t)
}

func TestPostEditComment_Forbidden(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.PUT("/posts/:id/comments/:commentId", h.EditComment)

	comments.On("FindCommentByID", mock.Anything, uint(3)).
		Return(&model.Comment{ID: 3, UserID: 999, Content: "old"}, nil)

	w := doRequest(r, http.MethodPut, "/posts/5/comments/3", map[string]string{"content": "new content"})
	assertStatus(t, w, http.StatusForbidden)
	comments.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// 存在しないコメントは移行前と同じく内部エラーになる。
func TestPostEditComment_NotFound(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.PUT("/posts/:id/comments/:commentId", h.EditComment)

	comments.On("FindCommentByID", mock.Anything, uint(3)).Return(nil, gorm.ErrRecordNotFound)

	w := doRequest(r, http.MethodPut, "/posts/5/comments/3", map[string]string{"content": "new content"})
	assertStatus(t, w, http.StatusInternalServerError)
	comments.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// 所有権チェックの後に本文を検証する。
func TestPostEditComment_BlankContent(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.PUT("/posts/:id/comments/:commentId", h.EditComment)

	comments.On("FindCommentByID", mock.Anything, uint(3)).
		Return(&model.Comment{ID: 3, UserID: 1, Content: "old"}, nil)

	w := doRequest(r, http.MethodPut, "/posts/5/comments/3", map[string]string{"content": "   "})
	assertStatus(t, w, http.StatusBadRequest)
	comments.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestPostEditComment_InvalidID(t *testing.T) {
	h, _ := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.PUT("/posts/:id/comments/:commentId", h.EditComment)

	w := doRequest(r, http.MethodPut, "/posts/5/comments/abc", map[string]string{"content": "new"})
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- DeleteComment ----------

func TestPostDeleteComment_Success(t *testing.T) {
	h, ports := setupPostHandler()
	comments := ports.Comments
	r := newRouter(1)
	r.DELETE("/posts/:id/comments/:commentId", h.DeleteComment)

	comments.On("FindCommentByID", mock.Anything, uint(3)).Return(&model.Comment{ID: 3, UserID: 1}, nil)
	comments.On("Delete", mock.Anything, uint(3)).Return(nil)
	// コメントが消えたら、そのコメントを指すメンションも消す
	ports.Mentions.On("DeleteByCommentID", mock.Anything, uint(3)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/posts/5/comments/3", nil)
	assertStatus(t, w, http.StatusOK)
	comments.AssertExpectations(t)
	ports.Mentions.AssertExpectations(t)
}

func TestPostDeleteComment_Forbidden(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.DELETE("/posts/:id/comments/:commentId", h.DeleteComment)

	comments.On("FindCommentByID", mock.Anything, uint(3)).Return(&model.Comment{ID: 3, UserID: 999}, nil)

	w := doRequest(r, http.MethodDelete, "/posts/5/comments/3", nil)
	assertStatus(t, w, http.StatusForbidden)
	comments.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

func TestPostDeleteComment_InvalidID(t *testing.T) {
	h, _ := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.DELETE("/posts/:id/comments/:commentId", h.DeleteComment)

	w := doRequest(r, http.MethodDelete, "/posts/5/comments/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostDeleteComment_RepositoryError(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.DELETE("/posts/:id/comments/:commentId", h.DeleteComment)

	comments.On("FindCommentByID", mock.Anything, uint(3)).Return(&model.Comment{ID: 3, UserID: 1}, nil)
	comments.On("Delete", mock.Anything, uint(3)).Return(errors.New("db error"))

	w := doRequest(r, http.MethodDelete, "/posts/5/comments/3", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	comments.AssertExpectations(t)
}

// ---------- HideComment / UnhideComment ----------

func TestPostHideComment_Success(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments/:commentId/hide", h.HideComment)

	comments.On("FindCommentByID", mock.Anything, uint(10)).Return(&model.Comment{ID: 10, UserID: 1}, nil)
	comments.On("Update", mock.Anything, mock.MatchedBy(func(c *model.Comment) bool {
		return c.ID == 10 && c.IsHidden
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/comments/10/hide", nil)
	assertStatus(t, w, http.StatusOK)
	comments.AssertExpectations(t)
}

func TestPostHideComment_Forbidden(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments/:commentId/hide", h.HideComment)

	comments.On("FindCommentByID", mock.Anything, uint(10)).Return(&model.Comment{ID: 10, UserID: 999}, nil)

	w := doRequest(r, http.MethodPost, "/posts/5/comments/10/hide", nil)
	assertStatus(t, w, http.StatusForbidden)
	comments.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestPostHideComment_InvalidID(t *testing.T) {
	h, _ := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments/:commentId/hide", h.HideComment)

	w := doRequest(r, http.MethodPost, "/posts/5/comments/abc/hide", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

func TestPostUnhideComment_Success(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments/:commentId/unhide", h.UnhideComment)

	comments.On("FindCommentByID", mock.Anything, uint(10)).
		Return(&model.Comment{ID: 10, UserID: 1, IsHidden: true}, nil)
	comments.On("Update", mock.Anything, mock.MatchedBy(func(c *model.Comment) bool {
		return c.ID == 10 && !c.IsHidden
	})).Return(nil)

	w := doRequest(r, http.MethodPost, "/posts/5/comments/10/unhide", nil)
	assertStatus(t, w, http.StatusOK)
	comments.AssertExpectations(t)
}

func TestPostUnhideComment_Forbidden(t *testing.T) {
	h, comments := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments/:commentId/unhide", h.UnhideComment)

	comments.On("FindCommentByID", mock.Anything, uint(10)).
		Return(&model.Comment{ID: 10, UserID: 999, IsHidden: true}, nil)

	w := doRequest(r, http.MethodPost, "/posts/5/comments/10/unhide", nil)
	assertStatus(t, w, http.StatusForbidden)
	comments.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestPostUnhideComment_InvalidID(t *testing.T) {
	h, _ := setupPostHandlerWithCommentPort()
	r := newRouter(1)
	r.POST("/posts/:id/comments/:commentId/unhide", h.UnhideComment)

	w := doRequest(r, http.MethodPost, "/posts/5/comments/abc/unhide", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

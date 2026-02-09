package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// PostHandler は投稿関連のHTTPハンドラ。
// 投稿のCRUD・いいね・コメント・タイムラインを処理する。
type PostHandler struct {
	service        *service.PostService
	snippetService *service.CodeSnippetService
}

// NewPostHandler は新しいPostHandlerインスタンスを生成する。
func NewPostHandler(s *service.PostService, snippetService *service.CodeSnippetService) *PostHandler {
	return &PostHandler{service: s, snippetService: snippetService}
}

// Create は新しい投稿を作成する。
func (h *PostHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	var input struct {
		Title     string `json:"title" binding:"required"`
		Content   string `json:"content" binding:"required"`
		ImageURLs string `json:"image_urls"`
		CodeSnippets []struct {
			Language string `json:"language"`
			FileName string `json:"file_name"`
			Code     string `json:"code"`
		} `json:"code_snippets"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := &model.Post{
		UserID:    userID,
		Title:     input.Title,
		Content:   input.Content,
		ImageURLs: input.ImageURLs,
	}
	created, err := h.service.Create(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// コードスニペットを一括作成
	if len(input.CodeSnippets) > 0 && h.snippetService != nil {
		for _, s := range input.CodeSnippets {
			if s.Language == "" || s.Code == "" {
				continue
			}
			snippet := &model.CodeSnippet{
				PostID:   created.ID,
				UserID:   userID,
				Language: s.Language,
				FileName: s.FileName,
				Code:     s.Code,
			}
			h.snippetService.Create(snippet)
		}
		// スニペット付きで再取得
		if updated, err := h.service.GetByID(created.ID); err == nil {
			created = updated
		}
	}

	c.JSON(http.StatusCreated, created)
}

// GetAll は投稿一覧をページネーション付きで返す。
func (h *PostHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}

	posts, err := h.service.GetAll(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// GetByID は指定IDの投稿を返す。いいね済みフラグも付与する。
func (h *PostHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	post, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	userID := c.GetUint("userID")
	type postWithLiked struct {
		model.Post
		Liked bool `json:"liked"`
	}
	c.JSON(http.StatusOK, postWithLiked{
		Post:  *post,
		Liked: h.service.HasLiked(userID, post.ID),
	})
}

// Update は投稿を更新する。所有者のみ更新可能。
func (h *PostHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userID := c.GetUint("userID")

	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.service.Update(uint(id), userID, input.Title, input.Content)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not your post"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, post)
}

// Delete は投稿を削除する。所有者のみ削除可能。
func (h *PostHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Delete(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not your post"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Timeline はフォロー中ユーザーの投稿タイムラインを返す。
func (h *PostHandler) Timeline(c *gin.Context) {
	userID := c.GetUint("userID")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}

	posts, err := h.service.Timeline(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// GetUserPosts は指定ユーザーの投稿一覧を返す。
func (h *PostHandler) GetUserPosts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	posts, err := h.service.GetByUserID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// Like は投稿にいいねする。
func (h *PostHandler) Like(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Like(userID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "liked"})
}

// Unlike は投稿のいいねを取り消す。
func (h *PostHandler) Unlike(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Unlike(userID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "unliked"})
}

// GetComments は投稿のコメント一覧を返す。
func (h *PostHandler) GetComments(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	comments, err := h.service.GetComments(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comments)
}

// CreateComment は投稿にコメントを作成する。
func (h *PostHandler) CreateComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userID := c.GetUint("userID")

	var input struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := &model.Comment{UserID: userID, PostID: uint(id), Content: input.Content}
	if err := h.service.CreateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, comment)
}

// DeleteComment はコメントを削除する。所有者のみ削除可能。
func (h *PostHandler) DeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("commentId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.DeleteComment(uint(commentID), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

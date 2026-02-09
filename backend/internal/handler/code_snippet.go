package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// CodeSnippetHandler はコードスニペット関連のHTTPハンドラ。
type CodeSnippetHandler struct {
	service *service.CodeSnippetService
}

// NewCodeSnippetHandler は新しいCodeSnippetHandlerインスタンスを生成する。
func NewCodeSnippetHandler(s *service.CodeSnippetService) *CodeSnippetHandler {
	return &CodeSnippetHandler{service: s}
}

// Create は投稿にコードスニペットを追加する。
func (h *CodeSnippetHandler) Create(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}
	userID := c.GetUint("userID")

	var input struct {
		Language string `json:"language" binding:"required"`
		FileName string `json:"file_name"`
		Code     string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	snippet := &model.CodeSnippet{
		PostID:   uint(postID),
		UserID:   userID,
		Language: input.Language,
		FileName: input.FileName,
		Code:     input.Code,
	}
	created, err := h.service.Create(snippet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// GetByPostID は投稿に紐づくスニペット一覧を返す。
func (h *CodeSnippetHandler) GetByPostID(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	snippets, err := h.service.GetByPostID(uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, snippets)
}

// GetByID は指定IDのスニペットを返す。
func (h *CodeSnippetHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	snippets, err := h.service.GetByPostID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "snippet not found"})
		return
	}
	c.JSON(http.StatusOK, snippets)
}

// Update はスニペットを更新する。所有者のみ更新可能。
func (h *CodeSnippetHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userID := c.GetUint("userID")

	var input struct {
		Language string `json:"language"`
		FileName string `json:"file_name"`
		Code     string `json:"code"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	snippet, err := h.service.Update(uint(id), userID, input.Language, input.FileName, input.Code)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not your snippet"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "snippet not found"})
		return
	}
	c.JSON(http.StatusOK, snippet)
}

// Delete はスニペットを削除する。所有者のみ削除可能。
func (h *CodeSnippetHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Delete(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not your snippet"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "snippet not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GetComments はスニペットのインラインコメント一覧を返す。
func (h *CodeSnippetHandler) GetComments(c *gin.Context) {
	snippetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid snippet id"})
		return
	}

	comments, err := h.service.GetComments(uint(snippetID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comments)
}

// CreateComment はスニペットにインラインコメントを作成する。
func (h *CodeSnippetHandler) CreateComment(c *gin.Context) {
	snippetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid snippet id"})
		return
	}
	userID := c.GetUint("userID")

	var input struct {
		LineNumber int    `json:"line_number" binding:"required"`
		Content    string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := &model.SnippetComment{
		SnippetID:  uint(snippetID),
		UserID:     userID,
		LineNumber: input.LineNumber,
		Content:    input.Content,
	}
	if err := h.service.CreateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, comment)
}

// DeleteComment はインラインコメントを削除する。所有者のみ削除可能。
func (h *CodeSnippetHandler) DeleteComment(c *gin.Context) {
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

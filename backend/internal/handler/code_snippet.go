package handler

import (
	"net/http"

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
	postID, ok := parseID(c, "id")
	if !ok {
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
		PostID:   postID,
		UserID:   userID,
		Language: input.Language,
		FileName: input.FileName,
		Code:     input.Code,
	}
	created, err := h.service.Create(snippet)
	if err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, created)
}

// GetByPostID は投稿に紐づくスニペット一覧を返す。
func (h *CodeSnippetHandler) GetByPostID(c *gin.Context) {
	postID, ok := parseID(c, "id")
	if !ok {
		return
	}

	snippets, err := h.service.GetByPostID(postID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, snippets)
}

// GetByID は指定IDのスニペットを返す。
func (h *CodeSnippetHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	snippets, err := h.service.GetByPostID(id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, snippets)
}

// Update はスニペットを更新する。所有者のみ更新可能。
func (h *CodeSnippetHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
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

	snippet, err := h.service.Update(id, userID, input.Language, input.FileName, input.Code)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, snippet)
}

// Delete はスニペットを削除する。所有者のみ削除可能。
func (h *CodeSnippetHandler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Delete(id, userID); err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

// GetComments はスニペットのインラインコメント一覧を返す。
func (h *CodeSnippetHandler) GetComments(c *gin.Context) {
	snippetID, ok := parseID(c, "id")
	if !ok {
		return
	}

	comments, err := h.service.GetComments(snippetID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, comments)
}

// CreateComment はスニペットにインラインコメントを作成する。
func (h *CodeSnippetHandler) CreateComment(c *gin.Context) {
	snippetID, ok := parseID(c, "id")
	if !ok {
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
		SnippetID:  snippetID,
		UserID:     userID,
		LineNumber: input.LineNumber,
		Content:    input.Content,
	}
	if err := h.service.CreateComment(comment); err != nil {
		respondError(c, err)
		return
	}
	respondCreated(c, comment)
}

// DeleteComment はインラインコメントを削除する。所有者のみ削除可能。
func (h *CodeSnippetHandler) DeleteComment(c *gin.Context) {
	commentID, ok := parseID(c, "commentId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.DeleteComment(commentID, userID); err != nil {
		respondError(c, err)
		return
	}
	respondDeleted(c)
}

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// CodeSnippetHandlerServiceInterface はCodeSnippetHandlerが依存するサービスのインターフェース。
type CodeSnippetHandlerServiceInterface interface {
	Create(snippet *model.CodeSnippet) (*model.CodeSnippet, error)
	GetByPostID(postID uint) ([]model.CodeSnippet, error)
	GetByUserLanguage(userID uint, language string) ([]model.CodeSnippet, error)
	Update(id, userID uint, language, fileName, code string) (*model.CodeSnippet, error)
	Delete(id, userID uint) error
	GetComments(snippetID uint) ([]model.SnippetComment, error)
	CreateComment(comment *model.SnippetComment) error
	DeleteComment(id, userID uint) error
	Fork(userID, snippetID, targetPostID uint) (*model.CodeSnippet, error)
}

// CodeSnippetHandler はコードスニペット関連のHTTPハンドラ。
type CodeSnippetHandler struct {
	service CodeSnippetHandlerServiceInterface
}

// NewCodeSnippetHandler は新しいCodeSnippetHandlerインスタンスを生成する。
func NewCodeSnippetHandler(s CodeSnippetHandlerServiceInterface) *CodeSnippetHandler {
	return &CodeSnippetHandler{service: s}
}

// Create は投稿にコードスニペットを追加する。
func (h *CodeSnippetHandler) Create(c *gin.Context) {
	postID, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[dto.CreateCodeSnippetRequest](c)
	if input == nil {
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
	respondOK(c, ensureSlice(snippets))
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
	respondOK(c, ensureSlice(snippets))
}

// Update はスニペットを更新する。所有者のみ更新可能。
func (h *CodeSnippetHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[dto.UpdateCodeSnippetRequest](c)
	if input == nil {
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
	handleDelete(c, h.service.Delete)
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
	respondOK(c, ensureSlice(comments))
}

// CreateComment はスニペットにインラインコメントを作成する。
func (h *CodeSnippetHandler) CreateComment(c *gin.Context) {
	snippetID, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[dto.CreateSnippetCommentRequest](c)
	if input == nil {
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

// Fork はスニペットをフォークして指定投稿にコピーする。
func (h *CodeSnippetHandler) Fork(c *gin.Context) {
	snippetID, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	req := bindJSON[dto.ForkSnippetRequest](c)
	if req == nil {
		return
	}

	forked, err := h.service.Fork(userID, snippetID, req.TargetPostID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, forked)
}

// GetByUserLanguage は認証ユーザーのスニペットを言語でフィルタリングして取得する。
func (h *CodeSnippetHandler) GetByUserLanguage(c *gin.Context) {
	userID := c.GetUint("userID")
	language := c.Param("language")

	snippets, err := h.service.GetByUserLanguage(userID, language)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(snippets))
}

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// NoteTemplateHandler はノートテンプレート関連のHTTPハンドラ。
type NoteTemplateHandler struct {
	create     *usecase.CreateNoteTemplateUseCase
	get        *usecase.GetNoteTemplateUseCase
	list       *usecase.ListNoteTemplatesUseCase
	getDefault *usecase.GetDefaultNoteTemplateUseCase
	update     *usecase.UpdateNoteTemplateUseCase
	remove     *usecase.DeleteNoteTemplateUseCase
	createNote *usecase.CreateNoteFromTemplateUseCase
	count      *usecase.CountNoteTemplatesUseCase
}

// NewNoteTemplateHandler は新しいNoteTemplateHandlerインスタンスを生成する。
func NewNoteTemplateHandler(
	create *usecase.CreateNoteTemplateUseCase,
	get *usecase.GetNoteTemplateUseCase,
	list *usecase.ListNoteTemplatesUseCase,
	getDefault *usecase.GetDefaultNoteTemplateUseCase,
	update *usecase.UpdateNoteTemplateUseCase,
	remove *usecase.DeleteNoteTemplateUseCase,
	createNote *usecase.CreateNoteFromTemplateUseCase,
	count *usecase.CountNoteTemplatesUseCase,
) *NoteTemplateHandler {
	return &NoteTemplateHandler{
		create: create, get: get, list: list, getDefault: getDefault,
		update: update, remove: remove, createNote: createNote, count: count,
	}
}

// Create は新しいテンプレートを作成する。
func (h *NoteTemplateHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	input := bindJSON[dto.CreateNoteTemplateRequest](c)
	if input == nil {
		return
	}

	template := &model.NoteTemplate{
		UserID:          userID,
		Name:            input.Name,
		Description:     input.Description,
		DefaultTitle:    input.DefaultTitle,
		ContentTemplate: input.ContentTemplate,
		DefaultTags:     input.DefaultTags,
		IsDefault:       input.IsDefault,
	}

	if err := h.create.Execute(c.Request.Context(), template); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, template)
}

// GetByID は指定IDのテンプレートを取得する。所有者のみ取得可能。
func (h *NoteTemplateHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	template, err := h.get.Execute(c.Request.Context(), id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, template)
}

// GetByUserID は現在のユーザーのテンプレート一覧を取得する。
func (h *NoteTemplateHandler) GetByUserID(c *gin.Context) {
	userID := c.GetUint("userID")

	templates, err := h.list.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(templates))
}

// GetDefault は現在のユーザーのデフォルトテンプレートを取得する。
func (h *NoteTemplateHandler) GetDefault(c *gin.Context) {
	userID := c.GetUint("userID")

	template, err := h.getDefault.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, template)
}

// Update はテンプレートを更新する。
func (h *NoteTemplateHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[dto.UpdateNoteTemplateRequest](c)
	if input == nil {
		return
	}

	template, err := h.update.Execute(c.Request.Context(), usecase.UpdateNoteTemplateInput{
		ID:              id,
		UserID:          userID,
		Name:            input.Name,
		Description:     input.Description,
		DefaultTitle:    input.DefaultTitle,
		ContentTemplate: input.ContentTemplate,
		DefaultTags:     input.DefaultTags,
		IsDefault:       input.IsDefault,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, template)
}

// Delete はテンプレートを削除する。
func (h *NoteTemplateHandler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.remove.Execute(c.Request.Context(), id, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

// UseTemplate はテンプレートからノートを作成する。
func (h *NoteTemplateHandler) UseTemplate(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	note, err := h.createNote.Execute(c.Request.Context(), id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, note)
}

// GetMyCount は認証ユーザーのノートテンプレート総数を返す。
func (h *NoteTemplateHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")
	count, err := h.count.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"count": count})
}

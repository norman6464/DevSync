package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// NoteTemplateServiceInterface はNoteTemplateServiceのインターフェース。
type NoteTemplateServiceInterface interface {
	Create(template *model.NoteTemplate) error
	GetByID(id uint) (*model.NoteTemplate, error)
	GetByUserID(userID uint) ([]model.NoteTemplate, error)
	GetDefaultByUserID(userID uint) (*model.NoteTemplate, error)
	Update(id, userID uint, name, description, defaultTitle, contentTemplate, defaultTags string, isDefault *bool) (*model.NoteTemplate, error)
	Delete(id, userID uint) error
	UseTemplate(id, userID uint) (*model.Note, error)
}

// NoteTemplateHandler はノートテンプレート関連のHTTPハンドラ。
type NoteTemplateHandler struct {
	service NoteTemplateServiceInterface
}

// NewNoteTemplateHandler は新しいNoteTemplateHandlerインスタンスを生成する。
func NewNoteTemplateHandler(s NoteTemplateServiceInterface) *NoteTemplateHandler {
	return &NoteTemplateHandler{service: s}
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

	if err := h.service.Create(template); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, template)
}

// GetByID は指定IDのテンプレートを取得する。
func (h *NoteTemplateHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	template, err := h.service.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, template)
}

// GetByUserID は現在のユーザーのテンプレート一覧を取得する。
func (h *NoteTemplateHandler) GetByUserID(c *gin.Context) {
	userID := c.GetUint("userID")

	templates, err := h.service.GetByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(templates))
}

// GetDefault は現在のユーザーのデフォルトテンプレートを取得する。
func (h *NoteTemplateHandler) GetDefault(c *gin.Context) {
	userID := c.GetUint("userID")

	template, err := h.service.GetDefaultByUserID(userID)
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

	template, err := h.service.Update(id, userID, input.Name, input.Description, input.DefaultTitle, input.ContentTemplate, input.DefaultTags, input.IsDefault)
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

	if err := h.service.Delete(id, userID); err != nil {
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

	note, err := h.service.UseTemplate(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, note)
}

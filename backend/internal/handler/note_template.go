package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// NoteTemplateServiceInterface はNoteTemplateServiceのインターフェース。
type NoteTemplateServiceInterface interface {
	Create(template *model.NoteTemplate) error
	GetByID(id uint) (*model.NoteTemplate, error)
	GetByUserID(userID uint) ([]model.NoteTemplate, error)
	GetDefaultByUserID(userID uint) (*model.NoteTemplate, error)
	Update(template *model.NoteTemplate) error
	Delete(id uint) error
}

// NoteTemplateHandler はノートテンプレート関連のHTTPハンドラ。
type NoteTemplateHandler struct {
	service     NoteTemplateServiceInterface
	noteService NoteServiceInterface
}

// NewNoteTemplateHandler は新しいNoteTemplateHandlerインスタンスを生成する。
func NewNoteTemplateHandler(s NoteTemplateServiceInterface, noteService NoteServiceInterface) *NoteTemplateHandler {
	return &NoteTemplateHandler{
		service:     s,
		noteService: noteService,
	}
}

// CreateTemplateInput はテンプレート作成のリクエストボディ。
type CreateTemplateInput struct {
	Name            string `json:"name" binding:"required"`
	Description     string `json:"description"`
	DefaultTitle    string `json:"default_title"`
	ContentTemplate string `json:"content_template" binding:"required"`
	DefaultTags     string `json:"default_tags"`
	IsDefault       bool   `json:"is_default"`
}

// Create は新しいテンプレートを作成する。
func (h *NoteTemplateHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	input := bindJSON[CreateTemplateInput](c)
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

	respondOK(c, templates)
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

// UpdateTemplateInput はテンプレート更新のリクエストボディ。
type UpdateTemplateInput struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	DefaultTitle    string `json:"default_title"`
	ContentTemplate string `json:"content_template"`
	DefaultTags     string `json:"default_tags"`
	IsDefault       *bool  `json:"is_default"`
}

// Update はテンプレートを更新する。
func (h *NoteTemplateHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	input := bindJSON[UpdateTemplateInput](c)
	if input == nil {
		return
	}

	// 既存のテンプレートを取得して所有者確認
	template, err := h.service.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	// 所有者チェック
	if template.UserID != userID {
		respondForbidden(c, "この操作を行う権限がありません")
		return
	}

	// 更新フィールドを適用
	if input.Name != "" {
		template.Name = input.Name
	}
	if input.Description != "" {
		template.Description = input.Description
	}
	if input.DefaultTitle != "" {
		template.DefaultTitle = input.DefaultTitle
	}
	if input.ContentTemplate != "" {
		template.ContentTemplate = input.ContentTemplate
	}
	if input.DefaultTags != "" {
		template.DefaultTags = input.DefaultTags
	}
	if input.IsDefault != nil {
		template.IsDefault = *input.IsDefault
	}

	if err := h.service.Update(template); err != nil {
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

	// 所有者確認
	template, err := h.service.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	if template.UserID != userID {
		respondForbidden(c, "この操作を行う権限がありません")
		return
	}

	if err := h.service.Delete(id); err != nil {
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

	// NoteServiceからNoteRepositoryを取得する必要があるが、
	// 簡便化のため、NoteServiceのインターフェースを使用する想定
	// 実装上は、noteServiceを経由してノートを作成する
	template, err := h.service.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	// 所有者チェック
	if template.UserID != userID {
		respondForbidden(c, "この操作を行う権限がありません")
		return
	}

	// テンプレートからノートを作成
	note := &model.Note{
		UserID:  userID,
		Title:   template.DefaultTitle,
		Content: template.ContentTemplate,
		Tags:    template.DefaultTags,
	}

	// タイトルが空の場合、デフォルト値を設定
	if note.Title == "" {
		note.Title = "新しいノート"
	}

	// NoteServiceを使ってノートを作成（バリデーション含む）
	if err := h.noteService.Create(note); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, note)
}

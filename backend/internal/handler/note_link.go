package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// NoteLinkServiceInterface はNoteLinkServiceのインターフェース。
type NoteLinkServiceInterface interface {
	CreateLink(sourceNoteID, targetNoteID uint) error
	GetLinks(sourceNoteID uint) ([]model.NoteLink, error)
	GetBacklinks(targetNoteID uint) ([]model.NoteLink, error)
	DeleteLink(sourceNoteID, targetNoteID uint) error
}

// NoteLinkHandler はノート間リンク関連のHTTPハンドラ。
type NoteLinkHandler struct {
	service NoteLinkServiceInterface
}

// NewNoteLinkHandler は新しいNoteLinkHandlerインスタンスを生成する。
func NewNoteLinkHandler(s NoteLinkServiceInterface) *NoteLinkHandler {
	return &NoteLinkHandler{service: s}
}

// CreateLinkInput はリンク作成のリクエストボディ。
type CreateLinkInput struct {
	TargetNoteID uint `json:"target_note_id" binding:"required"`
}

// CreateLink は新しいリンクを作成する。
func (h *NoteLinkHandler) CreateLink(c *gin.Context) {
	sourceNoteID, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[CreateLinkInput](c)
	if input == nil {
		return
	}

	if err := h.service.CreateLink(sourceNoteID, input.TargetNoteID); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, gin.H{"message": "リンクを作成しました"})
}

// GetLinks は指定ノートからのリンク一覧を取得する。
func (h *NoteLinkHandler) GetLinks(c *gin.Context) {
	noteID, ok := parseID(c, "id")
	if !ok {
		return
	}

	links, err := h.service.GetLinks(noteID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, links)
}

// GetBacklinks は指定ノートへのリンク一覧（バックリンク）を取得する。
func (h *NoteLinkHandler) GetBacklinks(c *gin.Context) {
	noteID, ok := parseID(c, "id")
	if !ok {
		return
	}

	backlinks, err := h.service.GetBacklinks(noteID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, backlinks)
}

// DeleteLink はリンクを削除する。
func (h *NoteLinkHandler) DeleteLink(c *gin.Context) {
	sourceNoteID, ok := parseID(c, "id")
	if !ok {
		return
	}
	targetNoteID, ok := parseID(c, "targetId")
	if !ok {
		return
	}

	if err := h.service.DeleteLink(sourceNoteID, targetNoteID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

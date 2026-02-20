package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// NoteLinkServiceInterface はNoteLinkServiceのインターフェース。
type NoteLinkServiceInterface interface {
	CreateLink(sourceNoteID, targetNoteID, userID uint) error
	GetLinks(sourceNoteID uint) ([]model.NoteLink, error)
	GetBacklinks(targetNoteID uint) ([]model.NoteLink, error)
	DeleteLink(sourceNoteID, targetNoteID, userID uint) error
	GetLinkStats(noteID, userID uint) (*model.NoteLinkStats, error)
}

// NoteLinkHandler はノート間リンク関連のHTTPハンドラ。
type NoteLinkHandler struct {
	service NoteLinkServiceInterface
}

// NewNoteLinkHandler は新しいNoteLinkHandlerインスタンスを生成する。
func NewNoteLinkHandler(s NoteLinkServiceInterface) *NoteLinkHandler {
	return &NoteLinkHandler{service: s}
}

// CreateLink は新しいリンクを作成する。
func (h *NoteLinkHandler) CreateLink(c *gin.Context) {
	sourceNoteID, ok := parseID(c, "id")
	if !ok {
		return
	}

	userID := c.GetUint("userID")

	input := bindJSON[dto.CreateNoteLinkRequest](c)
	if input == nil {
		return
	}

	if err := h.service.CreateLink(sourceNoteID, input.TargetNoteID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, domain.NewMessageResponse("リンクを作成しました"))
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

	respondOK(c, ensureSlice(links))
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

	respondOK(c, ensureSlice(backlinks))
}

// GetLinkStats はノートのリンク統計を取得する。
func (h *NoteLinkHandler) GetLinkStats(c *gin.Context) {
	noteID, ok := parseID(c, "id")
	if !ok {
		return
	}

	userID := c.GetUint("userID")

	stats, err := h.service.GetLinkStats(noteID, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
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

	userID := c.GetUint("userID")

	if err := h.service.DeleteLink(sourceNoteID, targetNoteID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// NoteLinkHandler はノート間リンク関連のHTTPハンドラ。
type NoteLinkHandler struct {
	createLink   *usecase.CreateNoteLinkUseCase
	listLinks    *usecase.ListNoteLinksUseCase
	listBacklink *usecase.ListNoteBacklinksUseCase
	stats        *usecase.GetNoteLinkStatsUseCase
	deleteLink   *usecase.DeleteNoteLinkUseCase
}

// NewNoteLinkHandler は新しいNoteLinkHandlerインスタンスを生成する。
func NewNoteLinkHandler(
	createLink *usecase.CreateNoteLinkUseCase,
	listLinks *usecase.ListNoteLinksUseCase,
	listBacklink *usecase.ListNoteBacklinksUseCase,
	stats *usecase.GetNoteLinkStatsUseCase,
	deleteLink *usecase.DeleteNoteLinkUseCase,
) *NoteLinkHandler {
	return &NoteLinkHandler{
		createLink:   createLink,
		listLinks:    listLinks,
		listBacklink: listBacklink,
		stats:        stats,
		deleteLink:   deleteLink,
	}
}

// CreateLink は新しいリンクを作成する。
func (h *NoteLinkHandler) CreateLink(c *gin.Context) {
	sourceNoteID, ok := parseID(c, "id")
	if !ok {
		return
	}

	input := bindJSON[dto.CreateNoteLinkRequest](c)
	if input == nil {
		return
	}

	if err := h.createLink.Execute(c.Request.Context(), sourceNoteID, input.TargetNoteID, c.GetUint("userID")); err != nil {
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

	links, err := h.listLinks.Execute(c.Request.Context(), noteID)
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

	backlinks, err := h.listBacklink.Execute(c.Request.Context(), noteID)
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

	stats, err := h.stats.Execute(c.Request.Context(), noteID, c.GetUint("userID"))
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

	if err := h.deleteLink.Execute(c.Request.Context(), sourceNoteID, targetNoteID, c.GetUint("userID")); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

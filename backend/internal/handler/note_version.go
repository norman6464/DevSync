package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// NoteVersionHandler はノートバージョン履歴関連のHTTPハンドラ。
type NoteVersionHandler struct {
	list    *usecase.ListNoteVersionsUseCase
	get     *usecase.GetNoteVersionUseCase
	restore *usecase.RestoreNoteVersionUseCase
}

// NewNoteVersionHandler は新しいNoteVersionHandlerインスタンスを生成する。
func NewNoteVersionHandler(
	list *usecase.ListNoteVersionsUseCase,
	get *usecase.GetNoteVersionUseCase,
	restore *usecase.RestoreNoteVersionUseCase,
) *NoteVersionHandler {
	return &NoteVersionHandler{list: list, get: get, restore: restore}
}

// noteVersionListResponse はノートバージョン一覧レスポンス（ページネーション付き）。
type noteVersionListResponse struct {
	Versions []model.NoteVersion `json:"versions"`
	Total    int64               `json:"total"`
	Limit    int                 `json:"limit"`
	Offset   int                 `json:"offset"`
}

// GetVersions はノートのバージョン履歴一覧を返す。
func (h *NoteVersionHandler) GetVersions(c *gin.Context) {
	noteID, ok := parseID(c, "id")
	if !ok {
		return
	}
	limit, offset := parseLimitOffset(c)

	versions, total, err := h.list.Execute(c.Request.Context(), noteID, c.GetUint("userID"), limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, noteVersionListResponse{
		Versions: ensureSlice(versions),
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	})
}

// GetVersion はバージョンの詳細を返す。
func (h *NoteVersionHandler) GetVersion(c *gin.Context) {
	noteID, ok := parseID(c, "id")
	if !ok {
		return
	}
	versionID, ok := parseID(c, "versionId")
	if !ok {
		return
	}

	version, err := h.get.Execute(c.Request.Context(), noteID, versionID, c.GetUint("userID"))
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, version)
}

// RestoreVersion はバージョンの内容でノートを復元する。
func (h *NoteVersionHandler) RestoreVersion(c *gin.Context) {
	noteID, ok := parseID(c, "id")
	if !ok {
		return
	}
	versionID, ok := parseID(c, "versionId")
	if !ok {
		return
	}

	note, err := h.restore.Execute(c.Request.Context(), noteID, versionID, c.GetUint("userID"))
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, note)
}

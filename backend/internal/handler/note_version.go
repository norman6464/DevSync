package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// NoteVersionServiceInterface はNoteVersionServiceのインターフェース。
type NoteVersionServiceInterface interface {
	SaveVersion(noteID, userID uint) error
	GetVersions(noteID, userID uint, limit, offset int) ([]model.NoteVersion, int64, error)
	GetVersion(noteID, versionID, userID uint) (*model.NoteVersion, error)
	RestoreVersion(noteID, versionID, userID uint) (*model.Note, error)
}

// NoteVersionHandler はノートバージョン履歴関連のHTTPハンドラ。
type NoteVersionHandler struct {
	service NoteVersionServiceInterface
}

// NewNoteVersionHandler は新しいNoteVersionHandlerインスタンスを生成する。
func NewNoteVersionHandler(s NoteVersionServiceInterface) *NoteVersionHandler {
	return &NoteVersionHandler{service: s}
}

// GetVersions はノートのバージョン履歴一覧を返す。
func (h *NoteVersionHandler) GetVersions(c *gin.Context) {
	userID := c.GetUint("userID")
	noteID, ok := parseID(c, "id")
	if !ok {
		return
	}
	limit, offset := parseLimitOffset(c)

	versions, total, err := h.service.GetVersions(noteID, userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.NoteVersionListResponse{
		Versions: ensureSlice(versions),
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	})
}

// GetVersion はバージョンの詳細を返す。
func (h *NoteVersionHandler) GetVersion(c *gin.Context) {
	userID := c.GetUint("userID")
	noteID, ok := parseID(c, "id")
	if !ok {
		return
	}
	versionID, ok := parseID(c, "versionId")
	if !ok {
		return
	}

	version, err := h.service.GetVersion(noteID, versionID, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, version)
}

// RestoreVersion はバージョンの内容でノートを復元する。
func (h *NoteVersionHandler) RestoreVersion(c *gin.Context) {
	userID := c.GetUint("userID")
	noteID, ok := parseID(c, "id")
	if !ok {
		return
	}
	versionID, ok := parseID(c, "versionId")
	if !ok {
		return
	}

	note, err := h.service.RestoreVersion(noteID, versionID, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, note)
}

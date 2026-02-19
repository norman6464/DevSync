package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// NoteStatsServiceInterface はNoteStatsHandlerが依存するサービスメソッドを定義する。
type NoteStatsServiceInterface interface {
	GetNoteStats(userID uint) (*model.NoteStats, error)
}

// NoteStatsHandler はノート統計関連のHTTPハンドラ。
type NoteStatsHandler struct {
	service NoteStatsServiceInterface
}

// NewNoteStatsHandler は新しいNoteStatsHandlerインスタンスを生成する。
func NewNoteStatsHandler(s NoteStatsServiceInterface) *NoteStatsHandler {
	return &NoteStatsHandler{service: s}
}

// GetStats は指定ユーザーのノート集計統計を返す。
func (h *NoteStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.service.GetNoteStats(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}

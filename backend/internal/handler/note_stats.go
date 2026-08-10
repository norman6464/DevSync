package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// NoteStatsHandler はノート統計関連の HTTP ハンドラ。
type NoteStatsHandler struct {
	getStats *usecase.GetNoteStatsUseCase
}

// NewNoteStatsHandler は NoteStatsHandler を生成する。
func NewNoteStatsHandler(getStats *usecase.GetNoteStatsUseCase) *NoteStatsHandler {
	return &NoteStatsHandler{getStats: getStats}
}

// GetStats は指定ユーザーのノート集計統計を返す。
func (h *NoteStatsHandler) GetStats(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	stats, err := h.getStats.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, stats)
}

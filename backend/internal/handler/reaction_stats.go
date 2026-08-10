package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// ReactionStatsHandler はユーザーリアクション統計関連の HTTP ハンドラ。
type ReactionStatsHandler struct {
	getStats   *usecase.GetReactionStatsUseCase
	getSummary *usecase.GetReactionSummaryUseCase
}

// NewReactionStatsHandler は ReactionStatsHandler を生成する。
func NewReactionStatsHandler(
	getStats *usecase.GetReactionStatsUseCase,
	getSummary *usecase.GetReactionSummaryUseCase,
) *ReactionStatsHandler {
	return &ReactionStatsHandler{getStats: getStats, getSummary: getSummary}
}

// GetStats は指定ユーザーのリアクション集計統計を返す。
func (h *ReactionStatsHandler) GetStats(c *gin.Context) {
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

// GetSummary は指定ユーザーのリアクションサマリー（絵文字別集計＋トップ投稿）を返す。
func (h *ReactionStatsHandler) GetSummary(c *gin.Context) {
	userID, ok := parseID(c, "id")
	if !ok {
		return
	}

	summary, err := h.getSummary.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, summary)
}

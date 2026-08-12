package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// LevelHandler はレベルシステム関連のHTTPハンドラ。
// ユーザーのレベル情報とXP内訳の取得を処理する。
type LevelHandler struct {
	levelInfo *usecase.GetLevelInfoUseCase
	breakdown *usecase.GetXPBreakdownUseCase
}

// NewLevelHandler は新しいLevelHandlerインスタンスを生成する。
func NewLevelHandler(
	levelInfo *usecase.GetLevelInfoUseCase,
	breakdown *usecase.GetXPBreakdownUseCase,
) *LevelHandler {
	return &LevelHandler{levelInfo: levelInfo, breakdown: breakdown}
}

// GetMyLevelInfo は認証済みユーザー自身のレベル情報を返す。
func (h *LevelHandler) GetMyLevelInfo(c *gin.Context) {
	userID := c.GetUint("userID")

	info, err := h.levelInfo.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, info)
}

// GetLevelInfo は指定ユーザーのレベル情報を返す。
func (h *LevelHandler) GetLevelInfo(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	info, err := h.levelInfo.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, info)
}

// GetXPBreakdown は指定ユーザーのXP内訳を返す。
func (h *LevelHandler) GetXPBreakdown(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	breakdown, err := h.breakdown.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, breakdown)
}

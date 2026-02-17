package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// LevelServiceInterface はLevelHandlerが依存するサービスメソッドを定義する。
type LevelServiceInterface interface {
	GetLevelInfo(userID uint) (*model.LevelInfo, error)
	GetXPBreakdown(userID uint) (*model.XPBreakdown, error)
}

// LevelHandler はレベルシステム関連のHTTPハンドラ。
// ユーザーのレベル情報とXP内訳の取得を処理する。
type LevelHandler struct {
	service LevelServiceInterface
}

// NewLevelHandler は新しいLevelHandlerインスタンスを生成する。
func NewLevelHandler(s LevelServiceInterface) *LevelHandler {
	return &LevelHandler{service: s}
}

// GetMyLevelInfo は認証済みユーザー自身のレベル情報を返す。
func (h *LevelHandler) GetMyLevelInfo(c *gin.Context) {
	userID := c.GetUint("userID")

	info, err := h.service.GetLevelInfo(userID)
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

	info, err := h.service.GetLevelInfo(userID)
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

	breakdown, err := h.service.GetXPBreakdown(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, breakdown)
}

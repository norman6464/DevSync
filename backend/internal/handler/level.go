package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/service"
)

// LevelHandler はレベルシステム関連のHTTPハンドラ。
// ユーザーのレベル情報とXP内訳の取得を処理する。
type LevelHandler struct {
	service *service.LevelService
}

// NewLevelHandler は新しいLevelHandlerインスタンスを生成する。
func NewLevelHandler(s *service.LevelService) *LevelHandler {
	return &LevelHandler{service: s}
}

// GetMyLevelInfo は認証済みユーザー自身のレベル情報を返す。
func (h *LevelHandler) GetMyLevelInfo(c *gin.Context) {
	userID := c.GetUint("userID")

	info, err := h.service.GetLevelInfo(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

// GetLevelInfo は指定ユーザーのレベル情報を返す。
func (h *LevelHandler) GetLevelInfo(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	info, err := h.service.GetLevelInfo(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

// GetXPBreakdown は指定ユーザーのXP内訳を返す。
func (h *LevelHandler) GetXPBreakdown(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	breakdown, err := h.service.GetXPBreakdown(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, breakdown)
}

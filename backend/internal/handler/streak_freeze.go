package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
)

// StreakFreezeServiceInterface はStreakFreezeHandlerが依存するサービスのインターフェース。
type StreakFreezeServiceInterface interface {
	UseFreeze(userID uint) error
	GetFreezeStatus(userID uint) (*model.StreakFreezeStatus, error)
	GetFreezeDates(userID uint) ([]string, error)
}

// StreakFreezeHandler はストリークフリーズ関連のHTTPハンドラ。
type StreakFreezeHandler struct {
	service StreakFreezeServiceInterface
}

// NewStreakFreezeHandler は新しいStreakFreezeHandlerインスタンスを生成する。
func NewStreakFreezeHandler(s StreakFreezeServiceInterface) *StreakFreezeHandler {
	return &StreakFreezeHandler{service: s}
}

// UseFreeze は今日のストリークフリーズを使用する。
func (h *StreakFreezeHandler) UseFreeze(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.service.UseFreeze(userID); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, domain.NewMessageResponse("ストリークフリーズを使用しました"))
}

// GetStatus は今月のフリーズ使用状況を返す。
func (h *StreakFreezeHandler) GetStatus(c *gin.Context) {
	userID := c.GetUint("userID")

	status, err := h.service.GetFreezeStatus(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, status)
}

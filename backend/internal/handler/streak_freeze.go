package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// StreakFreezeHandler はストリークフリーズ関連の HTTP ハンドラ。
type StreakFreezeHandler struct {
	useFreeze *usecase.UseStreakFreezeUseCase
	getStatus *usecase.GetStreakFreezeStatusUseCase
}

// NewStreakFreezeHandler は StreakFreezeHandler を生成する。
func NewStreakFreezeHandler(
	useFreeze *usecase.UseStreakFreezeUseCase,
	getStatus *usecase.GetStreakFreezeStatusUseCase,
) *StreakFreezeHandler {
	return &StreakFreezeHandler{useFreeze: useFreeze, getStatus: getStatus}
}

// UseFreeze は今日のストリークフリーズを使用する。
func (h *StreakFreezeHandler) UseFreeze(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.useFreeze.Execute(c.Request.Context(), userID); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, domain.NewMessageResponse("ストリークフリーズを使用しました"))
}

// GetStatus は今月のフリーズ使用状況を返す。
func (h *StreakFreezeHandler) GetStatus(c *gin.Context) {
	userID := c.GetUint("userID")

	status, err := h.getStatus.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, status)
}

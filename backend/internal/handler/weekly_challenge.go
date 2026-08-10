package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// WeeklyChallengeHandler はウィークリーチャレンジの HTTP ハンドラ。
type WeeklyChallengeHandler struct {
	getCurrent     *usecase.GetCurrentWeeklyChallengeUseCase
	updateProgress *usecase.UpdateWeeklyChallengeProgressUseCase
}

// NewWeeklyChallengeHandler は WeeklyChallengeHandler を生成する。
func NewWeeklyChallengeHandler(
	getCurrent *usecase.GetCurrentWeeklyChallengeUseCase,
	updateProgress *usecase.UpdateWeeklyChallengeProgressUseCase,
) *WeeklyChallengeHandler {
	return &WeeklyChallengeHandler{getCurrent: getCurrent, updateProgress: updateProgress}
}

// GetCurrent は今週のチャレンジを返す。
func (h *WeeklyChallengeHandler) GetCurrent(c *gin.Context) {
	userID := c.GetUint("userID")

	challenge, err := h.getCurrent.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, challenge)
}

type updateProgressRequest struct {
	Value int `json:"value" binding:"required"`
}

// UpdateProgress はチャレンジの進捗を更新する。
func (h *WeeklyChallengeHandler) UpdateProgress(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[updateProgressRequest](c)
	if input == nil {
		return
	}

	challenge, err := h.updateProgress.Execute(c.Request.Context(), userID, input.Value)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, challenge)
}

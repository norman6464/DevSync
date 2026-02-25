package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// WeeklyChallengeServiceInterface はWeeklyChallengeHandlerが依存するサービスのインターフェース。
type WeeklyChallengeServiceInterface interface {
	GetCurrentChallenge(userID uint) (*model.WeeklyChallenge, error)
	UpdateProgress(userID uint, value int) (*model.WeeklyChallenge, error)
}

// WeeklyChallengeHandler はウィークリーチャレンジのHTTPハンドラ。
type WeeklyChallengeHandler struct {
	service WeeklyChallengeServiceInterface
}

// NewWeeklyChallengeHandler は新しいWeeklyChallengeHandlerを生成する。
func NewWeeklyChallengeHandler(s WeeklyChallengeServiceInterface) *WeeklyChallengeHandler {
	return &WeeklyChallengeHandler{service: s}
}

// GetCurrent は今週のチャレンジを返す。
func (h *WeeklyChallengeHandler) GetCurrent(c *gin.Context) {
	userID := c.GetUint("userID")

	challenge, err := h.service.GetCurrentChallenge(userID)
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

	challenge, err := h.service.UpdateProgress(userID, input.Value)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, challenge)
}

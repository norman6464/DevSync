package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// AtCoderHandler はAtCoder関連のHTTPハンドラ。
type AtCoderHandler struct {
	getRating  *usecase.GetAtCoderRatingUseCase
	connect    *usecase.ConnectAtCoderUseCase
	disconnect *usecase.DisconnectAtCoderUseCase
}

// NewAtCoderHandler は新しいAtCoderHandlerインスタンスを生成する。
func NewAtCoderHandler(
	getRating *usecase.GetAtCoderRatingUseCase,
	connect *usecase.ConnectAtCoderUseCase,
	disconnect *usecase.DisconnectAtCoderUseCase,
) *AtCoderHandler {
	return &AtCoderHandler{
		getRating:  getRating,
		connect:    connect,
		disconnect: disconnect,
	}
}

// GetRating は指定ユーザーのAtCoderレーティング情報を返す。
func (h *AtCoderHandler) GetRating(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		respondBadRequest(c, "ユーザー名は必須です")
		return
	}

	info, err := h.getRating.Execute(c.Request.Context(), username)
	if err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	respondOK(c, info)
}

// Connect はAtCoderユーザー名を検証し、ユーザープロフィールに保存する。
func (h *AtCoderHandler) Connect(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[connectUsernameRequest](c)
	if input == nil {
		return
	}

	user, err := h.connect.Execute(c.Request.Context(), userID, input.Username)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, user)
}

// Disconnect はAtCoderユーザー名をクリアする。
func (h *AtCoderHandler) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := h.disconnect.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, user)
}

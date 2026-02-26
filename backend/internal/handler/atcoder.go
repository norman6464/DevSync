package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// AtCoderServiceInterface はAtCoderサービスの抽象インターフェース。
type AtCoderServiceInterface interface {
	GetRating(username string) (*service.AtCoderRatingInfo, error)
	ConnectAtCoder(userID uint, username string) (*model.User, error)
	DisconnectAtCoder(userID uint) (*model.User, error)
}

// AtCoderHandler はAtCoder関連のHTTPハンドラ。
type AtCoderHandler struct {
	atcoderService AtCoderServiceInterface
}

// NewAtCoderHandler は新しいAtCoderHandlerインスタンスを生成する。
func NewAtCoderHandler(atcoderService AtCoderServiceInterface) *AtCoderHandler {
	return &AtCoderHandler{
		atcoderService: atcoderService,
	}
}

// GetRating は指定ユーザーのAtCoderレーティング情報を返す。
func (h *AtCoderHandler) GetRating(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		respondBadRequest(c, "ユーザー名は必須です")
		return
	}

	info, err := h.atcoderService.GetRating(username)
	if err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	respondOK(c, info)
}

// Connect はAtCoderユーザー名を検証し、ユーザープロフィールに保存する。
func (h *AtCoderHandler) Connect(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[dto.ConnectUsernameRequest](c)
	if input == nil {
		return
	}

	user, err := h.atcoderService.ConnectAtCoder(userID, input.Username)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, user)
}

// Disconnect はAtCoderユーザー名をクリアする。
func (h *AtCoderHandler) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := h.atcoderService.DisconnectAtCoder(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, user)
}

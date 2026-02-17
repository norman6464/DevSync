package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// AtCoderServiceInterface はAtCoderサービスの抽象インターフェース。
type AtCoderServiceInterface interface {
	GetRating(username string) (*service.AtCoderRatingInfo, error)
	ValidateUsername(username string) bool
}

// AtCoderUserServiceInterface はAtCoderハンドラーが必要とするユーザーサービスの抽象インターフェース。
type AtCoderUserServiceInterface interface {
	GetByID(id uint) (*model.User, error)
	Update(user *model.User) error
}

// AtCoderHandler はAtCoder関連のHTTPハンドラ。
type AtCoderHandler struct {
	atcoderService AtCoderServiceInterface
	userService    AtCoderUserServiceInterface
}

// NewAtCoderHandler は新しいAtCoderHandlerインスタンスを生成する。
func NewAtCoderHandler(atcoderService AtCoderServiceInterface, userService AtCoderUserServiceInterface) *AtCoderHandler {
	return &AtCoderHandler{
		atcoderService: atcoderService,
		userService:    userService,
	}
}

// GetRating は指定ユーザーのAtCoderレーティング情報を返す。
func (h *AtCoderHandler) GetRating(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		respondBadRequest(c, "username is required")
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

	var input struct {
		Username string `json:"username" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		respondBadRequest(c, "username is required")
		return
	}

	if !h.atcoderService.ValidateUsername(input.Username) {
		respondBadRequest(c, "invalid AtCoder username")
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		respondNotFound(c, "user not found")
		return
	}

	user.AtCoderUsername = input.Username
	if err := h.userService.Update(user); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, user)
}

// Disconnect はAtCoderユーザー名をクリアする。
func (h *AtCoderHandler) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := h.userService.GetByID(userID)
	if err != nil {
		respondNotFound(c, "user not found")
		return
	}

	user.AtCoderUsername = ""
	if err := h.userService.Update(user); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, user)
}

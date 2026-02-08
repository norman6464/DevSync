package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/service"
)

// AtCoderHandler はAtCoder関連のHTTPハンドラ。
type AtCoderHandler struct {
	atcoderService *service.AtCoderService
	userService    *service.UserService
}

// NewAtCoderHandler は新しいAtCoderHandlerインスタンスを生成する。
func NewAtCoderHandler(atcoderService *service.AtCoderService, userService *service.UserService) *AtCoderHandler {
	return &AtCoderHandler{
		atcoderService: atcoderService,
		userService:    userService,
	}
}

// GetRating は指定ユーザーのAtCoderレーティング情報を返す。
func (h *AtCoderHandler) GetRating(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	info, err := h.atcoderService.GetRating(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

// Connect はAtCoderユーザー名を検証し、ユーザープロフィールに保存する。
func (h *AtCoderHandler) Connect(c *gin.Context) {
	userID := c.GetUint("userID")

	var input struct {
		Username string `json:"username" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	if !h.atcoderService.ValidateUsername(input.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid AtCoder username"})
		return
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.AtCoderUsername = input.Username
	if err := h.userService.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Disconnect はAtCoderユーザー名をクリアする。
func (h *AtCoderHandler) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := h.userService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.AtCoderUsername = ""
	if err := h.userService.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

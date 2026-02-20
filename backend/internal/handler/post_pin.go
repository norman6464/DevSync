package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// PostPinServiceInterface はPostPinHandlerが依存するサービスインターフェース。
type PostPinServiceInterface interface {
	Pin(userID, postID uint) error
	Unpin(userID, postID uint) error
	GetByUserID(userID uint) ([]model.PostPin, error)
	Reorder(userID uint, postIDs []uint) error
	IsPinned(userID, postID uint) (bool, error)
}

// PostPinHandler は投稿ピン留めのHTTPハンドラー。
type PostPinHandler struct {
	service PostPinServiceInterface
}

// NewPostPinHandler は新しいPostPinHandlerを生成する。
func NewPostPinHandler(service PostPinServiceInterface) *PostPinHandler {
	return &PostPinHandler{service: service}
}

// Pin は投稿をピン留めする。
func (h *PostPinHandler) Pin(c *gin.Context) {
	postID, ok := parseID(c, "postId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Pin(userID, postID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("ピン留めしました"))
}

// Unpin は投稿のピン留めを解除する。
func (h *PostPinHandler) Unpin(c *gin.Context) {
	postID, ok := parseID(c, "postId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Unpin(userID, postID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("ピン留めを解除しました"))
}

// GetByUserID はユーザーのピン留め投稿を取得する。
func (h *PostPinHandler) GetByUserID(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	pins, err := h.service.GetByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, dto.PinsResponse{Pins: pins})
}

// Reorder はピン留めの表示順序を変更する。
func (h *PostPinHandler) Reorder(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[dto.ReorderPinsRequest](c)
	if req == nil {
		return
	}

	if err := h.service.Reorder(userID, req.PostIDs); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("順序を更新しました"))
}

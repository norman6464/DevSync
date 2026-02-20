package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// PostViewServiceInterface はPostViewHandlerが依存するサービスインターフェース。
type PostViewServiceInterface interface {
	RecordView(userID, postID uint) error
	GetViewCount(postID uint) (int64, error)
	HasViewed(userID, postID uint) (bool, error)
	GetMostViewed(limit int) ([]model.ViewCount, error)
}

// PostViewHandler は投稿閲覧数のHTTPハンドラー。
type PostViewHandler struct {
	service PostViewServiceInterface
}

// NewPostViewHandler は新しいPostViewHandlerを生成する。
func NewPostViewHandler(service PostViewServiceInterface) *PostViewHandler {
	return &PostViewHandler{service: service}
}

// RecordView は投稿の閲覧を記録する。
func (h *PostViewHandler) RecordView(c *gin.Context) {
	postID, ok := parseID(c, "postId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.RecordView(userID, postID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("記録しました"))
}

// GetViewCount は投稿の閲覧数を取得する。
func (h *PostViewHandler) GetViewCount(c *gin.Context) {
	postID, ok := parseID(c, "postId")
	if !ok {
		return
	}

	count, err := h.service.GetViewCount(postID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, dto.ViewCountResponse{PostID: postID, ViewCount: count})
}

// GetMostViewed は閲覧数の多い投稿ランキングを取得する。
func (h *PostViewHandler) GetMostViewed(c *gin.Context) {
	result, err := h.service.GetMostViewed(20)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, result)
}

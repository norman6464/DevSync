package handler

import (
	"github.com/gin-gonic/gin"
)

// CommentLikeServiceInterface はCommentLikeHandlerが依存するサービスインターフェース。
type CommentLikeServiceInterface interface {
	Like(userID, commentID uint) error
	Unlike(userID, commentID uint) error
	GetStatus(userID, commentID uint) (bool, int64, error)
}

// CommentLikeHandler はコメントへのいいね機能のHTTPハンドラー。
type CommentLikeHandler struct {
	service CommentLikeServiceInterface
}

// NewCommentLikeHandler は新しいCommentLikeHandlerを生成する。
func NewCommentLikeHandler(service CommentLikeServiceInterface) *CommentLikeHandler {
	return &CommentLikeHandler{service: service}
}

// Like はコメントにいいねする。
func (h *CommentLikeHandler) Like(c *gin.Context) {
	commentID, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Like(userID, commentID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"message": "いいねしました"})
}

// Unlike はコメントのいいねを取り消す。
func (h *CommentLikeHandler) Unlike(c *gin.Context) {
	commentID, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.service.Unlike(userID, commentID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"message": "いいねを取り消しました"})
}

// GetStatus はコメントのいいね状態（いいねしているか・いいね数）を返す。
func (h *CommentLikeHandler) GetStatus(c *gin.Context) {
	commentID, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	liked, count, err := h.service.GetStatus(userID, commentID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{
		"liked": liked,
		"count": count,
	})
}

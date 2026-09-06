package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// CommentLikeHandler はコメントへのいいね機能の HTTP ハンドラー。
// 各操作は 1 責務の usecase に委譲する。
type CommentLikeHandler struct {
	likeComment   *usecase.LikeCommentUseCase
	unlikeComment *usecase.UnlikeCommentUseCase
	getStatus     *usecase.GetCommentLikeStatusUseCase
}

// NewCommentLikeHandler は CommentLikeHandler を生成する。
func NewCommentLikeHandler(
	likeComment *usecase.LikeCommentUseCase,
	unlikeComment *usecase.UnlikeCommentUseCase,
	getStatus *usecase.GetCommentLikeStatusUseCase,
) *CommentLikeHandler {
	return &CommentLikeHandler{
		likeComment:   likeComment,
		unlikeComment: unlikeComment,
		getStatus:     getStatus,
	}
}

// Like はコメントにいいねする。
func (h *CommentLikeHandler) Like(c *gin.Context) {
	handleToggleAction(c, func(userID, id uint) error {
		return h.likeComment.Execute(c.Request.Context(), userID, id)
	}, "いいねしました")
}

// Unlike はコメントのいいねを取り消す。
func (h *CommentLikeHandler) Unlike(c *gin.Context) {
	handleToggleAction(c, func(userID, id uint) error {
		return h.unlikeComment.Execute(c.Request.Context(), userID, id)
	}, "いいねを取り消しました")
}

// likeStatusResponse はいいね状態レスポンス
type likeStatusResponse struct {
	Liked bool  `json:"liked"`
	Count int64 `json:"count"`
}

// GetStatus はコメントのいいね状態（いいねしているか・いいね数）を返す。
func (h *CommentLikeHandler) GetStatus(c *gin.Context) {
	commentID, ok := parseID(c, "id")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	liked, count, err := h.getStatus.Execute(c.Request.Context(), userID, commentID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, likeStatusResponse{Liked: liked, Count: count})
}

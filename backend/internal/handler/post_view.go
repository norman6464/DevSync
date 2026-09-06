package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// PostViewHandler は投稿閲覧数の HTTP ハンドラー。各操作は 1 責務の usecase に委譲する。
type PostViewHandler struct {
	recordView    *usecase.RecordPostViewUseCase
	getViewCount  *usecase.GetPostViewCountUseCase
	getMostViewed *usecase.GetMostViewedPostsUseCase
}

// NewPostViewHandler は PostViewHandler を生成する。
func NewPostViewHandler(
	recordView *usecase.RecordPostViewUseCase,
	getViewCount *usecase.GetPostViewCountUseCase,
	getMostViewed *usecase.GetMostViewedPostsUseCase,
) *PostViewHandler {
	return &PostViewHandler{
		recordView:    recordView,
		getViewCount:  getViewCount,
		getMostViewed: getMostViewed,
	}
}

// RecordView は投稿の閲覧を記録する。
func (h *PostViewHandler) RecordView(c *gin.Context) {
	postID, ok := parseID(c, "postId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.recordView.Execute(c.Request.Context(), userID, postID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("記録しました"))
}

// viewCountResponse は閲覧数レスポンス
type viewCountResponse struct {
	PostID    uint  `json:"post_id"`
	ViewCount int64 `json:"view_count"`
}

// GetViewCount は投稿の閲覧数を取得する。
func (h *PostViewHandler) GetViewCount(c *gin.Context) {
	postID, ok := parseID(c, "postId")
	if !ok {
		return
	}

	count, err := h.getViewCount.Execute(c.Request.Context(), postID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, viewCountResponse{PostID: postID, ViewCount: count})
}

// GetMostViewed は閲覧数の多い投稿ランキングを取得する。
func (h *PostViewHandler) GetMostViewed(c *gin.Context) {
	result, err := h.getMostViewed.Execute(c.Request.Context(), 20)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, result)
}

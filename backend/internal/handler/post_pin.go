package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// PostPinHandler は投稿ピン留めの HTTP ハンドラー。各操作は 1 責務の usecase に委譲する。
type PostPinHandler struct {
	pinPost      *usecase.PinPostUseCase
	unpinPost    *usecase.UnpinPostUseCase
	listPinned   *usecase.ListPinnedPostsUseCase
	countPinned  *usecase.CountPinnedPostsUseCase
	reorderPinned *usecase.ReorderPinnedPostsUseCase
}

// NewPostPinHandler は PostPinHandler を生成する。
func NewPostPinHandler(
	pinPost *usecase.PinPostUseCase,
	unpinPost *usecase.UnpinPostUseCase,
	listPinned *usecase.ListPinnedPostsUseCase,
	countPinned *usecase.CountPinnedPostsUseCase,
	reorderPinned *usecase.ReorderPinnedPostsUseCase,
) *PostPinHandler {
	return &PostPinHandler{
		pinPost:       pinPost,
		unpinPost:     unpinPost,
		listPinned:    listPinned,
		countPinned:   countPinned,
		reorderPinned: reorderPinned,
	}
}

// Pin は投稿をピン留めする。
func (h *PostPinHandler) Pin(c *gin.Context) {
	postID, ok := parseID(c, "postId")
	if !ok {
		return
	}
	userID := c.GetUint("userID")

	if err := h.pinPost.Execute(c.Request.Context(), userID, postID); err != nil {
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

	if err := h.unpinPost.Execute(c.Request.Context(), userID, postID); err != nil {
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

	pins, err := h.listPinned.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, dto.PinsResponse{Pins: pins})
}

// GetMyCount は認証ユーザー自身のピン留め投稿数を返す。
func (h *PostPinHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")
	count, err := h.countPinned.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, gin.H{"count": count})
}

// Reorder はピン留めの表示順序を変更する。
func (h *PostPinHandler) Reorder(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[dto.ReorderPinsRequest](c)
	if req == nil {
		return
	}

	if err := h.reorderPinned.Execute(c.Request.Context(), userID, req.PostIDs); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("順序を更新しました"))
}

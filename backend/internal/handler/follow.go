package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
)

// FollowServiceInterface はFollowServiceが実装すべきインターフェース。
type FollowServiceInterface interface {
	Follow(followerID, followeeID uint) error
	Unfollow(followerID, followeeID uint) error
	GetFollowers(userID uint) ([]model.User, error)
	GetFollowing(userID uint) ([]model.User, error)
}

// FollowHandler はフォロー関連のHTTPハンドラ。
// フォロー・アンフォロー・フォロワー/フォロー中一覧の取得を処理する。
type FollowHandler struct {
	service FollowServiceInterface
}

// NewFollowHandler は新しいFollowHandlerインスタンスを生成する。
func NewFollowHandler(s FollowServiceInterface) *FollowHandler {
	return &FollowHandler{service: s}
}

// Follow は指定ユーザーをフォローする。
func (h *FollowHandler) Follow(c *gin.Context) {
	userID := c.GetUint("userID")
	targetID, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.service.Follow(userID, targetID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("followed"))
}

// Unfollow は指定ユーザーのフォローを解除する。
func (h *FollowHandler) Unfollow(c *gin.Context) {
	userID := c.GetUint("userID")
	targetID, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.service.Unfollow(userID, targetID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("unfollowed"))
}

// GetFollowers は指定ユーザーのフォロワー一覧を返す。
func (h *FollowHandler) GetFollowers(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	users, err := h.service.GetFollowers(id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(users))
}

// GetFollowing は指定ユーザーのフォロー中一覧を返す。
func (h *FollowHandler) GetFollowing(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	users, err := h.service.GetFollowing(id)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, ensureSlice(users))
}

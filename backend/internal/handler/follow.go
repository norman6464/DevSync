package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// FollowHandler はフォロー関連の HTTP ハンドラ。
// 各操作は 1 責務の usecase に委譲する（フォロー / 解除 / フォロワー一覧 / フォロー中一覧）。
type FollowHandler struct {
	followUser    *usecase.FollowUserUseCase
	unfollowUser  *usecase.UnfollowUserUseCase
	listFollowers *usecase.ListFollowersUseCase
	listFollowing *usecase.ListFollowingUseCase
}

// NewFollowHandler は FollowHandler を生成する。
func NewFollowHandler(
	followUser *usecase.FollowUserUseCase,
	unfollowUser *usecase.UnfollowUserUseCase,
	listFollowers *usecase.ListFollowersUseCase,
	listFollowing *usecase.ListFollowingUseCase,
) *FollowHandler {
	return &FollowHandler{
		followUser:    followUser,
		unfollowUser:  unfollowUser,
		listFollowers: listFollowers,
		listFollowing: listFollowing,
	}
}

// Follow は指定ユーザーをフォローする。
func (h *FollowHandler) Follow(c *gin.Context) {
	userID := c.GetUint("userID")
	targetID, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.followUser.Execute(c.Request.Context(), userID, targetID); err != nil {
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
	if err := h.unfollowUser.Execute(c.Request.Context(), userID, targetID); err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, domain.NewMessageResponse("unfollowed"))
}

// followListResponse はフォロワー/フォロー中一覧レスポンス（ページネーション付き）。
type followListResponse struct {
	Users  []model.User `json:"users"`
	Total  int64        `json:"total"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

// GetFollowers は指定ユーザーのフォロワー一覧をページネーション付きで返す。
func (h *FollowHandler) GetFollowers(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	limit, offset := parseLimitOffset(c)
	users, total, err := h.listFollowers.Execute(c.Request.Context(), id, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, followListResponse{
		Users:  ensureSlice(users),
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// GetFollowing は指定ユーザーのフォロー中一覧をページネーション付きで返す。
func (h *FollowHandler) GetFollowing(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	limit, offset := parseLimitOffset(c)
	users, total, err := h.listFollowing.Execute(c.Request.Context(), id, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, followListResponse{
		Users:  ensureSlice(users),
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

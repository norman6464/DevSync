package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// FollowHandler はフォロー関連のHTTPハンドラ。
// フォロー・アンフォロー・フォロワー/フォロー中一覧の取得を処理する。
type FollowHandler struct {
	service *service.FollowService
}

// NewFollowHandler は新しいFollowHandlerインスタンスを生成する。
func NewFollowHandler(s *service.FollowService) *FollowHandler {
	return &FollowHandler{service: s}
}

// Follow は指定ユーザーをフォローする。
func (h *FollowHandler) Follow(c *gin.Context) {
	userID := c.GetUint("userID")
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.Follow(userID, uint(targetID)); err != nil {
		if errors.Is(err, service.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot follow yourself"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "followed"})
}

// Unfollow は指定ユーザーのフォローを解除する。
func (h *FollowHandler) Unfollow(c *gin.Context) {
	userID := c.GetUint("userID")
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.Unfollow(userID, uint(targetID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "unfollowed"})
}

// GetFollowers は指定ユーザーのフォロワー一覧を返す。
func (h *FollowHandler) GetFollowers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	users, err := h.service.GetFollowers(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if users == nil {
		users = []model.User{}
	}
	c.JSON(http.StatusOK, users)
}

// GetFollowing は指定ユーザーのフォロー中一覧を返す。
func (h *FollowHandler) GetFollowing(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	users, err := h.service.GetFollowing(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if users == nil {
		users = []model.User{}
	}
	c.JSON(http.StatusOK, users)
}

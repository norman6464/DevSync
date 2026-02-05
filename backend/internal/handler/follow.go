package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

type FollowHandler struct {
	repo *repository.FollowRepository
}

func NewFollowHandler(repo *repository.FollowRepository) *FollowHandler {
	return &FollowHandler{repo: repo}
}

func (h *FollowHandler) Follow(c *gin.Context) {
	userID := c.GetUint("userID")
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if userID == uint(targetID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot follow yourself"})
		return
	}
	if err := h.repo.Follow(userID, uint(targetID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "followed"})
}

func (h *FollowHandler) Unfollow(c *gin.Context) {
	userID := c.GetUint("userID")
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.repo.Unfollow(userID, uint(targetID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "unfollowed"})
}

func (h *FollowHandler) GetFollowers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	users, err := h.repo.GetFollowers(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if users == nil {
		users = []model.User{}
	}
	c.JSON(http.StatusOK, users)
}

func (h *FollowHandler) GetFollowing(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	users, err := h.repo.GetFollowing(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if users == nil {
		users = []model.User{}
	}
	c.JSON(http.StatusOK, users)
}

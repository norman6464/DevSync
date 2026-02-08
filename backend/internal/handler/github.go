package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/service"
)

type GitHubHandler struct {
	githubService *service.GitHubService
	authService   *service.AuthService
}

func NewGitHubHandler(
	githubService *service.GitHubService,
	authService *service.AuthService,
) *GitHubHandler {
	return &GitHubHandler{
		githubService: githubService,
		authService:   authService,
	}
}

func (h *GitHubHandler) Connect(c *gin.Context) {
	userID := c.GetUint("userID")
	state, err := h.authService.GenerateOAuthState(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})
		return
	}
	url := h.githubService.GetOAuthURL(state)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *GitHubHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code or state"})
		return
	}

	userID, err := h.authService.ValidateOAuthState(state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}

	if err := h.githubService.ConnectGitHub(userID, code, state); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect github"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "github connected"})
}

func (h *GitHubHandler) GetContributions(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	contributions, err := h.githubService.GetContributions(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, contributions)
}

func (h *GitHubHandler) GetLanguages(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	stats, err := h.githubService.GetLanguages(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *GitHubHandler) GetRepos(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	repos, err := h.githubService.GetRepos(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, repos)
}

func (h *GitHubHandler) Sync(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.githubService.SyncUserData(userID); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sync complete"})
}

func (h *GitHubHandler) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")

	if err := h.githubService.DisconnectGitHub(userID); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "github disconnected"})
}

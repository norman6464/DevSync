package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/repository"
	"github.com/norman6464/devsync/backend/internal/service"
)

type GitHubHandler struct {
	githubService *service.GitHubService
	authService   *service.AuthService
	userRepo      *repository.UserRepository
	githubRepo    *repository.GitHubRepository
}

func NewGitHubHandler(
	githubService *service.GitHubService,
	authService *service.AuthService,
	userRepo *repository.UserRepository,
	githubRepo *repository.GitHubRepository,
) *GitHubHandler {
	return &GitHubHandler{
		githubService: githubService,
		authService:   authService,
		userRepo:      userRepo,
		githubRepo:    githubRepo,
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

	accessToken, err := h.githubService.ExchangeCode(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange code"})
		return
	}

	username, avatarURL, err := h.githubService.GetGitHubUser(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get github user"})
		return
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.GitHubToken = accessToken
	user.GitHubUsername = username
	user.GitHubConnected = true
	if user.AvatarURL == "" {
		user.AvatarURL = avatarURL
	}

	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	// Sync data in the background-ish (synchronous for MVP)
	go h.githubService.SyncData(user)

	c.JSON(http.StatusOK, gin.H{"message": "github connected"})
}

func (h *GitHubHandler) GetContributions(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	contributions, err := h.githubRepo.GetContributions(uint(userID))
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

	stats, err := h.githubRepo.GetLanguageStats(uint(userID))
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

	repos, err := h.githubRepo.GetRepos(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, repos)
}

func (h *GitHubHandler) Sync(c *gin.Context) {
	userID := c.GetUint("userID")
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := h.githubService.SyncData(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sync complete"})
}

func (h *GitHubHandler) Disconnect(c *gin.Context) {
	userID := c.GetUint("userID")
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.GitHubToken = ""
	user.GitHubUsername = ""
	user.GitHubConnected = false

	if err := h.userRepo.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.githubRepo.DeleteUserData(userID)

	c.JSON(http.StatusOK, gin.H{"message": "github disconnected"})
}

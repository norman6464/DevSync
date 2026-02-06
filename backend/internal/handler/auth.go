package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
	"github.com/norman6464/devsync/backend/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authService       *service.AuthService
	githubService     *service.GitHubService
	userRepo          *repository.UserRepository
	passwordResetRepo *repository.PasswordResetRepository
}

func NewAuthHandler(authService *service.AuthService, githubService *service.GitHubService, userRepo *repository.UserRepository, passwordResetRepo *repository.PasswordResetRepository) *AuthHandler {
	return &AuthHandler{
		authService:       authService,
		githubService:     githubService,
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input service.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Register(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input service.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) GitHubLogin(c *gin.Context) {
	state, err := h.authService.GenerateLoginState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})
		return
	}
	url := h.githubService.GetLoginOAuthURL(state)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *AuthHandler) GitHubLoginCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code or state"})
		return
	}

	if err := h.authService.ValidateLoginState(state); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}

	accessToken, err := h.githubService.ExchangeCode(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange code"})
		return
	}

	ghUser, err := h.githubService.GetGitHubUser(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get github user"})
		return
	}

	resp, err := h.authService.GitHubLogin(ghUser, accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Sync GitHub data in background
	user, _ := h.userRepo.FindByID(resp.User.ID)
	if user != nil && user.GitHubConnected {
		go h.githubService.SyncData(user)
	}

	c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetUint("userID")
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// RequestPasswordReset generates a password reset token and returns it
// In production, this should send an email instead of returning the token
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.FindByEmail(input.Email)
	if err != nil {
		// Don't reveal if email exists or not for security
		c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset link has been sent"})
		return
	}

	// Invalidate any existing tokens for this user
	h.passwordResetRepo.InvalidateUserTokens(user.ID)

	// Generate secure random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	token := hex.EncodeToString(tokenBytes)

	// Create reset token (expires in 1 hour)
	resetToken := &model.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	if err := h.passwordResetRepo.Create(resetToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reset token"})
		return
	}

	// In a real app, send email here
	// For demo purposes, return the token (remove in production!)
	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset token generated",
		"token":   token, // Remove this in production - should be sent via email
	})
}

// ResetPassword resets the password using a valid token
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var input struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find and validate token
	resetToken, err := h.passwordResetRepo.FindByToken(input.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}

	if !resetToken.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired token"})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Update password
	if err := h.userRepo.UpdatePassword(resetToken.UserID, string(hashedPassword)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	// Mark token as used
	h.passwordResetRepo.MarkAsUsed(resetToken.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
}

// DeleteAccount permanently deletes the user's account and all related data
func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	userID := c.GetUint("userID")

	var input struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify user exists
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// If user has a password (not GitHub-only), verify it
	if user.Password != "" {
		if input.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password required"})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password"})
			return
		}
	}

	// Delete user and all related data
	if err := h.userRepo.DeleteWithRelatedData(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

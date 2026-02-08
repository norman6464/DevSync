package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	q := c.Query("q")
	users, err := h.service.GetAll(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID := c.GetUint("userID")
	if userID != uint(id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot update other user's profile"})
		return
	}

	existing, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var input struct {
		Name                string  `json:"name"`
		Bio                 string  `json:"bio"`
		AvatarURL           string  `json:"avatar_url"`
		SkillsLanguages     *string `json:"skills_languages"`
		SkillsFrameworks    *string `json:"skills_frameworks"`
		OnboardingCompleted *bool   `json:"onboarding_completed"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" {
		existing.Name = input.Name
	}
	existing.Bio = input.Bio
	existing.AvatarURL = input.AvatarURL
	if input.SkillsLanguages != nil {
		existing.SkillsLanguages = *input.SkillsLanguages
	}
	if input.SkillsFrameworks != nil {
		existing.SkillsFrameworks = *input.SkillsFrameworks
	}
	if input.OnboardingCompleted != nil {
		existing.OnboardingCompleted = *input.OnboardingCompleted
	}

	if err := h.service.Update(existing); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, existing)
}

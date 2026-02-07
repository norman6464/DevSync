package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/repository"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	q := c.Query("q")
	if q != "" {
		users, err := h.repo.Search(q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
		return
	}

	users, err := h.repo.FindAll()
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

	user, err := h.repo.FindByID(uint(id))
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

	existing, err := h.repo.FindByID(uint(id))
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

	if err := h.repo.Update(existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, existing)
}

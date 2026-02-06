package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

type LearningResourceHandler struct {
	repo *repository.LearningResourceRepository
}

func NewLearningResourceHandler(repo *repository.LearningResourceRepository) *LearningResourceHandler {
	return &LearningResourceHandler{repo: repo}
}

type CreateResourceRequest struct {
	Title       string `json:"title" binding:"required,max=300"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Category    string `json:"category" binding:"required"`
	Difficulty  string `json:"difficulty"`
	Tags        string `json:"tags"`
	ImageURL    string `json:"image_url"`
	IsPublic    *bool  `json:"is_public"`
}

type UpdateResourceRequest struct {
	Title       string `json:"title" binding:"max=300"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Category    string `json:"category"`
	Difficulty  string `json:"difficulty"`
	Tags        string `json:"tags"`
	ImageURL    string `json:"image_url"`
	IsPublic    *bool  `json:"is_public"`
}

func (h *LearningResourceHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	resource := &model.LearningResource{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		URL:         req.URL,
		Category:    model.ResourceCategory(req.Category),
		Difficulty:  model.ResourceDifficulty(req.Difficulty),
		Tags:        req.Tags,
		ImageURL:    req.ImageURL,
		IsPublic:    isPublic,
	}

	if err := h.repo.Create(resource); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create resource"})
		return
	}

	c.JSON(http.StatusCreated, resource)
}

func (h *LearningResourceHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	resource, err := h.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	// Check if resource is private and user is not owner
	userID := c.GetUint("userID")
	if !resource.IsPublic && resource.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to view this resource"})
		return
	}

	// Check if current user has liked/saved
	hasLiked, _ := h.repo.HasLiked(userID, uint(id))
	hasSaved, _ := h.repo.HasSaved(userID, uint(id))

	c.JSON(http.StatusOK, gin.H{
		"resource":  resource,
		"has_liked": hasLiked,
		"has_saved": hasSaved,
	})
}

func (h *LearningResourceHandler) GetByUserID(c *gin.Context) {
	targetUserID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	currentUserID := c.GetUint("userID")
	includePrivate := currentUserID == uint(targetUserID)

	resources, err := h.repo.FindByUserID(uint(targetUserID), includePrivate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resources"})
		return
	}

	c.JSON(http.StatusOK, resources)
}

func (h *LearningResourceHandler) GetPublic(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	category := c.Query("category")
	difficulty := c.Query("difficulty")

	if limit > 100 {
		limit = 100
	}

	resources, total, err := h.repo.FindPublic(limit, offset, category, difficulty)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resources"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resources": resources,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

func (h *LearningResourceHandler) Search(c *gin.Context) {
	query := c.Query("q")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	resources, total, err := h.repo.Search(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search resources"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resources": resources,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

func (h *LearningResourceHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	resource, err := h.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	if resource.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this resource"})
		return
	}

	var req UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != "" {
		resource.Title = req.Title
	}
	if req.Description != "" {
		resource.Description = req.Description
	}
	if req.URL != "" {
		resource.URL = req.URL
	}
	if req.Category != "" {
		resource.Category = model.ResourceCategory(req.Category)
	}
	if req.Difficulty != "" {
		resource.Difficulty = model.ResourceDifficulty(req.Difficulty)
	}
	if req.Tags != "" {
		resource.Tags = req.Tags
	}
	if req.ImageURL != "" {
		resource.ImageURL = req.ImageURL
	}
	if req.IsPublic != nil {
		resource.IsPublic = *req.IsPublic
	}

	if err := h.repo.Update(resource); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update resource"})
		return
	}

	c.JSON(http.StatusOK, resource)
}

func (h *LearningResourceHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	resource, err := h.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	if resource.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this resource"})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete resource"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
}

func (h *LearningResourceHandler) Like(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	if err := h.repo.Like(userID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like resource"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource liked"})
}

func (h *LearningResourceHandler) Unlike(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	if err := h.repo.Unlike(userID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike resource"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource unliked"})
}

func (h *LearningResourceHandler) SaveResource(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	if err := h.repo.Save(userID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save resource"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource saved"})
}

func (h *LearningResourceHandler) UnsaveResource(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resource ID"})
		return
	}

	if err := h.repo.Unsave(userID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsave resource"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource unsaved"})
}

func (h *LearningResourceHandler) GetSaved(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	resources, total, err := h.repo.FindSavedByUserID(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch saved resources"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resources": resources,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

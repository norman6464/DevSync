package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

type LearningResourceHandler struct {
	service *service.LearningResourceService
}

func NewLearningResourceHandler(s *service.LearningResourceService) *LearningResourceHandler {
	return &LearningResourceHandler{service: s}
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

	if err := h.service.Create(resource); err != nil {
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

	userID := c.GetUint("userID")

	resource, err := h.service.GetByID(uint(id), userID)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to view this resource"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	// Check if current user has liked/saved
	hasLiked, _ := h.service.HasLiked(userID, uint(id))
	hasSaved, _ := h.service.HasSaved(userID, uint(id))

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

	resources, err := h.service.GetByUserID(uint(targetUserID), currentUserID)
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

	resources, total, err := h.service.GetPublic(limit, offset, category, difficulty)
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

	resources, total, err := h.service.Search(query, limit, offset)
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

	var req UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := &model.LearningResource{}
	if req.Title != "" {
		updates.Title = req.Title
	}
	if req.Description != "" {
		updates.Description = req.Description
	}
	if req.URL != "" {
		updates.URL = req.URL
	}
	if req.Category != "" {
		updates.Category = model.ResourceCategory(req.Category)
	}
	if req.Difficulty != "" {
		updates.Difficulty = model.ResourceDifficulty(req.Difficulty)
	}
	if req.Tags != "" {
		updates.Tags = req.Tags
	}
	if req.ImageURL != "" {
		updates.ImageURL = req.ImageURL
	}

	resource, err := h.service.Update(uint(id), userID, updates)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this resource"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
		return
	}

	// Handle IsPublic separately if provided
	if req.IsPublic != nil {
		resource, err = h.service.UpdateVisibility(uint(id), userID, *req.IsPublic)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update resource"})
			return
		}
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

	if err := h.service.Delete(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this resource"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
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

	if err := h.service.Like(userID, uint(id)); err != nil {
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

	if err := h.service.Unlike(userID, uint(id)); err != nil {
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

	if err := h.service.Save(userID, uint(id)); err != nil {
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

	if err := h.service.Unsave(userID, uint(id)); err != nil {
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

	resources, total, err := h.service.GetSavedByUserID(userID, limit, offset)
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

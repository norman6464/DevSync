package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

type RoadmapHandler struct {
	roadmapRepo *repository.RoadmapRepository
}

func NewRoadmapHandler(roadmapRepo *repository.RoadmapRepository) *RoadmapHandler {
	return &RoadmapHandler{roadmapRepo: roadmapRepo}
}

// === Roadmap Endpoints ===

// Create creates a new roadmap
func (h *RoadmapHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Category    string `json:"category"`
		IsPublic    bool   `json:"is_public"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	roadmap := &model.Roadmap{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Category:    model.RoadmapCategory(req.Category),
		IsPublic:    req.IsPublic,
		Status:      model.RoadmapStatusActive,
	}

	if req.Category == "" {
		roadmap.Category = model.RoadmapCategoryOther
	}

	if err := h.roadmapRepo.Create(roadmap); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create roadmap"})
		return
	}

	c.JSON(http.StatusCreated, roadmap)
}

// GetMyRoadmaps gets all roadmaps for the current user
func (h *RoadmapHandler) GetMyRoadmaps(c *gin.Context) {
	userID := c.GetUint("userID")

	roadmaps, err := h.roadmapRepo.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get roadmaps"})
		return
	}

	c.JSON(http.StatusOK, roadmaps)
}

// GetPublicRoadmaps gets all public roadmaps
func (h *RoadmapHandler) GetPublicRoadmaps(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	roadmaps, total, err := h.roadmapRepo.GetPublicRoadmaps(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get roadmaps"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roadmaps": roadmaps,
		"total":    total,
	})
}

// GetByID gets a roadmap by ID (with steps)
func (h *RoadmapHandler) GetByID(c *gin.Context) {
	roadmapID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid roadmap ID"})
		return
	}

	roadmap, err := h.roadmapRepo.FindByID(uint(roadmapID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "roadmap not found"})
		return
	}

	userID := c.GetUint("userID")
	if roadmap.UserID != userID && !roadmap.IsPublic {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	c.JSON(http.StatusOK, roadmap)
}

// Update updates a roadmap
func (h *RoadmapHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid roadmap ID"})
		return
	}

	roadmap, err := h.roadmapRepo.FindByID(uint(roadmapID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "roadmap not found"})
		return
	}

	if roadmap.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Category    *string `json:"category"`
		IsPublic    *bool   `json:"is_public"`
		Status      *string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.Title != nil {
		roadmap.Title = *req.Title
	}
	if req.Description != nil {
		roadmap.Description = *req.Description
	}
	if req.Category != nil {
		roadmap.Category = model.RoadmapCategory(*req.Category)
	}
	if req.IsPublic != nil {
		roadmap.IsPublic = *req.IsPublic
	}
	if req.Status != nil {
		roadmap.Status = model.RoadmapStatus(*req.Status)
		if roadmap.Status == model.RoadmapStatusCompleted && roadmap.CompletedAt == nil {
			now := time.Now()
			roadmap.CompletedAt = &now
		}
	}

	if err := h.roadmapRepo.Update(roadmap); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update roadmap"})
		return
	}

	c.JSON(http.StatusOK, roadmap)
}

// Delete deletes a roadmap
func (h *RoadmapHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid roadmap ID"})
		return
	}

	roadmap, err := h.roadmapRepo.FindByID(uint(roadmapID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "roadmap not found"})
		return
	}

	if roadmap.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	if err := h.roadmapRepo.Delete(uint(roadmapID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete roadmap"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "roadmap deleted"})
}

// CopyRoadmap copies a public roadmap as a template
func (h *RoadmapHandler) CopyRoadmap(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid roadmap ID"})
		return
	}

	original, err := h.roadmapRepo.FindByID(uint(roadmapID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "roadmap not found"})
		return
	}

	if !original.IsPublic && original.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	copied, err := h.roadmapRepo.CopyRoadmap(uint(roadmapID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to copy roadmap"})
		return
	}

	c.JSON(http.StatusCreated, copied)
}

// === RoadmapStep Endpoints ===

// CreateStep creates a new step
func (h *RoadmapHandler) CreateStep(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid roadmap ID"})
		return
	}

	roadmap, err := h.roadmapRepo.FindByID(uint(roadmapID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "roadmap not found"})
		return
	}

	if roadmap.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		ResourceURL string `json:"resource_url"`
		OrderIndex  *int   `json:"order_index"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	orderIndex := roadmap.StepCount
	if req.OrderIndex != nil {
		orderIndex = *req.OrderIndex
	}

	step := &model.RoadmapStep{
		RoadmapID:   uint(roadmapID),
		Title:       req.Title,
		Description: req.Description,
		ResourceURL: req.ResourceURL,
		OrderIndex:  orderIndex,
	}

	if err := h.roadmapRepo.CreateStep(step); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create step"})
		return
	}

	c.JSON(http.StatusCreated, step)
}

// UpdateStep updates a step
func (h *RoadmapHandler) UpdateStep(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid roadmap ID"})
		return
	}
	stepID, err := strconv.ParseUint(c.Param("stepId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step ID"})
		return
	}

	roadmap, err := h.roadmapRepo.FindByID(uint(roadmapID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "roadmap not found"})
		return
	}

	if roadmap.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	step, err := h.roadmapRepo.FindStepByID(uint(stepID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "step not found"})
		return
	}

	if step.RoadmapID != uint(roadmapID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "step does not belong to roadmap"})
		return
	}

	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		ResourceURL *string `json:"resource_url"`
		IsCompleted *bool   `json:"is_completed"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.Title != nil {
		step.Title = *req.Title
	}
	if req.Description != nil {
		step.Description = *req.Description
	}
	if req.ResourceURL != nil {
		step.ResourceURL = *req.ResourceURL
	}
	if req.IsCompleted != nil {
		step.IsCompleted = *req.IsCompleted
		if step.IsCompleted && step.CompletedAt == nil {
			now := time.Now()
			step.CompletedAt = &now
		} else if !step.IsCompleted {
			step.CompletedAt = nil
		}
	}

	if err := h.roadmapRepo.UpdateStep(step); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update step"})
		return
	}

	c.JSON(http.StatusOK, step)
}

// DeleteStep deletes a step
func (h *RoadmapHandler) DeleteStep(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid roadmap ID"})
		return
	}
	stepID, err := strconv.ParseUint(c.Param("stepId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid step ID"})
		return
	}

	roadmap, err := h.roadmapRepo.FindByID(uint(roadmapID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "roadmap not found"})
		return
	}

	if roadmap.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	step, err := h.roadmapRepo.FindStepByID(uint(stepID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "step not found"})
		return
	}

	if step.RoadmapID != uint(roadmapID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "step does not belong to roadmap"})
		return
	}

	if err := h.roadmapRepo.DeleteStep(uint(stepID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete step"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "step deleted"})
}

// ReorderSteps reorders steps within a roadmap
func (h *RoadmapHandler) ReorderSteps(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid roadmap ID"})
		return
	}

	roadmap, err := h.roadmapRepo.FindByID(uint(roadmapID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "roadmap not found"})
		return
	}

	if roadmap.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
		return
	}

	var req struct {
		Orders []repository.StepOrder `json:"orders" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.roadmapRepo.ReorderSteps(uint(roadmapID), req.Orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reorder steps"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "steps reordered"})
}

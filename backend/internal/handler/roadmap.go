package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

type RoadmapHandler struct {
	service *service.RoadmapService
}

func NewRoadmapHandler(s *service.RoadmapService) *RoadmapHandler {
	return &RoadmapHandler{service: s}
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

	if err := h.service.Create(roadmap); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create roadmap"})
		return
	}

	c.JSON(http.StatusCreated, roadmap)
}

// GetMyRoadmaps gets all roadmaps for the current user
func (h *RoadmapHandler) GetMyRoadmaps(c *gin.Context) {
	userID := c.GetUint("userID")

	roadmaps, err := h.service.GetByUserID(userID)
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

	roadmaps, total, err := h.service.GetPublicRoadmaps(limit, offset)
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

	userID := c.GetUint("userID")

	roadmap, err := h.service.GetByID(uint(roadmapID), userID)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "roadmap not found"})
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

	updates := &model.Roadmap{}
	if req.Title != nil {
		updates.Title = *req.Title
	}
	if req.Description != nil {
		updates.Description = *req.Description
	}
	if req.Category != nil {
		updates.Category = model.RoadmapCategory(*req.Category)
	}
	if req.Status != nil {
		updates.Status = model.RoadmapStatus(*req.Status)
	}

	roadmap, err := h.service.Update(uint(roadmapID), userID, updates)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "roadmap not found"})
		return
	}

	// Handle IsPublic separately if provided
	if req.IsPublic != nil {
		roadmap, err = h.service.UpdateVisibility(uint(roadmapID), userID, *req.IsPublic)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update roadmap"})
			return
		}
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

	if err := h.service.Delete(uint(roadmapID), userID); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "roadmap not found"})
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

	copied, err := h.service.CopyRoadmap(uint(roadmapID), userID)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "roadmap not found"})
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

	step := &model.RoadmapStep{
		Title:       req.Title,
		Description: req.Description,
		ResourceURL: req.ResourceURL,
	}
	if req.OrderIndex != nil {
		step.OrderIndex = *req.OrderIndex
	}

	if err := h.service.CreateStep(uint(roadmapID), userID, step); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
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

	// Handle completion status change separately
	if req.IsCompleted != nil {
		step, err := h.service.UpdateStepCompletion(uint(roadmapID), uint(stepID), userID, *req.IsCompleted)
		if err != nil {
			if errors.Is(err, service.ErrForbidden) {
				c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
				return
			}
			if errors.Is(err, service.ErrBadRequest) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "step does not belong to roadmap"})
				return
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "step not found"})
			return
		}
		// If only completion was updated, return early
		if req.Title == nil && req.Description == nil && req.ResourceURL == nil {
			c.JSON(http.StatusOK, step)
			return
		}
	}

	updates := &model.RoadmapStep{}
	if req.Title != nil {
		updates.Title = *req.Title
	}
	if req.Description != nil {
		updates.Description = *req.Description
	}
	if req.ResourceURL != nil {
		updates.ResourceURL = *req.ResourceURL
	}

	step, err := h.service.UpdateStep(uint(roadmapID), uint(stepID), userID, updates)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		if errors.Is(err, service.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "step does not belong to roadmap"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "step not found"})
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

	if err := h.service.DeleteStep(uint(roadmapID), uint(stepID), userID); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		if errors.Is(err, service.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "step does not belong to roadmap"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "step not found"})
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

	var req struct {
		Orders []service.StepOrder `json:"orders" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.service.ReorderSteps(uint(roadmapID), userID, req.Orders); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reorder steps"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "steps reordered"})
}

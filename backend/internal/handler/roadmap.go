package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// RoadmapHandler はロードマップ関連のHTTPハンドラ。
// ロードマップとステップのCRUD・公開一覧・コピー・並べ替えを処理する。
type RoadmapHandler struct {
	service *service.RoadmapService
}

// NewRoadmapHandler は新しいRoadmapHandlerインスタンスを生成する。
func NewRoadmapHandler(s *service.RoadmapService) *RoadmapHandler {
	return &RoadmapHandler{service: s}
}

// === ロードマップエンドポイント ===

// Create は新しいロードマップを作成する。
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
		respondError(c, err)
		return
	}

	respondCreated(c, roadmap)
}

// GetMyRoadmaps は現在のユーザーのロードマップ一覧を取得する。
func (h *RoadmapHandler) GetMyRoadmaps(c *gin.Context) {
	userID := c.GetUint("userID")

	roadmaps, err := h.service.GetByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, roadmaps)
}

// GetPublicRoadmaps は公開ロードマップの一覧をページネーション付きで取得する。
func (h *RoadmapHandler) GetPublicRoadmaps(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	roadmaps, total, err := h.service.GetPublicRoadmaps(limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{
		"roadmaps": roadmaps,
		"total":    total,
	})
}

// GetByID は指定IDのロードマップをステップ付きで取得する。
func (h *RoadmapHandler) GetByID(c *gin.Context) {
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}

	userID := c.GetUint("userID")

	roadmap, err := h.service.GetByID(roadmapID, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, roadmap)
}

// Update は指定IDのロードマップを更新する。
func (h *RoadmapHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
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

	roadmap, err := h.service.Update(roadmapID, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	// IsPublicが指定されている場合は別途処理する
	if req.IsPublic != nil {
		roadmap, err = h.service.UpdateVisibility(roadmapID, userID, *req.IsPublic)
		if err != nil {
			respondError(c, err)
			return
		}
	}

	respondOK(c, roadmap)
}

// Delete は指定IDのロードマップを削除する。
func (h *RoadmapHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(roadmapID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

// CopyRoadmap は公開ロードマップをテンプレートとしてコピーする。
func (h *RoadmapHandler) CopyRoadmap(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}

	copied, err := h.service.CopyRoadmap(roadmapID, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, copied)
}

// GetTemplates はテンプレートロードマップの一覧を取得する。
func (h *RoadmapHandler) GetTemplates(c *gin.Context) {
	templates, err := h.service.GetTemplates()
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, templates)
}

// CreateFromTemplate はテンプレートからユーザー用ロードマップを作成する。
func (h *RoadmapHandler) CreateFromTemplate(c *gin.Context) {
	userID := c.GetUint("userID")
	templateID, ok := parseID(c, "id")
	if !ok {
		return
	}

	roadmap, err := h.service.CreateFromTemplate(templateID, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, roadmap)
}

// === ロードマップステップエンドポイント ===

// CreateStep はロードマップに新しいステップを作成する。
func (h *RoadmapHandler) CreateStep(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
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

	if err := h.service.CreateStep(roadmapID, userID, step); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, step)
}

// UpdateStep はロードマップのステップを更新する。
func (h *RoadmapHandler) UpdateStep(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}
	stepID, ok := parseID(c, "stepId")
	if !ok {
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

	// 完了ステータスの変更を別途処理する
	if req.IsCompleted != nil {
		step, err := h.service.UpdateStepCompletion(roadmapID, stepID, userID, *req.IsCompleted)
		if err != nil {
			respondError(c, err)
			return
		}
		// 完了ステータスのみの更新の場合は早期リターンする
		if req.Title == nil && req.Description == nil && req.ResourceURL == nil {
			respondOK(c, step)
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

	step, err := h.service.UpdateStep(roadmapID, stepID, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, step)
}

// DeleteStep はロードマップのステップを削除する。
func (h *RoadmapHandler) DeleteStep(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}
	stepID, ok := parseID(c, "stepId")
	if !ok {
		return
	}

	if err := h.service.DeleteStep(roadmapID, stepID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

// ReorderSteps はロードマップ内のステップの並び順を変更する。
func (h *RoadmapHandler) ReorderSteps(c *gin.Context) {
	userID := c.GetUint("userID")
	roadmapID, ok := parseID(c, "id")
	if !ok {
		return
	}

	var req struct {
		Orders []service.StepOrder `json:"orders" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.service.ReorderSteps(roadmapID, userID, req.Orders); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"message": "steps reordered"})
}

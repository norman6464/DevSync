package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// LearningResourceHandler は学習リソース関連のHTTPハンドラ。
// 学習リソースのCRUD・検索・いいね・保存を処理する。
type LearningResourceHandler struct {
	service *service.LearningResourceService
}

// NewLearningResourceHandler は新しいLearningResourceHandlerインスタンスを生成する。
func NewLearningResourceHandler(s *service.LearningResourceService) *LearningResourceHandler {
	return &LearningResourceHandler{service: s}
}

// CreateResourceRequest は学習リソース作成のリクエストボディ。
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

// UpdateResourceRequest は学習リソース更新のリクエストボディ。
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

// Create は新しい学習リソースを作成する。
func (h *LearningResourceHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 公開設定のデフォルト値はtrue
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

// GetByID は指定されたIDの学習リソースを取得する。
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

	// 現在のユーザーがいいね・保存済みかを確認
	hasLiked, _ := h.service.HasLiked(userID, uint(id))
	hasSaved, _ := h.service.HasSaved(userID, uint(id))

	c.JSON(http.StatusOK, gin.H{
		"resource":  resource,
		"has_liked": hasLiked,
		"has_saved": hasSaved,
	})
}

// GetByUserID は指定されたユーザーの学習リソース一覧を取得する。
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

// GetPublic は公開学習リソース一覧をページネーション・フィルター付きで取得する。
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

// Search はキーワードで学習リソースを検索する。
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

// Update は指定された学習リソースを更新する。
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

	// 公開設定が指定されている場合は別途更新
	if req.IsPublic != nil {
		resource, err = h.service.UpdateVisibility(uint(id), userID, *req.IsPublic)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update resource"})
			return
		}
	}

	c.JSON(http.StatusOK, resource)
}

// Delete は指定された学習リソースを削除する。
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

// Like は学習リソースにいいねする。
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

// Unlike は学習リソースのいいねを取り消す。
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

// SaveResource は学習リソースを保存する。
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

// UnsaveResource は学習リソースの保存を取り消す。
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

// GetSaved は認証ユーザーの保存済み学習リソース一覧を取得する。
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

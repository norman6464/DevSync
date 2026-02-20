package handler

import (

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// LearningResourceServiceInterface は学習リソースサービスの抽象インターフェース。
type LearningResourceServiceInterface interface {
	Create(resource *model.LearningResource) error
	GetByID(id, userID uint) (*model.LearningResource, error)
	HasLiked(userID, resourceID uint) (bool, error)
	HasSaved(userID, resourceID uint) (bool, error)
	GetByUserID(targetUserID, currentUserID uint) ([]model.LearningResource, error)
	GetPublic(limit, offset int, category, difficulty string) ([]model.LearningResource, int64, error)
	Search(query string, limit, offset int) ([]model.LearningResource, int64, error)
	Update(id, userID uint, updates *model.LearningResource) (*model.LearningResource, error)
	UpdateVisibility(id, userID uint, isPublic bool) (*model.LearningResource, error)
	Delete(id, userID uint) error
	Like(userID, resourceID uint) error
	Unlike(userID, resourceID uint) error
	Save(userID, resourceID uint) error
	Unsave(userID, resourceID uint) error
	GetSavedByUserID(userID uint, limit, offset int) ([]model.LearningResource, int64, error)
	GetByDifficulty(difficulty string, limit, offset int) ([]model.LearningResource, int64, error)
}

// LearningResourceHandler は学習リソース関連のHTTPハンドラ。
// 学習リソースのCRUD・検索・いいね・保存を処理する。
type LearningResourceHandler struct {
	service LearningResourceServiceInterface
}

// NewLearningResourceHandler は新しいLearningResourceHandlerインスタンスを生成する。
func NewLearningResourceHandler(s LearningResourceServiceInterface) *LearningResourceHandler {
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

	req := bindJSON[CreateResourceRequest](c)
	if req == nil {
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
		respondError(c, err)
		return
	}

	respondCreated(c, resource)
}

// GetByID は指定されたIDの学習リソースを取得する。
func (h *LearningResourceHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	userID := c.GetUint("userID")

	resource, err := h.service.GetByID(id, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	// 現在のユーザーがいいね・保存済みかを確認
	hasLiked, _ := h.service.HasLiked(userID, id)
	hasSaved, _ := h.service.HasSaved(userID, id)

	respondOK(c, dto.ResourceDetailResponse{
		Resource: *resource,
		HasLiked: hasLiked,
		HasSaved: hasSaved,
	})
}

// GetByUserID は指定されたユーザーの学習リソース一覧を取得する。
func (h *LearningResourceHandler) GetByUserID(c *gin.Context) {
	targetUserID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	currentUserID := c.GetUint("userID")

	resources, err := h.service.GetByUserID(targetUserID, currentUserID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, resources)
}

// GetPublic は公開学習リソース一覧をページネーション・フィルター付きで取得する。
func (h *LearningResourceHandler) GetPublic(c *gin.Context) {
	limit, offset := parseLimitOffset(c)
	category := c.Query("category")
	difficulty := c.Query("difficulty")

	resources, total, err := h.service.GetPublic(limit, offset, category, difficulty)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ResourceListResponse{
		Resources: resources,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// Search はキーワードで学習リソースを検索する。
func (h *LearningResourceHandler) Search(c *gin.Context) {
	query := c.Query("q")
	limit, offset := parseLimitOffset(c)

	resources, total, err := h.service.Search(query, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ResourceListResponse{
		Resources: resources,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// Update は指定された学習リソースを更新する。
func (h *LearningResourceHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[UpdateResourceRequest](c)
	if req == nil {
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

	resource, err := h.service.Update(id, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	// 公開設定が指定されている場合は別途更新
	if req.IsPublic != nil {
		resource, err = h.service.UpdateVisibility(id, userID, *req.IsPublic)
		if err != nil {
			respondError(c, err)
			return
		}
	}

	respondOK(c, resource)
}

// Delete は指定された学習リソースを削除する。
func (h *LearningResourceHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(id, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

// Like は学習リソースにいいねする。
func (h *LearningResourceHandler) Like(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.Like(userID, id); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("Resource liked"))
}

// Unlike は学習リソースのいいねを取り消す。
func (h *LearningResourceHandler) Unlike(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.Unlike(userID, id); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("Resource unliked"))
}

// SaveResource は学習リソースを保存する。
func (h *LearningResourceHandler) SaveResource(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.Save(userID, id); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("Resource saved"))
}

// UnsaveResource は学習リソースの保存を取り消す。
func (h *LearningResourceHandler) UnsaveResource(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.Unsave(userID, id); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("Resource unsaved"))
}

// GetByDifficulty は難易度別の公開学習リソースを取得する。
func (h *LearningResourceHandler) GetByDifficulty(c *gin.Context) {
	difficulty := c.Param("difficulty")
	limit, offset := parseLimitOffset(c)

	resources, total, err := h.service.GetByDifficulty(difficulty, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ResourceListResponse{
		Resources: resources,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

// GetSaved は認証ユーザーの保存済み学習リソース一覧を取得する。
func (h *LearningResourceHandler) GetSaved(c *gin.Context) {
	userID := c.GetUint("userID")
	limit, offset := parseLimitOffset(c)

	resources, total, err := h.service.GetSavedByUserID(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.ResourceListResponse{
		Resources: resources,
		Total:     total,
		Limit:     limit,
		Offset:    offset,
	})
}

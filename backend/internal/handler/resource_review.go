package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// ResourceReviewServiceInterface はResourceReviewHandlerが依存するサービスのインターフェース。
type ResourceReviewServiceInterface interface {
	Create(review *model.ResourceReview) error
	GetByResourceID(resourceID uint, limit, offset int) ([]model.ResourceReview, int64, error)
	Update(id, userID uint, rating int, comment string) (*model.ResourceReview, error)
	Delete(id, userID uint) error
}

// ResourceReviewHandler は学習リソースレビュー関連のHTTPハンドラ。
type ResourceReviewHandler struct {
	service ResourceReviewServiceInterface
}

// NewResourceReviewHandler は新しいResourceReviewHandlerインスタンスを生成する。
func NewResourceReviewHandler(s ResourceReviewServiceInterface) *ResourceReviewHandler {
	return &ResourceReviewHandler{service: s}
}

// Create は新しいレビューを作成する。
func (h *ResourceReviewHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	resourceID, ok := parseID(c, "id")
	if !ok {
		return
	}

	var req struct {
		Rating  int    `json:"rating" binding:"required"`
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "リクエストが不正です")
		return
	}

	review := &model.ResourceReview{
		UserID:     userID,
		ResourceID: resourceID,
		Rating:     req.Rating,
		Comment:    req.Comment,
	}

	if err := h.service.Create(review); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, review)
}

// GetByResourceID は指定リソースのレビュー一覧を取得する。
func (h *ResourceReviewHandler) GetByResourceID(c *gin.Context) {
	resourceID, ok := parseID(c, "id")
	if !ok {
		return
	}

	limit, offset := parseLimitOffset(c)

	reviews, total, err := h.service.GetByResourceID(resourceID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{
		"reviews": ensureSlice(reviews),
		"total":   total,
	})
}

// Update はレビューを更新する。
func (h *ResourceReviewHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	reviewID, ok := parseID(c, "reviewId")
	if !ok {
		return
	}

	var req struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, "リクエストが不正です")
		return
	}

	review, err := h.service.Update(reviewID, userID, req.Rating, req.Comment)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, review)
}

// Delete はレビューを削除する。
func (h *ResourceReviewHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	reviewID, ok := parseID(c, "reviewId")
	if !ok {
		return
	}

	if err := h.service.Delete(reviewID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// ResourceReviewHandler は学習リソースレビュー関連の HTTP ハンドラー。各操作は 1 責務の usecase に委譲する。
type ResourceReviewHandler struct {
	createReview *usecase.CreateResourceReviewUseCase
	listReviews  *usecase.ListResourceReviewsUseCase
	updateReview *usecase.UpdateResourceReviewUseCase
	deleteReview *usecase.DeleteResourceReviewUseCase
}

// NewResourceReviewHandler は ResourceReviewHandler を生成する。
func NewResourceReviewHandler(
	createReview *usecase.CreateResourceReviewUseCase,
	listReviews *usecase.ListResourceReviewsUseCase,
	updateReview *usecase.UpdateResourceReviewUseCase,
	deleteReview *usecase.DeleteResourceReviewUseCase,
) *ResourceReviewHandler {
	return &ResourceReviewHandler{
		createReview: createReview,
		listReviews:  listReviews,
		updateReview: updateReview,
		deleteReview: deleteReview,
	}
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

	if err := h.createReview.Execute(c.Request.Context(), review); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, review)
}

// resourceReviewListResponse はリソースレビュー一覧レスポンス。
type resourceReviewListResponse struct {
	Reviews []model.ResourceReview `json:"reviews"`
	Total   int64                  `json:"total"`
}

// GetByResourceID は指定リソースのレビュー一覧を取得する。
func (h *ResourceReviewHandler) GetByResourceID(c *gin.Context) {
	resourceID, ok := parseID(c, "id")
	if !ok {
		return
	}

	limit, offset := parseLimitOffset(c)

	reviews, total, err := h.listReviews.Execute(c.Request.Context(), resourceID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, resourceReviewListResponse{
		Reviews: ensureSlice(reviews),
		Total:   total,
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

	review, err := h.updateReview.Execute(c.Request.Context(), reviewID, userID, req.Rating, req.Comment)
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

	if err := h.deleteReview.Execute(c.Request.Context(), reviewID, userID); err != nil {
		respondError(c, err)
		return
	}

	respondDeleted(c)
}

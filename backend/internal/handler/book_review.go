package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/model"
)

// BookReviewServiceInterface はBookReviewServiceが実装すべきインターフェース。
type BookReviewServiceInterface interface {
	Create(review *model.BookReview) error
	GetByID(id uint) (*model.BookReview, error)
	GetByUserID(userID uint, limit, offset int) ([]model.BookReview, int64, error)
	GetAll(limit, offset int) ([]model.BookReview, int64, error)
	Update(id, userID uint, updates *model.BookReview) (*model.BookReview, error)
	Delete(id, userID uint) error
	GetByRating(userID uint, minRating, maxRating int) ([]model.BookReview, error)
}

// BookReviewHandler は書籍レビュー関連のHTTPハンドラ。
// 書籍レビューのCRUD・一覧取得を処理する。
type BookReviewHandler struct {
	service BookReviewServiceInterface
}

// NewBookReviewHandler は新しいBookReviewHandlerインスタンスを生成する。
func NewBookReviewHandler(s BookReviewServiceInterface) *BookReviewHandler {
	return &BookReviewHandler{service: s}
}

// Create は新しい書籍レビューを作成する。
func (h *BookReviewHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[dto.CreateBookReviewRequest](c)
	if req == nil {
		return
	}

	review := &model.BookReview{
		UserID:   userID,
		Title:    req.Title,
		Author:   req.Author,
		ISBN:     req.ISBN,
		Rating:   req.Rating,
		Review:   req.Review,
		ImageURL: req.ImageURL,
	}

	if err := h.service.Create(review); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, review)
}

// GetByID は指定IDの書籍レビューを取得する。
func (h *BookReviewHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	review, err := h.service.GetByID(id)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, review)
}

// GetByUserID は指定ユーザーの書籍レビュー一覧をページネーション付きで取得する。
func (h *BookReviewHandler) GetByUserID(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	limit, offset := parseLimitOffset(c)

	reviews, total, err := h.service.GetByUserID(userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.BookReviewListResponse{
		Reviews: reviews,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	})
}

// GetAll は書籍レビューの一覧をページネーション付きで取得する。
func (h *BookReviewHandler) GetAll(c *gin.Context) {
	limit, offset := parseLimitOffset(c)

	reviews, total, err := h.service.GetAll(limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.BookReviewListResponse{
		Reviews: reviews,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	})
}

// Update は指定IDの書籍レビューを更新する。
func (h *BookReviewHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[dto.UpdateBookReviewRequest](c)
	if req == nil {
		return
	}

	updates := &model.BookReview{}
	if req.Title != "" {
		updates.Title = req.Title
	}
	if req.Author != "" {
		updates.Author = req.Author
	}
	if req.ISBN != "" {
		updates.ISBN = req.ISBN
	}
	if req.Rating != nil {
		updates.Rating = *req.Rating
	}
	if req.Review != "" {
		updates.Review = req.Review
	}
	if req.ImageURL != "" {
		updates.ImageURL = req.ImageURL
	}

	review, err := h.service.Update(id, userID, updates)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, review)
}

// GetByRating は評価範囲で書籍レビューをフィルタリングして取得する。
// クエリパラメータ: min_rating, max_rating（1〜5）
func (h *BookReviewHandler) GetByRating(c *gin.Context) {
	userID := c.GetUint("userID")

	minRating, ok := parseQueryInt(c, "min_rating", "")
	if !ok {
		return
	}
	maxRating, ok := parseQueryInt(c, "max_rating", "")
	if !ok {
		return
	}

	reviews, err := h.service.GetByRating(userID, minRating, maxRating)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, reviews)
}

// Delete は指定IDの書籍レビューを削除する。
func (h *BookReviewHandler) Delete(c *gin.Context) {
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

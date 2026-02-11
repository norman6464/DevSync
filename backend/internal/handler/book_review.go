package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// BookReviewHandler は書籍レビュー関連のHTTPハンドラ。
// 書籍レビューのCRUD・一覧取得を処理する。
type BookReviewHandler struct {
	service *service.BookReviewService
}

// NewBookReviewHandler は新しいBookReviewHandlerインスタンスを生成する。
func NewBookReviewHandler(s *service.BookReviewService) *BookReviewHandler {
	return &BookReviewHandler{service: s}
}

// CreateBookReviewRequest は書籍レビュー作成のリクエストボディ。
type CreateBookReviewRequest struct {
	Title    string `json:"title" binding:"required,max=300"`
	Author   string `json:"author" binding:"max=200"`
	ISBN     string `json:"isbn" binding:"max=20"`
	Rating   int    `json:"rating" binding:"required,min=1,max=5"`
	Review   string `json:"review"`
	ImageURL string `json:"image_url"`
}

// UpdateBookReviewRequest は書籍レビュー更新のリクエストボディ。
type UpdateBookReviewRequest struct {
	Title    string `json:"title" binding:"max=300"`
	Author   string `json:"author" binding:"max=200"`
	ISBN     string `json:"isbn" binding:"max=20"`
	Rating   *int   `json:"rating" binding:"omitempty,min=1,max=5"`
	Review   string `json:"review"`
	ImageURL string `json:"image_url"`
}

// Create は新しい書籍レビューを作成する。
func (h *BookReviewHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[CreateBookReviewRequest](c)
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

// GetByUserID は指定ユーザーの書籍レビュー一覧を取得する。
func (h *BookReviewHandler) GetByUserID(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	reviews, err := h.service.GetByUserID(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, reviews)
}

// GetAll は書籍レビューの一覧をページネーション付きで取得する。
func (h *BookReviewHandler) GetAll(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	reviews, total, err := h.service.GetAll(limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{
		"reviews": reviews,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// Update は指定IDの書籍レビューを更新する。
func (h *BookReviewHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[UpdateBookReviewRequest](c)
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

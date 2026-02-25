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
	ArchiveReview(id, userID uint) error
	UnarchiveReview(id, userID uint) error
	UpdateStatus(id, userID uint, status model.ReviewStatus) error
	Search(query string, limit, offset int) ([]model.BookReview, int64, error)
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
	handleGetByIDPublic(c, h.service.GetByID)
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

	respondOK(c, ensureSlice(reviews))
}

// Archive は指定IDの書籍レビューをアーカイブする。
func (h *BookReviewHandler) Archive(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.ArchiveReview(id, userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"message": "書籍レビューをアーカイブしました"})
}

// Unarchive は指定IDの書籍レビューのアーカイブを解除する。
func (h *BookReviewHandler) Unarchive(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	if err := h.service.UnarchiveReview(id, userID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"message": "書籍レビューのアーカイブを解除しました"})
}

// UpdateStatus は書籍レビューの読書状態を更新する。
func (h *BookReviewHandler) UpdateStatus(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[dto.UpdateBookReviewStatusRequest](c)
	if req == nil {
		return
	}

	if err := h.service.UpdateStatus(id, userID, model.ReviewStatus(req.Status)); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"message": "読書状態を更新しました"})
}

// Search は書籍レビューをキーワード検索する。
func (h *BookReviewHandler) Search(c *gin.Context) {
	query := c.Query("q")
	limit, offset := parseLimitOffset(c)

	reviews, total, err := h.service.Search(query, limit, offset)
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

// Delete は指定IDの書籍レビューを削除する。
func (h *BookReviewHandler) Delete(c *gin.Context) {
	handleDelete(c, h.service.Delete)
}

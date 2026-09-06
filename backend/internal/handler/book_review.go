package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// BookReviewHandler は書籍レビュー関連のHTTPハンドラ。
type BookReviewHandler struct {
	create      *usecase.CreateBookReviewUseCase
	get         *usecase.GetBookReviewUseCase
	listByUser  *usecase.ListBookReviewsByUserUseCase
	listAll     *usecase.ListAllBookReviewsUseCase
	byRating    *usecase.ListBookReviewsByRatingUseCase
	search      *usecase.SearchBookReviewsUseCase
	update      *usecase.UpdateBookReviewUseCase
	setStatus   *usecase.UpdateBookReviewStatusUseCase
	setArchived *usecase.ArchiveBookReviewUseCase
	setProgress *usecase.UpdateBookReviewProgressUseCase
	remove      *usecase.DeleteBookReviewUseCase
	count       *usecase.CountBookReviewsUseCase
}

// NewBookReviewHandler は新しいBookReviewHandlerインスタンスを生成する。
func NewBookReviewHandler(
	create *usecase.CreateBookReviewUseCase,
	get *usecase.GetBookReviewUseCase,
	listByUser *usecase.ListBookReviewsByUserUseCase,
	listAll *usecase.ListAllBookReviewsUseCase,
	byRating *usecase.ListBookReviewsByRatingUseCase,
	search *usecase.SearchBookReviewsUseCase,
	update *usecase.UpdateBookReviewUseCase,
	setStatus *usecase.UpdateBookReviewStatusUseCase,
	setArchived *usecase.ArchiveBookReviewUseCase,
	setProgress *usecase.UpdateBookReviewProgressUseCase,
	remove *usecase.DeleteBookReviewUseCase,
	count *usecase.CountBookReviewsUseCase,
) *BookReviewHandler {
	return &BookReviewHandler{
		create: create, get: get, listByUser: listByUser, listAll: listAll,
		byRating: byRating, search: search, update: update,
		setStatus: setStatus, setArchived: setArchived, setProgress: setProgress,
		remove: remove, count: count,
	}
}

// createBookReviewRequest は書籍レビュー作成のリクエストボディ。
type createBookReviewRequest struct {
	Title      string `json:"title" binding:"required,min=1,max=300" validate:"required,min=1,max=300"`
	Author     string `json:"author" binding:"max=200" validate:"max=200"`
	ISBN       string `json:"isbn" binding:"max=20" validate:"max=20"`
	Rating     int    `json:"rating" binding:"required,min=1,max=5" validate:"required,min=1,max=5"`
	Review     string `json:"review" binding:"omitempty,max=5000"`
	ImageURL   string `json:"image_url" binding:"omitempty,http_url,max=2000"`
	TotalPages *int   `json:"total_pages" binding:"omitempty,min=0"`
}

// Create は新しい書籍レビューを作成する。
func (h *BookReviewHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	req := bindJSON[createBookReviewRequest](c)
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
	if req.TotalPages != nil {
		review.TotalPages = *req.TotalPages
	}

	if err := h.create.Execute(c.Request.Context(), review); err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, review)
}

// GetByID は指定IDの書籍レビューを取得する。
func (h *BookReviewHandler) GetByID(c *gin.Context) {
	handleGetByIDPublic(c, func(id uint) (*model.BookReview, error) {
		return h.get.Execute(c.Request.Context(), id)
	})
}

// bookReviewListResponse は書籍レビュー一覧レスポンス。
type bookReviewListResponse struct {
	Reviews []model.BookReview `json:"reviews"`
	Total   int64              `json:"total"`
	Limit   int                `json:"limit"`
	Offset  int                `json:"offset"`
}

// GetByUserID は指定ユーザーの書籍レビュー一覧をページネーション付きで取得する。
func (h *BookReviewHandler) GetByUserID(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	limit, offset := parseLimitOffset(c)

	reviews, total, err := h.listByUser.Execute(c.Request.Context(), userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, bookReviewListResponse{
		Reviews: reviews,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	})
}

// GetMyReviews は認証ユーザー自身の書籍レビュー一覧を取得する。
func (h *BookReviewHandler) GetMyReviews(c *gin.Context) {
	userID := c.GetUint("userID")

	limit, offset := parseLimitOffset(c)

	reviews, total, err := h.listByUser.Execute(c.Request.Context(), userID, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, bookReviewListResponse{
		Reviews: reviews,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	})
}

// GetAll は書籍レビューの一覧をページネーション付きで取得する。
func (h *BookReviewHandler) GetAll(c *gin.Context) {
	limit, offset := parseLimitOffset(c)

	reviews, total, err := h.listAll.Execute(c.Request.Context(), limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, bookReviewListResponse{
		Reviews: reviews,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	})
}

// updateBookReviewRequest は書籍レビュー更新のリクエストボディ。
type updateBookReviewRequest struct {
	Title    string `json:"title" binding:"omitempty,min=1,max=300" validate:"omitempty,min=1,max=300"`
	Author   string `json:"author" binding:"omitempty,max=200" validate:"omitempty,max=200"`
	ISBN     string `json:"isbn" binding:"omitempty,max=20" validate:"omitempty,max=20"`
	Rating   *int   `json:"rating" binding:"omitempty,min=1,max=5" validate:"omitempty,min=1,max=5"`
	Review   string `json:"review" binding:"omitempty,min=1,max=5000"`
	ImageURL string `json:"image_url" binding:"omitempty,http_url,max=2000"`
}

// Update は指定IDの書籍レビューを更新する。
func (h *BookReviewHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[updateBookReviewRequest](c)
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

	review, err := h.update.Execute(c.Request.Context(), id, userID, updates)
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

	reviews, err := h.byRating.Execute(c.Request.Context(), userID, minRating, maxRating)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, ensureSlice(reviews))
}

// Archive は指定IDの書籍レビューをアーカイブする。
func (h *BookReviewHandler) Archive(c *gin.Context) {
	handleAction(c, func(id, userID uint) error {
		return h.setArchived.Execute(c.Request.Context(), id, userID, true)
	}, "書籍レビューをアーカイブしました")
}

// Unarchive は指定IDの書籍レビューのアーカイブを解除する。
func (h *BookReviewHandler) Unarchive(c *gin.Context) {
	handleAction(c, func(id, userID uint) error {
		return h.setArchived.Execute(c.Request.Context(), id, userID, false)
	}, "書籍レビューのアーカイブを解除しました")
}

// updateBookReviewStatusRequest は書籍レビューの読書状態更新リクエスト。
type updateBookReviewStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateStatus は書籍レビューの読書状態を更新する。
func (h *BookReviewHandler) UpdateStatus(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[updateBookReviewStatusRequest](c)
	if req == nil {
		return
	}

	if err := h.setStatus.Execute(c.Request.Context(), id, userID, model.ReviewStatus(req.Status)); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("読書状態を更新しました"))
}

// updateReadingProgressRequest は読書進捗更新リクエスト。
type updateReadingProgressRequest struct {
	CurrentPage int `json:"current_page" binding:"required,min=0"`
}

// UpdateProgress は書籍レビューの読書進捗を更新する。
func (h *BookReviewHandler) UpdateProgress(c *gin.Context) {
	userID := c.GetUint("userID")
	id, ok := parseID(c, "id")
	if !ok {
		return
	}

	req := bindJSON[updateReadingProgressRequest](c)
	if req == nil {
		return
	}

	review, err := h.setProgress.Execute(c.Request.Context(), id, userID, req.CurrentPage)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, review)
}

// Search は書籍レビューをキーワード検索する。
func (h *BookReviewHandler) Search(c *gin.Context) {
	query, ok := parseSearchQuery(c)
	if !ok {
		return
	}
	limit, offset := parseLimitOffset(c)

	reviews, total, err := h.search.Execute(c.Request.Context(), query, limit, offset)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, bookReviewListResponse{
		Reviews: reviews,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	})
}

// Delete は指定IDの書籍レビューを削除する。
func (h *BookReviewHandler) Delete(c *gin.Context) {
	handleDelete(c, func(id, userID uint) error {
		return h.remove.Execute(c.Request.Context(), id, userID)
	})
}

// GetMyCount は認証ユーザーの書籍レビュー総数を取得する。
func (h *BookReviewHandler) GetMyCount(c *gin.Context) {
	userID := c.GetUint("userID")

	count, err := h.count.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, gin.H{"count": count})
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

type BookReviewHandler struct {
	repo *repository.BookReviewRepository
}

func NewBookReviewHandler(repo *repository.BookReviewRepository) *BookReviewHandler {
	return &BookReviewHandler{repo: repo}
}

type CreateBookReviewRequest struct {
	Title    string `json:"title" binding:"required,max=255"`
	Author   string `json:"author" binding:"required,max=255"`
	CoverURL string `json:"cover_url" binding:"omitempty,url,max=500"`
	Rating   int    `json:"rating" binding:"required,min=1,max=5"`
	Content  string `json:"content"`
}

type UpdateBookReviewRequest struct {
	Title    string `json:"title" binding:"omitempty,max=255"`
	Author   string `json:"author" binding:"omitempty,max=255"`
	CoverURL string `json:"cover_url" binding:"omitempty,max=500"`
	Rating   int    `json:"rating" binding:"omitempty,min=1,max=5"`
	Content  string `json:"content"`
}

type BookReviewResponse struct {
	model.BookReview
	HasLiked bool `json:"has_liked"`
}

func (h *BookReviewHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateBookReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review := &model.BookReview{
		UserID:   userID,
		Title:    req.Title,
		Author:   req.Author,
		CoverURL: req.CoverURL,
		Rating:   req.Rating,
		Content:  req.Content,
	}

	if err := h.repo.Create(review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func (h *BookReviewHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	review, err := h.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book review not found"})
		return
	}

	userID := c.GetUint("userID")
	hasLiked, _ := h.repo.HasLiked(userID, uint(id))

	c.JSON(http.StatusOK, BookReviewResponse{
		BookReview: *review,
		HasLiked:   hasLiked,
	})
}

func (h *BookReviewHandler) GetByUserID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	reviews, err := h.repo.FindByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func (h *BookReviewHandler) GetAll(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	reviews, err := h.repo.FindAll(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch book reviews"})
		return
	}

	userID := c.GetUint("userID")

	// Get liked status for all reviews
	reviewIDs := make([]uint, len(reviews))
	for i, r := range reviews {
		reviewIDs[i] = r.ID
	}

	likedIDs, _ := h.repo.GetUserLikedReviewIDs(userID, reviewIDs)
	likedMap := make(map[uint]bool)
	for _, id := range likedIDs {
		likedMap[id] = true
	}

	response := make([]BookReviewResponse, len(reviews))
	for i, review := range reviews {
		response[i] = BookReviewResponse{
			BookReview: review,
			HasLiked:   likedMap[review.ID],
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *BookReviewHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID := c.GetUint("userID")

	review, err := h.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book review not found"})
		return
	}

	if review.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	var req UpdateBookReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != "" {
		review.Title = req.Title
	}
	if req.Author != "" {
		review.Author = req.Author
	}
	if req.CoverURL != "" {
		review.CoverURL = req.CoverURL
	}
	if req.Rating > 0 {
		review.Rating = req.Rating
	}
	if req.Content != "" {
		review.Content = req.Content
	}

	if err := h.repo.Update(review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book review"})
		return
	}

	c.JSON(http.StatusOK, review)
}

func (h *BookReviewHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID := c.GetUint("userID")

	review, err := h.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book review not found"})
		return
	}

	if review.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book review deleted"})
}

func (h *BookReviewHandler) Like(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID := c.GetUint("userID")

	// Check if already liked
	hasLiked, _ := h.repo.HasLiked(userID, uint(id))
	if hasLiked {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already liked"})
		return
	}

	like := &model.BookReviewLike{
		UserID:       userID,
		BookReviewID: uint(id),
	}

	if err := h.repo.CreateLike(like); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like book review"})
		return
	}

	h.repo.IncrementLikeCount(uint(id))

	c.JSON(http.StatusOK, gin.H{"message": "Liked"})
}

func (h *BookReviewHandler) Unlike(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID := c.GetUint("userID")

	if err := h.repo.DeleteLike(userID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike book review"})
		return
	}

	h.repo.DecrementLikeCount(uint(id))

	c.JSON(http.StatusOK, gin.H{"message": "Unliked"})
}

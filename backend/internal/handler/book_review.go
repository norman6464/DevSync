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
	Title    string `json:"title" binding:"required,max=300"`
	Author   string `json:"author" binding:"max=200"`
	ISBN     string `json:"isbn" binding:"max=20"`
	Rating   int    `json:"rating" binding:"required,min=1,max=5"`
	Review   string `json:"review"`
	ImageURL string `json:"image_url"`
}

type UpdateBookReviewRequest struct {
	Title    string `json:"title" binding:"max=300"`
	Author   string `json:"author" binding:"max=200"`
	ISBN     string `json:"isbn" binding:"max=20"`
	Rating   *int   `json:"rating" binding:"omitempty,min=1,max=5"`
	Review   string `json:"review"`
	ImageURL string `json:"image_url"`
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
		ISBN:     req.ISBN,
		Rating:   req.Rating,
		Review:   req.Review,
		ImageURL: req.ImageURL,
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	review, err := h.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	c.JSON(http.StatusOK, review)
}

func (h *BookReviewHandler) GetByUserID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	reviews, err := h.repo.FindByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
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

	reviews, total, err := h.repo.FindAll(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews": reviews,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

func (h *BookReviewHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	review, err := h.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	if review.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this review"})
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
	if req.ISBN != "" {
		review.ISBN = req.ISBN
	}
	if req.Rating != nil {
		review.Rating = *req.Rating
	}
	if req.Review != "" {
		review.Review = req.Review
	}
	if req.ImageURL != "" {
		review.ImageURL = req.ImageURL
	}

	if err := h.repo.Update(review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	c.JSON(http.StatusOK, review)
}

func (h *BookReviewHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	review, err := h.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	if review.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this review"})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}

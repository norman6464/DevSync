package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

type QuestionHandler struct {
	repo *repository.QuestionRepository
}

func NewQuestionHandler(repo *repository.QuestionRepository) *QuestionHandler {
	return &QuestionHandler{repo: repo}
}

type CreateQuestionRequest struct {
	Title string `json:"title" binding:"required,max=500"`
	Body  string `json:"body" binding:"required"`
	Tags  string `json:"tags"`
}

type UpdateQuestionRequest struct {
	Title string `json:"title" binding:"max=500"`
	Body  string `json:"body"`
	Tags  string `json:"tags"`
}

type VoteRequest struct {
	Value int `json:"value" binding:"required,oneof=1 -1"`
}

func (h *QuestionHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	var req CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question := &model.Question{
		UserID: userID,
		Title:  req.Title,
		Body:   req.Body,
		Tags:   req.Tags,
	}

	if err := h.repo.Create(question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question"})
		return
	}

	c.JSON(http.StatusCreated, question)
}

func (h *QuestionHandler) GetAll(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	tag := c.Query("tag")
	sort := c.DefaultQuery("sort", "newest")

	if limit > 100 {
		limit = 100
	}

	questions, total, err := h.repo.FindAll(limit, offset, tag, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"questions": questions,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

func (h *QuestionHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit > 100 {
		limit = 100
	}

	questions, total, err := h.repo.Search(q, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search questions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"questions": questions,
		"total":     total,
		"limit":     limit,
		"offset":    offset,
	})
}

func (h *QuestionHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	question, err := h.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	userVote, _ := h.repo.GetUserVote(userID, uint(id))

	c.JSON(http.StatusOK, gin.H{
		"question":  question,
		"user_vote": userVote,
	})
}

func (h *QuestionHandler) GetByUserID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	questions, err := h.repo.FindByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}

	c.JSON(http.StatusOK, questions)
}

func (h *QuestionHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	question, err := h.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	if question.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this question"})
		return
	}

	var req UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != "" {
		question.Title = req.Title
	}
	if req.Body != "" {
		question.Body = req.Body
	}
	if req.Tags != "" {
		question.Tags = req.Tags
	}

	if err := h.repo.Update(question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update question"})
		return
	}

	c.JSON(http.StatusOK, question)
}

func (h *QuestionHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	question, err := h.repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	if question.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this question"})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete question"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}

func (h *QuestionHandler) Vote(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Vote(userID, uint(id), req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Voted successfully"})
}

func (h *QuestionHandler) RemoveVote(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	if err := h.repo.RemoveVote(userID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote removed successfully"})
}

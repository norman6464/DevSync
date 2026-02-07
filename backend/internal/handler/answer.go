package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

type AnswerHandler struct {
	answerRepo   *repository.AnswerRepository
	questionRepo *repository.QuestionRepository
}

func NewAnswerHandler(answerRepo *repository.AnswerRepository, questionRepo *repository.QuestionRepository) *AnswerHandler {
	return &AnswerHandler{answerRepo: answerRepo, questionRepo: questionRepo}
}

func (h *AnswerHandler) GetByQuestionID(c *gin.Context) {
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	answers, err := h.answerRepo.FindByQuestionID(uint(questionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch answers"})
		return
	}

	c.JSON(http.StatusOK, answers)
}

type CreateAnswerRequest struct {
	Body string `json:"body" binding:"required"`
}

type UpdateAnswerRequest struct {
	Body string `json:"body" binding:"required"`
}

func (h *AnswerHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	// Verify question exists
	if _, err := h.questionRepo.FindByID(uint(questionID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	var req CreateAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer := &model.Answer{
		UserID:     userID,
		QuestionID: uint(questionID),
		Body:       req.Body,
	}

	if err := h.answerRepo.Create(answer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create answer"})
		return
	}

	c.JSON(http.StatusCreated, answer)
}

func (h *AnswerHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	answerID, err := strconv.ParseUint(c.Param("answerId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
		return
	}

	answer, err := h.answerRepo.FindByID(uint(answerID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
		return
	}

	if answer.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this answer"})
		return
	}

	var req UpdateAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer.Body = req.Body

	if err := h.answerRepo.Update(answer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update answer"})
		return
	}

	c.JSON(http.StatusOK, answer)
}

func (h *AnswerHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	answerID, err := strconv.ParseUint(c.Param("answerId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
		return
	}

	answer, err := h.answerRepo.FindByID(uint(answerID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
		return
	}

	if answer.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this answer"})
		return
	}

	if err := h.answerRepo.Delete(answer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete answer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Answer deleted successfully"})
}

func (h *AnswerHandler) SetBestAnswer(c *gin.Context) {
	userID := c.GetUint("userID")
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	answerID, err := strconv.ParseUint(c.Param("answerId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
		return
	}

	// Verify question exists and user is the owner
	question, err := h.questionRepo.FindByID(uint(questionID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	if question.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the question owner can select the best answer"})
		return
	}

	// Verify answer belongs to this question
	answer, err := h.answerRepo.FindByID(uint(answerID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
		return
	}
	if answer.QuestionID != uint(questionID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Answer does not belong to this question"})
		return
	}

	if err := h.answerRepo.SetBestAnswer(uint(questionID), uint(answerID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set best answer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Best answer set successfully"})
}

func (h *AnswerHandler) Vote(c *gin.Context) {
	userID := c.GetUint("userID")
	answerID, err := strconv.ParseUint(c.Param("answerId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
		return
	}

	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.answerRepo.Vote(userID, uint(answerID), req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Voted successfully"})
}

func (h *AnswerHandler) RemoveVote(c *gin.Context) {
	userID := c.GetUint("userID")
	answerID, err := strconv.ParseUint(c.Param("answerId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid answer ID"})
		return
	}

	if err := h.answerRepo.RemoveVote(userID, uint(answerID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote removed successfully"})
}

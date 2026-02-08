package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

type AnswerHandler struct {
	service *service.AnswerService
}

func NewAnswerHandler(s *service.AnswerService) *AnswerHandler {
	return &AnswerHandler{service: s}
}

func (h *AnswerHandler) GetByQuestionID(c *gin.Context) {
	questionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	answers, err := h.service.GetByQuestionID(uint(questionID))
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

	if err := h.service.Create(answer); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
			return
		}
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

	var req UpdateAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	answer, err := h.service.Update(uint(answerID), userID, req.Body)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this answer"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
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

	if err := h.service.Delete(uint(answerID), userID); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this answer"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Answer not found"})
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

	if err := h.service.SetBestAnswer(uint(questionID), uint(answerID), userID); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only the question owner can select the best answer"})
			return
		}
		if errors.Is(err, service.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Answer does not belong to this question"})
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question or answer not found"})
			return
		}
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

	if err := h.service.Vote(userID, uint(answerID), req.Value); err != nil {
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

	if err := h.service.RemoveVote(userID, uint(answerID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote removed successfully"})
}

package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// QuestionHandler は質問関連のHTTPハンドラ。
// 質問のCRUD・検索・投票を処理する。
type QuestionHandler struct {
	service *service.QuestionService
}

// NewQuestionHandler は新しいQuestionHandlerインスタンスを生成する。
func NewQuestionHandler(s *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{service: s}
}

// CreateQuestionRequest は質問作成のリクエストボディ。
type CreateQuestionRequest struct {
	Title string `json:"title" binding:"required,max=500"`
	Body  string `json:"body" binding:"required"`
	Tags  string `json:"tags"`
}

// UpdateQuestionRequest は質問更新のリクエストボディ。
type UpdateQuestionRequest struct {
	Title string `json:"title" binding:"max=500"`
	Body  string `json:"body"`
	Tags  string `json:"tags"`
}

// VoteRequest は投票のリクエストボディ。
type VoteRequest struct {
	Value int `json:"value" binding:"required,oneof=1 -1"`
}

// Create は新しい質問を作成する。
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

	if err := h.service.Create(question); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question"})
		return
	}

	c.JSON(http.StatusCreated, question)
}

// GetAll は質問一覧をページネーション・タグ・ソート付きで取得する。
func (h *QuestionHandler) GetAll(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	tag := c.Query("tag")
	sort := c.DefaultQuery("sort", "newest")

	if limit > 100 {
		limit = 100
	}

	questions, total, err := h.service.GetAll(limit, offset, tag, sort)
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

// Search はキーワードで質問を検索する。
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

	questions, total, err := h.service.Search(q, limit, offset)
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

// GetByID は指定されたIDの質問を取得する。
func (h *QuestionHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	question, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	userVote, _ := h.service.GetUserVote(userID, uint(id))

	c.JSON(http.StatusOK, gin.H{
		"question":  question,
		"user_vote": userVote,
	})
}

// GetByUserID は指定されたユーザーの質問一覧を取得する。
func (h *QuestionHandler) GetByUserID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	questions, err := h.service.GetByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}

	c.JSON(http.StatusOK, questions)
}

// Update は指定された質問を更新する。
func (h *QuestionHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	var req UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question, err := h.service.Update(uint(id), userID, req.Title, req.Body, req.Tags)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to update this question"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, question)
}

// Delete は指定された質問を削除する。
func (h *QuestionHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	if err := h.service.Delete(uint(id), userID); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this question"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
}

// Vote は質問に投票する。
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

	if err := h.service.Vote(userID, uint(id), req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Voted successfully"})
}

// RemoveVote は質問への投票を取り消す。
func (h *QuestionHandler) RemoveVote(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}

	if err := h.service.RemoveVote(userID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote removed successfully"})
}

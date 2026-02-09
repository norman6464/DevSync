package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// LearningLogHandler は学習ログ関連のHTTPハンドラ。
// 学習ログのCRUD・ストリーク・カレンダーデータの取得を処理する。
type LearningLogHandler struct {
	service *service.LearningLogService
}

// NewLearningLogHandler は新しいLearningLogHandlerインスタンスを生成する。
func NewLearningLogHandler(s *service.LearningLogService) *LearningLogHandler {
	return &LearningLogHandler{service: s}
}

// Create は新しい学習ログを作成する。
func (h *LearningLogHandler) Create(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content" binding:"required"`
		Category string `json:"category"`
		Duration int    `json:"duration"`
		Source   string `json:"source"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title and content are required"})
		return
	}

	log := &model.LearningLog{
		UserID:   userID,
		Title:    req.Title,
		Content:  req.Content,
		Category: model.LogCategory(req.Category),
		Duration: req.Duration,
		Source:   model.LogSource(req.Source),
	}

	// カテゴリが未指定の場合はデフォルト値を設定
	if req.Category == "" {
		log.Category = model.LogCategoryOther
	}

	if err := h.service.Create(log); err != nil {
		if errors.Is(err, service.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request parameters"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create log"})
		return
	}

	c.JSON(http.StatusCreated, log)
}

// Update は指定された学習ログを更新する。
func (h *LearningLogHandler) Update(c *gin.Context) {
	userID := c.GetUint("userID")
	logID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid log ID"})
		return
	}

	var req struct {
		Title    *string `json:"title"`
		Content  *string `json:"content"`
		Category *string `json:"category"`
		Duration *int    `json:"duration"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	updates := &model.LearningLog{}
	if req.Title != nil {
		updates.Title = *req.Title
	}
	if req.Content != nil {
		updates.Content = *req.Content
	}
	if req.Category != nil {
		updates.Category = model.LogCategory(*req.Category)
	}
	if req.Duration != nil {
		updates.Duration = *req.Duration
	}

	log, err := h.service.Update(uint(logID), userID, updates)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// Delete は指定された学習ログを削除する。
func (h *LearningLogHandler) Delete(c *gin.Context) {
	userID := c.GetUint("userID")
	logID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid log ID"})
		return
	}

	if err := h.service.Delete(uint(logID), userID); err != nil {
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "log not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "log deleted"})
}

// GetByID は指定されたIDの学習ログを取得する。
func (h *LearningLogHandler) GetByID(c *gin.Context) {
	logID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid log ID"})
		return
	}

	log, err := h.service.GetByID(uint(logID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "log not found"})
		return
	}

	c.JSON(http.StatusOK, log)
}

// GetMyLogs は認証ユーザー自身の学習ログ一覧を取得する。
func (h *LearningLogHandler) GetMyLogs(c *gin.Context) {
	userID := c.GetUint("userID")

	logs, err := h.service.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetByUserID は指定されたユーザーの学習ログ一覧を取得する。
func (h *LearningLogHandler) GetByUserID(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	logs, err := h.service.GetByUserID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// GetStreakInfo は指定されたユーザーのストリーク情報を取得する。
func (h *LearningLogHandler) GetStreakInfo(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	info, err := h.service.GetStreakInfo(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get streak info"})
		return
	}

	c.JSON(http.StatusOK, info)
}

// GetCalendarData はカレンダー表示用の日別学習ログ件数を取得する。
func (h *LearningLogHandler) GetCalendarData(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	entries, err := h.service.GetCalendarData(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get calendar data"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

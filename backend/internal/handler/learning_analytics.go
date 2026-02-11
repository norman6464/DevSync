package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/service"
)

// LearningAnalyticsHandler は学習分析関連のHTTPハンドラ。
// ヒートマップ、カテゴリ別、トレンド、生産性スコア、AIインサイトの取得を処理する。
type LearningAnalyticsHandler struct {
	service *service.LearningAnalyticsService
}

// NewLearningAnalyticsHandler は新しいLearningAnalyticsHandlerインスタンスを生成する。
func NewLearningAnalyticsHandler(s *service.LearningAnalyticsService) *LearningAnalyticsHandler {
	return &LearningAnalyticsHandler{service: s}
}

// GetHeatmap は指定ユーザーの学習時間ヒートマップを返す。
func (h *LearningAnalyticsHandler) GetHeatmap(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	data, err := h.service.GetHeatmap(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetCategoryBreakdown は指定ユーザーのカテゴリ別学習時間を返す。
func (h *LearningAnalyticsHandler) GetCategoryBreakdown(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	data, err := h.service.GetCategoryBreakdown(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetProductivityScore は指定ユーザーの生産性スコアを返す。
func (h *LearningAnalyticsHandler) GetProductivityScore(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	score, err := h.service.GetProductivityScore(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, score)
}

// GetWeeklyTrends は指定ユーザーの週間学習トレンドを返す。
func (h *LearningAnalyticsHandler) GetWeeklyTrends(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	weeks := 12
	if w := c.Query("weeks"); w != "" {
		if parsed, err := strconv.Atoi(w); err == nil && parsed > 0 {
			weeks = parsed
		}
	}

	data, err := h.service.GetWeeklyTrends(uint(userID), weeks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetInsights は認証済みユーザーのAIインサイトを返す。
func (h *LearningAnalyticsHandler) GetInsights(c *gin.Context) {
	userID := c.GetUint("userID")

	insights, err := h.service.GetInsights(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, insights)
}

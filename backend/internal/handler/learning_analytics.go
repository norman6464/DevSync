package handler

import (
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
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	data, err := h.service.GetHeatmap(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, data)
}

// GetCategoryBreakdown は指定ユーザーのカテゴリ別学習時間を返す。
func (h *LearningAnalyticsHandler) GetCategoryBreakdown(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	data, err := h.service.GetCategoryBreakdown(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, data)
}

// GetProductivityScore は指定ユーザーの生産性スコアを返す。
func (h *LearningAnalyticsHandler) GetProductivityScore(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	score, err := h.service.GetProductivityScore(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, score)
}

// GetWeeklyTrends は指定ユーザーの週間学習トレンドを返す。
func (h *LearningAnalyticsHandler) GetWeeklyTrends(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	weeks := 12
	if w := c.Query("weeks"); w != "" {
		if parsed, err := strconv.Atoi(w); err == nil && parsed > 0 {
			weeks = parsed
		}
	}

	data, err := h.service.GetWeeklyTrends(userID, weeks)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, data)
}

// GetInsights は認証済みユーザーのAIインサイトを返す。
func (h *LearningAnalyticsHandler) GetInsights(c *gin.Context) {
	userID := c.GetUint("userID")

	insights, err := h.service.GetInsights(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, insights)
}

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

// LearningAnalyticsServiceInterface はLearningAnalyticsHandlerが依存するサービスのインターフェース。
type LearningAnalyticsServiceInterface interface {
	GetHeatmap(userID uint) ([]model.HeatmapEntry, error)
	GetCategoryBreakdown(userID uint) ([]model.CategoryBreakdown, error)
	GetWeeklyTrends(userID uint, weeks int) ([]model.WeeklyTrend, error)
	GetProductivityScore(userID uint) (*model.ProductivityScore, error)
	GetInsights(userID uint) ([]model.AIInsight, error)
}

// LearningAnalyticsHandler は学習分析関連のHTTPハンドラ。
// ヒートマップ、カテゴリ別、トレンド、生産性スコア、AIインサイトの取得を処理する。
type LearningAnalyticsHandler struct {
	service LearningAnalyticsServiceInterface
}

// NewLearningAnalyticsHandler は新しいLearningAnalyticsHandlerインスタンスを生成する。
func NewLearningAnalyticsHandler(s LearningAnalyticsServiceInterface) *LearningAnalyticsHandler {
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

	weeks := parseQueryIntSilent(c, "weeks", 12)

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

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// LearningAnalyticsHandler は学習分析関連のHTTPハンドラ。
// ヒートマップ、カテゴリ別、トレンド、生産性スコア、AIインサイトの取得を処理する。
type LearningAnalyticsHandler struct {
	heatmap      *usecase.GetLearningHeatmapUseCase
	categories   *usecase.GetCategoryBreakdownUseCase
	weeklyTrends *usecase.GetWeeklyTrendsUseCase
	dayOfWeek    *usecase.GetDayOfWeekSummaryUseCase
	productivity *usecase.GetProductivityScoreUseCase
	insights     *usecase.GetLearningInsightsUseCase
}

// NewLearningAnalyticsHandler は新しいLearningAnalyticsHandlerインスタンスを生成する。
func NewLearningAnalyticsHandler(
	heatmap *usecase.GetLearningHeatmapUseCase,
	categories *usecase.GetCategoryBreakdownUseCase,
	weeklyTrends *usecase.GetWeeklyTrendsUseCase,
	dayOfWeek *usecase.GetDayOfWeekSummaryUseCase,
	productivity *usecase.GetProductivityScoreUseCase,
	insights *usecase.GetLearningInsightsUseCase,
) *LearningAnalyticsHandler {
	return &LearningAnalyticsHandler{
		heatmap: heatmap, categories: categories, weeklyTrends: weeklyTrends,
		dayOfWeek: dayOfWeek, productivity: productivity, insights: insights,
	}
}

// GetHeatmap は指定ユーザーの学習時間ヒートマップを返す。
func (h *LearningAnalyticsHandler) GetHeatmap(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	data, err := h.heatmap.Execute(c.Request.Context(), userID)
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

	data, err := h.categories.Execute(c.Request.Context(), userID)
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

	score, err := h.productivity.Execute(c.Request.Context(), userID)
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

	data, err := h.weeklyTrends.Execute(c.Request.Context(), userID, weeks)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, data)
}

// GetDayOfWeekSummary は指定ユーザーの曜日別学習サマリーを返す。
func (h *LearningAnalyticsHandler) GetDayOfWeekSummary(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	data, err := h.dayOfWeek.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, data)
}

// GetInsights は認証済みユーザーのAIインサイトを返す。
func (h *LearningAnalyticsHandler) GetInsights(c *gin.Context) {
	userID := c.GetUint("userID")

	insights, err := h.insights.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, insights)
}

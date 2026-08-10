package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// ActivityReportHandler はアクティビティレポート関連の HTTP ハンドラ。
// 週次・月次レポートの取得および期間比較を処理する。
type ActivityReportHandler struct {
	getWeekly     *usecase.GetWeeklyActivityReportUseCase
	getMonthly    *usecase.GetMonthlyActivityReportUseCase
	getComparison *usecase.GetActivityReportComparisonUseCase
}

// NewActivityReportHandler は ActivityReportHandler を生成する。
func NewActivityReportHandler(
	getWeekly *usecase.GetWeeklyActivityReportUseCase,
	getMonthly *usecase.GetMonthlyActivityReportUseCase,
	getComparison *usecase.GetActivityReportComparisonUseCase,
) *ActivityReportHandler {
	return &ActivityReportHandler{
		getWeekly:     getWeekly,
		getMonthly:    getMonthly,
		getComparison: getComparison,
	}
}

// GetWeeklyReport は指定ユーザーの週次アクティビティレポートを返す。
func (h *ActivityReportHandler) GetWeeklyReport(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	report, err := h.getWeekly.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, report)
}

// GetMonthlyReport は指定ユーザーの月次アクティビティレポートを返す。
func (h *ActivityReportHandler) GetMonthlyReport(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	report, err := h.getMonthly.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, report)
}

// GetMyWeeklyReport は現在のユーザーの週次レポートを返す。
func (h *ActivityReportHandler) GetMyWeeklyReport(c *gin.Context) {
	userID := c.GetUint("userID")

	report, err := h.getWeekly.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, report)
}

// GetMyMonthlyReport は現在のユーザーの月次レポートを返す。
func (h *ActivityReportHandler) GetMyMonthlyReport(c *gin.Context) {
	userID := c.GetUint("userID")

	report, err := h.getMonthly.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, report)
}

// GetComparison は現在の期間と前期間のアクティビティ比較を返す。
func (h *ActivityReportHandler) GetComparison(c *gin.Context) {
	userID := c.GetUint("userID")
	periodParam := c.Query("period")

	period := model.ReportPeriodWeekly
	if periodParam == "monthly" {
		period = model.ReportPeriodMonthly
	}

	comparison, err := h.getComparison.Execute(c.Request.Context(), userID, period)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, comparison)
}

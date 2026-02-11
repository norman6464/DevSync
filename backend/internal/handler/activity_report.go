package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

// ActivityReportHandler はアクティビティレポート関連のHTTPハンドラ。
// 週次・月次レポートの取得および期間比較を処理する。
type ActivityReportHandler struct {
	service *service.ActivityReportService
}

// NewActivityReportHandler は新しいActivityReportHandlerインスタンスを生成する。
func NewActivityReportHandler(s *service.ActivityReportService) *ActivityReportHandler {
	return &ActivityReportHandler{service: s}
}

// GetWeeklyReport は指定ユーザーの週次アクティビティレポートを返す。
func (h *ActivityReportHandler) GetWeeklyReport(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	report, err := h.service.GetWeeklyReport(userID)
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

	report, err := h.service.GetMonthlyReport(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, report)
}

// GetMyWeeklyReport は現在のユーザーの週次レポートを返す。
func (h *ActivityReportHandler) GetMyWeeklyReport(c *gin.Context) {
	userID := c.GetUint("userID")

	report, err := h.service.GetWeeklyReport(userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, report)
}

// GetMyMonthlyReport は現在のユーザーの月次レポートを返す。
func (h *ActivityReportHandler) GetMyMonthlyReport(c *gin.Context) {
	userID := c.GetUint("userID")

	report, err := h.service.GetMonthlyReport(userID)
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

	comparison, err := h.service.GetComparison(userID, period)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, comparison)
}

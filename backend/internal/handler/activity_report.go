package handler

import (
	"net/http"
	"strconv"

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
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	report, err := h.service.GetWeeklyReport(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetMonthlyReport は指定ユーザーの月次アクティビティレポートを返す。
func (h *ActivityReportHandler) GetMonthlyReport(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	report, err := h.service.GetMonthlyReport(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetMyWeeklyReport は現在のユーザーの週次レポートを返す。
func (h *ActivityReportHandler) GetMyWeeklyReport(c *gin.Context) {
	userID := c.GetUint("userID")

	report, err := h.service.GetWeeklyReport(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetMyMonthlyReport は現在のユーザーの月次レポートを返す。
func (h *ActivityReportHandler) GetMyMonthlyReport(c *gin.Context) {
	userID := c.GetUint("userID")

	report, err := h.service.GetMonthlyReport(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report"})
		return
	}

	c.JSON(http.StatusOK, report)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate comparison"})
		return
	}

	c.JSON(http.StatusOK, comparison)
}

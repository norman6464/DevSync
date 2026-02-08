package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
)

type ActivityReportHandler struct {
	service *service.ActivityReportService
}

func NewActivityReportHandler(s *service.ActivityReportService) *ActivityReportHandler {
	return &ActivityReportHandler{service: s}
}

// GetWeeklyReport returns the weekly activity report for a user
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

// GetMonthlyReport returns the monthly activity report for a user
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

// GetMyWeeklyReport returns the weekly report for the current user
func (h *ActivityReportHandler) GetMyWeeklyReport(c *gin.Context) {
	userID := c.GetUint("userID")

	report, err := h.service.GetWeeklyReport(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetMyMonthlyReport returns the monthly report for the current user
func (h *ActivityReportHandler) GetMyMonthlyReport(c *gin.Context) {
	userID := c.GetUint("userID")

	report, err := h.service.GetMonthlyReport(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetComparison returns the comparison between current and previous period
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

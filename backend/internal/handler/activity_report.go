package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

type ActivityReportHandler struct {
	reportRepo *repository.ActivityReportRepository
}

func NewActivityReportHandler(reportRepo *repository.ActivityReportRepository) *ActivityReportHandler {
	return &ActivityReportHandler{reportRepo: reportRepo}
}

// GetWeeklyReport returns the weekly activity report for a user
func (h *ActivityReportHandler) GetWeeklyReport(c *gin.Context) {
	userIDParam := c.Param("userId")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	report, err := h.reportRepo.GetWeeklyReport(uint(userID))
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

	report, err := h.reportRepo.GetMonthlyReport(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetMyWeeklyReport returns the weekly report for the current user
func (h *ActivityReportHandler) GetMyWeeklyReport(c *gin.Context) {
	userID := c.GetUint("userID")

	report, err := h.reportRepo.GetWeeklyReport(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetMyMonthlyReport returns the monthly report for the current user
func (h *ActivityReportHandler) GetMyMonthlyReport(c *gin.Context) {
	userID := c.GetUint("userID")

	report, err := h.reportRepo.GetMonthlyReport(userID)
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

	comparison, err := h.reportRepo.GetComparison(userID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate comparison"})
		return
	}

	c.JSON(http.StatusOK, comparison)
}

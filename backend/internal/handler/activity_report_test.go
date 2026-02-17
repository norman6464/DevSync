package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
)

func TestActivityReport_GetWeeklyReport_Success(t *testing.T) {
	h, svc := setupActivityReportHandler()
	report := &model.ActivityReport{Period: model.ReportPeriodWeekly, TotalContributions: 10}
	svc.On("GetWeeklyReport", uint(1)).Return(report, nil)

	r := gin.New()
	r.GET("/users/:userId/reports/weekly", h.GetWeeklyReport)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1/reports/weekly", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestActivityReport_GetWeeklyReport_InvalidID(t *testing.T) {
	h, _ := setupActivityReportHandler()

	r := gin.New()
	r.GET("/users/:userId/reports/weekly", h.GetWeeklyReport)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/abc/reports/weekly", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusBadRequest)
}

func TestActivityReport_GetWeeklyReport_ServiceError(t *testing.T) {
	h, svc := setupActivityReportHandler()
	svc.On("GetWeeklyReport", uint(1)).Return(nil, errors.New("db error"))

	r := gin.New()
	r.GET("/users/:userId/reports/weekly", h.GetWeeklyReport)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1/reports/weekly", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestActivityReport_GetMonthlyReport_Success(t *testing.T) {
	h, svc := setupActivityReportHandler()
	report := &model.ActivityReport{Period: model.ReportPeriodMonthly, TotalContributions: 42}
	svc.On("GetMonthlyReport", uint(1)).Return(report, nil)

	r := gin.New()
	r.GET("/users/:userId/reports/monthly", h.GetMonthlyReport)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1/reports/monthly", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestActivityReport_GetMyWeeklyReport_Success(t *testing.T) {
	h, svc := setupActivityReportHandler()
	report := &model.ActivityReport{Period: model.ReportPeriodWeekly}
	svc.On("GetWeeklyReport", uint(5)).Return(report, nil)

	r := gin.New()
	r.GET("/me/reports/weekly", authMiddleware(5), h.GetMyWeeklyReport)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me/reports/weekly", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestActivityReport_GetMyMonthlyReport_Success(t *testing.T) {
	h, svc := setupActivityReportHandler()
	report := &model.ActivityReport{Period: model.ReportPeriodMonthly}
	svc.On("GetMonthlyReport", uint(5)).Return(report, nil)

	r := gin.New()
	r.GET("/me/reports/monthly", authMiddleware(5), h.GetMyMonthlyReport)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me/reports/monthly", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestActivityReport_GetComparison_Success(t *testing.T) {
	h, svc := setupActivityReportHandler()
	comp := &model.ReportComparison{ContributionsDiff: 5}
	svc.On("GetComparison", uint(3), model.ReportPeriodWeekly).Return(comp, nil)

	r := gin.New()
	r.GET("/me/reports/comparison", authMiddleware(3), h.GetComparison)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me/reports/comparison", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

func TestActivityReport_GetComparison_Monthly(t *testing.T) {
	h, svc := setupActivityReportHandler()
	comp := &model.ReportComparison{ContributionsDiff: -2}
	svc.On("GetComparison", uint(3), model.ReportPeriodMonthly).Return(comp, nil)

	r := gin.New()
	r.GET("/me/reports/comparison", authMiddleware(3), h.GetComparison)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me/reports/comparison?period=monthly", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	svc.AssertExpectations(t)
}

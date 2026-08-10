package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/mock"
)

// mockActivityReportRepo は usecase/repository.ActivityReportRepository のモック（ctx 付き）。
type mockActivityReportRepo struct{ mock.Mock }

func (m *mockActivityReportRepo) GetWeeklyReport(ctx context.Context, userID uint) (*model.ActivityReport, error) {
	args := m.Called(ctx, userID)
	r, _ := args.Get(0).(*model.ActivityReport)
	return r, args.Error(1)
}

func (m *mockActivityReportRepo) GetMonthlyReport(ctx context.Context, userID uint) (*model.ActivityReport, error) {
	args := m.Called(ctx, userID)
	r, _ := args.Get(0).(*model.ActivityReport)
	return r, args.Error(1)
}

func (m *mockActivityReportRepo) GetComparison(ctx context.Context, userID uint, period model.ReportPeriod) (*model.ReportComparison, error) {
	args := m.Called(ctx, userID, period)
	c, _ := args.Get(0).(*model.ReportComparison)
	return c, args.Error(1)
}

// setupActivityReportHandler は本物の usecase と port モックで ActivityReportHandler を組む。
func setupActivityReportHandler() (*ActivityReportHandler, *mockActivityReportRepo) {
	repo := new(mockActivityReportRepo)
	h := NewActivityReportHandler(
		usecase.NewGetWeeklyActivityReportUseCase(repo),
		usecase.NewGetMonthlyActivityReportUseCase(repo),
		usecase.NewGetActivityReportComparisonUseCase(repo),
	)
	return h, repo
}

func TestActivityReport_GetWeeklyReport_Success(t *testing.T) {
	h, repo := setupActivityReportHandler()
	report := &model.ActivityReport{Period: model.ReportPeriodWeekly, TotalContributions: 10}
	repo.On("GetWeeklyReport", mock.Anything, uint(1)).Return(report, nil)

	r := gin.New()
	r.GET("/users/:userId/reports/weekly", h.GetWeeklyReport)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1/reports/weekly", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
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
	h, repo := setupActivityReportHandler()
	repo.On("GetWeeklyReport", mock.Anything, uint(1)).Return((*model.ActivityReport)(nil), errors.New("db error"))

	r := gin.New()
	r.GET("/users/:userId/reports/weekly", h.GetWeeklyReport)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1/reports/weekly", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

func TestActivityReport_GetMonthlyReport_Success(t *testing.T) {
	h, repo := setupActivityReportHandler()
	report := &model.ActivityReport{Period: model.ReportPeriodMonthly, TotalContributions: 42}
	repo.On("GetMonthlyReport", mock.Anything, uint(1)).Return(report, nil)

	r := gin.New()
	r.GET("/users/:userId/reports/monthly", h.GetMonthlyReport)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1/reports/monthly", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestActivityReport_GetMyWeeklyReport_Success(t *testing.T) {
	h, repo := setupActivityReportHandler()
	report := &model.ActivityReport{Period: model.ReportPeriodWeekly}
	repo.On("GetWeeklyReport", mock.Anything, uint(5)).Return(report, nil)

	r := gin.New()
	r.GET("/me/reports/weekly", authMiddleware(5), h.GetMyWeeklyReport)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me/reports/weekly", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestActivityReport_GetMyMonthlyReport_Success(t *testing.T) {
	h, repo := setupActivityReportHandler()
	report := &model.ActivityReport{Period: model.ReportPeriodMonthly}
	repo.On("GetMonthlyReport", mock.Anything, uint(5)).Return(report, nil)

	r := gin.New()
	r.GET("/me/reports/monthly", authMiddleware(5), h.GetMyMonthlyReport)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me/reports/monthly", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

func TestActivityReport_GetComparison_Success(t *testing.T) {
	h, repo := setupActivityReportHandler()
	comp := &model.ReportComparison{ContributionsDiff: 5}
	repo.On("GetComparison", mock.Anything, uint(3), model.ReportPeriodWeekly).Return(comp, nil)

	r := gin.New()
	r.GET("/me/reports/comparison", authMiddleware(3), h.GetComparison)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me/reports/comparison", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

// ============================================================
// GetMonthlyReport テスト
// ============================================================

func TestActivityReport_GetMonthlyReport_ServiceError(t *testing.T) {
	h, repo := setupActivityReportHandler()
	repo.On("GetMonthlyReport", mock.Anything, uint(1)).Return((*model.ActivityReport)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/users/:userId/reports/monthly", h.GetMonthlyReport)

	w := doRequest(r, http.MethodGet, "/users/1/reports/monthly", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

func TestActivityReport_GetMonthlyReport_InvalidID(t *testing.T) {
	h, _ := setupActivityReportHandler()
	r := newRouter(1)
	r.GET("/users/:userId/reports/monthly", h.GetMonthlyReport)

	w := doRequest(r, http.MethodGet, "/users/abc/reports/monthly", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ============================================================
// GetMyWeeklyReport / GetMyMonthlyReport テスト
// ============================================================

func TestActivityReport_GetMyWeeklyReport_ServiceError(t *testing.T) {
	h, repo := setupActivityReportHandler()
	repo.On("GetWeeklyReport", mock.Anything, uint(1)).Return((*model.ActivityReport)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/me/reports/weekly", h.GetMyWeeklyReport)

	w := doRequest(r, http.MethodGet, "/me/reports/weekly", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

func TestActivityReport_GetMyMonthlyReport_ServiceError(t *testing.T) {
	h, repo := setupActivityReportHandler()
	repo.On("GetMonthlyReport", mock.Anything, uint(1)).Return((*model.ActivityReport)(nil), errors.New("db error"))

	r := newRouter(1)
	r.GET("/me/reports/monthly", h.GetMyMonthlyReport)

	w := doRequest(r, http.MethodGet, "/me/reports/monthly", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	repo.AssertExpectations(t)
}

func TestActivityReport_GetComparison_Monthly(t *testing.T) {
	h, repo := setupActivityReportHandler()
	comp := &model.ReportComparison{ContributionsDiff: -2}
	repo.On("GetComparison", mock.Anything, uint(3), model.ReportPeriodMonthly).Return(comp, nil)

	r := gin.New()
	r.GET("/me/reports/comparison", authMiddleware(3), h.GetComparison)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/me/reports/comparison?period=monthly", nil)
	r.ServeHTTP(w, req)

	assertStatus(t, w, http.StatusOK)
	repo.AssertExpectations(t)
}

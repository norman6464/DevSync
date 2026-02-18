package service

import (
	"errors"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestGetWeeklyReport_Success(t *testing.T) {
	repo := new(MockActivityReportRepository)
	svc := NewActivityReportService(repo)

	expected := &model.ActivityReport{
		Period:             model.ReportPeriodWeekly,
		UserID:             1,
		TotalContributions: 42,
		PostsCreated:       5,
		StartDate:          time.Now().AddDate(0, 0, -7),
		EndDate:            time.Now(),
	}
	repo.On("GetWeeklyReport", uint(1)).Return(expected, nil)

	result, err := svc.GetWeeklyReport(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetWeeklyReport_Error(t *testing.T) {
	repo := new(MockActivityReportRepository)
	svc := NewActivityReportService(repo)

	repo.On("GetWeeklyReport", uint(1)).Return(nil, errors.New("db error"))

	result, err := svc.GetWeeklyReport(1)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestGetMonthlyReport_Success(t *testing.T) {
	repo := new(MockActivityReportRepository)
	svc := NewActivityReportService(repo)

	expected := &model.ActivityReport{
		Period:             model.ReportPeriodMonthly,
		UserID:             2,
		TotalContributions: 150,
		PostsCreated:       20,
	}
	repo.On("GetMonthlyReport", uint(2)).Return(expected, nil)

	result, err := svc.GetMonthlyReport(2)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetMonthlyReport_Error(t *testing.T) {
	repo := new(MockActivityReportRepository)
	svc := NewActivityReportService(repo)

	repo.On("GetMonthlyReport", uint(2)).Return(nil, errors.New("db error"))

	result, err := svc.GetMonthlyReport(2)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestGetComparison_Weekly_Success(t *testing.T) {
	repo := new(MockActivityReportRepository)
	svc := NewActivityReportService(repo)

	expected := &model.ReportComparison{
		ContributionsDiff: 10,
		PostsDiff:         3,
		FollowersDiff:     2,
		GoalsDiff:         1,
		TrendPercentage:   15.5,
	}
	repo.On("GetComparison", uint(1), model.ReportPeriodWeekly).Return(expected, nil)

	result, err := svc.GetComparison(1, model.ReportPeriodWeekly)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetComparison_Monthly_Success(t *testing.T) {
	repo := new(MockActivityReportRepository)
	svc := NewActivityReportService(repo)

	expected := &model.ReportComparison{
		ContributionsDiff: -5,
		TrendPercentage:   -8.2,
	}
	repo.On("GetComparison", uint(3), model.ReportPeriodMonthly).Return(expected, nil)

	result, err := svc.GetComparison(3, model.ReportPeriodMonthly)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestGetComparison_Error(t *testing.T) {
	repo := new(MockActivityReportRepository)
	svc := NewActivityReportService(repo)

	repo.On("GetComparison", uint(1), model.ReportPeriodWeekly).Return(nil, errors.New("db error"))

	result, err := svc.GetComparison(1, model.ReportPeriodWeekly)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// エッジケーステスト
// ============================================================

func TestNewActivityReportService(t *testing.T) {
	repo := new(MockActivityReportRepository)
	svc := NewActivityReportService(repo)
	assert.NotNil(t, svc)
}

func TestGetWeeklyReport_EmptyReport(t *testing.T) {
	repo := new(MockActivityReportRepository)
	svc := NewActivityReportService(repo)

	expected := &model.ActivityReport{
		Period:             model.ReportPeriodWeekly,
		UserID:             1,
		TotalContributions: 0,
		PostsCreated:       0,
	}
	repo.On("GetWeeklyReport", uint(1)).Return(expected, nil)

	result, err := svc.GetWeeklyReport(1)
	assert.NoError(t, err)
	assert.Equal(t, 0, result.TotalContributions)
	assert.Equal(t, 0, result.PostsCreated)
	repo.AssertExpectations(t)
}

func TestGetMonthlyReport_LargeValues(t *testing.T) {
	repo := new(MockActivityReportRepository)
	svc := NewActivityReportService(repo)

	expected := &model.ActivityReport{
		Period:             model.ReportPeriodMonthly,
		UserID:             1,
		TotalContributions: 999999,
		PostsCreated:       10000,
	}
	repo.On("GetMonthlyReport", uint(1)).Return(expected, nil)

	result, err := svc.GetMonthlyReport(1)
	assert.NoError(t, err)
	assert.Equal(t, 999999, result.TotalContributions)
	assert.Equal(t, 10000, result.PostsCreated)
	repo.AssertExpectations(t)
}

func TestGetComparison_Monthly_Error(t *testing.T) {
	repo := new(MockActivityReportRepository)
	svc := NewActivityReportService(repo)

	repo.On("GetComparison", uint(1), model.ReportPeriodMonthly).Return(nil, errors.New("db error"))

	result, err := svc.GetComparison(1, model.ReportPeriodMonthly)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestGetComparison_ZeroDifferences(t *testing.T) {
	repo := new(MockActivityReportRepository)
	svc := NewActivityReportService(repo)

	expected := &model.ReportComparison{
		ContributionsDiff: 0,
		PostsDiff:         0,
		FollowersDiff:     0,
		GoalsDiff:         0,
		TrendPercentage:   0.0,
	}
	repo.On("GetComparison", uint(1), model.ReportPeriodWeekly).Return(expected, nil)

	result, err := svc.GetComparison(1, model.ReportPeriodWeekly)
	assert.NoError(t, err)
	assert.Equal(t, 0, result.ContributionsDiff)
	assert.Equal(t, 0.0, result.TrendPercentage)
	repo.AssertExpectations(t)
}

func TestGetComparison_LargeNegativeTrend(t *testing.T) {
	repo := new(MockActivityReportRepository)
	svc := NewActivityReportService(repo)

	expected := &model.ReportComparison{
		ContributionsDiff: -100,
		PostsDiff:         -50,
		TrendPercentage:   -99.9,
	}
	repo.On("GetComparison", uint(1), model.ReportPeriodMonthly).Return(expected, nil)

	result, err := svc.GetComparison(1, model.ReportPeriodMonthly)
	assert.NoError(t, err)
	assert.Equal(t, -100, result.ContributionsDiff)
	assert.Equal(t, -99.9, result.TrendPercentage)
	repo.AssertExpectations(t)
}

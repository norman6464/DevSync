package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockActivityReportRepo は usecase/repository.ActivityReportRepository のモック。
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

func TestGetWeeklyActivityReportUseCase_Execute(t *testing.T) {
	t.Run("週次レポートを返す", func(t *testing.T) {
		repo := new(mockActivityReportRepo)
		expected := &model.ActivityReport{Period: model.ReportPeriodWeekly, TotalContributions: 10}
		repo.On("GetWeeklyReport", mock.Anything, uint(7)).Return(expected, nil)
		uc := usecase.NewGetWeeklyActivityReportUseCase(repo)

		got, err := uc.Execute(context.Background(), 7)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		repo.AssertExpectations(t)
	})

	t.Run("repo のエラーをそのまま返す", func(t *testing.T) {
		repo := new(mockActivityReportRepo)
		repo.On("GetWeeklyReport", mock.Anything, uint(7)).
			Return((*model.ActivityReport)(nil), errors.New("db error"))
		uc := usecase.NewGetWeeklyActivityReportUseCase(repo)

		_, err := uc.Execute(context.Background(), 7)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestGetMonthlyActivityReportUseCase_Execute(t *testing.T) {
	t.Run("月次レポートを返す", func(t *testing.T) {
		repo := new(mockActivityReportRepo)
		expected := &model.ActivityReport{Period: model.ReportPeriodMonthly, TotalContributions: 42}
		repo.On("GetMonthlyReport", mock.Anything, uint(7)).Return(expected, nil)
		uc := usecase.NewGetMonthlyActivityReportUseCase(repo)

		got, err := uc.Execute(context.Background(), 7)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		repo.AssertExpectations(t)
	})

	t.Run("repo のエラーをそのまま返す", func(t *testing.T) {
		repo := new(mockActivityReportRepo)
		repo.On("GetMonthlyReport", mock.Anything, uint(7)).
			Return((*model.ActivityReport)(nil), errors.New("db error"))
		uc := usecase.NewGetMonthlyActivityReportUseCase(repo)

		_, err := uc.Execute(context.Background(), 7)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestGetActivityReportComparisonUseCase_Execute(t *testing.T) {
	t.Run("指定した期間区分をそのまま repo へ渡す", func(t *testing.T) {
		repo := new(mockActivityReportRepo)
		expected := &model.ReportComparison{ContributionsDiff: 5, TrendPercentage: 12.5}
		repo.On("GetComparison", mock.Anything, uint(7), model.ReportPeriodMonthly).Return(expected, nil)
		uc := usecase.NewGetActivityReportComparisonUseCase(repo)

		got, err := uc.Execute(context.Background(), 7, model.ReportPeriodMonthly)

		assert.NoError(t, err)
		assert.Equal(t, expected, got)
		repo.AssertExpectations(t)
	})

	t.Run("repo のエラーをそのまま返す", func(t *testing.T) {
		repo := new(mockActivityReportRepo)
		repo.On("GetComparison", mock.Anything, uint(7), model.ReportPeriodWeekly).
			Return((*model.ReportComparison)(nil), errors.New("db error"))
		uc := usecase.NewGetActivityReportComparisonUseCase(repo)

		_, err := uc.Execute(context.Background(), 7, model.ReportPeriodWeekly)

		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

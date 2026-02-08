package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// ActivityReportService handles activity report business logic.
type ActivityReportService struct {
	repo repository.ActivityReportRepositoryInterface
}

// NewActivityReportService creates a new ActivityReportService.
func NewActivityReportService(repo repository.ActivityReportRepositoryInterface) *ActivityReportService {
	return &ActivityReportService{repo: repo}
}

// GetWeeklyReport returns the weekly activity report for a user.
func (s *ActivityReportService) GetWeeklyReport(userID uint) (*model.ActivityReport, error) {
	return s.repo.GetWeeklyReport(userID)
}

// GetMonthlyReport returns the monthly activity report for a user.
func (s *ActivityReportService) GetMonthlyReport(userID uint) (*model.ActivityReport, error) {
	return s.repo.GetMonthlyReport(userID)
}

// GetComparison returns the comparison between current and previous period.
func (s *ActivityReportService) GetComparison(userID uint, period model.ReportPeriod) (*model.ReportComparison, error) {
	return s.repo.GetComparison(userID, period)
}

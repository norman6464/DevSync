package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// ActivityReportService はアクティビティレポートのビジネスロジックを提供する。
// リポジトリ層に処理を委譲する薄いラッパー。
type ActivityReportService struct {
	repo repository.ActivityReportRepositoryInterface
}

// NewActivityReportService は新しいActivityReportServiceインスタンスを生成する。
func NewActivityReportService(repo repository.ActivityReportRepositoryInterface) *ActivityReportService {
	return &ActivityReportService{repo: repo}
}

// GetWeeklyReport は指定ユーザーの週次アクティビティレポートを返す。
func (s *ActivityReportService) GetWeeklyReport(userID uint) (*model.ActivityReport, error) {
	return s.repo.GetWeeklyReport(userID)
}

// GetMonthlyReport は指定ユーザーの月次アクティビティレポートを返す。
func (s *ActivityReportService) GetMonthlyReport(userID uint) (*model.ActivityReport, error) {
	return s.repo.GetMonthlyReport(userID)
}

// GetComparison は現在期間と前期間の比較データを返す。
func (s *ActivityReportService) GetComparison(userID uint, period model.ReportPeriod) (*model.ReportComparison, error) {
	return s.repo.GetComparison(userID, period)
}

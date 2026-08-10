package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetActivityReportComparisonUseCase は現在期間と前期間のアクティビティ比較を取得する。
type GetActivityReportComparisonUseCase struct {
	reports repository.ActivityReportRepository
}

// NewGetActivityReportComparisonUseCase は GetActivityReportComparisonUseCase を生成する。
func NewGetActivityReportComparisonUseCase(reports repository.ActivityReportRepository) *GetActivityReportComparisonUseCase {
	return &GetActivityReportComparisonUseCase{reports: reports}
}

// Execute は指定期間区分での比較データを返す。
func (uc *GetActivityReportComparisonUseCase) Execute(ctx context.Context, userID uint, period model.ReportPeriod) (*model.ReportComparison, error) {
	return uc.reports.GetComparison(ctx, userID, period)
}

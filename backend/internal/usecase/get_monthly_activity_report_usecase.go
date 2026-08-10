package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetMonthlyActivityReportUseCase は指定ユーザーの月次アクティビティレポートを取得する。
type GetMonthlyActivityReportUseCase struct {
	reports repository.ActivityReportRepository
}

// NewGetMonthlyActivityReportUseCase は GetMonthlyActivityReportUseCase を生成する。
func NewGetMonthlyActivityReportUseCase(reports repository.ActivityReportRepository) *GetMonthlyActivityReportUseCase {
	return &GetMonthlyActivityReportUseCase{reports: reports}
}

// Execute は月次アクティビティレポートを返す。
func (uc *GetMonthlyActivityReportUseCase) Execute(ctx context.Context, userID uint) (*model.ActivityReport, error) {
	return uc.reports.GetMonthlyReport(ctx, userID)
}

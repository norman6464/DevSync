package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// GetWeeklyActivityReportUseCase は指定ユーザーの週次アクティビティレポートを取得する。
type GetWeeklyActivityReportUseCase struct {
	reports repository.ActivityReportRepository
}

// NewGetWeeklyActivityReportUseCase は GetWeeklyActivityReportUseCase を生成する。
func NewGetWeeklyActivityReportUseCase(reports repository.ActivityReportRepository) *GetWeeklyActivityReportUseCase {
	return &GetWeeklyActivityReportUseCase{reports: reports}
}

// Execute は週次アクティビティレポートを返す。
func (uc *GetWeeklyActivityReportUseCase) Execute(ctx context.Context, userID uint) (*model.ActivityReport, error) {
	return uc.reports.GetWeeklyReport(ctx, userID)
}

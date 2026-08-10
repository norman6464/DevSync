package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// WeeklyActivityReportReader は週次レポートの取得だけを要求する最小の契約。
// 週次レポートメールの一括送信のように、レポート生成の一部しか必要としない消費者が依存する。
type WeeklyActivityReportReader interface {
	GetWeeklyReport(ctx context.Context, userID uint) (*model.ActivityReport, error)
}

// ActivityReportRepository はアクティビティレポートの生成に対する、usecase 側が要求する契約。
type ActivityReportRepository interface {
	WeeklyActivityReportReader

	GetMonthlyReport(ctx context.Context, userID uint) (*model.ActivityReport, error)
	GetComparison(ctx context.Context, userID uint, period model.ReportPeriod) (*model.ReportComparison, error)
}

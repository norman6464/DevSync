package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// RoadmapRepository は学習ロードマップの永続化に対する、usecase 側が要求する契約。
type RoadmapRepository interface {
	Create(ctx context.Context, roadmap *model.Roadmap) error
	// Update はタイトル・説明・カテゴリ・公開設定・テンプレート設定を上書きする。
	// status/completed_at は対象外（UpdateStatus を使う）。ステップ完了による
	// 自動遷移（recalcRoadmapStatusFromMetrics）と経路を分け、ロストアップデートを防ぐため。
	Update(ctx context.Context, roadmap *model.Roadmap) error
	// UpdateStatus はユーザーによる明示的なステータス変更だけに使う専用の更新。
	UpdateStatus(ctx context.Context, roadmap *model.Roadmap) error
	Delete(ctx context.Context, id uint) error
	// FindByID はステップ（表示順）とユーザーを含めてロードマップを返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindByID(ctx context.Context, id uint) (*model.Roadmap, error)

	GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Roadmap, int64, error)
	GetByStatus(ctx context.Context, userID uint, status string) ([]model.Roadmap, error)
	GetPublicRoadmaps(ctx context.Context, limit, offset int) ([]model.Roadmap, int64, error)
	GetTemplates(ctx context.Context) ([]model.Roadmap, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)

	// CopyRoadmap は元のロードマップとそのステップを複製する。
	// 複製は非公開・アクティブで作られ、ステップの完了状態は引き継がない。
	CopyRoadmap(ctx context.Context, originalID, newUserID uint) (*model.Roadmap, error)

	// CreateStep はステップを追加し、ロードマップのステップ数を 1 増やす。
	CreateStep(ctx context.Context, step *model.RoadmapStep) error
	// UpdateStep はステップを更新する。完了状態が変わった場合はロードマップの完了
	// ステップ数を増減し、進捗率（ステップ数から都度算出される）が100%を跨いだときは
	// ステータスも自動遷移させる。
	UpdateStep(ctx context.Context, step *model.RoadmapStep) error
	// DeleteStep はステップを削除し、ロードマップのステップ数・完了ステップ数を
	// 増減する。UpdateStepと異なりステータスの自動遷移は行わない。
	DeleteStep(ctx context.Context, stepID uint) error
	// FindStepByID は指定 ID のステップを返す。不在の場合は (nil, nil) を返す。
	FindStepByID(ctx context.Context, stepID uint) (*model.RoadmapStep, error)
	ReorderSteps(ctx context.Context, roadmapID uint, stepOrders []model.StepOrder) error
}

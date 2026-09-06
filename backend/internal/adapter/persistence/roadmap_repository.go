package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// roadmapRepository は [repository.RoadmapRepository] の sqlc(pgx) 実装。
// CopyRoadmap（ロードマップ＋ステップの複製）や CreateStep/UpdateStep/DeleteStep/
// ReorderSteps（ステップ変更とロードマップ集計列の同時更新）を1トランザクションで
// 行うため、Queries だけでなくトランザクションを開始できる *pgxpool.Pool を直接保持する。
type roadmapRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewRoadmapRepository は RoadmapRepository の sqlc(pgx) 実装を返す。
func NewRoadmapRepository(pool *pgxpool.Pool) repository.RoadmapRepository {
	return &roadmapRepository{pool: pool, q: sqlcgen.New(pool)}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.RoadmapRepository = (*roadmapRepository)(nil)

func toModelRoadmap(row sqlcgen.Roadmap) model.Roadmap {
	return model.Roadmap{
		ID:                 uint(row.ID),
		UserID:             uint(row.UserID),
		Title:              row.Title,
		Description:        fromStringPtr(row.Description),
		Category:           model.RoadmapCategory(fromStringPtr(row.Category)),
		IsPublic:           fromBoolPtr(row.IsPublic),
		IsTemplate:         fromBoolPtr(row.IsTemplate),
		StepCount:          int(fromInt64PtrValue(row.StepCount)),
		CompletedStepCount: int(fromInt64PtrValue(row.CompletedStepCount)),
		Progress:           int(fromInt64PtrValue(row.Progress)),
		Status:             model.RoadmapStatus(fromStringPtr(row.Status)),
		CreatedAt:          timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:          timeValue(fromTimestamptz(row.UpdatedAt)),
		CompletedAt:        fromTimestamptz(row.CompletedAt),
	}
}

func toModelRoadmapStep(row sqlcgen.RoadmapStep) model.RoadmapStep {
	return model.RoadmapStep{
		ID:          uint(row.ID),
		RoadmapID:   uint(row.RoadmapID),
		Title:       row.Title,
		Description: fromStringPtr(row.Description),
		OrderIndex:  int(row.OrderIndex),
		IsCompleted: fromBoolPtr(row.IsCompleted),
		CompletedAt: fromTimestamptz(row.CompletedAt),
		ResourceURL: fromStringPtr(row.ResourceUrl),
		CreatedAt:   timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:   timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

func toModelRoadmapSteps(rows []sqlcgen.RoadmapStep) []model.RoadmapStep {
	steps := make([]model.RoadmapStep, len(rows))
	for i, row := range rows {
		steps[i] = toModelRoadmapStep(row)
	}
	return steps
}

// Create は新しいロードマップを作成する。
// Category/Statusは移行前のGORMの `gorm:"default:..."` に相当し、
// 未指定（ゼロ値）ならそれぞれ other / active を補う。
func (r *roadmapRepository) Create(ctx context.Context, roadmap *model.Roadmap) error {
	category := roadmap.Category
	if category == "" {
		category = model.RoadmapCategoryOther
	}
	status := roadmap.Status
	if status == "" {
		status = model.RoadmapStatusActive
	}

	row, err := r.q.CreateRoadmap(ctx, sqlcgen.CreateRoadmapParams{
		UserID:             int64(roadmap.UserID),
		Title:              roadmap.Title,
		Description:        &roadmap.Description,
		Category:           (*string)(&category),
		IsPublic:           &roadmap.IsPublic,
		IsTemplate:         &roadmap.IsTemplate,
		StepCount:          toInt64Ptr(roadmap.StepCount),
		CompletedStepCount: toInt64Ptr(roadmap.CompletedStepCount),
		Progress:           toInt64Ptr(roadmap.Progress),
		Status:             (*string)(&status),
	})
	if err != nil {
		return err
	}
	*roadmap = toModelRoadmap(row)
	return nil
}

// Update は既存のロードマップを更新する（GORMのSave＝全カラム上書きに相当）。
func (r *roadmapRepository) Update(ctx context.Context, roadmap *model.Roadmap) error {
	row, err := r.q.UpdateRoadmap(ctx, sqlcgen.UpdateRoadmapParams{
		ID:          int64(roadmap.ID),
		Title:       roadmap.Title,
		Description: &roadmap.Description,
		Category:    (*string)(&roadmap.Category),
		IsPublic:    &roadmap.IsPublic,
		IsTemplate:  &roadmap.IsTemplate,
		Status:      (*string)(&roadmap.Status),
		CompletedAt: toTimestamptz(roadmap.CompletedAt),
	})
	if err != nil {
		return err
	}
	*roadmap = toModelRoadmap(row)
	return nil
}

// Delete はロードマップを削除する（ステップは FK の ON DELETE CASCADE で消える）。
func (r *roadmapRepository) Delete(ctx context.Context, id uint) error {
	return r.q.DeleteRoadmap(ctx, int64(id))
}

// FindByID はステップ（表示順）とユーザーを含めてロードマップを取得する。
// 不在の場合は (nil, nil) を返す。
func (r *roadmapRepository) FindByID(ctx context.Context, id uint) (*model.Roadmap, error) {
	row, err := r.q.GetRoadmapWithUserByID(ctx, int64(id))
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	roadmap := toModelRoadmap(row.Roadmap)
	roadmap.User = toModelUser(row.User)

	steps, err := r.q.ListRoadmapStepsByRoadmap(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	roadmap.Steps = toModelRoadmapSteps(steps)

	return &roadmap, nil
}

// GetByUserID は指定ユーザーのロードマップをページネーション付きで取得する（新しい順・ステップなし）。
func (r *roadmapRepository) GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Roadmap, int64, error) {
	total, err := r.q.CountRoadmapsByUser(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.q.ListRoadmapsByUser(ctx, sqlcgen.ListRoadmapsByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	roadmaps := make([]model.Roadmap, len(rows))
	for i, row := range rows {
		roadmaps[i] = toModelRoadmap(row)
	}
	return roadmaps, total, nil
}

// GetByStatus は指定ユーザーのロードマップをステータスで絞り込んで取得する（新しい順）。
func (r *roadmapRepository) GetByStatus(ctx context.Context, userID uint, status string) ([]model.Roadmap, error) {
	rows, err := r.q.ListRoadmapsByStatus(ctx, sqlcgen.ListRoadmapsByStatusParams{
		UserID: int64(userID),
		Status: &status,
	})
	if err != nil {
		return nil, err
	}
	roadmaps := make([]model.Roadmap, len(rows))
	for i, row := range rows {
		roadmaps[i] = toModelRoadmap(row)
	}
	return roadmaps, nil
}

// GetPublicRoadmaps は公開ロードマップをページネーション付きで取得する（新しい順）。
func (r *roadmapRepository) GetPublicRoadmaps(ctx context.Context, limit, offset int) ([]model.Roadmap, int64, error) {
	total, err := r.q.CountPublicRoadmaps(ctx)
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.q.ListPublicRoadmapsWithUser(ctx, sqlcgen.ListPublicRoadmapsWithUserParams{
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	roadmaps := make([]model.Roadmap, len(rows))
	for i, row := range rows {
		roadmaps[i] = toModelRoadmap(row.Roadmap)
		roadmaps[i].User = toModelUser(row.User)
	}
	return roadmaps, total, nil
}

// GetTemplates はテンプレートのロードマップをステップ付きで取得する（古い順）。
func (r *roadmapRepository) GetTemplates(ctx context.Context) ([]model.Roadmap, error) {
	rows, err := r.q.ListTemplateRoadmaps(ctx)
	if err != nil {
		return nil, err
	}

	roadmaps := make([]model.Roadmap, len(rows))
	roadmapIDs := make([]int64, len(rows))
	for i, row := range rows {
		roadmaps[i] = toModelRoadmap(row)
		roadmapIDs[i] = row.ID
	}

	if len(roadmapIDs) > 0 {
		stepRows, err := r.q.ListRoadmapStepsByRoadmapIDs(ctx, roadmapIDs)
		if err != nil {
			return nil, err
		}
		stepsByRoadmapID := make(map[uint][]model.RoadmapStep)
		for _, row := range stepRows {
			roadmapID := uint(row.RoadmapID)
			stepsByRoadmapID[roadmapID] = append(stepsByRoadmapID[roadmapID], toModelRoadmapStep(row))
		}
		for i := range roadmaps {
			roadmaps[i].Steps = stepsByRoadmapID[roadmaps[i].ID]
		}
	}

	return roadmaps, nil
}

// CountByUserID は指定ユーザーのロードマップ総数を返す。
func (r *roadmapRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountRoadmapsByUser(ctx, int64(userID))
}

// CopyRoadmap は元のロードマップとステップを複製する。複製は非公開・アクティブで作られる。
// ロードマップ本体とステップ群の作成を1トランザクションで行い、
// 途中で失敗したときに不完全な複製が残らないようにする
// （移行前のGORM実装はトランザクションで括っていなかったが、複数行の作成が
// 部分的に失敗すると不整合な複製が残るため、安全側に倒して括った）。
func (r *roadmapRepository) CopyRoadmap(ctx context.Context, originalID, newUserID uint) (*model.Roadmap, error) {
	original, err := r.FindByID(ctx, originalID)
	if err != nil {
		return nil, err
	}
	if original == nil {
		// 不在は (nil, nil)。エラーは DB 障害だけを表す。
		return nil, nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)

	status := model.RoadmapStatusActive
	newRoadmapRow, err := q.CreateRoadmap(ctx, sqlcgen.CreateRoadmapParams{
		UserID:      int64(newUserID),
		Title:       original.Title + " (コピー)",
		Description: &original.Description,
		Category:    (*string)(&original.Category),
		IsPublic:    new(bool),
		IsTemplate:  new(bool),
		StepCount:   toInt64Ptr(original.StepCount),
		Status:      (*string)(&status),
	})
	if err != nil {
		return nil, err
	}
	newRoadmapID := newRoadmapRow.ID

	for _, step := range original.Steps {
		if _, err := q.CreateRoadmapStep(ctx, sqlcgen.CreateRoadmapStepParams{
			RoadmapID:   newRoadmapID,
			Title:       step.Title,
			Description: &step.Description,
			OrderIndex:  int64(step.OrderIndex),
			ResourceUrl: &step.ResourceURL,
		}); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.FindByID(ctx, uint(newRoadmapID))
}

// CreateStep はステップを追加し、ロードマップのステップ数を 1 増やす。
func (r *roadmapRepository) CreateStep(ctx context.Context, step *model.RoadmapStep) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	row, err := q.CreateRoadmapStep(ctx, sqlcgen.CreateRoadmapStepParams{
		RoadmapID:   int64(step.RoadmapID),
		Title:       step.Title,
		Description: &step.Description,
		OrderIndex:  int64(step.OrderIndex),
		IsCompleted: &step.IsCompleted,
		CompletedAt: toTimestamptz(step.CompletedAt),
		ResourceUrl: &step.ResourceURL,
	})
	if err != nil {
		return err
	}
	if err := q.IncrementRoadmapStepCount(ctx, int64(step.RoadmapID)); err != nil {
		return err
	}
	*step = toModelRoadmapStep(row)
	return tx.Commit(ctx)
}

// UpdateStep はステップを更新する。完了状態が変わった場合は進捗も再計算する。
func (r *roadmapRepository) UpdateStep(ctx context.Context, step *model.RoadmapStep) error {
	oldStepRow, err := r.q.GetRoadmapStepByID(ctx, int64(step.ID))
	if err != nil {
		return err
	}
	oldIsCompleted := fromBoolPtr(oldStepRow.IsCompleted)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	row, err := q.UpdateRoadmapStep(ctx, sqlcgen.UpdateRoadmapStepParams{
		ID:          int64(step.ID),
		Title:       step.Title,
		Description: &step.Description,
		OrderIndex:  int64(step.OrderIndex),
		IsCompleted: &step.IsCompleted,
		CompletedAt: toTimestamptz(step.CompletedAt),
		ResourceUrl: &step.ResourceURL,
	})
	if err != nil {
		return err
	}
	*step = toModelRoadmapStep(row)

	if oldIsCompleted == step.IsCompleted {
		return tx.Commit(ctx)
	}

	delta := int64(1)
	if !step.IsCompleted {
		delta = -1
	}
	if err := q.AdjustRoadmapCompletedStepCount(ctx, sqlcgen.AdjustRoadmapCompletedStepCountParams{
		ID:    int64(step.RoadmapID),
		Delta: delta,
	}); err != nil {
		return err
	}
	if err := recalcRoadmapProgress(ctx, q, int64(step.RoadmapID), true); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// DeleteStep はステップを削除し、ステップ数・完了ステップ数・進捗率を再計算する。
func (r *roadmapRepository) DeleteStep(ctx context.Context, stepID uint) error {
	stepRow, err := r.q.GetRoadmapStepByID(ctx, int64(stepID))
	if err != nil {
		return err
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	if err := q.DeleteRoadmapStep(ctx, int64(stepID)); err != nil {
		return err
	}
	if err := q.DecrementRoadmapStepCount(ctx, stepRow.RoadmapID); err != nil {
		return err
	}
	if fromBoolPtr(stepRow.IsCompleted) {
		if err := q.AdjustRoadmapCompletedStepCount(ctx, sqlcgen.AdjustRoadmapCompletedStepCountParams{
			ID:    stepRow.RoadmapID,
			Delta: -1,
		}); err != nil {
			return err
		}
	}
	if err := recalcRoadmapProgress(ctx, q, stepRow.RoadmapID, false); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// recalcRoadmapProgress は進捗率を再計算して保存する。
// withStatus が true のときは 100% 到達での自動完了と、100% 未満へ戻ったときの
// アクティブ復帰も行う（ステップ更新のみが対象で、削除時は進捗率だけを更新する）。
func recalcRoadmapProgress(ctx context.Context, q *sqlcgen.Queries, roadmapID int64, withStatus bool) error {
	roadmapRow, err := q.GetRoadmapByID(ctx, roadmapID)
	if err != nil {
		return err
	}
	stepCount := fromInt64PtrValue(roadmapRow.StepCount)
	completedStepCount := fromInt64PtrValue(roadmapRow.CompletedStepCount)

	progress := int64(0)
	if stepCount > 0 {
		progress = (completedStepCount * 100) / stepCount
	}

	if !withStatus {
		return q.UpdateRoadmapProgress(ctx, sqlcgen.UpdateRoadmapProgressParams{
			ID:       roadmapID,
			Progress: &progress,
		})
	}

	status := model.RoadmapStatus(fromStringPtr(roadmapRow.Status))
	switch {
	case progress == 100 && status == model.RoadmapStatusActive:
		return q.UpdateRoadmapProgressCompleted(ctx, sqlcgen.UpdateRoadmapProgressCompletedParams{
			ID:       roadmapID,
			Progress: &progress,
		})
	case progress < 100 && status == model.RoadmapStatusCompleted:
		return q.UpdateRoadmapProgressReactivated(ctx, sqlcgen.UpdateRoadmapProgressReactivatedParams{
			ID:       roadmapID,
			Progress: &progress,
		})
	default:
		return q.UpdateRoadmapProgress(ctx, sqlcgen.UpdateRoadmapProgressParams{
			ID:       roadmapID,
			Progress: &progress,
		})
	}
}

// FindStepByID は指定 ID のステップを取得する。不在の場合は (nil, nil) を返す。
func (r *roadmapRepository) FindStepByID(ctx context.Context, stepID uint) (*model.RoadmapStep, error) {
	row, err := r.q.GetRoadmapStepByID(ctx, int64(stepID))
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	step := toModelRoadmapStep(row)
	return &step, nil
}

// ReorderSteps はステップの表示順序をまとめて更新する。
func (r *roadmapRepository) ReorderSteps(ctx context.Context, roadmapID uint, stepOrders []model.StepOrder) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	for _, order := range stepOrders {
		if err := q.ReorderRoadmapStep(ctx, sqlcgen.ReorderRoadmapStepParams{
			ID:         int64(order.StepID),
			RoadmapID:  int64(roadmapID),
			OrderIndex: int64(order.OrderIndex),
		}); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

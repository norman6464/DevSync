package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// roadmapRepository は [repository.RoadmapRepository] の sqlc(pgx) 実装。
// CopyRoadmap（ロードマップ＋ステップの複製）を1トランザクションで行うため、
// Queries だけでなくトランザクションを開始できる *pgxpool.Pool を直接保持する。
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
		ID:          uint(row.ID),
		UserID:      uint(row.UserID),
		Title:       row.Title,
		Description: fromStringPtr(row.Description),
		Category:    model.RoadmapCategory(fromStringPtr(row.Category)),
		IsPublic:    row.IsPublic,
		IsTemplate:  row.IsTemplate,
		Status:      model.RoadmapStatus(fromStringPtr(row.Status)),
		CreatedAt:   timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:   timeValue(fromTimestamptz(row.UpdatedAt)),
		CompletedAt: fromTimestamptz(row.CompletedAt),
	}
}

func toModelRoadmapStep(row sqlcgen.RoadmapStep) model.RoadmapStep {
	return model.RoadmapStep{
		ID:          uint(row.ID),
		RoadmapID:   uint(row.RoadmapID),
		Title:       row.Title,
		Description: fromStringPtr(row.Description),
		OrderIndex:  int(row.OrderIndex),
		IsCompleted: row.IsCompleted,
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

// attachRoadmapMetrics は複数のロードマップへstep_count/completed_step_count/progressを
// まとめて取得して付与する（DEVSYNC-159でroadmap_metrics側テーブルへ分離済み）。
// progressは列として持たず、GetRoadmapMetricsByRoadmapIDsのSELECT側の算出式で都度
// 計算される。1件もステップが無いロードマップはroadmap_metrics行が存在しないため
// （遅延生成）0のまま。
func attachRoadmapMetrics(ctx context.Context, q *sqlcgen.Queries, roadmaps []model.Roadmap) error {
	if len(roadmaps) == 0 {
		return nil
	}
	roadmapIDs := make([]int64, len(roadmaps))
	for i, roadmap := range roadmaps {
		roadmapIDs[i] = int64(roadmap.ID)
	}

	metricsRows, err := q.GetRoadmapMetricsByRoadmapIDs(ctx, roadmapIDs)
	if err != nil {
		return err
	}
	metricsByRoadmapID := make(map[uint]sqlcgen.GetRoadmapMetricsByRoadmapIDsRow, len(metricsRows))
	for _, row := range metricsRows {
		metricsByRoadmapID[uint(row.RoadmapID)] = row
	}
	for i := range roadmaps {
		m := metricsByRoadmapID[roadmaps[i].ID]
		roadmaps[i].StepCount = int(m.StepCount)
		roadmaps[i].CompletedStepCount = int(m.CompletedStepCount)
		roadmaps[i].Progress = int(m.Progress)
	}
	return nil
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
		UserID:      int64(roadmap.UserID),
		Title:       roadmap.Title,
		Description: &roadmap.Description,
		Category:    (*string)(&category),
		IsPublic:    roadmap.IsPublic,
		IsTemplate:  roadmap.IsTemplate,
		Status:      (*string)(&status),
	})
	if err != nil {
		return err
	}
	*roadmap = toModelRoadmap(row)
	return nil
}

// Update は既存のロードマップを更新する（GORMのSave＝全カラム上書きに相当）。
// status/completed_atは対象外（UpdateStatusを使う）。
func (r *roadmapRepository) Update(ctx context.Context, roadmap *model.Roadmap) error {
	row, err := r.q.UpdateRoadmap(ctx, sqlcgen.UpdateRoadmapParams{
		ID:          int64(roadmap.ID),
		Title:       roadmap.Title,
		Description: &roadmap.Description,
		Category:    (*string)(&roadmap.Category),
		IsPublic:    roadmap.IsPublic,
		IsTemplate:  roadmap.IsTemplate,
	})
	if err != nil {
		return err
	}
	*roadmap = toModelRoadmap(row)
	roadmaps := []model.Roadmap{*roadmap}
	if err := attachRoadmapMetrics(ctx, r.q, roadmaps); err != nil {
		return err
	}
	*roadmap = roadmaps[0]
	return nil
}

// UpdateStatus はステータスと完了日時だけを更新する
// （ユーザーによる明示的なステータス変更専用。他のフィールドは触らない）。
func (r *roadmapRepository) UpdateStatus(ctx context.Context, roadmap *model.Roadmap) error {
	row, err := r.q.UpdateRoadmapStatus(ctx, sqlcgen.UpdateRoadmapStatusParams{
		ID:          int64(roadmap.ID),
		Status:      (*string)(&roadmap.Status),
		CompletedAt: toTimestamptz(roadmap.CompletedAt),
	})
	if err != nil {
		return err
	}
	*roadmap = toModelRoadmap(row)
	roadmaps := []model.Roadmap{*roadmap}
	if err := attachRoadmapMetrics(ctx, r.q, roadmaps); err != nil {
		return err
	}
	*roadmap = roadmaps[0]
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

	roadmaps := []model.Roadmap{roadmap}
	if err := attachRoadmapMetrics(ctx, r.q, roadmaps); err != nil {
		return nil, err
	}
	return &roadmaps[0], nil
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
	if err := attachRoadmapMetrics(ctx, r.q, roadmaps); err != nil {
		return nil, 0, err
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
	if err := attachRoadmapMetrics(ctx, r.q, roadmaps); err != nil {
		return nil, err
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
	if err := attachRoadmapMetrics(ctx, r.q, roadmaps); err != nil {
		return nil, 0, err
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

	if err := attachRoadmapMetrics(ctx, r.q, roadmaps); err != nil {
		return nil, err
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
// CreateRoadmapStepは1呼び出しごとにroadmap_metrics.step_countを加算するCTEを含むため
// （DEVSYNC-159）、ここでstep_countを別途シードする必要はない。
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
		IsPublic:    false,
		IsTemplate:  false,
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

// CreateStep はステップを追加し、roadmap_metrics.step_countを加算する
// （同一SQL文、DEVSYNC-159）。
func (r *roadmapRepository) CreateStep(ctx context.Context, step *model.RoadmapStep) error {
	row, err := r.q.CreateRoadmapStep(ctx, sqlcgen.CreateRoadmapStepParams{
		RoadmapID:   int64(step.RoadmapID),
		Title:       step.Title,
		Description: &step.Description,
		OrderIndex:  int64(step.OrderIndex),
		IsCompleted: step.IsCompleted,
		CompletedAt: toTimestamptz(step.CompletedAt),
		ResourceUrl: &step.ResourceURL,
	})
	if err != nil {
		return err
	}
	*step = toModelRoadmapStep(sqlcgen.RoadmapStep(row))
	return nil
}

// UpdateStep はステップを更新する。完了状態が変わった場合はroadmap_metrics.
// completed_step_countも同一SQL文で加減算する。旧値の読み取り（GetRoadmapStepByIDForUpdate）
// でFOR UPDATEにより対象行をロックし、そのままUPDATE・ステータス自動遷移まで
// 1トランザクションで括ることで、同じステップへの同時更新がcompleted_step_countを
// 二重に加減算したり、途中失敗でステータスだけ古いまま残ったりしないようにする。
func (r *roadmapRepository) UpdateStep(ctx context.Context, step *model.RoadmapStep) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	oldStepRow, err := q.GetRoadmapStepByIDForUpdate(ctx, int64(step.ID))
	if err != nil {
		return err
	}
	oldIsCompleted := oldStepRow.IsCompleted

	var completedDelta int64
	switch {
	case step.IsCompleted && !oldIsCompleted:
		completedDelta = 1
	case !step.IsCompleted && oldIsCompleted:
		completedDelta = -1
	}

	row, err := q.UpdateRoadmapStep(ctx, sqlcgen.UpdateRoadmapStepParams{
		ID:             int64(step.ID),
		Title:          step.Title,
		Description:    &step.Description,
		OrderIndex:     int64(step.OrderIndex),
		IsCompleted:    step.IsCompleted,
		CompletedAt:    toTimestamptz(step.CompletedAt),
		ResourceUrl:    &step.ResourceURL,
		CompletedDelta: completedDelta,
	})
	if err != nil {
		return err
	}
	*step = toModelRoadmapStep(sqlcgen.RoadmapStep(row))

	if completedDelta != 0 {
		if err := recalcRoadmapStatusFromMetrics(ctx, q, step.RoadmapID); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// DeleteStep はステップを削除し、roadmap_metrics.step_count/completed_step_countを
// 同一SQL文で調整する（削除されたステップの完了状態はDELETEのRETURNING自体から
// 取得するため、削除前の別読み取りは不要）。削除では移行前と同じくステータスの
// 自動遷移は行わない。
func (r *roadmapRepository) DeleteStep(ctx context.Context, stepID uint) error {
	return r.q.DeleteRoadmapStep(ctx, int64(stepID))
}

// recalcRoadmapStatusFromMetrics はroadmap_metricsから算出される進捗率が100%へ到達/
// 未達に変化したときにステータスを自動遷移させる（進捗率が列でなくなった後も、移行前の
// GORM実装が持っていた自動完了・アクティブ復帰の挙動を再現するため）。
func recalcRoadmapStatusFromMetrics(ctx context.Context, q *sqlcgen.Queries, roadmapID uint) error {
	roadmapRow, err := q.GetRoadmapByID(ctx, int64(roadmapID))
	if err != nil {
		return err
	}
	metricsRows, err := q.GetRoadmapMetricsByRoadmapIDs(ctx, []int64{int64(roadmapID)})
	if err != nil {
		return err
	}
	var progress int64
	if len(metricsRows) > 0 {
		progress = metricsRows[0].Progress
	}

	status := model.RoadmapStatus(fromStringPtr(roadmapRow.Status))
	completedStatus := string(model.RoadmapStatusCompleted)
	activeStatus := string(model.RoadmapStatusActive)

	switch {
	case progress >= 100 && status == model.RoadmapStatusActive:
		now := time.Now()
		_, err = q.UpdateRoadmapStatus(ctx, sqlcgen.UpdateRoadmapStatusParams{
			ID:          int64(roadmapID),
			Status:      &completedStatus,
			CompletedAt: toTimestamptz(&now),
		})
	case progress < 100 && status == model.RoadmapStatusCompleted:
		_, err = q.UpdateRoadmapStatus(ctx, sqlcgen.UpdateRoadmapStatusParams{
			ID:          int64(roadmapID),
			Status:      &activeStatus,
			CompletedAt: toTimestamptz(nil),
		})
	}
	return err
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

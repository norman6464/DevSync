package persistence

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRoadmapMetrics_StepLifecycle_UpdatesInSameStatementAndAutoTransitionsStatus は、
// ステップの作成/完了切替/削除がroadmap_metrics（step_count/completed_step_count）を
// 同一SQL文で正しく増減させ、進捗率（SELECT側の算出式）が100%を跨いだときだけ
// ロードマップのステータスが自動遷移することを検証する（DEVSYNC-159）。
// 削除では移行前と同じくステータスの自動遷移を行わない非対称性も確認する。
func TestRoadmapMetrics_StepLifecycle_UpdatesInSameStatementAndAutoTransitionsStatus(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()
	repo := NewRoadmapRepository(pool)

	var userID int64
	err := pool.QueryRow(ctx, `
		INSERT INTO users (username, name, email, created_at, updated_at)
		VALUES ('rm_owner', 'rm_owner', 'rm_owner@example.com', now(), now())
		RETURNING id
	`).Scan(&userID)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	})

	roadmap := &model.Roadmap{UserID: uint(userID), Title: "roadmap"}
	require.NoError(t, repo.Create(ctx, roadmap))

	assertMetrics := func(wantStep, wantCompleted int64) {
		t.Helper()
		var step, completed int64
		row := pool.QueryRow(ctx, `
			SELECT COALESCE(step_count, 0), COALESCE(completed_step_count, 0)
			FROM roadmap_metrics WHERE roadmap_id = $1
		`, roadmap.ID)
		err := row.Scan(&step, &completed)
		if err != nil {
			// roadmap_metrics行がまだ無い（遅延生成前）場合は全て0として扱う。
			step, completed = 0, 0
		}
		assert.Equal(t, wantStep, step, "step_count")
		assert.Equal(t, wantCompleted, completed, "completed_step_count")
	}

	// 最初はroadmap_metrics行が存在しない。
	assertMetrics(0, 0)

	step1 := &model.RoadmapStep{RoadmapID: roadmap.ID, Title: "step1"}
	require.NoError(t, repo.CreateStep(ctx, step1))
	assertMetrics(1, 0)

	step2 := &model.RoadmapStep{RoadmapID: roadmap.ID, Title: "step2"}
	require.NoError(t, repo.CreateStep(ctx, step2))
	assertMetrics(2, 0)

	// step1を完了にする: completed_step_countが増えるが、進捗50%ではステータスは変わらない。
	step1.IsCompleted = true
	require.NoError(t, repo.UpdateStep(ctx, step1))
	assertMetrics(2, 1)

	half, err := repo.FindByID(ctx, roadmap.ID)
	require.NoError(t, err)
	assert.Equal(t, 50, half.Progress)
	assert.Equal(t, model.RoadmapStatusActive, half.Status)

	// step2も完了にする: 進捗100%に到達しステータスがcompletedへ自動遷移する。
	step2.IsCompleted = true
	require.NoError(t, repo.UpdateStep(ctx, step2))
	assertMetrics(2, 2)

	full, err := repo.FindByID(ctx, roadmap.ID)
	require.NoError(t, err)
	assert.Equal(t, 100, full.Progress)
	assert.Equal(t, model.RoadmapStatusCompleted, full.Status)
	assert.NotNil(t, full.CompletedAt)

	// step1を未完了に戻す: 進捗が100%未満に戻りステータスがactiveへ復帰する。
	step1.IsCompleted = false
	require.NoError(t, repo.UpdateStep(ctx, step1))
	assertMetrics(2, 1)

	reactivated, err := repo.FindByID(ctx, roadmap.ID)
	require.NoError(t, err)
	assert.Equal(t, 50, reactivated.Progress)
	assert.Equal(t, model.RoadmapStatusActive, reactivated.Status)
	assert.Nil(t, reactivated.CompletedAt)

	// step1（未完了）を削除する: 残るstep2は完了済みのため進捗は100%になるが、
	// 削除ではステータスの自動遷移を行わない（移行前の挙動を保持）。
	require.NoError(t, repo.DeleteStep(ctx, step1.ID))
	assertMetrics(1, 1)

	afterDelete, err := repo.FindByID(ctx, roadmap.ID)
	require.NoError(t, err)
	assert.Equal(t, 100, afterDelete.Progress)
	assert.Equal(t, model.RoadmapStatusActive, afterDelete.Status, "delete should not auto-transition status even at 100% progress")
}

// TestReconcileAllRoadmapMetrics_CorrectsDrift は、roadmap_stepsの直接操作等で
// roadmap_metricsが実件数からずれたケースを意図的に作り、reconcileジョブが正しく
// 補正することを検証する（DEVSYNC-159のチケット本文で明示されているテスト要件）。
func TestReconcileAllRoadmapMetrics_CorrectsDrift(t *testing.T) {
	pool := cascadeTestDB(t)
	ctx := context.Background()
	q := sqlcgen.New(pool)
	repo := NewRoadmapRepository(pool)
	metricsRepo := NewRoadmapMetricsRepository(q)

	var userID int64
	err := pool.QueryRow(ctx, `
		INSERT INTO users (username, name, email, created_at, updated_at)
		VALUES ('rm_reconcile_owner', 'rm_reconcile_owner', 'rm_reconcile_owner@example.com', now(), now())
		RETURNING id
	`).Scan(&userID)
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	})

	roadmap := &model.Roadmap{UserID: uint(userID), Title: "roadmap"}
	require.NoError(t, repo.Create(ctx, roadmap))

	step1 := &model.RoadmapStep{RoadmapID: roadmap.ID, Title: "step1"}
	require.NoError(t, repo.CreateStep(ctx, step1))
	step2 := &model.RoadmapStep{RoadmapID: roadmap.ID, Title: "step2"}
	require.NoError(t, repo.CreateStep(ctx, step2))
	step1.IsCompleted = true
	require.NoError(t, repo.UpdateStep(ctx, step1))

	var stepCount, completedStepCount int64
	require.NoError(t, pool.QueryRow(ctx, `
		SELECT step_count, completed_step_count FROM roadmap_metrics WHERE roadmap_id = $1
	`, roadmap.ID).Scan(&stepCount, &completedStepCount))
	require.EqualValues(t, 2, stepCount)
	require.EqualValues(t, 1, completedStepCount)

	// 意図的にroadmap_metricsをずらす（roadmap_stepsの直接削除等で、CTEベースの
	// 加減算を経由しない変化が起きた状況を模している）。
	_, err = pool.Exec(ctx, `DELETE FROM roadmap_steps WHERE id = $1`, step2.ID)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `
		UPDATE roadmap_metrics SET step_count = 999, completed_step_count = 999 WHERE roadmap_id = $1
	`, roadmap.ID)
	require.NoError(t, err)

	require.NoError(t, pool.QueryRow(ctx, `
		SELECT step_count, completed_step_count FROM roadmap_metrics WHERE roadmap_id = $1
	`, roadmap.ID).Scan(&stepCount, &completedStepCount))
	require.EqualValues(t, 999, stepCount, "drift is set up correctly before reconcile")
	require.EqualValues(t, 999, completedStepCount, "drift is set up correctly before reconcile")

	require.NoError(t, metricsRepo.Reconcile(ctx))

	require.NoError(t, pool.QueryRow(ctx, `
		SELECT step_count, completed_step_count FROM roadmap_metrics WHERE roadmap_id = $1
	`, roadmap.ID).Scan(&stepCount, &completedStepCount))
	assert.EqualValues(t, 1, stepCount, "reconcile should correct step_count back to the real roadmap_steps row count")
	assert.EqualValues(t, 1, completedStepCount, "reconcile should correct completed_step_count back to the real completed row count")
}

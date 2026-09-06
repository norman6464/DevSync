package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// streakFreezeRepository は [repository.StreakFreezeRepository] の sqlc(pgx) 実装。
// CreateWithinLimits は行ロック・判定・作成を1トランザクションで行うため、
// Queries だけでなくトランザクションを開始できる *pgxpool.Pool を直接保持する。
type streakFreezeRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewStreakFreezeRepository は StreakFreezeRepository の sqlc(pgx) 実装を返す。
func NewStreakFreezeRepository(pool *pgxpool.Pool) repository.StreakFreezeRepository {
	return &streakFreezeRepository{pool: pool, q: sqlcgen.New(pool)}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.StreakFreezeRepository = (*streakFreezeRepository)(nil)

// toModelStreakFreeze は sqlc の生成行を model.StreakFreeze へ変換する。
func toModelStreakFreeze(row sqlcgen.StreakFreeze) model.StreakFreeze {
	freeze := model.StreakFreeze{
		ID:       uint(row.ID),
		UserID:   uint(row.UserID),
		UsedDate: fromDateToDateString(row.UsedOn),
		Month:    int(row.Month),
		Year:     int(row.Year),
	}
	if t := fromTimestamptz(row.CreatedAt); t != nil {
		freeze.CreatedAt = *t
	}
	return freeze
}

// CreateWithinLimits は当日重複・月次上限の判定とフリーズ作成を 1 トランザクションで行う。
// ユーザー行の行ロックで同一ユーザーの同時実行を直列化し、判定と作成の間に
// 他のリクエストが差し込めないようにする（月次上限は件数ベースで一意制約では守れないため）。
func (r *streakFreezeRepository) CreateWithinLimits(ctx context.Context, freeze *model.StreakFreeze, maxPerMonth int) (repository.FreezeUseOutcome, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return repository.FreezeUseCreated, err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	if err := q.LockUserForStreakFreeze(ctx, int64(freeze.UserID)); err != nil {
		return repository.FreezeUseCreated, err
	}

	dayCount, err := q.CountStreakFreezesOnDate(ctx, sqlcgen.CountStreakFreezesOnDateParams{
		UserID: int64(freeze.UserID),
		UsedOn: toDateFromDateString(freeze.UsedDate),
	})
	if err != nil {
		return repository.FreezeUseCreated, err
	}
	if dayCount > 0 {
		return repository.FreezeUseDuplicateDay, tx.Commit(ctx)
	}

	monthCount, err := q.CountStreakFreezesInMonth(ctx, sqlcgen.CountStreakFreezesInMonthParams{
		UserID: int64(freeze.UserID),
		Year:   int64(freeze.Year),
		Month:  int64(freeze.Month),
	})
	if err != nil {
		return repository.FreezeUseCreated, err
	}
	if monthCount >= int64(maxPerMonth) {
		return repository.FreezeUseMonthlyLimitReached, tx.Commit(ctx)
	}

	row, err := q.CreateStreakFreeze(ctx, sqlcgen.CreateStreakFreezeParams{
		UserID: int64(freeze.UserID),
		UsedOn: toDateFromDateString(freeze.UsedDate),
		Month:  int64(freeze.Month),
		Year:   int64(freeze.Year),
	})
	if err != nil {
		return repository.FreezeUseCreated, err
	}
	*freeze = toModelStreakFreeze(row)
	return repository.FreezeUseCreated, tx.Commit(ctx)
}

// GetByUserIDAndMonth は指定ユーザーの指定月のフリーズ一覧を取得する。
func (r *streakFreezeRepository) GetByUserIDAndMonth(ctx context.Context, userID uint, year, month int) ([]model.StreakFreeze, error) {
	rows, err := r.q.ListStreakFreezesByUserAndMonth(ctx, sqlcgen.ListStreakFreezesByUserAndMonthParams{
		UserID: int64(userID),
		Year:   int64(year),
		Month:  int64(month),
	})
	if err != nil {
		return nil, err
	}
	freezes := make([]model.StreakFreeze, len(rows))
	for i, row := range rows {
		freezes[i] = toModelStreakFreeze(row)
	}
	return freezes, nil
}

// HasFreezeOnDate は指定日にフリーズが使用されているかを返す。
func (r *streakFreezeRepository) HasFreezeOnDate(ctx context.Context, userID uint, date string) (bool, error) {
	return r.q.HasStreakFreezeOnDate(ctx, sqlcgen.HasStreakFreezeOnDateParams{
		UserID: int64(userID),
		UsedOn: toDateFromDateString(date),
	})
}

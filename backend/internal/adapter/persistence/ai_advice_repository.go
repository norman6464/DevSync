package persistence

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// errAIAdviceNotFound は既読にする対象のアドバイスが無いときに返すエラー。
var errAIAdviceNotFound = errors.New("アドバイスが見つかりません")

// aiAdviceRepository は [repository.AIAdviceRepository] の sqlc(pgx) 実装。
// CreateBatch は複数件の作成を1トランザクションで行うため、
// Queries だけでなくトランザクションを開始できる *pgxpool.Pool を直接保持する。
type aiAdviceRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewAIAdviceRepository は AIAdviceRepository の sqlc(pgx) 実装を返す。
func NewAIAdviceRepository(pool *pgxpool.Pool) repository.AIAdviceRepository {
	return &aiAdviceRepository{pool: pool, q: sqlcgen.New(pool)}
}

var _ repository.AIAdviceRepository = (*aiAdviceRepository)(nil)

// toModelAIAdvice は sqlc の生成行を model.AIAdvice へ変換する。
func toModelAIAdvice(row sqlcgen.AiAdvice) model.AIAdvice {
	return model.AIAdvice{
		ID:         uint(row.ID),
		UserID:     uint(row.UserID),
		Type:       model.AdviceType(row.Type),
		Priority:   model.AdvicePriority(row.Priority),
		TitleKey:   row.TitleKey,
		MessageKey: row.MessageKey,
		Params:     fromStringPtr(row.Params),
		ActionURL:  fromStringPtr(row.ActionUrl),
		IsRead:     row.IsRead,
		ExpiresAt:  fromTimestamptz(row.ExpiresAt),
		CreatedAt:  timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:  timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// CreateBatch は複数のアドバイスを一括作成する。
// 移行前のGORMバッチ作成（単一INSERT文）と同じ原子性を保つため、1トランザクションで行う。
func (r *aiAdviceRepository) CreateBatch(ctx context.Context, advices []*model.AIAdvice) error {
	if len(advices) == 0 {
		return nil
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	for _, advice := range advices {
		row, err := q.CreateAIAdvice(ctx, sqlcgen.CreateAIAdviceParams{
			UserID:     int64(advice.UserID),
			Type:       string(advice.Type),
			Priority:   int64(advice.Priority),
			TitleKey:   advice.TitleKey,
			MessageKey: advice.MessageKey,
			Params:     &advice.Params,
			ActionUrl:  &advice.ActionURL,
			ExpiresAt:  toTimestamptz(advice.ExpiresAt),
		})
		if err != nil {
			return err
		}
		*advice = toModelAIAdvice(row)
	}
	return tx.Commit(ctx)
}

// FindByUserID は優先度の高い順・作成の新しい順にアドバイスを返す。
func (r *aiAdviceRepository) FindByUserID(ctx context.Context, userID uint, limit int) ([]model.AIAdvice, error) {
	rows, err := r.q.ListAIAdvicesByUser(ctx, sqlcgen.ListAIAdvicesByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
	})
	if err != nil {
		return nil, err
	}
	advices := make([]model.AIAdvice, len(rows))
	for i, row := range rows {
		advices[i] = toModelAIAdvice(row)
	}
	return advices, nil
}

// FindUnreadByUserID は未読のアドバイスを優先度順で返す。
func (r *aiAdviceRepository) FindUnreadByUserID(ctx context.Context, userID uint) ([]model.AIAdvice, error) {
	rows, err := r.q.ListUnreadAIAdvicesByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	advices := make([]model.AIAdvice, len(rows))
	for i, row := range rows {
		advices[i] = toModelAIAdvice(row)
	}
	return advices, nil
}

// MarkAsRead は本人のアドバイス 1 件を既読にする。対象が無ければエラーを返す。
func (r *aiAdviceRepository) MarkAsRead(ctx context.Context, id, userID uint) error {
	rowsAffected, err := r.q.MarkAIAdviceAsRead(ctx, sqlcgen.MarkAIAdviceAsReadParams{
		ID:     int64(id),
		UserID: int64(userID),
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errAIAdviceNotFound
	}
	return nil
}

// DeleteByUserID は指定ユーザーのアドバイスをすべて削除する。
func (r *aiAdviceRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	return r.q.DeleteAIAdvicesByUser(ctx, int64(userID))
}

package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// aiConversationRepository は [repository.AIConversationRepository] の sqlc(pgx) 実装。
// DeleteConversation は複数テーブルの削除を1トランザクションで行うため、
// Queries だけでなくトランザクションを開始できる *pgxpool.Pool を直接保持する。
type aiConversationRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewAIConversationRepository は AIConversationRepository の sqlc(pgx) 実装を返す。
func NewAIConversationRepository(pool *pgxpool.Pool) repository.AIConversationRepository {
	return &aiConversationRepository{pool: pool, q: sqlcgen.New(pool)}
}

var _ repository.AIConversationRepository = (*aiConversationRepository)(nil)

func toModelAIConversation(row sqlcgen.AiConversation) model.AIConversation {
	return model.AIConversation{
		ID:        uint(row.ID),
		UserID:    uint(row.UserID),
		Title:     row.Title,
		CreatedAt: timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt: timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

func toModelAIMessage(row sqlcgen.AiMessage) model.AIMessage {
	return model.AIMessage{
		ID:             uint(row.ID),
		ConversationID: uint(row.ConversationID),
		Role:           model.AIMessageRole(row.Role),
		Content:        row.Content,
		TokensUsed:     int(fromInt64PtrValue(row.TokensUsed)),
		CreatedAt:      timeValue(fromTimestamptz(row.CreatedAt)),
	}
}

// CreateConversation は会話を作成する。
func (r *aiConversationRepository) CreateConversation(ctx context.Context, conv *model.AIConversation) error {
	row, err := r.q.CreateAIConversation(ctx, sqlcgen.CreateAIConversationParams{
		UserID: int64(conv.UserID),
		Title:  conv.Title,
	})
	if err != nil {
		return err
	}
	*conv = toModelAIConversation(row)
	return nil
}

// FindConversationsByUserID は会話を更新の新しい順にメッセージ付きで返す。
func (r *aiConversationRepository) FindConversationsByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.AIConversation, error) {
	rows, err := r.q.ListAIConversationsByUser(ctx, sqlcgen.ListAIConversationsByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, err
	}

	conversations := make([]model.AIConversation, len(rows))
	convIDs := make([]int64, len(rows))
	for i, row := range rows {
		conversations[i] = toModelAIConversation(row)
		convIDs[i] = row.ID
	}

	if len(convIDs) > 0 {
		msgRows, err := r.q.ListAIMessagesByConversationIDs(ctx, convIDs)
		if err != nil {
			return nil, err
		}
		messagesByConvID := make(map[uint][]model.AIMessage)
		for _, row := range msgRows {
			convID := uint(row.ConversationID)
			messagesByConvID[convID] = append(messagesByConvID[convID], toModelAIMessage(row))
		}
		for i := range conversations {
			conversations[i].Messages = messagesByConvID[conversations[i].ID]
		}
	}

	return conversations, nil
}

// FindConversationByID は会話をメッセージ付き（古い順）で返す。存在しなければ (nil, nil) を返す。
func (r *aiConversationRepository) FindConversationByID(ctx context.Context, id uint) (*model.AIConversation, error) {
	row, err := r.q.GetAIConversationByID(ctx, int64(id))
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	conv := toModelAIConversation(row)

	msgRows, err := r.q.ListAIMessagesByConversationIDOrderedByCreatedAt(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	conv.Messages = make([]model.AIMessage, len(msgRows))
	for i, msgRow := range msgRows {
		conv.Messages[i] = toModelAIMessage(msgRow)
	}

	return &conv, nil
}

// AddMessage はメッセージを追加し、会話の更新日時を進める。
func (r *aiConversationRepository) AddMessage(ctx context.Context, msg *model.AIMessage) error {
	row, err := r.q.CreateAIMessage(ctx, sqlcgen.CreateAIMessageParams{
		ConversationID: int64(msg.ConversationID),
		Role:           string(msg.Role),
		Content:        msg.Content,
		TokensUsed:     toInt64Ptr(msg.TokensUsed),
	})
	if err != nil {
		return err
	}
	*msg = toModelAIMessage(row)
	return r.q.TouchAIConversation(ctx, sqlcgen.TouchAIConversationParams{
		ID:        int64(msg.ConversationID),
		UpdatedAt: toTimestamptzNotNull(time.Now()),
	})
}

// CountTodayMessages は当日のユーザー発言数を返す。
func (r *aiConversationRepository) CountTodayMessages(ctx context.Context, userID uint) (int64, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return r.q.CountTodayAIMessagesByUser(ctx, sqlcgen.CountTodayAIMessagesByUserParams{
		UserID:    int64(userID),
		Role:      string(model.AIMessageRoleUser),
		CreatedAt: toTimestamptzNotNull(startOfDay),
	})
}

// DeleteConversation は本人の会話をメッセージごと削除する。
// 所有権の判定は usecase 側で済んでいる前提で、本人の会話だけを対象にする。
// 既に無ければ何もしない（冪等）。エラーは DB 障害だけを表す。
// メッセージと会話は 1 トランザクションで消し、途中で失敗しても片方だけ消えない。
func (r *aiConversationRepository) DeleteConversation(ctx context.Context, id, userID uint) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	_, err = q.GetAIConversationByIDAndUser(ctx, sqlcgen.GetAIConversationByIDAndUserParams{
		ID:     int64(id),
		UserID: int64(userID),
	})
	if err != nil {
		if isNoRows(err) {
			return nil
		}
		return err
	}

	if err := q.DeleteAIMessagesByConversationID(ctx, int64(id)); err != nil {
		return err
	}
	if err := q.DeleteAIConversation(ctx, int64(id)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

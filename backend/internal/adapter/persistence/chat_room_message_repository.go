package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// chatRoomMessageRepository は [repository.ChatRoomMessageRepository] の sqlc(pgx) 実装。
type chatRoomMessageRepository struct {
	q *sqlcgen.Queries
}

// NewChatRoomMessageRepository は ChatRoomMessageRepository の sqlc(pgx) 実装を返す。
func NewChatRoomMessageRepository(q *sqlcgen.Queries) repository.ChatRoomMessageRepository {
	return &chatRoomMessageRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ChatRoomMessageRepository = (*chatRoomMessageRepository)(nil)

// toModelGroupMessage は sqlc の生成行を model.GroupMessage へ変換する（Sender は含まない）。
func toModelGroupMessage(row sqlcgen.GroupMessage) model.GroupMessage {
	return model.GroupMessage{
		ID:         uint(row.ID),
		ChatRoomID: uint(row.ChatRoomID),
		SenderID:   uint(row.SenderID),
		Content:    row.Content,
		CreatedAt:  timeValue(fromTimestamptz(row.CreatedAt)),
	}
}

// Create はグループメッセージを保存する。
func (r *chatRoomMessageRepository) Create(ctx context.Context, msg *model.GroupMessage) error {
	row, err := r.q.CreateGroupMessage(ctx, sqlcgen.CreateGroupMessageParams{
		ChatRoomID: int64(msg.ChatRoomID),
		SenderID:   int64(msg.SenderID),
		Content:    msg.Content,
	})
	if err != nil {
		return err
	}
	*msg = toModelGroupMessage(row)
	return nil
}

// FindByRoomID は指定ルームのメッセージを送信者情報付きで、作成日時の昇順に取得する。
func (r *chatRoomMessageRepository) FindByRoomID(ctx context.Context, roomID uint, page, limit int) ([]model.GroupMessage, error) {
	offset := (page - 1) * limit
	rows, err := r.q.ListGroupMessagesByRoom(ctx, sqlcgen.ListGroupMessagesByRoomParams{
		ChatRoomID: int64(roomID),
		Limit:      int32Param(limit),
		Offset:     int32Param(offset),
	})
	if err != nil {
		return nil, err
	}
	messages := make([]model.GroupMessage, len(rows))
	for i, row := range rows {
		messages[i] = toModelGroupMessage(row.GroupMessage)
		sender := toModelUser(row.User)
		messages[i].Sender = &sender
	}
	return messages, nil
}

// FindSender は送信者のユーザー情報を取得する。不在の場合は (nil, nil) を返す。
func (r *chatRoomMessageRepository) FindSender(ctx context.Context, senderID uint) (*model.User, error) {
	row, err := r.q.GetUserByIDForChatSender(ctx, int64(senderID))
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	user := toModelUser(row)
	return &user, nil
}

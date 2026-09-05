package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// chatRoomRepository は [repository.ChatRoomRepository] の sqlc(pgx) 実装。
// Delete はメッセージ・メンバーごとの削除を1トランザクションで行うため、
// Queries だけでなくトランザクションを開始できる *pgxpool.Pool を直接保持する。
type chatRoomRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewChatRoomRepository は ChatRoomRepository の sqlc(pgx) 実装を返す。
func NewChatRoomRepository(pool *pgxpool.Pool) repository.ChatRoomRepository {
	return &chatRoomRepository{pool: pool, q: sqlcgen.New(pool)}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ChatRoomRepository = (*chatRoomRepository)(nil)

func toModelChatRoom(row sqlcgen.ChatRoom) model.ChatRoom {
	return model.ChatRoom{
		ID:          uint(row.ID),
		Name:        row.Name,
		Description: fromStringPtr(row.Description),
		OwnerID:     uint(row.OwnerID),
		CreatedAt:   timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:   timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// Create はチャットルームを作成する。
func (r *chatRoomRepository) Create(ctx context.Context, room *model.ChatRoom) error {
	row, err := r.q.CreateChatRoom(ctx, sqlcgen.CreateChatRoomParams{
		Name:        room.Name,
		Description: &room.Description,
		OwnerID:     int64(room.OwnerID),
	})
	if err != nil {
		return err
	}
	*room = toModelChatRoom(row)
	return nil
}

// FindByID は指定 ID のチャットルームをオーナー情報付きで取得する。不在の場合は (nil, nil) を返す。
func (r *chatRoomRepository) FindByID(ctx context.Context, id uint) (*model.ChatRoom, error) {
	row, err := r.q.GetChatRoomWithOwnerByID(ctx, int64(id))
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	room := toModelChatRoom(row.ChatRoom)
	owner := toModelUser(row.User)
	room.Owner = &owner
	return &room, nil
}

// FindByUserID は指定ユーザーが参加しているチャットルームを更新日時の降順で取得する。
func (r *chatRoomRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.ChatRoom, int64, error) {
	// 総件数の取得エラーは移行前と同じく無視し、一覧取得のエラーだけを返す。
	total, _ := r.q.CountChatRoomsByUser(ctx, int64(userID))

	rows, err := r.q.ListChatRoomsByUser(ctx, sqlcgen.ListChatRoomsByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	rooms := make([]model.ChatRoom, len(rows))
	for i, row := range rows {
		rooms[i] = toModelChatRoom(row.ChatRoom)
		owner := toModelUser(row.User)
		rooms[i].Owner = &owner
	}
	return rooms, total, nil
}

// Update はチャットルーム情報を更新する（GORMのSave＝全カラム上書きに相当）。
func (r *chatRoomRepository) Update(ctx context.Context, room *model.ChatRoom) error {
	row, err := r.q.UpdateChatRoom(ctx, sqlcgen.UpdateChatRoomParams{
		ID:          int64(room.ID),
		Name:        room.Name,
		Description: &room.Description,
	})
	if err != nil {
		return err
	}
	*room = toModelChatRoom(row)
	return nil
}

// Delete はチャットルームをメッセージ・メンバーごとトランザクション内で削除する。
func (r *chatRoomRepository) Delete(ctx context.Context, roomID uint) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	if err := q.DeleteGroupMessagesByRoom(ctx, int64(roomID)); err != nil {
		return err
	}
	if err := q.DeleteChatRoomMembersByRoom(ctx, int64(roomID)); err != nil {
		return err
	}
	if err := q.DeleteChatRoom(ctx, int64(roomID)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// AddMember はチャットルームにメンバーを追加する。
func (r *chatRoomRepository) AddMember(ctx context.Context, roomID, userID uint) error {
	return r.q.CreateChatRoomMember(ctx, sqlcgen.CreateChatRoomMemberParams{
		ChatRoomID: int64(roomID),
		UserID:     int64(userID),
		JoinedAt:   toTimestamptzNotNull(time.Now()),
	})
}

// RemoveMember はチャットルームからメンバーを除外する。
func (r *chatRoomRepository) RemoveMember(ctx context.Context, roomID, userID uint) error {
	return r.q.DeleteChatRoomMember(ctx, sqlcgen.DeleteChatRoomMemberParams{
		ChatRoomID: int64(roomID),
		UserID:     int64(userID),
	})
}

// GetMembers はチャットルームの全メンバーをユーザー情報付きで取得する。
func (r *chatRoomRepository) GetMembers(ctx context.Context, roomID uint) ([]model.ChatRoomMember, error) {
	rows, err := r.q.ListChatRoomMembersWithUser(ctx, int64(roomID))
	if err != nil {
		return nil, err
	}
	members := make([]model.ChatRoomMember, len(rows))
	for i, row := range rows {
		user := toModelUser(row.User)
		members[i] = model.ChatRoomMember{
			ID:         uint(row.ChatRoomMember.ID),
			ChatRoomID: uint(row.ChatRoomMember.ChatRoomID),
			UserID:     uint(row.ChatRoomMember.UserID),
			User:       &user,
			JoinedAt:   timeValue(fromTimestamptz(row.ChatRoomMember.JoinedAt)),
		}
	}
	return members, nil
}

// IsMember は指定ユーザーがチャットルームのメンバーかを判定する。
func (r *chatRoomRepository) IsMember(ctx context.Context, roomID, userID uint) (bool, error) {
	count, err := r.q.CountChatRoomMembership(ctx, sqlcgen.CountChatRoomMembershipParams{
		ChatRoomID: int64(roomID),
		UserID:     int64(userID),
	})
	return count > 0, err
}

// CountByUserID は指定ユーザーが参加しているチャットルーム総数を返す。
func (r *chatRoomRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountChatRoomsByUser(ctx, int64(userID))
}

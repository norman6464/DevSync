package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
)

// RoomMemberLookup は WebSocket の配信先を決めるためのルームメンバー取得。
// 契約は利用側（ハブ）が宣言するため、ここでは実装だけを提供する。
type RoomMemberLookup struct {
	q *sqlcgen.Queries
}

// NewRoomMemberLookup は RoomMemberLookup を生成する。
func NewRoomMemberLookup(q *sqlcgen.Queries) *RoomMemberLookup {
	return &RoomMemberLookup{q: q}
}

// MemberUserIDs は指定チャットルームの全メンバーのユーザー ID を返す。
func (r *RoomMemberLookup) MemberUserIDs(ctx context.Context, roomID uint) ([]uint, error) {
	rows, err := r.q.ListChatRoomMemberUserIDs(ctx, int64(roomID))
	if err != nil {
		return nil, err
	}
	userIDs := make([]uint, len(rows))
	for i, id := range rows {
		userIDs[i] = uint(id)
	}
	return userIDs, nil
}

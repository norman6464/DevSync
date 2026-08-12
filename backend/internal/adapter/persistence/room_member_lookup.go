package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// RoomMemberLookup は WebSocket の配信先を決めるためのルームメンバー取得。
// 契約は利用側（ハブ）が宣言するため、ここでは実装だけを提供する。
type RoomMemberLookup struct {
	db *gorm.DB
}

// NewRoomMemberLookup は RoomMemberLookup を生成する。
func NewRoomMemberLookup(db *gorm.DB) *RoomMemberLookup {
	return &RoomMemberLookup{db: db}
}

// MemberUserIDs は指定チャットルームの全メンバーのユーザー ID を返す。
func (r *RoomMemberLookup) MemberUserIDs(ctx context.Context, roomID uint) ([]uint, error) {
	var userIDs []uint
	err := r.db.WithContext(ctx).Model(&model.ChatRoomMember{}).
		Where("chat_room_id = ?", roomID).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}

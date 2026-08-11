package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// GroupMessageRepository はWebSocket配信に必要なグループチャットのメンバー情報を提供する。
// メッセージ本体の永続化はクリーンアーキテクチャ（DIP）へ移行済みで adapter/persistence が担う。
type GroupMessageRepository struct {
	db *gorm.DB
}

// NewGroupMessageRepository は新しいGroupMessageRepositoryインスタンスを生成する。
func NewGroupMessageRepository(db *gorm.DB) *GroupMessageRepository {
	return &GroupMessageRepository{db: db}
}

// GetMemberUserIDs は指定チャットルームの全メンバーのユーザーIDリストを取得する。
// WebSocket配信先の特定に使用される。
func (r *GroupMessageRepository) GetMemberUserIDs(roomID uint) []uint {
	var userIDs []uint
	r.db.Model(&model.ChatRoomMember{}).
		Where("chat_room_id = ?", roomID).
		Pluck("user_id", &userIDs)
	return userIDs
}

package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// GroupMessageRepository はグループチャットメッセージデータへのアクセスを提供するリポジトリ実装。
type GroupMessageRepository struct {
	db *gorm.DB
}

// NewGroupMessageRepository は新しいGroupMessageRepositoryインスタンスを生成する。
func NewGroupMessageRepository(db *gorm.DB) *GroupMessageRepository {
	return &GroupMessageRepository{db: db}
}

// Create は新しいグループメッセージをデータベースに作成する。
func (r *GroupMessageRepository) Create(msg *model.GroupMessage) error {
	return r.db.Create(msg).Error
}

// FindByRoomID は指定チャットルームのメッセージをページネーション付きで取得する。
// 送信者情報をPreloadし、作成日時の昇順（古い順）でソートされる。
func (r *GroupMessageRepository) FindByRoomID(roomID uint, page, limit int) ([]model.GroupMessage, error) {
	var messages []model.GroupMessage
	offset := (page - 1) * limit
	err := r.db.Preload("Sender").
		Where("chat_room_id = ?", roomID).
		Order("created_at ASC").
		Offset(offset).Limit(limit).
		Find(&messages).Error
	return messages, err
}

// FindSenderByID はメッセージの送信者情報をデータベースから取得してセットする。
func (r *GroupMessageRepository) FindSenderByID(msg *model.GroupMessage) {
	r.db.Model(&model.User{}).First(msg.Sender, msg.SenderID)
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

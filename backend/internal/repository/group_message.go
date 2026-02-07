package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

type GroupMessageRepository struct {
	db *gorm.DB
}

func NewGroupMessageRepository(db *gorm.DB) *GroupMessageRepository {
	return &GroupMessageRepository{db: db}
}

func (r *GroupMessageRepository) Create(msg *model.GroupMessage) error {
	return r.db.Create(msg).Error
}

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

func (r *GroupMessageRepository) GetMemberUserIDs(roomID uint) []uint {
	var userIDs []uint
	r.db.Model(&model.ChatRoomMember{}).
		Where("chat_room_id = ?", roomID).
		Pluck("user_id", &userIDs)
	return userIDs
}

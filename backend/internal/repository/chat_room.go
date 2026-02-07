package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

type ChatRoomRepository struct {
	db *gorm.DB
}

func NewChatRoomRepository(db *gorm.DB) *ChatRoomRepository {
	return &ChatRoomRepository{db: db}
}

func (r *ChatRoomRepository) Create(room *model.ChatRoom) error {
	return r.db.Create(room).Error
}

func (r *ChatRoomRepository) FindByID(id uint) (*model.ChatRoom, error) {
	var room model.ChatRoom
	err := r.db.Preload("Owner").First(&room, id).Error
	return &room, err
}

func (r *ChatRoomRepository) FindByUserID(userID uint) ([]model.ChatRoom, error) {
	var rooms []model.ChatRoom
	err := r.db.Joins("JOIN chat_room_members ON chat_room_members.chat_room_id = chat_rooms.id").
		Where("chat_room_members.user_id = ?", userID).
		Preload("Owner").
		Order("chat_rooms.updated_at DESC").
		Find(&rooms).Error
	return rooms, err
}

func (r *ChatRoomRepository) Update(room *model.ChatRoom) error {
	return r.db.Save(room).Error
}

func (r *ChatRoomRepository) Delete(roomID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("chat_room_id = ?", roomID).Delete(&model.GroupMessage{}).Error; err != nil {
			return err
		}
		if err := tx.Where("chat_room_id = ?", roomID).Delete(&model.ChatRoomMember{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.ChatRoom{}, roomID).Error
	})
}

func (r *ChatRoomRepository) AddMember(roomID, userID uint) error {
	member := model.ChatRoomMember{
		ChatRoomID: roomID,
		UserID:     userID,
		JoinedAt:   time.Now(),
	}
	return r.db.Create(&member).Error
}

func (r *ChatRoomRepository) RemoveMember(roomID, userID uint) error {
	return r.db.Where("chat_room_id = ? AND user_id = ?", roomID, userID).
		Delete(&model.ChatRoomMember{}).Error
}

func (r *ChatRoomRepository) GetMembers(roomID uint) ([]model.ChatRoomMember, error) {
	var members []model.ChatRoomMember
	err := r.db.Preload("User").Where("chat_room_id = ?", roomID).Find(&members).Error
	return members, err
}

func (r *ChatRoomRepository) IsMember(roomID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.ChatRoomMember{}).
		Where("chat_room_id = ? AND user_id = ?", roomID, userID).
		Count(&count).Error
	return count > 0, err
}

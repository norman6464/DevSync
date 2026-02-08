package model

import "time"

type ChatRoom struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"size:500"`
	OwnerID     uint      `json:"owner_id" gorm:"not null;index"`
	Owner       *User     `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ChatRoomMember struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	ChatRoomID uint      `json:"chat_room_id" gorm:"not null;index;uniqueIndex:idx_room_user"`
	ChatRoom   *ChatRoom `json:"-" gorm:"foreignKey:ChatRoomID"`
	UserID     uint      `json:"user_id" gorm:"not null;index;uniqueIndex:idx_room_user"`
	User       *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	JoinedAt   time.Time `json:"joined_at"`
}

type GroupMessage struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	ChatRoomID uint      `json:"chat_room_id" gorm:"not null;index"`
	ChatRoom   *ChatRoom `json:"-" gorm:"foreignKey:ChatRoomID"`
	SenderID   uint      `json:"sender_id" gorm:"not null;index"`
	Sender     *User     `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	CreatedAt  time.Time `json:"created_at"`
}

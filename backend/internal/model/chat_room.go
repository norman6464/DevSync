package model

import "time"

// ChatRoom はグループチャットのルームを表す。
// OwnerID はルーム作成者を示し、ルームの管理権限を持つ。
type ChatRoom struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`   // ルーム名
	Description string    `json:"description" gorm:"size:500"`     // ルームの説明
	OwnerID     uint      `json:"owner_id" gorm:"not null;index"`  // ルーム作成者のユーザーID
	Owner       *User     `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ChatRoomMember はチャットルームのメンバーシップを表す。
// ChatRoomID と UserID の組み合わせでユニークインデックスを持ち、
// 同一ユーザーが同じルームに重複参加することを防ぐ。
type ChatRoomMember struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	ChatRoomID uint      `json:"chat_room_id" gorm:"not null;index;uniqueIndex:idx_room_user"`
	ChatRoom   *ChatRoom `json:"-" gorm:"foreignKey:ChatRoomID"`
	UserID     uint      `json:"user_id" gorm:"not null;index;uniqueIndex:idx_room_user"`
	User       *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	JoinedAt   time.Time `json:"joined_at"` // メンバー参加日時
}

// GroupMessage はグループチャットルーム内の個別メッセージを表す。
// SenderID で送信者を識別し、ChatRoomID で所属するルームを示す。
type GroupMessage struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	ChatRoomID uint      `json:"chat_room_id" gorm:"not null;index"` // 送信先チャットルームID
	ChatRoom   *ChatRoom `json:"-" gorm:"foreignKey:ChatRoomID"`
	SenderID   uint      `json:"sender_id" gorm:"not null;index"` // 送信者のユーザーID
	Sender     *User     `json:"sender,omitempty" gorm:"foreignKey:SenderID"`
	Content    string    `json:"content" gorm:"type:text;not null"` // メッセージ本文
	CreatedAt  time.Time `json:"created_at"`
}

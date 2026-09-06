package model

import "time"

// ChatRoom はグループチャットのルームを表す。
// OwnerID はルーム作成者を示し、ルームの管理権限を持つ。
type ChatRoom struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`        // ルーム名
	Description string    `json:"description"` // ルームの説明
	OwnerID     uint      `json:"owner_id"`    // ルーム作成者のユーザーID
	Owner       *User     `json:"owner,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ChatRoomMember はチャットルームのメンバーシップを表す。
// ChatRoomID と UserID の組み合わせでユニークインデックスを持ち、
// 同一ユーザーが同じルームに重複参加することを防ぐ。
type ChatRoomMember struct {
	ID         uint      `json:"id"`
	ChatRoomID uint      `json:"chat_room_id"`
	ChatRoom   *ChatRoom `json:"-"`
	UserID     uint      `json:"user_id"`
	User       *User     `json:"user,omitempty"`
	JoinedAt   time.Time `json:"joined_at"` // メンバー参加日時
}

// GroupMessage はグループチャットルーム内の個別メッセージを表す。
// SenderID で送信者を識別し、ChatRoomID で所属するルームを示す。
type GroupMessage struct {
	ID         uint      `json:"id"`
	ChatRoomID uint      `json:"chat_room_id"` // 送信先チャットルームID
	ChatRoom   *ChatRoom `json:"-"`
	SenderID   uint      `json:"sender_id"` // 送信者のユーザーID
	Sender     *User     `json:"sender,omitempty"`
	Content    string    `json:"content"` // メッセージ本文
	CreatedAt  time.Time `json:"created_at"`
}

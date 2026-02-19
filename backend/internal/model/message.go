package model

import "time"

// ConversationSummary は会話一覧表示用のサマリー情報を表す。
type ConversationSummary struct {
	UserID      uint   `json:"user_id"`      // 会話相手のユーザーID
	Name        string `json:"name"`         // 会話相手の名前
	AvatarURL   string `json:"avatar_url"`   // 会話相手のアバターURL
	LastMessage string `json:"last_message"` // 最新メッセージの内容
	LastTime    string `json:"last_time"`    // 最新メッセージの日時
	UnreadCount int    `json:"unread_count"` // 未読メッセージ数
}

// Message はユーザー間のダイレクトメッセージを表す。
type Message struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	SenderID   uint      `json:"sender_id" gorm:"not null;index"`   // 送信者のユーザーID
	Sender     User      `json:"sender" gorm:"foreignKey:SenderID"`
	ReceiverID uint      `json:"receiver_id" gorm:"not null;index"` // 受信者のユーザーID
	Receiver   User      `json:"receiver" gorm:"foreignKey:ReceiverID"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	Read       bool      `json:"read" gorm:"default:false"` // 既読フラグ
	CreatedAt  time.Time `json:"created_at"`
}

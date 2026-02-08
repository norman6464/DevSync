package model

import "time"

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

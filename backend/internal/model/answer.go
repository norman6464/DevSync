package model

import (
	"time"

	"gorm.io/gorm"
)

type Answer struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"not null;index"`
	User       User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	QuestionID uint           `json:"question_id" gorm:"not null;index"`
	Body       string         `json:"body" gorm:"type:text;not null"`
	VoteCount  int            `json:"vote_count" gorm:"default:0"`
	IsBest     bool           `json:"is_best" gorm:"default:false"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type AnswerVote struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_answer_vote"`
	AnswerID  uint      `json:"answer_id" gorm:"not null;uniqueIndex:idx_answer_vote"`
	Value     int       `json:"value" gorm:"not null"` // +1 or -1
	CreatedAt time.Time `json:"created_at"`
}

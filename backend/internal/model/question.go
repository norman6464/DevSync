package model

import (
	"time"

	"gorm.io/gorm"
)

type Question struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	User        User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Title       string         `json:"title" gorm:"not null;size:500"`
	Body        string         `json:"body" gorm:"type:text;not null"`
	Tags        string         `json:"tags" gorm:"type:text"` // JSON array of tags
	VoteCount   int            `json:"vote_count" gorm:"default:0"`
	AnswerCount int            `json:"answer_count" gorm:"default:0"`
	IsSolved    bool           `json:"is_solved" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type QuestionVote struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_question_vote"`
	QuestionID uint      `json:"question_id" gorm:"not null;uniqueIndex:idx_question_vote"`
	Value      int       `json:"value" gorm:"not null"` // +1 or -1
	CreatedAt  time.Time `json:"created_at"`
}

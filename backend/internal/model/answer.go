package model

import (
	"time"

	"gorm.io/gorm"
)

// Answer はQ&A機能における質問への回答を表す。
// IsBest がtrueの場合、その回答は質問投稿者によってベストアンサーに選ばれたことを示す。
type Answer struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"not null;index"`
	User       User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	QuestionID uint           `json:"question_id" gorm:"not null;index"`        // 回答先の質問ID
	Body       string         `json:"body" gorm:"type:text;not null"`           // 回答本文
	VoteCount  int            `json:"vote_count" gorm:"default:0"`              // 投票数（賛成−反対の合計）
	IsBest     bool           `json:"is_best" gorm:"default:false"`             // ベストアンサーフラグ
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"` // 論理削除用
}

// AnswerVote は回答への投票（賛成/反対）を記録する。
// Value は +1（賛成）または -1（反対）の値を取る。
// ユーザーと回答の組み合わせでユニークインデックスを持つ。
type AnswerVote struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_answer_vote"`
	AnswerID  uint      `json:"answer_id" gorm:"not null;uniqueIndex:idx_answer_vote"`
	Value     int       `json:"value" gorm:"not null"` // +1（賛成）または -1（反対）
	CreatedAt time.Time `json:"created_at"`
}

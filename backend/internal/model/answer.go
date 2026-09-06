package model

import (
	"time"
)

// Answer はQ&A機能における質問への回答を表す。
// IsBest がtrueの場合、その回答は質問投稿者によってベストアンサーに選ばれたことを示す。
type Answer struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	User       User      `json:"user,omitempty"`
	QuestionID uint      `json:"question_id"` // 回答先の質問ID
	Body       string    `json:"body"`        // 回答本文
	VoteCount  int       `json:"vote_count"`  // 投票数（賛成−反対の合計）
	IsBest     bool      `json:"is_best"`     // ベストアンサーフラグ
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// AnswerVote は回答への投票（賛成/反対）を記録する。
// Value は +1（賛成）または -1（反対）の値を取る。
// ユーザーと回答の組み合わせでユニークインデックスを持つ。
type AnswerVote struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	AnswerID  uint      `json:"answer_id"`
	Value     int       `json:"value"` // +1（賛成）または -1（反対）
	CreatedAt time.Time `json:"created_at"`
}

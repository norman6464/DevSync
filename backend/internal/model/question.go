package model

import (
	"time"
)

// Question はQ&A機能における質問を表す。
// Tags にはJSON配列形式のタグ情報を格納し、IsSolved で解決済みかどうかを管理する。
type Question struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	User        User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Title       string    `json:"title" gorm:"not null;size:500"`           // 質問タイトル
	Body        string    `json:"body" gorm:"type:text;not null"`           // 質問本文
	Tags        string    `json:"tags" gorm:"type:text"`                    // JSON配列形式のタグ
	VoteCount   int       `json:"vote_count" gorm:"default:0"`              // 投票数（賛成−反対の合計）
	AnswerCount int       `json:"answer_count" gorm:"default:0"`            // 回答数
	IsSolved    bool      `json:"is_solved" gorm:"default:false"`           // ベストアンサーが選ばれたか
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// QuestionBookmark は質問のブックマークを記録する。
// ユーザーと質問の組み合わせでユニークインデックスを持つ。
type QuestionBookmark struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_question_bookmark"`
	QuestionID uint      `json:"question_id" gorm:"not null;uniqueIndex:idx_question_bookmark"`
	Question   Question  `json:"question,omitempty" gorm:"foreignKey:QuestionID"`
	CreatedAt  time.Time `json:"created_at"`
}

// QuestionVote は質問への投票（賛成/反対）を記録する。
// Value は +1（賛成）または -1（反対）の値を取る。
// ユーザーと質問の組み合わせでユニークインデックスを持つ。
type QuestionVote struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_question_vote"`
	QuestionID uint      `json:"question_id" gorm:"not null;uniqueIndex:idx_question_vote"`
	Value      int       `json:"value" gorm:"not null"` // +1（賛成）または -1（反対）
	CreatedAt  time.Time `json:"created_at"`
}

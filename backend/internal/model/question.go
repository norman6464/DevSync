package model

import (
	"time"
)

// Question はQ&A機能における質問を表す。
// Tags にはJSON配列形式のタグ情報を格納し、IsSolved で解決済みかどうかを管理する。
type Question struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	User        User      `json:"user,omitempty"`
	Title       string    `json:"title"`        // 質問タイトル
	Body        string    `json:"body"`         // 質問本文
	Tags        string    `json:"tags"`         // JSON配列形式のタグ
	VoteCount   int       `json:"vote_count"`   // 投票数（賛成−反対の合計）
	AnswerCount int       `json:"answer_count"` // 回答数
	IsSolved    bool      `json:"is_solved"`    // ベストアンサーが選ばれたか
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// QuestionBookmark は質問のブックマークを記録する。
// ユーザーと質問の組み合わせでユニークインデックスを持つ。
type QuestionBookmark struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	QuestionID uint      `json:"question_id"`
	Question   Question  `json:"question,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// QuestionVote は質問への投票（賛成/反対）を記録する。
// Value は +1（賛成）または -1（反対）の値を取る。
// ユーザーと質問の組み合わせでユニークインデックスを持つ。
type QuestionVote struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	QuestionID uint      `json:"question_id"`
	Value      int       `json:"value"` // +1（賛成）または -1（反対）
	CreatedAt  time.Time `json:"created_at"`
}

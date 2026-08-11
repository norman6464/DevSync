package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// AnswerRepository は Q&A 回答の永続化に対する、usecase 側が要求する契約。
type AnswerRepository interface {
	// Create は回答を作成し、質問の回答数を 1 増やす。
	Create(ctx context.Context, answer *model.Answer) error
	// FindByID は指定 ID の回答を返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindByID(ctx context.Context, id uint) (*model.Answer, error)
	Update(ctx context.Context, answer *model.Answer) error
	// Delete は回答を論理削除し、質問の回答数を 1 減らす。
	Delete(ctx context.Context, answer *model.Answer) error

	FindByQuestionID(ctx context.Context, questionID uint) ([]model.Answer, error)
	FindByVoteRange(ctx context.Context, questionID uint, minVote, maxVote int) ([]model.Answer, error)

	// SetBestAnswer は既存のベストアンサーを解除して指定回答を設定し、質問を解決済みにする。
	SetBestAnswer(ctx context.Context, questionID, answerID uint) error

	Vote(ctx context.Context, userID, answerID uint, value int) error
	RemoveVote(ctx context.Context, userID, answerID uint) error
}

// QuestionReader は回答の作成・ベストアンサー設定で質問を参照するための最小の契約。
// 不在の場合は (nil, nil) を返す。
type QuestionReader interface {
	FindByID(ctx context.Context, id uint) (*model.Question, error)
}

package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// qaStatsRepository は [repository.QAStatsRepository] の sqlc(pgx) 実装。
type qaStatsRepository struct {
	q *sqlcgen.Queries
}

// NewQAStatsRepository は QAStatsRepository の sqlc(pgx) 実装を返す。
func NewQAStatsRepository(q *sqlcgen.Queries) repository.QAStatsRepository {
	return &qaStatsRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.QAStatsRepository = (*qaStatsRepository)(nil)

// GetQAStats は指定ユーザーの Q&A 活動集計統計を返す。
func (r *qaStatsRepository) GetQAStats(ctx context.Context, userID uint) (*model.QAStats, error) {
	questions, err := r.q.CountQuestionsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	answers, err := r.q.CountAnswersByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	bestAnswers, err := r.q.CountBestAnswersByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	questionVotes, err := r.q.SumQuestionVotesByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	answerVotes, err := r.q.SumAnswerVotesByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	return &model.QAStats{
		TotalQuestions:     questions,
		TotalAnswers:       answers,
		BestAnswerCount:    bestAnswers,
		TotalVotesReceived: questionVotes + answerVotes,
	}, nil
}

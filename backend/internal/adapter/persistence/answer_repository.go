package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// answerRepository は [repository.AnswerRepository] の sqlc(pgx) 実装。
// Create/Delete/SetBestAnswer/Vote/RemoveVote は行ロック付きのトランザクションで
// 複数文を実行するため、Queries だけでなくトランザクションを開始できる
// *pgxpool.Pool を直接保持する。
type answerRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewAnswerRepository は AnswerRepository の sqlc(pgx) 実装を返す。
func NewAnswerRepository(pool *pgxpool.Pool) repository.AnswerRepository {
	return &answerRepository{pool: pool, q: sqlcgen.New(pool)}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.AnswerRepository = (*answerRepository)(nil)

func toModelAnswer(row sqlcgen.Answer) model.Answer {
	return model.Answer{
		ID:         uint(row.ID),
		UserID:     uint(row.UserID),
		QuestionID: uint(row.QuestionID),
		Body:       row.Body,
		VoteCount:  int(fromInt64PtrValue(row.VoteCount)),
		IsBest:     row.IsBest,
		CreatedAt:  timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:  timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// Create は回答の作成と質問の回答数の加算を同一トランザクションで行う。
func (r *answerRepository) Create(ctx context.Context, answer *model.Answer) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	row, err := q.CreateAnswer(ctx, sqlcgen.CreateAnswerParams{
		UserID:     int64(answer.UserID),
		QuestionID: int64(answer.QuestionID),
		Body:       answer.Body,
		VoteCount:  toInt64Ptr(answer.VoteCount),
		IsBest:     answer.IsBest,
	})
	if err != nil {
		return err
	}
	if err := q.IncrementQuestionAnswerCount(ctx, int64(answer.QuestionID)); err != nil {
		return err
	}
	*answer = toModelAnswer(row)
	return tx.Commit(ctx)
}

// FindByID は指定 ID の回答をユーザー情報付きで取得する。不在の場合は (nil, nil) を返す。
func (r *answerRepository) FindByID(ctx context.Context, id uint) (*model.Answer, error) {
	row, err := r.q.GetAnswerWithUserByID(ctx, int64(id))
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	answer := toModelAnswer(row.Answer)
	answer.User = toModelUser(row.User)
	return &answer, nil
}

// Update は既存の回答を更新する（GORMのSave＝全カラム上書きに相当）。
func (r *answerRepository) Update(ctx context.Context, answer *model.Answer) error {
	row, err := r.q.UpdateAnswer(ctx, sqlcgen.UpdateAnswerParams{
		ID:        int64(answer.ID),
		Body:      answer.Body,
		VoteCount: toInt64Ptr(answer.VoteCount),
		IsBest:    answer.IsBest,
	})
	if err != nil {
		return err
	}
	*answer = toModelAnswer(row)
	return nil
}

// Delete は回答の論理削除と質問の回答数の減算（0 未満にはしない）を
// 同一トランザクションで行う。SetBestAnswer とロック順序（質問 → 回答）を
// 揃えるため、先に質問行を FOR UPDATE でロックしてデッドロックを防ぐ。
func (r *answerRepository) Delete(ctx context.Context, answer *model.Answer) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	if _, err := q.LockQuestionForAnswerChange(ctx, int64(answer.QuestionID)); err != nil {
		return err
	}
	if err := q.SoftDeleteAnswer(ctx, int64(answer.ID)); err != nil {
		return err
	}
	if err := q.DecrementQuestionAnswerCountFloored(ctx, int64(answer.QuestionID)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// FindByQuestionID は指定質問の回答を取得する（ベストアンサー優先 → 投票数降順 → 作成日昇順）。
func (r *answerRepository) FindByQuestionID(ctx context.Context, questionID uint) ([]model.Answer, error) {
	rows, err := r.q.ListAnswersByQuestionID(ctx, int64(questionID))
	if err != nil {
		return nil, err
	}
	answers := make([]model.Answer, len(rows))
	for i, row := range rows {
		answers[i] = toModelAnswer(row.Answer)
		answers[i].User = toModelUser(row.User)
	}
	return answers, nil
}

// FindByVoteRange は指定質問の回答を投票数の範囲で絞り込んで取得する（投票数降順 → 作成日昇順）。
func (r *answerRepository) FindByVoteRange(ctx context.Context, questionID uint, minVote, maxVote int) ([]model.Answer, error) {
	rows, err := r.q.ListAnswersByVoteRange(ctx, sqlcgen.ListAnswersByVoteRangeParams{
		QuestionID:  int64(questionID),
		VoteCount:   toInt64Ptr(minVote),
		VoteCount_2: toInt64Ptr(maxVote),
	})
	if err != nil {
		return nil, err
	}
	answers := make([]model.Answer, len(rows))
	for i, row := range rows {
		answers[i] = toModelAnswer(row.Answer)
		answers[i].User = toModelUser(row.User)
	}
	return answers, nil
}

// SetBestAnswer は既存のベストアンサーの解除・指定回答の設定・質問の解決済み化を
// 同一トランザクションで行う。質問行を FOR UPDATE でロックして直列化し、
// 同時設定で複数の回答に is_best が立つ競合を防ぐ。
func (r *answerRepository) SetBestAnswer(ctx context.Context, questionID, answerID uint) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	if _, err := q.LockQuestionForAnswerChange(ctx, int64(questionID)); err != nil {
		return err
	}
	if err := q.ClearBestAnswer(ctx, int64(questionID)); err != nil {
		return err
	}
	if err := q.SetAnswerBest(ctx, int64(answerID)); err != nil {
		return err
	}
	if err := q.SetQuestionSolved(ctx, int64(questionID)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// Vote は投票行の作成・更新と vote_count の増減を同一トランザクションで行う。
// 既に投票済みなら値を更新し、回答の投票数も差分だけ増減させる。
// 差分計算が並行実行で古い値を読まないよう、先に回答行をロックして直列化する。
func (r *answerRepository) Vote(ctx context.Context, userID, answerID uint, value int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	if _, err := q.LockAnswerForVoteChange(ctx, int64(answerID)); err != nil {
		return err
	}

	existing, err := q.GetAnswerVoteByUserAndAnswer(ctx, sqlcgen.GetAnswerVoteByUserAndAnswerParams{
		UserID:   int64(userID),
		AnswerID: int64(answerID),
	})
	switch {
	case err == nil:
		diff := int64(value) - existing.Value
		if err := q.UpdateAnswerVoteValue(ctx, sqlcgen.UpdateAnswerVoteValueParams{
			UserID:   int64(userID),
			AnswerID: int64(answerID),
			Value:    int64(value),
		}); err != nil {
			return err
		}
		if err := q.AdjustAnswerVoteCount(ctx, sqlcgen.AdjustAnswerVoteCountParams{
			ID:   int64(answerID),
			Diff: diff,
		}); err != nil {
			return err
		}
	case isNoRows(err):
		if _, err := q.CreateAnswerVote(ctx, sqlcgen.CreateAnswerVoteParams{
			UserID:   int64(userID),
			AnswerID: int64(answerID),
			Value:    int64(value),
		}); err != nil {
			return err
		}
		if err := q.AdjustAnswerVoteCount(ctx, sqlcgen.AdjustAnswerVoteCountParams{
			ID:   int64(answerID),
			Diff: int64(value),
		}); err != nil {
			return err
		}
	default:
		return err
	}
	return tx.Commit(ctx)
}

// RemoveVote は投票の削除と vote_count の減算を同一トランザクションで行う。
// Vote と同じく回答行を先にロックし、並行する投票変更との差分ずれを防ぐ。
func (r *answerRepository) RemoveVote(ctx context.Context, userID, answerID uint) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	if _, err := q.LockAnswerForVoteChange(ctx, int64(answerID)); err != nil {
		return err
	}

	existing, err := q.GetAnswerVoteByUserAndAnswer(ctx, sqlcgen.GetAnswerVoteByUserAndAnswerParams{
		UserID:   int64(userID),
		AnswerID: int64(answerID),
	})
	if err != nil {
		return err
	}

	if err := q.DeleteAnswerVote(ctx, sqlcgen.DeleteAnswerVoteParams{
		UserID:   int64(userID),
		AnswerID: int64(answerID),
	}); err != nil {
		return err
	}
	if err := q.AdjustAnswerVoteCount(ctx, sqlcgen.AdjustAnswerVoteCountParams{
		ID:   int64(answerID),
		Diff: -existing.Value,
	}); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

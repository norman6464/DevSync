package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// questionRepository は [repository.QuestionRepository] の sqlc(pgx) 実装。
// Vote/RemoveVote は投票行の作成・更新・削除と vote_count の増減を
// 1トランザクションで行うため、Queries だけでなくトランザクションを開始できる
// *pgxpool.Pool を直接保持する。
type questionRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewQuestionRepository は QuestionRepository の sqlc(pgx) 実装を返す。
func NewQuestionRepository(pool *pgxpool.Pool) repository.QuestionRepository {
	return &questionRepository{pool: pool, q: sqlcgen.New(pool)}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.QuestionRepository = (*questionRepository)(nil)

func toModelQuestion(row sqlcgen.Question) model.Question {
	return model.Question{
		ID:          uint(row.ID),
		UserID:      uint(row.UserID),
		Title:       row.Title,
		Body:        row.Body,
		Tags:        fromStringPtr(row.Tags),
		VoteCount:   int(fromInt64PtrValue(row.VoteCount)),
		AnswerCount: int(fromInt64PtrValue(row.AnswerCount)),
		IsSolved:    fromBoolPtr(row.IsSolved),
		CreatedAt:   timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:   timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// Create は新しい質問を作成する。
func (r *questionRepository) Create(ctx context.Context, question *model.Question) error {
	row, err := r.q.CreateQuestion(ctx, sqlcgen.CreateQuestionParams{
		UserID:      int64(question.UserID),
		Title:       question.Title,
		Body:        question.Body,
		Tags:        &question.Tags,
		VoteCount:   toInt64Ptr(question.VoteCount),
		AnswerCount: toInt64Ptr(question.AnswerCount),
		IsSolved:    &question.IsSolved,
	})
	if err != nil {
		return err
	}
	*question = toModelQuestion(row)
	return nil
}

// FindByID は指定 ID の質問をユーザー情報付きで取得する。不在の場合は (nil, nil) を返す。
func (r *questionRepository) FindByID(ctx context.Context, id uint) (*model.Question, error) {
	row, err := r.q.GetQuestionWithUserByID(ctx, int64(id))
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	question := toModelQuestion(row.Question)
	question.User = toModelUser(row.User)
	return &question, nil
}

// Update は既存の質問を更新する（GORMのSave＝全カラム上書きに相当）。
func (r *questionRepository) Update(ctx context.Context, question *model.Question) error {
	row, err := r.q.UpdateQuestion(ctx, sqlcgen.UpdateQuestionParams{
		ID:          int64(question.ID),
		Title:       question.Title,
		Body:        question.Body,
		Tags:        &question.Tags,
		VoteCount:   toInt64Ptr(question.VoteCount),
		AnswerCount: toInt64Ptr(question.AnswerCount),
		IsSolved:    &question.IsSolved,
	})
	if err != nil {
		return err
	}
	*question = toModelQuestion(row)
	return nil
}

// Delete は質問を論理削除する。
func (r *questionRepository) Delete(ctx context.Context, id uint) error {
	return r.q.DeleteQuestion(ctx, int64(id))
}

// FindAll は質問一覧をタグ絞り込み・ソート・ページネーション付きで取得する。
// sort が "votes" なら投票数降順、"unanswered" なら回答が 0 件のものだけに絞る。
func (r *questionRepository) FindAll(ctx context.Context, limit, offset int, tag, sort string) ([]model.Question, int64, error) {
	var tagPattern *string
	if tag != "" {
		// タグは JSON 配列形式の文字列なので、引用符ごと含む部分一致で絞り込む。
		pattern := "%\"" + escapeLikeChars(tag) + "\"%"
		tagPattern = &pattern
	}

	total, err := r.q.CountQuestions(ctx, sqlcgen.CountQuestionsParams{
		TagPattern: tagPattern,
		Sort:       sort,
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListQuestionsWithUser(ctx, sqlcgen.ListQuestionsWithUserParams{
		Limit:      int32Param(limit),
		Offset:     int32Param(offset),
		TagPattern: tagPattern,
		Sort:       sort,
	})
	if err != nil {
		return nil, 0, err
	}

	questions := make([]model.Question, len(rows))
	for i, row := range rows {
		questions[i] = toModelQuestion(row.Question)
		questions[i].User = toModelUser(row.User)
	}
	return questions, total, nil
}

// Search は質問をタイトル・本文・タグで部分一致検索する（投票数降順）。
func (r *questionRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.Question, int64, error) {
	pattern := escapeLikePattern(query)

	total, err := r.q.CountSearchQuestions(ctx, pattern)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.SearchQuestionsWithUser(ctx, sqlcgen.SearchQuestionsWithUserParams{
		Title:  pattern,
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	questions := make([]model.Question, len(rows))
	for i, row := range rows {
		questions[i] = toModelQuestion(row.Question)
		questions[i].User = toModelUser(row.User)
	}
	return questions, total, nil
}

// FindByUserID は指定ユーザーの質問をページネーション付きで取得する（新しい順）。
// 一覧系の中でこれだけユーザー情報をプリロードしない（移行前からの挙動）。
func (r *questionRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Question, int64, error) {
	total, err := r.q.CountQuestionsByUser(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListQuestionsByUser(ctx, sqlcgen.ListQuestionsByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	questions := make([]model.Question, len(rows))
	for i, row := range rows {
		questions[i] = toModelQuestion(row)
	}
	return questions, total, nil
}

// FindSolved は解決済みの質問をページネーション付きで取得する（新しい順）。
func (r *questionRepository) FindSolved(ctx context.Context, limit, offset int) ([]model.Question, int64, error) {
	total, err := r.q.CountSolvedQuestions(ctx)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListSolvedQuestionsWithUser(ctx, sqlcgen.ListSolvedQuestionsWithUserParams{
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	questions := make([]model.Question, len(rows))
	for i, row := range rows {
		questions[i] = toModelQuestion(row.Question)
		questions[i].User = toModelUser(row.User)
	}
	return questions, total, nil
}

// FindUnanswered は回答が 0 件の質問をページネーション付きで取得する（新しい順）。
func (r *questionRepository) FindUnanswered(ctx context.Context, limit, offset int) ([]model.Question, int64, error) {
	total, err := r.q.CountUnansweredQuestions(ctx)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListUnansweredQuestionsWithUser(ctx, sqlcgen.ListUnansweredQuestionsWithUserParams{
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	questions := make([]model.Question, len(rows))
	for i, row := range rows {
		questions[i] = toModelQuestion(row.Question)
		questions[i].User = toModelUser(row.User)
	}
	return questions, total, nil
}

// FindBookmarkedByUserID は指定ユーザーがブックマークした質問を取得する（新しい順）。
func (r *questionRepository) FindBookmarkedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Question, int64, error) {
	total, err := r.q.CountBookmarkedQuestions(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListBookmarkedQuestionsWithUser(ctx, sqlcgen.ListBookmarkedQuestionsWithUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	questions := make([]model.Question, len(rows))
	for i, row := range rows {
		questions[i] = toModelQuestion(row.Question)
		questions[i].User = toModelUser(row.User)
	}
	return questions, total, nil
}

// CountByUserID は指定ユーザーの質問総数を返す。
func (r *questionRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountQuestionsByUser(ctx, int64(userID))
}

// Vote は質問に投票する。既に投票済みなら値を更新し、質問の投票数も差分だけ増減させる。
// 投票行の作成・更新と vote_count の増減を同一トランザクションで行い、
// 片方だけが反映されて投票行と vote_count がずれた状態を残さない。
func (r *questionRepository) Vote(ctx context.Context, userID, questionID uint, value int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	existing, err := q.GetQuestionVoteByUserAndQuestion(ctx, sqlcgen.GetQuestionVoteByUserAndQuestionParams{
		UserID:     int64(userID),
		QuestionID: int64(questionID),
	})
	switch {
	case err == nil:
		diff := int64(value) - existing.Value
		if err := q.UpdateQuestionVoteValue(ctx, sqlcgen.UpdateQuestionVoteValueParams{
			UserID:     int64(userID),
			QuestionID: int64(questionID),
			Value:      int64(value),
		}); err != nil {
			return err
		}
		if err := q.AdjustQuestionVoteCount(ctx, sqlcgen.AdjustQuestionVoteCountParams{
			ID:   int64(questionID),
			Diff: diff,
		}); err != nil {
			return err
		}
	case isNoRows(err):
		if _, err := q.CreateQuestionVote(ctx, sqlcgen.CreateQuestionVoteParams{
			UserID:     int64(userID),
			QuestionID: int64(questionID),
			Value:      int64(value),
		}); err != nil {
			return err
		}
		if err := q.AdjustQuestionVoteCount(ctx, sqlcgen.AdjustQuestionVoteCountParams{
			ID:   int64(questionID),
			Diff: int64(value),
		}); err != nil {
			return err
		}
	default:
		return err
	}
	return tx.Commit(ctx)
}

// RemoveVote は投票を取り消し、質問の投票数から元の値を差し引く。
// 未投票の場合は何もせず成功する（冪等）。削除と vote_count の減算は
// 同一トランザクションで行い、ずれを残さない。
func (r *questionRepository) RemoveVote(ctx context.Context, userID, questionID uint) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	existing, err := q.GetQuestionVoteByUserAndQuestion(ctx, sqlcgen.GetQuestionVoteByUserAndQuestionParams{
		UserID:     int64(userID),
		QuestionID: int64(questionID),
	})
	if isNoRows(err) {
		return nil
	}
	if err != nil {
		return err
	}

	if err := q.DeleteQuestionVote(ctx, sqlcgen.DeleteQuestionVoteParams{
		UserID:     int64(userID),
		QuestionID: int64(questionID),
	}); err != nil {
		return err
	}
	if err := q.AdjustQuestionVoteCount(ctx, sqlcgen.AdjustQuestionVoteCountParams{
		ID:   int64(questionID),
		Diff: -existing.Value,
	}); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// GetUserVote は指定ユーザーの投票値を返す。未投票の場合は 0 を返す。
func (r *questionRepository) GetUserVote(ctx context.Context, userID, questionID uint) (int, error) {
	vote, err := r.q.GetQuestionVoteByUserAndQuestion(ctx, sqlcgen.GetQuestionVoteByUserAndQuestionParams{
		UserID:     int64(userID),
		QuestionID: int64(questionID),
	})
	if isNoRows(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return int(vote.Value), nil
}

// Bookmark は質問をブックマークする。
func (r *questionRepository) Bookmark(ctx context.Context, userID, questionID uint) error {
	return r.q.CreateQuestionBookmark(ctx, sqlcgen.CreateQuestionBookmarkParams{
		UserID:     int64(userID),
		QuestionID: int64(questionID),
	})
}

// Unbookmark は質問のブックマークを解除する。
func (r *questionRepository) Unbookmark(ctx context.Context, userID, questionID uint) error {
	return r.q.DeleteQuestionBookmark(ctx, sqlcgen.DeleteQuestionBookmarkParams{
		UserID:     int64(userID),
		QuestionID: int64(questionID),
	})
}

// HasBookmarked は指定ユーザーがブックマーク済みかを返す。
func (r *questionRepository) HasBookmarked(ctx context.Context, userID, questionID uint) (bool, error) {
	count, err := r.q.CountQuestionBookmark(ctx, sqlcgen.CountQuestionBookmarkParams{
		UserID:     int64(userID),
		QuestionID: int64(questionID),
	})
	return count > 0, err
}

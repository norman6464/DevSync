package usecase

import (
	"context"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/domain/validator"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// requireVotableQuestion は投票対象の質問を取得し、自分の質問でないことを検証する。
// 不在の場合は 404、自分の質問なら 403 を返す。
func requireVotableQuestion(ctx context.Context, questions repository.QuestionRepository, userID, questionID uint) error {
	question, err := questions.FindByID(ctx, questionID)
	if err != nil || question == nil {
		return domain.ErrNotFound
	}
	if question.UserID == userID {
		return domain.ErrForbidden
	}
	return nil
}

// CreateQuestionUseCase は質問を作成する。
type CreateQuestionUseCase struct {
	questions repository.QuestionRepository
}

// NewCreateQuestionUseCase は CreateQuestionUseCase を生成する。
func NewCreateQuestionUseCase(questions repository.QuestionRepository) *CreateQuestionUseCase {
	return &CreateQuestionUseCase{questions: questions}
}

// Execute はタイトル・本文・タグを検証したうえで質問を作成する。
func (uc *CreateQuestionUseCase) Execute(ctx context.Context, question *model.Question) error {
	v := validator.NewQuestionValidator()
	if err := v.ValidateCreateQuestion(question.Title, question.Body, question.Tags); err != nil {
		return err
	}
	return uc.questions.Create(ctx, question)
}

// ListQuestionsUseCase は質問一覧を取得する。
type ListQuestionsUseCase struct {
	questions repository.QuestionRepository
}

// NewListQuestionsUseCase は ListQuestionsUseCase を生成する。
func NewListQuestionsUseCase(questions repository.QuestionRepository) *ListQuestionsUseCase {
	return &ListQuestionsUseCase{questions: questions}
}

// Execute はタグ絞り込み・ソート付きで質問一覧を返す。
func (uc *ListQuestionsUseCase) Execute(ctx context.Context, limit, offset int, tag, sort string) ([]model.Question, int64, error) {
	return uc.questions.FindAll(ctx, limit, offset, tag, sort)
}

// SearchQuestionsUseCase は質問をキーワード検索する。
type SearchQuestionsUseCase struct {
	questions repository.QuestionRepository
}

// NewSearchQuestionsUseCase は SearchQuestionsUseCase を生成する。
func NewSearchQuestionsUseCase(questions repository.QuestionRepository) *SearchQuestionsUseCase {
	return &SearchQuestionsUseCase{questions: questions}
}

// Execute はタイトル・本文・タグへの部分一致で質問を検索する。
func (uc *SearchQuestionsUseCase) Execute(ctx context.Context, query string, limit, offset int) ([]model.Question, int64, error) {
	return uc.questions.Search(ctx, query, limit, offset)
}

// GetQuestionUseCase は質問を 1 件取得する。
type GetQuestionUseCase struct {
	questions repository.QuestionRepository
}

// NewGetQuestionUseCase は GetQuestionUseCase を生成する。
func NewGetQuestionUseCase(questions repository.QuestionRepository) *GetQuestionUseCase {
	return &GetQuestionUseCase{questions: questions}
}

// Execute は指定 ID の質問を返す。不在なら 404 を返す。
func (uc *GetQuestionUseCase) Execute(ctx context.Context, id uint) (*model.Question, error) {
	question, err := uc.questions.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if question == nil {
		return nil, domain.ErrNotFound
	}
	return question, nil
}

// ListQuestionsByUserUseCase は指定ユーザーの質問一覧を取得する。
type ListQuestionsByUserUseCase struct {
	questions repository.QuestionRepository
}

// NewListQuestionsByUserUseCase は ListQuestionsByUserUseCase を生成する。
func NewListQuestionsByUserUseCase(questions repository.QuestionRepository) *ListQuestionsByUserUseCase {
	return &ListQuestionsByUserUseCase{questions: questions}
}

// Execute は指定ユーザーの質問を新しい順で返す。
func (uc *ListQuestionsByUserUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.Question, int64, error) {
	return uc.questions.FindByUserID(ctx, userID, limit, offset)
}

// GetQuestionUserVoteUseCase は指定ユーザーの投票値を取得する。
type GetQuestionUserVoteUseCase struct {
	questions repository.QuestionRepository
}

// NewGetQuestionUserVoteUseCase は GetQuestionUserVoteUseCase を生成する。
func NewGetQuestionUserVoteUseCase(questions repository.QuestionRepository) *GetQuestionUserVoteUseCase {
	return &GetQuestionUserVoteUseCase{questions: questions}
}

// Execute は投票値を返す。未投票の場合は 0 を返す。
func (uc *GetQuestionUserVoteUseCase) Execute(ctx context.Context, userID, questionID uint) (int, error) {
	return uc.questions.GetUserVote(ctx, userID, questionID)
}

// UpdateQuestionUseCase は質問を更新する。
type UpdateQuestionUseCase struct {
	questions repository.QuestionRepository
}

// NewUpdateQuestionUseCase は UpdateQuestionUseCase を生成する。
func NewUpdateQuestionUseCase(questions repository.QuestionRepository) *UpdateQuestionUseCase {
	return &UpdateQuestionUseCase{questions: questions}
}

// Execute はタイトル・本文・タグを部分更新する。所有者のみ。
// 前後の空白を除いて空になった項目は「変更なし」として扱う。
func (uc *UpdateQuestionUseCase) Execute(ctx context.Context, id, userID uint, title, body, tags string) (*model.Question, error) {
	question, err := ensureOwner(ctx, uc.questions.FindByID, id, userID,
		func(q *model.Question) uint { return q.UserID })
	if err != nil {
		return nil, err
	}

	title = strings.TrimSpace(title)
	body = strings.TrimSpace(body)
	tags = strings.TrimSpace(tags)

	v := validator.NewQuestionValidator()
	if err := v.ValidateUpdateQuestion(title, body, tags); err != nil {
		return nil, err
	}

	if title != "" {
		question.Title = title
	}
	if body != "" {
		question.Body = body
	}
	if tags != "" {
		question.Tags = tags
	}

	if err := uc.questions.Update(ctx, question); err != nil {
		return nil, err
	}
	return question, nil
}

// DeleteQuestionUseCase は質問を削除する。
type DeleteQuestionUseCase struct {
	questions repository.QuestionRepository
}

// NewDeleteQuestionUseCase は DeleteQuestionUseCase を生成する。
func NewDeleteQuestionUseCase(questions repository.QuestionRepository) *DeleteQuestionUseCase {
	return &DeleteQuestionUseCase{questions: questions}
}

// Execute は質問を削除する。所有者のみ。
func (uc *DeleteQuestionUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensureOwner(ctx, uc.questions.FindByID, id, userID,
		func(q *model.Question) uint { return q.UserID }); err != nil {
		return err
	}
	return uc.questions.Delete(ctx, id)
}

// VoteQuestionUseCase は質問に投票する。
type VoteQuestionUseCase struct {
	questions repository.QuestionRepository
}

// NewVoteQuestionUseCase は VoteQuestionUseCase を生成する。
func NewVoteQuestionUseCase(questions repository.QuestionRepository) *VoteQuestionUseCase {
	return &VoteQuestionUseCase{questions: questions}
}

// Execute は質問に投票する。投票値は +1 か -1 のみで、自分の質問には投票できない。
func (uc *VoteQuestionUseCase) Execute(ctx context.Context, userID, questionID uint, value int) error {
	v := validator.NewQuestionValidator()
	if err := v.ValidateVote(value); err != nil {
		return err
	}
	if err := requireVotableQuestion(ctx, uc.questions, userID, questionID); err != nil {
		return err
	}
	return uc.questions.Vote(ctx, userID, questionID, value)
}

// RemoveQuestionVoteUseCase は質問への投票を取り消す。
type RemoveQuestionVoteUseCase struct {
	questions repository.QuestionRepository
}

// NewRemoveQuestionVoteUseCase は RemoveQuestionVoteUseCase を生成する。
func NewRemoveQuestionVoteUseCase(questions repository.QuestionRepository) *RemoveQuestionVoteUseCase {
	return &RemoveQuestionVoteUseCase{questions: questions}
}

// Execute は投票を取り消す。自分の質問はそもそも投票できないため 403 を返す。
func (uc *RemoveQuestionVoteUseCase) Execute(ctx context.Context, userID, questionID uint) error {
	if err := requireVotableQuestion(ctx, uc.questions, userID, questionID); err != nil {
		return err
	}
	return uc.questions.RemoveVote(ctx, userID, questionID)
}

// ListSolvedQuestionsUseCase は解決済みの質問一覧を取得する。
type ListSolvedQuestionsUseCase struct {
	questions repository.QuestionRepository
}

// NewListSolvedQuestionsUseCase は ListSolvedQuestionsUseCase を生成する。
func NewListSolvedQuestionsUseCase(questions repository.QuestionRepository) *ListSolvedQuestionsUseCase {
	return &ListSolvedQuestionsUseCase{questions: questions}
}

// Execute は解決済みの質問を新しい順で返す。
func (uc *ListSolvedQuestionsUseCase) Execute(ctx context.Context, limit, offset int) ([]model.Question, int64, error) {
	return uc.questions.FindSolved(ctx, limit, offset)
}

// ListUnansweredQuestionsUseCase は未回答の質問一覧を取得する。
type ListUnansweredQuestionsUseCase struct {
	questions repository.QuestionRepository
}

// NewListUnansweredQuestionsUseCase は ListUnansweredQuestionsUseCase を生成する。
func NewListUnansweredQuestionsUseCase(questions repository.QuestionRepository) *ListUnansweredQuestionsUseCase {
	return &ListUnansweredQuestionsUseCase{questions: questions}
}

// Execute は回答が 0 件の質問を新しい順で返す。
func (uc *ListUnansweredQuestionsUseCase) Execute(ctx context.Context, limit, offset int) ([]model.Question, int64, error) {
	return uc.questions.FindUnanswered(ctx, limit, offset)
}

// BookmarkQuestionUseCase は質問をブックマークする。
type BookmarkQuestionUseCase struct {
	questions repository.QuestionRepository
}

// NewBookmarkQuestionUseCase は BookmarkQuestionUseCase を生成する。
func NewBookmarkQuestionUseCase(questions repository.QuestionRepository) *BookmarkQuestionUseCase {
	return &BookmarkQuestionUseCase{questions: questions}
}

// Execute は質問をブックマークする。不在なら 404、既にブックマーク済みなら 409。
func (uc *BookmarkQuestionUseCase) Execute(ctx context.Context, userID, questionID uint) error {
	question, err := uc.questions.FindByID(ctx, questionID)
	if err != nil || question == nil {
		return domain.ErrNotFound
	}

	has, err := uc.questions.HasBookmarked(ctx, userID, questionID)
	if err != nil {
		return err
	}
	if has {
		return domain.ErrConflict
	}
	return uc.questions.Bookmark(ctx, userID, questionID)
}

// UnbookmarkQuestionUseCase は質問のブックマークを解除する。
type UnbookmarkQuestionUseCase struct {
	questions repository.QuestionRepository
}

// NewUnbookmarkQuestionUseCase は UnbookmarkQuestionUseCase を生成する。
func NewUnbookmarkQuestionUseCase(questions repository.QuestionRepository) *UnbookmarkQuestionUseCase {
	return &UnbookmarkQuestionUseCase{questions: questions}
}

// Execute はブックマークを解除する。対象が無くても成功として扱う（移行前からの挙動）。
func (uc *UnbookmarkQuestionUseCase) Execute(ctx context.Context, userID, questionID uint) error {
	return uc.questions.Unbookmark(ctx, userID, questionID)
}

// ListBookmarkedQuestionsUseCase はブックマーク済みの質問一覧を取得する。
type ListBookmarkedQuestionsUseCase struct {
	questions repository.QuestionRepository
}

// NewListBookmarkedQuestionsUseCase は ListBookmarkedQuestionsUseCase を生成する。
func NewListBookmarkedQuestionsUseCase(questions repository.QuestionRepository) *ListBookmarkedQuestionsUseCase {
	return &ListBookmarkedQuestionsUseCase{questions: questions}
}

// Execute はブックマーク済みの質問を新しい順で返す。
func (uc *ListBookmarkedQuestionsUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.Question, int64, error) {
	return uc.questions.FindBookmarkedByUserID(ctx, userID, limit, offset)
}

// CountQuestionsUseCase は指定ユーザーの質問数を取得する。
type CountQuestionsUseCase struct {
	questions repository.QuestionRepository
}

// NewCountQuestionsUseCase は CountQuestionsUseCase を生成する。
func NewCountQuestionsUseCase(questions repository.QuestionRepository) *CountQuestionsUseCase {
	return &CountQuestionsUseCase{questions: questions}
}

// Execute は指定ユーザーの質問総数を返す。
func (uc *CountQuestionsUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.questions.CountByUserID(ctx, userID)
}

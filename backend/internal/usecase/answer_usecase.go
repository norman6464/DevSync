package usecase

import (
	"context"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/domain/validator"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// requireVotableAnswer は投票対象の回答を取得し、自分の回答でないことを検証する。
// 不在の場合は 404、自分の回答なら 403 を返す。
func requireVotableAnswer(ctx context.Context, answers repository.AnswerRepository, userID, answerID uint) error {
	answer, err := answers.FindByID(ctx, answerID)
	if err != nil || answer == nil {
		return domain.ErrNotFound
	}
	if answer.UserID == userID {
		return domain.ErrForbidden
	}
	return nil
}

// ListAnswersUseCase は質問に紐づく回答一覧を取得する。
type ListAnswersUseCase struct {
	answers repository.AnswerRepository
}

// NewListAnswersUseCase は ListAnswersUseCase を生成する。
func NewListAnswersUseCase(answers repository.AnswerRepository) *ListAnswersUseCase {
	return &ListAnswersUseCase{answers: answers}
}

// Execute は指定質問の回答をベストアンサー優先・投票数降順で返す。
func (uc *ListAnswersUseCase) Execute(ctx context.Context, questionID uint) ([]model.Answer, error) {
	return uc.answers.FindByQuestionID(ctx, questionID)
}

// CreateAnswerUseCase は質問に回答を投稿する。
type CreateAnswerUseCase struct {
	answers   repository.AnswerRepository
	questions repository.QuestionReader
}

// NewCreateAnswerUseCase は CreateAnswerUseCase を生成する。
func NewCreateAnswerUseCase(answers repository.AnswerRepository, questions repository.QuestionReader) *CreateAnswerUseCase {
	return &CreateAnswerUseCase{answers: answers, questions: questions}
}

// Execute は本文を検証し、質問の存在を確認したうえで回答を作成する。
func (uc *CreateAnswerUseCase) Execute(ctx context.Context, answer *model.Answer) error {
	if err := domain.ValidateStringLength(answer.Body, 1, 10000, "回答内容"); err != nil {
		return err
	}
	answer.Body = strings.TrimSpace(answer.Body)

	question, err := uc.questions.FindByID(ctx, answer.QuestionID)
	if err != nil || question == nil {
		return domain.ErrNotFound
	}
	return uc.answers.Create(ctx, answer)
}

// UpdateAnswerUseCase は回答を更新する。
type UpdateAnswerUseCase struct {
	answers repository.AnswerRepository
}

// NewUpdateAnswerUseCase は UpdateAnswerUseCase を生成する。
func NewUpdateAnswerUseCase(answers repository.AnswerRepository) *UpdateAnswerUseCase {
	return &UpdateAnswerUseCase{answers: answers}
}

// Execute は回答の本文を更新する。所有者のみ。
func (uc *UpdateAnswerUseCase) Execute(ctx context.Context, answerID, userID uint, body string) (*model.Answer, error) {
	answer, err := ensureOwner(ctx, uc.answers.FindByID, answerID, userID,
		func(a *model.Answer) uint { return a.UserID })
	if err != nil {
		return nil, err
	}

	if err := domain.ValidateStringLength(body, 1, 10000, "回答内容"); err != nil {
		return nil, err
	}
	answer.Body = strings.TrimSpace(body)

	if err := uc.answers.Update(ctx, answer); err != nil {
		return nil, err
	}
	return answer, nil
}

// DeleteAnswerUseCase は回答を削除する。
type DeleteAnswerUseCase struct {
	answers repository.AnswerRepository
}

// NewDeleteAnswerUseCase は DeleteAnswerUseCase を生成する。
func NewDeleteAnswerUseCase(answers repository.AnswerRepository) *DeleteAnswerUseCase {
	return &DeleteAnswerUseCase{answers: answers}
}

// Execute は回答を論理削除する。所有者のみ。
func (uc *DeleteAnswerUseCase) Execute(ctx context.Context, answerID, userID uint) error {
	answer, err := ensureOwner(ctx, uc.answers.FindByID, answerID, userID,
		func(a *model.Answer) uint { return a.UserID })
	if err != nil {
		return err
	}
	return uc.answers.Delete(ctx, answer)
}

// SetBestAnswerUseCase はベストアンサーを設定する。
type SetBestAnswerUseCase struct {
	answers   repository.AnswerRepository
	questions repository.QuestionReader
}

// NewSetBestAnswerUseCase は SetBestAnswerUseCase を生成する。
func NewSetBestAnswerUseCase(answers repository.AnswerRepository, questions repository.QuestionReader) *SetBestAnswerUseCase {
	return &SetBestAnswerUseCase{answers: answers, questions: questions}
}

// Execute はベストアンサーを設定する。質問の投稿者のみが設定でき、
// 指定された回答がその質問のものであることも確認する。
func (uc *SetBestAnswerUseCase) Execute(ctx context.Context, questionID, answerID, userID uint) error {
	question, err := uc.questions.FindByID(ctx, questionID)
	if err != nil || question == nil {
		return domain.ErrNotFound
	}
	if question.UserID != userID {
		return domain.ErrForbidden
	}

	answer, err := uc.answers.FindByID(ctx, answerID)
	if err != nil || answer == nil {
		return domain.ErrNotFound
	}
	if answer.QuestionID != questionID {
		return domain.ErrBadRequest
	}

	return uc.answers.SetBestAnswer(ctx, questionID, answerID)
}

// VoteAnswerUseCase は回答に投票する。
type VoteAnswerUseCase struct {
	answers repository.AnswerRepository
}

// NewVoteAnswerUseCase は VoteAnswerUseCase を生成する。
func NewVoteAnswerUseCase(answers repository.AnswerRepository) *VoteAnswerUseCase {
	return &VoteAnswerUseCase{answers: answers}
}

// Execute は回答に投票する。投票値は +1 か -1 のみで、自分の回答には投票できない。
func (uc *VoteAnswerUseCase) Execute(ctx context.Context, userID, answerID uint, value int) error {
	v := validator.NewQuestionValidator()
	if err := v.ValidateVote(value); err != nil {
		return err
	}
	if err := requireVotableAnswer(ctx, uc.answers, userID, answerID); err != nil {
		return err
	}
	return uc.answers.Vote(ctx, userID, answerID, value)
}

// RemoveAnswerVoteUseCase は回答への投票を取り消す。
type RemoveAnswerVoteUseCase struct {
	answers repository.AnswerRepository
}

// NewRemoveAnswerVoteUseCase は RemoveAnswerVoteUseCase を生成する。
func NewRemoveAnswerVoteUseCase(answers repository.AnswerRepository) *RemoveAnswerVoteUseCase {
	return &RemoveAnswerVoteUseCase{answers: answers}
}

// Execute は投票を取り消す。自分の回答はそもそも投票できないため 403 を返す。
func (uc *RemoveAnswerVoteUseCase) Execute(ctx context.Context, userID, answerID uint) error {
	if err := requireVotableAnswer(ctx, uc.answers, userID, answerID); err != nil {
		return err
	}
	return uc.answers.RemoveVote(ctx, userID, answerID)
}

// ListAnswersByVoteRangeUseCase は回答を投票数の範囲で絞り込んで取得する。
type ListAnswersByVoteRangeUseCase struct {
	answers repository.AnswerRepository
}

// NewListAnswersByVoteRangeUseCase は ListAnswersByVoteRangeUseCase を生成する。
func NewListAnswersByVoteRangeUseCase(answers repository.AnswerRepository) *ListAnswersByVoteRangeUseCase {
	return &ListAnswersByVoteRangeUseCase{answers: answers}
}

// Execute は投票数が範囲内の回答を返す。下限が上限を上回る場合は 400。
func (uc *ListAnswersByVoteRangeUseCase) Execute(ctx context.Context, questionID uint, minVote, maxVote int) ([]model.Answer, error) {
	if minVote > maxVote {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "投票範囲が無効です", nil)
	}
	return uc.answers.FindByVoteRange(ctx, questionID, minVote, maxVote)
}

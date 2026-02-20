package service

import (
	"github.com/norman6464/devsync/backend/internal/domain/validator"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// AnswerService はQ&A回答のビジネスロジックを提供する。
// 回答のCRUD操作、ベストアンサー設定、投票機能を担当する。
type AnswerService struct {
	answerRepo   repository.AnswerRepositoryInterface
	questionRepo repository.QuestionRepositoryInterface
}

// NewAnswerService は新しいAnswerServiceインスタンスを生成する。
func NewAnswerService(answerRepo repository.AnswerRepositoryInterface, questionRepo repository.QuestionRepositoryInterface) *AnswerService {
	return &AnswerService{answerRepo: answerRepo, questionRepo: questionRepo}
}

// GetByQuestionID は指定質問の全回答を取得する。
func (s *AnswerService) GetByQuestionID(questionID uint) ([]model.Answer, error) {
	return s.answerRepo.FindByQuestionID(questionID)
}

// Create は質問の存在を確認した後、新しい回答を作成する。
func (s *AnswerService) Create(answer *model.Answer) error {
	if _, err := s.questionRepo.FindByID(answer.QuestionID); err != nil {
		return ErrNotFound
	}
	return s.answerRepo.Create(answer)
}

// findAndCheckOwnership は回答を取得し、指定ユーザーが所有者かを検証する。
func (s *AnswerService) findAndCheckOwnership(answerID, userID uint) (*model.Answer, error) {
	answer, err := s.answerRepo.FindByID(answerID)
	if err != nil {
		return nil, err
	}
	if answer.UserID != userID {
		return nil, ErrForbidden
	}
	return answer, nil
}

// Update は所有権を検証した後、回答を更新する。
func (s *AnswerService) Update(answerID, userID uint, body string) (*model.Answer, error) {
	answer, err := s.findAndCheckOwnership(answerID, userID)
	if err != nil {
		return nil, err
	}
	answer.Body = body
	if err := s.answerRepo.Update(answer); err != nil {
		return nil, err
	}
	return answer, nil
}

// Delete は所有権を検証した後、回答を削除する。
func (s *AnswerService) Delete(answerID, userID uint) error {
	answer, err := s.findAndCheckOwnership(answerID, userID)
	if err != nil {
		return err
	}
	return s.answerRepo.Delete(answer)
}

// SetBestAnswer は質問の所有権を検証し、指定回答をベストアンサーに設定する。
// 回答が対象質問に属していることも検証する。
func (s *AnswerService) SetBestAnswer(questionID, answerID, userID uint) error {
	question, err := s.questionRepo.FindByID(questionID)
	if err != nil {
		return ErrNotFound
	}
	if question.UserID != userID {
		return ErrForbidden
	}

	answer, err := s.answerRepo.FindByID(answerID)
	if err != nil {
		return ErrNotFound
	}
	if answer.QuestionID != questionID {
		return ErrBadRequest
	}

	return s.answerRepo.SetBestAnswer(questionID, answerID)
}

// findAndPreventSelfVote は回答を取得し、自己投票でないことを検証する。
func (s *AnswerService) findAndPreventSelfVote(userID, answerID uint) error {
	answer, err := s.answerRepo.FindByID(answerID)
	if err != nil {
		return ErrNotFound
	}
	if answer.UserID == userID {
		return ErrForbidden
	}
	return nil
}

// Vote は投票値を検証した後、回答に投票する。
// valueは1（賛成）または-1（反対）のみ許可される。
// 自分の回答への自己投票は禁止する。
func (s *AnswerService) Vote(userID, answerID uint, value int) error {
	v := validator.NewQuestionValidator()
	if err := v.ValidateVote(value); err != nil {
		return err
	}
	if err := s.findAndPreventSelfVote(userID, answerID); err != nil {
		return err
	}
	return s.answerRepo.Vote(userID, answerID, value)
}

// RemoveVote は回答への投票を取り消す。
// 自分の回答への投票削除は禁止する（そもそも投票できないため）。
func (s *AnswerService) RemoveVote(userID, answerID uint) error {
	if err := s.findAndPreventSelfVote(userID, answerID); err != nil {
		return err
	}
	return s.answerRepo.RemoveVote(userID, answerID)
}

// GetUserVotes は指定ユーザーの複数回答への投票値をマップで取得する。
func (s *AnswerService) GetUserVotes(userID uint, answerIDs []uint) (map[uint]int, error) {
	return s.answerRepo.GetUserVotes(userID, answerIDs)
}

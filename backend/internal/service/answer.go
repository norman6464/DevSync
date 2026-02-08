package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// AnswerService handles answer business logic.
type AnswerService struct {
	answerRepo   repository.AnswerRepositoryInterface
	questionRepo repository.QuestionRepositoryInterface
}

// NewAnswerService creates a new AnswerService.
func NewAnswerService(answerRepo repository.AnswerRepositoryInterface, questionRepo repository.QuestionRepositoryInterface) *AnswerService {
	return &AnswerService{answerRepo: answerRepo, questionRepo: questionRepo}
}

// GetByQuestionID returns answers for a question.
func (s *AnswerService) GetByQuestionID(questionID uint) ([]model.Answer, error) {
	return s.answerRepo.FindByQuestionID(questionID)
}

// Create creates a new answer after verifying the question exists.
func (s *AnswerService) Create(answer *model.Answer) error {
	if _, err := s.questionRepo.FindByID(answer.QuestionID); err != nil {
		return ErrNotFound
	}
	return s.answerRepo.Create(answer)
}

// Update updates an answer after verifying ownership.
func (s *AnswerService) Update(answerID, userID uint, body string) (*model.Answer, error) {
	answer, err := s.answerRepo.FindByID(answerID)
	if err != nil {
		return nil, err
	}
	if answer.UserID != userID {
		return nil, ErrForbidden
	}
	answer.Body = body
	if err := s.answerRepo.Update(answer); err != nil {
		return nil, err
	}
	return answer, nil
}

// Delete deletes an answer after verifying ownership.
func (s *AnswerService) Delete(answerID, userID uint) error {
	answer, err := s.answerRepo.FindByID(answerID)
	if err != nil {
		return err
	}
	if answer.UserID != userID {
		return ErrForbidden
	}
	return s.answerRepo.Delete(answer)
}

// SetBestAnswer sets the best answer for a question after verifying ownership.
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

// Vote votes on an answer.
func (s *AnswerService) Vote(userID, answerID uint, value int) error {
	return s.answerRepo.Vote(userID, answerID, value)
}

// RemoveVote removes a vote from an answer.
func (s *AnswerService) RemoveVote(userID, answerID uint) error {
	return s.answerRepo.RemoveVote(userID, answerID)
}

// GetUserVotes returns user votes for given answer IDs.
func (s *AnswerService) GetUserVotes(userID uint, answerIDs []uint) (map[uint]int, error) {
	return s.answerRepo.GetUserVotes(userID, answerIDs)
}

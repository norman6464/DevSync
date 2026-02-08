package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// QuestionService handles question business logic.
type QuestionService struct {
	repo repository.QuestionRepositoryInterface
}

// NewQuestionService creates a new QuestionService.
func NewQuestionService(repo repository.QuestionRepositoryInterface) *QuestionService {
	return &QuestionService{repo: repo}
}

// Create creates a new question.
func (s *QuestionService) Create(question *model.Question) error {
	return s.repo.Create(question)
}

// GetAll returns paginated questions.
func (s *QuestionService) GetAll(limit, offset int, tag, sort string) ([]model.Question, int64, error) {
	return s.repo.FindAll(limit, offset, tag, sort)
}

// Search searches questions.
func (s *QuestionService) Search(q string, limit, offset int) ([]model.Question, int64, error) {
	return s.repo.Search(q, limit, offset)
}

// GetByID returns a question by ID.
func (s *QuestionService) GetByID(id uint) (*model.Question, error) {
	return s.repo.FindByID(id)
}

// GetByUserID returns questions by user ID.
func (s *QuestionService) GetByUserID(userID uint) ([]model.Question, error) {
	return s.repo.FindByUserID(userID)
}

// GetUserVote returns the user's vote on a question.
func (s *QuestionService) GetUserVote(userID, questionID uint) (int, error) {
	return s.repo.GetUserVote(userID, questionID)
}

// Update updates a question after verifying ownership.
func (s *QuestionService) Update(id, userID uint, title, body, tags string) (*model.Question, error) {
	question, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if question.UserID != userID {
		return nil, ErrForbidden
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

	if err := s.repo.Update(question); err != nil {
		return nil, err
	}
	return question, nil
}

// Delete deletes a question after verifying ownership.
func (s *QuestionService) Delete(id, userID uint) error {
	question, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if question.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}

// Vote votes on a question.
func (s *QuestionService) Vote(userID, questionID uint, value int) error {
	return s.repo.Vote(userID, questionID, value)
}

// RemoveVote removes a vote from a question.
func (s *QuestionService) RemoveVote(userID, questionID uint) error {
	return s.repo.RemoveVote(userID, questionID)
}

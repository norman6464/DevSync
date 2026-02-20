package service

import (
	"github.com/norman6464/devsync/backend/internal/domain/validator"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// QuestionService はQ&A質問のビジネスロジックを提供する。
// 質問のCRUD操作と投票機能を担当する。
type QuestionService struct {
	repo repository.QuestionRepositoryInterface
}

// NewQuestionService は新しいQuestionServiceインスタンスを生成する。
func NewQuestionService(repo repository.QuestionRepositoryInterface) *QuestionService {
	return &QuestionService{repo: repo}
}

// Create は新しい質問を作成する。
func (s *QuestionService) Create(question *model.Question) error {
	v := validator.NewQuestionValidator()
	if err := v.ValidateCreateQuestion(question.Title, question.Body, question.Tags); err != nil {
		return err
	}
	return s.repo.Create(question)
}

// GetAll は質問一覧をフィルタ・ソート・ページネーション付きで取得する。
func (s *QuestionService) GetAll(limit, offset int, tag, sort string) ([]model.Question, int64, error) {
	return s.repo.FindAll(limit, offset, tag, sort)
}

// Search は質問をキーワードで全文検索する。
func (s *QuestionService) Search(q string, limit, offset int) ([]model.Question, int64, error) {
	return s.repo.Search(q, limit, offset)
}

// GetByID は指定IDの質問を取得する。
func (s *QuestionService) GetByID(id uint) (*model.Question, error) {
	return s.repo.FindByID(id)
}

// GetByUserID は指定ユーザーの全質問を取得する。
func (s *QuestionService) GetByUserID(userID uint) ([]model.Question, error) {
	return s.repo.FindByUserID(userID)
}

// GetUserVote は指定ユーザーの質問への投票値を取得する。
func (s *QuestionService) GetUserVote(userID, questionID uint) (int, error) {
	return s.repo.GetUserVote(userID, questionID)
}

// findAndCheckOwnership は質問を取得し、指定ユーザーが所有者かを検証する。
func (s *QuestionService) findAndCheckOwnership(id, userID uint) (*model.Question, error) {
	question, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if question.UserID != userID {
		return nil, ErrForbidden
	}
	return question, nil
}

// Update は所有権を検証した後、質問を更新する。
func (s *QuestionService) Update(id, userID uint, title, body, tags string) (*model.Question, error) {
	question, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

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

	if err := s.repo.Update(question); err != nil {
		return nil, err
	}
	return question, nil
}

// Delete は所有権を検証した後、質問を削除する。
func (s *QuestionService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// findAndPreventSelfVote は質問を取得し、自分の質問への投票を防止する。
func (s *QuestionService) findAndPreventSelfVote(userID, questionID uint) error {
	question, err := s.repo.FindByID(questionID)
	if err != nil {
		return ErrNotFound
	}
	if question.UserID == userID {
		return ErrForbidden
	}
	return nil
}

// Vote は質問に投票する。
// 自分の質問への自己投票は禁止する。
func (s *QuestionService) Vote(userID, questionID uint, value int) error {
	v := validator.NewQuestionValidator()
	if err := v.ValidateVote(value); err != nil {
		return err
	}
	if err := s.findAndPreventSelfVote(userID, questionID); err != nil {
		return err
	}
	return s.repo.Vote(userID, questionID, value)
}

// RemoveVote は質問への投票を取り消す。
// 自分の質問への投票削除は禁止する（そもそも投票できないため）。
func (s *QuestionService) RemoveVote(userID, questionID uint) error {
	if err := s.findAndPreventSelfVote(userID, questionID); err != nil {
		return err
	}
	return s.repo.RemoveVote(userID, questionID)
}

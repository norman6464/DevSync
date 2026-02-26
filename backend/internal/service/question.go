package service

import (
	"strings"

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

// GetByUserID は指定ユーザーの質問をページネーション付きで取得する。
func (s *QuestionService) GetByUserID(userID uint, limit, offset int) ([]model.Question, int64, error) {
	return s.repo.FindByUserID(userID, limit, offset)
}

// GetUserVote は指定ユーザーの質問への投票値を取得する。
func (s *QuestionService) GetUserVote(userID, questionID uint) (int, error) {
	return s.repo.GetUserVote(userID, questionID)
}

// findAndCheckOwnership は質問を取得し、指定ユーザーが所有者かを検証する。
func (s *QuestionService) findAndCheckOwnership(id, userID uint) (*model.Question, error) {
	return checkOwnership(s.repo.FindByID, id, userID, func(q *model.Question) uint { return q.UserID })
}

// Update は所有権を検証した後、質問を更新する。
func (s *QuestionService) Update(id, userID uint, title, body, tags string) (*model.Question, error) {
	question, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	// 空白のみの入力を正規化（空白のみ→変更なしとして扱う）
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

// GetSolved は解決済みの質問一覧をページネーション付きで取得する。
func (s *QuestionService) GetSolved(limit, offset int) ([]model.Question, int64, error) {
	return s.repo.FindSolved(limit, offset)
}

// GetUnanswered は未回答の質問一覧をページネーション付きで取得する。
func (s *QuestionService) GetUnanswered(limit, offset int) ([]model.Question, int64, error) {
	return s.repo.FindUnanswered(limit, offset)
}

// RemoveVote は質問への投票を取り消す。
// 自分の質問への投票削除は禁止する（そもそも投票できないため）。
func (s *QuestionService) RemoveVote(userID, questionID uint) error {
	if err := s.findAndPreventSelfVote(userID, questionID); err != nil {
		return err
	}
	return s.repo.RemoveVote(userID, questionID)
}

// Bookmark は質問をブックマークする。
// 既にブックマーク済みの場合はErrConflictを返す。
func (s *QuestionService) Bookmark(userID, questionID uint) error {
	if _, err := s.repo.FindByID(questionID); err != nil {
		return ErrNotFound
	}
	has, err := s.repo.HasBookmarked(userID, questionID)
	if err != nil {
		return err
	}
	if has {
		return ErrConflict
	}
	return s.repo.Bookmark(userID, questionID)
}

// Unbookmark は質問のブックマークを解除する。
func (s *QuestionService) Unbookmark(userID, questionID uint) error {
	return s.repo.Unbookmark(userID, questionID)
}

// GetBookmarkedByUserID はブックマーク済み質問一覧を取得する。
func (s *QuestionService) GetBookmarkedByUserID(userID uint, limit, offset int) ([]model.Question, int64, error) {
	return s.repo.FindBookmarkedByUserID(userID, limit, offset)
}

// CountByUserID は指定ユーザーの質問総数を返す。
func (s *QuestionService) CountByUserID(userID uint) (int64, error) {
	return s.repo.CountByUserID(userID)
}

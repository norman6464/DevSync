package service

import (
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// CodeSnippetService はコードスニペットのビジネスロジックを提供する。
type CodeSnippetService struct {
	repo     repository.CodeSnippetRepositoryInterface
	postRepo repository.PostRepositoryInterface
}

// NewCodeSnippetService は新しいCodeSnippetServiceインスタンスを生成する。
func NewCodeSnippetService(
	repo repository.CodeSnippetRepositoryInterface,
	postRepo repository.PostRepositoryInterface,
) *CodeSnippetService {
	return &CodeSnippetService{repo: repo, postRepo: postRepo}
}

// Create は新しいコードスニペットを作成する。
// 投稿の存在を確認してからスニペットを作成する。
func (s *CodeSnippetService) Create(snippet *model.CodeSnippet) (*model.CodeSnippet, error) {
	// 投稿の存在確認
	if _, err := s.postRepo.FindByID(snippet.PostID); err != nil {
		return nil, err
	}

	if err := s.repo.Create(snippet); err != nil {
		return nil, err
	}

	created, err := s.repo.FindByID(snippet.ID)
	if err != nil {
		return snippet, nil
	}
	return created, nil
}

// GetByPostID は指定投稿IDに紐づくスニペット一覧を返す。
func (s *CodeSnippetService) GetByPostID(postID uint) ([]model.CodeSnippet, error) {
	return s.repo.FindByPostID(postID)
}

// GetByUserLanguage は指定ユーザーのスニペットをプログラミング言語でフィルタリングして取得する。
func (s *CodeSnippetService) GetByUserLanguage(userID uint, language string) ([]model.CodeSnippet, error) {
	if language == "" {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "言語の指定は必須です", nil)
	}
	return s.repo.FindByUserIDAndLanguage(userID, language)
}

// findAndCheckOwnership はスニペットを取得し、指定ユーザーが所有者かを検証する。
func (s *CodeSnippetService) findAndCheckOwnership(id, userID uint) (*model.CodeSnippet, error) {
	snippet, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if snippet.UserID != userID {
		return nil, ErrForbidden
	}
	return snippet, nil
}

// Update はスニペットを更新する。所有者のみ更新可能。
// 空文字列のフィールドは更新しない（部分更新対応）。
func (s *CodeSnippetService) Update(id, userID uint, language, fileName, code string) (*model.CodeSnippet, error) {
	snippet, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if language != "" {
		snippet.Language = language
	}
	if fileName != "" {
		snippet.FileName = fileName
	}
	if code != "" {
		snippet.Code = code
	}

	if err := s.repo.Update(snippet); err != nil {
		return nil, err
	}
	return snippet, nil
}

// Delete はスニペットを削除する。所有者のみ削除可能。
func (s *CodeSnippetService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// CreateComment はスニペットへのインラインコメントを作成する。
// スニペットの存在を確認してからコメントを作成する。
func (s *CodeSnippetService) CreateComment(comment *model.SnippetComment) error {
	if strings.TrimSpace(comment.Content) == "" {
		return domain.NewError(domain.ErrCodeBadRequest, "コメント内容は必須です", nil)
	}
	if _, err := s.repo.FindByID(comment.SnippetID); err != nil {
		return err
	}
	return s.repo.CreateComment(comment)
}

// GetComments は指定スニペットのインラインコメント一覧を返す。
func (s *CodeSnippetService) GetComments(snippetID uint) ([]model.SnippetComment, error) {
	return s.repo.GetComments(snippetID)
}

// DeleteComment はインラインコメントを削除する。
func (s *CodeSnippetService) DeleteComment(id, userID uint) error {
	return s.repo.DeleteComment(id, userID)
}

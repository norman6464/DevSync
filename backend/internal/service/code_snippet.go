package service

import (
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

// Update はスニペットを更新する。所有者のみ更新可能。
// 空文字列のフィールドは更新しない（部分更新対応）。
func (s *CodeSnippetService) Update(id, userID uint, language, fileName, code string) (*model.CodeSnippet, error) {
	snippet, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if snippet.UserID != userID {
		return nil, ErrForbidden
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
	snippet, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if snippet.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}

// CreateComment はスニペットへのインラインコメントを作成する。
// スニペットの存在を確認してからコメントを作成する。
func (s *CodeSnippetService) CreateComment(comment *model.SnippetComment) error {
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

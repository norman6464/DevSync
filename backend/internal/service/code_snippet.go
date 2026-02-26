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
	if err := domain.ValidateStringLength(snippet.Language, 1, 100, "言語"); err != nil {
		return nil, err
	}
	if err := domain.ValidateStringLength(snippet.Code, 1, 50000, "コード"); err != nil {
		return nil, err
	}
	if err := domain.ValidateStringLength(snippet.FileName, 0, 200, "ファイル名"); err != nil {
		return nil, err
	}
	snippet.Language = strings.TrimSpace(snippet.Language)
	snippet.Code = strings.TrimSpace(snippet.Code)

	// 投稿の存在確認
	if _, err := s.postRepo.FindByID(snippet.PostID); err != nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "投稿が見つかりません", err)
	}

	if err := s.repo.Create(snippet); err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "スニペットの作成に失敗しました", err)
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

// Search はコードスニペットをキーワード検索する。
func (s *CodeSnippetService) Search(query string, limit, offset int) ([]model.CodeSnippet, int64, error) {
	q, err := validateSearchQuery(query)
	if err != nil {
		return nil, 0, err
	}
	return s.repo.Search(q, limit, offset)
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
	return checkOwnership(s.repo.FindByID, id, userID, func(cs *model.CodeSnippet) uint { return cs.UserID })
}

// Update はスニペットを更新する。所有者のみ更新可能。
// 空文字列のフィールドは更新しない（部分更新対応）。
func (s *CodeSnippetService) Update(id, userID uint, language, fileName, code string) (*model.CodeSnippet, error) {
	snippet, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(language) != "" {
		if err := domain.ValidateStringLength(language, 1, 100, "言語"); err != nil {
			return nil, err
		}
		snippet.Language = strings.TrimSpace(language)
	}
	if strings.TrimSpace(fileName) != "" {
		if err := domain.ValidateStringLength(fileName, 1, 200, "ファイル名"); err != nil {
			return nil, err
		}
		snippet.FileName = strings.TrimSpace(fileName)
	}
	if strings.TrimSpace(code) != "" {
		if err := domain.ValidateStringLength(code, 1, 50000, "コード"); err != nil {
			return nil, err
		}
		snippet.Code = strings.TrimSpace(code)
	}

	if err := s.repo.Update(snippet); err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "スニペットの更新に失敗しました", err)
	}
	return snippet, nil
}

// Delete はスニペットを削除する。所有者のみ削除可能。
func (s *CodeSnippetService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		return domain.NewError(domain.ErrCodeInternal, "スニペットの削除に失敗しました", err)
	}
	return nil
}

// CreateComment はスニペットへのインラインコメントを作成する。
// スニペットの存在を確認してからコメントを作成する。
func (s *CodeSnippetService) CreateComment(comment *model.SnippetComment) error {
	if err := domain.ValidateStringLength(comment.Content, 1, 2000, "コメント内容"); err != nil {
		return err
	}
	comment.Content = strings.TrimSpace(comment.Content)
	if _, err := s.repo.FindByID(comment.SnippetID); err != nil {
		return domain.NewError(domain.ErrCodeNotFound, "スニペットが見つかりません", err)
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

// Fork は既存スニペットをコピーして指定投稿に新しいスニペットとして作成する。
func (s *CodeSnippetService) Fork(userID, snippetID, targetPostID uint) (*model.CodeSnippet, error) {
	original, err := s.repo.FindByID(snippetID)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "スニペットが見つかりません", err)
	}

	// 対象投稿の存在確認と所有権チェック
	post, err := s.postRepo.FindByID(targetPostID)
	if err != nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "投稿が見つかりません", err)
	}
	if post.UserID != userID {
		return nil, domain.NewError(domain.ErrCodeForbidden, "自分の投稿にのみフォークできます。投稿の編集権限がありません", nil)
	}

	forked := &model.CodeSnippet{
		PostID:       targetPostID,
		UserID:       userID,
		Language:     original.Language,
		FileName:     original.FileName,
		Code:         original.Code,
		ForkedFromID: &snippetID,
	}

	if err := s.repo.Create(forked); err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "スニペットのフォークに失敗しました", err)
	}

	// フォーク元のカウンターをインクリメント
	_ = s.repo.IncrementForkCount(snippetID)

	return forked, nil
}

// Favorite はスニペットをお気に入りに追加する。
// 既にお気に入り済みの場合はErrConflictを返す。
func (s *CodeSnippetService) Favorite(userID, snippetID uint) error {
	if _, err := s.repo.FindByID(snippetID); err != nil {
		return ErrNotFound
	}
	has, err := s.repo.HasFavorited(userID, snippetID)
	if err != nil {
		return err
	}
	if has {
		return ErrConflict
	}
	return s.repo.Favorite(userID, snippetID)
}

// Unfavorite はスニペットのお気に入りを解除する。
func (s *CodeSnippetService) Unfavorite(userID, snippetID uint) error {
	return s.repo.Unfavorite(userID, snippetID)
}

// GetFavoritedByUserID はお気に入りスニペット一覧を取得する。
func (s *CodeSnippetService) GetFavoritedByUserID(userID uint, limit, offset int) ([]model.CodeSnippet, int64, error) {
	return s.repo.FindFavoritedByUserID(userID, limit, offset)
}

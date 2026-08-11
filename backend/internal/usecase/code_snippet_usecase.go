package usecase

import (
	"context"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// codeSnippetOwnerOf はスニペットの所有者 ID を返す。
func codeSnippetOwnerOf(cs *model.CodeSnippet) uint { return cs.UserID }

// validateSearchQuery は検索キーワードを検証して前後空白を除いた値を返す。
func validateSearchQuery(query string) (string, error) {
	q := strings.TrimSpace(query)
	if q == "" {
		return "", domain.NewError(domain.ErrCodeBadRequest, "検索キーワードは必須です", nil)
	}
	return q, nil
}

// requireSnippet はスニペットを取得する。不在なら 404 を返す。
func requireSnippet(ctx context.Context, snippets repository.CodeSnippetRepository, id uint) (*model.CodeSnippet, error) {
	snippet, err := snippets.FindByID(ctx, id)
	if err != nil || snippet == nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "スニペットが見つかりません", err)
	}
	return snippet, nil
}

// CreateCodeSnippetUseCase はコードスニペットを作成する。
type CreateCodeSnippetUseCase struct {
	snippets repository.CodeSnippetRepository
	posts    repository.PostReader
}

// NewCreateCodeSnippetUseCase は CreateCodeSnippetUseCase を生成する。
func NewCreateCodeSnippetUseCase(snippets repository.CodeSnippetRepository, posts repository.PostReader) *CreateCodeSnippetUseCase {
	return &CreateCodeSnippetUseCase{snippets: snippets, posts: posts}
}

// Execute は入力を検証し、投稿の存在を確認したうえでスニペットを作成する。
func (uc *CreateCodeSnippetUseCase) Execute(ctx context.Context, snippet *model.CodeSnippet) (*model.CodeSnippet, error) {
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
	snippet.FileName = strings.TrimSpace(snippet.FileName)

	post, err := uc.posts.FindByID(ctx, snippet.PostID)
	if err != nil || post == nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "投稿が見つかりません", err)
	}

	if err := uc.snippets.Create(ctx, snippet); err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "スニペットの作成に失敗しました", err)
	}

	// 作成後の再取得に失敗しても、作成した値をそのまま返す（移行前の挙動を維持している）。
	created, err := uc.snippets.FindByID(ctx, snippet.ID)
	if err != nil || created == nil {
		return snippet, nil
	}
	return created, nil
}

// ListCodeSnippetsByPostUseCase は投稿に紐づくスニペット一覧を取得する。
type ListCodeSnippetsByPostUseCase struct {
	snippets repository.CodeSnippetRepository
}

// NewListCodeSnippetsByPostUseCase は ListCodeSnippetsByPostUseCase を生成する。
func NewListCodeSnippetsByPostUseCase(snippets repository.CodeSnippetRepository) *ListCodeSnippetsByPostUseCase {
	return &ListCodeSnippetsByPostUseCase{snippets: snippets}
}

// Execute はスニペット一覧を返す。
func (uc *ListCodeSnippetsByPostUseCase) Execute(ctx context.Context, postID uint) ([]model.CodeSnippet, error) {
	return uc.snippets.FindByPostID(ctx, postID)
}

// SearchCodeSnippetsUseCase はスニペットをキーワード検索する。
type SearchCodeSnippetsUseCase struct {
	snippets repository.CodeSnippetRepository
}

// NewSearchCodeSnippetsUseCase は SearchCodeSnippetsUseCase を生成する。
func NewSearchCodeSnippetsUseCase(snippets repository.CodeSnippetRepository) *SearchCodeSnippetsUseCase {
	return &SearchCodeSnippetsUseCase{snippets: snippets}
}

// Execute は検索結果と総件数を返す。
func (uc *SearchCodeSnippetsUseCase) Execute(ctx context.Context, query string, limit, offset int) ([]model.CodeSnippet, int64, error) {
	q, err := validateSearchQuery(query)
	if err != nil {
		return nil, 0, err
	}
	return uc.snippets.Search(ctx, q, limit, offset)
}

// ListCodeSnippetsByLanguageUseCase は指定ユーザーのスニペットを言語で絞り込んで取得する。
type ListCodeSnippetsByLanguageUseCase struct {
	snippets repository.CodeSnippetRepository
}

// NewListCodeSnippetsByLanguageUseCase は ListCodeSnippetsByLanguageUseCase を生成する。
func NewListCodeSnippetsByLanguageUseCase(snippets repository.CodeSnippetRepository) *ListCodeSnippetsByLanguageUseCase {
	return &ListCodeSnippetsByLanguageUseCase{snippets: snippets}
}

// Execute は言語で絞り込んだスニペット一覧を返す。
func (uc *ListCodeSnippetsByLanguageUseCase) Execute(ctx context.Context, userID uint, language string) ([]model.CodeSnippet, error) {
	if language == "" {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "言語の指定は必須です", nil)
	}
	return uc.snippets.FindByUserIDAndLanguage(ctx, userID, language)
}

// UpdateCodeSnippetInput はスニペット更新の入力。
// 前後の空白を除いて空になる項目は据え置く部分更新。
type UpdateCodeSnippetInput struct {
	ID       uint
	UserID   uint
	Language string
	FileName string
	Code     string
}

// UpdateCodeSnippetUseCase は所有者のスニペットを更新する。
type UpdateCodeSnippetUseCase struct {
	snippets repository.CodeSnippetRepository
}

// NewUpdateCodeSnippetUseCase は UpdateCodeSnippetUseCase を生成する。
func NewUpdateCodeSnippetUseCase(snippets repository.CodeSnippetRepository) *UpdateCodeSnippetUseCase {
	return &UpdateCodeSnippetUseCase{snippets: snippets}
}

// Execute は所有権を検証したうえでスニペットを部分更新する。
func (uc *UpdateCodeSnippetUseCase) Execute(ctx context.Context, in UpdateCodeSnippetInput) (*model.CodeSnippet, error) {
	snippet, err := ensureOwner(ctx, uc.snippets.FindByID, in.ID, in.UserID, codeSnippetOwnerOf)
	if err != nil {
		return nil, err
	}

	if lang := strings.TrimSpace(in.Language); lang != "" {
		if err := domain.ValidateStringLength(lang, 1, 100, "言語"); err != nil {
			return nil, err
		}
		snippet.Language = lang
	}
	if fn := strings.TrimSpace(in.FileName); fn != "" {
		if err := domain.ValidateStringLength(fn, 1, 200, "ファイル名"); err != nil {
			return nil, err
		}
		snippet.FileName = fn
	}
	if c := strings.TrimSpace(in.Code); c != "" {
		if err := domain.ValidateStringLength(c, 1, 50000, "コード"); err != nil {
			return nil, err
		}
		snippet.Code = c
	}

	if err := uc.snippets.Update(ctx, snippet); err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "スニペットの更新に失敗しました", err)
	}
	return snippet, nil
}

// DeleteCodeSnippetUseCase は所有者のスニペットを削除する。
type DeleteCodeSnippetUseCase struct {
	snippets repository.CodeSnippetRepository
}

// NewDeleteCodeSnippetUseCase は DeleteCodeSnippetUseCase を生成する。
func NewDeleteCodeSnippetUseCase(snippets repository.CodeSnippetRepository) *DeleteCodeSnippetUseCase {
	return &DeleteCodeSnippetUseCase{snippets: snippets}
}

// Execute は所有権を検証したうえでスニペットを削除する。
func (uc *DeleteCodeSnippetUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensureOwner(ctx, uc.snippets.FindByID, id, userID, codeSnippetOwnerOf); err != nil {
		return err
	}
	if err := uc.snippets.Delete(ctx, id); err != nil {
		return domain.NewError(domain.ErrCodeInternal, "スニペットの削除に失敗しました", err)
	}
	return nil
}

// CreateSnippetCommentUseCase はスニペットへのインラインコメントを作成する。
type CreateSnippetCommentUseCase struct {
	snippets repository.CodeSnippetRepository
}

// NewCreateSnippetCommentUseCase は CreateSnippetCommentUseCase を生成する。
func NewCreateSnippetCommentUseCase(snippets repository.CodeSnippetRepository) *CreateSnippetCommentUseCase {
	return &CreateSnippetCommentUseCase{snippets: snippets}
}

// Execute は内容を検証し、スニペットの存在を確認したうえでコメントを作成する。
func (uc *CreateSnippetCommentUseCase) Execute(ctx context.Context, comment *model.SnippetComment) error {
	if err := domain.ValidateStringLength(comment.Content, 1, 2000, "コメント内容"); err != nil {
		return err
	}
	comment.Content = strings.TrimSpace(comment.Content)

	if _, err := requireSnippet(ctx, uc.snippets, comment.SnippetID); err != nil {
		return err
	}
	return uc.snippets.CreateComment(ctx, comment)
}

// ListSnippetCommentsUseCase はスニペットのインラインコメント一覧を取得する。
type ListSnippetCommentsUseCase struct {
	snippets repository.CodeSnippetRepository
}

// NewListSnippetCommentsUseCase は ListSnippetCommentsUseCase を生成する。
func NewListSnippetCommentsUseCase(snippets repository.CodeSnippetRepository) *ListSnippetCommentsUseCase {
	return &ListSnippetCommentsUseCase{snippets: snippets}
}

// Execute はコメント一覧を返す。
func (uc *ListSnippetCommentsUseCase) Execute(ctx context.Context, snippetID uint) ([]model.SnippetComment, error) {
	return uc.snippets.GetComments(ctx, snippetID)
}

// DeleteSnippetCommentUseCase はインラインコメントを削除する。
type DeleteSnippetCommentUseCase struct {
	snippets repository.CodeSnippetRepository
}

// NewDeleteSnippetCommentUseCase は DeleteSnippetCommentUseCase を生成する。
func NewDeleteSnippetCommentUseCase(snippets repository.CodeSnippetRepository) *DeleteSnippetCommentUseCase {
	return &DeleteSnippetCommentUseCase{snippets: snippets}
}

// Execute はコメントを削除する。所有権の判定は adapter 側で行う（移行前の挙動を維持している）。
func (uc *DeleteSnippetCommentUseCase) Execute(ctx context.Context, id, userID uint) error {
	return uc.snippets.DeleteComment(ctx, id, userID)
}

// ForkCodeSnippetUseCase は既存スニペットを自分の投稿へコピーする。
type ForkCodeSnippetUseCase struct {
	snippets repository.CodeSnippetRepository
	posts    repository.PostReader
}

// NewForkCodeSnippetUseCase は ForkCodeSnippetUseCase を生成する。
func NewForkCodeSnippetUseCase(snippets repository.CodeSnippetRepository, posts repository.PostReader) *ForkCodeSnippetUseCase {
	return &ForkCodeSnippetUseCase{snippets: snippets, posts: posts}
}

// Execute はフォーク元と対象投稿を検証したうえでスニペットを複製する。
func (uc *ForkCodeSnippetUseCase) Execute(ctx context.Context, userID, snippetID, targetPostID uint) (*model.CodeSnippet, error) {
	original, err := requireSnippet(ctx, uc.snippets, snippetID)
	if err != nil {
		return nil, err
	}

	post, findErr := uc.posts.FindByID(ctx, targetPostID)
	if findErr != nil || post == nil {
		return nil, domain.NewError(domain.ErrCodeNotFound, "投稿が見つかりません", findErr)
	}
	// 所有権の判定は専用メッセージを返すため、共通 helper の ensureOwner は使わない。
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
	if err := uc.snippets.Create(ctx, forked); err != nil {
		return nil, domain.NewError(domain.ErrCodeInternal, "スニペットのフォークに失敗しました", err)
	}

	// カウンター加算の失敗はフォーク自体を失敗させない（移行前の挙動を維持している）。
	_ = uc.snippets.IncrementForkCount(ctx, snippetID)

	return forked, nil
}

// FavoriteCodeSnippetUseCase はスニペットをお気に入りに追加する。
type FavoriteCodeSnippetUseCase struct {
	snippets repository.CodeSnippetRepository
}

// NewFavoriteCodeSnippetUseCase は FavoriteCodeSnippetUseCase を生成する。
func NewFavoriteCodeSnippetUseCase(snippets repository.CodeSnippetRepository) *FavoriteCodeSnippetUseCase {
	return &FavoriteCodeSnippetUseCase{snippets: snippets}
}

// Execute は存在確認と重複確認をしたうえでお気に入りに追加する。
func (uc *FavoriteCodeSnippetUseCase) Execute(ctx context.Context, userID, snippetID uint) error {
	snippet, err := uc.snippets.FindByID(ctx, snippetID)
	if err != nil || snippet == nil {
		// 移行前は取得エラーを一律 ErrNotFound に変換していたため、それを維持する。
		return domain.ErrNotFound
	}

	has, err := uc.snippets.HasFavorited(ctx, userID, snippetID)
	if err != nil {
		return err
	}
	if has {
		return domain.ErrConflict
	}
	return uc.snippets.Favorite(ctx, userID, snippetID)
}

// UnfavoriteCodeSnippetUseCase はスニペットのお気に入りを解除する。
type UnfavoriteCodeSnippetUseCase struct {
	snippets repository.CodeSnippetRepository
}

// NewUnfavoriteCodeSnippetUseCase は UnfavoriteCodeSnippetUseCase を生成する。
func NewUnfavoriteCodeSnippetUseCase(snippets repository.CodeSnippetRepository) *UnfavoriteCodeSnippetUseCase {
	return &UnfavoriteCodeSnippetUseCase{snippets: snippets}
}

// Execute はお気に入りを解除する。存在確認はしない（移行前の挙動を維持している）。
func (uc *UnfavoriteCodeSnippetUseCase) Execute(ctx context.Context, userID, snippetID uint) error {
	return uc.snippets.Unfavorite(ctx, userID, snippetID)
}

// ListFavoritedCodeSnippetsUseCase はお気に入りスニペット一覧を取得する。
type ListFavoritedCodeSnippetsUseCase struct {
	snippets repository.CodeSnippetRepository
}

// NewListFavoritedCodeSnippetsUseCase は ListFavoritedCodeSnippetsUseCase を生成する。
func NewListFavoritedCodeSnippetsUseCase(snippets repository.CodeSnippetRepository) *ListFavoritedCodeSnippetsUseCase {
	return &ListFavoritedCodeSnippetsUseCase{snippets: snippets}
}

// Execute はお気に入り一覧と総件数を返す。
func (uc *ListFavoritedCodeSnippetsUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.CodeSnippet, int64, error) {
	return uc.snippets.FindFavoritedByUserID(ctx, userID, limit, offset)
}

// CountCodeSnippetsUseCase は指定ユーザーのスニペット総数を返す。
type CountCodeSnippetsUseCase struct {
	snippets repository.CodeSnippetRepository
}

// NewCountCodeSnippetsUseCase は CountCodeSnippetsUseCase を生成する。
func NewCountCodeSnippetsUseCase(snippets repository.CodeSnippetRepository) *CountCodeSnippetsUseCase {
	return &CountCodeSnippetsUseCase{snippets: snippets}
}

// Execute はスニペット総数を返す。
func (uc *CountCodeSnippetsUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.snippets.CountByUserID(ctx, userID)
}

package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// codeSnippetRepository は [repository.CodeSnippetRepository] の sqlc(pgx) 実装。
type codeSnippetRepository struct {
	q *sqlcgen.Queries
}

// NewCodeSnippetRepository は CodeSnippetRepository の sqlc(pgx) 実装を返す。
func NewCodeSnippetRepository(q *sqlcgen.Queries) repository.CodeSnippetRepository {
	return &codeSnippetRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.CodeSnippetRepository = (*codeSnippetRepository)(nil)

// Create は新しいコードスニペットを作成する。
func (r *codeSnippetRepository) Create(ctx context.Context, snippet *model.CodeSnippet) error {
	row, err := r.q.CreateCodeSnippet(ctx, sqlcgen.CreateCodeSnippetParams{
		PostID:       int64(snippet.PostID),
		UserID:       int64(snippet.UserID),
		Language:     snippet.Language,
		FileName:     &snippet.FileName,
		Code:         snippet.Code,
		CommentCount: toInt64Ptr(snippet.CommentCount),
		ForkedFromID: toInt64PtrFromUintPtr(snippet.ForkedFromID),
		ForkCount:    toInt64Ptr(snippet.ForkCount),
	})
	if err != nil {
		return err
	}
	*snippet = toModelCodeSnippet(row)
	return nil
}

// FindByID は指定 ID のスニペットを取得する。不在の場合は (nil, nil) を返す。
func (r *codeSnippetRepository) FindByID(ctx context.Context, id uint) (*model.CodeSnippet, error) {
	row, err := r.q.GetCodeSnippetByID(ctx, int64(id))
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	snippet := toModelCodeSnippet(row)
	return &snippet, nil
}

// FindByPostID は指定投稿のスニペット一覧を取得する。
func (r *codeSnippetRepository) FindByPostID(ctx context.Context, postID uint) ([]model.CodeSnippet, error) {
	rows, err := r.q.ListCodeSnippetsByPostIDs(ctx, []int64{int64(postID)})
	if err != nil {
		return nil, err
	}
	snippets := make([]model.CodeSnippet, len(rows))
	for i, row := range rows {
		snippets[i] = toModelCodeSnippet(row)
	}
	return snippets, nil
}

// FindByUserIDAndLanguage は指定ユーザーのスニペットを言語で絞り込んで取得する。
func (r *codeSnippetRepository) FindByUserIDAndLanguage(ctx context.Context, userID uint, language string) ([]model.CodeSnippet, error) {
	rows, err := r.q.FindCodeSnippetsByUserAndLanguage(ctx, sqlcgen.FindCodeSnippetsByUserAndLanguageParams{
		UserID:   int64(userID),
		Language: language,
	})
	if err != nil {
		return nil, err
	}
	snippets := make([]model.CodeSnippet, len(rows))
	for i, row := range rows {
		snippets[i] = toModelCodeSnippet(row)
	}
	return snippets, nil
}

// Search は言語・ファイル名・コード内容からスニペットをキーワード検索する。
// CodeSnippet は投稿者の ID しか持たず User の関連を張っていないため、Preload しない
// （移行前のGORM実装のコメントの通り、Preloadするとunsupported relationsで失敗する）。
func (r *codeSnippetRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.CodeSnippet, int64, error) {
	like := escapeLikePattern(query)

	total, err := r.q.CountSearchCodeSnippets(ctx, like)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.SearchCodeSnippets(ctx, sqlcgen.SearchCodeSnippetsParams{
		Language: like,
		Limit:    int32Param(limit),
		Offset:   int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	snippets := make([]model.CodeSnippet, len(rows))
	for i, row := range rows {
		snippets[i] = toModelCodeSnippet(row)
	}
	return snippets, total, nil
}

// Update は既存のスニペットを更新する（GORMのSave＝全カラム上書きに相当）。
func (r *codeSnippetRepository) Update(ctx context.Context, snippet *model.CodeSnippet) error {
	row, err := r.q.UpdateCodeSnippet(ctx, sqlcgen.UpdateCodeSnippetParams{
		ID:           int64(snippet.ID),
		Language:     snippet.Language,
		FileName:     &snippet.FileName,
		Code:         snippet.Code,
		CommentCount: toInt64Ptr(snippet.CommentCount),
		ForkedFromID: toInt64PtrFromUintPtr(snippet.ForkedFromID),
		ForkCount:    toInt64Ptr(snippet.ForkCount),
	})
	if err != nil {
		return err
	}
	*snippet = toModelCodeSnippet(row)
	return nil
}

// Delete はスニペットを削除する。
func (r *codeSnippetRepository) Delete(ctx context.Context, id uint) error {
	return r.q.DeleteCodeSnippetByID(ctx, int64(id))
}

func toModelSnippetComment(row sqlcgen.SnippetComment) model.SnippetComment {
	return model.SnippetComment{
		ID:         uint(row.ID),
		SnippetID:  uint(row.SnippetID),
		UserID:     uint(row.UserID),
		LineNumber: int(row.LineNumber),
		Content:    row.Content,
		CreatedAt:  timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:  timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// CreateComment はインラインコメントを作成し、スニペットのコメント数を加算する。
// 移行前の GORM 実装と同じくトランザクションでは括らない（元実装も2つの独立した操作だったため）。
func (r *codeSnippetRepository) CreateComment(ctx context.Context, comment *model.SnippetComment) error {
	row, err := r.q.CreateSnippetComment(ctx, sqlcgen.CreateSnippetCommentParams{
		SnippetID:  int64(comment.SnippetID),
		UserID:     int64(comment.UserID),
		LineNumber: int64(comment.LineNumber),
		Content:    comment.Content,
	})
	if err != nil {
		return err
	}
	*comment = toModelSnippetComment(row)
	return r.q.IncrementSnippetCommentCount(ctx, int64(comment.SnippetID))
}

// GetComments は指定スニペットのインラインコメントを行番号順で取得する。
func (r *codeSnippetRepository) GetComments(ctx context.Context, snippetID uint) ([]model.SnippetComment, error) {
	rows, err := r.q.GetSnippetCommentsWithUser(ctx, int64(snippetID))
	if err != nil {
		return nil, err
	}
	comments := make([]model.SnippetComment, len(rows))
	for i, row := range rows {
		comments[i] = toModelSnippetComment(row.SnippetComment)
		comments[i].User = toModelUser(row.User)
	}
	return comments, nil
}

// FindCommentByID は指定 ID のインラインコメントを取得する。不在の場合は (nil, nil) を返す。
func (r *codeSnippetRepository) FindCommentByID(ctx context.Context, id uint) (*model.SnippetComment, error) {
	row, err := r.q.GetSnippetCommentByID(ctx, int64(id))
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	comment := toModelSnippetComment(row)
	return &comment, nil
}

// DeleteComment はインラインコメントを削除し、スニペットのコメント数を減算する。
// 所有権の判定は usecase 側で済んでいる前提。既に無ければ何もしない（冪等）。
// 移行前の GORM 実装と同じくトランザクションでは括らない（元実装も2つの独立した操作だったため）。
func (r *codeSnippetRepository) DeleteComment(ctx context.Context, id uint) error {
	row, err := r.q.GetSnippetCommentByID(ctx, int64(id))
	if isNoRows(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if err := r.q.DecrementSnippetCommentCountFloored(ctx, row.SnippetID); err != nil {
		return err
	}
	return r.q.DeleteSnippetCommentByID(ctx, int64(id))
}

// IncrementForkCount はスニペットのフォーク数を加算する。
func (r *codeSnippetRepository) IncrementForkCount(ctx context.Context, id uint) error {
	return r.q.IncrementSnippetForkCount(ctx, int64(id))
}

// Favorite はスニペットをお気に入りに追加する。
func (r *codeSnippetRepository) Favorite(ctx context.Context, userID, snippetID uint) error {
	return r.q.CreateSnippetFavorite(ctx, sqlcgen.CreateSnippetFavoriteParams{
		UserID:    int64(userID),
		SnippetID: int64(snippetID),
	})
}

// Unfavorite はスニペットのお気に入りを解除する。
func (r *codeSnippetRepository) Unfavorite(ctx context.Context, userID, snippetID uint) error {
	return r.q.DeleteSnippetFavorite(ctx, sqlcgen.DeleteSnippetFavoriteParams{
		UserID:    int64(userID),
		SnippetID: int64(snippetID),
	})
}

// HasFavorited は指定ユーザーが指定スニペットをお気に入りしているかを返す。
func (r *codeSnippetRepository) HasFavorited(ctx context.Context, userID, snippetID uint) (bool, error) {
	count, err := r.q.CountSnippetFavorite(ctx, sqlcgen.CountSnippetFavoriteParams{
		UserID:    int64(userID),
		SnippetID: int64(snippetID),
	})
	return count > 0, err
}

// FindFavoritedByUserID はお気に入りスニペットをページネーション付きで取得する。
func (r *codeSnippetRepository) FindFavoritedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.CodeSnippet, int64, error) {
	total, err := r.q.CountFavoritedCodeSnippets(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListFavoritedCodeSnippets(ctx, sqlcgen.ListFavoritedCodeSnippetsParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	snippets := make([]model.CodeSnippet, len(rows))
	for i, row := range rows {
		snippets[i] = toModelCodeSnippet(row)
	}
	return snippets, total, nil
}

// CountByUserID は指定ユーザーのスニペット総数を返す。
func (r *codeSnippetRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountCodeSnippetsByUser(ctx, int64(userID))
}

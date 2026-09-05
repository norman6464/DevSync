package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postReader は [repository.PostReader] の sqlc(pgx) 実装。
// 所有権チェックに必要な投稿読み取りだけを提供する。
type postReader struct {
	q *sqlcgen.Queries
}

// NewPostReader は PostReader の sqlc(pgx) 実装を返す。
func NewPostReader(q *sqlcgen.Queries) repository.PostReader {
	return &postReader{q: q}
}

var _ repository.PostReader = (*postReader)(nil)

// FindByID は ID で投稿を取得する（既存 PostRepository.FindByID と同じ preload）。
//
// 呼び出し側（PinPostUseCase・CreateCodeSnippetUseCase・ForkCodeSnippetUseCase・
// ensureOwner 経由の SetPostTagsUseCase）が全て「不在は非 nil の error」という
// 移行前の GORM 実装の挙動に依存しているため、CommentReader と同じく (nil, nil) には
// 変換せず pgx.ErrNoRows をそのまま返す。
func (r *postReader) FindByID(ctx context.Context, id uint) (*model.Post, error) {
	row, err := r.q.GetPostWithUserByID(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	post := toModelPost(row.Post)
	post.User = toModelUser(row.User)

	snippetRows, err := r.q.ListCodeSnippetsByPostIDs(ctx, []int64{row.Post.ID})
	if err != nil {
		return nil, err
	}
	post.CodeSnippets = make([]model.CodeSnippet, len(snippetRows))
	for i, snippetRow := range snippetRows {
		post.CodeSnippets[i] = toModelCodeSnippet(snippetRow)
	}

	return &post, nil
}

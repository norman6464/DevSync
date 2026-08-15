package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// CodeSnippetRepository はコードスニペットの永続化に対する、usecase 側が要求する契約。
type CodeSnippetRepository interface {
	Create(ctx context.Context, snippet *model.CodeSnippet) error
	// FindByID は指定 ID のスニペットを返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindByID(ctx context.Context, id uint) (*model.CodeSnippet, error)
	FindByPostID(ctx context.Context, postID uint) ([]model.CodeSnippet, error)
	FindByUserIDAndLanguage(ctx context.Context, userID uint, language string) ([]model.CodeSnippet, error)
	Search(ctx context.Context, query string, limit, offset int) ([]model.CodeSnippet, int64, error)
	Update(ctx context.Context, snippet *model.CodeSnippet) error
	Delete(ctx context.Context, id uint) error

	CreateComment(ctx context.Context, comment *model.SnippetComment) error
	GetComments(ctx context.Context, snippetID uint) ([]model.SnippetComment, error)
	// FindCommentByID は指定 ID のインラインコメントを返す。不在の場合は (nil, nil) を返す。
	FindCommentByID(ctx context.Context, id uint) (*model.SnippetComment, error)
	// DeleteComment はインラインコメントを削除し、スニペットのコメント数を減算する。
	// 所有権の判定は usecase 側で行う。既に無ければ何もしない（冪等）。
	DeleteComment(ctx context.Context, id uint) error

	IncrementForkCount(ctx context.Context, id uint) error

	Favorite(ctx context.Context, userID, snippetID uint) error
	Unfavorite(ctx context.Context, userID, snippetID uint) error
	HasFavorited(ctx context.Context, userID, snippetID uint) (bool, error)
	FindFavoritedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.CodeSnippet, int64, error)

	CountByUserID(ctx context.Context, userID uint) (int64, error)
}

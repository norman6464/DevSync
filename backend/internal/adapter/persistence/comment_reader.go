package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// commentReader は [repository.CommentReader] の sqlc(pgx) 実装。
// コメントの所有権・存在チェックに必要な読み取りだけを提供する。
type commentReader struct {
	q *sqlcgen.Queries
}

// NewCommentReader は CommentReader の sqlc(pgx) 実装を返す。
func NewCommentReader(q *sqlcgen.Queries) repository.CommentReader {
	return &commentReader{q: q}
}

var _ repository.CommentReader = (*commentReader)(nil)

// toModelComment は sqlc の生成行を model.Comment へ変換する（関連の User は含まない）。
func toModelComment(row sqlcgen.Comment) model.Comment {
	return model.Comment{
		ID:        uint(row.ID),
		UserID:    uint(row.UserID),
		PostID:    uint(row.PostID),
		ParentID:  fromInt64PtrToUintPtr(row.ParentID),
		Content:   row.Content,
		LikeCount: int(fromInt64PtrValue(row.LikeCount)),
		IsHidden:  fromBoolPtr(row.IsHidden),
		CreatedAt: timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt: timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// FindCommentByID は ID でコメントを取得する。
//
// 呼び出し側（ensureNotSelfComment 等）が「不在は非 nil の error」という
// 移行前の GORM 実装の挙動にそのまま依存しているため、ここでは (nil, nil) には
// 変換せず pgx.ErrNoRows をそのまま返す（呼び出し側を壊さないための意図的な非対称）。
func (r *commentReader) FindCommentByID(ctx context.Context, id uint) (*model.Comment, error) {
	row, err := r.q.GetCommentByID(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	comment := toModelComment(row)
	return &comment, nil
}

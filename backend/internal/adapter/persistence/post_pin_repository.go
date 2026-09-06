package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postPinRepository は [repository.PostPinRepository] の sqlc(pgx) 実装。
// UpdateOrder は複数件の更新を1トランザクションで行うため、
// Queries だけでなくトランザクションを開始できる *pgxpool.Pool を直接保持する。
type postPinRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewPostPinRepository は PostPinRepository の sqlc(pgx) 実装を返す。
func NewPostPinRepository(pool *pgxpool.Pool) repository.PostPinRepository {
	return &postPinRepository{pool: pool, q: sqlcgen.New(pool)}
}

var _ repository.PostPinRepository = (*postPinRepository)(nil)

// toModelPost は sqlc の生成行を model.Post へ変換する（Post.User の Preload 相当分のみ）。
// LikeCount/CommentCount/ViewCount/BookmarkCountはpost.postsに列を持たないため、
// 呼び出し側でattachMetricsToPosts/attachBookmarkCountsToPostsを使って別途付与する。
func toModelPost(row sqlcgen.Post) model.Post {
	return model.Post{
		ID:                uint(row.ID),
		UserID:            uint(row.UserID),
		Title:             row.Title,
		Content:           row.Content,
		ImageURLs:         fromStringPtr(row.ImageUrls),
		IsDraft:           row.IsDraft,
		EstimatedReadTime: int(fromInt64PtrValue(row.EstimatedReadTime)),
		ScheduledAt:       fromTimestamptz(row.ScheduledAt),
		CreatedAt:         timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:         timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

func (r *postPinRepository) Pin(ctx context.Context, pin *model.PostPin) error {
	return r.q.CreatePostPin(ctx, sqlcgen.CreatePostPinParams{
		UserID:   int64(pin.UserID),
		PostID:   int64(pin.PostID),
		PinOrder: int64(pin.PinOrder),
	})
}

func (r *postPinRepository) Unpin(ctx context.Context, userID, postID uint) error {
	return r.q.DeletePostPin(ctx, sqlcgen.DeletePostPinParams{
		UserID: int64(userID),
		PostID: int64(postID),
	})
}

func (r *postPinRepository) GetByUserID(ctx context.Context, userID uint) ([]model.PostPin, error) {
	rows, err := r.q.ListPostPinsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	posts := make([]model.Post, len(rows))
	for i, row := range rows {
		posts[i] = toModelPost(row.Post)
		posts[i].User = toModelUser(row.User)
	}
	if err := attachBookmarkCountsToPosts(ctx, r.q, posts); err != nil {
		return nil, err
	}
	if err := attachMetricsToPosts(ctx, r.q, posts); err != nil {
		return nil, err
	}

	pins := make([]model.PostPin, len(rows))
	for i, row := range rows {
		pins[i] = model.PostPin{
			ID:        uint(row.PostPin.ID),
			UserID:    uint(row.PostPin.UserID),
			PostID:    uint(row.PostPin.PostID),
			Post:      posts[i],
			PinOrder:  int(row.PostPin.PinOrder),
			CreatedAt: timeValue(fromTimestamptz(row.PostPin.CreatedAt)),
		}
	}
	return pins, nil
}

func (r *postPinRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountPostPinsByUser(ctx, int64(userID))
}

func (r *postPinRepository) IsPinned(ctx context.Context, userID, postID uint) (bool, error) {
	count, err := r.q.CountPostPinsByUserAndPost(ctx, sqlcgen.CountPostPinsByUserAndPostParams{
		UserID: int64(userID),
		PostID: int64(postID),
	})
	return count > 0, err
}

// UpdateOrder は複数件の pin_order 更新を1トランザクションで行う。
// 途中で失敗した場合に一部だけ並び順が更新された不整合な状態を残さないため、
// GORMのTransactionと同じくすべて成功したときだけコミットする。
func (r *postPinRepository) UpdateOrder(ctx context.Context, userID uint, postIDs []uint) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	for i, postID := range postIDs {
		if err := q.UpdatePostPinOrder(ctx, sqlcgen.UpdatePostPinOrderParams{
			UserID:   int64(userID),
			PostID:   int64(postID),
			PinOrder: int64(i),
		}); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

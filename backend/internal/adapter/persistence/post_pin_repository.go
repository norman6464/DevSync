package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postPinRepository は [repository.PostPinRepository] の sqlc(pgx) 実装。
type postPinRepository struct {
	q *sqlcgen.Queries
}

// NewPostPinRepository は PostPinRepository の sqlc(pgx) 実装を返す。
func NewPostPinRepository(q *sqlcgen.Queries) repository.PostPinRepository {
	return &postPinRepository{q: q}
}

var _ repository.PostPinRepository = (*postPinRepository)(nil)

// toModelPost は sqlc の生成行を model.Post へ変換する（Post.User の Preload 相当分のみ）。
func toModelPost(row sqlcgen.Post) model.Post {
	return model.Post{
		ID:                uint(row.ID),
		UserID:            uint(row.UserID),
		Title:             row.Title,
		Content:           row.Content,
		ImageURLs:         fromStringPtr(row.ImageUrls),
		IsDraft:           fromBoolPtr(row.IsDraft),
		LikeCount:         int(fromInt64PtrValue(row.LikeCount)),
		CommentCount:      int(fromInt64PtrValue(row.CommentCount)),
		BookmarkCount:     int(fromInt64PtrValue(row.BookmarkCount)),
		ViewCount:         int(fromInt64PtrValue(row.ViewCount)),
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
	pins := make([]model.PostPin, len(rows))
	for i, row := range rows {
		post := toModelPost(row.Post)
		post.User = toModelUser(row.User)
		pins[i] = model.PostPin{
			ID:        uint(row.PostPin.ID),
			UserID:    uint(row.PostPin.UserID),
			PostID:    uint(row.PostPin.PostID),
			Post:      post,
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

func (r *postPinRepository) UpdateOrder(ctx context.Context, userID uint, postIDs []uint) error {
	for i, postID := range postIDs {
		if err := r.q.UpdatePostPinOrder(ctx, sqlcgen.UpdatePostPinOrderParams{
			UserID:   int64(userID),
			PostID:   int64(postID),
			PinOrder: int64(i),
		}); err != nil {
			return err
		}
	}
	return nil
}

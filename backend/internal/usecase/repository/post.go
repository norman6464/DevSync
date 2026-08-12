package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// PostRepository は投稿本体の永続化に対する、usecase 側が要求する契約。
type PostRepository interface {
	Create(ctx context.Context, post *model.Post) error
	// FindByID は投稿を投稿者・コードスニペット付きで返す。存在しなければ (nil, nil) を返す。
	FindByID(ctx context.Context, id uint) (*model.Post, error)
	Update(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, id uint) error

	// FindAll は公開済み投稿を新しい順にページネーションして返す。
	FindAll(ctx context.Context, page, limit int) ([]model.Post, error)
	// CountAll は公開済み投稿の総数を返す。
	CountAll(ctx context.Context) (int64, error)
	// FindByUserID は指定ユーザーの公開済み投稿を新しい順に返す（総件数も返す）。
	FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Post, int64, error)
	// FindDraftsByUserID は指定ユーザーの下書きを更新の新しい順に返す。
	FindDraftsByUserID(ctx context.Context, userID uint) ([]model.Post, error)
	// FindScheduledByUserID は指定ユーザーの公開予約済み投稿を公開予定日時順に返す。
	FindScheduledByUserID(ctx context.Context, userID uint) ([]model.Post, error)
	// Timeline はフォロー中ユーザーと自分の公開済み投稿を新しい順に返す。
	Timeline(ctx context.Context, userID uint, page, limit int) ([]model.Post, error)

	CountByUserID(ctx context.Context, userID uint) (int64, error)
	CountDraftsByUserID(ctx context.Context, userID uint) (int64, error)
	CountScheduledByUserID(ctx context.Context, userID uint) (int64, error)
}

// PostLikeRepository は投稿への「いいね」の永続化に対する、usecase 側が要求する契約。
type PostLikeRepository interface {
	Like(ctx context.Context, userID, postID uint) error
	Unlike(ctx context.Context, userID, postID uint) error
	// HasLiked は指定ユーザーが投稿にいいね済みかを返す。
	HasLiked(ctx context.Context, userID, postID uint) (bool, error)
}

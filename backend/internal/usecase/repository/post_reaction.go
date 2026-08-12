package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// PostReactionRepository は投稿へのリアクションの永続化に対する、usecase 側が要求する契約。
type PostReactionRepository interface {
	AddReaction(ctx context.Context, userID, postID uint, emoji string) error
	RemoveReaction(ctx context.Context, userID, postID uint, emoji string) error
	// GetReactionsByPostID は指定投稿のリアクションを絵文字ごとに集計して返す（件数の降順）。
	GetReactionsByPostID(ctx context.Context, postID uint) ([]model.ReactionCount, error)
	// GetUserReactions は指定ユーザーが投稿に付けた絵文字の一覧を返す。
	GetUserReactions(ctx context.Context, userID, postID uint) ([]string, error)
	// GetReactionsBatch は複数投稿のリアクション集計をまとめて返す。
	GetReactionsBatch(ctx context.Context, postIDs []uint) (map[uint][]model.ReactionCount, error)
	// GetUserReactionsBatch は複数投稿に対するユーザーのリアクションをまとめて返す。
	GetUserReactionsBatch(ctx context.Context, userID uint, postIDs []uint) (map[uint][]string, error)
}

// PostAuthorReader は投稿の投稿者を判定するための最小の契約。
// 「自分の投稿には反応できない」の判定に使う。
type PostAuthorReader interface {
	// FindAuthorID は指定投稿の投稿者 ID を返す。
	// 投稿が存在しない場合は「不在」を表す (0, nil) を返し、DB 障害だけを error として返す。
	FindAuthorID(ctx context.Context, postID uint) (uint, error)
}

// Package repository は usecase が依存する repository の契約（port）を定義する。
// 実装は adapter/persistence に置き、依存の向きを usecase ← persistence に逆転させる（DIP）。
package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// FollowRepository はフォロー関係の永続化に対する、usecase 側が要求する契約。
// 使う側（usecase）がこの interface を宣言し、実装（persistence）がそれに従う。
type FollowRepository interface {
	Follow(ctx context.Context, followerID, followeeID uint) error
	Unfollow(ctx context.Context, followerID, followeeID uint) error
	GetFollowers(ctx context.Context, userID uint, limit, offset int) ([]model.User, int64, error)
	GetFollowing(ctx context.Context, userID uint, limit, offset int) ([]model.User, int64, error)
}

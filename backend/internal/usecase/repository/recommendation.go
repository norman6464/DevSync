package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// RecommendationRepository はレコメンドのデータ取得に対する、usecase 側が要求する契約。
type RecommendationRepository interface {
	// GetRecommendedUsers は指定スキルと共通点を持つユーザーを、共通スキル数の降順で limit 件返す。
	// 自分自身とフォロー済みのユーザーは除外する。
	GetRecommendedUsers(ctx context.Context, userID uint, skills []string, limit int) ([]model.RecommendedUser, error)
	// GetTrendingPosts は直近 days 日の人気投稿を limit 件返す。
	GetTrendingPosts(ctx context.Context, limit, days int) ([]model.Post, error)
	// GetTrendingResources は直近 days 日の公開学習リソースを人気順に limit 件返す。
	GetTrendingResources(ctx context.Context, limit, days int) ([]model.LearningResource, error)
}

// UserSkillsReader はおすすめユーザーの算出でプロフィールのスキルを読むための最小の契約。
// 共有の user リポジトリ全体には依存しない。
type UserSkillsReader interface {
	// FindByID は指定 ID のユーザーを返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindByID(ctx context.Context, id uint) (*model.User, error)
}

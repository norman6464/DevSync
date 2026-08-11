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

// おすすめユーザーの算出で使うプロフィール参照は、user スライスの [UserSkillsReader] を使う。

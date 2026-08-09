package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// maxMostViewedLimit は人気投稿ランキングで取得できる最大件数。
const maxMostViewedLimit = 100

// GetMostViewedPostsUseCase は閲覧数が多い投稿のランキングを取得する。
type GetMostViewedPostsUseCase struct {
	views repository.PostViewRepository
}

// NewGetMostViewedPostsUseCase は GetMostViewedPostsUseCase を生成する。
func NewGetMostViewedPostsUseCase(views repository.PostViewRepository) *GetMostViewedPostsUseCase {
	return &GetMostViewedPostsUseCase{views: views}
}

// Execute は閲覧数上位の投稿ランキングを返す。limit は 1〜100 の範囲。
func (uc *GetMostViewedPostsUseCase) Execute(ctx context.Context, limit int) ([]model.ViewCount, error) {
	if limit <= 0 || limit > maxMostViewedLimit {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "limitは1〜100の範囲で指定してください", nil)
	}
	return uc.views.GetMostViewed(ctx, limit)
}

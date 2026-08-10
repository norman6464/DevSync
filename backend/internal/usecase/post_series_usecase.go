package usecase

import (
	"context"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// seriesOwnerOf は所有権チェック用にシリーズの所有者 ID を取り出す。
func seriesOwnerOf(s *model.PostSeries) uint { return s.UserID }

// CreatePostSeriesUseCase は投稿シリーズを作成する。
type CreatePostSeriesUseCase struct {
	series repository.PostSeriesRepository
}

// NewCreatePostSeriesUseCase は CreatePostSeriesUseCase を生成する。
func NewCreatePostSeriesUseCase(series repository.PostSeriesRepository) *CreatePostSeriesUseCase {
	return &CreatePostSeriesUseCase{series: series}
}

// Execute はタイトルと説明を検証し、前後の空白を落として作成する。
func (uc *CreatePostSeriesUseCase) Execute(ctx context.Context, series *model.PostSeries) error {
	if err := domain.ValidateStringLength(series.Title, 1, 200, "タイトル"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(series.Description, 0, 1000, "説明"); err != nil {
		return err
	}
	series.Title = strings.TrimSpace(series.Title)
	series.Description = strings.TrimSpace(series.Description)
	return uc.series.Create(ctx, series)
}

// GetPostSeriesUseCase は指定 ID のシリーズを取得する。
type GetPostSeriesUseCase struct {
	series repository.PostSeriesRepository
}

// NewGetPostSeriesUseCase は GetPostSeriesUseCase を生成する。
func NewGetPostSeriesUseCase(series repository.PostSeriesRepository) *GetPostSeriesUseCase {
	return &GetPostSeriesUseCase{series: series}
}

// Execute はシリーズを返す。
func (uc *GetPostSeriesUseCase) Execute(ctx context.Context, id uint) (*model.PostSeries, error) {
	return uc.series.FindByID(ctx, id)
}

// ListPostSeriesUseCase は指定ユーザーのシリーズ一覧を取得する。
type ListPostSeriesUseCase struct {
	series repository.PostSeriesRepository
}

// NewListPostSeriesUseCase は ListPostSeriesUseCase を生成する。
func NewListPostSeriesUseCase(series repository.PostSeriesRepository) *ListPostSeriesUseCase {
	return &ListPostSeriesUseCase{series: series}
}

// Execute はページ番号から offset を求めて一覧を返す。
func (uc *ListPostSeriesUseCase) Execute(ctx context.Context, userID uint, page, limit int) ([]model.PostSeries, error) {
	offset := (page - 1) * limit
	return uc.series.FindByUserID(ctx, userID, offset, limit)
}

// CountPostSeriesUseCase は指定ユーザーのシリーズ数を数える。
type CountPostSeriesUseCase struct {
	series repository.PostSeriesRepository
}

// NewCountPostSeriesUseCase は CountPostSeriesUseCase を生成する。
func NewCountPostSeriesUseCase(series repository.PostSeriesRepository) *CountPostSeriesUseCase {
	return &CountPostSeriesUseCase{series: series}
}

// Execute はシリーズ数を返す。
func (uc *CountPostSeriesUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.series.CountByUser(ctx, userID)
}

// UpdatePostSeriesUseCase は所有者本人のシリーズを更新する。
type UpdatePostSeriesUseCase struct {
	series repository.PostSeriesRepository
}

// NewUpdatePostSeriesUseCase は UpdatePostSeriesUseCase を生成する。
func NewUpdatePostSeriesUseCase(series repository.PostSeriesRepository) *UpdatePostSeriesUseCase {
	return &UpdatePostSeriesUseCase{series: series}
}

// Execute は所有権を検証し、空でないフィールドだけを更新する。
func (uc *UpdatePostSeriesUseCase) Execute(ctx context.Context, id, userID uint, updates *model.PostSeries) (*model.PostSeries, error) {
	series, err := ensureOwner(ctx, uc.series.FindByID, id, userID, seriesOwnerOf)
	if err != nil {
		return nil, err
	}

	if updates.Title != "" {
		if err := domain.ValidateStringLength(updates.Title, 1, 200, "タイトル"); err != nil {
			return nil, err
		}
		series.Title = strings.TrimSpace(updates.Title)
	}
	if updates.Description != "" {
		if err := domain.ValidateStringLength(updates.Description, 1, 1000, "説明"); err != nil {
			return nil, err
		}
		series.Description = strings.TrimSpace(updates.Description)
	}

	if err := uc.series.Update(ctx, series); err != nil {
		return nil, err
	}
	return series, nil
}

// DeletePostSeriesUseCase は所有者本人のシリーズを削除する。
type DeletePostSeriesUseCase struct {
	series repository.PostSeriesRepository
}

// NewDeletePostSeriesUseCase は DeletePostSeriesUseCase を生成する。
func NewDeletePostSeriesUseCase(series repository.PostSeriesRepository) *DeletePostSeriesUseCase {
	return &DeletePostSeriesUseCase{series: series}
}

// Execute は所有権を検証してから削除する。
func (uc *DeletePostSeriesUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensureOwner(ctx, uc.series.FindByID, id, userID, seriesOwnerOf); err != nil {
		return err
	}
	return uc.series.Delete(ctx, id)
}

// AddPostToSeriesUseCase はシリーズに投稿を追加する。
type AddPostToSeriesUseCase struct {
	series repository.PostSeriesRepository
}

// NewAddPostToSeriesUseCase は AddPostToSeriesUseCase を生成する。
func NewAddPostToSeriesUseCase(series repository.PostSeriesRepository) *AddPostToSeriesUseCase {
	return &AddPostToSeriesUseCase{series: series}
}

// Execute は所有権を検証し、未追加の投稿だけを追加する。
func (uc *AddPostToSeriesUseCase) Execute(ctx context.Context, seriesID, postID uint, orderIndex int, userID uint) error {
	if _, err := ensureOwner(ctx, uc.series.FindByID, seriesID, userID, seriesOwnerOf); err != nil {
		return err
	}

	exists, err := uc.series.HasPost(ctx, seriesID, postID)
	if err != nil {
		return err
	}
	if exists {
		return domain.NewError(domain.ErrCodeBadRequest, "すでに追加済みの投稿です", nil)
	}

	return uc.series.AddPost(ctx, &model.PostSeriesItem{
		SeriesID:   seriesID,
		PostID:     postID,
		OrderIndex: orderIndex,
	})
}

// RemovePostFromSeriesUseCase はシリーズから投稿を取り除く。
type RemovePostFromSeriesUseCase struct {
	series repository.PostSeriesRepository
}

// NewRemovePostFromSeriesUseCase は RemovePostFromSeriesUseCase を生成する。
func NewRemovePostFromSeriesUseCase(series repository.PostSeriesRepository) *RemovePostFromSeriesUseCase {
	return &RemovePostFromSeriesUseCase{series: series}
}

// Execute は所有権を検証してから取り除く。
func (uc *RemovePostFromSeriesUseCase) Execute(ctx context.Context, seriesID, postID, userID uint) error {
	if _, err := ensureOwner(ctx, uc.series.FindByID, seriesID, userID, seriesOwnerOf); err != nil {
		return err
	}
	return uc.series.RemovePost(ctx, seriesID, postID)
}

// ListPostSeriesPostsUseCase はシリーズに含まれる投稿一覧を取得する。
type ListPostSeriesPostsUseCase struct {
	series repository.PostSeriesRepository
}

// NewListPostSeriesPostsUseCase は ListPostSeriesPostsUseCase を生成する。
func NewListPostSeriesPostsUseCase(series repository.PostSeriesRepository) *ListPostSeriesPostsUseCase {
	return &ListPostSeriesPostsUseCase{series: series}
}

// Execute は順序付きの投稿一覧を返す。
func (uc *ListPostSeriesPostsUseCase) Execute(ctx context.Context, seriesID uint) ([]model.PostSeriesItem, error) {
	return uc.series.GetPostsBySeriesID(ctx, seriesID)
}

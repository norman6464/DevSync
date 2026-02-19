package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// PostSeriesService は投稿シリーズに関するビジネスロジックを提供する。
// 関連する投稿をシリーズとしてグループ化するCRUD操作を担当する。
type PostSeriesService struct {
	repo repository.PostSeriesRepositoryInterface
}

// NewPostSeriesService は新しいPostSeriesServiceインスタンスを生成する。
func NewPostSeriesService(repo repository.PostSeriesRepositoryInterface) *PostSeriesService {
	return &PostSeriesService{repo: repo}
}

// Create は新しい投稿シリーズを作成する。
func (s *PostSeriesService) Create(series *model.PostSeries) error {
	if series.Title == "" {
		return domain.NewError(domain.ErrCodeValidation, "タイトルは必須です", nil)
	}
	return s.repo.Create(series)
}

// GetByID は指定IDのシリーズを取得する。
func (s *PostSeriesService) GetByID(id uint) (*model.PostSeries, error) {
	return s.repo.FindByID(id)
}

// GetByUserID は指定ユーザーの全シリーズを取得する。
func (s *PostSeriesService) GetByUserID(userID uint) ([]model.PostSeries, error) {
	return s.repo.FindByUserID(userID)
}

// findAndCheckOwnership はシリーズを取得し、指定ユーザーが所有者かを検証する。
func (s *PostSeriesService) findAndCheckOwnership(id, userID uint) (*model.PostSeries, error) {
	series, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if series.UserID != userID {
		return nil, ErrForbidden
	}
	return series, nil
}

// Update は所有権を検証した後、シリーズを更新する。
func (s *PostSeriesService) Update(id, userID uint, updates *model.PostSeries) (*model.PostSeries, error) {
	series, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if updates.Title != "" {
		series.Title = updates.Title
	}
	if updates.Description != "" {
		series.Description = updates.Description
	}

	if err := s.repo.Update(series); err != nil {
		return nil, err
	}
	return series, nil
}

// Delete は所有権を検証した後、シリーズを削除する。
func (s *PostSeriesService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// AddPost は所有権を検証した後、シリーズに投稿を追加する。
func (s *PostSeriesService) AddPost(seriesID, postID uint, orderIndex int, userID uint) error {
	if _, err := s.findAndCheckOwnership(seriesID, userID); err != nil {
		return err
	}

	item := &model.PostSeriesItem{
		SeriesID:   seriesID,
		PostID:     postID,
		OrderIndex: orderIndex,
	}
	return s.repo.AddPost(item)
}

// RemovePost は所有権を検証した後、シリーズから投稿を削除する。
func (s *PostSeriesService) RemovePost(seriesID, postID, userID uint) error {
	if _, err := s.findAndCheckOwnership(seriesID, userID); err != nil {
		return err
	}
	return s.repo.RemovePost(seriesID, postID)
}

// GetPosts は指定シリーズの投稿一覧を取得する。
func (s *PostSeriesService) GetPosts(seriesID uint) ([]model.PostSeriesItem, error) {
	return s.repo.GetPostsBySeriesID(seriesID)
}

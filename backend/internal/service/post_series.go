package service

import (
	"strings"

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
	if err := domain.ValidateStringLength(series.Title, 1, 200, "タイトル"); err != nil {
		return err
	}
	if len([]rune(series.Description)) > 1000 {
		return domain.NewError(domain.ErrCodeValidation, "説明は1000文字以下である必要があります", nil)
	}
	series.Title = strings.TrimSpace(series.Title)
	return s.repo.Create(series)
}

// GetByID は指定IDのシリーズを取得する。
func (s *PostSeriesService) GetByID(id uint) (*model.PostSeries, error) {
	return s.repo.FindByID(id)
}

// GetByUserID は指定ユーザーのシリーズをページネーション付きで取得する。
func (s *PostSeriesService) GetByUserID(userID uint, page, limit int) ([]model.PostSeries, error) {
	offset := (page - 1) * limit
	return s.repo.FindByUserID(userID, offset, limit)
}

// CountByUser は指定ユーザーのシリーズ数をカウントする。
func (s *PostSeriesService) CountByUser(userID uint) (int64, error) {
	return s.repo.CountByUser(userID)
}

// findAndCheckOwnership はシリーズを取得し、指定ユーザーが所有者かを検証する。
func (s *PostSeriesService) findAndCheckOwnership(id, userID uint) (*model.PostSeries, error) {
	return checkOwnership(s.repo.FindByID, id, userID, func(ps *model.PostSeries) uint { return ps.UserID })
}

// Update は所有権を検証した後、シリーズを更新する。
func (s *PostSeriesService) Update(id, userID uint, updates *model.PostSeries) (*model.PostSeries, error) {
	series, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if updates.Title != "" {
		if strings.TrimSpace(updates.Title) == "" {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "タイトルは空白のみでは入力できません", nil)
		}
		if len(strings.TrimSpace(updates.Title)) > 200 {
			return nil, domain.NewError(domain.ErrCodeValidation, "タイトルは200文字以下である必要があります", nil)
		}
		series.Title = strings.TrimSpace(updates.Title)
	}
	if updates.Description != "" {
		if strings.TrimSpace(updates.Description) == "" {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "説明は空白のみでは入力できません", nil)
		}
		if len(strings.TrimSpace(updates.Description)) > 1000 {
			return nil, domain.NewError(domain.ErrCodeValidation, "説明は1000文字以下である必要があります", nil)
		}
		series.Description = strings.TrimSpace(updates.Description)
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
// 同じ投稿がすでに追加されている場合はエラーを返す。
func (s *PostSeriesService) AddPost(seriesID, postID uint, orderIndex int, userID uint) error {
	if _, err := s.findAndCheckOwnership(seriesID, userID); err != nil {
		return err
	}

	exists, err := s.repo.HasPost(seriesID, postID)
	if err != nil {
		return err
	}
	if exists {
		return domain.NewError(domain.ErrCodeBadRequest, "すでに追加済みの投稿です", nil)
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

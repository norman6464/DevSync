package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// PostViewService は投稿閲覧数のビジネスロジックを提供する。
type PostViewService struct {
	repo repository.PostViewRepositoryInterface
}

// NewPostViewService は新しいPostViewServiceインスタンスを生成する。
func NewPostViewService(repo repository.PostViewRepositoryInterface) *PostViewService {
	return &PostViewService{repo: repo}
}

// RecordView はユーザーの投稿閲覧を記録する。
// 既に閲覧済みの場合は何もしない（ユニーク閲覧のみカウント）。
func (s *PostViewService) RecordView(userID, postID uint) error {
	if err := validateRequiredID(userID, "userID"); err != nil {
		return err
	}
	if err := validateRequiredID(postID, "postID"); err != nil {
		return err
	}
	viewed, err := s.repo.HasViewed(userID, postID)
	if err != nil {
		return err
	}
	if viewed {
		return nil
	}
	return s.repo.RecordView(&model.PostView{
		UserID: userID,
		PostID: postID,
	})
}

// GetViewCount は投稿のユニーク閲覧数を取得する。
func (s *PostViewService) GetViewCount(postID uint) (int64, error) {
	if err := validateRequiredID(postID, "postID"); err != nil {
		return 0, err
	}
	return s.repo.GetViewCount(postID)
}

// HasViewed はユーザーが投稿を閲覧済みかどうかを判定する。
func (s *PostViewService) HasViewed(userID, postID uint) (bool, error) {
	if err := validateRequiredID(userID, "userID"); err != nil {
		return false, err
	}
	if err := validateRequiredID(postID, "postID"); err != nil {
		return false, err
	}
	return s.repo.HasViewed(userID, postID)
}

const maxMostViewedLimit = 100

// GetMostViewed は閲覧数が多い投稿のランキングを取得する。
func (s *PostViewService) GetMostViewed(limit int) ([]model.ViewCount, error) {
	if limit <= 0 || limit > maxMostViewedLimit {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "limitは1〜100の範囲で指定してください", nil)
	}
	return s.repo.GetMostViewed(limit)
}

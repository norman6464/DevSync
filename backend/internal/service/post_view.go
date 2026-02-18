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
	if userID == 0 {
		return domain.NewError(domain.ErrCodeBadRequest, "userIDは必須です", nil)
	}
	if postID == 0 {
		return domain.NewError(domain.ErrCodeBadRequest, "postIDは必須です", nil)
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
	if postID == 0 {
		return 0, domain.NewError(domain.ErrCodeBadRequest, "postIDは必須です", nil)
	}
	return s.repo.GetViewCount(postID)
}

// HasViewed はユーザーが投稿を閲覧済みかどうかを判定する。
func (s *PostViewService) HasViewed(userID, postID uint) (bool, error) {
	if userID == 0 {
		return false, domain.NewError(domain.ErrCodeBadRequest, "userIDは必須です", nil)
	}
	if postID == 0 {
		return false, domain.NewError(domain.ErrCodeBadRequest, "postIDは必須です", nil)
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

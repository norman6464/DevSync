package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

const maxPinsPerUser = 3

// PostPinService は投稿ピン留めのビジネスロジックを提供する。
type PostPinService struct {
	pinRepo  repository.PostPinRepositoryInterface
	postRepo repository.PostRepositoryInterface
}

// NewPostPinService は新しいPostPinServiceインスタンスを生成する。
func NewPostPinService(pinRepo repository.PostPinRepositoryInterface, postRepo repository.PostRepositoryInterface) *PostPinService {
	return &PostPinService{pinRepo: pinRepo, postRepo: postRepo}
}

// Pin は投稿をプロフィールにピン留めする。
// 自分の投稿のみピン留め可能、最大3件まで。
func (s *PostPinService) Pin(userID, postID uint) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return ErrNotFound
	}
	if post.UserID != userID {
		return domain.NewError(domain.ErrCodeForbidden, "自分の投稿のみピン留めできます", nil)
	}

	count, err := s.pinRepo.CountByUserID(userID)
	if err != nil {
		return err
	}
	if count >= maxPinsPerUser {
		return domain.NewError(domain.ErrCodeBadRequest, "ピン留めは最大3件までです", nil)
	}

	pin := &model.PostPin{
		UserID:   userID,
		PostID:   postID,
		PinOrder: int(count),
	}
	return s.pinRepo.Pin(pin)
}

// Unpin は投稿のピン留めを解除する。
func (s *PostPinService) Unpin(userID, postID uint) error {
	pinned, err := s.pinRepo.IsPinned(userID, postID)
	if err != nil {
		return err
	}
	if !pinned {
		return domain.NewError(domain.ErrCodeNotFound, "ピン留めされていません", nil)
	}
	return s.pinRepo.Unpin(userID, postID)
}

// GetByUserID はユーザーのピン留め投稿一覧を取得する。
func (s *PostPinService) GetByUserID(userID uint) ([]model.PostPin, error) {
	return s.pinRepo.GetByUserID(userID)
}

// Reorder はピン留め投稿の表示順序を変更する。
// 渡されたpostIDsが全てuserIDのピン留め済み投稿であることを検証する。
func (s *PostPinService) Reorder(userID uint, postIDs []uint) error {
	if len(postIDs) > maxPinsPerUser {
		return domain.NewError(domain.ErrCodeBadRequest, "ピン留めは最大3件までです", nil)
	}

	// 現在のピン留め投稿を取得して所有権を検証
	pins, err := s.pinRepo.GetByUserID(userID)
	if err != nil {
		return err
	}
	pinnedSet := make(map[uint]bool, len(pins))
	for _, pin := range pins {
		pinnedSet[pin.PostID] = true
	}
	for _, postID := range postIDs {
		if !pinnedSet[postID] {
			return domain.NewError(domain.ErrCodeForbidden, "自分のピン留め投稿のみ順序変更できます", nil)
		}
	}

	return s.pinRepo.UpdateOrder(userID, postIDs)
}

// IsPinned は投稿がピン留めされているか確認する。
func (s *PostPinService) IsPinned(userID, postID uint) (bool, error) {
	return s.pinRepo.IsPinned(userID, postID)
}

package service

import (
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// ResourceReviewService は学習リソースレビューのビジネスロジックを提供する。
type ResourceReviewService struct {
	repo         repository.ResourceReviewRepositoryInterface
	resourceRepo repository.LearningResourceRepositoryInterface
}

// NewResourceReviewService は新しいResourceReviewServiceインスタンスを生成する。
func NewResourceReviewService(
	repo repository.ResourceReviewRepositoryInterface,
	resourceRepo repository.LearningResourceRepositoryInterface,
) *ResourceReviewService {
	return &ResourceReviewService{repo: repo, resourceRepo: resourceRepo}
}

// Create は新しいレビューを作成する。
// リソースの存在確認、評価値バリデーション、重複チェックを行う。
func (s *ResourceReviewService) Create(review *model.ResourceReview) error {
	// リソースの存在確認
	if _, err := s.resourceRepo.FindByID(review.ResourceID); err != nil {
		return ErrNotFound
	}

	// 評価値バリデーション（1-5）
	if err := validateRating(review.Rating); err != nil {
		return err
	}

	// コメントバリデーション
	if err := domain.ValidateStringLength(review.Comment, 0, 5000, "コメント"); err != nil {
		return err
	}

	// 重複チェック
	existing, _ := s.repo.FindByUserAndResource(review.UserID, review.ResourceID)
	if existing != nil {
		return ErrConflict
	}

	return s.repo.Create(review)
}

// GetByResourceID は指定リソースのレビュー一覧を取得する。
func (s *ResourceReviewService) GetByResourceID(resourceID uint, limit, offset int) ([]model.ResourceReview, int64, error) {
	return s.repo.FindByResourceID(resourceID, limit, offset)
}

// findAndCheckOwnership はレビューを取得し、指定ユーザーが所有者かを検証する。
func (s *ResourceReviewService) findAndCheckOwnership(id, userID uint) (*model.ResourceReview, error) {
	return checkOwnership(s.repo.FindByID, id, userID, func(r *model.ResourceReview) uint { return r.UserID })
}

// Update はレビューを更新する。所有者のみ更新可能。
func (s *ResourceReviewService) Update(id, userID uint, rating int, comment string) (*model.ResourceReview, error) {
	review, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if rating != 0 {
		if err := validateRating(rating); err != nil {
			return nil, err
		}
		review.Rating = rating
	}
	if strings.TrimSpace(comment) != "" {
		comment = strings.TrimSpace(comment)
		if err := domain.ValidateStringLength(comment, 1, 5000, "コメント"); err != nil {
			return nil, err
		}
		review.Comment = comment
	}

	if err := s.repo.Update(review); err != nil {
		return nil, err
	}
	return review, nil
}

// Delete はレビューを削除する。所有者のみ削除可能。
func (s *ResourceReviewService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

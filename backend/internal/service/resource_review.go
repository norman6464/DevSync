package service

import (
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
	if review.Rating < 1 || review.Rating > 5 {
		return domain.NewError(domain.ErrCodeValidation, "評価は1〜5の範囲で指定してください", nil)
	}

	// コメントバリデーション
	if len(review.Comment) > 5000 {
		return domain.NewError(domain.ErrCodeValidation, "コメントは5000文字以下で入力してください", nil)
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

// Update はレビューを更新する。所有者のみ更新可能。
func (s *ResourceReviewService) Update(id, userID uint, rating int, comment string) (*model.ResourceReview, error) {
	review, err := s.repo.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	if review.UserID != userID {
		return nil, ErrForbidden
	}

	if rating != 0 {
		if rating < 1 || rating > 5 {
			return nil, domain.NewError(domain.ErrCodeValidation, "評価は1〜5の範囲で指定してください", nil)
		}
		review.Rating = rating
	}
	if comment != "" {
		if len(comment) > 5000 {
			return nil, domain.NewError(domain.ErrCodeValidation, "コメントは5000文字以下で入力してください", nil)
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
	review, err := s.repo.FindByID(id)
	if err != nil {
		return ErrNotFound
	}
	if review.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}

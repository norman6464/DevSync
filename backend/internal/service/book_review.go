package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// BookReviewService は書籍レビューのビジネスロジックを提供する。
type BookReviewService struct {
	repo repository.BookReviewRepositoryInterface
}

// NewBookReviewService は新しいBookReviewServiceインスタンスを生成する。
func NewBookReviewService(repo repository.BookReviewRepositoryInterface) *BookReviewService {
	return &BookReviewService{repo: repo}
}

// Create は新しい書籍レビューを作成する。
func (s *BookReviewService) Create(review *model.BookReview) error {
	if review.Rating < 1 || review.Rating > 5 {
		return domain.NewError(domain.ErrCodeBadRequest, "評価は1〜5の範囲で指定してください", nil)
	}
	return s.repo.Create(review)
}

// GetByID は指定IDの書籍レビューを取得する。
func (s *BookReviewService) GetByID(id uint) (*model.BookReview, error) {
	return s.repo.FindByID(id)
}

// GetByUserID は指定ユーザーの全書籍レビューを取得する。
func (s *BookReviewService) GetByUserID(userID uint) ([]model.BookReview, error) {
	return s.repo.FindByUserID(userID)
}

// GetAll は書籍レビュー一覧をページネーション付きで取得する。
func (s *BookReviewService) GetAll(limit, offset int) ([]model.BookReview, int64, error) {
	return s.repo.FindAll(limit, offset)
}

// findAndCheckOwnership は書籍レビューを取得し、指定ユーザーが所有者かを検証する。
func (s *BookReviewService) findAndCheckOwnership(id, userID uint) (*model.BookReview, error) {
	review, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if review.UserID != userID {
		return nil, ErrForbidden
	}
	return review, nil
}

// Update は所有権を検証した後、書籍レビューを更新する。
func (s *BookReviewService) Update(id, userID uint, updates *model.BookReview) (*model.BookReview, error) {
	review, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if updates.Title != "" {
		review.Title = updates.Title
	}
	if updates.Author != "" {
		review.Author = updates.Author
	}
	if updates.ISBN != "" {
		review.ISBN = updates.ISBN
	}
	if updates.Rating != 0 {
		if updates.Rating < 1 || updates.Rating > 5 {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "評価は1〜5の範囲で指定してください", nil)
		}
		review.Rating = updates.Rating
	}
	if updates.Review != "" {
		review.Review = updates.Review
	}
	if updates.ImageURL != "" {
		review.ImageURL = updates.ImageURL
	}

	if err := s.repo.Update(review); err != nil {
		return nil, err
	}
	return review, nil
}

// Delete は所有権を検証した後、書籍レビューを削除する。
func (s *BookReviewService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

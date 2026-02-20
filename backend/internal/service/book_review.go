package service

import (
	"strings"

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

// validateRating は評価値が1〜5の範囲内かを検証する。
func validateRating(rating int) error {
	if rating < 1 || rating > 5 {
		return domain.NewError(domain.ErrCodeBadRequest, "評価は1〜5の範囲で指定してください", nil)
	}
	return nil
}

// Create は新しい書籍レビューを作成する。
func (s *BookReviewService) Create(review *model.BookReview) error {
	if strings.TrimSpace(review.Title) == "" {
		return domain.NewError(domain.ErrCodeBadRequest, "タイトルは必須です", nil)
	}
	if err := validateRating(review.Rating); err != nil {
		return err
	}
	return s.repo.Create(review)
}

// GetByID は指定IDの書籍レビューを取得する。
func (s *BookReviewService) GetByID(id uint) (*model.BookReview, error) {
	return s.repo.FindByID(id)
}

// GetByUserID は指定ユーザーの書籍レビューをページネーション付きで取得する。
func (s *BookReviewService) GetByUserID(userID uint, limit, offset int) ([]model.BookReview, int64, error) {
	return s.repo.FindByUserID(userID, limit, offset)
}

// GetAll は書籍レビュー一覧をページネーション付きで取得する。
func (s *BookReviewService) GetAll(limit, offset int) ([]model.BookReview, int64, error) {
	return s.repo.FindAll(limit, offset)
}

// GetByRating は指定ユーザーの書籍レビューを評価範囲でフィルタリングして取得する。
func (s *BookReviewService) GetByRating(userID uint, minRating, maxRating int) ([]model.BookReview, error) {
	if minRating < 1 || minRating > 5 || maxRating < 1 || maxRating > 5 || minRating > maxRating {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "評価範囲が無効です", nil)
	}
	return s.repo.FindByRating(userID, minRating, maxRating)
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

	if strings.TrimSpace(updates.Title) != "" {
		review.Title = strings.TrimSpace(updates.Title)
	}
	if strings.TrimSpace(updates.Author) != "" {
		review.Author = strings.TrimSpace(updates.Author)
	}
	if strings.TrimSpace(updates.ISBN) != "" {
		review.ISBN = strings.TrimSpace(updates.ISBN)
	}
	if updates.Rating != 0 {
		if err := validateRating(updates.Rating); err != nil {
			return nil, err
		}
		review.Rating = updates.Rating
	}
	if strings.TrimSpace(updates.Review) != "" {
		review.Review = strings.TrimSpace(updates.Review)
	}
	if strings.TrimSpace(updates.ImageURL) != "" {
		review.ImageURL = strings.TrimSpace(updates.ImageURL)
	}

	if err := s.repo.Update(review); err != nil {
		return nil, err
	}
	return review, nil
}

// ArchiveReview は所有権を検証した後、書籍レビューをアーカイブする。
func (s *BookReviewService) ArchiveReview(id, userID uint) error {
	review, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return err
	}
	review.IsArchived = true
	return s.repo.Update(review)
}

// UnarchiveReview は所有権を検証した後、書籍レビューのアーカイブを解除する。
func (s *BookReviewService) UnarchiveReview(id, userID uint) error {
	review, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return err
	}
	review.IsArchived = false
	return s.repo.Update(review)
}

// Delete は所有権を検証した後、書籍レビューを削除する。
func (s *BookReviewService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

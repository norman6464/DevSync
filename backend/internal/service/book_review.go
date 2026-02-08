package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// BookReviewService handles book review business logic.
type BookReviewService struct {
	repo repository.BookReviewRepositoryInterface
}

// NewBookReviewService creates a new BookReviewService.
func NewBookReviewService(repo repository.BookReviewRepositoryInterface) *BookReviewService {
	return &BookReviewService{repo: repo}
}

// Create creates a new book review.
func (s *BookReviewService) Create(review *model.BookReview) error {
	return s.repo.Create(review)
}

// GetByID returns a book review by ID.
func (s *BookReviewService) GetByID(id uint) (*model.BookReview, error) {
	return s.repo.FindByID(id)
}

// GetByUserID returns all book reviews for a user.
func (s *BookReviewService) GetByUserID(userID uint) ([]model.BookReview, error) {
	return s.repo.FindByUserID(userID)
}

// GetAll returns paginated book reviews.
func (s *BookReviewService) GetAll(limit, offset int) ([]model.BookReview, int64, error) {
	return s.repo.FindAll(limit, offset)
}

// Update updates a book review after verifying ownership.
func (s *BookReviewService) Update(id, userID uint, updates *model.BookReview) (*model.BookReview, error) {
	review, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if review.UserID != userID {
		return nil, ErrForbidden
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

// Delete deletes a book review after verifying ownership.
func (s *BookReviewService) Delete(id, userID uint) error {
	review, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if review.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}

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
	if err := domain.ValidateStringLength(review.Title, 1, 200, "タイトル"); err != nil {
		return err
	}
	if err := validateRating(review.Rating); err != nil {
		return err
	}
	if len([]rune(strings.TrimSpace(review.Author))) > 200 {
		return domain.NewError(domain.ErrCodeValidation, "著者名は200文字以下である必要があります", nil)
	}
	if len(strings.TrimSpace(review.ISBN)) > 20 {
		return domain.NewError(domain.ErrCodeValidation, "ISBNは20文字以下である必要があります", nil)
	}
	if len(strings.TrimSpace(review.Review)) > 10000 {
		return domain.NewError(domain.ErrCodeValidation, "レビュー本文は10000文字以下である必要があります", nil)
	}
	review.Title = strings.TrimSpace(review.Title)
	review.Review = strings.TrimSpace(review.Review)
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

// Search は書籍レビューをキーワード検索する。
func (s *BookReviewService) Search(query string, limit, offset int) ([]model.BookReview, int64, error) {
	if strings.TrimSpace(query) == "" {
		return nil, 0, domain.NewError(domain.ErrCodeBadRequest, "検索キーワードは必須です", nil)
	}
	return s.repo.Search(strings.TrimSpace(query), limit, offset)
}

// findAndCheckOwnership は書籍レビューを取得し、指定ユーザーが所有者かを検証する。
func (s *BookReviewService) findAndCheckOwnership(id, userID uint) (*model.BookReview, error) {
	return checkOwnership(s.repo.FindByID, id, userID, func(r *model.BookReview) uint { return r.UserID })
}

// Update は所有権を検証した後、書籍レビューを更新する。
func (s *BookReviewService) Update(id, userID uint, updates *model.BookReview) (*model.BookReview, error) {
	review, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(updates.Title) != "" {
		if len(strings.TrimSpace(updates.Title)) > 200 {
			return nil, domain.NewError(domain.ErrCodeValidation, "タイトルは200文字以下である必要があります", nil)
		}
		review.Title = strings.TrimSpace(updates.Title)
	}
	if strings.TrimSpace(updates.Author) != "" {
		if len([]rune(strings.TrimSpace(updates.Author))) > 200 {
			return nil, domain.NewError(domain.ErrCodeValidation, "著者名は200文字以下である必要があります", nil)
		}
		review.Author = strings.TrimSpace(updates.Author)
	}
	if strings.TrimSpace(updates.ISBN) != "" {
		if len(strings.TrimSpace(updates.ISBN)) > 20 {
			return nil, domain.NewError(domain.ErrCodeValidation, "ISBNは20文字以下である必要があります", nil)
		}
		review.ISBN = strings.TrimSpace(updates.ISBN)
	}
	if updates.Rating != 0 {
		if err := validateRating(updates.Rating); err != nil {
			return nil, err
		}
		review.Rating = updates.Rating
	}
	if strings.TrimSpace(updates.Review) != "" {
		if len(strings.TrimSpace(updates.Review)) > 10000 {
			return nil, domain.NewError(domain.ErrCodeValidation, "レビュー本文は10000文字以下である必要があります", nil)
		}
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

// UpdateStatus は所有権を検証した後、書籍レビューの読書状態を更新する。
func (s *BookReviewService) UpdateStatus(id, userID uint, status model.ReviewStatus) error {
	if !model.ValidReviewStatuses[status] {
		return domain.NewError(domain.ErrCodeBadRequest, "無効なステータスです", nil)
	}
	review, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return err
	}
	review.Status = status
	return s.repo.Update(review)
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

// UpdateProgress は所有権を検証した後、読書進捗を更新する。
// 総ページ数に到達した場合は自動的にステータスを「読了」に変更する。
// 未読状態で進捗が1以上になった場合は自動的に「読中」に変更する。
func (s *BookReviewService) UpdateProgress(id, userID uint, currentPage int) (*model.BookReview, error) {
	if currentPage < 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "ページ数は0以上で指定してください", nil)
	}

	review, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if review.TotalPages == 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "総ページ数が設定されていません", nil)
	}

	if currentPage > review.TotalPages {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "総ページ数を超えることはできません", nil)
	}

	review.CurrentPage = currentPage

	// ステータス自動更新
	if currentPage >= review.TotalPages {
		review.Status = model.ReviewStatusCompleted
	} else if currentPage > 0 && review.Status == model.ReviewStatusNotStarted {
		review.Status = model.ReviewStatusReading
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

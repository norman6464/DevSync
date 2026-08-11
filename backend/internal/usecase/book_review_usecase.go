package usecase

import (
	"context"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// maxBookTotalPages は総ページ数の上限。
const maxBookTotalPages = 99999

// bookReviewOwnerOf は書籍レビューの所有者 ID を返す。
func bookReviewOwnerOf(r *model.BookReview) uint { return r.UserID }

// CreateBookReviewUseCase は書籍レビューを作成する。
type CreateBookReviewUseCase struct {
	reviews repository.BookReviewRepository
}

// NewCreateBookReviewUseCase は CreateBookReviewUseCase を生成する。
func NewCreateBookReviewUseCase(reviews repository.BookReviewRepository) *CreateBookReviewUseCase {
	return &CreateBookReviewUseCase{reviews: reviews}
}

// Execute は入力を検証し、前後の空白を除いてレビューを作成する。
func (uc *CreateBookReviewUseCase) Execute(ctx context.Context, review *model.BookReview) error {
	if err := domain.ValidateStringLength(review.Title, 1, 200, "タイトル"); err != nil {
		return err
	}
	if err := domain.ValidateRating(review.Rating); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(review.Author, 0, 200, "著者名"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(review.ISBN, 0, 20, "ISBN"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(review.Review, 0, 10000, "レビュー本文"); err != nil {
		return err
	}
	if review.TotalPages < 0 || review.TotalPages > maxBookTotalPages {
		return domain.NewError(domain.ErrCodeValidation, "総ページ数は0〜99999の範囲で指定してください", nil)
	}

	review.Title = strings.TrimSpace(review.Title)
	review.Author = strings.TrimSpace(review.Author)
	review.ISBN = strings.TrimSpace(review.ISBN)
	review.Review = strings.TrimSpace(review.Review)

	return uc.reviews.Create(ctx, review)
}

// GetBookReviewUseCase は指定 ID の書籍レビューを取得する。
type GetBookReviewUseCase struct {
	reviews repository.BookReviewRepository
}

// NewGetBookReviewUseCase は GetBookReviewUseCase を生成する。
func NewGetBookReviewUseCase(reviews repository.BookReviewRepository) *GetBookReviewUseCase {
	return &GetBookReviewUseCase{reviews: reviews}
}

// Execute はレビューを返す。所有権は検証しない（移行前の挙動を維持している）。
func (uc *GetBookReviewUseCase) Execute(ctx context.Context, id uint) (*model.BookReview, error) {
	review, err := uc.reviews.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if review == nil {
		// 不在は DomainError にせず 500 のままにする（移行前の挙動を維持している）。
		return nil, errOwnedEntityNotFound
	}
	return review, nil
}

// ListBookReviewsByUserUseCase は指定ユーザーのレビューをページネーション付きで取得する。
type ListBookReviewsByUserUseCase struct {
	reviews repository.BookReviewRepository
}

// NewListBookReviewsByUserUseCase は ListBookReviewsByUserUseCase を生成する。
func NewListBookReviewsByUserUseCase(reviews repository.BookReviewRepository) *ListBookReviewsByUserUseCase {
	return &ListBookReviewsByUserUseCase{reviews: reviews}
}

// Execute はレビュー一覧と総件数を返す。
func (uc *ListBookReviewsByUserUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.BookReview, int64, error) {
	return uc.reviews.FindByUserID(ctx, userID, limit, offset)
}

// ListAllBookReviewsUseCase は全レビューをページネーション付きで取得する。
type ListAllBookReviewsUseCase struct {
	reviews repository.BookReviewRepository
}

// NewListAllBookReviewsUseCase は ListAllBookReviewsUseCase を生成する。
func NewListAllBookReviewsUseCase(reviews repository.BookReviewRepository) *ListAllBookReviewsUseCase {
	return &ListAllBookReviewsUseCase{reviews: reviews}
}

// Execute はレビュー一覧と総件数を返す。
func (uc *ListAllBookReviewsUseCase) Execute(ctx context.Context, limit, offset int) ([]model.BookReview, int64, error) {
	return uc.reviews.FindAll(ctx, limit, offset)
}

// ListBookReviewsByRatingUseCase は評価範囲でレビューを絞り込む。
type ListBookReviewsByRatingUseCase struct {
	reviews repository.BookReviewRepository
}

// NewListBookReviewsByRatingUseCase は ListBookReviewsByRatingUseCase を生成する。
func NewListBookReviewsByRatingUseCase(reviews repository.BookReviewRepository) *ListBookReviewsByRatingUseCase {
	return &ListBookReviewsByRatingUseCase{reviews: reviews}
}

// Execute は評価範囲を検証したうえでレビュー一覧を返す。
func (uc *ListBookReviewsByRatingUseCase) Execute(ctx context.Context, userID uint, minRating, maxRating int) ([]model.BookReview, error) {
	if minRating < 1 || minRating > 5 || maxRating < 1 || maxRating > 5 || minRating > maxRating {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "評価範囲が無効です", nil)
	}
	return uc.reviews.FindByRating(ctx, userID, minRating, maxRating)
}

// SearchBookReviewsUseCase はレビューをキーワード検索する。
type SearchBookReviewsUseCase struct {
	reviews repository.BookReviewRepository
}

// NewSearchBookReviewsUseCase は SearchBookReviewsUseCase を生成する。
func NewSearchBookReviewsUseCase(reviews repository.BookReviewRepository) *SearchBookReviewsUseCase {
	return &SearchBookReviewsUseCase{reviews: reviews}
}

// Execute は検索結果と総件数を返す。
func (uc *SearchBookReviewsUseCase) Execute(ctx context.Context, query string, limit, offset int) ([]model.BookReview, int64, error) {
	q, err := validateSearchQuery(query)
	if err != nil {
		return nil, 0, err
	}
	return uc.reviews.Search(ctx, q, limit, offset)
}

// UpdateBookReviewUseCase は所有者の書籍レビューを更新する。
type UpdateBookReviewUseCase struct {
	reviews repository.BookReviewRepository
}

// NewUpdateBookReviewUseCase は UpdateBookReviewUseCase を生成する。
func NewUpdateBookReviewUseCase(reviews repository.BookReviewRepository) *UpdateBookReviewUseCase {
	return &UpdateBookReviewUseCase{reviews: reviews}
}

// Execute は所有権を検証したうえでレビューを部分更新する。
// 前後の空白を除いて空になる文字列と 0 の数値は据え置く。
func (uc *UpdateBookReviewUseCase) Execute(ctx context.Context, id, userID uint, updates *model.BookReview) (*model.BookReview, error) {
	review, err := ensureOwner(ctx, uc.reviews.FindByID, id, userID, bookReviewOwnerOf)
	if err != nil {
		return nil, err
	}

	if title := strings.TrimSpace(updates.Title); title != "" {
		if err := domain.ValidateStringLength(title, 1, 200, "タイトル"); err != nil {
			return nil, err
		}
		review.Title = title
	}
	if author := strings.TrimSpace(updates.Author); author != "" {
		if err := domain.ValidateStringLength(author, 1, 200, "著者名"); err != nil {
			return nil, err
		}
		review.Author = author
	}
	if isbn := strings.TrimSpace(updates.ISBN); isbn != "" {
		if err := domain.ValidateStringLength(isbn, 1, 20, "ISBN"); err != nil {
			return nil, err
		}
		review.ISBN = isbn
	}
	if updates.Rating != 0 {
		if err := domain.ValidateRating(updates.Rating); err != nil {
			return nil, err
		}
		review.Rating = updates.Rating
	}
	if reviewText := strings.TrimSpace(updates.Review); reviewText != "" {
		if err := domain.ValidateStringLength(reviewText, 1, 10000, "レビュー本文"); err != nil {
			return nil, err
		}
		review.Review = reviewText
	}
	if updates.TotalPages != 0 {
		if updates.TotalPages < 0 || updates.TotalPages > maxBookTotalPages {
			return nil, domain.NewError(domain.ErrCodeValidation, "総ページ数は0〜99999の範囲で指定してください", nil)
		}
		review.TotalPages = updates.TotalPages
	}
	if imageURL := strings.TrimSpace(updates.ImageURL); imageURL != "" {
		if err := domain.ValidateStringLength(imageURL, 1, 2000, "画像URL"); err != nil {
			return nil, err
		}
		review.ImageURL = imageURL
	}

	if err := uc.reviews.Update(ctx, review); err != nil {
		return nil, err
	}
	return review, nil
}

// UpdateBookReviewStatusUseCase は書籍レビューの読書状態を更新する。
type UpdateBookReviewStatusUseCase struct {
	reviews repository.BookReviewRepository
}

// NewUpdateBookReviewStatusUseCase は UpdateBookReviewStatusUseCase を生成する。
func NewUpdateBookReviewStatusUseCase(reviews repository.BookReviewRepository) *UpdateBookReviewStatusUseCase {
	return &UpdateBookReviewStatusUseCase{reviews: reviews}
}

// Execute はステータスを検証し、所有権を確認したうえで更新する。
func (uc *UpdateBookReviewStatusUseCase) Execute(ctx context.Context, id, userID uint, status model.ReviewStatus) error {
	if !model.ValidReviewStatuses[status] {
		return domain.NewError(domain.ErrCodeBadRequest, "無効なステータスです", nil)
	}
	review, err := ensureOwner(ctx, uc.reviews.FindByID, id, userID, bookReviewOwnerOf)
	if err != nil {
		return err
	}
	review.Status = status
	return uc.reviews.Update(ctx, review)
}

// ArchiveBookReviewUseCase は書籍レビューのアーカイブ状態を切り替える。
type ArchiveBookReviewUseCase struct {
	reviews repository.BookReviewRepository
}

// NewArchiveBookReviewUseCase は ArchiveBookReviewUseCase を生成する。
func NewArchiveBookReviewUseCase(reviews repository.BookReviewRepository) *ArchiveBookReviewUseCase {
	return &ArchiveBookReviewUseCase{reviews: reviews}
}

// Execute は所有権を検証したうえでアーカイブ状態を設定する。
func (uc *ArchiveBookReviewUseCase) Execute(ctx context.Context, id, userID uint, archived bool) error {
	review, err := ensureOwner(ctx, uc.reviews.FindByID, id, userID, bookReviewOwnerOf)
	if err != nil {
		return err
	}
	review.IsArchived = archived
	return uc.reviews.Update(ctx, review)
}

// UpdateBookReviewProgressUseCase は読書進捗を更新する。
type UpdateBookReviewProgressUseCase struct {
	reviews repository.BookReviewRepository
}

// NewUpdateBookReviewProgressUseCase は UpdateBookReviewProgressUseCase を生成する。
func NewUpdateBookReviewProgressUseCase(reviews repository.BookReviewRepository) *UpdateBookReviewProgressUseCase {
	return &UpdateBookReviewProgressUseCase{reviews: reviews}
}

// Execute は進捗を更新する。総ページ数に到達したら読了、未読から 1 以上になったら読中へ自動遷移する。
func (uc *UpdateBookReviewProgressUseCase) Execute(ctx context.Context, id, userID uint, currentPage int) (*model.BookReview, error) {
	// 検証の順序（負数 → 所有権 → 総ページ数未設定 → 超過）は移行前と同じ。
	if currentPage < 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "ページ数は0以上で指定してください", nil)
	}

	review, err := ensureOwner(ctx, uc.reviews.FindByID, id, userID, bookReviewOwnerOf)
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
	if currentPage >= review.TotalPages {
		review.Status = model.ReviewStatusCompleted
	} else if currentPage > 0 && review.Status == model.ReviewStatusNotStarted {
		review.Status = model.ReviewStatusReading
	}

	if err := uc.reviews.Update(ctx, review); err != nil {
		return nil, err
	}
	return review, nil
}

// DeleteBookReviewUseCase は所有者の書籍レビューを削除する。
type DeleteBookReviewUseCase struct {
	reviews repository.BookReviewRepository
}

// NewDeleteBookReviewUseCase は DeleteBookReviewUseCase を生成する。
func NewDeleteBookReviewUseCase(reviews repository.BookReviewRepository) *DeleteBookReviewUseCase {
	return &DeleteBookReviewUseCase{reviews: reviews}
}

// Execute は所有権を検証したうえでレビューを削除する。
func (uc *DeleteBookReviewUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensureOwner(ctx, uc.reviews.FindByID, id, userID, bookReviewOwnerOf); err != nil {
		return err
	}
	return uc.reviews.Delete(ctx, id)
}

// CountBookReviewsUseCase は指定ユーザーのレビュー総数を返す。
type CountBookReviewsUseCase struct {
	reviews repository.BookReviewRepository
}

// NewCountBookReviewsUseCase は CountBookReviewsUseCase を生成する。
func NewCountBookReviewsUseCase(reviews repository.BookReviewRepository) *CountBookReviewsUseCase {
	return &CountBookReviewsUseCase{reviews: reviews}
}

// Execute はレビュー総数を返す。
func (uc *CountBookReviewsUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.reviews.CountByUserID(ctx, userID)
}

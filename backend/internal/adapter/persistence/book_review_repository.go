package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// bookReviewRepository は [repository.BookReviewRepository] の sqlc(pgx) 実装。
// BookReview は論理削除を使うため、全クエリで deleted_at IS NULL を明示する
// （GORM が自動的に付与していたスコープ相当）。
type bookReviewRepository struct {
	q *sqlcgen.Queries
}

// NewBookReviewRepository は BookReviewRepository の sqlc(pgx) 実装を返す。
func NewBookReviewRepository(q *sqlcgen.Queries) repository.BookReviewRepository {
	return &bookReviewRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.BookReviewRepository = (*bookReviewRepository)(nil)

func toModelBookReview(row sqlcgen.BookReview) model.BookReview {
	return model.BookReview{
		ID:          uint(row.ID),
		UserID:      uint(row.UserID),
		Title:       row.Title,
		Author:      fromStringPtr(row.Author),
		ISBN:        fromStringPtr(row.Isbn),
		Rating:      int(row.Rating),
		Review:      fromStringPtr(row.Review),
		TotalPages:  int(fromInt64PtrValue(row.TotalPages)),
		CurrentPage: int(fromInt64PtrValue(row.CurrentPage)),
		ImageURL:    fromStringPtr(row.ImageUrl),
		Status:      model.ReviewStatus(fromStringPtr(row.Status)),
		IsArchived:  row.IsArchived,
		CreatedAt:   timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:   timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// Create は新しい書籍レビューを作成する。
// Status はGORMの `gorm:"default:'not_started'"` に相当し、未指定（ゼロ値）なら not_started を補う。
func (r *bookReviewRepository) Create(ctx context.Context, review *model.BookReview) error {
	status := review.Status
	if status == "" {
		status = model.ReviewStatusNotStarted
	}

	row, err := r.q.CreateBookReview(ctx, sqlcgen.CreateBookReviewParams{
		UserID:      int64(review.UserID),
		Title:       review.Title,
		Author:      &review.Author,
		Isbn:        &review.ISBN,
		Rating:      int64(review.Rating),
		Review:      &review.Review,
		TotalPages:  toInt64Ptr(review.TotalPages),
		CurrentPage: toInt64Ptr(review.CurrentPage),
		ImageUrl:    &review.ImageURL,
		Status:      (*string)(&status),
		IsArchived:  review.IsArchived,
	})
	if err != nil {
		return err
	}
	*review = toModelBookReview(row)
	return nil
}

// FindByID は指定 ID のレビューをユーザー情報付きで取得する。不在の場合は (nil, nil) を返す。
func (r *bookReviewRepository) FindByID(ctx context.Context, id uint) (*model.BookReview, error) {
	row, err := r.q.GetBookReviewWithUserByID(ctx, int64(id))
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	review := toModelBookReview(row.BookReview)
	review.User = toModelUser(row.User)
	return &review, nil
}

// FindByUserID は指定ユーザーのレビューをページネーション付きで取得する（新しい順）。
func (r *bookReviewRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.BookReview, int64, error) {
	total, err := r.q.CountBookReviewsByUser(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListBookReviewsByUser(ctx, sqlcgen.ListBookReviewsByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	reviews := make([]model.BookReview, len(rows))
	for i, row := range rows {
		reviews[i] = toModelBookReview(row)
	}
	return reviews, total, nil
}

// FindAll は全レビューをページネーション付きで取得する（新しい順）。
func (r *bookReviewRepository) FindAll(ctx context.Context, limit, offset int) ([]model.BookReview, int64, error) {
	total, err := r.q.CountAllBookReviews(ctx)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListAllBookReviewsWithUser(ctx, sqlcgen.ListAllBookReviewsWithUserParams{
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	reviews := make([]model.BookReview, len(rows))
	for i, row := range rows {
		reviews[i] = toModelBookReview(row.BookReview)
		reviews[i].User = toModelUser(row.User)
	}
	return reviews, total, nil
}

// FindByRating は指定ユーザーのレビューを評価範囲で絞り込んで取得する（新しい順）。
func (r *bookReviewRepository) FindByRating(ctx context.Context, userID uint, minRating, maxRating int) ([]model.BookReview, error) {
	rows, err := r.q.ListBookReviewsByRating(ctx, sqlcgen.ListBookReviewsByRatingParams{
		UserID:   int64(userID),
		Rating:   int64(minRating),
		Rating_2: int64(maxRating),
	})
	if err != nil {
		return nil, err
	}
	reviews := make([]model.BookReview, len(rows))
	for i, row := range rows {
		reviews[i] = toModelBookReview(row)
	}
	return reviews, nil
}

// Search はタイトル・著者名・ISBN からレビューをキーワード検索する（新しい順）。
func (r *bookReviewRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.BookReview, int64, error) {
	like := escapeLikePattern(query)

	total, err := r.q.CountSearchBookReviews(ctx, like)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.SearchBookReviews(ctx, sqlcgen.SearchBookReviewsParams{
		Title:  like,
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	reviews := make([]model.BookReview, len(rows))
	for i, row := range rows {
		reviews[i] = toModelBookReview(row.BookReview)
		reviews[i].User = toModelUser(row.User)
	}
	return reviews, total, nil
}

// Update は既存のレビューを更新する（GORMのSave＝全カラム上書きに相当）。
func (r *bookReviewRepository) Update(ctx context.Context, review *model.BookReview) error {
	row, err := r.q.UpdateBookReview(ctx, sqlcgen.UpdateBookReviewParams{
		ID:          int64(review.ID),
		Title:       review.Title,
		Author:      &review.Author,
		Isbn:        &review.ISBN,
		Rating:      int64(review.Rating),
		Review:      &review.Review,
		TotalPages:  toInt64Ptr(review.TotalPages),
		CurrentPage: toInt64Ptr(review.CurrentPage),
		ImageUrl:    &review.ImageURL,
		Status:      (*string)(&review.Status),
		IsArchived:  review.IsArchived,
	})
	if err != nil {
		return err
	}
	*review = toModelBookReview(row)
	return nil
}

// Delete はレビューを削除する（論理削除）。
func (r *bookReviewRepository) Delete(ctx context.Context, id uint) error {
	return r.q.DeleteBookReview(ctx, int64(id))
}

// CountByUserID は指定ユーザーのレビュー総数を返す。
func (r *bookReviewRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountBookReviewsByUser(ctx, int64(userID))
}

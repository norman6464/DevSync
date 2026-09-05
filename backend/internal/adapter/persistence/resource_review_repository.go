package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// resourceReviewRepository は [repository.ResourceReviewRepository] の sqlc(pgx) 実装。
type resourceReviewRepository struct {
	q *sqlcgen.Queries
}

// NewResourceReviewRepository は ResourceReviewRepository の sqlc(pgx) 実装を返す。
func NewResourceReviewRepository(q *sqlcgen.Queries) repository.ResourceReviewRepository {
	return &resourceReviewRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ResourceReviewRepository = (*resourceReviewRepository)(nil)

// toModelResourceReview は sqlc の生成行を model.ResourceReview へ変換する（User は含まない）。
func toModelResourceReview(row sqlcgen.ResourceReview) model.ResourceReview {
	return model.ResourceReview{
		ID:         uint(row.ID),
		UserID:     uint(row.UserID),
		ResourceID: uint(row.ResourceID),
		Rating:     int(row.Rating),
		Comment:    fromStringPtr(row.Comment),
		CreatedAt:  timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:  timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// Create はレビューを保存する。複合ユニーク索引（user_id, resource_id）に
// 衝突した場合は domain.ErrConflict を返す。usecase の重複チェックは同時実行を
// すり抜けるため、最後の砦は DB の制約に委ねる。
func (r *resourceReviewRepository) Create(ctx context.Context, review *model.ResourceReview) error {
	row, err := r.q.CreateResourceReview(ctx, sqlcgen.CreateResourceReviewParams{
		UserID:     int64(review.UserID),
		ResourceID: int64(review.ResourceID),
		Rating:     int64(review.Rating),
		Comment:    &review.Comment,
	})
	if isUniqueViolation(err) {
		return domain.ErrConflict
	}
	if err != nil {
		return err
	}
	*review = toModelResourceReview(row)
	return nil
}

// FindByID は指定 ID のレビューをユーザー情報付きで取得する。不在の場合は (nil, nil) を返す。
func (r *resourceReviewRepository) FindByID(ctx context.Context, id uint) (*model.ResourceReview, error) {
	row, err := r.q.GetResourceReviewWithUserByID(ctx, int64(id))
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	review := toModelResourceReview(row.ResourceReview)
	review.User = toModelUser(row.User)
	return &review, nil
}

// FindByResourceID は指定リソースのレビュー一覧をページネーション付きで取得する。
func (r *resourceReviewRepository) FindByResourceID(ctx context.Context, resourceID uint, limit, offset int) ([]model.ResourceReview, int64, error) {
	total, err := r.q.CountResourceReviewsByResource(ctx, int64(resourceID))
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListResourceReviewsByResource(ctx, sqlcgen.ListResourceReviewsByResourceParams{
		ResourceID: int64(resourceID),
		Limit:      int32Param(limit),
		Offset:     int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	reviews := make([]model.ResourceReview, len(rows))
	for i, row := range rows {
		reviews[i] = toModelResourceReview(row.ResourceReview)
		reviews[i].User = toModelUser(row.User)
	}
	return reviews, total, nil
}

// FindByUserAndResource は指定ユーザーの指定リソースへのレビューを取得する。不在の場合は (nil, nil) を返す。
func (r *resourceReviewRepository) FindByUserAndResource(ctx context.Context, userID, resourceID uint) (*model.ResourceReview, error) {
	row, err := r.q.GetResourceReviewByUserAndResource(ctx, sqlcgen.GetResourceReviewByUserAndResourceParams{
		UserID:     int64(userID),
		ResourceID: int64(resourceID),
	})
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	review := toModelResourceReview(row)
	return &review, nil
}

// Update はレビューを更新する（GORMのSave＝全カラム上書きに相当）。
func (r *resourceReviewRepository) Update(ctx context.Context, review *model.ResourceReview) error {
	row, err := r.q.UpdateResourceReview(ctx, sqlcgen.UpdateResourceReviewParams{
		ID:      int64(review.ID),
		Rating:  int64(review.Rating),
		Comment: &review.Comment,
	})
	if err != nil {
		return err
	}
	*review = toModelResourceReview(row)
	return nil
}

// Delete はレビューを削除する。
func (r *resourceReviewRepository) Delete(ctx context.Context, id uint) error {
	return r.q.DeleteResourceReview(ctx, int64(id))
}

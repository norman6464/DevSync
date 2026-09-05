package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postSeriesRepository は [repository.PostSeriesRepository] の sqlc(pgx) 実装。
type postSeriesRepository struct {
	q *sqlcgen.Queries
}

// NewPostSeriesRepository は PostSeriesRepository の sqlc(pgx) 実装を返す。
func NewPostSeriesRepository(q *sqlcgen.Queries) repository.PostSeriesRepository {
	return &postSeriesRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.PostSeriesRepository = (*postSeriesRepository)(nil)

// toModelPostSeries は sqlc の生成行を model.PostSeries へ変換する（User は含まない）。
func toModelPostSeries(row sqlcgen.PostSeries) model.PostSeries {
	return model.PostSeries{
		ID:          uint(row.ID),
		UserID:      uint(row.UserID),
		Title:       row.Title,
		Description: fromStringPtr(row.Description),
		CreatedAt:   timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:   timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// Create は新しい投稿シリーズをデータベースに作成する。
func (r *postSeriesRepository) Create(ctx context.Context, series *model.PostSeries) error {
	row, err := r.q.CreatePostSeries(ctx, sqlcgen.CreatePostSeriesParams{
		UserID:      int64(series.UserID),
		Title:       series.Title,
		Description: &series.Description,
	})
	if err != nil {
		return err
	}
	*series = toModelPostSeries(row)
	return nil
}

// FindByID は指定IDのシリーズをユーザー情報付きで取得する。
//
// 呼び出し側（GetPostSeriesUseCase）が「不在は非 nil の error」という
// 移行前の GORM 実装の挙動にそのまま依存しているため、(nil, nil) には
// 変換せず pgx.ErrNoRows をそのまま返す。
func (r *postSeriesRepository) FindByID(ctx context.Context, id uint) (*model.PostSeries, error) {
	row, err := r.q.GetPostSeriesWithUserByID(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	series := toModelPostSeries(row.PostSeries)
	series.User = toModelUser(row.User)
	return &series, nil
}

// FindByUserID は指定ユーザーのシリーズをページネーション付きで取得する（新しい順）。
func (r *postSeriesRepository) FindByUserID(ctx context.Context, userID uint, offset, limit int) ([]model.PostSeries, error) {
	rows, err := r.q.ListPostSeriesByUser(ctx, sqlcgen.ListPostSeriesByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, err
	}
	series := make([]model.PostSeries, len(rows))
	for i, row := range rows {
		series[i] = toModelPostSeries(row)
	}
	return series, nil
}

// CountByUser は指定ユーザーのシリーズ数をカウントする。
func (r *postSeriesRepository) CountByUser(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountPostSeriesByUser(ctx, int64(userID))
}

// Update は既存のシリーズを更新する（GORMのSave＝全カラム上書きに相当）。
func (r *postSeriesRepository) Update(ctx context.Context, series *model.PostSeries) error {
	row, err := r.q.UpdatePostSeries(ctx, sqlcgen.UpdatePostSeriesParams{
		ID:          int64(series.ID),
		Title:       series.Title,
		Description: &series.Description,
	})
	if err != nil {
		return err
	}
	*series = toModelPostSeries(row)
	return nil
}

// Delete は指定IDのシリーズとその関連アイテムを削除する。
func (r *postSeriesRepository) Delete(ctx context.Context, id uint) error {
	if err := r.q.DeletePostSeriesItemsBySeriesID(ctx, int64(id)); err != nil {
		return err
	}
	return r.q.DeletePostSeries(ctx, int64(id))
}

// AddPost はシリーズに投稿を追加する。
func (r *postSeriesRepository) AddPost(ctx context.Context, item *model.PostSeriesItem) error {
	row, err := r.q.CreatePostSeriesItem(ctx, sqlcgen.CreatePostSeriesItemParams{
		SeriesID:   int64(item.SeriesID),
		PostID:     int64(item.PostID),
		OrderIndex: int64(item.OrderIndex),
	})
	if err != nil {
		return err
	}
	*item = model.PostSeriesItem{
		ID:         uint(row.ID),
		SeriesID:   uint(row.SeriesID),
		PostID:     uint(row.PostID),
		OrderIndex: int(row.OrderIndex),
	}
	return nil
}

// HasPost は指定シリーズに指定投稿が存在するかを確認する。
func (r *postSeriesRepository) HasPost(ctx context.Context, seriesID, postID uint) (bool, error) {
	count, err := r.q.CountPostSeriesItemsBySeriesAndPost(ctx, sqlcgen.CountPostSeriesItemsBySeriesAndPostParams{
		SeriesID: int64(seriesID),
		PostID:   int64(postID),
	})
	return count > 0, err
}

// RemovePost はシリーズから投稿を削除する。
func (r *postSeriesRepository) RemovePost(ctx context.Context, seriesID, postID uint) error {
	return r.q.DeletePostSeriesItem(ctx, sqlcgen.DeletePostSeriesItemParams{
		SeriesID: int64(seriesID),
		PostID:   int64(postID),
	})
}

// GetPostsBySeriesID は指定シリーズの投稿一覧を順序付きで取得する。
func (r *postSeriesRepository) GetPostsBySeriesID(ctx context.Context, seriesID uint) ([]model.PostSeriesItem, error) {
	rows, err := r.q.ListPostSeriesItemsWithPostBySeriesID(ctx, int64(seriesID))
	if err != nil {
		return nil, err
	}
	items := make([]model.PostSeriesItem, len(rows))
	for i, row := range rows {
		post := toModelPost(row.Post)
		post.User = toModelUser(row.User)
		items[i] = model.PostSeriesItem{
			ID:         uint(row.PostSeriesItem.ID),
			SeriesID:   uint(row.PostSeriesItem.SeriesID),
			PostID:     uint(row.PostSeriesItem.PostID),
			Post:       post,
			OrderIndex: int(row.PostSeriesItem.OrderIndex),
		}
	}
	return items, nil
}

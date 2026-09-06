package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postTagRepository は [repository.PostTagRepository] の sqlc(pgx) 実装。
// SetTags は削除と挿入を1トランザクションで行うため、
// Queries だけでなくトランザクションを開始できる *pgxpool.Pool を直接保持する。
type postTagRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewPostTagRepository は PostTagRepository の sqlc(pgx) 実装を返す。
func NewPostTagRepository(pool *pgxpool.Pool) repository.PostTagRepository {
	return &postTagRepository{pool: pool, q: sqlcgen.New(pool)}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.PostTagRepository = (*postTagRepository)(nil)

// SetTags は投稿のタグを全て置き換える（削除と挿入を 1 トランザクションで行う）。
func (r *postTagRepository) SetTags(ctx context.Context, postID uint, tags []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	if err := q.DeletePostTagsByPostID(ctx, int64(postID)); err != nil {
		return err
	}
	for _, tag := range tags {
		if err := q.CreatePostTag(ctx, sqlcgen.CreatePostTagParams{
			PostID: int64(postID),
			Tag:    tag,
		}); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// GetByPostID は投稿のタグ一覧を取得する。
func (r *postTagRepository) GetByPostID(ctx context.Context, postID uint) ([]string, error) {
	return r.q.ListPostTagsByPostID(ctx, int64(postID))
}

// FindPostsByTag はタグで投稿を検索する（下書きは除外する）。
func (r *postTagRepository) FindPostsByTag(ctx context.Context, tag string, limit, offset int) ([]model.Post, int64, error) {
	count, err := r.q.CountPostsByTag(ctx, tag)
	if err != nil {
		return nil, 0, err
	}

	postIDs, err := r.q.ListPostIDsByTag(ctx, sqlcgen.ListPostIDsByTagParams{
		Tag:    tag,
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	if len(postIDs) == 0 {
		// 範囲外の offset でも算出済みの総件数を返す（ページネーション表示を壊さない）。
		return []model.Post{}, count, nil
	}

	rows, err := r.q.ListPostsByIDs(ctx, postIDs)
	if err != nil {
		return nil, 0, err
	}
	posts := make([]model.Post, len(rows))
	for i, row := range rows {
		posts[i] = toModelPost(row.Post)
		posts[i].User = toModelUser(row.User)
	}
	if err := attachBookmarkCountsToPosts(ctx, r.q, posts); err != nil {
		return nil, 0, err
	}
	if err := attachMetricsToPosts(ctx, r.q, posts); err != nil {
		return nil, 0, err
	}
	return posts, count, nil
}

// GetPopularTags は使用回数の多いタグ一覧を取得する。
func (r *postTagRepository) GetPopularTags(ctx context.Context, limit int) ([]model.TagCount, error) {
	rows, err := r.q.ListPopularTags(ctx, int32Param(limit))
	if err != nil {
		return nil, err
	}
	results := make([]model.TagCount, len(rows))
	for i, row := range rows {
		results[i] = model.TagCount{Tag: row.Tag, Count: int(row.Count)}
	}
	return results, nil
}

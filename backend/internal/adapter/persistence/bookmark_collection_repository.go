package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// bookmarkCollectionRepository は [repository.BookmarkCollectionRepository] の sqlc(pgx) 実装。
type bookmarkCollectionRepository struct {
	q *sqlcgen.Queries
}

// NewBookmarkCollectionRepository は BookmarkCollectionRepository の sqlc(pgx) 実装を返す。
func NewBookmarkCollectionRepository(q *sqlcgen.Queries) repository.BookmarkCollectionRepository {
	return &bookmarkCollectionRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.BookmarkCollectionRepository = (*bookmarkCollectionRepository)(nil)

func toModelBookmarkCollection(row sqlcgen.BookmarkCollection) model.BookmarkCollection {
	return model.BookmarkCollection{
		ID:          uint(row.ID),
		UserID:      uint(row.UserID),
		Name:        row.Name,
		Description: fromStringPtr(row.Description),
		Color:       fromStringPtr(row.Color),
		CreatedAt:   timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:   timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// Create は新しいブックマークコレクションを作成する。
func (r *bookmarkCollectionRepository) Create(ctx context.Context, collection *model.BookmarkCollection) error {
	row, err := r.q.CreateBookmarkCollection(ctx, sqlcgen.CreateBookmarkCollectionParams{
		UserID:      int64(collection.UserID),
		Name:        collection.Name,
		Description: &collection.Description,
		Color:       &collection.Color,
	})
	if err != nil {
		return err
	}
	*collection = toModelBookmarkCollection(row)
	return nil
}

// FindByID は指定IDのコレクションを取得する。
//
// 呼び出し側が「不在は非 nil の error」という移行前の GORM 実装の挙動に
// そのまま依存しているため、(nil, nil) には変換せず pgx.ErrNoRows をそのまま返す。
func (r *bookmarkCollectionRepository) FindByID(ctx context.Context, id uint) (*model.BookmarkCollection, error) {
	row, err := r.q.GetBookmarkCollectionByID(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	collection := toModelBookmarkCollection(row)
	return &collection, nil
}

// FindByUserID は指定ユーザーのコレクション一覧を取得する。
func (r *bookmarkCollectionRepository) FindByUserID(ctx context.Context, userID uint) ([]model.BookmarkCollection, error) {
	rows, err := r.q.ListBookmarkCollectionsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	collections := make([]model.BookmarkCollection, len(rows))
	for i, row := range rows {
		collections[i] = toModelBookmarkCollection(row)
	}
	return collections, nil
}

// Update はコレクションを更新する（GORMのSave＝全カラム上書きに相当）。
func (r *bookmarkCollectionRepository) Update(ctx context.Context, collection *model.BookmarkCollection) error {
	row, err := r.q.UpdateBookmarkCollection(ctx, sqlcgen.UpdateBookmarkCollectionParams{
		ID:          int64(collection.ID),
		Name:        collection.Name,
		Description: &collection.Description,
		Color:       &collection.Color,
	})
	if err != nil {
		return err
	}
	*collection = toModelBookmarkCollection(row)
	return nil
}

// Delete はコレクションとそのアイテムを削除する。
func (r *bookmarkCollectionRepository) Delete(ctx context.Context, id uint) error {
	if err := r.q.DeleteBookmarkCollectionItemsByCollectionID(ctx, int64(id)); err != nil {
		return err
	}
	return r.q.DeleteBookmarkCollection(ctx, int64(id))
}

// AddPost はコレクションに投稿を追加する。(collection_id, post_id) の一意制約に任せて
// ON CONFLICT DO NOTHING で挿入し、既に入っていた場合は (false, nil) を返す。
// 「存在確認 → 挿入」の 2 クエリと違い、同時実行でも重複行は作られない。
func (r *bookmarkCollectionRepository) AddPost(ctx context.Context, item *model.BookmarkCollectionItem) (bool, error) {
	rowsAffected, err := r.q.CreateBookmarkCollectionItem(ctx, sqlcgen.CreateBookmarkCollectionItemParams{
		CollectionID: int64(item.CollectionID),
		PostID:       int64(item.PostID),
	})
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}

// RemovePost はコレクションから投稿を削除する。
func (r *bookmarkCollectionRepository) RemovePost(ctx context.Context, collectionID, postID uint) error {
	return r.q.DeleteBookmarkCollectionItem(ctx, sqlcgen.DeleteBookmarkCollectionItemParams{
		CollectionID: int64(collectionID),
		PostID:       int64(postID),
	})
}

// GetPosts はコレクション内の投稿一覧を取得する。
func (r *bookmarkCollectionRepository) GetPosts(ctx context.Context, collectionID uint, limit, offset int) ([]model.Post, int64, error) {
	total, err := r.q.CountBookmarkCollectionItemsByCollection(ctx, int64(collectionID))
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListBookmarkCollectionItemsWithPostByCollection(ctx, sqlcgen.ListBookmarkCollectionItemsWithPostByCollectionParams{
		CollectionID: int64(collectionID),
		Limit:        int32Param(limit),
		Offset:       int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	posts := make([]model.Post, 0, len(rows))
	for _, row := range rows {
		posts = append(posts, toModelPost(row.Post))
	}
	return posts, total, nil
}

// CountByUserID は指定ユーザーのコレクション総数を返す。
func (r *bookmarkCollectionRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountBookmarkCollectionsByUser(ctx, int64(userID))
}

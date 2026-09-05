package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postCollectionRepository は [repository.PostCollectionRepository] の sqlc(pgx) 実装。
type postCollectionRepository struct {
	q *sqlcgen.Queries
}

// NewPostCollectionRepository は PostCollectionRepository の sqlc(pgx) 実装を返す。
func NewPostCollectionRepository(q *sqlcgen.Queries) repository.PostCollectionRepository {
	return &postCollectionRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.PostCollectionRepository = (*postCollectionRepository)(nil)

// toModelPostCollection は sqlc の生成行を model.PostCollection へ変換する（User は含まない）。
func toModelPostCollection(row sqlcgen.PostCollection) model.PostCollection {
	return model.PostCollection{
		ID:          uint(row.ID),
		UserID:      uint(row.UserID),
		Title:       row.Title,
		Description: fromStringPtr(row.Description),
		IsPublic:    fromBoolPtr(row.IsPublic),
		CreatedAt:   timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:   timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// Create は新しい投稿コレクションをデータベースに作成する。
func (r *postCollectionRepository) Create(ctx context.Context, collection *model.PostCollection) error {
	row, err := r.q.CreatePostCollection(ctx, sqlcgen.CreatePostCollectionParams{
		UserID:      int64(collection.UserID),
		Title:       collection.Title,
		Description: &collection.Description,
		IsPublic:    &collection.IsPublic,
	})
	if err != nil {
		return err
	}
	*collection = toModelPostCollection(row)
	return nil
}

// FindByID は指定IDのコレクションをユーザー情報付きで取得する。
//
// 呼び出し側が「不在は非 nil の error」という移行前の GORM 実装の挙動に
// そのまま依存しているため、(nil, nil) には変換せず pgx.ErrNoRows をそのまま返す。
func (r *postCollectionRepository) FindByID(ctx context.Context, id uint) (*model.PostCollection, error) {
	row, err := r.q.GetPostCollectionWithUserByID(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	collection := toModelPostCollection(row.PostCollection)
	collection.User = toModelUser(row.User)
	return &collection, nil
}

// FindByUserID は指定ユーザーの全コレクションをページネーション付きで取得する（新しい順）。
func (r *postCollectionRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.PostCollection, int64, error) {
	total, err := r.q.CountPostCollectionsByUser(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListPostCollectionsByUser(ctx, sqlcgen.ListPostCollectionsByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	collections := make([]model.PostCollection, len(rows))
	for i, row := range rows {
		collections[i] = toModelPostCollection(row)
	}
	return collections, total, nil
}

// FindPublicByUserID は指定ユーザーの公開コレクションを取得する（新しい順）。
func (r *postCollectionRepository) FindPublicByUserID(ctx context.Context, userID uint) ([]model.PostCollection, error) {
	rows, err := r.q.ListPublicPostCollectionsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	collections := make([]model.PostCollection, len(rows))
	for i, row := range rows {
		collections[i] = toModelPostCollection(row)
	}
	return collections, nil
}

// Update は既存のコレクションを更新する（GORMのSave＝全カラム上書きに相当）。
func (r *postCollectionRepository) Update(ctx context.Context, collection *model.PostCollection) error {
	row, err := r.q.UpdatePostCollection(ctx, sqlcgen.UpdatePostCollectionParams{
		ID:          int64(collection.ID),
		Title:       collection.Title,
		Description: &collection.Description,
		IsPublic:    &collection.IsPublic,
	})
	if err != nil {
		return err
	}
	*collection = toModelPostCollection(row)
	return nil
}

// Delete は指定IDのコレクションとその関連アイテムを削除する。
func (r *postCollectionRepository) Delete(ctx context.Context, id uint) error {
	if err := r.q.DeletePostCollectionItemsByCollectionID(ctx, int64(id)); err != nil {
		return err
	}
	return r.q.DeletePostCollection(ctx, int64(id))
}

// AddPost はコレクションに投稿を追加する。
func (r *postCollectionRepository) AddPost(ctx context.Context, item *model.PostCollectionItem) error {
	row, err := r.q.CreatePostCollectionItem(ctx, sqlcgen.CreatePostCollectionItemParams{
		CollectionID: int64(item.CollectionID),
		PostID:       int64(item.PostID),
		Note:         &item.Note,
		OrderIndex:   int64(item.OrderIndex),
	})
	if err != nil {
		return err
	}
	*item = model.PostCollectionItem{
		ID:           uint(row.ID),
		CollectionID: uint(row.CollectionID),
		PostID:       uint(row.PostID),
		Note:         fromStringPtr(row.Note),
		OrderIndex:   int(row.OrderIndex),
		CreatedAt:    timeValue(fromTimestamptz(row.CreatedAt)),
	}
	return nil
}

// HasPost は指定コレクションに指定投稿が存在するかを確認する。
func (r *postCollectionRepository) HasPost(ctx context.Context, collectionID, postID uint) (bool, error) {
	count, err := r.q.CountPostCollectionItemsByCollectionAndPost(ctx, sqlcgen.CountPostCollectionItemsByCollectionAndPostParams{
		CollectionID: int64(collectionID),
		PostID:       int64(postID),
	})
	return count > 0, err
}

// RemovePost はコレクションから投稿を削除する。
func (r *postCollectionRepository) RemovePost(ctx context.Context, collectionID, postID uint) error {
	return r.q.DeletePostCollectionItem(ctx, sqlcgen.DeletePostCollectionItemParams{
		CollectionID: int64(collectionID),
		PostID:       int64(postID),
	})
}

// GetPostsByCollectionID は指定コレクションの投稿一覧を順序付きで取得する。
func (r *postCollectionRepository) GetPostsByCollectionID(ctx context.Context, collectionID uint) ([]model.PostCollectionItem, error) {
	rows, err := r.q.ListPostCollectionItemsWithPostByCollectionID(ctx, int64(collectionID))
	if err != nil {
		return nil, err
	}
	items := make([]model.PostCollectionItem, len(rows))
	for i, row := range rows {
		post := toModelPost(row.Post)
		post.User = toModelUser(row.User)
		items[i] = model.PostCollectionItem{
			ID:           uint(row.PostCollectionItem.ID),
			CollectionID: uint(row.PostCollectionItem.CollectionID),
			PostID:       uint(row.PostCollectionItem.PostID),
			Post:         post,
			Note:         fromStringPtr(row.PostCollectionItem.Note),
			OrderIndex:   int(row.PostCollectionItem.OrderIndex),
			CreatedAt:    timeValue(fromTimestamptz(row.PostCollectionItem.CreatedAt)),
		}
	}
	return items, nil
}

// CountByUserID は指定ユーザーのコレクション総数を返す。
func (r *postCollectionRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountPostCollectionsByUser(ctx, int64(userID))
}

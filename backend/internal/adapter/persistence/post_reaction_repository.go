package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postReactionRepository は [repository.PostReactionRepository] と
// [repository.PostAuthorReader] の sqlc(pgx) 実装。
type postReactionRepository struct {
	q *sqlcgen.Queries
}

// NewPostReactionRepository は PostReactionRepository の sqlc(pgx) 実装を返す。
func NewPostReactionRepository(q *sqlcgen.Queries) repository.PostReactionRepository {
	return &postReactionRepository{q: q}
}

// NewPostAuthorReader は PostAuthorReader の sqlc(pgx) 実装を返す。
func NewPostAuthorReader(q *sqlcgen.Queries) repository.PostAuthorReader {
	return &postReactionRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var (
	_ repository.PostReactionRepository = (*postReactionRepository)(nil)
	_ repository.PostAuthorReader       = (*postReactionRepository)(nil)
)

// toUint64Slice は []uint を sqlc の ANY($1) パラメータ用に []int64 へ変換する。
func toInt64Slice(ids []uint) []int64 {
	result := make([]int64, len(ids))
	for i, id := range ids {
		result[i] = int64(id)
	}
	return result
}

// AddReaction は投稿にリアクションを追加する。
func (r *postReactionRepository) AddReaction(ctx context.Context, userID, postID uint, emoji string) error {
	return r.q.CreateReaction(ctx, sqlcgen.CreateReactionParams{
		UserID: int64(userID),
		PostID: int64(postID),
		Value:  emoji,
	})
}

// RemoveReaction は投稿のリアクションを削除する。
func (r *postReactionRepository) RemoveReaction(ctx context.Context, userID, postID uint, emoji string) error {
	return r.q.DeleteReaction(ctx, sqlcgen.DeleteReactionParams{
		UserID: int64(userID),
		PostID: int64(postID),
		Value:  emoji,
	})
}

// GetReactionsByPostID は指定投稿のリアクション集計を絵文字ごとに返す。
func (r *postReactionRepository) GetReactionsByPostID(ctx context.Context, postID uint) ([]model.ReactionCount, error) {
	rows, err := r.q.ListReactionCountsByPost(ctx, int64(postID))
	if err != nil {
		return nil, err
	}
	counts := make([]model.ReactionCount, len(rows))
	for i, row := range rows {
		counts[i] = model.ReactionCount{Emoji: row.Emoji, Count: int(row.Count)}
	}
	return counts, nil
}

// GetUserReactions は指定ユーザーが投稿に付けたリアクション絵文字一覧を返す。
func (r *postReactionRepository) GetUserReactions(ctx context.Context, userID, postID uint) ([]string, error) {
	return r.q.ListUserReactionEmojisByPost(ctx, sqlcgen.ListUserReactionEmojisByPostParams{
		UserID: int64(userID),
		PostID: int64(postID),
	})
}

// GetReactionsBatch は複数投稿のリアクション集計を一括取得する。
func (r *postReactionRepository) GetReactionsBatch(ctx context.Context, postIDs []uint) (map[uint][]model.ReactionCount, error) {
	rows, err := r.q.ListReactionCountsByPosts(ctx, toInt64Slice(postIDs))
	if err != nil {
		return nil, err
	}
	m := make(map[uint][]model.ReactionCount)
	for _, row := range rows {
		postID := uint(row.PostID)
		m[postID] = append(m[postID], model.ReactionCount{Emoji: row.Emoji, Count: int(row.Count)})
	}
	return m, nil
}

// GetUserReactionsBatch は複数投稿に対するユーザーのリアクションを一括取得する。
func (r *postReactionRepository) GetUserReactionsBatch(ctx context.Context, userID uint, postIDs []uint) (map[uint][]string, error) {
	rows, err := r.q.ListUserReactionsByPosts(ctx, sqlcgen.ListUserReactionsByPostsParams{
		UserID:  int64(userID),
		Column2: toInt64Slice(postIDs),
	})
	if err != nil {
		return nil, err
	}
	m := make(map[uint][]string)
	for _, row := range rows {
		postID := uint(row.PostID)
		m[postID] = append(m[postID], row.Emoji)
	}
	return m, nil
}

// FindAuthorID は指定投稿の投稿者 ID を返す。投稿が存在しない場合は (0, nil) を返す。
func (r *postReactionRepository) FindAuthorID(ctx context.Context, postID uint) (uint, error) {
	authorID, err := r.q.GetPostAuthorID(ctx, int64(postID))
	if err != nil {
		if isNoRows(err) {
			return 0, nil
		}
		return 0, err
	}
	return uint(authorID), nil
}

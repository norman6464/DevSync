package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postRepository は [repository.PostRepository] の sqlc(pgx) 実装。
// Delete は投稿を参照する行ごとの削除を1トランザクションで行うため、
// Queries だけでなくトランザクションを開始できる *pgxpool.Pool を直接保持する。
type postRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewPostRepository は PostRepository の sqlc(pgx) 実装を返す。
func NewPostRepository(pool *pgxpool.Pool) repository.PostRepository {
	return &postRepository{pool: pool, q: sqlcgen.New(pool)}
}

var _ repository.PostRepository = (*postRepository)(nil)

// attachCodeSnippetsToPosts は複数の投稿へ CodeSnippets をまとめて取得して付与する
// （GORMの Preload("CodeSnippets") に相当。1対多のため投稿IDのまとめ取りと
// Go側でのグルーピングで再現する）。
func attachCodeSnippetsToPosts(ctx context.Context, q *sqlcgen.Queries, posts []model.Post) error {
	if len(posts) == 0 {
		return nil
	}
	postIDs := make([]int64, len(posts))
	for i, post := range posts {
		postIDs[i] = int64(post.ID)
	}

	snippetRows, err := q.ListCodeSnippetsByPostIDs(ctx, postIDs)
	if err != nil {
		return err
	}
	snippetsByPostID := make(map[uint][]model.CodeSnippet)
	for _, row := range snippetRows {
		postID := uint(row.PostID)
		snippetsByPostID[postID] = append(snippetsByPostID[postID], toModelCodeSnippet(row))
	}
	for i := range posts {
		posts[i].CodeSnippets = snippetsByPostID[posts[i].ID]
	}
	return nil
}

// Create は投稿を作成する。
func (r *postRepository) Create(ctx context.Context, post *model.Post) error {
	row, err := r.q.CreatePost(ctx, sqlcgen.CreatePostParams{
		UserID:            int64(post.UserID),
		Title:             post.Title,
		Content:           post.Content,
		ImageUrls:         &post.ImageURLs,
		IsDraft:           &post.IsDraft,
		LikeCount:         toInt64Ptr(post.LikeCount),
		CommentCount:      toInt64Ptr(post.CommentCount),
		BookmarkCount:     toInt64Ptr(post.BookmarkCount),
		ViewCount:         toInt64Ptr(post.ViewCount),
		EstimatedReadTime: toInt64Ptr(post.EstimatedReadTime),
		ScheduledAt:       toTimestamptz(post.ScheduledAt),
	})
	if err != nil {
		return err
	}
	*post = toModelPost(row)
	return nil
}

// FindByID は投稿を投稿者・コードスニペット付きで取得する。存在しなければ (nil, nil) を返す。
func (r *postRepository) FindByID(ctx context.Context, id uint) (*model.Post, error) {
	row, err := r.q.GetPostWithUserByID(ctx, int64(id))
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	post := toModelPost(row.Post)
	post.User = toModelUser(row.User)

	posts := []model.Post{post}
	if err := attachCodeSnippetsToPosts(ctx, r.q, posts); err != nil {
		return nil, err
	}
	return &posts[0], nil
}

// Update は投稿を更新する（GORMのSave＝全カラム上書きに相当）。
func (r *postRepository) Update(ctx context.Context, post *model.Post) error {
	row, err := r.q.UpdatePost(ctx, sqlcgen.UpdatePostParams{
		ID:                int64(post.ID),
		Title:             post.Title,
		Content:           post.Content,
		ImageUrls:         &post.ImageURLs,
		IsDraft:           &post.IsDraft,
		LikeCount:         toInt64Ptr(post.LikeCount),
		CommentCount:      toInt64Ptr(post.CommentCount),
		BookmarkCount:     toInt64Ptr(post.BookmarkCount),
		ViewCount:         toInt64Ptr(post.ViewCount),
		EstimatedReadTime: toInt64Ptr(post.EstimatedReadTime),
		ScheduledAt:       toTimestamptz(post.ScheduledAt),
	})
	if err != nil {
		return err
	}
	*post = toModelPost(row)
	return nil
}

// Delete は投稿を、投稿を参照する行ごとトランザクション内で削除する。
// 参照する行（通知・ブックマーク・スニペット等）には外部キー制約があり、
// 先に消さないと投稿本体の削除が拒否される。途中で失敗しても何も消えない。
func (r *postRepository) Delete(ctx context.Context, id uint) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)

	// 先に投稿行を排他ロックする。FK 付きの参照行の挿入は親行への共有ロックを
	// 要求するため、ここで直列化され、掃除の最中に新しい参照行が入り込んで
	// 最後の本体削除が失敗することを防ぐ。既に無ければ何もしない（冪等）。
	if _, err := q.LockPostForDelete(ctx, int64(id)); err != nil {
		if isNoRows(err) {
			return nil
		}
		return err
	}

	postID := int64(id)

	// コメントに従属する行（コメントいいね・コメント由来のメンション）を先に消す
	if err := q.DeleteCommentLikesByPostComments(ctx, postID); err != nil {
		return err
	}
	if err := q.DeleteMentionsByPostComments(ctx, postID); err != nil {
		return err
	}
	if err := q.DeleteCommentsByPost(ctx, postID); err != nil {
		return err
	}

	// スニペットに従属する行を先に消す
	if err := q.DeleteSnippetCommentsByPostSnippets(ctx, postID); err != nil {
		return err
	}
	if err := q.DeleteCodeSnippetsByPost(ctx, postID); err != nil {
		return err
	}

	// 投稿を直接参照する行を消す
	if err := q.DeleteLikesByPost(ctx, postID); err != nil {
		return err
	}
	if err := q.DeleteReactionsByPost(ctx, postID); err != nil {
		return err
	}
	if err := q.DeleteBookmarksByPost(ctx, postID); err != nil {
		return err
	}
	if err := q.DeleteBookmarkCollectionItemsByPost(ctx, postID); err != nil {
		return err
	}
	if err := q.DeletePostSeriesItemsByPost(ctx, postID); err != nil {
		return err
	}
	if err := q.DeletePostCollectionItemsByPost(ctx, postID); err != nil {
		return err
	}
	if err := q.DeletePostTagsByPost(ctx, postID); err != nil {
		return err
	}
	if err := q.DeletePostPinsByPost(ctx, postID); err != nil {
		return err
	}
	if err := q.DeletePostViewsByPost(ctx, postID); err != nil {
		return err
	}
	if err := q.DeleteNotificationsByPost(ctx, &postID); err != nil {
		return err
	}
	if err := q.DeleteMentionsByPost(ctx, &postID); err != nil {
		return err
	}

	if err := q.DeletePost(ctx, postID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// FindAll は公開済み投稿をページネーション付きで取得する（新しい順）。
func (r *postRepository) FindAll(ctx context.Context, page, limit int) ([]model.Post, error) {
	offset := (page - 1) * limit
	rows, err := r.q.ListPublicPostsWithUser(ctx, sqlcgen.ListPublicPostsWithUserParams{
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, err
	}

	posts := make([]model.Post, len(rows))
	for i, row := range rows {
		posts[i] = toModelPost(row.Post)
		posts[i].User = toModelUser(row.User)
	}
	if err := attachCodeSnippetsToPosts(ctx, r.q, posts); err != nil {
		return nil, err
	}
	return posts, nil
}

// CountAll は公開済み投稿の総数を返す。
func (r *postRepository) CountAll(ctx context.Context) (int64, error) {
	return r.q.CountPublicPosts(ctx)
}

// FindByUserID は指定ユーザーの公開済み投稿をページネーション付きで取得する（新しい順）。
func (r *postRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Post, int64, error) {
	total, err := r.q.CountPublishedPostsByUser(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListPublishedPostsByUserWithUser(ctx, sqlcgen.ListPublishedPostsByUserWithUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	posts := make([]model.Post, len(rows))
	for i, row := range rows {
		posts[i] = toModelPost(row.Post)
		posts[i].User = toModelUser(row.User)
	}
	if err := attachCodeSnippetsToPosts(ctx, r.q, posts); err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

// FindDraftsByUserID は指定ユーザーの下書きを取得する（更新の新しい順）。
func (r *postRepository) FindDraftsByUserID(ctx context.Context, userID uint) ([]model.Post, error) {
	rows, err := r.q.ListDraftPostsByUserWithUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	posts := make([]model.Post, len(rows))
	for i, row := range rows {
		posts[i] = toModelPost(row.Post)
		posts[i].User = toModelUser(row.User)
	}
	if err := attachCodeSnippetsToPosts(ctx, r.q, posts); err != nil {
		return nil, err
	}
	return posts, nil
}

// FindScheduledByUserID は指定ユーザーの公開予約済み投稿を取得する（公開予定日時順）。
func (r *postRepository) FindScheduledByUserID(ctx context.Context, userID uint) ([]model.Post, error) {
	rows, err := r.q.ListScheduledPostsByUserWithUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	posts := make([]model.Post, len(rows))
	for i, row := range rows {
		posts[i] = toModelPost(row.Post)
		posts[i].User = toModelUser(row.User)
	}
	if err := attachCodeSnippetsToPosts(ctx, r.q, posts); err != nil {
		return nil, err
	}
	return posts, nil
}

// Timeline はフォロー中ユーザーと自分の公開済み投稿を取得する（新しい順）。
func (r *postRepository) Timeline(ctx context.Context, userID uint, page, limit int) ([]model.Post, error) {
	offset := (page - 1) * limit
	rows, err := r.q.ListTimelinePostsWithUser(ctx, sqlcgen.ListTimelinePostsWithUserParams{
		FollowerID: int64(userID),
		Limit:      int32Param(limit),
		Offset:     int32Param(offset),
	})
	if err != nil {
		return nil, err
	}

	posts := make([]model.Post, len(rows))
	for i, row := range rows {
		posts[i] = toModelPost(row.Post)
		posts[i].User = toModelUser(row.User)
	}
	if err := attachCodeSnippetsToPosts(ctx, r.q, posts); err != nil {
		return nil, err
	}
	return posts, nil
}

// CountByUserID は指定ユーザーの公開済み投稿数を返す。
func (r *postRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountPublishedPostsByUser(ctx, int64(userID))
}

// CountDraftsByUserID は指定ユーザーの下書き数を返す。
func (r *postRepository) CountDraftsByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountDraftPostsByUser(ctx, int64(userID))
}

// CountScheduledByUserID は指定ユーザーの公開予約済み投稿数を返す。
func (r *postRepository) CountScheduledByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountScheduledPostsByUser(ctx, int64(userID))
}

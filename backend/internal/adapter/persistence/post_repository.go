package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postRepository は [repository.PostRepository] の sqlc(pgx) 実装。
type postRepository struct {
	q *sqlcgen.Queries
}

// NewPostRepository は PostRepository の sqlc(pgx) 実装を返す。
func NewPostRepository(pool *pgxpool.Pool) repository.PostRepository {
	return &postRepository{q: sqlcgen.New(pool)}
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

// attachBookmarkCountsToPosts は複数の投稿へブックマーク数をまとめて取得して付与する。
// bookmark_countはORDER BYに使われない表示専用の値のため列として持たず、都度
// post_reactions（kind='bookmark'）からCOUNT(*)する
// （queries/post_bookmark.sqlのCountBookmarksByPostIDs参照）。
func attachBookmarkCountsToPosts(ctx context.Context, q *sqlcgen.Queries, posts []model.Post) error {
	if len(posts) == 0 {
		return nil
	}
	postIDs := make([]int64, len(posts))
	for i, post := range posts {
		postIDs[i] = int64(post.ID)
	}

	countRows, err := q.CountBookmarksByPostIDs(ctx, postIDs)
	if err != nil {
		return err
	}
	countByPostID := make(map[uint]int, len(countRows))
	for _, row := range countRows {
		countByPostID[uint(row.PostID)] = int(row.BookmarkCount)
	}
	for i := range posts {
		posts[i].BookmarkCount = countByPostID[posts[i].ID]
	}
	return nil
}

// attachMetricsToPosts は複数の投稿へlike_count/comment_count/view_countをまとめて
// 取得して付与する（DEVSYNC-159でpost_metrics側テーブルへ分離済み。1件もいいね/
// コメント/閲覧が無い投稿はpost_metrics行が存在しないため0のまま）。
func attachMetricsToPosts(ctx context.Context, q *sqlcgen.Queries, posts []model.Post) error {
	if len(posts) == 0 {
		return nil
	}
	postIDs := make([]int64, len(posts))
	for i, post := range posts {
		postIDs[i] = int64(post.ID)
	}

	metricsRows, err := q.GetPostMetricsByPostIDs(ctx, postIDs)
	if err != nil {
		return err
	}
	metricsByPostID := make(map[uint]sqlcgen.PostMetric, len(metricsRows))
	for _, row := range metricsRows {
		metricsByPostID[uint(row.PostID)] = row
	}
	for i := range posts {
		m := metricsByPostID[posts[i].ID]
		posts[i].LikeCount = int(m.LikeCount)
		posts[i].CommentCount = int(m.CommentCount)
		posts[i].ViewCount = int(m.ViewCount)
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
		IsDraft:           post.IsDraft,
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
	if err := attachBookmarkCountsToPosts(ctx, r.q, posts); err != nil {
		return nil, err
	}
	if err := attachMetricsToPosts(ctx, r.q, posts); err != nil {
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
		IsDraft:           post.IsDraft,
		EstimatedReadTime: toInt64Ptr(post.EstimatedReadTime),
		ScheduledAt:       toTimestamptz(post.ScheduledAt),
	})
	if err != nil {
		return err
	}
	*post = toModelPost(row)
	posts := []model.Post{*post}
	if err := attachBookmarkCountsToPosts(ctx, r.q, posts); err != nil {
		return err
	}
	if err := attachMetricsToPosts(ctx, r.q, posts); err != nil {
		return err
	}
	*post = posts[0]
	return nil
}

// Delete は投稿を削除する。投稿を参照する行（コメント・いいね・通知等）の削除は
// FKのON DELETE CASCADE宣言（internal/infra/database/schema/schema.hcl）に委ねる
// （DEVSYNC-156でFKを全テーブルへ投入済み）。かつては参照テーブルを1つずつ
// 手書きで削除する16手順のトランザクションだったが、新しい参照テーブルを
// 追加するたびにここへの追記が必要という壊れやすい運用だった。
func (r *postRepository) Delete(ctx context.Context, id uint) error {
	return r.q.DeletePost(ctx, int64(id))
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
	if err := attachBookmarkCountsToPosts(ctx, r.q, posts); err != nil {
		return nil, err
	}
	if err := attachMetricsToPosts(ctx, r.q, posts); err != nil {
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
	if err := attachBookmarkCountsToPosts(ctx, r.q, posts); err != nil {
		return nil, 0, err
	}
	if err := attachMetricsToPosts(ctx, r.q, posts); err != nil {
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
	if err := attachBookmarkCountsToPosts(ctx, r.q, posts); err != nil {
		return nil, err
	}
	if err := attachMetricsToPosts(ctx, r.q, posts); err != nil {
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
	if err := attachBookmarkCountsToPosts(ctx, r.q, posts); err != nil {
		return nil, err
	}
	if err := attachMetricsToPosts(ctx, r.q, posts); err != nil {
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
	if err := attachBookmarkCountsToPosts(ctx, r.q, posts); err != nil {
		return nil, err
	}
	if err := attachMetricsToPosts(ctx, r.q, posts); err != nil {
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

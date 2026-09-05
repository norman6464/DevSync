package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postCommentRepository は [repository.PostCommentRepository] の sqlc(pgx) 実装。
// 単体取得は sqlc(pgx) 実装の commentReader を埋め込んで再利用する。
type postCommentRepository struct {
	repository.CommentReader
	q *sqlcgen.Queries
}

// NewPostCommentRepository は PostCommentRepository の sqlc(pgx) 実装を返す。
func NewPostCommentRepository(q *sqlcgen.Queries) repository.PostCommentRepository {
	return &postCommentRepository{CommentReader: NewCommentReader(q), q: q}
}

var _ repository.PostCommentRepository = (*postCommentRepository)(nil)

// Create はコメントを作成し、投稿の comment_count を加算する。
func (r *postCommentRepository) Create(ctx context.Context, comment *model.Comment) error {
	row, err := r.q.CreatePostComment(ctx, sqlcgen.CreatePostCommentParams{
		UserID:   int64(comment.UserID),
		PostID:   int64(comment.PostID),
		ParentID: toInt64PtrFromUintPtr(comment.ParentID),
		Content:  comment.Content,
	})
	if err != nil {
		return err
	}
	*comment = toModelComment(row)
	return r.q.IncrementPostCommentCount(ctx, int64(comment.PostID))
}

// Update はコメントを更新する。
func (r *postCommentRepository) Update(ctx context.Context, comment *model.Comment) error {
	isHidden := comment.IsHidden
	row, err := r.q.UpdatePostComment(ctx, sqlcgen.UpdatePostCommentParams{
		ID:       int64(comment.ID),
		Content:  comment.Content,
		IsHidden: &isHidden,
	})
	if err != nil {
		return err
	}
	*comment = toModelComment(row)
	return nil
}

// Delete はコメントを削除し、投稿の comment_count をデクリメントする。
// 所有権チェックは usecase 層で実施済みであること。
func (r *postCommentRepository) Delete(ctx context.Context, id uint) error {
	row, err := r.q.GetCommentByID(ctx, int64(id))
	if err != nil {
		return err
	}
	comment := toModelComment(row)
	// 移行前の GORM 実装と同じく、カウンタ減算自体のエラーは呼び出し元へ返さない。
	_ = r.q.DecrementPostCommentCount(ctx, int64(comment.PostID))
	return r.q.DeletePostComment(ctx, int64(id))
}

// ListByPostID は指定投稿のトップレベルコメントをユーザー情報・返信付きで取得する（古い順）。
func (r *postCommentRepository) ListByPostID(ctx context.Context, postID uint) ([]model.Comment, error) {
	rows, err := r.q.ListTopLevelCommentsByPost(ctx, int64(postID))
	if err != nil {
		return nil, err
	}
	comments := make([]model.Comment, len(rows))
	parentIDs := make([]int64, len(rows))
	for i, row := range rows {
		comments[i] = toModelComment(row.Comment)
		comments[i].User = toModelUser(row.User)
		parentIDs[i] = row.Comment.ID
	}

	if len(parentIDs) > 0 {
		replyRows, err := r.q.ListCommentRepliesByParentIDs(ctx, parentIDs)
		if err != nil {
			return nil, err
		}
		repliesByParentID := make(map[uint][]model.Comment)
		for _, row := range replyRows {
			reply := toModelComment(row.Comment)
			reply.User = toModelUser(row.User)
			parentID := *reply.ParentID
			repliesByParentID[parentID] = append(repliesByParentID[parentID], reply)
		}
		for i := range comments {
			comments[i].Replies = repliesByParentID[comments[i].ID]
		}
	}

	return comments, nil
}

// ListReplies は指定コメントへの返信をユーザー情報付きで取得する（古い順）。
func (r *postCommentRepository) ListReplies(ctx context.Context, parentID uint) ([]model.Comment, error) {
	rows, err := r.q.ListCommentRepliesByParentIDs(ctx, []int64{int64(parentID)})
	if err != nil {
		return nil, err
	}
	replies := make([]model.Comment, len(rows))
	for i, row := range rows {
		replies[i] = toModelComment(row.Comment)
		replies[i].User = toModelUser(row.User)
	}
	return replies, nil
}

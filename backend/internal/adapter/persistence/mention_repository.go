package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// mentionRepository は [repository.MentionRepository] の sqlc(pgx) 実装。
type mentionRepository struct {
	q *sqlcgen.Queries
}

// NewMentionRepository は MentionRepository の sqlc(pgx) 実装を返す。
func NewMentionRepository(q *sqlcgen.Queries) repository.MentionRepository {
	return &mentionRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.MentionRepository = (*mentionRepository)(nil)

// toModelMention は sqlc の生成行を model.Mention へ変換する（関連の User/Actor は含まない）。
func toModelMention(row sqlcgen.Mention) model.Mention {
	return model.Mention{
		ID:        uint(row.ID),
		UserID:    uint(row.UserID),
		ActorID:   uint(row.ActorID),
		PostID:    fromInt64PtrToUintPtr(row.PostID),
		CommentID: fromInt64PtrToUintPtr(row.CommentID),
		CreatedAt: timeValue(fromTimestamptz(row.CreatedAt)),
	}
}

// Create はメンションを保存する。同じ相手が同じ投稿・コメントで既にメンションされている
// 場合は何もせず（ON CONFLICT DO NOTHING）、作成したかどうかを返す。
func (r *mentionRepository) Create(ctx context.Context, mention *model.Mention) (bool, error) {
	row, err := r.q.CreateMention(ctx, sqlcgen.CreateMentionParams{
		UserID:    int64(mention.UserID),
		ActorID:   int64(mention.ActorID),
		PostID:    toInt64PtrFromUintPtr(mention.PostID),
		CommentID: toInt64PtrFromUintPtr(mention.CommentID),
	})
	if err != nil {
		if isNoRows(err) {
			return false, nil
		}
		return false, err
	}
	*mention = toModelMention(row)
	return true, nil
}

// FindByUserID は指定ユーザー宛のメンションを作成日時の降順で取得する。
func (r *mentionRepository) FindByUserID(ctx context.Context, userID uint, page, limit int) ([]model.Mention, error) {
	offset := (page - 1) * limit
	rows, err := r.q.ListMentionsByUser(ctx, sqlcgen.ListMentionsByUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, err
	}
	mentions := make([]model.Mention, len(rows))
	for i, row := range rows {
		mentions[i] = toModelMention(row.Mention)
		mentions[i].Actor = toModelUser(row.User)
	}
	return mentions, nil
}

// FindByPostID は指定投稿に紐づくメンションを取得する。
func (r *mentionRepository) FindByPostID(ctx context.Context, postID uint) ([]model.Mention, error) {
	id := int64(postID)
	rows, err := r.q.ListMentionsByPostID(ctx, &id)
	if err != nil {
		return nil, err
	}
	mentions := make([]model.Mention, len(rows))
	for i, row := range rows {
		mentions[i] = toModelMention(row.Mention)
		mentions[i].User = toModelUser(row.User)
		mentions[i].Actor = toModelUser(row.User_2)
	}
	return mentions, nil
}

// FindByCommentID は指定コメントに紐づくメンションを取得する。
func (r *mentionRepository) FindByCommentID(ctx context.Context, commentID uint) ([]model.Mention, error) {
	id := int64(commentID)
	rows, err := r.q.ListMentionsByCommentID(ctx, &id)
	if err != nil {
		return nil, err
	}
	mentions := make([]model.Mention, len(rows))
	for i, row := range rows {
		mentions[i] = toModelMention(row.Mention)
		mentions[i].User = toModelUser(row.User)
		mentions[i].Actor = toModelUser(row.User_2)
	}
	return mentions, nil
}

// DeleteByPostID は指定投稿に紐づくメンションをすべて削除する。
func (r *mentionRepository) DeleteByPostID(ctx context.Context, postID uint) error {
	id := int64(postID)
	return r.q.DeleteMentionsByPostID(ctx, &id)
}

// DeleteByCommentID は指定コメントに紐づくメンションをすべて削除する。
func (r *mentionRepository) DeleteByCommentID(ctx context.Context, commentID uint) error {
	id := int64(commentID)
	return r.q.DeleteMentionsByCommentID(ctx, &id)
}

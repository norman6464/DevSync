package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// postReactionRepository は [repository.PostReactionRepository] と
// [repository.PostAuthorReader] の GORM 実装。
type postReactionRepository struct {
	db *gorm.DB
}

// NewPostReactionRepository は PostReactionRepository の GORM 実装を返す。
func NewPostReactionRepository(db *gorm.DB) repository.PostReactionRepository {
	return &postReactionRepository{db: db}
}

// NewPostAuthorReader は PostAuthorReader の GORM 実装を返す。
func NewPostAuthorReader(db *gorm.DB) repository.PostAuthorReader {
	return &postReactionRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.PostReactionRepository = (*postReactionRepository)(nil)
var _ repository.PostAuthorReader = (*postReactionRepository)(nil)

// AddReaction は投稿にリアクションを追加する。
func (r *postReactionRepository) AddReaction(ctx context.Context, userID, postID uint, emoji string) error {
	return r.db.WithContext(ctx).
		Create(&model.Reaction{UserID: userID, PostID: postID, Emoji: emoji}).Error
}

// RemoveReaction は投稿のリアクションを削除する。
func (r *postReactionRepository) RemoveReaction(ctx context.Context, userID, postID uint, emoji string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND post_id = ? AND emoji = ?", userID, postID, emoji).
		Delete(&model.Reaction{}).Error
}

// GetReactionsByPostID は指定投稿のリアクション集計を絵文字ごとに返す。
func (r *postReactionRepository) GetReactionsByPostID(ctx context.Context, postID uint) ([]model.ReactionCount, error) {
	var counts []model.ReactionCount
	err := r.db.WithContext(ctx).Model(&model.Reaction{}).
		Select("emoji, COUNT(*) as count").
		Where("post_id = ?", postID).
		Group("emoji").
		Order("count DESC").
		Find(&counts).Error
	return counts, err
}

// GetUserReactions は指定ユーザーが投稿に付けたリアクション絵文字一覧を返す。
func (r *postReactionRepository) GetUserReactions(ctx context.Context, userID, postID uint) ([]string, error) {
	var emojis []string
	err := r.db.WithContext(ctx).Model(&model.Reaction{}).
		Select("emoji").
		Where("user_id = ? AND post_id = ?", userID, postID).
		Pluck("emoji", &emojis).Error
	return emojis, err
}

// GetReactionsBatch は複数投稿のリアクション集計を一括取得する。
func (r *postReactionRepository) GetReactionsBatch(ctx context.Context, postIDs []uint) (map[uint][]model.ReactionCount, error) {
	type row struct {
		PostID uint   `gorm:"column:post_id"`
		Emoji  string `gorm:"column:emoji"`
		Count  int    `gorm:"column:count"`
	}
	var rows []row
	err := r.db.WithContext(ctx).Model(&model.Reaction{}).
		Select("post_id, emoji, COUNT(*) as count").
		Where("post_id IN ?", postIDs).
		Group("post_id, emoji").
		Order("post_id, count DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	m := make(map[uint][]model.ReactionCount)
	for _, res := range rows {
		m[res.PostID] = append(m[res.PostID], model.ReactionCount{Emoji: res.Emoji, Count: res.Count})
	}
	return m, nil
}

// GetUserReactionsBatch は複数投稿に対するユーザーのリアクションを一括取得する。
func (r *postReactionRepository) GetUserReactionsBatch(ctx context.Context, userID uint, postIDs []uint) (map[uint][]string, error) {
	type row struct {
		PostID uint   `gorm:"column:post_id"`
		Emoji  string `gorm:"column:emoji"`
	}
	var rows []row
	err := r.db.WithContext(ctx).Model(&model.Reaction{}).
		Select("post_id, emoji").
		Where("user_id = ? AND post_id IN ?", userID, postIDs).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	m := make(map[uint][]string)
	for _, res := range rows {
		m[res.PostID] = append(m[res.PostID], res.Emoji)
	}
	return m, nil
}

// FindAuthorID は指定投稿の投稿者 ID を返す。投稿が存在しない場合は (0, nil) を返す。
func (r *postReactionRepository) FindAuthorID(ctx context.Context, postID uint) (uint, error) {
	var post model.Post
	err := r.db.WithContext(ctx).Select("id", "user_id").First(&post, postID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return post.UserID, nil
}

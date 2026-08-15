package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// postRepository は [repository.PostRepository] の GORM 実装。
type postRepository struct {
	db *gorm.DB
}

// NewPostRepository は PostRepository の GORM 実装を返す。
func NewPostRepository(db *gorm.DB) repository.PostRepository {
	return &postRepository{db: db}
}

var _ repository.PostRepository = (*postRepository)(nil)

// Create は投稿を作成する。
func (r *postRepository) Create(ctx context.Context, post *model.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

// FindByID は投稿を投稿者・コードスニペット付きで取得する。存在しなければ (nil, nil) を返す。
func (r *postRepository) FindByID(ctx context.Context, id uint) (*model.Post, error) {
	var post model.Post
	if err := r.db.WithContext(ctx).Preload("User").Preload("CodeSnippets").First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

// Update は投稿を更新する。
func (r *postRepository) Update(ctx context.Context, post *model.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

// Delete は投稿を削除する。
// Delete は投稿を、投稿を参照する行ごとトランザクション内で削除する。
// 参照する行（通知・ブックマーク・スニペット等）には外部キー制約があり、
// 先に消さないと投稿本体の削除が拒否される。途中で失敗しても何も消えない。
func (r *postRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// コメントに従属する行（コメントいいね・コメント由来のメンション）を先に消す
		commentIDs := tx.Model(&model.Comment{}).Select("id").Where("post_id = ?", id)
		if err := tx.Where("comment_id IN (?)", commentIDs).Delete(&model.CommentLike{}).Error; err != nil {
			return err
		}
		if err := tx.Where("comment_id IN (?)", commentIDs).Delete(&model.Mention{}).Error; err != nil {
			return err
		}
		if err := tx.Where("post_id = ?", id).Delete(&model.Comment{}).Error; err != nil {
			return err
		}

		// スニペットに従属する行を先に消す
		snippetIDs := tx.Model(&model.CodeSnippet{}).Select("id").Where("post_id = ?", id)
		if err := tx.Where("snippet_id IN (?)", snippetIDs).Delete(&model.SnippetComment{}).Error; err != nil {
			return err
		}
		if err := tx.Where("post_id = ?", id).Delete(&model.CodeSnippet{}).Error; err != nil {
			return err
		}

		// 投稿を直接参照する行を消す
		for _, target := range []interface{}{
			&model.Like{}, &model.Reaction{}, &model.Bookmark{},
			&model.BookmarkCollectionItem{}, &model.PostSeriesItem{}, &model.PostCollectionItem{},
			&model.PostTag{}, &model.PostPin{}, &model.PostView{},
			&model.Notification{}, &model.Mention{},
		} {
			if err := tx.Where("post_id = ?", id).Delete(target).Error; err != nil {
				return err
			}
		}

		return tx.Delete(&model.Post{}, id).Error
	})
}

// FindAll は公開済み投稿をページネーション付きで取得する（新しい順）。
func (r *postRepository) FindAll(ctx context.Context, page, limit int) ([]model.Post, error) {
	var posts []model.Post
	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).Preload("User").Preload("CodeSnippets").
		Where("is_draft = ?", false).
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error
	return posts, err
}

// CountAll は公開済み投稿の総数を返す。
func (r *postRepository) CountAll(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Post{}).Where("is_draft = ?", false).Count(&count).Error
	return count, err
}

// FindByUserID は指定ユーザーの公開済み投稿をページネーション付きで取得する（新しい順）。
func (r *postRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64
	query := r.db.WithContext(ctx).Where("user_id = ? AND is_draft = ?", userID, false)
	query.Model(&model.Post{}).Count(&total)
	err := query.Preload("User").Preload("CodeSnippets").
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&posts).Error
	return posts, total, err
}

// FindDraftsByUserID は指定ユーザーの下書きを取得する（更新の新しい順）。
func (r *postRepository) FindDraftsByUserID(ctx context.Context, userID uint) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.WithContext(ctx).Preload("User").Preload("CodeSnippets").
		Where("user_id = ? AND is_draft = ?", userID, true).
		Order("updated_at DESC").Find(&posts).Error
	return posts, err
}

// FindScheduledByUserID は指定ユーザーの公開予約済み投稿を取得する（公開予定日時順）。
func (r *postRepository) FindScheduledByUserID(ctx context.Context, userID uint) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.WithContext(ctx).Preload("User").Preload("CodeSnippets").
		Where("user_id = ? AND is_draft = ? AND scheduled_at IS NOT NULL", userID, true).
		Order("scheduled_at ASC").Find(&posts).Error
	return posts, err
}

// Timeline はフォロー中ユーザーと自分の公開済み投稿を取得する（新しい順）。
func (r *postRepository) Timeline(ctx context.Context, userID uint, page, limit int) ([]model.Post, error) {
	var posts []model.Post
	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).Preload("User").Preload("CodeSnippets").
		Where("(user_id IN (SELECT followee_id FROM follows WHERE follower_id = ?) OR user_id = ?) AND is_draft = ?", userID, userID, false).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&posts).Error
	return posts, err
}

// CountByUserID は指定ユーザーの公開済み投稿数を返す。
func (r *postRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Post{}).
		Where("user_id = ? AND is_draft = ?", userID, false).Count(&count).Error
	return count, err
}

// CountDraftsByUserID は指定ユーザーの下書き数を返す。
func (r *postRepository) CountDraftsByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Post{}).
		Where("user_id = ? AND is_draft = ?", userID, true).Count(&count).Error
	return count, err
}

// CountScheduledByUserID は指定ユーザーの公開予約済み投稿数を返す。
func (r *postRepository) CountScheduledByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Post{}).
		Where("user_id = ? AND is_draft = ? AND scheduled_at IS NOT NULL", userID, true).Count(&count).Error
	return count, err
}

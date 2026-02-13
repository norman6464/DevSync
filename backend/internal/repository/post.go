package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// PostRepository は投稿データへのアクセスを提供するリポジトリ実装。
type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository は新しいPostRepositoryインスタンスを生成する。
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Create は新しい投稿をデータベースに作成する。
func (r *PostRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

// FindByID は指定IDの投稿をユーザー情報付きで取得する。
func (r *PostRepository) FindByID(id uint) (*model.Post, error) {
	var post model.Post
	err := r.db.Preload("User").Preload("CodeSnippets").First(&post, id).Error
	return &post, err
}

// FindAll はページネーション付きで全投稿を取得する（新しい順）。下書きは除外。
func (r *PostRepository) FindAll(page, limit int) ([]model.Post, error) {
	var posts []model.Post
	offset := (page - 1) * limit
	err := r.db.Preload("User").Preload("CodeSnippets").
		Where("is_draft = ?", false).
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error
	return posts, err
}

// FindByUserID は指定ユーザーの全投稿を取得する（新しい順）。下書きは除外。
func (r *PostRepository) FindByUserID(userID uint) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Preload("User").Preload("CodeSnippets").
		Where("user_id = ? AND is_draft = ?", userID, false).
		Order("created_at DESC").Find(&posts).Error
	return posts, err
}

// Timeline はフォロー中ユーザーと自分の投稿をタイムライン形式で取得する。下書きは除外。
// サブクエリでフォロー中のユーザーIDを取得し、自分のIDと合わせてフィルタする。
func (r *PostRepository) Timeline(userID uint, page, limit int) ([]model.Post, error) {
	var posts []model.Post
	offset := (page - 1) * limit
	err := r.db.Preload("User").Preload("CodeSnippets").
		Where("(user_id IN (SELECT followee_id FROM follows WHERE follower_id = ?) OR user_id = ?) AND is_draft = ?", userID, userID, false).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&posts).Error
	return posts, err
}

// Update は既存の投稿情報を更新する。
func (r *PostRepository) Update(post *model.Post) error {
	return r.db.Save(post).Error
}

// Delete は指定IDの投稿を削除する。
func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&model.Post{}, id).Error
}

// Like は投稿にいいねを追加し、投稿のlike_countをインクリメントする。
func (r *PostRepository) Like(userID, postID uint) error {
	err := r.db.Create(&model.Like{UserID: userID, PostID: postID}).Error
	if err != nil {
		return err
	}
	return r.db.Model(&model.Post{}).Where("id = ?", postID).UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

// Unlike は投稿のいいねを取り消し、投稿のlike_countをデクリメントする。
func (r *PostRepository) Unlike(userID, postID uint) error {
	result := r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&model.Like{})
	if result.RowsAffected > 0 {
		r.db.Model(&model.Post{}).Where("id = ?", postID).UpdateColumn("like_count", gorm.Expr("GREATEST(like_count - 1, 0)"))
	}
	return result.Error
}

// HasLiked は指定ユーザーが投稿にいいね済みかどうかを判定する。
func (r *PostRepository) HasLiked(userID, postID uint) bool {
	var count int64
	r.db.Model(&model.Like{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count)
	return count > 0
}

// CreateComment は投稿にコメントを追加し、投稿のcomment_countをインクリメントする。
func (r *PostRepository) CreateComment(comment *model.Comment) error {
	err := r.db.Create(comment).Error
	if err != nil {
		return err
	}
	return r.db.Model(&model.Post{}).Where("id = ?", comment.PostID).UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error
}

// GetComments は指定投稿の全コメントをユーザー情報付きで取得する（古い順）。
func (r *PostRepository) GetComments(postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Preload("User").Where("post_id = ?", postID).Order("created_at ASC").Find(&comments).Error
	return comments, err
}

// DeleteComment はコメントを削除し、投稿のcomment_countをデクリメントする。
// userIDによる所有権チェックを行う。
func (r *PostRepository) DeleteComment(id, userID uint) error {
	var comment model.Comment
	if err := r.db.First(&comment, id).Error; err != nil {
		return err
	}
	if comment.UserID != userID {
		return gorm.ErrRecordNotFound
	}
	r.db.Model(&model.Post{}).Where("id = ?", comment.PostID).UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)"))
	return r.db.Delete(&comment).Error
}

// Search はキーワードで投稿を検索する（タイトルまたは本文に部分一致）。下書きは除外。
func (r *PostRepository) Search(query string, limit, offset int) (interface{}, int64, error) {
	var posts []model.Post
	var total int64

	searchPattern := "%" + query + "%"
	db := r.db.Preload("User").Preload("CodeSnippets").
		Where("(title LIKE ? OR content LIKE ?) AND is_draft = ?", searchPattern, searchPattern, false).
		Order("created_at DESC")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(offset).Limit(limit).Find(&posts).Error
	return posts, total, err
}

// FindDraftsByUserID は指定ユーザーの下書き一覧を取得する（新しい順）。
func (r *PostRepository) FindDraftsByUserID(userID uint) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Preload("User").Preload("CodeSnippets").
		Where("user_id = ? AND is_draft = ?", userID, true).
		Order("updated_at DESC").Find(&posts).Error
	return posts, err
}

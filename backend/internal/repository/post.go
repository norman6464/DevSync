package repository

import (
	"time"

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

// CountAll は公開済み投稿の総数を取得する（下書きを除く）。
func (r *PostRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&model.Post{}).Where("is_draft = ?", false).Count(&count).Error
	return count, err
}

// FindByUserID は指定ユーザーの投稿をページネーション付きで取得する（新しい順）。下書きは除外。
func (r *PostRepository) FindByUserID(userID uint, limit, offset int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64
	query := r.db.Where("user_id = ? AND is_draft = ?", userID, false)
	query.Model(&model.Post{}).Count(&total)
	err := query.Preload("User").Preload("CodeSnippets").
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&posts).Error
	return posts, total, err
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

// FindCommentByID はコメントをIDで取得する。
func (r *PostRepository) FindCommentByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	if err := r.db.First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetComments は指定投稿の全コメントをユーザー情報付きで取得する（古い順）。
func (r *PostRepository) GetComments(postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := r.db.Preload("User").Preload("Replies").Preload("Replies.User").
		Where("post_id = ? AND parent_id IS NULL", postID).
		Order("created_at ASC").Find(&comments).Error
	return comments, err
}

func (r *PostRepository) GetReplies(parentID uint) ([]model.Comment, error) {
	var replies []model.Comment
	err := r.db.Preload("User").Where("parent_id = ?", parentID).Order("created_at ASC").Find(&replies).Error
	return replies, err
}

// DeleteComment はコメントを削除し、投稿のcomment_countをデクリメントする。
// 所有権チェックはservice層で実施済みであること。
func (r *PostRepository) DeleteComment(id uint) error {
	var comment model.Comment
	if err := r.db.First(&comment, id).Error; err != nil {
		return err
	}
	r.db.Model(&model.Post{}).Where("id = ?", comment.PostID).UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)"))
	return r.db.Delete(&comment).Error
}

// Search はキーワードで投稿を検索する（タイトルまたは本文に部分一致）。下書きは除外。
func (r *PostRepository) Search(query string, limit, offset int) (interface{}, int64, error) {
	var posts []model.Post
	var total int64

	searchPattern := EscapeLikePattern(query)
	db := r.db.Preload("User").Preload("CodeSnippets").
		Where("(title LIKE ? OR content LIKE ?) AND is_draft = ?", searchPattern, searchPattern, false).
		Order("created_at DESC")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(offset).Limit(limit).Find(&posts).Error
	return posts, total, err
}

// Bookmark は投稿をブックマークする。
func (r *PostRepository) Bookmark(userID, postID uint) error {
	err := r.db.Create(&model.Bookmark{UserID: userID, PostID: postID}).Error
	if err != nil {
		return err
	}
	return r.db.Model(&model.Post{}).Where("id = ?", postID).UpdateColumn("bookmark_count", gorm.Expr("bookmark_count + 1")).Error
}

// Unbookmark は投稿のブックマークを解除し、bookmark_countをデクリメントする。
func (r *PostRepository) Unbookmark(userID, postID uint) error {
	result := r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&model.Bookmark{})
	if result.RowsAffected > 0 {
		r.db.Model(&model.Post{}).Where("id = ?", postID).UpdateColumn("bookmark_count", gorm.Expr("GREATEST(bookmark_count - 1, 0)"))
	}
	return result.Error
}

// HasBookmarked は指定ユーザーが投稿をブックマーク済みかどうかを判定する。
func (r *PostRepository) HasBookmarked(userID, postID uint) bool {
	var count int64
	r.db.Model(&model.Bookmark{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count)
	return count > 0
}

// FindBookmarkedByUserID は指定ユーザーのブックマーク済み投稿をページネーション付きで取得する。
func (r *PostRepository) FindBookmarkedByUserID(userID uint, page, limit int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64
	offset := (page - 1) * limit

	subQuery := r.db.Model(&model.Bookmark{}).Select("post_id").Where("user_id = ?", userID)

	if err := r.db.Model(&model.Post{}).Where("id IN (?)", subQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("User").Preload("CodeSnippets").
		Where("id IN (?)", subQuery).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&posts).Error
	return posts, total, err
}

// AddReaction は投稿にリアクション（絵文字）を追加する。
func (r *PostRepository) AddReaction(userID, postID uint, emoji string) error {
	return r.db.Create(&model.Reaction{UserID: userID, PostID: postID, Emoji: emoji}).Error
}

// RemoveReaction は投稿のリアクションを削除する。
func (r *PostRepository) RemoveReaction(userID, postID uint, emoji string) error {
	return r.db.Where("user_id = ? AND post_id = ? AND emoji = ?", userID, postID, emoji).Delete(&model.Reaction{}).Error
}

// GetReactionsByPostID は指定投稿のリアクション集計を絵文字ごとに返す。
func (r *PostRepository) GetReactionsByPostID(postID uint) ([]model.ReactionCount, error) {
	var counts []model.ReactionCount
	err := r.db.Model(&model.Reaction{}).
		Select("emoji, COUNT(*) as count").
		Where("post_id = ?", postID).
		Group("emoji").
		Order("count DESC").
		Find(&counts).Error
	return counts, err
}

// GetUserReactions は指定ユーザーが投稿に付けたリアクション絵文字一覧を返す。
func (r *PostRepository) GetUserReactions(userID, postID uint) ([]string, error) {
	var emojis []string
	err := r.db.Model(&model.Reaction{}).
		Select("emoji").
		Where("user_id = ? AND post_id = ?", userID, postID).
		Pluck("emoji", &emojis).Error
	return emojis, err
}

// FindDraftsByUserID は指定ユーザーの下書き一覧を取得する（新しい順）。
func (r *PostRepository) FindDraftsByUserID(userID uint) ([]model.Post, error) {
	var posts []model.Post
	err := r.db.Preload("User").Preload("CodeSnippets").
		Where("user_id = ? AND is_draft = ?", userID, true).
		Order("updated_at DESC").Find(&posts).Error
	return posts, err
}

// SearchWithFilter はタグ・日付範囲・ソート順による高度な投稿検索を実行する。
// 下書きは検索対象外。タグはAND条件で絞り込む。
func (r *PostRepository) SearchWithFilter(
	query string,
	tags []string,
	sortBy string,
	dateFrom, dateTo *time.Time,
	limit, offset int,
) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	searchPattern := EscapeLikePattern(query)
	db := r.db.Preload("User").Preload("CodeSnippets").
		Where("(title LIKE ? OR content LIKE ?) AND is_draft = ?", searchPattern, searchPattern, false)

	// 日付範囲フィルター
	if dateFrom != nil {
		db = db.Where("created_at >= ?", dateFrom)
	}
	if dateTo != nil {
		db = db.Where("created_at <= ?", dateTo)
	}

	// タグフィルター（AND条件：全タグが付与されている投稿のみ）
	for _, tag := range tags {
		db = db.Where("id IN (SELECT post_id FROM post_tags WHERE tag = ?)", tag)
	}

	// ソート順
	switch sortBy {
	case "popular":
		db = db.Order("like_count DESC")
	case "views":
		db = db.Order("view_count DESC")
	default:
		db = db.Order("created_at DESC")
	}

	if err := db.Model(&model.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(offset).Limit(limit).Find(&posts).Error
	return posts, total, err
}

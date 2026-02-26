package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// CodeSnippetRepository はコードスニペットデータのDB操作を実装する。
type CodeSnippetRepository struct {
	db *gorm.DB
}

// NewCodeSnippetRepository は新しいCodeSnippetRepositoryインスタンスを生成する。
func NewCodeSnippetRepository(db *gorm.DB) *CodeSnippetRepository {
	return &CodeSnippetRepository{db: db}
}

// Create は新しいコードスニペットをDBに作成する。
func (r *CodeSnippetRepository) Create(snippet *model.CodeSnippet) error {
	return r.db.Create(snippet).Error
}

// FindByID は指定IDのコードスニペットを取得する。
func (r *CodeSnippetRepository) FindByID(id uint) (*model.CodeSnippet, error) {
	var snippet model.CodeSnippet
	err := r.db.First(&snippet, id).Error
	return &snippet, err
}

// FindByPostID は指定投稿IDに紐づくコードスニペット一覧を取得する。
func (r *CodeSnippetRepository) FindByPostID(postID uint) ([]model.CodeSnippet, error) {
	var snippets []model.CodeSnippet
	err := r.db.Where("post_id = ?", postID).Order("created_at ASC").Find(&snippets).Error
	return snippets, err
}

// FindByUserIDAndLanguage は指定ユーザーのスニペットをプログラミング言語でフィルタリングして取得する。
func (r *CodeSnippetRepository) FindByUserIDAndLanguage(userID uint, language string) ([]model.CodeSnippet, error) {
	var snippets []model.CodeSnippet
	err := r.db.Where("user_id = ? AND language = ?", userID, language).Order("created_at DESC").Find(&snippets).Error
	return snippets, err
}

// Update は既存のコードスニペットを更新する。
func (r *CodeSnippetRepository) Update(snippet *model.CodeSnippet) error {
	return r.db.Save(snippet).Error
}

// Delete はコードスニペットを削除する。
func (r *CodeSnippetRepository) Delete(id uint) error {
	return r.db.Delete(&model.CodeSnippet{}, id).Error
}

// CreateComment はスニペットへのインラインコメントを作成する。
// 同時にスニペットのcomment_countをインクリメントする。
func (r *CodeSnippetRepository) CreateComment(comment *model.SnippetComment) error {
	err := r.db.Create(comment).Error
	if err != nil {
		return err
	}
	return r.db.Model(&model.CodeSnippet{}).Where("id = ?", comment.SnippetID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error
}

// GetComments は指定スニペットIDのインラインコメント一覧を取得する。
// 行番号→作成日時の順でソートし、Userアソシエーションもプリロードする。
func (r *CodeSnippetRepository) GetComments(snippetID uint) ([]model.SnippetComment, error) {
	var comments []model.SnippetComment
	err := r.db.Preload("User").Where("snippet_id = ?", snippetID).
		Order("line_number ASC, created_at ASC").Find(&comments).Error
	return comments, err
}

// DeleteComment はインラインコメントを削除する。
// 所有権を検証し、スニペットのcomment_countをデクリメントする。
func (r *CodeSnippetRepository) DeleteComment(id, userID uint) error {
	var comment model.SnippetComment
	if err := r.db.First(&comment, id).Error; err != nil {
		return err
	}
	if comment.UserID != userID {
		return gorm.ErrRecordNotFound
	}
	r.db.Model(&model.CodeSnippet{}).Where("id = ?", comment.SnippetID).
		UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)"))
	return r.db.Delete(&comment).Error
}

// IncrementForkCount はスニペットのfork_countをインクリメントする。
func (r *CodeSnippetRepository) IncrementForkCount(id uint) error {
	return r.db.Model(&model.CodeSnippet{}).Where("id = ?", id).
		UpdateColumn("fork_count", gorm.Expr("fork_count + 1")).Error
}

// Favorite はスニペットをお気に入りに追加する。
func (r *CodeSnippetRepository) Favorite(userID, snippetID uint) error {
	fav := &model.CodeSnippetFavorite{
		UserID:    userID,
		SnippetID: snippetID,
	}
	return r.db.Create(fav).Error
}

// Unfavorite はスニペットのお気に入りを解除する。
func (r *CodeSnippetRepository) Unfavorite(userID, snippetID uint) error {
	return r.db.Where("user_id = ? AND snippet_id = ?", userID, snippetID).
		Delete(&model.CodeSnippetFavorite{}).Error
}

// HasFavorited は指定ユーザーが指定スニペットをお気に入りしているかを返す。
func (r *CodeSnippetRepository) HasFavorited(userID, snippetID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.CodeSnippetFavorite{}).
		Where("user_id = ? AND snippet_id = ?", userID, snippetID).
		Count(&count).Error
	return count > 0, err
}

// FindFavoritedByUserID は指定ユーザーのお気に入りスニペットをページネーション付きで取得する。
func (r *CodeSnippetRepository) FindFavoritedByUserID(userID uint, limit, offset int) ([]model.CodeSnippet, int64, error) {
	var snippets []model.CodeSnippet
	var total int64

	subQuery := r.db.Model(&model.CodeSnippetFavorite{}).
		Select("snippet_id").
		Where("user_id = ?", userID)

	query := r.db.Model(&model.CodeSnippet{}).Where("id IN (?)", subQuery)
	query.Count(&total)

	err := query.Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&snippets).Error

	return snippets, total, err
}

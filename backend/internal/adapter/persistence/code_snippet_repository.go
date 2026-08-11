package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// codeSnippetRepository は [repository.CodeSnippetRepository] の GORM 実装。
type codeSnippetRepository struct {
	db *gorm.DB
}

// NewCodeSnippetRepository は CodeSnippetRepository の GORM 実装を返す。
func NewCodeSnippetRepository(db *gorm.DB) repository.CodeSnippetRepository {
	return &codeSnippetRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.CodeSnippetRepository = (*codeSnippetRepository)(nil)

// Create は新しいコードスニペットを作成する。
func (r *codeSnippetRepository) Create(ctx context.Context, snippet *model.CodeSnippet) error {
	return r.db.WithContext(ctx).Create(snippet).Error
}

// FindByID は指定 ID のスニペットを取得する。不在の場合は (nil, nil) を返す。
func (r *codeSnippetRepository) FindByID(ctx context.Context, id uint) (*model.CodeSnippet, error) {
	var snippet model.CodeSnippet
	err := r.db.WithContext(ctx).First(&snippet, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &snippet, nil
}

// FindByPostID は指定投稿のスニペット一覧を取得する。
func (r *codeSnippetRepository) FindByPostID(ctx context.Context, postID uint) ([]model.CodeSnippet, error) {
	var snippets []model.CodeSnippet
	err := r.db.WithContext(ctx).Where("post_id = ?", postID).
		Order("created_at ASC").Find(&snippets).Error
	return snippets, err
}

// FindByUserIDAndLanguage は指定ユーザーのスニペットを言語で絞り込んで取得する。
func (r *codeSnippetRepository) FindByUserIDAndLanguage(ctx context.Context, userID uint, language string) ([]model.CodeSnippet, error) {
	var snippets []model.CodeSnippet
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND language = ?", userID, language).
		Order("created_at DESC").Find(&snippets).Error
	return snippets, err
}

// Search は言語・ファイル名・コード内容からスニペットをキーワード検索する。
func (r *codeSnippetRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.CodeSnippet, int64, error) {
	db := r.db.WithContext(ctx)
	like := escapeLikePattern(query)
	const cond = "language ILIKE ? OR file_name ILIKE ? OR code ILIKE ?"

	var total int64
	if err := db.Model(&model.CodeSnippet{}).Where(cond, like, like, like).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var snippets []model.CodeSnippet
	err := db.Where(cond, like, like, like).
		Preload("User").Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&snippets).Error
	return snippets, total, err
}

// Update は既存のスニペットを更新する。
func (r *codeSnippetRepository) Update(ctx context.Context, snippet *model.CodeSnippet) error {
	return r.db.WithContext(ctx).Save(snippet).Error
}

// Delete はスニペットを削除する。
func (r *codeSnippetRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.CodeSnippet{}, id).Error
}

// CreateComment はインラインコメントを作成し、スニペットのコメント数を加算する。
func (r *codeSnippetRepository) CreateComment(ctx context.Context, comment *model.SnippetComment) error {
	db := r.db.WithContext(ctx)
	if err := db.Create(comment).Error; err != nil {
		return err
	}
	return db.Model(&model.CodeSnippet{}).Where("id = ?", comment.SnippetID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + 1")).Error
}

// GetComments は指定スニペットのインラインコメントを行番号順で取得する。
func (r *codeSnippetRepository) GetComments(ctx context.Context, snippetID uint) ([]model.SnippetComment, error) {
	var comments []model.SnippetComment
	err := r.db.WithContext(ctx).Preload("User").
		Where("snippet_id = ?", snippetID).
		Order("line_number ASC, created_at ASC").Find(&comments).Error
	return comments, err
}

// DeleteComment は所有者のインラインコメントを削除し、スニペットのコメント数を減算する。
// 所有者でない場合は gorm.ErrRecordNotFound を返す（移行前の挙動を維持している）。
func (r *codeSnippetRepository) DeleteComment(ctx context.Context, id, userID uint) error {
	db := r.db.WithContext(ctx)
	var comment model.SnippetComment
	if err := db.First(&comment, id).Error; err != nil {
		return err
	}
	if comment.UserID != userID {
		return gorm.ErrRecordNotFound
	}
	db.Model(&model.CodeSnippet{}).Where("id = ?", comment.SnippetID).
		UpdateColumn("comment_count", gorm.Expr("GREATEST(comment_count - 1, 0)"))
	return db.Delete(&comment).Error
}

// IncrementForkCount はスニペットのフォーク数を加算する。
func (r *codeSnippetRepository) IncrementForkCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.CodeSnippet{}).Where("id = ?", id).
		UpdateColumn("fork_count", gorm.Expr("fork_count + 1")).Error
}

// Favorite はスニペットをお気に入りに追加する。
func (r *codeSnippetRepository) Favorite(ctx context.Context, userID, snippetID uint) error {
	return r.db.WithContext(ctx).Create(&model.CodeSnippetFavorite{
		UserID:    userID,
		SnippetID: snippetID,
	}).Error
}

// Unfavorite はスニペットのお気に入りを解除する。
func (r *codeSnippetRepository) Unfavorite(ctx context.Context, userID, snippetID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND snippet_id = ?", userID, snippetID).
		Delete(&model.CodeSnippetFavorite{}).Error
}

// HasFavorited は指定ユーザーが指定スニペットをお気に入りしているかを返す。
func (r *codeSnippetRepository) HasFavorited(ctx context.Context, userID, snippetID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CodeSnippetFavorite{}).
		Where("user_id = ? AND snippet_id = ?", userID, snippetID).
		Count(&count).Error
	return count > 0, err
}

// FindFavoritedByUserID はお気に入りスニペットをページネーション付きで取得する。
func (r *codeSnippetRepository) FindFavoritedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.CodeSnippet, int64, error) {
	db := r.db.WithContext(ctx)
	subQuery := db.Model(&model.CodeSnippetFavorite{}).
		Select("snippet_id").Where("user_id = ?", userID)

	var total int64
	if err := db.Model(&model.CodeSnippet{}).Where("id IN (?)", subQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var snippets []model.CodeSnippet
	err := db.Model(&model.CodeSnippet{}).Where("id IN (?)", subQuery).
		Order("created_at DESC").Limit(limit).Offset(offset).
		Find(&snippets).Error
	return snippets, total, err
}

// CountByUserID は指定ユーザーのスニペット総数を返す。
func (r *codeSnippetRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CodeSnippet{}).
		Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

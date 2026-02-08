package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ZennRepository はZenn連携データへのアクセスを提供するリポジトリ実装。
type ZennRepository struct {
	db *gorm.DB
}

// NewZennRepository は新しいZennRepositoryインスタンスを生成する。
func NewZennRepository(db *gorm.DB) *ZennRepository {
	return &ZennRepository{db: db}
}

// UpsertArticles はZenn記事データを挿入または更新する。
// zenn_idで重複判定を行い、全記事にuserIDを設定する。
func (r *ZennRepository) UpsertArticles(userID uint, articles []model.ZennArticle) error {
	if len(articles) == 0 {
		return nil
	}

	// 全記事にユーザーIDを設定
	for i := range articles {
		articles[i].UserID = userID
	}

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "zenn_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "slug", "emoji", "article_type", "liked_count", "comments_count", "published_at", "updated_at"}),
	}).Create(&articles).Error
}

// GetArticles は指定ユーザーのZenn記事を公開日降順で取得する。
func (r *ZennRepository) GetArticles(userID uint) ([]model.ZennArticle, error) {
	var articles []model.ZennArticle
	err := r.db.Where("user_id = ?", userID).Order("published_at DESC").Find(&articles).Error
	return articles, err
}

// GetStats は指定ユーザーのZenn記事統計情報を算出する。
func (r *ZennRepository) GetStats(userID uint) (*model.ZennStats, error) {
	var stats model.ZennStats

	err := r.db.Model(&model.ZennArticle{}).
		Where("user_id = ?", userID).
		Select("COUNT(*) as total_articles, COALESCE(SUM(liked_count), 0) as total_likes, COALESCE(SUM(comments_count), 0) as total_comments").
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// DeleteUserArticles は指定ユーザーの全Zenn記事を削除する。
func (r *ZennRepository) DeleteUserArticles(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.ZennArticle{}).Error
}

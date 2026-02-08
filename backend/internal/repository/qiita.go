package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// QiitaRepository はQiita連携データへのアクセスを提供するリポジトリ実装。
type QiitaRepository struct {
	db *gorm.DB
}

// NewQiitaRepository は新しいQiitaRepositoryインスタンスを生成する。
func NewQiitaRepository(db *gorm.DB) *QiitaRepository {
	return &QiitaRepository{db: db}
}

// UpsertArticles はQiita記事データを挿入または更新する。
// qiita_idで重複判定を行い、全記事にuserIDを設定する。
func (r *QiitaRepository) UpsertArticles(userID uint, articles []model.QiitaArticle) error {
	if len(articles) == 0 {
		return nil
	}

	// 全記事にユーザーIDを設定
	for i := range articles {
		articles[i].UserID = userID
	}

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "qiita_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "url", "likes_count", "comments_count", "tags", "published_at", "updated_at"}),
	}).Create(&articles).Error
}

// GetArticles は指定ユーザーのQiita記事を公開日降順で取得する。
func (r *QiitaRepository) GetArticles(userID uint) ([]model.QiitaArticle, error) {
	var articles []model.QiitaArticle
	err := r.db.Where("user_id = ?", userID).Order("published_at DESC").Find(&articles).Error
	return articles, err
}

// GetStats は指定ユーザーのQiita記事統計情報を算出する。
func (r *QiitaRepository) GetStats(userID uint) (*model.QiitaStats, error) {
	var stats model.QiitaStats

	err := r.db.Model(&model.QiitaArticle{}).
		Where("user_id = ?", userID).
		Select("COUNT(*) as total_articles, COALESCE(SUM(likes_count), 0) as total_likes, COALESCE(SUM(comments_count), 0) as total_comments").
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// DeleteUserArticles は指定ユーザーの全Qiita記事を削除する。
func (r *QiitaRepository) DeleteUserArticles(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.QiitaArticle{}).Error
}

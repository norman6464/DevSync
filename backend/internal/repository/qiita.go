package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QiitaRepository struct {
	db *gorm.DB
}

func NewQiitaRepository(db *gorm.DB) *QiitaRepository {
	return &QiitaRepository{db: db}
}

// UpsertArticles inserts or updates Qiita articles
func (r *QiitaRepository) UpsertArticles(userID uint, articles []model.QiitaArticle) error {
	if len(articles) == 0 {
		return nil
	}

	// Set user ID for all articles
	for i := range articles {
		articles[i].UserID = userID
	}

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "qiita_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "url", "likes_count", "comments_count", "tags", "published_at", "updated_at"}),
	}).Create(&articles).Error
}

// GetArticles retrieves all Qiita articles for a user
func (r *QiitaRepository) GetArticles(userID uint) ([]model.QiitaArticle, error) {
	var articles []model.QiitaArticle
	err := r.db.Where("user_id = ?", userID).Order("published_at DESC").Find(&articles).Error
	return articles, err
}

// GetStats calculates Qiita statistics for a user
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

// DeleteUserArticles removes all Qiita articles for a user
func (r *QiitaRepository) DeleteUserArticles(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.QiitaArticle{}).Error
}

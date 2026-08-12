package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// qiitaRepository は [repository.QiitaArticleRepository] の GORM 実装。
type qiitaRepository struct {
	db *gorm.DB
}

// NewQiitaRepository は QiitaArticleRepository の GORM 実装を返す。
func NewQiitaRepository(db *gorm.DB) repository.QiitaArticleRepository {
	return &qiitaRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.QiitaArticleRepository = (*qiitaRepository)(nil)

// UpsertArticles は Qiita 記事を挿入または更新する。qiita_id で重複判定し、全記事に userID を設定する。
func (r *qiitaRepository) UpsertArticles(ctx context.Context, userID uint, articles []model.QiitaArticle) error {
	if len(articles) == 0 {
		return nil
	}

	for i := range articles {
		articles[i].UserID = userID
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "qiita_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"title", "url", "likes_count", "comments_count", "tags", "published_at", "updated_at",
		}),
	}).Create(&articles).Error
}

// GetArticles は指定ユーザーの Qiita 記事を公開日の降順で取得する。
func (r *qiitaRepository) GetArticles(ctx context.Context, userID uint) ([]model.QiitaArticle, error) {
	var articles []model.QiitaArticle
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("published_at DESC").
		Find(&articles).Error
	return articles, err
}

// GetStats は指定ユーザーの Qiita 記事統計を算出する。
func (r *qiitaRepository) GetStats(ctx context.Context, userID uint) (*model.QiitaStats, error) {
	var stats model.QiitaStats
	err := r.db.WithContext(ctx).Model(&model.QiitaArticle{}).
		Where("user_id = ?", userID).
		Select("COUNT(*) as total_articles, COALESCE(SUM(likes_count), 0) as total_likes, COALESCE(SUM(comments_count), 0) as total_comments").
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// DeleteUserArticles は指定ユーザーの Qiita 記事をすべて削除する。
func (r *qiitaRepository) DeleteUserArticles(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.QiitaArticle{}).Error
}

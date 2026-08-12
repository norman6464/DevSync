package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// zennRepository は [repository.ZennArticleRepository] の GORM 実装。
type zennRepository struct {
	db *gorm.DB
}

// NewZennRepository は ZennArticleRepository の GORM 実装を返す。
func NewZennRepository(db *gorm.DB) repository.ZennArticleRepository {
	return &zennRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ZennArticleRepository = (*zennRepository)(nil)

// UpsertArticles は Zenn 記事を挿入または更新する。zenn_id で重複判定し、全記事に userID を設定する。
func (r *zennRepository) UpsertArticles(ctx context.Context, userID uint, articles []model.ZennArticle) error {
	if len(articles) == 0 {
		return nil
	}

	for i := range articles {
		articles[i].UserID = userID
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "zenn_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"title", "slug", "emoji", "article_type", "liked_count", "comments_count", "published_at", "updated_at",
		}),
	}).Create(&articles).Error
}

// GetArticles は指定ユーザーの Zenn 記事を公開日の降順で取得する。
func (r *zennRepository) GetArticles(ctx context.Context, userID uint) ([]model.ZennArticle, error) {
	var articles []model.ZennArticle
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("published_at DESC").
		Find(&articles).Error
	return articles, err
}

// GetStats は指定ユーザーの Zenn 記事統計を算出する。
func (r *zennRepository) GetStats(ctx context.Context, userID uint) (*model.ZennStats, error) {
	var stats model.ZennStats
	err := r.db.WithContext(ctx).Model(&model.ZennArticle{}).
		Where("user_id = ?", userID).
		Select("COUNT(*) as total_articles, COALESCE(SUM(liked_count), 0) as total_likes, COALESCE(SUM(comments_count), 0) as total_comments").
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// DeleteUserArticles は指定ユーザーの Zenn 記事をすべて削除する。
func (r *zennRepository) DeleteUserArticles(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.ZennArticle{}).Error
}

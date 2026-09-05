package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// zennRepository は [repository.ZennArticleRepository] の sqlc(pgx) 実装。
// UpsertArticles は複数件のupsertを1トランザクションで行うため、
// Queries だけでなくトランザクションを開始できる *pgxpool.Pool を直接保持する。
type zennRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewZennRepository は ZennArticleRepository の sqlc(pgx) 実装を返す。
func NewZennRepository(pool *pgxpool.Pool) repository.ZennArticleRepository {
	return &zennRepository{pool: pool, q: sqlcgen.New(pool)}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ZennArticleRepository = (*zennRepository)(nil)

// toModelZennArticle は sqlc の生成行を model.ZennArticle へ変換する。
func toModelZennArticle(row sqlcgen.ZennArticle) model.ZennArticle {
	t := fromTimestamptz(row.PublishedAt)
	var publishedAt time.Time
	if t != nil {
		publishedAt = *t
	}
	return model.ZennArticle{
		ID:            uint(row.ID),
		UserID:        uint(row.UserID),
		ZennID:        row.ZennID,
		Title:         row.Title,
		Slug:          row.Slug,
		Emoji:         fromStringPtr(row.Emoji),
		ArticleType:   fromStringPtr(row.ArticleType),
		LikedCount:    int(fromInt64PtrValue(row.LikedCount)),
		CommentsCount: int(fromInt64PtrValue(row.CommentsCount)),
		PublishedAt:   publishedAt,
		UpdatedAt:     timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// UpsertArticles は Zenn 記事を挿入または更新する。zenn_id で重複判定し、全記事に userID を設定する。
// 移行前のGORMバッチ作成（単一INSERT文）と同じ原子性を保つため、1トランザクションで行う。
func (r *zennRepository) UpsertArticles(ctx context.Context, userID uint, articles []model.ZennArticle) error {
	if len(articles) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	for i, article := range articles {
		row, err := q.UpsertZennArticle(ctx, sqlcgen.UpsertZennArticleParams{
			UserID:        int64(userID),
			ZennID:        article.ZennID,
			Title:         article.Title,
			Slug:          article.Slug,
			Emoji:         &article.Emoji,
			ArticleType:   &article.ArticleType,
			LikedCount:    toInt64Ptr(article.LikedCount),
			CommentsCount: toInt64Ptr(article.CommentsCount),
			PublishedAt:   toTimestamptzNotNull(article.PublishedAt),
		})
		if err != nil {
			return err
		}
		articles[i] = toModelZennArticle(row)
	}
	return tx.Commit(ctx)
}

// GetArticles は指定ユーザーの Zenn 記事を公開日の降順で取得する。
func (r *zennRepository) GetArticles(ctx context.Context, userID uint) ([]model.ZennArticle, error) {
	rows, err := r.q.ListZennArticlesByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	articles := make([]model.ZennArticle, len(rows))
	for i, row := range rows {
		articles[i] = toModelZennArticle(row)
	}
	return articles, nil
}

// GetStats は指定ユーザーの Zenn 記事統計を算出する。
func (r *zennRepository) GetStats(ctx context.Context, userID uint) (*model.ZennStats, error) {
	row, err := r.q.GetZennStatsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	return &model.ZennStats{
		TotalArticles: int(row.TotalArticles),
		TotalLikes:    int(row.TotalLikes),
		TotalComments: int(row.TotalComments),
	}, nil
}

// DeleteUserArticles は指定ユーザーの Zenn 記事をすべて削除する。
func (r *zennRepository) DeleteUserArticles(ctx context.Context, userID uint) error {
	return r.q.DeleteZennArticlesByUser(ctx, int64(userID))
}

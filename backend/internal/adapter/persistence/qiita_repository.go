package persistence

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// qiitaRepository は [repository.QiitaArticleRepository] の sqlc(pgx) 実装。
// UpsertArticles は複数件のupsertを1トランザクションで行うため、
// Queries だけでなくトランザクションを開始できる *pgxpool.Pool を直接保持する。
type qiitaRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewQiitaRepository は QiitaArticleRepository の sqlc(pgx) 実装を返す。
func NewQiitaRepository(pool *pgxpool.Pool) repository.QiitaArticleRepository {
	return &qiitaRepository{pool: pool, q: sqlcgen.New(pool)}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.QiitaArticleRepository = (*qiitaRepository)(nil)

// toModelQiitaArticle は sqlc の生成行を model.QiitaArticle へ変換する。
func toModelQiitaArticle(row sqlcgen.QiitaArticle) model.QiitaArticle {
	t := fromTimestamptz(row.PublishedAt)
	var publishedAt time.Time
	if t != nil {
		publishedAt = *t
	}
	return model.QiitaArticle{
		ID:            uint(row.ID),
		UserID:        uint(row.UserID),
		QiitaID:       row.QiitaID,
		Title:         row.Title,
		URL:           row.Url,
		LikesCount:    int(fromInt64PtrValue(row.LikesCount)),
		CommentsCount: int(fromInt64PtrValue(row.CommentsCount)),
		Tags:          fromStringPtr(row.Tags),
		PublishedAt:   publishedAt,
		UpdatedAt:     timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// UpsertArticles は Qiita 記事を挿入または更新する。qiita_id で重複判定し、全記事に userID を設定する。
// 移行前のGORMバッチ作成（単一INSERT文）と同じ原子性を保つため、1トランザクションで行う。
func (r *qiitaRepository) UpsertArticles(ctx context.Context, userID uint, articles []model.QiitaArticle) error {
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
		row, err := q.UpsertQiitaArticle(ctx, sqlcgen.UpsertQiitaArticleParams{
			UserID:        int64(userID),
			QiitaID:       article.QiitaID,
			Title:         article.Title,
			Url:           article.URL,
			LikesCount:    toInt64Ptr(article.LikesCount),
			CommentsCount: toInt64Ptr(article.CommentsCount),
			Tags:          &article.Tags,
			PublishedAt:   toTimestamptzNotNull(article.PublishedAt),
		})
		if err != nil {
			return err
		}
		articles[i] = toModelQiitaArticle(row)
	}
	return tx.Commit(ctx)
}

// GetArticles は指定ユーザーの Qiita 記事を公開日の降順で取得する。
func (r *qiitaRepository) GetArticles(ctx context.Context, userID uint) ([]model.QiitaArticle, error) {
	rows, err := r.q.ListQiitaArticlesByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	articles := make([]model.QiitaArticle, len(rows))
	for i, row := range rows {
		articles[i] = toModelQiitaArticle(row)
	}
	return articles, nil
}

// GetStats は指定ユーザーの Qiita 記事統計を算出する。
func (r *qiitaRepository) GetStats(ctx context.Context, userID uint) (*model.QiitaStats, error) {
	row, err := r.q.GetQiitaStatsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	return &model.QiitaStats{
		TotalArticles: int(row.TotalArticles),
		TotalLikes:    int(row.TotalLikes),
		TotalComments: int(row.TotalComments),
	}, nil
}

// DeleteUserArticles は指定ユーザーの Qiita 記事をすべて削除する。
func (r *qiitaRepository) DeleteUserArticles(ctx context.Context, userID uint) error {
	return r.q.DeleteQiitaArticlesByUser(ctx, int64(userID))
}

package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// githubRepository は [repository.GitHubRepository] の sqlc(pgx) 実装。
// 各Upsertは複数件のupsertを1トランザクションで行うため、
// Queries だけでなくトランザクションを開始できる *pgxpool.Pool を直接保持する。
type githubRepository struct {
	pool *pgxpool.Pool
	q    *sqlcgen.Queries
}

// NewGitHubRepository は GitHubRepository の sqlc(pgx) 実装を返す。
func NewGitHubRepository(pool *pgxpool.Pool) repository.GitHubRepository {
	return &githubRepository{pool: pool, q: sqlcgen.New(pool)}
}

var _ repository.GitHubRepository = (*githubRepository)(nil)

func toModelGitHubContribution(row sqlcgen.GitHubContribution) model.GitHubContribution {
	return model.GitHubContribution{
		ID:        uint(row.ID),
		UserID:    uint(row.UserID),
		Date:      timeValue(fromTimestamptz(row.Date)),
		Count:     int(row.Count),
		CreatedAt: timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt: timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

func toModelGitHubLanguageStat(row sqlcgen.GitHubLanguageStat) model.GitHubLanguageStat {
	return model.GitHubLanguageStat{
		ID:        uint(row.ID),
		UserID:    uint(row.UserID),
		Language:  row.Language,
		Bytes:     row.Bytes,
		RepoCount: int(row.RepoCount),
		UpdatedAt: timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

func toModelGitHubRepository(row sqlcgen.GitHubRepository) model.GitHubRepository {
	return model.GitHubRepository{
		ID:           uint(row.ID),
		UserID:       uint(row.UserID),
		GitHubRepoID: row.GitHubRepoID,
		Name:         row.Name,
		FullName:     fromStringPtr(row.FullName),
		Description:  fromStringPtr(row.Description),
		Language:     fromStringPtr(row.Language),
		Stars:        int(fromInt64PtrValue(row.Stars)),
		Forks:        int(fromInt64PtrValue(row.Forks)),
		IsPrivate:    row.IsPrivate,
		UpdatedAt:    timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// UpsertContributions は日別コントリビューションを (user_id, date) で重複判定して保存する。
// 移行前のGORMバッチ作成と同じ原子性を保つため、1トランザクションで行う。
func (r *githubRepository) UpsertContributions(ctx context.Context, contributions []model.GitHubContribution) error {
	if len(contributions) == 0 {
		return nil
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	for i, c := range contributions {
		row, err := q.UpsertGitHubContribution(ctx, sqlcgen.UpsertGitHubContributionParams{
			UserID: int64(c.UserID),
			Date:   toTimestamptzNotNull(c.Date),
			Count:  int64(c.Count),
		})
		if err != nil {
			return err
		}
		contributions[i] = toModelGitHubContribution(row)
	}
	return tx.Commit(ctx)
}

// GetContributions は日付の昇順でコントリビューションを返す。
func (r *githubRepository) GetContributions(ctx context.Context, userID uint) ([]model.GitHubContribution, error) {
	rows, err := r.q.ListAllGitHubContributionsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	contributions := make([]model.GitHubContribution, len(rows))
	for i, row := range rows {
		contributions[i] = toModelGitHubContribution(row)
	}
	return contributions, nil
}

// UpsertLanguageStats は言語統計を (user_id, language) で重複判定して保存する。
// 移行前のGORMバッチ作成と同じ原子性を保つため、1トランザクションで行う。
func (r *githubRepository) UpsertLanguageStats(ctx context.Context, stats []model.GitHubLanguageStat) error {
	if len(stats) == 0 {
		return nil
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	for i, s := range stats {
		row, err := q.UpsertGitHubLanguageStat(ctx, sqlcgen.UpsertGitHubLanguageStatParams{
			UserID:    int64(s.UserID),
			Language:  s.Language,
			Bytes:     s.Bytes,
			RepoCount: int64(s.RepoCount),
		})
		if err != nil {
			return err
		}
		stats[i] = toModelGitHubLanguageStat(row)
	}
	return tx.Commit(ctx)
}

// GetLanguageStats はバイト数の降順で言語統計を返す。
func (r *githubRepository) GetLanguageStats(ctx context.Context, userID uint) ([]model.GitHubLanguageStat, error) {
	rows, err := r.q.ListGitHubLanguageStatsByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	stats := make([]model.GitHubLanguageStat, len(rows))
	for i, row := range rows {
		stats[i] = toModelGitHubLanguageStat(row)
	}
	return stats, nil
}

// UpsertRepos はリポジトリを GitHub 側のリポジトリ ID で重複判定して保存する。
// 移行前のGORMバッチ作成と同じ原子性を保つため、1トランザクションで行う。
func (r *githubRepository) UpsertRepos(ctx context.Context, repos []model.GitHubRepository) error {
	if len(repos) == 0 {
		return nil
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := r.q.WithTx(tx)
	for i, repo := range repos {
		row, err := q.UpsertGitHubRepo(ctx, sqlcgen.UpsertGitHubRepoParams{
			UserID:       int64(repo.UserID),
			GitHubRepoID: repo.GitHubRepoID,
			Name:         repo.Name,
			FullName:     &repo.FullName,
			Description:  &repo.Description,
			Language:     &repo.Language,
			Stars:        toInt64Ptr(repo.Stars),
			Forks:        toInt64Ptr(repo.Forks),
			IsPrivate:    repo.IsPrivate,
		})
		if err != nil {
			return err
		}
		repos[i] = toModelGitHubRepository(row)
	}
	return tx.Commit(ctx)
}

// GetRepos はスター数の降順でリポジトリを返す。
func (r *githubRepository) GetRepos(ctx context.Context, userID uint) ([]model.GitHubRepository, error) {
	rows, err := r.q.ListGitHubReposByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	repos := make([]model.GitHubRepository, len(rows))
	for i, row := range rows {
		repos[i] = toModelGitHubRepository(row)
	}
	return repos, nil
}

// DeleteUserData は指定ユーザーの GitHub 連携データ（コントリビューション・言語統計・リポジトリ）を削除する。
func (r *githubRepository) DeleteUserData(ctx context.Context, userID uint) error {
	uid := int64(userID)
	if err := r.q.DeleteGitHubContributionsByUser(ctx, uid); err != nil {
		return err
	}
	if err := r.q.DeleteGitHubLanguageStatsByUser(ctx, uid); err != nil {
		return err
	}
	return r.q.DeleteGitHubReposByUser(ctx, uid)
}

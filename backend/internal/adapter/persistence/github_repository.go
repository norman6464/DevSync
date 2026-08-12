package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// githubRepository は [repository.GitHubRepository] の GORM 実装。
type githubRepository struct {
	db *gorm.DB
}

// NewGitHubRepository は GitHubRepository の GORM 実装を返す。
func NewGitHubRepository(db *gorm.DB) repository.GitHubRepository {
	return &githubRepository{db: db}
}

var _ repository.GitHubRepository = (*githubRepository)(nil)

// UpsertContributions は日別コントリビューションを (user_id, date) で重複判定して保存する。
func (r *githubRepository) UpsertContributions(ctx context.Context, contributions []model.GitHubContribution) error {
	if len(contributions) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"count", "updated_at"}),
	}).Create(&contributions).Error
}

// GetContributions は日付の昇順でコントリビューションを返す。
func (r *githubRepository) GetContributions(ctx context.Context, userID uint) ([]model.GitHubContribution, error) {
	var contributions []model.GitHubContribution
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("date ASC").Find(&contributions).Error
	return contributions, err
}

// UpsertLanguageStats は言語統計を (user_id, language) で重複判定して保存する。
func (r *githubRepository) UpsertLanguageStats(ctx context.Context, stats []model.GitHubLanguageStat) error {
	if len(stats) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "language"}},
		DoUpdates: clause.AssignmentColumns([]string{"bytes", "repo_count", "updated_at"}),
	}).Create(&stats).Error
}

// GetLanguageStats はバイト数の降順で言語統計を返す。
func (r *githubRepository) GetLanguageStats(ctx context.Context, userID uint) ([]model.GitHubLanguageStat, error) {
	var stats []model.GitHubLanguageStat
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("bytes DESC").Find(&stats).Error
	return stats, err
}

// UpsertRepos はリポジトリを GitHub 側のリポジトリ ID で重複判定して保存する。
func (r *githubRepository) UpsertRepos(ctx context.Context, repos []model.GitHubRepository) error {
	if len(repos) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "git_hub_repo_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "full_name", "description", "language", "stars", "forks", "is_private", "updated_at"}),
	}).Create(&repos).Error
}

// GetRepos はスター数の降順でリポジトリを返す。
func (r *githubRepository) GetRepos(ctx context.Context, userID uint) ([]model.GitHubRepository, error) {
	var repos []model.GitHubRepository
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("stars DESC").Find(&repos).Error
	return repos, err
}

// DeleteUserData は指定ユーザーの GitHub 連携データ（コントリビューション・言語統計・リポジトリ）を削除する。
func (r *githubRepository) DeleteUserData(ctx context.Context, userID uint) error {
	db := r.db.WithContext(ctx)
	if err := db.Where("user_id = ?", userID).Delete(&model.GitHubContribution{}).Error; err != nil {
		return err
	}
	if err := db.Where("user_id = ?", userID).Delete(&model.GitHubLanguageStat{}).Error; err != nil {
		return err
	}
	return db.Where("user_id = ?", userID).Delete(&model.GitHubRepository{}).Error
}

package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GitHubRepository はGitHub連携データへのアクセスを提供するリポジトリ実装。
// コントリビューション、言語統計、リポジトリのUpsert（挿入or更新）操作を提供する。
type GitHubRepository struct {
	db *gorm.DB
}

// NewGitHubRepository は新しいGitHubRepositoryインスタンスを生成する。
func NewGitHubRepository(db *gorm.DB) *GitHubRepository {
	return &GitHubRepository{db: db}
}

// UpsertContributions はGitHubコントリビューションデータを挿入または更新する。
// user_idとdateの組み合わせで重複判定を行う。
func (r *GitHubRepository) UpsertContributions(contributions []model.GitHubContribution) error {
	if len(contributions) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"count", "updated_at"}),
	}).Create(&contributions).Error
}

// GetContributions は指定ユーザーのGitHubコントリビューションを日付昇順で取得する。
func (r *GitHubRepository) GetContributions(userID uint) ([]model.GitHubContribution, error) {
	var contributions []model.GitHubContribution
	err := r.db.Where("user_id = ?", userID).Order("date ASC").Find(&contributions).Error
	return contributions, err
}

// UpsertLanguageStats はGitHub言語統計データを挿入または更新する。
// user_idとlanguageの組み合わせで重複判定を行う。
func (r *GitHubRepository) UpsertLanguageStats(stats []model.GitHubLanguageStat) error {
	if len(stats) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "language"}},
		DoUpdates: clause.AssignmentColumns([]string{"bytes", "repo_count", "updated_at"}),
	}).Create(&stats).Error
}

// GetLanguageStats は指定ユーザーのGitHub言語統計をバイト数降順で取得する。
func (r *GitHubRepository) GetLanguageStats(userID uint) ([]model.GitHubLanguageStat, error) {
	var stats []model.GitHubLanguageStat
	err := r.db.Where("user_id = ?", userID).Order("bytes DESC").Find(&stats).Error
	return stats, err
}

// UpsertRepos はGitHubリポジトリデータを挿入または更新する。
// git_hub_repo_idで重複判定を行う。
func (r *GitHubRepository) UpsertRepos(repos []model.GitHubRepository) error {
	if len(repos) == 0 {
		return nil
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "git_hub_repo_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "full_name", "description", "language", "stars", "forks", "is_private", "updated_at"}),
	}).Create(&repos).Error
}

// GetRepos は指定ユーザーのGitHubリポジトリをスター数降順で取得する。
func (r *GitHubRepository) GetRepos(userID uint) ([]model.GitHubRepository, error) {
	var repos []model.GitHubRepository
	err := r.db.Where("user_id = ?", userID).Order("stars DESC").Find(&repos).Error
	return repos, err
}

// DeleteUserData は指定ユーザーの全GitHub連携データ（コントリビューション、言語統計、リポジトリ）を削除する。
func (r *GitHubRepository) DeleteUserData(userID uint) error {
	r.db.Where("user_id = ?", userID).Delete(&model.GitHubContribution{})
	r.db.Where("user_id = ?", userID).Delete(&model.GitHubLanguageStat{})
	r.db.Where("user_id = ?", userID).Delete(&model.GitHubRepository{})
	return nil
}

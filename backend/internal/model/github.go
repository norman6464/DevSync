package model

import "time"

// GitHubContribution はGitHubから同期した日別コントリビューション（草）データを表す。
// uniqueIndex制約でユーザーごとに1日1レコードを保証する。
type GitHubContribution struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Date      time.Time `json:"date"`  // コントリビューション日
	Count     int       `json:"count"` // その日のコントリビューション数
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GitHubLanguageStat はGitHubリポジトリの言語別統計データを表す。
type GitHubLanguageStat struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Language  string    `json:"language"`   // プログラミング言語名
	Bytes     int64     `json:"bytes"`      // 使用バイト数
	RepoCount int       `json:"repo_count"` // 使用リポジトリ数
	UpdatedAt time.Time `json:"updated_at"`
}

// GitHubRepository はGitHubから同期したリポジトリ情報を表す。
type GitHubRepository struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	GitHubRepoID int64     `json:"github_repo_id"` // GitHub側のリポジトリID
	Name         string    `json:"name"`
	FullName     string    `json:"full_name"` // "owner/repo" 形式
	Description  string    `json:"description"`
	Language     string    `json:"language"` // メイン言語
	Stars        int       `json:"stars"`
	Forks        int       `json:"forks"`
	IsPrivate    bool      `json:"is_private"` // プライベートリポジトリかどうか
	UpdatedAt    time.Time `json:"updated_at"`
}

// GitHubUserInfo は GitHub API から取得するユーザー情報を表す。
type GitHubUserInfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// GitHubContributionDay は GitHub から取得した 1 日分のコントリビューション数を表す。
type GitHubContributionDay struct {
	Date  time.Time
	Count int
}

// GitHubRepoSummary は GitHub API から取得したリポジトリ情報を表す。
type GitHubRepoSummary struct {
	ID          int64
	Name        string
	FullName    string
	Description string
	Language    string
	Stars       int
	Forks       int
	Private     bool
}

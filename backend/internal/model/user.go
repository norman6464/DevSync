// Package model はDevSyncアプリケーションのデータモデルを定義する。
// GORMのstructタグでDB制約を、jsonタグでAPIレスポンス形式を制御する。
package model

import "time"

// User はDevSyncのユーザーアカウント情報を表す。
// 認証情報、プロフィール、外部サービス連携、スキル情報を保持する。
type User struct {
	ID                  uint      `json:"id" gorm:"primaryKey"`
	Username            string    `json:"username" gorm:"uniqueIndex;not null"` // 一意のユーザー名（プロフィールURL用）
	Name                string    `json:"name" gorm:"not null"`
	Email               string    `json:"email" gorm:"uniqueIndex;not null"`
	Password            string    `json:"-"` // json:"-" でAPIレスポンスから除外
	AvatarURL           string    `json:"avatar_url"`
	Bio                 string    `json:"bio"`
	GitHubID            int64     `json:"github_id"` // GitHub OAuth連携時のユーザーID（未連携は 0）
	GitHubUsername      string    `json:"github_username"`
	GitHubToken         string    `json:"-"`                                      // json:"-" でAPIレスポンスから除外
	GitHubConnected     bool      `json:"github_connected" gorm:"default:false"`  // GitHub連携済みフラグ
	SpotifyConnected    bool      `json:"spotify_connected" gorm:"default:false"` // Spotify連携済みフラグ
	SpotifyToken        string    `json:"-"`                                      // Spotifyアクセストークン（APIレスポンスから除外）
	SpotifyRefreshToken string    `json:"-"`                                      // Spotifyリフレッシュトークン（APIレスポンスから除外）
	SpotifyTokenExpiry  time.Time `json:"-"`                                      // Spotifyトークン有効期限
	ZennUsername        string    `json:"zenn_username"`
	QiitaUsername       string    `json:"qiita_username"`
	AtCoderUsername     string    `json:"atcoder_username"`                          // AtCoderユーザー名
	PaizaRank           string    `json:"paiza_rank"`                                // paizaランク（S/A/B/C/D/E、自己申告）
	SkillsLanguages     string    `json:"skills_languages"`                          // プログラミング言語スキル（カンマ区切り）
	SkillsFrameworks    string    `json:"skills_frameworks"`                         // フレームワークスキル（カンマ区切り）
	OnboardingCompleted bool      `json:"onboarding_completed" gorm:"default:false"` // 初回セットアップ完了フラグ
	EmailWeeklyReport   bool      `json:"email_weekly_report" gorm:"default:true"`   // ウィークリーレポートメール配信フラグ
	EmailLanguage       string    `json:"email_language" gorm:"default:'ja'"`        // メール配信言語（ja, en, ko等）
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

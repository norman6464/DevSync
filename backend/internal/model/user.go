// Package model はDevSyncアプリケーションのデータモデルを定義する。
// GORMのstructタグでDB制約を、jsonタグでAPIレスポンス形式を制御する。
package model

import "time"

// User はDevSyncのユーザーアカウント情報を表す。
// 認証情報、プロフィール、外部サービス連携、スキル情報を保持する。
type User struct {
	ID                  uint      `json:"id" gorm:"primaryKey"`
	Name                string    `json:"name" gorm:"not null"`
	Email               string    `json:"email" gorm:"uniqueIndex;not null"`
	Password            string    `json:"-"`                                         // json:"-" でAPIレスポンスから除外
	AvatarURL           string    `json:"avatar_url"`
	Bio                 string    `json:"bio"`
	GitHubID            int64     `json:"github_id" gorm:"uniqueIndex"`              // GitHub OAuth連携時のユーザーID
	GitHubUsername      string    `json:"github_username"`
	GitHubToken         string    `json:"-"`                                         // json:"-" でAPIレスポンスから除外
	GitHubConnected     bool      `json:"github_connected" gorm:"default:false"`     // GitHub連携済みフラグ
	ZennUsername        string    `json:"zenn_username"`
	QiitaUsername       string    `json:"qiita_username"`
	SkillsLanguages     string    `json:"skills_languages"`                          // プログラミング言語スキル（カンマ区切り）
	SkillsFrameworks    string    `json:"skills_frameworks"`                         // フレームワークスキル（カンマ区切り）
	OnboardingCompleted bool      `json:"onboarding_completed" gorm:"default:false"` // 初回セットアップ完了フラグ
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

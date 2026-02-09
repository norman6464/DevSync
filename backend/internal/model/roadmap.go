package model

import "time"

// RoadmapCategory は学習ロードマップのカテゴリを表す型。
type RoadmapCategory string

// 学習ロードマップのカテゴリ定数群。
const (
	RoadmapCategoryLanguage  RoadmapCategory = "language"  // プログラミング言語
	RoadmapCategoryFramework RoadmapCategory = "framework" // フレームワーク
	RoadmapCategorySkill     RoadmapCategory = "skill"     // スキル
	RoadmapCategoryProject   RoadmapCategory = "project"   // プロジェクト
	RoadmapCategoryOther     RoadmapCategory = "other"     // その他
)

// RoadmapStatus は学習ロードマップの進行状態を表す型。
type RoadmapStatus string

// 学習ロードマップのステータス定数群。
const (
	RoadmapStatusActive    RoadmapStatus = "active"    // 進行中
	RoadmapStatusCompleted RoadmapStatus = "completed" // 完了
)

// Roadmap はユーザーの学習ロードマップを表す。
// 複数の RoadmapStep を持ち、Progress は全ステップの完了率から自動計算される。
// IsPublic がtrueの場合、他のユーザーにも公開される。
type Roadmap struct {
	ID                 uint            `json:"id" gorm:"primaryKey"`
	UserID             uint            `json:"user_id" gorm:"not null;index"`
	User               User            `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Title              string          `json:"title" gorm:"not null;size:200"`
	Description        string          `json:"description" gorm:"type:text"`
	Category           RoadmapCategory `json:"category" gorm:"default:'other'"`
	IsPublic           bool            `json:"is_public" gorm:"default:false;index"`        // 公開フラグ
	IsTemplate         bool            `json:"is_template" gorm:"default:false;index"`      // テンプレートフラグ
	StepCount          int             `json:"step_count" gorm:"default:0"`                  // 全ステップ数
	CompletedStepCount int             `json:"completed_step_count" gorm:"default:0"`        // 完了済みステップ数
	Progress           int             `json:"progress" gorm:"default:0"`                    // 進捗率（0〜100、自動計算）
	Status             RoadmapStatus   `json:"status" gorm:"default:'active'"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	CompletedAt        *time.Time      `json:"completed_at"` // 完了日時（完了時に自動設定）

	// リレーション
	Steps []RoadmapStep `json:"steps,omitempty" gorm:"foreignKey:RoadmapID;constraint:OnDelete:CASCADE"`
}

// RoadmapStep は学習ロードマップ内の個別ステップを表す。
// OrderIndex で表示順序を管理し、IsCompleted で完了状態を追跡する。
type RoadmapStep struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	RoadmapID   uint       `json:"roadmap_id" gorm:"not null;index"`           // 所属するロードマップのID
	Title       string     `json:"title" gorm:"not null;size:200"`
	Description string     `json:"description" gorm:"type:text"`
	OrderIndex  int        `json:"order_index" gorm:"not null;default:0"`      // 表示順序（0始まり）
	IsCompleted bool       `json:"is_completed" gorm:"default:false"`          // 完了フラグ
	CompletedAt *time.Time `json:"completed_at"`                               // 完了日時
	ResourceURL string     `json:"resource_url" gorm:"size:500"`               // 参考リソースのURL
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// RoadmapStats はユーザーのロードマップ統計情報を表す。
// DBテーブルには対応せず、集計結果を格納するDTO。
type RoadmapStats struct {
	TotalRoadmaps     int `json:"total_roadmaps"`     // ロードマップ総数
	ActiveRoadmaps    int `json:"active_roadmaps"`    // 進行中のロードマップ数
	CompletedRoadmaps int `json:"completed_roadmaps"` // 完了済みロードマップ数
	TotalSteps        int `json:"total_steps"`        // 全ステップ総数
	CompletedSteps    int `json:"completed_steps"`    // 完了済みステップ総数
}

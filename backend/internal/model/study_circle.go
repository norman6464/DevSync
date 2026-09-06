package model

import "time"

// StudyCircleStatus はスタディサークルの状態を表す。
type StudyCircleStatus string

const (
	StudyCircleStatusActive    StudyCircleStatus = "active"
	StudyCircleStatusCompleted StudyCircleStatus = "completed"
	StudyCircleStatusArchived  StudyCircleStatus = "archived"
)

// StudyCircleMemberRole はサークルメンバーの役割を表す。
type StudyCircleMemberRole string

const (
	StudyCircleRoleOwner  StudyCircleMemberRole = "owner"
	StudyCircleRoleMember StudyCircleMemberRole = "member"
)

// StudyCircle はスタディサークル（学習グループ）を表す。
type StudyCircle struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Topic       string            `json:"topic"`
	Description string            `json:"description"`
	OwnerID     uint              `json:"owner_id"`
	Owner       *User             `json:"owner,omitempty"`
	MaxMembers  int               `json:"max_members"`
	Status      StudyCircleStatus `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`

	Steps   []StudyCircleStep   `json:"steps,omitempty"`
	Members []StudyCircleMember `json:"members,omitempty"`
}

// StudyCircleMember はサークルのメンバーシップを表す。
type StudyCircleMember struct {
	ID       uint                  `json:"id"`
	CircleID uint                  `json:"circle_id"`
	UserID   uint                  `json:"user_id"`
	User     *User                 `json:"user,omitempty"`
	Role     StudyCircleMemberRole `json:"role"`
	JoinedAt time.Time             `json:"joined_at"`
}

// StudyCircleStep はサークルの共有ロードマップのステップを表す。
type StudyCircleStep struct {
	ID          uint      `json:"id"`
	CircleID    uint      `json:"circle_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	OrderIndex  int       `json:"order_index"`
	ResourceURL string    `json:"resource_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StudyCircleMemberProgress はメンバー別のステップ進捗を表す。
type StudyCircleMemberProgress struct {
	ID          uint       `json:"id"`
	CircleID    uint       `json:"circle_id"`
	StepID      uint       `json:"step_id"`
	UserID      uint       `json:"user_id"`
	User        *User      `json:"user,omitempty"`
	IsCompleted bool       `json:"is_completed"`
	CompletedAt *time.Time `json:"completed_at"`
}

// StudyCircleCheckin は日次チェックインを表す。
type StudyCircleCheckin struct {
	ID        uint      `json:"id"`
	CircleID  uint      `json:"circle_id"`
	UserID    uint      `json:"user_id"`
	User      *User     `json:"user,omitempty"`
	Date      string    `json:"date"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// CircleMemberStreak はサークル内のメンバーストリーク情報を表すDTO。
type CircleMemberStreak struct {
	UserID        uint   `json:"user_id"`
	UserName      string `json:"user_name"`
	AvatarURL     string `json:"avatar_url"`
	CurrentStreak int    `json:"current_streak"`
	TotalCheckins int    `json:"total_checkins"`
}

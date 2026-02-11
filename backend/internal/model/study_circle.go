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
	ID          uint              `json:"id" gorm:"primaryKey"`
	Name        string            `json:"name" gorm:"size:200;not null"`
	Topic       string            `json:"topic" gorm:"size:200;not null"`
	Description string            `json:"description" gorm:"type:text"`
	OwnerID     uint              `json:"owner_id" gorm:"not null;index"`
	Owner       *User             `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	MaxMembers  int               `json:"max_members" gorm:"default:5;not null"`
	Status      StudyCircleStatus `json:"status" gorm:"default:'active'"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`

	Steps   []StudyCircleStep   `json:"steps,omitempty" gorm:"foreignKey:CircleID;constraint:OnDelete:CASCADE"`
	Members []StudyCircleMember `json:"members,omitempty" gorm:"foreignKey:CircleID;constraint:OnDelete:CASCADE"`
}

// StudyCircleMember はサークルのメンバーシップを表す。
type StudyCircleMember struct {
	ID       uint                  `json:"id" gorm:"primaryKey"`
	CircleID uint                  `json:"circle_id" gorm:"not null;uniqueIndex:idx_circle_user"`
	UserID   uint                  `json:"user_id" gorm:"not null;uniqueIndex:idx_circle_user"`
	User     *User                 `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Role     StudyCircleMemberRole `json:"role" gorm:"default:'member'"`
	JoinedAt time.Time             `json:"joined_at"`
}

// StudyCircleStep はサークルの共有ロードマップのステップを表す。
type StudyCircleStep struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CircleID    uint      `json:"circle_id" gorm:"not null;index"`
	Title       string    `json:"title" gorm:"size:200;not null"`
	Description string    `json:"description" gorm:"type:text"`
	OrderIndex  int       `json:"order_index" gorm:"default:0"`
	ResourceURL string    `json:"resource_url" gorm:"size:500"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StudyCircleMemberProgress はメンバー別のステップ進捗を表す。
type StudyCircleMemberProgress struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	CircleID    uint       `json:"circle_id" gorm:"not null;uniqueIndex:idx_circle_step_user"`
	StepID      uint       `json:"step_id" gorm:"not null;uniqueIndex:idx_circle_step_user"`
	UserID      uint       `json:"user_id" gorm:"not null;uniqueIndex:idx_circle_step_user"`
	User        *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	IsCompleted bool       `json:"is_completed" gorm:"default:false"`
	CompletedAt *time.Time `json:"completed_at"`
}

// StudyCircleCheckin は日次チェックインを表す。
type StudyCircleCheckin struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CircleID  uint      `json:"circle_id" gorm:"not null;uniqueIndex:idx_checkin_unique"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_checkin_unique"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Date      string    `json:"date" gorm:"size:10;not null;uniqueIndex:idx_checkin_unique"`
	Content   string    `json:"content" gorm:"size:500;not null"`
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

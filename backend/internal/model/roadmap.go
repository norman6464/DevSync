package model

import "time"

// RoadmapCategory represents the category of a learning roadmap
type RoadmapCategory string

const (
	RoadmapCategoryLanguage  RoadmapCategory = "language"
	RoadmapCategoryFramework RoadmapCategory = "framework"
	RoadmapCategorySkill     RoadmapCategory = "skill"
	RoadmapCategoryProject   RoadmapCategory = "project"
	RoadmapCategoryOther     RoadmapCategory = "other"
)

// RoadmapStatus represents the status of a learning roadmap
type RoadmapStatus string

const (
	RoadmapStatusActive    RoadmapStatus = "active"
	RoadmapStatusCompleted RoadmapStatus = "completed"
)

// Roadmap represents a user's learning roadmap with multiple steps
type Roadmap struct {
	ID                 uint            `json:"id" gorm:"primaryKey"`
	UserID             uint            `json:"user_id" gorm:"not null;index"`
	User               User            `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Title              string          `json:"title" gorm:"not null;size:200"`
	Description        string          `json:"description" gorm:"type:text"`
	Category           RoadmapCategory `json:"category" gorm:"default:'other'"`
	IsPublic           bool            `json:"is_public" gorm:"default:false;index"`
	StepCount          int             `json:"step_count" gorm:"default:0"`
	CompletedStepCount int             `json:"completed_step_count" gorm:"default:0"`
	Progress           int             `json:"progress" gorm:"default:0"` // 0-100, auto-calculated
	Status             RoadmapStatus   `json:"status" gorm:"default:'active'"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	CompletedAt        *time.Time      `json:"completed_at"`

	// Relations
	Steps []RoadmapStep `json:"steps,omitempty" gorm:"foreignKey:RoadmapID;constraint:OnDelete:CASCADE"`
}

// RoadmapStep represents a single step in a learning roadmap
type RoadmapStep struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	RoadmapID   uint       `json:"roadmap_id" gorm:"not null;index"`
	Title       string     `json:"title" gorm:"not null;size:200"`
	Description string     `json:"description" gorm:"type:text"`
	OrderIndex  int        `json:"order_index" gorm:"not null;default:0"`
	IsCompleted bool       `json:"is_completed" gorm:"default:false"`
	CompletedAt *time.Time `json:"completed_at"`
	ResourceURL string     `json:"resource_url" gorm:"size:500"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// RoadmapStats represents aggregated roadmap statistics for a user
type RoadmapStats struct {
	TotalRoadmaps     int `json:"total_roadmaps"`
	ActiveRoadmaps    int `json:"active_roadmaps"`
	CompletedRoadmaps int `json:"completed_roadmaps"`
	TotalSteps        int `json:"total_steps"`
	CompletedSteps    int `json:"completed_steps"`
}

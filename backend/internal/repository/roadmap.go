package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

type RoadmapRepository struct {
	db *gorm.DB
}

func NewRoadmapRepository(db *gorm.DB) *RoadmapRepository {
	return &RoadmapRepository{db: db}
}

// === Roadmap CRUD ===

// Create creates a new roadmap
func (r *RoadmapRepository) Create(roadmap *model.Roadmap) error {
	return r.db.Create(roadmap).Error
}

// Update updates an existing roadmap
func (r *RoadmapRepository) Update(roadmap *model.Roadmap) error {
	return r.db.Save(roadmap).Error
}

// Delete deletes a roadmap (cascade deletes steps)
func (r *RoadmapRepository) Delete(id uint) error {
	return r.db.Delete(&model.Roadmap{}, id).Error
}

// FindByID finds a roadmap by ID with steps preloaded
func (r *RoadmapRepository) FindByID(id uint) (*model.Roadmap, error) {
	var roadmap model.Roadmap
	err := r.db.Preload("Steps", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_index ASC")
	}).Preload("User").First(&roadmap, id).Error
	if err != nil {
		return nil, err
	}
	return &roadmap, nil
}

// GetByUserID gets all roadmaps for a user (without steps)
func (r *RoadmapRepository) GetByUserID(userID uint) ([]model.Roadmap, error) {
	var roadmaps []model.Roadmap
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&roadmaps).Error
	return roadmaps, err
}

// GetPublicRoadmaps gets all public roadmaps with pagination
func (r *RoadmapRepository) GetPublicRoadmaps(limit, offset int) ([]model.Roadmap, int64, error) {
	var roadmaps []model.Roadmap
	var total int64

	r.db.Model(&model.Roadmap{}).Where("is_public = ?", true).Count(&total)

	err := r.db.Preload("User").
		Where("is_public = ?", true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&roadmaps).Error

	return roadmaps, total, err
}

// CopyRoadmap creates a copy of an existing roadmap for a new user
func (r *RoadmapRepository) CopyRoadmap(originalID, newUserID uint) (*model.Roadmap, error) {
	original, err := r.FindByID(originalID)
	if err != nil {
		return nil, err
	}

	newRoadmap := &model.Roadmap{
		UserID:      newUserID,
		Title:       original.Title + " (コピー)",
		Description: original.Description,
		Category:    original.Category,
		IsPublic:    false,
		StepCount:   original.StepCount,
		Status:      model.RoadmapStatusActive,
	}

	if err := r.db.Create(newRoadmap).Error; err != nil {
		return nil, err
	}

	for _, step := range original.Steps {
		newStep := model.RoadmapStep{
			RoadmapID:   newRoadmap.ID,
			Title:       step.Title,
			Description: step.Description,
			OrderIndex:  step.OrderIndex,
			ResourceURL: step.ResourceURL,
		}
		if err := r.db.Create(&newStep).Error; err != nil {
			return nil, err
		}
	}

	return r.FindByID(newRoadmap.ID)
}

// GetStats gets roadmap statistics for a user
func (r *RoadmapRepository) GetStats(userID uint) (*model.RoadmapStats, error) {
	var stats model.RoadmapStats

	var totalCount int64
	r.db.Model(&model.Roadmap{}).Where("user_id = ?", userID).Count(&totalCount)
	stats.TotalRoadmaps = int(totalCount)

	var activeCount int64
	r.db.Model(&model.Roadmap{}).Where("user_id = ? AND status = ?", userID, model.RoadmapStatusActive).Count(&activeCount)
	stats.ActiveRoadmaps = int(activeCount)

	var completedCount int64
	r.db.Model(&model.Roadmap{}).Where("user_id = ? AND status = ?", userID, model.RoadmapStatusCompleted).Count(&completedCount)
	stats.CompletedRoadmaps = int(completedCount)

	var totalSteps int64
	r.db.Model(&model.Roadmap{}).Where("user_id = ?", userID).Select("COALESCE(SUM(step_count), 0)").Scan(&totalSteps)
	stats.TotalSteps = int(totalSteps)

	var completedSteps int64
	r.db.Model(&model.Roadmap{}).Where("user_id = ?", userID).Select("COALESCE(SUM(completed_step_count), 0)").Scan(&completedSteps)
	stats.CompletedSteps = int(completedSteps)

	return &stats, nil
}

// === RoadmapStep CRUD ===

// CreateStep creates a new step and increments roadmap step_count
func (r *RoadmapRepository) CreateStep(step *model.RoadmapStep) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(step).Error; err != nil {
			return err
		}
		return tx.Model(&model.Roadmap{}).Where("id = ?", step.RoadmapID).
			UpdateColumn("step_count", gorm.Expr("step_count + 1")).Error
	})
}

// UpdateStep updates a step and recalculates roadmap progress if completion changed
func (r *RoadmapRepository) UpdateStep(step *model.RoadmapStep) error {
	oldStep := &model.RoadmapStep{}
	if err := r.db.First(oldStep, step.ID).Error; err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(step).Error; err != nil {
			return err
		}

		if oldStep.IsCompleted != step.IsCompleted {
			delta := 1
			if !step.IsCompleted {
				delta = -1
			}

			if err := tx.Model(&model.Roadmap{}).Where("id = ?", step.RoadmapID).
				UpdateColumn("completed_step_count", gorm.Expr("completed_step_count + ?", delta)).Error; err != nil {
				return err
			}

			// Recalculate progress
			var roadmap model.Roadmap
			if err := tx.First(&roadmap, step.RoadmapID).Error; err != nil {
				return err
			}

			progress := 0
			if roadmap.StepCount > 0 {
				progress = (roadmap.CompletedStepCount * 100) / roadmap.StepCount
			}

			updates := map[string]interface{}{"progress": progress}
			if progress == 100 && roadmap.Status == model.RoadmapStatusActive {
				now := time.Now()
				updates["status"] = model.RoadmapStatusCompleted
				updates["completed_at"] = now
			} else if progress < 100 && roadmap.Status == model.RoadmapStatusCompleted {
				updates["status"] = model.RoadmapStatusActive
				updates["completed_at"] = nil
			}

			return tx.Model(&roadmap).Updates(updates).Error
		}

		return nil
	})
}

// DeleteStep deletes a step and decrements roadmap step_count
func (r *RoadmapRepository) DeleteStep(stepID uint) error {
	step := &model.RoadmapStep{}
	if err := r.db.First(step, stepID).Error; err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(step).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Roadmap{}).Where("id = ?", step.RoadmapID).
			UpdateColumn("step_count", gorm.Expr("step_count - 1")).Error; err != nil {
			return err
		}

		if step.IsCompleted {
			if err := tx.Model(&model.Roadmap{}).Where("id = ?", step.RoadmapID).
				UpdateColumn("completed_step_count", gorm.Expr("completed_step_count - 1")).Error; err != nil {
				return err
			}
		}

		// Recalculate progress
		var roadmap model.Roadmap
		if err := tx.First(&roadmap, step.RoadmapID).Error; err != nil {
			return err
		}

		progress := 0
		if roadmap.StepCount > 0 {
			progress = (roadmap.CompletedStepCount * 100) / roadmap.StepCount
		}

		return tx.Model(&roadmap).Update("progress", progress).Error
	})
}

// FindStepByID finds a step by ID
func (r *RoadmapRepository) FindStepByID(stepID uint) (*model.RoadmapStep, error) {
	var step model.RoadmapStep
	err := r.db.First(&step, stepID).Error
	return &step, err
}

// ReorderSteps updates the order_index of multiple steps
func (r *RoadmapRepository) ReorderSteps(roadmapID uint, stepOrders []StepOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, order := range stepOrders {
			if err := tx.Model(&model.RoadmapStep{}).
				Where("id = ? AND roadmap_id = ?", order.StepID, roadmapID).
				Update("order_index", order.OrderIndex).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// StepOrder represents the order of a step for reordering
type StepOrder struct {
	StepID     uint `json:"step_id"`
	OrderIndex int  `json:"order_index"`
}

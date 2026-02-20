package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// RoadmapRepository は学習ロードマップデータへのアクセスを提供するリポジトリ実装。
type RoadmapRepository struct {
	db *gorm.DB
}

// NewRoadmapRepository は新しいRoadmapRepositoryインスタンスを生成する。
func NewRoadmapRepository(db *gorm.DB) *RoadmapRepository {
	return &RoadmapRepository{db: db}
}

// === ロードマップCRUD ===

// Create は新しいロードマップをデータベースに作成する。
func (r *RoadmapRepository) Create(roadmap *model.Roadmap) error {
	return r.db.Create(roadmap).Error
}

// Update は既存のロードマップを更新する。
func (r *RoadmapRepository) Update(roadmap *model.Roadmap) error {
	return r.db.Save(roadmap).Error
}

// Delete は指定IDのロードマップを削除する（ステップはCASCADE削除される）。
func (r *RoadmapRepository) Delete(id uint) error {
	return r.db.Delete(&model.Roadmap{}, id).Error
}

// FindByID は指定IDのロードマップをステップ・ユーザー情報付きで取得する。
// ステップはorder_index昇順でソートされる。
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

// GetByUserID は指定ユーザーの全ロードマップを取得する（ステップなし、新しい順）。
func (r *RoadmapRepository) GetByUserID(userID uint, limit, offset int) ([]model.Roadmap, int64, error) {
	var roadmaps []model.Roadmap
	var total int64
	query := r.db.Where("user_id = ?", userID)
	query.Model(&model.Roadmap{}).Count(&total)
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&roadmaps).Error
	return roadmaps, total, err
}

// GetByStatus は指定ユーザーのロードマップをステータスでフィルタリングして取得する（新しい順）。
func (r *RoadmapRepository) GetByStatus(userID uint, status string) ([]model.Roadmap, error) {
	var roadmaps []model.Roadmap
	err := r.db.Where("user_id = ? AND status = ?", userID, status).
		Order("created_at DESC").
		Find(&roadmaps).Error
	return roadmaps, err
}

// GetPublicRoadmaps は公開ロードマップをページネーション付きで取得する。
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

// CopyRoadmap は既存のロードマップを新しいユーザー用にコピーする。
// ステップもコピーされるが、完了状態はリセットされる。タイトルに「(コピー)」が付与される。
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

// GetStats は指定ユーザーのロードマップ統計情報を算出する。
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

// === ロードマップステップCRUD ===

// CreateStep は新しいステップを作成し、ロードマップのstep_countをインクリメントする。
func (r *RoadmapRepository) CreateStep(step *model.RoadmapStep) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(step).Error; err != nil {
			return err
		}
		return tx.Model(&model.Roadmap{}).Where("id = ?", step.RoadmapID).
			UpdateColumn("step_count", gorm.Expr("step_count + 1")).Error
	})
}

// UpdateStep はステップを更新する。
// 完了状態が変更された場合、ロードマップのcompleted_step_count・progress・statusも再計算する。
func (r *RoadmapRepository) UpdateStep(step *model.RoadmapStep) error {
	oldStep := &model.RoadmapStep{}
	if err := r.db.First(oldStep, step.ID).Error; err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(step).Error; err != nil {
			return err
		}

		// 完了状態が変更された場合、ロードマップの進捗を再計算
		if oldStep.IsCompleted != step.IsCompleted {
			delta := 1
			if !step.IsCompleted {
				delta = -1
			}

			if err := tx.Model(&model.Roadmap{}).Where("id = ?", step.RoadmapID).
				UpdateColumn("completed_step_count", gorm.Expr("completed_step_count + ?", delta)).Error; err != nil {
				return err
			}

			// 進捗率を再計算
			var roadmap model.Roadmap
			if err := tx.First(&roadmap, step.RoadmapID).Error; err != nil {
				return err
			}

			progress := 0
			if roadmap.StepCount > 0 {
				progress = (roadmap.CompletedStepCount * 100) / roadmap.StepCount
			}

			updates := map[string]interface{}{"progress": progress}
			// 100%達成時に自動完了、100%未満時に自動アクティブ化
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

// DeleteStep はステップを削除し、ロードマップのstep_count・completed_step_count・progressを再計算する。
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

		// 完了済みステップの場合、completed_step_countもデクリメント
		if step.IsCompleted {
			if err := tx.Model(&model.Roadmap{}).Where("id = ?", step.RoadmapID).
				UpdateColumn("completed_step_count", gorm.Expr("completed_step_count - 1")).Error; err != nil {
				return err
			}
		}

		// 進捗率を再計算
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

// FindStepByID は指定IDのステップを取得する。
func (r *RoadmapRepository) FindStepByID(stepID uint) (*model.RoadmapStep, error) {
	var step model.RoadmapStep
	err := r.db.First(&step, stepID).Error
	return &step, err
}

// ReorderSteps は複数ステップの表示順序を一括で更新する。
func (r *RoadmapRepository) ReorderSteps(roadmapID uint, stepOrders []model.StepOrder) error {
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

// GetTemplates はテンプレートとしてマークされた全ロードマップをステップ付きで取得する。
func (r *RoadmapRepository) GetTemplates() ([]model.Roadmap, error) {
	var templates []model.Roadmap
	err := r.db.Where("is_template = ?", true).
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Order("created_at ASC").
		Find(&templates).Error
	return templates, err
}


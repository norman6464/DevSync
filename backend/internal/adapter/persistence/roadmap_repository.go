package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// roadmapRepository は [repository.RoadmapRepository] の GORM 実装。
//
// 旧 repository パッケージにも同じテーブルを扱う実装が残っている。ai_advice がまだそちらを
// 使っているため、移行が一巡するまで新旧のアダプタが並存する。
type roadmapRepository struct {
	db *gorm.DB
}

// NewRoadmapRepository は RoadmapRepository の GORM 実装を返す。
func NewRoadmapRepository(db *gorm.DB) repository.RoadmapRepository {
	return &roadmapRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.RoadmapRepository = (*roadmapRepository)(nil)

// Create は新しいロードマップを作成する。
func (r *roadmapRepository) Create(ctx context.Context, roadmap *model.Roadmap) error {
	return r.db.WithContext(ctx).Create(roadmap).Error
}

// Update は既存のロードマップを更新する。
func (r *roadmapRepository) Update(ctx context.Context, roadmap *model.Roadmap) error {
	return r.db.WithContext(ctx).Save(roadmap).Error
}

// Delete はロードマップを削除する（ステップは CASCADE で消える）。
func (r *roadmapRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Roadmap{}, id).Error
}

// FindByID はステップ（表示順）とユーザーを含めてロードマップを取得する。
// 不在の場合は (nil, nil) を返す。
func (r *roadmapRepository) FindByID(ctx context.Context, id uint) (*model.Roadmap, error) {
	var roadmap model.Roadmap
	err := r.db.WithContext(ctx).
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Preload("User").
		First(&roadmap, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &roadmap, nil
}

// GetByUserID は指定ユーザーのロードマップをページネーション付きで取得する（新しい順・ステップなし）。
func (r *roadmapRepository) GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Roadmap, int64, error) {
	scope := r.db.WithContext(ctx).Model(&model.Roadmap{}).Where("user_id = ?", userID)

	var total int64
	if err := scope.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var roadmaps []model.Roadmap
	err := scope.Session(&gorm.Session{}).
		Order("created_at DESC").Limit(limit).Offset(offset).
		Find(&roadmaps).Error
	return roadmaps, total, err
}

// GetByStatus は指定ユーザーのロードマップをステータスで絞り込んで取得する（新しい順）。
func (r *roadmapRepository) GetByStatus(ctx context.Context, userID uint, status string) ([]model.Roadmap, error) {
	var roadmaps []model.Roadmap
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, status).
		Order("created_at DESC").
		Find(&roadmaps).Error
	return roadmaps, err
}

// GetPublicRoadmaps は公開ロードマップをページネーション付きで取得する（新しい順）。
func (r *roadmapRepository) GetPublicRoadmaps(ctx context.Context, limit, offset int) ([]model.Roadmap, int64, error) {
	scope := r.db.WithContext(ctx).Model(&model.Roadmap{}).Where("is_public = ?", true)

	var total int64
	if err := scope.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var roadmaps []model.Roadmap
	err := scope.Session(&gorm.Session{}).
		Preload("User").
		Order("created_at DESC").Limit(limit).Offset(offset).
		Find(&roadmaps).Error
	return roadmaps, total, err
}

// GetTemplates はテンプレートのロードマップをステップ付きで取得する（古い順）。
func (r *roadmapRepository) GetTemplates(ctx context.Context) ([]model.Roadmap, error) {
	var templates []model.Roadmap
	err := r.db.WithContext(ctx).
		Where("is_template = ?", true).
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Order("created_at ASC").
		Find(&templates).Error
	return templates, err
}

// CountByUserID は指定ユーザーのロードマップ総数を返す。
func (r *roadmapRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Roadmap{}).
		Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CopyRoadmap は元のロードマップとステップを複製する。複製は非公開・アクティブで作られる。
func (r *roadmapRepository) CopyRoadmap(ctx context.Context, originalID, newUserID uint) (*model.Roadmap, error) {
	original, err := r.FindByID(ctx, originalID)
	if err != nil {
		return nil, err
	}
	if original == nil {
		return nil, gorm.ErrRecordNotFound
	}

	db := r.db.WithContext(ctx)
	newRoadmap := &model.Roadmap{
		UserID:      newUserID,
		Title:       original.Title + " (コピー)",
		Description: original.Description,
		Category:    original.Category,
		IsPublic:    false,
		StepCount:   original.StepCount,
		Status:      model.RoadmapStatusActive,
	}
	if err := db.Create(newRoadmap).Error; err != nil {
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
		if err := db.Create(&newStep).Error; err != nil {
			return nil, err
		}
	}

	return r.FindByID(ctx, newRoadmap.ID)
}

// CreateStep はステップを追加し、ロードマップのステップ数を 1 増やす。
func (r *roadmapRepository) CreateStep(ctx context.Context, step *model.RoadmapStep) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(step).Error; err != nil {
			return err
		}
		return tx.Model(&model.Roadmap{}).Where("id = ?", step.RoadmapID).
			UpdateColumn("step_count", gorm.Expr("step_count + 1")).Error
	})
}

// UpdateStep はステップを更新する。完了状態が変わった場合は進捗も再計算する。
func (r *roadmapRepository) UpdateStep(ctx context.Context, step *model.RoadmapStep) error {
	db := r.db.WithContext(ctx)

	var oldStep model.RoadmapStep
	if err := db.First(&oldStep, step.ID).Error; err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(step).Error; err != nil {
			return err
		}
		if oldStep.IsCompleted == step.IsCompleted {
			return nil
		}

		delta := 1
		if !step.IsCompleted {
			delta = -1
		}
		if err := tx.Model(&model.Roadmap{}).Where("id = ?", step.RoadmapID).
			UpdateColumn("completed_step_count", gorm.Expr("completed_step_count + ?", delta)).Error; err != nil {
			return err
		}
		return recalcRoadmapProgress(tx, step.RoadmapID, true)
	})
}

// DeleteStep はステップを削除し、ステップ数・完了ステップ数・進捗率を再計算する。
func (r *roadmapRepository) DeleteStep(ctx context.Context, stepID uint) error {
	db := r.db.WithContext(ctx)

	var step model.RoadmapStep
	if err := db.First(&step, stepID).Error; err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&step).Error; err != nil {
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
		return recalcRoadmapProgress(tx, step.RoadmapID, false)
	})
}

// recalcRoadmapProgress は進捗率を再計算して保存する。
// withStatus が true のときは 100% 到達での自動完了と、100% 未満へ戻ったときの
// アクティブ復帰も行う（ステップ更新のみが対象で、削除時は進捗率だけを更新する）。
func recalcRoadmapProgress(tx *gorm.DB, roadmapID uint, withStatus bool) error {
	var roadmap model.Roadmap
	if err := tx.First(&roadmap, roadmapID).Error; err != nil {
		return err
	}

	progress := 0
	if roadmap.StepCount > 0 {
		progress = (roadmap.CompletedStepCount * 100) / roadmap.StepCount
	}
	if !withStatus {
		return tx.Model(&roadmap).Update("progress", progress).Error
	}

	updates := map[string]interface{}{"progress": progress}
	switch {
	case progress == 100 && roadmap.Status == model.RoadmapStatusActive:
		updates["status"] = model.RoadmapStatusCompleted
		updates["completed_at"] = time.Now()
	case progress < 100 && roadmap.Status == model.RoadmapStatusCompleted:
		updates["status"] = model.RoadmapStatusActive
		updates["completed_at"] = nil
	}
	return tx.Model(&roadmap).Updates(updates).Error
}

// FindStepByID は指定 ID のステップを取得する。不在の場合は (nil, nil) を返す。
func (r *roadmapRepository) FindStepByID(ctx context.Context, stepID uint) (*model.RoadmapStep, error) {
	var step model.RoadmapStep
	err := r.db.WithContext(ctx).First(&step, stepID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &step, nil
}

// ReorderSteps はステップの表示順序をまとめて更新する。
func (r *roadmapRepository) ReorderSteps(ctx context.Context, roadmapID uint, stepOrders []model.StepOrder) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

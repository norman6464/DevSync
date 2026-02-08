package service

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// StepOrder はステップの並び替え情報を表す型エイリアス。
type StepOrder = repository.StepOrder

// RoadmapService は学習ロードマップのビジネスロジックを提供する。
// ロードマップとステップのCRUD操作、可視性制御、コピー機能を担当する。
type RoadmapService struct {
	repo repository.RoadmapRepositoryInterface
}

// NewRoadmapService は新しいRoadmapServiceインスタンスを生成する。
func NewRoadmapService(repo repository.RoadmapRepositoryInterface) *RoadmapService {
	return &RoadmapService{repo: repo}
}

// Create は新しいロードマップを作成する。
func (s *RoadmapService) Create(roadmap *model.Roadmap) error {
	return s.repo.Create(roadmap)
}

// GetByID は指定IDのロードマップを可視性チェック付きで取得する。
// 非公開ロードマップはオーナー以外アクセスできない。
func (s *RoadmapService) GetByID(id, userID uint) (*model.Roadmap, error) {
	roadmap, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if roadmap.UserID != userID && !roadmap.IsPublic {
		return nil, ErrForbidden
	}
	return roadmap, nil
}

// GetByUserID は指定ユーザーの全ロードマップを取得する。
func (s *RoadmapService) GetByUserID(userID uint) ([]model.Roadmap, error) {
	return s.repo.GetByUserID(userID)
}

// GetPublicRoadmaps は公開ロードマップをページネーション付きで取得する。
func (s *RoadmapService) GetPublicRoadmaps(limit, offset int) ([]model.Roadmap, int64, error) {
	return s.repo.GetPublicRoadmaps(limit, offset)
}

// GetStats は指定ユーザーのロードマップ統計情報を取得する。
func (s *RoadmapService) GetStats(userID uint) (*model.RoadmapStats, error) {
	return s.repo.GetStats(userID)
}

// Update は所有権を検証した後、ロードマップを更新する。
// ステータスが「完了」に変更された場合、完了日時を自動設定する。
func (s *RoadmapService) Update(id, userID uint, updates *model.Roadmap) (*model.Roadmap, error) {
	roadmap, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if roadmap.UserID != userID {
		return nil, ErrForbidden
	}

	if updates.Title != "" {
		roadmap.Title = updates.Title
	}
	if updates.Description != "" {
		roadmap.Description = updates.Description
	}
	if updates.Category != "" {
		roadmap.Category = updates.Category
	}
	if updates.Status != "" {
		roadmap.Status = updates.Status
		if roadmap.Status == model.RoadmapStatusCompleted && roadmap.CompletedAt == nil {
			now := time.Now()
			roadmap.CompletedAt = &now
		}
	}

	if err := s.repo.Update(roadmap); err != nil {
		return nil, err
	}
	return roadmap, nil
}

// UpdateVisibility は所有権を検証した後、ロードマップの公開/非公開状態を更新する。
func (s *RoadmapService) UpdateVisibility(id, userID uint, isPublic bool) (*model.Roadmap, error) {
	roadmap, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if roadmap.UserID != userID {
		return nil, ErrForbidden
	}

	roadmap.IsPublic = isPublic

	if err := s.repo.Update(roadmap); err != nil {
		return nil, err
	}
	return roadmap, nil
}

// Delete は所有権を検証した後、ロードマップを削除する。
func (s *RoadmapService) Delete(id, userID uint) error {
	roadmap, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if roadmap.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}

// CopyRoadmap は公開ロードマップをテンプレートとしてコピーする。
// 非公開かつ自分のものでないロードマップはコピーできない。
func (s *RoadmapService) CopyRoadmap(roadmapID, userID uint) (*model.Roadmap, error) {
	original, err := s.repo.FindByID(roadmapID)
	if err != nil {
		return nil, err
	}
	if !original.IsPublic && original.UserID != userID {
		return nil, ErrForbidden
	}
	return s.repo.CopyRoadmap(roadmapID, userID)
}

// CreateStep は所有権を検証した後、ロードマップにステップを追加する。
func (s *RoadmapService) CreateStep(roadmapID, userID uint, step *model.RoadmapStep) error {
	roadmap, err := s.repo.FindByID(roadmapID)
	if err != nil {
		return err
	}
	if roadmap.UserID != userID {
		return ErrForbidden
	}

	step.RoadmapID = roadmapID
	return s.repo.CreateStep(step)
}

// UpdateStep はロードマップの所有権とステップの所属を検証した後、ステップを更新する。
func (s *RoadmapService) UpdateStep(roadmapID, stepID, userID uint, updates *model.RoadmapStep) (*model.RoadmapStep, error) {
	roadmap, err := s.repo.FindByID(roadmapID)
	if err != nil {
		return nil, err
	}
	if roadmap.UserID != userID {
		return nil, ErrForbidden
	}

	step, err := s.repo.FindStepByID(stepID)
	if err != nil {
		return nil, err
	}
	if step.RoadmapID != roadmapID {
		return nil, ErrBadRequest
	}

	if updates.Title != "" {
		step.Title = updates.Title
	}
	if updates.Description != "" {
		step.Description = updates.Description
	}
	if updates.ResourceURL != "" {
		step.ResourceURL = updates.ResourceURL
	}

	if err := s.repo.UpdateStep(step); err != nil {
		return nil, err
	}
	return step, nil
}

// UpdateStepCompletion はステップの完了状態を更新する。
// 完了時にはCompletedAtを設定し、未完了に戻す場合はnilにリセットする。
func (s *RoadmapService) UpdateStepCompletion(roadmapID, stepID, userID uint, isCompleted bool) (*model.RoadmapStep, error) {
	roadmap, err := s.repo.FindByID(roadmapID)
	if err != nil {
		return nil, err
	}
	if roadmap.UserID != userID {
		return nil, ErrForbidden
	}

	step, err := s.repo.FindStepByID(stepID)
	if err != nil {
		return nil, err
	}
	if step.RoadmapID != roadmapID {
		return nil, ErrBadRequest
	}

	step.IsCompleted = isCompleted
	if isCompleted && step.CompletedAt == nil {
		now := time.Now()
		step.CompletedAt = &now
	} else if !isCompleted {
		step.CompletedAt = nil
	}

	if err := s.repo.UpdateStep(step); err != nil {
		return nil, err
	}
	return step, nil
}

// DeleteStep はロードマップの所有権とステップの所属を検証した後、ステップを削除する。
func (s *RoadmapService) DeleteStep(roadmapID, stepID, userID uint) error {
	roadmap, err := s.repo.FindByID(roadmapID)
	if err != nil {
		return err
	}
	if roadmap.UserID != userID {
		return ErrForbidden
	}

	step, err := s.repo.FindStepByID(stepID)
	if err != nil {
		return err
	}
	if step.RoadmapID != roadmapID {
		return ErrBadRequest
	}

	return s.repo.DeleteStep(stepID)
}

// ReorderSteps は所有権を検証した後、ステップの表示順序を一括更新する。
func (s *RoadmapService) ReorderSteps(roadmapID, userID uint, orders []repository.StepOrder) error {
	roadmap, err := s.repo.FindByID(roadmapID)
	if err != nil {
		return err
	}
	if roadmap.UserID != userID {
		return ErrForbidden
	}
	return s.repo.ReorderSteps(roadmapID, orders)
}

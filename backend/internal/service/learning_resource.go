package service

import (
	"github.com/norman6464/devsync/backend/internal/domain/validator"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// LearningResourceService は学習リソースのビジネスロジックを提供する。
// リソースのCRUD操作、公開/非公開制御、いいね・保存機能を担当する。
type LearningResourceService struct {
	repo repository.LearningResourceRepositoryInterface
}

// NewLearningResourceService は新しいLearningResourceServiceインスタンスを生成する。
func NewLearningResourceService(repo repository.LearningResourceRepositoryInterface) *LearningResourceService {
	return &LearningResourceService{repo: repo}
}

// Create は新しい学習リソースを作成する。
func (s *LearningResourceService) Create(resource *model.LearningResource) error {
	v := validator.NewResourceValidator()
	if err := v.ValidateCreateResource(resource.Title, resource.Description, resource.URL, string(resource.Category), string(resource.Difficulty)); err != nil {
		return err
	}
	return s.repo.Create(resource)
}

// GetByID は指定IDの学習リソースを可視性チェック付きで取得する。
// 非公開リソースはオーナー以外アクセスできない。
func (s *LearningResourceService) GetByID(id, userID uint) (*model.LearningResource, error) {
	resource, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 非公開リソースで、かつオーナーでない場合はアクセス拒否
	if !resource.IsPublic && resource.UserID != userID {
		return nil, ErrForbidden
	}

	return resource, nil
}

// HasLiked は指定ユーザーがリソースにいいね済みかを判定する。
func (s *LearningResourceService) HasLiked(userID, resourceID uint) (bool, error) {
	return s.repo.HasLiked(userID, resourceID)
}

// HasSaved は指定ユーザーがリソースを保存済みかを判定する。
func (s *LearningResourceService) HasSaved(userID, resourceID uint) (bool, error) {
	return s.repo.HasSaved(userID, resourceID)
}

// GetByUserID は指定ユーザーの学習リソースを取得する。
// 自分のリソースの場合は非公開も含め、他ユーザーの場合は公開のみ返す。
func (s *LearningResourceService) GetByUserID(targetUserID, currentUserID uint) ([]model.LearningResource, error) {
	includePrivate := currentUserID == targetUserID
	return s.repo.FindByUserID(targetUserID, includePrivate)
}

// GetPublic は公開学習リソースをページネーション・フィルタ付きで取得する。
func (s *LearningResourceService) GetPublic(limit, offset int, category, difficulty string) ([]model.LearningResource, int64, error) {
	return s.repo.FindPublic(limit, offset, category, difficulty)
}

// Search は学習リソースをキーワードで検索する。
func (s *LearningResourceService) Search(query string, limit, offset int) ([]model.LearningResource, int64, error) {
	return s.repo.Search(query, limit, offset)
}

// Update は所有権を検証した後、学習リソースを更新する。
func (s *LearningResourceService) Update(id, userID uint, updates *model.LearningResource) (*model.LearningResource, error) {
	resource, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if resource.UserID != userID {
		return nil, ErrForbidden
	}

	v := validator.NewResourceValidator()
	if err := v.ValidateUpdateResource(updates.Title, updates.Description, updates.URL, string(updates.Category), string(updates.Difficulty)); err != nil {
		return nil, err
	}

	if updates.Title != "" {
		resource.Title = updates.Title
	}
	if updates.Description != "" {
		resource.Description = updates.Description
	}
	if updates.URL != "" {
		resource.URL = updates.URL
	}
	if updates.Category != "" {
		resource.Category = updates.Category
	}
	if updates.Difficulty != "" {
		resource.Difficulty = updates.Difficulty
	}
	if updates.Tags != "" {
		resource.Tags = updates.Tags
	}
	if updates.ImageURL != "" {
		resource.ImageURL = updates.ImageURL
	}

	if err := s.repo.Update(resource); err != nil {
		return nil, err
	}
	return resource, nil
}

// UpdateVisibility は所有権を検証した後、リソースの公開/非公開状態を更新する。
func (s *LearningResourceService) UpdateVisibility(id, userID uint, isPublic bool) (*model.LearningResource, error) {
	resource, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if resource.UserID != userID {
		return nil, ErrForbidden
	}

	resource.IsPublic = isPublic

	if err := s.repo.Update(resource); err != nil {
		return nil, err
	}
	return resource, nil
}

// Delete は所有権を検証した後、学習リソースを削除する。
func (s *LearningResourceService) Delete(id, userID uint) error {
	resource, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if resource.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}

// Like は学習リソースにいいねを追加する。
func (s *LearningResourceService) Like(userID, resourceID uint) error {
	return s.repo.Like(userID, resourceID)
}

// Unlike は学習リソースのいいねを取り消す。
func (s *LearningResourceService) Unlike(userID, resourceID uint) error {
	return s.repo.Unlike(userID, resourceID)
}

// Save は学習リソースを保存する。
func (s *LearningResourceService) Save(userID, resourceID uint) error {
	return s.repo.Save(userID, resourceID)
}

// Unsave は学習リソースの保存を取り消す。
func (s *LearningResourceService) Unsave(userID, resourceID uint) error {
	return s.repo.Unsave(userID, resourceID)
}

// GetSavedByUserID は指定ユーザーの保存済みリソースをページネーション付きで取得する。
func (s *LearningResourceService) GetSavedByUserID(userID uint, limit, offset int) ([]model.LearningResource, int64, error) {
	return s.repo.FindSavedByUserID(userID, limit, offset)
}

package service

import (
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
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
	if err := domain.ValidateStringLength(resource.Tags, 0, 1000, "タグ"); err != nil {
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

// GetByUserID は指定ユーザーの学習リソースをページネーション付きで取得する。
// 自分のリソースの場合は非公開も含め、他ユーザーの場合は公開のみ返す。
func (s *LearningResourceService) GetByUserID(targetUserID, currentUserID uint, limit, offset int) ([]model.LearningResource, int64, error) {
	includePrivate := currentUserID == targetUserID
	return s.repo.FindByUserID(targetUserID, includePrivate, limit, offset)
}

// GetPublic は公開学習リソースをページネーション・フィルタ付きで取得する。
func (s *LearningResourceService) GetPublic(limit, offset int, category, difficulty string) ([]model.LearningResource, int64, error) {
	return s.repo.FindPublic(limit, offset, category, difficulty)
}

// validDifficulties は有効な難易度の一覧。
var validDifficulties = map[string]bool{
	string(model.ResourceDifficultyBeginner):     true,
	string(model.ResourceDifficultyIntermediate): true,
	string(model.ResourceDifficultyAdvanced):     true,
}

// GetByDifficulty は公開学習リソースを難易度でフィルタリングして取得する。
func (s *LearningResourceService) GetByDifficulty(difficulty string, limit, offset int) ([]model.LearningResource, int64, error) {
	if !validDifficulties[difficulty] {
		return nil, 0, domain.NewError(domain.ErrCodeBadRequest, "無効な難易度です", nil)
	}
	return s.repo.FindByDifficulty(difficulty, limit, offset)
}

// Search は学習リソースをキーワードで検索する。
func (s *LearningResourceService) Search(query string, limit, offset int) ([]model.LearningResource, int64, error) {
	return s.repo.Search(query, limit, offset)
}

// findAndCheckOwnership はリソースを取得し、指定ユーザーが所有者かを検証する。
func (s *LearningResourceService) findAndCheckOwnership(id, userID uint) (*model.LearningResource, error) {
	return checkOwnership(s.repo.FindByID, id, userID, func(r *model.LearningResource) uint { return r.UserID })
}

// Update は所有権を検証した後、学習リソースを更新する。
func (s *LearningResourceService) Update(id, userID uint, updates *model.LearningResource) (*model.LearningResource, error) {
	resource, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	// 空白のみの値を空文字列に正規化（空白バイパス防止）
	updates.Title = strings.TrimSpace(updates.Title)
	updates.Description = strings.TrimSpace(updates.Description)
	updates.URL = strings.TrimSpace(updates.URL)
	updates.Tags = strings.TrimSpace(updates.Tags)
	updates.ImageURL = strings.TrimSpace(updates.ImageURL)
	trimmedCategory := strings.TrimSpace(string(updates.Category))
	trimmedDifficulty := strings.TrimSpace(string(updates.Difficulty))

	v := validator.NewResourceValidator()
	if err := v.ValidateUpdateResource(updates.Title, updates.Description, updates.URL, trimmedCategory, trimmedDifficulty); err != nil {
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
	if trimmedCategory != "" {
		updates.Category = model.ResourceCategory(trimmedCategory)
		resource.Category = updates.Category
	}
	if trimmedDifficulty != "" {
		updates.Difficulty = model.ResourceDifficulty(trimmedDifficulty)
		resource.Difficulty = updates.Difficulty
	}
	if updates.Tags != "" {
		if err := domain.ValidateStringLength(updates.Tags, 1, 1000, "タグ"); err != nil {
			return nil, err
		}
		resource.Tags = updates.Tags
	}
	if updates.ImageURL != "" {
		if err := domain.ValidateStringLength(updates.ImageURL, 1, 2000, "画像URL"); err != nil {
			return nil, err
		}
		resource.ImageURL = updates.ImageURL
	}

	if err := s.repo.Update(resource); err != nil {
		return nil, err
	}
	return resource, nil
}

// UpdateVisibility は所有権を検証した後、リソースの公開/非公開状態を更新する。
func (s *LearningResourceService) UpdateVisibility(id, userID uint, isPublic bool) (*model.LearningResource, error) {
	resource, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	resource.IsPublic = isPublic

	if err := s.repo.Update(resource); err != nil {
		return nil, err
	}
	return resource, nil
}

// Delete は所有権を検証した後、学習リソースを削除する。
func (s *LearningResourceService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// findAndPreventSelfAction はリソースを取得し、自分のリソースへの操作を防止する。
func (s *LearningResourceService) findAndPreventSelfAction(userID, resourceID uint) error {
	resource, err := s.repo.FindByID(resourceID)
	if err != nil {
		return ErrNotFound
	}
	if resource.UserID == userID {
		return ErrForbidden
	}
	return nil
}

// Like は学習リソースにいいねを追加する。
// 自分のリソースにはいいねできない。
func (s *LearningResourceService) Like(userID, resourceID uint) error {
	if err := s.findAndPreventSelfAction(userID, resourceID); err != nil {
		return err
	}
	return s.repo.Like(userID, resourceID)
}

// Unlike は学習リソースのいいねを取り消す。
// 自分のリソースのいいねは取り消せない（そもそもいいねできないため）。
func (s *LearningResourceService) Unlike(userID, resourceID uint) error {
	if err := s.findAndPreventSelfAction(userID, resourceID); err != nil {
		return err
	}
	return s.repo.Unlike(userID, resourceID)
}

// Save は学習リソースを保存する。
// 自分のリソースは保存できない。
func (s *LearningResourceService) Save(userID, resourceID uint) error {
	if err := s.findAndPreventSelfAction(userID, resourceID); err != nil {
		return err
	}
	return s.repo.Save(userID, resourceID)
}

// Unsave は学習リソースの保存を取り消す。
// 自分のリソースの保存は取り消せない（そもそも保存できないため）。
func (s *LearningResourceService) Unsave(userID, resourceID uint) error {
	if err := s.findAndPreventSelfAction(userID, resourceID); err != nil {
		return err
	}
	return s.repo.Unsave(userID, resourceID)
}

// GetSavedByUserID は指定ユーザーの保存済みリソースをページネーション付きで取得する。
func (s *LearningResourceService) GetSavedByUserID(userID uint, limit, offset int) ([]model.LearningResource, int64, error) {
	return s.repo.FindSavedByUserID(userID, limit, offset)
}

// CountByUserID は指定ユーザーの学習リソース総数を返す。
func (s *LearningResourceService) CountByUserID(userID uint) (int64, error) {
	return s.repo.CountByUserID(userID)
}

package service

import (
	"github.com/norman6464/devsync/backend/internal/domain/validator"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// NoteFolderService はノートフォルダのビジネスロジックを提供する。
type NoteFolderService struct {
	repo repository.NoteFolderRepositoryInterface
}

// NewNoteFolderService は新しいNoteFolderServiceインスタンスを生成する。
func NewNoteFolderService(repo repository.NoteFolderRepositoryInterface) *NoteFolderService {
	return &NoteFolderService{repo: repo}
}

// Create は新しいフォルダを作成する。
func (s *NoteFolderService) Create(folder *model.NoteFolder) error {
	// バリデーション
	v := validator.NewNoteFolderValidator()
	if err := v.ValidateCreate(folder.Name); err != nil {
		return err
	}

	return s.repo.Create(folder)
}

// GetByID は指定IDのフォルダを取得する。
func (s *NoteFolderService) GetByID(id uint) (*model.NoteFolder, error) {
	return s.repo.FindByID(id)
}

// GetByUserID は指定ユーザーの全フォルダを取得する。
func (s *NoteFolderService) GetByUserID(userID uint) ([]model.NoteFolder, error) {
	return s.repo.FindByUserID(userID)
}

// GetChildren は指定フォルダの子フォルダを取得する。
func (s *NoteFolderService) GetChildren(parentID uint) ([]model.NoteFolder, error) {
	return s.repo.FindByParentID(parentID)
}

// GetRootFolders は指定ユーザーのルートフォルダを取得する。
func (s *NoteFolderService) GetRootFolders(userID uint) ([]model.NoteFolder, error) {
	return s.repo.GetRootFolders(userID)
}

// Update は所有権を検証した後、フォルダを更新する。
func (s *NoteFolderService) Update(id, userID uint, name string, parentID *uint) (*model.NoteFolder, error) {
	folder, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if folder.UserID != userID {
		return nil, ErrForbidden
	}

	if name != "" {
		folder.Name = name
	}
	if parentID != nil {
		folder.ParentID = parentID
	}

	v := validator.NewNoteFolderValidator()
	if err := v.ValidateUpdate(folder.Name); err != nil {
		return nil, err
	}

	if err := s.repo.Update(folder); err != nil {
		return nil, err
	}
	return folder, nil
}

// Delete は所有権を検証した後、フォルダを削除する。
func (s *NoteFolderService) Delete(id, userID uint) error {
	folder, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if folder.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}

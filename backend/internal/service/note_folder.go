package service

import (
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
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

// GetByUserID は指定ユーザーのフォルダをページネーション付きで取得する。
func (s *NoteFolderService) GetByUserID(userID uint, limit, offset int) ([]model.NoteFolder, int64, error) {
	return s.repo.FindByUserID(userID, limit, offset)
}

// GetChildren は指定フォルダの子フォルダを取得する。
func (s *NoteFolderService) GetChildren(parentID uint) ([]model.NoteFolder, error) {
	return s.repo.FindByParentID(parentID)
}

// GetRootFolders は指定ユーザーのルートフォルダを取得する。
func (s *NoteFolderService) GetRootFolders(userID uint) ([]model.NoteFolder, error) {
	return s.repo.GetRootFolders(userID)
}

// findAndCheckOwnership はフォルダを取得し、指定ユーザーが所有者かを検証する。
func (s *NoteFolderService) findAndCheckOwnership(id, userID uint) (*model.NoteFolder, error) {
	return checkOwnership(s.repo.FindByID, id, userID, func(f *model.NoteFolder) uint { return f.UserID })
}

// isDescendant は targetID が ancestorID の子孫かを再帰的にチェックする。
func (s *NoteFolderService) isDescendant(ancestorID, targetID uint) (bool, error) {
	children, err := s.repo.FindByParentID(ancestorID)
	if err != nil {
		return false, err
	}
	for _, child := range children {
		if child.ID == targetID {
			return true, nil
		}
		desc, err := s.isDescendant(child.ID, targetID)
		if err != nil {
			return false, err
		}
		if desc {
			return true, nil
		}
	}
	return false, nil
}

// Update は所有権を検証した後、フォルダを更新する。
func (s *NoteFolderService) Update(id, userID uint, name string, parentID *uint) (*model.NoteFolder, error) {
	folder, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if name != "" {
		if strings.TrimSpace(name) == "" {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "フォルダ名は空白のみにできません", nil)
		}
		folder.Name = name
	}
	if parentID != nil {
		// 自己参照チェック
		if *parentID == id {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "フォルダを自分自身の子にすることはできません", nil)
		}
		// 循環参照チェック: parentIDが自分の子孫でないことを確認
		isDesc, err := s.isDescendant(id, *parentID)
		if err != nil {
			return nil, err
		}
		if isDesc {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "循環参照が発生するため、この親フォルダは設定できません", nil)
		}
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
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

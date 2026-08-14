package usecase

import (
	"context"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/domain/validator"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// noteFolderOwnerOf はノートフォルダの所有者 ID を返す。
func noteFolderOwnerOf(f *model.NoteFolder) uint { return f.UserID }

// CreateNoteFolderInput はフォルダ作成の入力。
type CreateNoteFolderInput struct {
	UserID   uint
	Name     string
	ParentID *uint
}

// CreateNoteFolderUseCase はノートフォルダを作成する。
type CreateNoteFolderUseCase struct {
	folders repository.NoteFolderRepository
}

// NewCreateNoteFolderUseCase は CreateNoteFolderUseCase を生成する。
func NewCreateNoteFolderUseCase(folders repository.NoteFolderRepository) *CreateNoteFolderUseCase {
	return &CreateNoteFolderUseCase{folders: folders}
}

// Execute はフォルダ名を検証したうえでフォルダを作成する。
func (uc *CreateNoteFolderUseCase) Execute(ctx context.Context, in CreateNoteFolderInput) (*model.NoteFolder, error) {
	if err := domain.ValidateRequiredID(in.UserID, "ユーザーID"); err != nil {
		return nil, err
	}
	if err := validator.NewNoteFolderValidator().ValidateCreate(in.Name); err != nil {
		return nil, err
	}

	folder := &model.NoteFolder{UserID: in.UserID, Name: in.Name, ParentID: in.ParentID}
	if err := uc.folders.Create(ctx, folder); err != nil {
		return nil, err
	}
	return folder, nil
}

// GetNoteFolderUseCase は指定 ID のノートフォルダを取得する。
type GetNoteFolderUseCase struct {
	folders repository.NoteFolderRepository
}

// NewGetNoteFolderUseCase は GetNoteFolderUseCase を生成する。
func NewGetNoteFolderUseCase(folders repository.NoteFolderRepository) *GetNoteFolderUseCase {
	return &GetNoteFolderUseCase{folders: folders}
}

// Execute は所有者本人のフォルダだけを返す。他ユーザーのフォルダには 403 を返す。
func (uc *GetNoteFolderUseCase) Execute(ctx context.Context, id, userID uint) (*model.NoteFolder, error) {
	return ensureOwner(ctx, uc.folders.FindByID, id, userID, noteFolderOwnerOf)
}

// ListNoteFoldersUseCase は指定ユーザーのフォルダ一覧をページネーション付きで取得する。
type ListNoteFoldersUseCase struct {
	folders repository.NoteFolderRepository
}

// NewListNoteFoldersUseCase は ListNoteFoldersUseCase を生成する。
func NewListNoteFoldersUseCase(folders repository.NoteFolderRepository) *ListNoteFoldersUseCase {
	return &ListNoteFoldersUseCase{folders: folders}
}

// Execute はフォルダ一覧と総件数を返す。
func (uc *ListNoteFoldersUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.NoteFolder, int64, error) {
	if err := domain.ValidateRequiredID(userID, "ユーザーID"); err != nil {
		return nil, 0, err
	}
	return uc.folders.FindByUserID(ctx, userID, limit, offset)
}

// ListChildNoteFoldersUseCase は指定フォルダ直下の子フォルダを取得する。
type ListChildNoteFoldersUseCase struct {
	folders repository.NoteFolderRepository
}

// NewListChildNoteFoldersUseCase は ListChildNoteFoldersUseCase を生成する。
func NewListChildNoteFoldersUseCase(folders repository.NoteFolderRepository) *ListChildNoteFoldersUseCase {
	return &ListChildNoteFoldersUseCase{folders: folders}
}

// Execute は所有者本人の子フォルダ一覧を返す。親が他ユーザーのものなら 403 を返す。
//
// 親の所有者を確認したうえで、返す子も本人のものだけに絞る。
// 親の検証が無かった頃に他ユーザーのフォルダを親に付け替えられたため、
// 既存データには持ち主の違う子が混ざり得る。
func (uc *ListChildNoteFoldersUseCase) Execute(ctx context.Context, parentID, userID uint) ([]model.NoteFolder, error) {
	if _, err := ensureOwner(ctx, uc.folders.FindByID, parentID, userID, noteFolderOwnerOf); err != nil {
		return nil, err
	}

	children, err := uc.folders.FindByParentID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	owned := make([]model.NoteFolder, 0, len(children))
	for _, child := range children {
		if child.UserID == userID {
			owned = append(owned, child)
		}
	}
	return owned, nil
}

// ListRootNoteFoldersUseCase は指定ユーザーのルートフォルダ（親なし）を取得する。
type ListRootNoteFoldersUseCase struct {
	folders repository.NoteFolderRepository
}

// NewListRootNoteFoldersUseCase は ListRootNoteFoldersUseCase を生成する。
func NewListRootNoteFoldersUseCase(folders repository.NoteFolderRepository) *ListRootNoteFoldersUseCase {
	return &ListRootNoteFoldersUseCase{folders: folders}
}

// Execute はルートフォルダ一覧を返す。
func (uc *ListRootNoteFoldersUseCase) Execute(ctx context.Context, userID uint) ([]model.NoteFolder, error) {
	if err := domain.ValidateRequiredID(userID, "ユーザーID"); err != nil {
		return nil, err
	}
	return uc.folders.FindRootsByUserID(ctx, userID)
}

// CountNoteFoldersUseCase は指定ユーザーのフォルダ総数を返す。
type CountNoteFoldersUseCase struct {
	folders repository.NoteFolderRepository
}

// NewCountNoteFoldersUseCase は CountNoteFoldersUseCase を生成する。
func NewCountNoteFoldersUseCase(folders repository.NoteFolderRepository) *CountNoteFoldersUseCase {
	return &CountNoteFoldersUseCase{folders: folders}
}

// Execute はフォルダ総数を返す。
func (uc *CountNoteFoldersUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	if err := domain.ValidateRequiredID(userID, "ユーザーID"); err != nil {
		return 0, err
	}
	return uc.folders.CountByUserID(ctx, userID)
}

// UpdateNoteFolderInput はフォルダ更新の入力。
// Name が空文字なら名前は据え置き、ParentID が nil なら親は据え置きの部分更新。
type UpdateNoteFolderInput struct {
	ID       uint
	UserID   uint
	Name     string
	ParentID *uint
}

// UpdateNoteFolderUseCase は所有権を検証したうえでノートフォルダを更新する。
type UpdateNoteFolderUseCase struct {
	folders repository.NoteFolderRepository
}

// NewUpdateNoteFolderUseCase は UpdateNoteFolderUseCase を生成する。
func NewUpdateNoteFolderUseCase(folders repository.NoteFolderRepository) *UpdateNoteFolderUseCase {
	return &UpdateNoteFolderUseCase{folders: folders}
}

// Execute は名前と親フォルダを更新する。親の変更は自己参照・循環参照を拒否する。
func (uc *UpdateNoteFolderUseCase) Execute(ctx context.Context, in UpdateNoteFolderInput) (*model.NoteFolder, error) {
	folder, err := ensureOwner(ctx, uc.folders.FindByID, in.ID, in.UserID, noteFolderOwnerOf)
	if err != nil {
		return nil, err
	}

	if in.Name != "" {
		if strings.TrimSpace(in.Name) == "" {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "フォルダ名は空白のみにできません", nil)
		}
		folder.Name = strings.TrimSpace(in.Name)
	}

	if in.ParentID != nil {
		if *in.ParentID == in.ID {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "フォルダを自分自身の子にすることはできません", nil)
		}
		// 自分の子孫を親に設定すると木構造が閉路になるため拒否する。
		isDesc, err := uc.isDescendant(ctx, in.ID, *in.ParentID)
		if err != nil {
			return nil, err
		}
		if isDesc {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "循環参照が発生するため、この親フォルダは設定できません", nil)
		}
		folder.ParentID = in.ParentID
	}

	if err := validator.NewNoteFolderValidator().ValidateUpdate(folder.Name); err != nil {
		return nil, err
	}

	if err := uc.folders.Update(ctx, folder); err != nil {
		return nil, err
	}
	return folder, nil
}

// isDescendant は targetID が ancestorID の子孫かを判定する。
//
// 既に閉路を含むデータに対しても必ず停止するよう、訪問済みの ID を記録しながら
// 幅優先で辿る。以前は再帰で辿っており、閉路があると停止せずスタックを食い潰して
// プロセスごと落ちていた（閉路は木の更新が同時に走ったときに生じ得る）。
// 閉路の無い正常なデータに対する判定結果は従来と同じ。
func (uc *UpdateNoteFolderUseCase) isDescendant(ctx context.Context, ancestorID, targetID uint) (bool, error) {
	visited := map[uint]bool{ancestorID: true}
	queue := []uint{ancestorID}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		children, err := uc.folders.FindByParentID(ctx, current)
		if err != nil {
			return false, err
		}
		for _, child := range children {
			if child.ID == targetID {
				return true, nil
			}
			if visited[child.ID] {
				continue
			}
			visited[child.ID] = true
			queue = append(queue, child.ID)
		}
	}
	return false, nil
}

// DeleteNoteFolderUseCase は所有権を検証したうえでノートフォルダを削除する。
type DeleteNoteFolderUseCase struct {
	folders repository.NoteFolderRepository
}

// NewDeleteNoteFolderUseCase は DeleteNoteFolderUseCase を生成する。
func NewDeleteNoteFolderUseCase(folders repository.NoteFolderRepository) *DeleteNoteFolderUseCase {
	return &DeleteNoteFolderUseCase{folders: folders}
}

// Execute はフォルダを削除する。
func (uc *DeleteNoteFolderUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensureOwner(ctx, uc.folders.FindByID, id, userID, noteFolderOwnerOf); err != nil {
		return err
	}
	return uc.folders.Delete(ctx, id)
}

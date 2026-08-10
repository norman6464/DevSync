package usecase

import (
	"context"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postTemplateOwnerOf は所有権チェック用にテンプレートの所有者 ID を取り出す。
func postTemplateOwnerOf(t *model.PostTemplate) uint { return t.UserID }

// CreatePostTemplateUseCase は投稿テンプレートを作成する。
type CreatePostTemplateUseCase struct {
	templates repository.PostTemplateRepository
}

// NewCreatePostTemplateUseCase は CreatePostTemplateUseCase を生成する。
func NewCreatePostTemplateUseCase(templates repository.PostTemplateRepository) *CreatePostTemplateUseCase {
	return &CreatePostTemplateUseCase{templates: templates}
}

// Execute は各項目を検証し、前後の空白を落として作成する。
func (uc *CreatePostTemplateUseCase) Execute(ctx context.Context, tmpl *model.PostTemplate) error {
	tmpl.Name = strings.TrimSpace(tmpl.Name)
	tmpl.ContentTemplate = strings.TrimSpace(tmpl.ContentTemplate)
	tmpl.TitleTemplate = strings.TrimSpace(tmpl.TitleTemplate)

	if err := domain.ValidateStringLength(tmpl.Name, 1, 100, "テンプレート名"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(tmpl.ContentTemplate, 1, 50000, "テンプレート内容"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(tmpl.TitleTemplate, 0, 200, "タイトルテンプレート"); err != nil {
		return err
	}
	return uc.templates.Create(ctx, tmpl)
}

// GetPostTemplateUseCase は所有者本人のテンプレートを取得する。
type GetPostTemplateUseCase struct {
	templates repository.PostTemplateRepository
}

// NewGetPostTemplateUseCase は GetPostTemplateUseCase を生成する。
func NewGetPostTemplateUseCase(templates repository.PostTemplateRepository) *GetPostTemplateUseCase {
	return &GetPostTemplateUseCase{templates: templates}
}

// Execute は所有権を検証したうえでテンプレートを返す。
func (uc *GetPostTemplateUseCase) Execute(ctx context.Context, id, userID uint) (*model.PostTemplate, error) {
	return ensureOwner(ctx, uc.templates.FindByID, id, userID, postTemplateOwnerOf)
}

// ListPostTemplatesUseCase は指定ユーザーのテンプレート一覧を取得する。
type ListPostTemplatesUseCase struct {
	templates repository.PostTemplateRepository
}

// NewListPostTemplatesUseCase は ListPostTemplatesUseCase を生成する。
func NewListPostTemplatesUseCase(templates repository.PostTemplateRepository) *ListPostTemplatesUseCase {
	return &ListPostTemplatesUseCase{templates: templates}
}

// Execute はページネーション付きの一覧と総数を返す。
func (uc *ListPostTemplatesUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.PostTemplate, int64, error) {
	return uc.templates.FindByUserID(ctx, userID, limit, offset)
}

// UpdatePostTemplateUseCase は所有者本人のテンプレートを更新する。
type UpdatePostTemplateUseCase struct {
	templates repository.PostTemplateRepository
}

// NewUpdatePostTemplateUseCase は UpdatePostTemplateUseCase を生成する。
func NewUpdatePostTemplateUseCase(templates repository.PostTemplateRepository) *UpdatePostTemplateUseCase {
	return &UpdatePostTemplateUseCase{templates: templates}
}

// Execute は所有権を検証し、トリム後に空でないフィールドだけを更新する。
func (uc *UpdatePostTemplateUseCase) Execute(ctx context.Context, id, userID uint, updates *model.PostTemplate) (*model.PostTemplate, error) {
	tmpl, err := ensureOwner(ctx, uc.templates.FindByID, id, userID, postTemplateOwnerOf)
	if err != nil {
		return nil, err
	}

	if name := strings.TrimSpace(updates.Name); name != "" {
		if err := domain.ValidateStringLength(name, 1, 100, "テンプレート名"); err != nil {
			return nil, err
		}
		tmpl.Name = name
	}
	if tt := strings.TrimSpace(updates.TitleTemplate); tt != "" {
		if err := domain.ValidateStringLength(tt, 1, 200, "タイトルテンプレート"); err != nil {
			return nil, err
		}
		tmpl.TitleTemplate = tt
	}
	if ct := strings.TrimSpace(updates.ContentTemplate); ct != "" {
		if err := domain.ValidateStringLength(ct, 1, 50000, "テンプレート内容"); err != nil {
			return nil, err
		}
		tmpl.ContentTemplate = ct
	}

	if err := uc.templates.Update(ctx, tmpl); err != nil {
		return nil, err
	}
	return tmpl, nil
}

// DeletePostTemplateUseCase は所有者本人のテンプレートを削除する。
type DeletePostTemplateUseCase struct {
	templates repository.PostTemplateRepository
}

// NewDeletePostTemplateUseCase は DeletePostTemplateUseCase を生成する。
func NewDeletePostTemplateUseCase(templates repository.PostTemplateRepository) *DeletePostTemplateUseCase {
	return &DeletePostTemplateUseCase{templates: templates}
}

// Execute は所有権を検証してから削除する。
func (uc *DeletePostTemplateUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensureOwner(ctx, uc.templates.FindByID, id, userID, postTemplateOwnerOf); err != nil {
		return err
	}
	return uc.templates.Delete(ctx, id)
}

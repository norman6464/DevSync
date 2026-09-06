package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain/validator"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// defaultNoteTitleFromTemplate はテンプレートにデフォルトタイトルが無いときに使うノート名。
const defaultNoteTitleFromTemplate = "新しいノート"

// errDefaultNoteTemplateNotFound はデフォルトテンプレートが未設定のときに返すエラー。
// DomainError ではないため handler では 500 になり、不在を素の DB エラーとして扱っていた
// 移行前の挙動と一致する。
var errDefaultNoteTemplateNotFound = errors.New("デフォルトテンプレートが見つかりません")

// noteTemplateOwnerOf は所有権チェック用にテンプレートの所有者 ID を取り出す。
func noteTemplateOwnerOf(t *model.NoteTemplate) uint { return t.UserID }

// CreateNoteTemplateUseCase はノートテンプレートを作成する。
type CreateNoteTemplateUseCase struct {
	templates repository.NoteTemplateRepository
}

// NewCreateNoteTemplateUseCase は CreateNoteTemplateUseCase を生成する。
func NewCreateNoteTemplateUseCase(templates repository.NoteTemplateRepository) *CreateNoteTemplateUseCase {
	return &CreateNoteTemplateUseCase{templates: templates}
}

// Execute は前後の空白を落として各項目を検証し、テンプレートを作成する。
// デフォルト指定つきの場合は、先に同一ユーザーの既存のデフォルト指定を外す。
func (uc *CreateNoteTemplateUseCase) Execute(ctx context.Context, tmpl *model.NoteTemplate) error {
	tmpl.Name = strings.TrimSpace(tmpl.Name)
	tmpl.ContentTemplate = strings.TrimSpace(tmpl.ContentTemplate)
	tmpl.Description = strings.TrimSpace(tmpl.Description)
	tmpl.DefaultTitle = strings.TrimSpace(tmpl.DefaultTitle)

	v := validator.NewNoteTemplateValidator()
	if err := v.ValidateCreateTemplate(tmpl.Name, tmpl.ContentTemplate); err != nil {
		return err
	}
	if tmpl.Description != "" {
		if err := v.ValidateDescription(tmpl.Description); err != nil {
			return err
		}
	}
	if tmpl.DefaultTitle != "" {
		if err := v.ValidateDefaultTitle(tmpl.DefaultTitle); err != nil {
			return err
		}
	}

	return uc.templates.Create(ctx, tmpl)
}

// GetNoteTemplateUseCase は所有者本人のテンプレートを 1 件取得する。
type GetNoteTemplateUseCase struct {
	templates repository.NoteTemplateRepository
}

// NewGetNoteTemplateUseCase は GetNoteTemplateUseCase を生成する。
func NewGetNoteTemplateUseCase(templates repository.NoteTemplateRepository) *GetNoteTemplateUseCase {
	return &GetNoteTemplateUseCase{templates: templates}
}

// Execute は所有権を検証したうえでテンプレートを返す。
func (uc *GetNoteTemplateUseCase) Execute(ctx context.Context, id, userID uint) (*model.NoteTemplate, error) {
	return ensureOwner(ctx, uc.templates.FindByID, id, userID, noteTemplateOwnerOf)
}

// ListNoteTemplatesUseCase は指定ユーザーのテンプレート一覧を取得する。
type ListNoteTemplatesUseCase struct {
	templates repository.NoteTemplateRepository
}

// NewListNoteTemplatesUseCase は ListNoteTemplatesUseCase を生成する。
func NewListNoteTemplatesUseCase(templates repository.NoteTemplateRepository) *ListNoteTemplatesUseCase {
	return &ListNoteTemplatesUseCase{templates: templates}
}

// Execute はテンプレートを作成日の新しい順で返す。
func (uc *ListNoteTemplatesUseCase) Execute(ctx context.Context, userID uint) ([]model.NoteTemplate, error) {
	return uc.templates.FindByUserID(ctx, userID)
}

// GetDefaultNoteTemplateUseCase はデフォルトに設定されたテンプレートを取得する。
type GetDefaultNoteTemplateUseCase struct {
	templates repository.NoteTemplateRepository
}

// NewGetDefaultNoteTemplateUseCase は GetDefaultNoteTemplateUseCase を生成する。
func NewGetDefaultNoteTemplateUseCase(templates repository.NoteTemplateRepository) *GetDefaultNoteTemplateUseCase {
	return &GetDefaultNoteTemplateUseCase{templates: templates}
}

// Execute はデフォルトテンプレートを返す。未設定なら DomainError ではないエラーを返す
// （移行前と同じく 500 になる）。
func (uc *GetDefaultNoteTemplateUseCase) Execute(ctx context.Context, userID uint) (*model.NoteTemplate, error) {
	tmpl, err := uc.templates.FindDefaultByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if tmpl == nil {
		return nil, errDefaultNoteTemplateNotFound
	}
	return tmpl, nil
}

// UpdateNoteTemplateInput はテンプレート更新の入力。
// 文字列は空文字列なら「変更なし」、IsDefault は nil なら「変更なし」を表す。
type UpdateNoteTemplateInput struct {
	ID              uint
	UserID          uint
	Name            string
	Description     string
	DefaultTitle    string
	ContentTemplate string
	DefaultTags     string
	IsDefault       *bool
}

// UpdateNoteTemplateUseCase は所有者本人のテンプレートを更新する。
type UpdateNoteTemplateUseCase struct {
	templates repository.NoteTemplateRepository
}

// NewUpdateNoteTemplateUseCase は UpdateNoteTemplateUseCase を生成する。
func NewUpdateNoteTemplateUseCase(templates repository.NoteTemplateRepository) *UpdateNoteTemplateUseCase {
	return &UpdateNoteTemplateUseCase{templates: templates}
}

// Execute は所有権を検証し、トリム後に空でないフィールドだけを更新する。
// 更新後の値に対してまとめて検証をかけ、デフォルト指定つきなら書き込み前に既存の指定を外す。
func (uc *UpdateNoteTemplateUseCase) Execute(ctx context.Context, in UpdateNoteTemplateInput) (*model.NoteTemplate, error) {
	tmpl, err := ensureOwner(ctx, uc.templates.FindByID, in.ID, in.UserID, noteTemplateOwnerOf)
	if err != nil {
		return nil, err
	}

	if name := strings.TrimSpace(in.Name); name != "" {
		tmpl.Name = name
	}
	if description := strings.TrimSpace(in.Description); description != "" {
		tmpl.Description = description
	}
	if defaultTitle := strings.TrimSpace(in.DefaultTitle); defaultTitle != "" {
		tmpl.DefaultTitle = defaultTitle
	}
	if contentTemplate := strings.TrimSpace(in.ContentTemplate); contentTemplate != "" {
		tmpl.ContentTemplate = contentTemplate
	}
	if tags := strings.TrimSpace(in.DefaultTags); tags != "" {
		tmpl.DefaultTags = tags
	}
	if in.IsDefault != nil {
		tmpl.IsDefault = *in.IsDefault
	}

	v := validator.NewNoteTemplateValidator()
	if tmpl.Name != "" {
		if err := v.ValidateName(tmpl.Name); err != nil {
			return nil, err
		}
	}
	if tmpl.ContentTemplate != "" {
		if err := v.ValidateContentTemplate(tmpl.ContentTemplate); err != nil {
			return nil, err
		}
	}
	if tmpl.Description != "" {
		if err := v.ValidateDescription(tmpl.Description); err != nil {
			return nil, err
		}
	}
	if tmpl.DefaultTitle != "" {
		if err := v.ValidateDefaultTitle(tmpl.DefaultTitle); err != nil {
			return nil, err
		}
	}

	if err := uc.templates.Update(ctx, tmpl); err != nil {
		return nil, err
	}
	return tmpl, nil
}

// DeleteNoteTemplateUseCase は所有者本人のテンプレートを削除する。
type DeleteNoteTemplateUseCase struct {
	templates repository.NoteTemplateRepository
}

// NewDeleteNoteTemplateUseCase は DeleteNoteTemplateUseCase を生成する。
func NewDeleteNoteTemplateUseCase(templates repository.NoteTemplateRepository) *DeleteNoteTemplateUseCase {
	return &DeleteNoteTemplateUseCase{templates: templates}
}

// Execute は所有権を検証したうえでテンプレートを削除する。
func (uc *DeleteNoteTemplateUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensureOwner(ctx, uc.templates.FindByID, id, userID, noteTemplateOwnerOf); err != nil {
		return err
	}
	return uc.templates.Delete(ctx, id)
}

// CreateNoteFromTemplateUseCase は所有者本人のテンプレートからノートを作成する。
// ノートの検証と書き込みは既存の [CreateNoteUseCase] に委ねる。
type CreateNoteFromTemplateUseCase struct {
	templates  repository.NoteTemplateRepository
	createNote *CreateNoteUseCase
}

// NewCreateNoteFromTemplateUseCase は CreateNoteFromTemplateUseCase を生成する。
func NewCreateNoteFromTemplateUseCase(
	templates repository.NoteTemplateRepository,
	createNote *CreateNoteUseCase,
) *CreateNoteFromTemplateUseCase {
	return &CreateNoteFromTemplateUseCase{templates: templates, createNote: createNote}
}

// Execute は所有権を検証し、テンプレートの内容でノートを作成して返す。
// デフォルトタイトルが空のときは既定のノート名を使う。
func (uc *CreateNoteFromTemplateUseCase) Execute(ctx context.Context, id, userID uint) (*model.Note, error) {
	tmpl, err := ensureOwner(ctx, uc.templates.FindByID, id, userID, noteTemplateOwnerOf)
	if err != nil {
		return nil, err
	}

	note := &model.Note{
		UserID:  userID,
		Title:   tmpl.DefaultTitle,
		Content: tmpl.ContentTemplate,
		Tags:    tmpl.DefaultTags,
	}
	if note.Title == "" {
		note.Title = defaultNoteTitleFromTemplate
	}

	if err := uc.createNote.Execute(ctx, note); err != nil {
		return nil, err
	}
	return note, nil
}

// CountNoteTemplatesUseCase は指定ユーザーのテンプレート総数を取得する。
type CountNoteTemplatesUseCase struct {
	templates repository.NoteTemplateRepository
}

// NewCountNoteTemplatesUseCase は CountNoteTemplatesUseCase を生成する。
func NewCountNoteTemplatesUseCase(templates repository.NoteTemplateRepository) *CountNoteTemplatesUseCase {
	return &CountNoteTemplatesUseCase{templates: templates}
}

// Execute はテンプレート総数を返す。
func (uc *CountNoteTemplatesUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.templates.CountByUserID(ctx, userID)
}

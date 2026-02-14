package validator

import (
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
)

// NoteTemplateValidator はノートテンプレートのバリデーションロジック。
type NoteTemplateValidator struct{}

// NewNoteTemplateValidator は新しいNoteTemplateValidatorインスタンスを生成する。
func NewNoteTemplateValidator() *NoteTemplateValidator {
	return &NoteTemplateValidator{}
}

// ValidateCreateTemplate はテンプレート作成時のバリデーション。
func (v *NoteTemplateValidator) ValidateCreateTemplate(name, contentTemplate string) error {
	if err := v.ValidateName(name); err != nil {
		return err
	}
	if err := v.ValidateContentTemplate(contentTemplate); err != nil {
		return err
	}
	return nil
}

// ValidateUpdateTemplate はテンプレート更新時のバリデーション。
func (v *NoteTemplateValidator) ValidateUpdateTemplate(name, contentTemplate string) error {
	if name != "" {
		if err := v.ValidateName(name); err != nil {
			return err
		}
	}
	if contentTemplate != "" {
		if err := v.ValidateContentTemplate(contentTemplate); err != nil {
			return err
		}
	}
	return nil
}

// ValidateName はテンプレート名のバリデーション。
func (v *NoteTemplateValidator) ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.NewError(domain.ErrCodeValidation, "テンプレート名を入力してください", nil)
	}
	if len(name) > 100 {
		return domain.NewError(domain.ErrCodeValidation, "テンプレート名は100文字以下である必要があります", nil)
	}
	return nil
}

// ValidateContentTemplate は本文テンプレートのバリデーション。
func (v *NoteTemplateValidator) ValidateContentTemplate(content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return domain.NewError(domain.ErrCodeValidation, "本文テンプレートを入力してください", nil)
	}
	return nil
}

// ValidateDescription は説明のバリデーション。
func (v *NoteTemplateValidator) ValidateDescription(description string) error {
	if len(description) > 500 {
		return domain.NewError(domain.ErrCodeValidation, "説明は500文字以下である必要があります", nil)
	}
	return nil
}

// ValidateDefaultTitle はデフォルトタイトルのバリデーション。
func (v *NoteTemplateValidator) ValidateDefaultTitle(title string) error {
	if len(title) > 200 {
		return domain.NewError(domain.ErrCodeValidation, "デフォルトタイトルは200文字以下である必要があります", nil)
	}
	return nil
}

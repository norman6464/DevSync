package validator

import (
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
	return domain.ValidateStringLength(name, 1, 100, "テンプレート名")
}

// ValidateContentTemplate は本文テンプレートのバリデーション。
func (v *NoteTemplateValidator) ValidateContentTemplate(content string) error {
	return domain.ValidateStringLength(content, 1, 50000, "本文テンプレート")
}

// ValidateDescription は説明のバリデーション。
func (v *NoteTemplateValidator) ValidateDescription(description string) error {
	return domain.ValidateStringLength(description, 0, 500, "説明")
}

// ValidateDefaultTitle はデフォルトタイトルのバリデーション。
func (v *NoteTemplateValidator) ValidateDefaultTitle(title string) error {
	return domain.ValidateStringLength(title, 0, 200, "デフォルトタイトル")
}

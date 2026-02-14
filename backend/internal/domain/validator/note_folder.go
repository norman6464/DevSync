package validator

import (
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
)

// NoteFolderValidator はノートフォルダのバリデーションを提供する。
type NoteFolderValidator struct{}

// NewNoteFolderValidator は新しいNoteFolderValidatorインスタンスを生成する。
func NewNoteFolderValidator() *NoteFolderValidator {
	return &NoteFolderValidator{}
}

// ValidateName はフォルダ名のバリデーションを行う。
func (v *NoteFolderValidator) ValidateName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return domain.NewError(domain.ErrCodeValidation, "フォルダ名を入力してください", nil)
	}

	if len(name) > 100 {
		return domain.NewError(domain.ErrCodeValidation, "フォルダ名は100文字以下である必要があります", nil)
	}

	return nil
}

// ValidateCreate はフォルダ作成時のバリデーションを行う。
func (v *NoteFolderValidator) ValidateCreate(name string) error {
	return v.ValidateName(name)
}

// ValidateUpdate はフォルダ更新時のバリデーションを行う。
// 更新時は空文字を許容（部分更新対応）
func (v *NoteFolderValidator) ValidateUpdate(name string) error {
	name = strings.TrimSpace(name)

	// 空文字の場合は更新しないとみなしてOK
	if name == "" {
		return nil
	}

	if len(name) > 100 {
		return domain.NewError(domain.ErrCodeValidation, "フォルダ名は100文字以下である必要があります", nil)
	}

	return nil
}

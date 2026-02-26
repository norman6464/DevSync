package validator

import (
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
	return domain.ValidateStringLength(name, 1, 100, "フォルダ名")
}

// ValidateCreate はフォルダ作成時のバリデーションを行う。
func (v *NoteFolderValidator) ValidateCreate(name string) error {
	return v.ValidateName(name)
}

// ValidateUpdate はフォルダ更新時のバリデーションを行う。
// 更新時は空文字を許容（部分更新対応）
func (v *NoteFolderValidator) ValidateUpdate(name string) error {
	return domain.ValidateStringLength(name, 0, 100, "フォルダ名")
}

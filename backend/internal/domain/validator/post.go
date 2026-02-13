// Package validator provides validation logic for domain entities.
package validator

import (
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
)

// PostValidator handles validation for Post entities.
type PostValidator struct{}

// NewPostValidator creates a new PostValidator instance.
func NewPostValidator() *PostValidator {
	return &PostValidator{}
}

// ValidateTitle validates the post title.
func (v *PostValidator) ValidateTitle(title string) error {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return domain.NewError(domain.ErrCodeValidation, "タイトルは必須です", nil)
	}
	if len(title) > 200 {
		return domain.NewError(domain.ErrCodeValidation, "タイトルは200文字以内で入力してください", nil)
	}
	return nil
}

// ValidateContent validates the post content.
func (v *PostValidator) ValidateContent(content string) error {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return domain.NewError(domain.ErrCodeValidation, "本文は必須です", nil)
	}
	if len(content) > 10000 {
		return domain.NewError(domain.ErrCodeValidation, "本文は10000文字以内で入力してください", nil)
	}
	return nil
}

// ValidatePost validates both title and content.
func (v *PostValidator) ValidatePost(title, content string) error {
	if err := v.ValidateTitle(title); err != nil {
		return err
	}
	if err := v.ValidateContent(content); err != nil {
		return err
	}
	return nil
}

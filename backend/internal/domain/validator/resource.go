// Package validator provides validation logic for domain entities.
package validator

import (
	"github.com/norman6464/devsync/backend/internal/domain"
)

// ResourceValidator handles validation for LearningResource entities.
type ResourceValidator struct{}

// NewResourceValidator creates a new ResourceValidator instance.
func NewResourceValidator() *ResourceValidator {
	return &ResourceValidator{}
}

var (
	// 許可されたカテゴリー
	allowedCategories = []string{"book", "video", "article", "course", "tutorial", "podcast", "tool", "other"}
	// 許可された難易度
	allowedDifficulties = []string{"beginner", "intermediate", "advanced"}
)

// ValidateCreateResource validates inputs for creating a new learning resource.
func (v *ResourceValidator) ValidateCreateResource(title, description, url, category, difficulty string) error {
	// タイトルのバリデーション
	if err := domain.ValidateTitle(title); err != nil {
		return err
	}

	// 説明のバリデーション（オプショナル）
	if description != "" {
		if err := domain.ValidateStringLength(description, 0, 1000, "説明"); err != nil {
			return err
		}
	}

	// URLのバリデーション（必須）
	if err := domain.ValidateURL(url); err != nil {
		return err
	}
	if url == "" {
		return domain.NewError(domain.ErrCodeValidation, "URLは必須です", nil)
	}

	// カテゴリーのバリデーション（必須）
	if category == "" {
		return domain.NewError(domain.ErrCodeValidation, "カテゴリーは必須です", nil)
	}
	if err := domain.ValidateEnum(category, allowedCategories, "カテゴリー"); err != nil {
		return err
	}

	// 難易度のバリデーション（オプショナル）
	if difficulty != "" {
		if err := domain.ValidateEnum(difficulty, allowedDifficulties, "難易度"); err != nil {
			return err
		}
	}

	return nil
}

// ValidateUpdateResource validates inputs for updating an existing learning resource.
func (v *ResourceValidator) ValidateUpdateResource(title, description, url, category, difficulty string) error {
	// 作成時と同じバリデーション
	return v.ValidateCreateResource(title, description, url, category, difficulty)
}

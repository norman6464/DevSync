// Package validator provides validation logic for domain entities.
package validator

import (
	"github.com/norman6464/devsync/backend/internal/domain"
)

// ProjectValidator handles validation for Project entities.
type ProjectValidator struct{}

// NewProjectValidator creates a new ProjectValidator instance.
func NewProjectValidator() *ProjectValidator {
	return &ProjectValidator{}
}

// ValidateCreateProject validates inputs for creating a new project.
func (v *ProjectValidator) ValidateCreateProject(title, description, demoURL, githubURL string) error {
	// タイトルのバリデーション
	if err := domain.ValidateTitle(title); err != nil {
		return err
	}

	// 説明のバリデーション
	if err := domain.ValidateContent(description); err != nil {
		return err
	}

	// デモURLのバリデーション（オプショナル）
	if demoURL != "" {
		if err := domain.ValidateURL(demoURL); err != nil {
			return err
		}
	}

	// GitHub URLのバリデーション（オプショナル）
	if githubURL != "" {
		if err := domain.ValidateURL(githubURL); err != nil {
			return err
		}
	}

	return nil
}

// ValidateUpdateProject validates inputs for updating an existing project.
// 更新では部分更新をサポートするため、各フィールドが空でない場合のみバリデーション
func (v *ProjectValidator) ValidateUpdateProject(title, description, demoURL, githubURL string) error {
	// タイトルのバリデーション（空でない場合のみ）
	if title != "" {
		if err := domain.ValidateTitle(title); err != nil {
			return err
		}
	}

	// 説明のバリデーション（空でない場合のみ）
	if description != "" {
		if err := domain.ValidateContent(description); err != nil {
			return err
		}
	}

	// デモURLのバリデーション（空でない場合のみ）
	if demoURL != "" {
		if err := domain.ValidateURL(demoURL); err != nil {
			return err
		}
	}

	// GitHub URLのバリデーション（空でない場合のみ）
	if githubURL != "" {
		if err := domain.ValidateURL(githubURL); err != nil {
			return err
		}
	}

	return nil
}

// Package validator provides validation logic for domain entities.
package validator

import (
	"github.com/norman6464/devsync/backend/internal/domain"
)

// QuestionValidator handles validation for Question entities.
type QuestionValidator struct{}

// NewQuestionValidator creates a new QuestionValidator instance.
func NewQuestionValidator() *QuestionValidator {
	return &QuestionValidator{}
}

// ValidateCreateQuestion validates inputs for creating a new question.
func (v *QuestionValidator) ValidateCreateQuestion(title, body, tags string) error {
	// タイトルのバリデーション（最大500文字）
	if err := domain.ValidateStringLength(title, 1, 500, "タイトル"); err != nil {
		return err
	}

	// 本文のバリデーション
	if err := domain.ValidateContent(body); err != nil {
		return err
	}

	// タグのバリデーション（オプショナル・カンマ区切り文字列）
	// tagsが空文字列でない場合のみバリデーション
	if tags != "" {
		if err := domain.ValidateStringLength(tags, 1, 300, "タグ"); err != nil {
			return err
		}
	}

	return nil
}

// ValidateUpdateQuestion validates inputs for updating an existing question.
func (v *QuestionValidator) ValidateUpdateQuestion(title, body, tags string) error {
	// 作成時と同じバリデーション
	return v.ValidateCreateQuestion(title, body, tags)
}

// ValidateVote validates vote value (should be 1 or -1).
func (v *QuestionValidator) ValidateVote(value int) error {
	if value != 1 && value != -1 {
		return domain.NewError(domain.ErrCodeValidation, "投票値は1または-1である必要があります", nil)
	}
	return nil
}

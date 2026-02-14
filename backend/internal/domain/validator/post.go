// Package validator provides validation logic for domain entities.
package validator

import (
	"encoding/json"

	"github.com/norman6464/devsync/backend/internal/domain"
)

// PostValidator handles validation for Post entities.
type PostValidator struct{}

// NewPostValidator creates a new PostValidator instance.
func NewPostValidator() *PostValidator {
	return &PostValidator{}
}

// ValidateTitle validates the post title using domain.ValidateTitle.
func (v *PostValidator) ValidateTitle(title string) error {
	return domain.ValidateTitle(title)
}

// ValidateContent validates the post content using domain.ValidateContent.
func (v *PostValidator) ValidateContent(content string) error {
	return domain.ValidateContent(content)
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

// ValidateCreatePost validates inputs for creating a new post.
func (v *PostValidator) ValidateCreatePost(title, content string, imageURLs string, tags []string) error {
	// タイトルと本文のバリデーション
	if err := v.ValidatePost(title, content); err != nil {
		return err
	}

	// 画像URLのバリデーション（オプショナル）
	if imageURLs != "" {
		var urls []string
		if err := json.Unmarshal([]byte(imageURLs), &urls); err != nil {
			return domain.NewError(domain.ErrCodeValidation, "画像URLの形式が不正です", err)
		}
		for _, url := range urls {
			if err := domain.ValidateURL(url); err != nil {
				return err
			}
		}
	}

	// タグのバリデーション（オプショナル）
	if len(tags) > 0 {
		if err := domain.ValidateTags(tags); err != nil {
			return err
		}
	}

	return nil
}

// ValidateUpdatePost validates inputs for updating an existing post.
// 更新では部分更新をサポートするため、各フィールドが空でない場合のみバリデーション
func (v *PostValidator) ValidateUpdatePost(title, content string, imageURLs string) error {
	// タイトルのバリデーション（空でない場合のみ）
	if title != "" {
		if err := v.ValidateTitle(title); err != nil {
			return err
		}
	}

	// 本文のバリデーション（空でない場合のみ）
	if content != "" {
		if err := v.ValidateContent(content); err != nil {
			return err
		}
	}

	// 画像URLのバリデーション（空でない場合のみ）
	if imageURLs != "" {
		var urls []string
		if err := json.Unmarshal([]byte(imageURLs), &urls); err != nil {
			return domain.NewError(domain.ErrCodeValidation, "画像URLの形式が不正です", err)
		}
		for _, url := range urls {
			if err := domain.ValidateURL(url); err != nil {
				return err
			}
		}
	}

	return nil
}

// ValidateComment validates comment content.
func (v *PostValidator) ValidateComment(content string) error {
	return domain.ValidateStringLength(content, 1, 1000, "コメント")
}

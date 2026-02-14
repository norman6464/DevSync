package validator

import (
	"github.com/norman6464/devsync/backend/internal/domain"
)

// NoteValidator handles validation for Note entities.
type NoteValidator struct{}

// NewNoteValidator creates a new NoteValidator instance.
func NewNoteValidator() *NoteValidator {
	return &NoteValidator{}
}

// ValidateTitle validates the note title.
// タイトルは1〜200文字、空文字・スペースのみは不可。
func (v *NoteValidator) ValidateTitle(title string) error {
	return domain.ValidateTitle(title)
}

// ValidateContent validates the note content.
// 本文は0〜100,000文字（約100KB）。空文字OK（マークダウンエディタで空の場合あり）。
func (v *NoteValidator) ValidateContent(content string) error {
	return domain.ValidateStringLength(content, 0, 100000, "本文")
}

// ValidateTags validates the note tags.
// タグは0〜500文字。空文字OK（タグなしノート可）。
func (v *NoteValidator) ValidateTags(tags string) error {
	return domain.ValidateStringLength(tags, 0, 500, "タグ")
}

// ValidateCreateNote validates inputs for creating a new note.
func (v *NoteValidator) ValidateCreateNote(title, content, tags string) error {
	// タイトルのバリデーション（必須）
	if err := v.ValidateTitle(title); err != nil {
		return err
	}

	// 本文のバリデーション（オプショナル）
	if err := v.ValidateContent(content); err != nil {
		return err
	}

	// タグのバリデーション（オプショナル）
	if err := v.ValidateTags(tags); err != nil {
		return err
	}

	return nil
}

// ValidateUpdateNote validates inputs for updating an existing note.
// 更新では部分更新をサポートするため、各フィールドが空でない場合のみバリデーション。
func (v *NoteValidator) ValidateUpdateNote(title, content, tags string) error {
	// タイトルのバリデーション（空でない場合のみ）
	if title != "" {
		if err := v.ValidateTitle(title); err != nil {
			return err
		}
	}

	// 本文のバリデーション（空でない場合のみ）
	// 注：部分更新時に明示的に空にする場合は、別途処理が必要
	if content != "" {
		if err := v.ValidateContent(content); err != nil {
			return err
		}
	}

	// タグのバリデーション（空でない場合のみ）
	if tags != "" {
		if err := v.ValidateTags(tags); err != nil {
			return err
		}
	}

	return nil
}

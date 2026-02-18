package service

import (
	"github.com/norman6464/devsync/backend/internal/domain/validator"
	"github.com/norman6464/devsync/backend/internal/model"
)

// NoteTemplateRepositoryInterface はNoteTemplateRepositoryのインターフェース。
type NoteTemplateRepositoryInterface interface {
	Create(template *model.NoteTemplate) error
	FindByID(id uint) (*model.NoteTemplate, error)
	FindByUserID(userID uint) ([]model.NoteTemplate, error)
	FindDefaultByUserID(userID uint) (*model.NoteTemplate, error)
	Update(template *model.NoteTemplate) error
	Delete(id uint) error
	ClearDefaultFlag(userID uint) error
}

// NoteCreatorInterface はノート作成機能のインターフェース。
type NoteCreatorInterface interface {
	Create(note *model.Note) error
}

// NoteTemplateService はノートテンプレートのビジネスロジック。
type NoteTemplateService struct {
	repo        NoteTemplateRepositoryInterface
	noteCreator NoteCreatorInterface
}

// NewNoteTemplateService は新しいNoteTemplateServiceインスタンスを生成する。
func NewNoteTemplateService(repo NoteTemplateRepositoryInterface, noteCreator NoteCreatorInterface) *NoteTemplateService {
	return &NoteTemplateService{repo: repo, noteCreator: noteCreator}
}

// Create は新しいテンプレートを作成する。
func (s *NoteTemplateService) Create(template *model.NoteTemplate) error {
	v := validator.NewNoteTemplateValidator()
	if err := v.ValidateCreateTemplate(template.Name, template.ContentTemplate); err != nil {
		return err
	}
	if template.Description != "" {
		if err := v.ValidateDescription(template.Description); err != nil {
			return err
		}
	}
	if template.DefaultTitle != "" {
		if err := v.ValidateDefaultTitle(template.DefaultTitle); err != nil {
			return err
		}
	}

	// is_default=trueの場合、既存のデフォルトフラグをクリア
	if template.IsDefault {
		if err := s.repo.ClearDefaultFlag(template.UserID); err != nil {
			return err
		}
	}

	return s.repo.Create(template)
}

// GetByID は指定IDのテンプレートを取得する。
func (s *NoteTemplateService) GetByID(id uint) (*model.NoteTemplate, error) {
	return s.repo.FindByID(id)
}

// GetByUserID は指定ユーザーの全テンプレートを取得する。
func (s *NoteTemplateService) GetByUserID(userID uint) ([]model.NoteTemplate, error) {
	return s.repo.FindByUserID(userID)
}

// GetDefaultByUserID は指定ユーザーのデフォルトテンプレートを取得する。
func (s *NoteTemplateService) GetDefaultByUserID(userID uint) (*model.NoteTemplate, error) {
	return s.repo.FindDefaultByUserID(userID)
}

// Update は所有権を検証した後、テンプレートを更新する。
func (s *NoteTemplateService) Update(id, userID uint, name, description, defaultTitle, contentTemplate, defaultTags string, isDefault *bool) (*model.NoteTemplate, error) {
	template, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if template.UserID != userID {
		return nil, ErrForbidden
	}

	if name != "" {
		template.Name = name
	}
	if description != "" {
		template.Description = description
	}
	if defaultTitle != "" {
		template.DefaultTitle = defaultTitle
	}
	if contentTemplate != "" {
		template.ContentTemplate = contentTemplate
	}
	if defaultTags != "" {
		template.DefaultTags = defaultTags
	}
	if isDefault != nil {
		template.IsDefault = *isDefault
	}

	v := validator.NewNoteTemplateValidator()
	if template.Name != "" {
		if err := v.ValidateName(template.Name); err != nil {
			return nil, err
		}
	}
	if template.ContentTemplate != "" {
		if err := v.ValidateContentTemplate(template.ContentTemplate); err != nil {
			return nil, err
		}
	}
	if template.Description != "" {
		if err := v.ValidateDescription(template.Description); err != nil {
			return nil, err
		}
	}
	if template.DefaultTitle != "" {
		if err := v.ValidateDefaultTitle(template.DefaultTitle); err != nil {
			return nil, err
		}
	}

	if template.IsDefault {
		if err := s.repo.ClearDefaultFlag(template.UserID); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Update(template); err != nil {
		return nil, err
	}
	return template, nil
}

// Delete は所有権を検証した後、テンプレートを削除する。
func (s *NoteTemplateService) Delete(id, userID uint) error {
	template, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if template.UserID != userID {
		return ErrForbidden
	}
	return s.repo.Delete(id)
}

// UseTemplate は所有権を検証した後、テンプレートからノートを作成する。
func (s *NoteTemplateService) UseTemplate(id, userID uint) (*model.Note, error) {
	template, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if template.UserID != userID {
		return nil, ErrForbidden
	}

	note := &model.Note{
		UserID:  userID,
		Title:   template.DefaultTitle,
		Content: template.ContentTemplate,
		Tags:    template.DefaultTags,
	}
	if note.Title == "" {
		note.Title = "新しいノート"
	}

	if err := s.noteCreator.Create(note); err != nil {
		return nil, err
	}
	return note, nil
}

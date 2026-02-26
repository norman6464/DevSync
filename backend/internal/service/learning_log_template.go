package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// LearningLogCreatorInterface は学習ログ作成機能のインターフェース。
type LearningLogCreatorInterface interface {
	Create(log *model.LearningLog) error
}

// LearningLogTemplateService は学習ログテンプレートのビジネスロジック。
type LearningLogTemplateService struct {
	repo       repository.LearningLogTemplateRepositoryInterface
	logCreator LearningLogCreatorInterface
}

// NewLearningLogTemplateService は新しいLearningLogTemplateServiceインスタンスを生成する。
func NewLearningLogTemplateService(repo repository.LearningLogTemplateRepositoryInterface, logCreator LearningLogCreatorInterface) *LearningLogTemplateService {
	return &LearningLogTemplateService{repo: repo, logCreator: logCreator}
}

// Create は新しいテンプレートを作成する。
func (s *LearningLogTemplateService) Create(template *model.LearningLogTemplate) error {
	if template.Name == "" {
		return domain.NewError(domain.ErrCodeValidation, "テンプレート名は必須です", nil)
	}
	if len([]rune(template.Name)) > 100 {
		return domain.NewError(domain.ErrCodeValidation, "テンプレート名は100文字以下である必要があります", nil)
	}
	if len([]rune(template.DefaultTitle)) > 200 {
		return domain.NewError(domain.ErrCodeValidation, "デフォルトタイトルは200文字以下である必要があります", nil)
	}
	if len([]rune(template.DefaultContent)) > 50000 {
		return domain.NewError(domain.ErrCodeValidation, "デフォルト本文は50000文字以下である必要があります", nil)
	}
	if template.DefaultCategory != "" && !model.ValidCategories[template.DefaultCategory] {
		return domain.NewError(domain.ErrCodeValidation, "無効なカテゴリです", nil)
	}
	if template.DefaultDuration < 0 || template.DefaultDuration > 1440 {
		return domain.NewError(domain.ErrCodeValidation, "デフォルト時間は0〜1440分の範囲で指定してください", nil)
	}

	if template.IsDefault {
		if err := s.repo.ClearDefaultFlag(template.UserID); err != nil {
			return err
		}
	}

	return s.repo.Create(template)
}

// GetByID は指定IDのテンプレートを取得する。所有権を検証する。
func (s *LearningLogTemplateService) GetByID(id, userID uint) (*model.LearningLogTemplate, error) {
	return s.findAndCheckOwnership(id, userID)
}

// GetByUserID は指定ユーザーの全テンプレートを取得する。
func (s *LearningLogTemplateService) GetByUserID(userID uint) ([]model.LearningLogTemplate, error) {
	return s.repo.FindByUserID(userID)
}

// GetDefaultByUserID は指定ユーザーのデフォルトテンプレートを取得する。
func (s *LearningLogTemplateService) GetDefaultByUserID(userID uint) (*model.LearningLogTemplate, error) {
	return s.repo.FindDefaultByUserID(userID)
}

// findAndCheckOwnership はテンプレートを取得し、指定ユーザーが所有者かを検証する。
func (s *LearningLogTemplateService) findAndCheckOwnership(id, userID uint) (*model.LearningLogTemplate, error) {
	return checkOwnership(s.repo.FindByID, id, userID, func(t *model.LearningLogTemplate) uint { return t.UserID })
}

// Update は所有権を検証した後、テンプレートを更新する。
func (s *LearningLogTemplateService) Update(id, userID uint, name, defaultTitle, defaultContent string, defaultCategory model.LogCategory, defaultDuration *int, isDefault *bool) (*model.LearningLogTemplate, error) {
	template, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if name != "" {
		if len([]rune(name)) > 100 {
			return nil, domain.NewError(domain.ErrCodeValidation, "テンプレート名は100文字以下である必要があります", nil)
		}
		template.Name = name
	}
	if defaultTitle != "" {
		if len([]rune(defaultTitle)) > 200 {
			return nil, domain.NewError(domain.ErrCodeValidation, "デフォルトタイトルは200文字以下である必要があります", nil)
		}
		template.DefaultTitle = defaultTitle
	}
	if defaultContent != "" {
		if len([]rune(defaultContent)) > 50000 {
			return nil, domain.NewError(domain.ErrCodeValidation, "デフォルト本文は50000文字以下である必要があります", nil)
		}
		template.DefaultContent = defaultContent
	}
	if defaultCategory != "" {
		if !model.ValidCategories[defaultCategory] {
			return nil, domain.NewError(domain.ErrCodeValidation, "無効なカテゴリです", nil)
		}
		template.DefaultCategory = defaultCategory
	}
	if defaultDuration != nil {
		if *defaultDuration < 0 || *defaultDuration > 1440 {
			return nil, domain.NewError(domain.ErrCodeValidation, "デフォルト時間は0〜1440分の範囲で指定してください", nil)
		}
		template.DefaultDuration = *defaultDuration
	}
	if isDefault != nil {
		template.IsDefault = *isDefault
		if *isDefault {
			if err := s.repo.ClearDefaultFlag(template.UserID); err != nil {
				return nil, err
			}
		}
	}

	if err := s.repo.Update(template); err != nil {
		return nil, err
	}
	return template, nil
}

// Delete は所有権を検証した後、テンプレートを削除する。
func (s *LearningLogTemplateService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// UseTemplate は所有権を検証した後、テンプレートから学習ログを作成する。
func (s *LearningLogTemplateService) UseTemplate(id, userID uint) (*model.LearningLog, error) {
	template, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	log := &model.LearningLog{
		UserID:   userID,
		Title:    template.DefaultTitle,
		Content:  template.DefaultContent,
		Category: template.DefaultCategory,
		Duration: template.DefaultDuration,
		Source:   model.LogSourceManual,
	}
	if log.Title == "" {
		log.Title = template.Name
	}
	if log.Category == "" {
		log.Category = model.LogCategoryOther
	}

	if err := s.logCreator.Create(log); err != nil {
		return nil, err
	}
	return log, nil
}

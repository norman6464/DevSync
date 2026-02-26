package service

import (
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// PostTemplateService は投稿テンプレートのビジネスロジックを提供する。
type PostTemplateService struct {
	repo repository.PostTemplateRepositoryInterface
}

// NewPostTemplateService は新しいPostTemplateServiceインスタンスを生成する。
func NewPostTemplateService(repo repository.PostTemplateRepositoryInterface) *PostTemplateService {
	return &PostTemplateService{repo: repo}
}

// Create は新しい投稿テンプレートを作成する。
func (s *PostTemplateService) Create(tmpl *model.PostTemplate) error {
	if strings.TrimSpace(tmpl.Name) == "" {
		return domain.NewError(domain.ErrCodeBadRequest, "テンプレート名は必須です", nil)
	}
	if strings.TrimSpace(tmpl.ContentTemplate) == "" {
		return domain.NewError(domain.ErrCodeBadRequest, "テンプレート内容は必須です", nil)
	}
	return s.repo.Create(tmpl)
}

// GetByID は指定IDの投稿テンプレートを取得する。所有権を検証する。
func (s *PostTemplateService) GetByID(id, userID uint) (*model.PostTemplate, error) {
	return s.findAndCheckOwnership(id, userID)
}

// GetByUserID は指定ユーザーの投稿テンプレート一覧を取得する。
func (s *PostTemplateService) GetByUserID(userID uint, limit, offset int) ([]model.PostTemplate, int64, error) {
	return s.repo.FindByUserID(userID, limit, offset)
}

// Update は所有権を検証した後、投稿テンプレートを更新する。
func (s *PostTemplateService) Update(id, userID uint, updates *model.PostTemplate) (*model.PostTemplate, error) {
	tmpl, err := s.findAndCheckOwnership(id, userID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(updates.Name) != "" {
		tmpl.Name = updates.Name
	}
	if strings.TrimSpace(updates.TitleTemplate) != "" {
		tmpl.TitleTemplate = updates.TitleTemplate
	}
	if strings.TrimSpace(updates.ContentTemplate) != "" {
		tmpl.ContentTemplate = updates.ContentTemplate
	}

	if err := s.repo.Update(tmpl); err != nil {
		return nil, err
	}
	return tmpl, nil
}

// Delete は所有権を検証した後、投稿テンプレートを削除する。
func (s *PostTemplateService) Delete(id, userID uint) error {
	if _, err := s.findAndCheckOwnership(id, userID); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

// findAndCheckOwnership は投稿テンプレートを取得し、所有権を検証する。
func (s *PostTemplateService) findAndCheckOwnership(id, userID uint) (*model.PostTemplate, error) {
	return checkOwnership(s.repo.FindByID, id, userID, func(t *model.PostTemplate) uint { return t.UserID })
}

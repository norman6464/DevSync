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
	tmpl.Name = strings.TrimSpace(tmpl.Name)
	tmpl.ContentTemplate = strings.TrimSpace(tmpl.ContentTemplate)
	tmpl.TitleTemplate = strings.TrimSpace(tmpl.TitleTemplate)

	if err := domain.ValidateStringLength(tmpl.Name, 1, 100, "テンプレート名"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(tmpl.ContentTemplate, 1, 50000, "テンプレート内容"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(tmpl.TitleTemplate, 0, 200, "タイトルテンプレート"); err != nil {
		return err
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

	if name := strings.TrimSpace(updates.Name); name != "" {
		if err := domain.ValidateStringLength(name, 1, 100, "テンプレート名"); err != nil {
			return nil, err
		}
		tmpl.Name = name
	}
	if tt := strings.TrimSpace(updates.TitleTemplate); tt != "" {
		if err := domain.ValidateStringLength(tt, 1, 200, "タイトルテンプレート"); err != nil {
			return nil, err
		}
		tmpl.TitleTemplate = tt
	}
	if ct := strings.TrimSpace(updates.ContentTemplate); ct != "" {
		if err := domain.ValidateStringLength(ct, 1, 50000, "テンプレート内容"); err != nil {
			return nil, err
		}
		tmpl.ContentTemplate = ct
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

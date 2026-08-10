package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// PostTemplateRepository は投稿テンプレートの永続化に対する、usecase 側が要求する契約。
type PostTemplateRepository interface {
	Create(ctx context.Context, template *model.PostTemplate) error
	FindByID(ctx context.Context, id uint) (*model.PostTemplate, error)
	FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.PostTemplate, int64, error)
	Update(ctx context.Context, template *model.PostTemplate) error
	Delete(ctx context.Context, id uint) error
}

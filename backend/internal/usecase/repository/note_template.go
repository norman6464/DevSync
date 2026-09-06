package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// NoteTemplateRepository はノートテンプレートの永続化に対する、usecase 側が要求する契約。
type NoteTemplateRepository interface {
	Create(ctx context.Context, template *model.NoteTemplate) error
	Update(ctx context.Context, template *model.NoteTemplate) error
	Delete(ctx context.Context, id uint) error

	// FindByID は指定 ID のテンプレートを返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindByID(ctx context.Context, id uint) (*model.NoteTemplate, error)
	// FindByUserID は指定ユーザーの全テンプレートを作成日の新しい順で返す。
	FindByUserID(ctx context.Context, userID uint) ([]model.NoteTemplate, error)
	// FindDefaultByUserID はデフォルトに設定されたテンプレートを返す。
	// 未設定の場合は「不在」を表す (nil, nil) を返す。
	FindDefaultByUserID(ctx context.Context, userID uint) (*model.NoteTemplate, error)

	CountByUserID(ctx context.Context, userID uint) (int64, error)
}

package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// noteTemplateRepository は [repository.NoteTemplateRepository] の sqlc(pgx) 実装。
type noteTemplateRepository struct {
	q *sqlcgen.Queries
}

// NewNoteTemplateRepository は NoteTemplateRepository の sqlc(pgx) 実装を返す。
func NewNoteTemplateRepository(q *sqlcgen.Queries) repository.NoteTemplateRepository {
	return &noteTemplateRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NoteTemplateRepository = (*noteTemplateRepository)(nil)

// toModelNoteTemplate は sqlc の生成行を model.NoteTemplate へ変換する。
func toModelNoteTemplate(row sqlcgen.NoteTemplate) model.NoteTemplate {
	return model.NoteTemplate{
		ID:              uint(row.ID),
		UserID:          uint(row.UserID),
		Name:            row.Name,
		Description:     fromStringPtr(row.Description),
		DefaultTitle:    fromStringPtr(row.DefaultTitle),
		ContentTemplate: row.ContentTemplate,
		DefaultTags:     fromStringPtr(row.DefaultTags),
		IsDefault:       fromBoolPtr(row.IsDefault),
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
	}
}

// Create は新しいテンプレートを作成する。
func (r *noteTemplateRepository) Create(ctx context.Context, template *model.NoteTemplate) error {
	row, err := r.q.CreateNoteTemplate(ctx, sqlcgen.CreateNoteTemplateParams{
		UserID:          int64(template.UserID),
		Name:            template.Name,
		Description:     &template.Description,
		DefaultTitle:    &template.DefaultTitle,
		ContentTemplate: template.ContentTemplate,
		DefaultTags:     &template.DefaultTags,
		IsDefault:       &template.IsDefault,
	})
	if err != nil {
		return err
	}
	*template = toModelNoteTemplate(row)
	return nil
}

// Update は既存のテンプレートを更新する。
func (r *noteTemplateRepository) Update(ctx context.Context, template *model.NoteTemplate) error {
	row, err := r.q.UpdateNoteTemplate(ctx, sqlcgen.UpdateNoteTemplateParams{
		ID:              int64(template.ID),
		Name:            template.Name,
		Description:     &template.Description,
		DefaultTitle:    &template.DefaultTitle,
		ContentTemplate: template.ContentTemplate,
		DefaultTags:     &template.DefaultTags,
		IsDefault:       &template.IsDefault,
	})
	if err != nil {
		return err
	}
	*template = toModelNoteTemplate(row)
	return nil
}

// Delete はテンプレートを削除する。
func (r *noteTemplateRepository) Delete(ctx context.Context, id uint) error {
	return r.q.DeleteNoteTemplate(ctx, int64(id))
}

// FindByID は指定 ID のテンプレートを取得する。不在の場合は (nil, nil) を返す。
func (r *noteTemplateRepository) FindByID(ctx context.Context, id uint) (*model.NoteTemplate, error) {
	row, err := r.q.GetNoteTemplateByID(ctx, int64(id))
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	template := toModelNoteTemplate(row)
	return &template, nil
}

// FindByUserID は指定ユーザーの全テンプレートを作成日の新しい順で取得する。
func (r *noteTemplateRepository) FindByUserID(ctx context.Context, userID uint) ([]model.NoteTemplate, error) {
	rows, err := r.q.ListNoteTemplatesByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	templates := make([]model.NoteTemplate, len(rows))
	for i, row := range rows {
		templates[i] = toModelNoteTemplate(row)
	}
	return templates, nil
}

// FindDefaultByUserID はデフォルトに設定されたテンプレートを取得する。未設定の場合は (nil, nil) を返す。
func (r *noteTemplateRepository) FindDefaultByUserID(ctx context.Context, userID uint) (*model.NoteTemplate, error) {
	row, err := r.q.GetDefaultNoteTemplateByUser(ctx, int64(userID))
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	template := toModelNoteTemplate(row)
	return &template, nil
}

// ClearDefaultFlag は指定ユーザーの全テンプレートのデフォルト指定を外す。
func (r *noteTemplateRepository) ClearDefaultFlag(ctx context.Context, userID uint) error {
	return r.q.ClearNoteTemplateDefaultFlag(ctx, int64(userID))
}

// CountByUserID は指定ユーザーのテンプレート総数を返す。
func (r *noteTemplateRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountNoteTemplatesByUser(ctx, int64(userID))
}

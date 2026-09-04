package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// postTemplateRepository は [repository.PostTemplateRepository] の sqlc(pgx) 実装。
type postTemplateRepository struct {
	q *sqlcgen.Queries
}

// NewPostTemplateRepository は PostTemplateRepository の sqlc(pgx) 実装を返す。
func NewPostTemplateRepository(q *sqlcgen.Queries) repository.PostTemplateRepository {
	return &postTemplateRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.PostTemplateRepository = (*postTemplateRepository)(nil)

// toModelPostTemplate は sqlc の生成行を model.PostTemplate へ変換する。
func toModelPostTemplate(row sqlcgen.PostTemplate) model.PostTemplate {
	titleTemplate := ""
	if row.TitleTemplate != nil {
		titleTemplate = *row.TitleTemplate
	}
	return model.PostTemplate{
		ID:              uint(row.ID),
		UserID:          uint(row.UserID),
		Name:            row.Name,
		TitleTemplate:   titleTemplate,
		ContentTemplate: row.ContentTemplate,
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       row.UpdatedAt.Time,
	}
}

// Create は新しい投稿テンプレートをデータベースに作成する。
func (r *postTemplateRepository) Create(ctx context.Context, template *model.PostTemplate) error {
	row, err := r.q.CreatePostTemplate(ctx, sqlcgen.CreatePostTemplateParams{
		UserID:          int64(template.UserID),
		Name:            template.Name,
		TitleTemplate:   &template.TitleTemplate,
		ContentTemplate: template.ContentTemplate,
	})
	if err != nil {
		return err
	}
	*template = toModelPostTemplate(row)
	return nil
}

// FindByID は指定IDの投稿テンプレートを取得する。不在は (nil, nil) を返す。
func (r *postTemplateRepository) FindByID(ctx context.Context, id uint) (*model.PostTemplate, error) {
	row, err := r.q.GetPostTemplateByID(ctx, int64(id))
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	template := toModelPostTemplate(row)
	return &template, nil
}

// FindByUserID は指定ユーザーの投稿テンプレートをページネーション付きで取得する（新しい順）。
func (r *postTemplateRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.PostTemplate, int64, error) {
	rows, err := r.q.ListPostTemplatesByUserID(ctx, sqlcgen.ListPostTemplatesByUserIDParams{
		UserID: int64(userID),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := r.q.CountPostTemplatesByUserID(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}

	templates := make([]model.PostTemplate, len(rows))
	for i, row := range rows {
		templates[i] = toModelPostTemplate(row)
	}
	return templates, total, nil
}

// Update は既存の投稿テンプレートを更新する。
func (r *postTemplateRepository) Update(ctx context.Context, template *model.PostTemplate) error {
	row, err := r.q.UpdatePostTemplate(ctx, sqlcgen.UpdatePostTemplateParams{
		ID:              int64(template.ID),
		Name:            template.Name,
		TitleTemplate:   &template.TitleTemplate,
		ContentTemplate: template.ContentTemplate,
	})
	if err != nil {
		return err
	}
	*template = toModelPostTemplate(row)
	return nil
}

// Delete は指定IDの投稿テンプレートを削除する。
func (r *postTemplateRepository) Delete(ctx context.Context, id uint) error {
	return r.q.DeletePostTemplate(ctx, int64(id))
}

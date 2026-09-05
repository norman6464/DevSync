package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// learningLogTemplateRepository は [repository.LearningLogTemplateRepository] の sqlc(pgx) 実装。
type learningLogTemplateRepository struct {
	q *sqlcgen.Queries
}

// NewLearningLogTemplateRepository は LearningLogTemplateRepository の sqlc(pgx) 実装を返す。
func NewLearningLogTemplateRepository(q *sqlcgen.Queries) repository.LearningLogTemplateRepository {
	return &learningLogTemplateRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.LearningLogTemplateRepository = (*learningLogTemplateRepository)(nil)

// toModelLearningLogTemplate は sqlc の生成行を model.LearningLogTemplate へ変換する。
func toModelLearningLogTemplate(row sqlcgen.LearningLogTemplate) model.LearningLogTemplate {
	return model.LearningLogTemplate{
		ID:              uint(row.ID),
		UserID:          uint(row.UserID),
		Name:            row.Name,
		DefaultTitle:    fromStringPtr(row.DefaultTitle),
		DefaultContent:  fromStringPtr(row.DefaultContent),
		DefaultCategory: model.LogCategory(fromStringPtr(row.DefaultCategory)),
		DefaultDuration: int(fromInt64PtrValue(row.DefaultDuration)),
		IsDefault:       fromBoolPtr(row.IsDefault),
		CreatedAt:       timeValue(fromTimestamptz(row.CreatedAt)),
		UpdatedAt:       timeValue(fromTimestamptz(row.UpdatedAt)),
	}
}

// Create は新しいテンプレートを作成する。
func (r *learningLogTemplateRepository) Create(ctx context.Context, template *model.LearningLogTemplate) error {
	category := string(template.DefaultCategory)
	row, err := r.q.CreateLearningLogTemplate(ctx, sqlcgen.CreateLearningLogTemplateParams{
		UserID:          int64(template.UserID),
		Name:            template.Name,
		DefaultTitle:    &template.DefaultTitle,
		DefaultContent:  &template.DefaultContent,
		DefaultCategory: &category,
		DefaultDuration: toInt64Ptr(template.DefaultDuration),
		IsDefault:       &template.IsDefault,
	})
	if err != nil {
		return err
	}
	*template = toModelLearningLogTemplate(row)
	return nil
}

// Update は既存のテンプレートを更新する（GORMのSave＝全カラム上書きに相当）。
func (r *learningLogTemplateRepository) Update(ctx context.Context, template *model.LearningLogTemplate) error {
	category := string(template.DefaultCategory)
	row, err := r.q.UpdateLearningLogTemplate(ctx, sqlcgen.UpdateLearningLogTemplateParams{
		ID:              int64(template.ID),
		Name:            template.Name,
		DefaultTitle:    &template.DefaultTitle,
		DefaultContent:  &template.DefaultContent,
		DefaultCategory: &category,
		DefaultDuration: toInt64Ptr(template.DefaultDuration),
		IsDefault:       &template.IsDefault,
	})
	if err != nil {
		return err
	}
	*template = toModelLearningLogTemplate(row)
	return nil
}

// Delete はテンプレートを削除する。
func (r *learningLogTemplateRepository) Delete(ctx context.Context, id uint) error {
	return r.q.DeleteLearningLogTemplate(ctx, int64(id))
}

// FindByID は指定 ID のテンプレートを取得する。不在の場合は (nil, nil) を返す。
func (r *learningLogTemplateRepository) FindByID(ctx context.Context, id uint) (*model.LearningLogTemplate, error) {
	row, err := r.q.GetLearningLogTemplateByID(ctx, int64(id))
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	template := toModelLearningLogTemplate(row)
	return &template, nil
}

// FindByUserID は指定ユーザーの全テンプレートを作成日の新しい順で取得する。
func (r *learningLogTemplateRepository) FindByUserID(ctx context.Context, userID uint) ([]model.LearningLogTemplate, error) {
	rows, err := r.q.ListLearningLogTemplatesByUser(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	templates := make([]model.LearningLogTemplate, len(rows))
	for i, row := range rows {
		templates[i] = toModelLearningLogTemplate(row)
	}
	return templates, nil
}

// FindDefaultByUserID はデフォルトに設定されたテンプレートを取得する。未設定の場合は (nil, nil) を返す。
func (r *learningLogTemplateRepository) FindDefaultByUserID(ctx context.Context, userID uint) (*model.LearningLogTemplate, error) {
	row, err := r.q.GetDefaultLearningLogTemplateByUser(ctx, int64(userID))
	if err != nil {
		if isNoRows(err) {
			return nil, nil
		}
		return nil, err
	}
	template := toModelLearningLogTemplate(row)
	return &template, nil
}

// ClearDefaultFlag は指定ユーザーの全テンプレートのデフォルト指定を外す。
func (r *learningLogTemplateRepository) ClearDefaultFlag(ctx context.Context, userID uint) error {
	return r.q.ClearLearningLogTemplateDefaultFlag(ctx, int64(userID))
}

// CountByUserID は指定ユーザーのテンプレート総数を返す。
func (r *learningLogTemplateRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountLearningLogTemplatesByUser(ctx, int64(userID))
}

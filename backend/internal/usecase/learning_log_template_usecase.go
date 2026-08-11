package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// msgInvalidLearningLogDuration はデフォルト時間が範囲外のときのメッセージ。
const msgInvalidLearningLogDuration = "デフォルト時間は0〜1440分の範囲で指定してください"

// errDefaultLearningLogTemplateNotFound はデフォルトテンプレートが未設定のときに返すエラー。
// DomainError ではないため handler では 500 になり、不在を素の DB エラーとして扱っていた
// 移行前の挙動と一致する。
var errDefaultLearningLogTemplateNotFound = errors.New("デフォルトテンプレートが見つかりません")

// learningLogTemplateOwnerOf は所有権チェック用にテンプレートの所有者 ID を取り出す。
func learningLogTemplateOwnerOf(t *model.LearningLogTemplate) uint { return t.UserID }

// validateLearningLogTemplateCategory はデフォルトカテゴリを検証する。空文字は未指定として許容する。
func validateLearningLogTemplateCategory(category model.LogCategory) error {
	if category != "" && !model.ValidCategories[category] {
		return domain.NewError(domain.ErrCodeValidation, msgInvalidLogCategory, nil)
	}
	return nil
}

// validateLearningLogTemplateDuration はデフォルト時間が 0〜1440 分の範囲かを検証する。
func validateLearningLogTemplateDuration(duration int) error {
	if duration < 0 || duration > 1440 {
		return domain.NewError(domain.ErrCodeValidation, msgInvalidLearningLogDuration, nil)
	}
	return nil
}

// CreateLearningLogTemplateUseCase は学習ログテンプレートを作成する。
type CreateLearningLogTemplateUseCase struct {
	templates repository.LearningLogTemplateRepository
}

// NewCreateLearningLogTemplateUseCase は CreateLearningLogTemplateUseCase を生成する。
func NewCreateLearningLogTemplateUseCase(templates repository.LearningLogTemplateRepository) *CreateLearningLogTemplateUseCase {
	return &CreateLearningLogTemplateUseCase{templates: templates}
}

// Execute は前後の空白を落として各項目を検証し、テンプレートを作成する。
// デフォルト指定つきの場合は、先に同一ユーザーの既存のデフォルト指定を外す。
func (uc *CreateLearningLogTemplateUseCase) Execute(ctx context.Context, tmpl *model.LearningLogTemplate) error {
	tmpl.Name = strings.TrimSpace(tmpl.Name)
	tmpl.DefaultTitle = strings.TrimSpace(tmpl.DefaultTitle)
	tmpl.DefaultContent = strings.TrimSpace(tmpl.DefaultContent)

	if err := domain.ValidateStringLength(tmpl.Name, 1, 100, "テンプレート名"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(tmpl.DefaultTitle, 0, 200, "デフォルトタイトル"); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(tmpl.DefaultContent, 0, 50000, "デフォルト本文"); err != nil {
		return err
	}
	if err := validateLearningLogTemplateCategory(tmpl.DefaultCategory); err != nil {
		return err
	}
	if err := validateLearningLogTemplateDuration(tmpl.DefaultDuration); err != nil {
		return err
	}

	if tmpl.IsDefault {
		if err := uc.templates.ClearDefaultFlag(ctx, tmpl.UserID); err != nil {
			return err
		}
	}
	return uc.templates.Create(ctx, tmpl)
}

// GetLearningLogTemplateUseCase は学習ログテンプレートを 1 件取得する。
type GetLearningLogTemplateUseCase struct {
	templates repository.LearningLogTemplateRepository
}

// NewGetLearningLogTemplateUseCase は GetLearningLogTemplateUseCase を生成する。
func NewGetLearningLogTemplateUseCase(templates repository.LearningLogTemplateRepository) *GetLearningLogTemplateUseCase {
	return &GetLearningLogTemplateUseCase{templates: templates}
}

// Execute は所有権を検証したうえでテンプレートを返す。
func (uc *GetLearningLogTemplateUseCase) Execute(ctx context.Context, id, userID uint) (*model.LearningLogTemplate, error) {
	return ensureOwner(ctx, uc.templates.FindByID, id, userID, learningLogTemplateOwnerOf)
}

// ListLearningLogTemplatesUseCase は指定ユーザーのテンプレート一覧を取得する。
type ListLearningLogTemplatesUseCase struct {
	templates repository.LearningLogTemplateRepository
}

// NewListLearningLogTemplatesUseCase は ListLearningLogTemplatesUseCase を生成する。
func NewListLearningLogTemplatesUseCase(templates repository.LearningLogTemplateRepository) *ListLearningLogTemplatesUseCase {
	return &ListLearningLogTemplatesUseCase{templates: templates}
}

// Execute はテンプレートを作成日の新しい順で返す。
func (uc *ListLearningLogTemplatesUseCase) Execute(ctx context.Context, userID uint) ([]model.LearningLogTemplate, error) {
	return uc.templates.FindByUserID(ctx, userID)
}

// GetDefaultLearningLogTemplateUseCase はデフォルトに設定されたテンプレートを取得する。
type GetDefaultLearningLogTemplateUseCase struct {
	templates repository.LearningLogTemplateRepository
}

// NewGetDefaultLearningLogTemplateUseCase は GetDefaultLearningLogTemplateUseCase を生成する。
func NewGetDefaultLearningLogTemplateUseCase(templates repository.LearningLogTemplateRepository) *GetDefaultLearningLogTemplateUseCase {
	return &GetDefaultLearningLogTemplateUseCase{templates: templates}
}

// Execute はデフォルトテンプレートを返す。未設定なら DomainError ではないエラーを返す
// （移行前と同じく 500 になる）。
func (uc *GetDefaultLearningLogTemplateUseCase) Execute(ctx context.Context, userID uint) (*model.LearningLogTemplate, error) {
	tmpl, err := uc.templates.FindDefaultByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if tmpl == nil {
		return nil, errDefaultLearningLogTemplateNotFound
	}
	return tmpl, nil
}

// UpdateLearningLogTemplateInput はテンプレート更新の入力。
// 文字列は空文字列なら「変更なし」、ポインタは nil なら「変更なし」を表す。
type UpdateLearningLogTemplateInput struct {
	ID              uint
	UserID          uint
	Name            string
	DefaultTitle    string
	DefaultContent  string
	DefaultCategory model.LogCategory
	DefaultDuration *int
	IsDefault       *bool
}

// UpdateLearningLogTemplateUseCase は学習ログテンプレートを更新する。
type UpdateLearningLogTemplateUseCase struct {
	templates repository.LearningLogTemplateRepository
}

// NewUpdateLearningLogTemplateUseCase は UpdateLearningLogTemplateUseCase を生成する。
func NewUpdateLearningLogTemplateUseCase(templates repository.LearningLogTemplateRepository) *UpdateLearningLogTemplateUseCase {
	return &UpdateLearningLogTemplateUseCase{templates: templates}
}

// Execute は所有権を検証し、トリム後に空でないフィールドだけを更新する。
// デフォルト指定に true が渡されたときだけ、書き込み前に既存の指定を外す。
func (uc *UpdateLearningLogTemplateUseCase) Execute(ctx context.Context, in UpdateLearningLogTemplateInput) (*model.LearningLogTemplate, error) {
	tmpl, err := ensureOwner(ctx, uc.templates.FindByID, in.ID, in.UserID, learningLogTemplateOwnerOf)
	if err != nil {
		return nil, err
	}

	if name := strings.TrimSpace(in.Name); name != "" {
		if err := domain.ValidateStringLength(name, 1, 100, "テンプレート名"); err != nil {
			return nil, err
		}
		tmpl.Name = name
	}
	if title := strings.TrimSpace(in.DefaultTitle); title != "" {
		if err := domain.ValidateStringLength(title, 1, 200, "デフォルトタイトル"); err != nil {
			return nil, err
		}
		tmpl.DefaultTitle = title
	}
	if content := strings.TrimSpace(in.DefaultContent); content != "" {
		if err := domain.ValidateStringLength(content, 1, 50000, "デフォルト本文"); err != nil {
			return nil, err
		}
		tmpl.DefaultContent = content
	}
	if in.DefaultCategory != "" {
		if err := validateLearningLogTemplateCategory(in.DefaultCategory); err != nil {
			return nil, err
		}
		tmpl.DefaultCategory = in.DefaultCategory
	}
	if in.DefaultDuration != nil {
		if err := validateLearningLogTemplateDuration(*in.DefaultDuration); err != nil {
			return nil, err
		}
		tmpl.DefaultDuration = *in.DefaultDuration
	}
	if in.IsDefault != nil {
		tmpl.IsDefault = *in.IsDefault
		if *in.IsDefault {
			if err := uc.templates.ClearDefaultFlag(ctx, tmpl.UserID); err != nil {
				return nil, err
			}
		}
	}

	if err := uc.templates.Update(ctx, tmpl); err != nil {
		return nil, err
	}
	return tmpl, nil
}

// DeleteLearningLogTemplateUseCase は学習ログテンプレートを削除する。
type DeleteLearningLogTemplateUseCase struct {
	templates repository.LearningLogTemplateRepository
}

// NewDeleteLearningLogTemplateUseCase は DeleteLearningLogTemplateUseCase を生成する。
func NewDeleteLearningLogTemplateUseCase(templates repository.LearningLogTemplateRepository) *DeleteLearningLogTemplateUseCase {
	return &DeleteLearningLogTemplateUseCase{templates: templates}
}

// Execute は所有権を検証したうえでテンプレートを削除する。
func (uc *DeleteLearningLogTemplateUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensureOwner(ctx, uc.templates.FindByID, id, userID, learningLogTemplateOwnerOf); err != nil {
		return err
	}
	return uc.templates.Delete(ctx, id)
}

// CreateLearningLogFromTemplateUseCase はテンプレートから学習ログを作成する。
// 学習ログの検証と書き込みは既存の [CreateLearningLogUseCase] に委ねる。
type CreateLearningLogFromTemplateUseCase struct {
	templates repository.LearningLogTemplateRepository
	createLog *CreateLearningLogUseCase
}

// NewCreateLearningLogFromTemplateUseCase は CreateLearningLogFromTemplateUseCase を生成する。
func NewCreateLearningLogFromTemplateUseCase(
	templates repository.LearningLogTemplateRepository,
	createLog *CreateLearningLogUseCase,
) *CreateLearningLogFromTemplateUseCase {
	return &CreateLearningLogFromTemplateUseCase{templates: templates, createLog: createLog}
}

// Execute は所有権を検証し、テンプレートの内容で学習ログを作成して返す。
// デフォルトタイトルが空のときはテンプレート名を、カテゴリが空のときはその他を使う。
func (uc *CreateLearningLogFromTemplateUseCase) Execute(ctx context.Context, id, userID uint) (*model.LearningLog, error) {
	tmpl, err := ensureOwner(ctx, uc.templates.FindByID, id, userID, learningLogTemplateOwnerOf)
	if err != nil {
		return nil, err
	}

	log := &model.LearningLog{
		UserID:   userID,
		Title:    tmpl.DefaultTitle,
		Content:  tmpl.DefaultContent,
		Category: tmpl.DefaultCategory,
		Duration: tmpl.DefaultDuration,
		Source:   model.LogSourceManual,
	}
	if log.Title == "" {
		log.Title = tmpl.Name
	}
	if log.Category == "" {
		log.Category = model.LogCategoryOther
	}

	if err := uc.createLog.Execute(ctx, log); err != nil {
		return nil, err
	}
	return log, nil
}

// CountLearningLogTemplatesUseCase は学習ログテンプレートの総数を取得する。
type CountLearningLogTemplatesUseCase struct {
	templates repository.LearningLogTemplateRepository
}

// NewCountLearningLogTemplatesUseCase は CountLearningLogTemplatesUseCase を生成する。
func NewCountLearningLogTemplatesUseCase(templates repository.LearningLogTemplateRepository) *CountLearningLogTemplatesUseCase {
	return &CountLearningLogTemplatesUseCase{templates: templates}
}

// Execute はテンプレート総数を返す。
func (uc *CountLearningLogTemplatesUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.templates.CountByUserID(ctx, userID)
}

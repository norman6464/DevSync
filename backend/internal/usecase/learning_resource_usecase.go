package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/domain/validator"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// errLearningResourceNotFound は port が「不在」を表す nil を返したときに返すエラー。
// DomainError ではないため handler では 500 になり、リポジトリの生エラーが
// そのまま返っていた移行前の挙動と一致する。
var errLearningResourceNotFound = errors.New("学習リソースが見つかりません")

// validResourceDifficulties は有効な難易度のマップ。
var validResourceDifficulties = map[string]bool{
	string(model.ResourceDifficultyBeginner):     true,
	string(model.ResourceDifficultyIntermediate): true,
	string(model.ResourceDifficultyAdvanced):     true,
}

// ownerOfLearningResource は所有権判定に使う所有者 ID の取り出し。
func ownerOfLearningResource(r *model.LearningResource) uint { return r.UserID }

// requireOthersResource は対象リソースが自分のものでないことを検証する。
// 不在の場合は 404、自分のリソースなら 403 を返す。
func requireOthersResource(ctx context.Context, resources repository.LearningResourceRepository, userID, resourceID uint) error {
	resource, err := resources.FindByID(ctx, resourceID)
	if err != nil || resource == nil {
		return domain.ErrNotFound
	}
	if resource.UserID == userID {
		return domain.ErrForbidden
	}
	return nil
}

// CreateLearningResourceUseCase は学習リソースを作成する。
type CreateLearningResourceUseCase struct {
	resources repository.LearningResourceRepository
}

// NewCreateLearningResourceUseCase は CreateLearningResourceUseCase を生成する。
func NewCreateLearningResourceUseCase(resources repository.LearningResourceRepository) *CreateLearningResourceUseCase {
	return &CreateLearningResourceUseCase{resources: resources}
}

// Execute は入力を検証したうえで学習リソースを作成する。
func (uc *CreateLearningResourceUseCase) Execute(ctx context.Context, resource *model.LearningResource) error {
	v := validator.NewResourceValidator()
	if err := v.ValidateCreateResource(resource.Title, resource.Description, resource.URL,
		string(resource.Category), string(resource.Difficulty)); err != nil {
		return err
	}
	if err := domain.ValidateStringLength(resource.Tags, 0, 1000, "タグ"); err != nil {
		return err
	}
	return uc.resources.Create(ctx, resource)
}

// GetLearningResourceUseCase は学習リソースを 1 件取得する。
type GetLearningResourceUseCase struct {
	resources repository.LearningResourceRepository
}

// NewGetLearningResourceUseCase は GetLearningResourceUseCase を生成する。
func NewGetLearningResourceUseCase(resources repository.LearningResourceRepository) *GetLearningResourceUseCase {
	return &GetLearningResourceUseCase{resources: resources}
}

// Execute はリソースを返す。非公開のものは所有者しか取得できない。
func (uc *GetLearningResourceUseCase) Execute(ctx context.Context, id, userID uint) (*model.LearningResource, error) {
	resource, err := uc.resources.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if resource == nil {
		return nil, errLearningResourceNotFound
	}
	if !resource.IsPublic && resource.UserID != userID {
		return nil, domain.ErrForbidden
	}
	return resource, nil
}

// ListLearningResourcesByUserUseCase は指定ユーザーの学習リソース一覧を取得する。
type ListLearningResourcesByUserUseCase struct {
	resources repository.LearningResourceRepository
}

// NewListLearningResourcesByUserUseCase は ListLearningResourcesByUserUseCase を生成する。
func NewListLearningResourcesByUserUseCase(resources repository.LearningResourceRepository) *ListLearningResourcesByUserUseCase {
	return &ListLearningResourcesByUserUseCase{resources: resources}
}

// Execute は指定ユーザーのリソースを返す。自分の一覧なら非公開も含める。
func (uc *ListLearningResourcesByUserUseCase) Execute(ctx context.Context, targetUserID, currentUserID uint, limit, offset int) ([]model.LearningResource, int64, error) {
	includePrivate := currentUserID == targetUserID
	return uc.resources.FindByUserID(ctx, targetUserID, includePrivate, limit, offset)
}

// ListPublicLearningResourcesUseCase は公開リソース一覧を取得する。
type ListPublicLearningResourcesUseCase struct {
	resources repository.LearningResourceRepository
}

// NewListPublicLearningResourcesUseCase は ListPublicLearningResourcesUseCase を生成する。
func NewListPublicLearningResourcesUseCase(resources repository.LearningResourceRepository) *ListPublicLearningResourcesUseCase {
	return &ListPublicLearningResourcesUseCase{resources: resources}
}

// Execute はカテゴリ・難易度で絞り込んだ公開リソースを返す。
func (uc *ListPublicLearningResourcesUseCase) Execute(ctx context.Context, limit, offset int, category, difficulty string) ([]model.LearningResource, int64, error) {
	return uc.resources.FindPublic(ctx, limit, offset, category, difficulty)
}

// ListLearningResourcesByDifficultyUseCase は難易度別に公開リソースを取得する。
type ListLearningResourcesByDifficultyUseCase struct {
	resources repository.LearningResourceRepository
}

// NewListLearningResourcesByDifficultyUseCase は ListLearningResourcesByDifficultyUseCase を生成する。
func NewListLearningResourcesByDifficultyUseCase(resources repository.LearningResourceRepository) *ListLearningResourcesByDifficultyUseCase {
	return &ListLearningResourcesByDifficultyUseCase{resources: resources}
}

// Execute は指定難易度の公開リソースを返す。未知の難易度は 400。
func (uc *ListLearningResourcesByDifficultyUseCase) Execute(ctx context.Context, difficulty string, limit, offset int) ([]model.LearningResource, int64, error) {
	if !validResourceDifficulties[difficulty] {
		return nil, 0, domain.NewError(domain.ErrCodeBadRequest, "無効な難易度です", nil)
	}
	return uc.resources.FindByDifficulty(ctx, difficulty, limit, offset)
}

// SearchLearningResourcesUseCase は学習リソースをキーワード検索する。
type SearchLearningResourcesUseCase struct {
	resources repository.LearningResourceRepository
}

// NewSearchLearningResourcesUseCase は SearchLearningResourcesUseCase を生成する。
func NewSearchLearningResourcesUseCase(resources repository.LearningResourceRepository) *SearchLearningResourcesUseCase {
	return &SearchLearningResourcesUseCase{resources: resources}
}

// Execute は公開リソースをタイトル・説明・タグへの部分一致で検索する。
func (uc *SearchLearningResourcesUseCase) Execute(ctx context.Context, query string, limit, offset int) ([]model.LearningResource, int64, error) {
	return uc.resources.Search(ctx, query, limit, offset)
}

// UpdateLearningResourceUseCase は学習リソースを更新する。
type UpdateLearningResourceUseCase struct {
	resources repository.LearningResourceRepository
}

// NewUpdateLearningResourceUseCase は UpdateLearningResourceUseCase を生成する。
func NewUpdateLearningResourceUseCase(resources repository.LearningResourceRepository) *UpdateLearningResourceUseCase {
	return &UpdateLearningResourceUseCase{resources: resources}
}

// Execute は学習リソースを部分更新する。所有者のみ。
// 前後の空白を除いて空になった項目は「変更なし」として扱う。
func (uc *UpdateLearningResourceUseCase) Execute(ctx context.Context, id, userID uint, updates *model.LearningResource) (*model.LearningResource, error) {
	resource, err := ensureOwner(ctx, uc.resources.FindByID, id, userID, ownerOfLearningResource)
	if err != nil {
		return nil, err
	}

	// 空白だけの入力で検証をすり抜けないよう、先に正規化する。
	updates.Title = strings.TrimSpace(updates.Title)
	updates.Description = strings.TrimSpace(updates.Description)
	updates.URL = strings.TrimSpace(updates.URL)
	updates.Tags = strings.TrimSpace(updates.Tags)
	updates.ImageURL = strings.TrimSpace(updates.ImageURL)
	category := strings.TrimSpace(string(updates.Category))
	difficulty := strings.TrimSpace(string(updates.Difficulty))

	v := validator.NewResourceValidator()
	if err := v.ValidateUpdateResource(updates.Title, updates.Description, updates.URL, category, difficulty); err != nil {
		return nil, err
	}

	if updates.Title != "" {
		resource.Title = updates.Title
	}
	if updates.Description != "" {
		resource.Description = updates.Description
	}
	if updates.URL != "" {
		resource.URL = updates.URL
	}
	if category != "" {
		resource.Category = model.ResourceCategory(category)
	}
	if difficulty != "" {
		resource.Difficulty = model.ResourceDifficulty(difficulty)
	}
	if updates.Tags != "" {
		if err := domain.ValidateStringLength(updates.Tags, 1, 1000, "タグ"); err != nil {
			return nil, err
		}
		resource.Tags = updates.Tags
	}
	if updates.ImageURL != "" {
		if err := domain.ValidateStringLength(updates.ImageURL, 1, 2000, "画像URL"); err != nil {
			return nil, err
		}
		resource.ImageURL = updates.ImageURL
	}

	if err := uc.resources.Update(ctx, resource); err != nil {
		return nil, err
	}
	return resource, nil
}

// UpdateLearningResourceVisibilityUseCase は学習リソースの公開/非公開を切り替える。
type UpdateLearningResourceVisibilityUseCase struct {
	resources repository.LearningResourceRepository
}

// NewUpdateLearningResourceVisibilityUseCase は UpdateLearningResourceVisibilityUseCase を生成する。
func NewUpdateLearningResourceVisibilityUseCase(resources repository.LearningResourceRepository) *UpdateLearningResourceVisibilityUseCase {
	return &UpdateLearningResourceVisibilityUseCase{resources: resources}
}

// Execute は公開状態を更新する。所有者のみ。
func (uc *UpdateLearningResourceVisibilityUseCase) Execute(ctx context.Context, id, userID uint, isPublic bool) (*model.LearningResource, error) {
	resource, err := ensureOwner(ctx, uc.resources.FindByID, id, userID, ownerOfLearningResource)
	if err != nil {
		return nil, err
	}

	resource.IsPublic = isPublic
	if err := uc.resources.Update(ctx, resource); err != nil {
		return nil, err
	}
	return resource, nil
}

// DeleteLearningResourceUseCase は学習リソースを削除する。
type DeleteLearningResourceUseCase struct {
	resources repository.LearningResourceRepository
}

// NewDeleteLearningResourceUseCase は DeleteLearningResourceUseCase を生成する。
func NewDeleteLearningResourceUseCase(resources repository.LearningResourceRepository) *DeleteLearningResourceUseCase {
	return &DeleteLearningResourceUseCase{resources: resources}
}

// Execute は学習リソースを削除する。所有者のみ。
func (uc *DeleteLearningResourceUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensureOwner(ctx, uc.resources.FindByID, id, userID, ownerOfLearningResource); err != nil {
		return err
	}
	return uc.resources.Delete(ctx, id)
}

// LikeLearningResourceUseCase は学習リソースにいいねする。
type LikeLearningResourceUseCase struct {
	resources repository.LearningResourceRepository
}

// NewLikeLearningResourceUseCase は LikeLearningResourceUseCase を生成する。
func NewLikeLearningResourceUseCase(resources repository.LearningResourceRepository) *LikeLearningResourceUseCase {
	return &LikeLearningResourceUseCase{resources: resources}
}

// Execute はいいねを追加する。自分のリソースにはいいねできない。
func (uc *LikeLearningResourceUseCase) Execute(ctx context.Context, userID, resourceID uint) error {
	if err := requireOthersResource(ctx, uc.resources, userID, resourceID); err != nil {
		return err
	}
	return uc.resources.Like(ctx, userID, resourceID)
}

// UnlikeLearningResourceUseCase は学習リソースのいいねを取り消す。
type UnlikeLearningResourceUseCase struct {
	resources repository.LearningResourceRepository
}

// NewUnlikeLearningResourceUseCase は UnlikeLearningResourceUseCase を生成する。
func NewUnlikeLearningResourceUseCase(resources repository.LearningResourceRepository) *UnlikeLearningResourceUseCase {
	return &UnlikeLearningResourceUseCase{resources: resources}
}

// Execute はいいねを取り消す。自分のリソースはそもそもいいねできないため 403 を返す。
func (uc *UnlikeLearningResourceUseCase) Execute(ctx context.Context, userID, resourceID uint) error {
	if err := requireOthersResource(ctx, uc.resources, userID, resourceID); err != nil {
		return err
	}
	return uc.resources.Unlike(ctx, userID, resourceID)
}

// HasLikedLearningResourceUseCase はいいね済みかを判定する。
type HasLikedLearningResourceUseCase struct {
	resources repository.LearningResourceRepository
}

// NewHasLikedLearningResourceUseCase は HasLikedLearningResourceUseCase を生成する。
func NewHasLikedLearningResourceUseCase(resources repository.LearningResourceRepository) *HasLikedLearningResourceUseCase {
	return &HasLikedLearningResourceUseCase{resources: resources}
}

// Execute は指定ユーザーがいいね済みかを返す。
func (uc *HasLikedLearningResourceUseCase) Execute(ctx context.Context, userID, resourceID uint) (bool, error) {
	return uc.resources.HasLiked(ctx, userID, resourceID)
}

// SaveLearningResourceUseCase は学習リソースを保存する。
type SaveLearningResourceUseCase struct {
	resources repository.LearningResourceRepository
}

// NewSaveLearningResourceUseCase は SaveLearningResourceUseCase を生成する。
func NewSaveLearningResourceUseCase(resources repository.LearningResourceRepository) *SaveLearningResourceUseCase {
	return &SaveLearningResourceUseCase{resources: resources}
}

// Execute は保存を追加する。自分のリソースは保存できない。
func (uc *SaveLearningResourceUseCase) Execute(ctx context.Context, userID, resourceID uint) error {
	if err := requireOthersResource(ctx, uc.resources, userID, resourceID); err != nil {
		return err
	}
	return uc.resources.Save(ctx, userID, resourceID)
}

// UnsaveLearningResourceUseCase は学習リソースの保存を取り消す。
type UnsaveLearningResourceUseCase struct {
	resources repository.LearningResourceRepository
}

// NewUnsaveLearningResourceUseCase は UnsaveLearningResourceUseCase を生成する。
func NewUnsaveLearningResourceUseCase(resources repository.LearningResourceRepository) *UnsaveLearningResourceUseCase {
	return &UnsaveLearningResourceUseCase{resources: resources}
}

// Execute は保存を取り消す。自分のリソースはそもそも保存できないため 403 を返す。
func (uc *UnsaveLearningResourceUseCase) Execute(ctx context.Context, userID, resourceID uint) error {
	if err := requireOthersResource(ctx, uc.resources, userID, resourceID); err != nil {
		return err
	}
	return uc.resources.Unsave(ctx, userID, resourceID)
}

// HasSavedLearningResourceUseCase は保存済みかを判定する。
type HasSavedLearningResourceUseCase struct {
	resources repository.LearningResourceRepository
}

// NewHasSavedLearningResourceUseCase は HasSavedLearningResourceUseCase を生成する。
func NewHasSavedLearningResourceUseCase(resources repository.LearningResourceRepository) *HasSavedLearningResourceUseCase {
	return &HasSavedLearningResourceUseCase{resources: resources}
}

// Execute は指定ユーザーが保存済みかを返す。
func (uc *HasSavedLearningResourceUseCase) Execute(ctx context.Context, userID, resourceID uint) (bool, error) {
	return uc.resources.HasSaved(ctx, userID, resourceID)
}

// ListSavedLearningResourcesUseCase は保存済みリソース一覧を取得する。
type ListSavedLearningResourcesUseCase struct {
	resources repository.LearningResourceRepository
}

// NewListSavedLearningResourcesUseCase は ListSavedLearningResourcesUseCase を生成する。
func NewListSavedLearningResourcesUseCase(resources repository.LearningResourceRepository) *ListSavedLearningResourcesUseCase {
	return &ListSavedLearningResourcesUseCase{resources: resources}
}

// Execute は保存済みリソースを新しい順で返す。
func (uc *ListSavedLearningResourcesUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.LearningResource, int64, error) {
	return uc.resources.FindSavedByUserID(ctx, userID, limit, offset)
}

// CountLearningResourcesUseCase は学習リソース数を取得する。
type CountLearningResourcesUseCase struct {
	resources repository.LearningResourceRepository
}

// NewCountLearningResourcesUseCase は CountLearningResourcesUseCase を生成する。
func NewCountLearningResourcesUseCase(resources repository.LearningResourceRepository) *CountLearningResourcesUseCase {
	return &CountLearningResourcesUseCase{resources: resources}
}

// Execute は指定ユーザーの学習リソース総数を返す。
func (uc *CountLearningResourcesUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.resources.CountByUserID(ctx, userID)
}

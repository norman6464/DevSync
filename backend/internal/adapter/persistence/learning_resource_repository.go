package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// learningResourceRepository は [repository.LearningResourceRepository] の sqlc(pgx) 実装。
type learningResourceRepository struct {
	q *sqlcgen.Queries
}

// NewLearningResourceRepository は LearningResourceRepository の sqlc(pgx) 実装を返す。
func NewLearningResourceRepository(q *sqlcgen.Queries) repository.LearningResourceRepository {
	return &learningResourceRepository{q: q}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.LearningResourceRepository = (*learningResourceRepository)(nil)

// Create は新しい学習リソースを作成する。
func (r *learningResourceRepository) Create(ctx context.Context, resource *model.LearningResource) error {
	row, err := r.q.CreateLearningResource(ctx, sqlcgen.CreateLearningResourceParams{
		UserID:      int64(resource.UserID),
		Title:       resource.Title,
		Description: &resource.Description,
		Url:         &resource.URL,
		Category:    string(resource.Category),
		Difficulty:  (*string)(&resource.Difficulty),
		Tags:        &resource.Tags,
		ImageUrl:    &resource.ImageURL,
		IsPublic:    resource.IsPublic,
		LikeCount:   toInt64Ptr(resource.LikeCount),
		SaveCount:   toInt64Ptr(resource.SaveCount),
	})
	if err != nil {
		return err
	}
	*resource = toModelLearningResource(row)
	return nil
}

// Update は既存の学習リソースを更新する（GORMのSave＝全カラム上書きに相当）。
func (r *learningResourceRepository) Update(ctx context.Context, resource *model.LearningResource) error {
	row, err := r.q.UpdateLearningResource(ctx, sqlcgen.UpdateLearningResourceParams{
		ID:          int64(resource.ID),
		Title:       resource.Title,
		Description: &resource.Description,
		Url:         &resource.URL,
		Category:    string(resource.Category),
		Difficulty:  (*string)(&resource.Difficulty),
		Tags:        &resource.Tags,
		ImageUrl:    &resource.ImageURL,
		IsPublic:    resource.IsPublic,
	})
	if err != nil {
		return err
	}
	*resource = toModelLearningResource(row)
	return nil
}

// Delete は学習リソースを削除する。依存するいいね・保存等はFKのON DELETE CASCADEで
// DBが自動的に削除する。
func (r *learningResourceRepository) Delete(ctx context.Context, id uint) error {
	return r.q.DeleteLearningResource(ctx, int64(id))
}

// FindByID は指定 ID の学習リソースをユーザー情報付きで取得する。不在の場合は (nil, nil) を返す。
func (r *learningResourceRepository) FindByID(ctx context.Context, id uint) (*model.LearningResource, error) {
	row, err := r.q.GetLearningResourceWithUserByID(ctx, int64(id))
	if isNoRows(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	resource := toModelLearningResource(row.LearningResource)
	resource.User = toModelUser(row.User)
	return &resource, nil
}

// FindByUserID は指定ユーザーのリソースを取得する（新しい順）。
// 一覧系の中でこれだけユーザー情報をプリロードしない（移行前からの挙動）。
func (r *learningResourceRepository) FindByUserID(ctx context.Context, userID uint, includePrivate bool, limit, offset int) ([]model.LearningResource, int64, error) {
	total, err := r.q.CountUserVisibleLearningResources(ctx, sqlcgen.CountUserVisibleLearningResourcesParams{
		UserID:         int64(userID),
		IncludePrivate: includePrivate,
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListLearningResourcesByUser(ctx, sqlcgen.ListLearningResourcesByUserParams{
		UserID:         int64(userID),
		Limit:          int32Param(limit),
		Offset:         int32Param(offset),
		IncludePrivate: includePrivate,
	})
	if err != nil {
		return nil, 0, err
	}

	resources := make([]model.LearningResource, len(rows))
	for i, row := range rows {
		resources[i] = toModelLearningResource(row)
	}
	return resources, total, nil
}

// FindPublic は公開リソースをカテゴリ・難易度で絞り込んで取得する（いいね数降順 → 新しい順）。
func (r *learningResourceRepository) FindPublic(ctx context.Context, limit, offset int, category, difficulty string) ([]model.LearningResource, int64, error) {
	total, err := r.q.CountPublicLearningResources(ctx, sqlcgen.CountPublicLearningResourcesParams{
		Category:   nilIfEmpty(category),
		Difficulty: nilIfEmpty(difficulty),
	})
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListPublicLearningResourcesWithUser(ctx, sqlcgen.ListPublicLearningResourcesWithUserParams{
		Limit:      int32Param(limit),
		Offset:     int32Param(offset),
		Category:   nilIfEmpty(category),
		Difficulty: nilIfEmpty(difficulty),
	})
	if err != nil {
		return nil, 0, err
	}

	resources := make([]model.LearningResource, len(rows))
	for i, row := range rows {
		resources[i] = toModelLearningResource(row.LearningResource)
		resources[i].User = toModelUser(row.User)
	}
	return resources, total, nil
}

// FindByDifficulty は公開リソースを難易度で絞り込んで取得する（新しい順）。
func (r *learningResourceRepository) FindByDifficulty(ctx context.Context, difficulty string, limit, offset int) ([]model.LearningResource, int64, error) {
	total, err := r.q.CountLearningResourcesByDifficulty(ctx, &difficulty)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListLearningResourcesByDifficultyWithUser(ctx, sqlcgen.ListLearningResourcesByDifficultyWithUserParams{
		Difficulty: &difficulty,
		Limit:      int32Param(limit),
		Offset:     int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	resources := make([]model.LearningResource, len(rows))
	for i, row := range rows {
		resources[i] = toModelLearningResource(row.LearningResource)
		resources[i].User = toModelUser(row.User)
	}
	return resources, total, nil
}

// Search は公開リソースをタイトル・説明・タグで部分一致検索する（いいね数降順 → 新しい順）。
func (r *learningResourceRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.LearningResource, int64, error) {
	pattern := escapeLikePattern(query)

	total, err := r.q.CountSearchLearningResources(ctx, pattern)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.SearchLearningResourcesWithUser(ctx, sqlcgen.SearchLearningResourcesWithUserParams{
		Title:  pattern,
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	resources := make([]model.LearningResource, len(rows))
	for i, row := range rows {
		resources[i] = toModelLearningResource(row.LearningResource)
		resources[i].User = toModelUser(row.User)
	}
	return resources, total, nil
}

// FindSavedByUserID は指定ユーザーが保存したリソースを取得する（新しい順）。
func (r *learningResourceRepository) FindSavedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.LearningResource, int64, error) {
	total, err := r.q.CountSavedLearningResources(ctx, int64(userID))
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.q.ListSavedLearningResourcesWithUser(ctx, sqlcgen.ListSavedLearningResourcesWithUserParams{
		UserID: int64(userID),
		Limit:  int32Param(limit),
		Offset: int32Param(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	resources := make([]model.LearningResource, len(rows))
	for i, row := range rows {
		resources[i] = toModelLearningResource(row.LearningResource)
		resources[i].User = toModelUser(row.User)
	}
	return resources, total, nil
}

// CountByUserID は指定ユーザーのリソース総数を返す。
func (r *learningResourceRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	return r.q.CountLearningResourcesByUser(ctx, int64(userID))
}

// Like はいいねを追加し、リソースのいいね数を 1 増やす。
// 移行前の GORM 実装と同じくトランザクションでは括らない（元実装も2つの独立した操作だったため）。
func (r *learningResourceRepository) Like(ctx context.Context, userID, resourceID uint) error {
	if err := r.q.CreateResourceLike(ctx, sqlcgen.CreateResourceLikeParams{
		UserID:     int64(userID),
		ResourceID: int64(resourceID),
	}); err != nil {
		return err
	}
	return r.q.IncrementResourceLikeCount(ctx, int64(resourceID))
}

// Unlike はいいねを取り消し、リソースのいいね数を 1 減らす（0 未満にはしない）。
func (r *learningResourceRepository) Unlike(ctx context.Context, userID, resourceID uint) error {
	if err := r.q.DeleteResourceLike(ctx, sqlcgen.DeleteResourceLikeParams{
		UserID:     int64(userID),
		ResourceID: int64(resourceID),
	}); err != nil {
		return err
	}
	return r.q.DecrementResourceLikeCountFloored(ctx, int64(resourceID))
}

// HasLiked は指定ユーザーがいいね済みかを返す。
func (r *learningResourceRepository) HasLiked(ctx context.Context, userID, resourceID uint) (bool, error) {
	count, err := r.q.CountResourceLike(ctx, sqlcgen.CountResourceLikeParams{
		UserID:     int64(userID),
		ResourceID: int64(resourceID),
	})
	return count > 0, err
}

// Save は保存を追加し、リソースの保存数を 1 増やす。
func (r *learningResourceRepository) Save(ctx context.Context, userID, resourceID uint) error {
	if err := r.q.CreateResourceSave(ctx, sqlcgen.CreateResourceSaveParams{
		UserID:     int64(userID),
		ResourceID: int64(resourceID),
	}); err != nil {
		return err
	}
	return r.q.IncrementResourceSaveCount(ctx, int64(resourceID))
}

// Unsave は保存を取り消し、リソースの保存数を 1 減らす（0 未満にはしない）。
func (r *learningResourceRepository) Unsave(ctx context.Context, userID, resourceID uint) error {
	if err := r.q.DeleteResourceSave(ctx, sqlcgen.DeleteResourceSaveParams{
		UserID:     int64(userID),
		ResourceID: int64(resourceID),
	}); err != nil {
		return err
	}
	return r.q.DecrementResourceSaveCountFloored(ctx, int64(resourceID))
}

// HasSaved は指定ユーザーが保存済みかを返す。
func (r *learningResourceRepository) HasSaved(ctx context.Context, userID, resourceID uint) (bool, error) {
	count, err := r.q.CountResourceSave(ctx, sqlcgen.CountResourceSaveParams{
		UserID:     int64(userID),
		ResourceID: int64(resourceID),
	})
	return count > 0, err
}

package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// maxTagsPerPost は 1 投稿に設定できるタグの上限。
const maxTagsPerPost = 10

// postOwnerOf は投稿の所有者 ID を返す。
func postTagPostOwnerOf(p *model.Post) uint { return p.UserID }

// SetPostTagsUseCase は投稿のタグを設定する。
type SetPostTagsUseCase struct {
	tags  repository.PostTagRepository
	posts repository.PostReader
}

// NewSetPostTagsUseCase は SetPostTagsUseCase を生成する。
func NewSetPostTagsUseCase(tags repository.PostTagRepository, posts repository.PostReader) *SetPostTagsUseCase {
	return &SetPostTagsUseCase{tags: tags, posts: posts}
}

// Execute は所有権とタグの妥当性を検査したうえで、投稿のタグを置き換える。
func (uc *SetPostTagsUseCase) Execute(ctx context.Context, postID, userID uint, tags []string) error {
	if _, err := ensureOwner(ctx, uc.posts.FindByID, postID, userID, postTagPostOwnerOf); err != nil {
		return err
	}

	normalized := domain.NormalizeTags(tags)

	if len(normalized) > maxTagsPerPost {
		return domain.NewError(domain.ErrCodeBadRequest, "タグは最大10個までです", nil)
	}
	for _, tag := range normalized {
		if err := domain.ValidateStringLength(tag, 1, 50, "タグ"); err != nil {
			return err
		}
	}

	return uc.tags.SetTags(ctx, postID, normalized)
}

// SetAutoPostTagsUseCase は本文から抽出したハッシュタグを投稿のタグとして設定する。
type SetAutoPostTagsUseCase struct {
	setTags *SetPostTagsUseCase
}

// NewSetAutoPostTagsUseCase は SetAutoPostTagsUseCase を生成する。
func NewSetAutoPostTagsUseCase(setTags *SetPostTagsUseCase) *SetAutoPostTagsUseCase {
	return &SetAutoPostTagsUseCase{setTags: setTags}
}

// Execute は本文からハッシュタグを抽出して設定する。
// ハッシュタグが 1 つも無ければ、既存のタグには触れずに何もしない。
func (uc *SetAutoPostTagsUseCase) Execute(ctx context.Context, postID, userID uint, content string) error {
	tags := domain.ExtractHashtags(content)
	if len(tags) == 0 {
		return nil
	}
	return uc.setTags.Execute(ctx, postID, userID, tags)
}

// GetPostTagsUseCase は投稿のタグ一覧を取得する。
type GetPostTagsUseCase struct {
	tags repository.PostTagRepository
}

// NewGetPostTagsUseCase は GetPostTagsUseCase を生成する。
func NewGetPostTagsUseCase(tags repository.PostTagRepository) *GetPostTagsUseCase {
	return &GetPostTagsUseCase{tags: tags}
}

// Execute はタグ一覧を返す。所有権は検証しない（移行前の挙動を維持している）。
func (uc *GetPostTagsUseCase) Execute(ctx context.Context, postID uint) ([]string, error) {
	return uc.tags.GetByPostID(ctx, postID)
}

// FindPostsByTagUseCase はタグで投稿を検索する。
type FindPostsByTagUseCase struct {
	tags repository.PostTagRepository
}

// NewFindPostsByTagUseCase は FindPostsByTagUseCase を生成する。
func NewFindPostsByTagUseCase(tags repository.PostTagRepository) *FindPostsByTagUseCase {
	return &FindPostsByTagUseCase{tags: tags}
}

// Execute はタグに紐づく投稿と総件数を返す。
func (uc *FindPostsByTagUseCase) Execute(ctx context.Context, tag string, limit, offset int) ([]model.Post, int64, error) {
	return uc.tags.FindPostsByTag(ctx, tag, limit, offset)
}

// GetPopularTagsUseCase は人気タグ一覧を取得する。
type GetPopularTagsUseCase struct {
	tags repository.PostTagRepository
}

// NewGetPopularTagsUseCase は GetPopularTagsUseCase を生成する。
func NewGetPopularTagsUseCase(tags repository.PostTagRepository) *GetPopularTagsUseCase {
	return &GetPopularTagsUseCase{tags: tags}
}

// Execute は使用回数の多い順にタグを返す。
func (uc *GetPopularTagsUseCase) Execute(ctx context.Context, limit int) ([]model.TagCount, error) {
	return uc.tags.GetPopularTags(ctx, limit)
}

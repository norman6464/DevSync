package service

import (
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

const maxTagsPerPost = 10

// PostTagService は投稿タグのビジネスロジックを提供する。
type PostTagService struct {
	tagRepo  repository.PostTagRepositoryInterface
	postRepo repository.PostRepositoryInterface
}

// NewPostTagService は新しいPostTagServiceインスタンスを生成する。
func NewPostTagService(tagRepo repository.PostTagRepositoryInterface, postRepo repository.PostRepositoryInterface) *PostTagService {
	return &PostTagService{tagRepo: tagRepo, postRepo: postRepo}
}

// SetTags は投稿のタグを設定する。所有権チェック・バリデーション付き。
func (s *PostTagService) SetTags(postID, userID uint, tags []string) error {
	if _, err := checkOwnership(s.postRepo.FindByID, postID, userID, func(p *model.Post) uint { return p.UserID }); err != nil {
		return err
	}

	// 正規化: 小文字変換・トリム・空文字除外・重複除外
	normalized := normalizeTags(tags)

	if len(normalized) > maxTagsPerPost {
		return domain.NewError(domain.ErrCodeBadRequest, "タグは最大10個までです", nil)
	}

	for _, tag := range normalized {
		if err := domain.ValidateStringLength(tag, 1, 50, "タグ"); err != nil {
			return err
		}
	}

	return s.tagRepo.SetTags(postID, normalized)
}

// GetByPostID は投稿のタグ一覧を取得する。
func (s *PostTagService) GetByPostID(postID uint) ([]string, error) {
	return s.tagRepo.GetByPostID(postID)
}

// FindPostsByTag はタグで投稿を検索する。
func (s *PostTagService) FindPostsByTag(tag string, limit, offset int) ([]model.Post, int64, error) {
	return s.tagRepo.FindPostsByTag(tag, limit, offset)
}

// GetPopularTags は人気タグ一覧を取得する。
func (s *PostTagService) GetPopularTags(limit int) ([]model.TagCount, error) {
	return s.tagRepo.GetPopularTags(limit)
}

// SetAutoTags はコンテンツからハッシュタグを自動抽出し、投稿のタグとして設定する。
// ハッシュタグが見つからない場合は何もしない。
func (s *PostTagService) SetAutoTags(postID, userID uint, content string) error {
	tags := ExtractHashtags(content)
	if len(tags) == 0 {
		return nil
	}
	return s.SetTags(postID, userID, tags)
}

// normalizeTags はタグを正規化する（小文字変換・トリム・空文字除外・重複除外）。
func normalizeTags(tags []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, tag := range tags {
		t := strings.ToLower(strings.TrimSpace(tag))
		if t == "" || seen[t] {
			continue
		}
		seen[t] = true
		result = append(result, t)
	}
	return result
}

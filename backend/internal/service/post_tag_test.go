package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newTestPostTagService はPostTagServiceのテスト用インスタンスを生成するヘルパー。
func newTestPostTagService() (*PostTagService, *MockPostTagRepository, *MockPostRepository) {
	tagRepo := new(MockPostTagRepository)
	postRepo := new(MockPostRepository)
	svc := NewPostTagService(tagRepo, postRepo)
	return svc, tagRepo, postRepo
}

// ============================================================
// SetTags（タグ設定 — 所有権チェック + バリデーション）
// ============================================================

func TestPostTagSetTags_Success(t *testing.T) {
	svc, tagRepo, postRepo := newTestPostTagService()

	post := &model.Post{UserID: 1}
	post.ID = 10

	postRepo.On("FindByID", uint(10)).Return(post, nil)
	tagRepo.On("SetTags", uint(10), []string{"go", "react"}).Return(nil)

	err := svc.SetTags(10, 1, []string{"Go", "React"})
	assert.NoError(t, err)
	tagRepo.AssertExpectations(t)
	postRepo.AssertExpectations(t)
}

func TestPostTagSetTags_Forbidden(t *testing.T) {
	svc, _, postRepo := newTestPostTagService()

	post := &model.Post{UserID: 1}
	post.ID = 10

	postRepo.On("FindByID", uint(10)).Return(post, nil)

	err := svc.SetTags(10, 999, []string{"go"})
	assert.ErrorIs(t, err, ErrForbidden)
	postRepo.AssertExpectations(t)
}

func TestPostTagSetTags_PostNotFound(t *testing.T) {
	svc, _, postRepo := newTestPostTagService()

	postRepo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	err := svc.SetTags(999, 1, []string{"go"})
	assert.Error(t, err)
	postRepo.AssertExpectations(t)
}

func TestPostTagSetTags_TooManyTags(t *testing.T) {
	svc, _, postRepo := newTestPostTagService()

	post := &model.Post{UserID: 1}
	post.ID = 10
	postRepo.On("FindByID", uint(10)).Return(post, nil)

	tags := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}
	err := svc.SetTags(10, 1, tags)
	assert.Error(t, err) // 最大10個
}

func TestPostTagSetTags_EmptyTagFiltered(t *testing.T) {
	svc, tagRepo, postRepo := newTestPostTagService()

	post := &model.Post{UserID: 1}
	post.ID = 10

	postRepo.On("FindByID", uint(10)).Return(post, nil)
	tagRepo.On("SetTags", uint(10), []string{"go"}).Return(nil)

	// 空文字タグはフィルタされる
	err := svc.SetTags(10, 1, []string{"go", "", "  "})
	assert.NoError(t, err)
	tagRepo.AssertExpectations(t)
}

func TestPostTagSetTags_DuplicateTagDeduped(t *testing.T) {
	svc, tagRepo, postRepo := newTestPostTagService()

	post := &model.Post{UserID: 1}
	post.ID = 10

	postRepo.On("FindByID", uint(10)).Return(post, nil)
	tagRepo.On("SetTags", uint(10), []string{"go"}).Return(nil)

	// 重複タグは1つにまとめられる
	err := svc.SetTags(10, 1, []string{"Go", "go", "GO"})
	assert.NoError(t, err)
	tagRepo.AssertExpectations(t)
}

func TestPostTagSetTags_RepoError(t *testing.T) {
	svc, tagRepo, postRepo := newTestPostTagService()

	post := &model.Post{UserID: 1}
	post.ID = 10

	postRepo.On("FindByID", uint(10)).Return(post, nil)
	tagRepo.On("SetTags", uint(10), mock.Anything).Return(errors.New("db error"))

	err := svc.SetTags(10, 1, []string{"go"})
	assert.Error(t, err)
}

// ============================================================
// GetByPostID
// ============================================================

func TestPostTagGetByPostID_Success(t *testing.T) {
	svc, tagRepo, _ := newTestPostTagService()

	tagRepo.On("GetByPostID", uint(10)).Return([]string{"go", "react"}, nil)

	tags, err := svc.GetByPostID(10)
	assert.NoError(t, err)
	assert.Equal(t, []string{"go", "react"}, tags)
	tagRepo.AssertExpectations(t)
}

func TestPostTagGetByPostID_Empty(t *testing.T) {
	svc, tagRepo, _ := newTestPostTagService()

	tagRepo.On("GetByPostID", uint(10)).Return([]string{}, nil)

	tags, err := svc.GetByPostID(10)
	assert.NoError(t, err)
	assert.Empty(t, tags)
	tagRepo.AssertExpectations(t)
}

func TestPostTagGetByPostID_Error(t *testing.T) {
	svc, tagRepo, _ := newTestPostTagService()

	tagRepo.On("GetByPostID", uint(10)).Return([]string{}, errors.New("db error"))

	tags, err := svc.GetByPostID(10)
	assert.Error(t, err)
	assert.Empty(t, tags)
	tagRepo.AssertExpectations(t)
}

// ============================================================
// FindPostsByTag
// ============================================================

func TestPostTagFindPostsByTag_Success(t *testing.T) {
	svc, tagRepo, _ := newTestPostTagService()

	posts := []model.Post{{Title: "Go Tips"}}
	tagRepo.On("FindPostsByTag", "go", 10, 0).Return(posts, int64(1), nil)

	result, total, err := svc.FindPostsByTag("go", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)
	tagRepo.AssertExpectations(t)
}

func TestPostTagFindPostsByTag_Empty(t *testing.T) {
	svc, tagRepo, _ := newTestPostTagService()

	tagRepo.On("FindPostsByTag", "nonexistent", 10, 0).Return([]model.Post{}, int64(0), nil)

	result, total, err := svc.FindPostsByTag("nonexistent", 10, 0)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
	tagRepo.AssertExpectations(t)
}

// ============================================================
// GetPopularTags
// ============================================================

func TestPostTagGetPopularTags_Success(t *testing.T) {
	svc, tagRepo, _ := newTestPostTagService()

	tags := []model.TagCount{
		{Tag: "go", Count: 15},
		{Tag: "react", Count: 10},
	}
	tagRepo.On("GetPopularTags", 10).Return(tags, nil)

	result, err := svc.GetPopularTags(10)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "go", result[0].Tag)
	assert.Equal(t, 15, result[0].Count)
	tagRepo.AssertExpectations(t)
}

func TestPostTagGetPopularTags_Empty(t *testing.T) {
	svc, tagRepo, _ := newTestPostTagService()

	tagRepo.On("GetPopularTags", 10).Return([]model.TagCount{}, nil)

	result, err := svc.GetPopularTags(10)
	assert.NoError(t, err)
	assert.Empty(t, result)
	tagRepo.AssertExpectations(t)
}

func TestPostTagGetPopularTags_Error(t *testing.T) {
	svc, tagRepo, _ := newTestPostTagService()

	tagRepo.On("GetPopularTags", 10).Return([]model.TagCount{}, errors.New("db error"))

	result, err := svc.GetPopularTags(10)
	assert.Error(t, err)
	assert.Empty(t, result)
	tagRepo.AssertExpectations(t)
}

// ============================================================
// SetAutoTags（コンテンツからハッシュタグ自動抽出・設定）
// ============================================================

func TestPostTagSetAutoTags_ExtractsAndSets(t *testing.T) {
	svc, tagRepo, postRepo := newTestPostTagService()

	post := &model.Post{UserID: 1}
	post.ID = 10

	postRepo.On("FindByID", uint(10)).Return(post, nil)
	tagRepo.On("SetTags", uint(10), []string{"golang", "react"}).Return(nil)

	err := svc.SetAutoTags(10, 1, "今日は #golang と #React を学んだ")
	assert.NoError(t, err)
	tagRepo.AssertExpectations(t)
	postRepo.AssertExpectations(t)
}

func TestPostTagSetAutoTags_NoHashtags(t *testing.T) {
	svc, tagRepo, _ := newTestPostTagService()

	// ハッシュタグがない場合はSetTagsを呼ばない
	err := svc.SetAutoTags(10, 1, "ハッシュタグなしの通常テキスト")
	assert.NoError(t, err)
	tagRepo.AssertNotCalled(t, "SetTags")
}

func TestPostTagSetAutoTags_IgnoresCodeBlocks(t *testing.T) {
	svc, tagRepo, postRepo := newTestPostTagService()

	post := &model.Post{UserID: 1}
	post.ID = 10

	postRepo.On("FindByID", uint(10)).Return(post, nil)
	tagRepo.On("SetTags", uint(10), []string{"validtag"}).Return(nil)

	content := "#validTag を紹介\n```\n#notATag\n```"
	err := svc.SetAutoTags(10, 1, content)
	assert.NoError(t, err)
	tagRepo.AssertExpectations(t)
}

func TestPostTagSetAutoTags_Forbidden(t *testing.T) {
	svc, _, postRepo := newTestPostTagService()

	post := &model.Post{UserID: 1}
	post.ID = 10

	postRepo.On("FindByID", uint(10)).Return(post, nil)

	err := svc.SetAutoTags(10, 999, "#golang を学んだ")
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestPostTagSetAutoTags_EmptyContent(t *testing.T) {
	svc, tagRepo, _ := newTestPostTagService()

	err := svc.SetAutoTags(10, 1, "")
	assert.NoError(t, err)
	tagRepo.AssertNotCalled(t, "SetTags")
}

// ============================================================
// SetTags — タグ文字数上限バリデーション
// ============================================================

func TestPostTagSetTags_TagTooLong(t *testing.T) {
	svc, _, postRepo := newTestPostTagService()

	post := &model.Post{UserID: 1}
	post.ID = 10
	postRepo.On("FindByID", uint(10)).Return(post, nil)

	longTag := strings.Repeat("a", 51)
	err := svc.SetTags(10, 1, []string{longTag})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "タグは50文字以下である必要があります")
}

func TestPostTagSetTags_TagExactly50Chars(t *testing.T) {
	svc, tagRepo, postRepo := newTestPostTagService()

	post := &model.Post{UserID: 1}
	post.ID = 10
	postRepo.On("FindByID", uint(10)).Return(post, nil)

	tag50 := strings.Repeat("a", 50)
	tagRepo.On("SetTags", uint(10), []string{tag50}).Return(nil)

	err := svc.SetTags(10, 1, []string{tag50})
	assert.NoError(t, err)
	tagRepo.AssertExpectations(t)
}

func TestPostTagSetTags_UnicodeTagWithin50Runes(t *testing.T) {
	svc, tagRepo, postRepo := newTestPostTagService()

	post := &model.Post{UserID: 1}
	post.ID = 10
	postRepo.On("FindByID", uint(10)).Return(post, nil)

	// 50文字の日本語タグ（150バイトだがルーン数は50なので許容される）
	unicodeTag := strings.Repeat("あ", 50)
	tagRepo.On("SetTags", uint(10), []string{unicodeTag}).Return(nil)

	err := svc.SetTags(10, 1, []string{unicodeTag})
	assert.NoError(t, err)
	tagRepo.AssertExpectations(t)
}

func TestPostTagSetTags_UnicodeTagExceeds50Runes(t *testing.T) {
	svc, _, postRepo := newTestPostTagService()

	post := &model.Post{UserID: 1}
	post.ID = 10
	postRepo.On("FindByID", uint(10)).Return(post, nil)

	// 51文字の日本語タグ（ルーン数が50を超えるのでエラー）
	unicodeTag := strings.Repeat("あ", 51)
	err := svc.SetTags(10, 1, []string{unicodeTag})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "タグは50文字以下である必要があります")
}

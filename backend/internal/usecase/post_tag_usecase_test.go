package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockPostTagRepo は usecase/repository.PostTagRepository のモック。
type mockPostTagRepo struct{ mock.Mock }

func (m *mockPostTagRepo) SetTags(ctx context.Context, postID uint, tags []string) error {
	return m.Called(ctx, postID, tags).Error(0)
}

func (m *mockPostTagRepo) GetByPostID(ctx context.Context, postID uint) ([]string, error) {
	args := m.Called(ctx, postID)
	t, _ := args.Get(0).([]string)
	return t, args.Error(1)
}

func (m *mockPostTagRepo) FindPostsByTag(ctx context.Context, tag string, limit, offset int) ([]model.Post, int64, error) {
	args := m.Called(ctx, tag, limit, offset)
	p, _ := args.Get(0).([]model.Post)
	return p, args.Get(1).(int64), args.Error(2)
}

func (m *mockPostTagRepo) GetPopularTags(ctx context.Context, limit int) ([]model.TagCount, error) {
	args := m.Called(ctx, limit)
	t, _ := args.Get(0).([]model.TagCount)
	return t, args.Error(1)
}

// ownedPostOf は所有者が userID=1 の投稿を返す。
func ownedPostOf(id uint) *model.Post {
	p := &model.Post{}
	p.ID = id
	p.UserID = 1
	return p
}

func TestSetPostTagsUseCase_Execute(t *testing.T) {
	t.Run("正規化したタグを保存する", func(t *testing.T) {
		tags, posts := new(mockPostTagRepo), new(mockPostReader)
		posts.On("FindByID", mock.Anything, uint(5)).Return(ownedPostOf(5), nil)
		// 小文字化・前後空白除去・重複除外・空文字除外
		tags.On("SetTags", mock.Anything, uint(5), []string{"go", "web"}).Return(nil)
		uc := usecase.NewSetPostTagsUseCase(tags, posts)

		err := uc.Execute(context.Background(), 5, 1, []string{"Go", " web ", "GO", "  "})

		assert.NoError(t, err)
		tags.AssertExpectations(t)
	})

	t.Run("所有者以外は Forbidden（保存しない）", func(t *testing.T) {
		tags, posts := new(mockPostTagRepo), new(mockPostReader)
		other := &model.Post{}
		other.ID = 5
		other.UserID = 999
		posts.On("FindByID", mock.Anything, uint(5)).Return(other, nil)
		uc := usecase.NewSetPostTagsUseCase(tags, posts)

		assertDomainCode(t, uc.Execute(context.Background(), 5, 1, []string{"go"}), domain.ErrCodeForbidden)
		tags.AssertNotCalled(t, "SetTags")
	})

	// 正規化後の件数で判定する（重複を含めて 11 個でも、正規化後 10 個なら通る）。
	t.Run("正規化後に 10 個なら通る", func(t *testing.T) {
		tags, posts := new(mockPostTagRepo), new(mockPostReader)
		posts.On("FindByID", mock.Anything, uint(5)).Return(ownedPostOf(5), nil)
		tags.On("SetTags", mock.Anything, uint(5), mock.MatchedBy(func(ts []string) bool {
			return len(ts) == 10
		})).Return(nil)
		uc := usecase.NewSetPostTagsUseCase(tags, posts)

		in := []string{"t1", "t2", "t3", "t4", "t5", "t6", "t7", "t8", "t9", "t10", "T1"}
		assert.NoError(t, uc.Execute(context.Background(), 5, 1, in))
		tags.AssertExpectations(t)
	})

	t.Run("正規化後に 11 個なら BadRequest（保存しない）", func(t *testing.T) {
		tags, posts := new(mockPostTagRepo), new(mockPostReader)
		posts.On("FindByID", mock.Anything, uint(5)).Return(ownedPostOf(5), nil)
		uc := usecase.NewSetPostTagsUseCase(tags, posts)

		in := []string{"t1", "t2", "t3", "t4", "t5", "t6", "t7", "t8", "t9", "t10", "t11"}
		assertDomainCode(t, uc.Execute(context.Background(), 5, 1, in), domain.ErrCodeBadRequest)
		tags.AssertNotCalled(t, "SetTags")
	})

	t.Run("51 文字のタグは保存しない", func(t *testing.T) {
		tags, posts := new(mockPostTagRepo), new(mockPostReader)
		posts.On("FindByID", mock.Anything, uint(5)).Return(ownedPostOf(5), nil)
		uc := usecase.NewSetPostTagsUseCase(tags, posts)

		err := uc.Execute(context.Background(), 5, 1, []string{strings.Repeat("a", 51)})

		assert.Error(t, err)
		tags.AssertNotCalled(t, "SetTags")
	})

	// 空配列を渡すとタグが全て消える（置き換えの仕様）。
	t.Run("空配列は空で置き換える", func(t *testing.T) {
		tags, posts := new(mockPostTagRepo), new(mockPostReader)
		posts.On("FindByID", mock.Anything, uint(5)).Return(ownedPostOf(5), nil)
		tags.On("SetTags", mock.Anything, uint(5), []string(nil)).Return(nil)
		uc := usecase.NewSetPostTagsUseCase(tags, posts)

		assert.NoError(t, uc.Execute(context.Background(), 5, 1, []string{}))
		tags.AssertExpectations(t)
	})

	t.Run("DB 障害を伝播する", func(t *testing.T) {
		tags, posts := new(mockPostTagRepo), new(mockPostReader)
		posts.On("FindByID", mock.Anything, uint(5)).Return(ownedPostOf(5), nil)
		tags.On("SetTags", mock.Anything, uint(5), mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewSetPostTagsUseCase(tags, posts)

		assert.Error(t, uc.Execute(context.Background(), 5, 1, []string{"go"}))
		tags.AssertExpectations(t)
	})
}

func TestSetAutoPostTagsUseCase_Execute(t *testing.T) {
	t.Run("本文から抽出したタグを設定する", func(t *testing.T) {
		tags, posts := new(mockPostTagRepo), new(mockPostReader)
		posts.On("FindByID", mock.Anything, uint(5)).Return(ownedPostOf(5), nil)
		tags.On("SetTags", mock.Anything, uint(5), []string{"golang", "web"}).Return(nil)
		uc := usecase.NewSetAutoPostTagsUseCase(usecase.NewSetPostTagsUseCase(tags, posts))

		err := uc.Execute(context.Background(), 5, 1, "今日は #golang と #Web を学んだ")

		assert.NoError(t, err)
		tags.AssertExpectations(t)
	})

	// ハッシュタグが無ければ既存のタグには触れない。
	t.Run("ハッシュタグが無ければ何もしない", func(t *testing.T) {
		tags, posts := new(mockPostTagRepo), new(mockPostReader)
		uc := usecase.NewSetAutoPostTagsUseCase(usecase.NewSetPostTagsUseCase(tags, posts))

		assert.NoError(t, uc.Execute(context.Background(), 5, 1, "タグのない本文"))
		posts.AssertNotCalled(t, "FindByID")
		tags.AssertNotCalled(t, "SetTags")
	})

	// コードブロック内の # はタグとして扱わない。
	t.Run("コードブロック内のタグは無視する", func(t *testing.T) {
		tags, posts := new(mockPostTagRepo), new(mockPostReader)
		uc := usecase.NewSetAutoPostTagsUseCase(usecase.NewSetPostTagsUseCase(tags, posts))

		err := uc.Execute(context.Background(), 5, 1, "```\n# comment\n#golang\n```")

		assert.NoError(t, err)
		tags.AssertNotCalled(t, "SetTags")
	})
}

func TestGetPostTagsUseCase_Execute(t *testing.T) {
	tags := new(mockPostTagRepo)
	tags.On("GetByPostID", mock.Anything, uint(5)).Return([]string{"go"}, nil)
	uc := usecase.NewGetPostTagsUseCase(tags)

	got, err := uc.Execute(context.Background(), 5)

	assert.NoError(t, err)
	assert.Equal(t, []string{"go"}, got)
	tags.AssertExpectations(t)
}

func TestFindPostsByTagUseCase_Execute(t *testing.T) {
	t.Run("タグに紐づく投稿と総件数を返す", func(t *testing.T) {
		tags := new(mockPostTagRepo)
		tags.On("FindPostsByTag", mock.Anything, "go", 20, 0).
			Return([]model.Post{{Title: "t"}}, int64(1), nil)
		uc := usecase.NewFindPostsByTagUseCase(tags)

		got, total, err := uc.Execute(context.Background(), "go", 20, 0)

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, int64(1), total)
		tags.AssertExpectations(t)
	})

	// タグは保存時に小文字化・トリムされるため、検索語も同じ規則で正規化しないと一致しない。
	t.Run("大文字・前後空白の検索語を正規化してから検索する", func(t *testing.T) {
		tags := new(mockPostTagRepo)
		tags.On("FindPostsByTag", mock.Anything, "golang", 20, 0).
			Return([]model.Post{{Title: "t"}}, int64(1), nil)
		uc := usecase.NewFindPostsByTagUseCase(tags)

		_, total, err := uc.Execute(context.Background(), "  GoLang  ", 20, 0)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		tags.AssertExpectations(t)
	})
}

func TestGetPopularTagsUseCase_Execute(t *testing.T) {
	tags := new(mockPostTagRepo)
	tags.On("GetPopularTags", mock.Anything, 20).
		Return([]model.TagCount{{Tag: "go", Count: 3}}, nil)
	uc := usecase.NewGetPopularTagsUseCase(tags)

	got, err := uc.Execute(context.Background(), 20)

	assert.NoError(t, err)
	assert.Len(t, got, 1)
	tags.AssertExpectations(t)
}

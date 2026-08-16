package usecase_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockPostRepo は usecase/repository.PostRepository のモック。
type mockPostRepo struct{ mock.Mock }

func (m *mockPostRepo) Create(ctx context.Context, post *model.Post) error {
	return m.Called(ctx, post).Error(0)
}

func (m *mockPostRepo) FindByID(ctx context.Context, id uint) (*model.Post, error) {
	args := m.Called(ctx, id)
	p, _ := args.Get(0).(*model.Post)
	return p, args.Error(1)
}

func (m *mockPostRepo) Update(ctx context.Context, post *model.Post) error {
	return m.Called(ctx, post).Error(0)
}

func (m *mockPostRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

func (m *mockPostRepo) FindAll(ctx context.Context, page, limit int) ([]model.Post, error) {
	args := m.Called(ctx, page, limit)
	posts, _ := args.Get(0).([]model.Post)
	return posts, args.Error(1)
}

func (m *mockPostRepo) CountAll(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockPostRepo) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Post, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	posts, _ := args.Get(0).([]model.Post)
	return posts, args.Get(1).(int64), args.Error(2)
}

func (m *mockPostRepo) FindDraftsByUserID(ctx context.Context, userID uint) ([]model.Post, error) {
	args := m.Called(ctx, userID)
	posts, _ := args.Get(0).([]model.Post)
	return posts, args.Error(1)
}

func (m *mockPostRepo) FindScheduledByUserID(ctx context.Context, userID uint) ([]model.Post, error) {
	args := m.Called(ctx, userID)
	posts, _ := args.Get(0).([]model.Post)
	return posts, args.Error(1)
}

func (m *mockPostRepo) Timeline(ctx context.Context, userID uint, page, limit int) ([]model.Post, error) {
	args := m.Called(ctx, userID, page, limit)
	posts, _ := args.Get(0).([]model.Post)
	return posts, args.Error(1)
}

func (m *mockPostRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockPostRepo) CountDraftsByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockPostRepo) CountScheduledByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// mockPostLikeRepo は usecase/repository.PostLikeRepository のモック。
type mockPostLikeRepo struct{ mock.Mock }

func (m *mockPostLikeRepo) Like(ctx context.Context, userID, postID uint) error {
	return m.Called(ctx, userID, postID).Error(0)
}

func (m *mockPostLikeRepo) Unlike(ctx context.Context, userID, postID uint) error {
	return m.Called(ctx, userID, postID).Error(0)
}

func (m *mockPostLikeRepo) HasLiked(ctx context.Context, userID, postID uint) (bool, error) {
	args := m.Called(ctx, userID, postID)
	return args.Bool(0), args.Error(1)
}

// fakeFollowerNotifier は FollowerNotifier の fake。
// フォロワー通知は goroutine で実行されるため、呼び出し内容は mutex で保護して記録し、
// 完了はテスト側の assert.Eventually で待つ。
type fakeFollowerNotifier struct {
	mu            sync.Mutex
	followerIDs   []uint
	findErr       error
	createErr     error
	notifications []*model.Notification
	called        int
}

func (f *fakeFollowerNotifier) FindFollowerIDs(ctx context.Context, userID uint) ([]uint, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.called++
	return f.followerIDs, f.findErr
}

func (f *fakeFollowerNotifier) CreateBatch(ctx context.Context, notifications []*model.Notification) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.notifications = append(f.notifications, notifications...)
	return f.createErr
}

func (f *fakeFollowerNotifier) snapshot() (int, []*model.Notification) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.called, append([]*model.Notification(nil), f.notifications...)
}

// ============================================================
// 読了時間
// ============================================================

func TestEstimateReadTime(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{"短い本文は最低 1 分", "Hello World", 1},
		{"空文字も 1 分", "", 1},
		{"500 文字ちょうどで 1 分", strings.Repeat("あ", 500), 1},
		{"1500 文字で 3 分", strings.Repeat("あ", 1500), 3},
		{"マルチバイトは文字数で数える", strings.Repeat("あ", 1000) + strings.Repeat("a", 500), 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, usecase.EstimateReadTime(tt.content))
		})
	}
}

// ============================================================
// 作成
// ============================================================

func TestCreatePostUseCase(t *testing.T) {
	t.Run("投稿を作成してフォロワーへ通知する", func(t *testing.T) {
		posts := new(mockPostRepo)
		notifier := &fakeFollowerNotifier{followerIDs: []uint{2, 3}}
		uc := usecase.NewCreatePostUseCase(posts, usecase.NewNotifyFollowersUseCase(notifier))

		posts.On("Create", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
			return p.Title == "タイトル" && p.Content == "本文" && p.EstimatedReadTime == 1
		})).Return(nil).Run(func(args mock.Arguments) {
			args.Get(1).(*model.Post).ID = 10
		})
		posts.On("FindByID", mock.Anything, uint(10)).Return(&model.Post{ID: 10, Title: "タイトル"}, nil)

		got, err := uc.Execute(context.Background(), &model.Post{UserID: 1, Title: " タイトル ", Content: " 本文 "})
		require.NoError(t, err)
		assert.Equal(t, uint(10), got.ID)

		// 通知は goroutine で作られるため完了を待つ
		assert.Eventually(t, func() bool {
			_, notifications := notifier.snapshot()
			return len(notifications) == 2
		}, time.Second, 10*time.Millisecond)

		_, notifications := notifier.snapshot()
		assert.Equal(t, uint(2), notifications[0].UserID)
		assert.Equal(t, uint(1), notifications[0].ActorID)
		assert.Equal(t, model.NotificationTypePost, notifications[0].Type)
		require.NotNil(t, notifications[0].PostID)
		assert.Equal(t, uint(10), *notifications[0].PostID)
		posts.AssertExpectations(t)
	})

	t.Run("下書きはフォロワーへ通知しない", func(t *testing.T) {
		posts := new(mockPostRepo)
		notifier := &fakeFollowerNotifier{followerIDs: []uint{2}}
		uc := usecase.NewCreatePostUseCase(posts, usecase.NewNotifyFollowersUseCase(notifier))

		posts.On("Create", mock.Anything, mock.Anything).Return(nil)
		posts.On("FindByID", mock.Anything, mock.Anything).Return(nil, nil)

		_, err := uc.Execute(context.Background(), &model.Post{UserID: 1, Title: "下書き", Content: "本文", IsDraft: true})
		require.NoError(t, err)

		time.Sleep(50 * time.Millisecond)
		called, _ := notifier.snapshot()
		assert.Zero(t, called)
	})

	t.Run("本文が空ならバリデーションエラー", func(t *testing.T) {
		posts := new(mockPostRepo)
		uc := usecase.NewCreatePostUseCase(posts, usecase.NewNotifyFollowersUseCase(&fakeFollowerNotifier{}))

		_, err := uc.Execute(context.Background(), &model.Post{UserID: 1, Title: "タイトル", Content: "   "})
		assert.Error(t, err)
		posts.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("再取得に失敗しても作成した投稿を返す", func(t *testing.T) {
		posts := new(mockPostRepo)
		uc := usecase.NewCreatePostUseCase(posts, usecase.NewNotifyFollowersUseCase(&fakeFollowerNotifier{}))

		posts.On("Create", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			args.Get(1).(*model.Post).ID = 7
		})
		posts.On("FindByID", mock.Anything, uint(7)).Return(nil, errors.New("db error"))

		got, err := uc.Execute(context.Background(), &model.Post{UserID: 1, Title: "T", Content: "C", IsDraft: true})
		require.NoError(t, err)
		assert.Equal(t, uint(7), got.ID)
	})

	t.Run("保存に失敗したらエラーを返す", func(t *testing.T) {
		posts := new(mockPostRepo)
		uc := usecase.NewCreatePostUseCase(posts, usecase.NewNotifyFollowersUseCase(&fakeFollowerNotifier{}))

		createErr := errors.New("db error")
		posts.On("Create", mock.Anything, mock.Anything).Return(createErr)

		_, err := uc.Execute(context.Background(), &model.Post{UserID: 1, Title: "T", Content: "C"})
		assert.ErrorIs(t, err, createErr)
	})
}

// ============================================================
// 参照
// ============================================================

func TestGetPostUseCase(t *testing.T) {
	t.Run("投稿を返す", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("FindByID", mock.Anything, uint(1)).Return(&model.Post{ID: 1}, nil)

		got, err := usecase.NewGetPostUseCase(posts).Execute(context.Background(), 1)
		require.NoError(t, err)
		assert.Equal(t, uint(1), got.ID)
	})

	t.Run("存在しなければ 404 を返す", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

		_, err := usecase.NewGetPostUseCase(posts).Execute(context.Background(), 1)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})

	t.Run("取得エラーはそのまま返す", func(t *testing.T) {
		posts := new(mockPostRepo)
		dbErr := errors.New("db error")
		posts.On("FindByID", mock.Anything, uint(1)).Return(nil, dbErr)

		_, err := usecase.NewGetPostUseCase(posts).Execute(context.Background(), 1)
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestPostQueryUseCases(t *testing.T) {
	posts := new(mockPostRepo)
	ctx := context.Background()

	posts.On("FindAll", mock.Anything, 1, 20).Return([]model.Post{{ID: 1}}, nil)
	list, err := usecase.NewListPostsUseCase(posts).Execute(ctx, 1, 20)
	require.NoError(t, err)
	assert.Len(t, list, 1)

	posts.On("CountAll", mock.Anything).Return(int64(5), nil)
	total, err := usecase.NewCountPostsUseCase(posts).Execute(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(5), total)

	posts.On("FindByUserID", mock.Anything, uint(2), 20, 0).Return([]model.Post{{ID: 2}}, int64(1), nil)
	byUser, userTotal, err := usecase.NewListUserPostsUseCase(posts).Execute(ctx, 2, 20, 0)
	require.NoError(t, err)
	assert.Len(t, byUser, 1)
	assert.Equal(t, int64(1), userTotal)

	posts.On("FindDraftsByUserID", mock.Anything, uint(1)).Return([]model.Post{{ID: 3, IsDraft: true}}, nil)
	drafts, err := usecase.NewListDraftPostsUseCase(posts).Execute(ctx, 1)
	require.NoError(t, err)
	assert.Len(t, drafts, 1)

	posts.On("FindScheduledByUserID", mock.Anything, uint(1)).Return([]model.Post{{ID: 4}}, nil)
	scheduled, err := usecase.NewListScheduledPostsUseCase(posts).Execute(ctx, 1)
	require.NoError(t, err)
	assert.Len(t, scheduled, 1)

	posts.On("Timeline", mock.Anything, uint(1), 1, 20).Return([]model.Post{{ID: 5}}, nil)
	timeline, err := usecase.NewGetTimelineUseCase(posts).Execute(ctx, 1, 1, 20)
	require.NoError(t, err)
	assert.Len(t, timeline, 1)

	posts.On("CountByUserID", mock.Anything, uint(1)).Return(int64(3), nil)
	countByUser, err := usecase.NewCountUserPostsUseCase(posts).Execute(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(3), countByUser)

	posts.On("CountDraftsByUserID", mock.Anything, uint(1)).Return(int64(2), nil)
	countDrafts, err := usecase.NewCountUserDraftsUseCase(posts).Execute(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(2), countDrafts)

	posts.On("CountScheduledByUserID", mock.Anything, uint(1)).Return(int64(1), nil)
	countScheduled, err := usecase.NewCountUserScheduledPostsUseCase(posts).Execute(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(1), countScheduled)

	posts.AssertExpectations(t)
}

// ============================================================
// 更新・削除
// ============================================================

func TestUpdatePostUseCase(t *testing.T) {
	t.Run("指定した項目だけ更新する", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("FindByID", mock.Anything, uint(1)).
			Return(&model.Post{ID: 1, UserID: 1, Title: "旧", Content: "旧本文", ImageURLs: "old"}, nil)
		posts.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
			return p.Title == "新" && p.Content == "旧本文" && p.ImageURLs == "old"
		})).Return(nil)

		got, err := usecase.NewUpdatePostUseCase(posts).Execute(context.Background(), 1, 1, " 新 ", "  ", "")
		require.NoError(t, err)
		assert.Equal(t, "新", got.Title)
		posts.AssertExpectations(t)
	})

	t.Run("本文を更新すると読了時間も再計算する", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("FindByID", mock.Anything, uint(1)).Return(&model.Post{ID: 1, UserID: 1, EstimatedReadTime: 1}, nil)
		posts.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
			return p.EstimatedReadTime == 2
		})).Return(nil)

		_, err := usecase.NewUpdatePostUseCase(posts).Execute(context.Background(), 1, 1, "", strings.Repeat("あ", 1000), "")
		require.NoError(t, err)
		posts.AssertExpectations(t)
	})

	t.Run("所有者以外は 403", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("FindByID", mock.Anything, uint(1)).Return(&model.Post{ID: 1, UserID: 999}, nil)

		_, err := usecase.NewUpdatePostUseCase(posts).Execute(context.Background(), 1, 1, "新", "", "")
		assert.ErrorIs(t, err, domain.ErrForbidden)
		posts.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("存在しなければ 404 を返す", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)

		_, err := usecase.NewUpdatePostUseCase(posts).Execute(context.Background(), 1, 1, "新", "", "")
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestDeletePostUseCase(t *testing.T) {
	t.Run("所有者は削除できる", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("FindByID", mock.Anything, uint(1)).Return(&model.Post{ID: 1, UserID: 1}, nil)
		posts.On("Delete", mock.Anything, uint(1)).Return(nil)

		require.NoError(t, usecase.NewDeletePostUseCase(posts).Execute(context.Background(), 1, 1))
		posts.AssertExpectations(t)
	})

	t.Run("所有者以外は 403", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("FindByID", mock.Anything, uint(1)).Return(&model.Post{ID: 1, UserID: 999}, nil)

		assert.ErrorIs(t, usecase.NewDeletePostUseCase(posts).Execute(context.Background(), 1, 1), domain.ErrForbidden)
		posts.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})
}

// ============================================================
// 公開・非公開
// ============================================================

func TestPublishPostUseCase(t *testing.T) {
	t.Run("下書きを公開してフォロワーへ通知する", func(t *testing.T) {
		posts := new(mockPostRepo)
		notifier := &fakeFollowerNotifier{followerIDs: []uint{2}}
		uc := usecase.NewPublishPostUseCase(posts, usecase.NewNotifyFollowersUseCase(notifier))

		posts.On("FindByID", mock.Anything, uint(1)).Return(&model.Post{ID: 1, UserID: 1, IsDraft: true}, nil)
		posts.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Post) bool { return !p.IsDraft })).Return(nil)

		got, err := uc.Execute(context.Background(), 1, 1)
		require.NoError(t, err)
		assert.False(t, got.IsDraft)

		assert.Eventually(t, func() bool {
			_, notifications := notifier.snapshot()
			return len(notifications) == 1
		}, time.Second, 10*time.Millisecond)
	})

	t.Run("公開済みは 400", func(t *testing.T) {
		posts := new(mockPostRepo)
		uc := usecase.NewPublishPostUseCase(posts, usecase.NewNotifyFollowersUseCase(&fakeFollowerNotifier{}))
		posts.On("FindByID", mock.Anything, uint(1)).Return(&model.Post{ID: 1, UserID: 1, IsDraft: false}, nil)

		_, err := uc.Execute(context.Background(), 1, 1)
		assert.ErrorIs(t, err, domain.ErrBadRequest)
		posts.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("所有者以外は 403", func(t *testing.T) {
		posts := new(mockPostRepo)
		uc := usecase.NewPublishPostUseCase(posts, usecase.NewNotifyFollowersUseCase(&fakeFollowerNotifier{}))
		posts.On("FindByID", mock.Anything, uint(1)).Return(&model.Post{ID: 1, UserID: 999, IsDraft: true}, nil)

		_, err := uc.Execute(context.Background(), 1, 1)
		assert.ErrorIs(t, err, domain.ErrForbidden)
	})
}

func TestUnpublishPostUseCase(t *testing.T) {
	t.Run("公開済みを下書きに戻す", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("FindByID", mock.Anything, uint(1)).Return(&model.Post{ID: 1, UserID: 1, IsDraft: false}, nil)
		posts.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Post) bool { return p.IsDraft })).Return(nil)

		got, err := usecase.NewUnpublishPostUseCase(posts).Execute(context.Background(), 1, 1)
		require.NoError(t, err)
		assert.True(t, got.IsDraft)
		posts.AssertExpectations(t)
	})

	t.Run("すでに下書きなら 400", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("FindByID", mock.Anything, uint(1)).Return(&model.Post{ID: 1, UserID: 1, IsDraft: true}, nil)

		_, err := usecase.NewUnpublishPostUseCase(posts).Execute(context.Background(), 1, 1)
		assert.ErrorIs(t, err, domain.ErrBadRequest)
	})
}

// ============================================================
// 公開予約
// ============================================================

func TestSchedulePostPublishUseCase(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)

	t.Run("下書きに公開予定日時を設定する", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("FindByID", mock.Anything, uint(1)).Return(&model.Post{ID: 1, UserID: 1, IsDraft: true}, nil)
		posts.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
			return p.ScheduledAt != nil
		})).Return(nil)

		got, err := usecase.NewSchedulePostPublishUseCase(posts).Execute(context.Background(), 1, 1, future)
		require.NoError(t, err)
		require.NotNil(t, got.ScheduledAt)
		posts.AssertExpectations(t)
	})

	t.Run("公開済みは 400", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("FindByID", mock.Anything, uint(1)).Return(&model.Post{ID: 1, UserID: 1, IsDraft: false}, nil)

		_, err := usecase.NewSchedulePostPublishUseCase(posts).Execute(context.Background(), 1, 1, future)
		var domainErr *domain.DomainError
		require.ErrorAs(t, err, &domainErr)
		assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
	})

	t.Run("過去の日時は 400", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("FindByID", mock.Anything, uint(1)).Return(&model.Post{ID: 1, UserID: 1, IsDraft: true}, nil)

		_, err := usecase.NewSchedulePostPublishUseCase(posts).Execute(context.Background(), 1, 1, time.Now().Add(-time.Hour))
		var domainErr *domain.DomainError
		require.ErrorAs(t, err, &domainErr)
		assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
		posts.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})
}

func TestCancelPostScheduleUseCase(t *testing.T) {
	t.Run("公開予約を解除する", func(t *testing.T) {
		posts := new(mockPostRepo)
		scheduled := time.Now().Add(time.Hour)
		posts.On("FindByID", mock.Anything, uint(1)).
			Return(&model.Post{ID: 1, UserID: 1, IsDraft: true, ScheduledAt: &scheduled}, nil)
		posts.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
			return p.ScheduledAt == nil
		})).Return(nil)

		got, err := usecase.NewCancelPostScheduleUseCase(posts).Execute(context.Background(), 1, 1)
		require.NoError(t, err)
		assert.Nil(t, got.ScheduledAt)
		posts.AssertExpectations(t)
	})

	t.Run("予約されていなければ 400", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("FindByID", mock.Anything, uint(1)).Return(&model.Post{ID: 1, UserID: 1, IsDraft: true}, nil)

		_, err := usecase.NewCancelPostScheduleUseCase(posts).Execute(context.Background(), 1, 1)
		var domainErr *domain.DomainError
		require.ErrorAs(t, err, &domainErr)
		assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
	})
}

// ============================================================
// 下書きの自動保存
// ============================================================

func TestAutoSaveDraftUseCase(t *testing.T) {
	t.Run("draftID が 0 なら新規作成する", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("Create", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
			return p.UserID == 1 && p.IsDraft && p.Title == "下書き"
		})).Return(nil)

		got, err := usecase.NewAutoSaveDraftUseCase(posts).Execute(context.Background(), 1, 0, "下書き", "本文", "")
		require.NoError(t, err)
		assert.True(t, got.IsDraft)
		posts.AssertExpectations(t)
	})

	t.Run("空のタイトル・本文でも保存できる", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("Create", mock.Anything, mock.Anything).Return(nil)

		_, err := usecase.NewAutoSaveDraftUseCase(posts).Execute(context.Background(), 1, 0, "", "", "")
		require.NoError(t, err)
	})

	t.Run("タイトルが長すぎるとエラー", func(t *testing.T) {
		posts := new(mockPostRepo)

		_, err := usecase.NewAutoSaveDraftUseCase(posts).Execute(context.Background(), 1, 0, strings.Repeat("a", 201), "", "")
		assert.Error(t, err)
		posts.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("既存の下書きを更新する", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("FindByID", mock.Anything, uint(5)).
			Return(&model.Post{ID: 5, UserID: 1, IsDraft: true, ImageURLs: "old"}, nil)
		posts.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Post) bool {
			// 画像 URL が空なら既存の値を維持する
			return p.Title == "新" && p.Content == "新本文" && p.ImageURLs == "old"
		})).Return(nil)

		_, err := usecase.NewAutoSaveDraftUseCase(posts).Execute(context.Background(), 1, 5, "新", "新本文", "")
		require.NoError(t, err)
		posts.AssertExpectations(t)
	})

	t.Run("公開済みの投稿は 400", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("FindByID", mock.Anything, uint(5)).Return(&model.Post{ID: 5, UserID: 1, IsDraft: false}, nil)

		_, err := usecase.NewAutoSaveDraftUseCase(posts).Execute(context.Background(), 1, 5, "新", "", "")
		var domainErr *domain.DomainError
		require.ErrorAs(t, err, &domainErr)
		assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
	})

	t.Run("所有者以外は 403", func(t *testing.T) {
		posts := new(mockPostRepo)
		posts.On("FindByID", mock.Anything, uint(5)).Return(&model.Post{ID: 5, UserID: 999, IsDraft: true}, nil)

		_, err := usecase.NewAutoSaveDraftUseCase(posts).Execute(context.Background(), 1, 5, "新", "", "")
		assert.ErrorIs(t, err, domain.ErrForbidden)
	})
}

// ============================================================
// いいね
// ============================================================

func TestLikePostUseCase(t *testing.T) {
	t.Run("他人の投稿にはいいねできる", func(t *testing.T) {
		likes := new(mockPostLikeRepo)
		authors := new(mockPostAuthorReader)
		authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(99), nil)
		likes.On("Like", mock.Anything, uint(1), uint(5)).Return(nil)

		require.NoError(t, usecase.NewLikePostUseCase(likes, authors).Execute(context.Background(), 1, 5))
		likes.AssertExpectations(t)
		authors.AssertExpectations(t)
	})

	t.Run("自分の投稿は 403", func(t *testing.T) {
		likes := new(mockPostLikeRepo)
		authors := new(mockPostAuthorReader)
		authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(1), nil)

		err := usecase.NewLikePostUseCase(likes, authors).Execute(context.Background(), 1, 5)
		assert.ErrorIs(t, err, domain.ErrForbidden)
		likes.AssertNotCalled(t, "Like", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("投稿が存在しなければ 404", func(t *testing.T) {
		likes := new(mockPostLikeRepo)
		authors := new(mockPostAuthorReader)
		authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(0), nil)

		err := usecase.NewLikePostUseCase(likes, authors).Execute(context.Background(), 1, 5)
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestUnlikeAndHasLikedUseCases(t *testing.T) {
	likes := new(mockPostLikeRepo)
	authors := new(mockPostAuthorReader)

	authors.On("FindAuthorID", mock.Anything, uint(5)).Return(uint(99), nil)
	likes.On("Unlike", mock.Anything, uint(1), uint(5)).Return(nil)
	require.NoError(t, usecase.NewUnlikePostUseCase(likes, authors).Execute(context.Background(), 1, 5))

	likes.On("HasLiked", mock.Anything, uint(1), uint(5)).Return(true, nil)
	liked, err := usecase.NewHasLikedPostUseCase(likes).Execute(context.Background(), 1, 5)
	require.NoError(t, err)
	assert.True(t, liked)

	likes.AssertExpectations(t)
	authors.AssertExpectations(t)
}

// ============================================================
// フォロワー通知
// ============================================================

func TestNotifyFollowersUseCase(t *testing.T) {
	t.Run("フォロワー全員分の通知を作る", func(t *testing.T) {
		notifier := &fakeFollowerNotifier{followerIDs: []uint{2, 3, 4}}

		require.NoError(t, usecase.NewNotifyFollowersUseCase(notifier).
			Execute(context.Background(), 1, 10, model.NotificationTypePost))

		_, notifications := notifier.snapshot()
		require.Len(t, notifications, 3)
		for i, n := range notifications {
			assert.Equal(t, []uint{2, 3, 4}[i], n.UserID)
			assert.Equal(t, uint(1), n.ActorID)
			require.NotNil(t, n.PostID)
			assert.Equal(t, uint(10), *n.PostID)
		}
	})

	t.Run("フォロワーがいなければ何もしない", func(t *testing.T) {
		notifier := &fakeFollowerNotifier{}

		require.NoError(t, usecase.NewNotifyFollowersUseCase(notifier).
			Execute(context.Background(), 1, 10, model.NotificationTypePost))

		_, notifications := notifier.snapshot()
		assert.Empty(t, notifications)
	})

	t.Run("フォロワー取得に失敗したらエラーを返す", func(t *testing.T) {
		findErr := errors.New("db error")
		notifier := &fakeFollowerNotifier{findErr: findErr}

		err := usecase.NewNotifyFollowersUseCase(notifier).
			Execute(context.Background(), 1, 10, model.NotificationTypePost)
		assert.ErrorIs(t, err, findErr)
	})

	// 非同期実行はリクエストの ctx がキャンセルされても止まらない。
	t.Run("Notify は ctx のキャンセルに影響されない", func(t *testing.T) {
		notifier := &fakeFollowerNotifier{followerIDs: []uint{2}}
		ctx, cancel := context.WithCancel(context.Background())

		usecase.NewNotifyFollowersUseCase(notifier).Notify(ctx, 1, 10, model.NotificationTypePost)
		cancel()

		assert.Eventually(t, func() bool {
			_, notifications := notifier.snapshot()
			return len(notifications) == 1
		}, time.Second, 10*time.Millisecond)
	})
}

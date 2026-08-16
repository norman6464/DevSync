package usecase_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockRoadmapRepo は usecase/repository.RoadmapRepository のモック。
type mockRoadmapRepo struct{ mock.Mock }

func (m *mockRoadmapRepo) Create(ctx context.Context, r *model.Roadmap) error {
	return m.Called(ctx, r).Error(0)
}
func (m *mockRoadmapRepo) Update(ctx context.Context, r *model.Roadmap) error {
	return m.Called(ctx, r).Error(0)
}
func (m *mockRoadmapRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockRoadmapRepo) FindByID(ctx context.Context, id uint) (*model.Roadmap, error) {
	args := m.Called(ctx, id)
	r, _ := args.Get(0).(*model.Roadmap)
	return r, args.Error(1)
}
func (m *mockRoadmapRepo) GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Roadmap, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	r, _ := args.Get(0).([]model.Roadmap)
	return r, args.Get(1).(int64), args.Error(2)
}
func (m *mockRoadmapRepo) GetByStatus(ctx context.Context, userID uint, status string) ([]model.Roadmap, error) {
	args := m.Called(ctx, userID, status)
	r, _ := args.Get(0).([]model.Roadmap)
	return r, args.Error(1)
}
func (m *mockRoadmapRepo) GetPublicRoadmaps(ctx context.Context, limit, offset int) ([]model.Roadmap, int64, error) {
	args := m.Called(ctx, limit, offset)
	r, _ := args.Get(0).([]model.Roadmap)
	return r, args.Get(1).(int64), args.Error(2)
}
func (m *mockRoadmapRepo) GetTemplates(ctx context.Context) ([]model.Roadmap, error) {
	args := m.Called(ctx)
	r, _ := args.Get(0).([]model.Roadmap)
	return r, args.Error(1)
}
func (m *mockRoadmapRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockRoadmapRepo) CopyRoadmap(ctx context.Context, originalID, newUserID uint) (*model.Roadmap, error) {
	args := m.Called(ctx, originalID, newUserID)
	r, _ := args.Get(0).(*model.Roadmap)
	return r, args.Error(1)
}
func (m *mockRoadmapRepo) CreateStep(ctx context.Context, step *model.RoadmapStep) error {
	return m.Called(ctx, step).Error(0)
}
func (m *mockRoadmapRepo) UpdateStep(ctx context.Context, step *model.RoadmapStep) error {
	return m.Called(ctx, step).Error(0)
}
func (m *mockRoadmapRepo) DeleteStep(ctx context.Context, stepID uint) error {
	return m.Called(ctx, stepID).Error(0)
}
func (m *mockRoadmapRepo) FindStepByID(ctx context.Context, stepID uint) (*model.RoadmapStep, error) {
	args := m.Called(ctx, stepID)
	s, _ := args.Get(0).(*model.RoadmapStep)
	return s, args.Error(1)
}
func (m *mockRoadmapRepo) ReorderSteps(ctx context.Context, roadmapID uint, stepOrders []model.StepOrder) error {
	return m.Called(ctx, roadmapID, stepOrders).Error(0)
}

// assertRoadmapStatus は err が期待の HTTP ステータスに対応する DomainError であることを検証する。
func assertRoadmapStatus(t *testing.T, err error, want int) {
	t.Helper()
	require.Error(t, err)
	domainErr := domain.GetDomainError(err)
	require.NotNil(t, domainErr, "DomainError であること")
	assert.Equal(t, want, domainErr.HTTPStatus())
}

func TestCreateRoadmapUseCase_Execute(t *testing.T) {
	t.Run("検証を通れば前後空白を除いて作成する", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Roadmap")).Return(nil)
		uc := usecase.NewCreateRoadmapUseCase(repo)

		roadmap := &model.Roadmap{Title: "  Go 学習  ", Description: "  説明  "}
		require.NoError(t, uc.Execute(context.Background(), roadmap))
		assert.Equal(t, "Go 学習", roadmap.Title)
		assert.Equal(t, "説明", roadmap.Description)
	})

	t.Run("タイトルが空なら 400 で作成しない", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		uc := usecase.NewCreateRoadmapUseCase(repo)

		assertRoadmapStatus(t, uc.Execute(context.Background(), &model.Roadmap{Title: "  "}), http.StatusBadRequest)
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("説明が 1001 文字なら 400", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		uc := usecase.NewCreateRoadmapUseCase(repo)

		err := uc.Execute(context.Background(), &model.Roadmap{
			Title: "題", Description: strings.Repeat("あ", 1001),
		})
		assertRoadmapStatus(t, err, http.StatusBadRequest)
	})
}

func TestGetRoadmapUseCase_Execute(t *testing.T) {
	t.Run("自分のロードマップは非公開でも取得できる", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.Roadmap{ID: 1, UserID: 1, IsPublic: false}, nil)
		uc := usecase.NewGetRoadmapUseCase(repo)

		got, err := uc.Execute(context.Background(), 1, 1)
		require.NoError(t, err)
		assert.Equal(t, uint(1), got.ID)
	})

	t.Run("他人の公開ロードマップは取得できる", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.Roadmap{ID: 1, UserID: 99, IsPublic: true}, nil)
		uc := usecase.NewGetRoadmapUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1)
		assert.NoError(t, err)
	})

	t.Run("他人の非公開ロードマップは 403", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.Roadmap{ID: 1, UserID: 99, IsPublic: false}, nil)
		uc := usecase.NewGetRoadmapUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1)
		assertRoadmapStatus(t, err, http.StatusForbidden)
	})

	t.Run("不在は DomainError ではないエラー（handler で 500）", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(nil, nil)
		uc := usecase.NewGetRoadmapUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1)
		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
	})
}

func TestListRoadmapsByStatusUseCase_Execute(t *testing.T) {
	t.Run("未知のステータスは 400 で repo を引かない", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		uc := usecase.NewListRoadmapsByStatusUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, "unknown")
		assertRoadmapStatus(t, err, http.StatusBadRequest)
		assert.Equal(t, "無効なステータスです", domain.GetDomainError(err).Message)
		repo.AssertNotCalled(t, "GetByStatus", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("有効なステータスは repo に委譲する", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("GetByStatus", mock.Anything, uint(1), "completed").Return([]model.Roadmap{{ID: 1}}, nil)
		uc := usecase.NewListRoadmapsByStatusUseCase(repo)

		got, err := uc.Execute(context.Background(), 1, "completed")
		assert.NoError(t, err)
		assert.Len(t, got, 1)
	})
}

func TestUpdateRoadmapUseCase_Execute(t *testing.T) {
	newRoadmap := func() *model.Roadmap {
		return &model.Roadmap{
			ID: 1, UserID: 1, Title: "旧題", Description: "旧説明",
			Category: model.RoadmapCategorySkill, Status: model.RoadmapStatusActive,
		}
	}

	t.Run("指定した項目だけ更新する", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(newRoadmap(), nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.Roadmap")).Return(nil)
		uc := usecase.NewUpdateRoadmapUseCase(repo)

		got, err := uc.Execute(context.Background(), 1, 1, &model.Roadmap{Title: "  新題  "})
		require.NoError(t, err)
		assert.Equal(t, "新題", got.Title)
		assert.Equal(t, "旧説明", got.Description, "未指定の項目は据え置き")
		assert.Equal(t, model.RoadmapCategorySkill, got.Category)
	})

	t.Run("完了へ変えると完了日時が入る", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(newRoadmap(), nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.Roadmap")).Return(nil)
		uc := usecase.NewUpdateRoadmapUseCase(repo)

		got, err := uc.Execute(context.Background(), 1, 1, &model.Roadmap{Status: model.RoadmapStatusCompleted})
		require.NoError(t, err)
		assert.Equal(t, model.RoadmapStatusCompleted, got.Status)
		assert.NotNil(t, got.CompletedAt)
	})

	t.Run("所有者でなければ 403", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.Roadmap{ID: 1, UserID: 99}, nil)
		uc := usecase.NewUpdateRoadmapUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, &model.Roadmap{Title: "新題"})
		assertRoadmapStatus(t, err, http.StatusForbidden)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("タイトルが 201 文字なら 400", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(newRoadmap(), nil)
		uc := usecase.NewUpdateRoadmapUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, &model.Roadmap{Title: strings.Repeat("あ", 201)})
		assertRoadmapStatus(t, err, http.StatusBadRequest)
	})
}

func TestCopyRoadmapUseCase_Execute(t *testing.T) {
	t.Run("公開ロードマップは他人でも複製できる", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.Roadmap{ID: 5, UserID: 99, IsPublic: true}, nil)
		repo.On("CopyRoadmap", mock.Anything, uint(5), uint(1)).Return(&model.Roadmap{ID: 6}, nil)
		uc := usecase.NewCopyRoadmapUseCase(repo)

		got, err := uc.Execute(context.Background(), 5, 1)
		require.NoError(t, err)
		assert.Equal(t, uint(6), got.ID)
	})

	t.Run("他人の非公開ロードマップは 403", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.Roadmap{ID: 5, UserID: 99, IsPublic: false}, nil)
		uc := usecase.NewCopyRoadmapUseCase(repo)

		_, err := uc.Execute(context.Background(), 5, 1)
		assertRoadmapStatus(t, err, http.StatusForbidden)
		repo.AssertNotCalled(t, "CopyRoadmap", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("自分の非公開ロードマップは複製できる", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.Roadmap{ID: 5, UserID: 1, IsPublic: false}, nil)
		repo.On("CopyRoadmap", mock.Anything, uint(5), uint(1)).Return(&model.Roadmap{ID: 6}, nil)
		uc := usecase.NewCopyRoadmapUseCase(repo)

		_, err := uc.Execute(context.Background(), 5, 1)
		assert.NoError(t, err)
	})

	// 存在確認との間に削除された場合、adapter は (nil, nil) を返す。nil を成功として返さない。
	t.Run("複製結果が nil なら不在として扱う", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(5)).
			Return(&model.Roadmap{ID: 5, UserID: 1, IsPublic: true}, nil)
		repo.On("CopyRoadmap", mock.Anything, uint(5), uint(1)).Return((*model.Roadmap)(nil), nil)
		uc := usecase.NewCopyRoadmapUseCase(repo)

		got, err := uc.Execute(context.Background(), 5, 1)

		assert.Nil(t, got)
		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err), "既存の不在と同じく handler で 500 になる")
	})
}

// テンプレート複製も同じ経路を通るため、nil を成功として返さないことを固定する。
func TestCreateRoadmapFromTemplateUseCase_CopyVanished(t *testing.T) {
	repo := new(mockRoadmapRepo)
	repo.On("FindByID", mock.Anything, uint(5)).
		Return(&model.Roadmap{ID: 5, UserID: 2, IsTemplate: true, IsPublic: true}, nil)
	repo.On("CopyRoadmap", mock.Anything, uint(5), uint(1)).Return((*model.Roadmap)(nil), nil)
	uc := usecase.NewCreateRoadmapFromTemplateUseCase(repo)

	got, err := uc.Execute(context.Background(), 5, 1)

	assert.Nil(t, got)
	require.Error(t, err)
	assert.Nil(t, domain.GetDomainError(err))
}

func TestCreateRoadmapFromTemplateUseCase_Execute(t *testing.T) {
	t.Run("テンプレートなら複製する", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.Roadmap{ID: 5, IsTemplate: true}, nil)
		repo.On("CopyRoadmap", mock.Anything, uint(5), uint(1)).Return(&model.Roadmap{ID: 6}, nil)
		uc := usecase.NewCreateRoadmapFromTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), 5, 1)
		assert.NoError(t, err)
	})

	t.Run("テンプレートでなければ 400", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(5)).Return(&model.Roadmap{ID: 5, IsTemplate: false}, nil)
		uc := usecase.NewCreateRoadmapFromTemplateUseCase(repo)

		_, err := uc.Execute(context.Background(), 5, 1)
		assertRoadmapStatus(t, err, http.StatusBadRequest)
		repo.AssertNotCalled(t, "CopyRoadmap", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestRoadmapStepUseCases(t *testing.T) {
	owned := &model.Roadmap{ID: 1, UserID: 1}

	t.Run("ステップ作成はロードマップ ID を上書きする", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(owned, nil)
		repo.On("CreateStep", mock.Anything, mock.AnythingOfType("*model.RoadmapStep")).Return(nil)
		uc := usecase.NewCreateRoadmapStepUseCase(repo)

		step := &model.RoadmapStep{RoadmapID: 999, Title: "ステップ"}
		require.NoError(t, uc.Execute(context.Background(), 1, 1, step))
		assert.Equal(t, uint(1), step.RoadmapID)
	})

	t.Run("別ロードマップのステップ ID は 400", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(owned, nil)
		repo.On("FindStepByID", mock.Anything, uint(2)).Return(&model.RoadmapStep{ID: 2, RoadmapID: 77}, nil)
		uc := usecase.NewDeleteRoadmapStepUseCase(repo)

		assertRoadmapStatus(t, uc.Execute(context.Background(), 1, 2, 1), http.StatusBadRequest)
		repo.AssertNotCalled(t, "DeleteStep", mock.Anything, mock.Anything)
	})

	t.Run("ステップ更新は指定した項目だけ変える", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(owned, nil)
		repo.On("FindStepByID", mock.Anything, uint(2)).
			Return(&model.RoadmapStep{ID: 2, RoadmapID: 1, Title: "旧", Description: "旧説明"}, nil)
		repo.On("UpdateStep", mock.Anything, mock.AnythingOfType("*model.RoadmapStep")).Return(nil)
		uc := usecase.NewUpdateRoadmapStepUseCase(repo)

		got, err := uc.Execute(context.Background(), 1, 2, 1, &model.RoadmapStep{Title: "  新  "})
		require.NoError(t, err)
		assert.Equal(t, "新", got.Title)
		assert.Equal(t, "旧説明", got.Description)
	})

	t.Run("完了にすると完了日時が入り、戻すと消える", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(owned, nil)
		repo.On("FindStepByID", mock.Anything, uint(2)).Return(&model.RoadmapStep{ID: 2, RoadmapID: 1}, nil)
		repo.On("UpdateStep", mock.Anything, mock.AnythingOfType("*model.RoadmapStep")).Return(nil)
		uc := usecase.NewUpdateRoadmapStepCompletionUseCase(repo)

		done, err := uc.Execute(context.Background(), 1, 2, 1, true)
		require.NoError(t, err)
		assert.True(t, done.IsCompleted)
		assert.NotNil(t, done.CompletedAt)

		undone, err := uc.Execute(context.Background(), 1, 2, 1, false)
		require.NoError(t, err)
		assert.False(t, undone.IsCompleted)
		assert.Nil(t, undone.CompletedAt)
	})

	t.Run("並べ替えは所有者以外 403", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.Roadmap{ID: 1, UserID: 99}, nil)
		uc := usecase.NewReorderRoadmapStepsUseCase(repo)

		err := uc.Execute(context.Background(), 1, 1, []model.StepOrder{{StepID: 2, OrderIndex: 0}})
		assertRoadmapStatus(t, err, http.StatusForbidden)
	})
}

func TestBatchCompleteRoadmapStepsUseCase_Execute(t *testing.T) {
	owned := &model.Roadmap{ID: 1, UserID: 1}

	t.Run("未完了のステップだけ更新する", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(owned, nil)
		repo.On("FindStepByID", mock.Anything, uint(2)).Return(&model.RoadmapStep{ID: 2, RoadmapID: 1}, nil)
		repo.On("FindStepByID", mock.Anything, uint(3)).
			Return(&model.RoadmapStep{ID: 3, RoadmapID: 1, IsCompleted: true}, nil)
		repo.On("UpdateStep", mock.Anything, mock.AnythingOfType("*model.RoadmapStep")).Return(nil).Once()
		uc := usecase.NewBatchCompleteRoadmapStepsUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, []uint{2, 3})
		require.NoError(t, err)
		repo.AssertNumberOfCalls(t, "UpdateStep", 1)
	})

	t.Run("別ロードマップのステップが混ざっていたら 400 で中断する", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(owned, nil)
		repo.On("FindStepByID", mock.Anything, uint(2)).Return(&model.RoadmapStep{ID: 2, RoadmapID: 1}, nil)
		repo.On("FindStepByID", mock.Anything, uint(9)).Return(&model.RoadmapStep{ID: 9, RoadmapID: 77}, nil)
		repo.On("UpdateStep", mock.Anything, mock.AnythingOfType("*model.RoadmapStep")).Return(nil)
		uc := usecase.NewBatchCompleteRoadmapStepsUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, []uint{2, 9})
		assertRoadmapStatus(t, err, http.StatusBadRequest)
		// 中断前に処理した分の更新は残る（移行前からの挙動）。
		repo.AssertNumberOfCalls(t, "UpdateStep", 1)
	})

	t.Run("所有者でなければ 403", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.Roadmap{ID: 1, UserID: 99}, nil)
		uc := usecase.NewBatchCompleteRoadmapStepsUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, 1, []uint{2})
		assertRoadmapStatus(t, err, http.StatusForbidden)
	})
}

func TestSeedRoadmapTemplatesUseCase_Execute(t *testing.T) {
	t.Run("既にテンプレートがあれば何もしない", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("GetTemplates", mock.Anything).Return([]model.Roadmap{{ID: 1}}, nil)
		uc := usecase.NewSeedRoadmapTemplatesUseCase(repo)

		require.NoError(t, uc.Execute(context.Background(), 1))
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("テンプレートが無ければ全件を作る", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("GetTemplates", mock.Anything).Return([]model.Roadmap{}, nil)
		// CreateStep が 1 件ごとに step_count を加算するため、Roadmap は 0 で作られる
		// （実数をセットすると二重加算で 2 倍になる）。
		repo.On("Create", mock.Anything, mock.MatchedBy(func(r *model.Roadmap) bool {
			return r.StepCount == 0
		})).Return(nil)
		repo.On("CreateStep", mock.Anything, mock.AnythingOfType("*model.RoadmapStep")).Return(nil)
		uc := usecase.NewSeedRoadmapTemplatesUseCase(repo)

		require.NoError(t, uc.Execute(context.Background(), 1))
		// テンプレートは 5 種類。ステップは全部で 48 件。
		repo.AssertNumberOfCalls(t, "Create", 5)
		repo.AssertNumberOfCalls(t, "CreateStep", 48)
	})

	t.Run("取得に失敗したらそのまま返す", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("GetTemplates", mock.Anything).Return([]model.Roadmap(nil), errors.New("db error"))
		uc := usecase.NewSeedRoadmapTemplatesUseCase(repo)

		require.Error(t, uc.Execute(context.Background(), 1))
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})
}

func TestRoadmapPassThroughUseCases(t *testing.T) {
	ctx := context.Background()

	t.Run("ユーザー別一覧は repo に委譲する", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("GetByUserID", mock.Anything, uint(1), 20, 0).Return([]model.Roadmap{{ID: 1}}, int64(1), nil)
		got, total, err := usecase.NewListRoadmapsByUserUseCase(repo).Execute(ctx, 1, 20, 0)
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, int64(1), total)
	})

	t.Run("公開一覧は repo に委譲する", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("GetPublicRoadmaps", mock.Anything, 20, 0).Return([]model.Roadmap{}, int64(0), nil)
		_, total, err := usecase.NewListPublicRoadmapsUseCase(repo).Execute(ctx, 20, 0)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
	})

	t.Run("テンプレート一覧は repo に委譲する", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("GetTemplates", mock.Anything).Return([]model.Roadmap{{ID: 1}}, nil)
		got, err := usecase.NewListRoadmapTemplatesUseCase(repo).Execute(ctx)
		assert.NoError(t, err)
		assert.Len(t, got, 1)
	})

	t.Run("件数は repo に委譲する", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(3), nil)
		got, err := usecase.NewCountRoadmapsUseCase(repo).Execute(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), got)
	})

	t.Run("公開切替は所有者のみ", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.Roadmap{ID: 1, UserID: 99}, nil)
		_, err := usecase.NewUpdateRoadmapVisibilityUseCase(repo).Execute(ctx, 1, 1, true)
		assertRoadmapStatus(t, err, http.StatusForbidden)
	})

	t.Run("削除は所有者のみ", func(t *testing.T) {
		repo := new(mockRoadmapRepo)
		repo.On("FindByID", mock.Anything, uint(1)).Return(&model.Roadmap{ID: 1, UserID: 1}, nil)
		repo.On("Delete", mock.Anything, uint(1)).Return(nil)
		assert.NoError(t, usecase.NewDeleteRoadmapUseCase(repo).Execute(ctx, 1, 1))
		repo.AssertExpectations(t)
	})
}

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

// mockStudyCircleRepo は usecase/repository.StudyCircleRepository のモック。
type mockStudyCircleRepo struct{ mock.Mock }

func (m *mockStudyCircleRepo) CreateWithOwner(ctx context.Context, circle *model.StudyCircle) error {
	return m.Called(ctx, circle).Error(0)
}
func (m *mockStudyCircleRepo) FindByID(ctx context.Context, id uint) (*model.StudyCircle, error) {
	args := m.Called(ctx, id)
	c, _ := args.Get(0).(*model.StudyCircle)
	return c, args.Error(1)
}
func (m *mockStudyCircleRepo) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.StudyCircle, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	c, _ := args.Get(0).([]model.StudyCircle)
	return c, args.Get(1).(int64), args.Error(2)
}
func (m *mockStudyCircleRepo) Update(ctx context.Context, circle *model.StudyCircle) error {
	return m.Called(ctx, circle).Error(0)
}
func (m *mockStudyCircleRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockStudyCircleRepo) GetByStatus(ctx context.Context, userID uint, status string) ([]model.StudyCircle, error) {
	args := m.Called(ctx, userID, status)
	c, _ := args.Get(0).([]model.StudyCircle)
	return c, args.Error(1)
}
func (m *mockStudyCircleRepo) Search(ctx context.Context, query string, limit, offset int) ([]model.StudyCircle, int64, error) {
	args := m.Called(ctx, query, limit, offset)
	c, _ := args.Get(0).([]model.StudyCircle)
	return c, args.Get(1).(int64), args.Error(2)
}
func (m *mockStudyCircleRepo) AddMember(ctx context.Context, circleID, userID uint, role model.StudyCircleMemberRole) error {
	return m.Called(ctx, circleID, userID, role).Error(0)
}
func (m *mockStudyCircleRepo) RemoveMember(ctx context.Context, circleID, userID uint) error {
	return m.Called(ctx, circleID, userID).Error(0)
}
func (m *mockStudyCircleRepo) GetMembers(ctx context.Context, circleID uint) ([]model.StudyCircleMember, error) {
	args := m.Called(ctx, circleID)
	c, _ := args.Get(0).([]model.StudyCircleMember)
	return c, args.Error(1)
}
func (m *mockStudyCircleRepo) IsMember(ctx context.Context, circleID, userID uint) (bool, error) {
	args := m.Called(ctx, circleID, userID)
	return args.Bool(0), args.Error(1)
}
func (m *mockStudyCircleRepo) GetMemberCount(ctx context.Context, circleID uint) (int, error) {
	args := m.Called(ctx, circleID)
	return args.Int(0), args.Error(1)
}
func (m *mockStudyCircleRepo) UpdateMemberRole(ctx context.Context, circleID, userID uint, role model.StudyCircleMemberRole) error {
	return m.Called(ctx, circleID, userID, role).Error(0)
}
func (m *mockStudyCircleRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockStudyCircleRepo) CreateStep(ctx context.Context, step *model.StudyCircleStep) error {
	return m.Called(ctx, step).Error(0)
}
func (m *mockStudyCircleRepo) UpdateStep(ctx context.Context, step *model.StudyCircleStep) error {
	return m.Called(ctx, step).Error(0)
}
func (m *mockStudyCircleRepo) DeleteStep(ctx context.Context, stepID uint) error {
	return m.Called(ctx, stepID).Error(0)
}
func (m *mockStudyCircleRepo) FindStepByID(ctx context.Context, stepID uint) (*model.StudyCircleStep, error) {
	args := m.Called(ctx, stepID)
	s, _ := args.Get(0).(*model.StudyCircleStep)
	return s, args.Error(1)
}
func (m *mockStudyCircleRepo) ReorderSteps(ctx context.Context, circleID uint, stepOrders []model.StepOrder) error {
	return m.Called(ctx, circleID, stepOrders).Error(0)
}
func (m *mockStudyCircleRepo) UpsertProgress(ctx context.Context, progress *model.StudyCircleMemberProgress) error {
	return m.Called(ctx, progress).Error(0)
}
func (m *mockStudyCircleRepo) GetProgress(ctx context.Context, circleID uint) ([]model.StudyCircleMemberProgress, error) {
	args := m.Called(ctx, circleID)
	p, _ := args.Get(0).([]model.StudyCircleMemberProgress)
	return p, args.Error(1)
}
func (m *mockStudyCircleRepo) CreateCheckin(ctx context.Context, checkin *model.StudyCircleCheckin) error {
	return m.Called(ctx, checkin).Error(0)
}
func (m *mockStudyCircleRepo) GetCheckins(ctx context.Context, circleID uint) ([]model.StudyCircleCheckin, error) {
	args := m.Called(ctx, circleID)
	c, _ := args.Get(0).([]model.StudyCircleCheckin)
	return c, args.Error(1)
}
func (m *mockStudyCircleRepo) HasCheckedInToday(ctx context.Context, circleID, userID uint) (bool, error) {
	args := m.Called(ctx, circleID, userID)
	return args.Bool(0), args.Error(1)
}
func (m *mockStudyCircleRepo) GetStreakRanking(ctx context.Context, circleID uint) ([]model.CircleMemberStreak, error) {
	args := m.Called(ctx, circleID)
	r, _ := args.Get(0).([]model.CircleMemberStreak)
	return r, args.Error(1)
}

// assertStudyCircleStatus は err が期待の HTTP ステータスに対応する DomainError であることを検証する。
func assertStudyCircleStatus(t *testing.T, err error, want int) {
	t.Helper()
	require.Error(t, err)
	domainErr := domain.GetDomainError(err)
	require.NotNil(t, domainErr, "DomainError であること")
	assert.Equal(t, want, domainErr.HTTPStatus())
}

func TestCreateStudyCircleUseCase_Execute(t *testing.T) {
	t.Run("オーナー込みで作成し、招待メンバーも追加する", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("CreateWithOwner", mock.Anything, mock.AnythingOfType("*model.StudyCircle")).Return(nil)
		repo.On("AddMember", mock.Anything, uint(0), uint(2), model.StudyCircleRoleMember).Return(nil)
		uc := usecase.NewCreateStudyCircleUseCase(repo)

		circle := &model.StudyCircle{Name: "  Go 勉強会  ", Topic: " 入門 ", OwnerID: 1, MaxMembers: 99}
		err := uc.Execute(context.Background(), circle, []uint{1, 2})

		assert.NoError(t, err)
		assert.Equal(t, "Go 勉強会", circle.Name, "前後の空白は除去される")
		assert.Equal(t, "入門", circle.Topic)
		assert.Equal(t, 5, circle.MaxMembers, "範囲外の上限は 5 に補正される")
		assert.Equal(t, model.StudyCircleStatusActive, circle.Status)
		repo.AssertExpectations(t)
		// オーナー自身は member_ids に含まれていても二重追加しない（オーナー登録は CreateWithOwner 側）。
		repo.AssertNumberOfCalls(t, "AddMember", 1)
	})

	t.Run("名前が空なら 400 で作成しない", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		uc := usecase.NewCreateStudyCircleUseCase(repo)

		err := uc.Execute(context.Background(), &model.StudyCircle{Name: "   ", Topic: "x"}, nil)

		assertStudyCircleStatus(t, err, http.StatusBadRequest)
		repo.AssertNotCalled(t, "CreateWithOwner", mock.Anything, mock.Anything)
	})

	t.Run("説明が 1000 文字を超えたら 400", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		uc := usecase.NewCreateStudyCircleUseCase(repo)

		err := uc.Execute(context.Background(), &model.StudyCircle{
			Name: "x", Topic: "y", Description: strings.Repeat("あ", 1001),
		}, nil)

		assertStudyCircleStatus(t, err, http.StatusBadRequest)
	})

	t.Run("招待メンバーの追加に失敗しても全体は成功する", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("CreateWithOwner", mock.Anything, mock.Anything).Return(nil)
		repo.On("AddMember", mock.Anything, uint(0), uint(9), model.StudyCircleRoleMember).
			Return(errors.New("db error"))
		uc := usecase.NewCreateStudyCircleUseCase(repo)

		err := uc.Execute(context.Background(), &model.StudyCircle{Name: "x", Topic: "y", OwnerID: 1}, []uint{9})

		assert.NoError(t, err)
	})

	t.Run("オーナー込みの作成に失敗したらエラーを返し、招待メンバーを追加しない", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("CreateWithOwner", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewCreateStudyCircleUseCase(repo)

		err := uc.Execute(context.Background(), &model.StudyCircle{Name: "x", Topic: "y", OwnerID: 1}, []uint{9})

		assert.Error(t, err)
		repo.AssertNotCalled(t, "AddMember", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestListStudyCirclesByStatusUseCase_Execute(t *testing.T) {
	t.Run("未知のステータスは 400", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		uc := usecase.NewListStudyCirclesByStatusUseCase(repo)

		_, err := uc.Execute(context.Background(), 1, "unknown")

		assertStudyCircleStatus(t, err, http.StatusBadRequest)
		repo.AssertNotCalled(t, "GetByStatus", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("有効なステータスは repo に委譲する", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("GetByStatus", mock.Anything, uint(1), "active").
			Return([]model.StudyCircle{{Name: "A"}}, nil)
		uc := usecase.NewListStudyCirclesByStatusUseCase(repo)

		got, err := uc.Execute(context.Background(), 1, "active")

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		repo.AssertExpectations(t)
	})
}

func TestGetStudyCircleUseCase_Execute(t *testing.T) {
	t.Run("不在なら 404", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)
		uc := usecase.NewGetStudyCircleUseCase(repo)

		_, err := uc.Execute(context.Background(), 10, 1)

		assertStudyCircleStatus(t, err, http.StatusNotFound)
	})

	t.Run("非メンバーは 403", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(&model.StudyCircle{OwnerID: 2}, nil)
		repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(false, nil)
		uc := usecase.NewGetStudyCircleUseCase(repo)

		_, err := uc.Execute(context.Background(), 10, 1)

		assertStudyCircleStatus(t, err, http.StatusForbidden)
	})

	t.Run("メンバー判定の DB 障害はそのまま伝播する", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(&model.StudyCircle{OwnerID: 2}, nil)
		repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(false, errors.New("db error"))
		uc := usecase.NewGetStudyCircleUseCase(repo)

		_, err := uc.Execute(context.Background(), 10, 1)

		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err), "素の DB エラーのまま返る")
	})
}

func TestUpdateStudyCircleUseCase_Execute(t *testing.T) {
	newCircle := func() *model.StudyCircle {
		return &model.StudyCircle{ID: 10, Name: "旧名", Topic: "旧トピック", Description: "旧説明", OwnerID: 1}
	}

	t.Run("指定した項目だけ更新する", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(newCircle(), nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*model.StudyCircle")).Return(nil)
		uc := usecase.NewUpdateStudyCircleUseCase(repo)

		name := "  新名  "
		got, err := uc.Execute(context.Background(), 10, 1, &name, nil, nil)

		require.NoError(t, err)
		assert.Equal(t, "新名", got.Name)
		assert.Equal(t, "旧トピック", got.Topic, "未指定の項目は据え置き")
		assert.Equal(t, "旧説明", got.Description)
	})

	t.Run("オーナー以外は 403", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(&model.StudyCircle{ID: 10, OwnerID: 99}, nil)
		uc := usecase.NewUpdateStudyCircleUseCase(repo)

		name := "新名"
		_, err := uc.Execute(context.Background(), 10, 1, &name, nil, nil)

		assertStudyCircleStatus(t, err, http.StatusForbidden)
	})

	t.Run("取得の DB 障害も 404 に潰れる", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(nil, errors.New("db error"))
		uc := usecase.NewUpdateStudyCircleUseCase(repo)

		_, err := uc.Execute(context.Background(), 10, 1, nil, nil, nil)

		assertStudyCircleStatus(t, err, http.StatusNotFound)
	})

	t.Run("空トピックは 400", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(newCircle(), nil)
		uc := usecase.NewUpdateStudyCircleUseCase(repo)

		topic := "   "
		_, err := uc.Execute(context.Background(), 10, 1, nil, &topic, nil)

		assertStudyCircleStatus(t, err, http.StatusBadRequest)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})
}

func TestAddStudyCircleMemberUseCase_Execute(t *testing.T) {
	t.Run("非メンバーからの追加は 403", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(false, nil)
		uc := usecase.NewAddStudyCircleMemberUseCase(repo)

		err := uc.Execute(context.Background(), 10, 1, 5)

		assertStudyCircleStatus(t, err, http.StatusForbidden)
	})

	t.Run("既にメンバーなら 400", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
		repo.On("IsMember", mock.Anything, uint(10), uint(5)).Return(true, nil)
		uc := usecase.NewAddStudyCircleMemberUseCase(repo)

		err := uc.Execute(context.Background(), 10, 1, 5)

		assertStudyCircleStatus(t, err, http.StatusBadRequest)
	})

	t.Run("上限到達なら専用メッセージで 400", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
		repo.On("IsMember", mock.Anything, uint(10), uint(5)).Return(false, nil)
		repo.On("FindByID", mock.Anything, uint(10)).Return(&model.StudyCircle{MaxMembers: 3}, nil)
		repo.On("GetMemberCount", mock.Anything, uint(10)).Return(3, nil)
		uc := usecase.NewAddStudyCircleMemberUseCase(repo)

		err := uc.Execute(context.Background(), 10, 1, 5)

		assertStudyCircleStatus(t, err, http.StatusBadRequest)
		assert.Equal(t, "メンバー上限に達しました", domain.GetDomainError(err).Message)
		repo.AssertNotCalled(t, "AddMember", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("サークルが不在なら 404", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
		repo.On("IsMember", mock.Anything, uint(10), uint(5)).Return(false, nil)
		repo.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)
		uc := usecase.NewAddStudyCircleMemberUseCase(repo)

		err := uc.Execute(context.Background(), 10, 1, 5)

		assertStudyCircleStatus(t, err, http.StatusNotFound)
	})
}

func TestUpdateStudyCircleMemberRoleUseCase_Execute(t *testing.T) {
	t.Run("未知の役割は 400 で所有権も見ない", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		uc := usecase.NewUpdateStudyCircleMemberRoleUseCase(repo)

		err := uc.Execute(context.Background(), 10, 1, 5, "admin")

		assertStudyCircleStatus(t, err, http.StatusBadRequest)
		repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	})

	t.Run("オーナー自身の役割は変更できない", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(&model.StudyCircle{OwnerID: 1}, nil)
		uc := usecase.NewUpdateStudyCircleMemberRoleUseCase(repo)

		err := uc.Execute(context.Background(), 10, 1, 1, "member")

		assertStudyCircleStatus(t, err, http.StatusBadRequest)
		assert.Equal(t, "オーナー自身の役割は変更できません", domain.GetDomainError(err).Message)
	})

	t.Run("対象が非メンバーなら 404", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(&model.StudyCircle{OwnerID: 1}, nil)
		repo.On("IsMember", mock.Anything, uint(10), uint(5)).Return(false, nil)
		uc := usecase.NewUpdateStudyCircleMemberRoleUseCase(repo)

		err := uc.Execute(context.Background(), 10, 1, 5, "owner")

		assertStudyCircleStatus(t, err, http.StatusNotFound)
		assert.Equal(t, "指定されたユーザーはメンバーではありません", domain.GetDomainError(err).Message)
	})

	t.Run("オーナーは他メンバーの役割を変更できる", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(&model.StudyCircle{OwnerID: 1}, nil)
		repo.On("IsMember", mock.Anything, uint(10), uint(5)).Return(true, nil)
		repo.On("UpdateMemberRole", mock.Anything, uint(10), uint(5), model.StudyCircleRoleOwner).Return(nil)
		uc := usecase.NewUpdateStudyCircleMemberRoleUseCase(repo)

		err := uc.Execute(context.Background(), 10, 1, 5, "owner")

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestRemoveStudyCircleMemberUseCase_Execute(t *testing.T) {
	t.Run("自分自身ならオーナーでなくても退出できる", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(&model.StudyCircle{OwnerID: 99}, nil)
		repo.On("RemoveMember", mock.Anything, uint(10), uint(1)).Return(nil)
		uc := usecase.NewRemoveStudyCircleMemberUseCase(repo)

		err := uc.Execute(context.Background(), 10, 1, 1)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("他人の除外はオーナーのみ", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(&model.StudyCircle{OwnerID: 99}, nil)
		uc := usecase.NewRemoveStudyCircleMemberUseCase(repo)

		err := uc.Execute(context.Background(), 10, 1, 5)

		assertStudyCircleStatus(t, err, http.StatusForbidden)
		repo.AssertNotCalled(t, "RemoveMember", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("サークルが不在なら 404", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)
		uc := usecase.NewRemoveStudyCircleMemberUseCase(repo)

		err := uc.Execute(context.Background(), 10, 1, 1)

		assertStudyCircleStatus(t, err, http.StatusNotFound)
	})
}

func TestStudyCircleStepUseCases(t *testing.T) {
	ownedCircle := &model.StudyCircle{ID: 10, OwnerID: 1}

	t.Run("ステップ作成はサークル ID を上書きする", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(ownedCircle, nil)
		repo.On("CreateStep", mock.Anything, mock.AnythingOfType("*model.StudyCircleStep")).Return(nil)
		uc := usecase.NewCreateStudyCircleStepUseCase(repo)

		step := &model.StudyCircleStep{CircleID: 999, Title: "1章"}
		err := uc.Execute(context.Background(), 10, 1, step)

		assert.NoError(t, err)
		assert.Equal(t, uint(10), step.CircleID)
	})

	t.Run("参考URLが2000文字を超えるとリポジトリを呼ばず検証エラー", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		uc := usecase.NewCreateStudyCircleStepUseCase(repo)

		step := &model.StudyCircleStep{Title: "1章", ResourceURL: "https://example.com/" + strings.Repeat("a", 2000)}
		err := uc.Execute(context.Background(), 10, 1, step)

		assert.Error(t, err)
		assert.True(t, domain.IsDomainError(err), "DomainError（400系）として返す: %v", err)
		repo.AssertNotCalled(t, "CreateStep", mock.Anything, mock.Anything)
	})

	t.Run("2000文字ちょうどの参考URLでステップを作成できる", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(ownedCircle, nil)
		repo.On("CreateStep", mock.Anything, mock.AnythingOfType("*model.StudyCircleStep")).Return(nil)
		uc := usecase.NewCreateStudyCircleStepUseCase(repo)

		url2000 := "https://example.com/" + strings.Repeat("a", 2000-len("https://example.com/"))
		step := &model.StudyCircleStep{Title: "1章", ResourceURL: url2000}
		err := uc.Execute(context.Background(), 10, 1, step)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("別サークルのステップ ID は 404", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(ownedCircle, nil)
		repo.On("FindStepByID", mock.Anything, uint(5)).
			Return(&model.StudyCircleStep{ID: 5, CircleID: 77}, nil)
		uc := usecase.NewDeleteStudyCircleStepUseCase(repo)

		err := uc.Execute(context.Background(), 10, 1, 5)

		assertStudyCircleStatus(t, err, http.StatusNotFound)
		repo.AssertNotCalled(t, "DeleteStep", mock.Anything, mock.Anything)
	})

	t.Run("ステップ更新は指定した項目だけ変える", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(ownedCircle, nil)
		repo.On("FindStepByID", mock.Anything, uint(5)).
			Return(&model.StudyCircleStep{ID: 5, CircleID: 10, Title: "旧", Description: "旧説明"}, nil)
		repo.On("UpdateStep", mock.Anything, mock.AnythingOfType("*model.StudyCircleStep")).Return(nil)
		uc := usecase.NewUpdateStudyCircleStepUseCase(repo)

		title := " 新 "
		got, err := uc.Execute(context.Background(), 10, 1, 5, &title, nil)

		require.NoError(t, err)
		assert.Equal(t, "新", got.Title)
		assert.Equal(t, "旧説明", got.Description)
	})

	t.Run("ステップ更新のタイトルが空なら 400", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(ownedCircle, nil)
		repo.On("FindStepByID", mock.Anything, uint(5)).
			Return(&model.StudyCircleStep{ID: 5, CircleID: 10}, nil)
		uc := usecase.NewUpdateStudyCircleStepUseCase(repo)

		title := ""
		_, err := uc.Execute(context.Background(), 10, 1, 5, &title, nil)

		assertStudyCircleStatus(t, err, http.StatusBadRequest)
	})

	t.Run("並べ替えはオーナー以外 403", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(&model.StudyCircle{ID: 10, OwnerID: 99}, nil)
		uc := usecase.NewReorderStudyCircleStepsUseCase(repo)

		err := uc.Execute(context.Background(), 10, 1, []model.StepOrder{{StepID: 1, OrderIndex: 0}})

		assertStudyCircleStatus(t, err, http.StatusForbidden)
	})
}

func TestUpdateStudyCircleProgressUseCase_Execute(t *testing.T) {
	t.Run("完了時のみ完了日時を記録する", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
		var saved *model.StudyCircleMemberProgress
		repo.On("UpsertProgress", mock.Anything, mock.AnythingOfType("*model.StudyCircleMemberProgress")).
			Run(func(args mock.Arguments) {
				saved = args.Get(1).(*model.StudyCircleMemberProgress)
			}).Return(nil)
		uc := usecase.NewUpdateStudyCircleProgressUseCase(repo)

		require.NoError(t, uc.Execute(context.Background(), 10, 1, 3, true))
		require.NotNil(t, saved)
		assert.True(t, saved.IsCompleted)
		assert.NotNil(t, saved.CompletedAt)

		require.NoError(t, uc.Execute(context.Background(), 10, 1, 3, false))
		assert.False(t, saved.IsCompleted)
		assert.Nil(t, saved.CompletedAt, "未完了に戻すときは完了日時を送らない")
	})

	t.Run("非メンバーは 403", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(false, nil)
		uc := usecase.NewUpdateStudyCircleProgressUseCase(repo)

		err := uc.Execute(context.Background(), 10, 1, 3, true)

		assertStudyCircleStatus(t, err, http.StatusForbidden)
	})
}

func TestCreateStudyCircleCheckinUseCase_Execute(t *testing.T) {
	t.Run("同じ日の 2 回目は 409", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
		repo.On("HasCheckedInToday", mock.Anything, uint(10), uint(1)).Return(true, nil)
		uc := usecase.NewCreateStudyCircleCheckinUseCase(repo)

		_, err := uc.Execute(context.Background(), 10, 1, "今日やったこと")

		assertStudyCircleStatus(t, err, http.StatusConflict)
	})

	t.Run("空の内容は 400 でメンバー判定もしない", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		uc := usecase.NewCreateStudyCircleCheckinUseCase(repo)

		_, err := uc.Execute(context.Background(), 10, 1, "   ")

		assertStudyCircleStatus(t, err, http.StatusBadRequest)
		repo.AssertNotCalled(t, "IsMember", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("作成に成功したら当日日付で保存する", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
		repo.On("HasCheckedInToday", mock.Anything, uint(10), uint(1)).Return(false, nil)
		repo.On("CreateCheckin", mock.Anything, mock.AnythingOfType("*model.StudyCircleCheckin")).Return(nil)
		uc := usecase.NewCreateStudyCircleCheckinUseCase(repo)

		got, err := uc.Execute(context.Background(), 10, 1, "  進捗  ")

		require.NoError(t, err)
		assert.Equal(t, "進捗", got.Content)
		assert.NotEmpty(t, got.Date)
	})
}

func TestStudyCircleMemberOnlyReads(t *testing.T) {
	cases := []struct {
		name string
		call func(repo *mockStudyCircleRepo) error
	}{
		{"メンバー一覧", func(repo *mockStudyCircleRepo) error {
			_, err := usecase.NewListStudyCircleMembersUseCase(repo).Execute(context.Background(), 10, 1)
			return err
		}},
		{"進捗一覧", func(repo *mockStudyCircleRepo) error {
			_, err := usecase.NewListStudyCircleProgressUseCase(repo).Execute(context.Background(), 10, 1)
			return err
		}},
		{"チェックイン一覧", func(repo *mockStudyCircleRepo) error {
			_, err := usecase.NewListStudyCircleCheckinsUseCase(repo).Execute(context.Background(), 10, 1)
			return err
		}},
		{"ストリークランキング", func(repo *mockStudyCircleRepo) error {
			_, err := usecase.NewGetStudyCircleStreakRankingUseCase(repo).Execute(context.Background(), 10, 1)
			return err
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name+"は非メンバーなら 403", func(t *testing.T) {
			repo := new(mockStudyCircleRepo)
			repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(false, nil)

			assertStudyCircleStatus(t, tc.call(repo), http.StatusForbidden)
		})
	}
}

func TestStudyCirclePassThroughUseCases(t *testing.T) {
	t.Run("検索は repo に委譲する", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("Search", mock.Anything, "go", 20, 0).
			Return([]model.StudyCircle{{Name: "Go"}}, int64(1), nil)
		uc := usecase.NewSearchStudyCirclesUseCase(repo)

		got, total, err := uc.Execute(context.Background(), "go", 20, 0)

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, int64(1), total)
	})

	t.Run("参加数は repo に委譲する", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(3), nil)
		uc := usecase.NewCountStudyCirclesUseCase(repo)

		got, err := uc.Execute(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, int64(3), got)
	})

	t.Run("参加一覧は repo に委譲する", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByUserID", mock.Anything, uint(1), 20, 0).
			Return([]model.StudyCircle{{Name: "A"}}, int64(1), nil)
		uc := usecase.NewListMyStudyCirclesUseCase(repo)

		got, total, err := uc.Execute(context.Background(), 1, 20, 0)

		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, int64(1), total)
	})

	t.Run("削除はオーナーのみ", func(t *testing.T) {
		repo := new(mockStudyCircleRepo)
		repo.On("FindByID", mock.Anything, uint(10)).Return(&model.StudyCircle{ID: 10, OwnerID: 1}, nil)
		repo.On("Delete", mock.Anything, uint(10)).Return(nil)
		uc := usecase.NewDeleteStudyCircleUseCase(repo)

		assert.NoError(t, uc.Execute(context.Background(), 10, 1))
		repo.AssertExpectations(t)
	})
}

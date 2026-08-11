package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockProjectMilestoneRepo は usecase/repository.ProjectMilestoneRepository のモック。
type mockProjectMilestoneRepo struct{ mock.Mock }

func (m *mockProjectMilestoneRepo) Create(ctx context.Context, milestone *model.ProjectMilestone) error {
	return m.Called(ctx, milestone).Error(0)
}

func (m *mockProjectMilestoneRepo) FindByID(ctx context.Context, id uint) (*model.ProjectMilestone, error) {
	args := m.Called(ctx, id)
	ms, _ := args.Get(0).(*model.ProjectMilestone)
	return ms, args.Error(1)
}

func (m *mockProjectMilestoneRepo) FindByProjectID(ctx context.Context, projectID uint) ([]model.ProjectMilestone, error) {
	args := m.Called(ctx, projectID)
	ms, _ := args.Get(0).([]model.ProjectMilestone)
	return ms, args.Error(1)
}

func (m *mockProjectMilestoneRepo) Update(ctx context.Context, milestone *model.ProjectMilestone) error {
	return m.Called(ctx, milestone).Error(0)
}

func (m *mockProjectMilestoneRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}

// mockProjectReader は usecase/repository.ProjectReader のモック。
type mockProjectReader struct{ mock.Mock }

func (m *mockProjectReader) FindByID(ctx context.Context, id uint) (*model.Project, error) {
	args := m.Called(ctx, id)
	p, _ := args.Get(0).(*model.Project)
	return p, args.Error(1)
}

// ownedProjectOf は所有者が userID=1 のプロジェクトを返す。
func ownedProjectOf(id uint) *model.Project {
	p := &model.Project{}
	p.ID = id
	p.UserID = 1
	return p
}

// otherProjectOf は他ユーザー所有のプロジェクトを返す。
func otherProjectOf(id uint) *model.Project {
	p := &model.Project{}
	p.ID = id
	p.UserID = 999
	return p
}

// milestoneOf はプロジェクト 10 に属するマイルストーンを返す。
func milestoneOf() *model.ProjectMilestone {
	ms := &model.ProjectMilestone{ProjectID: 10, Title: "旧題", Description: "旧説明"}
	ms.ID = 7
	return ms
}

func TestCreateProjectMilestoneUseCase_Execute(t *testing.T) {
	t.Run("前後空白を除いて未着手で作成する", func(t *testing.T) {
		milestones, projects := new(mockProjectMilestoneRepo), new(mockProjectReader)
		projects.On("FindByID", mock.Anything, uint(10)).Return(ownedProjectOf(10), nil)
		milestones.On("Create", mock.Anything, mock.MatchedBy(func(ms *model.ProjectMilestone) bool {
			return ms.Title == "リリース" && ms.Description == "説明" &&
				ms.Status == model.MilestoneNotStarted && ms.ProjectID == 10
		})).Return(nil)
		uc := usecase.NewCreateProjectMilestoneUseCase(milestones, projects)

		err := uc.Execute(context.Background(), usecase.CreateProjectMilestoneInput{
			UserID: 1, ProjectID: 10, Title: "  リリース  ", Description: "  説明  ",
		})

		assert.NoError(t, err)
		milestones.AssertExpectations(t)
	})

	t.Run("プロジェクト不在は NotFound（作成しない）", func(t *testing.T) {
		milestones, projects := new(mockProjectMilestoneRepo), new(mockProjectReader)
		projects.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)
		uc := usecase.NewCreateProjectMilestoneUseCase(milestones, projects)

		err := uc.Execute(context.Background(), usecase.CreateProjectMilestoneInput{
			UserID: 1, ProjectID: 10, Title: "t",
		})

		assertDomainCode(t, err, domain.ErrCodeNotFound)
		milestones.AssertNotCalled(t, "Create")
	})

	// 403 はこのスライス専用のメッセージを返す（共通 helper の汎用文言ではない）。
	t.Run("所有者以外は Forbidden で専用メッセージ", func(t *testing.T) {
		milestones, projects := new(mockProjectMilestoneRepo), new(mockProjectReader)
		projects.On("FindByID", mock.Anything, uint(10)).Return(otherProjectOf(10), nil)
		uc := usecase.NewCreateProjectMilestoneUseCase(milestones, projects)

		err := uc.Execute(context.Background(), usecase.CreateProjectMilestoneInput{
			UserID: 1, ProjectID: 10, Title: "t",
		})

		assertDomainCode(t, err, domain.ErrCodeForbidden)
		var de *domain.DomainError
		if assert.ErrorAs(t, err, &de) {
			assert.Equal(t, "このプロジェクトを編集する権限がありません", de.Message)
		}
		milestones.AssertNotCalled(t, "Create")
	})

	t.Run("タイトルが不正なら作成しない", func(t *testing.T) {
		for name, title := range map[string]string{
			"空文字":    "",
			"空白のみ":   "   ",
			"201 文字": strings.Repeat("a", 201),
		} {
			t.Run(name, func(t *testing.T) {
				milestones, projects := new(mockProjectMilestoneRepo), new(mockProjectReader)
				projects.On("FindByID", mock.Anything, uint(10)).Return(ownedProjectOf(10), nil)
				uc := usecase.NewCreateProjectMilestoneUseCase(milestones, projects)

				err := uc.Execute(context.Background(), usecase.CreateProjectMilestoneInput{
					UserID: 1, ProjectID: 10, Title: title,
				})

				assert.Error(t, err)
				milestones.AssertNotCalled(t, "Create")
			})
		}
	})

	// 作成時の説明は 0 文字を許す（更新時は 1 文字以上を要求する点と異なる）。
	t.Run("説明は空でもよい", func(t *testing.T) {
		milestones, projects := new(mockProjectMilestoneRepo), new(mockProjectReader)
		projects.On("FindByID", mock.Anything, uint(10)).Return(ownedProjectOf(10), nil)
		milestones.On("Create", mock.Anything, mock.Anything).Return(nil)
		uc := usecase.NewCreateProjectMilestoneUseCase(milestones, projects)

		err := uc.Execute(context.Background(), usecase.CreateProjectMilestoneInput{
			UserID: 1, ProjectID: 10, Title: "t", Description: "",
		})

		assert.NoError(t, err)
		milestones.AssertExpectations(t)
	})
}

func TestListProjectMilestonesUseCase_Execute(t *testing.T) {
	milestones := new(mockProjectMilestoneRepo)
	milestones.On("FindByProjectID", mock.Anything, uint(10)).
		Return([]model.ProjectMilestone{{ProjectID: 10}}, nil)
	uc := usecase.NewListProjectMilestonesUseCase(milestones)

	got, err := uc.Execute(context.Background(), 10)

	assert.NoError(t, err)
	assert.Len(t, got, 1)
	milestones.AssertExpectations(t)
}

func TestUpdateProjectMilestoneUseCase_Execute(t *testing.T) {
	t.Run("空の項目は据え置く部分更新", func(t *testing.T) {
		milestones, projects := new(mockProjectMilestoneRepo), new(mockProjectReader)
		milestones.On("FindByID", mock.Anything, uint(7)).Return(milestoneOf(), nil)
		projects.On("FindByID", mock.Anything, uint(10)).Return(ownedProjectOf(10), nil)
		milestones.On("Update", mock.Anything, mock.MatchedBy(func(ms *model.ProjectMilestone) bool {
			return ms.Title == "新題" && ms.Description == "旧説明" && ms.DueDate == nil
		})).Return(nil)
		uc := usecase.NewUpdateProjectMilestoneUseCase(milestones, projects)

		got, err := uc.Execute(context.Background(), usecase.UpdateProjectMilestoneInput{
			UserID: 1, MilestoneID: 7, Title: "  新題  ",
		})

		assert.NoError(t, err)
		assert.Equal(t, "新題", got.Title)
		milestones.AssertExpectations(t)
	})

	t.Run("完了へ遷移すると完了日時が入る", func(t *testing.T) {
		milestones, projects := new(mockProjectMilestoneRepo), new(mockProjectReader)
		milestones.On("FindByID", mock.Anything, uint(7)).Return(milestoneOf(), nil)
		projects.On("FindByID", mock.Anything, uint(10)).Return(ownedProjectOf(10), nil)
		milestones.On("Update", mock.Anything, mock.MatchedBy(func(ms *model.ProjectMilestone) bool {
			return ms.Status == model.MilestoneCompleted && ms.CompletedAt != nil
		})).Return(nil)
		uc := usecase.NewUpdateProjectMilestoneUseCase(milestones, projects)

		_, err := uc.Execute(context.Background(), usecase.UpdateProjectMilestoneInput{
			UserID: 1, MilestoneID: 7, Status: string(model.MilestoneCompleted),
		})

		assert.NoError(t, err)
		milestones.AssertExpectations(t)
	})

	t.Run("完了以外へ遷移すると完了日時が消える", func(t *testing.T) {
		milestones, projects := new(mockProjectMilestoneRepo), new(mockProjectReader)
		now := time.Now()
		existing := milestoneOf()
		existing.Status = model.MilestoneCompleted
		existing.CompletedAt = &now
		milestones.On("FindByID", mock.Anything, uint(7)).Return(existing, nil)
		projects.On("FindByID", mock.Anything, uint(10)).Return(ownedProjectOf(10), nil)
		milestones.On("Update", mock.Anything, mock.MatchedBy(func(ms *model.ProjectMilestone) bool {
			return ms.Status == model.MilestoneInProgress && ms.CompletedAt == nil
		})).Return(nil)
		uc := usecase.NewUpdateProjectMilestoneUseCase(milestones, projects)

		_, err := uc.Execute(context.Background(), usecase.UpdateProjectMilestoneInput{
			UserID: 1, MilestoneID: 7, Status: string(model.MilestoneInProgress),
		})

		assert.NoError(t, err)
		milestones.AssertExpectations(t)
	})

	t.Run("無効なステータスは Validation エラー（保存しない）", func(t *testing.T) {
		milestones, projects := new(mockProjectMilestoneRepo), new(mockProjectReader)
		milestones.On("FindByID", mock.Anything, uint(7)).Return(milestoneOf(), nil)
		projects.On("FindByID", mock.Anything, uint(10)).Return(ownedProjectOf(10), nil)
		uc := usecase.NewUpdateProjectMilestoneUseCase(milestones, projects)

		_, err := uc.Execute(context.Background(), usecase.UpdateProjectMilestoneInput{
			UserID: 1, MilestoneID: 7, Status: "done",
		})

		assertDomainCode(t, err, domain.ErrCodeValidation)
		milestones.AssertNotCalled(t, "Update")
	})

	t.Run("マイルストーン不在は NotFound（所有権も見ない）", func(t *testing.T) {
		milestones, projects := new(mockProjectMilestoneRepo), new(mockProjectReader)
		milestones.On("FindByID", mock.Anything, uint(7)).Return(nil, nil)
		uc := usecase.NewUpdateProjectMilestoneUseCase(milestones, projects)

		_, err := uc.Execute(context.Background(), usecase.UpdateProjectMilestoneInput{UserID: 1, MilestoneID: 7})

		assertDomainCode(t, err, domain.ErrCodeNotFound)
		projects.AssertNotCalled(t, "FindByID")
		milestones.AssertNotCalled(t, "Update")
	})

	t.Run("所有者以外は Forbidden（保存しない）", func(t *testing.T) {
		milestones, projects := new(mockProjectMilestoneRepo), new(mockProjectReader)
		milestones.On("FindByID", mock.Anything, uint(7)).Return(milestoneOf(), nil)
		projects.On("FindByID", mock.Anything, uint(10)).Return(otherProjectOf(10), nil)
		uc := usecase.NewUpdateProjectMilestoneUseCase(milestones, projects)

		_, err := uc.Execute(context.Background(), usecase.UpdateProjectMilestoneInput{UserID: 1, MilestoneID: 7})

		assertDomainCode(t, err, domain.ErrCodeForbidden)
		milestones.AssertNotCalled(t, "Update")
	})

	t.Run("保存時の DB 障害を伝播する", func(t *testing.T) {
		milestones, projects := new(mockProjectMilestoneRepo), new(mockProjectReader)
		milestones.On("FindByID", mock.Anything, uint(7)).Return(milestoneOf(), nil)
		projects.On("FindByID", mock.Anything, uint(10)).Return(ownedProjectOf(10), nil)
		milestones.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
		uc := usecase.NewUpdateProjectMilestoneUseCase(milestones, projects)

		_, err := uc.Execute(context.Background(), usecase.UpdateProjectMilestoneInput{UserID: 1, MilestoneID: 7})

		assert.Error(t, err)
		milestones.AssertExpectations(t)
	})
}

func TestDeleteProjectMilestoneUseCase_Execute(t *testing.T) {
	t.Run("所有者なら削除する", func(t *testing.T) {
		milestones, projects := new(mockProjectMilestoneRepo), new(mockProjectReader)
		milestones.On("FindByID", mock.Anything, uint(7)).Return(milestoneOf(), nil)
		projects.On("FindByID", mock.Anything, uint(10)).Return(ownedProjectOf(10), nil)
		milestones.On("Delete", mock.Anything, uint(7)).Return(nil)
		uc := usecase.NewDeleteProjectMilestoneUseCase(milestones, projects)

		assert.NoError(t, uc.Execute(context.Background(), 1, 7))
		milestones.AssertExpectations(t)
	})

	t.Run("所有者以外は Forbidden（削除しない）", func(t *testing.T) {
		milestones, projects := new(mockProjectMilestoneRepo), new(mockProjectReader)
		milestones.On("FindByID", mock.Anything, uint(7)).Return(milestoneOf(), nil)
		projects.On("FindByID", mock.Anything, uint(10)).Return(otherProjectOf(10), nil)
		uc := usecase.NewDeleteProjectMilestoneUseCase(milestones, projects)

		assertDomainCode(t, uc.Execute(context.Background(), 1, 7), domain.ErrCodeForbidden)
		milestones.AssertNotCalled(t, "Delete")
	})

	t.Run("マイルストーン不在は NotFound（削除しない）", func(t *testing.T) {
		milestones, projects := new(mockProjectMilestoneRepo), new(mockProjectReader)
		milestones.On("FindByID", mock.Anything, uint(7)).Return(nil, nil)
		uc := usecase.NewDeleteProjectMilestoneUseCase(milestones, projects)

		assertDomainCode(t, uc.Execute(context.Background(), 1, 7), domain.ErrCodeNotFound)
		milestones.AssertNotCalled(t, "Delete")
	})
}

package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestStudyCircleService() (*StudyCircleService, *MockStudyCircleRepository) {
	repo := new(MockStudyCircleRepository)
	svc := NewStudyCircleService(repo)
	return svc, repo
}

// --- Create ---

func TestStudyCircleCreate_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{Name: "Go勉強会", Topic: "Go", OwnerID: 1, MaxMembers: 5}
	repo.On("Create", circle).Return(nil)
	repo.On("AddMember", uint(0), uint(1), model.StudyCircleRoleOwner).Return(nil)

	err := svc.Create(circle, nil)
	assert.NoError(t, err)
	assert.Equal(t, model.StudyCircleStatusActive, circle.Status)
	repo.AssertExpectations(t)
}

func TestStudyCircleCreate_WhitespaceName(t *testing.T) {
	svc, _ := newTestStudyCircleService()

	circle := &model.StudyCircle{Name: "   \t  ", Topic: "Go", OwnerID: 1, MaxMembers: 5}
	err := svc.Create(circle, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "サークル名は必須です")
}

func TestStudyCircleCreate_EmptyName(t *testing.T) {
	svc, _ := newTestStudyCircleService()

	circle := &model.StudyCircle{Name: "", Topic: "Go", OwnerID: 1, MaxMembers: 5}
	err := svc.Create(circle, nil)
	assert.Error(t, err)
}

func TestStudyCircleCreate_WithMembers(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{Name: "React勉強会", Topic: "React", OwnerID: 1, MaxMembers: 5}
	repo.On("Create", circle).Return(nil)
	repo.On("AddMember", uint(0), uint(1), model.StudyCircleRoleOwner).Return(nil)
	repo.On("AddMember", uint(0), uint(2), model.StudyCircleRoleMember).Return(nil)
	repo.On("AddMember", uint(0), uint(3), model.StudyCircleRoleMember).Return(nil)

	err := svc.Create(circle, []uint{1, 2, 3}) // 1はオーナー自身→スキップ
	assert.NoError(t, err)
	repo.AssertNumberOfCalls(t, "AddMember", 3) // owner + 2 members
}

func TestStudyCircleCreate_NormalizesMaxMembers(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{Name: "Test", Topic: "Test", OwnerID: 1, MaxMembers: 100}
	repo.On("Create", circle).Return(nil)
	repo.On("AddMember", uint(0), uint(1), model.StudyCircleRoleOwner).Return(nil)

	err := svc.Create(circle, nil)
	assert.NoError(t, err)
	assert.Equal(t, 5, circle.MaxMembers) // 3〜10の範囲外 → 5にリセット
}

func TestStudyCircleCreate_CreateError(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{Name: "Test", Topic: "Test", OwnerID: 1, MaxMembers: 5}
	repo.On("Create", circle).Return(errors.New("db error"))

	err := svc.Create(circle, nil)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

// --- GetByID ---

func TestStudyCircleGetByID_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	expected := &model.StudyCircle{ID: 1, Name: "Go勉強会", OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(expected, nil)
	repo.On("IsMember", uint(1), uint(2)).Return(true, nil)

	result, err := svc.GetByID(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, "Go勉強会", result.Name)
}

func TestStudyCircleGetByID_NotFound(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	result, err := svc.GetByID(99, 1)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestStudyCircleGetByID_Forbidden(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("IsMember", uint(1), uint(99)).Return(false, nil)

	result, err := svc.GetByID(1, 99)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrForbidden)
}

// --- Update ---

func TestStudyCircleUpdate_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, Name: "旧名", Topic: "Go", OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("Update", circle).Return(nil)

	newName := "新名"
	result, err := svc.Update(1, 1, &newName, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, "新名", result.Name)
}

func TestStudyCircleUpdate_Forbidden(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)

	newName := "変更"
	result, err := svc.Update(1, 99, &newName, nil, nil)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestStudyCircleUpdate_RepoError(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, Name: "旧名", OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("Update", circle).Return(errors.New("db error"))

	newName := "新名"
	result, err := svc.Update(1, 1, &newName, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// --- Delete ---

func TestStudyCircleDelete_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
}

func TestStudyCircleDelete_Forbidden(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)

	err := svc.Delete(1, 99)
	assert.ErrorIs(t, err, ErrForbidden)
}

// --- AddMember ---

func TestStudyCircleAddMember_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, MaxMembers: 5}
	repo.On("IsMember", uint(1), uint(1)).Return(true, nil)
	repo.On("IsMember", uint(1), uint(5)).Return(false, nil)
	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("GetMemberCount", uint(1)).Return(3, nil)
	repo.On("AddMember", uint(1), uint(5), model.StudyCircleRoleMember).Return(nil)

	err := svc.AddMember(1, 1, 5)
	assert.NoError(t, err)
}

func TestStudyCircleAddMember_AlreadyMember(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("IsMember", uint(1), uint(1)).Return(true, nil)
	repo.On("IsMember", uint(1), uint(5)).Return(true, nil)

	err := svc.AddMember(1, 1, 5)
	assert.ErrorIs(t, err, ErrBadRequest)
	repo.AssertNotCalled(t, "AddMember")
}

func TestStudyCircleAddMember_MemberLimitReached(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, MaxMembers: 5}
	repo.On("IsMember", uint(1), uint(1)).Return(true, nil)
	repo.On("IsMember", uint(1), uint(10)).Return(false, nil)
	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("GetMemberCount", uint(1)).Return(5, nil)

	err := svc.AddMember(1, 1, 10)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
}

func TestStudyCircleAddMember_Forbidden(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("IsMember", uint(1), uint(99)).Return(false, nil)

	err := svc.AddMember(1, 99, 5)
	assert.ErrorIs(t, err, ErrForbidden)
}

// --- RemoveMember ---

func TestStudyCircleRemoveMember_SelfLeave(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("RemoveMember", uint(1), uint(2)).Return(nil)

	err := svc.RemoveMember(1, 2, 2)
	assert.NoError(t, err)
}

func TestStudyCircleRemoveMember_OwnerRemoves(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("RemoveMember", uint(1), uint(3)).Return(nil)

	err := svc.RemoveMember(1, 1, 3)
	assert.NoError(t, err)
}

func TestStudyCircleRemoveMember_Forbidden(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)

	err := svc.RemoveMember(1, 2, 3) // 非オーナーが他人を除外
	assert.ErrorIs(t, err, ErrForbidden)
}

// --- CreateCheckin ---

func TestStudyCircleCreateCheckin_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("IsMember", uint(1), uint(2)).Return(true, nil)
	repo.On("HasCheckedInToday", uint(1), uint(2)).Return(false, nil)
	repo.On("CreateCheckin", mock.AnythingOfType("*model.StudyCircleCheckin")).Return(nil)

	checkin, err := svc.CreateCheckin(1, 2, "今日はGo学習した")
	assert.NoError(t, err)
	assert.Equal(t, "今日はGo学習した", checkin.Content)
	assert.Equal(t, uint(1), checkin.CircleID)
	assert.Equal(t, uint(2), checkin.UserID)
}

func TestStudyCircleCreateCheckin_AlreadyCheckedIn(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("IsMember", uint(1), uint(2)).Return(true, nil)
	repo.On("HasCheckedInToday", uint(1), uint(2)).Return(true, nil)

	checkin, err := svc.CreateCheckin(1, 2, "二回目")
	assert.Nil(t, checkin)
	assert.Error(t, err)

	var domainErr *domain.DomainError
	assert.True(t, errors.As(err, &domainErr))
	assert.Equal(t, domain.ErrCodeConflict, domainErr.Code)
}

func TestStudyCircleCreateCheckin_Forbidden(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("IsMember", uint(1), uint(99)).Return(false, nil)

	checkin, err := svc.CreateCheckin(1, 99, "test")
	assert.Nil(t, checkin)
	assert.ErrorIs(t, err, ErrForbidden)
}

// --- CreateStep ---

func TestStudyCircleCreateStep_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("CreateStep", mock.AnythingOfType("*model.StudyCircleStep")).Return(nil)

	step := &model.StudyCircleStep{Title: "ステップ1"}
	err := svc.CreateStep(1, 1, step)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), step.CircleID)
}

func TestStudyCircleCreateStep_Forbidden(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)

	step := &model.StudyCircleStep{Title: "ステップ1"}
	err := svc.CreateStep(1, 99, step)
	assert.ErrorIs(t, err, ErrForbidden)
}

// --- GetStreakRanking ---

func TestStudyCircleGetStreakRanking_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("IsMember", uint(1), uint(2)).Return(true, nil)
	ranking := []model.CircleMemberStreak{
		{UserID: 1, UserName: "Alice", CurrentStreak: 10},
		{UserID: 2, UserName: "Bob", CurrentStreak: 5},
	}
	repo.On("GetStreakRanking", uint(1)).Return(ranking, nil)

	result, err := svc.GetStreakRanking(1, 2)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 10, result[0].CurrentStreak)
}

func TestStudyCircleGetStreakRanking_Forbidden(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("IsMember", uint(1), uint(99)).Return(false, nil)

	result, err := svc.GetStreakRanking(1, 99)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrForbidden)
}

// --- GetMyCircles ---

func TestStudyCircleGetMyCircles_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circles := []model.StudyCircle{
		{ID: 1, Name: "Go勉強会"},
		{ID: 2, Name: "React勉強会"},
	}
	repo.On("FindByUserID", uint(1), 20, 0).Return(circles, int64(2), nil)

	result, total, err := svc.GetMyCircles(1, 20, 0)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(2), total)
	repo.AssertExpectations(t)
}

func TestStudyCircleGetMyCircles_Empty(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("FindByUserID", uint(99), 20, 0).Return([]model.StudyCircle{}, int64(0), nil)

	result, total, err := svc.GetMyCircles(99, 20, 0)
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
	repo.AssertExpectations(t)
}

// --- GetMembers ---

func TestStudyCircleGetMembers_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	members := []model.StudyCircleMember{
		{CircleID: 1, UserID: 1, Role: model.StudyCircleRoleOwner},
		{CircleID: 1, UserID: 2, Role: model.StudyCircleRoleMember},
	}
	repo.On("IsMember", uint(1), uint(1)).Return(true, nil)
	repo.On("GetMembers", uint(1)).Return(members, nil)

	result, err := svc.GetMembers(1, 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestStudyCircleGetMembers_Forbidden(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("IsMember", uint(1), uint(99)).Return(false, nil)

	result, err := svc.GetMembers(1, 99)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrForbidden)
}

// --- UpdateStep ---

func TestStudyCircleUpdateStep_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	step := &model.StudyCircleStep{CircleID: 1, Title: "旧タイトル"}
	step.ID = 5

	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("UpdateStep", step).Return(nil)

	newTitle := "新タイトル"
	result, err := svc.UpdateStep(1, 1, 5, &newTitle, nil)
	assert.NoError(t, err)
	assert.Equal(t, "新タイトル", result.Title)
	repo.AssertExpectations(t)
}

func TestStudyCircleUpdateStep_Forbidden(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)

	newTitle := "変更"
	result, err := svc.UpdateStep(1, 99, 5, &newTitle, nil)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestStudyCircleUpdateStep_StepBelongsToDifferentCircle(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	step := &model.StudyCircleStep{CircleID: 2, Title: "別サークル"}
	step.ID = 5

	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)

	newTitle := "変更"
	result, err := svc.UpdateStep(1, 1, 5, &newTitle, nil)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrNotFound)
}

// --- DeleteStep ---

func TestStudyCircleDeleteStep_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	step := &model.StudyCircleStep{CircleID: 1}
	step.ID = 5

	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("DeleteStep", uint(5)).Return(nil)

	err := svc.DeleteStep(1, 1, 5)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestStudyCircleDeleteStep_Forbidden(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)

	err := svc.DeleteStep(1, 99, 5)
	assert.ErrorIs(t, err, ErrForbidden)
}

// --- ReorderSteps ---

func TestStudyCircleReorderSteps_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	orders := []model.StepOrder{
		{StepID: 1, OrderIndex: 0},
		{StepID: 2, OrderIndex: 1},
	}

	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("ReorderSteps", uint(1), orders).Return(nil)

	err := svc.ReorderSteps(1, 1, orders)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestStudyCircleReorderSteps_Forbidden(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)

	orders := []model.StepOrder{{StepID: 1, OrderIndex: 0}}
	err := svc.ReorderSteps(1, 99, orders)
	assert.ErrorIs(t, err, ErrForbidden)
}

// --- UpdateProgress ---

func TestStudyCircleUpdateProgress_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("IsMember", uint(1), uint(2)).Return(true, nil)
	repo.On("UpsertProgress", mock.AnythingOfType("*model.StudyCircleMemberProgress")).Return(nil)

	err := svc.UpdateProgress(1, 2, 5, true)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestStudyCircleUpdateProgress_Forbidden(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("IsMember", uint(1), uint(99)).Return(false, nil)

	err := svc.UpdateProgress(1, 99, 5, true)
	assert.ErrorIs(t, err, ErrForbidden)
}

// --- GetProgress ---

func TestStudyCircleGetProgress_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	progress := []model.StudyCircleMemberProgress{
		{CircleID: 1, UserID: 1, StepID: 1, IsCompleted: true},
		{CircleID: 1, UserID: 2, StepID: 1, IsCompleted: false},
	}
	repo.On("IsMember", uint(1), uint(1)).Return(true, nil)
	repo.On("GetProgress", uint(1)).Return(progress, nil)

	result, err := svc.GetProgress(1, 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestStudyCircleGetProgress_Forbidden(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("IsMember", uint(1), uint(99)).Return(false, nil)

	result, err := svc.GetProgress(1, 99)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrForbidden)
}

// --- GetCheckins ---

func TestStudyCircleGetCheckins_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	checkins := []model.StudyCircleCheckin{
		{CircleID: 1, UserID: 1, Content: "Go学習"},
		{CircleID: 1, UserID: 2, Content: "React学習"},
	}
	repo.On("IsMember", uint(1), uint(1)).Return(true, nil)
	repo.On("GetCheckins", uint(1)).Return(checkins, nil)

	result, err := svc.GetCheckins(1, 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertExpectations(t)
}

func TestStudyCircleGetCheckins_Forbidden(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("IsMember", uint(1), uint(99)).Return(false, nil)

	result, err := svc.GetCheckins(1, 99)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrForbidden)
}

// --- エラーパステスト ---

func TestStudyCircleCreate_AddMemberOwnerError(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	circle := &model.StudyCircle{Name: "Test", Topic: "Go", OwnerID: 1, MaxMembers: 5}
	repo.On("Create", circle).Return(nil)
	repo.On("AddMember", uint(0), uint(1), model.StudyCircleRoleOwner).Return(errors.New("db error"))
	err := svc.Create(circle, nil)
	assert.Error(t, err)
}

func TestStudyCircleGetByID_IsMemberError(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("IsMember", uint(1), uint(2)).Return(false, errors.New("db error"))
	result, err := svc.GetByID(1, 2)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestStudyCircleUpdate_NotFound(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))
	name := "新名"
	result, err := svc.Update(99, 1, &name, nil, nil)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestStudyCircleDelete_NotFound(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))
	err := svc.Delete(99, 1)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestStudyCircleDeleteStep_StepNotFound(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("FindStepByID", uint(99)).Return(nil, errors.New("not found"))
	err := svc.DeleteStep(1, 1, 99)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestStudyCircleDeleteStep_CircleNotFound(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))
	err := svc.DeleteStep(99, 1, 5)
	assert.ErrorIs(t, err, ErrNotFound)
	repo.AssertExpectations(t)
}

func TestStudyCircleDeleteStep_StepBelongsToDifferentCircle(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	step := &model.StudyCircleStep{CircleID: 2}
	step.ID = 5
	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	err := svc.DeleteStep(1, 1, 5)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestStudyCircleUpdateStep_StepNotFound(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("FindStepByID", uint(99)).Return(nil, errors.New("not found"))
	title := "新タイトル"
	result, err := svc.UpdateStep(1, 1, 99, &title, nil)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestStudyCircleUpdateStep_RepoError(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	step := &model.StudyCircleStep{CircleID: 1, Title: "旧"}
	step.ID = 5
	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("UpdateStep", step).Return(errors.New("db error"))
	title := "新"
	result, err := svc.UpdateStep(1, 1, 5, &title, nil)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestStudyCircleAddMember_IsMemberError(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	repo.On("IsMember", uint(1), uint(1)).Return(false, errors.New("db error"))
	err := svc.AddMember(1, 1, 5)
	assert.Error(t, err)
}

func TestStudyCircleAddMember_FindByIDError(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	repo.On("IsMember", uint(1), uint(1)).Return(true, nil)
	repo.On("IsMember", uint(1), uint(5)).Return(false, nil)
	repo.On("FindByID", uint(1)).Return(nil, errors.New("not found"))
	err := svc.AddMember(1, 1, 5)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestStudyCircleAddMember_GetMemberCountError(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	circle := &model.StudyCircle{ID: 1, MaxMembers: 5}
	repo.On("IsMember", uint(1), uint(1)).Return(true, nil)
	repo.On("IsMember", uint(1), uint(5)).Return(false, nil)
	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("GetMemberCount", uint(1)).Return(0, errors.New("db error"))
	err := svc.AddMember(1, 1, 5)
	assert.Error(t, err)
}

func TestStudyCircleAddMember_TargetIsMemberError(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	repo.On("IsMember", uint(1), uint(1)).Return(true, nil)
	repo.On("IsMember", uint(1), uint(5)).Return(false, errors.New("db error"))
	err := svc.AddMember(1, 1, 5)
	assert.Error(t, err)
	repo.AssertNotCalled(t, "FindByID")
	repo.AssertExpectations(t)
}

func TestStudyCircleRemoveMember_NotFound(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))
	err := svc.RemoveMember(99, 1, 2)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestStudyCircleGetMembers_IsMemberError(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	repo.On("IsMember", uint(1), uint(1)).Return(false, errors.New("db error"))
	result, err := svc.GetMembers(1, 1)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestStudyCircleGetProgress_IsMemberError(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	repo.On("IsMember", uint(1), uint(1)).Return(false, errors.New("db error"))
	result, err := svc.GetProgress(1, 1)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestStudyCircleGetCheckins_IsMemberError(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	repo.On("IsMember", uint(1), uint(1)).Return(false, errors.New("db error"))
	result, err := svc.GetCheckins(1, 1)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestStudyCircleGetStreakRanking_IsMemberError(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	repo.On("IsMember", uint(1), uint(1)).Return(false, errors.New("db error"))
	result, err := svc.GetStreakRanking(1, 1)
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestStudyCircleCreateCheckin_HasCheckedInTodayError(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	repo.On("IsMember", uint(1), uint(2)).Return(true, nil)
	repo.On("HasCheckedInToday", uint(1), uint(2)).Return(false, errors.New("db error"))
	result, err := svc.CreateCheckin(1, 2, "test")
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestStudyCircleCreateCheckin_RepoError(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	repo.On("IsMember", uint(1), uint(2)).Return(true, nil)
	repo.On("HasCheckedInToday", uint(1), uint(2)).Return(false, nil)
	repo.On("CreateCheckin", mock.AnythingOfType("*model.StudyCircleCheckin")).Return(errors.New("db error"))
	result, err := svc.CreateCheckin(1, 2, "test")
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestStudyCircleCreateStep_NotFound(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))
	step := &model.StudyCircleStep{Title: "ステップ1"}
	err := svc.CreateStep(99, 1, step)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestStudyCircleUpdateProgress_IsMemberError(t *testing.T) {
	svc, repo := newTestStudyCircleService()
	repo.On("IsMember", uint(1), uint(1)).Return(false, errors.New("db error"))
	err := svc.UpdateProgress(1, 1, 5, true)
	assert.Error(t, err)
}

// ============================================================
// SearchCircles テスト
// ============================================================

func TestSearchCircles_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	expected := []model.StudyCircle{
		{Name: "Golang勉強会", Topic: "プログラミング"},
		{Name: "Go入門", Topic: "バックエンド"},
	}
	repo.On("Search", "Go", 20, 0).Return(expected, int64(2), nil)

	result, total, err := svc.SearchCircles("Go", 20, 0)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	circles, ok := result.([]model.StudyCircle)
	assert.True(t, ok)
	assert.Len(t, circles, 2)
	repo.AssertExpectations(t)
}

func TestSearchCircles_EmptyResult(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("Search", "存在しない", 20, 0).Return([]model.StudyCircle{}, int64(0), nil)

	result, total, err := svc.SearchCircles("存在しない", 20, 0)

	assert.NoError(t, err)
	assert.Equal(t, int64(0), total)
	circles, ok := result.([]model.StudyCircle)
	assert.True(t, ok)
	assert.Len(t, circles, 0)
}

func TestSearchCircles_RepoError(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("Search", "Go", 20, 0).Return(nil, int64(0), errors.New("db error"))

	result, _, err := svc.SearchCircles("Go", 20, 0)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestStudyCircleUpdateStep_WithDescription(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{OwnerID: 1}
	circle.ID = 1
	step := &model.StudyCircleStep{CircleID: 1, Title: "Title", Description: "Old Desc"}
	step.ID = 5

	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("FindStepByID", uint(5)).Return(step, nil)
	repo.On("UpdateStep", step).Return(nil)

	newDesc := "New Description"
	result, err := svc.UpdateStep(1, 1, 5, nil, &newDesc)
	assert.NoError(t, err)
	assert.Equal(t, "New Description", result.Description)
	repo.AssertExpectations(t)
}

func TestStudyCircleReorderSteps_NotFound(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	err := svc.ReorderSteps(99, 1, []model.StepOrder{})
	assert.ErrorIs(t, err, ErrNotFound)
	repo.AssertExpectations(t)
}

func TestStudyCircleUpdate_FindByIDError(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	newName := "変更"
	result, err := svc.Update(99, 1, &newName, nil, nil)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrNotFound)
	repo.AssertExpectations(t)
}

func TestStudyCircleUpdate_TopicAndDescription(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("Update", circle).Return(nil)

	newTopic := "新トピック"
	newDesc := "新説明"
	result, err := svc.Update(1, 1, nil, &newTopic, &newDesc)
	assert.NoError(t, err)
	assert.Equal(t, "新トピック", result.Topic)
	assert.Equal(t, "新説明", result.Description)
	repo.AssertExpectations(t)
}

func TestStudyCircleUpdateStep_CircleNotFound(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))

	newTitle := "Title"
	result, err := svc.UpdateStep(99, 1, 5, &newTitle, nil)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// --- Update 空白バイパス ---

func TestStudyCircleUpdate_WhitespaceName(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, Name: "旧名", OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)

	name := "   "
	result, err := svc.Update(1, 1, &name, nil, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "サークル名")
}

func TestStudyCircleUpdate_WhitespaceTopic(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, Name: "旧名", Topic: "Go", OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)

	topic := "  \t  "
	result, err := svc.Update(1, 1, nil, &topic, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "トピック")
}

func TestStudyCircleUpdateStep_WhitespaceTitle(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, OwnerID: 1}
	repo.On("FindByID", uint(1)).Return(circle, nil)
	step := &model.StudyCircleStep{ID: 5, CircleID: 1, Title: "旧タイトル"}
	repo.On("FindStepByID", uint(5)).Return(step, nil)

	title := "   "
	result, err := svc.UpdateStep(1, 1, 5, &title, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "タイトル")
}

func TestStudyCircleCreateCheckin_IsMemberError(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("IsMember", uint(1), uint(1)).Return(false, errors.New("db error"))

	result, err := svc.CreateCheckin(1, 1, "content")
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// GetByStatus テスト
// ============================================================

func TestStudyCircleGetByStatus_Success(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circles := []model.StudyCircle{
		{ID: 1, Name: "Go勉強会", Status: model.StudyCircleStatusActive},
	}
	repo.On("GetByStatus", uint(1), "active").Return(circles, nil)

	result, err := svc.GetByStatus(1, "active")
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, model.StudyCircleStatusActive, result[0].Status)
	repo.AssertExpectations(t)
}

func TestStudyCircleGetByStatus_EmptyResult(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("GetByStatus", uint(1), "archived").Return([]model.StudyCircle{}, nil)

	result, err := svc.GetByStatus(1, "archived")
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestStudyCircleGetByStatus_InvalidStatus(t *testing.T) {
	svc, _ := newTestStudyCircleService()

	result, err := svc.GetByStatus(1, "invalid")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "無効なステータス")
}

func TestStudyCircleGetByStatus_RepoError(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	repo.On("GetByStatus", uint(1), "active").Return([]model.StudyCircle{}, errors.New("db error"))

	result, err := svc.GetByStatus(1, "active")
	assert.Error(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

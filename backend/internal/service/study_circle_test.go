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
	repo.On("FindByID", uint(1)).Return(circle, nil)
	repo.On("GetMemberCount", uint(1)).Return(3, nil)
	repo.On("AddMember", uint(1), uint(5), model.StudyCircleRoleMember).Return(nil)

	err := svc.AddMember(1, 1, 5)
	assert.NoError(t, err)
}

func TestStudyCircleAddMember_MemberLimitReached(t *testing.T) {
	svc, repo := newTestStudyCircleService()

	circle := &model.StudyCircle{ID: 1, MaxMembers: 5}
	repo.On("IsMember", uint(1), uint(1)).Return(true, nil)
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

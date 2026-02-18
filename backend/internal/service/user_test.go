package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// newTestUserService はUserServiceのテスト用インスタンスを生成するヘルパー。
func newTestUserService() (*UserService, *MockUserRepository) {
	repo := new(MockUserRepository)
	svc := NewUserService(repo)
	return svc, repo
}

// ============================================================
// GetAll（クエリ有→Search / クエリ空→FindAll）
// ============================================================

func TestUserGetAll_WithQuery(t *testing.T) {
	svc, repo := newTestUserService()

	users := []model.User{{Name: "Alice"}}
	repo.On("Search", "alice").Return(users, nil)

	result, err := svc.GetAll("alice")
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Alice", result[0].Name)
	repo.AssertCalled(t, "Search", "alice")
	repo.AssertNotCalled(t, "FindAll")
}

func TestUserGetAll_EmptyQuery(t *testing.T) {
	svc, repo := newTestUserService()

	users := []model.User{{Name: "Alice"}, {Name: "Bob"}}
	repo.On("FindAll").Return(users, nil)

	result, err := svc.GetAll("")
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	repo.AssertCalled(t, "FindAll")
	repo.AssertNotCalled(t, "Search")
}

// ============================================================
// ユーザーID検索テスト
// ============================================================

func TestUserGetByID_Success(t *testing.T) {
	svc, repo := newTestUserService()

	user := &model.User{Name: "Alice"}
	user.ID = 1
	repo.On("FindByID", uint(1)).Return(user, nil)

	result, err := svc.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", result.Name)
	repo.AssertExpectations(t)
}

// ============================================================
// ユーザー更新テスト
// ============================================================

func TestUserUpdate_Success(t *testing.T) {
	svc, repo := newTestUserService()

	user := &model.User{Name: "Alice"}
	repo.On("Update", user).Return(nil)

	err := svc.Update(user)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

// ============================================================
// GetProfileCompleteness テスト
// ============================================================

func TestGetProfileCompleteness_FullProfile(t *testing.T) {
	svc, repo := newTestUserService()

	user := &model.User{
		ID:              1,
		Name:            "Alice",
		Username:        "alice",
		Email:           "alice@example.com",
		AvatarURL:       "https://example.com/avatar.png",
		Bio:             "I love coding",
		GitHubConnected: true,
		SkillsLanguages: "Go,TypeScript",
	}
	repo.On("FindByID", uint(1)).Return(user, nil)

	result, err := svc.GetProfileCompleteness(1)
	assert.NoError(t, err)
	assert.Equal(t, 100, result.Percentage)
	assert.Empty(t, result.MissingFields)
}

func TestGetProfileCompleteness_EmptyProfile(t *testing.T) {
	svc, repo := newTestUserService()

	user := &model.User{
		ID:       1,
		Name:     "Alice",
		Username: "alice",
		Email:    "alice@example.com",
	}
	repo.On("FindByID", uint(1)).Return(user, nil)

	result, err := svc.GetProfileCompleteness(1)
	assert.NoError(t, err)
	assert.Less(t, result.Percentage, 100)
	assert.Contains(t, result.MissingFields, "avatar")
	assert.Contains(t, result.MissingFields, "bio")
	assert.Contains(t, result.MissingFields, "github")
	assert.Contains(t, result.MissingFields, "skills")
}

func TestGetProfileCompleteness_PartialProfile(t *testing.T) {
	svc, repo := newTestUserService()

	user := &model.User{
		ID:        1,
		Name:      "Alice",
		Username:  "alice",
		Email:     "alice@example.com",
		AvatarURL: "https://example.com/avatar.png",
		Bio:       "Hello",
	}
	repo.On("FindByID", uint(1)).Return(user, nil)

	result, err := svc.GetProfileCompleteness(1)
	assert.NoError(t, err)
	assert.Greater(t, result.Percentage, 0)
	assert.Less(t, result.Percentage, 100)
	assert.NotContains(t, result.MissingFields, "avatar")
	assert.NotContains(t, result.MissingFields, "bio")
	assert.Contains(t, result.MissingFields, "github")
	assert.Contains(t, result.MissingFields, "skills")
}

func TestGetProfileCompleteness_UserNotFound(t *testing.T) {
	svc, repo := newTestUserService()

	repo.On("FindByID", uint(999)).Return((*model.User)(nil), errors.New("not found"))

	_, err := svc.GetProfileCompleteness(999)
	assert.Error(t, err)
}

func TestGetProfileCompleteness_WithSkillsFrameworks(t *testing.T) {
	svc, repo := newTestUserService()

	user := &model.User{
		ID:               1,
		Name:             "Alice",
		Username:         "alice",
		Email:            "alice@example.com",
		AvatarURL:        "https://example.com/avatar.png",
		Bio:              "Developer",
		GitHubConnected:  true,
		SkillsFrameworks: "React,Gin",
	}
	repo.On("FindByID", uint(1)).Return(user, nil)

	result, err := svc.GetProfileCompleteness(1)
	assert.NoError(t, err)
	assert.Equal(t, 100, result.Percentage)
	assert.Empty(t, result.MissingFields)
}

// ============================================================
// GetByUsername テスト
// ============================================================

func TestUserGetByUsername_Success(t *testing.T) {
	svc, repo := newTestUserService()

	user := &model.User{Name: "Alice", Username: "alice"}
	repo.On("FindByUsername", "alice").Return(user, nil)

	result, err := svc.GetByUsername("alice")
	assert.NoError(t, err)
	assert.Equal(t, "Alice", result.Name)
	repo.AssertExpectations(t)
}

func TestUserGetByUsername_NotFound(t *testing.T) {
	svc, repo := newTestUserService()

	repo.On("FindByUsername", "unknown").Return((*model.User)(nil), errors.New("not found"))

	_, err := svc.GetByUsername("unknown")
	assert.Error(t, err)
}

// ============================================================
// FindByID テスト
// ============================================================

func TestUserFindByID_Success(t *testing.T) {
	svc, repo := newTestUserService()

	user := &model.User{Name: "Alice"}
	user.ID = 1
	repo.On("FindByID", uint(1)).Return(user, nil)

	result, err := svc.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", result.Name)
	repo.AssertExpectations(t)
}

func TestUserGetAll_SearchError(t *testing.T) {
	svc, repo := newTestUserService()

	repo.On("Search", "fail").Return([]model.User(nil), errors.New("db error"))

	result, err := svc.GetAll("fail")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserGetAll_FindAllError(t *testing.T) {
	svc, repo := newTestUserService()

	repo.On("FindAll").Return([]model.User(nil), errors.New("db error"))

	result, err := svc.GetAll("")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserGetByID_NotFound(t *testing.T) {
	svc, repo := newTestUserService()

	repo.On("FindByID", uint(999)).Return((*model.User)(nil), errors.New("not found"))

	result, err := svc.GetByID(999)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserUpdate_Error(t *testing.T) {
	svc, repo := newTestUserService()

	user := &model.User{Name: "Alice"}
	repo.On("Update", user).Return(errors.New("db error"))

	err := svc.Update(user)
	assert.Error(t, err)
}

func TestUserFindByID_NotFound(t *testing.T) {
	svc, repo := newTestUserService()

	repo.On("FindByID", uint(999)).Return((*model.User)(nil), errors.New("not found"))

	result, err := svc.FindByID(999)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetProfileCompleteness_OnlyFrameworks(t *testing.T) {
	svc, repo := newTestUserService()

	user := &model.User{
		ID:               1,
		Name:             "Alice",
		Username:         "alice",
		Email:            "alice@example.com",
		SkillsFrameworks: "React,Next.js",
	}
	repo.On("FindByID", uint(1)).Return(user, nil)

	result, err := svc.GetProfileCompleteness(1)
	assert.NoError(t, err)
	assert.NotContains(t, result.MissingFields, "skills")
	assert.Contains(t, result.MissingFields, "avatar")
	assert.Contains(t, result.MissingFields, "bio")
	assert.Contains(t, result.MissingFields, "github")
}

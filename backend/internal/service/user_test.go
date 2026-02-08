package service

import (
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

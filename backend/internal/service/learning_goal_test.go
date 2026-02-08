package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// newTestLearningGoalService はLearningGoalServiceのテスト用インスタンスを生成するヘルパー。
func newTestLearningGoalService() (*LearningGoalService, *MockLearningGoalRepository) {
	repo := new(MockLearningGoalRepository)
	svc := NewLearningGoalService(repo)
	return svc, repo
}

// ============================================================
// 学習目標更新テスト
// ============================================================

func TestLearningGoalUpdate_Success(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	existing := &model.LearningGoal{
		Title:    "Old Title",
		UserID:   1,
		Status:   model.GoalStatusActive,
		Progress: 50,
	}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	updates := &model.LearningGoal{Title: "New Title"}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	repo.AssertExpectations(t)
}

func TestLearningGoalUpdate_Forbidden(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	existing := &model.LearningGoal{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	updates := &model.LearningGoal{Title: "New Title"}
	result, err := svc.Update(1, 999, updates)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

func TestLearningGoalUpdate_AutoCompleteAt100(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	existing := &model.LearningGoal{
		Title:    "Goal",
		UserID:   1,
		Status:   model.GoalStatusActive,
		Progress: 80,
	}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	// 進捗100%で自動完了
	updates := &model.LearningGoal{Progress: 100}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, 100, result.Progress)
	assert.Equal(t, model.GoalStatusCompleted, result.Status)
	assert.NotNil(t, result.CompletedAt)
	repo.AssertExpectations(t)
}

func TestLearningGoalUpdate_ProgressClampTo100(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	existing := &model.LearningGoal{
		Title:    "Goal",
		UserID:   1,
		Status:   model.GoalStatusActive,
		Progress: 90,
	}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	// 進捗101はクランプされて100になる
	updates := &model.LearningGoal{Progress: 150}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, 100, result.Progress)
	assert.Equal(t, model.GoalStatusCompleted, result.Status)
	repo.AssertExpectations(t)
}

func TestLearningGoalUpdate_SetCompletedStatus(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	existing := &model.LearningGoal{
		Title:  "Goal",
		UserID: 1,
		Status: model.GoalStatusActive,
	}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", existing).Return(nil)

	// ステータスを直接completedに変更
	updates := &model.LearningGoal{Status: model.GoalStatusCompleted}
	result, err := svc.Update(1, 1, updates)
	assert.NoError(t, err)
	assert.Equal(t, model.GoalStatusCompleted, result.Status)
	assert.NotNil(t, result.CompletedAt)
	repo.AssertExpectations(t)
}

func TestLearningGoalUpdate_NotFound(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	repo.On("FindByID", uint(999)).Return(nil, errors.New("not found"))

	updates := &model.LearningGoal{Title: "New Title"}
	result, err := svc.Update(999, 1, updates)
	assert.Error(t, err)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}

// ============================================================
// 学習目標削除テスト
// ============================================================

func TestLearningGoalDelete_Success(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	existing := &model.LearningGoal{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Delete", uint(1)).Return(nil)

	err := svc.Delete(1, 1)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLearningGoalDelete_Forbidden(t *testing.T) {
	svc, repo := newTestLearningGoalService()

	existing := &model.LearningGoal{UserID: 1}
	existing.ID = 1

	repo.On("FindByID", uint(1)).Return(existing, nil)

	err := svc.Delete(1, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	repo.AssertExpectations(t)
}

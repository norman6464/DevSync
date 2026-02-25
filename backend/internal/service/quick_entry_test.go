package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================
// クイックエントリー: 最近のカテゴリ取得テスト
// ============================================================

func TestGetRecentCategories_Success(t *testing.T) {
	svc, repo := newTestLearningLogService()

	categories := []string{"coding", "reading", "course"}
	repo.On("GetRecentCategories", uint(1), 5).Return(categories, nil)

	result, err := svc.GetRecentCategories(uint(1))
	assert.NoError(t, err)
	assert.Equal(t, categories, result)
	repo.AssertExpectations(t)
}

func TestGetRecentCategories_Empty(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetRecentCategories", uint(1), 5).Return([]string{}, nil)

	result, err := svc.GetRecentCategories(uint(1))
	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}

func TestGetRecentCategories_RepoError(t *testing.T) {
	svc, repo := newTestLearningLogService()

	repo.On("GetRecentCategories", uint(1), 5).Return([]string{}, ErrNotFound)

	_, err := svc.GetRecentCategories(uint(1))
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

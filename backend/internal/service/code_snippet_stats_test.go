package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

func newCodeSnippetStatsTestService() (*CodeSnippetStatsService, *MockCodeSnippetStatsRepository) {
	repo := new(MockCodeSnippetStatsRepository)
	svc := NewCodeSnippetStatsService(repo)
	return svc, repo
}

func TestCodeSnippetStatsService_GetCodeSnippetStats_Success(t *testing.T) {
	svc, repo := newCodeSnippetStatsTestService()
	expected := &model.CodeSnippetStats{
		TotalSnippets: 15,
		TotalComments: 42,
		LanguageCount: 5,
	}
	repo.On("GetCodeSnippetStats", uint(1)).Return(expected, nil)

	result, err := svc.GetCodeSnippetStats(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestCodeSnippetStatsService_GetCodeSnippetStats_InvalidUserID(t *testing.T) {
	svc, _ := newCodeSnippetStatsTestService()

	_, err := svc.GetCodeSnippetStats(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "userIDは必須です")
}

func TestCodeSnippetStatsService_GetCodeSnippetStats_RepoError(t *testing.T) {
	svc, repo := newCodeSnippetStatsTestService()
	repo.On("GetCodeSnippetStats", uint(1)).Return(nil, errors.New("db error"))

	_, err := svc.GetCodeSnippetStats(1)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestCodeSnippetStatsService_GetCodeSnippetStats_NoActivity(t *testing.T) {
	svc, repo := newCodeSnippetStatsTestService()
	expected := &model.CodeSnippetStats{}
	repo.On("GetCodeSnippetStats", uint(99)).Return(expected, nil)

	result, err := svc.GetCodeSnippetStats(99)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result.TotalSnippets)
	assert.Equal(t, int64(0), result.TotalComments)
	assert.Equal(t, int64(0), result.LanguageCount)
	repo.AssertExpectations(t)
}

package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
)

// ============================================================
// 純粋関数テスト
// ============================================================

func TestTrimString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"通常の文字列", "Go", "Go"},
		{"前後に空白", "  Go  ", "Go"},
		{"前後にタブ", "\tGo\t", "Go"},
		{"空白のみ", "   ", ""},
		{"空文字列", "", ""},
		{"内部の空白は保持", "Go Lang", "Go Lang"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trimString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSplitSkills(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"空文字列", "", nil},
		{"単一スキル", "Go", []string{"Go"}},
		{"複数スキル", "Go,Python,Rust", []string{"Go", "Python", "Rust"}},
		{"空白付き", " Go , Python , Rust ", []string{"Go", "Python", "Rust"}},
		{"末尾カンマ", "Go,Python,", []string{"Go", "Python"}},
		{"連続カンマ", "Go,,Python", []string{"Go", "Python"}},
		{"カンマのみ", ",,,", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitSkills(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseSkills(t *testing.T) {
	tests := []struct {
		name       string
		languages  string
		frameworks string
		expected   []string
	}{
		{"両方空", "", "", nil},
		{"言語のみ", "Go,Python", "", []string{"Go", "Python"}},
		{"フレームワークのみ", "", "React,Vue", []string{"React", "Vue"}},
		{"両方あり", "Go,Python", "React,Vue", []string{"Go", "Python", "React", "Vue"}},
		{"重複含む空白", " Go , ", " React , ", []string{"Go", "React"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseSkills(tt.languages, tt.frameworks)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// ============================================================
// サービスメソッドテスト
// ============================================================

func TestGetRecommendedUsers_Success(t *testing.T) {
	mockRepo := new(MockRecommendationRepository)
	mockUserRepo := new(MockUserRepository)
	svc := NewRecommendationService(mockRepo, mockUserRepo)

	user := &model.User{SkillsLanguages: "Go,Python", SkillsFrameworks: "React"}
	mockUserRepo.On("FindByID", uint(1)).Return(user, nil)

	expected := []model.RecommendedUser{
		{User: model.User{Username: "user2"}, CommonSkills: []string{"Go"}, MatchScore: 80},
	}
	mockRepo.On("GetRecommendedUsers", uint(1), []string{"Go", "Python", "React"}, 10).Return(expected, nil)

	result, err := svc.GetRecommendedUsers(1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockUserRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestGetRecommendedUsers_UserNotFound(t *testing.T) {
	mockRepo := new(MockRecommendationRepository)
	mockUserRepo := new(MockUserRepository)
	svc := NewRecommendationService(mockRepo, mockUserRepo)

	mockUserRepo.On("FindByID", uint(1)).Return(nil, errors.New("user not found"))

	result, err := svc.GetRecommendedUsers(1)
	assert.Error(t, err)
	assert.Nil(t, result)
	mockUserRepo.AssertExpectations(t)
}

func TestGetRecommendedUsers_NoSkills(t *testing.T) {
	mockRepo := new(MockRecommendationRepository)
	mockUserRepo := new(MockUserRepository)
	svc := NewRecommendationService(mockRepo, mockUserRepo)

	user := &model.User{SkillsLanguages: "", SkillsFrameworks: ""}
	mockUserRepo.On("FindByID", uint(1)).Return(user, nil)

	result, err := svc.GetRecommendedUsers(1)
	assert.NoError(t, err)
	assert.Empty(t, result)
	mockUserRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "GetRecommendedUsers")
}

func TestGetRecommendedUsers_RepoError(t *testing.T) {
	mockRepo := new(MockRecommendationRepository)
	mockUserRepo := new(MockUserRepository)
	svc := NewRecommendationService(mockRepo, mockUserRepo)

	user := &model.User{SkillsLanguages: "Go", SkillsFrameworks: ""}
	mockUserRepo.On("FindByID", uint(1)).Return(user, nil)
	mockRepo.On("GetRecommendedUsers", uint(1), []string{"Go"}, 10).Return([]model.RecommendedUser(nil), errors.New("db error"))

	result, err := svc.GetRecommendedUsers(1)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTrendingPosts_Success(t *testing.T) {
	mockRepo := new(MockRecommendationRepository)
	mockUserRepo := new(MockUserRepository)
	svc := NewRecommendationService(mockRepo, mockUserRepo)

	expected := []model.Post{{Title: "Trending Post"}}
	mockRepo.On("GetTrendingPosts", 10, 7).Return(expected, nil)

	result, err := svc.GetTrendingPosts()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetTrendingPosts_Error(t *testing.T) {
	mockRepo := new(MockRecommendationRepository)
	mockUserRepo := new(MockUserRepository)
	svc := NewRecommendationService(mockRepo, mockUserRepo)

	mockRepo.On("GetTrendingPosts", 10, 7).Return([]model.Post(nil), errors.New("db error"))

	result, err := svc.GetTrendingPosts()
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetTrendingResources_Success(t *testing.T) {
	mockRepo := new(MockRecommendationRepository)
	mockUserRepo := new(MockUserRepository)
	svc := NewRecommendationService(mockRepo, mockUserRepo)

	expected := []model.LearningResource{{Title: "Trending Resource"}}
	mockRepo.On("GetTrendingResources", 10, 30).Return(expected, nil)

	result, err := svc.GetTrendingResources()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetTrendingResources_Error(t *testing.T) {
	mockRepo := new(MockRecommendationRepository)
	mockUserRepo := new(MockUserRepository)
	svc := NewRecommendationService(mockRepo, mockUserRepo)

	mockRepo.On("GetTrendingResources", 10, 30).Return([]model.LearningResource(nil), errors.New("db error"))

	result, err := svc.GetTrendingResources()
	assert.Error(t, err)
	assert.Nil(t, result)
}

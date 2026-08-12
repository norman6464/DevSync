// mock_test.go はService層テスト用のモックリポジトリを定義する。
package service

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	usecaserepo "github.com/norman6464/devsync/backend/internal/usecase/repository"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/mock"
)

// ============================================================
// MockUserRepository は repository.UserRepositoryInterface のテスト用モック実装。
// ============================================================

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Search(query string) ([]model.User, error) {
	args := m.Called(query)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) FindByGitHubID(githubID int64) (*model.User, error) {
	args := m.Called(githubID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteWithRelatedData(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(userID uint, hashedPassword string) error {
	args := m.Called(userID, hashedPassword)
	return args.Error(0)
}

// ============================================================
// MockQuestionRepository は repository.QuestionRepositoryInterface のテスト用モック実装。
// ============================================================

type MockQuestionRepository struct {
	mock.Mock
}

func (m *MockQuestionRepository) Create(question *model.Question) error {
	args := m.Called(question)
	return args.Error(0)
}

func (m *MockQuestionRepository) FindByID(id uint) (*model.Question, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Question), args.Error(1)
}

func (m *MockQuestionRepository) FindAll(limit, offset int, tag string, sort string) ([]model.Question, int64, error) {
	args := m.Called(limit, offset, tag, sort)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}

func (m *MockQuestionRepository) Search(q string, limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(q, limit, offset)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}

func (m *MockQuestionRepository) FindByUserID(userID uint, limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}

func (m *MockQuestionRepository) Update(question *model.Question) error {
	args := m.Called(question)
	return args.Error(0)
}

func (m *MockQuestionRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockQuestionRepository) Vote(userID, questionID uint, value int) error {
	args := m.Called(userID, questionID, value)
	return args.Error(0)
}

func (m *MockQuestionRepository) RemoveVote(userID, questionID uint) error {
	args := m.Called(userID, questionID)
	return args.Error(0)
}

func (m *MockQuestionRepository) GetUserVote(userID, questionID uint) (int, error) {
	args := m.Called(userID, questionID)
	return args.Int(0), args.Error(1)
}

func (m *MockQuestionRepository) FindSolved(limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}

func (m *MockQuestionRepository) Bookmark(userID, questionID uint) error {
	args := m.Called(userID, questionID)
	return args.Error(0)
}

func (m *MockQuestionRepository) Unbookmark(userID, questionID uint) error {
	args := m.Called(userID, questionID)
	return args.Error(0)
}

func (m *MockQuestionRepository) HasBookmarked(userID, questionID uint) (bool, error) {
	args := m.Called(userID, questionID)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuestionRepository) FindBookmarkedByUserID(userID uint, limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}

func (m *MockQuestionRepository) FindUnanswered(limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}

func (m *MockQuestionRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// ============================================================
// MockAnswerRepository は repository.AnswerRepositoryInterface のテスト用モック実装。
// ============================================================

type MockAnswerRepository struct {
	mock.Mock
}

func (m *MockAnswerRepository) Create(answer *model.Answer) error {
	args := m.Called(answer)
	return args.Error(0)
}

func (m *MockAnswerRepository) FindByQuestionID(questionID uint) ([]model.Answer, error) {
	args := m.Called(questionID)
	return args.Get(0).([]model.Answer), args.Error(1)
}

func (m *MockAnswerRepository) FindByID(id uint) (*model.Answer, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Answer), args.Error(1)
}

func (m *MockAnswerRepository) Update(answer *model.Answer) error {
	args := m.Called(answer)
	return args.Error(0)
}

func (m *MockAnswerRepository) Delete(answer *model.Answer) error {
	args := m.Called(answer)
	return args.Error(0)
}

func (m *MockAnswerRepository) SetBestAnswer(questionID, answerID uint) error {
	args := m.Called(questionID, answerID)
	return args.Error(0)
}

func (m *MockAnswerRepository) Vote(userID, answerID uint, value int) error {
	args := m.Called(userID, answerID, value)
	return args.Error(0)
}

func (m *MockAnswerRepository) RemoveVote(userID, answerID uint) error {
	args := m.Called(userID, answerID)
	return args.Error(0)
}

func (m *MockAnswerRepository) GetUserVotes(userID uint, answerIDs []uint) (map[uint]int, error) {
	args := m.Called(userID, answerIDs)
	return args.Get(0).(map[uint]int), args.Error(1)
}

func (m *MockAnswerRepository) FindByVoteRange(questionID uint, minVote, maxVote int) ([]model.Answer, error) {
	args := m.Called(questionID, minVote, maxVote)
	return args.Get(0).([]model.Answer), args.Error(1)
}

// (MockLearningLogRepository は DIP 移行に伴い撤去)

// (MockLearningGoalRepository は DIP 移行に伴い撤去)

// ============================================================
// MockBookReviewRepository は repository.BookReviewRepositoryInterface のテスト用モック実装。
// ============================================================

type MockBookReviewRepository struct {
	mock.Mock
}

func (m *MockBookReviewRepository) Create(review *model.BookReview) error {
	args := m.Called(review)
	return args.Error(0)
}

func (m *MockBookReviewRepository) FindByID(id uint) (*model.BookReview, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.BookReview), args.Error(1)
}

func (m *MockBookReviewRepository) FindByUserID(userID uint, limit, offset int) ([]model.BookReview, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.BookReview), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookReviewRepository) FindAll(limit, offset int) ([]model.BookReview, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]model.BookReview), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookReviewRepository) Update(review *model.BookReview) error {
	args := m.Called(review)
	return args.Error(0)
}

func (m *MockBookReviewRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBookReviewRepository) FindByRating(userID uint, minRating, maxRating int) ([]model.BookReview, error) {
	args := m.Called(userID, minRating, maxRating)
	return args.Get(0).([]model.BookReview), args.Error(1)
}

func (m *MockBookReviewRepository) Search(query string, limit, offset int) ([]model.BookReview, int64, error) {
	args := m.Called(query, limit, offset)
	return args.Get(0).([]model.BookReview), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookReviewRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// (MockLearningResourceRepository は DIP 移行に伴い撤去)

// ============================================================
// MockPasswordResetRepository は repository.PasswordResetRepositoryInterface のテスト用モック実装。
// ============================================================

type MockPasswordResetRepository struct {
	mock.Mock
}

func (m *MockPasswordResetRepository) Create(token *model.PasswordResetToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) FindByToken(token string) (*model.PasswordResetToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PasswordResetToken), args.Error(1)
}

func (m *MockPasswordResetRepository) MarkAsUsed(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) InvalidateUserTokens(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) DeleteExpired() error {
	args := m.Called()
	return args.Error(0)
}

// (MockRoadmapRepository は DIP 移行に伴い撤去)

// (MockAIAdviceRepository は DIP 移行に伴い撤去)

// (MockAIConversationRepository は DIP 移行に伴い撤去)

// ============================================================
// MockWeeklyActivityReportReader は usecase/repository.WeeklyActivityReportReader のテスト用モック実装。
// ============================================================

type MockWeeklyActivityReportReader struct {
	mock.Mock
}

func (m *MockWeeklyActivityReportReader) GetWeeklyReport(ctx context.Context, userID uint) (*model.ActivityReport, error) {
	args := m.Called(ctx, userID)
	r, _ := args.Get(0).(*model.ActivityReport)
	return r, args.Error(1)
}

// インターフェース適合チェック
var _ usecaserepo.WeeklyActivityReportReader = (*MockWeeklyActivityReportReader)(nil)

// ============================================================
// MockCronScheduler は CronScheduler のテスト用モック実装。
// ============================================================

type MockCronScheduler struct {
	mock.Mock
}

var _ CronScheduler = (*MockCronScheduler)(nil)

func (m *MockCronScheduler) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	args := m.Called(spec, cmd)
	return args.Get(0).(cron.EntryID), args.Error(1)
}

func (m *MockCronScheduler) Start() {
	m.Called()
}

func (m *MockCronScheduler) Stop() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

// ============================================================
// MockWeeklyReportSender は WeeklyReportSender のテスト用モック実装。
// ============================================================

type MockWeeklyReportSender struct {
	mock.Mock
}

var _ WeeklyReportSender = (*MockWeeklyReportSender)(nil)

func (m *MockWeeklyReportSender) Execute(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

// (MockSpotifyRepository は spotify の DIP 移行に伴い撤去)

// (MockProjectMilestoneRepository は project_milestone の DIP 移行に伴い撤去)

// (MockUserActivityRepository は user_activity の DIP 移行に伴い撤去)

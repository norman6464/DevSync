// mock_test.go はService層テスト用のモックリポジトリを定義する。
package service

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
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
// MockPostRepository は repository.PostRepositoryInterface のテスト用モック実装。
// ============================================================

type MockPostRepository struct {
	mock.Mock
}

func (m *MockPostRepository) Create(post *model.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockPostRepository) FindByID(id uint) (*model.Post, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Post), args.Error(1)
}

func (m *MockPostRepository) FindAll(page, limit int) ([]model.Post, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]model.Post), args.Error(1)
}

func (m *MockPostRepository) CountAll() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPostRepository) FindByUserID(userID uint, limit, offset int) ([]model.Post, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.Post), args.Get(1).(int64), args.Error(2)
}

func (m *MockPostRepository) FindDraftsByUserID(userID uint) ([]model.Post, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Post), args.Error(1)
}

func (m *MockPostRepository) Timeline(userID uint, page, limit int) ([]model.Post, error) {
	args := m.Called(userID, page, limit)
	return args.Get(0).([]model.Post), args.Error(1)
}

func (m *MockPostRepository) Update(post *model.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockPostRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPostRepository) Like(userID, postID uint) error {
	args := m.Called(userID, postID)
	return args.Error(0)
}

func (m *MockPostRepository) Unlike(userID, postID uint) error {
	args := m.Called(userID, postID)
	return args.Error(0)
}

func (m *MockPostRepository) HasLiked(userID, postID uint) bool {
	args := m.Called(userID, postID)
	return args.Bool(0)
}

func (m *MockPostRepository) Bookmark(userID, postID uint) error {
	args := m.Called(userID, postID)
	return args.Error(0)
}

func (m *MockPostRepository) Unbookmark(userID, postID uint) error {
	args := m.Called(userID, postID)
	return args.Error(0)
}

func (m *MockPostRepository) HasBookmarked(userID, postID uint) bool {
	args := m.Called(userID, postID)
	return args.Bool(0)
}

func (m *MockPostRepository) FindBookmarkedByUserID(userID uint, page, limit int) ([]model.Post, int64, error) {
	args := m.Called(userID, page, limit)
	return args.Get(0).([]model.Post), args.Get(1).(int64), args.Error(2)
}
func (m *MockPostRepository) CountBookmarkedByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPostRepository) AddReaction(userID, postID uint, emoji string) error {
	return m.Called(userID, postID, emoji).Error(0)
}

func (m *MockPostRepository) RemoveReaction(userID, postID uint, emoji string) error {
	return m.Called(userID, postID, emoji).Error(0)
}

func (m *MockPostRepository) GetReactionsByPostID(postID uint) ([]model.ReactionCount, error) {
	args := m.Called(postID)
	return args.Get(0).([]model.ReactionCount), args.Error(1)
}

func (m *MockPostRepository) GetUserReactions(userID, postID uint) ([]string, error) {
	args := m.Called(userID, postID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockPostRepository) GetReactionsBatch(postIDs []uint) (map[uint][]model.ReactionCount, error) {
	args := m.Called(postIDs)
	return args.Get(0).(map[uint][]model.ReactionCount), args.Error(1)
}

func (m *MockPostRepository) GetUserReactionsBatch(userID uint, postIDs []uint) (map[uint][]string, error) {
	args := m.Called(userID, postIDs)
	return args.Get(0).(map[uint][]string), args.Error(1)
}

func (m *MockPostRepository) FindScheduledByUserID(userID uint) ([]model.Post, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Post), args.Error(1)
}
func (m *MockPostRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockPostRepository) CountDraftsByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockPostRepository) CountScheduledByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// ============================================================
// MockNotificationRepository は repository.NotificationRepositoryInterface のテスト用モック実装。
// ============================================================

type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) Create(notification *model.Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}

func (m *MockNotificationRepository) CreateBatch(notifications []*model.Notification) error {
	args := m.Called(notifications)
	return args.Error(0)
}

func (m *MockNotificationRepository) FindByUserID(userID uint, page, limit int, notificationType string) ([]model.Notification, error) {
	args := m.Called(userID, page, limit, notificationType)
	return args.Get(0).([]model.Notification), args.Error(1)
}

func (m *MockNotificationRepository) CountByUserID(userID uint, notificationType string) (int64, error) {
	args := m.Called(userID, notificationType)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockNotificationRepository) CountUnread(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockNotificationRepository) MarkAsRead(id, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockNotificationRepository) MarkAllAsRead(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockNotificationRepository) Delete(id, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockNotificationRepository) GetFollowerIDs(userID uint) ([]uint, error) {
	args := m.Called(userID)
	return args.Get(0).([]uint), args.Error(1)
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

// ============================================================
// MockLearningLogRepository は repository.LearningLogRepositoryInterface のテスト用モック実装。
// ============================================================

type MockLearningLogRepository struct {
	mock.Mock
}

func (m *MockLearningLogRepository) Create(log *model.LearningLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockLearningLogRepository) CreateBatch(logs []model.LearningLog) error {
	args := m.Called(logs)
	return args.Error(0)
}

func (m *MockLearningLogRepository) Update(log *model.LearningLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockLearningLogRepository) Delete(id, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockLearningLogRepository) FindByID(id uint) (*model.LearningLog, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.LearningLog), args.Error(1)
}

func (m *MockLearningLogRepository) GetByUserID(userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.LearningLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockLearningLogRepository) GetByCategory(userID uint, category string) ([]model.LearningLog, error) {
	args := m.Called(userID, category)
	return args.Get(0).([]model.LearningLog), args.Error(1)
}

func (m *MockLearningLogRepository) GetBySource(userID uint, source string) ([]model.LearningLog, error) {
	args := m.Called(userID, source)
	return args.Get(0).([]model.LearningLog), args.Error(1)
}

func (m *MockLearningLogRepository) GetByPeriod(userID uint, days int) ([]model.LearningLog, error) {
	args := m.Called(userID, days)
	return args.Get(0).([]model.LearningLog), args.Error(1)
}

func (m *MockLearningLogRepository) SumDurationByPeriod(userID uint, days int) (int, error) {
	args := m.Called(userID, days)
	return args.Int(0), args.Error(1)
}

func (m *MockLearningLogRepository) GetStreakInfo(userID uint) (*model.StreakInfo, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.StreakInfo), args.Error(1)
}

func (m *MockLearningLogRepository) GetCalendarData(userID uint) ([]model.CalendarEntry, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.CalendarEntry), args.Error(1)
}

func (m *MockLearningLogRepository) GetRecentCategories(userID uint, limit int) ([]string, error) {
	args := m.Called(userID, limit)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockLearningLogRepository) GetByGoalID(goalID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	args := m.Called(goalID, limit, offset)
	return args.Get(0).([]model.LearningLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockLearningLogRepository) SumDurationByGoalID(goalID uint) (int, error) {
	args := m.Called(goalID)
	return args.Int(0), args.Error(1)
}

func (m *MockLearningLogRepository) GetFavorites(userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.LearningLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockLearningLogRepository) GetMonthlySummary(userID uint, months int) ([]model.MonthlySummary, error) {
	args := m.Called(userID, months)
	return args.Get(0).([]model.MonthlySummary), args.Error(1)
}

func (m *MockLearningLogRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// ============================================================
// MockLearningGoalRepository は repository.LearningGoalRepositoryInterface のテスト用モック実装。
// ============================================================

type MockLearningGoalRepository struct {
	mock.Mock
}

func (m *MockLearningGoalRepository) Create(goal *model.LearningGoal) error {
	args := m.Called(goal)
	return args.Error(0)
}

func (m *MockLearningGoalRepository) Update(goal *model.LearningGoal) error {
	args := m.Called(goal)
	return args.Error(0)
}

func (m *MockLearningGoalRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockLearningGoalRepository) FindByID(id uint) (*model.LearningGoal, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.LearningGoal), args.Error(1)
}

func (m *MockLearningGoalRepository) GetByUserID(userID uint, limit, offset int) ([]model.LearningGoal, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.LearningGoal), args.Get(1).(int64), args.Error(2)
}

func (m *MockLearningGoalRepository) GetActiveByUserID(userID uint) ([]model.LearningGoal, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.LearningGoal), args.Error(1)
}

func (m *MockLearningGoalRepository) GetByCategory(userID uint, category string) ([]model.LearningGoal, error) {
	args := m.Called(userID, category)
	return args.Get(0).([]model.LearningGoal), args.Error(1)
}

func (m *MockLearningGoalRepository) GetByStatus(userID uint, status string) ([]model.LearningGoal, error) {
	args := m.Called(userID, status)
	return args.Get(0).([]model.LearningGoal), args.Error(1)
}

func (m *MockLearningGoalRepository) GetPublicByUserID(userID uint, limit, offset int) ([]model.LearningGoal, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.LearningGoal), args.Get(1).(int64), args.Error(2)
}

func (m *MockLearningGoalRepository) GetPublicGoals(limit, offset int) ([]model.LearningGoal, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]model.LearningGoal), args.Get(1).(int64), args.Error(2)
}

func (m *MockLearningGoalRepository) GetStats(userID uint) (*model.LearningGoalStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.LearningGoalStats), args.Error(1)
}

func (m *MockLearningGoalRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

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

// ============================================================
// MockLearningResourceRepository は repository.LearningResourceRepositoryInterface のテスト用モック実装。
// ============================================================

type MockLearningResourceRepository struct {
	mock.Mock
}

func (m *MockLearningResourceRepository) Create(resource *model.LearningResource) error {
	args := m.Called(resource)
	return args.Error(0)
}

func (m *MockLearningResourceRepository) FindByID(id uint) (*model.LearningResource, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.LearningResource), args.Error(1)
}

func (m *MockLearningResourceRepository) FindByUserID(userID uint, includePrivate bool, limit, offset int) ([]model.LearningResource, int64, error) {
	args := m.Called(userID, includePrivate, limit, offset)
	return args.Get(0).([]model.LearningResource), args.Get(1).(int64), args.Error(2)
}

func (m *MockLearningResourceRepository) FindPublic(limit, offset int, category string, difficulty string) ([]model.LearningResource, int64, error) {
	args := m.Called(limit, offset, category, difficulty)
	return args.Get(0).([]model.LearningResource), args.Get(1).(int64), args.Error(2)
}

func (m *MockLearningResourceRepository) Update(resource *model.LearningResource) error {
	args := m.Called(resource)
	return args.Error(0)
}

func (m *MockLearningResourceRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockLearningResourceRepository) Search(query string, limit, offset int) ([]model.LearningResource, int64, error) {
	args := m.Called(query, limit, offset)
	return args.Get(0).([]model.LearningResource), args.Get(1).(int64), args.Error(2)
}

func (m *MockLearningResourceRepository) Like(userID, resourceID uint) error {
	args := m.Called(userID, resourceID)
	return args.Error(0)
}

func (m *MockLearningResourceRepository) Unlike(userID, resourceID uint) error {
	args := m.Called(userID, resourceID)
	return args.Error(0)
}

func (m *MockLearningResourceRepository) HasLiked(userID, resourceID uint) (bool, error) {
	args := m.Called(userID, resourceID)
	return args.Bool(0), args.Error(1)
}

func (m *MockLearningResourceRepository) Save(userID, resourceID uint) error {
	args := m.Called(userID, resourceID)
	return args.Error(0)
}

func (m *MockLearningResourceRepository) Unsave(userID, resourceID uint) error {
	args := m.Called(userID, resourceID)
	return args.Error(0)
}

func (m *MockLearningResourceRepository) HasSaved(userID, resourceID uint) (bool, error) {
	args := m.Called(userID, resourceID)
	return args.Bool(0), args.Error(1)
}

func (m *MockLearningResourceRepository) FindSavedByUserID(userID uint, limit, offset int) ([]model.LearningResource, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.LearningResource), args.Get(1).(int64), args.Error(2)
}

func (m *MockLearningResourceRepository) FindByDifficulty(difficulty string, limit, offset int) ([]model.LearningResource, int64, error) {
	args := m.Called(difficulty, limit, offset)
	return args.Get(0).([]model.LearningResource), args.Get(1).(int64), args.Error(2)
}

func (m *MockLearningResourceRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

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

// ============================================================
// MockRoadmapRepository は repository.RoadmapRepositoryInterface のテスト用モック実装。
// ============================================================

type MockRoadmapRepository struct {
	mock.Mock
}

func (m *MockRoadmapRepository) Create(roadmap *model.Roadmap) error {
	args := m.Called(roadmap)
	return args.Error(0)
}

func (m *MockRoadmapRepository) Update(roadmap *model.Roadmap) error {
	args := m.Called(roadmap)
	return args.Error(0)
}

func (m *MockRoadmapRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoadmapRepository) FindByID(id uint) (*model.Roadmap, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Roadmap), args.Error(1)
}

func (m *MockRoadmapRepository) GetByUserID(userID uint, limit, offset int) ([]model.Roadmap, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.Roadmap), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoadmapRepository) GetByStatus(userID uint, status string) ([]model.Roadmap, error) {
	args := m.Called(userID, status)
	return args.Get(0).([]model.Roadmap), args.Error(1)
}

func (m *MockRoadmapRepository) GetPublicRoadmaps(limit, offset int) ([]model.Roadmap, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]model.Roadmap), args.Get(1).(int64), args.Error(2)
}

func (m *MockRoadmapRepository) CopyRoadmap(originalID, newUserID uint) (*model.Roadmap, error) {
	args := m.Called(originalID, newUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Roadmap), args.Error(1)
}

func (m *MockRoadmapRepository) GetStats(userID uint) (*model.RoadmapStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RoadmapStats), args.Error(1)
}

func (m *MockRoadmapRepository) CreateStep(step *model.RoadmapStep) error {
	args := m.Called(step)
	return args.Error(0)
}

func (m *MockRoadmapRepository) UpdateStep(step *model.RoadmapStep) error {
	args := m.Called(step)
	return args.Error(0)
}

func (m *MockRoadmapRepository) DeleteStep(stepID uint) error {
	args := m.Called(stepID)
	return args.Error(0)
}

func (m *MockRoadmapRepository) FindStepByID(stepID uint) (*model.RoadmapStep, error) {
	args := m.Called(stepID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RoadmapStep), args.Error(1)
}

func (m *MockRoadmapRepository) ReorderSteps(roadmapID uint, stepOrders []model.StepOrder) error {
	args := m.Called(roadmapID, stepOrders)
	return args.Error(0)
}

func (m *MockRoadmapRepository) GetTemplates() ([]model.Roadmap, error) {
	args := m.Called()
	return args.Get(0).([]model.Roadmap), args.Error(1)
}

func (m *MockRoadmapRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// (MockCodeSnippetRepository は code_snippet の DIP 移行に伴い撤去)

var _ repository.AIConversationRepositoryInterface = (*MockAIConversationRepository)(nil)
var _ repository.GitHubRepositoryInterface = (*MockGitHubRepository)(nil)
var _ LLMClientInterface = (*MockLLMClient)(nil)

// ============================================================
// MockAIAdviceRepository は repository.AIAdviceRepositoryInterface のテスト用モック実装。
// ============================================================

type MockAIAdviceRepository struct {
	mock.Mock
}

func (m *MockAIAdviceRepository) Create(advice *model.AIAdvice) error {
	args := m.Called(advice)
	return args.Error(0)
}

func (m *MockAIAdviceRepository) CreateBatch(advices []*model.AIAdvice) error {
	args := m.Called(advices)
	return args.Error(0)
}

func (m *MockAIAdviceRepository) FindByUserID(userID uint, limit int) ([]model.AIAdvice, error) {
	args := m.Called(userID, limit)
	return args.Get(0).([]model.AIAdvice), args.Error(1)
}

func (m *MockAIAdviceRepository) FindUnreadByUserID(userID uint) ([]model.AIAdvice, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.AIAdvice), args.Error(1)
}

func (m *MockAIAdviceRepository) MarkAsRead(id, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockAIAdviceRepository) MarkAllAsRead(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockAIAdviceRepository) DeleteExpired() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockAIAdviceRepository) DeleteByUserID(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

// ============================================================
// MockAIConversationRepository は repository.AIConversationRepositoryInterface のテスト用モック実装。
// ============================================================

type MockAIConversationRepository struct {
	mock.Mock
}

func (m *MockAIConversationRepository) CreateConversation(conv *model.AIConversation) error {
	args := m.Called(conv)
	return args.Error(0)
}

func (m *MockAIConversationRepository) FindConversationsByUserID(userID uint, limit, offset int) ([]model.AIConversation, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.AIConversation), args.Error(1)
}

func (m *MockAIConversationRepository) FindConversationByID(id uint) (*model.AIConversation, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AIConversation), args.Error(1)
}

func (m *MockAIConversationRepository) AddMessage(msg *model.AIMessage) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockAIConversationRepository) GetMessages(conversationID uint) ([]model.AIMessage, error) {
	args := m.Called(conversationID)
	return args.Get(0).([]model.AIMessage), args.Error(1)
}

func (m *MockAIConversationRepository) CountTodayMessages(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAIConversationRepository) DeleteConversation(id, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

// ============================================================
// MockLLMClient は LLMClientInterface のテスト用モック実装。
// ============================================================

type MockLLMClient struct {
	mock.Mock
}

func (m *MockLLMClient) Complete(messages []ChatMessage) (*ChatResponse, error) {
	args := m.Called(messages)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ChatResponse), args.Error(1)
}

// ============================================================
// MockGitHubRepository は repository.GitHubRepositoryInterface のテスト用モック実装。
// ============================================================

type MockGitHubRepository struct {
	mock.Mock
}

func (m *MockGitHubRepository) UpsertContributions(contributions []model.GitHubContribution) error {
	args := m.Called(contributions)
	return args.Error(0)
}

func (m *MockGitHubRepository) GetContributions(userID uint) ([]model.GitHubContribution, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.GitHubContribution), args.Error(1)
}

func (m *MockGitHubRepository) UpsertLanguageStats(stats []model.GitHubLanguageStat) error {
	args := m.Called(stats)
	return args.Error(0)
}

func (m *MockGitHubRepository) GetLanguageStats(userID uint) ([]model.GitHubLanguageStat, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.GitHubLanguageStat), args.Error(1)
}

func (m *MockGitHubRepository) UpsertRepos(repos []model.GitHubRepository) error {
	args := m.Called(repos)
	return args.Error(0)
}

func (m *MockGitHubRepository) GetRepos(userID uint) ([]model.GitHubRepository, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.GitHubRepository), args.Error(1)
}

func (m *MockGitHubRepository) DeleteUserData(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

// ============================================================
// MockEmailSender は EmailSenderInterface のテスト用モック実装。
// ============================================================

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) Send(to, subject, htmlBody string) error {
	args := m.Called(to, subject, htmlBody)
	return args.Error(0)
}

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
var _ EmailSenderInterface = (*MockEmailSender)(nil)
var _ usecaserepo.WeeklyActivityReportReader = (*MockWeeklyActivityReportReader)(nil)
var _ NotificationServiceInterface = (*MockNotificationService)(nil)

// ============================================================
// MockNotificationService は NotificationServiceInterface のテスト用モック実装。
// ============================================================

type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) CreateNotification(notification *model.Notification) error {
	args := m.Called(notification)
	return args.Error(0)
}

func (m *MockNotificationService) NotifyFollowers(actorID uint, postID uint, notificationType model.NotificationType) {
	m.Called(actorID, postID, notificationType)
}

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

func (m *MockWeeklyReportSender) SendAllWeeklyReports() error {
	args := m.Called()
	return args.Error(0)
}

// ============================================================
// MockSpotifyRepository は repository.SpotifyRepositoryInterface のテスト用モック実装。
// ============================================================

type MockSpotifyRepository struct {
	mock.Mock
}

var _ repository.SpotifyRepositoryInterface = (*MockSpotifyRepository)(nil)

func (m *MockSpotifyRepository) DeleteUserData(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

// (MockProjectMilestoneRepository は project_milestone の DIP 移行に伴い撤去)

// (MockUserActivityRepository は user_activity の DIP 移行に伴い撤去)

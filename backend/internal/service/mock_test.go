// mock_test.go はService層テスト用のモックリポジトリを定義する。
package service

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
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

func (m *MockPostRepository) FindByUserID(userID uint) ([]model.Post, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Post), args.Error(1)
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

func (m *MockPostRepository) CreateComment(comment *model.Comment) error {
	args := m.Called(comment)
	return args.Error(0)
}

func (m *MockPostRepository) FindCommentByID(id uint) (*model.Comment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Comment), args.Error(1)
}

func (m *MockPostRepository) GetComments(postID uint) ([]model.Comment, error) {
	args := m.Called(postID)
	return args.Get(0).([]model.Comment), args.Error(1)
}

func (m *MockPostRepository) GetReplies(parentID uint) ([]model.Comment, error) {
	args := m.Called(parentID)
	return args.Get(0).([]model.Comment), args.Error(1)
}

func (m *MockPostRepository) DeleteComment(id uint) error {
	args := m.Called(id)
	return args.Error(0)
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

// ============================================================
// MockFollowRepository は repository.FollowRepositoryInterface のテスト用モック実装。
// ============================================================

type MockFollowRepository struct {
	mock.Mock
}

func (m *MockFollowRepository) Follow(followerID, followeeID uint) error {
	args := m.Called(followerID, followeeID)
	return args.Error(0)
}

func (m *MockFollowRepository) Unfollow(followerID, followeeID uint) error {
	args := m.Called(followerID, followeeID)
	return args.Error(0)
}

func (m *MockFollowRepository) IsFollowing(followerID, followeeID uint) bool {
	args := m.Called(followerID, followeeID)
	return args.Bool(0)
}

func (m *MockFollowRepository) GetFollowers(userID uint) ([]model.User, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockFollowRepository) GetFollowing(userID uint) ([]model.User, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.User), args.Error(1)
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
// MockMessageRepository は repository.MessageRepositoryInterface のテスト用モック実装。
// ============================================================

type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Create(msg *model.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockMessageRepository) GetConversation(userID, otherUserID uint, page, limit int) ([]model.Message, error) {
	args := m.Called(userID, otherUserID, page, limit)
	return args.Get(0).([]model.Message), args.Error(1)
}

func (m *MockMessageRepository) GetConversations(userID uint) ([]model.ConversationSummary, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.ConversationSummary), args.Error(1)
}

func (m *MockMessageRepository) MarkAsRead(senderID, receiverID uint) error {
	args := m.Called(senderID, receiverID)
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

func (m *MockLearningGoalRepository) GetStats(userID uint) (*model.LearningGoalStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.LearningGoalStats), args.Error(1)
}

// ============================================================
// MockProjectRepository は repository.ProjectRepositoryInterface のテスト用モック実装。
// ============================================================

type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(project *model.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectRepository) FindByID(id uint) (*model.Project, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Project), args.Error(1)
}

func (m *MockProjectRepository) FindByUserID(userID uint, limit, offset int) ([]model.Project, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.Project), args.Get(1).(int64), args.Error(2)
}

func (m *MockProjectRepository) FindFeaturedByUserID(userID uint) ([]model.Project, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Project), args.Error(1)
}

func (m *MockProjectRepository) Update(project *model.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProjectRepository) FindAll(limit, offset int) ([]model.Project, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]model.Project), args.Get(1).(int64), args.Error(2)
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

// ============================================================
// MockChatRoomRepository は repository.ChatRoomRepositoryInterface のテスト用モック実装。
// ============================================================

type MockChatRoomRepository struct {
	mock.Mock
}

func (m *MockChatRoomRepository) Create(room *model.ChatRoom) error {
	args := m.Called(room)
	return args.Error(0)
}

func (m *MockChatRoomRepository) FindByID(id uint) (*model.ChatRoom, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ChatRoom), args.Error(1)
}

func (m *MockChatRoomRepository) FindByUserID(userID uint) ([]model.ChatRoom, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.ChatRoom), args.Error(1)
}

func (m *MockChatRoomRepository) Update(room *model.ChatRoom) error {
	args := m.Called(room)
	return args.Error(0)
}

func (m *MockChatRoomRepository) Delete(roomID uint) error {
	args := m.Called(roomID)
	return args.Error(0)
}

func (m *MockChatRoomRepository) AddMember(roomID, userID uint) error {
	args := m.Called(roomID, userID)
	return args.Error(0)
}

func (m *MockChatRoomRepository) RemoveMember(roomID, userID uint) error {
	args := m.Called(roomID, userID)
	return args.Error(0)
}

func (m *MockChatRoomRepository) GetMembers(roomID uint) ([]model.ChatRoomMember, error) {
	args := m.Called(roomID)
	return args.Get(0).([]model.ChatRoomMember), args.Error(1)
}

func (m *MockChatRoomRepository) IsMember(roomID, userID uint) (bool, error) {
	args := m.Called(roomID, userID)
	return args.Bool(0), args.Error(1)
}

// ============================================================
// MockGroupMessageRepository は repository.GroupMessageRepositoryInterface のテスト用モック実装。
// ============================================================

type MockGroupMessageRepository struct {
	mock.Mock
}

func (m *MockGroupMessageRepository) Create(msg *model.GroupMessage) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockGroupMessageRepository) FindByRoomID(roomID uint, page, limit int) ([]model.GroupMessage, error) {
	args := m.Called(roomID, page, limit)
	return args.Get(0).([]model.GroupMessage), args.Error(1)
}

func (m *MockGroupMessageRepository) FindSenderByID(msg *model.GroupMessage) {
	m.Called(msg)
}

func (m *MockGroupMessageRepository) GetMemberUserIDs(roomID uint) []uint {
	args := m.Called(roomID)
	return args.Get(0).([]uint)
}

// ============================================================
// MockCodeSnippetRepository は repository.CodeSnippetRepositoryInterface のテスト用モック実装。
// ============================================================

type MockCodeSnippetRepository struct {
	mock.Mock
}

func (m *MockCodeSnippetRepository) Create(snippet *model.CodeSnippet) error {
	args := m.Called(snippet)
	return args.Error(0)
}

func (m *MockCodeSnippetRepository) FindByID(id uint) (*model.CodeSnippet, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CodeSnippet), args.Error(1)
}

func (m *MockCodeSnippetRepository) FindByPostID(postID uint) ([]model.CodeSnippet, error) {
	args := m.Called(postID)
	return args.Get(0).([]model.CodeSnippet), args.Error(1)
}

func (m *MockCodeSnippetRepository) FindByUserIDAndLanguage(userID uint, language string) ([]model.CodeSnippet, error) {
	args := m.Called(userID, language)
	return args.Get(0).([]model.CodeSnippet), args.Error(1)
}

func (m *MockCodeSnippetRepository) Update(snippet *model.CodeSnippet) error {
	args := m.Called(snippet)
	return args.Error(0)
}

func (m *MockCodeSnippetRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCodeSnippetRepository) CreateComment(comment *model.SnippetComment) error {
	args := m.Called(comment)
	return args.Error(0)
}

func (m *MockCodeSnippetRepository) GetComments(snippetID uint) ([]model.SnippetComment, error) {
	args := m.Called(snippetID)
	return args.Get(0).([]model.SnippetComment), args.Error(1)
}

func (m *MockCodeSnippetRepository) DeleteComment(id, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

// ============================================================
// インターフェース適合チェック（コンパイル時検証）
// ============================================================

var _ repository.UserRepositoryInterface = (*MockUserRepository)(nil)
var _ repository.PostRepositoryInterface = (*MockPostRepository)(nil)
var _ repository.FollowRepositoryInterface = (*MockFollowRepository)(nil)
var _ repository.NotificationRepositoryInterface = (*MockNotificationRepository)(nil)
var _ repository.MessageRepositoryInterface = (*MockMessageRepository)(nil)
var _ repository.QuestionRepositoryInterface = (*MockQuestionRepository)(nil)
var _ repository.AnswerRepositoryInterface = (*MockAnswerRepository)(nil)
var _ repository.LearningLogRepositoryInterface = (*MockLearningLogRepository)(nil)
var _ repository.LearningGoalRepositoryInterface = (*MockLearningGoalRepository)(nil)
var _ repository.ProjectRepositoryInterface = (*MockProjectRepository)(nil)
var _ repository.BookReviewRepositoryInterface = (*MockBookReviewRepository)(nil)
var _ repository.LearningResourceRepositoryInterface = (*MockLearningResourceRepository)(nil)
var _ repository.PasswordResetRepositoryInterface = (*MockPasswordResetRepository)(nil)
var _ repository.RoadmapRepositoryInterface = (*MockRoadmapRepository)(nil)
var _ repository.ChatRoomRepositoryInterface = (*MockChatRoomRepository)(nil)
var _ repository.GroupMessageRepositoryInterface = (*MockGroupMessageRepository)(nil)
var _ repository.CodeSnippetRepositoryInterface = (*MockCodeSnippetRepository)(nil)
var _ repository.RecommendationRepositoryInterface = (*MockRecommendationRepository)(nil)
var _ repository.AIAdviceRepositoryInterface = (*MockAIAdviceRepository)(nil)

// ============================================================
// MockRecommendationRepository は repository.RecommendationRepositoryInterface のテスト用モック実装。
// ============================================================

type MockRecommendationRepository struct {
	mock.Mock
}

func (m *MockRecommendationRepository) GetRecommendedUsers(userID uint, skills []string, limit int) ([]model.RecommendedUser, error) {
	args := m.Called(userID, skills, limit)
	return args.Get(0).([]model.RecommendedUser), args.Error(1)
}

func (m *MockRecommendationRepository) GetTrendingPosts(limit int, days int) ([]model.Post, error) {
	args := m.Called(limit, days)
	return args.Get(0).([]model.Post), args.Error(1)
}

func (m *MockRecommendationRepository) GetTrendingResources(limit int, days int) ([]model.LearningResource, error) {
	args := m.Called(limit, days)
	return args.Get(0).([]model.LearningResource), args.Error(1)
}
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
// MockQiitaRepository は repository.QiitaRepositoryInterface のテスト用モック実装。
// ============================================================

type MockQiitaRepository struct {
	mock.Mock
}

var _ repository.QiitaRepositoryInterface = (*MockQiitaRepository)(nil)

func (m *MockQiitaRepository) UpsertArticles(userID uint, articles []model.QiitaArticle) error {
	args := m.Called(userID, articles)
	return args.Error(0)
}

func (m *MockQiitaRepository) GetArticles(userID uint) ([]model.QiitaArticle, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.QiitaArticle), args.Error(1)
}

func (m *MockQiitaRepository) GetStats(userID uint) (*model.QiitaStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.QiitaStats), args.Error(1)
}

func (m *MockQiitaRepository) DeleteUserArticles(userID uint) error {
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
// MockActivityReportRepository は repository.ActivityReportRepositoryInterface のテスト用モック実装。
// ============================================================

type MockActivityReportRepository struct {
	mock.Mock
}

func (m *MockActivityReportRepository) GetWeeklyReport(userID uint) (*model.ActivityReport, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ActivityReport), args.Error(1)
}

func (m *MockActivityReportRepository) GetMonthlyReport(userID uint) (*model.ActivityReport, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ActivityReport), args.Error(1)
}

func (m *MockActivityReportRepository) GetComparison(userID uint, period model.ReportPeriod) (*model.ReportComparison, error) {
	args := m.Called(userID, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ReportComparison), args.Error(1)
}

// ============================================================
// MockLevelRepository は repository.LevelRepositoryInterface のテスト用モック実装。
// ============================================================

type MockLevelRepository struct {
	mock.Mock
}

func (m *MockLevelRepository) GetXPStats(userID uint) (*model.XPStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.XPStats), args.Error(1)
}

// ============================================================
// MockLearningAnalyticsRepository は repository.LearningAnalyticsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockLearningAnalyticsRepository struct {
	mock.Mock
}

func (m *MockLearningAnalyticsRepository) GetHeatmapData(userID uint) ([]model.HeatmapEntry, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.HeatmapEntry), args.Error(1)
}

func (m *MockLearningAnalyticsRepository) GetCategoryBreakdown(userID uint) ([]model.CategoryBreakdown, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.CategoryBreakdown), args.Error(1)
}

func (m *MockLearningAnalyticsRepository) GetWeeklyTrends(userID uint, weeks int) ([]model.WeeklyTrend, error) {
	args := m.Called(userID, weeks)
	return args.Get(0).([]model.WeeklyTrend), args.Error(1)
}

func (m *MockLearningAnalyticsRepository) GetProductivityStats(userID uint) (*model.ProductivityStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ProductivityStats), args.Error(1)
}

// ============================================================
// MockStudyCircleRepository は repository.StudyCircleRepositoryInterface のテスト用モック実装。
// ============================================================

type MockStudyCircleRepository struct {
	mock.Mock
}

func (m *MockStudyCircleRepository) Create(circle *model.StudyCircle) error {
	args := m.Called(circle)
	return args.Error(0)
}

func (m *MockStudyCircleRepository) FindByID(id uint) (*model.StudyCircle, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.StudyCircle), args.Error(1)
}

func (m *MockStudyCircleRepository) FindByUserID(userID uint) ([]model.StudyCircle, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.StudyCircle), args.Error(1)
}

func (m *MockStudyCircleRepository) Update(circle *model.StudyCircle) error {
	args := m.Called(circle)
	return args.Error(0)
}

func (m *MockStudyCircleRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockStudyCircleRepository) AddMember(circleID, userID uint, role model.StudyCircleMemberRole) error {
	args := m.Called(circleID, userID, role)
	return args.Error(0)
}

func (m *MockStudyCircleRepository) RemoveMember(circleID, userID uint) error {
	args := m.Called(circleID, userID)
	return args.Error(0)
}

func (m *MockStudyCircleRepository) GetMembers(circleID uint) ([]model.StudyCircleMember, error) {
	args := m.Called(circleID)
	return args.Get(0).([]model.StudyCircleMember), args.Error(1)
}

func (m *MockStudyCircleRepository) IsMember(circleID, userID uint) (bool, error) {
	args := m.Called(circleID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockStudyCircleRepository) GetMemberCount(circleID uint) (int, error) {
	args := m.Called(circleID)
	return args.Int(0), args.Error(1)
}

func (m *MockStudyCircleRepository) CreateStep(step *model.StudyCircleStep) error {
	args := m.Called(step)
	return args.Error(0)
}

func (m *MockStudyCircleRepository) UpdateStep(step *model.StudyCircleStep) error {
	args := m.Called(step)
	return args.Error(0)
}

func (m *MockStudyCircleRepository) DeleteStep(stepID uint) error {
	args := m.Called(stepID)
	return args.Error(0)
}

func (m *MockStudyCircleRepository) FindStepByID(stepID uint) (*model.StudyCircleStep, error) {
	args := m.Called(stepID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.StudyCircleStep), args.Error(1)
}

func (m *MockStudyCircleRepository) ReorderSteps(circleID uint, stepOrders []model.StepOrder) error {
	args := m.Called(circleID, stepOrders)
	return args.Error(0)
}

func (m *MockStudyCircleRepository) UpsertProgress(progress *model.StudyCircleMemberProgress) error {
	args := m.Called(progress)
	return args.Error(0)
}

func (m *MockStudyCircleRepository) GetProgress(circleID uint) ([]model.StudyCircleMemberProgress, error) {
	args := m.Called(circleID)
	return args.Get(0).([]model.StudyCircleMemberProgress), args.Error(1)
}

func (m *MockStudyCircleRepository) CreateCheckin(checkin *model.StudyCircleCheckin) error {
	args := m.Called(checkin)
	return args.Error(0)
}

func (m *MockStudyCircleRepository) GetCheckins(circleID uint) ([]model.StudyCircleCheckin, error) {
	args := m.Called(circleID)
	return args.Get(0).([]model.StudyCircleCheckin), args.Error(1)
}

func (m *MockStudyCircleRepository) HasCheckedInToday(circleID, userID uint) (bool, error) {
	args := m.Called(circleID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockStudyCircleRepository) GetStreakRanking(circleID uint) ([]model.CircleMemberStreak, error) {
	args := m.Called(circleID)
	return args.Get(0).([]model.CircleMemberStreak), args.Error(1)
}

func (m *MockStudyCircleRepository) Search(query string, limit, offset int) (interface{}, int64, error) {
	args := m.Called(query, limit, offset)
	return args.Get(0), args.Get(1).(int64), args.Error(2)
}

func (m *MockStudyCircleRepository) GetByStatus(userID uint, status string) ([]model.StudyCircle, error) {
	args := m.Called(userID, status)
	return args.Get(0).([]model.StudyCircle), args.Error(1)
}

// インターフェース適合チェック
var _ EmailSenderInterface = (*MockEmailSender)(nil)
var _ repository.ActivityReportRepositoryInterface = (*MockActivityReportRepository)(nil)
var _ repository.LevelRepositoryInterface = (*MockLevelRepository)(nil)
var _ repository.LearningAnalyticsRepositoryInterface = (*MockLearningAnalyticsRepository)(nil)
var _ repository.StudyCircleRepositoryInterface = (*MockStudyCircleRepository)(nil)
var _ repository.BadgeRepositoryInterface = (*MockBadgeRepository)(nil)
var _ NotificationServiceInterface = (*MockNotificationService)(nil)

// ============================================================
// MockBadgeRepository は repository.BadgeRepositoryInterface のテスト用モック実装。
// ============================================================

type MockBadgeRepository struct {
	mock.Mock
}

func (m *MockBadgeRepository) GetBadgeStats(userID uint) (*model.BadgeStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.BadgeStats), args.Error(1)
}

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
// MockRankingRepository は repository.RankingRepositoryInterface のテスト用モック実装。
// ============================================================

type MockRankingRepository struct {
	mock.Mock
}

var _ repository.RankingRepositoryInterface = (*MockRankingRepository)(nil)

func (m *MockRankingRepository) ContributionRanking(period string) ([]model.RankingEntry, error) {
	args := m.Called(period)
	return args.Get(0).([]model.RankingEntry), args.Error(1)
}

func (m *MockRankingRepository) LanguageRanking(language, period string) ([]model.RankingEntry, error) {
	args := m.Called(language, period)
	return args.Get(0).([]model.RankingEntry), args.Error(1)
}

func (m *MockRankingRepository) LevelRanking() ([]model.RankingEntry, error) {
	args := m.Called()
	return args.Get(0).([]model.RankingEntry), args.Error(1)
}

func (m *MockRankingRepository) AvailableLanguages() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

// ============================================================
// MockPostSeriesRepository
// ============================================================

type MockPostSeriesRepository struct {
	mock.Mock
}

var _ repository.PostSeriesRepositoryInterface = (*MockPostSeriesRepository)(nil)

func (m *MockPostSeriesRepository) Create(series *model.PostSeries) error {
	args := m.Called(series)
	return args.Error(0)
}

func (m *MockPostSeriesRepository) FindByID(id uint) (*model.PostSeries, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PostSeries), args.Error(1)
}

func (m *MockPostSeriesRepository) FindByUserID(userID uint, offset, limit int) ([]model.PostSeries, error) {
	args := m.Called(userID, offset, limit)
	return args.Get(0).([]model.PostSeries), args.Error(1)
}

func (m *MockPostSeriesRepository) CountByUser(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPostSeriesRepository) Update(series *model.PostSeries) error {
	args := m.Called(series)
	return args.Error(0)
}

func (m *MockPostSeriesRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPostSeriesRepository) AddPost(item *model.PostSeriesItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockPostSeriesRepository) HasPost(seriesID, postID uint) (bool, error) {
	args := m.Called(seriesID, postID)
	return args.Bool(0), args.Error(1)
}

func (m *MockPostSeriesRepository) RemovePost(seriesID, postID uint) error {
	args := m.Called(seriesID, postID)
	return args.Error(0)
}

func (m *MockPostSeriesRepository) GetPostsBySeriesID(seriesID uint) ([]model.PostSeriesItem, error) {
	args := m.Called(seriesID)
	return args.Get(0).([]model.PostSeriesItem), args.Error(1)
}

// MockPostCollectionRepository は repository.PostCollectionRepositoryInterface のテスト用モック実装。
type MockPostCollectionRepository struct {
	mock.Mock
}

var _ repository.PostCollectionRepositoryInterface = (*MockPostCollectionRepository)(nil)

func (m *MockPostCollectionRepository) Create(collection *model.PostCollection) error {
	args := m.Called(collection)
	return args.Error(0)
}

func (m *MockPostCollectionRepository) FindByID(id uint) (*model.PostCollection, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PostCollection), args.Error(1)
}

func (m *MockPostCollectionRepository) FindByUserID(userID uint) ([]model.PostCollection, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.PostCollection), args.Error(1)
}

func (m *MockPostCollectionRepository) FindPublicByUserID(userID uint) ([]model.PostCollection, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.PostCollection), args.Error(1)
}

func (m *MockPostCollectionRepository) Update(collection *model.PostCollection) error {
	args := m.Called(collection)
	return args.Error(0)
}

func (m *MockPostCollectionRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPostCollectionRepository) AddPost(item *model.PostCollectionItem) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *MockPostCollectionRepository) HasPost(collectionID, postID uint) (bool, error) {
	args := m.Called(collectionID, postID)
	return args.Bool(0), args.Error(1)
}

func (m *MockPostCollectionRepository) RemovePost(collectionID, postID uint) error {
	args := m.Called(collectionID, postID)
	return args.Error(0)
}

func (m *MockPostCollectionRepository) GetPostsByCollectionID(collectionID uint) ([]model.PostCollectionItem, error) {
	args := m.Called(collectionID)
	return args.Get(0).([]model.PostCollectionItem), args.Error(1)
}

// ============================================================
// MockPostPinRepository は repository.PostPinRepositoryInterface のテスト用モック実装。
// ============================================================

type MockPostPinRepository struct {
	mock.Mock
}

var _ repository.PostPinRepositoryInterface = (*MockPostPinRepository)(nil)

func (m *MockPostPinRepository) Pin(pin *model.PostPin) error {
	args := m.Called(pin)
	return args.Error(0)
}

func (m *MockPostPinRepository) Unpin(userID, postID uint) error {
	args := m.Called(userID, postID)
	return args.Error(0)
}

func (m *MockPostPinRepository) GetByUserID(userID uint) ([]model.PostPin, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.PostPin), args.Error(1)
}

func (m *MockPostPinRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPostPinRepository) IsPinned(userID, postID uint) (bool, error) {
	args := m.Called(userID, postID)
	return args.Bool(0), args.Error(1)
}

func (m *MockPostPinRepository) UpdateOrder(userID uint, postIDs []uint) error {
	args := m.Called(userID, postIDs)
	return args.Error(0)
}

// ============================================================
// MockPostTagRepository は repository.PostTagRepositoryInterface のテスト用モック実装。
// ============================================================

type MockPostTagRepository struct {
	mock.Mock
}

var _ repository.PostTagRepositoryInterface = (*MockPostTagRepository)(nil)

func (m *MockPostTagRepository) SetTags(postID uint, tags []string) error {
	args := m.Called(postID, tags)
	return args.Error(0)
}

func (m *MockPostTagRepository) GetByPostID(postID uint) ([]string, error) {
	args := m.Called(postID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockPostTagRepository) FindPostsByTag(tag string, limit, offset int) ([]model.Post, int64, error) {
	args := m.Called(tag, limit, offset)
	return args.Get(0).([]model.Post), args.Get(1).(int64), args.Error(2)
}

func (m *MockPostTagRepository) GetPopularTags(limit int) ([]model.TagCount, error) {
	args := m.Called(limit)
	return args.Get(0).([]model.TagCount), args.Error(1)
}

// ============================================================
// MockPostViewRepository は repository.PostViewRepositoryInterface のテスト用モック実装。
// ============================================================

type MockPostViewRepository struct {
	mock.Mock
}

var _ repository.PostViewRepositoryInterface = (*MockPostViewRepository)(nil)

func (m *MockPostViewRepository) RecordView(view *model.PostView) error {
	args := m.Called(view)
	return args.Error(0)
}

func (m *MockPostViewRepository) GetViewCount(postID uint) (int64, error) {
	args := m.Called(postID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockPostViewRepository) HasViewed(userID, postID uint) (bool, error) {
	args := m.Called(userID, postID)
	return args.Bool(0), args.Error(1)
}

func (m *MockPostViewRepository) GetMostViewed(limit int) ([]model.ViewCount, error) {
	args := m.Called(limit)
	return args.Get(0).([]model.ViewCount), args.Error(1)
}

// ============================================================
// MockMentionRepository は repository.MentionRepositoryInterface のテスト用モック実装。
// ============================================================

type MockMentionRepository struct {
	mock.Mock
}

var _ repository.MentionRepositoryInterface = (*MockMentionRepository)(nil)

func (m *MockMentionRepository) Create(mention *model.Mention) error {
	args := m.Called(mention)
	return args.Error(0)
}

func (m *MockMentionRepository) FindByUserID(userID uint, page, limit int) ([]model.Mention, error) {
	args := m.Called(userID, page, limit)
	return args.Get(0).([]model.Mention), args.Error(1)
}

func (m *MockMentionRepository) FindByPostID(postID uint) ([]model.Mention, error) {
	args := m.Called(postID)
	return args.Get(0).([]model.Mention), args.Error(1)
}

func (m *MockMentionRepository) FindByCommentID(commentID uint) ([]model.Mention, error) {
	args := m.Called(commentID)
	return args.Get(0).([]model.Mention), args.Error(1)
}

func (m *MockMentionRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMentionRepository) DeleteByPostID(postID uint) error {
	args := m.Called(postID)
	return args.Error(0)
}

func (m *MockMentionRepository) DeleteByCommentID(commentID uint) error {
	args := m.Called(commentID)
	return args.Error(0)
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

// ============================================================
// MockCommentLikeRepository は repository.CommentLikeRepositoryInterface のテスト用モック実装。
// ============================================================

type MockCommentLikeRepository struct {
	mock.Mock
}

var _ repository.CommentLikeRepositoryInterface = (*MockCommentLikeRepository)(nil)

func (m *MockCommentLikeRepository) Like(userID, commentID uint) error {
	args := m.Called(userID, commentID)
	return args.Error(0)
}

func (m *MockCommentLikeRepository) Unlike(userID, commentID uint) error {
	args := m.Called(userID, commentID)
	return args.Error(0)
}

func (m *MockCommentLikeRepository) HasLiked(userID, commentID uint) (bool, error) {
	args := m.Called(userID, commentID)
	return args.Bool(0), args.Error(1)
}

func (m *MockCommentLikeRepository) CountByCommentID(commentID uint) (int64, error) {
	args := m.Called(commentID)
	return args.Get(0).(int64), args.Error(1)
}

// ============================================================
// MockUserDashboardRepository は repository.UserDashboardRepositoryInterface のテスト用モック実装。
// ============================================================

type MockUserDashboardRepository struct {
	mock.Mock
}

func (m *MockUserDashboardRepository) GetDashboardStats(userID uint) (*model.UserDashboardStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserDashboardStats), args.Error(1)
}

// ============================================================
// MockNoteStatsRepository は repository.NoteStatsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockNoteStatsRepository struct {
	mock.Mock
}

func (m *MockNoteStatsRepository) GetNoteStats(userID uint) (*model.NoteStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.NoteStats), args.Error(1)
}

// ============================================================
// MockStudyCircleStatsRepository は repository.StudyCircleStatsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockStudyCircleStatsRepository struct {
	mock.Mock
}

func (m *MockStudyCircleStatsRepository) GetCircleStats(circleID uint) (*model.StudyCircleStats, error) {
	args := m.Called(circleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.StudyCircleStats), args.Error(1)
}

// ============================================================
// MockPostStatsRepository は repository.PostStatsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockPostStatsRepository struct {
	mock.Mock
}

func (m *MockPostStatsRepository) GetPostStats(userID uint) (*model.PostStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PostStats), args.Error(1)
}

// ============================================================
// MockBookReviewStatsRepository は repository.BookReviewStatsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockBookReviewStatsRepository struct {
	mock.Mock
}

func (m *MockBookReviewStatsRepository) GetBookReviewStats(userID uint) (*model.BookReviewStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.BookReviewStats), args.Error(1)
}

// ============================================================
// MockQAStatsRepository は repository.QAStatsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockQAStatsRepository struct {
	mock.Mock
}

func (m *MockQAStatsRepository) GetQAStats(userID uint) (*model.QAStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.QAStats), args.Error(1)
}

// MockCodeSnippetStatsRepository は repository.CodeSnippetStatsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockCodeSnippetStatsRepository struct {
	mock.Mock
}

func (m *MockCodeSnippetStatsRepository) GetCodeSnippetStats(userID uint) (*model.CodeSnippetStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CodeSnippetStats), args.Error(1)
}

// MockLearningResourceStatsRepository は repository.LearningResourceStatsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockLearningResourceStatsRepository struct {
	mock.Mock
}

func (m *MockLearningResourceStatsRepository) GetLearningResourceStats(userID uint) (*model.LearningResourceStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.LearningResourceStats), args.Error(1)
}

// MockProjectStatsRepository は repository.ProjectStatsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockProjectStatsRepository struct {
	mock.Mock
}

func (m *MockProjectStatsRepository) GetProjectStats(userID uint) (*model.ProjectStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ProjectStats), args.Error(1)
}

// MockFollowStatsRepository は repository.FollowStatsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockFollowStatsRepository struct {
	mock.Mock
}

func (m *MockFollowStatsRepository) GetFollowStats(userID uint) (*model.FollowStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.FollowStats), args.Error(1)
}

// MockRoadmapStatsRepository は repository.RoadmapStatsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockRoadmapStatsRepository struct {
	mock.Mock
}

func (m *MockRoadmapStatsRepository) GetRoadmapStats(userID uint) (*model.RoadmapStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RoadmapStats), args.Error(1)
}

// MockLearningLogStatsRepository は repository.LearningLogStatsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockLearningLogStatsRepository struct {
	mock.Mock
}

func (m *MockLearningLogStatsRepository) GetLearningLogStats(userID uint) (*model.LearningLogStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.LearningLogStats), args.Error(1)
}

// MockCommentStatsRepository は repository.CommentStatsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockCommentStatsRepository struct {
	mock.Mock
}

func (m *MockCommentStatsRepository) GetCommentStats(userID uint) (*model.CommentStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CommentStats), args.Error(1)
}

// MockNotificationStatsRepository は repository.NotificationStatsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockNotificationStatsRepository struct {
	mock.Mock
}

func (m *MockNotificationStatsRepository) GetNotificationStats(userID uint) (*model.NotificationStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.NotificationStats), args.Error(1)
}

// MockMessageStatsRepository は repository.MessageStatsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockMessageStatsRepository struct {
	mock.Mock
}

func (m *MockMessageStatsRepository) GetMessageStats(userID uint) (*model.MessageStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MessageStats), args.Error(1)
}

// MockMentionStatsRepository は repository.MentionStatsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockMentionStatsRepository struct {
	mock.Mock
}

func (m *MockMentionStatsRepository) GetMentionStats(userID uint) (*model.MentionStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.MentionStats), args.Error(1)
}

// MockReactionStatsRepository は repository.ReactionStatsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockReactionStatsRepository struct {
	mock.Mock
}

func (m *MockReactionStatsRepository) GetReactionStats(userID uint) (*model.ReactionStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.ReactionStats), args.Error(1)
}

// MockBookmarkStatsRepository は repository.BookmarkStatsRepositoryInterface のテスト用モック実装。
// ============================================================

type MockBookmarkStatsRepository struct {
	mock.Mock
}

func (m *MockBookmarkStatsRepository) GetBookmarkStats(userID uint) (*model.BookmarkStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.BookmarkStats), args.Error(1)
}

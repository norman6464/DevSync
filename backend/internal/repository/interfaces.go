package repository

import "github.com/norman6464/devsync/backend/internal/model"

// UserRepositoryInterface defines the contract for user data operations.
type UserRepositoryInterface interface {
	FindAll() ([]model.User, error)
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Search(query string) ([]model.User, error)
	FindByGitHubID(githubID int64) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uint) error
	DeleteWithRelatedData(id uint) error
	UpdatePassword(userID uint, hashedPassword string) error
}

// PostRepositoryInterface defines the contract for post data operations.
type PostRepositoryInterface interface {
	Create(post *model.Post) error
	FindByID(id uint) (*model.Post, error)
	FindAll(page, limit int) ([]model.Post, error)
	FindByUserID(userID uint) ([]model.Post, error)
	Timeline(userID uint, page, limit int) ([]model.Post, error)
	Update(post *model.Post) error
	Delete(id uint) error
	Like(userID, postID uint) error
	Unlike(userID, postID uint) error
	HasLiked(userID, postID uint) bool
	CreateComment(comment *model.Comment) error
	GetComments(postID uint) ([]model.Comment, error)
	DeleteComment(id, userID uint) error
}

// FollowRepositoryInterface defines the contract for follow data operations.
type FollowRepositoryInterface interface {
	Follow(followerID, followeeID uint) error
	Unfollow(followerID, followeeID uint) error
	IsFollowing(followerID, followeeID uint) bool
	GetFollowers(userID uint) ([]model.User, error)
	GetFollowing(userID uint) ([]model.User, error)
}

// NotificationRepositoryInterface defines the contract for notification data operations.
type NotificationRepositoryInterface interface {
	Create(notification *model.Notification) error
	CreateBatch(notifications []*model.Notification) error
	FindByUserID(userID uint, page, limit int, notificationType string) ([]model.Notification, error)
	CountByUserID(userID uint, notificationType string) (int64, error)
	CountUnread(userID uint) (int64, error)
	MarkAsRead(id, userID uint) error
	MarkAllAsRead(userID uint) error
	Delete(id, userID uint) error
	GetFollowerIDs(userID uint) ([]uint, error)
}

// MessageRepositoryInterface defines the contract for message data operations.
type MessageRepositoryInterface interface {
	Create(msg *model.Message) error
	GetConversation(userID, otherUserID uint, page, limit int) ([]model.Message, error)
	GetConversations(userID uint) ([]ConversationSummary, error)
	MarkAsRead(senderID, receiverID uint) error
}

// QuestionRepositoryInterface defines the contract for question data operations.
type QuestionRepositoryInterface interface {
	Create(question *model.Question) error
	FindByID(id uint) (*model.Question, error)
	FindAll(limit, offset int, tag string, sort string) ([]model.Question, int64, error)
	Search(q string, limit, offset int) ([]model.Question, int64, error)
	FindByUserID(userID uint) ([]model.Question, error)
	Update(question *model.Question) error
	Delete(id uint) error
	Vote(userID, questionID uint, value int) error
	RemoveVote(userID, questionID uint) error
	GetUserVote(userID, questionID uint) (int, error)
}

// AnswerRepositoryInterface defines the contract for answer data operations.
type AnswerRepositoryInterface interface {
	Create(answer *model.Answer) error
	FindByQuestionID(questionID uint) ([]model.Answer, error)
	FindByID(id uint) (*model.Answer, error)
	Update(answer *model.Answer) error
	Delete(answer *model.Answer) error
	SetBestAnswer(questionID, answerID uint) error
	Vote(userID, answerID uint, value int) error
	RemoveVote(userID, answerID uint) error
	GetUserVotes(userID uint, answerIDs []uint) (map[uint]int, error)
}

// LearningLogRepositoryInterface defines the contract for learning log data operations.
type LearningLogRepositoryInterface interface {
	Create(log *model.LearningLog) error
	Update(log *model.LearningLog) error
	Delete(id, userID uint) error
	FindByID(id uint) (*model.LearningLog, error)
	GetByUserID(userID uint) ([]model.LearningLog, error)
	GetStreakInfo(userID uint) (*model.StreakInfo, error)
	GetCalendarData(userID uint) ([]model.CalendarEntry, error)
}

// LearningGoalRepositoryInterface defines the contract for learning goal data operations.
type LearningGoalRepositoryInterface interface {
	Create(goal *model.LearningGoal) error
	Update(goal *model.LearningGoal) error
	Delete(id uint) error
	FindByID(id uint) (*model.LearningGoal, error)
	GetByUserID(userID uint) ([]model.LearningGoal, error)
	GetActiveByUserID(userID uint) ([]model.LearningGoal, error)
	GetStats(userID uint) (*model.LearningGoalStats, error)
}

// RankingRepositoryInterface defines the contract for ranking data operations.
type RankingRepositoryInterface interface {
	ContributionRanking(period string) ([]RankingEntry, error)
	LanguageRanking(language, period string) ([]RankingEntry, error)
	AvailableLanguages() ([]string, error)
}

// ProjectRepositoryInterface defines the contract for project data operations.
type ProjectRepositoryInterface interface {
	Create(project *model.Project) error
	FindByID(id uint) (*model.Project, error)
	FindByUserID(userID uint) ([]model.Project, error)
	FindFeaturedByUserID(userID uint) ([]model.Project, error)
	Update(project *model.Project) error
	Delete(id uint) error
	FindAll(limit, offset int) ([]model.Project, int64, error)
}

// BookReviewRepositoryInterface defines the contract for book review data operations.
type BookReviewRepositoryInterface interface {
	Create(review *model.BookReview) error
	FindByID(id uint) (*model.BookReview, error)
	FindByUserID(userID uint) ([]model.BookReview, error)
	FindAll(limit, offset int) ([]model.BookReview, int64, error)
	Update(review *model.BookReview) error
	Delete(id uint) error
}

// LearningResourceRepositoryInterface defines the contract for learning resource data operations.
type LearningResourceRepositoryInterface interface {
	Create(resource *model.LearningResource) error
	FindByID(id uint) (*model.LearningResource, error)
	FindByUserID(userID uint, includePrivate bool) ([]model.LearningResource, error)
	FindPublic(limit, offset int, category string, difficulty string) ([]model.LearningResource, int64, error)
	Update(resource *model.LearningResource) error
	Delete(id uint) error
	Search(query string, limit, offset int) ([]model.LearningResource, int64, error)
	Like(userID, resourceID uint) error
	Unlike(userID, resourceID uint) error
	HasLiked(userID, resourceID uint) (bool, error)
	Save(userID, resourceID uint) error
	Unsave(userID, resourceID uint) error
	HasSaved(userID, resourceID uint) (bool, error)
	FindSavedByUserID(userID uint, limit, offset int) ([]model.LearningResource, int64, error)
}

// RoadmapRepositoryInterface defines the contract for roadmap data operations.
type RoadmapRepositoryInterface interface {
	Create(roadmap *model.Roadmap) error
	Update(roadmap *model.Roadmap) error
	Delete(id uint) error
	FindByID(id uint) (*model.Roadmap, error)
	GetByUserID(userID uint) ([]model.Roadmap, error)
	GetPublicRoadmaps(limit, offset int) ([]model.Roadmap, int64, error)
	CopyRoadmap(originalID, newUserID uint) (*model.Roadmap, error)
	GetStats(userID uint) (*model.RoadmapStats, error)
	CreateStep(step *model.RoadmapStep) error
	UpdateStep(step *model.RoadmapStep) error
	DeleteStep(stepID uint) error
	FindStepByID(stepID uint) (*model.RoadmapStep, error)
	ReorderSteps(roadmapID uint, stepOrders []StepOrder) error
}

// ChatRoomRepositoryInterface defines the contract for chat room data operations.
type ChatRoomRepositoryInterface interface {
	Create(room *model.ChatRoom) error
	FindByID(id uint) (*model.ChatRoom, error)
	FindByUserID(userID uint) ([]model.ChatRoom, error)
	Update(room *model.ChatRoom) error
	Delete(roomID uint) error
	AddMember(roomID, userID uint) error
	RemoveMember(roomID, userID uint) error
	GetMembers(roomID uint) ([]model.ChatRoomMember, error)
	IsMember(roomID, userID uint) (bool, error)
}

// GroupMessageRepositoryInterface defines the contract for group message data operations.
type GroupMessageRepositoryInterface interface {
	Create(msg *model.GroupMessage) error
	FindByRoomID(roomID uint, page, limit int) ([]model.GroupMessage, error)
	FindSenderByID(msg *model.GroupMessage)
	GetMemberUserIDs(roomID uint) []uint
}

// GitHubRepositoryInterface defines the contract for GitHub data operations.
type GitHubRepositoryInterface interface {
	UpsertContributions(contributions []model.GitHubContribution) error
	GetContributions(userID uint) ([]model.GitHubContribution, error)
	UpsertLanguageStats(stats []model.GitHubLanguageStat) error
	GetLanguageStats(userID uint) ([]model.GitHubLanguageStat, error)
	UpsertRepos(repos []model.GitHubRepository) error
	GetRepos(userID uint) ([]model.GitHubRepository, error)
	DeleteUserData(userID uint) error
}

// QiitaRepositoryInterface defines the contract for Qiita data operations.
type QiitaRepositoryInterface interface {
	UpsertArticles(userID uint, articles []model.QiitaArticle) error
	GetArticles(userID uint) ([]model.QiitaArticle, error)
	GetStats(userID uint) (*model.QiitaStats, error)
	DeleteUserArticles(userID uint) error
}

// ZennRepositoryInterface defines the contract for Zenn data operations.
type ZennRepositoryInterface interface {
	UpsertArticles(userID uint, articles []model.ZennArticle) error
	GetArticles(userID uint) ([]model.ZennArticle, error)
	GetStats(userID uint) (*model.ZennStats, error)
	DeleteUserArticles(userID uint) error
}

// PasswordResetRepositoryInterface defines the contract for password reset data operations.
type PasswordResetRepositoryInterface interface {
	Create(token *model.PasswordResetToken) error
	FindByToken(token string) (*model.PasswordResetToken, error)
	MarkAsUsed(id uint) error
	InvalidateUserTokens(userID uint) error
	DeleteExpired() error
}

// ActivityReportRepositoryInterface defines the contract for activity report data operations.
type ActivityReportRepositoryInterface interface {
	GetWeeklyReport(userID uint) (*model.ActivityReport, error)
	GetMonthlyReport(userID uint) (*model.ActivityReport, error)
	GetComparison(userID uint, period model.ReportPeriod) (*model.ReportComparison, error)
}

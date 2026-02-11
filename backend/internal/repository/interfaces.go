// Package repository はDevSyncアプリケーションのデータアクセス層を提供する。
// 各リポジトリはGORMを使用してPostgreSQLに対するCRUD操作を実装する。
package repository

import "github.com/norman6464/devsync/backend/internal/model"

// UserRepositoryInterface はユーザーデータ操作の契約を定義する。
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

// PostRepositoryInterface は投稿データ操作の契約を定義する。
// いいね・コメント操作も含む。
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

// FollowRepositoryInterface はフォロー関係データ操作の契約を定義する。
type FollowRepositoryInterface interface {
	Follow(followerID, followeeID uint) error
	Unfollow(followerID, followeeID uint) error
	IsFollowing(followerID, followeeID uint) bool
	GetFollowers(userID uint) ([]model.User, error)
	GetFollowing(userID uint) ([]model.User, error)
}

// NotificationRepositoryInterface は通知データ操作の契約を定義する。
// 一括作成、未読カウント、既読マーク等を含む。
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

// MessageRepositoryInterface はDMメッセージデータ操作の契約を定義する。
type MessageRepositoryInterface interface {
	Create(msg *model.Message) error
	GetConversation(userID, otherUserID uint, page, limit int) ([]model.Message, error)
	GetConversations(userID uint) ([]ConversationSummary, error)
	MarkAsRead(senderID, receiverID uint) error
}

// QuestionRepositoryInterface はQ&A質問データ操作の契約を定義する。
// 投票操作も含む。
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

// AnswerRepositoryInterface はQ&A回答データ操作の契約を定義する。
// ベストアンサー設定や投票操作も含む。
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

// LearningLogRepositoryInterface は学習ログデータ操作の契約を定義する。
// ストリーク計算やカレンダーデータ取得も含む。
type LearningLogRepositoryInterface interface {
	Create(log *model.LearningLog) error
	Update(log *model.LearningLog) error
	Delete(id, userID uint) error
	FindByID(id uint) (*model.LearningLog, error)
	GetByUserID(userID uint) ([]model.LearningLog, error)
	GetStreakInfo(userID uint) (*model.StreakInfo, error)
	GetCalendarData(userID uint) ([]model.CalendarEntry, error)
}

// LearningGoalRepositoryInterface は学習目標データ操作の契約を定義する。
// アクティブ目標の取得や統計情報の算出も含む。
type LearningGoalRepositoryInterface interface {
	Create(goal *model.LearningGoal) error
	Update(goal *model.LearningGoal) error
	Delete(id uint) error
	FindByID(id uint) (*model.LearningGoal, error)
	GetByUserID(userID uint) ([]model.LearningGoal, error)
	GetActiveByUserID(userID uint) ([]model.LearningGoal, error)
	GetStats(userID uint) (*model.LearningGoalStats, error)
}

// RankingRepositoryInterface はランキングデータ操作の契約を定義する。
// GitHubコントリビューションと言語別ランキングを提供する。
type RankingRepositoryInterface interface {
	ContributionRanking(period string) ([]RankingEntry, error)
	LanguageRanking(language, period string) ([]RankingEntry, error)
	LevelRanking() ([]RankingEntry, error)
	AvailableLanguages() ([]string, error)
}

// ProjectRepositoryInterface はプロジェクトショーケースデータ操作の契約を定義する。
type ProjectRepositoryInterface interface {
	Create(project *model.Project) error
	FindByID(id uint) (*model.Project, error)
	FindByUserID(userID uint) ([]model.Project, error)
	FindFeaturedByUserID(userID uint) ([]model.Project, error)
	Update(project *model.Project) error
	Delete(id uint) error
	FindAll(limit, offset int) ([]model.Project, int64, error)
}

// BookReviewRepositoryInterface は書籍レビューデータ操作の契約を定義する。
type BookReviewRepositoryInterface interface {
	Create(review *model.BookReview) error
	FindByID(id uint) (*model.BookReview, error)
	FindByUserID(userID uint) ([]model.BookReview, error)
	FindAll(limit, offset int) ([]model.BookReview, int64, error)
	Update(review *model.BookReview) error
	Delete(id uint) error
}

// LearningResourceRepositoryInterface は学習リソースデータ操作の契約を定義する。
// いいね・ブックマーク（保存）操作も含む。
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

// RoadmapRepositoryInterface は学習ロードマップデータ操作の契約を定義する。
// ステップのCRUD操作やコピー機能も含む。
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
	GetTemplates() ([]model.Roadmap, error)
}

// ChatRoomRepositoryInterface はチャットルームデータ操作の契約を定義する。
// メンバー管理操作も含む。
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

// GroupMessageRepositoryInterface はグループメッセージデータ操作の契約を定義する。
type GroupMessageRepositoryInterface interface {
	Create(msg *model.GroupMessage) error
	FindByRoomID(roomID uint, page, limit int) ([]model.GroupMessage, error)
	FindSenderByID(msg *model.GroupMessage)
	GetMemberUserIDs(roomID uint) []uint
}

// CodeSnippetRepositoryInterface はコードスニペットデータ操作の契約を定義する。
// スニペットのCRUDおよびインラインコメント操作を含む。
type CodeSnippetRepositoryInterface interface {
	Create(snippet *model.CodeSnippet) error
	FindByID(id uint) (*model.CodeSnippet, error)
	FindByPostID(postID uint) ([]model.CodeSnippet, error)
	Update(snippet *model.CodeSnippet) error
	Delete(id uint) error
	CreateComment(comment *model.SnippetComment) error
	GetComments(snippetID uint) ([]model.SnippetComment, error)
	DeleteComment(id, userID uint) error
}

// GitHubRepositoryInterface はGitHub連携データ操作の契約を定義する。
// コントリビューション、言語統計、リポジトリのUpsert操作を提供する。
type GitHubRepositoryInterface interface {
	UpsertContributions(contributions []model.GitHubContribution) error
	GetContributions(userID uint) ([]model.GitHubContribution, error)
	UpsertLanguageStats(stats []model.GitHubLanguageStat) error
	GetLanguageStats(userID uint) ([]model.GitHubLanguageStat, error)
	UpsertRepos(repos []model.GitHubRepository) error
	GetRepos(userID uint) ([]model.GitHubRepository, error)
	DeleteUserData(userID uint) error
}

// QiitaRepositoryInterface はQiita連携データ操作の契約を定義する。
type QiitaRepositoryInterface interface {
	UpsertArticles(userID uint, articles []model.QiitaArticle) error
	GetArticles(userID uint) ([]model.QiitaArticle, error)
	GetStats(userID uint) (*model.QiitaStats, error)
	DeleteUserArticles(userID uint) error
}

// ZennRepositoryInterface はZenn連携データ操作の契約を定義する。
type ZennRepositoryInterface interface {
	UpsertArticles(userID uint, articles []model.ZennArticle) error
	GetArticles(userID uint) ([]model.ZennArticle, error)
	GetStats(userID uint) (*model.ZennStats, error)
	DeleteUserArticles(userID uint) error
}

// PasswordResetRepositoryInterface はパスワードリセットトークンデータ操作の契約を定義する。
type PasswordResetRepositoryInterface interface {
	Create(token *model.PasswordResetToken) error
	FindByToken(token string) (*model.PasswordResetToken, error)
	MarkAsUsed(id uint) error
	InvalidateUserTokens(userID uint) error
	DeleteExpired() error
}

// ActivityReportRepositoryInterface はアクティビティレポートデータ操作の契約を定義する。
// 週次・月次レポートの生成と前期間比較を提供する。
type ActivityReportRepositoryInterface interface {
	GetWeeklyReport(userID uint) (*model.ActivityReport, error)
	GetMonthlyReport(userID uint) (*model.ActivityReport, error)
	GetComparison(userID uint, period model.ReportPeriod) (*model.ReportComparison, error)
}

// AIAdviceRepositoryInterface はAIアドバイスデータ操作の契約を定義する。
// ルールエンジンが生成したアドバイスのCRUDと既読管理を含む。
type AIAdviceRepositoryInterface interface {
	Create(advice *model.AIAdvice) error
	CreateBatch(advices []*model.AIAdvice) error
	FindByUserID(userID uint, limit int) ([]model.AIAdvice, error)
	FindUnreadByUserID(userID uint) ([]model.AIAdvice, error)
	MarkAsRead(id, userID uint) error
	MarkAllAsRead(userID uint) error
	DeleteExpired() error
	DeleteByUserID(userID uint) error
}

// RecommendationRepositoryInterface はレコメンド機能のデータ操作の契約を定義する。
// スキルマッチングによるユーザー推薦とトレンドコンテンツ取得を提供する。
type RecommendationRepositoryInterface interface {
	GetRecommendedUsers(userID uint, skills []string, limit int) ([]model.RecommendedUser, error)
	GetTrendingPosts(limit int, days int) ([]model.Post, error)
	GetTrendingResources(limit int, days int) ([]model.LearningResource, error)
}

// LevelRepositoryInterface はレベルシステムのデータ操作の契約を定義する。
// 複数テーブルからXP計算に必要な統計を集計する。
type LevelRepositoryInterface interface {
	GetXPStats(userID uint) (*model.XPStats, error)
}

// LearningAnalyticsRepositoryInterface は学習分析データ操作の契約を定義する。
// 学習ログから集計したヒートマップ・カテゴリ別・週間トレンド・生産性統計を提供する。
type LearningAnalyticsRepositoryInterface interface {
	GetHeatmapData(userID uint) ([]model.HeatmapEntry, error)
	GetCategoryBreakdown(userID uint) ([]model.CategoryBreakdown, error)
	GetWeeklyTrends(userID uint, weeks int) ([]model.WeeklyTrend, error)
	GetProductivityStats(userID uint) (*model.ProductivityStats, error)
}

// BadgeRepositoryInterface はバッジ判定に必要な統計データ操作の契約を定義する。
// 複数テーブルからの集計をカプセル化する。
type BadgeRepositoryInterface interface {
	GetBadgeStats(userID uint) (*model.BadgeStats, error)
}

// AIConversationRepositoryInterface はAI会話セッションデータ操作の契約を定義する。
// LLMとの会話管理とレート制限用のメッセージカウントを含む。
type AIConversationRepositoryInterface interface {
	CreateConversation(conv *model.AIConversation) error
	FindConversationsByUserID(userID uint, limit, offset int) ([]model.AIConversation, error)
	FindConversationByID(id uint) (*model.AIConversation, error)
	AddMessage(msg *model.AIMessage) error
	GetMessages(conversationID uint) ([]model.AIMessage, error)
	CountTodayMessages(userID uint) (int64, error)
	DeleteConversation(id, userID uint) error
}

// StudyCircleRepositoryInterface はスタディサークルのデータ操作の契約を定義する。
// サークルCRUD、メンバー管理、ステップ管理、進捗、チェックイン、ストリークランキングを提供する。
type StudyCircleRepositoryInterface interface {
	// サークルCRUD
	Create(circle *model.StudyCircle) error
	FindByID(id uint) (*model.StudyCircle, error)
	FindByUserID(userID uint) ([]model.StudyCircle, error)
	Update(circle *model.StudyCircle) error
	Delete(id uint) error

	// メンバー管理
	AddMember(circleID, userID uint, role model.StudyCircleMemberRole) error
	RemoveMember(circleID, userID uint) error
	GetMembers(circleID uint) ([]model.StudyCircleMember, error)
	IsMember(circleID, userID uint) (bool, error)
	GetMemberCount(circleID uint) (int, error)

	// ステップCRUD
	CreateStep(step *model.StudyCircleStep) error
	UpdateStep(step *model.StudyCircleStep) error
	DeleteStep(stepID uint) error
	FindStepByID(stepID uint) (*model.StudyCircleStep, error)
	ReorderSteps(circleID uint, stepOrders []StepOrder) error

	// 進捗管理
	UpsertProgress(progress *model.StudyCircleMemberProgress) error
	GetProgress(circleID uint) ([]model.StudyCircleMemberProgress, error)

	// チェックイン
	CreateCheckin(checkin *model.StudyCircleCheckin) error
	GetCheckins(circleID uint) ([]model.StudyCircleCheckin, error)
	HasCheckedInToday(circleID, userID uint) (bool, error)

	// ストリークランキング
	GetStreakRanking(circleID uint) ([]model.CircleMemberStreak, error)
}

// Package repository はDevSyncアプリケーションのデータアクセス層を提供する。
// 各リポジトリはGORMを使用してPostgreSQLに対するCRUD操作を実装する。
package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
)



// PostAdvancedSearchRepositoryInterface は投稿の高度な検索フィルター機能の契約を定義する。
// タグ・日付範囲・ソート順による絞り込みを提供する。
type PostAdvancedSearchRepositoryInterface interface {
	SearchWithFilter(query string, tags []string, sortBy string, dateFrom, dateTo *time.Time, limit, offset int) ([]model.Post, int64, error)
}


// UserRepositoryInterface はユーザーデータ操作の契約を定義する。
type UserRepositoryInterface interface {
	FindAll() ([]model.User, error)
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
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
	CountAll() (int64, error)
	FindByUserID(userID uint, limit, offset int) ([]model.Post, int64, error)
	FindDraftsByUserID(userID uint) ([]model.Post, error)
	Timeline(userID uint, page, limit int) ([]model.Post, error)
	Update(post *model.Post) error
	Delete(id uint) error
	Like(userID, postID uint) error
	Unlike(userID, postID uint) error
	HasLiked(userID, postID uint) bool
	CreateComment(comment *model.Comment) error
	FindCommentByID(id uint) (*model.Comment, error)
	GetComments(postID uint) ([]model.Comment, error)
	GetReplies(parentID uint) ([]model.Comment, error)
	UpdateComment(comment *model.Comment) error
	DeleteComment(id uint) error
	Bookmark(userID, postID uint) error
	Unbookmark(userID, postID uint) error
	HasBookmarked(userID, postID uint) bool
	FindBookmarkedByUserID(userID uint, page, limit int) ([]model.Post, int64, error)
	CountBookmarkedByUserID(userID uint) (int64, error)
	AddReaction(userID, postID uint, emoji string) error
	RemoveReaction(userID, postID uint, emoji string) error
	GetReactionsByPostID(postID uint) ([]model.ReactionCount, error)
	GetUserReactions(userID, postID uint) ([]string, error)
	GetReactionsBatch(postIDs []uint) (map[uint][]model.ReactionCount, error)
	GetUserReactionsBatch(userID uint, postIDs []uint) (map[uint][]string, error)
	FindScheduledByUserID(userID uint) ([]model.Post, error)
	CountByUserID(userID uint) (int64, error)
	CountDraftsByUserID(userID uint) (int64, error)
	CountScheduledByUserID(userID uint) (int64, error)
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
	GetConversations(userID uint) ([]model.ConversationSummary, error)
	MarkAsRead(senderID, receiverID uint) error
}

// LearningLogRepositoryInterface は学習ログデータ操作の契約を定義する。
// ストリーク計算やカレンダーデータ取得も含む。
type LearningLogRepositoryInterface interface {
	Create(log *model.LearningLog) error
	CreateBatch(logs []model.LearningLog) error
	Update(log *model.LearningLog) error
	Delete(id, userID uint) error
	FindByID(id uint) (*model.LearningLog, error)
	GetByUserID(userID uint, limit, offset int) ([]model.LearningLog, int64, error)
	GetByCategory(userID uint, category string) ([]model.LearningLog, error)
	GetByPeriod(userID uint, days int) ([]model.LearningLog, error)
	GetBySource(userID uint, source string) ([]model.LearningLog, error)
	SumDurationByPeriod(userID uint, days int) (int, error)
	GetStreakInfo(userID uint) (*model.StreakInfo, error)
	GetCalendarData(userID uint) ([]model.CalendarEntry, error)
	GetRecentCategories(userID uint, limit int) ([]string, error)
	GetByGoalID(goalID uint, limit, offset int) ([]model.LearningLog, int64, error)
	SumDurationByGoalID(goalID uint) (int, error)
	GetFavorites(userID uint, limit, offset int) ([]model.LearningLog, int64, error)
	GetMonthlySummary(userID uint, months int) ([]model.MonthlySummary, error)
	CountByUserID(userID uint) (int64, error)
}

// LearningGoalRepositoryInterface は学習目標データ操作の契約を定義する。
// アクティブ目標の取得や統計情報の算出も含む。
type LearningGoalRepositoryInterface interface {
	Create(goal *model.LearningGoal) error
	Update(goal *model.LearningGoal) error
	Delete(id uint) error
	FindByID(id uint) (*model.LearningGoal, error)
	GetByUserID(userID uint, limit, offset int) ([]model.LearningGoal, int64, error)
	GetActiveByUserID(userID uint) ([]model.LearningGoal, error)
	GetByCategory(userID uint, category string) ([]model.LearningGoal, error)
	GetByStatus(userID uint, status string) ([]model.LearningGoal, error)
	GetPublicByUserID(userID uint, limit, offset int) ([]model.LearningGoal, int64, error)
	GetPublicGoals(limit, offset int) ([]model.LearningGoal, int64, error)
	GetStats(userID uint) (*model.LearningGoalStats, error)
	CountByUserID(userID uint) (int64, error)
}

// LearningResourceRepositoryInterface は学習リソースデータ操作の契約を定義する。
// いいね・ブックマーク（保存）操作も含む。
type LearningResourceRepositoryInterface interface {
	Create(resource *model.LearningResource) error
	FindByID(id uint) (*model.LearningResource, error)
	FindByUserID(userID uint, includePrivate bool, limit, offset int) ([]model.LearningResource, int64, error)
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
	FindByDifficulty(difficulty string, limit, offset int) ([]model.LearningResource, int64, error)
	CountByUserID(userID uint) (int64, error)
}

// RoadmapRepositoryInterface は学習ロードマップデータ操作の契約を定義する。
// ステップのCRUD操作やコピー機能も含む。
type RoadmapRepositoryInterface interface {
	Create(roadmap *model.Roadmap) error
	Update(roadmap *model.Roadmap) error
	Delete(id uint) error
	FindByID(id uint) (*model.Roadmap, error)
	GetByUserID(userID uint, limit, offset int) ([]model.Roadmap, int64, error)
	GetByStatus(userID uint, status string) ([]model.Roadmap, error)
	GetPublicRoadmaps(limit, offset int) ([]model.Roadmap, int64, error)
	CopyRoadmap(originalID, newUserID uint) (*model.Roadmap, error)
	GetStats(userID uint) (*model.RoadmapStats, error)
	CreateStep(step *model.RoadmapStep) error
	UpdateStep(step *model.RoadmapStep) error
	DeleteStep(stepID uint) error
	FindStepByID(stepID uint) (*model.RoadmapStep, error)
	ReorderSteps(roadmapID uint, stepOrders []model.StepOrder) error
	GetTemplates() ([]model.Roadmap, error)
	CountByUserID(userID uint) (int64, error)
}

// ChatRoomRepositoryInterface はチャットルームデータ操作の契約を定義する。
// メンバー管理操作も含む。
type ChatRoomRepositoryInterface interface {
	Create(room *model.ChatRoom) error
	FindByID(id uint) (*model.ChatRoom, error)
	FindByUserID(userID uint, limit, offset int) ([]model.ChatRoom, int64, error)
	Update(room *model.ChatRoom) error
	Delete(roomID uint) error
	AddMember(roomID, userID uint) error
	RemoveMember(roomID, userID uint) error
	GetMembers(roomID uint) ([]model.ChatRoomMember, error)
	IsMember(roomID, userID uint) (bool, error)
	CountByUserID(userID uint) (int64, error)
}

// GroupMessageRepositoryInterface はグループメッセージデータ操作の契約を定義する。
type GroupMessageRepositoryInterface interface {
	Create(msg *model.GroupMessage) error
	FindByRoomID(roomID uint, page, limit int) ([]model.GroupMessage, error)
	FindSenderByID(msg *model.GroupMessage)
	GetMemberUserIDs(roomID uint) []uint
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

// LevelRepositoryInterface はレベルシステムのデータ操作の契約を定義する。
// 複数テーブルからXP計算に必要な統計を集計する。
type LevelRepositoryInterface interface {
	GetXPStats(userID uint) (*model.XPStats, error)
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

// MentionRepositoryInterface はメンションデータ操作の契約を定義する。
type MentionRepositoryInterface interface {
	Create(mention *model.Mention) error
	FindByUserID(userID uint, page, limit int) ([]model.Mention, error)
	FindByPostID(postID uint) ([]model.Mention, error)
	FindByCommentID(commentID uint) ([]model.Mention, error)
	Delete(id uint) error
	DeleteByPostID(postID uint) error
	DeleteByCommentID(commentID uint) error
}

// SpotifyRepositoryInterface はSpotify連携データ操作の契約を定義する。
type SpotifyRepositoryInterface interface {
	DeleteUserData(userID uint) error
}

// YouTubeVideoRepositoryInterface はYouTube動画キャッシュデータ操作の契約を定義する。
type YouTubeVideoRepositoryInterface interface {
	UpsertVideos(videos []model.YouTubeVideo) error
	FindByVideoIDs(videoIDs []string) ([]model.YouTubeVideo, error)
	FindCachedSearch(query, language string) (*model.YouTubeSearchCache, error)
	SaveSearchCache(cache *model.YouTubeSearchCache) error
}














// Package router はDevSyncアプリケーションのルーティング設定を提供する。
// DIコンテナからハンドラを受け取り、Ginルーターへのエンドポイント登録を行う。
package router

import (
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/config"
	"github.com/norman6464/devsync/backend/internal/di"
	"github.com/norman6464/devsync/backend/internal/handler"
	"github.com/norman6464/devsync/backend/internal/middleware"
	"github.com/norman6464/devsync/backend/internal/service"
	"gorm.io/gorm"
)

// Setup はGinルーターを構築し、全エンドポイントを登録して返す。
// DIコンテナを利用して依存関係を解決し、ルーティングのみに集中する。
func Setup(db *gorm.DB, cfg *config.Config, hub *service.Hub) *gin.Engine {
	// DIコンテナ構築
	c := di.NewContainer(db, cfg, hub)

	r := gin.Default()

	// 信頼するプロキシを設定（X-Forwarded-For偽装によるレート制限回避を防止）
	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies != "" {
		proxies := strings.Split(trustedProxies, ",")
		for i := range proxies {
			proxies[i] = strings.TrimSpace(proxies[i])
		}
		if err := r.SetTrustedProxies(proxies); err != nil {
			log.Fatalf("信頼するプロキシの設定に失敗: %v", err)
		}
	} else {
		// プロキシ未設定の場合はローカルホストのみ信頼
		_ = r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	}

	// CORS設定（AllowCredentials=true の場合、ワイルドカードオリジンを拒否）
	origins := strings.Split(cfg.CORSOrigins, ",")
	for _, origin := range origins {
		if strings.TrimSpace(origin) == "*" {
			log.Fatal("CORS: AllowCredentials=true の場合、ワイルドカード '*' オリジンは使用できません")
		}
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// リクエストボディサイズ制限（1MB: JSON API用）
	r.Use(middleware.BodyLimit(1 << 20))

	// セキュリティヘッダー（全エンドポイント）
	r.Use(func(ctx *gin.Context) {
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("X-XSS-Protection", "1; mode=block")
		ctx.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		ctx.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		if strings.HasPrefix(ctx.Request.URL.Path, "/uploads/") {
			ctx.Header("Cache-Control", "public, max-age=86400")
			ctx.Header("Content-Security-Policy", "default-src 'none'; img-src 'self'; style-src 'none'; script-src 'none'")
			ctx.Header("Content-Disposition", "inline")
			ctx.Header("X-Content-Type-Options", "nosniff")
		} else {
			ctx.Header("Cache-Control", "no-store")
			ctx.Header("Content-Security-Policy", "default-src 'none'")
		}
		if os.Getenv("ENVIRONMENT") == "production" {
			ctx.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		ctx.Next()
	})
	r.Static("/uploads", "./uploads")

	// パブリックルート
	r.GET("/health", handler.HealthCheck)
	r.GET("/ws", c.WebSocketHandler.HandleWebSocket)

	api := r.Group("/api/v1")

	// レート制限ストア
	authRateLimitStore := middleware.NewRateLimitStore()
	apiRateLimitStore := middleware.NewRateLimitStore()

	// 定期クリーンアップ（メモリリーク対策）
	go middleware.StartCleanup(authRateLimitStore, 60)
	go middleware.StartCleanup(apiRateLimitStore, 60)

	// 認証ルート（公開: 10リクエスト/60秒）
	auth := api.Group("/auth")
	auth.Use(middleware.RateLimit(authRateLimitStore, 10, 60))
	{
		auth.POST("/register", c.AuthHandler.Register)
		auth.POST("/login", c.AuthHandler.Login)
		auth.GET("/github", c.AuthHandler.GitHubLogin)
		auth.GET("/github/callback", c.AuthHandler.GitHubLoginCallback)
		auth.POST("/password-reset/request", c.AuthHandler.RequestPasswordReset)
		auth.POST("/password-reset/confirm", c.AuthHandler.ResetPassword)
	}
	api.GET("/github/callback", c.GitHubHandler.Callback)
	api.GET("/spotify/callback", c.SpotifyHandler.Callback)

	// 認証必須ルート（100リクエスト/60秒）
	protected := api.Group("")
	protected.Use(middleware.AuthRequired(c.AuthService))
	protected.Use(middleware.RateLimit(apiRateLimitStore, 100, 60))
	{
		protected.GET("/auth/me", c.AuthHandler.Me)
		protected.POST("/auth/logout", c.AuthHandler.Logout)
		protected.DELETE("/auth/account", c.AuthHandler.DeleteAccount)

		registerUserRoutes(protected, c)
		registerGitHubRoutes(protected, c)
		registerPostRoutes(protected, c)
		registerPostSeriesRoutes(protected, c)
		registerPostCollectionRoutes(protected, c)
		registerPostTagRoutes(protected, c)
		registerPostPinRoutes(protected, c)
		registerPostViewRoutes(protected, c)
		registerMentionRoutes(protected, c)
		registerSnippetRoutes(protected, c)
		registerRankingRoutes(protected, c)
		registerMessageRoutes(protected, c)
		registerUploadRoutes(protected, c)
		registerNotificationRoutes(protected, c)
		registerIntegrationRoutes(protected, c)
		registerLearningRoutes(protected, c)
		registerBookmarkCollectionRoutes(protected, c)
		registerCommunityRoutes(protected, c)
		registerAnalyticsRoutes(protected, c)
		registerRecommendationRoutes(protected, c)
		registerStudyCircleRoutes(protected, c)
		registerSearchRoutes(protected, c)
		registerYouTubeRoutes(protected, c)
		registerSpotifyRoutes(protected, c)
		registerCommentLikeRoutes(protected, c)
		registerPostTemplateRoutes(protected, c)
		registerWidgetSettingsRoutes(protected, c)
		registerReminderSettingsRoutes(protected, c)
		registerNotificationSettingsRoutes(protected, c)
		registerUserStatsRoutes(protected, c)
		registerStudyCircleStatsRoutes(protected, c)
	}

	return r
}

func registerUserRoutes(g *gin.RouterGroup, c *di.Container) {
	users := g.Group("/users")
	{
		users.GET("", c.UserHandler.GetAll)
		users.GET("/by-username/:username", c.UserHandler.GetByUsername) // ユーザー名でユーザー取得（IDとの競合回避）
		users.GET("/:id", c.UserHandler.GetByID)
		users.PUT("/:id", c.UserHandler.Update)
		users.GET("/:id/followers", c.FollowHandler.GetFollowers)
		users.GET("/:id/following", c.FollowHandler.GetFollowing)
		users.POST("/:id/follow", c.FollowHandler.Follow)
		users.DELETE("/:id/follow", c.FollowHandler.Unfollow)
		users.GET("/:id/posts", c.PostHandler.GetUserPosts)
		users.GET("/me/profile-completeness", c.UserHandler.GetProfileCompleteness)
		users.GET("/:userId/activity", c.UserActivityHandler.GetTimeline)
	}
}

func registerGitHubRoutes(g *gin.RouterGroup, c *di.Container) {
	github := g.Group("/github")
	{
		github.GET("/connect", c.GitHubHandler.Connect)
		github.POST("/sync", c.GitHubHandler.Sync)
		github.DELETE("/disconnect", c.GitHubHandler.Disconnect)
		github.GET("/contributions/:userId", c.GitHubHandler.GetContributions)
		github.GET("/languages/:userId", c.GitHubHandler.GetLanguages)
		github.GET("/repos/:userId", c.GitHubHandler.GetRepos)
	}
}

func registerPostRoutes(g *gin.RouterGroup, c *di.Container) {
	posts := g.Group("/posts")
	{
		posts.POST("", c.PostHandler.Create)
		posts.GET("", c.PostHandler.GetAll)
		posts.GET("/timeline", c.PostHandler.Timeline)
		posts.GET("/drafts", c.PostHandler.GetDrafts)
		posts.PUT("/drafts/auto-save", c.PostHandler.AutoSaveDraft)
		posts.GET("/:id", c.PostHandler.GetByID)
		posts.PUT("/:id", c.PostHandler.Update)
		posts.PUT("/:id/publish", c.PostHandler.Publish)
		posts.PUT("/:id/unpublish", c.PostHandler.Unpublish)
		posts.PUT("/:id/schedule", c.PostHandler.SchedulePublish)
		posts.PUT("/:id/cancel-schedule", c.PostHandler.CancelSchedule)
		posts.GET("/scheduled", c.PostHandler.GetScheduled)
		posts.DELETE("/:id", c.PostHandler.Delete)
		posts.POST("/:id/like", c.PostHandler.Like)
		posts.DELETE("/:id/like", c.PostHandler.Unlike)
		posts.POST("/:id/bookmark", c.PostHandler.Bookmark)
		posts.DELETE("/:id/bookmark", c.PostHandler.Unbookmark)
		posts.GET("/bookmarks", c.PostHandler.GetBookmarks)
		posts.GET("/:id/comments", c.PostHandler.GetComments)
		posts.POST("/:id/comments", c.PostHandler.CreateComment)
		posts.GET("/:id/comments/:commentId/replies", c.PostHandler.GetReplies)
		posts.PUT("/:id/comments/:commentId", c.PostHandler.EditComment)
		posts.DELETE("/:id/comments/:commentId", c.PostHandler.DeleteComment)
		posts.POST("/:id/comments/:commentId/hide", c.PostHandler.HideComment)
		posts.POST("/:id/comments/:commentId/unhide", c.PostHandler.UnhideComment)
		posts.POST("/reactions/batch", c.PostHandler.GetReactionsBatch)
		posts.GET("/:id/reactions", c.PostHandler.GetReactions)
		posts.POST("/:id/reactions", c.PostHandler.AddReaction)
		posts.DELETE("/:id/reactions", c.PostHandler.RemoveReaction)
		posts.POST("/:id/snippets", c.CodeSnippetHandler.Create)
		posts.GET("/:id/snippets", c.CodeSnippetHandler.GetByPostID)
	}
}

func registerPostSeriesRoutes(g *gin.RouterGroup, c *di.Container) {
	series := g.Group("/post-series")
	{
		series.POST("", c.PostSeriesHandler.Create)
		series.GET("/:id", c.PostSeriesHandler.GetByID)
		series.PUT("/:id", c.PostSeriesHandler.Update)
		series.DELETE("/:id", c.PostSeriesHandler.Delete)
		series.GET("/:id/posts", c.PostSeriesHandler.GetPosts)
		series.POST("/:id/posts", c.PostSeriesHandler.AddPost)
		series.DELETE("/:id/posts/:postId", c.PostSeriesHandler.RemovePost)
		series.GET("/user/:userId", c.PostSeriesHandler.GetByUserID)
	}
}

func registerPostCollectionRoutes(g *gin.RouterGroup, c *di.Container) {
	collections := g.Group("/post-collections")
	{
		collections.POST("", c.PostCollectionHandler.Create)
		collections.GET("/:id", c.PostCollectionHandler.GetByID)
		collections.PUT("/:id", c.PostCollectionHandler.Update)
		collections.DELETE("/:id", c.PostCollectionHandler.Delete)
		collections.GET("/:id/posts", c.PostCollectionHandler.GetPosts)
		collections.POST("/:id/posts", c.PostCollectionHandler.AddPost)
		collections.DELETE("/:id/posts/:postId", c.PostCollectionHandler.RemovePost)
		collections.GET("/user/:userId", c.PostCollectionHandler.GetByUserID)
	}
}

func registerPostTagRoutes(g *gin.RouterGroup, c *di.Container) {
	tags := g.Group("/post-tags")
	{
		tags.PUT("/posts/:postId", c.PostTagHandler.SetTags)
		tags.GET("/posts/:postId", c.PostTagHandler.GetByPostID)
		tags.GET("/search", c.PostTagHandler.FindPostsByTag)
		tags.GET("/popular", c.PostTagHandler.GetPopularTags)
	}
}

func registerPostPinRoutes(g *gin.RouterGroup, c *di.Container) {
	pins := g.Group("/post-pins")
	{
		pins.POST("/posts/:postId", c.PostPinHandler.Pin)
		pins.DELETE("/posts/:postId", c.PostPinHandler.Unpin)
		pins.GET("/users/:userId", c.PostPinHandler.GetByUserID)
		pins.PUT("/reorder", c.PostPinHandler.Reorder)
	}
}

func registerPostViewRoutes(g *gin.RouterGroup, c *di.Container) {
	views := g.Group("/post-views")
	{
		views.POST("/posts/:postId", c.PostViewHandler.RecordView)
		views.GET("/posts/:postId", c.PostViewHandler.GetViewCount)
		views.GET("/popular", c.PostViewHandler.GetMostViewed)
	}
}

func registerSnippetRoutes(g *gin.RouterGroup, c *di.Container) {
	snippets := g.Group("/snippets")
	{
		snippets.GET("/:id", c.CodeSnippetHandler.GetByID)
		snippets.PUT("/:id", c.CodeSnippetHandler.Update)
		snippets.DELETE("/:id", c.CodeSnippetHandler.Delete)
		snippets.GET("/:id/comments", c.CodeSnippetHandler.GetComments)
		snippets.POST("/:id/comments", c.CodeSnippetHandler.CreateComment)
		snippets.DELETE("/:id/comments/:commentId", c.CodeSnippetHandler.DeleteComment)
		snippets.GET("/language/:language", c.CodeSnippetHandler.GetByUserLanguage)
		snippets.POST("/:id/fork", c.CodeSnippetHandler.Fork)
		snippets.POST("/:id/favorite", c.CodeSnippetHandler.Favorite)
		snippets.DELETE("/:id/favorite", c.CodeSnippetHandler.Unfavorite)
		snippets.GET("/favorites", c.CodeSnippetHandler.GetFavorites)
	}
}

func registerRankingRoutes(g *gin.RouterGroup, c *di.Container) {
	rankings := g.Group("/rankings")
	{
		rankings.GET("/contributions", c.RankingHandler.ContributionRanking)
		rankings.GET("/level", c.RankingHandler.LevelRanking)
		rankings.GET("/languages/:lang", c.RankingHandler.LanguageRanking)
		rankings.GET("/languages", c.RankingHandler.AvailableLanguages)
	}
}

func registerMessageRoutes(g *gin.RouterGroup, c *di.Container) {
	messages := g.Group("/messages")
	{
		messages.GET("", c.MessageHandler.GetConversations)
		messages.GET("/:userId", c.MessageHandler.GetMessages)
		messages.POST("/:userId", c.MessageHandler.SendMessage)
	}
}

func registerUploadRoutes(g *gin.RouterGroup, c *di.Container) {
	upload := g.Group("/upload")
	upload.Use(middleware.BodyLimit(10 << 20)) // 10MB: ファイルアップロード用
	{
		upload.POST("/image", c.UploadHandler.UploadImage)
		upload.POST("/images", c.UploadHandler.UploadMultipleImages)
	}
}

func registerNotificationRoutes(g *gin.RouterGroup, c *di.Container) {
	notifications := g.Group("/notifications")
	{
		notifications.GET("", c.NotificationHandler.GetAll)
		notifications.GET("/unread-count", c.NotificationHandler.GetUnreadCount)
		notifications.PUT("/:id/read", c.NotificationHandler.MarkAsRead)
		notifications.PUT("/read-all", c.NotificationHandler.MarkAllAsRead)
		notifications.DELETE("/:id", c.NotificationHandler.Delete)
	}
}

func registerIntegrationRoutes(g *gin.RouterGroup, c *di.Container) {
	// Zenn連携
	zenn := g.Group("/zenn")
	{
		zenn.POST("/connect", c.ZennHandler.Connect)
		zenn.DELETE("/disconnect", c.ZennHandler.Disconnect)
		zenn.POST("/sync", c.ZennHandler.Sync)
		zenn.GET("/articles/:userId", c.ZennHandler.GetArticles)
		zenn.GET("/stats/:userId", c.ZennHandler.GetStats)
	}

	// Qiita連携
	qiita := g.Group("/qiita")
	{
		qiita.POST("/connect", c.QiitaHandler.Connect)
		qiita.DELETE("/disconnect", c.QiitaHandler.Disconnect)
		qiita.POST("/sync", c.QiitaHandler.Sync)
		qiita.GET("/articles/:userId", c.QiitaHandler.GetArticles)
		qiita.GET("/stats/:userId", c.QiitaHandler.GetStats)
	}

	// AtCoder連携
	atcoder := g.Group("/atcoder")
	{
		atcoder.POST("/connect", c.AtCoderHandler.Connect)
		atcoder.DELETE("/disconnect", c.AtCoderHandler.Disconnect)
		atcoder.GET("/rating/:username", c.AtCoderHandler.GetRating)
	}
}

// registerLearningRoutes は学習関連のエンドポイントをまとめて登録する。
func registerLearningRoutes(g *gin.RouterGroup, c *di.Container) {
	registerGoalRoutes(g, c)
	registerReportRoutes(g, c)
	registerProjectRoutes(g, c)
	registerResourceRoutes(g, c)
	registerLearningLogRoutes(g, c)
	registerLearningLogTemplateRoutes(g, c)
	registerStreakFreezeRoutes(g, c)
	registerNoteRoutes(g, c)
	registerEmailPreferencesRoutes(g, c)
	registerWeeklyChallengeRoutes(g, c)
	registerWeeklyGoalRoutes(g, c)
}

func registerGoalRoutes(g *gin.RouterGroup, c *di.Container) {
	goals := g.Group("/goals")
	{
		goals.POST("", c.LearningGoalHandler.Create)
		goals.GET("", c.LearningGoalHandler.GetMyGoals)
		goals.GET("/:id", c.LearningGoalHandler.GetByID)
		goals.PUT("/:id", c.LearningGoalHandler.Update)
		goals.DELETE("/:id", c.LearningGoalHandler.Delete)
		goals.GET("/user/:userId", c.LearningGoalHandler.GetByUserID)
		goals.GET("/stats/:userId", c.LearningGoalHandler.GetStats)
		goals.GET("/deadline-alerts", c.LearningGoalHandler.GetDeadlineAlerts)
		goals.GET("/category/:category", c.LearningGoalHandler.GetByCategory)
		goals.GET("/status/:status", c.LearningGoalHandler.GetByStatus)
		goals.POST("/:id/duplicate", c.LearningGoalHandler.Duplicate)
		goals.PUT("/:id/share", c.LearningGoalHandler.ToggleShare)
		goals.GET("/public", c.LearningGoalHandler.GetPublicGoals)
		goals.GET("/public/user/:userId", c.LearningGoalHandler.GetPublicByUserID)
		goals.GET("/:id/linked-logs", c.LearningLogHandler.GetLinkedLogs)
		goals.GET("/:id/progress", c.LearningLogHandler.GetGoalProgress)
		goals.GET("/my-stats", c.LearningGoalHandler.GetMyStats)
		goals.GET("/active", c.LearningGoalHandler.GetActiveGoals)
		goals.PUT("/batch-progress", c.LearningGoalHandler.BatchUpdateProgress)
		goals.GET("/forecast", c.LearningGoalHandler.GetForecast)
	}
}

func registerReportRoutes(g *gin.RouterGroup, c *di.Container) {
	reports := g.Group("/reports")
	{
		reports.GET("/weekly", c.ActivityReportHandler.GetMyWeeklyReport)
		reports.GET("/monthly", c.ActivityReportHandler.GetMyMonthlyReport)
		reports.GET("/weekly/:userId", c.ActivityReportHandler.GetWeeklyReport)
		reports.GET("/monthly/:userId", c.ActivityReportHandler.GetMonthlyReport)
		reports.GET("/comparison", c.ActivityReportHandler.GetComparison)
	}
}

func registerProjectRoutes(g *gin.RouterGroup, c *di.Container) {
	projects := g.Group("/projects")
	{
		projects.POST("", c.ProjectHandler.Create)
		projects.GET("", c.ProjectHandler.GetAll)
		projects.GET("/:id", c.ProjectHandler.GetByID)
		projects.PUT("/:id", c.ProjectHandler.Update)
		projects.DELETE("/:id", c.ProjectHandler.Delete)
		projects.GET("/user/:userId", c.ProjectHandler.GetByUserID)
		projects.GET("/user/:userId/featured", c.ProjectHandler.GetFeatured)
		projects.PUT("/:id/archive", c.ProjectHandler.Archive)
		projects.PUT("/:id/unarchive", c.ProjectHandler.Unarchive)
		projects.GET("/archived", c.ProjectHandler.GetArchived)

		// マイルストーン
		projects.POST("/:projectId/milestones", c.ProjectMilestoneHandler.Create)
		projects.GET("/:projectId/milestones", c.ProjectMilestoneHandler.GetByProjectID)
		projects.PUT("/milestones/:milestoneId", c.ProjectMilestoneHandler.Update)
		projects.DELETE("/milestones/:milestoneId", c.ProjectMilestoneHandler.Delete)
	}
}

func registerResourceRoutes(g *gin.RouterGroup, c *di.Container) {
	resources := g.Group("/resources")
	{
		resources.POST("", c.LearningResourceHandler.Create)
		resources.GET("", c.LearningResourceHandler.GetPublic)
		resources.GET("/search", c.LearningResourceHandler.Search)
		resources.GET("/saved", c.LearningResourceHandler.GetSaved)
		resources.GET("/:id", c.LearningResourceHandler.GetByID)
		resources.PUT("/:id", c.LearningResourceHandler.Update)
		resources.DELETE("/:id", c.LearningResourceHandler.Delete)
		resources.POST("/:id/like", c.LearningResourceHandler.Like)
		resources.DELETE("/:id/like", c.LearningResourceHandler.Unlike)
		resources.POST("/:id/save", c.LearningResourceHandler.SaveResource)
		resources.DELETE("/:id/save", c.LearningResourceHandler.UnsaveResource)
		resources.GET("/user/:userId", c.LearningResourceHandler.GetByUserID)
		resources.GET("/difficulty/:difficulty", c.LearningResourceHandler.GetByDifficulty)

		// リソース進捗
		resources.PUT("/progress", c.ResourceProgressHandler.Upsert)
		resources.GET("/progress", c.ResourceProgressHandler.GetMyProgress)
		resources.GET("/:resourceId/progress", c.ResourceProgressHandler.GetByResource)

		// リソースレビュー
		resources.POST("/:id/reviews", c.ResourceReviewHandler.Create)
		resources.GET("/:id/reviews", c.ResourceReviewHandler.GetByResourceID)
		resources.PUT("/reviews/:reviewId", c.ResourceReviewHandler.Update)
		resources.DELETE("/reviews/:reviewId", c.ResourceReviewHandler.Delete)
	}
}

func registerStreakFreezeRoutes(g *gin.RouterGroup, c *di.Container) {
	streakFreezes := g.Group("/streak-freezes")
	{
		streakFreezes.POST("", c.StreakFreezeHandler.UseFreeze)
		streakFreezes.GET("/status", c.StreakFreezeHandler.GetStatus)
	}
}

func registerLearningLogTemplateRoutes(g *gin.RouterGroup, c *di.Container) {
	logTemplates := g.Group("/learning-log-templates")
	{
		logTemplates.POST("", c.LearningLogTemplateHandler.Create)
		logTemplates.GET("", c.LearningLogTemplateHandler.GetByUserID)
		logTemplates.GET("/default", c.LearningLogTemplateHandler.GetDefault)
		logTemplates.GET("/:id", c.LearningLogTemplateHandler.GetByID)
		logTemplates.PUT("/:id", c.LearningLogTemplateHandler.Update)
		logTemplates.DELETE("/:id", c.LearningLogTemplateHandler.Delete)
		logTemplates.POST("/:id/use", c.LearningLogTemplateHandler.UseTemplate)
	}
}

func registerLearningLogRoutes(g *gin.RouterGroup, c *di.Container) {
	learningLogs := g.Group("/learning-logs")
	{
		learningLogs.POST("", c.LearningLogHandler.Create)
		learningLogs.POST("/batch", c.LearningLogHandler.BatchCreate)
		learningLogs.POST("/import", c.LearningLogHandler.ImportCSV)
		learningLogs.GET("", c.LearningLogHandler.GetMyLogs)
		learningLogs.GET("/export", c.LearningLogHandler.ExportLogs)
		learningLogs.GET("/favorites", c.LearningLogHandler.GetFavorites)
		learningLogs.GET("/recent-categories", c.LearningLogHandler.GetRecentCategories)
		learningLogs.GET("/category/:category", c.LearningLogHandler.GetByCategory)
		learningLogs.GET("/source/:source", c.LearningLogHandler.GetBySource)
		learningLogs.GET("/user/:userId", c.LearningLogHandler.GetByUserID)
		learningLogs.GET("/calendar/:userId", c.LearningLogHandler.GetCalendarData)
		learningLogs.GET("/streak/:userId", c.LearningLogHandler.GetStreakInfo)
		learningLogs.GET("/weekly-duration/:userId", c.LearningLogHandler.GetWeeklyDuration)
		learningLogs.GET("/monthly-summary/:userId", c.LearningLogHandler.GetMonthlySummary)
		learningLogs.GET("/:id", c.LearningLogHandler.GetByID)
		learningLogs.PUT("/:id", c.LearningLogHandler.Update)
		learningLogs.PUT("/:id/favorite", c.LearningLogHandler.Favorite)
		learningLogs.PUT("/:id/unfavorite", c.LearningLogHandler.Unfavorite)
		learningLogs.DELETE("/:id", c.LearningLogHandler.Delete)
	}
}

func registerNoteRoutes(g *gin.RouterGroup, c *di.Container) {
	notes := g.Group("/notes")
	{
		notes.POST("", c.NoteHandler.Create)
		notes.GET("", c.NoteHandler.GetByUserID)
		notes.GET("/search", c.NoteHandler.Search)
		notes.GET("/favorites", c.NoteHandler.GetFavorites)
		notes.GET("/archived", c.NoteHandler.GetArchived)
		notes.GET("/folder/:folderId", c.NoteHandler.GetByFolderID)
		notes.GET("/:id", c.NoteHandler.GetByID)
		notes.PUT("/:id", c.NoteHandler.Update)
		notes.DELETE("/:id", c.NoteHandler.Delete)
		notes.PUT("/:id/favorite", c.NoteHandler.ToggleFavorite)
		notes.PUT("/:id/archive", c.NoteHandler.Archive)
		notes.PUT("/:id/unarchive", c.NoteHandler.Unarchive)
		notes.POST("/:id/duplicate", c.NoteHandler.Duplicate)
		notes.GET("/tags", c.NoteHandler.GetTags)
		notes.GET("/:id/export", c.NoteHandler.Export)
		notes.POST("/:id/links", c.NoteLinkHandler.CreateLink)
		notes.GET("/:id/links", c.NoteLinkHandler.GetLinks)
		notes.GET("/:id/backlinks", c.NoteLinkHandler.GetBacklinks)
		notes.GET("/:id/link-stats", c.NoteLinkHandler.GetLinkStats)
		notes.DELETE("/:id/links/:targetId", c.NoteLinkHandler.DeleteLink)
		notes.GET("/:id/versions", c.NoteVersionHandler.GetVersions)
		notes.GET("/:id/versions/:versionId", c.NoteVersionHandler.GetVersion)
		notes.POST("/:id/versions/:versionId/restore", c.NoteVersionHandler.RestoreVersion)
	}

	noteFolders := g.Group("/note-folders")
	{
		noteFolders.POST("", c.NoteFolderHandler.Create)
		noteFolders.GET("", c.NoteFolderHandler.GetByUserID)
		noteFolders.GET("/root", c.NoteFolderHandler.GetRootFolders)
		noteFolders.GET("/:id", c.NoteFolderHandler.GetByID)
		noteFolders.GET("/:id/children", c.NoteFolderHandler.GetChildren)
		noteFolders.PUT("/:id", c.NoteFolderHandler.Update)
		noteFolders.DELETE("/:id", c.NoteFolderHandler.Delete)
	}

	noteTemplates := g.Group("/note-templates")
	{
		noteTemplates.POST("", c.NoteTemplateHandler.Create)
		noteTemplates.GET("", c.NoteTemplateHandler.GetByUserID)
		noteTemplates.GET("/default", c.NoteTemplateHandler.GetDefault)
		noteTemplates.GET("/:id", c.NoteTemplateHandler.GetByID)
		noteTemplates.PUT("/:id", c.NoteTemplateHandler.Update)
		noteTemplates.DELETE("/:id", c.NoteTemplateHandler.Delete)
		noteTemplates.POST("/:id/use", c.NoteTemplateHandler.UseTemplate)
	}
}

func registerEmailPreferencesRoutes(g *gin.RouterGroup, c *di.Container) {
	g.GET("/email-preferences", c.EmailPreferencesHandler.GetPreferences)
	g.PUT("/email-preferences", c.EmailPreferencesHandler.UpdatePreferences)
}

func registerCommunityRoutes(g *gin.RouterGroup, c *di.Container) {
	// Q&A
	questions := g.Group("/questions")
	{
		questions.POST("", c.QuestionHandler.Create)
		questions.GET("", c.QuestionHandler.GetAll)
		questions.GET("/search", c.QuestionHandler.Search)
		questions.GET("/solved", c.QuestionHandler.GetSolved)
		questions.GET("/unanswered", c.QuestionHandler.GetUnanswered)
		questions.GET("/:id", c.QuestionHandler.GetByID)
		questions.PUT("/:id", c.QuestionHandler.Update)
		questions.DELETE("/:id", c.QuestionHandler.Delete)
		questions.POST("/:id/vote", c.QuestionHandler.Vote)
		questions.DELETE("/:id/vote", c.QuestionHandler.RemoveVote)
		questions.POST("/:id/bookmark", c.QuestionHandler.Bookmark)
		questions.DELETE("/:id/bookmark", c.QuestionHandler.Unbookmark)
		questions.GET("/bookmarks", c.QuestionHandler.GetBookmarks)
		questions.GET("/user/:userId", c.QuestionHandler.GetByUserID)

		questions.GET("/:id/answers", c.AnswerHandler.GetByQuestionID)
		questions.GET("/:id/answers/vote-range", c.AnswerHandler.GetByVoteRange)
		questions.POST("/:id/answers", c.AnswerHandler.Create)
		questions.PUT("/:id/answers/:answerId", c.AnswerHandler.Update)
		questions.DELETE("/:id/answers/:answerId", c.AnswerHandler.Delete)
		questions.PUT("/:id/answers/:answerId/best", c.AnswerHandler.SetBestAnswer)
		questions.POST("/:id/answers/:answerId/vote", c.AnswerHandler.Vote)
		questions.DELETE("/:id/answers/:answerId/vote", c.AnswerHandler.RemoveVote)
	}

	// ロードマップ
	roadmaps := g.Group("/roadmaps")
	{
		roadmaps.POST("", c.RoadmapHandler.Create)
		roadmaps.GET("", c.RoadmapHandler.GetMyRoadmaps)
		roadmaps.GET("/public", c.RoadmapHandler.GetPublicRoadmaps)
		roadmaps.GET("/status/:status", c.RoadmapHandler.GetByStatus)
		roadmaps.GET("/templates", c.RoadmapHandler.GetTemplates)
		roadmaps.POST("/templates/:id/use", c.RoadmapHandler.CreateFromTemplate)
		roadmaps.GET("/:id", c.RoadmapHandler.GetByID)
		roadmaps.PUT("/:id", c.RoadmapHandler.Update)
		roadmaps.DELETE("/:id", c.RoadmapHandler.Delete)
		roadmaps.POST("/:id/copy", c.RoadmapHandler.CopyRoadmap)
		roadmaps.POST("/:id/steps", c.RoadmapHandler.CreateStep)
		roadmaps.PUT("/:id/steps/:stepId", c.RoadmapHandler.UpdateStep)
		roadmaps.DELETE("/:id/steps/:stepId", c.RoadmapHandler.DeleteStep)
		roadmaps.POST("/:id/steps/batch-complete", c.RoadmapHandler.BatchCompleteSteps)
		roadmaps.PUT("/:id/steps/reorder", c.RoadmapHandler.ReorderSteps)
	}

	// チャットルーム
	chatRooms := g.Group("/chat-rooms")
	{
		chatRooms.POST("", c.ChatRoomHandler.Create)
		chatRooms.GET("", c.ChatRoomHandler.GetMyRooms)
		chatRooms.GET("/:id", c.ChatRoomHandler.GetByID)
		chatRooms.PUT("/:id", c.ChatRoomHandler.Update)
		chatRooms.DELETE("/:id", c.ChatRoomHandler.Delete)
		chatRooms.GET("/:id/members", c.ChatRoomHandler.GetMembers)
		chatRooms.POST("/:id/members", c.ChatRoomHandler.AddMember)
		chatRooms.DELETE("/:id/members/:userId", c.ChatRoomHandler.RemoveMember)
		chatRooms.GET("/:id/messages", c.ChatRoomHandler.GetMessages)
		chatRooms.POST("/:id/messages", c.ChatRoomHandler.SendMessage)
	}

	// 書籍レビュー
	bookReviews := g.Group("/book-reviews")
	{
		bookReviews.POST("", c.BookReviewHandler.Create)
		bookReviews.GET("", c.BookReviewHandler.GetAll)
		bookReviews.GET("/:id", c.BookReviewHandler.GetByID)
		bookReviews.PUT("/:id", c.BookReviewHandler.Update)
		bookReviews.DELETE("/:id", c.BookReviewHandler.Delete)
		bookReviews.GET("/user/:userId", c.BookReviewHandler.GetByUserID)
		bookReviews.GET("/search", c.BookReviewHandler.Search)
		bookReviews.GET("/rating", c.BookReviewHandler.GetByRating)
		bookReviews.PUT("/:id/archive", c.BookReviewHandler.Archive)
		bookReviews.PUT("/:id/unarchive", c.BookReviewHandler.Unarchive)
		bookReviews.PUT("/:id/status", c.BookReviewHandler.UpdateStatus)
		bookReviews.PUT("/:id/progress", c.BookReviewHandler.UpdateProgress)
	}

	// バッジ
	badges := g.Group("/badges")
	{
		badges.GET("/:userId", c.BadgeHandler.GetUserBadges)
		badges.POST("/notify", c.BadgeHandler.NotifyBadgeEarned)
	}

	// レベル
	level := g.Group("/level")
	{
		level.GET("/me", c.LevelHandler.GetMyLevelInfo)
		level.GET("/:userId", c.LevelHandler.GetLevelInfo)
		level.GET("/:userId/breakdown", c.LevelHandler.GetXPBreakdown)
	}
}

func registerAnalyticsRoutes(g *gin.RouterGroup, c *di.Container) {
	// 学習分析
	analytics := g.Group("/analytics")
	{
		analytics.GET("/heatmap/:userId", c.LearningAnalyticsHandler.GetHeatmap)
		analytics.GET("/categories/:userId", c.LearningAnalyticsHandler.GetCategoryBreakdown)
		analytics.GET("/productivity/:userId", c.LearningAnalyticsHandler.GetProductivityScore)
		analytics.GET("/trends/:userId", c.LearningAnalyticsHandler.GetWeeklyTrends)
		analytics.GET("/insights", c.LearningAnalyticsHandler.GetInsights)
		analytics.GET("/day-of-week/:userId", c.LearningAnalyticsHandler.GetDayOfWeekSummary)
		analytics.GET("/dashboard-summary", c.LearningDashboardHandler.GetSummary)
	}

	// AIアドバイス
	advice := g.Group("/advice")
	{
		advice.GET("", c.AIAdviceHandler.GetAdvice)
		advice.GET("/unread", c.AIAdviceHandler.GetUnreadAdvice)
		advice.PUT("/:id/read", c.AIAdviceHandler.MarkAsRead)
		advice.POST("/chat", c.AIAdviceHandler.Chat)
		advice.GET("/conversations", c.AIAdviceHandler.GetConversations)
		advice.GET("/conversations/:id", c.AIAdviceHandler.GetConversation)
		advice.DELETE("/conversations/:id", c.AIAdviceHandler.DeleteConversation)
	}
}

func registerStudyCircleRoutes(g *gin.RouterGroup, c *di.Container) {
	circles := g.Group("/study-circles")
	{
		circles.POST("", c.StudyCircleHandler.Create)
		circles.GET("", c.StudyCircleHandler.GetMyCircles)
		circles.GET("/status/:status", c.StudyCircleHandler.GetByStatus)
		circles.GET("/:id", c.StudyCircleHandler.GetByID)
		circles.PUT("/:id", c.StudyCircleHandler.Update)
		circles.DELETE("/:id", c.StudyCircleHandler.Delete)
		circles.GET("/:id/members", c.StudyCircleHandler.GetMembers)
		circles.POST("/:id/members", c.StudyCircleHandler.AddMember)
		circles.DELETE("/:id/members/:userId", c.StudyCircleHandler.RemoveMember)
		circles.PUT("/:id/members/:userId/role", c.StudyCircleHandler.UpdateMemberRole)
		circles.POST("/:id/steps", c.StudyCircleHandler.CreateStep)
		circles.PUT("/:id/steps/:stepId", c.StudyCircleHandler.UpdateStep)
		circles.DELETE("/:id/steps/:stepId", c.StudyCircleHandler.DeleteStep)
		circles.PUT("/:id/steps/reorder", c.StudyCircleHandler.ReorderSteps)
		circles.PUT("/:id/steps/:stepId/progress", c.StudyCircleHandler.UpdateProgress)
		circles.GET("/:id/progress", c.StudyCircleHandler.GetProgress)
		circles.POST("/:id/checkins", c.StudyCircleHandler.CreateCheckin)
		circles.GET("/:id/checkins", c.StudyCircleHandler.GetCheckins)
		circles.GET("/:id/streak-ranking", c.StudyCircleHandler.GetStreakRanking)
	}
}

func registerRecommendationRoutes(g *gin.RouterGroup, c *di.Container) {
	recommendations := g.Group("/recommendations")
	{
		recommendations.GET("/users", c.RecommendationHandler.GetRecommendedUsers)
		recommendations.GET("/posts", c.RecommendationHandler.GetTrendingPosts)
		recommendations.GET("/resources", c.RecommendationHandler.GetTrendingResources)
	}
}

func registerMentionRoutes(g *gin.RouterGroup, c *di.Container) {
	mentions := g.Group("/mentions")
	{
		mentions.GET("", c.MentionHandler.GetMyMentions)
		mentions.GET("/posts/:postId", c.MentionHandler.GetPostMentions)
	}
}

func registerSearchRoutes(g *gin.RouterGroup, c *di.Container) {
	search := g.Group("/search")
	{
		search.GET("/posts", c.SearchHandler.SearchPosts)
		search.GET("/circles", c.SearchHandler.SearchCircles)
	}
}

func registerYouTubeRoutes(g *gin.RouterGroup, c *di.Container) {
	youtube := g.Group("/youtube")
	{
		youtube.GET("/search", c.YouTubeHandler.Search)
		youtube.GET("/recommend", c.YouTubeHandler.Recommend)
		youtube.GET("/status", c.YouTubeHandler.Status)
	}
}

func registerSpotifyRoutes(g *gin.RouterGroup, c *di.Container) {
	spotify := g.Group("/spotify")
	{
		spotify.GET("/connect", c.SpotifyHandler.Connect)
		spotify.DELETE("/disconnect", c.SpotifyHandler.Disconnect)
		spotify.GET("/currently-playing/:userId", c.SpotifyHandler.GetCurrentlyPlaying)
		spotify.GET("/recently-played/:userId", c.SpotifyHandler.GetRecentlyPlayed)
	}
}

func registerCommentLikeRoutes(g *gin.RouterGroup, c *di.Container) {
	comments := g.Group("/comments")
	{
		comments.POST("/:id/likes", c.CommentLikeHandler.Like)
		comments.DELETE("/:id/likes", c.CommentLikeHandler.Unlike)
		comments.GET("/:id/likes", c.CommentLikeHandler.GetStatus)
	}
}

// registerUserStatsRoutes はユーザー統計系のエンドポイントをまとめて登録する。
// 全て /users/:id/xxx-stats の形式で統一されている。
func registerUserStatsRoutes(g *gin.RouterGroup, c *di.Container) {
	users := g.Group("/users")
	{
		users.GET("/:id/dashboard-stats", c.UserDashboardHandler.GetStats)
		users.GET("/:id/note-stats", c.NoteStatsHandler.GetStats)
		users.GET("/:id/post-stats", c.PostStatsHandler.GetStats)
		users.GET("/:id/book-review-stats", c.BookReviewStatsHandler.GetStats)
		users.GET("/:id/qa-stats", c.QAStatsHandler.GetStats)
		users.GET("/:id/code-snippet-stats", c.CodeSnippetStatsHandler.GetStats)
		users.GET("/:id/learning-resource-stats", c.LearningResourceStatsHandler.GetStats)
		users.GET("/:id/project-stats", c.ProjectStatsHandler.GetStats)
		users.GET("/:id/follow-stats", c.FollowStatsHandler.GetStats)
		users.GET("/:id/roadmap-stats", c.RoadmapStatsHandler.GetStats)
		users.GET("/:id/learning-log-stats", c.LearningLogStatsHandler.GetStats)
		users.GET("/:id/comment-stats", c.CommentStatsHandler.GetStats)
		users.GET("/:id/notification-stats", c.NotificationStatsHandler.GetStats)
		users.GET("/:id/message-stats", c.MessageStatsHandler.GetStats)
		users.GET("/:id/mention-stats", c.MentionStatsHandler.GetStats)
		users.GET("/:id/reaction-stats", c.ReactionStatsHandler.GetStats)
		users.GET("/:id/reaction-summary", c.ReactionStatsHandler.GetSummary)
		users.GET("/:id/bookmark-stats", c.BookmarkStatsHandler.GetStats)
	}
}

func registerBookmarkCollectionRoutes(g *gin.RouterGroup, c *di.Container) {
	collections := g.Group("/bookmark-collections")
	{
		collections.POST("", c.BookmarkCollectionHandler.Create)
		collections.GET("", c.BookmarkCollectionHandler.GetMyCollections)
		collections.PUT("/:id", c.BookmarkCollectionHandler.Update)
		collections.DELETE("/:id", c.BookmarkCollectionHandler.Delete)
		collections.POST("/:id/posts/:postId", c.BookmarkCollectionHandler.AddPost)
		collections.DELETE("/:id/posts/:postId", c.BookmarkCollectionHandler.RemovePost)
		collections.GET("/:id/posts", c.BookmarkCollectionHandler.GetPosts)
	}
}

func registerWeeklyChallengeRoutes(g *gin.RouterGroup, c *di.Container) {
	challenges := g.Group("/weekly-challenges")
	{
		challenges.GET("/current", c.WeeklyChallengeHandler.GetCurrent)
		challenges.PUT("/progress", c.WeeklyChallengeHandler.UpdateProgress)
	}
}

func registerWeeklyGoalRoutes(g *gin.RouterGroup, c *di.Container) {
	weeklyGoals := g.Group("/weekly-goals")
	{
		weeklyGoals.PUT("", c.WeeklyGoalHandler.SetGoal)
		weeklyGoals.GET("", c.WeeklyGoalHandler.GetGoals)
		weeklyGoals.GET("/progress", c.WeeklyGoalHandler.GetProgress)
	}
}

func registerPostTemplateRoutes(g *gin.RouterGroup, c *di.Container) {
	templates := g.Group("/post-templates")
	{
		templates.POST("", c.PostTemplateHandler.Create)
		templates.GET("", c.PostTemplateHandler.GetMyTemplates)
		templates.GET("/:id", c.PostTemplateHandler.GetByID)
		templates.PUT("/:id", c.PostTemplateHandler.Update)
		templates.DELETE("/:id", c.PostTemplateHandler.Delete)
	}
}

func registerWidgetSettingsRoutes(g *gin.RouterGroup, c *di.Container) {
	widgets := g.Group("/widget-settings")
	{
		widgets.GET("", c.WidgetSettingsHandler.GetSettings)
		widgets.PUT("", c.WidgetSettingsHandler.UpdateSettings)
	}
}

func registerReminderSettingsRoutes(g *gin.RouterGroup, c *di.Container) {
	reminders := g.Group("/reminder-settings")
	{
		reminders.GET("", c.ReminderSettingsHandler.GetSettings)
		reminders.PUT("", c.ReminderSettingsHandler.UpdateSettings)
	}
}

func registerNotificationSettingsRoutes(g *gin.RouterGroup, c *di.Container) {
	notifSettings := g.Group("/notification-settings")
	{
		notifSettings.GET("", c.NotificationSettingsHandler.GetSettings)
		notifSettings.PUT("", c.NotificationSettingsHandler.UpdateSettings)
	}
}

func registerStudyCircleStatsRoutes(g *gin.RouterGroup, c *di.Container) {
	circles := g.Group("/study-circles")
	{
		circles.GET("/:id/stats", c.StudyCircleStatsHandler.GetStats)
	}
}

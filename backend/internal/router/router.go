// Package router はDevSyncアプリケーションのルーティング設定を提供する。
// DI（依存性注入）によるリポジトリ・サービス・ハンドラの構築と、
// Ginルーターへのエンドポイント登録を行う。
package router

import (
	"log"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/config"
	"github.com/norman6464/devsync/backend/internal/handler"
	"github.com/norman6464/devsync/backend/internal/middleware"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
	"github.com/norman6464/devsync/backend/internal/service"
	"gorm.io/gorm"
)

// Setup はGinルーターを構築し、全エンドポイントを登録して返す。
// リポジトリ → サービス → ハンドラの依存関係を手動DIで構築する。
func Setup(db *gorm.DB, cfg *config.Config, hub *service.Hub) *gin.Engine {
	r := gin.Default()

	// CORS設定
	origins := strings.Split(cfg.CORSOrigins, ",")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(db)
	followRepo := repository.NewFollowRepository(db)
	githubRepo := repository.NewGitHubRepository(db)
	postRepo := repository.NewPostRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	rankingRepo := repository.NewRankingRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)
	zennRepo := repository.NewZennRepository(db)
	qiitaRepo := repository.NewQiitaRepository(db)
	learningGoalRepo := repository.NewLearningGoalRepository(db)
	activityReportRepo := repository.NewActivityReportRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	learningResourceRepo := repository.NewLearningResourceRepository(db)
	bookReviewRepo := repository.NewBookReviewRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	answerRepo := repository.NewAnswerRepository(db)
	roadmapRepo := repository.NewRoadmapRepository(db)
	chatRoomRepo := repository.NewChatRoomRepository(db)
	groupMessageRepo := repository.NewGroupMessageRepository(db)
	learningLogRepo := repository.NewLearningLogRepository(db)
	codeSnippetRepo := repository.NewCodeSnippetRepository(db)
	aiAdviceRepo := repository.NewAIAdviceRepository(db)
	aiConversationRepo := repository.NewAIConversationRepository(db)

	// サービスの初期化
	authService := service.NewAuthService(userRepo, passwordResetRepo, cfg.JWTSecret)
	githubService := service.NewGitHubService(cfg, userRepo, githubRepo)
	zennService := service.NewZennService(userRepo, zennRepo)
	qiitaService := service.NewQiitaService(userRepo, qiitaRepo)
	notificationService := service.NewNotificationService(notificationRepo)
	postService := service.NewPostService(postRepo, notificationService)
	userService := service.NewUserService(userRepo)
	followService := service.NewFollowService(followRepo)
	questionService := service.NewQuestionService(questionRepo)
	answerService := service.NewAnswerService(answerRepo, questionRepo)
	learningLogService := service.NewLearningLogService(learningLogRepo)
	learningGoalService := service.NewLearningGoalService(learningGoalRepo)
	messageService := service.NewMessageService(messageRepo, notificationService)
	rankingService := service.NewRankingService(rankingRepo)
	projectService := service.NewProjectService(projectRepo)
	bookReviewService := service.NewBookReviewService(bookReviewRepo)
	learningResourceService := service.NewLearningResourceService(learningResourceRepo)
	codeSnippetService := service.NewCodeSnippetService(codeSnippetRepo, postRepo)
	roadmapService := service.NewRoadmapService(roadmapRepo)
	// テンプレートロードマップの初期登録（システムユーザーを取得/作成して使用）
	go func() {
		systemUser, err := getOrCreateSystemUser(db)
		if err != nil {
			log.Printf("テンプレートシード用システムユーザー作成失敗: %v", err)
			return
		}
		if err := roadmapService.SeedTemplates(systemUser.ID); err != nil {
			log.Printf("テンプレートシード失敗: %v", err)
		}
	}()
	chatRoomService := service.NewChatRoomService(chatRoomRepo, groupMessageRepo, hub)
	activityReportService := service.NewActivityReportService(activityReportRepo)
	atcoderService := service.NewAtCoderService()
	badgeService := service.NewBadgeService(db, notificationService)

	// AIアドバイスサービス（LLMクライアントはAPIキー設定時のみ初期化）
	var llmClient service.LLMClientInterface
	if cfg.OpenAIAPIKey != "" {
		llmClient = service.NewOpenAIClient(cfg.OpenAIAPIKey)
		log.Println("OpenAI APIキーが設定されています。LLMチャット機能が有効です。")
	} else {
		log.Println("OpenAI APIキー未設定。ルールベース推薦のみ有効です。")
	}
	aiAdviceService := service.NewAIAdviceService(
		aiAdviceRepo, aiConversationRepo,
		learningGoalRepo, learningLogRepo, roadmapRepo,
		githubRepo, learningResourceRepo, userRepo,
		llmClient,
	)

	// ハンドラの初期化
	authHandler := handler.NewAuthHandler(authService, githubService)
	userHandler := handler.NewUserHandler(userService)
	followHandler := handler.NewFollowHandler(followService)
	githubHandler := handler.NewGitHubHandler(githubService, authService)
	postHandler := handler.NewPostHandler(postService, codeSnippetService)
	snippetHandler := handler.NewCodeSnippetHandler(codeSnippetService)
	rankingHandler := handler.NewRankingHandler(rankingService)
	messageHandler := handler.NewMessageHandler(messageService)
	wsHandler := handler.NewWebSocketHandler(hub, authService)
	uploadHandler := handler.NewUploadHandler()
	notificationHandler := handler.NewNotificationHandler(notificationService)
	zennHandler := handler.NewZennHandler(zennService)
	qiitaHandler := handler.NewQiitaHandler(qiitaService)
	learningGoalHandler := handler.NewLearningGoalHandler(learningGoalService)
	activityReportHandler := handler.NewActivityReportHandler(activityReportService)
	projectHandler := handler.NewProjectHandler(projectService)
	learningResourceHandler := handler.NewLearningResourceHandler(learningResourceService)
	bookReviewHandler := handler.NewBookReviewHandler(bookReviewService)
	questionHandler := handler.NewQuestionHandler(questionService)
	answerHandler := handler.NewAnswerHandler(answerService)
	roadmapHandler := handler.NewRoadmapHandler(roadmapService)
	chatRoomHandler := handler.NewChatRoomHandler(chatRoomService)
	atcoderHandler := handler.NewAtCoderHandler(atcoderService, userService)
	badgeHandler := handler.NewBadgeHandler(badgeService)
	learningLogHandler := handler.NewLearningLogHandler(learningLogService)
	aiAdviceHandler := handler.NewAIAdviceHandler(aiAdviceService)

	// HubのGetRoomMembersコールバックを設定
	hub.GetRoomMembers = groupMessageRepo.GetMemberUserIDs

	// アップロードファイルの静的配信
	r.Static("/uploads", "./uploads")

	// パブリックルート
	r.GET("/health", handler.HealthCheck)

	// WebSocket（クエリパラメータで認証）
	r.GET("/ws", wsHandler.HandleWebSocket)

	api := r.Group("/api/v1")

	// 認証ルート（公開）
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/github", authHandler.GitHubLogin)
		auth.GET("/github/callback", authHandler.GitHubLoginCallback)
		auth.POST("/password-reset/request", authHandler.RequestPasswordReset)
		auth.POST("/password-reset/confirm", authHandler.ResetPassword)
	}

	// GitHubデータ連携コールバック（公開 - フロントエンドのOAuthリダイレクト後に呼ばれる）
	api.GET("/github/callback", githubHandler.Callback)

	// 認証必須ルート
	protected := api.Group("")
	protected.Use(middleware.AuthRequired(authService))
	{
		// 認証
		protected.GET("/auth/me", authHandler.Me)
		protected.POST("/auth/logout", authHandler.Logout)
		protected.DELETE("/auth/account", authHandler.DeleteAccount)

		// ユーザー
		users := protected.Group("/users")
		{
			users.GET("", userHandler.GetAll)
			users.GET("/:id", userHandler.GetByID)
			users.PUT("/:id", userHandler.Update)
			users.GET("/:id/followers", followHandler.GetFollowers)
			users.GET("/:id/following", followHandler.GetFollowing)
			users.POST("/:id/follow", followHandler.Follow)
			users.DELETE("/:id/follow", followHandler.Unfollow)
			users.GET("/:id/posts", postHandler.GetUserPosts)
		}

		// GitHub連携
		github := protected.Group("/github")
		{
			github.GET("/connect", githubHandler.Connect)
			github.POST("/sync", githubHandler.Sync)
			github.DELETE("/disconnect", githubHandler.Disconnect)
			github.GET("/contributions/:userId", githubHandler.GetContributions)
			github.GET("/languages/:userId", githubHandler.GetLanguages)
			github.GET("/repos/:userId", githubHandler.GetRepos)
		}

		// 投稿
		posts := protected.Group("/posts")
		{
			posts.POST("", postHandler.Create)
			posts.GET("", postHandler.GetAll)
			posts.GET("/timeline", postHandler.Timeline)
			posts.GET("/:id", postHandler.GetByID)
			posts.PUT("/:id", postHandler.Update)
			posts.DELETE("/:id", postHandler.Delete)
			posts.POST("/:id/like", postHandler.Like)
			posts.DELETE("/:id/like", postHandler.Unlike)
			posts.GET("/:id/comments", postHandler.GetComments)
			posts.POST("/:id/comments", postHandler.CreateComment)
			posts.DELETE("/:id/comments/:commentId", postHandler.DeleteComment)
			// コードスニペット（投稿の子リソース）
			posts.POST("/:id/snippets", snippetHandler.Create)
			posts.GET("/:id/snippets", snippetHandler.GetByPostID)
		}

		// コードスニペット（個別操作）
		snippets := protected.Group("/snippets")
		{
			snippets.GET("/:id", snippetHandler.GetByID)
			snippets.PUT("/:id", snippetHandler.Update)
			snippets.DELETE("/:id", snippetHandler.Delete)
			snippets.GET("/:id/comments", snippetHandler.GetComments)
			snippets.POST("/:id/comments", snippetHandler.CreateComment)
			snippets.DELETE("/:id/comments/:commentId", snippetHandler.DeleteComment)
		}

		// ランキング
		rankings := protected.Group("/rankings")
		{
			rankings.GET("/contributions", rankingHandler.ContributionRanking)
			rankings.GET("/languages/:lang", rankingHandler.LanguageRanking)
			rankings.GET("/languages", rankingHandler.AvailableLanguages)
		}

		// メッセージ
		messages := protected.Group("/messages")
		{
			messages.GET("", messageHandler.GetConversations)
			messages.GET("/:userId", messageHandler.GetMessages)
			messages.POST("/:userId", messageHandler.SendMessage)
		}

		// アップロード
		upload := protected.Group("/upload")
		{
			upload.POST("/image", uploadHandler.UploadImage)
			upload.POST("/images", uploadHandler.UploadMultipleImages)
		}

		// 通知
		notifications := protected.Group("/notifications")
		{
			notifications.GET("", notificationHandler.GetAll)
			notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
			notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
			notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
			notifications.DELETE("/:id", notificationHandler.Delete)
		}

		// Zenn連携
		zenn := protected.Group("/zenn")
		{
			zenn.POST("/connect", zennHandler.Connect)
			zenn.DELETE("/disconnect", zennHandler.Disconnect)
			zenn.POST("/sync", zennHandler.Sync)
			zenn.GET("/articles/:userId", zennHandler.GetArticles)
			zenn.GET("/stats/:userId", zennHandler.GetStats)
		}

		// Qiita連携
		qiita := protected.Group("/qiita")
		{
			qiita.POST("/connect", qiitaHandler.Connect)
			qiita.DELETE("/disconnect", qiitaHandler.Disconnect)
			qiita.POST("/sync", qiitaHandler.Sync)
			qiita.GET("/articles/:userId", qiitaHandler.GetArticles)
			qiita.GET("/stats/:userId", qiitaHandler.GetStats)
		}

		// AtCoder連携
		atcoder := protected.Group("/atcoder")
		{
			atcoder.POST("/connect", atcoderHandler.Connect)
			atcoder.DELETE("/disconnect", atcoderHandler.Disconnect)
			atcoder.GET("/rating/:username", atcoderHandler.GetRating)
		}

		// 学習目標
		goals := protected.Group("/goals")
		{
			goals.POST("", learningGoalHandler.Create)
			goals.GET("", learningGoalHandler.GetMyGoals)
			goals.GET("/:id", learningGoalHandler.GetByID)
			goals.PUT("/:id", learningGoalHandler.Update)
			goals.DELETE("/:id", learningGoalHandler.Delete)
			goals.GET("/user/:userId", learningGoalHandler.GetByUserID)
			goals.GET("/stats/:userId", learningGoalHandler.GetStats)
		}

		// アクティビティレポート
		reports := protected.Group("/reports")
		{
			reports.GET("/weekly", activityReportHandler.GetMyWeeklyReport)
			reports.GET("/monthly", activityReportHandler.GetMyMonthlyReport)
			reports.GET("/weekly/:userId", activityReportHandler.GetWeeklyReport)
			reports.GET("/monthly/:userId", activityReportHandler.GetMonthlyReport)
			reports.GET("/comparison", activityReportHandler.GetComparison)
		}

		// プロジェクトショーケース
		projects := protected.Group("/projects")
		{
			projects.POST("", projectHandler.Create)
			projects.GET("", projectHandler.GetAll)
			projects.GET("/:id", projectHandler.GetByID)
			projects.PUT("/:id", projectHandler.Update)
			projects.DELETE("/:id", projectHandler.Delete)
			projects.GET("/user/:userId", projectHandler.GetByUserID)
			projects.GET("/user/:userId/featured", projectHandler.GetFeatured)
		}

		// 学習リソース
		resources := protected.Group("/resources")
		{
			resources.POST("", learningResourceHandler.Create)
			resources.GET("", learningResourceHandler.GetPublic)
			resources.GET("/search", learningResourceHandler.Search)
			resources.GET("/saved", learningResourceHandler.GetSaved)
			resources.GET("/:id", learningResourceHandler.GetByID)
			resources.PUT("/:id", learningResourceHandler.Update)
			resources.DELETE("/:id", learningResourceHandler.Delete)
			resources.POST("/:id/like", learningResourceHandler.Like)
			resources.DELETE("/:id/like", learningResourceHandler.Unlike)
			resources.POST("/:id/save", learningResourceHandler.SaveResource)
			resources.DELETE("/:id/save", learningResourceHandler.UnsaveResource)
			resources.GET("/user/:userId", learningResourceHandler.GetByUserID)
		}

		// Q&A（質問）
		questions := protected.Group("/questions")
		{
			questions.POST("", questionHandler.Create)
			questions.GET("", questionHandler.GetAll)
			questions.GET("/search", questionHandler.Search)
			questions.GET("/:id", questionHandler.GetByID)
			questions.PUT("/:id", questionHandler.Update)
			questions.DELETE("/:id", questionHandler.Delete)
			questions.POST("/:id/vote", questionHandler.Vote)
			questions.DELETE("/:id/vote", questionHandler.RemoveVote)
			questions.GET("/user/:userId", questionHandler.GetByUserID)

			// 回答（質問の子リソース）
			questions.GET("/:id/answers", answerHandler.GetByQuestionID)
			questions.POST("/:id/answers", answerHandler.Create)
			questions.PUT("/:id/answers/:answerId", answerHandler.Update)
			questions.DELETE("/:id/answers/:answerId", answerHandler.Delete)
			questions.PUT("/:id/answers/:answerId/best", answerHandler.SetBestAnswer)
			questions.POST("/:id/answers/:answerId/vote", answerHandler.Vote)
			questions.DELETE("/:id/answers/:answerId/vote", answerHandler.RemoveVote)
		}

		// 学習ロードマップ
		roadmaps := protected.Group("/roadmaps")
		{
			roadmaps.POST("", roadmapHandler.Create)
			roadmaps.GET("", roadmapHandler.GetMyRoadmaps)
			roadmaps.GET("/public", roadmapHandler.GetPublicRoadmaps)
			roadmaps.GET("/templates", roadmapHandler.GetTemplates)
			roadmaps.POST("/templates/:id/use", roadmapHandler.CreateFromTemplate)
			roadmaps.GET("/:id", roadmapHandler.GetByID)
			roadmaps.PUT("/:id", roadmapHandler.Update)
			roadmaps.DELETE("/:id", roadmapHandler.Delete)
			roadmaps.POST("/:id/copy", roadmapHandler.CopyRoadmap)

			// ステップ
			roadmaps.POST("/:id/steps", roadmapHandler.CreateStep)
			roadmaps.PUT("/:id/steps/:stepId", roadmapHandler.UpdateStep)
			roadmaps.DELETE("/:id/steps/:stepId", roadmapHandler.DeleteStep)
			roadmaps.PUT("/:id/steps/reorder", roadmapHandler.ReorderSteps)
		}

		// チャットルーム
		chatRooms := protected.Group("/chat-rooms")
		{
			chatRooms.POST("", chatRoomHandler.Create)
			chatRooms.GET("", chatRoomHandler.GetMyRooms)
			chatRooms.GET("/:id", chatRoomHandler.GetByID)
			chatRooms.PUT("/:id", chatRoomHandler.Update)
			chatRooms.DELETE("/:id", chatRoomHandler.Delete)
			chatRooms.GET("/:id/members", chatRoomHandler.GetMembers)
			chatRooms.POST("/:id/members", chatRoomHandler.AddMember)
			chatRooms.DELETE("/:id/members/:userId", chatRoomHandler.RemoveMember)
			chatRooms.GET("/:id/messages", chatRoomHandler.GetMessages)
			chatRooms.POST("/:id/messages", chatRoomHandler.SendMessage)
		}

		// 書籍レビュー
		bookReviews := protected.Group("/book-reviews")
		{
			bookReviews.POST("", bookReviewHandler.Create)
			bookReviews.GET("", bookReviewHandler.GetAll)
			bookReviews.GET("/:id", bookReviewHandler.GetByID)
			bookReviews.PUT("/:id", bookReviewHandler.Update)
			bookReviews.DELETE("/:id", bookReviewHandler.Delete)
			bookReviews.GET("/user/:userId", bookReviewHandler.GetByUserID)
		}

		// バッジ
		badges := protected.Group("/badges")
		{
			badges.GET("/:userId", badgeHandler.GetUserBadges)
			badges.POST("/notify", badgeHandler.NotifyBadgeEarned)
		}

		// 学習ログ
		learningLogs := protected.Group("/learning-logs")
		{
			learningLogs.POST("", learningLogHandler.Create)
			learningLogs.GET("", learningLogHandler.GetMyLogs)
			learningLogs.GET("/user/:userId", learningLogHandler.GetByUserID)
			learningLogs.GET("/calendar/:userId", learningLogHandler.GetCalendarData)
			learningLogs.GET("/streak/:userId", learningLogHandler.GetStreakInfo)
			learningLogs.GET("/:id", learningLogHandler.GetByID)
			learningLogs.PUT("/:id", learningLogHandler.Update)
			learningLogs.DELETE("/:id", learningLogHandler.Delete)
		}

		// AIアドバイス
		advice := protected.Group("/advice")
		{
			advice.GET("", aiAdviceHandler.GetAdvice)
			advice.PUT("/:id/read", aiAdviceHandler.MarkAsRead)
			advice.POST("/chat", aiAdviceHandler.Chat)
			advice.GET("/conversations", aiAdviceHandler.GetConversations)
			advice.GET("/conversations/:id", aiAdviceHandler.GetConversation)
		}
	}

	return r
}

// getOrCreateSystemUser はテンプレート用のシステムユーザーを取得または作成する。
// "system@devsync.local" のメールアドレスで検索し、存在しなければ新規作成する。
func getOrCreateSystemUser(db *gorm.DB) (*model.User, error) {
	const systemEmail = "system@devsync.local"
	var user model.User
	err := db.Where("email = ?", systemEmail).First(&user).Error
	if err == nil {
		return &user, nil
	}
	// ユーザーが存在しない場合は作成
	// ユニーク制約（git_hub_id, git_hub_username）の衝突を回避するため専用値を設定
	user = model.User{
		Name:           "DevSync System",
		Email:          systemEmail,
		GitHubID:       -1,
		GitHubUsername:  "__system__",
	}
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

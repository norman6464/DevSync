package router

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/config"
	"github.com/norman6464/devsync/backend/internal/handler"
	"github.com/norman6464/devsync/backend/internal/middleware"
	"github.com/norman6464/devsync/backend/internal/repository"
	"github.com/norman6464/devsync/backend/internal/service"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config, hub *service.Hub) *gin.Engine {
	r := gin.Default()

	origins := strings.Split(cfg.CORSOrigins, ",")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Repositories
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

	// Services
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	githubService := service.NewGitHubService(cfg, userRepo, githubRepo)
	zennService := service.NewZennService()
	qiitaService := service.NewQiitaService()

	// Handlers
	authHandler := handler.NewAuthHandler(authService, githubService, userRepo, passwordResetRepo)
	userHandler := handler.NewUserHandler(userRepo)
	followHandler := handler.NewFollowHandler(followRepo)
	githubHandler := handler.NewGitHubHandler(githubService, authService, userRepo, githubRepo)
	postHandler := handler.NewPostHandler(postRepo, notificationRepo)
	rankingHandler := handler.NewRankingHandler(rankingRepo)
	messageHandler := handler.NewMessageHandler(messageRepo, notificationRepo)
	wsHandler := handler.NewWebSocketHandler(hub, authService)
	uploadHandler := handler.NewUploadHandler()
	notificationHandler := handler.NewNotificationHandler(notificationRepo)
	zennHandler := handler.NewZennHandler(zennRepo, userRepo, zennService)
	qiitaHandler := handler.NewQiitaHandler(qiitaRepo, userRepo, qiitaService)
	learningGoalHandler := handler.NewLearningGoalHandler(learningGoalRepo)
	activityReportHandler := handler.NewActivityReportHandler(activityReportRepo)
	projectHandler := handler.NewProjectHandler(projectRepo)
	learningResourceHandler := handler.NewLearningResourceHandler(learningResourceRepo)
	bookReviewHandler := handler.NewBookReviewHandler(bookReviewRepo)
	questionHandler := handler.NewQuestionHandler(questionRepo)
	answerHandler := handler.NewAnswerHandler(answerRepo, questionRepo)
	roadmapHandler := handler.NewRoadmapHandler(roadmapRepo)

	// Static file serving for uploads
	r.Static("/uploads", "./uploads")

	// Public routes
	r.GET("/health", handler.HealthCheck)

	// WebSocket (auth via query param)
	r.GET("/ws", wsHandler.HandleWebSocket)

	api := r.Group("/api/v1")

	// Auth routes (public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/github", authHandler.GitHubLogin)
		auth.GET("/github/callback", authHandler.GitHubLoginCallback)
		auth.POST("/password-reset/request", authHandler.RequestPasswordReset)
		auth.POST("/password-reset/confirm", authHandler.ResetPassword)
	}

	// GitHub data-connect callback (public - called by frontend after OAuth redirect)
	api.GET("/github/callback", githubHandler.Callback)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthRequired(authService))
	{
		// Auth
		protected.GET("/auth/me", authHandler.Me)
		protected.DELETE("/auth/account", authHandler.DeleteAccount)

		// Users
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

		// GitHub
		github := protected.Group("/github")
		{
			github.GET("/connect", githubHandler.Connect)
			github.POST("/sync", githubHandler.Sync)
			github.DELETE("/disconnect", githubHandler.Disconnect)
			github.GET("/contributions/:userId", githubHandler.GetContributions)
			github.GET("/languages/:userId", githubHandler.GetLanguages)
			github.GET("/repos/:userId", githubHandler.GetRepos)
		}

		// Posts
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
		}

		// Rankings
		rankings := protected.Group("/rankings")
		{
			rankings.GET("/contributions", rankingHandler.ContributionRanking)
			rankings.GET("/languages/:lang", rankingHandler.LanguageRanking)
			rankings.GET("/languages", rankingHandler.AvailableLanguages)
		}

		// Messages
		messages := protected.Group("/messages")
		{
			messages.GET("", messageHandler.GetConversations)
			messages.GET("/:userId", messageHandler.GetMessages)
			messages.POST("/:userId", messageHandler.SendMessage)
		}

		// Upload
		upload := protected.Group("/upload")
		{
			upload.POST("/image", uploadHandler.UploadImage)
			upload.POST("/images", uploadHandler.UploadMultipleImages)
		}

		// Notifications
		notifications := protected.Group("/notifications")
		{
			notifications.GET("", notificationHandler.GetAll)
			notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
			notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
			notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
			notifications.DELETE("/:id", notificationHandler.Delete)
		}

		// Zenn
		zenn := protected.Group("/zenn")
		{
			zenn.POST("/connect", zennHandler.Connect)
			zenn.DELETE("/disconnect", zennHandler.Disconnect)
			zenn.POST("/sync", zennHandler.Sync)
			zenn.GET("/articles/:userId", zennHandler.GetArticles)
			zenn.GET("/stats/:userId", zennHandler.GetStats)
		}

		// Qiita
		qiita := protected.Group("/qiita")
		{
			qiita.POST("/connect", qiitaHandler.Connect)
			qiita.DELETE("/disconnect", qiitaHandler.Disconnect)
			qiita.POST("/sync", qiitaHandler.Sync)
			qiita.GET("/articles/:userId", qiitaHandler.GetArticles)
			qiita.GET("/stats/:userId", qiitaHandler.GetStats)
		}

		// Learning Goals
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

		// Activity Reports
		reports := protected.Group("/reports")
		{
			reports.GET("/weekly", activityReportHandler.GetMyWeeklyReport)
			reports.GET("/monthly", activityReportHandler.GetMyMonthlyReport)
			reports.GET("/weekly/:userId", activityReportHandler.GetWeeklyReport)
			reports.GET("/monthly/:userId", activityReportHandler.GetMonthlyReport)
			reports.GET("/comparison", activityReportHandler.GetComparison)
		}

		// Projects (Showcase)
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

		// Learning Resources
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

		// Questions (Q&A)
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

			// Answers (nested under questions)
			questions.GET("/:id/answers", answerHandler.GetByQuestionID)
			questions.POST("/:id/answers", answerHandler.Create)
			questions.PUT("/:id/answers/:answerId", answerHandler.Update)
			questions.DELETE("/:id/answers/:answerId", answerHandler.Delete)
			questions.PUT("/:id/answers/:answerId/best", answerHandler.SetBestAnswer)
			questions.POST("/:id/answers/:answerId/vote", answerHandler.Vote)
			questions.DELETE("/:id/answers/:answerId/vote", answerHandler.RemoveVote)
		}

		// Learning Roadmaps
		roadmaps := protected.Group("/roadmaps")
		{
			roadmaps.POST("", roadmapHandler.Create)
			roadmaps.GET("", roadmapHandler.GetMyRoadmaps)
			roadmaps.GET("/public", roadmapHandler.GetPublicRoadmaps)
			roadmaps.GET("/:id", roadmapHandler.GetByID)
			roadmaps.PUT("/:id", roadmapHandler.Update)
			roadmaps.DELETE("/:id", roadmapHandler.Delete)
			roadmaps.POST("/:id/copy", roadmapHandler.CopyRoadmap)

			// Steps
			roadmaps.POST("/:id/steps", roadmapHandler.CreateStep)
			roadmaps.PUT("/:id/steps/:stepId", roadmapHandler.UpdateStep)
			roadmaps.DELETE("/:id/steps/:stepId", roadmapHandler.DeleteStep)
			roadmaps.PUT("/:id/steps/reorder", roadmapHandler.ReorderSteps)
		}

		// Book Reviews
		bookReviews := protected.Group("/book-reviews")
		{
			bookReviews.POST("", bookReviewHandler.Create)
			bookReviews.GET("", bookReviewHandler.GetAll)
			bookReviews.GET("/:id", bookReviewHandler.GetByID)
			bookReviews.PUT("/:id", bookReviewHandler.Update)
			bookReviews.DELETE("/:id", bookReviewHandler.Delete)
			bookReviews.GET("/user/:userId", bookReviewHandler.GetByUserID)
		}
	}

	return r
}

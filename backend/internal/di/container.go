// Package di はDevSyncアプリケーションの依存性注入コンテナを提供する。
// リポジトリ→サービス→ハンドラの依存関係を構築し、ルーターに公開する。
package di

import (
	"log"

	"github.com/norman6464/devsync/backend/internal/config"
	"github.com/norman6464/devsync/backend/internal/handler"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
	"github.com/norman6464/devsync/backend/internal/service"
	"gorm.io/gorm"
)

// Container はDI（依存性注入）コンテナ。
// 全ハンドラとルーティングに必要な公開フィールドを保持する。
type Container struct {
	// ハンドラ
	AuthHandler              *handler.AuthHandler
	UserHandler              *handler.UserHandler
	FollowHandler            *handler.FollowHandler
	GitHubHandler            *handler.GitHubHandler
	PostHandler              *handler.PostHandler
	CodeSnippetHandler       *handler.CodeSnippetHandler
	RankingHandler           *handler.RankingHandler
	MessageHandler           *handler.MessageHandler
	WebSocketHandler         *handler.WebSocketHandler
	UploadHandler            *handler.UploadHandler
	NotificationHandler      *handler.NotificationHandler
	ZennHandler              *handler.ZennHandler
	QiitaHandler             *handler.QiitaHandler
	LearningGoalHandler      *handler.LearningGoalHandler
	ActivityReportHandler    *handler.ActivityReportHandler
	ProjectHandler           *handler.ProjectHandler
	LearningResourceHandler  *handler.LearningResourceHandler
	BookReviewHandler        *handler.BookReviewHandler
	QuestionHandler          *handler.QuestionHandler
	AnswerHandler            *handler.AnswerHandler
	RoadmapHandler           *handler.RoadmapHandler
	ChatRoomHandler          *handler.ChatRoomHandler
	AtCoderHandler           *handler.AtCoderHandler
	BadgeHandler             *handler.BadgeHandler
	LearningLogHandler       *handler.LearningLogHandler
	AIAdviceHandler          *handler.AIAdviceHandler
	EmailPreferencesHandler  *handler.EmailPreferencesHandler
	LevelHandler             *handler.LevelHandler
	LearningAnalyticsHandler *handler.LearningAnalyticsHandler
	RecommendationHandler    *handler.RecommendationHandler
	StudyCircleHandler       *handler.StudyCircleHandler
	SearchHandler            *handler.SearchHandler
	NoteHandler              *handler.NoteHandler
	NoteFolderHandler        *handler.NoteFolderHandler
	NoteTemplateHandler      *handler.NoteTemplateHandler
	NoteLinkHandler          *handler.NoteLinkHandler
	PostSeriesHandler        *handler.PostSeriesHandler
	PostCollectionHandler    *handler.PostCollectionHandler
	PostTagHandler           *handler.PostTagHandler
	PostPinHandler           *handler.PostPinHandler
	PostViewHandler          *handler.PostViewHandler
	CommentLikeHandler       *handler.CommentLikeHandler
	MentionHandler           *handler.MentionHandler
	UserDashboardHandler     *handler.UserDashboardHandler
	NoteStatsHandler              *handler.NoteStatsHandler
	StudyCircleStatsHandler       *handler.StudyCircleStatsHandler
	PostStatsHandler              *handler.PostStatsHandler
	BookReviewStatsHandler        *handler.BookReviewStatsHandler
	QAStatsHandler                *handler.QAStatsHandler
	CodeSnippetStatsHandler       *handler.CodeSnippetStatsHandler
	LearningResourceStatsHandler  *handler.LearningResourceStatsHandler
	ProjectStatsHandler           *handler.ProjectStatsHandler
	FollowStatsHandler            *handler.FollowStatsHandler
	RoadmapStatsHandler           *handler.RoadmapStatsHandler
	LearningLogStatsHandler       *handler.LearningLogStatsHandler
	CommentStatsHandler           *handler.CommentStatsHandler
	NotificationStatsHandler      *handler.NotificationStatsHandler
	MessageStatsHandler           *handler.MessageStatsHandler
	YouTubeHandler                *handler.YouTubeHandler
	SpotifyHandler           *handler.SpotifyHandler

	// ミドルウェア・コールバック用
	AuthService      *service.AuthService
	Hub              *service.Hub
	GroupMessageRepo *repository.GroupMessageRepository
}

// NewContainer はDIコンテナを構築する。
// リポジトリ→サービス→ハンドラの順で依存関係を解決する。
func NewContainer(db *gorm.DB, cfg *config.Config, hub *service.Hub) *Container {
	c := &Container{Hub: hub}

	// リポジトリ
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
	levelRepo := repository.NewLevelRepository(db)
	analyticsRepo := repository.NewLearningAnalyticsRepository(db)
	badgeRepo := repository.NewBadgeRepository(db)
	recommendationRepo := repository.NewRecommendationRepository(db)
	studyCircleRepo := repository.NewStudyCircleRepository(db)
	noteRepo := repository.NewNoteRepository(db)
	noteFolderRepo := repository.NewNoteFolderRepository(db)
	noteTemplateRepo := repository.NewNoteTemplateRepository(db)
	noteLinkRepo := repository.NewNoteLinkRepository(db)
	postSeriesRepo := repository.NewPostSeriesRepository(db)

	c.GroupMessageRepo = groupMessageRepo

	// 共通サービス
	authService := service.NewAuthService(userRepo, passwordResetRepo, cfg.JWTSecret)
	c.AuthService = authService

	notificationService := service.NewNotificationService(notificationRepo)
	userService := service.NewUserService(userRepo)

	// ドメインサービス
	githubService := service.NewGitHubService(cfg, userRepo, githubRepo)
	zennService := service.NewZennService(userRepo, zennRepo)
	qiitaService := service.NewQiitaService(userRepo, qiitaRepo)
	postService := service.NewPostService(postRepo, notificationService)
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
	chatRoomService := service.NewChatRoomService(chatRoomRepo, groupMessageRepo, hub)
	activityReportService := service.NewActivityReportService(activityReportRepo)
	atcoderService := service.NewAtCoderService()
	badgeService := service.NewBadgeService(badgeRepo, notificationService)
	levelService := service.NewLevelService(levelRepo, notificationService)
	analyticsService := service.NewLearningAnalyticsService(analyticsRepo)
	recommendationService := service.NewRecommendationService(recommendationRepo, userRepo)
	studyCircleService := service.NewStudyCircleService(studyCircleRepo)
	noteService := service.NewNoteService(noteRepo)
	noteFolderService := service.NewNoteFolderService(noteFolderRepo)
	noteTemplateService := service.NewNoteTemplateService(noteTemplateRepo, noteService)
	noteLinkService := service.NewNoteLinkService(noteLinkRepo, noteRepo)
	postSeriesService := service.NewPostSeriesService(postSeriesRepo)

	// テンプレートロードマップの初期登録
	go seedTemplateRoadmaps(db, roadmapService)

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

	// メールサービスの初期化（SMTP設定がある場合のみ有効化）
	var emailSender service.EmailSenderInterface
	if cfg.SMTPHost != "" {
		emailSender = service.NewSMTPEmailSender(cfg)
		log.Println("SMTP設定が検出されました。メール機能が有効です。")
	}
	weeklyReportEmailService := service.NewWeeklyReportEmailService(emailSender, activityReportService, userRepo)
	weeklyReportEmailService.SetAppURL(cfg.AppURL)

	// スケジューラの初期化
	if cfg.SMTPHost != "" {
		scheduler := service.NewScheduler(weeklyReportEmailService)
		go scheduler.Start()
	} else {
		log.Println("SMTP未設定。ウィークリーレポートメールのスケジューラは無効です。")
	}

	// ハンドラ
	origins := cfg.CORSOrigins
	c.AuthHandler = handler.NewAuthHandler(authService, githubService)
	c.UserHandler = handler.NewUserHandler(userService)
	c.FollowHandler = handler.NewFollowHandler(followService)
	c.GitHubHandler = handler.NewGitHubHandler(githubService, authService)
	c.PostHandler = handler.NewPostHandler(postService, codeSnippetService)
	c.CodeSnippetHandler = handler.NewCodeSnippetHandler(codeSnippetService)
	c.RankingHandler = handler.NewRankingHandler(rankingService)
	c.MessageHandler = handler.NewMessageHandler(messageService)
	c.WebSocketHandler = handler.NewWebSocketHandler(hub, authService, parseOrigins(origins))
	c.UploadHandler = handler.NewUploadHandler()
	c.NotificationHandler = handler.NewNotificationHandler(notificationService)
	c.ZennHandler = handler.NewZennHandler(zennService)
	c.QiitaHandler = handler.NewQiitaHandler(qiitaService)
	c.LearningGoalHandler = handler.NewLearningGoalHandler(learningGoalService)
	c.ActivityReportHandler = handler.NewActivityReportHandler(activityReportService)
	c.ProjectHandler = handler.NewProjectHandler(projectService)
	c.LearningResourceHandler = handler.NewLearningResourceHandler(learningResourceService)
	c.BookReviewHandler = handler.NewBookReviewHandler(bookReviewService)
	c.QuestionHandler = handler.NewQuestionHandler(questionService)
	c.AnswerHandler = handler.NewAnswerHandler(answerService)
	c.RoadmapHandler = handler.NewRoadmapHandler(roadmapService)
	c.ChatRoomHandler = handler.NewChatRoomHandler(chatRoomService)
	c.AtCoderHandler = handler.NewAtCoderHandler(atcoderService, userService)
	c.BadgeHandler = handler.NewBadgeHandler(badgeService)
	c.LearningLogHandler = handler.NewLearningLogHandler(learningLogService)
	c.AIAdviceHandler = handler.NewAIAdviceHandler(aiAdviceService)
	c.EmailPreferencesHandler = handler.NewEmailPreferencesHandler(userService)
	c.LevelHandler = handler.NewLevelHandler(levelService)
	c.LearningAnalyticsHandler = handler.NewLearningAnalyticsHandler(analyticsService)
	c.RecommendationHandler = handler.NewRecommendationHandler(recommendationService)
	c.StudyCircleHandler = handler.NewStudyCircleHandler(studyCircleService)
	searchService := service.NewSearchService(postRepo)
	c.SearchHandler = handler.NewSearchHandler(searchService, studyCircleService)
	c.NoteHandler = handler.NewNoteHandler(noteService)
	c.NoteFolderHandler = handler.NewNoteFolderHandler(noteFolderService)
	c.NoteTemplateHandler = handler.NewNoteTemplateHandler(noteTemplateService)
	c.NoteLinkHandler = handler.NewNoteLinkHandler(noteLinkService)
	c.PostSeriesHandler = handler.NewPostSeriesHandler(postSeriesService)

	postCollectionRepo := repository.NewPostCollectionRepository(db)
	postCollectionService := service.NewPostCollectionService(postCollectionRepo)
	c.PostCollectionHandler = handler.NewPostCollectionHandler(postCollectionService)

	postTagRepo := repository.NewPostTagRepository(db)
	postTagService := service.NewPostTagService(postTagRepo, postRepo)
	c.PostTagHandler = handler.NewPostTagHandler(postTagService)

	postPinRepo := repository.NewPostPinRepository(db)
	postPinService := service.NewPostPinService(postPinRepo, postRepo)
	c.PostPinHandler = handler.NewPostPinHandler(postPinService)

	commentLikeRepo := repository.NewCommentLikeRepository(db)
	commentLikeService := service.NewCommentLikeService(commentLikeRepo, postRepo)
	c.CommentLikeHandler = handler.NewCommentLikeHandler(commentLikeService)

	postViewRepo := repository.NewPostViewRepository(db)
	postViewService := service.NewPostViewService(postViewRepo)
	c.PostViewHandler = handler.NewPostViewHandler(postViewService)

	mentionRepo := repository.NewMentionRepository(db)
	mentionService := service.NewMentionService(mentionRepo, userRepo, notificationService)
	c.MentionHandler = handler.NewMentionHandler(mentionService)

	userDashboardRepo := repository.NewUserDashboardRepository(db)
	userDashboardService := service.NewUserDashboardService(userDashboardRepo)
	c.UserDashboardHandler = handler.NewUserDashboardHandler(userDashboardService)

	noteStatsRepo := repository.NewNoteStatsRepository(db)
	noteStatsService := service.NewNoteStatsService(noteStatsRepo)
	c.NoteStatsHandler = handler.NewNoteStatsHandler(noteStatsService)

	studyCircleStatsRepo := repository.NewStudyCircleStatsRepository(db)
	studyCircleStatsService := service.NewStudyCircleStatsService(studyCircleStatsRepo)
	c.StudyCircleStatsHandler = handler.NewStudyCircleStatsHandler(studyCircleStatsService)

	postStatsRepo := repository.NewPostStatsRepository(db)
	postStatsService := service.NewPostStatsService(postStatsRepo)
	c.PostStatsHandler = handler.NewPostStatsHandler(postStatsService)

	bookReviewStatsRepo := repository.NewBookReviewStatsRepository(db)
	bookReviewStatsService := service.NewBookReviewStatsService(bookReviewStatsRepo)
	c.BookReviewStatsHandler = handler.NewBookReviewStatsHandler(bookReviewStatsService)

	qaStatsRepo := repository.NewQAStatsRepository(db)
	qaStatsService := service.NewQAStatsService(qaStatsRepo)
	c.QAStatsHandler = handler.NewQAStatsHandler(qaStatsService)

	codeSnippetStatsRepo := repository.NewCodeSnippetStatsRepository(db)
	codeSnippetStatsService := service.NewCodeSnippetStatsService(codeSnippetStatsRepo)
	c.CodeSnippetStatsHandler = handler.NewCodeSnippetStatsHandler(codeSnippetStatsService)

	learningResourceStatsRepo := repository.NewLearningResourceStatsRepository(db)
	learningResourceStatsService := service.NewLearningResourceStatsService(learningResourceStatsRepo)
	c.LearningResourceStatsHandler = handler.NewLearningResourceStatsHandler(learningResourceStatsService)

	projectStatsRepo := repository.NewProjectStatsRepository(db)
	projectStatsService := service.NewProjectStatsService(projectStatsRepo)
	c.ProjectStatsHandler = handler.NewProjectStatsHandler(projectStatsService)

	followStatsRepo := repository.NewFollowStatsRepository(db)
	followStatsService := service.NewFollowStatsService(followStatsRepo)
	c.FollowStatsHandler = handler.NewFollowStatsHandler(followStatsService)

	roadmapStatsRepo := repository.NewRoadmapStatsRepository(db)
	roadmapStatsService := service.NewRoadmapStatsService(roadmapStatsRepo)
	c.RoadmapStatsHandler = handler.NewRoadmapStatsHandler(roadmapStatsService)

	learningLogStatsRepo := repository.NewLearningLogStatsRepository(db)
	learningLogStatsService := service.NewLearningLogStatsService(learningLogStatsRepo)
	c.LearningLogStatsHandler = handler.NewLearningLogStatsHandler(learningLogStatsService)

	commentStatsRepo := repository.NewCommentStatsRepository(db)
	commentStatsService := service.NewCommentStatsService(commentStatsRepo)
	c.CommentStatsHandler = handler.NewCommentStatsHandler(commentStatsService)

	notificationStatsRepo := repository.NewNotificationStatsRepository(db)
	notificationStatsService := service.NewNotificationStatsService(notificationStatsRepo)
	c.NotificationStatsHandler = handler.NewNotificationStatsHandler(notificationStatsService)

	messageStatsRepo := repository.NewMessageStatsRepository(db)
	messageStatsService := service.NewMessageStatsService(messageStatsRepo)
	c.MessageStatsHandler = handler.NewMessageStatsHandler(messageStatsService)

	// Spotifyサービス
	spotifyRepo := repository.NewSpotifyRepository(db)
	spotifyService := service.NewSpotifyService(cfg, userRepo, spotifyRepo)
	c.SpotifyHandler = handler.NewSpotifyHandler(spotifyService, authService)

	// YouTubeサービス（APIキー設定時のみクライアント初期化）
	youtubeVideoRepo := repository.NewYouTubeVideoRepository(db)
	var youtubeClient service.YouTubeClientInterface
	if cfg.YouTubeAPIKey != "" {
		youtubeClient = service.NewYouTubeClient(cfg.YouTubeAPIKey)
		log.Println("YouTube APIキーが設定されています。YouTube動画検索機能が有効です。")
	} else {
		log.Println("YouTube APIキー未設定。YouTube動画検索機能は無効です。")
	}
	youtubeService := service.NewYouTubeService(youtubeVideoRepo, userRepo, youtubeClient)
	c.YouTubeHandler = handler.NewYouTubeHandler(youtubeService)

	// HubのGetRoomMembersコールバックを設定
	hub.GetRoomMembers = groupMessageRepo.GetMemberUserIDs

	return c
}

// parseOrigins はカンマ区切りのCORSオリジン文字列をスライスに変換する。
func parseOrigins(origins string) []string {
	var result []string
	for _, o := range splitAndTrim(origins) {
		if o != "" {
			result = append(result, o)
		}
	}
	return result
}

// splitAndTrim はカンマ区切り文字列を分割してトリムする。
func splitAndTrim(s string) []string {
	parts := make([]string, 0)
	for _, p := range splitByComma(s) {
		trimmed := trimSpace(p)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}

func splitByComma(s string) []string {
	result := []string{""}
	for _, c := range s {
		if c == ',' {
			result = append(result, "")
		} else {
			result[len(result)-1] += string(c)
		}
	}
	return result
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

// seedTemplateRoadmaps はシステムユーザーを取得/作成し、テンプレートロードマップを登録する。
func seedTemplateRoadmaps(db *gorm.DB, roadmapService *service.RoadmapService) {
	const systemEmail = "system@devsync.local"
	var user model.User
	err := db.Where("email = ?", systemEmail).First(&user).Error
	if err != nil {
		user = model.User{
			Name:           "DevSync System",
			Email:          systemEmail,
			Username:       "__system__",
			GitHubID:       -1,
			GitHubUsername: "__system__",
		}
		if err := db.Create(&user).Error; err != nil {
			log.Printf("テンプレートシード用システムユーザー作成失敗: %v", err)
			return
		}
	}
	if err := roadmapService.SeedTemplates(user.ID); err != nil {
		log.Printf("テンプレートシード失敗: %v", err)
	}
}

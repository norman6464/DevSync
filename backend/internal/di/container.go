// Package di はDevSyncアプリケーションの依存性注入コンテナを提供する。
// リポジトリ→サービス→ハンドラの依存関係を構築し、ルーターに公開する。
package di

import (
	"log"

	"github.com/norman6464/devsync/backend/internal/adapter/persistence"
	"github.com/norman6464/devsync/backend/internal/config"
	"github.com/norman6464/devsync/backend/internal/handler"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/norman6464/devsync/backend/internal/usecase"
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
	ZennHandler              *handler.ArticlePlatformHandler[model.ZennArticle, model.ZennStats]
	QiitaHandler             *handler.ArticlePlatformHandler[model.QiitaArticle, model.QiitaStats]
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
	MentionStatsHandler           *handler.MentionStatsHandler
	ReactionStatsHandler          *handler.ReactionStatsHandler
	BookmarkStatsHandler          *handler.BookmarkStatsHandler
	YouTubeHandler                *handler.YouTubeHandler
	SpotifyHandler           *handler.SpotifyHandler
	StreakFreezeHandler            *handler.StreakFreezeHandler
	BookmarkCollectionHandler     *handler.BookmarkCollectionHandler
	WeeklyChallengeHandler        *handler.WeeklyChallengeHandler
	PostTemplateHandler           *handler.PostTemplateHandler
	WidgetSettingsHandler         *handler.WidgetSettingsHandler
	WeeklyGoalHandler            *handler.WeeklyGoalHandler
	NoteVersionHandler           *handler.NoteVersionHandler
	ResourceProgressHandler      *handler.ResourceProgressHandler
	ProjectMilestoneHandler      *handler.ProjectMilestoneHandler
	UserActivityHandler          *handler.UserActivityHandler
	LearningLogTemplateHandler   *handler.LearningLogTemplateHandler
	ResourceReviewHandler        *handler.ResourceReviewHandler
	LearningDashboardHandler     *handler.LearningDashboardHandler
	ReminderSettingsHandler      *handler.ReminderSettingsHandler
	NotificationSettingsHandler  *handler.NotificationSettingsHandler

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
	// follow はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	followRepo := persistence.NewFollowRepository(db)
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
	learningLogTemplateRepo := repository.NewLearningLogTemplateRepository(db)
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
	questionService := service.NewQuestionService(questionRepo)
	answerService := service.NewAnswerService(answerRepo, questionRepo)
	learningLogService := service.NewLearningLogService(learningLogRepo, learningGoalRepo)
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
	atcoderService := service.NewAtCoderService(userRepo)
	badgeService := service.NewBadgeService(badgeRepo, notificationService)
	levelService := service.NewLevelService(levelRepo, notificationService)
	analyticsService := service.NewLearningAnalyticsService(analyticsRepo)
	recommendationService := service.NewRecommendationService(recommendationRepo, userRepo)
	studyCircleService := service.NewStudyCircleService(studyCircleRepo)
	noteService := service.NewNoteService(noteRepo)
	noteFolderService := service.NewNoteFolderService(noteFolderRepo)
	noteTemplateService := service.NewNoteTemplateService(noteTemplateRepo, noteService)
	learningLogTemplateService := service.NewLearningLogTemplateService(learningLogTemplateRepo, learningLogService)
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
	c.FollowHandler = handler.NewFollowHandler(
		usecase.NewFollowUserUseCase(followRepo),
		usecase.NewUnfollowUserUseCase(followRepo),
		usecase.NewListFollowersUseCase(followRepo),
		usecase.NewListFollowingUseCase(followRepo),
	)
	c.GitHubHandler = handler.NewGitHubHandler(githubService, authService)
	c.PostHandler = handler.NewPostHandler(postService, codeSnippetService)
	c.CodeSnippetHandler = handler.NewCodeSnippetHandler(codeSnippetService)
	c.RankingHandler = handler.NewRankingHandler(rankingService)
	c.MessageHandler = handler.NewMessageHandler(messageService)
	c.WebSocketHandler = handler.NewWebSocketHandler(hub, authService, parseOrigins(origins))
	uploadHandler, err := handler.NewUploadHandler()
	if err != nil {
		log.Fatalf("アップロードハンドラの初期化に失敗: %v", err)
	}
	c.UploadHandler = uploadHandler
	c.NotificationHandler = handler.NewNotificationHandler(notificationService)
	c.ZennHandler = handler.NewArticlePlatformHandler[model.ZennArticle, model.ZennStats](zennService, "Zenn")
	c.QiitaHandler = handler.NewArticlePlatformHandler[model.QiitaArticle, model.QiitaStats](qiitaService, "Qiita")
	c.LearningGoalHandler = handler.NewLearningGoalHandler(learningGoalService)
	c.ActivityReportHandler = handler.NewActivityReportHandler(activityReportService)
	c.ProjectHandler = handler.NewProjectHandler(projectService)
	c.LearningResourceHandler = handler.NewLearningResourceHandler(learningResourceService)
	c.BookReviewHandler = handler.NewBookReviewHandler(bookReviewService)
	c.QuestionHandler = handler.NewQuestionHandler(questionService)
	c.AnswerHandler = handler.NewAnswerHandler(answerService)
	c.RoadmapHandler = handler.NewRoadmapHandler(roadmapService)
	c.ChatRoomHandler = handler.NewChatRoomHandler(chatRoomService)
	c.AtCoderHandler = handler.NewAtCoderHandler(atcoderService)
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
	noteVersionRepo := repository.NewNoteVersionRepository(db)
	noteVersionService := service.NewNoteVersionService(noteRepo, noteVersionRepo)
	c.NoteVersionHandler = handler.NewNoteVersionHandler(noteVersionService)
	c.NoteFolderHandler = handler.NewNoteFolderHandler(noteFolderService)
	c.NoteTemplateHandler = handler.NewNoteTemplateHandler(noteTemplateService)
	c.LearningLogTemplateHandler = handler.NewLearningLogTemplateHandler(learningLogTemplateService)
	c.NoteLinkHandler = handler.NewNoteLinkHandler(noteLinkService)
	c.PostSeriesHandler = handler.NewPostSeriesHandler(postSeriesService)

	postCollectionRepo := repository.NewPostCollectionRepository(db)
	postCollectionService := service.NewPostCollectionService(postCollectionRepo)
	c.PostCollectionHandler = handler.NewPostCollectionHandler(postCollectionService)

	postTagRepo := repository.NewPostTagRepository(db)
	postTagService := service.NewPostTagService(postTagRepo, postRepo)
	c.PostTagHandler = handler.NewPostTagHandler(postTagService)
	c.PostHandler.SetTagService(postTagService)

	// post_pin はクリーンアーキテクチャ(DIP)へ移行済み。
	postPinRepo := persistence.NewPostPinRepository(db)
	postReader := persistence.NewPostReader(db)
	c.PostPinHandler = handler.NewPostPinHandler(
		usecase.NewPinPostUseCase(postPinRepo, postReader),
		usecase.NewUnpinPostUseCase(postPinRepo),
		usecase.NewListPinnedPostsUseCase(postPinRepo),
		usecase.NewCountPinnedPostsUseCase(postPinRepo),
		usecase.NewReorderPinnedPostsUseCase(postPinRepo),
	)

	// comment_like はクリーンアーキテクチャ(DIP)へ移行済み。
	commentLikeRepo := persistence.NewCommentLikeRepository(db)
	commentReader := persistence.NewCommentReader(db)
	c.CommentLikeHandler = handler.NewCommentLikeHandler(
		usecase.NewLikeCommentUseCase(commentLikeRepo, commentReader),
		usecase.NewUnlikeCommentUseCase(commentLikeRepo, commentReader),
		usecase.NewGetCommentLikeStatusUseCase(commentLikeRepo, commentReader),
	)

	// 投稿閲覧数はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	postViewRepo := persistence.NewPostViewRepository(db)
	c.PostViewHandler = handler.NewPostViewHandler(
		usecase.NewRecordPostViewUseCase(postViewRepo),
		usecase.NewGetPostViewCountUseCase(postViewRepo),
		usecase.NewGetMostViewedPostsUseCase(postViewRepo),
	)

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

	// 投稿統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	postStatsRepo := persistence.NewPostStatsRepository(db)
	c.PostStatsHandler = handler.NewPostStatsHandler(
		usecase.NewGetPostStatsUseCase(postStatsRepo),
	)

	// 書籍レビュー統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	bookReviewStatsRepo := persistence.NewBookReviewStatsRepository(db)
	c.BookReviewStatsHandler = handler.NewBookReviewStatsHandler(
		usecase.NewGetBookReviewStatsUseCase(bookReviewStatsRepo),
	)

	qaStatsRepo := repository.NewQAStatsRepository(db)
	qaStatsService := service.NewQAStatsService(qaStatsRepo)
	c.QAStatsHandler = handler.NewQAStatsHandler(qaStatsService)

	// コードスニペット統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	codeSnippetStatsRepo := persistence.NewCodeSnippetStatsRepository(db)
	c.CodeSnippetStatsHandler = handler.NewCodeSnippetStatsHandler(
		usecase.NewGetCodeSnippetStatsUseCase(codeSnippetStatsRepo),
	)

	// 学習リソース統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	learningResourceStatsRepo := persistence.NewLearningResourceStatsRepository(db)
	c.LearningResourceStatsHandler = handler.NewLearningResourceStatsHandler(
		usecase.NewGetLearningResourceStatsUseCase(learningResourceStatsRepo),
	)

	projectStatsRepo := repository.NewProjectStatsRepository(db)
	projectStatsService := service.NewProjectStatsService(projectStatsRepo)
	c.ProjectStatsHandler = handler.NewProjectStatsHandler(projectStatsService)

	// フォロー統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	followStatsRepo := persistence.NewFollowStatsRepository(db)
	c.FollowStatsHandler = handler.NewFollowStatsHandler(
		usecase.NewGetFollowStatsUseCase(followStatsRepo),
	)

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

	mentionStatsRepo := repository.NewMentionStatsRepository(db)
	mentionStatsService := service.NewMentionStatsService(mentionStatsRepo)
	c.MentionStatsHandler = handler.NewMentionStatsHandler(mentionStatsService)

	reactionStatsRepo := repository.NewReactionStatsRepository(db)
	reactionStatsService := service.NewReactionStatsService(reactionStatsRepo)
	c.ReactionStatsHandler = handler.NewReactionStatsHandler(reactionStatsService)

	bookmarkStatsRepo := repository.NewBookmarkStatsRepository(db)
	bookmarkStatsService := service.NewBookmarkStatsService(bookmarkStatsRepo)
	c.BookmarkStatsHandler = handler.NewBookmarkStatsHandler(bookmarkStatsService)

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

	// ストリークフリーズサービス
	streakFreezeRepo := repository.NewStreakFreezeRepository(db)
	streakFreezeService := service.NewStreakFreezeService(streakFreezeRepo)
	c.StreakFreezeHandler = handler.NewStreakFreezeHandler(streakFreezeService)

	// ブックマークコレクションサービス
	bookmarkCollectionRepo := repository.NewBookmarkCollectionRepository(db)
	bookmarkCollectionService := service.NewBookmarkCollectionService(bookmarkCollectionRepo)
	c.BookmarkCollectionHandler = handler.NewBookmarkCollectionHandler(bookmarkCollectionService)

	// ウィークリーチャレンジサービス
	weeklyChallengeRepo := repository.NewWeeklyChallengeRepository(db)
	weeklyChallengeService := service.NewWeeklyChallengeService(weeklyChallengeRepo)
	c.WeeklyChallengeHandler = handler.NewWeeklyChallengeHandler(weeklyChallengeService)

	// 投稿テンプレートサービス
	postTemplateRepo := repository.NewPostTemplateRepository(db)
	postTemplateService := service.NewPostTemplateService(postTemplateRepo)
	c.PostTemplateHandler = handler.NewPostTemplateHandler(postTemplateService)

	// ウィジェット設定サービス
	widgetSettingsRepo := repository.NewWidgetSettingsRepository(db)
	widgetSettingsService := service.NewWidgetSettingsService(widgetSettingsRepo)
	c.WidgetSettingsHandler = handler.NewWidgetSettingsHandler(widgetSettingsService)

	// カテゴリ別週間学習目標はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	weeklyGoalRepo := persistence.NewWeeklyGoalRepository(db)
	c.WeeklyGoalHandler = handler.NewWeeklyGoalHandler(
		usecase.NewSetWeeklyGoalUseCase(weeklyGoalRepo),
		usecase.NewListWeeklyGoalsUseCase(weeklyGoalRepo),
		usecase.NewGetWeeklyGoalProgressUseCase(weeklyGoalRepo),
	)

	// リソース進捗サービス
	// リソース進捗はクリーンアーキテクチャ（DIP）へ移行済み。リソース存在確認は最小 port LearningResourceReader を再利用。
	resourceProgressRepo := persistence.NewResourceProgressRepository(db)
	resourceProgressResourceReader := persistence.NewLearningResourceReader(db)
	c.ResourceProgressHandler = handler.NewResourceProgressHandler(
		usecase.NewUpsertResourceProgressUseCase(resourceProgressRepo, resourceProgressResourceReader),
		usecase.NewGetResourceProgressUseCase(resourceProgressRepo),
		usecase.NewListResourceProgressUseCase(resourceProgressRepo),
	)

	// プロジェクトマイルストーンサービス
	projectMilestoneRepo := repository.NewProjectMilestoneRepository(db)
	projectMilestoneService := service.NewProjectMilestoneService(projectMilestoneRepo, projectRepo)
	c.ProjectMilestoneHandler = handler.NewProjectMilestoneHandler(projectMilestoneService)

	// ユーザーアクティビティサービス
	// ユーザーアクティビティはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	userActivityRepo := persistence.NewUserActivityRepository(db)
	c.UserActivityHandler = handler.NewUserActivityHandler(
		usecase.NewGetActivityTimelineUseCase(userActivityRepo),
	)

	// リソースレビューはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	resourceReviewRepo := persistence.NewResourceReviewRepository(db)
	learningResourceReader := persistence.NewLearningResourceReader(db)
	c.ResourceReviewHandler = handler.NewResourceReviewHandler(
		usecase.NewCreateResourceReviewUseCase(resourceReviewRepo, learningResourceReader),
		usecase.NewListResourceReviewsUseCase(resourceReviewRepo),
		usecase.NewUpdateResourceReviewUseCase(resourceReviewRepo),
		usecase.NewDeleteResourceReviewUseCase(resourceReviewRepo),
	)

	// 学習ダッシュボード統合サマリーサービス
	learningDashboardService := service.NewLearningDashboardService(learningLogRepo, learningGoalRepo, analyticsRepo)
	c.LearningDashboardHandler = handler.NewLearningDashboardHandler(learningDashboardService)

	// リマインダー設定サービス
	reminderSettingsRepo := repository.NewReminderSettingsRepository(db)
	reminderSettingsService := service.NewReminderSettingsService(reminderSettingsRepo)
	c.ReminderSettingsHandler = handler.NewReminderSettingsHandler(reminderSettingsService)

	// 通知設定サービス
	notificationSettingsRepo := repository.NewNotificationSettingsRepository(db)
	notificationSettingsService := service.NewNotificationSettingsService(notificationSettingsRepo)
	c.NotificationSettingsHandler = handler.NewNotificationSettingsHandler(notificationSettingsService)

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

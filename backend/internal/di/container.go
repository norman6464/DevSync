// Package di はDevSyncアプリケーションの依存性注入コンテナを提供する。
// リポジトリ→サービス→ハンドラの依存関係を構築し、ルーターに公開する。
package di

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/external"
	"github.com/norman6464/devsync/backend/internal/adapter/notify"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/config"
	"github.com/norman6464/devsync/backend/internal/handler"
	"github.com/norman6464/devsync/backend/internal/infra/scheduler"
	"github.com/norman6464/devsync/backend/internal/infra/ws"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	usecaserepo "github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// Container はDI（依存性注入）コンテナ。
// 全ハンドラとルーティングに必要な公開フィールドを保持する。
type Container struct {
	// ハンドラ
	AuthHandler                  *handler.AuthHandler
	UserHandler                  *handler.UserHandler
	FollowHandler                *handler.FollowHandler
	GitHubHandler                *handler.GitHubHandler
	PostHandler                  *handler.PostHandler
	CodeSnippetHandler           *handler.CodeSnippetHandler
	RankingHandler               *handler.RankingHandler
	MessageHandler               *handler.MessageHandler
	WebSocketHandler             *handler.WebSocketHandler
	UploadHandler                *handler.UploadHandler
	NotificationHandler          *handler.NotificationHandler
	ZennHandler                  *handler.ArticlePlatformHandler[model.ZennArticle, model.ZennStats]
	QiitaHandler                 *handler.ArticlePlatformHandler[model.QiitaArticle, model.QiitaStats]
	LearningGoalHandler          *handler.LearningGoalHandler
	ActivityReportHandler        *handler.ActivityReportHandler
	ProjectHandler               *handler.ProjectHandler
	LearningResourceHandler      *handler.LearningResourceHandler
	BookReviewHandler            *handler.BookReviewHandler
	QuestionHandler              *handler.QuestionHandler
	AnswerHandler                *handler.AnswerHandler
	RoadmapHandler               *handler.RoadmapHandler
	ChatRoomHandler              *handler.ChatRoomHandler
	AtCoderHandler               *handler.AtCoderHandler
	BadgeHandler                 *handler.BadgeHandler
	LearningLogHandler           *handler.LearningLogHandler
	AIAdviceHandler              *handler.AIAdviceHandler
	EmailPreferencesHandler      *handler.EmailPreferencesHandler
	LevelHandler                 *handler.LevelHandler
	LearningAnalyticsHandler     *handler.LearningAnalyticsHandler
	RecommendationHandler        *handler.RecommendationHandler
	StudyCircleHandler           *handler.StudyCircleHandler
	SearchHandler                *handler.SearchHandler
	NoteHandler                  *handler.NoteHandler
	NoteFolderHandler            *handler.NoteFolderHandler
	NoteTemplateHandler          *handler.NoteTemplateHandler
	NoteLinkHandler              *handler.NoteLinkHandler
	PostSeriesHandler            *handler.PostSeriesHandler
	PostCollectionHandler        *handler.PostCollectionHandler
	PostTagHandler               *handler.PostTagHandler
	PostPinHandler               *handler.PostPinHandler
	PostViewHandler              *handler.PostViewHandler
	CommentLikeHandler           *handler.CommentLikeHandler
	MentionHandler               *handler.MentionHandler
	UserDashboardHandler         *handler.UserDashboardHandler
	NoteStatsHandler             *handler.NoteStatsHandler
	StudyCircleStatsHandler      *handler.StudyCircleStatsHandler
	PostStatsHandler             *handler.PostStatsHandler
	BookReviewStatsHandler       *handler.BookReviewStatsHandler
	QAStatsHandler               *handler.QAStatsHandler
	CodeSnippetStatsHandler      *handler.CodeSnippetStatsHandler
	LearningResourceStatsHandler *handler.LearningResourceStatsHandler
	ProjectStatsHandler          *handler.ProjectStatsHandler
	FollowStatsHandler           *handler.FollowStatsHandler
	RoadmapStatsHandler          *handler.RoadmapStatsHandler
	LearningLogStatsHandler      *handler.LearningLogStatsHandler
	CommentStatsHandler          *handler.CommentStatsHandler
	NotificationStatsHandler     *handler.NotificationStatsHandler
	MessageStatsHandler          *handler.MessageStatsHandler
	MentionStatsHandler          *handler.MentionStatsHandler
	ReactionStatsHandler         *handler.ReactionStatsHandler
	BookmarkStatsHandler         *handler.BookmarkStatsHandler
	YouTubeHandler               *handler.YouTubeHandler
	SpotifyHandler               *handler.SpotifyHandler
	StreakFreezeHandler          *handler.StreakFreezeHandler
	BookmarkCollectionHandler    *handler.BookmarkCollectionHandler
	WeeklyChallengeHandler       *handler.WeeklyChallengeHandler
	PostTemplateHandler          *handler.PostTemplateHandler
	WidgetSettingsHandler        *handler.WidgetSettingsHandler
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
	ValidateAuthToken *usecase.ValidateAuthTokenUseCase
	Hub               *ws.Hub
}

// NewContainer はDIコンテナを構築する。
// リポジトリ→サービス→ハンドラの順で依存関係を解決する。
// sqlPool は sqlc(pgx) へ移行済みのリポジトリ用の接続。GORMからの移行が完了するまで db と併存する。
func NewContainer(db *gorm.DB, sqlPool *pgxpool.Pool, cfg *config.Config, hub *ws.Hub) *Container {
	c := &Container{Hub: hub}

	// follow はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	followRepo := persistence.NewFollowRepository(sqlcgen.New(sqlPool))
	// ダイレクトメッセージはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	messagePort := persistence.NewMessageRepository(sqlcgen.New(sqlPool))
	// ランキングはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	rankingRepo := persistence.NewRankingRepository(sqlcgen.New(sqlPool))
	// Zenn 連携はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence と adapter/external。
	zennPort := persistence.NewZennRepository(sqlPool)
	zennClient := external.NewZennClient()
	// Qiita 連携はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence と adapter/external。
	qiitaPort := persistence.NewQiitaRepository(sqlPool)
	qiitaClient := external.NewQiitaClient()
	// アクティビティレポートはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	activityReportRepo := persistence.NewActivityReportRepository(db)
	// プロジェクトはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	projectRepo := persistence.NewProjectRepository(sqlcgen.New(sqlPool))
	// 学習リソースはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	// learningResourcePort は aiAdviceService もこの port 経由で使う。
	learningResourcePort := persistence.NewLearningResourceRepository(sqlcgen.New(sqlPool))
	// 書籍レビューはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	bookReviewRepo := persistence.NewBookReviewRepository(sqlcgen.New(sqlPool))
	// 質問はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	questionPort := persistence.NewQuestionRepository(sqlPool)
	// 回答はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	answerPort := persistence.NewAnswerRepository(sqlPool)
	// ロードマップはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	// 旧 roadmapRepo は aiAdviceService がまだ使うため残している。
	roadmapPort := persistence.NewRoadmapRepository(db)
	// チャットルームはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	chatRoomPort := persistence.NewChatRoomRepository(sqlPool)
	chatRoomMessagePort := persistence.NewChatRoomMessageRepository(sqlcgen.New(sqlPool))
	// コードスニペットはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	codeSnippetRepo := persistence.NewCodeSnippetRepository(sqlcgen.New(sqlPool))
	codeSnippetPostReader := persistence.NewPostReader(sqlcgen.New(sqlPool))
	createCodeSnippet := usecase.NewCreateCodeSnippetUseCase(codeSnippetRepo, codeSnippetPostReader)
	// レベル / XP はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	levelPort := persistence.NewLevelRepository(sqlcgen.New(sqlPool))
	// 学習分析はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	analyticsPort := persistence.NewLearningAnalyticsRepository(sqlcgen.New(sqlPool))
	// バッジはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	badgePort := persistence.NewBadgeRepository(sqlcgen.New(sqlPool))
	// レコメンドはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	recommendationRepo := persistence.NewRecommendationRepository(sqlcgen.New(sqlPool))

	// スタディサークルはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	studyCirclePort := persistence.NewStudyCircleRepository(db)
	searchStudyCircles := usecase.NewSearchStudyCirclesUseCase(studyCirclePort)
	// ノートはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	notePort := persistence.NewNoteRepository(sqlcgen.New(sqlPool))
	// ノートフォルダはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	noteFolderRepo := persistence.NewNoteFolderRepository(sqlcgen.New(sqlPool))
	// ノートテンプレートはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	noteTemplateRepo := persistence.NewNoteTemplateRepository(sqlcgen.New(sqlPool))
	// 学習ログテンプレートはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	learningLogTemplateRepo := persistence.NewLearningLogTemplateRepository(sqlcgen.New(sqlPool))
	// ノート間リンクはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	noteLinkRepo := persistence.NewNoteLinkRepository(sqlcgen.New(sqlPool))
	noteReader := persistence.NewNoteReader(sqlcgen.New(sqlPool))
	// 投稿シリーズはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	postSeriesRepo := persistence.NewPostSeriesRepository(sqlcgen.New(sqlPool))

	// 通知は保存したうえで受信者へ WebSocket 配信する。
	// 配信の有無で作成側の呼び出しは変わらないよう、port をラップして注入する。
	notificationCreator := notificationCreatorWith(sqlPool, hub)
	followerNotifier := followerNotifierWith(sqlPool, hub)

	// 共通サービス
	// 認証はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	authUserPort := persistence.NewAuthUserRepository(sqlPool)
	passwordResetPort := persistence.NewPasswordResetTokenRepository(sqlcgen.New(sqlPool))
	validateAuthToken := usecase.NewValidateAuthTokenUseCase(cfg.JWTSecret)
	githubOAuthState := usecase.NewOAuthStateUseCase(cfg.JWTSecret, usecase.OAuthProviderGitHub)
	spotifyOAuthState := usecase.NewOAuthStateUseCase(cfg.JWTSecret, usecase.OAuthProviderSpotify)
	authUseCases := handler.AuthUseCases{
		Register:             usecase.NewRegisterUserUseCase(authUserPort, cfg.JWTSecret),
		Login:                usecase.NewLoginUseCase(authUserPort, cfg.JWTSecret),
		GitHubLogin:          usecase.NewGitHubLoginUseCase(authUserPort, cfg.JWTSecret),
		LoginState:           usecase.NewGitHubLoginStateUseCase(cfg.JWTSecret),
		GetMe:                usecase.NewGetMeUseCase(authUserPort),
		RequestPasswordReset: usecase.NewRequestPasswordResetUseCase(authUserPort, passwordResetPort),
		ResetPassword:        usecase.NewResetPasswordUseCase(authUserPort, passwordResetPort),
		DeleteAccount:        usecase.NewDeleteAccountUseCase(authUserPort),
	}
	c.ValidateAuthToken = validateAuthToken

	// ユーザー情報はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	// 旧 userRepo は認証・GitHub・Zenn・Qiita・AtCoder・YouTube・メンションがまだ使うため残している。
	userPort := persistence.NewUserRepository(sqlPool)

	// ドメインサービス
	// GitHub 連携はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、
	// 実装は adapter/persistence（永続化）と adapter/external（GitHub API）。
	// 旧 githubRepo は ai_chat / ai_rule_engine がまだ使うため残している。
	githubPort := persistence.NewGitHubRepository(sqlPool)
	githubClient := external.NewGitHubClient(cfg.GitHubClientID, cfg.GitHubClientSecret, cfg.GitHubRedirectURL)
	syncGitHubData := usecase.NewSyncGitHubDataUseCase(userPort, githubPort, githubClient)
	githubUseCases := handler.GitHubUseCases{
		OAuthURL:      usecase.NewGetGitHubOAuthURLUseCase(githubClient),
		Connect:       usecase.NewConnectGitHubUseCase(userPort, githubClient, syncGitHubData),
		Disconnect:    usecase.NewDisconnectGitHubUseCase(userPort, githubPort),
		Sync:          syncGitHubData,
		Contributions: usecase.NewGetGitHubContributionsUseCase(githubPort),
		Languages:     usecase.NewGetGitHubLanguagesUseCase(githubPort),
		Repos:         usecase.NewGetGitHubReposUseCase(githubPort),
	}
	authGitHubUseCases := handler.AuthGitHubUseCases{
		LoginURL:     usecase.NewGetGitHubLoginURLUseCase(githubClient),
		ExchangeCode: usecase.NewExchangeGitHubCodeUseCase(githubClient),
		GetUser:      usecase.NewGetGitHubUserUseCase(githubClient),
		Sync:         syncGitHubData,
	}
	connectZenn := usecase.NewConnectZennUseCase(userPort, zennPort, zennClient)
	disconnectZenn := usecase.NewDisconnectZennUseCase(userPort, zennPort)
	syncZenn := usecase.NewSyncZennUseCase(userPort, zennPort, zennClient)
	connectQiita := usecase.NewConnectQiitaUseCase(userPort, qiitaPort, qiitaClient)
	disconnectQiita := usecase.NewDisconnectQiitaUseCase(userPort, qiitaPort)
	syncQiita := usecase.NewSyncQiitaUseCase(userPort, qiitaPort, qiitaClient)
	// 学習目標はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	// 旧 learningGoalRepo は recommendation / learning_dashboard がまだ使うため残している。
	learningGoalPort := persistence.NewLearningGoalRepository(sqlcgen.New(sqlPool))
	updateLearningGoal := usecase.NewUpdateLearningGoalUseCase(learningGoalPort)
	// AtCoder 連携はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、HTTP 実装は adapter/external。
	atcoderClient := external.NewAtCoderClient()
	createNote := usecase.NewCreateNoteUseCase(notePort)
	// 学習ログはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	// 旧 learningLogRepo は learning_dashboard と AI アドバイスがまだ使うため残している。
	learningLogPort := persistence.NewLearningLogRepository(db)
	createLearningLog := usecase.NewCreateLearningLogUseCase(learningLogPort, learningGoalPort)

	// テンプレートロードマップの初期登録
	go seedTemplateRoadmaps(db, usecase.NewSeedRoadmapTemplatesUseCase(roadmapPort))

	// AI 機能はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、
	// 実装は adapter/persistence（永続化）と adapter/external（OpenAI）。
	// LLM クライアントは API キー設定時のみ初期化し、未設定ならチャットは 503 になる。
	var llmClient usecaserepo.LLMClient
	if cfg.OpenAIAPIKey != "" {
		llmClient = external.NewOpenAIClient(cfg.OpenAIAPIKey)
		log.Println("OpenAI APIキーが設定されています。LLMチャット機能が有効です。")
	} else {
		log.Println("OpenAI APIキー未設定。ルールベース推薦のみ有効です。")
	}
	aiAdvicePort := persistence.NewAIAdviceRepository(sqlPool)
	aiConversationPort := persistence.NewAIConversationRepository(sqlPool)
	generateAIAdvice := usecase.NewGenerateAIAdviceUseCase(
		learningLogPort, learningGoalPort, roadmapPort,
		githubPort, learningResourcePort, userPort,
	)
	aiChatPrompt := usecase.NewBuildAIChatPromptUseCase(learningGoalPort, learningLogPort, roadmapPort, githubPort)
	aiUseCases := handler.AIAdviceUseCases{
		Generate:           generateAIAdvice,
		MarkAsRead:         usecase.NewMarkAIAdviceAsReadUseCase(aiAdvicePort),
		Unread:             usecase.NewGetUnreadAIAdviceUseCase(aiAdvicePort),
		DailyChatRemaining: usecase.NewGetDailyChatRemainingUseCase(aiConversationPort),
		Chat:               usecase.NewChatWithAIUseCase(aiConversationPort, llmClient, aiChatPrompt),
		ListConversations:  usecase.NewListAIConversationsUseCase(aiConversationPort),
		GetConversation:    usecase.NewGetAIConversationUseCase(aiConversationPort),
		DeleteConversation: usecase.NewDeleteAIConversationUseCase(aiConversationPort),
	}

	// ウィークリーレポートメールはクリーンアーキテクチャ（DIP）へ移行済み。
	// port は usecase/repository、SMTP 実装は adapter/external。
	// SMTP 未設定のときはスケジューラを起動しないため送信 usecase も組み立てない。
	if cfg.SMTPHost != "" {
		log.Println("SMTP設定が検出されました。メール機能が有効です。")
		sendWeeklyReports := usecase.NewSendAllWeeklyReportsUseCase(
			userPort,
			activityReportRepo,
			usecase.NewSendWeeklyReportUseCase(external.NewSMTPEmailSender(cfg), cfg.AppURL),
		)
		weeklyScheduler := scheduler.New(sendWeeklyReports)
		go weeklyScheduler.Start()
	} else {
		log.Println("SMTP未設定。ウィークリーレポートメールのスケジューラは無効です。")
	}

	// ハンドラ
	origins := cfg.CORSOrigins
	c.AuthHandler = handler.NewAuthHandler(authUseCases, authGitHubUseCases)
	c.UserHandler = handler.NewUserHandler(
		usecase.NewListUsersUseCase(userPort),
		usecase.NewGetUserUseCase(userPort),
		usecase.NewGetUserByUsernameUseCase(userPort),
		usecase.NewUpdateUserProfileUseCase(userPort),
		usecase.NewGetProfileCompletenessUseCase(userPort),
	)
	c.FollowHandler = handler.NewFollowHandler(
		usecase.NewFollowUserUseCase(followRepo, notificationCreator),
		usecase.NewUnfollowUserUseCase(followRepo),
		usecase.NewListFollowersUseCase(followRepo),
		usecase.NewListFollowingUseCase(followRepo),
	)
	c.GitHubHandler = handler.NewGitHubHandler(githubUseCases, githubOAuthState)
	// 投稿スライスはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	postPort := persistence.NewPostRepository(sqlPool)
	postReactionPort := persistence.NewPostReactionRepository(sqlcgen.New(sqlPool))
	postAuthorPort := persistence.NewPostAuthorReader(sqlcgen.New(sqlPool))
	postCommentPort := persistence.NewPostCommentRepository(sqlcgen.New(sqlPool))
	postBookmarkPort := persistence.NewPostBookmarkRepository(sqlcgen.New(sqlPool))
	postLikePort := persistence.NewPostLikeRepository(sqlcgen.New(sqlPool))
	notifyFollowers := usecase.NewNotifyFollowersUseCase(followerNotifier)
	// メンションは投稿・コメントの本文から解決するため、投稿スライスから呼ぶ。
	mentionPort := persistence.NewMentionRepository(sqlcgen.New(sqlPool))
	processMentions := usecase.NewProcessMentionsUseCase(mentionPort, persistence.NewUserRepository(sqlPool), notificationCreator)
	c.PostHandler = handler.NewPostHandler(handler.PostUseCases{
		ProcessMentions:       processMentions,
		DeleteCommentMentions: usecase.NewDeleteCommentMentionsUseCase(mentionPort),
		Create:                usecase.NewCreatePostUseCase(postPort, notifyFollowers),
		Get:                   usecase.NewGetPostUseCase(postPort),
		List:                  usecase.NewListPostsUseCase(postPort),
		Count:                 usecase.NewCountPostsUseCase(postPort),
		ListByUser:            usecase.NewListUserPostsUseCase(postPort),
		ListDrafts:            usecase.NewListDraftPostsUseCase(postPort),
		ListScheduled:         usecase.NewListScheduledPostsUseCase(postPort),
		Timeline:              usecase.NewGetTimelineUseCase(postPort),
		Update:                usecase.NewUpdatePostUseCase(postPort),
		Delete:                usecase.NewDeletePostUseCase(postPort),
		Publish:               usecase.NewPublishPostUseCase(postPort, notifyFollowers),
		Unpublish:             usecase.NewUnpublishPostUseCase(postPort),
		Schedule:              usecase.NewSchedulePostPublishUseCase(postPort),
		CancelSchedule:        usecase.NewCancelPostScheduleUseCase(postPort),
		AutoSaveDraft:         usecase.NewAutoSaveDraftUseCase(postPort),
		CountByUser:           usecase.NewCountUserPostsUseCase(postPort),
		CountDrafts:           usecase.NewCountUserDraftsUseCase(postPort),
		CountScheduled:        usecase.NewCountUserScheduledPostsUseCase(postPort),

		Like:     usecase.NewLikePostUseCase(postLikePort, postAuthorPort),
		Unlike:   usecase.NewUnlikePostUseCase(postLikePort, postAuthorPort),
		HasLiked: usecase.NewHasLikedPostUseCase(postLikePort),

		CreateSnippet: createCodeSnippet,

		AddReaction:    usecase.NewAddPostReactionUseCase(postReactionPort, postAuthorPort),
		RemoveReaction: usecase.NewRemovePostReactionUseCase(postReactionPort, postAuthorPort),
		GetReactions:   usecase.NewGetPostReactionsUseCase(postReactionPort),
		ReactionsBatch: usecase.NewGetPostReactionsBatchUseCase(postReactionPort),

		CreateComment: usecase.NewCreatePostCommentUseCase(postCommentPort),
		ListComments:  usecase.NewListPostCommentsUseCase(postCommentPort),
		ListReplies:   usecase.NewListCommentRepliesUseCase(postCommentPort),
		EditComment:   usecase.NewEditPostCommentUseCase(postCommentPort),
		DeleteComment: usecase.NewDeletePostCommentUseCase(postCommentPort),
		HideComment:   usecase.NewHidePostCommentUseCase(postCommentPort),
		UnhideComment: usecase.NewUnhidePostCommentUseCase(postCommentPort),

		Bookmark:       usecase.NewBookmarkPostUseCase(postBookmarkPort, postAuthorPort),
		Unbookmark:     usecase.NewUnbookmarkPostUseCase(postBookmarkPort, postAuthorPort),
		HasBookmarked:  usecase.NewHasBookmarkedPostUseCase(postBookmarkPort),
		ListBookmarks:  usecase.NewListBookmarkedPostsUseCase(postBookmarkPort),
		CountBookmarks: usecase.NewCountBookmarkedPostsUseCase(postBookmarkPort),
	})
	c.CodeSnippetHandler = handler.NewCodeSnippetHandler(
		createCodeSnippet,
		usecase.NewListCodeSnippetsByPostUseCase(codeSnippetRepo),
		usecase.NewListCodeSnippetsByLanguageUseCase(codeSnippetRepo),
		usecase.NewUpdateCodeSnippetUseCase(codeSnippetRepo),
		usecase.NewDeleteCodeSnippetUseCase(codeSnippetRepo),
		usecase.NewListSnippetCommentsUseCase(codeSnippetRepo),
		usecase.NewCreateSnippetCommentUseCase(codeSnippetRepo),
		usecase.NewDeleteSnippetCommentUseCase(codeSnippetRepo),
		usecase.NewSearchCodeSnippetsUseCase(codeSnippetRepo),
		usecase.NewForkCodeSnippetUseCase(codeSnippetRepo, codeSnippetPostReader),
		usecase.NewFavoriteCodeSnippetUseCase(codeSnippetRepo),
		usecase.NewUnfavoriteCodeSnippetUseCase(codeSnippetRepo),
		usecase.NewListFavoritedCodeSnippetsUseCase(codeSnippetRepo),
		usecase.NewCountCodeSnippetsUseCase(codeSnippetRepo),
	)
	c.RankingHandler = handler.NewRankingHandler(
		usecase.NewGetContributionRankingUseCase(rankingRepo),
		usecase.NewGetLanguageRankingUseCase(rankingRepo),
		usecase.NewGetLevelRankingUseCase(rankingRepo),
		usecase.NewListRankingLanguagesUseCase(rankingRepo),
	)
	c.MessageHandler = handler.NewMessageHandler(
		usecase.NewListConversationsUseCase(messagePort),
		usecase.NewGetConversationUseCase(messagePort),
		usecase.NewSendMessageUseCase(messagePort, notificationCreator),
		usecase.NewMarkMessagesAsReadUseCase(messagePort),
	)
	c.WebSocketHandler = handler.NewWebSocketHandler(hub, validateAuthToken, parseOrigins(origins))
	uploadHandler, err := handler.NewUploadHandler()
	if err != nil {
		log.Fatalf("アップロードハンドラの初期化に失敗: %v", err)
	}
	c.UploadHandler = uploadHandler
	// 通知の参照・既読・削除はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	// 通知の作成（WebSocket 配信を含む）は post / badge / level / mention / message がまだ service 経由で使うため残している。
	notificationPort := persistence.NewNotificationRepository(sqlPool)
	c.NotificationHandler = handler.NewNotificationHandler(
		usecase.NewListNotificationsUseCase(notificationPort),
		usecase.NewCountUnreadNotificationsUseCase(notificationPort),
		usecase.NewMarkNotificationAsReadUseCase(notificationPort),
		usecase.NewMarkAllNotificationsAsReadUseCase(notificationPort),
		usecase.NewDeleteNotificationUseCase(notificationPort),
	)
	c.ZennHandler = handler.NewArticlePlatformHandler("Zenn", handler.ArticlePlatformOps[model.ZennArticle, model.ZennStats]{
		Connect:     connectZenn.Execute,
		Disconnect:  disconnectZenn.Execute,
		Sync:        syncZenn.Execute,
		GetArticles: usecase.NewListZennArticlesUseCase(zennPort).Execute,
		GetStats:    usecase.NewGetZennStatsUseCase(zennPort).Execute,
	})
	c.QiitaHandler = handler.NewArticlePlatformHandler("Qiita", handler.ArticlePlatformOps[model.QiitaArticle, model.QiitaStats]{
		Connect:     connectQiita.Execute,
		Disconnect:  disconnectQiita.Execute,
		Sync:        syncQiita.Execute,
		GetArticles: usecase.NewListQiitaArticlesUseCase(qiitaPort).Execute,
		GetStats:    usecase.NewGetQiitaStatsUseCase(qiitaPort).Execute,
	})
	c.LearningGoalHandler = handler.NewLearningGoalHandler(
		usecase.NewCreateLearningGoalUseCase(learningGoalPort),
		usecase.NewGetLearningGoalUseCase(learningGoalPort),
		usecase.NewListLearningGoalsUseCase(learningGoalPort),
		usecase.NewListActiveLearningGoalsUseCase(learningGoalPort),
		usecase.NewListLearningGoalsByCategoryUseCase(learningGoalPort),
		usecase.NewListLearningGoalsByStatusUseCase(learningGoalPort),
		usecase.NewGetLearningGoalStatsUseCase(learningGoalPort),
		updateLearningGoal,
		usecase.NewGetGoalDeadlineAlertsUseCase(learningGoalPort),
		usecase.NewDuplicateLearningGoalUseCase(learningGoalPort),
		usecase.NewToggleLearningGoalShareUseCase(learningGoalPort),
		usecase.NewListPublicLearningGoalsUseCase(learningGoalPort),
		usecase.NewListPublicLearningGoalsByUserUseCase(learningGoalPort),
		usecase.NewCountLearningGoalsUseCase(learningGoalPort),
		usecase.NewDeleteLearningGoalUseCase(learningGoalPort),
		usecase.NewBatchUpdateGoalProgressUseCase(updateLearningGoal),
		usecase.NewGetGoalForecastUseCase(learningGoalPort),
	)
	c.ActivityReportHandler = handler.NewActivityReportHandler(
		usecase.NewGetWeeklyActivityReportUseCase(activityReportRepo),
		usecase.NewGetMonthlyActivityReportUseCase(activityReportRepo),
		usecase.NewGetActivityReportComparisonUseCase(activityReportRepo),
	)
	c.ProjectHandler = handler.NewProjectHandler(
		usecase.NewCreateProjectUseCase(projectRepo),
		usecase.NewGetProjectUseCase(projectRepo),
		usecase.NewListProjectsByUserUseCase(projectRepo),
		usecase.NewListFeaturedProjectsUseCase(projectRepo),
		usecase.NewListAllProjectsUseCase(projectRepo),
		usecase.NewListArchivedProjectsUseCase(projectRepo),
		usecase.NewSearchProjectsUseCase(projectRepo),
		usecase.NewUpdateProjectUseCase(projectRepo),
		usecase.NewUpdateProjectFeaturedUseCase(projectRepo),
		usecase.NewArchiveProjectUseCase(projectRepo),
		usecase.NewUnarchiveProjectUseCase(projectRepo),
		usecase.NewDeleteProjectUseCase(projectRepo),
		usecase.NewCountProjectsUseCase(projectRepo),
	)
	c.LearningResourceHandler = handler.NewLearningResourceHandler(
		usecase.NewCreateLearningResourceUseCase(learningResourcePort),
		usecase.NewGetLearningResourceUseCase(learningResourcePort),
		usecase.NewListLearningResourcesByUserUseCase(learningResourcePort),
		usecase.NewListPublicLearningResourcesUseCase(learningResourcePort),
		usecase.NewListLearningResourcesByDifficultyUseCase(learningResourcePort),
		usecase.NewSearchLearningResourcesUseCase(learningResourcePort),
		usecase.NewUpdateLearningResourceUseCase(learningResourcePort),
		usecase.NewUpdateLearningResourceVisibilityUseCase(learningResourcePort),
		usecase.NewDeleteLearningResourceUseCase(learningResourcePort),
		usecase.NewLikeLearningResourceUseCase(learningResourcePort),
		usecase.NewUnlikeLearningResourceUseCase(learningResourcePort),
		usecase.NewHasLikedLearningResourceUseCase(learningResourcePort),
		usecase.NewSaveLearningResourceUseCase(learningResourcePort),
		usecase.NewUnsaveLearningResourceUseCase(learningResourcePort),
		usecase.NewHasSavedLearningResourceUseCase(learningResourcePort),
		usecase.NewListSavedLearningResourcesUseCase(learningResourcePort),
		usecase.NewCountLearningResourcesUseCase(learningResourcePort),
	)
	c.BookReviewHandler = handler.NewBookReviewHandler(
		usecase.NewCreateBookReviewUseCase(bookReviewRepo),
		usecase.NewGetBookReviewUseCase(bookReviewRepo),
		usecase.NewListBookReviewsByUserUseCase(bookReviewRepo),
		usecase.NewListAllBookReviewsUseCase(bookReviewRepo),
		usecase.NewListBookReviewsByRatingUseCase(bookReviewRepo),
		usecase.NewSearchBookReviewsUseCase(bookReviewRepo),
		usecase.NewUpdateBookReviewUseCase(bookReviewRepo),
		usecase.NewUpdateBookReviewStatusUseCase(bookReviewRepo),
		usecase.NewArchiveBookReviewUseCase(bookReviewRepo),
		usecase.NewUpdateBookReviewProgressUseCase(bookReviewRepo),
		usecase.NewDeleteBookReviewUseCase(bookReviewRepo),
		usecase.NewCountBookReviewsUseCase(bookReviewRepo),
	)
	c.QuestionHandler = handler.NewQuestionHandler(
		usecase.NewCreateQuestionUseCase(questionPort),
		usecase.NewListQuestionsUseCase(questionPort),
		usecase.NewSearchQuestionsUseCase(questionPort),
		usecase.NewGetQuestionUseCase(questionPort),
		usecase.NewListQuestionsByUserUseCase(questionPort),
		usecase.NewGetQuestionUserVoteUseCase(questionPort),
		usecase.NewUpdateQuestionUseCase(questionPort),
		usecase.NewDeleteQuestionUseCase(questionPort),
		usecase.NewVoteQuestionUseCase(questionPort),
		usecase.NewRemoveQuestionVoteUseCase(questionPort),
		usecase.NewListSolvedQuestionsUseCase(questionPort),
		usecase.NewListUnansweredQuestionsUseCase(questionPort),
		usecase.NewBookmarkQuestionUseCase(questionPort),
		usecase.NewUnbookmarkQuestionUseCase(questionPort),
		usecase.NewListBookmarkedQuestionsUseCase(questionPort),
		usecase.NewCountQuestionsUseCase(questionPort),
	)
	c.AnswerHandler = handler.NewAnswerHandler(
		usecase.NewListAnswersUseCase(answerPort),
		usecase.NewCreateAnswerUseCase(answerPort, questionPort),
		usecase.NewUpdateAnswerUseCase(answerPort),
		usecase.NewDeleteAnswerUseCase(answerPort),
		usecase.NewSetBestAnswerUseCase(answerPort, questionPort),
		usecase.NewVoteAnswerUseCase(answerPort),
		usecase.NewRemoveAnswerVoteUseCase(answerPort),
		usecase.NewListAnswersByVoteRangeUseCase(answerPort),
	)
	c.RoadmapHandler = handler.NewRoadmapHandler(
		usecase.NewCreateRoadmapUseCase(roadmapPort),
		usecase.NewGetRoadmapUseCase(roadmapPort),
		usecase.NewListRoadmapsByUserUseCase(roadmapPort),
		usecase.NewListRoadmapsByStatusUseCase(roadmapPort),
		usecase.NewListPublicRoadmapsUseCase(roadmapPort),
		usecase.NewUpdateRoadmapUseCase(roadmapPort),
		usecase.NewUpdateRoadmapVisibilityUseCase(roadmapPort),
		usecase.NewDeleteRoadmapUseCase(roadmapPort),
		usecase.NewCopyRoadmapUseCase(roadmapPort),
		usecase.NewListRoadmapTemplatesUseCase(roadmapPort),
		usecase.NewCreateRoadmapFromTemplateUseCase(roadmapPort),
		usecase.NewCreateRoadmapStepUseCase(roadmapPort),
		usecase.NewUpdateRoadmapStepUseCase(roadmapPort),
		usecase.NewUpdateRoadmapStepCompletionUseCase(roadmapPort),
		usecase.NewBatchCompleteRoadmapStepsUseCase(roadmapPort),
		usecase.NewDeleteRoadmapStepUseCase(roadmapPort),
		usecase.NewReorderRoadmapStepsUseCase(roadmapPort),
		usecase.NewGetRoadmapStatsUseCase(persistence.NewRoadmapStatsRepository(sqlcgen.New(sqlPool))),
		usecase.NewCountRoadmapsUseCase(roadmapPort),
	)
	c.ChatRoomHandler = handler.NewChatRoomHandler(
		usecase.NewCreateChatRoomUseCase(chatRoomPort),
		usecase.NewListMyChatRoomsUseCase(chatRoomPort),
		usecase.NewGetChatRoomUseCase(chatRoomPort),
		usecase.NewUpdateChatRoomUseCase(chatRoomPort),
		usecase.NewDeleteChatRoomUseCase(chatRoomPort),
		usecase.NewListChatRoomMembersUseCase(chatRoomPort),
		usecase.NewAddChatRoomMemberUseCase(chatRoomPort),
		usecase.NewRemoveChatRoomMemberUseCase(chatRoomPort),
		usecase.NewListChatRoomMessagesUseCase(chatRoomPort, chatRoomMessagePort),
		usecase.NewSendChatRoomMessageUseCase(chatRoomPort, chatRoomMessagePort, hub),
		usecase.NewCountMyChatRoomsUseCase(chatRoomPort),
	)
	c.AtCoderHandler = handler.NewAtCoderHandler(
		usecase.NewGetAtCoderRatingUseCase(atcoderClient),
		usecase.NewConnectAtCoderUseCase(userPort, atcoderClient),
		usecase.NewDisconnectAtCoderUseCase(userPort),
	)
	c.BadgeHandler = handler.NewBadgeHandler(
		usecase.NewGetUserBadgesUseCase(badgePort),
		usecase.NewNotifyBadgeEarnedUseCase(notificationCreator),
	)
	c.LearningLogHandler = handler.NewLearningLogHandler(
		createLearningLog,
		usecase.NewBatchCreateLearningLogsUseCase(learningLogPort),
		usecase.NewImportLearningLogsCSVUseCase(learningLogPort),
		usecase.NewGetLearningLogUseCase(learningLogPort),
		usecase.NewListLearningLogsUseCase(learningLogPort),
		usecase.NewUpdateLearningLogUseCase(learningLogPort),
		usecase.NewDeleteLearningLogUseCase(learningLogPort),
		usecase.NewGetLearningStreakUseCase(learningLogPort),
		usecase.NewGetLearningCalendarUseCase(learningLogPort),
		usecase.NewExportLearningLogsCSVUseCase(learningLogPort),
		usecase.NewExportLearningLogsJSONUseCase(learningLogPort),
		usecase.NewListLearningLogsByCategoryUseCase(learningLogPort),
		usecase.NewListLearningLogsBySourceUseCase(learningLogPort),
		usecase.NewGetWeeklyLearningDurationUseCase(learningLogPort),
		usecase.NewFavoriteLearningLogUseCase(learningLogPort),
		usecase.NewUnfavoriteLearningLogUseCase(learningLogPort),
		usecase.NewListRecentLearningCategoriesUseCase(learningLogPort),
		usecase.NewListGoalLinkedLogsUseCase(learningLogPort, learningGoalPort),
		usecase.NewGetGoalProgressUseCase(learningLogPort, learningGoalPort),
		usecase.NewListFavoriteLearningLogsUseCase(learningLogPort),
		usecase.NewGetLearningLogMonthlySummaryUseCase(learningLogPort),
		usecase.NewCountLearningLogsUseCase(learningLogPort),
	)
	c.AIAdviceHandler = handler.NewAIAdviceHandler(aiUseCases)
	c.EmailPreferencesHandler = handler.NewEmailPreferencesHandler(
		usecase.NewGetEmailPreferencesUseCase(userPort),
		usecase.NewUpdateEmailPreferencesUseCase(userPort),
	)
	c.LevelHandler = handler.NewLevelHandler(
		usecase.NewGetLevelInfoUseCase(levelPort),
		usecase.NewGetXPBreakdownUseCase(levelPort),
	)
	c.LearningAnalyticsHandler = handler.NewLearningAnalyticsHandler(
		usecase.NewGetLearningHeatmapUseCase(analyticsPort),
		usecase.NewGetCategoryBreakdownUseCase(analyticsPort),
		usecase.NewGetWeeklyTrendsUseCase(analyticsPort),
		usecase.NewGetDayOfWeekSummaryUseCase(analyticsPort),
		usecase.NewGetProductivityScoreUseCase(analyticsPort),
		usecase.NewGetLearningInsightsUseCase(analyticsPort),
	)
	c.RecommendationHandler = handler.NewRecommendationHandler(
		usecase.NewGetRecommendedUsersUseCase(recommendationRepo, userPort),
		usecase.NewGetTrendingPostsUseCase(recommendationRepo),
		usecase.NewGetTrendingResourcesUseCase(recommendationRepo),
	)
	c.StudyCircleHandler = handler.NewStudyCircleHandler(
		usecase.NewCreateStudyCircleUseCase(studyCirclePort),
		usecase.NewListMyStudyCirclesUseCase(studyCirclePort),
		usecase.NewListStudyCirclesByStatusUseCase(studyCirclePort),
		usecase.NewGetStudyCircleUseCase(studyCirclePort),
		usecase.NewUpdateStudyCircleUseCase(studyCirclePort),
		usecase.NewDeleteStudyCircleUseCase(studyCirclePort),
		usecase.NewListStudyCircleMembersUseCase(studyCirclePort),
		usecase.NewAddStudyCircleMemberUseCase(studyCirclePort),
		usecase.NewUpdateStudyCircleMemberRoleUseCase(studyCirclePort),
		usecase.NewRemoveStudyCircleMemberUseCase(studyCirclePort),
		usecase.NewCreateStudyCircleStepUseCase(studyCirclePort),
		usecase.NewUpdateStudyCircleStepUseCase(studyCirclePort),
		usecase.NewDeleteStudyCircleStepUseCase(studyCirclePort),
		usecase.NewReorderStudyCircleStepsUseCase(studyCirclePort),
		usecase.NewUpdateStudyCircleProgressUseCase(studyCirclePort),
		usecase.NewListStudyCircleProgressUseCase(studyCirclePort),
		usecase.NewCreateStudyCircleCheckinUseCase(studyCirclePort),
		usecase.NewListStudyCircleCheckinsUseCase(studyCirclePort),
		usecase.NewGetStudyCircleStreakRankingUseCase(studyCirclePort),
		searchStudyCircles,
		usecase.NewCountStudyCirclesUseCase(studyCirclePort),
	)
	// 投稿検索はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	c.SearchHandler = handler.NewSearchHandler(
		usecase.NewSearchPostsUseCase(persistence.NewPostSearchRepository(sqlcgen.New(sqlPool))),
		searchStudyCircles,
	)
	c.NoteHandler = handler.NewNoteHandler(
		createNote,
		usecase.NewGetNoteUseCase(notePort),
		usecase.NewListNotesUseCase(notePort),
		usecase.NewListNotesByFolderUseCase(notePort),
		usecase.NewUpdateNoteUseCase(notePort),
		usecase.NewDeleteNoteUseCase(notePort),
		usecase.NewSearchNotesUseCase(notePort),
		usecase.NewCountNotesUseCase(notePort),
		usecase.NewToggleNoteFavoriteUseCase(notePort),
		usecase.NewListFavoriteNotesUseCase(notePort),
		usecase.NewCountFavoriteNotesUseCase(notePort),
		usecase.NewArchiveNoteUseCase(notePort),
		usecase.NewUnarchiveNoteUseCase(notePort),
		usecase.NewListArchivedNotesUseCase(notePort),
		usecase.NewCountArchivedNotesUseCase(notePort),
		usecase.NewListNoteTagsUseCase(notePort),
		usecase.NewExportNoteMarkdownUseCase(notePort),
		usecase.NewDuplicateNoteUseCase(notePort),
	)
	// ノートのバージョン履歴はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	noteVersionRepo := persistence.NewNoteVersionRepository(sqlcgen.New(sqlPool))
	noteUpdater := persistence.NewNoteUpdater(sqlcgen.New(sqlPool))
	c.NoteVersionHandler = handler.NewNoteVersionHandler(
		usecase.NewListNoteVersionsUseCase(noteVersionRepo, noteReader),
		usecase.NewGetNoteVersionUseCase(noteVersionRepo, noteReader),
		usecase.NewRestoreNoteVersionUseCase(noteVersionRepo, noteReader, noteUpdater),
	)
	c.NoteFolderHandler = handler.NewNoteFolderHandler(
		usecase.NewCreateNoteFolderUseCase(noteFolderRepo),
		usecase.NewGetNoteFolderUseCase(noteFolderRepo),
		usecase.NewListNoteFoldersUseCase(noteFolderRepo),
		usecase.NewListChildNoteFoldersUseCase(noteFolderRepo),
		usecase.NewListRootNoteFoldersUseCase(noteFolderRepo),
		usecase.NewUpdateNoteFolderUseCase(noteFolderRepo),
		usecase.NewCountNoteFoldersUseCase(noteFolderRepo),
		usecase.NewDeleteNoteFolderUseCase(noteFolderRepo),
	)
	c.NoteTemplateHandler = handler.NewNoteTemplateHandler(
		usecase.NewCreateNoteTemplateUseCase(noteTemplateRepo),
		usecase.NewGetNoteTemplateUseCase(noteTemplateRepo),
		usecase.NewListNoteTemplatesUseCase(noteTemplateRepo),
		usecase.NewGetDefaultNoteTemplateUseCase(noteTemplateRepo),
		usecase.NewUpdateNoteTemplateUseCase(noteTemplateRepo),
		usecase.NewDeleteNoteTemplateUseCase(noteTemplateRepo),
		usecase.NewCreateNoteFromTemplateUseCase(noteTemplateRepo, createNote),
		usecase.NewCountNoteTemplatesUseCase(noteTemplateRepo),
	)
	c.LearningLogTemplateHandler = handler.NewLearningLogTemplateHandler(
		usecase.NewCreateLearningLogTemplateUseCase(learningLogTemplateRepo),
		usecase.NewGetLearningLogTemplateUseCase(learningLogTemplateRepo),
		usecase.NewListLearningLogTemplatesUseCase(learningLogTemplateRepo),
		usecase.NewGetDefaultLearningLogTemplateUseCase(learningLogTemplateRepo),
		usecase.NewUpdateLearningLogTemplateUseCase(learningLogTemplateRepo),
		usecase.NewDeleteLearningLogTemplateUseCase(learningLogTemplateRepo),
		usecase.NewCreateLearningLogFromTemplateUseCase(learningLogTemplateRepo, createLearningLog),
		usecase.NewCountLearningLogTemplatesUseCase(learningLogTemplateRepo),
	)
	c.NoteLinkHandler = handler.NewNoteLinkHandler(
		usecase.NewCreateNoteLinkUseCase(noteLinkRepo, noteReader),
		usecase.NewListNoteLinksUseCase(noteLinkRepo),
		usecase.NewListNoteBacklinksUseCase(noteLinkRepo),
		usecase.NewGetNoteLinkStatsUseCase(noteLinkRepo, noteReader),
		usecase.NewDeleteNoteLinkUseCase(noteLinkRepo, noteReader),
	)
	c.PostSeriesHandler = handler.NewPostSeriesHandler(
		usecase.NewCreatePostSeriesUseCase(postSeriesRepo),
		usecase.NewGetPostSeriesUseCase(postSeriesRepo),
		usecase.NewListPostSeriesUseCase(postSeriesRepo),
		usecase.NewCountPostSeriesUseCase(postSeriesRepo),
		usecase.NewUpdatePostSeriesUseCase(postSeriesRepo),
		usecase.NewDeletePostSeriesUseCase(postSeriesRepo),
		usecase.NewAddPostToSeriesUseCase(postSeriesRepo),
		usecase.NewRemovePostFromSeriesUseCase(postSeriesRepo),
		usecase.NewListPostSeriesPostsUseCase(postSeriesRepo),
	)

	// 投稿コレクションはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	postCollectionRepo := persistence.NewPostCollectionRepository(sqlcgen.New(sqlPool))
	c.PostCollectionHandler = handler.NewPostCollectionHandler(
		usecase.NewCreatePostCollectionUseCase(postCollectionRepo),
		usecase.NewGetPostCollectionUseCase(postCollectionRepo),
		usecase.NewListPostCollectionsForViewerUseCase(postCollectionRepo),
		usecase.NewCountPostCollectionsUseCase(postCollectionRepo),
		usecase.NewUpdatePostCollectionUseCase(postCollectionRepo),
		usecase.NewDeletePostCollectionUseCase(postCollectionRepo),
		usecase.NewAddPostToCollectionUseCase(postCollectionRepo),
		usecase.NewRemovePostFromCollectionUseCase(postCollectionRepo),
		usecase.NewListPostCollectionPostsUseCase(postCollectionRepo),
	)

	// 投稿タグはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	postTagRepo := persistence.NewPostTagRepository(sqlPool)
	setPostTags := usecase.NewSetPostTagsUseCase(postTagRepo, persistence.NewPostReader(sqlcgen.New(sqlPool)))
	c.PostTagHandler = handler.NewPostTagHandler(
		setPostTags,
		usecase.NewGetPostTagsUseCase(postTagRepo),
		usecase.NewFindPostsByTagUseCase(postTagRepo),
		usecase.NewGetPopularTagsUseCase(postTagRepo),
	)
	c.PostHandler.SetAutoTagsUseCase(usecase.NewSetAutoPostTagsUseCase(setPostTags))

	// post_pin はクリーンアーキテクチャ(DIP)へ移行済み。
	postPinRepo := persistence.NewPostPinRepository(sqlPool)
	postReader := persistence.NewPostReader(sqlcgen.New(sqlPool))
	c.PostPinHandler = handler.NewPostPinHandler(
		usecase.NewPinPostUseCase(postPinRepo, postReader),
		usecase.NewUnpinPostUseCase(postPinRepo),
		usecase.NewListPinnedPostsUseCase(postPinRepo),
		usecase.NewCountPinnedPostsUseCase(postPinRepo),
		usecase.NewReorderPinnedPostsUseCase(postPinRepo),
	)

	// comment_like はクリーンアーキテクチャ(DIP)へ移行済み。
	commentLikeRepo := persistence.NewCommentLikeRepository(sqlPool)
	commentReader := persistence.NewCommentReader(sqlcgen.New(sqlPool))
	c.CommentLikeHandler = handler.NewCommentLikeHandler(
		usecase.NewLikeCommentUseCase(commentLikeRepo, commentReader),
		usecase.NewUnlikeCommentUseCase(commentLikeRepo, commentReader),
		usecase.NewGetCommentLikeStatusUseCase(commentLikeRepo, commentReader),
	)

	// 投稿閲覧数はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	postViewRepo := persistence.NewPostViewRepository(sqlPool)
	c.PostViewHandler = handler.NewPostViewHandler(
		usecase.NewRecordPostViewUseCase(postViewRepo),
		usecase.NewGetPostViewCountUseCase(postViewRepo),
		usecase.NewGetMostViewedPostsUseCase(postViewRepo),
	)

	// メンションはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	// 記録側（投稿・コメント作成時の処理）は投稿スライスの配線でつないでいる。
	c.MentionHandler = handler.NewMentionHandler(
		usecase.NewListUserMentionsUseCase(mentionPort),
		usecase.NewListPostMentionsUseCase(mentionPort),
	)

	// ユーザーダッシュボード統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	userDashboardRepo := persistence.NewUserDashboardRepository(sqlcgen.New(sqlPool))
	c.UserDashboardHandler = handler.NewUserDashboardHandler(
		usecase.NewGetUserDashboardStatsUseCase(userDashboardRepo),
	)

	// ノート統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	noteStatsRepo := persistence.NewNoteStatsRepository(sqlcgen.New(sqlPool))
	c.NoteStatsHandler = handler.NewNoteStatsHandler(
		usecase.NewGetNoteStatsUseCase(noteStatsRepo),
	)

	// スタディサークル統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	studyCircleStatsRepo := persistence.NewStudyCircleStatsRepository(sqlcgen.New(sqlPool))
	c.StudyCircleStatsHandler = handler.NewStudyCircleStatsHandler(
		usecase.NewGetStudyCircleStatsUseCase(studyCircleStatsRepo),
	)

	// 投稿統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	postStatsRepo := persistence.NewPostStatsRepository(sqlcgen.New(sqlPool))
	c.PostStatsHandler = handler.NewPostStatsHandler(
		usecase.NewGetPostStatsUseCase(postStatsRepo),
	)

	// 書籍レビュー統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	bookReviewStatsRepo := persistence.NewBookReviewStatsRepository(sqlcgen.New(sqlPool))
	c.BookReviewStatsHandler = handler.NewBookReviewStatsHandler(
		usecase.NewGetBookReviewStatsUseCase(bookReviewStatsRepo),
	)

	// Q&A 統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	qaStatsRepo := persistence.NewQAStatsRepository(sqlcgen.New(sqlPool))
	c.QAStatsHandler = handler.NewQAStatsHandler(
		usecase.NewGetQAStatsUseCase(qaStatsRepo),
	)

	// コードスニペット統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	codeSnippetStatsRepo := persistence.NewCodeSnippetStatsRepository(sqlcgen.New(sqlPool))
	c.CodeSnippetStatsHandler = handler.NewCodeSnippetStatsHandler(
		usecase.NewGetCodeSnippetStatsUseCase(codeSnippetStatsRepo),
	)

	// 学習リソース統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	learningResourceStatsRepo := persistence.NewLearningResourceStatsRepository(sqlcgen.New(sqlPool))
	c.LearningResourceStatsHandler = handler.NewLearningResourceStatsHandler(
		usecase.NewGetLearningResourceStatsUseCase(learningResourceStatsRepo),
	)

	// プロジェクト統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	projectStatsRepo := persistence.NewProjectStatsRepository(sqlcgen.New(sqlPool))
	c.ProjectStatsHandler = handler.NewProjectStatsHandler(
		usecase.NewGetProjectStatsUseCase(projectStatsRepo),
	)

	// フォロー統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	followStatsRepo := persistence.NewFollowStatsRepository(sqlcgen.New(sqlPool))
	c.FollowStatsHandler = handler.NewFollowStatsHandler(
		usecase.NewGetFollowStatsUseCase(followStatsRepo),
	)

	// ロードマップ統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	roadmapStatsRepo := persistence.NewRoadmapStatsRepository(sqlcgen.New(sqlPool))
	c.RoadmapStatsHandler = handler.NewRoadmapStatsHandler(
		usecase.NewGetRoadmapStatsUseCase(roadmapStatsRepo),
	)

	// 学習ログ統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	learningLogStatsRepo := persistence.NewLearningLogStatsRepository(sqlcgen.New(sqlPool))
	c.LearningLogStatsHandler = handler.NewLearningLogStatsHandler(
		usecase.NewGetLearningLogStatsUseCase(learningLogStatsRepo),
	)

	// コメント統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	commentStatsRepo := persistence.NewCommentStatsRepository(sqlcgen.New(sqlPool))
	c.CommentStatsHandler = handler.NewCommentStatsHandler(
		usecase.NewGetCommentStatsUseCase(commentStatsRepo),
	)

	// 通知統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	notificationStatsRepo := persistence.NewNotificationStatsRepository(sqlcgen.New(sqlPool))
	c.NotificationStatsHandler = handler.NewNotificationStatsHandler(
		usecase.NewGetNotificationStatsUseCase(notificationStatsRepo),
	)

	// メッセージ統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	messageStatsRepo := persistence.NewMessageStatsRepository(sqlcgen.New(sqlPool))
	c.MessageStatsHandler = handler.NewMessageStatsHandler(
		usecase.NewGetMessageStatsUseCase(messageStatsRepo),
	)

	// メンション統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	mentionStatsRepo := persistence.NewMentionStatsRepository(sqlcgen.New(sqlPool))
	c.MentionStatsHandler = handler.NewMentionStatsHandler(
		usecase.NewGetMentionStatsUseCase(mentionStatsRepo),
	)

	// リアクション統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	reactionStatsRepo := persistence.NewReactionStatsRepository(sqlcgen.New(sqlPool))
	c.ReactionStatsHandler = handler.NewReactionStatsHandler(
		usecase.NewGetReactionStatsUseCase(reactionStatsRepo),
		usecase.NewGetReactionSummaryUseCase(reactionStatsRepo),
	)

	// ブックマーク統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	bookmarkStatsRepo := persistence.NewBookmarkStatsRepository(sqlcgen.New(sqlPool))
	c.BookmarkStatsHandler = handler.NewBookmarkStatsHandler(
		usecase.NewGetBookmarkStatsUseCase(bookmarkStatsRepo),
	)

	// Spotify 連携はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、
	// 実装は adapter/persistence（永続化）と adapter/external（Spotify API）。
	spotifyPort := persistence.NewSpotifyRepository(sqlcgen.New(sqlPool))
	spotifyClient := external.NewSpotifyClient(cfg.SpotifyClientID, cfg.SpotifyClientSecret, cfg.SpotifyRedirectURL)
	spotifyUseCases := handler.SpotifyUseCases{
		OAuthURL:         usecase.NewGetSpotifyOAuthURLUseCase(spotifyClient),
		Connect:          usecase.NewConnectSpotifyUseCase(userPort, spotifyClient),
		Disconnect:       usecase.NewDisconnectSpotifyUseCase(userPort, spotifyPort),
		CurrentlyPlaying: usecase.NewGetSpotifyCurrentlyPlayingUseCase(userPort, spotifyClient),
		RecentlyPlayed:   usecase.NewGetSpotifyRecentlyPlayedUseCase(userPort, spotifyClient),
	}
	c.SpotifyHandler = handler.NewSpotifyHandler(spotifyUseCases, spotifyOAuthState)

	// YouTube 連携はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence と adapter/external。
	// APIキー未設定のときは検索クライアントを nil のままにし、利用不可（503）として扱う。
	youtubeVideoPort := persistence.NewYouTubeVideoRepository(sqlPool)
	var youtubeSearcher usecaserepo.YouTubeVideoSearcher
	if cfg.YouTubeAPIKey != "" {
		youtubeSearcher = external.NewYouTubeClient(cfg.YouTubeAPIKey)
		log.Println("YouTube APIキーが設定されています。YouTube動画検索機能が有効です。")
	} else {
		log.Println("YouTube APIキー未設定。YouTube動画検索機能は無効です。")
	}
	c.YouTubeHandler = handler.NewYouTubeHandler(
		usecase.NewSearchYouTubeVideosUseCase(youtubeVideoPort, youtubeSearcher),
		usecase.NewRecommendYouTubeVideosUseCase(userPort, youtubeVideoPort, youtubeSearcher),
		usecase.NewCheckYouTubeAvailabilityUseCase(youtubeSearcher),
	)

	// ストリークフリーズサービス
	// ストリークフリーズはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	streakFreezeRepo := persistence.NewStreakFreezeRepository(sqlPool)
	c.StreakFreezeHandler = handler.NewStreakFreezeHandler(
		usecase.NewUseStreakFreezeUseCase(streakFreezeRepo),
		usecase.NewGetStreakFreezeStatusUseCase(streakFreezeRepo),
	)

	// ブックマークコレクションサービス
	// ブックマークコレクションはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	bookmarkCollectionRepo := persistence.NewBookmarkCollectionRepository(sqlcgen.New(sqlPool))
	c.BookmarkCollectionHandler = handler.NewBookmarkCollectionHandler(
		usecase.NewCreateBookmarkCollectionUseCase(bookmarkCollectionRepo),
		usecase.NewListBookmarkCollectionsUseCase(bookmarkCollectionRepo),
		usecase.NewUpdateBookmarkCollectionUseCase(bookmarkCollectionRepo),
		usecase.NewDeleteBookmarkCollectionUseCase(bookmarkCollectionRepo),
		usecase.NewAddPostToBookmarkCollectionUseCase(bookmarkCollectionRepo),
		usecase.NewRemovePostFromBookmarkCollectionUseCase(bookmarkCollectionRepo),
		usecase.NewListBookmarkCollectionPostsUseCase(bookmarkCollectionRepo),
		usecase.NewCountBookmarkCollectionsUseCase(bookmarkCollectionRepo),
	)

	// ウィークリーチャレンジはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	weeklyChallengeRepo := persistence.NewWeeklyChallengeRepository(sqlcgen.New(sqlPool))
	c.WeeklyChallengeHandler = handler.NewWeeklyChallengeHandler(
		usecase.NewGetCurrentWeeklyChallengeUseCase(weeklyChallengeRepo),
		usecase.NewUpdateWeeklyChallengeProgressUseCase(weeklyChallengeRepo),
	)

	// 投稿テンプレートはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence（sqlc/pgx）。
	postTemplateRepo := persistence.NewPostTemplateRepository(sqlcgen.New(sqlPool))
	c.PostTemplateHandler = handler.NewPostTemplateHandler(
		usecase.NewCreatePostTemplateUseCase(postTemplateRepo),
		usecase.NewGetPostTemplateUseCase(postTemplateRepo),
		usecase.NewListPostTemplatesUseCase(postTemplateRepo),
		usecase.NewUpdatePostTemplateUseCase(postTemplateRepo),
		usecase.NewDeletePostTemplateUseCase(postTemplateRepo),
	)

	// ウィジェット設定はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	widgetSettingsRepo := persistence.NewWidgetSettingsRepository(sqlcgen.New(sqlPool))
	c.WidgetSettingsHandler = handler.NewWidgetSettingsHandler(
		usecase.NewGetWidgetSettingsUseCase(widgetSettingsRepo),
		usecase.NewUpdateWidgetSettingsUseCase(widgetSettingsRepo),
	)

	// カテゴリ別週間学習目標はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	weeklyGoalRepo := persistence.NewWeeklyGoalRepository(sqlcgen.New(sqlPool))
	c.WeeklyGoalHandler = handler.NewWeeklyGoalHandler(
		usecase.NewSetWeeklyGoalUseCase(weeklyGoalRepo),
		usecase.NewListWeeklyGoalsUseCase(weeklyGoalRepo),
		usecase.NewGetWeeklyGoalProgressUseCase(weeklyGoalRepo),
	)

	// リソース進捗サービス
	// リソース進捗はクリーンアーキテクチャ（DIP）へ移行済み。リソース存在確認は最小 port LearningResourceReader を再利用。
	resourceProgressRepo := persistence.NewResourceProgressRepository(sqlcgen.New(sqlPool))
	resourceProgressResourceReader := persistence.NewLearningResourceReader(sqlcgen.New(sqlPool))
	c.ResourceProgressHandler = handler.NewResourceProgressHandler(
		usecase.NewUpsertResourceProgressUseCase(resourceProgressRepo, resourceProgressResourceReader),
		usecase.NewGetResourceProgressUseCase(resourceProgressRepo),
		usecase.NewListResourceProgressUseCase(resourceProgressRepo),
	)

	// プロジェクトマイルストーンサービス
	// プロジェクトマイルストーンはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	projectMilestoneRepo := persistence.NewProjectMilestoneRepository(sqlcgen.New(sqlPool))
	projectReader := persistence.NewProjectReader(sqlcgen.New(sqlPool))
	c.ProjectMilestoneHandler = handler.NewProjectMilestoneHandler(
		usecase.NewCreateProjectMilestoneUseCase(projectMilestoneRepo, projectReader),
		usecase.NewListProjectMilestonesUseCase(projectMilestoneRepo),
		usecase.NewUpdateProjectMilestoneUseCase(projectMilestoneRepo, projectReader),
		usecase.NewDeleteProjectMilestoneUseCase(projectMilestoneRepo, projectReader),
	)

	// ユーザーアクティビティサービス
	// ユーザーアクティビティはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	userActivityRepo := persistence.NewUserActivityRepository(sqlcgen.New(sqlPool))
	c.UserActivityHandler = handler.NewUserActivityHandler(
		usecase.NewGetActivityTimelineUseCase(userActivityRepo),
	)

	// リソースレビューはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	resourceReviewRepo := persistence.NewResourceReviewRepository(sqlcgen.New(sqlPool))
	learningResourceReader := persistence.NewLearningResourceReader(sqlcgen.New(sqlPool))
	c.ResourceReviewHandler = handler.NewResourceReviewHandler(
		usecase.NewCreateResourceReviewUseCase(resourceReviewRepo, learningResourceReader),
		usecase.NewListResourceReviewsUseCase(resourceReviewRepo),
		usecase.NewUpdateResourceReviewUseCase(resourceReviewRepo),
		usecase.NewDeleteResourceReviewUseCase(resourceReviewRepo),
	)

	// 学習ダッシュボード統合サマリーサービス
	// 学習ダッシュボードはクリーンアーキテクチャ（DIP）へ移行済み。
	// 学習ログ・目標・分析の移行済み実装を、それぞれ最小 port として受け取る。
	c.LearningDashboardHandler = handler.NewLearningDashboardHandler(
		usecase.NewGetLearningDashboardSummaryUseCase(learningLogPort, learningGoalPort, analyticsPort),
	)

	// リマインダー設定はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	reminderSettingsRepo := persistence.NewReminderSettingsRepository(sqlcgen.New(sqlPool))
	c.ReminderSettingsHandler = handler.NewReminderSettingsHandler(
		usecase.NewGetReminderSettingsUseCase(reminderSettingsRepo),
		usecase.NewUpdateReminderSettingsUseCase(reminderSettingsRepo),
	)

	// 通知設定サービス
	// 通知設定はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	notificationSettingsRepo := persistence.NewNotificationSettingsRepository(sqlcgen.New(sqlPool))
	c.NotificationSettingsHandler = handler.NewNotificationSettingsHandler(
		usecase.NewGetNotificationSettingsUseCase(notificationSettingsRepo),
		usecase.NewUpdateNotificationSettingsUseCase(notificationSettingsRepo),
	)

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
func seedTemplateRoadmaps(db *gorm.DB, seed *usecase.SeedRoadmapTemplatesUseCase) {
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
	if err := seed.Execute(context.Background(), user.ID); err != nil {
		log.Printf("テンプレートシード失敗: %v", err)
	}
}

// notificationCreatorWith は通知作成の port を組み立てる。
// hub があれば保存後の配信を上乗せする。
func notificationCreatorWith(sqlPool *pgxpool.Pool, hub *ws.Hub) usecaserepo.NotificationCreator {
	creator := persistence.NewNotificationCreator(sqlPool)
	if hub == nil {
		return creator
	}
	return notify.NewBroadcastingCreator(creator, hub)
}

// followerNotifierWith はフォロワー一括通知の port を組み立てる。
// hub があれば保存後の配信を上乗せする。
func followerNotifierWith(sqlPool *pgxpool.Pool, hub *ws.Hub) usecaserepo.FollowerNotifier {
	notifier := persistence.NewFollowerNotifier(sqlPool)
	if hub == nil {
		return notifier
	}
	return notify.NewBroadcastingFollowerNotifier(notifier, hub)
}

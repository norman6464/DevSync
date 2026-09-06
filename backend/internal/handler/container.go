// Package handler はDevSyncアプリケーションのHTTPハンドラを提供する。
// container.go・router.go は依存関係の組み立てとルーティング登録を担う配線専用ファイル
// （FreStyle同様、DIコンテナ・ルーターを独立パッケージに分けず handler 内に集約している）。
//
// archlint:ignore-file
// このファイルは配線専用（旧 internal/di）のため、handler層の
// 「adapter・usecase/repository(port)を直接importできない」制約から意図的に除外する。
// 独立パッケージだった頃の internal/di・internal/router は元々検証対象外だったため、
// FreStyle同様に配線をhandlerパッケージへ集約した後もその位置づけを維持する
// （通常のhandlerファイル本体の制約は変えない）。
package handler

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/norman6464/devsync/backend/internal/adapter/external"
	"github.com/norman6464/devsync/backend/internal/adapter/notify"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence"
	"github.com/norman6464/devsync/backend/internal/adapter/persistence/sqlcgen"
	"github.com/norman6464/devsync/backend/internal/infra/config"
	"github.com/norman6464/devsync/backend/internal/infra/scheduler"
	"github.com/norman6464/devsync/backend/internal/infra/ws"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	usecaserepo "github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// Container はDI（依存性注入）コンテナ。
// 全ハンドラとルーティングに必要な公開フィールドを保持する。
type Container struct {
	// ハンドラ
	AuthHandler                  *AuthHandler
	UserHandler                  *UserHandler
	FollowHandler                *FollowHandler
	GitHubHandler                *GitHubHandler
	PostHandler                  *PostHandler
	CodeSnippetHandler           *CodeSnippetHandler
	RankingHandler               *RankingHandler
	MessageHandler               *MessageHandler
	WebSocketHandler             *WebSocketHandler
	UploadHandler                *UploadHandler
	NotificationHandler          *NotificationHandler
	ZennHandler                  *ArticlePlatformHandler[model.ZennArticle, model.ZennStats]
	QiitaHandler                 *ArticlePlatformHandler[model.QiitaArticle, model.QiitaStats]
	LearningGoalHandler          *LearningGoalHandler
	ActivityReportHandler        *ActivityReportHandler
	ProjectHandler               *ProjectHandler
	LearningResourceHandler      *LearningResourceHandler
	BookReviewHandler            *BookReviewHandler
	QuestionHandler              *QuestionHandler
	AnswerHandler                *AnswerHandler
	RoadmapHandler               *RoadmapHandler
	ChatRoomHandler              *ChatRoomHandler
	AtCoderHandler               *AtCoderHandler
	BadgeHandler                 *BadgeHandler
	LearningLogHandler           *LearningLogHandler
	AIAdviceHandler              *AIAdviceHandler
	EmailPreferencesHandler      *EmailPreferencesHandler
	LevelHandler                 *LevelHandler
	LearningAnalyticsHandler     *LearningAnalyticsHandler
	RecommendationHandler        *RecommendationHandler
	StudyCircleHandler           *StudyCircleHandler
	SearchHandler                *SearchHandler
	NoteHandler                  *NoteHandler
	NoteFolderHandler            *NoteFolderHandler
	NoteTemplateHandler          *NoteTemplateHandler
	NoteLinkHandler              *NoteLinkHandler
	PostSeriesHandler            *PostSeriesHandler
	PostCollectionHandler        *PostCollectionHandler
	PostTagHandler               *PostTagHandler
	PostPinHandler               *PostPinHandler
	PostViewHandler              *PostViewHandler
	CommentLikeHandler           *CommentLikeHandler
	MentionHandler               *MentionHandler
	UserDashboardHandler         *UserDashboardHandler
	NoteStatsHandler             *NoteStatsHandler
	StudyCircleStatsHandler      *StudyCircleStatsHandler
	PostStatsHandler             *PostStatsHandler
	BookReviewStatsHandler       *BookReviewStatsHandler
	QAStatsHandler               *QAStatsHandler
	CodeSnippetStatsHandler      *CodeSnippetStatsHandler
	LearningResourceStatsHandler *LearningResourceStatsHandler
	ProjectStatsHandler          *ProjectStatsHandler
	FollowStatsHandler           *FollowStatsHandler
	RoadmapStatsHandler          *RoadmapStatsHandler
	LearningLogStatsHandler      *LearningLogStatsHandler
	CommentStatsHandler          *CommentStatsHandler
	NotificationStatsHandler     *NotificationStatsHandler
	MessageStatsHandler          *MessageStatsHandler
	MentionStatsHandler          *MentionStatsHandler
	ReactionStatsHandler         *ReactionStatsHandler
	BookmarkStatsHandler         *BookmarkStatsHandler
	YouTubeHandler               *YouTubeHandler
	SpotifyHandler               *SpotifyHandler
	StreakFreezeHandler          *StreakFreezeHandler
	BookmarkCollectionHandler    *BookmarkCollectionHandler
	WeeklyChallengeHandler       *WeeklyChallengeHandler
	PostTemplateHandler          *PostTemplateHandler
	WidgetSettingsHandler        *WidgetSettingsHandler
	WeeklyGoalHandler            *WeeklyGoalHandler
	NoteVersionHandler           *NoteVersionHandler
	ResourceProgressHandler      *ResourceProgressHandler
	ProjectMilestoneHandler      *ProjectMilestoneHandler
	UserActivityHandler          *UserActivityHandler
	LearningLogTemplateHandler   *LearningLogTemplateHandler
	ResourceReviewHandler        *ResourceReviewHandler
	LearningDashboardHandler     *LearningDashboardHandler
	ReminderSettingsHandler      *ReminderSettingsHandler
	NotificationSettingsHandler  *NotificationSettingsHandler

	// ミドルウェア・コールバック用
	ValidateAuthToken *usecase.ValidateAuthTokenUseCase
	Hub               *ws.Hub
}

// NewContainer はDIコンテナを構築する。
// リポジトリ→サービス→ハンドラの順で依存関係を解決する。
func NewContainer(sqlPool *pgxpool.Pool, cfg *config.Config, hub *ws.Hub) *Container {
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
	activityReportRepo := persistence.NewActivityReportRepository(sqlcgen.New(sqlPool))
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
	roadmapPort := persistence.NewRoadmapRepository(sqlPool)
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
	studyCirclePort := persistence.NewStudyCircleRepository(sqlPool)
	searchStudyCircles := usecase.NewSearchStudyCirclesUseCase(studyCirclePort)
	// ノートはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	notePort := persistence.NewNoteRepository(sqlcgen.New(sqlPool))
	// ノートフォルダはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	noteFolderRepo := persistence.NewNoteFolderRepository(sqlcgen.New(sqlPool))
	// ノートテンプレートはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	// デフォルト解除とCreate/Updateを1トランザクションで行うためPoolを直接渡す。
	noteTemplateRepo := persistence.NewNoteTemplateRepository(sqlPool)
	// 学習ログテンプレートはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	// デフォルト解除とCreate/Updateを1トランザクションで行うためPoolを直接渡す。
	learningLogTemplateRepo := persistence.NewLearningLogTemplateRepository(sqlPool)
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
	authUseCases := AuthUseCases{
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
	githubUseCases := GitHubUseCases{
		OAuthURL:      usecase.NewGetGitHubOAuthURLUseCase(githubClient),
		Connect:       usecase.NewConnectGitHubUseCase(userPort, githubClient, syncGitHubData),
		Disconnect:    usecase.NewDisconnectGitHubUseCase(userPort, githubPort),
		Sync:          syncGitHubData,
		Contributions: usecase.NewGetGitHubContributionsUseCase(githubPort),
		Languages:     usecase.NewGetGitHubLanguagesUseCase(githubPort),
		Repos:         usecase.NewGetGitHubReposUseCase(githubPort),
	}
	authGitHubUseCases := AuthGitHubUseCases{
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
	learningLogPort := persistence.NewLearningLogRepository(sqlPool)
	createLearningLog := usecase.NewCreateLearningLogUseCase(learningLogPort, learningGoalPort)

	// テンプレートロードマップの初期登録
	go seedTemplateRoadmaps(authUserPort, usecase.NewSeedRoadmapTemplatesUseCase(roadmapPort))

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
	aiUseCases := AIAdviceUseCases{
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
	// SMTP 未設定のときは送信 usecase を組み立てず nil のままにする
	// （scheduler.Start は emailSvc が nil ならそのジョブだけ登録しない）。
	// scheduler.WeeklyReportSender インターフェース型で宣言する（具象型の nil ポインタを
	// インターフェースへ代入すると非 nil インターフェースになりパニックの原因になるため）。
	var sendWeeklyReports scheduler.WeeklyReportSender
	if cfg.SMTPHost != "" {
		log.Println("SMTP設定が検出されました。メール機能が有効です。")
		sendWeeklyReports = usecase.NewSendAllWeeklyReportsUseCase(
			userPort,
			activityReportRepo,
			usecase.NewSendWeeklyReportUseCase(external.NewSMTPEmailSender(cfg), cfg.AppURL),
		)
	} else {
		log.Println("SMTP未設定。ウィークリーレポートメールのスケジューラは無効です。")
	}
	// 各ドメインのカウンタ（post_metrics / learning_resource_metrics）のreconcileは
	// SMTP設定に関わらず常に有効。port は usecase/repository、sqlc(pgx) 実装は
	// adapter/persistence。ドメインが増えてもscheduler.Newへの登録ポイントは
	// ReconcileAllMetricsUseCase 1つのままにするため、ここでまとめる。
	postMetricsPort := persistence.NewPostMetricsRepository(sqlcgen.New(sqlPool))
	learningResourceMetricsPort := persistence.NewLearningResourceMetricsRepository(sqlcgen.New(sqlPool))
	reconcileAllMetrics := usecase.NewReconcileAllMetricsUseCase(
		usecase.NewReconcilePostMetricsUseCase(postMetricsPort),
		usecase.NewReconcileLearningResourceMetricsUseCase(learningResourceMetricsPort),
	)
	appScheduler := scheduler.New(sendWeeklyReports, reconcileAllMetrics)
	go appScheduler.Start()

	// ハンドラ
	origins := cfg.CORSOrigins
	c.AuthHandler = NewAuthHandler(authUseCases, authGitHubUseCases)
	c.UserHandler = NewUserHandler(
		usecase.NewListUsersUseCase(userPort),
		usecase.NewGetUserUseCase(userPort),
		usecase.NewGetUserByUsernameUseCase(userPort),
		usecase.NewUpdateUserProfileUseCase(userPort),
		usecase.NewGetProfileCompletenessUseCase(userPort),
	)
	c.FollowHandler = NewFollowHandler(
		usecase.NewFollowUserUseCase(followRepo, notificationCreator),
		usecase.NewUnfollowUserUseCase(followRepo),
		usecase.NewListFollowersUseCase(followRepo),
		usecase.NewListFollowingUseCase(followRepo),
	)
	c.GitHubHandler = NewGitHubHandler(githubUseCases, githubOAuthState)
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
	c.PostHandler = NewPostHandler(PostUseCases{
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
	c.CodeSnippetHandler = NewCodeSnippetHandler(
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
	c.RankingHandler = NewRankingHandler(
		usecase.NewGetContributionRankingUseCase(rankingRepo),
		usecase.NewGetLanguageRankingUseCase(rankingRepo),
		usecase.NewGetLevelRankingUseCase(rankingRepo),
		usecase.NewListRankingLanguagesUseCase(rankingRepo),
	)
	c.MessageHandler = NewMessageHandler(
		usecase.NewListConversationsUseCase(messagePort),
		usecase.NewGetConversationUseCase(messagePort),
		usecase.NewSendMessageUseCase(messagePort, notificationCreator),
		usecase.NewMarkMessagesAsReadUseCase(messagePort),
	)
	c.WebSocketHandler = NewWebSocketHandler(hub, validateAuthToken, parseOrigins(origins))
	uploadHandler, err := NewUploadHandler()
	if err != nil {
		log.Fatalf("アップロードハンドラの初期化に失敗: %v", err)
	}
	c.UploadHandler = uploadHandler
	// 通知の参照・既読・削除はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	// 通知の作成（WebSocket 配信を含む）は post / badge / level / mention / message がまだ service 経由で使うため残している。
	notificationPort := persistence.NewNotificationRepository(sqlPool)
	c.NotificationHandler = NewNotificationHandler(
		usecase.NewListNotificationsUseCase(notificationPort),
		usecase.NewCountUnreadNotificationsUseCase(notificationPort),
		usecase.NewMarkNotificationAsReadUseCase(notificationPort),
		usecase.NewMarkAllNotificationsAsReadUseCase(notificationPort),
		usecase.NewDeleteNotificationUseCase(notificationPort),
	)
	c.ZennHandler = NewArticlePlatformHandler("Zenn", ArticlePlatformOps[model.ZennArticle, model.ZennStats]{
		Connect:     connectZenn.Execute,
		Disconnect:  disconnectZenn.Execute,
		Sync:        syncZenn.Execute,
		GetArticles: usecase.NewListZennArticlesUseCase(zennPort).Execute,
		GetStats:    usecase.NewGetZennStatsUseCase(zennPort).Execute,
	})
	c.QiitaHandler = NewArticlePlatformHandler("Qiita", ArticlePlatformOps[model.QiitaArticle, model.QiitaStats]{
		Connect:     connectQiita.Execute,
		Disconnect:  disconnectQiita.Execute,
		Sync:        syncQiita.Execute,
		GetArticles: usecase.NewListQiitaArticlesUseCase(qiitaPort).Execute,
		GetStats:    usecase.NewGetQiitaStatsUseCase(qiitaPort).Execute,
	})
	c.LearningGoalHandler = NewLearningGoalHandler(
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
	c.ActivityReportHandler = NewActivityReportHandler(
		usecase.NewGetWeeklyActivityReportUseCase(activityReportRepo),
		usecase.NewGetMonthlyActivityReportUseCase(activityReportRepo),
		usecase.NewGetActivityReportComparisonUseCase(activityReportRepo),
	)
	c.ProjectHandler = NewProjectHandler(
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
	c.LearningResourceHandler = NewLearningResourceHandler(
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
	c.BookReviewHandler = NewBookReviewHandler(
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
	c.QuestionHandler = NewQuestionHandler(
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
	c.AnswerHandler = NewAnswerHandler(
		usecase.NewListAnswersUseCase(answerPort),
		usecase.NewCreateAnswerUseCase(answerPort, questionPort),
		usecase.NewUpdateAnswerUseCase(answerPort),
		usecase.NewDeleteAnswerUseCase(answerPort),
		usecase.NewSetBestAnswerUseCase(answerPort, questionPort),
		usecase.NewVoteAnswerUseCase(answerPort),
		usecase.NewRemoveAnswerVoteUseCase(answerPort),
		usecase.NewListAnswersByVoteRangeUseCase(answerPort),
	)
	c.RoadmapHandler = NewRoadmapHandler(
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
	c.ChatRoomHandler = NewChatRoomHandler(
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
	c.AtCoderHandler = NewAtCoderHandler(
		usecase.NewGetAtCoderRatingUseCase(atcoderClient),
		usecase.NewConnectAtCoderUseCase(userPort, atcoderClient),
		usecase.NewDisconnectAtCoderUseCase(userPort),
	)
	c.BadgeHandler = NewBadgeHandler(
		usecase.NewGetUserBadgesUseCase(badgePort),
		usecase.NewNotifyBadgeEarnedUseCase(notificationCreator),
	)
	c.LearningLogHandler = NewLearningLogHandler(
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
	c.AIAdviceHandler = NewAIAdviceHandler(aiUseCases)
	c.EmailPreferencesHandler = NewEmailPreferencesHandler(
		usecase.NewGetEmailPreferencesUseCase(userPort),
		usecase.NewUpdateEmailPreferencesUseCase(userPort),
	)
	c.LevelHandler = NewLevelHandler(
		usecase.NewGetLevelInfoUseCase(levelPort),
		usecase.NewGetXPBreakdownUseCase(levelPort),
	)
	c.LearningAnalyticsHandler = NewLearningAnalyticsHandler(
		usecase.NewGetLearningHeatmapUseCase(analyticsPort),
		usecase.NewGetCategoryBreakdownUseCase(analyticsPort),
		usecase.NewGetWeeklyTrendsUseCase(analyticsPort),
		usecase.NewGetDayOfWeekSummaryUseCase(analyticsPort),
		usecase.NewGetProductivityScoreUseCase(analyticsPort),
		usecase.NewGetLearningInsightsUseCase(analyticsPort),
	)
	c.RecommendationHandler = NewRecommendationHandler(
		usecase.NewGetRecommendedUsersUseCase(recommendationRepo, userPort),
		usecase.NewGetTrendingPostsUseCase(recommendationRepo),
		usecase.NewGetTrendingResourcesUseCase(recommendationRepo),
	)
	c.StudyCircleHandler = NewStudyCircleHandler(
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
	c.SearchHandler = NewSearchHandler(
		usecase.NewSearchPostsUseCase(persistence.NewPostSearchRepository(sqlcgen.New(sqlPool))),
		searchStudyCircles,
	)
	c.NoteHandler = NewNoteHandler(
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
	c.NoteVersionHandler = NewNoteVersionHandler(
		usecase.NewListNoteVersionsUseCase(noteVersionRepo, noteReader),
		usecase.NewGetNoteVersionUseCase(noteVersionRepo, noteReader),
		usecase.NewRestoreNoteVersionUseCase(noteVersionRepo, noteReader, noteUpdater),
	)
	c.NoteFolderHandler = NewNoteFolderHandler(
		usecase.NewCreateNoteFolderUseCase(noteFolderRepo),
		usecase.NewGetNoteFolderUseCase(noteFolderRepo),
		usecase.NewListNoteFoldersUseCase(noteFolderRepo),
		usecase.NewListChildNoteFoldersUseCase(noteFolderRepo),
		usecase.NewListRootNoteFoldersUseCase(noteFolderRepo),
		usecase.NewUpdateNoteFolderUseCase(noteFolderRepo),
		usecase.NewCountNoteFoldersUseCase(noteFolderRepo),
		usecase.NewDeleteNoteFolderUseCase(noteFolderRepo),
	)
	c.NoteTemplateHandler = NewNoteTemplateHandler(
		usecase.NewCreateNoteTemplateUseCase(noteTemplateRepo),
		usecase.NewGetNoteTemplateUseCase(noteTemplateRepo),
		usecase.NewListNoteTemplatesUseCase(noteTemplateRepo),
		usecase.NewGetDefaultNoteTemplateUseCase(noteTemplateRepo),
		usecase.NewUpdateNoteTemplateUseCase(noteTemplateRepo),
		usecase.NewDeleteNoteTemplateUseCase(noteTemplateRepo),
		usecase.NewCreateNoteFromTemplateUseCase(noteTemplateRepo, createNote),
		usecase.NewCountNoteTemplatesUseCase(noteTemplateRepo),
	)
	c.LearningLogTemplateHandler = NewLearningLogTemplateHandler(
		usecase.NewCreateLearningLogTemplateUseCase(learningLogTemplateRepo),
		usecase.NewGetLearningLogTemplateUseCase(learningLogTemplateRepo),
		usecase.NewListLearningLogTemplatesUseCase(learningLogTemplateRepo),
		usecase.NewGetDefaultLearningLogTemplateUseCase(learningLogTemplateRepo),
		usecase.NewUpdateLearningLogTemplateUseCase(learningLogTemplateRepo),
		usecase.NewDeleteLearningLogTemplateUseCase(learningLogTemplateRepo),
		usecase.NewCreateLearningLogFromTemplateUseCase(learningLogTemplateRepo, createLearningLog),
		usecase.NewCountLearningLogTemplatesUseCase(learningLogTemplateRepo),
	)
	c.NoteLinkHandler = NewNoteLinkHandler(
		usecase.NewCreateNoteLinkUseCase(noteLinkRepo, noteReader),
		usecase.NewListNoteLinksUseCase(noteLinkRepo),
		usecase.NewListNoteBacklinksUseCase(noteLinkRepo),
		usecase.NewGetNoteLinkStatsUseCase(noteLinkRepo, noteReader),
		usecase.NewDeleteNoteLinkUseCase(noteLinkRepo, noteReader),
	)
	c.PostSeriesHandler = NewPostSeriesHandler(
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
	c.PostCollectionHandler = NewPostCollectionHandler(
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
	c.PostTagHandler = NewPostTagHandler(
		setPostTags,
		usecase.NewGetPostTagsUseCase(postTagRepo),
		usecase.NewFindPostsByTagUseCase(postTagRepo),
		usecase.NewGetPopularTagsUseCase(postTagRepo),
	)
	c.PostHandler.SetAutoTagsUseCase(usecase.NewSetAutoPostTagsUseCase(setPostTags))

	// post_pin はクリーンアーキテクチャ(DIP)へ移行済み。
	postPinRepo := persistence.NewPostPinRepository(sqlPool)
	postReader := persistence.NewPostReader(sqlcgen.New(sqlPool))
	c.PostPinHandler = NewPostPinHandler(
		usecase.NewPinPostUseCase(postPinRepo, postReader),
		usecase.NewUnpinPostUseCase(postPinRepo),
		usecase.NewListPinnedPostsUseCase(postPinRepo),
		usecase.NewCountPinnedPostsUseCase(postPinRepo),
		usecase.NewReorderPinnedPostsUseCase(postPinRepo),
	)

	// comment_like はクリーンアーキテクチャ(DIP)へ移行済み。
	commentLikeRepo := persistence.NewCommentLikeRepository(sqlPool)
	commentReader := persistence.NewCommentReader(sqlcgen.New(sqlPool))
	c.CommentLikeHandler = NewCommentLikeHandler(
		usecase.NewLikeCommentUseCase(commentLikeRepo, commentReader),
		usecase.NewUnlikeCommentUseCase(commentLikeRepo, commentReader),
		usecase.NewGetCommentLikeStatusUseCase(commentLikeRepo, commentReader),
	)

	// 投稿閲覧数はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	postViewRepo := persistence.NewPostViewRepository(sqlPool)
	c.PostViewHandler = NewPostViewHandler(
		usecase.NewRecordPostViewUseCase(postViewRepo),
		usecase.NewGetPostViewCountUseCase(postViewRepo),
		usecase.NewGetMostViewedPostsUseCase(postViewRepo),
	)

	// メンションはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	// 記録側（投稿・コメント作成時の処理）は投稿スライスの配線でつないでいる。
	c.MentionHandler = NewMentionHandler(
		usecase.NewListUserMentionsUseCase(mentionPort),
		usecase.NewListPostMentionsUseCase(mentionPort),
	)

	// ユーザーダッシュボード統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	userDashboardRepo := persistence.NewUserDashboardRepository(sqlcgen.New(sqlPool))
	c.UserDashboardHandler = NewUserDashboardHandler(
		usecase.NewGetUserDashboardStatsUseCase(userDashboardRepo),
	)

	// ノート統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	noteStatsRepo := persistence.NewNoteStatsRepository(sqlcgen.New(sqlPool))
	c.NoteStatsHandler = NewNoteStatsHandler(
		usecase.NewGetNoteStatsUseCase(noteStatsRepo),
	)

	// スタディサークル統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	studyCircleStatsRepo := persistence.NewStudyCircleStatsRepository(sqlcgen.New(sqlPool))
	c.StudyCircleStatsHandler = NewStudyCircleStatsHandler(
		usecase.NewGetStudyCircleStatsUseCase(studyCircleStatsRepo),
	)

	// 投稿統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	postStatsRepo := persistence.NewPostStatsRepository(sqlcgen.New(sqlPool))
	c.PostStatsHandler = NewPostStatsHandler(
		usecase.NewGetPostStatsUseCase(postStatsRepo),
	)

	// 書籍レビュー統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	bookReviewStatsRepo := persistence.NewBookReviewStatsRepository(sqlcgen.New(sqlPool))
	c.BookReviewStatsHandler = NewBookReviewStatsHandler(
		usecase.NewGetBookReviewStatsUseCase(bookReviewStatsRepo),
	)

	// Q&A 統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	qaStatsRepo := persistence.NewQAStatsRepository(sqlcgen.New(sqlPool))
	c.QAStatsHandler = NewQAStatsHandler(
		usecase.NewGetQAStatsUseCase(qaStatsRepo),
	)

	// コードスニペット統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	codeSnippetStatsRepo := persistence.NewCodeSnippetStatsRepository(sqlcgen.New(sqlPool))
	c.CodeSnippetStatsHandler = NewCodeSnippetStatsHandler(
		usecase.NewGetCodeSnippetStatsUseCase(codeSnippetStatsRepo),
	)

	// 学習リソース統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	learningResourceStatsRepo := persistence.NewLearningResourceStatsRepository(sqlcgen.New(sqlPool))
	c.LearningResourceStatsHandler = NewLearningResourceStatsHandler(
		usecase.NewGetLearningResourceStatsUseCase(learningResourceStatsRepo),
	)

	// プロジェクト統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	projectStatsRepo := persistence.NewProjectStatsRepository(sqlcgen.New(sqlPool))
	c.ProjectStatsHandler = NewProjectStatsHandler(
		usecase.NewGetProjectStatsUseCase(projectStatsRepo),
	)

	// フォロー統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	followStatsRepo := persistence.NewFollowStatsRepository(sqlcgen.New(sqlPool))
	c.FollowStatsHandler = NewFollowStatsHandler(
		usecase.NewGetFollowStatsUseCase(followStatsRepo),
	)

	// ロードマップ統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	roadmapStatsRepo := persistence.NewRoadmapStatsRepository(sqlcgen.New(sqlPool))
	c.RoadmapStatsHandler = NewRoadmapStatsHandler(
		usecase.NewGetRoadmapStatsUseCase(roadmapStatsRepo),
	)

	// 学習ログ統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	learningLogStatsRepo := persistence.NewLearningLogStatsRepository(sqlcgen.New(sqlPool))
	c.LearningLogStatsHandler = NewLearningLogStatsHandler(
		usecase.NewGetLearningLogStatsUseCase(learningLogStatsRepo),
	)

	// コメント統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	commentStatsRepo := persistence.NewCommentStatsRepository(sqlcgen.New(sqlPool))
	c.CommentStatsHandler = NewCommentStatsHandler(
		usecase.NewGetCommentStatsUseCase(commentStatsRepo),
	)

	// 通知統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	notificationStatsRepo := persistence.NewNotificationStatsRepository(sqlcgen.New(sqlPool))
	c.NotificationStatsHandler = NewNotificationStatsHandler(
		usecase.NewGetNotificationStatsUseCase(notificationStatsRepo),
	)

	// メッセージ統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	messageStatsRepo := persistence.NewMessageStatsRepository(sqlcgen.New(sqlPool))
	c.MessageStatsHandler = NewMessageStatsHandler(
		usecase.NewGetMessageStatsUseCase(messageStatsRepo),
	)

	// メンション統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	mentionStatsRepo := persistence.NewMentionStatsRepository(sqlcgen.New(sqlPool))
	c.MentionStatsHandler = NewMentionStatsHandler(
		usecase.NewGetMentionStatsUseCase(mentionStatsRepo),
	)

	// リアクション統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	reactionStatsRepo := persistence.NewReactionStatsRepository(sqlcgen.New(sqlPool))
	c.ReactionStatsHandler = NewReactionStatsHandler(
		usecase.NewGetReactionStatsUseCase(reactionStatsRepo),
		usecase.NewGetReactionSummaryUseCase(reactionStatsRepo),
	)

	// ブックマーク統計はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	bookmarkStatsRepo := persistence.NewBookmarkStatsRepository(sqlcgen.New(sqlPool))
	c.BookmarkStatsHandler = NewBookmarkStatsHandler(
		usecase.NewGetBookmarkStatsUseCase(bookmarkStatsRepo),
	)

	// Spotify 連携はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、
	// 実装は adapter/persistence（永続化）と adapter/external（Spotify API）。
	spotifyPort := persistence.NewSpotifyRepository(sqlcgen.New(sqlPool))
	spotifyClient := external.NewSpotifyClient(cfg.SpotifyClientID, cfg.SpotifyClientSecret, cfg.SpotifyRedirectURL)
	spotifyUseCases := SpotifyUseCases{
		OAuthURL:         usecase.NewGetSpotifyOAuthURLUseCase(spotifyClient),
		Connect:          usecase.NewConnectSpotifyUseCase(userPort, spotifyClient),
		Disconnect:       usecase.NewDisconnectSpotifyUseCase(userPort, spotifyPort),
		CurrentlyPlaying: usecase.NewGetSpotifyCurrentlyPlayingUseCase(userPort, spotifyClient),
		RecentlyPlayed:   usecase.NewGetSpotifyRecentlyPlayedUseCase(userPort, spotifyClient),
	}
	c.SpotifyHandler = NewSpotifyHandler(spotifyUseCases, spotifyOAuthState)

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
	c.YouTubeHandler = NewYouTubeHandler(
		usecase.NewSearchYouTubeVideosUseCase(youtubeVideoPort, youtubeSearcher),
		usecase.NewRecommendYouTubeVideosUseCase(userPort, youtubeVideoPort, youtubeSearcher),
		usecase.NewCheckYouTubeAvailabilityUseCase(youtubeSearcher),
	)

	// ストリークフリーズサービス
	// ストリークフリーズはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	streakFreezeRepo := persistence.NewStreakFreezeRepository(sqlPool)
	c.StreakFreezeHandler = NewStreakFreezeHandler(
		usecase.NewUseStreakFreezeUseCase(streakFreezeRepo),
		usecase.NewGetStreakFreezeStatusUseCase(streakFreezeRepo),
	)

	// ブックマークコレクションサービス
	// ブックマークコレクションはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	bookmarkCollectionRepo := persistence.NewBookmarkCollectionRepository(sqlcgen.New(sqlPool))
	c.BookmarkCollectionHandler = NewBookmarkCollectionHandler(
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
	c.WeeklyChallengeHandler = NewWeeklyChallengeHandler(
		usecase.NewGetCurrentWeeklyChallengeUseCase(weeklyChallengeRepo),
		usecase.NewUpdateWeeklyChallengeProgressUseCase(weeklyChallengeRepo),
	)

	// 投稿テンプレートはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence（sqlc/pgx）。
	postTemplateRepo := persistence.NewPostTemplateRepository(sqlcgen.New(sqlPool))
	c.PostTemplateHandler = NewPostTemplateHandler(
		usecase.NewCreatePostTemplateUseCase(postTemplateRepo),
		usecase.NewGetPostTemplateUseCase(postTemplateRepo),
		usecase.NewListPostTemplatesUseCase(postTemplateRepo),
		usecase.NewUpdatePostTemplateUseCase(postTemplateRepo),
		usecase.NewDeletePostTemplateUseCase(postTemplateRepo),
	)

	// ウィジェット設定はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	widgetSettingsRepo := persistence.NewWidgetSettingsRepository(sqlcgen.New(sqlPool))
	c.WidgetSettingsHandler = NewWidgetSettingsHandler(
		usecase.NewGetWidgetSettingsUseCase(widgetSettingsRepo),
		usecase.NewUpdateWidgetSettingsUseCase(widgetSettingsRepo),
	)

	// カテゴリ別週間学習目標はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	weeklyGoalRepo := persistence.NewWeeklyGoalRepository(sqlcgen.New(sqlPool))
	c.WeeklyGoalHandler = NewWeeklyGoalHandler(
		usecase.NewSetWeeklyGoalUseCase(weeklyGoalRepo),
		usecase.NewListWeeklyGoalsUseCase(weeklyGoalRepo),
		usecase.NewGetWeeklyGoalProgressUseCase(weeklyGoalRepo),
	)

	// リソース進捗サービス
	// リソース進捗はクリーンアーキテクチャ（DIP）へ移行済み。リソース存在確認は最小 port LearningResourceReader を再利用。
	resourceProgressRepo := persistence.NewResourceProgressRepository(sqlcgen.New(sqlPool))
	resourceProgressResourceReader := persistence.NewLearningResourceReader(sqlcgen.New(sqlPool))
	c.ResourceProgressHandler = NewResourceProgressHandler(
		usecase.NewUpsertResourceProgressUseCase(resourceProgressRepo, resourceProgressResourceReader),
		usecase.NewGetResourceProgressUseCase(resourceProgressRepo),
		usecase.NewListResourceProgressUseCase(resourceProgressRepo),
	)

	// プロジェクトマイルストーンサービス
	// プロジェクトマイルストーンはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	projectMilestoneRepo := persistence.NewProjectMilestoneRepository(sqlcgen.New(sqlPool))
	projectReader := persistence.NewProjectReader(sqlcgen.New(sqlPool))
	c.ProjectMilestoneHandler = NewProjectMilestoneHandler(
		usecase.NewCreateProjectMilestoneUseCase(projectMilestoneRepo, projectReader),
		usecase.NewListProjectMilestonesUseCase(projectMilestoneRepo),
		usecase.NewUpdateProjectMilestoneUseCase(projectMilestoneRepo, projectReader),
		usecase.NewDeleteProjectMilestoneUseCase(projectMilestoneRepo, projectReader),
	)

	// ユーザーアクティビティサービス
	// ユーザーアクティビティはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	userActivityRepo := persistence.NewUserActivityRepository(sqlcgen.New(sqlPool))
	c.UserActivityHandler = NewUserActivityHandler(
		usecase.NewGetActivityTimelineUseCase(userActivityRepo),
	)

	// リソースレビューはクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	resourceReviewRepo := persistence.NewResourceReviewRepository(sqlcgen.New(sqlPool))
	learningResourceReader := persistence.NewLearningResourceReader(sqlcgen.New(sqlPool))
	c.ResourceReviewHandler = NewResourceReviewHandler(
		usecase.NewCreateResourceReviewUseCase(resourceReviewRepo, learningResourceReader),
		usecase.NewListResourceReviewsUseCase(resourceReviewRepo),
		usecase.NewUpdateResourceReviewUseCase(resourceReviewRepo),
		usecase.NewDeleteResourceReviewUseCase(resourceReviewRepo),
	)

	// 学習ダッシュボード統合サマリーサービス
	// 学習ダッシュボードはクリーンアーキテクチャ（DIP）へ移行済み。
	// 学習ログ・目標・分析の移行済み実装を、それぞれ最小 port として受け取る。
	c.LearningDashboardHandler = NewLearningDashboardHandler(
		usecase.NewGetLearningDashboardSummaryUseCase(learningLogPort, learningGoalPort, analyticsPort),
	)

	// リマインダー設定はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	reminderSettingsRepo := persistence.NewReminderSettingsRepository(sqlcgen.New(sqlPool))
	c.ReminderSettingsHandler = NewReminderSettingsHandler(
		usecase.NewGetReminderSettingsUseCase(reminderSettingsRepo),
		usecase.NewUpdateReminderSettingsUseCase(reminderSettingsRepo),
	)

	// 通知設定サービス
	// 通知設定はクリーンアーキテクチャ（DIP）へ移行済み。port は usecase/repository、実装は adapter/persistence。
	notificationSettingsRepo := persistence.NewNotificationSettingsRepository(sqlcgen.New(sqlPool))
	c.NotificationSettingsHandler = NewNotificationSettingsHandler(
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
func seedTemplateRoadmaps(authUserPort usecaserepo.AuthUserRepository, seed *usecase.SeedRoadmapTemplatesUseCase) {
	const systemEmail = "system@devsync.local"
	ctx := context.Background()
	user, err := authUserPort.FindByEmail(ctx, systemEmail)
	if err != nil {
		log.Printf("テンプレートシード用システムユーザー取得失敗: %v", err)
		return
	}
	if user == nil {
		user = &model.User{
			Name:           "DevSync System",
			Email:          systemEmail,
			Username:       "__system__",
			GitHubID:       -1,
			GitHubUsername: "__system__",
		}
		if err := authUserPort.Create(ctx, user); err != nil {
			log.Printf("テンプレートシード用システムユーザー作成失敗: %v", err)
			return
		}
	}
	if err := seed.Execute(ctx, user.ID); err != nil {
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

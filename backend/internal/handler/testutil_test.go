package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// ---------- テストユーティリティ ----------

// authMiddleware はテスト用認証ミドルウェア。userIDをコンテキストに設定する。
func authMiddleware(userID uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	}
}

// jsonBody はリクエストボディ用のJSONバッファを生成する。
func jsonBody(v interface{}) *bytes.Buffer {
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}

// parseJSON はレスポンスボディをmapにパースする。
func parseJSON(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var result map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err, "JSON parse failed: %s", w.Body.String())
	return result
}

// assertStatus はHTTPステータスコードを検証する。
func assertStatus(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	t.Helper()
	assert.Equal(t, expected, w.Code, "body: %s", w.Body.String())
}

// ---------- モックリポジトリ ----------

// MockPostRepository は PostRepositoryInterface のモック実装。
type MockPostRepository struct{ mock.Mock }

func (m *MockPostRepository) Create(post *model.Post) error {
	args := m.Called(post)
	return args.Error(0)
}
func (m *MockPostRepository) FindByID(id uint) (*model.Post, error) {
	args := m.Called(id)
	if p := args.Get(0); p != nil {
		return p.(*model.Post), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPostRepository) FindAll(page, limit int) ([]model.Post, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]model.Post), args.Error(1)
}
func (m *MockPostRepository) CountAll() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockPostRepository) FindByUserID(userID uint, limit, offset int) ([]model.Post, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.Post), args.Get(1).(int64), args.Error(2)
}
func (m *MockPostRepository) Timeline(userID uint, page, limit int) ([]model.Post, error) {
	args := m.Called(userID, page, limit)
	return args.Get(0).([]model.Post), args.Error(1)
}
func (m *MockPostRepository) Update(post *model.Post) error {
	return m.Called(post).Error(0)
}
func (m *MockPostRepository) Delete(id uint) error {
	return m.Called(id).Error(0)
}
func (m *MockPostRepository) Like(userID, postID uint) error {
	return m.Called(userID, postID).Error(0)
}
func (m *MockPostRepository) Unlike(userID, postID uint) error {
	return m.Called(userID, postID).Error(0)
}
func (m *MockPostRepository) HasLiked(userID, postID uint) bool {
	return m.Called(userID, postID).Bool(0)
}
func (m *MockPostRepository) AddReaction(userID, postID uint, emoji string) error {
	return m.Called(userID, postID, emoji).Error(0)
}
func (m *MockPostRepository) RemoveReaction(userID, postID uint, emoji string) error {
	return m.Called(userID, postID, emoji).Error(0)
}
func (m *MockPostRepository) GetReactionsByPostID(postID uint) ([]model.ReactionCount, error) {
	args := m.Called(postID)
	return args.Get(0).([]model.ReactionCount), args.Error(1)
}
func (m *MockPostRepository) GetUserReactions(userID, postID uint) ([]string, error) {
	args := m.Called(userID, postID)
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockPostRepository) GetReactionsBatch(postIDs []uint) (map[uint][]model.ReactionCount, error) {
	args := m.Called(postIDs)
	return args.Get(0).(map[uint][]model.ReactionCount), args.Error(1)
}
func (m *MockPostRepository) GetUserReactionsBatch(userID uint, postIDs []uint) (map[uint][]string, error) {
	args := m.Called(userID, postIDs)
	return args.Get(0).(map[uint][]string), args.Error(1)
}
func (m *MockPostRepository) Search(query string, limit, offset int) (interface{}, int64, error) {
	args := m.Called(query, limit, offset)
	return args.Get(0), args.Get(1).(int64), args.Error(2)
}
func (m *MockPostRepository) FindDraftsByUserID(userID uint) ([]model.Post, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Post), args.Error(1)
}
func (m *MockPostRepository) FindScheduledByUserID(userID uint) ([]model.Post, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Post), args.Error(1)
}
func (m *MockPostRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockPostRepository) CountDraftsByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockPostRepository) CountScheduledByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// MockNotificationRepository は NotificationRepositoryInterface のモック実装。
type MockNotificationRepository struct{ mock.Mock }

func (m *MockNotificationRepository) Create(n *model.Notification) error {
	return m.Called(n).Error(0)
}
func (m *MockNotificationRepository) CreateBatch(n []*model.Notification) error {
	return m.Called(n).Error(0)
}
func (m *MockNotificationRepository) FindByUserID(userID uint, page, limit int, t string) ([]model.Notification, error) {
	args := m.Called(userID, page, limit, t)
	return args.Get(0).([]model.Notification), args.Error(1)
}
func (m *MockNotificationRepository) CountByUserID(userID uint, t string) (int64, error) {
	args := m.Called(userID, t)
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockNotificationRepository) CountUnread(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockNotificationRepository) MarkAsRead(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}
func (m *MockNotificationRepository) MarkAllAsRead(userID uint) error {
	return m.Called(userID).Error(0)
}
func (m *MockNotificationRepository) Delete(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}
func (m *MockNotificationRepository) GetFollowerIDs(userID uint) ([]uint, error) {
	args := m.Called(userID)
	return args.Get(0).([]uint), args.Error(1)
}

// (MockCodeSnippetRepository は code_snippet の DIP 移行に伴い撤去)

// mockQuestionRepo は usecase/repository.QuestionRepository のモック（ctx 付き）。
type mockQuestionRepo struct{ mock.Mock }

func (m *mockQuestionRepo) Create(ctx context.Context, q *model.Question) error {
	return m.Called(ctx, q).Error(0)
}
func (m *mockQuestionRepo) FindByID(ctx context.Context, id uint) (*model.Question, error) {
	args := m.Called(ctx, id)
	if q := args.Get(0); q != nil {
		return q.(*model.Question), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockQuestionRepo) Update(ctx context.Context, q *model.Question) error {
	return m.Called(ctx, q).Error(0)
}
func (m *mockQuestionRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockQuestionRepo) FindAll(ctx context.Context, limit, offset int, tag, sort string) ([]model.Question, int64, error) {
	args := m.Called(ctx, limit, offset, tag, sort)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}
func (m *mockQuestionRepo) Search(ctx context.Context, q string, limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(ctx, q, limit, offset)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}
func (m *mockQuestionRepo) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}
func (m *mockQuestionRepo) FindSolved(ctx context.Context, limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}
func (m *mockQuestionRepo) FindUnanswered(ctx context.Context, limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}
func (m *mockQuestionRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockQuestionRepo) Vote(ctx context.Context, userID, questionID uint, value int) error {
	return m.Called(ctx, userID, questionID, value).Error(0)
}
func (m *mockQuestionRepo) RemoveVote(ctx context.Context, userID, questionID uint) error {
	return m.Called(ctx, userID, questionID).Error(0)
}
func (m *mockQuestionRepo) GetUserVote(ctx context.Context, userID, questionID uint) (int, error) {
	args := m.Called(ctx, userID, questionID)
	return args.Int(0), args.Error(1)
}
func (m *mockQuestionRepo) Bookmark(ctx context.Context, userID, questionID uint) error {
	return m.Called(ctx, userID, questionID).Error(0)
}
func (m *mockQuestionRepo) Unbookmark(ctx context.Context, userID, questionID uint) error {
	return m.Called(ctx, userID, questionID).Error(0)
}
func (m *mockQuestionRepo) HasBookmarked(ctx context.Context, userID, questionID uint) (bool, error) {
	args := m.Called(ctx, userID, questionID)
	return args.Bool(0), args.Error(1)
}
func (m *mockQuestionRepo) FindBookmarkedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}

// mockLearningResourceRepo は usecase/repository.LearningResourceRepository のモック（ctx 付き）。
type mockLearningResourceRepo struct{ mock.Mock }

func (m *mockLearningResourceRepo) Create(ctx context.Context, r *model.LearningResource) error {
	return m.Called(ctx, r).Error(0)
}
func (m *mockLearningResourceRepo) Update(ctx context.Context, r *model.LearningResource) error {
	return m.Called(ctx, r).Error(0)
}
func (m *mockLearningResourceRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockLearningResourceRepo) FindByID(ctx context.Context, id uint) (*model.LearningResource, error) {
	args := m.Called(ctx, id)
	if r := args.Get(0); r != nil {
		return r.(*model.LearningResource), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockLearningResourceRepo) FindByUserID(ctx context.Context, userID uint, includePrivate bool, limit, offset int) ([]model.LearningResource, int64, error) {
	args := m.Called(ctx, userID, includePrivate, limit, offset)
	r, _ := args.Get(0).([]model.LearningResource)
	return r, args.Get(1).(int64), args.Error(2)
}
func (m *mockLearningResourceRepo) FindPublic(ctx context.Context, limit, offset int, category, difficulty string) ([]model.LearningResource, int64, error) {
	args := m.Called(ctx, limit, offset, category, difficulty)
	r, _ := args.Get(0).([]model.LearningResource)
	return r, args.Get(1).(int64), args.Error(2)
}
func (m *mockLearningResourceRepo) FindByDifficulty(ctx context.Context, difficulty string, limit, offset int) ([]model.LearningResource, int64, error) {
	args := m.Called(ctx, difficulty, limit, offset)
	r, _ := args.Get(0).([]model.LearningResource)
	return r, args.Get(1).(int64), args.Error(2)
}
func (m *mockLearningResourceRepo) Search(ctx context.Context, query string, limit, offset int) ([]model.LearningResource, int64, error) {
	args := m.Called(ctx, query, limit, offset)
	r, _ := args.Get(0).([]model.LearningResource)
	return r, args.Get(1).(int64), args.Error(2)
}
func (m *mockLearningResourceRepo) FindSavedByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.LearningResource, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	r, _ := args.Get(0).([]model.LearningResource)
	return r, args.Get(1).(int64), args.Error(2)
}
func (m *mockLearningResourceRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockLearningResourceRepo) Like(ctx context.Context, userID, resourceID uint) error {
	return m.Called(ctx, userID, resourceID).Error(0)
}
func (m *mockLearningResourceRepo) Unlike(ctx context.Context, userID, resourceID uint) error {
	return m.Called(ctx, userID, resourceID).Error(0)
}
func (m *mockLearningResourceRepo) HasLiked(ctx context.Context, userID, resourceID uint) (bool, error) {
	args := m.Called(ctx, userID, resourceID)
	return args.Bool(0), args.Error(1)
}
func (m *mockLearningResourceRepo) Save(ctx context.Context, userID, resourceID uint) error {
	return m.Called(ctx, userID, resourceID).Error(0)
}
func (m *mockLearningResourceRepo) Unsave(ctx context.Context, userID, resourceID uint) error {
	return m.Called(ctx, userID, resourceID).Error(0)
}
func (m *mockLearningResourceRepo) HasSaved(ctx context.Context, userID, resourceID uint) (bool, error) {
	args := m.Called(ctx, userID, resourceID)
	return args.Bool(0), args.Error(1)
}

// mockRoadmapRepo は usecase/repository.RoadmapRepository のモック（ctx 付き）。
type mockRoadmapRepo struct{ mock.Mock }

func (m *mockRoadmapRepo) Create(ctx context.Context, r *model.Roadmap) error {
	return m.Called(ctx, r).Error(0)
}
func (m *mockRoadmapRepo) Update(ctx context.Context, r *model.Roadmap) error {
	return m.Called(ctx, r).Error(0)
}
func (m *mockRoadmapRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockRoadmapRepo) FindByID(ctx context.Context, id uint) (*model.Roadmap, error) {
	args := m.Called(ctx, id)
	if r := args.Get(0); r != nil {
		return r.(*model.Roadmap), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockRoadmapRepo) GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.Roadmap, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	r, _ := args.Get(0).([]model.Roadmap)
	return r, args.Get(1).(int64), args.Error(2)
}
func (m *mockRoadmapRepo) GetByStatus(ctx context.Context, userID uint, status string) ([]model.Roadmap, error) {
	args := m.Called(ctx, userID, status)
	r, _ := args.Get(0).([]model.Roadmap)
	return r, args.Error(1)
}
func (m *mockRoadmapRepo) GetPublicRoadmaps(ctx context.Context, limit, offset int) ([]model.Roadmap, int64, error) {
	args := m.Called(ctx, limit, offset)
	r, _ := args.Get(0).([]model.Roadmap)
	return r, args.Get(1).(int64), args.Error(2)
}
func (m *mockRoadmapRepo) GetTemplates(ctx context.Context) ([]model.Roadmap, error) {
	args := m.Called(ctx)
	r, _ := args.Get(0).([]model.Roadmap)
	return r, args.Error(1)
}
func (m *mockRoadmapRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockRoadmapRepo) CopyRoadmap(ctx context.Context, originalID, newUserID uint) (*model.Roadmap, error) {
	args := m.Called(ctx, originalID, newUserID)
	if r := args.Get(0); r != nil {
		return r.(*model.Roadmap), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockRoadmapRepo) CreateStep(ctx context.Context, step *model.RoadmapStep) error {
	return m.Called(ctx, step).Error(0)
}
func (m *mockRoadmapRepo) UpdateStep(ctx context.Context, step *model.RoadmapStep) error {
	return m.Called(ctx, step).Error(0)
}
func (m *mockRoadmapRepo) DeleteStep(ctx context.Context, stepID uint) error {
	return m.Called(ctx, stepID).Error(0)
}
func (m *mockRoadmapRepo) FindStepByID(ctx context.Context, stepID uint) (*model.RoadmapStep, error) {
	args := m.Called(ctx, stepID)
	if s := args.Get(0); s != nil {
		return s.(*model.RoadmapStep), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockRoadmapRepo) ReorderSteps(ctx context.Context, roadmapID uint, stepOrders []model.StepOrder) error {
	return m.Called(ctx, roadmapID, stepOrders).Error(0)
}

// Recommendation は DIP へ移行済み。テストは recommendation_test.go で
// 「本物の usecase + port モック」を組み立てる。

// ---------- リポジトリインターフェース適合チェック ----------
// import cycle を避けるため repository パッケージは使わないが、
// コンパイル時にインターフェース適合を検証したい場合は repository パッケージを
// import する必要がある。ここではハンドラ層テスト内でのみ使用する。

// ---------- ヘルパー関数 ----------

// postHandlerPorts は PostHandler に注入した port モックをまとめる。
type postHandlerPorts struct {
	Posts        *mockPostPort
	Likes        *mockPostLikePort
	Followers    *mockFollowerNotifierPort
	Snippets     *mockCodeSnippetRepo
	SnippetPosts *mockPostReader
	Reactions    *mockPostReactionPort
	Authors      *mockPostAuthorPort
	Comments     *mockPostCommentPort
	Bookmarks    *mockPostBookmarkPort
}

// setupPostHandler は PostHandler テスト用のセットアップを行う。
// 投稿スライスは DIP へ移行済みのため、本物の usecase に port モックを注入する。
func setupPostHandler() (*PostHandler, *postHandlerPorts) {
	ports := &postHandlerPorts{
		Posts:        new(mockPostPort),
		Likes:        new(mockPostLikePort),
		Followers:    new(mockFollowerNotifierPort),
		Snippets:     new(mockCodeSnippetRepo),
		SnippetPosts: new(mockPostReader),
		Reactions:    new(mockPostReactionPort),
		Authors:      new(mockPostAuthorPort),
		Comments:     new(mockPostCommentPort),
		Bookmarks:    new(mockPostBookmarkPort),
	}

	notifyFollowers := usecase.NewNotifyFollowersUseCase(ports.Followers)
	h := NewPostHandler(PostUseCases{
		Create:         usecase.NewCreatePostUseCase(ports.Posts, notifyFollowers),
		Get:            usecase.NewGetPostUseCase(ports.Posts),
		List:           usecase.NewListPostsUseCase(ports.Posts),
		Count:          usecase.NewCountPostsUseCase(ports.Posts),
		ListByUser:     usecase.NewListUserPostsUseCase(ports.Posts),
		ListDrafts:     usecase.NewListDraftPostsUseCase(ports.Posts),
		ListScheduled:  usecase.NewListScheduledPostsUseCase(ports.Posts),
		Timeline:       usecase.NewGetTimelineUseCase(ports.Posts),
		Update:         usecase.NewUpdatePostUseCase(ports.Posts),
		Delete:         usecase.NewDeletePostUseCase(ports.Posts),
		Publish:        usecase.NewPublishPostUseCase(ports.Posts, notifyFollowers),
		Unpublish:      usecase.NewUnpublishPostUseCase(ports.Posts),
		Schedule:       usecase.NewSchedulePostPublishUseCase(ports.Posts),
		CancelSchedule: usecase.NewCancelPostScheduleUseCase(ports.Posts),
		AutoSaveDraft:  usecase.NewAutoSaveDraftUseCase(ports.Posts),
		CountByUser:    usecase.NewCountUserPostsUseCase(ports.Posts),
		CountDrafts:    usecase.NewCountUserDraftsUseCase(ports.Posts),
		CountScheduled: usecase.NewCountUserScheduledPostsUseCase(ports.Posts),

		Like:     usecase.NewLikePostUseCase(ports.Likes, ports.Authors),
		Unlike:   usecase.NewUnlikePostUseCase(ports.Likes, ports.Authors),
		HasLiked: usecase.NewHasLikedPostUseCase(ports.Likes),

		CreateSnippet: usecase.NewCreateCodeSnippetUseCase(ports.Snippets, ports.SnippetPosts),

		AddReaction:    usecase.NewAddPostReactionUseCase(ports.Reactions, ports.Authors),
		RemoveReaction: usecase.NewRemovePostReactionUseCase(ports.Reactions, ports.Authors),
		GetReactions:   usecase.NewGetPostReactionsUseCase(ports.Reactions),
		ReactionsBatch: usecase.NewGetPostReactionsBatchUseCase(ports.Reactions),

		CreateComment: usecase.NewCreatePostCommentUseCase(ports.Comments),
		ListComments:  usecase.NewListPostCommentsUseCase(ports.Comments),
		ListReplies:   usecase.NewListCommentRepliesUseCase(ports.Comments),
		EditComment:   usecase.NewEditPostCommentUseCase(ports.Comments),
		DeleteComment: usecase.NewDeletePostCommentUseCase(ports.Comments),
		HideComment:   usecase.NewHidePostCommentUseCase(ports.Comments),
		UnhideComment: usecase.NewUnhidePostCommentUseCase(ports.Comments),

		Bookmark:       usecase.NewBookmarkPostUseCase(ports.Bookmarks, ports.Authors),
		Unbookmark:     usecase.NewUnbookmarkPostUseCase(ports.Bookmarks, ports.Authors),
		HasBookmarked:  usecase.NewHasBookmarkedPostUseCase(ports.Bookmarks),
		ListBookmarks:  usecase.NewListBookmarkedPostsUseCase(ports.Bookmarks),
		CountBookmarks: usecase.NewCountBookmarkedPostsUseCase(ports.Bookmarks),
	})
	return h, ports
}

// setupPostHandlerWithReactionPorts はリアクションの port モックを返すセットアップ。
func setupPostHandlerWithReactionPorts() (*PostHandler, *postHandlerReactionPorts) {
	h, ports := setupPostHandler()
	return h, &postHandlerReactionPorts{Reactions: ports.Reactions, Authors: ports.Authors}
}

// setupPostHandlerWithCommentPort はコメントの port モックを返すセットアップ。
func setupPostHandlerWithCommentPort() (*PostHandler, *mockPostCommentPort) {
	h, ports := setupPostHandler()
	return h, ports.Comments
}

// setupPostHandlerWithBookmarkPorts はブックマークと投稿者判定の port モックを返すセットアップ。
func setupPostHandlerWithBookmarkPorts() (*PostHandler, *mockPostBookmarkPort, *mockPostAuthorPort) {
	h, ports := setupPostHandler()
	return h, ports.Bookmarks, ports.Authors
}

// postHandlerReactionPorts は PostHandler のリアクションに注入した port モックをまとめる。
type postHandlerReactionPorts struct {
	Reactions *mockPostReactionPort
	Authors   *mockPostAuthorPort
}

// setupQuestionHandler はQuestionHandlerテスト用のセットアップを行う。
// 本物の usecase に port モックを注入する。
func setupQuestionHandler() (*QuestionHandler, *mockQuestionRepo) {
	repo := new(mockQuestionRepo)
	h := NewQuestionHandler(
		usecase.NewCreateQuestionUseCase(repo),
		usecase.NewListQuestionsUseCase(repo),
		usecase.NewSearchQuestionsUseCase(repo),
		usecase.NewGetQuestionUseCase(repo),
		usecase.NewListQuestionsByUserUseCase(repo),
		usecase.NewGetQuestionUserVoteUseCase(repo),
		usecase.NewUpdateQuestionUseCase(repo),
		usecase.NewDeleteQuestionUseCase(repo),
		usecase.NewVoteQuestionUseCase(repo),
		usecase.NewRemoveQuestionVoteUseCase(repo),
		usecase.NewListSolvedQuestionsUseCase(repo),
		usecase.NewListUnansweredQuestionsUseCase(repo),
		usecase.NewBookmarkQuestionUseCase(repo),
		usecase.NewUnbookmarkQuestionUseCase(repo),
		usecase.NewListBookmarkedQuestionsUseCase(repo),
		usecase.NewCountQuestionsUseCase(repo),
	)
	return h, repo
}

// setupLearningResourceHandler はLearningResourceHandlerテスト用のセットアップを行う。
// 本物の usecase に port モックを注入する。
func setupLearningResourceHandler() (*LearningResourceHandler, *mockLearningResourceRepo) {
	repo := new(mockLearningResourceRepo)
	h := NewLearningResourceHandler(
		usecase.NewCreateLearningResourceUseCase(repo),
		usecase.NewGetLearningResourceUseCase(repo),
		usecase.NewListLearningResourcesByUserUseCase(repo),
		usecase.NewListPublicLearningResourcesUseCase(repo),
		usecase.NewListLearningResourcesByDifficultyUseCase(repo),
		usecase.NewSearchLearningResourcesUseCase(repo),
		usecase.NewUpdateLearningResourceUseCase(repo),
		usecase.NewUpdateLearningResourceVisibilityUseCase(repo),
		usecase.NewDeleteLearningResourceUseCase(repo),
		usecase.NewLikeLearningResourceUseCase(repo),
		usecase.NewUnlikeLearningResourceUseCase(repo),
		usecase.NewHasLikedLearningResourceUseCase(repo),
		usecase.NewSaveLearningResourceUseCase(repo),
		usecase.NewUnsaveLearningResourceUseCase(repo),
		usecase.NewHasSavedLearningResourceUseCase(repo),
		usecase.NewListSavedLearningResourcesUseCase(repo),
		usecase.NewCountLearningResourcesUseCase(repo),
	)
	return h, repo
}

// MockNotificationSettingsService は NotificationSettingsServiceInterface のモック実装。
// NotificationSettings は DIP へ移行済み。テストは notification_settings_test.go で
// 「本物の usecase + port モック」を組み立てる。

// EmailPreferences は user スライスと一緒に DIP へ移行済み。テストは
// email_preferences_test.go で「本物の usecase + port モック」を組み立てる。

// roadmapHandlerPorts は RoadmapHandler に注入した port モックをまとめる。
type roadmapHandlerPorts struct {
	Roadmaps *mockRoadmapRepo
	Stats    *mockRoadmapStatsRepo
}

// setupRoadmapHandler はRoadmapHandlerテスト用のセットアップを行う。
// 本物の usecase に port モックを注入する。
func setupRoadmapHandler() (*RoadmapHandler, *roadmapHandlerPorts) {
	roadmaps := new(mockRoadmapRepo)
	stats := new(mockRoadmapStatsRepo)
	h := NewRoadmapHandler(
		usecase.NewCreateRoadmapUseCase(roadmaps),
		usecase.NewGetRoadmapUseCase(roadmaps),
		usecase.NewListRoadmapsByUserUseCase(roadmaps),
		usecase.NewListRoadmapsByStatusUseCase(roadmaps),
		usecase.NewListPublicRoadmapsUseCase(roadmaps),
		usecase.NewUpdateRoadmapUseCase(roadmaps),
		usecase.NewUpdateRoadmapVisibilityUseCase(roadmaps),
		usecase.NewDeleteRoadmapUseCase(roadmaps),
		usecase.NewCopyRoadmapUseCase(roadmaps),
		usecase.NewListRoadmapTemplatesUseCase(roadmaps),
		usecase.NewCreateRoadmapFromTemplateUseCase(roadmaps),
		usecase.NewCreateRoadmapStepUseCase(roadmaps),
		usecase.NewUpdateRoadmapStepUseCase(roadmaps),
		usecase.NewUpdateRoadmapStepCompletionUseCase(roadmaps),
		usecase.NewBatchCompleteRoadmapStepsUseCase(roadmaps),
		usecase.NewDeleteRoadmapStepUseCase(roadmaps),
		usecase.NewReorderRoadmapStepsUseCase(roadmaps),
		usecase.NewGetRoadmapStatsUseCase(stats),
		usecase.NewCountRoadmapsUseCase(roadmaps),
	)
	return h, &roadmapHandlerPorts{Roadmaps: roadmaps, Stats: stats}
}

// doRequest はHTTPリクエストを実行してレスポンスを返す。
func doRequest(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, path, jsonBody(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	r.ServeHTTP(w, req)
	return w
}

// doRequestRaw はStringボディでHTTPリクエストを実行する。
func doRequestRaw(r *gin.Engine, method, path, rawBody string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(rawBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

// newRouter はテスト用のGinルーターを生成する。
func newRouter(userID uint) *gin.Engine {
	r := gin.New()
	r.Use(authMiddleware(userID))
	return r
}

// fmtPath はURLパスをフォーマットする。
func fmtPath(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// mockStudyCircleRepo は usecase/repository.StudyCircleRepository のモック（ctx 付き）。
type mockStudyCircleRepo struct{ mock.Mock }

func (m *mockStudyCircleRepo) Create(ctx context.Context, circle *model.StudyCircle) error {
	return m.Called(ctx, circle).Error(0)
}
func (m *mockStudyCircleRepo) FindByID(ctx context.Context, id uint) (*model.StudyCircle, error) {
	args := m.Called(ctx, id)
	if c := args.Get(0); c != nil {
		return c.(*model.StudyCircle), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockStudyCircleRepo) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.StudyCircle, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]model.StudyCircle), args.Get(1).(int64), args.Error(2)
}
func (m *mockStudyCircleRepo) Update(ctx context.Context, circle *model.StudyCircle) error {
	return m.Called(ctx, circle).Error(0)
}
func (m *mockStudyCircleRepo) Delete(ctx context.Context, id uint) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockStudyCircleRepo) GetByStatus(ctx context.Context, userID uint, status string) ([]model.StudyCircle, error) {
	args := m.Called(ctx, userID, status)
	return args.Get(0).([]model.StudyCircle), args.Error(1)
}
func (m *mockStudyCircleRepo) Search(ctx context.Context, query string, limit, offset int) ([]model.StudyCircle, int64, error) {
	args := m.Called(ctx, query, limit, offset)
	if v := args.Get(0); v != nil {
		return v.([]model.StudyCircle), args.Get(1).(int64), args.Error(2)
	}
	return nil, args.Get(1).(int64), args.Error(2)
}
func (m *mockStudyCircleRepo) AddMember(ctx context.Context, circleID, userID uint, role model.StudyCircleMemberRole) error {
	return m.Called(ctx, circleID, userID, role).Error(0)
}
func (m *mockStudyCircleRepo) RemoveMember(ctx context.Context, circleID, userID uint) error {
	return m.Called(ctx, circleID, userID).Error(0)
}
func (m *mockStudyCircleRepo) GetMembers(ctx context.Context, circleID uint) ([]model.StudyCircleMember, error) {
	args := m.Called(ctx, circleID)
	return args.Get(0).([]model.StudyCircleMember), args.Error(1)
}
func (m *mockStudyCircleRepo) IsMember(ctx context.Context, circleID, userID uint) (bool, error) {
	args := m.Called(ctx, circleID, userID)
	return args.Bool(0), args.Error(1)
}
func (m *mockStudyCircleRepo) GetMemberCount(ctx context.Context, circleID uint) (int, error) {
	args := m.Called(ctx, circleID)
	return args.Int(0), args.Error(1)
}
func (m *mockStudyCircleRepo) UpdateMemberRole(ctx context.Context, circleID, userID uint, role model.StudyCircleMemberRole) error {
	return m.Called(ctx, circleID, userID, role).Error(0)
}
func (m *mockStudyCircleRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *mockStudyCircleRepo) CreateStep(ctx context.Context, step *model.StudyCircleStep) error {
	return m.Called(ctx, step).Error(0)
}
func (m *mockStudyCircleRepo) UpdateStep(ctx context.Context, step *model.StudyCircleStep) error {
	return m.Called(ctx, step).Error(0)
}
func (m *mockStudyCircleRepo) DeleteStep(ctx context.Context, stepID uint) error {
	return m.Called(ctx, stepID).Error(0)
}
func (m *mockStudyCircleRepo) FindStepByID(ctx context.Context, stepID uint) (*model.StudyCircleStep, error) {
	args := m.Called(ctx, stepID)
	if s := args.Get(0); s != nil {
		return s.(*model.StudyCircleStep), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockStudyCircleRepo) ReorderSteps(ctx context.Context, circleID uint, stepOrders []model.StepOrder) error {
	return m.Called(ctx, circleID, stepOrders).Error(0)
}
func (m *mockStudyCircleRepo) UpsertProgress(ctx context.Context, progress *model.StudyCircleMemberProgress) error {
	return m.Called(ctx, progress).Error(0)
}
func (m *mockStudyCircleRepo) GetProgress(ctx context.Context, circleID uint) ([]model.StudyCircleMemberProgress, error) {
	args := m.Called(ctx, circleID)
	return args.Get(0).([]model.StudyCircleMemberProgress), args.Error(1)
}
func (m *mockStudyCircleRepo) CreateCheckin(ctx context.Context, checkin *model.StudyCircleCheckin) error {
	return m.Called(ctx, checkin).Error(0)
}
func (m *mockStudyCircleRepo) GetCheckins(ctx context.Context, circleID uint) ([]model.StudyCircleCheckin, error) {
	args := m.Called(ctx, circleID)
	return args.Get(0).([]model.StudyCircleCheckin), args.Error(1)
}
func (m *mockStudyCircleRepo) HasCheckedInToday(ctx context.Context, circleID, userID uint) (bool, error) {
	args := m.Called(ctx, circleID, userID)
	return args.Bool(0), args.Error(1)
}
func (m *mockStudyCircleRepo) GetStreakRanking(ctx context.Context, circleID uint) ([]model.CircleMemberStreak, error) {
	args := m.Called(ctx, circleID)
	return args.Get(0).([]model.CircleMemberStreak), args.Error(1)
}

// setupStudyCircleHandler はStudyCircleHandlerテスト用のセットアップを行う。
// 本物の usecase に port モックを注入する。
func setupStudyCircleHandler() (*StudyCircleHandler, *mockStudyCircleRepo) {
	repo := new(mockStudyCircleRepo)
	h := NewStudyCircleHandler(
		usecase.NewCreateStudyCircleUseCase(repo),
		usecase.NewListMyStudyCirclesUseCase(repo),
		usecase.NewListStudyCirclesByStatusUseCase(repo),
		usecase.NewGetStudyCircleUseCase(repo),
		usecase.NewUpdateStudyCircleUseCase(repo),
		usecase.NewDeleteStudyCircleUseCase(repo),
		usecase.NewListStudyCircleMembersUseCase(repo),
		usecase.NewAddStudyCircleMemberUseCase(repo),
		usecase.NewUpdateStudyCircleMemberRoleUseCase(repo),
		usecase.NewRemoveStudyCircleMemberUseCase(repo),
		usecase.NewCreateStudyCircleStepUseCase(repo),
		usecase.NewUpdateStudyCircleStepUseCase(repo),
		usecase.NewDeleteStudyCircleStepUseCase(repo),
		usecase.NewReorderStudyCircleStepsUseCase(repo),
		usecase.NewUpdateStudyCircleProgressUseCase(repo),
		usecase.NewListStudyCircleProgressUseCase(repo),
		usecase.NewCreateStudyCircleCheckinUseCase(repo),
		usecase.NewListStudyCircleCheckinsUseCase(repo),
		usecase.NewGetStudyCircleStreakRankingUseCase(repo),
		usecase.NewSearchStudyCirclesUseCase(repo),
		usecase.NewCountStudyCirclesUseCase(repo),
	)
	return h, repo
}

// BookReview は DIP へ移行済み。テストは book_review_test.go で
// 「本物の usecase + port モック」を組み立てる。

// mockAnswerRepo は usecase/repository.AnswerRepository のモック（ctx 付き）。
type mockAnswerRepo struct{ mock.Mock }

func (m *mockAnswerRepo) Create(ctx context.Context, answer *model.Answer) error {
	return m.Called(ctx, answer).Error(0)
}
func (m *mockAnswerRepo) FindByID(ctx context.Context, id uint) (*model.Answer, error) {
	args := m.Called(ctx, id)
	if a := args.Get(0); a != nil {
		return a.(*model.Answer), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockAnswerRepo) Update(ctx context.Context, answer *model.Answer) error {
	return m.Called(ctx, answer).Error(0)
}
func (m *mockAnswerRepo) Delete(ctx context.Context, answer *model.Answer) error {
	return m.Called(ctx, answer).Error(0)
}
func (m *mockAnswerRepo) FindByQuestionID(ctx context.Context, questionID uint) ([]model.Answer, error) {
	args := m.Called(ctx, questionID)
	a, _ := args.Get(0).([]model.Answer)
	return a, args.Error(1)
}
func (m *mockAnswerRepo) FindByVoteRange(ctx context.Context, questionID uint, minVote, maxVote int) ([]model.Answer, error) {
	args := m.Called(ctx, questionID, minVote, maxVote)
	a, _ := args.Get(0).([]model.Answer)
	return a, args.Error(1)
}
func (m *mockAnswerRepo) SetBestAnswer(ctx context.Context, questionID, answerID uint) error {
	return m.Called(ctx, questionID, answerID).Error(0)
}
func (m *mockAnswerRepo) Vote(ctx context.Context, userID, answerID uint, value int) error {
	return m.Called(ctx, userID, answerID, value).Error(0)
}
func (m *mockAnswerRepo) RemoveVote(ctx context.Context, userID, answerID uint) error {
	return m.Called(ctx, userID, answerID).Error(0)
}

// answerHandlerPorts は AnswerHandler に注入した port モックをまとめる。
type answerHandlerPorts struct {
	Answers   *mockAnswerRepo
	Questions *mockQuestionRepo
}

// setupAnswerHandler はAnswerHandlerテスト用のセットアップを行う。
// 本物の usecase に port モックを注入する。
func setupAnswerHandler() (*AnswerHandler, *answerHandlerPorts) {
	answers := new(mockAnswerRepo)
	questions := new(mockQuestionRepo)
	h := NewAnswerHandler(
		usecase.NewListAnswersUseCase(answers),
		usecase.NewCreateAnswerUseCase(answers, questions),
		usecase.NewUpdateAnswerUseCase(answers),
		usecase.NewDeleteAnswerUseCase(answers),
		usecase.NewSetBestAnswerUseCase(answers, questions),
		usecase.NewVoteAnswerUseCase(answers),
		usecase.NewRemoveAnswerVoteUseCase(answers),
		usecase.NewListAnswersByVoteRangeUseCase(answers),
	)
	return h, &answerHandlerPorts{Answers: answers, Questions: questions}
}

// mockFollowRepo は usecase/repository.FollowRepository のモック実装（ctx 付き）。
// handler テストは「本物の usecase + port モック」で組む（FreStyle 流）。
type mockFollowRepo struct{ mock.Mock }

func (m *mockFollowRepo) Follow(ctx context.Context, followerID, followeeID uint) error {
	return m.Called(ctx, followerID, followeeID).Error(0)
}
func (m *mockFollowRepo) Unfollow(ctx context.Context, followerID, followeeID uint) error {
	return m.Called(ctx, followerID, followeeID).Error(0)
}
func (m *mockFollowRepo) GetFollowers(ctx context.Context, userID uint, limit, offset int) ([]model.User, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}
func (m *mockFollowRepo) GetFollowing(ctx context.Context, userID uint, limit, offset int) ([]model.User, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]model.User), args.Get(1).(int64), args.Error(2)
}

// setupFollowHandler は FollowHandler テスト用のセットアップを行う。
// 本物の usecase を組み、port(mockFollowRepo)の呼び出しを検証する。
func setupFollowHandler() (*FollowHandler, *mockFollowRepo) {
	repo := new(mockFollowRepo)
	h := NewFollowHandler(
		usecase.NewFollowUserUseCase(repo),
		usecase.NewUnfollowUserUseCase(repo),
		usecase.NewListFollowersUseCase(repo),
		usecase.NewListFollowingUseCase(repo),
	)
	return h, repo
}

// Notification（参照・既読・削除）は DIP へ移行済み。テストは notification_test.go で
// 「本物の usecase + port モック」を組み立てる。

// Project は DIP へ移行済み。テストは project_test.go で
// 「本物の usecase + port モック」を組み立てる。

// Badge は DIP へ移行済み。テストは badge_test.go で
// 「本物の usecase + port モック」を組み立てる。

// RankingHandler のテスト用モックは ranking_test.go（DIP 版・port モック）に置く。

// Level は DIP へ移行済み。テストは level_test.go で
// 「本物の usecase + port モック」を組み立てる。

// ActivityReportHandler のテスト用モックは activity_report_test.go（DIP 版・port モック）に置く。

// MockAIAdviceService は AIAdviceServiceInterface のモック実装。
type MockAIAdviceService struct{ mock.Mock }

func (m *MockAIAdviceService) GenerateAdvice(userID uint) []model.AIAdvice {
	args := m.Called(userID)
	return args.Get(0).([]model.AIAdvice)
}
func (m *MockAIAdviceService) IsLLMAvailable() bool {
	return m.Called().Bool(0)
}
func (m *MockAIAdviceService) GetDailyChatRemaining(userID uint) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}
func (m *MockAIAdviceService) MarkAsRead(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}
func (m *MockAIAdviceService) Chat(userID uint, message string, conversationID uint) (*model.AIConversation, error) {
	args := m.Called(userID, message, conversationID)
	if c := args.Get(0); c != nil {
		return c.(*model.AIConversation), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAIAdviceService) DeleteConversation(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}
func (m *MockAIAdviceService) GetConversations(userID uint, limit, offset int) ([]model.AIConversation, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.AIConversation), args.Error(1)
}
func (m *MockAIAdviceService) GetConversation(id, userID uint) (*model.AIConversation, error) {
	args := m.Called(id, userID)
	if c := args.Get(0); c != nil {
		return c.(*model.AIConversation), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAIAdviceService) GetUnreadAdvice(userID uint) ([]model.AIAdvice, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.AIAdvice), args.Error(1)
}

// setupAIAdviceHandler はAIAdviceHandlerテスト用のセットアップを行う。
func setupAIAdviceHandler() (*AIAdviceHandler, *MockAIAdviceService) {
	svc := new(MockAIAdviceService)
	h := NewAIAdviceHandler(svc)
	return h, svc
}

// LearningAnalytics は DIP へ移行済み。テストは learning_analytics_test.go で
// 「本物の usecase + port モック」を組み立てる。

// ChatRoom は DIP へ移行済み。テストは chat_room_test.go で
// 「本物の usecase + port モック」を組み立てる。

// LearningLog は DIP へ移行済み。テストは learning_log_test.go で
// 「本物の usecase + port モック」を組み立てる。

// Message は DIP へ移行済み。テストは message_test.go で
// 「本物の usecase + port モック」を組み立てる。

// ReminderSettings は DIP へ移行済み。テストは reminder_settings_test.go で
// 「本物の usecase + port モック」を組み立てる。

// Qiita は DIP へ移行済み。テストは qiita_test.go で
// 「本物の usecase + port モック」を組み立てる。

// Zenn は DIP へ移行済み。テストは zenn_test.go で
// 「本物の usecase + port モック」を組み立てる。

// AtCoder は DIP へ移行済み。テストは atcoder_test.go で
// 「本物の usecase + port モック」を組み立てる。

// MockAuthService は AuthServiceInterface のモック実装。
type MockAuthService struct{ mock.Mock }

func (m *MockAuthService) Register(input service.RegisterInput) (*service.AuthResponse, error) {
	args := m.Called(input)
	if r := args.Get(0); r != nil {
		return r.(*service.AuthResponse), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAuthService) Login(input service.LoginInput) (*service.AuthResponse, error) {
	args := m.Called(input)
	if r := args.Get(0); r != nil {
		return r.(*service.AuthResponse), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAuthService) GenerateLoginState() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
func (m *MockAuthService) ValidateLoginState(state string) error {
	return m.Called(state).Error(0)
}
func (m *MockAuthService) GitHubLogin(ghUser *model.GitHubUserInfo, accessToken string) (*service.AuthResponse, error) {
	args := m.Called(ghUser, accessToken)
	if r := args.Get(0); r != nil {
		return r.(*service.AuthResponse), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAuthService) GetMe(userID uint) (*model.User, error) {
	args := m.Called(userID)
	if u := args.Get(0); u != nil {
		return u.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAuthService) RequestPasswordReset(email string) (string, error) {
	args := m.Called(email)
	return args.String(0), args.Error(1)
}
func (m *MockAuthService) ResetPassword(token string, newPassword string) error {
	return m.Called(token, newPassword).Error(0)
}
func (m *MockAuthService) DeleteAccount(userID uint, password string) error {
	return m.Called(userID, password).Error(0)
}

// setupAuthHandlerMock はAuthHandlerテスト用のセットアップを行う。
// GitHub 連携は DIP へ移行済みのため、本物の usecase と port モックを注入する。
func setupAuthHandlerMock() (*AuthHandler, *MockAuthService, *githubPorts) {
	authSvc := new(MockAuthService)
	ports := newGitHubPorts()
	h := NewAuthHandler(authSvc, AuthGitHubUseCases{
		LoginURL:     usecase.NewGetGitHubLoginURLUseCase(ports.Client),
		ExchangeCode: usecase.NewExchangeGitHubCodeUseCase(ports.Client),
		GetUser:      usecase.NewGetGitHubUserUseCase(ports.Client),
		Sync:         usecase.NewSyncGitHubDataUseCase(ports.Users, ports.Repo, ports.Client),
	})
	return h, authSvc, ports
}

// newGitHubPorts は GitHub 連携の port モック一式を生成する。
func newGitHubPorts() *githubPorts {
	return &githubPorts{
		Users:  new(mockUserPort),
		Repo:   new(mockGitHubRepo),
		Client: new(mockGitHubAPIClient),
	}
}

// MockGHAuthService は GitHubAuthServiceInterface のモック実装。
type MockGHAuthService struct{ mock.Mock }

func (m *MockGHAuthService) GenerateOAuthState(userID uint) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}
func (m *MockGHAuthService) ValidateOAuthState(state string) (uint, error) {
	args := m.Called(state)
	return uint(args.Int(0)), args.Error(1)
}

// setupGitHubHandlerMock はGitHubHandlerテスト用のセットアップを行う。
// GitHub 連携は DIP へ移行済みのため、本物の usecase と port モックを注入する。
func setupGitHubHandlerMock() (*GitHubHandler, *githubPorts, *MockGHAuthService) {
	ports := newGitHubPorts()
	authSvc := new(MockGHAuthService)
	sync := usecase.NewSyncGitHubDataUseCase(ports.Users, ports.Repo, ports.Client)
	h := NewGitHubHandler(GitHubUseCases{
		OAuthURL:      usecase.NewGetGitHubOAuthURLUseCase(ports.Client),
		Connect:       usecase.NewConnectGitHubUseCase(ports.Users, ports.Client, sync),
		Disconnect:    usecase.NewDisconnectGitHubUseCase(ports.Users, ports.Repo),
		Sync:          sync,
		Contributions: usecase.NewGetGitHubContributionsUseCase(ports.Repo),
		Languages:     usecase.NewGetGitHubLanguagesUseCase(ports.Repo),
		Repos:         usecase.NewGetGitHubReposUseCase(ports.Repo),
	}, authSvc)
	return h, ports, authSvc
}

// ---------- CodeSnippetHandler モック ----------

// (MockCodeSnippetHandlerService は code_snippet の DIP 移行に伴い撤去。
// テストは code_snippet_test.go で「本物の usecase + port モック」を組み立てる)

// ---------- NoteLinkHandler モック ----------

// NoteLink は DIP へ移行済み。テストは note_link_test.go で
// 「本物の usecase + port モック」を組み立てる。

// ---------- NoteTemplateHandler モック ----------

// NoteTemplate は DIP へ移行済み。テストは note_template_test.go で
// 「本物の usecase + port モック」を組み立てる。

// ---------- CommentLikeHandler モック ----------

// mockCommentLikeRepo は usecase/repository.CommentLikeRepository のモック（ctx 付き）。
type mockCommentLikeRepo struct{ mock.Mock }

func (m *mockCommentLikeRepo) Like(ctx context.Context, userID, commentID uint) error {
	return m.Called(ctx, userID, commentID).Error(0)
}
func (m *mockCommentLikeRepo) Unlike(ctx context.Context, userID, commentID uint) error {
	return m.Called(ctx, userID, commentID).Error(0)
}
func (m *mockCommentLikeRepo) HasLiked(ctx context.Context, userID, commentID uint) (bool, error) {
	args := m.Called(ctx, userID, commentID)
	return args.Bool(0), args.Error(1)
}
func (m *mockCommentLikeRepo) CountByCommentID(ctx context.Context, commentID uint) (int64, error) {
	args := m.Called(ctx, commentID)
	return args.Get(0).(int64), args.Error(1)
}

// mockCommentReader は usecase/repository.CommentReader のモック（ctx 付き）。
type mockCommentReader struct{ mock.Mock }

func (m *mockCommentReader) FindCommentByID(ctx context.Context, id uint) (*model.Comment, error) {
	args := m.Called(ctx, id)
	c, _ := args.Get(0).(*model.Comment)
	return c, args.Error(1)
}

// setupCommentLikeHandler は本物の usecase + port モックで CommentLikeHandler を組む。
func setupCommentLikeHandler() (*CommentLikeHandler, *mockCommentLikeRepo, *mockCommentReader) {
	likes := new(mockCommentLikeRepo)
	reader := new(mockCommentReader)
	h := NewCommentLikeHandler(
		usecase.NewLikeCommentUseCase(likes, reader),
		usecase.NewUnlikeCommentUseCase(likes, reader),
		usecase.NewGetCommentLikeStatusUseCase(likes, reader),
	)
	return h, likes, reader
}

// YouTube は DIP へ移行済み。テストは youtube_test.go で
// 「本物の usecase + port モック」を組み立てる。

// UserDashboardHandler のテスト用モックは user_dashboard_test.go（DIP 版・port モック）に置く。

// ---------- PostTemplateHandler モック ----------

// PostTemplateHandler のテスト用モックは post_template_test.go（DIP 版・port モック）に置く。

// WidgetSettingsHandler のテスト用モックは widget_settings_test.go（DIP 版・port モック）に置く。

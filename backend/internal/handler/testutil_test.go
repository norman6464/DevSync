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
func (m *MockPostRepository) CreateComment(comment *model.Comment) error {
	return m.Called(comment).Error(0)
}
func (m *MockPostRepository) FindCommentByID(id uint) (*model.Comment, error) {
	args := m.Called(id)
	if p := args.Get(0); p != nil {
		return p.(*model.Comment), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockPostRepository) GetComments(postID uint) ([]model.Comment, error) {
	args := m.Called(postID)
	return args.Get(0).([]model.Comment), args.Error(1)
}
func (m *MockPostRepository) UpdateComment(comment *model.Comment) error {
	return m.Called(comment).Error(0)
}
func (m *MockPostRepository) DeleteComment(id uint) error {
	return m.Called(id).Error(0)
}
func (m *MockPostRepository) Bookmark(userID, postID uint) error {
	return m.Called(userID, postID).Error(0)
}
func (m *MockPostRepository) Unbookmark(userID, postID uint) error {
	return m.Called(userID, postID).Error(0)
}
func (m *MockPostRepository) HasBookmarked(userID, postID uint) bool {
	return m.Called(userID, postID).Bool(0)
}
func (m *MockPostRepository) FindBookmarkedByUserID(userID uint, page, limit int) ([]model.Post, int64, error) {
	args := m.Called(userID, page, limit)
	return args.Get(0).([]model.Post), args.Get(1).(int64), args.Error(2)
}
func (m *MockPostRepository) CountBookmarkedByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
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
func (m *MockPostRepository) GetReplies(parentID uint) ([]model.Comment, error) {
	args := m.Called(parentID)
	return args.Get(0).([]model.Comment), args.Error(1)
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

// MockQuestionRepository は QuestionRepositoryInterface のモック実装。
type MockQuestionRepository struct{ mock.Mock }

func (m *MockQuestionRepository) Create(q *model.Question) error {
	return m.Called(q).Error(0)
}
func (m *MockQuestionRepository) FindByID(id uint) (*model.Question, error) {
	args := m.Called(id)
	if q := args.Get(0); q != nil {
		return q.(*model.Question), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockQuestionRepository) FindAll(limit, offset int, tag string, sort string) ([]model.Question, int64, error) {
	args := m.Called(limit, offset, tag, sort)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}
func (m *MockQuestionRepository) Search(q string, limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(q, limit, offset)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}
func (m *MockQuestionRepository) FindByUserID(userID uint, limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}
func (m *MockQuestionRepository) Update(q *model.Question) error {
	return m.Called(q).Error(0)
}
func (m *MockQuestionRepository) Delete(id uint) error {
	return m.Called(id).Error(0)
}
func (m *MockQuestionRepository) Vote(userID, questionID uint, value int) error {
	return m.Called(userID, questionID, value).Error(0)
}
func (m *MockQuestionRepository) RemoveVote(userID, questionID uint) error {
	return m.Called(userID, questionID).Error(0)
}
func (m *MockQuestionRepository) GetUserVote(userID, questionID uint) (int, error) {
	args := m.Called(userID, questionID)
	return args.Int(0), args.Error(1)
}
func (m *MockQuestionRepository) FindSolved(limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}
func (m *MockQuestionRepository) Bookmark(userID, questionID uint) error {
	return m.Called(userID, questionID).Error(0)
}
func (m *MockQuestionRepository) Unbookmark(userID, questionID uint) error {
	return m.Called(userID, questionID).Error(0)
}
func (m *MockQuestionRepository) HasBookmarked(userID, questionID uint) (bool, error) {
	args := m.Called(userID, questionID)
	return args.Bool(0), args.Error(1)
}
func (m *MockQuestionRepository) FindBookmarkedByUserID(userID uint, limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}
func (m *MockQuestionRepository) FindUnanswered(limit, offset int) ([]model.Question, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]model.Question), args.Get(1).(int64), args.Error(2)
}
func (m *MockQuestionRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// MockLearningResourceRepository は LearningResourceRepositoryInterface のモック実装。
type MockLearningResourceRepository struct{ mock.Mock }

func (m *MockLearningResourceRepository) Create(r *model.LearningResource) error {
	return m.Called(r).Error(0)
}
func (m *MockLearningResourceRepository) FindByID(id uint) (*model.LearningResource, error) {
	args := m.Called(id)
	if r := args.Get(0); r != nil {
		return r.(*model.LearningResource), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockLearningResourceRepository) FindByUserID(userID uint, includePrivate bool, limit, offset int) ([]model.LearningResource, int64, error) {
	args := m.Called(userID, includePrivate, limit, offset)
	return args.Get(0).([]model.LearningResource), args.Get(1).(int64), args.Error(2)
}
func (m *MockLearningResourceRepository) FindPublic(limit, offset int, category string, difficulty string) ([]model.LearningResource, int64, error) {
	args := m.Called(limit, offset, category, difficulty)
	return args.Get(0).([]model.LearningResource), args.Get(1).(int64), args.Error(2)
}
func (m *MockLearningResourceRepository) Update(r *model.LearningResource) error {
	return m.Called(r).Error(0)
}
func (m *MockLearningResourceRepository) Delete(id uint) error {
	return m.Called(id).Error(0)
}
func (m *MockLearningResourceRepository) Search(query string, limit, offset int) ([]model.LearningResource, int64, error) {
	args := m.Called(query, limit, offset)
	return args.Get(0).([]model.LearningResource), args.Get(1).(int64), args.Error(2)
}
func (m *MockLearningResourceRepository) Like(userID, resourceID uint) error {
	return m.Called(userID, resourceID).Error(0)
}
func (m *MockLearningResourceRepository) Unlike(userID, resourceID uint) error {
	return m.Called(userID, resourceID).Error(0)
}
func (m *MockLearningResourceRepository) HasLiked(userID, resourceID uint) (bool, error) {
	args := m.Called(userID, resourceID)
	return args.Bool(0), args.Error(1)
}
func (m *MockLearningResourceRepository) Save(userID, resourceID uint) error {
	return m.Called(userID, resourceID).Error(0)
}
func (m *MockLearningResourceRepository) Unsave(userID, resourceID uint) error {
	return m.Called(userID, resourceID).Error(0)
}
func (m *MockLearningResourceRepository) HasSaved(userID, resourceID uint) (bool, error) {
	args := m.Called(userID, resourceID)
	return args.Bool(0), args.Error(1)
}
func (m *MockLearningResourceRepository) FindSavedByUserID(userID uint, limit, offset int) ([]model.LearningResource, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.LearningResource), args.Get(1).(int64), args.Error(2)
}
func (m *MockLearningResourceRepository) FindByDifficulty(difficulty string, limit, offset int) ([]model.LearningResource, int64, error) {
	args := m.Called(difficulty, limit, offset)
	return args.Get(0).([]model.LearningResource), args.Get(1).(int64), args.Error(2)
}
func (m *MockLearningResourceRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// MockRoadmapRepository は RoadmapRepositoryInterface のモック実装。
type MockRoadmapRepository struct{ mock.Mock }

func (m *MockRoadmapRepository) Create(r *model.Roadmap) error {
	return m.Called(r).Error(0)
}
func (m *MockRoadmapRepository) Update(r *model.Roadmap) error {
	return m.Called(r).Error(0)
}
func (m *MockRoadmapRepository) Delete(id uint) error {
	return m.Called(id).Error(0)
}
func (m *MockRoadmapRepository) FindByID(id uint) (*model.Roadmap, error) {
	args := m.Called(id)
	if r := args.Get(0); r != nil {
		return r.(*model.Roadmap), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRoadmapRepository) GetByUserID(userID uint, limit, offset int) ([]model.Roadmap, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.Roadmap), args.Get(1).(int64), args.Error(2)
}
func (m *MockRoadmapRepository) GetPublicRoadmaps(limit, offset int) ([]model.Roadmap, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]model.Roadmap), args.Get(1).(int64), args.Error(2)
}
func (m *MockRoadmapRepository) CopyRoadmap(originalID, newUserID uint) (*model.Roadmap, error) {
	args := m.Called(originalID, newUserID)
	if r := args.Get(0); r != nil {
		return r.(*model.Roadmap), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRoadmapRepository) GetStats(userID uint) (*model.RoadmapStats, error) {
	args := m.Called(userID)
	if s := args.Get(0); s != nil {
		return s.(*model.RoadmapStats), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRoadmapRepository) CreateStep(step *model.RoadmapStep) error {
	return m.Called(step).Error(0)
}
func (m *MockRoadmapRepository) UpdateStep(step *model.RoadmapStep) error {
	return m.Called(step).Error(0)
}
func (m *MockRoadmapRepository) DeleteStep(stepID uint) error {
	return m.Called(stepID).Error(0)
}
func (m *MockRoadmapRepository) FindStepByID(stepID uint) (*model.RoadmapStep, error) {
	args := m.Called(stepID)
	if s := args.Get(0); s != nil {
		return s.(*model.RoadmapStep), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRoadmapRepository) ReorderSteps(roadmapID uint, stepOrders []model.StepOrder) error {
	return m.Called(roadmapID, stepOrders).Error(0)
}
func (m *MockRoadmapRepository) GetTemplates() ([]model.Roadmap, error) {
	args := m.Called()
	return args.Get(0).([]model.Roadmap), args.Error(1)
}
func (m *MockRoadmapRepository) GetByStatus(userID uint, status string) ([]model.Roadmap, error) {
	args := m.Called(userID, status)
	return args.Get(0).([]model.Roadmap), args.Error(1)
}
func (m *MockRoadmapRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// MockChatRoomRepository は ChatRoomRepositoryInterface のモック実装。
type MockChatRoomRepository struct{ mock.Mock }

func (m *MockChatRoomRepository) Create(room *model.ChatRoom) error {
	return m.Called(room).Error(0)
}
func (m *MockChatRoomRepository) FindByID(id uint) (*model.ChatRoom, error) {
	args := m.Called(id)
	if r := args.Get(0); r != nil {
		return r.(*model.ChatRoom), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockChatRoomRepository) FindByUserID(userID uint, limit, offset int) ([]model.ChatRoom, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.ChatRoom), args.Get(1).(int64), args.Error(2)
}
func (m *MockChatRoomRepository) Update(room *model.ChatRoom) error {
	return m.Called(room).Error(0)
}
func (m *MockChatRoomRepository) Delete(roomID uint) error {
	return m.Called(roomID).Error(0)
}
func (m *MockChatRoomRepository) AddMember(roomID, userID uint) error {
	return m.Called(roomID, userID).Error(0)
}
func (m *MockChatRoomRepository) RemoveMember(roomID, userID uint) error {
	return m.Called(roomID, userID).Error(0)
}
func (m *MockChatRoomRepository) GetMembers(roomID uint) ([]model.ChatRoomMember, error) {
	args := m.Called(roomID)
	return args.Get(0).([]model.ChatRoomMember), args.Error(1)
}
func (m *MockChatRoomRepository) IsMember(roomID, userID uint) (bool, error) {
	args := m.Called(roomID, userID)
	return args.Bool(0), args.Error(1)
}
func (m *MockChatRoomRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// MockGroupMessageRepository は GroupMessageRepositoryInterface のモック実装。
type MockGroupMessageRepository struct{ mock.Mock }

func (m *MockGroupMessageRepository) Create(msg *model.GroupMessage) error {
	return m.Called(msg).Error(0)
}
func (m *MockGroupMessageRepository) FindByRoomID(roomID uint, page, limit int) ([]model.GroupMessage, error) {
	args := m.Called(roomID, page, limit)
	return args.Get(0).([]model.GroupMessage), args.Error(1)
}
func (m *MockGroupMessageRepository) FindSenderByID(msg *model.GroupMessage) {
	m.Called(msg)
}
func (m *MockGroupMessageRepository) GetMemberUserIDs(roomID uint) []uint {
	args := m.Called(roomID)
	return args.Get(0).([]uint)
}

// MockRecommendationRepository は RecommendationRepositoryInterface のモック実装。
type MockRecommendationRepository struct{ mock.Mock }

func (m *MockRecommendationRepository) GetRecommendedUsers(userID uint, skills []string, limit int) ([]model.RecommendedUser, error) {
	args := m.Called(userID, skills, limit)
	return args.Get(0).([]model.RecommendedUser), args.Error(1)
}
func (m *MockRecommendationRepository) GetTrendingPosts(limit int, days int) ([]model.Post, error) {
	args := m.Called(limit, days)
	return args.Get(0).([]model.Post), args.Error(1)
}
func (m *MockRecommendationRepository) GetTrendingResources(limit int, days int) ([]model.LearningResource, error) {
	args := m.Called(limit, days)
	return args.Get(0).([]model.LearningResource), args.Error(1)
}

// ---------- リポジトリインターフェース適合チェック ----------
// import cycle を避けるため repository パッケージは使わないが、
// コンパイル時にインターフェース適合を検証したい場合は repository パッケージを
// import する必要がある。ここではハンドラ層テスト内でのみ使用する。


// ---------- ヘルパー関数 ----------

// setupPostHandler はPostHandlerテスト用のセットアップを行う。
// スニペット作成は DIP へ移行済みのため、本物の usecase と port モックを注入する。
func setupPostHandler() (*PostHandler, *MockPostRepository, *MockNotificationRepository, *postHandlerSnippetPorts) {
	postRepo := new(MockPostRepository)
	notifRepo := new(MockNotificationRepository)
	snippets := new(mockCodeSnippetRepo)
	posts := new(mockPostReader)

	notifService := service.NewNotificationService(notifRepo)
	postService := service.NewPostService(postRepo, notifService)
	h := NewPostHandler(postService, usecase.NewCreateCodeSnippetUseCase(snippets, posts))

	return h, postRepo, notifRepo, &postHandlerSnippetPorts{Snippets: snippets, Posts: posts}
}

// postHandlerSnippetPorts は PostHandler のスニペット作成に注入した port モックをまとめる。
type postHandlerSnippetPorts struct {
	Snippets *mockCodeSnippetRepo
	Posts    *mockPostReader
}

// setupQuestionHandler はQuestionHandlerテスト用のセットアップを行う。
func setupQuestionHandler() (*QuestionHandler, *MockQuestionRepository) {
	repo := new(MockQuestionRepository)
	svc := service.NewQuestionService(repo)
	h := NewQuestionHandler(svc)
	return h, repo
}

// setupLearningResourceHandler はLearningResourceHandlerテスト用のセットアップを行う（リポジトリレベル）。
func setupLearningResourceHandler() (*LearningResourceHandler, *MockLearningResourceRepository) {
	repo := new(MockLearningResourceRepository)
	svc := service.NewLearningResourceService(repo)
	h := NewLearningResourceHandler(svc)
	return h, repo
}

// MockNotificationSettingsService は NotificationSettingsServiceInterface のモック実装。
// NotificationSettings は DIP へ移行済み。テストは notification_settings_test.go で
// 「本物の usecase + port モック」を組み立てる。

// MockEmailPreferencesService は EmailPreferencesServiceInterface のモック実装。
type MockEmailPreferencesService struct{ mock.Mock }

func (m *MockEmailPreferencesService) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if u := args.Get(0); u != nil {
		return u.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockEmailPreferencesService) UpdateEmailPreferences(userID uint, weeklyReport *bool, language *string) (*model.User, error) {
	args := m.Called(userID, weeklyReport, language)
	if u := args.Get(0); u != nil {
		return u.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// setupEmailPreferencesHandler はEmailPreferencesHandlerテスト用のセットアップを行う。
func setupEmailPreferencesHandler() (*EmailPreferencesHandler, *MockEmailPreferencesService) {
	svc := new(MockEmailPreferencesService)
	h := NewEmailPreferencesHandler(svc)
	return h, svc
}

// setupRoadmapHandler はRoadmapHandlerテスト用のセットアップを行う。
func setupRoadmapHandler() (*RoadmapHandler, *MockRoadmapRepository) {
	repo := new(MockRoadmapRepository)
	svc := service.NewRoadmapService(repo)
	h := NewRoadmapHandler(svc)
	return h, repo
}

// setupChatRoomHandlerRepo はChatRoomHandlerテスト用のリポジトリレベルセットアップを行う。
func setupChatRoomHandlerRepo() (*ChatRoomHandler, *MockChatRoomRepository, *MockGroupMessageRepository) {
	roomRepo := new(MockChatRoomRepository)
	msgRepo := new(MockGroupMessageRepository)
	hub := service.NewHub()
	svc := service.NewChatRoomService(roomRepo, msgRepo, hub)
	h := NewChatRoomHandler(svc)
	return h, roomRepo, msgRepo
}

// setupRecommendationHandlerRepo はRecommendationHandlerテスト用のセットアップを行う（リポジトリモック版）。
func setupRecommendationHandlerRepo() (*RecommendationHandler, *MockRecommendationRepository, *MockUserRepository) {
	recRepo := new(MockRecommendationRepository)
	userRepo := new(MockUserRepository)
	svc := service.NewRecommendationService(recRepo, userRepo)
	h := NewRecommendationHandler(svc)
	return h, recRepo, userRepo
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

// MockStudyCircleRepository は StudyCircleRepositoryInterface のモック実装。
type MockStudyCircleRepository struct{ mock.Mock }

func (m *MockStudyCircleRepository) Create(circle *model.StudyCircle) error {
	return m.Called(circle).Error(0)
}
func (m *MockStudyCircleRepository) FindByID(id uint) (*model.StudyCircle, error) {
	args := m.Called(id)
	if c := args.Get(0); c != nil {
		return c.(*model.StudyCircle), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockStudyCircleRepository) FindByUserID(userID uint, limit, offset int) ([]model.StudyCircle, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.StudyCircle), args.Get(1).(int64), args.Error(2)
}
func (m *MockStudyCircleRepository) Update(circle *model.StudyCircle) error {
	return m.Called(circle).Error(0)
}
func (m *MockStudyCircleRepository) Delete(id uint) error {
	return m.Called(id).Error(0)
}
func (m *MockStudyCircleRepository) AddMember(circleID, userID uint, role model.StudyCircleMemberRole) error {
	return m.Called(circleID, userID, role).Error(0)
}
func (m *MockStudyCircleRepository) RemoveMember(circleID, userID uint) error {
	return m.Called(circleID, userID).Error(0)
}
func (m *MockStudyCircleRepository) GetMembers(circleID uint) ([]model.StudyCircleMember, error) {
	args := m.Called(circleID)
	return args.Get(0).([]model.StudyCircleMember), args.Error(1)
}
func (m *MockStudyCircleRepository) IsMember(circleID, userID uint) (bool, error) {
	args := m.Called(circleID, userID)
	return args.Bool(0), args.Error(1)
}
func (m *MockStudyCircleRepository) GetMemberCount(circleID uint) (int, error) {
	args := m.Called(circleID)
	return args.Int(0), args.Error(1)
}
func (m *MockStudyCircleRepository) CreateStep(step *model.StudyCircleStep) error {
	return m.Called(step).Error(0)
}
func (m *MockStudyCircleRepository) UpdateStep(step *model.StudyCircleStep) error {
	return m.Called(step).Error(0)
}
func (m *MockStudyCircleRepository) DeleteStep(stepID uint) error {
	return m.Called(stepID).Error(0)
}
func (m *MockStudyCircleRepository) FindStepByID(stepID uint) (*model.StudyCircleStep, error) {
	args := m.Called(stepID)
	if s := args.Get(0); s != nil {
		return s.(*model.StudyCircleStep), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockStudyCircleRepository) ReorderSteps(circleID uint, stepOrders []model.StepOrder) error {
	return m.Called(circleID, stepOrders).Error(0)
}
func (m *MockStudyCircleRepository) UpsertProgress(progress *model.StudyCircleMemberProgress) error {
	return m.Called(progress).Error(0)
}
func (m *MockStudyCircleRepository) GetProgress(circleID uint) ([]model.StudyCircleMemberProgress, error) {
	args := m.Called(circleID)
	return args.Get(0).([]model.StudyCircleMemberProgress), args.Error(1)
}
func (m *MockStudyCircleRepository) CreateCheckin(checkin *model.StudyCircleCheckin) error {
	return m.Called(checkin).Error(0)
}
func (m *MockStudyCircleRepository) GetCheckins(circleID uint) ([]model.StudyCircleCheckin, error) {
	args := m.Called(circleID)
	return args.Get(0).([]model.StudyCircleCheckin), args.Error(1)
}
func (m *MockStudyCircleRepository) HasCheckedInToday(circleID, userID uint) (bool, error) {
	args := m.Called(circleID, userID)
	return args.Bool(0), args.Error(1)
}
func (m *MockStudyCircleRepository) GetStreakRanking(circleID uint) ([]model.CircleMemberStreak, error) {
	args := m.Called(circleID)
	return args.Get(0).([]model.CircleMemberStreak), args.Error(1)
}
func (m *MockStudyCircleRepository) Search(query string, limit, offset int) ([]model.StudyCircle, int64, error) {
	args := m.Called(query, limit, offset)
	if v := args.Get(0); v != nil {
		return v.([]model.StudyCircle), args.Get(1).(int64), args.Error(2)
	}
	return nil, args.Get(1).(int64), args.Error(2)
}
func (m *MockStudyCircleRepository) GetByStatus(userID uint, status string) ([]model.StudyCircle, error) {
	args := m.Called(userID, status)
	return args.Get(0).([]model.StudyCircle), args.Error(1)
}

func (m *MockStudyCircleRepository) UpdateMemberRole(circleID, userID uint, role model.StudyCircleMemberRole) error {
	args := m.Called(circleID, userID, role)
	return args.Error(0)
}

func (m *MockStudyCircleRepository) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// setupStudyCircleHandler はStudyCircleHandlerテスト用のセットアップを行う。
func setupStudyCircleHandler() (*StudyCircleHandler, *MockStudyCircleRepository) {
	repo := new(MockStudyCircleRepository)
	svc := service.NewStudyCircleService(repo)
	h := NewStudyCircleHandler(svc)
	return h, repo
}

// BookReview は DIP へ移行済み。テストは book_review_test.go で
// 「本物の usecase + port モック」を組み立てる。

// MockAnswerService は AnswerServiceInterface のモック実装。
type MockAnswerService struct{ mock.Mock }

func (m *MockAnswerService) GetByQuestionID(questionID uint) ([]model.Answer, error) {
	args := m.Called(questionID)
	return args.Get(0).([]model.Answer), args.Error(1)
}
func (m *MockAnswerService) Create(answer *model.Answer) error {
	return m.Called(answer).Error(0)
}
func (m *MockAnswerService) Update(answerID, userID uint, body string) (*model.Answer, error) {
	args := m.Called(answerID, userID, body)
	if a := args.Get(0); a != nil {
		return a.(*model.Answer), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAnswerService) Delete(answerID, userID uint) error {
	return m.Called(answerID, userID).Error(0)
}
func (m *MockAnswerService) SetBestAnswer(questionID, answerID, userID uint) error {
	return m.Called(questionID, answerID, userID).Error(0)
}
func (m *MockAnswerService) Vote(userID, answerID uint, value int) error {
	return m.Called(userID, answerID, value).Error(0)
}
func (m *MockAnswerService) RemoveVote(userID, answerID uint) error {
	return m.Called(userID, answerID).Error(0)
}
func (m *MockAnswerService) GetByVoteRange(questionID uint, minVote, maxVote int) ([]model.Answer, error) {
	args := m.Called(questionID, minVote, maxVote)
	return args.Get(0).([]model.Answer), args.Error(1)
}

// setupAnswerHandler はAnswerHandlerテスト用のセットアップを行う。
func setupAnswerHandler() (*AnswerHandler, *MockAnswerService) {
	svc := new(MockAnswerService)
	h := NewAnswerHandler(svc)
	return h, svc
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

// MockNotificationService は NotificationServiceInterface のモック実装。
type MockNotificationService struct{ mock.Mock }

func (m *MockNotificationService) GetByUserID(userID uint, page, limit int, notificationType string) ([]model.Notification, error) {
	args := m.Called(userID, page, limit, notificationType)
	return args.Get(0).([]model.Notification), args.Error(1)
}
func (m *MockNotificationService) CountByUserID(userID uint, notificationType string) (int64, error) {
	args := m.Called(userID, notificationType)
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockNotificationService) CountUnread(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}
func (m *MockNotificationService) MarkAsRead(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}
func (m *MockNotificationService) MarkAllAsRead(userID uint) error {
	return m.Called(userID).Error(0)
}
func (m *MockNotificationService) Delete(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}

// setupNotificationHandler はNotificationHandlerテスト用のセットアップを行う。
func setupNotificationHandler() (*NotificationHandler, *MockNotificationService) {
	svc := new(MockNotificationService)
	h := NewNotificationHandler(svc)
	return h, svc
}

// MockProjectService は ProjectServiceInterface のモック実装。
type MockProjectService struct{ mock.Mock }

func (m *MockProjectService) Create(project *model.Project) error {
	return m.Called(project).Error(0)
}
func (m *MockProjectService) GetByID(id, userID uint) (*model.Project, error) {
	args := m.Called(id, userID)
	if p := args.Get(0); p != nil {
		return p.(*model.Project), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockProjectService) GetByUserID(userID uint, limit, offset int) ([]model.Project, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.Project), args.Get(1).(int64), args.Error(2)
}
func (m *MockProjectService) GetFeaturedByUserID(userID uint) ([]model.Project, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Project), args.Error(1)
}
func (m *MockProjectService) GetAll(limit, offset int) ([]model.Project, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]model.Project), args.Get(1).(int64), args.Error(2)
}
func (m *MockProjectService) Update(id, userID uint, updates *model.Project) (*model.Project, error) {
	args := m.Called(id, userID, updates)
	if p := args.Get(0); p != nil {
		return p.(*model.Project), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockProjectService) UpdateFeatured(id, userID uint, featured bool) (*model.Project, error) {
	args := m.Called(id, userID, featured)
	if p := args.Get(0); p != nil {
		return p.(*model.Project), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockProjectService) Delete(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}
func (m *MockProjectService) Archive(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}
func (m *MockProjectService) Unarchive(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}
func (m *MockProjectService) GetArchivedByUserID(userID uint, limit, offset int) ([]model.Project, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.Project), args.Get(1).(int64), args.Error(2)
}
func (m *MockProjectService) Search(query string, limit, offset int) ([]model.Project, int64, error) {
	args := m.Called(query, limit, offset)
	return args.Get(0).([]model.Project), args.Get(1).(int64), args.Error(2)
}
func (m *MockProjectService) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// setupProjectHandler はProjectHandlerテスト用のセットアップを行う。
func setupProjectHandler() (*ProjectHandler, *MockProjectService) {
	svc := new(MockProjectService)
	h := NewProjectHandler(svc)
	return h, svc
}

// MockRoadmapService は RoadmapServiceInterface のモック実装。
type MockRoadmapService struct{ mock.Mock }

func (m *MockRoadmapService) Create(roadmap *model.Roadmap) error {
	return m.Called(roadmap).Error(0)
}
func (m *MockRoadmapService) GetByID(id, userID uint) (*model.Roadmap, error) {
	args := m.Called(id, userID)
	if r := args.Get(0); r != nil {
		return r.(*model.Roadmap), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRoadmapService) GetByUserID(userID uint, limit, offset int) ([]model.Roadmap, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.Roadmap), args.Get(1).(int64), args.Error(2)
}
func (m *MockRoadmapService) GetPublicRoadmaps(limit, offset int) ([]model.Roadmap, int64, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]model.Roadmap), args.Get(1).(int64), args.Error(2)
}
func (m *MockRoadmapService) Update(id, userID uint, updates *model.Roadmap) (*model.Roadmap, error) {
	args := m.Called(id, userID, updates)
	if r := args.Get(0); r != nil {
		return r.(*model.Roadmap), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRoadmapService) UpdateVisibility(id, userID uint, isPublic bool) (*model.Roadmap, error) {
	args := m.Called(id, userID, isPublic)
	if r := args.Get(0); r != nil {
		return r.(*model.Roadmap), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRoadmapService) Delete(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}
func (m *MockRoadmapService) CopyRoadmap(roadmapID, userID uint) (*model.Roadmap, error) {
	args := m.Called(roadmapID, userID)
	if r := args.Get(0); r != nil {
		return r.(*model.Roadmap), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRoadmapService) GetTemplates() ([]model.Roadmap, error) {
	args := m.Called()
	return args.Get(0).([]model.Roadmap), args.Error(1)
}
func (m *MockRoadmapService) CreateFromTemplate(templateID, userID uint) (*model.Roadmap, error) {
	args := m.Called(templateID, userID)
	if r := args.Get(0); r != nil {
		return r.(*model.Roadmap), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRoadmapService) CreateStep(roadmapID, userID uint, step *model.RoadmapStep) error {
	return m.Called(roadmapID, userID, step).Error(0)
}
func (m *MockRoadmapService) UpdateStep(roadmapID, stepID, userID uint, updates *model.RoadmapStep) (*model.RoadmapStep, error) {
	args := m.Called(roadmapID, stepID, userID, updates)
	if s := args.Get(0); s != nil {
		return s.(*model.RoadmapStep), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRoadmapService) UpdateStepCompletion(roadmapID, stepID, userID uint, isCompleted bool) (*model.RoadmapStep, error) {
	args := m.Called(roadmapID, stepID, userID, isCompleted)
	if s := args.Get(0); s != nil {
		return s.(*model.RoadmapStep), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRoadmapService) DeleteStep(roadmapID, stepID, userID uint) error {
	return m.Called(roadmapID, stepID, userID).Error(0)
}
func (m *MockRoadmapService) ReorderSteps(roadmapID, userID uint, orders []model.StepOrder) error {
	return m.Called(roadmapID, userID, orders).Error(0)
}
func (m *MockRoadmapService) GetByStatus(userID uint, status string) ([]model.Roadmap, error) {
	args := m.Called(userID, status)
	return args.Get(0).([]model.Roadmap), args.Error(1)
}
func (m *MockRoadmapService) BatchCompleteSteps(roadmapID, userID uint, stepIDs []uint) (*model.Roadmap, error) {
	args := m.Called(roadmapID, userID, stepIDs)
	if r := args.Get(0); r != nil {
		return r.(*model.Roadmap), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRoadmapService) GetStats(userID uint) (*model.RoadmapStats, error) {
	args := m.Called(userID)
	if s := args.Get(0); s != nil {
		return s.(*model.RoadmapStats), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockRoadmapService) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// setupRoadmapHandlerMock はRoadmapHandlerテスト用のモックセットアップを行う。
func setupRoadmapHandlerMock() (*RoadmapHandler, *MockRoadmapService) {
	svc := new(MockRoadmapService)
	h := NewRoadmapHandler(svc)
	return h, svc
}

// MockBadgeService は BadgeServiceInterface のモック実装。
type MockBadgeService struct{ mock.Mock }

func (m *MockBadgeService) GetUserBadges(userID uint) ([]service.BadgeResult, error) {
	args := m.Called(userID)
	return args.Get(0).([]service.BadgeResult), args.Error(1)
}
func (m *MockBadgeService) NotifyBadgeEarned(userID uint, badgeID string) error {
	return m.Called(userID, badgeID).Error(0)
}

// setupBadgeHandler はBadgeHandlerテスト用のセットアップを行う。
func setupBadgeHandler() (*BadgeHandler, *MockBadgeService) {
	svc := new(MockBadgeService)
	h := NewBadgeHandler(svc)
	return h, svc
}

// RankingHandler のテスト用モックは ranking_test.go（DIP 版・port モック）に置く。

// MockLevelService は LevelServiceInterface のモック実装。
type MockLevelService struct{ mock.Mock }

func (m *MockLevelService) GetLevelInfo(userID uint) (*model.LevelInfo, error) {
	args := m.Called(userID)
	if l := args.Get(0); l != nil {
		return l.(*model.LevelInfo), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockLevelService) GetXPBreakdown(userID uint) (*model.XPBreakdown, error) {
	args := m.Called(userID)
	if x := args.Get(0); x != nil {
		return x.(*model.XPBreakdown), args.Error(1)
	}
	return nil, args.Error(1)
}

// setupLevelHandler はLevelHandlerテスト用のセットアップを行う。
func setupLevelHandler() (*LevelHandler, *MockLevelService) {
	svc := new(MockLevelService)
	h := NewLevelHandler(svc)
	return h, svc
}

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

// MockLearningAnalyticsService は LearningAnalyticsServiceInterface のモック実装。
type MockLearningAnalyticsService struct{ mock.Mock }

func (m *MockLearningAnalyticsService) GetHeatmap(userID uint) ([]model.HeatmapEntry, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.HeatmapEntry), args.Error(1)
}
func (m *MockLearningAnalyticsService) GetCategoryBreakdown(userID uint) ([]model.CategoryBreakdown, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.CategoryBreakdown), args.Error(1)
}
func (m *MockLearningAnalyticsService) GetWeeklyTrends(userID uint, weeks int) ([]model.WeeklyTrend, error) {
	args := m.Called(userID, weeks)
	return args.Get(0).([]model.WeeklyTrend), args.Error(1)
}
func (m *MockLearningAnalyticsService) GetProductivityScore(userID uint) (*model.ProductivityScore, error) {
	args := m.Called(userID)
	if s := args.Get(0); s != nil {
		return s.(*model.ProductivityScore), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockLearningAnalyticsService) GetInsights(userID uint) ([]model.AIInsight, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.AIInsight), args.Error(1)
}
func (m *MockLearningAnalyticsService) GetDayOfWeekSummary(userID uint) ([]model.DayOfWeekSummary, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.DayOfWeekSummary), args.Error(1)
}

// setupLearningAnalyticsHandler はLearningAnalyticsHandlerテスト用のセットアップを行う。
func setupLearningAnalyticsHandler() (*LearningAnalyticsHandler, *MockLearningAnalyticsService) {
	svc := new(MockLearningAnalyticsService)
	h := NewLearningAnalyticsHandler(svc)
	return h, svc
}

// MockChatRoomService は ChatRoomServiceInterface のモック実装。
type MockChatRoomService struct{ mock.Mock }

func (m *MockChatRoomService) Create(room *model.ChatRoom, memberIDs []uint) (*model.ChatRoom, error) {
	args := m.Called(room, memberIDs)
	if r := args.Get(0); r != nil {
		return r.(*model.ChatRoom), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockChatRoomService) GetByUserID(userID uint, limit, offset int) ([]model.ChatRoom, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.ChatRoom), args.Get(1).(int64), args.Error(2)
}
func (m *MockChatRoomService) GetByID(roomID, userID uint) (*model.ChatRoom, error) {
	args := m.Called(roomID, userID)
	if r := args.Get(0); r != nil {
		return r.(*model.ChatRoom), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockChatRoomService) Update(roomID, userID uint, name, description string) (*model.ChatRoom, error) {
	args := m.Called(roomID, userID, name, description)
	if r := args.Get(0); r != nil {
		return r.(*model.ChatRoom), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockChatRoomService) Delete(roomID, userID uint) error {
	return m.Called(roomID, userID).Error(0)
}
func (m *MockChatRoomService) GetMembers(roomID, userID uint) ([]model.ChatRoomMember, error) {
	args := m.Called(roomID, userID)
	return args.Get(0).([]model.ChatRoomMember), args.Error(1)
}
func (m *MockChatRoomService) AddMember(roomID, userID, targetUserID uint) error {
	return m.Called(roomID, userID, targetUserID).Error(0)
}
func (m *MockChatRoomService) RemoveMember(roomID, userID, targetUserID uint) error {
	return m.Called(roomID, userID, targetUserID).Error(0)
}
func (m *MockChatRoomService) GetMessages(roomID, userID uint, page, limit int) ([]model.GroupMessage, error) {
	args := m.Called(roomID, userID, page, limit)
	return args.Get(0).([]model.GroupMessage), args.Error(1)
}
func (m *MockChatRoomService) SendMessage(roomID, userID uint, content string) (*model.GroupMessage, error) {
	args := m.Called(roomID, userID, content)
	if r := args.Get(0); r != nil {
		return r.(*model.GroupMessage), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockChatRoomService) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// setupChatRoomHandler はChatRoomHandlerテスト用のセットアップを行う。
func setupChatRoomHandler() (*ChatRoomHandler, *MockChatRoomService) {
	svc := new(MockChatRoomService)
	h := NewChatRoomHandler(svc)
	return h, svc
}

// MockLearningLogService は LearningLogServiceInterface のモック実装。
type MockLearningLogService struct{ mock.Mock }

func (m *MockLearningLogService) Create(log *model.LearningLog) error {
	return m.Called(log).Error(0)
}
func (m *MockLearningLogService) BatchCreate(userID uint, logs []model.LearningLog) ([]model.LearningLog, error) {
	args := m.Called(userID, logs)
	if l := args.Get(0); l != nil {
		return l.([]model.LearningLog), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockLearningLogService) ImportCSV(userID uint, data []byte) ([]model.LearningLog, error) {
	args := m.Called(userID, data)
	if l := args.Get(0); l != nil {
		return l.([]model.LearningLog), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockLearningLogService) GetByID(id, userID uint) (*model.LearningLog, error) {
	args := m.Called(id, userID)
	if l := args.Get(0); l != nil {
		return l.(*model.LearningLog), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockLearningLogService) GetByUserID(userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.LearningLog), args.Get(1).(int64), args.Error(2)
}
func (m *MockLearningLogService) Update(id, userID uint, updates *model.LearningLog) (*model.LearningLog, error) {
	args := m.Called(id, userID, updates)
	if l := args.Get(0); l != nil {
		return l.(*model.LearningLog), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockLearningLogService) Delete(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}
func (m *MockLearningLogService) GetStreakInfo(userID uint) (*model.StreakInfo, error) {
	args := m.Called(userID)
	if s := args.Get(0); s != nil {
		return s.(*model.StreakInfo), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockLearningLogService) GetCalendarData(userID uint) ([]model.CalendarEntry, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.CalendarEntry), args.Error(1)
}

func (m *MockLearningLogService) ExportCSV(userID uint, days int) ([]byte, error) {
	args := m.Called(userID, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}
func (m *MockLearningLogService) ExportJSON(userID uint, days int) ([]byte, error) {
	args := m.Called(userID, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}
func (m *MockLearningLogService) GetByCategory(userID uint, category string) ([]model.LearningLog, error) {
	args := m.Called(userID, category)
	return args.Get(0).([]model.LearningLog), args.Error(1)
}
func (m *MockLearningLogService) GetBySource(userID uint, source string) ([]model.LearningLog, error) {
	args := m.Called(userID, source)
	return args.Get(0).([]model.LearningLog), args.Error(1)
}
func (m *MockLearningLogService) GetWeeklyDuration(userID uint) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

func (m *MockLearningLogService) FavoriteLog(id, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockLearningLogService) UnfavoriteLog(id, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockLearningLogService) GetRecentCategories(userID uint) ([]string, error) {
	args := m.Called(userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockLearningLogService) GetLinkedLogs(goalID, userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	args := m.Called(goalID, userID, limit, offset)
	return args.Get(0).([]model.LearningLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockLearningLogService) GetFavorites(userID uint, limit, offset int) ([]model.LearningLog, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]model.LearningLog), args.Get(1).(int64), args.Error(2)
}
func (m *MockLearningLogService) GetMonthlySummary(userID uint, months int) ([]model.MonthlySummary, error) {
	args := m.Called(userID, months)
	return args.Get(0).([]model.MonthlySummary), args.Error(1)
}

func (m *MockLearningLogService) GetGoalProgress(goalID, userID uint) (*model.GoalProgress, error) {
	args := m.Called(goalID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.GoalProgress), args.Error(1)
}

func (m *MockLearningLogService) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

// setupLearningLogHandler はLearningLogHandlerテスト用のセットアップを行う。
func setupLearningLogHandler() (*LearningLogHandler, *MockLearningLogService) {
	svc := new(MockLearningLogService)
	h := NewLearningLogHandler(svc)
	return h, svc
}

// MockMessageService は MessageServiceInterface のモック実装。
type MockMessageService struct{ mock.Mock }

func (m *MockMessageService) GetConversations(userID uint) ([]model.ConversationSummary, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.ConversationSummary), args.Error(1)
}
func (m *MockMessageService) GetConversation(userID, otherUserID uint, page, limit int) ([]model.Message, error) {
	args := m.Called(userID, otherUserID, page, limit)
	return args.Get(0).([]model.Message), args.Error(1)
}
func (m *MockMessageService) SendMessage(msg *model.Message) error {
	return m.Called(msg).Error(0)
}
func (m *MockMessageService) MarkAsRead(senderID, receiverID uint) error {
	return m.Called(senderID, receiverID).Error(0)
}

// setupMessageHandler はMessageHandlerテスト用のセットアップを行う。
func setupMessageHandler() (*MessageHandler, *MockMessageService) {
	svc := new(MockMessageService)
	h := NewMessageHandler(svc)
	return h, svc
}

// ReminderSettings は DIP へ移行済み。テストは reminder_settings_test.go で
// 「本物の usecase + port モック」を組み立てる。

// MockRecommendationService は RecommendationServiceInterface のモック実装。
type MockRecommendationService struct{ mock.Mock }

func (m *MockRecommendationService) GetRecommendedUsers(userID uint) ([]model.RecommendedUser, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.RecommendedUser), args.Error(1)
}
func (m *MockRecommendationService) GetTrendingPosts() ([]model.Post, error) {
	args := m.Called()
	return args.Get(0).([]model.Post), args.Error(1)
}
func (m *MockRecommendationService) GetTrendingResources() ([]model.LearningResource, error) {
	args := m.Called()
	return args.Get(0).([]model.LearningResource), args.Error(1)
}

// setupRecommendationHandler はRecommendationHandlerテスト用のセットアップを行う。
func setupRecommendationHandler() (*RecommendationHandler, *MockRecommendationService) {
	svc := new(MockRecommendationService)
	h := NewRecommendationHandler(svc)
	return h, svc
}

// MockQiitaService は QiitaServiceInterface のモック実装。
type MockQiitaService struct{ mock.Mock }

func (m *MockQiitaService) Connect(userID uint, username string) (int, error) {
	args := m.Called(userID, username)
	return args.Int(0), args.Error(1)
}
func (m *MockQiitaService) Disconnect(userID uint) error {
	return m.Called(userID).Error(0)
}
func (m *MockQiitaService) Sync(userID uint) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}
func (m *MockQiitaService) GetArticles(userID uint) ([]model.QiitaArticle, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.QiitaArticle), args.Error(1)
}
func (m *MockQiitaService) GetStats(userID uint) (*model.QiitaStats, error) {
	args := m.Called(userID)
	if s := args.Get(0); s != nil {
		return s.(*model.QiitaStats), args.Error(1)
	}
	return nil, args.Error(1)
}

// setupQiitaHandler はQiitaHandlerテスト用のセットアップを行う。
func setupQiitaHandler() (*ArticlePlatformHandler[model.QiitaArticle, model.QiitaStats], *MockQiitaService) {
	svc := new(MockQiitaService)
	h := NewArticlePlatformHandler[model.QiitaArticle, model.QiitaStats](svc, "Qiita")
	return h, svc
}

// MockZennService は ZennServiceInterface のモック実装。
type MockZennService struct{ mock.Mock }

func (m *MockZennService) Connect(userID uint, username string) (int, error) {
	args := m.Called(userID, username)
	return args.Int(0), args.Error(1)
}
func (m *MockZennService) Disconnect(userID uint) error {
	return m.Called(userID).Error(0)
}
func (m *MockZennService) Sync(userID uint) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}
func (m *MockZennService) GetArticles(userID uint) ([]model.ZennArticle, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.ZennArticle), args.Error(1)
}
func (m *MockZennService) GetStats(userID uint) (*model.ZennStats, error) {
	args := m.Called(userID)
	if s := args.Get(0); s != nil {
		return s.(*model.ZennStats), args.Error(1)
	}
	return nil, args.Error(1)
}

// setupZennHandler はZennHandlerテスト用のセットアップを行う。
func setupZennHandler() (*ArticlePlatformHandler[model.ZennArticle, model.ZennStats], *MockZennService) {
	svc := new(MockZennService)
	h := NewArticlePlatformHandler[model.ZennArticle, model.ZennStats](svc, "Zenn")
	return h, svc
}

// MockAtCoderService は AtCoderServiceInterface のモック実装。
type MockAtCoderService struct{ mock.Mock }

func (m *MockAtCoderService) GetRating(username string) (*service.AtCoderRatingInfo, error) {
	args := m.Called(username)
	if r := args.Get(0); r != nil {
		return r.(*service.AtCoderRatingInfo), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAtCoderService) ConnectAtCoder(userID uint, username string) (*model.User, error) {
	args := m.Called(userID, username)
	if u := args.Get(0); u != nil {
		return u.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAtCoderService) DisconnectAtCoder(userID uint) (*model.User, error) {
	args := m.Called(userID)
	if u := args.Get(0); u != nil {
		return u.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// setupAtCoderHandler はAtCoderHandlerテスト用のセットアップを行う。
func setupAtCoderHandler() (*AtCoderHandler, *MockAtCoderService) {
	atcoderSvc := new(MockAtCoderService)
	h := NewAtCoderHandler(atcoderSvc)
	return h, atcoderSvc
}

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
func (m *MockAuthService) GitHubLogin(ghUser *service.GitHubUserInfo, accessToken string) (*service.AuthResponse, error) {
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

// MockAuthGitHubService は AuthGitHubServiceInterface のモック実装。
type MockAuthGitHubService struct{ mock.Mock }

func (m *MockAuthGitHubService) GetLoginOAuthURL(state string) string {
	return m.Called(state).String(0)
}
func (m *MockAuthGitHubService) ExchangeCode(code string) (string, error) {
	args := m.Called(code)
	return args.String(0), args.Error(1)
}
func (m *MockAuthGitHubService) GetGitHubUser(token string) (*service.GitHubUserInfo, error) {
	args := m.Called(token)
	if u := args.Get(0); u != nil {
		return u.(*service.GitHubUserInfo), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockAuthGitHubService) SyncUserData(userID uint) error {
	return m.Called(userID).Error(0)
}

// setupAuthHandler はAuthHandlerテスト用のセットアップを行う。
func setupAuthHandlerMock() (*AuthHandler, *MockAuthService, *MockAuthGitHubService) {
	authSvc := new(MockAuthService)
	ghSvc := new(MockAuthGitHubService)
	h := NewAuthHandler(authSvc, ghSvc)
	return h, authSvc, ghSvc
}

// MockGHService は GitHubServiceInterface のモック実装。
type MockGHService struct{ mock.Mock }

func (m *MockGHService) GetOAuthURL(state string) string {
	return m.Called(state).String(0)
}
func (m *MockGHService) ConnectGitHub(userID uint, code, state string) error {
	return m.Called(userID, code, state).Error(0)
}
func (m *MockGHService) GetContributions(userID uint) ([]model.GitHubContribution, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.GitHubContribution), args.Error(1)
}
func (m *MockGHService) GetLanguages(userID uint) ([]model.GitHubLanguageStat, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.GitHubLanguageStat), args.Error(1)
}
func (m *MockGHService) GetRepos(userID uint) ([]model.GitHubRepository, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.GitHubRepository), args.Error(1)
}
func (m *MockGHService) SyncUserData(userID uint) error {
	return m.Called(userID).Error(0)
}
func (m *MockGHService) DisconnectGitHub(userID uint) error {
	return m.Called(userID).Error(0)
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
func setupGitHubHandlerMock() (*GitHubHandler, *MockGHService, *MockGHAuthService) {
	ghSvc := new(MockGHService)
	authSvc := new(MockGHAuthService)
	h := NewGitHubHandler(ghSvc, authSvc)
	return h, ghSvc, authSvc
}

// ---------- CodeSnippetHandler モック ----------

// (MockCodeSnippetHandlerService は code_snippet の DIP 移行に伴い撤去。
// テストは code_snippet_test.go で「本物の usecase + port モック」を組み立てる)

// ---------- NoteLinkHandler モック ----------

// NoteLink は DIP へ移行済み。テストは note_link_test.go で
// 「本物の usecase + port モック」を組み立てる。

// ---------- NoteTemplateHandler モック ----------

// MockNoteTemplateService は NoteTemplateServiceInterface のモック実装。
type MockNoteTemplateService struct{ mock.Mock }

func (m *MockNoteTemplateService) Create(template *model.NoteTemplate) error {
	return m.Called(template).Error(0)
}
func (m *MockNoteTemplateService) GetByID(id, userID uint) (*model.NoteTemplate, error) {
	args := m.Called(id, userID)
	if t := args.Get(0); t != nil {
		return t.(*model.NoteTemplate), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockNoteTemplateService) GetByUserID(userID uint) ([]model.NoteTemplate, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.NoteTemplate), args.Error(1)
}
func (m *MockNoteTemplateService) GetDefaultByUserID(userID uint) (*model.NoteTemplate, error) {
	args := m.Called(userID)
	if t := args.Get(0); t != nil {
		return t.(*model.NoteTemplate), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockNoteTemplateService) Update(id, userID uint, name, description, defaultTitle, contentTemplate, defaultTags string, isDefault *bool) (*model.NoteTemplate, error) {
	args := m.Called(id, userID, name, description, defaultTitle, contentTemplate, defaultTags, isDefault)
	if t := args.Get(0); t != nil {
		return t.(*model.NoteTemplate), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockNoteTemplateService) Delete(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}
func (m *MockNoteTemplateService) UseTemplate(id, userID uint) (*model.Note, error) {
	args := m.Called(id, userID)
	if n := args.Get(0); n != nil {
		return n.(*model.Note), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockNoteTemplateService) CountByUserID(userID uint) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func setupNoteTemplateHandler() (*NoteTemplateHandler, *MockNoteTemplateService) {
	svc := new(MockNoteTemplateService)
	h := NewNoteTemplateHandler(svc)
	return h, svc
}

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

// ---------- YouTubeHandler モック ----------

// MockYouTubeService は YouTubeServiceInterface のモック実装。
type MockYouTubeService struct{ mock.Mock }

func (m *MockYouTubeService) Search(query, language string) ([]model.YouTubeVideo, bool, error) {
	args := m.Called(query, language)
	if v := args.Get(0); v != nil {
		return v.([]model.YouTubeVideo), args.Bool(1), args.Error(2)
	}
	return nil, args.Bool(1), args.Error(2)
}
func (m *MockYouTubeService) GetRecommendations(userID uint) ([]model.YouTubeVideo, []string, error) {
	args := m.Called(userID)
	if v := args.Get(0); v != nil {
		return v.([]model.YouTubeVideo), args.Get(1).([]string), args.Error(2)
	}
	return nil, nil, args.Error(2)
}
func (m *MockYouTubeService) IsAvailable() bool {
	return m.Called().Bool(0)
}

func setupYouTubeHandler() (*YouTubeHandler, *MockYouTubeService) {
	svc := new(MockYouTubeService)
	h := NewYouTubeHandler(svc)
	return h, svc
}

// UserDashboardHandler のテスト用モックは user_dashboard_test.go（DIP 版・port モック）に置く。

// ---------- PostTemplateHandler モック ----------

// PostTemplateHandler のテスト用モックは post_template_test.go（DIP 版・port モック）に置く。


// WidgetSettingsHandler のテスト用モックは widget_settings_test.go（DIP 版・port モック）に置く。


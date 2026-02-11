package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
	"github.com/norman6464/devsync/backend/internal/service"
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
func (m *MockPostRepository) FindByUserID(userID uint) ([]model.Post, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Post), args.Error(1)
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
func (m *MockPostRepository) GetComments(postID uint) ([]model.Comment, error) {
	args := m.Called(postID)
	return args.Get(0).([]model.Comment), args.Error(1)
}
func (m *MockPostRepository) DeleteComment(id, userID uint) error {
	return m.Called(id, userID).Error(0)
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

// MockCodeSnippetRepository は CodeSnippetRepositoryInterface のモック実装。
type MockCodeSnippetRepository struct{ mock.Mock }

func (m *MockCodeSnippetRepository) Create(s *model.CodeSnippet) error {
	return m.Called(s).Error(0)
}
func (m *MockCodeSnippetRepository) FindByID(id uint) (*model.CodeSnippet, error) {
	args := m.Called(id)
	if s := args.Get(0); s != nil {
		return s.(*model.CodeSnippet), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *MockCodeSnippetRepository) FindByPostID(postID uint) ([]model.CodeSnippet, error) {
	args := m.Called(postID)
	return args.Get(0).([]model.CodeSnippet), args.Error(1)
}
func (m *MockCodeSnippetRepository) Update(s *model.CodeSnippet) error {
	return m.Called(s).Error(0)
}
func (m *MockCodeSnippetRepository) Delete(id uint) error {
	return m.Called(id).Error(0)
}
func (m *MockCodeSnippetRepository) CreateComment(c *model.SnippetComment) error {
	return m.Called(c).Error(0)
}
func (m *MockCodeSnippetRepository) GetComments(snippetID uint) ([]model.SnippetComment, error) {
	args := m.Called(snippetID)
	return args.Get(0).([]model.SnippetComment), args.Error(1)
}
func (m *MockCodeSnippetRepository) DeleteComment(id, userID uint) error {
	return m.Called(id, userID).Error(0)
}

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
func (m *MockQuestionRepository) FindByUserID(userID uint) ([]model.Question, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Question), args.Error(1)
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
func (m *MockLearningResourceRepository) FindByUserID(userID uint, includePrivate bool) ([]model.LearningResource, error) {
	args := m.Called(userID, includePrivate)
	return args.Get(0).([]model.LearningResource), args.Error(1)
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
func (m *MockRoadmapRepository) GetByUserID(userID uint) ([]model.Roadmap, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Roadmap), args.Error(1)
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
func (m *MockRoadmapRepository) ReorderSteps(roadmapID uint, stepOrders []repository.StepOrder) error {
	return m.Called(roadmapID, stepOrders).Error(0)
}
func (m *MockRoadmapRepository) GetTemplates() ([]model.Roadmap, error) {
	args := m.Called()
	return args.Get(0).([]model.Roadmap), args.Error(1)
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
func (m *MockChatRoomRepository) FindByUserID(userID uint) ([]model.ChatRoom, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.ChatRoom), args.Error(1)
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

// ---------- リポジトリインターフェース適合チェック ----------
// import cycle を避けるため repository パッケージは使わないが、
// コンパイル時にインターフェース適合を検証したい場合は repository パッケージを
// import する必要がある。ここではハンドラ層テスト内でのみ使用する。


// ---------- ヘルパー関数 ----------

// setupPostHandler はPostHandlerテスト用のセットアップを行う。
func setupPostHandler() (*PostHandler, *MockPostRepository, *MockNotificationRepository, *MockCodeSnippetRepository) {
	postRepo := new(MockPostRepository)
	notifRepo := new(MockNotificationRepository)
	snippetRepo := new(MockCodeSnippetRepository)

	notifService := service.NewNotificationService(notifRepo)
	postService := service.NewPostService(postRepo, notifService)
	snippetService := service.NewCodeSnippetService(snippetRepo, postRepo)
	h := NewPostHandler(postService, snippetService)

	return h, postRepo, notifRepo, snippetRepo
}

// setupQuestionHandler はQuestionHandlerテスト用のセットアップを行う。
func setupQuestionHandler() (*QuestionHandler, *MockQuestionRepository) {
	repo := new(MockQuestionRepository)
	svc := service.NewQuestionService(repo)
	h := NewQuestionHandler(svc)
	return h, repo
}

// setupLearningResourceHandler はLearningResourceHandlerテスト用のセットアップを行う。
func setupLearningResourceHandler() (*LearningResourceHandler, *MockLearningResourceRepository) {
	repo := new(MockLearningResourceRepository)
	svc := service.NewLearningResourceService(repo)
	h := NewLearningResourceHandler(svc)
	return h, repo
}

// setupRoadmapHandler はRoadmapHandlerテスト用のセットアップを行う。
func setupRoadmapHandler() (*RoadmapHandler, *MockRoadmapRepository) {
	repo := new(MockRoadmapRepository)
	svc := service.NewRoadmapService(repo)
	h := NewRoadmapHandler(svc)
	return h, repo
}

// setupChatRoomHandler はChatRoomHandlerテスト用のセットアップを行う。
func setupChatRoomHandler() (*ChatRoomHandler, *MockChatRoomRepository, *MockGroupMessageRepository) {
	roomRepo := new(MockChatRoomRepository)
	msgRepo := new(MockGroupMessageRepository)
	hub := service.NewHub()
	svc := service.NewChatRoomService(roomRepo, msgRepo, hub)
	h := NewChatRoomHandler(svc)
	return h, roomRepo, msgRepo
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

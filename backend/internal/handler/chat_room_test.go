package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------- port モック ----------

// mockChatRoomPort は usecase/repository.ChatRoomRepository のモック（ctx 付き）。
type mockChatRoomPort struct{ mock.Mock }

func (m *mockChatRoomPort) Create(ctx context.Context, room *model.ChatRoom) error {
	return m.Called(ctx, room).Error(0)
}
func (m *mockChatRoomPort) FindByID(ctx context.Context, id uint) (*model.ChatRoom, error) {
	args := m.Called(ctx, id)
	if r := args.Get(0); r != nil {
		return r.(*model.ChatRoom), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockChatRoomPort) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.ChatRoom, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]model.ChatRoom), args.Get(1).(int64), args.Error(2)
}
func (m *mockChatRoomPort) Update(ctx context.Context, room *model.ChatRoom) error {
	return m.Called(ctx, room).Error(0)
}
func (m *mockChatRoomPort) Delete(ctx context.Context, roomID uint) error {
	return m.Called(ctx, roomID).Error(0)
}
func (m *mockChatRoomPort) AddMember(ctx context.Context, roomID, userID uint) error {
	return m.Called(ctx, roomID, userID).Error(0)
}
func (m *mockChatRoomPort) RemoveMember(ctx context.Context, roomID, userID uint) error {
	return m.Called(ctx, roomID, userID).Error(0)
}
func (m *mockChatRoomPort) GetMembers(ctx context.Context, roomID uint) ([]model.ChatRoomMember, error) {
	args := m.Called(ctx, roomID)
	return args.Get(0).([]model.ChatRoomMember), args.Error(1)
}
func (m *mockChatRoomPort) IsMember(ctx context.Context, roomID, userID uint) (bool, error) {
	args := m.Called(ctx, roomID, userID)
	return args.Bool(0), args.Error(1)
}
func (m *mockChatRoomPort) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// mockChatRoomMessagePort は usecase/repository.ChatRoomMessageRepository のモック。
type mockChatRoomMessagePort struct{ mock.Mock }

func (m *mockChatRoomMessagePort) Create(ctx context.Context, msg *model.GroupMessage) error {
	return m.Called(ctx, msg).Error(0)
}
func (m *mockChatRoomMessagePort) FindByRoomID(ctx context.Context, roomID uint, page, limit int) ([]model.GroupMessage, error) {
	args := m.Called(ctx, roomID, page, limit)
	return args.Get(0).([]model.GroupMessage), args.Error(1)
}
func (m *mockChatRoomMessagePort) FindSender(ctx context.Context, senderID uint) (*model.User, error) {
	args := m.Called(ctx, senderID)
	if u := args.Get(0); u != nil {
		return u.(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// mockRoomBroadcaster は usecase/repository.RoomBroadcaster のモック。
// 配信は goroutine 内で行われるため、受け取った内容をチャネルに積んでテスト側から待ち受ける。
type mockRoomBroadcaster struct {
	sent chan broadcastCall
}

// broadcastCall は SendToRoom に渡された引数。
type broadcastCall struct {
	roomID   uint
	senderID uint
	message  []byte
}

func newMockRoomBroadcaster() *mockRoomBroadcaster {
	return &mockRoomBroadcaster{sent: make(chan broadcastCall, 4)}
}

func (m *mockRoomBroadcaster) SendToRoom(roomID, senderID uint, message []byte) {
	m.sent <- broadcastCall{roomID: roomID, senderID: senderID, message: message}
}

// wait は配信されるまで待ち、配信が無ければ失敗させる。
func (m *mockRoomBroadcaster) wait(t *testing.T) broadcastCall {
	t.Helper()
	select {
	case call := <-m.sent:
		return call
	case <-time.After(2 * time.Second):
		t.Fatal("WebSocket への配信が行われなかった")
		return broadcastCall{}
	}
}

// newTestChatRoomHandler は本物の usecase に port モックを注入したハンドラーを生成する。
func newTestChatRoomHandler() (*ChatRoomHandler, *mockChatRoomPort, *mockChatRoomMessagePort, *mockRoomBroadcaster) {
	rooms := new(mockChatRoomPort)
	messages := new(mockChatRoomMessagePort)
	broadcaster := newMockRoomBroadcaster()
	h := NewChatRoomHandler(
		usecase.NewCreateChatRoomUseCase(rooms),
		usecase.NewListMyChatRoomsUseCase(rooms),
		usecase.NewGetChatRoomUseCase(rooms),
		usecase.NewUpdateChatRoomUseCase(rooms),
		usecase.NewDeleteChatRoomUseCase(rooms),
		usecase.NewListChatRoomMembersUseCase(rooms),
		usecase.NewAddChatRoomMemberUseCase(rooms),
		usecase.NewRemoveChatRoomMemberUseCase(rooms),
		usecase.NewListChatRoomMessagesUseCase(rooms, messages),
		usecase.NewSendChatRoomMessageUseCase(rooms, messages, broadcaster),
		usecase.NewCountMyChatRoomsUseCase(rooms),
	)
	return h, rooms, messages, broadcaster
}

// ---------- Create ----------

func TestChatRoomCreate_Success(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms", h.Create)

	rooms.On("Create", mock.Anything, mock.MatchedBy(func(room *model.ChatRoom) bool {
		// 前後の空白は落として保存する。
		return room.Name == "Test Room" && room.Description == "説明" && room.OwnerID == 1
	})).Run(func(args mock.Arguments) {
		args.Get(1).(*model.ChatRoom).ID = 10
	}).Return(nil)
	// オーナーと、オーナー以外のメンバーだけが追加される。
	rooms.On("AddMember", mock.Anything, uint(10), uint(1)).Return(nil)
	rooms.On("AddMember", mock.Anything, uint(10), uint(2)).Return(nil)
	rooms.On("AddMember", mock.Anything, uint(10), uint(3)).Return(nil)
	rooms.On("FindByID", mock.Anything, uint(10)).
		Return(&model.ChatRoom{Name: "Test Room", OwnerID: 1}, nil)

	w := doRequest(r, http.MethodPost, "/rooms", map[string]interface{}{
		"name": "  Test Room  ", "description": " 説明 ", "member_ids": []uint{1, 2, 3},
	})
	assertStatus(t, w, http.StatusCreated)
	// オーナーは member_ids に含まれていても 1 回しか追加しない。
	rooms.AssertNumberOfCalls(t, "AddMember", 3)
	rooms.AssertExpectations(t)
}

// 個々のメンバー追加の失敗は作成結果に影響させない（移行前と同じ）。
func TestChatRoomCreate_MemberAddFailureIsIgnored(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms", h.Create)

	rooms.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		args.Get(1).(*model.ChatRoom).ID = 10
	}).Return(nil)
	rooms.On("AddMember", mock.Anything, uint(10), uint(1)).Return(nil)
	rooms.On("AddMember", mock.Anything, uint(10), uint(2)).Return(errors.New("db error"))
	rooms.On("FindByID", mock.Anything, uint(10)).Return(&model.ChatRoom{Name: "Room"}, nil)

	w := doRequest(r, http.MethodPost, "/rooms", map[string]interface{}{
		"name": "Room", "member_ids": []uint{2},
	})
	assertStatus(t, w, http.StatusCreated)
}

// オーナーの追加に失敗した場合は作成自体を失敗として扱う。
func TestChatRoomCreate_OwnerAddFailure(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms", h.Create)

	rooms.On("Create", mock.Anything, mock.Anything).Return(nil)
	rooms.On("AddMember", mock.Anything, mock.Anything, uint(1)).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/rooms", map[string]interface{}{"name": "Room"})
	assertStatus(t, w, http.StatusInternalServerError)
	rooms.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// 再取得できなかった場合は作成直後のルームを返す。
func TestChatRoomCreate_RefetchFailureReturnsCreated(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms", h.Create)

	rooms.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		args.Get(1).(*model.ChatRoom).ID = 10
	}).Return(nil)
	rooms.On("AddMember", mock.Anything, uint(10), uint(1)).Return(nil)
	rooms.On("FindByID", mock.Anything, uint(10)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/rooms", map[string]interface{}{"name": "Room"})
	assertStatus(t, w, http.StatusCreated)
	assert.Contains(t, w.Body.String(), `"name":"Room"`)
}

func TestChatRoomCreate_RepositoryError(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms", h.Create)

	rooms.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/rooms", map[string]interface{}{"name": "Room"})
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestChatRoomCreate_ValidationError(t *testing.T) {
	h, _, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms", h.Create)

	// name は required
	w := doRequest(r, http.MethodPost, "/rooms", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

// リクエストの binding（200 文字）を通過しても、ルーム名は 100 文字までに制限される。
func TestChatRoomCreate_NameTooLong(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms", h.Create)

	w := doRequest(r, http.MethodPost, "/rooms", map[string]interface{}{
		"name": strings.Repeat("あ", 101),
	})
	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "チャットルーム名は100文字以下である必要があります")
	rooms.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestChatRoomCreate_DescriptionTooLong(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms", h.Create)

	w := doRequest(r, http.MethodPost, "/rooms", map[string]interface{}{
		"name": "Room", "description": strings.Repeat("あ", 501),
	})
	assertStatus(t, w, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), "説明は500文字以下である必要があります")
	rooms.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestChatRoomCreate_InvalidJSON(t *testing.T) {
	h, _, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/rooms", "bad json")
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- GetMyRooms ----------

func TestChatRoomGetMyRooms_Success(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.GET("/rooms", h.GetMyRooms)

	rooms.On("FindByUserID", mock.Anything, uint(1), 20, 0).Return([]model.ChatRoom{
		{Name: "Room A"}, {Name: "Room B"},
	}, int64(2), nil)

	w := doRequest(r, http.MethodGet, "/rooms", nil)
	assertStatus(t, w, http.StatusOK)

	var resp map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	rooms2 := resp["rooms"].([]interface{})
	assert.Len(t, rooms2, 2)
	assert.Equal(t, float64(2), resp["total"])
	assert.Equal(t, float64(20), resp["limit"])
	assert.Equal(t, float64(0), resp["offset"])
}

func TestChatRoomGetMyRooms_Pagination(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.GET("/rooms", h.GetMyRooms)

	rooms.On("FindByUserID", mock.Anything, uint(1), 5, 10).Return([]model.ChatRoom{}, int64(0), nil)

	w := doRequest(r, http.MethodGet, "/rooms?limit=5&offset=10", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"rooms":[]`)
	rooms.AssertExpectations(t)
}

func TestChatRoomGetMyRooms_RepositoryError(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.GET("/rooms", h.GetMyRooms)

	rooms.On("FindByUserID", mock.Anything, uint(1), 20, 0).
		Return([]model.ChatRoom{}, int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/rooms", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- GetByID ----------

func TestChatRoomGetByID_Success(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.GET("/rooms/:id", h.GetByID)

	rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
	rooms.On("FindByID", mock.Anything, uint(10)).Return(&model.ChatRoom{Name: "Found Room"}, nil)

	w := doRequest(r, http.MethodGet, "/rooms/10", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"name":"Found Room"`)
}

func TestChatRoomGetByID_NotMember(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.GET("/rooms/:id", h.GetByID)

	rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(false, nil)

	w := doRequest(r, http.MethodGet, "/rooms/10", nil)
	assertStatus(t, w, http.StatusForbidden)
	rooms.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
}

// メンバー判定に失敗した場合も 403 に倒す（判定エラーを漏らさない）。
func TestChatRoomGetByID_IsMemberError(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.GET("/rooms/:id", h.GetByID)

	rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(false, errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/rooms/10", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// メンバーだがルームが消えている場合は内部エラーとして扱う。
func TestChatRoomGetByID_RoomNotFound(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.GET("/rooms/:id", h.GetByID)

	rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
	rooms.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)

	w := doRequest(r, http.MethodGet, "/rooms/10", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestChatRoomGetByID_InvalidID(t *testing.T) {
	h, _, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.GET("/rooms/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/rooms/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- Update ----------

func TestChatRoomUpdate_Success(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.PUT("/rooms/:id", h.Update)

	room := &model.ChatRoom{Name: "Old Name", Description: "Old desc", OwnerID: 1}
	room.ID = 10
	rooms.On("FindByID", mock.Anything, uint(10)).Return(room, nil)
	rooms.On("Update", mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(r *model.ChatRoom) bool {
			return r.Name == "New Name" && r.Description == "Updated desc"
		})).Return(nil)

	w := doRequest(r, http.MethodPut, "/rooms/10", map[string]string{
		"name": " New Name ", "description": " Updated desc ",
	})
	assertStatus(t, w, http.StatusOK)
	rooms.AssertExpectations(t)
}

// 空文字の項目は据え置く。
func TestChatRoomUpdate_PartialKeepsExisting(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.PUT("/rooms/:id", h.Update)

	room := &model.ChatRoom{Name: "Old Name", Description: "Old desc", OwnerID: 1}
	room.ID = 10
	rooms.On("FindByID", mock.Anything, uint(10)).Return(room, nil)
	rooms.On("Update", mock.Anything, mock.MatchedBy(func(r *model.ChatRoom) bool {
		return r.Name == "Old Name" && r.Description == "New desc"
	})).Return(nil)

	w := doRequest(r, http.MethodPut, "/rooms/10", map[string]string{"description": "New desc"})
	assertStatus(t, w, http.StatusOK)
	rooms.AssertExpectations(t)
}

func TestChatRoomUpdate_Forbidden(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.PUT("/rooms/:id", h.Update)

	room := &model.ChatRoom{Name: "Room", OwnerID: 999} // 別ユーザーがオーナー
	room.ID = 10
	rooms.On("FindByID", mock.Anything, uint(10)).Return(room, nil)

	w := doRequest(r, http.MethodPut, "/rooms/10", map[string]string{"name": "Hacked"})
	assertStatus(t, w, http.StatusForbidden)
	rooms.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

// ルームが存在しない場合は内部エラーとして扱う（移行前も DB のエラーがそのまま 500 になっていた）。
func TestChatRoomUpdate_NotFound(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.PUT("/rooms/:id", h.Update)

	rooms.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)

	w := doRequest(r, http.MethodPut, "/rooms/10", map[string]string{"name": "X"})
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestChatRoomUpdate_NameTooLong(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.PUT("/rooms/:id", h.Update)

	room := &model.ChatRoom{Name: "Room", OwnerID: 1}
	room.ID = 10
	rooms.On("FindByID", mock.Anything, uint(10)).Return(room, nil)

	w := doRequest(r, http.MethodPut, "/rooms/10", map[string]string{
		"name": strings.Repeat("あ", 101),
	})
	assertStatus(t, w, http.StatusBadRequest)
	rooms.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestChatRoomUpdate_RepositoryError(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.PUT("/rooms/:id", h.Update)

	room := &model.ChatRoom{Name: "Room", OwnerID: 1}
	room.ID = 10
	rooms.On("FindByID", mock.Anything, uint(10)).Return(room, nil)
	rooms.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPut, "/rooms/10", map[string]string{"name": "New"})
	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- Delete ----------

func TestChatRoomDelete_Success(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.DELETE("/rooms/:id", h.Delete)

	room := &model.ChatRoom{OwnerID: 1}
	room.ID = 10
	rooms.On("FindByID", mock.Anything, uint(10)).Return(room, nil)
	rooms.On("Delete", mock.Anything, uint(10)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/rooms/10", nil)
	assertStatus(t, w, http.StatusOK)
	rooms.AssertExpectations(t)
}

func TestChatRoomDelete_Forbidden(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.DELETE("/rooms/:id", h.Delete)

	room := &model.ChatRoom{OwnerID: 999}
	room.ID = 10
	rooms.On("FindByID", mock.Anything, uint(10)).Return(room, nil)

	w := doRequest(r, http.MethodDelete, "/rooms/10", nil)
	assertStatus(t, w, http.StatusForbidden)
	rooms.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

func TestChatRoomDelete_NotFound(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.DELETE("/rooms/:id", h.Delete)

	rooms.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)

	w := doRequest(r, http.MethodDelete, "/rooms/10", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	rooms.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
}

// ---------- GetMembers ----------

func TestChatRoomGetMembers_Success(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.GET("/rooms/:id/members", h.GetMembers)

	rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
	rooms.On("GetMembers", mock.Anything, uint(10)).Return([]model.ChatRoomMember{
		{UserID: 1}, {UserID: 2},
	}, nil)

	w := doRequest(r, http.MethodGet, "/rooms/10/members", nil)
	assertStatus(t, w, http.StatusOK)
	rooms.AssertExpectations(t)
}

func TestChatRoomGetMembers_NotMember(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.GET("/rooms/:id/members", h.GetMembers)

	rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(false, nil)

	w := doRequest(r, http.MethodGet, "/rooms/10/members", nil)
	assertStatus(t, w, http.StatusForbidden)
	rooms.AssertNotCalled(t, "GetMembers", mock.Anything, mock.Anything)
}

// ---------- AddMember ----------

func TestChatRoomAddMember_Success(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms/:id/members", h.AddMember)

	rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
	rooms.On("IsMember", mock.Anything, uint(10), uint(5)).Return(false, nil)
	rooms.On("AddMember", mock.Anything, uint(10), uint(5)).Return(nil)

	w := doRequest(r, http.MethodPost, "/rooms/10/members", map[string]uint{"user_id": 5})
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "member added")
	rooms.AssertExpectations(t)
}

func TestChatRoomAddMember_NotMember(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms/:id/members", h.AddMember)

	rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(false, nil)

	w := doRequest(r, http.MethodPost, "/rooms/10/members", map[string]uint{"user_id": 5})
	assertStatus(t, w, http.StatusForbidden)
	rooms.AssertNotCalled(t, "AddMember", mock.Anything, mock.Anything, mock.Anything)
}

// 既に参加しているユーザーは追加できない。
func TestChatRoomAddMember_AlreadyMember(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms/:id/members", h.AddMember)

	rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
	rooms.On("IsMember", mock.Anything, uint(10), uint(5)).Return(true, nil)

	w := doRequest(r, http.MethodPost, "/rooms/10/members", map[string]uint{"user_id": 5})
	assertStatus(t, w, http.StatusBadRequest)
	rooms.AssertNotCalled(t, "AddMember", mock.Anything, mock.Anything, mock.Anything)
}

// 追加対象の参加判定に失敗した場合は 500（リクエスト者の判定とは扱いが異なる）。
func TestChatRoomAddMember_TargetIsMemberError(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms/:id/members", h.AddMember)

	rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
	rooms.On("IsMember", mock.Anything, uint(10), uint(5)).Return(false, errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/rooms/10/members", map[string]uint{"user_id": 5})
	assertStatus(t, w, http.StatusInternalServerError)
}

func TestChatRoomAddMember_ValidationError(t *testing.T) {
	h, _, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms/:id/members", h.AddMember)

	// user_id は required
	w := doRequest(r, http.MethodPost, "/rooms/10/members", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- RemoveMember ----------

func TestChatRoomRemoveMember_OwnerRemoves(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.DELETE("/rooms/:id/members/:userId", h.RemoveMember)

	room := &model.ChatRoom{OwnerID: 1}
	room.ID = 10
	rooms.On("FindByID", mock.Anything, uint(10)).Return(room, nil)
	rooms.On("RemoveMember", mock.Anything, uint(10), uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/rooms/10/members/5", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), "member removed")
	rooms.AssertExpectations(t)
}

func TestChatRoomRemoveMember_SelfLeave(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.DELETE("/rooms/:id/members/:userId", h.RemoveMember)

	room := &model.ChatRoom{OwnerID: 999} // 別ユーザーがオーナー
	room.ID = 10
	rooms.On("FindByID", mock.Anything, uint(10)).Return(room, nil)
	rooms.On("RemoveMember", mock.Anything, uint(10), uint(1)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/rooms/10/members/1", nil)
	assertStatus(t, w, http.StatusOK)
	rooms.AssertExpectations(t)
}

func TestChatRoomRemoveMember_Forbidden(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1) // userID=1
	r.DELETE("/rooms/:id/members/:userId", h.RemoveMember)

	room := &model.ChatRoom{OwnerID: 999}
	room.ID = 10
	rooms.On("FindByID", mock.Anything, uint(10)).Return(room, nil)

	// 非オーナーが他人を削除しようとする
	w := doRequest(r, http.MethodDelete, "/rooms/10/members/5", nil)
	assertStatus(t, w, http.StatusForbidden)
	rooms.AssertNotCalled(t, "RemoveMember", mock.Anything, mock.Anything, mock.Anything)
}

// オーナー自身は除外できない。
func TestChatRoomRemoveMember_TargetIsOwner(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.DELETE("/rooms/:id/members/:userId", h.RemoveMember)

	room := &model.ChatRoom{OwnerID: 1}
	room.ID = 10
	rooms.On("FindByID", mock.Anything, uint(10)).Return(room, nil)

	w := doRequest(r, http.MethodDelete, "/rooms/10/members/1", nil)
	assertStatus(t, w, http.StatusBadRequest)
	rooms.AssertNotCalled(t, "RemoveMember", mock.Anything, mock.Anything, mock.Anything)
}

func TestChatRoomRemoveMember_RoomNotFound(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.DELETE("/rooms/:id/members/:userId", h.RemoveMember)

	rooms.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)

	w := doRequest(r, http.MethodDelete, "/rooms/10/members/5", nil)
	assertStatus(t, w, http.StatusInternalServerError)
}

// ---------- GetMessages ----------

func TestChatRoomGetMessages_Success(t *testing.T) {
	h, rooms, messages, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.GET("/rooms/:id/messages", h.GetMessages)

	rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
	messages.On("FindByRoomID", mock.Anything, uint(10), 1, 20).Return([]model.GroupMessage{
		{Content: "Hello"},
	}, nil)

	w := doRequest(r, http.MethodGet, "/rooms/10/messages", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"content":"Hello"`)
	messages.AssertExpectations(t)
}

func TestChatRoomGetMessages_NotMember(t *testing.T) {
	h, rooms, messages, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.GET("/rooms/:id/messages", h.GetMessages)

	rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(false, nil)

	w := doRequest(r, http.MethodGet, "/rooms/10/messages", nil)
	assertStatus(t, w, http.StatusForbidden)
	messages.AssertNotCalled(t, "FindByRoomID", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestChatRoomGetMessages_WithPagination(t *testing.T) {
	h, rooms, messages, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.GET("/rooms/:id/messages", h.GetMessages)

	rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
	messages.On("FindByRoomID", mock.Anything, uint(10), 2, 30).Return([]model.GroupMessage{}, nil)

	w := doRequest(r, http.MethodGet, "/rooms/10/messages?page=2&limit=30", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Equal(t, "[]", w.Body.String())
	messages.AssertExpectations(t)
}

// ---------- SendMessage ----------

func TestChatRoomSendMessage_Success(t *testing.T) {
	h, rooms, messages, broadcaster := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms/:id/messages", h.SendMessage)

	rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
	messages.On("Create", mock.Anything, mock.MatchedBy(func(msg *model.GroupMessage) bool {
		// 前後の空白は落として保存する。
		return msg.ChatRoomID == 10 && msg.SenderID == 1 && msg.Content == "Hello everyone!"
	})).Return(nil)
	messages.On("FindSender", mock.Anything, uint(1)).Return(&model.User{Name: "送信者"}, nil)

	w := doRequest(r, http.MethodPost, "/rooms/10/messages", map[string]string{
		"content": "  Hello everyone!  ",
	})
	assertStatus(t, w, http.StatusCreated)
	assert.Contains(t, w.Body.String(), `"content":"Hello everyone!"`)

	// 送信者を除いたルーム参加者へ、WebSocket と同じ形式で配信する。
	call := broadcaster.wait(t)
	assert.Equal(t, uint(10), call.roomID)
	assert.Equal(t, uint(1), call.senderID)
	assert.JSONEq(t,
		`{"type":"group_message","sender_id":1,"receiver_id":0,"room_id":10,"content":"Hello everyone!","sender_name":"送信者"}`,
		string(call.message))
	messages.AssertExpectations(t)
}

// 送信者情報が引けなくてもメッセージ自体は返す。
func TestChatRoomSendMessage_SenderLookupFailure(t *testing.T) {
	h, rooms, messages, broadcaster := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms/:id/messages", h.SendMessage)

	rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
	messages.On("Create", mock.Anything, mock.Anything).Return(nil)
	messages.On("FindSender", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/rooms/10/messages", map[string]string{"content": "Hi"})
	assertStatus(t, w, http.StatusCreated)

	call := broadcaster.wait(t)
	assert.NotContains(t, string(call.message), "sender_name")
}

func TestChatRoomSendMessage_NotMember(t *testing.T) {
	h, rooms, messages, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms/:id/messages", h.SendMessage)

	rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(false, nil)

	w := doRequest(r, http.MethodPost, "/rooms/10/messages", map[string]string{"content": "Hello"})
	assertStatus(t, w, http.StatusForbidden)
	messages.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestChatRoomSendMessage_CreateError(t *testing.T) {
	h, rooms, messages, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms/:id/messages", h.SendMessage)

	rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
	messages.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	w := doRequest(r, http.MethodPost, "/rooms/10/messages", map[string]string{"content": "Hello"})
	assertStatus(t, w, http.StatusInternalServerError)
	messages.AssertNotCalled(t, "FindSender", mock.Anything, mock.Anything)
}

func TestChatRoomSendMessage_ValidationError(t *testing.T) {
	h, _, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms/:id/messages", h.SendMessage)

	// content は required
	w := doRequest(r, http.MethodPost, "/rooms/10/messages", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

// 空白だけの本文は空文字とみなして弾く。
func TestChatRoomSendMessage_BlankContent(t *testing.T) {
	h, rooms, messages, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.POST("/rooms/:id/messages", h.SendMessage)

	w := doRequest(r, http.MethodPost, "/rooms/10/messages", map[string]string{"content": "   "})
	assertStatus(t, w, http.StatusBadRequest)
	rooms.AssertNotCalled(t, "IsMember", mock.Anything, mock.Anything, mock.Anything)
	messages.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// ---------- GetMyCount ----------

func TestChatRoomGetMyCount_Success(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.GET("/chat-rooms/my/count", h.GetMyCount)

	rooms.On("CountByUserID", mock.Anything, uint(1)).Return(int64(3), nil)

	w := doRequest(r, http.MethodGet, "/chat-rooms/my/count", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"count":3`)
	rooms.AssertExpectations(t)
}

func TestChatRoomGetMyCount_RepositoryError(t *testing.T) {
	h, rooms, _, _ := newTestChatRoomHandler()
	r := newRouter(1)
	r.GET("/chat-rooms/my/count", h.GetMyCount)

	rooms.On("CountByUserID", mock.Anything, uint(1)).Return(int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/chat-rooms/my/count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	rooms.AssertExpectations(t)
}

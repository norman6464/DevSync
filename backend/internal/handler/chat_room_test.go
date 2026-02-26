package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------- Create ----------

func TestChatRoomCreate_Success(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.POST("/rooms", h.Create)

	roomRepo.On("Create", mock.AnythingOfType("*model.ChatRoom")).Return(nil)
	roomRepo.On("AddMember", mock.AnythingOfType("uint"), mock.AnythingOfType("uint")).Return(nil)
	roomRepo.On("FindByID", mock.AnythingOfType("uint")).Return(&model.ChatRoom{
		Name: "Test Room", OwnerID: 1,
	}, nil)

	w := doRequest(r, http.MethodPost, "/rooms", map[string]interface{}{
		"name": "Test Room", "member_ids": []uint{2, 3},
	})
	assertStatus(t, w, http.StatusCreated)
}

func TestChatRoomCreate_ValidationError(t *testing.T) {
	h, _, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.POST("/rooms", h.Create)

	// name は required
	w := doRequest(r, http.MethodPost, "/rooms", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

func TestChatRoomCreate_InvalidJSON(t *testing.T) {
	h, _, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.POST("/rooms", h.Create)

	w := doRequestRaw(r, http.MethodPost, "/rooms", "bad json")
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- GetMyRooms ----------

func TestChatRoomGetMyRooms_Success(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.GET("/rooms", h.GetMyRooms)

	roomRepo.On("FindByUserID", uint(1), 20, 0).Return([]model.ChatRoom{
		{Name: "Room A"}, {Name: "Room B"},
	}, int64(2), nil)

	w := doRequest(r, http.MethodGet, "/rooms", nil)
	assertStatus(t, w, http.StatusOK)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	rooms := resp["rooms"].([]interface{})
	assert.Len(t, rooms, 2)
	assert.Equal(t, float64(2), resp["total"])
	assert.Equal(t, float64(20), resp["limit"])
	assert.Equal(t, float64(0), resp["offset"])
}

// ---------- GetByID ----------

func TestChatRoomGetByID_Success(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.GET("/rooms/:id", h.GetByID)

	roomRepo.On("IsMember", uint(10), uint(1)).Return(true, nil)
	roomRepo.On("FindByID", uint(10)).Return(&model.ChatRoom{
		Name: "Found Room",
	}, nil)

	w := doRequest(r, http.MethodGet, "/rooms/10", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestChatRoomGetByID_NotMember(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.GET("/rooms/:id", h.GetByID)

	roomRepo.On("IsMember", uint(10), uint(1)).Return(false, nil)

	w := doRequest(r, http.MethodGet, "/rooms/10", nil)
	assertStatus(t, w, http.StatusForbidden)
}

func TestChatRoomGetByID_InvalidID(t *testing.T) {
	h, _, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.GET("/rooms/:id", h.GetByID)

	w := doRequest(r, http.MethodGet, "/rooms/abc", nil)
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- Update ----------

func TestChatRoomUpdate_Success(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.PUT("/rooms/:id", h.Update)

	room := &model.ChatRoom{Name: "Old Name", OwnerID: 1}
	room.ID = 10
	roomRepo.On("FindByID", uint(10)).Return(room, nil)
	roomRepo.On("Update", mock.AnythingOfType("*model.ChatRoom")).Return(nil)

	w := doRequest(r, http.MethodPut, "/rooms/10", map[string]string{
		"name": "New Name", "description": "Updated desc",
	})
	assertStatus(t, w, http.StatusOK)
}

func TestChatRoomUpdate_Forbidden(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.PUT("/rooms/:id", h.Update)

	room := &model.ChatRoom{Name: "Room", OwnerID: 999} // 別ユーザーがオーナー
	room.ID = 10
	roomRepo.On("FindByID", uint(10)).Return(room, nil)

	w := doRequest(r, http.MethodPut, "/rooms/10", map[string]string{
		"name": "Hacked",
	})
	assertStatus(t, w, http.StatusForbidden)
}

func TestChatRoomUpdate_NotFound(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.PUT("/rooms/:id", h.Update)

	roomRepo.On("FindByID", uint(10)).Return(nil, service.ErrNotFound)

	w := doRequest(r, http.MethodPut, "/rooms/10", map[string]string{"name": "X"})
	assertStatus(t, w, http.StatusNotFound)
}

// ---------- Delete ----------

func TestChatRoomDelete_Success(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.DELETE("/rooms/:id", h.Delete)

	room := &model.ChatRoom{OwnerID: 1}
	room.ID = 10
	roomRepo.On("FindByID", uint(10)).Return(room, nil)
	roomRepo.On("Delete", uint(10)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/rooms/10", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestChatRoomDelete_Forbidden(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.DELETE("/rooms/:id", h.Delete)

	room := &model.ChatRoom{OwnerID: 999}
	room.ID = 10
	roomRepo.On("FindByID", uint(10)).Return(room, nil)

	w := doRequest(r, http.MethodDelete, "/rooms/10", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- GetMembers ----------

func TestChatRoomGetMembers_Success(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.GET("/rooms/:id/members", h.GetMembers)

	roomRepo.On("IsMember", uint(10), uint(1)).Return(true, nil)
	roomRepo.On("GetMembers", uint(10)).Return([]model.ChatRoomMember{
		{UserID: 1}, {UserID: 2},
	}, nil)

	w := doRequest(r, http.MethodGet, "/rooms/10/members", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestChatRoomGetMembers_NotMember(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.GET("/rooms/:id/members", h.GetMembers)

	roomRepo.On("IsMember", uint(10), uint(1)).Return(false, nil)

	w := doRequest(r, http.MethodGet, "/rooms/10/members", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- AddMember ----------

func TestChatRoomAddMember_Success(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.POST("/rooms/:id/members", h.AddMember)

	roomRepo.On("IsMember", uint(10), uint(1)).Return(true, nil)
	roomRepo.On("IsMember", uint(10), uint(5)).Return(false, nil)
	roomRepo.On("AddMember", uint(10), uint(5)).Return(nil)

	w := doRequest(r, http.MethodPost, "/rooms/10/members", map[string]uint{
		"user_id": 5,
	})
	assertStatus(t, w, http.StatusOK)
}

func TestChatRoomAddMember_NotMember(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.POST("/rooms/:id/members", h.AddMember)

	roomRepo.On("IsMember", uint(10), uint(1)).Return(false, nil)

	w := doRequest(r, http.MethodPost, "/rooms/10/members", map[string]uint{
		"user_id": 5,
	})
	assertStatus(t, w, http.StatusForbidden)
}

func TestChatRoomAddMember_ValidationError(t *testing.T) {
	h, _, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.POST("/rooms/:id/members", h.AddMember)

	// user_id は required
	w := doRequest(r, http.MethodPost, "/rooms/10/members", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

// ---------- RemoveMember ----------

func TestChatRoomRemoveMember_OwnerRemoves(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.DELETE("/rooms/:id/members/:userId", h.RemoveMember)

	room := &model.ChatRoom{OwnerID: 1}
	room.ID = 10
	roomRepo.On("FindByID", uint(10)).Return(room, nil)
	roomRepo.On("RemoveMember", uint(10), uint(5)).Return(nil)

	w := doRequest(r, http.MethodDelete, "/rooms/10/members/5", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestChatRoomRemoveMember_SelfLeave(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.DELETE("/rooms/:id/members/:userId", h.RemoveMember)

	room := &model.ChatRoom{OwnerID: 999} // 別ユーザーがオーナー
	room.ID = 10
	roomRepo.On("FindByID", uint(10)).Return(room, nil)
	roomRepo.On("RemoveMember", uint(10), uint(1)).Return(nil)

	// 自分自身を退出
	w := doRequest(r, http.MethodDelete, "/rooms/10/members/1", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestChatRoomRemoveMember_Forbidden(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1) // userID=1
	r.DELETE("/rooms/:id/members/:userId", h.RemoveMember)

	room := &model.ChatRoom{OwnerID: 999}
	room.ID = 10
	roomRepo.On("FindByID", uint(10)).Return(room, nil)

	// 非オーナーが他人を削除しようとする
	w := doRequest(r, http.MethodDelete, "/rooms/10/members/5", nil)
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- GetMessages ----------

func TestChatRoomGetMessages_Success(t *testing.T) {
	h, roomRepo, msgRepo := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.GET("/rooms/:id/messages", h.GetMessages)

	roomRepo.On("IsMember", uint(10), uint(1)).Return(true, nil)
	msgRepo.On("FindByRoomID", uint(10), 1, 20).Return([]model.GroupMessage{
		{Content: "Hello"},
	}, nil)

	w := doRequest(r, http.MethodGet, "/rooms/10/messages", nil)
	assertStatus(t, w, http.StatusOK)
}

func TestChatRoomGetMessages_NotMember(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.GET("/rooms/:id/messages", h.GetMessages)

	roomRepo.On("IsMember", uint(10), uint(1)).Return(false, nil)

	w := doRequest(r, http.MethodGet, "/rooms/10/messages", nil)
	assertStatus(t, w, http.StatusForbidden)
}

func TestChatRoomGetMessages_WithPagination(t *testing.T) {
	h, roomRepo, msgRepo := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.GET("/rooms/:id/messages", h.GetMessages)

	roomRepo.On("IsMember", uint(10), uint(1)).Return(true, nil)
	msgRepo.On("FindByRoomID", uint(10), 2, 30).Return([]model.GroupMessage{}, nil)

	w := doRequest(r, http.MethodGet, "/rooms/10/messages?page=2&limit=30", nil)
	assertStatus(t, w, http.StatusOK)
}

// ---------- SendMessage ----------

func TestChatRoomSendMessage_Success(t *testing.T) {
	h, roomRepo, msgRepo := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.POST("/rooms/:id/messages", h.SendMessage)

	roomRepo.On("IsMember", uint(10), uint(1)).Return(true, nil)
	msgRepo.On("Create", mock.AnythingOfType("*model.GroupMessage")).Return(nil)
	msgRepo.On("FindSenderByID", mock.AnythingOfType("*model.GroupMessage")).Return()
	msgRepo.On("GetMemberUserIDs", uint(10)).Return([]uint{1, 2})

	w := doRequest(r, http.MethodPost, "/rooms/10/messages", map[string]string{
		"content": "Hello everyone!",
	})
	assertStatus(t, w, http.StatusCreated)
}

func TestChatRoomSendMessage_NotMember(t *testing.T) {
	h, roomRepo, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.POST("/rooms/:id/messages", h.SendMessage)

	roomRepo.On("IsMember", uint(10), uint(1)).Return(false, nil)

	w := doRequest(r, http.MethodPost, "/rooms/10/messages", map[string]string{
		"content": "Hello",
	})
	assertStatus(t, w, http.StatusForbidden)
}

// ---------- GetMyCount ----------

func TestChatRoomGetMyCount_Success(t *testing.T) {
	h, svc := setupChatRoomHandler()
	r := newRouter(1)
	r.GET("/chat-rooms/my/count", h.GetMyCount)

	svc.On("CountByUserID", uint(1)).Return(int64(3), nil)

	w := doRequest(r, http.MethodGet, "/chat-rooms/my/count", nil)
	assertStatus(t, w, http.StatusOK)
	assert.Contains(t, w.Body.String(), `"count":3`)
	svc.AssertExpectations(t)
}

func TestChatRoomGetMyCount_ServiceError(t *testing.T) {
	h, svc := setupChatRoomHandler()
	r := newRouter(1)
	r.GET("/chat-rooms/my/count", h.GetMyCount)

	svc.On("CountByUserID", uint(1)).Return(int64(0), errors.New("db error"))

	w := doRequest(r, http.MethodGet, "/chat-rooms/my/count", nil)
	assertStatus(t, w, http.StatusInternalServerError)
	svc.AssertExpectations(t)
}

func TestChatRoomSendMessage_ValidationError(t *testing.T) {
	h, _, _ := setupChatRoomHandlerRepo()
	r := newRouter(1)
	r.POST("/rooms/:id/messages", h.SendMessage)

	// content は required
	w := doRequest(r, http.MethodPost, "/rooms/10/messages", map[string]string{})
	assertStatus(t, w, http.StatusBadRequest)
}

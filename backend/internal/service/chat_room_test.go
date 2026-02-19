package service

import (
	"errors"
	"testing"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newTestChatRoomService はChatRoomServiceのテスト用インスタンスを生成するヘルパー。
func newTestChatRoomService() (*ChatRoomService, *MockChatRoomRepository, *MockGroupMessageRepository) {
	roomRepo := new(MockChatRoomRepository)
	msgRepo := new(MockGroupMessageRepository)
	hub := NewHub() // 実際のHubを生成（Run()は呼ばない）
	svc := NewChatRoomService(roomRepo, msgRepo, hub)
	return svc, roomRepo, msgRepo
}

// ============================================================
// チャットルーム作成テスト
// ============================================================

func TestChatRoomCreate_Success(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	room := &model.ChatRoom{Name: "Test Room", OwnerID: 1}

	roomRepo.On("Create", room).Run(func(args mock.Arguments) {
		r := args.Get(0).(*model.ChatRoom)
		r.ID = 10
	}).Return(nil)
	roomRepo.On("AddMember", uint(10), uint(1)).Return(nil) // owner
	roomRepo.On("AddMember", uint(10), uint(2)).Return(nil) // member

	createdRoom := &model.ChatRoom{Name: "Test Room", OwnerID: 1}
	createdRoom.ID = 10
	roomRepo.On("FindByID", uint(10)).Return(createdRoom, nil)

	result, err := svc.Create(room, []uint{2})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(10), result.ID)
	roomRepo.AssertExpectations(t)
}

func TestChatRoomCreate_SkipDuplicateOwner(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	room := &model.ChatRoom{Name: "Room", OwnerID: 1}

	roomRepo.On("Create", room).Run(func(args mock.Arguments) {
		r := args.Get(0).(*model.ChatRoom)
		r.ID = 10
	}).Return(nil)
	roomRepo.On("AddMember", uint(10), uint(1)).Return(nil)

	createdRoom := &model.ChatRoom{Name: "Room", OwnerID: 1}
	createdRoom.ID = 10
	roomRepo.On("FindByID", uint(10)).Return(createdRoom, nil)

	// memberIDsにオーナーを含めても重複追加されない
	result, err := svc.Create(room, []uint{1})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	// AddMember は ownerID=1 の1回のみ呼ばれる
	roomRepo.AssertNumberOfCalls(t, "AddMember", 1)
}

func TestChatRoomCreate_RepoError(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	room := &model.ChatRoom{Name: "Test Room", OwnerID: 1}
	roomRepo.On("Create", room).Return(errors.New("db error"))

	result, err := svc.Create(room, []uint{2})
	assert.Error(t, err)
	assert.Nil(t, result)
	roomRepo.AssertExpectations(t)
}

func TestChatRoomCreate_AddMemberError(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	room := &model.ChatRoom{Name: "Test Room", OwnerID: 1}
	roomRepo.On("Create", room).Run(func(args mock.Arguments) {
		r := args.Get(0).(*model.ChatRoom)
		r.ID = 10
	}).Return(nil)
	roomRepo.On("AddMember", uint(10), uint(1)).Return(errors.New("add member error"))

	result, err := svc.Create(room, []uint{2})
	assert.Error(t, err)
	assert.Nil(t, result)
	roomRepo.AssertExpectations(t)
}

func TestChatRoomCreate_FindByIDError(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	room := &model.ChatRoom{Name: "Test Room", OwnerID: 1}
	roomRepo.On("Create", room).Run(func(args mock.Arguments) {
		r := args.Get(0).(*model.ChatRoom)
		r.ID = 10
	}).Return(nil)
	roomRepo.On("AddMember", uint(10), uint(1)).Return(nil)
	roomRepo.On("FindByID", uint(10)).Return(nil, errors.New("not found"))

	// FindByID失敗時はフォールバックでroomを返す
	result, err := svc.Create(room, nil)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	roomRepo.AssertExpectations(t)
}

// ============================================================
// GetByID（メンバーシップチェック）
// ============================================================

func TestChatRoomGetByID_MemberSuccess(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	roomRepo.On("IsMember", uint(10), uint(1)).Return(true, nil)

	room := &model.ChatRoom{Name: "Room", OwnerID: 1}
	room.ID = 10
	roomRepo.On("FindByID", uint(10)).Return(room, nil)

	result, err := svc.GetByID(10, 1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	roomRepo.AssertExpectations(t)
}

func TestChatRoomGetByID_NotMember(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	roomRepo.On("IsMember", uint(10), uint(999)).Return(false, nil)

	result, err := svc.GetByID(10, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	roomRepo.AssertExpectations(t)
}

// ============================================================
// Update（オーナーチェック）
// ============================================================

func TestChatRoomUpdate_Success(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	room := &model.ChatRoom{Name: "Old Name", OwnerID: 1}
	room.ID = 10

	roomRepo.On("FindByID", uint(10)).Return(room, nil)
	roomRepo.On("Update", room).Return(nil)

	result, err := svc.Update(10, 1, "New Name", "New Description")
	assert.NoError(t, err)
	assert.Equal(t, "New Name", result.Name)
	assert.Equal(t, "New Description", result.Description)
	roomRepo.AssertExpectations(t)
}

func TestChatRoomUpdate_NotOwner(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	room := &model.ChatRoom{Name: "Room", OwnerID: 1}
	room.ID = 10

	roomRepo.On("FindByID", uint(10)).Return(room, nil)

	result, err := svc.Update(10, 999, "New Name", "")
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	roomRepo.AssertExpectations(t)
}

// ============================================================
// Delete（オーナーチェック）
// ============================================================

func TestChatRoomDelete_Success(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	room := &model.ChatRoom{OwnerID: 1}
	room.ID = 10

	roomRepo.On("FindByID", uint(10)).Return(room, nil)
	roomRepo.On("Delete", uint(10)).Return(nil)

	err := svc.Delete(10, 1)
	assert.NoError(t, err)
	roomRepo.AssertExpectations(t)
}

func TestChatRoomDelete_NotOwner(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	room := &model.ChatRoom{OwnerID: 1}
	room.ID = 10

	roomRepo.On("FindByID", uint(10)).Return(room, nil)

	err := svc.Delete(10, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	roomRepo.AssertExpectations(t)
}

// ============================================================
// RemoveMember（権限分岐）
// ============================================================

func TestChatRoomRemoveMember_OwnerCanRemoveAnyone(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	room := &model.ChatRoom{OwnerID: 1}
	room.ID = 10

	roomRepo.On("FindByID", uint(10)).Return(room, nil)
	roomRepo.On("RemoveMember", uint(10), uint(2)).Return(nil)

	// オーナー(1)が他のメンバー(2)を削除
	err := svc.RemoveMember(10, 1, 2)
	assert.NoError(t, err)
	roomRepo.AssertExpectations(t)
}

func TestChatRoomRemoveMember_MemberCanRemoveSelf(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	room := &model.ChatRoom{OwnerID: 1}
	room.ID = 10

	roomRepo.On("FindByID", uint(10)).Return(room, nil)
	roomRepo.On("RemoveMember", uint(10), uint(2)).Return(nil)

	// メンバー(2)が自分自身を削除（退室）
	err := svc.RemoveMember(10, 2, 2)
	assert.NoError(t, err)
	roomRepo.AssertExpectations(t)
}

func TestChatRoomRemoveMember_OwnerCannotRemoveSelf(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	room := &model.ChatRoom{OwnerID: 1}
	room.ID = 10

	roomRepo.On("FindByID", uint(10)).Return(room, nil)

	// オーナー(1)が自分自身を退出しようとする
	err := svc.RemoveMember(10, 1, 1)
	assert.ErrorIs(t, err, ErrBadRequest)
	roomRepo.AssertNotCalled(t, "RemoveMember")
}

func TestChatRoomRemoveMember_NonOwnerCannotRemoveOthers(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	room := &model.ChatRoom{OwnerID: 1}
	room.ID = 10

	roomRepo.On("FindByID", uint(10)).Return(room, nil)

	// 非オーナー(2)が他のメンバー(3)を削除しようとする
	err := svc.RemoveMember(10, 2, 3)
	assert.ErrorIs(t, err, ErrForbidden)
	roomRepo.AssertExpectations(t)
}

// ============================================================
// AddMember（メンバーシップチェック）
// ============================================================

func TestChatRoomAddMember_Success(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	roomRepo.On("IsMember", uint(10), uint(1)).Return(true, nil)
	roomRepo.On("AddMember", uint(10), uint(3)).Return(nil)

	err := svc.AddMember(10, 1, 3)
	assert.NoError(t, err)
	roomRepo.AssertExpectations(t)
}

func TestChatRoomAddMember_NotMember(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	roomRepo.On("IsMember", uint(10), uint(999)).Return(false, nil)

	err := svc.AddMember(10, 999, 3)
	assert.ErrorIs(t, err, ErrForbidden)
	roomRepo.AssertExpectations(t)
}

// ============================================================
// GetMessages（メンバーシップチェック）
// ============================================================

func TestChatRoomGetMessages_Success(t *testing.T) {
	svc, roomRepo, msgRepo := newTestChatRoomService()

	roomRepo.On("IsMember", uint(10), uint(1)).Return(true, nil)

	messages := []model.GroupMessage{{Content: "Hello"}}
	msgRepo.On("FindByRoomID", uint(10), 1, 20).Return(messages, nil)

	result, err := svc.GetMessages(10, 1, 1, 20)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	roomRepo.AssertExpectations(t)
	msgRepo.AssertExpectations(t)
}

func TestChatRoomGetMessages_NotMember(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	roomRepo.On("IsMember", uint(10), uint(999)).Return(false, nil)

	result, err := svc.GetMessages(10, 999, 1, 20)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	roomRepo.AssertExpectations(t)
}

// ============================================================
// GetByUserID
// ============================================================

func TestChatRoomGetByUserID_Success(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	rooms := []model.ChatRoom{
		{Name: "Room A", OwnerID: 1},
		{Name: "Room B", OwnerID: 2},
	}
	roomRepo.On("FindByUserID", uint(1)).Return(rooms, nil)

	result, err := svc.GetByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Room A", result[0].Name)
	roomRepo.AssertExpectations(t)
}

func TestChatRoomGetByUserID_Empty(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	roomRepo.On("FindByUserID", uint(999)).Return([]model.ChatRoom{}, nil)

	result, err := svc.GetByUserID(999)
	assert.NoError(t, err)
	assert.Empty(t, result)
	roomRepo.AssertExpectations(t)
}

func TestChatRoomGetByUserID_RepoError(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	roomRepo.On("FindByUserID", uint(1)).Return([]model.ChatRoom{}, assert.AnError)

	result, err := svc.GetByUserID(1)
	assert.Error(t, err)
	assert.Empty(t, result)
	roomRepo.AssertExpectations(t)
}

// ============================================================
// GetMembers（メンバーシップチェック）
// ============================================================

func TestChatRoomGetMembers_Success(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	roomRepo.On("IsMember", uint(10), uint(1)).Return(true, nil)

	members := []model.ChatRoomMember{
		{ChatRoomID: 10, UserID: 1},
		{ChatRoomID: 10, UserID: 2},
	}
	roomRepo.On("GetMembers", uint(10)).Return(members, nil)

	result, err := svc.GetMembers(10, 1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	roomRepo.AssertExpectations(t)
}

func TestChatRoomGetMembers_NotMember(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	roomRepo.On("IsMember", uint(10), uint(999)).Return(false, nil)

	result, err := svc.GetMembers(10, 999)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	roomRepo.AssertExpectations(t)
}

func TestChatRoomGetMembers_IsMemberError(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	roomRepo.On("IsMember", uint(10), uint(1)).Return(false, assert.AnError)

	result, err := svc.GetMembers(10, 1)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	roomRepo.AssertExpectations(t)
}

// ============================================================
// SendMessage（メンバーシップチェック + メッセージ作成 + WebSocket配信）
// ============================================================

func TestChatRoomSendMessage_Success(t *testing.T) {
	svc, roomRepo, msgRepo := newTestChatRoomService()

	roomRepo.On("IsMember", uint(10), uint(1)).Return(true, nil)

	msgRepo.On("Create", mock.MatchedBy(func(msg *model.GroupMessage) bool {
		return msg.ChatRoomID == 10 && msg.SenderID == 1 && msg.Content == "Hello!"
	})).Return(nil)

	msgRepo.On("FindSenderByID", mock.AnythingOfType("*model.GroupMessage")).Return()

	result, err := svc.SendMessage(10, 1, "Hello!")
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(10), result.ChatRoomID)
	assert.Equal(t, uint(1), result.SenderID)
	assert.Equal(t, "Hello!", result.Content)
	roomRepo.AssertExpectations(t)
	msgRepo.AssertExpectations(t)
}

func TestChatRoomSendMessage_NotMember(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	roomRepo.On("IsMember", uint(10), uint(999)).Return(false, nil)

	result, err := svc.SendMessage(10, 999, "Hello!")
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Nil(t, result)
	roomRepo.AssertExpectations(t)
}

func TestChatRoomSendMessage_CreateError(t *testing.T) {
	svc, roomRepo, msgRepo := newTestChatRoomService()

	roomRepo.On("IsMember", uint(10), uint(1)).Return(true, nil)

	msgRepo.On("Create", mock.MatchedBy(func(msg *model.GroupMessage) bool {
		return msg.ChatRoomID == 10 && msg.SenderID == 1
	})).Return(assert.AnError)

	result, err := svc.SendMessage(10, 1, "Hello!")
	assert.Error(t, err)
	assert.Nil(t, result)
	roomRepo.AssertExpectations(t)
	msgRepo.AssertExpectations(t)
}

// ============================================================
// IsMember
// ============================================================

func TestChatRoomIsMember_True(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	roomRepo.On("IsMember", uint(10), uint(1)).Return(true, nil)

	result, err := svc.IsMember(10, 1)
	assert.NoError(t, err)
	assert.True(t, result)
	roomRepo.AssertExpectations(t)
}

func TestChatRoomIsMember_False(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	roomRepo.On("IsMember", uint(10), uint(999)).Return(false, nil)

	result, err := svc.IsMember(10, 999)
	assert.NoError(t, err)
	assert.False(t, result)
	roomRepo.AssertExpectations(t)
}

func TestChatRoomIsMember_RepoError(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()

	roomRepo.On("IsMember", uint(10), uint(1)).Return(false, assert.AnError)

	result, err := svc.IsMember(10, 1)
	assert.Error(t, err)
	assert.False(t, result)
	roomRepo.AssertExpectations(t)
}

func TestChatRoomUpdate_NotFound(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()
	roomRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))
	result, err := svc.Update(99, 1, "Name", "Desc")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestChatRoomUpdate_RepoError(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()
	room := &model.ChatRoom{Name: "Room", OwnerID: 1}
	room.ID = 10
	roomRepo.On("FindByID", uint(10)).Return(room, nil)
	roomRepo.On("Update", room).Return(errors.New("db error"))
	result, err := svc.Update(10, 1, "New", "Desc")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestChatRoomDelete_NotFound(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()
	roomRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))
	err := svc.Delete(99, 1)
	assert.Error(t, err)
}

func TestChatRoomRemoveMember_NotFound(t *testing.T) {
	svc, roomRepo, _ := newTestChatRoomService()
	roomRepo.On("FindByID", uint(99)).Return(nil, errors.New("not found"))
	err := svc.RemoveMember(99, 1, 2)
	assert.Error(t, err)
}

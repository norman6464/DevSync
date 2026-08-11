package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockChatRoomRepo は usecase/repository.ChatRoomRepository のモック。
type mockChatRoomRepo struct{ mock.Mock }

func (m *mockChatRoomRepo) Create(ctx context.Context, room *model.ChatRoom) error {
	return m.Called(ctx, room).Error(0)
}

func (m *mockChatRoomRepo) FindByID(ctx context.Context, id uint) (*model.ChatRoom, error) {
	args := m.Called(ctx, id)
	r, _ := args.Get(0).(*model.ChatRoom)
	return r, args.Error(1)
}

func (m *mockChatRoomRepo) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.ChatRoom, int64, error) {
	args := m.Called(ctx, userID, limit, offset)
	rooms, _ := args.Get(0).([]model.ChatRoom)
	return rooms, args.Get(1).(int64), args.Error(2)
}

func (m *mockChatRoomRepo) Update(ctx context.Context, room *model.ChatRoom) error {
	return m.Called(ctx, room).Error(0)
}

func (m *mockChatRoomRepo) Delete(ctx context.Context, roomID uint) error {
	return m.Called(ctx, roomID).Error(0)
}

func (m *mockChatRoomRepo) AddMember(ctx context.Context, roomID, userID uint) error {
	return m.Called(ctx, roomID, userID).Error(0)
}

func (m *mockChatRoomRepo) RemoveMember(ctx context.Context, roomID, userID uint) error {
	return m.Called(ctx, roomID, userID).Error(0)
}

func (m *mockChatRoomRepo) GetMembers(ctx context.Context, roomID uint) ([]model.ChatRoomMember, error) {
	args := m.Called(ctx, roomID)
	members, _ := args.Get(0).([]model.ChatRoomMember)
	return members, args.Error(1)
}

func (m *mockChatRoomRepo) IsMember(ctx context.Context, roomID, userID uint) (bool, error) {
	args := m.Called(ctx, roomID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *mockChatRoomRepo) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

// mockChatRoomMessageRepo は usecase/repository.ChatRoomMessageRepository のモック。
type mockChatRoomMessageRepo struct{ mock.Mock }

func (m *mockChatRoomMessageRepo) Create(ctx context.Context, msg *model.GroupMessage) error {
	return m.Called(ctx, msg).Error(0)
}

func (m *mockChatRoomMessageRepo) FindByRoomID(ctx context.Context, roomID uint, page, limit int) ([]model.GroupMessage, error) {
	args := m.Called(ctx, roomID, page, limit)
	msgs, _ := args.Get(0).([]model.GroupMessage)
	return msgs, args.Error(1)
}

func (m *mockChatRoomMessageRepo) FindSender(ctx context.Context, senderID uint) (*model.User, error) {
	args := m.Called(ctx, senderID)
	u, _ := args.Get(0).(*model.User)
	return u, args.Error(1)
}

// stubBroadcaster は配信内容を記録するテスト用の RoomBroadcaster。
type stubBroadcaster struct {
	sent chan []byte
}

func newStubBroadcaster() *stubBroadcaster {
	return &stubBroadcaster{sent: make(chan []byte, 4)}
}

func (b *stubBroadcaster) SendToRoom(roomID, senderID uint, message []byte) {
	b.sent <- message
}

// wait は配信されるまで待つ。配信が無ければテストを失敗させる。
func (b *stubBroadcaster) wait(t *testing.T) []byte {
	t.Helper()
	select {
	case msg := <-b.sent:
		return msg
	case <-time.After(2 * time.Second):
		t.Fatal("配信が行われなかった")
		return nil
	}
}

// ownedChatRoom は指定ユーザーがオーナーのルームを返すテスト用ヘルパー。
func ownedChatRoom(id, ownerID uint) *model.ChatRoom {
	room := &model.ChatRoom{Name: "Room", OwnerID: ownerID}
	room.ID = id
	return room
}

// ============================================================
// 作成
// ============================================================

func TestCreateChatRoomUseCase(t *testing.T) {
	t.Run("オーナーを追加してから再取得したルームを返す", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewCreateChatRoomUseCase(repo)

		repo.On("Create", mock.Anything, mock.MatchedBy(func(room *model.ChatRoom) bool {
			return room.Name == "Room" && room.Description == "desc" && room.OwnerID == 1
		})).Run(func(args mock.Arguments) {
			args.Get(1).(*model.ChatRoom).ID = 10
		}).Return(nil)
		repo.On("AddMember", mock.Anything, uint(10), uint(1)).Return(nil)
		repo.On("AddMember", mock.Anything, uint(10), uint(2)).Return(nil)
		repo.On("FindByID", mock.Anything, uint(10)).Return(ownedChatRoom(10, 1), nil)

		created, err := uc.Execute(context.Background(), usecase.CreateChatRoomInput{
			Name: "  Room  ", Description: " desc ", OwnerID: 1, MemberIDs: []uint{1, 2},
		})
		require.NoError(t, err)
		assert.Equal(t, uint(10), created.ID)
		// オーナーは member_ids に含まれていても重複追加しない。
		repo.AssertNumberOfCalls(t, "AddMember", 2)
		repo.AssertExpectations(t)
	})

	t.Run("ルーム名が空なら検証エラー", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewCreateChatRoomUseCase(repo)

		_, err := uc.Execute(context.Background(), usecase.CreateChatRoomInput{Name: "   ", OwnerID: 1})
		assert.Error(t, err)
		assert.Equal(t, domain.ErrCodeValidation, domain.GetDomainError(err).Code)
		repo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("再取得が不在なら作成直後のルームを返す", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewCreateChatRoomUseCase(repo)

		repo.On("Create", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
			args.Get(1).(*model.ChatRoom).ID = 10
		}).Return(nil)
		repo.On("AddMember", mock.Anything, uint(10), uint(1)).Return(nil)
		repo.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)

		created, err := uc.Execute(context.Background(), usecase.CreateChatRoomInput{Name: "Room", OwnerID: 1})
		require.NoError(t, err)
		assert.Equal(t, uint(10), created.ID)
		assert.Equal(t, "Room", created.Name)
	})
}

// ============================================================
// 取得・一覧
// ============================================================

func TestGetChatRoomUseCase(t *testing.T) {
	t.Run("メンバーならルームを返す", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewGetChatRoomUseCase(repo)

		repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
		repo.On("FindByID", mock.Anything, uint(10)).Return(ownedChatRoom(10, 1), nil)

		room, err := uc.Execute(context.Background(), 10, 1)
		require.NoError(t, err)
		assert.Equal(t, uint(10), room.ID)
	})

	t.Run("メンバー判定の失敗も 403 に倒す", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewGetChatRoomUseCase(repo)

		repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(false, errors.New("db error"))

		_, err := uc.Execute(context.Background(), 10, 1)
		assert.ErrorIs(t, err, domain.ErrForbidden)
		repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	})

	t.Run("ルームが不在なら DomainError にしない", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewGetChatRoomUseCase(repo)

		repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
		repo.On("FindByID", mock.Anything, uint(10)).Return(nil, nil)

		_, err := uc.Execute(context.Background(), 10, 1)
		require.Error(t, err)
		assert.Nil(t, domain.GetDomainError(err))
	})
}

func TestListMyChatRoomsUseCase(t *testing.T) {
	repo := new(mockChatRoomRepo)
	uc := usecase.NewListMyChatRoomsUseCase(repo)

	repo.On("FindByUserID", mock.Anything, uint(1), 20, 0).
		Return([]model.ChatRoom{{Name: "A"}}, int64(1), nil)

	rooms, total, err := uc.Execute(context.Background(), 1, 20, 0)
	require.NoError(t, err)
	assert.Len(t, rooms, 1)
	assert.Equal(t, int64(1), total)
	repo.AssertExpectations(t)
}

func TestCountMyChatRoomsUseCase(t *testing.T) {
	repo := new(mockChatRoomRepo)
	uc := usecase.NewCountMyChatRoomsUseCase(repo)

	repo.On("CountByUserID", mock.Anything, uint(1)).Return(int64(3), nil)

	count, err := uc.Execute(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

// ============================================================
// 更新・削除
// ============================================================

func TestUpdateChatRoomUseCase(t *testing.T) {
	t.Run("空でない項目だけを更新する", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewUpdateChatRoomUseCase(repo)

		room := ownedChatRoom(10, 1)
		room.Description = "old desc"
		repo.On("FindByID", mock.Anything, uint(10)).Return(room, nil)
		repo.On("Update", mock.Anything, room).Return(nil)

		updated, err := uc.Execute(context.Background(), 10, 1, "", " new desc ")
		require.NoError(t, err)
		assert.Equal(t, "Room", updated.Name)
		assert.Equal(t, "new desc", updated.Description)
	})

	t.Run("オーナー以外は 403", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewUpdateChatRoomUseCase(repo)

		repo.On("FindByID", mock.Anything, uint(10)).Return(ownedChatRoom(10, 999), nil)

		_, err := uc.Execute(context.Background(), 10, 1, "New", "")
		assert.ErrorIs(t, err, domain.ErrForbidden)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("説明が長すぎる場合は検証エラー", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewUpdateChatRoomUseCase(repo)

		repo.On("FindByID", mock.Anything, uint(10)).Return(ownedChatRoom(10, 1), nil)

		_, err := uc.Execute(context.Background(), 10, 1, "", strings.Repeat("あ", 501))
		require.Error(t, err)
		assert.Equal(t, domain.ErrCodeValidation, domain.GetDomainError(err).Code)
		repo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})
}

func TestDeleteChatRoomUseCase(t *testing.T) {
	t.Run("オーナーなら削除できる", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewDeleteChatRoomUseCase(repo)

		repo.On("FindByID", mock.Anything, uint(10)).Return(ownedChatRoom(10, 1), nil)
		repo.On("Delete", mock.Anything, uint(10)).Return(nil)

		require.NoError(t, uc.Execute(context.Background(), 10, 1))
		repo.AssertExpectations(t)
	})

	t.Run("オーナー以外は 403", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewDeleteChatRoomUseCase(repo)

		repo.On("FindByID", mock.Anything, uint(10)).Return(ownedChatRoom(10, 999), nil)

		assert.ErrorIs(t, uc.Execute(context.Background(), 10, 1), domain.ErrForbidden)
		repo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})
}

// ============================================================
// メンバー管理
// ============================================================

func TestListChatRoomMembersUseCase(t *testing.T) {
	repo := new(mockChatRoomRepo)
	uc := usecase.NewListChatRoomMembersUseCase(repo)

	repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
	repo.On("GetMembers", mock.Anything, uint(10)).Return([]model.ChatRoomMember{{UserID: 1}}, nil)

	members, err := uc.Execute(context.Background(), 10, 1)
	require.NoError(t, err)
	assert.Len(t, members, 1)
}

func TestAddChatRoomMemberUseCase(t *testing.T) {
	t.Run("メンバーが未参加のユーザーを追加できる", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewAddChatRoomMemberUseCase(repo)

		repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
		repo.On("IsMember", mock.Anything, uint(10), uint(5)).Return(false, nil)
		repo.On("AddMember", mock.Anything, uint(10), uint(5)).Return(nil)

		require.NoError(t, uc.Execute(context.Background(), 10, 1, 5))
		repo.AssertExpectations(t)
	})

	t.Run("既に参加済みなら 400", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewAddChatRoomMemberUseCase(repo)

		repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
		repo.On("IsMember", mock.Anything, uint(10), uint(5)).Return(true, nil)

		assert.ErrorIs(t, uc.Execute(context.Background(), 10, 1, 5), domain.ErrBadRequest)
	})

	t.Run("追加対象の判定エラーはそのまま返す", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewAddChatRoomMemberUseCase(repo)

		dbErr := errors.New("db error")
		repo.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
		repo.On("IsMember", mock.Anything, uint(10), uint(5)).Return(false, dbErr)

		assert.ErrorIs(t, uc.Execute(context.Background(), 10, 1, 5), dbErr)
	})
}

func TestRemoveChatRoomMemberUseCase(t *testing.T) {
	t.Run("オーナーは他のメンバーを除外できる", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewRemoveChatRoomMemberUseCase(repo)

		repo.On("FindByID", mock.Anything, uint(10)).Return(ownedChatRoom(10, 1), nil)
		repo.On("RemoveMember", mock.Anything, uint(10), uint(5)).Return(nil)

		require.NoError(t, uc.Execute(context.Background(), 10, 1, 5))
	})

	t.Run("一般メンバーは自分だけ退出できる", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewRemoveChatRoomMemberUseCase(repo)

		repo.On("FindByID", mock.Anything, uint(10)).Return(ownedChatRoom(10, 999), nil)
		repo.On("RemoveMember", mock.Anything, uint(10), uint(1)).Return(nil)

		require.NoError(t, uc.Execute(context.Background(), 10, 1, 1))
	})

	t.Run("一般メンバーが他人を除外しようとすると 403", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewRemoveChatRoomMemberUseCase(repo)

		repo.On("FindByID", mock.Anything, uint(10)).Return(ownedChatRoom(10, 999), nil)

		assert.ErrorIs(t, uc.Execute(context.Background(), 10, 1, 5), domain.ErrForbidden)
		repo.AssertNotCalled(t, "RemoveMember", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("オーナー自身は除外できない", func(t *testing.T) {
		repo := new(mockChatRoomRepo)
		uc := usecase.NewRemoveChatRoomMemberUseCase(repo)

		repo.On("FindByID", mock.Anything, uint(10)).Return(ownedChatRoom(10, 1), nil)

		assert.ErrorIs(t, uc.Execute(context.Background(), 10, 1, 1), domain.ErrBadRequest)
	})
}

// ============================================================
// メッセージ
// ============================================================

func TestListChatRoomMessagesUseCase(t *testing.T) {
	t.Run("メンバーならメッセージを取得できる", func(t *testing.T) {
		rooms := new(mockChatRoomRepo)
		messages := new(mockChatRoomMessageRepo)
		uc := usecase.NewListChatRoomMessagesUseCase(rooms, messages)

		rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
		messages.On("FindByRoomID", mock.Anything, uint(10), 2, 30).
			Return([]model.GroupMessage{{Content: "Hi"}}, nil)

		got, err := uc.Execute(context.Background(), 10, 1, 2, 30)
		require.NoError(t, err)
		assert.Len(t, got, 1)
		messages.AssertExpectations(t)
	})

	t.Run("メンバーでなければ 403", func(t *testing.T) {
		rooms := new(mockChatRoomRepo)
		messages := new(mockChatRoomMessageRepo)
		uc := usecase.NewListChatRoomMessagesUseCase(rooms, messages)

		rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(false, nil)

		_, err := uc.Execute(context.Background(), 10, 1, 1, 20)
		assert.ErrorIs(t, err, domain.ErrForbidden)
		messages.AssertNotCalled(t, "FindByRoomID", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestSendChatRoomMessageUseCase(t *testing.T) {
	t.Run("保存したメッセージを返し、送信者を除く参加者へ配信する", func(t *testing.T) {
		rooms := new(mockChatRoomRepo)
		messages := new(mockChatRoomMessageRepo)
		broadcaster := newStubBroadcaster()
		uc := usecase.NewSendChatRoomMessageUseCase(rooms, messages, broadcaster)

		rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
		messages.On("Create", mock.Anything, mock.MatchedBy(func(msg *model.GroupMessage) bool {
			return msg.Content == "Hello"
		})).Return(nil)
		messages.On("FindSender", mock.Anything, uint(1)).Return(&model.User{Name: "送信者"}, nil)

		msg, err := uc.Execute(context.Background(), 10, 1, "  Hello  ")
		require.NoError(t, err)
		assert.Equal(t, "Hello", msg.Content)
		assert.Equal(t, "送信者", msg.Sender.Name)

		assert.JSONEq(t,
			`{"type":"group_message","sender_id":1,"receiver_id":0,"room_id":10,"content":"Hello","sender_name":"送信者"}`,
			string(broadcaster.wait(t)))
	})

	t.Run("送信者が引けなくても空のユーザーを埋めて返す", func(t *testing.T) {
		rooms := new(mockChatRoomRepo)
		messages := new(mockChatRoomMessageRepo)
		broadcaster := newStubBroadcaster()
		uc := usecase.NewSendChatRoomMessageUseCase(rooms, messages, broadcaster)

		rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
		messages.On("Create", mock.Anything, mock.Anything).Return(nil)
		messages.On("FindSender", mock.Anything, uint(1)).Return(nil, nil)

		msg, err := uc.Execute(context.Background(), 10, 1, "Hello")
		require.NoError(t, err)
		require.NotNil(t, msg.Sender)
		assert.Empty(t, msg.Sender.Name)
		assert.NotContains(t, string(broadcaster.wait(t)), "sender_name")
	})

	t.Run("空白だけの本文は検証エラー", func(t *testing.T) {
		rooms := new(mockChatRoomRepo)
		messages := new(mockChatRoomMessageRepo)
		uc := usecase.NewSendChatRoomMessageUseCase(rooms, messages, newStubBroadcaster())

		_, err := uc.Execute(context.Background(), 10, 1, "   ")
		require.Error(t, err)
		assert.Equal(t, domain.ErrCodeValidation, domain.GetDomainError(err).Code)
		rooms.AssertNotCalled(t, "IsMember", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("保存に失敗したら配信しない", func(t *testing.T) {
		rooms := new(mockChatRoomRepo)
		messages := new(mockChatRoomMessageRepo)
		broadcaster := newStubBroadcaster()
		uc := usecase.NewSendChatRoomMessageUseCase(rooms, messages, broadcaster)

		rooms.On("IsMember", mock.Anything, uint(10), uint(1)).Return(true, nil)
		messages.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

		_, err := uc.Execute(context.Background(), 10, 1, "Hello")
		require.Error(t, err)
		messages.AssertNotCalled(t, "FindSender", mock.Anything, mock.Anything)
		assert.Empty(t, broadcaster.sent)
	})
}

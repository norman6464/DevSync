package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewHub(t *testing.T) {
	hub := NewHub()
	assert.NotNil(t, hub)
	assert.NotNil(t, hub.clients)
	assert.NotNil(t, hub.register)
	assert.NotNil(t, hub.unregister)
	assert.Nil(t, hub.GetRoomMembers)
}

func TestHub_RegisterAndUnregister(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		Hub:    hub,
		UserID: 1,
		Send:   make(chan []byte, 256),
	}

	// Register
	hub.Register(client)
	time.Sleep(50 * time.Millisecond)

	hub.mu.RLock()
	_, ok := hub.clients[uint(1)]
	hub.mu.RUnlock()
	assert.True(t, ok, "クライアントが登録されているべき")

	// Unregister
	hub.Unregister(client)
	time.Sleep(50 * time.Millisecond)

	hub.mu.RLock()
	_, ok = hub.clients[uint(1)]
	hub.mu.RUnlock()
	assert.False(t, ok, "クライアントが解除されているべき")
}

func TestHub_UnregisterNonExistent(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		Hub:    hub,
		UserID: 99,
		Send:   make(chan []byte, 256),
	}

	// 未登録のクライアントをUnregisterしてもpanicしない
	hub.Unregister(client)
	time.Sleep(50 * time.Millisecond)

	hub.mu.RLock()
	count := len(hub.clients)
	hub.mu.RUnlock()
	assert.Equal(t, 0, count)
}

func TestHub_SendToUser_Success(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		Hub:    hub,
		UserID: 1,
		Send:   make(chan []byte, 256),
	}

	hub.Register(client)
	time.Sleep(50 * time.Millisecond)

	msg := []byte(`{"type":"message","content":"hello"}`)
	hub.SendToUser(1, msg)

	select {
	case received := <-client.Send:
		assert.Equal(t, msg, received)
	case <-time.After(time.Second):
		t.Fatal("メッセージを受信できなかった")
	}
}

func TestHub_SendToUser_NotFound(t *testing.T) {
	hub := NewHub()

	// 存在しないユーザーに送信してもpanicしない
	hub.SendToUser(999, []byte(`{"type":"message"}`))
}

func TestHub_SendToUser_QueueFull(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// バッファサイズ1のチャネルで即座に満杯になるクライアント
	client := &Client{
		Hub:    hub,
		UserID: 1,
		Send:   make(chan []byte, 1),
	}

	hub.Register(client)
	time.Sleep(50 * time.Millisecond)

	// チャネルを満杯にする
	client.Send <- []byte("fill")

	// 2つ目のメッセージで送信キュー満杯 → unregisterが発行される
	hub.SendToUser(1, []byte("overflow"))
	time.Sleep(100 * time.Millisecond)

	hub.mu.RLock()
	_, ok := hub.clients[uint(1)]
	hub.mu.RUnlock()
	assert.False(t, ok, "キュー満杯のクライアントはunregisterされるべき")
}

func TestHub_IsRoomMember_NoCallback(t *testing.T) {
	hub := NewHub()
	assert.False(t, hub.IsRoomMember(1, 1))
}

func TestHub_IsRoomMember_Found(t *testing.T) {
	hub := NewHub()
	hub.GetRoomMembers = func(roomID uint) []uint {
		return []uint{1, 2, 3}
	}

	assert.True(t, hub.IsRoomMember(1, 2))
}

func TestHub_IsRoomMember_NotFound(t *testing.T) {
	hub := NewHub()
	hub.GetRoomMembers = func(roomID uint) []uint {
		return []uint{1, 2, 3}
	}

	assert.False(t, hub.IsRoomMember(1, 99))
}

func TestHub_SendToRoom_NoCallback(t *testing.T) {
	hub := NewHub()

	// GetRoomMembersがnilの場合、panicしない
	hub.SendToRoom(1, 1, []byte(`{"type":"group_message"}`))
}

func TestHub_SendToRoom_Success(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client2 := &Client{Hub: hub, UserID: 2, Send: make(chan []byte, 256)}
	client3 := &Client{Hub: hub, UserID: 3, Send: make(chan []byte, 256)}

	hub.Register(client2)
	hub.Register(client3)
	time.Sleep(50 * time.Millisecond)

	hub.GetRoomMembers = func(roomID uint) []uint {
		return []uint{1, 2, 3}
	}

	msg := []byte(`{"type":"group_message","room_id":1,"content":"hello room"}`)
	hub.SendToRoom(1, 1, msg) // senderID=1なので、2と3に配信

	// client2が受信
	select {
	case received := <-client2.Send:
		assert.Equal(t, msg, received)
	case <-time.After(time.Second):
		t.Fatal("client2がメッセージを受信できなかった")
	}

	// client3が受信
	select {
	case received := <-client3.Send:
		assert.Equal(t, msg, received)
	case <-time.After(time.Second):
		t.Fatal("client3がメッセージを受信できなかった")
	}
}

func TestHub_SendToRoom_SkipsSender(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	sender := &Client{Hub: hub, UserID: 1, Send: make(chan []byte, 256)}
	hub.Register(sender)
	time.Sleep(50 * time.Millisecond)

	hub.GetRoomMembers = func(roomID uint) []uint {
		return []uint{1}
	}

	hub.SendToRoom(1, 1, []byte(`{"type":"group_message"}`))

	// 送信者自身は受信しない
	select {
	case <-sender.Send:
		t.Fatal("送信者にメッセージが送られるべきではない")
	case <-time.After(100 * time.Millisecond):
		// 期待通り受信なし
	}
}

package service

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// WSMessage はWebSocketで送受信されるメッセージ構造体。
type WSMessage struct {
	Type       string `json:"type"`                 // メッセージ種別（"message", "group_message"等）
	SenderID   uint   `json:"sender_id"`            // 送信者ユーザーID
	ReceiverID uint   `json:"receiver_id"`          // 受信者ユーザーID（DM用）
	RoomID     uint   `json:"room_id,omitempty"`    // チャットルームID（グループ用）
	Content    string `json:"content"`              // メッセージ本文
	SenderName string `json:"sender_name,omitempty"` // 送信者名
}

// Client はWebSocket接続を持つ個別クライアントを表す。
type Client struct {
	Hub    *Hub            // 所属するHub
	UserID uint            // クライアントのユーザーID
	Conn   *websocket.Conn // WebSocketコネクション
	Send   chan []byte      // 送信キュー
}

// Hub はWebSocketクライアントを一元管理する中央ハブ。
// クライアントの登録・解除と、ユーザー/ルームへのメッセージ配信を担当する。
type Hub struct {
	clients        map[uint]*Client          // ユーザーID → クライアントのマップ
	register       chan *Client              // クライアント登録チャネル
	unregister     chan *Client              // クライアント解除チャネル
	mu             sync.RWMutex             // clientsマップの排他制御
	GetRoomMembers func(roomID uint) []uint // ルームメンバー取得用コールバック
}

// NewHub は新しいHubインスタンスを生成する。
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run はHubのメインイベントループを起動する。
// register/unregisterチャネルからクライアントの登録・解除を処理する。
// 通常はgoroutineとして実行される。
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
		}
	}
}

// Register はクライアントをHubに登録する。
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister はクライアントをHubから解除する。
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// SendToUser は指定ユーザーにメッセージを送信する。
// 送信キューが満杯の場合はクライアントを切断する。
func (h *Hub) SendToUser(userID uint, message []byte) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()
	if ok {
		select {
		case client.Send <- message:
		default:
			// 送信キュー満杯 → クライアント切断
			h.unregister <- client
		}
	}
}

// IsRoomMember は指定ユーザーが指定ルームのメンバーかどうかを確認する。
func (h *Hub) IsRoomMember(roomID uint, userID uint) bool {
	if h.GetRoomMembers == nil {
		return false
	}
	memberIDs := h.GetRoomMembers(roomID)
	for _, memberID := range memberIDs {
		if memberID == userID {
			return true
		}
	}
	return false
}

// SendToRoom はチャットルームの全メンバー（送信者を除く）にメッセージを配信する。
func (h *Hub) SendToRoom(roomID uint, senderID uint, message []byte) {
	if h.GetRoomMembers == nil {
		return
	}
	memberIDs := h.GetRoomMembers(roomID)
	for _, memberID := range memberIDs {
		if memberID != senderID {
			h.SendToUser(memberID, message)
		}
	}
}

// ReadPump はWebSocketからメッセージを読み取り、適切な宛先に配信する。
// グループメッセージはルーム全体に、それ以外は個別ユーザーに送信する。
// 接続切断時にHubからの登録解除とコネクションクローズを行う。
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}
		// 送信者IDをサーバー側で上書き（なりすまし防止）
		msg.SenderID = c.UserID

		data, _ := json.Marshal(msg)
		if msg.Type == "group_message" && msg.RoomID > 0 {
			// 送信者がルームのメンバーか検証（認可チェック）
			if !c.Hub.IsRoomMember(msg.RoomID, c.UserID) {
				log.Printf("websocket: ユーザー %d はルーム %d のメンバーではないため、メッセージを拒否", c.UserID, msg.RoomID)
				continue
			}
			// グループメッセージ → ルーム全体に配信
			c.Hub.SendToRoom(msg.RoomID, c.UserID, data)
		} else {
			// DM → 個別ユーザーに送信
			c.Hub.SendToUser(msg.ReceiverID, data)
		}
	}
}

// WritePump はSendチャネルからメッセージを読み取り、WebSocketに書き込む。
// チャネルがクローズされると終了する。
func (c *Client) WritePump() {
	defer c.Conn.Close()

	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("websocket write error: %v", err)
			return
		}
	}
}

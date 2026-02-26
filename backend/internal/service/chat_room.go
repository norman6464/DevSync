package service

import (
	"encoding/json"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// ChatRoomService はグループチャットのビジネスロジックを提供する。
// チャットルームのCRUD操作、メンバーシップ管理、WebSocket経由のメッセージ配信を担当する。
type ChatRoomService struct {
	roomRepo    repository.ChatRoomRepositoryInterface
	messageRepo repository.GroupMessageRepositoryInterface
	hub         *Hub
}

// NewChatRoomService は新しいChatRoomServiceインスタンスを生成する。
func NewChatRoomService(roomRepo repository.ChatRoomRepositoryInterface, messageRepo repository.GroupMessageRepositoryInterface, hub *Hub) *ChatRoomService {
	return &ChatRoomService{roomRepo: roomRepo, messageRepo: messageRepo, hub: hub}
}

// Create は新しいチャットルームを作成し、オーナーと指定メンバーを追加する。
func (s *ChatRoomService) Create(room *model.ChatRoom, memberIDs []uint) (*model.ChatRoom, error) {
	if err := domain.ValidateStringLength(room.Name, 1, 100, "チャットルーム名"); err != nil {
		return nil, err
	}
	if err := domain.ValidateStringLength(room.Description, 0, 500, "説明"); err != nil {
		return nil, err
	}
	room.Name = strings.TrimSpace(room.Name)
	if err := s.roomRepo.Create(room); err != nil {
		return nil, err
	}

	// オーナーをメンバーとして追加
	if err := s.roomRepo.AddMember(room.ID, room.OwnerID); err != nil {
		return nil, err
	}

	// その他のメンバーを追加（オーナーの重複追加を防止）
	for _, memberID := range memberIDs {
		if memberID != room.OwnerID {
			s.roomRepo.AddMember(room.ID, memberID)
		}
	}

	// アソシエーション付きで再取得
	created, err := s.roomRepo.FindByID(room.ID)
	if err != nil {
		return room, nil
	}
	return created, nil
}

// checkMembership は指定ユーザーがチャットルームのメンバーかを検証する。
func (s *ChatRoomService) checkMembership(roomID, userID uint) error {
	isMember, err := s.roomRepo.IsMember(roomID, userID)
	if err != nil || !isMember {
		return ErrForbidden
	}
	return nil
}

// GetByUserID は指定ユーザーが参加しているチャットルームをページネーション付きで取得する。
func (s *ChatRoomService) GetByUserID(userID uint, limit, offset int) ([]model.ChatRoom, int64, error) {
	return s.roomRepo.FindByUserID(userID, limit, offset)
}

// GetByID はメンバーシップを検証した後、チャットルームを取得する。
func (s *ChatRoomService) GetByID(roomID, userID uint) (*model.ChatRoom, error) {
	if err := s.checkMembership(roomID, userID); err != nil {
		return nil, err
	}
	return s.roomRepo.FindByID(roomID)
}

// findAndCheckOwnership はチャットルームを取得し、指定ユーザーがオーナーかを検証する。
func (s *ChatRoomService) findAndCheckOwnership(roomID, userID uint) (*model.ChatRoom, error) {
	return checkOwnership(s.roomRepo.FindByID, roomID, userID, func(r *model.ChatRoom) uint { return r.OwnerID })
}

// Update はオーナー権限を検証した後、チャットルーム情報を更新する。
func (s *ChatRoomService) Update(roomID, userID uint, name, description string) (*model.ChatRoom, error) {
	room, err := s.findAndCheckOwnership(roomID, userID)
	if err != nil {
		return nil, err
	}

	if name != "" {
		if err := domain.ValidateStringLength(name, 1, 100, "チャットルーム名"); err != nil {
			return nil, err
		}
		room.Name = strings.TrimSpace(name)
	}
	if description != "" {
		if err := domain.ValidateStringLength(description, 1, 500, "説明"); err != nil {
			return nil, err
		}
		room.Description = strings.TrimSpace(description)
	}

	if err := s.roomRepo.Update(room); err != nil {
		return nil, err
	}
	return room, nil
}

// Delete はオーナー権限を検証した後、チャットルームを削除する。
func (s *ChatRoomService) Delete(roomID, userID uint) error {
	if _, err := s.findAndCheckOwnership(roomID, userID); err != nil {
		return err
	}
	return s.roomRepo.Delete(roomID)
}

// GetMembers はメンバーシップを検証した後、チャットルームの全メンバーを取得する。
func (s *ChatRoomService) GetMembers(roomID, userID uint) ([]model.ChatRoomMember, error) {
	if err := s.checkMembership(roomID, userID); err != nil {
		return nil, err
	}
	return s.roomRepo.GetMembers(roomID)
}

// AddMember はリクエスト者のメンバーシップを検証した後、新しいメンバーを追加する。
// 既にメンバーのユーザーは追加できない。
func (s *ChatRoomService) AddMember(roomID, userID, targetUserID uint) error {
	isMember, err := s.roomRepo.IsMember(roomID, userID)
	if err != nil || !isMember {
		return ErrForbidden
	}
	alreadyMember, err := s.roomRepo.IsMember(roomID, targetUserID)
	if err != nil {
		return err
	}
	if alreadyMember {
		return ErrBadRequest
	}
	return s.roomRepo.AddMember(roomID, targetUserID)
}

// RemoveMember はチャットルームからメンバーを除外する。
// オーナーは誰でも除外可能、一般メンバーは自分自身のみ退出可能。
func (s *ChatRoomService) RemoveMember(roomID, userID, targetUserID uint) error {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return err
	}
	if room.OwnerID == targetUserID {
		return ErrBadRequest
	}
	if room.OwnerID != userID && userID != targetUserID {
		return ErrForbidden
	}
	return s.roomRepo.RemoveMember(roomID, targetUserID)
}

// GetMessages はメンバーシップを検証した後、チャットルームのメッセージを取得する。
func (s *ChatRoomService) GetMessages(roomID, userID uint, page, limit int) ([]model.GroupMessage, error) {
	if err := s.checkMembership(roomID, userID); err != nil {
		return nil, err
	}
	return s.messageRepo.FindByRoomID(roomID, page, limit)
}

// SendMessage はメンバーシップを検証した後、メッセージを送信する。
// WebSocket経由でルーム内の他メンバーにリアルタイム配信する。
func (s *ChatRoomService) SendMessage(roomID, userID uint, content string) (*model.GroupMessage, error) {
	if err := domain.ValidateStringLength(content, 1, 5000, "メッセージ内容"); err != nil {
		return nil, err
	}
	content = strings.TrimSpace(content)
	if err := s.checkMembership(roomID, userID); err != nil {
		return nil, err
	}

	msg := &model.GroupMessage{
		ChatRoomID: roomID,
		SenderID:   userID,
		Content:    content,
	}
	if err := s.messageRepo.Create(msg); err != nil {
		return nil, err
	}

	// 送信者情報を取得してセット
	msg.Sender = &model.User{}
	s.messageRepo.FindSenderByID(msg)

	// WebSocket経由でルーム内の他メンバーに配信
	go func() {
		wsMsg := WSMessage{
			Type:       "group_message",
			SenderID:   userID,
			RoomID:     roomID,
			Content:    content,
			SenderName: msg.Sender.Name,
		}
		data, _ := json.Marshal(wsMsg)
		s.hub.SendToRoom(roomID, userID, data)
	}()

	return msg, nil
}

// CountByUserID は指定ユーザーが参加しているチャットルーム総数を返す。
func (s *ChatRoomService) CountByUserID(userID uint) (int64, error) {
	return s.roomRepo.CountByUserID(userID)
}

// IsMember は指定ユーザーがチャットルームのメンバーかを判定する。
func (s *ChatRoomService) IsMember(roomID, userID uint) (bool, error) {
	return s.roomRepo.IsMember(roomID, userID)
}

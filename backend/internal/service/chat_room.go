package service

import (
	"encoding/json"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// ChatRoomService handles chat room business logic.
type ChatRoomService struct {
	roomRepo    repository.ChatRoomRepositoryInterface
	messageRepo repository.GroupMessageRepositoryInterface
	hub         *Hub
}

// NewChatRoomService creates a new ChatRoomService.
func NewChatRoomService(roomRepo repository.ChatRoomRepositoryInterface, messageRepo repository.GroupMessageRepositoryInterface, hub *Hub) *ChatRoomService {
	return &ChatRoomService{roomRepo: roomRepo, messageRepo: messageRepo, hub: hub}
}

// Create creates a new chat room and adds the owner and specified members.
func (s *ChatRoomService) Create(room *model.ChatRoom, memberIDs []uint) (*model.ChatRoom, error) {
	if err := s.roomRepo.Create(room); err != nil {
		return nil, err
	}

	// Add owner as member
	if err := s.roomRepo.AddMember(room.ID, room.OwnerID); err != nil {
		return nil, err
	}

	// Add other members
	for _, memberID := range memberIDs {
		if memberID != room.OwnerID {
			s.roomRepo.AddMember(room.ID, memberID)
		}
	}

	// Reload with associations
	created, err := s.roomRepo.FindByID(room.ID)
	if err != nil {
		return room, nil
	}
	return created, nil
}

// GetByUserID returns all chat rooms for a user.
func (s *ChatRoomService) GetByUserID(userID uint) ([]model.ChatRoom, error) {
	return s.roomRepo.FindByUserID(userID)
}

// GetByID returns a chat room by ID after verifying membership.
func (s *ChatRoomService) GetByID(roomID, userID uint) (*model.ChatRoom, error) {
	isMember, err := s.roomRepo.IsMember(roomID, userID)
	if err != nil || !isMember {
		return nil, ErrForbidden
	}

	return s.roomRepo.FindByID(roomID)
}

// Update updates a chat room after verifying ownership.
func (s *ChatRoomService) Update(roomID, userID uint, name, description string) (*model.ChatRoom, error) {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return nil, err
	}
	if room.OwnerID != userID {
		return nil, ErrForbidden
	}

	if name != "" {
		room.Name = name
	}
	room.Description = description

	if err := s.roomRepo.Update(room); err != nil {
		return nil, err
	}
	return room, nil
}

// Delete deletes a chat room after verifying ownership.
func (s *ChatRoomService) Delete(roomID, userID uint) error {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return err
	}
	if room.OwnerID != userID {
		return ErrForbidden
	}
	return s.roomRepo.Delete(roomID)
}

// GetMembers returns all members of a chat room after verifying membership.
func (s *ChatRoomService) GetMembers(roomID, userID uint) ([]model.ChatRoomMember, error) {
	isMember, err := s.roomRepo.IsMember(roomID, userID)
	if err != nil || !isMember {
		return nil, ErrForbidden
	}
	return s.roomRepo.GetMembers(roomID)
}

// AddMember adds a member to a chat room after verifying the requester is a member.
func (s *ChatRoomService) AddMember(roomID, userID, targetUserID uint) error {
	isMember, err := s.roomRepo.IsMember(roomID, userID)
	if err != nil || !isMember {
		return ErrForbidden
	}
	return s.roomRepo.AddMember(roomID, targetUserID)
}

// RemoveMember removes a member from a chat room.
// Owner can remove anyone, others can only remove themselves.
func (s *ChatRoomService) RemoveMember(roomID, userID, targetUserID uint) error {
	room, err := s.roomRepo.FindByID(roomID)
	if err != nil {
		return err
	}
	if room.OwnerID != userID && userID != targetUserID {
		return ErrForbidden
	}
	return s.roomRepo.RemoveMember(roomID, targetUserID)
}

// GetMessages returns paginated messages for a chat room after verifying membership.
func (s *ChatRoomService) GetMessages(roomID, userID uint, page, limit int) ([]model.GroupMessage, error) {
	isMember, err := s.roomRepo.IsMember(roomID, userID)
	if err != nil || !isMember {
		return nil, ErrForbidden
	}
	return s.messageRepo.FindByRoomID(roomID, page, limit)
}

// SendMessage sends a message to a chat room after verifying membership.
func (s *ChatRoomService) SendMessage(roomID, userID uint, content string) (*model.GroupMessage, error) {
	isMember, err := s.roomRepo.IsMember(roomID, userID)
	if err != nil || !isMember {
		return nil, ErrForbidden
	}

	msg := &model.GroupMessage{
		ChatRoomID: roomID,
		SenderID:   userID,
		Content:    content,
	}
	if err := s.messageRepo.Create(msg); err != nil {
		return nil, err
	}

	// Reload with sender info
	msg.Sender = &model.User{}
	s.messageRepo.FindSenderByID(msg)

	// Send via WebSocket to all room members
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

// IsMember checks if a user is a member of a chat room.
func (s *ChatRoomService) IsMember(roomID, userID uint) (bool, error) {
	return s.roomRepo.IsMember(roomID, userID)
}

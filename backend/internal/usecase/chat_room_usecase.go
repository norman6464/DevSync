package usecase

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// チャットルームの入力値の上限。
const (
	chatRoomNameMaxLen        = 100
	chatRoomDescriptionMaxLen = 500
	chatRoomMessageMaxLen     = 5000
)

// chatRoomBroadcastType はルーム内配信であることを表すメッセージ種別。
const chatRoomBroadcastType = "group_message"

// chatRoomBroadcast はルームへリアルタイム配信するメッセージの表現。
// WebSocket クライアントが受け取る JSON を変えないため、フィールドと並びを配信形式に合わせている。
type chatRoomBroadcast struct {
	Type       string `json:"type"`
	SenderID   uint   `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	RoomID     uint   `json:"room_id,omitempty"`
	Content    string `json:"content"`
	SenderName string `json:"sender_name,omitempty"`
}

// chatRoomOwnerOf は所有権チェック用にルームの所有者 ID を取り出す。
func chatRoomOwnerOf(r *model.ChatRoom) uint { return r.OwnerID }

// ensureChatRoomMember はユーザーがルームのメンバーであることを検証する。
// 判定に失敗した場合もメンバーでない場合も 403 を返す（判定エラーを漏らさないため）。
func ensureChatRoomMember(ctx context.Context, rooms repository.ChatRoomRepository, roomID, userID uint) error {
	isMember, err := rooms.IsMember(ctx, roomID, userID)
	if err != nil || !isMember {
		return domain.ErrForbidden
	}
	return nil
}

// CreateChatRoomInput はチャットルーム作成の入力。
type CreateChatRoomInput struct {
	Name        string
	Description string
	OwnerID     uint
	MemberIDs   []uint
}

// CreateChatRoomUseCase はチャットルームを作成し、オーナーと指定メンバーを参加させる。
type CreateChatRoomUseCase struct {
	rooms repository.ChatRoomRepository
}

// NewCreateChatRoomUseCase は CreateChatRoomUseCase を生成する。
func NewCreateChatRoomUseCase(rooms repository.ChatRoomRepository) *CreateChatRoomUseCase {
	return &CreateChatRoomUseCase{rooms: rooms}
}

// Execute はルーム名と説明を検証してルームを作成し、オーナーと指定メンバーを追加する。
// メンバー追加に失敗しても作成自体は成功として扱い、再取得できなければ作成直後のルームを返す。
func (uc *CreateChatRoomUseCase) Execute(ctx context.Context, in CreateChatRoomInput) (*model.ChatRoom, error) {
	if err := domain.ValidateStringLength(in.Name, 1, chatRoomNameMaxLen, "チャットルーム名"); err != nil {
		return nil, err
	}
	if err := domain.ValidateStringLength(in.Description, 0, chatRoomDescriptionMaxLen, "説明"); err != nil {
		return nil, err
	}

	room := &model.ChatRoom{
		Name:        strings.TrimSpace(in.Name),
		Description: strings.TrimSpace(in.Description),
		OwnerID:     in.OwnerID,
	}
	if err := uc.rooms.Create(ctx, room); err != nil {
		return nil, err
	}
	if err := uc.rooms.AddMember(ctx, room.ID, room.OwnerID); err != nil {
		return nil, err
	}

	// オーナーの重複追加を避けつつメンバーを追加する。個々の失敗は作成結果に影響させない。
	for _, memberID := range in.MemberIDs {
		if memberID != room.OwnerID {
			_ = uc.rooms.AddMember(ctx, room.ID, memberID)
		}
	}

	created, err := uc.rooms.FindByID(ctx, room.ID)
	if err != nil || created == nil {
		return room, nil
	}
	return created, nil
}

// ListMyChatRoomsUseCase は指定ユーザーが参加しているチャットルームを一覧する。
type ListMyChatRoomsUseCase struct {
	rooms repository.ChatRoomRepository
}

// NewListMyChatRoomsUseCase は ListMyChatRoomsUseCase を生成する。
func NewListMyChatRoomsUseCase(rooms repository.ChatRoomRepository) *ListMyChatRoomsUseCase {
	return &ListMyChatRoomsUseCase{rooms: rooms}
}

// Execute は参加中のルームを更新日時の降順で返す。
func (uc *ListMyChatRoomsUseCase) Execute(ctx context.Context, userID uint, limit, offset int) ([]model.ChatRoom, int64, error) {
	return uc.rooms.FindByUserID(ctx, userID, limit, offset)
}

// GetChatRoomUseCase はメンバー本人にチャットルームを 1 件返す。
type GetChatRoomUseCase struct {
	rooms repository.ChatRoomRepository
}

// NewGetChatRoomUseCase は GetChatRoomUseCase を生成する。
func NewGetChatRoomUseCase(rooms repository.ChatRoomRepository) *GetChatRoomUseCase {
	return &GetChatRoomUseCase{rooms: rooms}
}

// Execute はメンバーシップを検証したうえでルームを返す。
func (uc *GetChatRoomUseCase) Execute(ctx context.Context, roomID, userID uint) (*model.ChatRoom, error) {
	if err := ensureChatRoomMember(ctx, uc.rooms, roomID, userID); err != nil {
		return nil, err
	}
	room, err := uc.rooms.FindByID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, domain.ErrNotFound
	}
	return room, nil
}

// UpdateChatRoomUseCase はオーナー本人のチャットルーム情報を更新する。
type UpdateChatRoomUseCase struct {
	rooms repository.ChatRoomRepository
}

// NewUpdateChatRoomUseCase は UpdateChatRoomUseCase を生成する。
func NewUpdateChatRoomUseCase(rooms repository.ChatRoomRepository) *UpdateChatRoomUseCase {
	return &UpdateChatRoomUseCase{rooms: rooms}
}

// Execute は所有権を検証し、空でない項目だけを検証して更新する。
func (uc *UpdateChatRoomUseCase) Execute(ctx context.Context, roomID, userID uint, name, description string) (*model.ChatRoom, error) {
	room, err := ensureOwner(ctx, uc.rooms.FindByID, roomID, userID, chatRoomOwnerOf)
	if err != nil {
		return nil, err
	}

	if name != "" {
		if err := domain.ValidateStringLength(name, 1, chatRoomNameMaxLen, "チャットルーム名"); err != nil {
			return nil, err
		}
		room.Name = strings.TrimSpace(name)
	}
	if description != "" {
		if err := domain.ValidateStringLength(description, 1, chatRoomDescriptionMaxLen, "説明"); err != nil {
			return nil, err
		}
		room.Description = strings.TrimSpace(description)
	}

	if err := uc.rooms.Update(ctx, room); err != nil {
		return nil, err
	}
	return room, nil
}

// DeleteChatRoomUseCase はオーナー本人のチャットルームを削除する。
type DeleteChatRoomUseCase struct {
	rooms repository.ChatRoomRepository
}

// NewDeleteChatRoomUseCase は DeleteChatRoomUseCase を生成する。
func NewDeleteChatRoomUseCase(rooms repository.ChatRoomRepository) *DeleteChatRoomUseCase {
	return &DeleteChatRoomUseCase{rooms: rooms}
}

// Execute は所有権を検証したうえでルームをメッセージ・メンバーごと削除する。
func (uc *DeleteChatRoomUseCase) Execute(ctx context.Context, roomID, userID uint) error {
	if _, err := ensureOwner(ctx, uc.rooms.FindByID, roomID, userID, chatRoomOwnerOf); err != nil {
		return err
	}
	return uc.rooms.Delete(ctx, roomID)
}

// ListChatRoomMembersUseCase はメンバー本人にチャットルームのメンバー一覧を返す。
type ListChatRoomMembersUseCase struct {
	rooms repository.ChatRoomRepository
}

// NewListChatRoomMembersUseCase は ListChatRoomMembersUseCase を生成する。
func NewListChatRoomMembersUseCase(rooms repository.ChatRoomRepository) *ListChatRoomMembersUseCase {
	return &ListChatRoomMembersUseCase{rooms: rooms}
}

// Execute はメンバーシップを検証したうえでメンバー一覧を返す。
func (uc *ListChatRoomMembersUseCase) Execute(ctx context.Context, roomID, userID uint) ([]model.ChatRoomMember, error) {
	if err := ensureChatRoomMember(ctx, uc.rooms, roomID, userID); err != nil {
		return nil, err
	}
	return uc.rooms.GetMembers(ctx, roomID)
}

// AddChatRoomMemberUseCase はチャットルームにメンバーを追加する。
type AddChatRoomMemberUseCase struct {
	rooms repository.ChatRoomRepository
}

// NewAddChatRoomMemberUseCase は AddChatRoomMemberUseCase を生成する。
func NewAddChatRoomMemberUseCase(rooms repository.ChatRoomRepository) *AddChatRoomMemberUseCase {
	return &AddChatRoomMemberUseCase{rooms: rooms}
}

// Execute はリクエスト者がメンバーであることを検証し、まだ参加していないユーザーを追加する。
func (uc *AddChatRoomMemberUseCase) Execute(ctx context.Context, roomID, userID, targetUserID uint) error {
	if err := ensureChatRoomMember(ctx, uc.rooms, roomID, userID); err != nil {
		return err
	}
	alreadyMember, err := uc.rooms.IsMember(ctx, roomID, targetUserID)
	if err != nil {
		return err
	}
	if alreadyMember {
		return domain.ErrBadRequest
	}
	return uc.rooms.AddMember(ctx, roomID, targetUserID)
}

// RemoveChatRoomMemberUseCase はチャットルームからメンバーを除外する。
type RemoveChatRoomMemberUseCase struct {
	rooms repository.ChatRoomRepository
}

// NewRemoveChatRoomMemberUseCase は RemoveChatRoomMemberUseCase を生成する。
func NewRemoveChatRoomMemberUseCase(rooms repository.ChatRoomRepository) *RemoveChatRoomMemberUseCase {
	return &RemoveChatRoomMemberUseCase{rooms: rooms}
}

// Execute はオーナーなら任意のメンバーを、一般メンバーなら自分自身だけを除外する。
// オーナー自身は除外できない。
func (uc *RemoveChatRoomMemberUseCase) Execute(ctx context.Context, roomID, userID, targetUserID uint) error {
	room, err := uc.rooms.FindByID(ctx, roomID)
	if err != nil {
		return err
	}
	if room == nil {
		return domain.ErrNotFound
	}
	if room.OwnerID == targetUserID {
		return domain.ErrBadRequest
	}
	if room.OwnerID != userID && userID != targetUserID {
		return domain.ErrForbidden
	}
	return uc.rooms.RemoveMember(ctx, roomID, targetUserID)
}

// ListChatRoomMessagesUseCase はメンバー本人にチャットルームのメッセージを返す。
type ListChatRoomMessagesUseCase struct {
	rooms    repository.ChatRoomRepository
	messages repository.ChatRoomMessageRepository
}

// NewListChatRoomMessagesUseCase は ListChatRoomMessagesUseCase を生成する。
func NewListChatRoomMessagesUseCase(
	rooms repository.ChatRoomRepository,
	messages repository.ChatRoomMessageRepository,
) *ListChatRoomMessagesUseCase {
	return &ListChatRoomMessagesUseCase{rooms: rooms, messages: messages}
}

// Execute はメンバーシップを検証したうえでメッセージを古い順に返す。
func (uc *ListChatRoomMessagesUseCase) Execute(ctx context.Context, roomID, userID uint, page, limit int) ([]model.GroupMessage, error) {
	if err := ensureChatRoomMember(ctx, uc.rooms, roomID, userID); err != nil {
		return nil, err
	}
	return uc.messages.FindByRoomID(ctx, roomID, page, limit)
}

// SendChatRoomMessageUseCase はチャットルームへメッセージを送信し、参加者へ配信する。
type SendChatRoomMessageUseCase struct {
	rooms       repository.ChatRoomRepository
	messages    repository.ChatRoomMessageRepository
	broadcaster repository.RoomBroadcaster
}

// NewSendChatRoomMessageUseCase は SendChatRoomMessageUseCase を生成する。
func NewSendChatRoomMessageUseCase(
	rooms repository.ChatRoomRepository,
	messages repository.ChatRoomMessageRepository,
	broadcaster repository.RoomBroadcaster,
) *SendChatRoomMessageUseCase {
	return &SendChatRoomMessageUseCase{rooms: rooms, messages: messages, broadcaster: broadcaster}
}

// Execute は本文を検証してメッセージを保存し、送信者を除く参加者へ非同期に配信する。
// 送信者情報が引けなかった場合も保存済みのメッセージは返す（配信名だけが空になる）。
func (uc *SendChatRoomMessageUseCase) Execute(ctx context.Context, roomID, userID uint, content string) (*model.GroupMessage, error) {
	if err := domain.ValidateStringLength(content, 1, chatRoomMessageMaxLen, "メッセージ内容"); err != nil {
		return nil, err
	}
	content = strings.TrimSpace(content)
	if err := ensureChatRoomMember(ctx, uc.rooms, roomID, userID); err != nil {
		return nil, err
	}

	msg := &model.GroupMessage{
		ChatRoomID: roomID,
		SenderID:   userID,
		Content:    content,
	}
	if err := uc.messages.Create(ctx, msg); err != nil {
		return nil, err
	}

	msg.Sender = &model.User{}
	if sender, err := uc.messages.FindSender(ctx, userID); err == nil && sender != nil {
		msg.Sender = sender
	}

	go uc.broadcast(roomID, userID, content, msg.Sender.Name)

	return msg, nil
}

// broadcast は送信者を除くルーム参加者へメッセージを配信する。
func (uc *SendChatRoomMessageUseCase) broadcast(roomID, senderID uint, content, senderName string) {
	data, err := json.Marshal(chatRoomBroadcast{
		Type:       chatRoomBroadcastType,
		SenderID:   senderID,
		RoomID:     roomID,
		Content:    content,
		SenderName: senderName,
	})
	if err != nil {
		return
	}
	uc.broadcaster.SendToRoom(roomID, senderID, data)
}

// CountMyChatRoomsUseCase は指定ユーザーが参加しているチャットルーム総数を返す。
type CountMyChatRoomsUseCase struct {
	rooms repository.ChatRoomRepository
}

// NewCountMyChatRoomsUseCase は CountMyChatRoomsUseCase を生成する。
func NewCountMyChatRoomsUseCase(rooms repository.ChatRoomRepository) *CountMyChatRoomsUseCase {
	return &CountMyChatRoomsUseCase{rooms: rooms}
}

// Execute は参加しているルームの総数を返す。
func (uc *CountMyChatRoomsUseCase) Execute(ctx context.Context, userID uint) (int64, error) {
	return uc.rooms.CountByUserID(ctx, userID)
}

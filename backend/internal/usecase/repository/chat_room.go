package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// ChatRoomRepository はグループチャットルームの永続化に対する、usecase 側が要求する契約。
// ルーム本体の CRUD とメンバーシップ操作をまとめて扱う。
type ChatRoomRepository interface {
	Create(ctx context.Context, room *model.ChatRoom) error
	// FindByID は指定 ID のルームをオーナー情報付きで返す。
	// 不在の場合は「不在」を表す (nil, nil) を返し、DB 障害だけを error として返す。
	FindByID(ctx context.Context, id uint) (*model.ChatRoom, error)
	// FindByUserID は指定ユーザーが参加しているルームを更新日時の降順で返す。
	FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.ChatRoom, int64, error)
	Update(ctx context.Context, room *model.ChatRoom) error
	// Delete はルームをメッセージ・メンバーごと削除する。
	Delete(ctx context.Context, roomID uint) error

	AddMember(ctx context.Context, roomID, userID uint) error
	RemoveMember(ctx context.Context, roomID, userID uint) error
	GetMembers(ctx context.Context, roomID uint) ([]model.ChatRoomMember, error)
	IsMember(ctx context.Context, roomID, userID uint) (bool, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
}

// ChatRoomMessageRepository はグループチャットのメッセージ永続化に対する契約。
type ChatRoomMessageRepository interface {
	Create(ctx context.Context, msg *model.GroupMessage) error
	// FindByRoomID は指定ルームのメッセージを作成日時の昇順でページネーションして返す。
	FindByRoomID(ctx context.Context, roomID uint, page, limit int) ([]model.GroupMessage, error)
	// FindSender は送信者のユーザー情報を返す。不在の場合は (nil, nil) を返す。
	FindSender(ctx context.Context, senderID uint) (*model.User, error)
}

// RoomBroadcaster はチャットルームの参加者へリアルタイム配信するための最小の契約。
// 送信者自身を除いたメンバーへ message をそのまま届ける。
type RoomBroadcaster interface {
	SendToRoom(roomID, senderID uint, message []byte)
}

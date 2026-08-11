package persistence

import (
	"context"
	"errors"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// chatRoomMessageRepository は [repository.ChatRoomMessageRepository] の GORM 実装。
type chatRoomMessageRepository struct {
	db *gorm.DB
}

// NewChatRoomMessageRepository は ChatRoomMessageRepository の GORM 実装を返す。
func NewChatRoomMessageRepository(db *gorm.DB) repository.ChatRoomMessageRepository {
	return &chatRoomMessageRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ChatRoomMessageRepository = (*chatRoomMessageRepository)(nil)

// Create はグループメッセージを保存する。
func (r *chatRoomMessageRepository) Create(ctx context.Context, msg *model.GroupMessage) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

// FindByRoomID は指定ルームのメッセージを送信者情報付きで、作成日時の昇順に取得する。
func (r *chatRoomMessageRepository) FindByRoomID(ctx context.Context, roomID uint, page, limit int) ([]model.GroupMessage, error) {
	var messages []model.GroupMessage
	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).Preload("Sender").
		Where("chat_room_id = ?", roomID).
		Order("created_at ASC").
		Offset(offset).Limit(limit).
		Find(&messages).Error
	return messages, err
}

// FindSender は送信者のユーザー情報を取得する。不在の場合は (nil, nil) を返す。
func (r *chatRoomMessageRepository) FindSender(ctx context.Context, senderID uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, senderID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

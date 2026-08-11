package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// chatRoomRepository は [repository.ChatRoomRepository] の GORM 実装。
type chatRoomRepository struct {
	db *gorm.DB
}

// NewChatRoomRepository は ChatRoomRepository の GORM 実装を返す。
func NewChatRoomRepository(db *gorm.DB) repository.ChatRoomRepository {
	return &chatRoomRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.ChatRoomRepository = (*chatRoomRepository)(nil)

// Create はチャットルームを作成する。
func (r *chatRoomRepository) Create(ctx context.Context, room *model.ChatRoom) error {
	return r.db.WithContext(ctx).Create(room).Error
}

// FindByID は指定 ID のチャットルームをオーナー情報付きで取得する。不在の場合は (nil, nil) を返す。
func (r *chatRoomRepository) FindByID(ctx context.Context, id uint) (*model.ChatRoom, error) {
	var room model.ChatRoom
	err := r.db.WithContext(ctx).Preload("Owner").First(&room, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// FindByUserID は指定ユーザーが参加しているチャットルームを更新日時の降順で取得する。
func (r *chatRoomRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]model.ChatRoom, int64, error) {
	base := r.db.WithContext(ctx).
		Joins("JOIN chat_room_members ON chat_room_members.chat_room_id = chat_rooms.id").
		Where("chat_room_members.user_id = ?", userID).
		Session(&gorm.Session{})

	// 総件数の取得エラーは移行前と同じく無視し、一覧取得のエラーだけを返す。
	var total int64
	base.Model(&model.ChatRoom{}).Count(&total)

	var rooms []model.ChatRoom
	err := base.Preload("Owner").
		Order("chat_rooms.updated_at DESC").
		Limit(limit).Offset(offset).
		Find(&rooms).Error
	return rooms, total, err
}

// Update はチャットルーム情報を更新する。
func (r *chatRoomRepository) Update(ctx context.Context, room *model.ChatRoom) error {
	return r.db.WithContext(ctx).Save(room).Error
}

// Delete はチャットルームをメッセージ・メンバーごとトランザクション内で削除する。
func (r *chatRoomRepository) Delete(ctx context.Context, roomID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("chat_room_id = ?", roomID).Delete(&model.GroupMessage{}).Error; err != nil {
			return err
		}
		if err := tx.Where("chat_room_id = ?", roomID).Delete(&model.ChatRoomMember{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.ChatRoom{}, roomID).Error
	})
}

// AddMember はチャットルームにメンバーを追加する。
func (r *chatRoomRepository) AddMember(ctx context.Context, roomID, userID uint) error {
	member := model.ChatRoomMember{
		ChatRoomID: roomID,
		UserID:     userID,
		JoinedAt:   time.Now(),
	}
	return r.db.WithContext(ctx).Create(&member).Error
}

// RemoveMember はチャットルームからメンバーを除外する。
func (r *chatRoomRepository) RemoveMember(ctx context.Context, roomID, userID uint) error {
	return r.db.WithContext(ctx).
		Where("chat_room_id = ? AND user_id = ?", roomID, userID).
		Delete(&model.ChatRoomMember{}).Error
}

// GetMembers はチャットルームの全メンバーをユーザー情報付きで取得する。
func (r *chatRoomRepository) GetMembers(ctx context.Context, roomID uint) ([]model.ChatRoomMember, error) {
	var members []model.ChatRoomMember
	err := r.db.WithContext(ctx).Preload("User").
		Where("chat_room_id = ?", roomID).
		Find(&members).Error
	return members, err
}

// IsMember は指定ユーザーがチャットルームのメンバーかを判定する。
func (r *chatRoomRepository) IsMember(ctx context.Context, roomID, userID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ChatRoomMember{}).
		Where("chat_room_id = ? AND user_id = ?", roomID, userID).
		Count(&count).Error
	return count > 0, err
}

// CountByUserID は指定ユーザーが参加しているチャットルーム総数を返す。
func (r *chatRoomRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ChatRoomMember{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

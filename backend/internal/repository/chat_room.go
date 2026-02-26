package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// ChatRoomRepository はグループチャットルームデータへのアクセスを提供するリポジトリ実装。
type ChatRoomRepository struct {
	db *gorm.DB
}

// NewChatRoomRepository は新しいChatRoomRepositoryインスタンスを生成する。
func NewChatRoomRepository(db *gorm.DB) *ChatRoomRepository {
	return &ChatRoomRepository{db: db}
}

// Create は新しいチャットルームをデータベースに作成する。
func (r *ChatRoomRepository) Create(room *model.ChatRoom) error {
	return r.db.Create(room).Error
}

// FindByID は指定IDのチャットルームをオーナー情報付きで取得する。
func (r *ChatRoomRepository) FindByID(id uint) (*model.ChatRoom, error) {
	var room model.ChatRoom
	err := r.db.Preload("Owner").First(&room, id).Error
	return &room, err
}

// FindByUserID は指定ユーザーが参加しているチャットルームをページネーション付きで取得する。
// chat_room_membersテーブルをJOINし、更新日時の降順でソートされる。
func (r *ChatRoomRepository) FindByUserID(userID uint, limit, offset int) ([]model.ChatRoom, int64, error) {
	var rooms []model.ChatRoom
	var total int64
	query := r.db.Joins("JOIN chat_room_members ON chat_room_members.chat_room_id = chat_rooms.id").
		Where("chat_room_members.user_id = ?", userID)
	query.Model(&model.ChatRoom{}).Count(&total)
	err := query.Preload("Owner").
		Order("chat_rooms.updated_at DESC").
		Limit(limit).Offset(offset).
		Find(&rooms).Error
	return rooms, total, err
}

// Update は既存のチャットルーム情報を更新する。
func (r *ChatRoomRepository) Update(room *model.ChatRoom) error {
	return r.db.Save(room).Error
}

// Delete はチャットルームを関連データ（メッセージ・メンバー）ごとトランザクション内で削除する。
func (r *ChatRoomRepository) Delete(roomID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// ルーム内のメッセージを先に削除
		if err := tx.Where("chat_room_id = ?", roomID).Delete(&model.GroupMessage{}).Error; err != nil {
			return err
		}
		// メンバーシップレコードを削除
		if err := tx.Where("chat_room_id = ?", roomID).Delete(&model.ChatRoomMember{}).Error; err != nil {
			return err
		}
		// チャットルーム本体を削除
		return tx.Delete(&model.ChatRoom{}, roomID).Error
	})
}

// AddMember はチャットルームに新しいメンバーを追加する。
func (r *ChatRoomRepository) AddMember(roomID, userID uint) error {
	member := model.ChatRoomMember{
		ChatRoomID: roomID,
		UserID:     userID,
		JoinedAt:   time.Now(),
	}
	return r.db.Create(&member).Error
}

// RemoveMember はチャットルームからメンバーを除外する。
func (r *ChatRoomRepository) RemoveMember(roomID, userID uint) error {
	return r.db.Where("chat_room_id = ? AND user_id = ?", roomID, userID).
		Delete(&model.ChatRoomMember{}).Error
}

// GetMembers は指定チャットルームの全メンバーをユーザー情報付きで取得する。
func (r *ChatRoomRepository) GetMembers(roomID uint) ([]model.ChatRoomMember, error) {
	var members []model.ChatRoomMember
	err := r.db.Preload("User").Where("chat_room_id = ?", roomID).Find(&members).Error
	return members, err
}

// CountByUserID は指定ユーザーが参加しているチャットルーム総数を返す。
func (r *ChatRoomRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.ChatRoomMember{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// IsMember は指定ユーザーがチャットルームのメンバーであるかを判定する。
func (r *ChatRoomRepository) IsMember(roomID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.ChatRoomMember{}).
		Where("chat_room_id = ? AND user_id = ?", roomID, userID).
		Count(&count).Error
	return count > 0, err
}

package repository

import (
	"time"

	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// NoteStatsRepository はノート集計統計の取得を担当するリポジトリ実装。
type NoteStatsRepository struct {
	db *gorm.DB
}

// NewNoteStatsRepository は新しいNoteStatsRepositoryインスタンスを生成する。
func NewNoteStatsRepository(db *gorm.DB) *NoteStatsRepository {
	return &NoteStatsRepository{db: db}
}

// GetNoteStats は指定ユーザーのノート集計統計を返す。
func (r *NoteStatsRepository) GetNoteStats(userID uint) (*model.NoteStats, error) {
	var stats model.NoteStats

	// 総ノート数
	if err := r.db.Model(&model.Note{}).Where("user_id = ?", userID).Count(&stats.TotalNotes).Error; err != nil {
		return nil, err
	}

	// アーカイブ済みノート数
	if err := r.db.Model(&model.Note{}).Where("user_id = ? AND is_archived = ?", userID, true).Count(&stats.ArchivedNotes).Error; err != nil {
		return nil, err
	}

	// お気に入りノート数
	if err := r.db.Model(&model.Note{}).Where("user_id = ? AND is_favorite = ?", userID, true).Count(&stats.FavoriteNotes).Error; err != nil {
		return nil, err
	}

	// フォルダ数
	if err := r.db.Model(&model.NoteFolder{}).Where("user_id = ?", userID).Count(&stats.TotalFolders).Error; err != nil {
		return nil, err
	}

	// 今週作成したノート数（過去7日間）
	weekAgo := time.Now().AddDate(0, 0, -7)
	if err := r.db.Model(&model.Note{}).Where("user_id = ? AND created_at >= ?", userID, weekAgo).Count(&stats.NotesThisWeek).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

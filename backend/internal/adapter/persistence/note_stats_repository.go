package persistence

import (
	"context"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// noteStatsRepository は [repository.NoteStatsRepository] の GORM 実装。
type noteStatsRepository struct {
	db *gorm.DB
}

// NewNoteStatsRepository は NoteStatsRepository の GORM 実装を返す。
func NewNoteStatsRepository(db *gorm.DB) repository.NoteStatsRepository {
	return &noteStatsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.NoteStatsRepository = (*noteStatsRepository)(nil)

// GetNoteStats は指定ユーザーのノート集計統計を返す。
func (r *noteStatsRepository) GetNoteStats(ctx context.Context, userID uint) (*model.NoteStats, error) {
	db := r.db.WithContext(ctx)
	var stats model.NoteStats

	// 総ノート数
	if err := db.Model(&model.Note{}).Where("user_id = ?", userID).Count(&stats.TotalNotes).Error; err != nil {
		return nil, err
	}

	// アーカイブ済みノート数
	if err := db.Model(&model.Note{}).Where("user_id = ? AND is_archived = ?", userID, true).Count(&stats.ArchivedNotes).Error; err != nil {
		return nil, err
	}

	// お気に入りノート数
	if err := db.Model(&model.Note{}).Where("user_id = ? AND is_favorite = ?", userID, true).Count(&stats.FavoriteNotes).Error; err != nil {
		return nil, err
	}

	// フォルダ数
	if err := db.Model(&model.NoteFolder{}).Where("user_id = ?", userID).Count(&stats.TotalFolders).Error; err != nil {
		return nil, err
	}

	// 今週作成したノート数（過去7日間）
	weekAgo := domain.DaysAgo(time.Now(), 7)
	if err := db.Model(&model.Note{}).Where("user_id = ? AND created_at >= ?", userID, weekAgo).Count(&stats.NotesThisWeek).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

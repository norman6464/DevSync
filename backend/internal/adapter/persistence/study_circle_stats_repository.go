package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// studyCircleStatsRepository は [repository.StudyCircleStatsRepository] の GORM 実装。
type studyCircleStatsRepository struct {
	db *gorm.DB
}

// NewStudyCircleStatsRepository は StudyCircleStatsRepository の GORM 実装を返す。
func NewStudyCircleStatsRepository(db *gorm.DB) repository.StudyCircleStatsRepository {
	return &studyCircleStatsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.StudyCircleStatsRepository = (*studyCircleStatsRepository)(nil)

// GetCircleStats は指定サークルの集計統計を返す。
func (r *studyCircleStatsRepository) GetCircleStats(ctx context.Context, circleID uint) (*model.StudyCircleStats, error) {
	db := r.db.WithContext(ctx)
	var stats model.StudyCircleStats

	// メンバー数
	if err := db.Model(&model.StudyCircleMember{}).Where("circle_id = ?", circleID).Count(&stats.MemberCount).Error; err != nil {
		return nil, err
	}

	// チェックイン数
	if err := db.Model(&model.StudyCircleCheckin{}).Where("circle_id = ?", circleID).Count(&stats.CheckinCount).Error; err != nil {
		return nil, err
	}

	// ステップ総数
	if err := db.Model(&model.StudyCircleStep{}).Where("circle_id = ?", circleID).Count(&stats.TotalSteps).Error; err != nil {
		return nil, err
	}

	// 完了済みステップ数（メンバー別進捗の完了エントリ数）
	if err := db.Model(&model.StudyCircleMemberProgress{}).Where("circle_id = ? AND is_completed = ?", circleID, true).Count(&stats.CompletedSteps).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

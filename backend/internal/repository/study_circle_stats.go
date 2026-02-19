package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// StudyCircleStatsRepository はスタディサークル集計統計の取得を担当するリポジトリ実装。
type StudyCircleStatsRepository struct {
	db *gorm.DB
}

// NewStudyCircleStatsRepository は新しいStudyCircleStatsRepositoryインスタンスを生成する。
func NewStudyCircleStatsRepository(db *gorm.DB) *StudyCircleStatsRepository {
	return &StudyCircleStatsRepository{db: db}
}

// GetCircleStats は指定サークルの集計統計を返す。
func (r *StudyCircleStatsRepository) GetCircleStats(circleID uint) (*model.StudyCircleStats, error) {
	var stats model.StudyCircleStats

	// メンバー数
	if err := r.db.Model(&model.StudyCircleMember{}).Where("circle_id = ?", circleID).Count(&stats.MemberCount).Error; err != nil {
		return nil, err
	}

	// チェックイン数
	if err := r.db.Model(&model.StudyCircleCheckin{}).Where("circle_id = ?", circleID).Count(&stats.CheckinCount).Error; err != nil {
		return nil, err
	}

	// ステップ総数
	if err := r.db.Model(&model.StudyCircleStep{}).Where("circle_id = ?", circleID).Count(&stats.TotalSteps).Error; err != nil {
		return nil, err
	}

	// 完了済みステップ数（メンバー別進捗の完了エントリ数）
	if err := r.db.Model(&model.StudyCircleMemberProgress{}).Where("circle_id = ? AND is_completed = ?", circleID, true).Count(&stats.CompletedSteps).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

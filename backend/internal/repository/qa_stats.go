package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// QAStatsRepository はユーザーQ&A活動集計統計の取得を担当するリポジトリ実装。
type QAStatsRepository struct {
	db *gorm.DB
}

// NewQAStatsRepository は新しいQAStatsRepositoryインスタンスを生成する。
func NewQAStatsRepository(db *gorm.DB) *QAStatsRepository {
	return &QAStatsRepository{db: db}
}

// GetQAStats は指定ユーザーのQ&A活動集計統計を返す。
func (r *QAStatsRepository) GetQAStats(userID uint) (*model.QAStats, error) {
	var stats model.QAStats

	// 質問数
	if err := r.db.Model(&model.Question{}).Where("user_id = ?", userID).Count(&stats.TotalQuestions).Error; err != nil {
		return nil, err
	}

	// 回答数
	if err := r.db.Model(&model.Answer{}).Where("user_id = ?", userID).Count(&stats.TotalAnswers).Error; err != nil {
		return nil, err
	}

	// ベストアンサー数
	if err := r.db.Model(&model.Answer{}).Where("user_id = ? AND is_best = ?", userID, true).Count(&stats.BestAnswerCount).Error; err != nil {
		return nil, err
	}

	// 受け取った投票総数（質問のvote_count + 回答のvote_countの合計）
	var questionVotes *int64
	if err := r.db.Model(&model.Question{}).Where("user_id = ?", userID).Select("SUM(vote_count)").Scan(&questionVotes).Error; err != nil {
		return nil, err
	}
	var answerVotes *int64
	if err := r.db.Model(&model.Answer{}).Where("user_id = ?", userID).Select("SUM(vote_count)").Scan(&answerVotes).Error; err != nil {
		return nil, err
	}
	if questionVotes != nil {
		stats.TotalVotesReceived += *questionVotes
	}
	if answerVotes != nil {
		stats.TotalVotesReceived += *answerVotes
	}

	return &stats, nil
}

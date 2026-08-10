package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// qaStatsRepository は [repository.QAStatsRepository] の GORM 実装。
type qaStatsRepository struct {
	db *gorm.DB
}

// NewQAStatsRepository は QAStatsRepository の GORM 実装を返す。
func NewQAStatsRepository(db *gorm.DB) repository.QAStatsRepository {
	return &qaStatsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.QAStatsRepository = (*qaStatsRepository)(nil)

// GetQAStats は指定ユーザーの Q&A 活動集計統計を返す。
func (r *qaStatsRepository) GetQAStats(ctx context.Context, userID uint) (*model.QAStats, error) {
	db := r.db.WithContext(ctx)
	var stats model.QAStats

	// 質問数
	if err := db.Model(&model.Question{}).Where("user_id = ?", userID).Count(&stats.TotalQuestions).Error; err != nil {
		return nil, err
	}

	// 回答数
	if err := db.Model(&model.Answer{}).Where("user_id = ?", userID).Count(&stats.TotalAnswers).Error; err != nil {
		return nil, err
	}

	// ベストアンサー数
	if err := db.Model(&model.Answer{}).Where("user_id = ? AND is_best = ?", userID, true).Count(&stats.BestAnswerCount).Error; err != nil {
		return nil, err
	}

	// 受け取った投票総数（質問の vote_count + 回答の vote_count の合計）
	var questionVotes *int64
	if err := db.Model(&model.Question{}).Where("user_id = ?", userID).Select("SUM(vote_count)").Scan(&questionVotes).Error; err != nil {
		return nil, err
	}
	var answerVotes *int64
	if err := db.Model(&model.Answer{}).Where("user_id = ?", userID).Select("SUM(vote_count)").Scan(&answerVotes).Error; err != nil {
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

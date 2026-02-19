package repository

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"gorm.io/gorm"
)

// CodeSnippetStatsRepository はユーザーコードスニペット活動集計統計の取得を担当するリポジトリ実装。
type CodeSnippetStatsRepository struct {
	db *gorm.DB
}

// NewCodeSnippetStatsRepository は新しいCodeSnippetStatsRepositoryインスタンスを生成する。
func NewCodeSnippetStatsRepository(db *gorm.DB) *CodeSnippetStatsRepository {
	return &CodeSnippetStatsRepository{db: db}
}

// GetCodeSnippetStats は指定ユーザーのコードスニペット活動集計統計を返す。
func (r *CodeSnippetStatsRepository) GetCodeSnippetStats(userID uint) (*model.CodeSnippetStats, error) {
	var stats model.CodeSnippetStats

	// 総スニペット数
	if err := r.db.Model(&model.CodeSnippet{}).Where("user_id = ?", userID).Count(&stats.TotalSnippets).Error; err != nil {
		return nil, err
	}

	// コメント総数（SUM of comment_count）
	var totalComments *int64
	if err := r.db.Model(&model.CodeSnippet{}).Where("user_id = ?", userID).Select("SUM(comment_count)").Scan(&totalComments).Error; err != nil {
		return nil, err
	}
	if totalComments != nil {
		stats.TotalComments = *totalComments
	}

	// 使用言語数（DISTINCT language）
	if err := r.db.Model(&model.CodeSnippet{}).Where("user_id = ?", userID).Distinct("language").Count(&stats.LanguageCount).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

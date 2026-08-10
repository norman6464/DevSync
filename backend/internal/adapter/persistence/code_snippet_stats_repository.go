package persistence

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
	"gorm.io/gorm"
)

// codeSnippetStatsRepository は [repository.CodeSnippetStatsRepository] の GORM 実装。
type codeSnippetStatsRepository struct {
	db *gorm.DB
}

// NewCodeSnippetStatsRepository は CodeSnippetStatsRepository の GORM 実装を返す。
func NewCodeSnippetStatsRepository(db *gorm.DB) repository.CodeSnippetStatsRepository {
	return &codeSnippetStatsRepository{db: db}
}

// コンパイル時に port を満たすことを保証する（メソッド追加漏れをビルドで検出）。
var _ repository.CodeSnippetStatsRepository = (*codeSnippetStatsRepository)(nil)

// GetCodeSnippetStats は指定ユーザーのコードスニペット活動集計統計を返す。
func (r *codeSnippetStatsRepository) GetCodeSnippetStats(ctx context.Context, userID uint) (*model.CodeSnippetStats, error) {
	db := r.db.WithContext(ctx)
	var stats model.CodeSnippetStats

	// 総スニペット数
	if err := db.Model(&model.CodeSnippet{}).Where("user_id = ?", userID).Count(&stats.TotalSnippets).Error; err != nil {
		return nil, err
	}

	// コメント総数（SUM of comment_count）
	var totalComments *int64
	if err := db.Model(&model.CodeSnippet{}).Where("user_id = ?", userID).Select("SUM(comment_count)").Scan(&totalComments).Error; err != nil {
		return nil, err
	}
	if totalComments != nil {
		stats.TotalComments = *totalComments
	}

	// 使用言語数（DISTINCT language）
	if err := db.Model(&model.CodeSnippet{}).Where("user_id = ?", userID).Distinct("language").Count(&stats.LanguageCount).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

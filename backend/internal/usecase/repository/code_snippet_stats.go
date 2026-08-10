package repository

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
)

// CodeSnippetStatsRepository はユーザーコードスニペット活動集計統計の取得に対する、usecase 側が要求する契約。
type CodeSnippetStatsRepository interface {
	GetCodeSnippetStats(ctx context.Context, userID uint) (*model.CodeSnippetStats, error)
}

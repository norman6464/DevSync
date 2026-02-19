package service

import (
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// CodeSnippetStatsService はユーザーコードスニペット活動集計統計のビジネスロジックを提供する。
type CodeSnippetStatsService struct {
	repo repository.CodeSnippetStatsRepositoryInterface
}

// NewCodeSnippetStatsService は新しいCodeSnippetStatsServiceインスタンスを生成する。
func NewCodeSnippetStatsService(repo repository.CodeSnippetStatsRepositoryInterface) *CodeSnippetStatsService {
	return &CodeSnippetStatsService{repo: repo}
}

// GetCodeSnippetStats は指定ユーザーのコードスニペット活動集計統計を取得する。
func (s *CodeSnippetStatsService) GetCodeSnippetStats(userID uint) (*model.CodeSnippetStats, error) {
	if userID == 0 {
		return nil, domain.NewError(domain.ErrCodeBadRequest, "userIDは必須です", nil)
	}
	return s.repo.GetCodeSnippetStats(userID)
}

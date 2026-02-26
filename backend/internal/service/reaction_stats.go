package service

import (
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// ReactionStatsService はユーザーリアクション集計統計のビジネスロジックを提供する。
type ReactionStatsService struct {
	repo repository.ReactionStatsRepositoryInterface
}

// NewReactionStatsService は新しいReactionStatsServiceインスタンスを生成する。
func NewReactionStatsService(repo repository.ReactionStatsRepositoryInterface) *ReactionStatsService {
	return &ReactionStatsService{repo: repo}
}

// GetReactionStats は指定ユーザーのリアクション集計統計を取得する。
func (s *ReactionStatsService) GetReactionStats(userID uint) (*model.ReactionStats, error) {
	if err := validateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}
	return s.repo.GetReactionStats(userID)
}

// GetReactionSummary は指定ユーザーのリアクションサマリー（絵文字別集計＋トップ投稿）を取得する。
func (s *ReactionStatsService) GetReactionSummary(userID uint) (*model.ReactionSummary, error) {
	if err := validateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}

	emojiCounts, err := s.repo.GetEmojiBreakdown(userID)
	if err != nil {
		return nil, err
	}

	topPosts, err := s.repo.GetTopReactedPosts(userID, 5)
	if err != nil {
		return nil, err
	}

	total := 0
	for _, ec := range emojiCounts {
		total += ec.Count
	}

	return &model.ReactionSummary{
		EmojiCounts:    emojiCounts,
		TopPosts:       topPosts,
		TotalReactions: total,
	}, nil
}

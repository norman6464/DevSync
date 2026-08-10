package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// topReactedPostsLimit はサマリーに載せるトップ投稿の件数。
const topReactedPostsLimit = 5

// GetReactionSummaryUseCase は指定ユーザーのリアクションサマリー（絵文字別集計＋トップ投稿）を取得する。
type GetReactionSummaryUseCase struct {
	stats repository.ReactionStatsRepository
}

// NewGetReactionSummaryUseCase は GetReactionSummaryUseCase を生成する。
func NewGetReactionSummaryUseCase(stats repository.ReactionStatsRepository) *GetReactionSummaryUseCase {
	return &GetReactionSummaryUseCase{stats: stats}
}

// Execute はユーザー ID を検証し、絵文字別集計・トップ投稿・合計値を組み立てて返す。
func (uc *GetReactionSummaryUseCase) Execute(ctx context.Context, userID uint) (*model.ReactionSummary, error) {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return nil, err
	}

	emojiCounts, err := uc.stats.GetEmojiBreakdown(ctx, userID)
	if err != nil {
		return nil, err
	}

	topPosts, err := uc.stats.GetTopReactedPosts(ctx, userID, topReactedPostsLimit)
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

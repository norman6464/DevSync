package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/constants"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// reactionBatchMaxPosts は一括取得で受け付ける投稿数の上限。
const reactionBatchMaxPosts = 50

// AddPostReactionUseCase は投稿にリアクションを追加する。
type AddPostReactionUseCase struct {
	reactions repository.PostReactionRepository
	posts     repository.PostAuthorReader
}

// NewAddPostReactionUseCase は AddPostReactionUseCase を生成する。
func NewAddPostReactionUseCase(
	reactions repository.PostReactionRepository,
	posts repository.PostAuthorReader,
) *AddPostReactionUseCase {
	return &AddPostReactionUseCase{reactions: reactions, posts: posts}
}

// Execute は絵文字と投稿者を検証したうえでリアクションを追加する。
func (uc *AddPostReactionUseCase) Execute(ctx context.Context, userID, postID uint, emoji string) error {
	if err := validateReactionEmoji(emoji); err != nil {
		return err
	}
	if err := ensureNotOwnPost(ctx, uc.posts, userID, postID); err != nil {
		return err
	}
	return uc.reactions.AddReaction(ctx, userID, postID, emoji)
}

// RemovePostReactionUseCase は投稿のリアクションを削除する。
type RemovePostReactionUseCase struct {
	reactions repository.PostReactionRepository
	posts     repository.PostAuthorReader
}

// NewRemovePostReactionUseCase は RemovePostReactionUseCase を生成する。
func NewRemovePostReactionUseCase(
	reactions repository.PostReactionRepository,
	posts repository.PostAuthorReader,
) *RemovePostReactionUseCase {
	return &RemovePostReactionUseCase{reactions: reactions, posts: posts}
}

// Execute は絵文字と投稿者を検証したうえでリアクションを削除する。
// 自分の投稿にはそもそもリアクションできないため、削除も同じ条件で弾く。
func (uc *RemovePostReactionUseCase) Execute(ctx context.Context, userID, postID uint, emoji string) error {
	if err := validateReactionEmoji(emoji); err != nil {
		return err
	}
	if err := ensureNotOwnPost(ctx, uc.posts, userID, postID); err != nil {
		return err
	}
	return uc.reactions.RemoveReaction(ctx, userID, postID, emoji)
}

// validateReactionEmoji は許可された絵文字かどうかを検証する。
func validateReactionEmoji(emoji string) error {
	if !constants.IsAllowedReactionEmoji(emoji) {
		return domain.NewError(domain.ErrCodeBadRequest, "許可されていない絵文字です: "+emoji, nil)
	}
	return nil
}

// ensureNotOwnPost は対象の投稿が自分のものでないことを検証する。
// 投稿が見つからない場合は 404、自分の投稿の場合は 403 を返す。
func ensureNotOwnPost(ctx context.Context, posts repository.PostAuthorReader, userID, postID uint) error {
	authorID, err := posts.FindAuthorID(ctx, postID)
	if err != nil || authorID == 0 {
		return domain.ErrNotFound
	}
	if authorID == userID {
		return domain.ErrForbidden
	}
	return nil
}

// GetPostReactionsUseCase は投稿のリアクション集計と自分のリアクションを返す。
type GetPostReactionsUseCase struct {
	reactions repository.PostReactionRepository
}

// NewGetPostReactionsUseCase は GetPostReactionsUseCase を生成する。
func NewGetPostReactionsUseCase(reactions repository.PostReactionRepository) *GetPostReactionsUseCase {
	return &GetPostReactionsUseCase{reactions: reactions}
}

// Execute は投稿のリアクション集計と、指定ユーザーが付けた絵文字を返す。
// どちらも nil にはせず空スライスへ正規化する。
func (uc *GetPostReactionsUseCase) Execute(ctx context.Context, userID, postID uint) ([]model.ReactionCount, []string, error) {
	reactions, err := uc.reactions.GetReactionsByPostID(ctx, postID)
	if err != nil {
		return nil, nil, err
	}

	userReactions, err := uc.reactions.GetUserReactions(ctx, userID, postID)
	if err != nil {
		return nil, nil, err
	}

	if reactions == nil {
		reactions = []model.ReactionCount{}
	}
	if userReactions == nil {
		userReactions = []string{}
	}
	return reactions, userReactions, nil
}

// GetPostReactionsBatchUseCase は複数投稿のリアクション情報を一括取得する。
type GetPostReactionsBatchUseCase struct {
	reactions repository.PostReactionRepository
}

// NewGetPostReactionsBatchUseCase は GetPostReactionsBatchUseCase を生成する。
func NewGetPostReactionsBatchUseCase(reactions repository.PostReactionRepository) *GetPostReactionsBatchUseCase {
	return &GetPostReactionsBatchUseCase{reactions: reactions}
}

// Execute は指定した投稿群のリアクション集計とユーザーのリアクションを返す。
// 一度に扱えるのは 50 件までで、リクエストされた全 ID に対してエントリを保証する。
func (uc *GetPostReactionsBatchUseCase) Execute(ctx context.Context, userID uint, postIDs []uint) (map[uint][]model.ReactionCount, map[uint][]string, error) {
	if len(postIDs) == 0 {
		return map[uint][]model.ReactionCount{}, map[uint][]string{}, nil
	}
	if len(postIDs) > reactionBatchMaxPosts {
		return nil, nil, domain.NewError(domain.ErrCodeBadRequest, "一度に取得できる投稿は50件までです", nil)
	}

	reactions, err := uc.reactions.GetReactionsBatch(ctx, postIDs)
	if err != nil {
		return nil, nil, err
	}

	userReactions, err := uc.reactions.GetUserReactionsBatch(ctx, userID, postIDs)
	if err != nil {
		return nil, nil, err
	}

	NormalizeReactionMaps(reactions, userReactions, postIDs)
	return reactions, userReactions, nil
}

// NormalizeReactionMaps はリアクションのマップを postIDs に対して正規化する純粋関数。
// nil スライスを空スライスに変換し、全 postID にエントリを保証する。
func NormalizeReactionMaps(reactions map[uint][]model.ReactionCount, userReactions map[uint][]string, postIDs []uint) {
	for _, id := range postIDs {
		if reactions[id] == nil {
			reactions[id] = []model.ReactionCount{}
		}
		if userReactions[id] == nil {
			userReactions[id] = []string{}
		}
	}
}

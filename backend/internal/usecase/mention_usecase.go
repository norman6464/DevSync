package usecase

import (
	"context"
	"regexp"
	"strings"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// mentionRegex はテキスト中の @username パターンを抽出する正規表現。
// メールアドレスを除外するため、@ の前が英数字でないことを確認する。
var mentionRegex = regexp.MustCompile(`(?:^|[^a-zA-Z0-9])@([a-zA-Z0-9_-]+)`)

// ParseMentions はテキストから @username のユーザー名一覧を重複なしで抽出する。
func ParseMentions(text string) []string {
	matches := mentionRegex.FindAllStringSubmatch(text, -1)
	seen := make(map[string]bool)
	var usernames []string
	for _, match := range matches {
		username := strings.ToLower(match[1])
		if !seen[username] {
			seen[username] = true
			usernames = append(usernames, username)
		}
	}
	return usernames
}

// ProcessMentionsInput はメンション処理の入力。
// PostID / CommentID はメンションの記録先で、どちらか一方を指定する。
// NotifyPostID は通知から辿らせる投稿。コメント由来のメンションでも投稿へ飛ばすために使う。
type ProcessMentionsInput struct {
	ActorID      uint
	Text         string
	PostID       *uint
	CommentID    *uint
	NotifyPostID *uint
}

// ProcessMentionsUseCase はテキスト中のメンションを解決し、記録と通知を行う。
type ProcessMentionsUseCase struct {
	mentions      repository.MentionRepository
	users         repository.UsernameLookup
	notifications repository.NotificationCreator
}

// NewProcessMentionsUseCase は ProcessMentionsUseCase を生成する。
func NewProcessMentionsUseCase(
	mentions repository.MentionRepository,
	users repository.UsernameLookup,
	notifications repository.NotificationCreator,
) *ProcessMentionsUseCase {
	return &ProcessMentionsUseCase{mentions: mentions, users: users, notifications: notifications}
}

// Execute は @username を解決してメンションを作成し、相手へ通知する。
// 存在しないユーザーと自分自身へのメンションはスキップし、通知の失敗は無視する。
//
// 同じ投稿・コメントで既にメンション済みのユーザーは作り直さない。
// 本文を編集するたびに同じ相手へ通知が飛ぶのを避けるため、追加された分だけを対象にする。
func (uc *ProcessMentionsUseCase) Execute(ctx context.Context, in ProcessMentionsInput) error {
	usernames := ParseMentions(in.Text)
	if len(usernames) == 0 {
		return nil
	}

	mentioned, err := uc.alreadyMentioned(ctx, in)
	if err != nil {
		return err
	}

	for _, username := range usernames {
		user, err := uc.users.FindByUsername(ctx, username)
		if err != nil || user == nil {
			continue
		}
		if user.ID == in.ActorID || mentioned[user.ID] {
			continue
		}
		mentioned[user.ID] = true

		created, err := uc.mentions.Create(ctx, &model.Mention{
			UserID:    user.ID,
			ActorID:   in.ActorID,
			PostID:    in.PostID,
			CommentID: in.CommentID,
		})
		if err != nil {
			return err
		}
		// 同時実行で先を越された場合は通知も送らない（DB の索引が二重作成を防ぐ）
		if !created {
			continue
		}

		_ = uc.notifications.Create(ctx, &model.Notification{
			UserID:  user.ID,
			ActorID: in.ActorID,
			Type:    model.NotificationTypeMention,
			PostID:  in.NotifyPostID,
		})
	}
	return nil
}

// alreadyMentioned は同じ記録先で既にメンション済みのユーザー ID を返す。
func (uc *ProcessMentionsUseCase) alreadyMentioned(ctx context.Context, in ProcessMentionsInput) (map[uint]bool, error) {
	var existing []model.Mention
	var err error
	switch {
	case in.CommentID != nil:
		existing, err = uc.mentions.FindByCommentID(ctx, *in.CommentID)
	case in.PostID != nil:
		existing, err = uc.mentions.FindByPostID(ctx, *in.PostID)
	default:
		return map[uint]bool{}, nil
	}
	if err != nil {
		return nil, err
	}

	mentioned := make(map[uint]bool, len(existing))
	for _, m := range existing {
		mentioned[m.UserID] = true
	}
	return mentioned, nil
}

// ListUserMentionsUseCase は指定ユーザー宛のメンション一覧を返す。
type ListUserMentionsUseCase struct {
	mentions repository.MentionRepository
}

// NewListUserMentionsUseCase は ListUserMentionsUseCase を生成する。
func NewListUserMentionsUseCase(mentions repository.MentionRepository) *ListUserMentionsUseCase {
	return &ListUserMentionsUseCase{mentions: mentions}
}

// Execute は自分宛のメンションを新しい順に返す。
func (uc *ListUserMentionsUseCase) Execute(ctx context.Context, userID uint, page, limit int) ([]model.Mention, error) {
	return uc.mentions.FindByUserID(ctx, userID, page, limit)
}

// ListPostMentionsUseCase は投稿に紐づくメンション一覧を返す。
type ListPostMentionsUseCase struct {
	mentions repository.MentionRepository
}

// NewListPostMentionsUseCase は ListPostMentionsUseCase を生成する。
func NewListPostMentionsUseCase(mentions repository.MentionRepository) *ListPostMentionsUseCase {
	return &ListPostMentionsUseCase{mentions: mentions}
}

// Execute は指定投稿に紐づくメンションを返す。
func (uc *ListPostMentionsUseCase) Execute(ctx context.Context, postID uint) ([]model.Mention, error) {
	return uc.mentions.FindByPostID(ctx, postID)
}

// DeletePostMentionsUseCase は投稿に紐づくメンションを削除する。
type DeletePostMentionsUseCase struct {
	mentions repository.MentionRepository
}

// NewDeletePostMentionsUseCase は DeletePostMentionsUseCase を生成する。
func NewDeletePostMentionsUseCase(mentions repository.MentionRepository) *DeletePostMentionsUseCase {
	return &DeletePostMentionsUseCase{mentions: mentions}
}

// Execute は指定投稿に紐づくメンションをすべて削除する。
func (uc *DeletePostMentionsUseCase) Execute(ctx context.Context, postID uint) error {
	return uc.mentions.DeleteByPostID(ctx, postID)
}

// DeleteCommentMentionsUseCase はコメントに紐づくメンションを削除する。
type DeleteCommentMentionsUseCase struct {
	mentions repository.MentionRepository
}

// NewDeleteCommentMentionsUseCase は DeleteCommentMentionsUseCase を生成する。
func NewDeleteCommentMentionsUseCase(mentions repository.MentionRepository) *DeleteCommentMentionsUseCase {
	return &DeleteCommentMentionsUseCase{mentions: mentions}
}

// Execute は指定コメントに紐づくメンションをすべて削除する。
func (uc *DeleteCommentMentionsUseCase) Execute(ctx context.Context, commentID uint) error {
	return uc.mentions.DeleteByCommentID(ctx, commentID)
}

package service

import (
	"regexp"
	"strings"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/repository"
)

// mentionRegex はテキスト中の @username パターンを抽出する正規表現。
// メールアドレスを除外するため、@の前が英数字でないことを確認する。
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

// MentionService はメンションのビジネスロジックを提供する。
type MentionService struct {
	repo                repository.MentionRepositoryInterface
	userRepo            repository.UserRepositoryInterface
	notificationService NotificationServiceInterface
}

// NewMentionService は新しいMentionServiceインスタンスを生成する。
func NewMentionService(
	repo repository.MentionRepositoryInterface,
	userRepo repository.UserRepositoryInterface,
	notificationService NotificationServiceInterface,
) *MentionService {
	return &MentionService{
		repo:                repo,
		userRepo:            userRepo,
		notificationService: notificationService,
	}
}

// ProcessMentions はテキストから @username を解析し、メンションレコード作成と通知送信を行う。
// 存在しないユーザーや自分自身へのメンションはスキップする。
func (s *MentionService) ProcessMentions(actorID uint, text string, postID *uint, commentID *uint) error {
	usernames := ParseMentions(text)
	if len(usernames) == 0 {
		return nil
	}

	for _, username := range usernames {
		user, err := s.userRepo.FindByUsername(username)
		if err != nil {
			// ユーザーが存在しない場合はスキップ
			continue
		}

		// 自分自身へのメンションはスキップ
		if user.ID == actorID {
			continue
		}

		mention := &model.Mention{
			UserID:    user.ID,
			ActorID:   actorID,
			PostID:    postID,
			CommentID: commentID,
		}
		if err := s.repo.Create(mention); err != nil {
			return err
		}

		// メンション通知を送信
		notification := &model.Notification{
			UserID:  user.ID,
			ActorID: actorID,
			Type:    model.NotificationTypeMention,
		}
		if postID != nil {
			notification.PostID = postID
		}
		_ = s.notificationService.CreateNotification(notification)
	}

	return nil
}

// GetMentionsByUserID はユーザーへのメンション一覧を取得する。
func (s *MentionService) GetMentionsByUserID(userID uint, page, limit int) ([]model.Mention, error) {
	return s.repo.FindByUserID(userID, page, limit)
}

// GetMentionsByPostID は投稿に関連するメンション一覧を取得する。
func (s *MentionService) GetMentionsByPostID(postID uint) ([]model.Mention, error) {
	return s.repo.FindByPostID(postID)
}

// DeleteMentionsByPostID は投稿に関連するメンションを全て削除する。
func (s *MentionService) DeleteMentionsByPostID(postID uint) error {
	return s.repo.DeleteByPostID(postID)
}

// DeleteMentionsByCommentID はコメントに関連するメンションを全て削除する。
func (s *MentionService) DeleteMentionsByCommentID(commentID uint) error {
	return s.repo.DeleteByCommentID(commentID)
}

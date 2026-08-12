package usecase

import (
	"context"
	"strings"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/domain/validator"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// CreatePostCommentUseCase は投稿にコメント（または返信）を作成する。
type CreatePostCommentUseCase struct {
	comments repository.PostCommentRepository
}

// NewCreatePostCommentUseCase は CreatePostCommentUseCase を生成する。
func NewCreatePostCommentUseCase(comments repository.PostCommentRepository) *CreatePostCommentUseCase {
	return &CreatePostCommentUseCase{comments: comments}
}

// Execute は本文を検証し、返信の場合は親コメントの存在と階層を確認したうえでコメントを作成する。
// 返信への返信は許可しない（1 階層まで）。
func (uc *CreatePostCommentUseCase) Execute(ctx context.Context, userID, postID uint, content string, parentID *uint) (*model.Comment, error) {
	content = strings.TrimSpace(content)
	if err := validator.NewPostValidator().ValidateComment(content); err != nil {
		return nil, err
	}

	if parentID != nil {
		parent, err := uc.comments.FindCommentByID(ctx, *parentID)
		if err != nil {
			return nil, domain.NewError(domain.ErrCodeNotFound, "親コメントが見つかりません", err)
		}
		if parent.PostID != postID {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "親コメントが別の投稿に属しています", nil)
		}
		if parent.ParentID != nil {
			return nil, domain.NewError(domain.ErrCodeBadRequest, "返信への返信はできません", nil)
		}
	}

	comment := &model.Comment{UserID: userID, PostID: postID, Content: content, ParentID: parentID}
	if err := uc.comments.Create(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

// ListPostCommentsUseCase は投稿のコメント一覧を取得する。
type ListPostCommentsUseCase struct {
	comments repository.PostCommentRepository
}

// NewListPostCommentsUseCase は ListPostCommentsUseCase を生成する。
func NewListPostCommentsUseCase(comments repository.PostCommentRepository) *ListPostCommentsUseCase {
	return &ListPostCommentsUseCase{comments: comments}
}

// Execute は指定投稿のトップレベルコメントを返信付きで返す。
func (uc *ListPostCommentsUseCase) Execute(ctx context.Context, postID uint) ([]model.Comment, error) {
	return uc.comments.ListByPostID(ctx, postID)
}

// ListCommentRepliesUseCase はコメントへの返信一覧を取得する。
type ListCommentRepliesUseCase struct {
	comments repository.PostCommentRepository
}

// NewListCommentRepliesUseCase は ListCommentRepliesUseCase を生成する。
func NewListCommentRepliesUseCase(comments repository.PostCommentRepository) *ListCommentRepliesUseCase {
	return &ListCommentRepliesUseCase{comments: comments}
}

// Execute は指定コメントへの返信を返す。
func (uc *ListCommentRepliesUseCase) Execute(ctx context.Context, parentID uint) ([]model.Comment, error) {
	return uc.comments.ListReplies(ctx, parentID)
}

// EditPostCommentUseCase はコメント本文を更新する。
type EditPostCommentUseCase struct {
	comments repository.PostCommentRepository
}

// NewEditPostCommentUseCase は EditPostCommentUseCase を生成する。
func NewEditPostCommentUseCase(comments repository.PostCommentRepository) *EditPostCommentUseCase {
	return &EditPostCommentUseCase{comments: comments}
}

// Execute は所有権と本文を検証したうえでコメントを更新する。
func (uc *EditPostCommentUseCase) Execute(ctx context.Context, id, userID uint, content string) (*model.Comment, error) {
	comment, err := ensureCommentOwner(ctx, uc.comments, id, userID)
	if err != nil {
		return nil, err
	}

	content = strings.TrimSpace(content)
	if err := validator.NewPostValidator().ValidateComment(content); err != nil {
		return nil, err
	}

	comment.Content = content
	if err := uc.comments.Update(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

// DeletePostCommentUseCase はコメントを削除する。
type DeletePostCommentUseCase struct {
	comments repository.PostCommentRepository
}

// NewDeletePostCommentUseCase は DeletePostCommentUseCase を生成する。
func NewDeletePostCommentUseCase(comments repository.PostCommentRepository) *DeletePostCommentUseCase {
	return &DeletePostCommentUseCase{comments: comments}
}

// Execute は所有権を検証したうえでコメントを削除する。
func (uc *DeletePostCommentUseCase) Execute(ctx context.Context, id, userID uint) error {
	if _, err := ensureCommentOwner(ctx, uc.comments, id, userID); err != nil {
		return err
	}
	return uc.comments.Delete(ctx, id)
}

// HidePostCommentUseCase はコメントを非表示にする。
type HidePostCommentUseCase struct {
	comments repository.PostCommentRepository
}

// NewHidePostCommentUseCase は HidePostCommentUseCase を生成する。
func NewHidePostCommentUseCase(comments repository.PostCommentRepository) *HidePostCommentUseCase {
	return &HidePostCommentUseCase{comments: comments}
}

// Execute は所有権を検証したうえでコメントを非表示にする。
func (uc *HidePostCommentUseCase) Execute(ctx context.Context, id, userID uint) error {
	return setCommentHidden(ctx, uc.comments, id, userID, true)
}

// UnhidePostCommentUseCase はコメントの非表示を解除する。
type UnhidePostCommentUseCase struct {
	comments repository.PostCommentRepository
}

// NewUnhidePostCommentUseCase は UnhidePostCommentUseCase を生成する。
func NewUnhidePostCommentUseCase(comments repository.PostCommentRepository) *UnhidePostCommentUseCase {
	return &UnhidePostCommentUseCase{comments: comments}
}

// Execute は所有権を検証したうえでコメントの非表示を解除する。
func (uc *UnhidePostCommentUseCase) Execute(ctx context.Context, id, userID uint) error {
	return setCommentHidden(ctx, uc.comments, id, userID, false)
}

// ensureCommentOwner はコメントを取得し、userID が投稿者であることを検証する。
func ensureCommentOwner(ctx context.Context, comments repository.PostCommentRepository, id, userID uint) (*model.Comment, error) {
	return ensureOwner(ctx, comments.FindCommentByID, id, userID, func(c *model.Comment) uint { return c.UserID })
}

// setCommentHidden は所有権を検証したうえで非表示フラグを更新する。Hide / Unhide で共有する。
func setCommentHidden(ctx context.Context, comments repository.PostCommentRepository, id, userID uint, hidden bool) error {
	comment, err := ensureCommentOwner(ctx, comments, id, userID)
	if err != nil {
		return err
	}
	comment.IsHidden = hidden
	return comments.Update(ctx, comment)
}

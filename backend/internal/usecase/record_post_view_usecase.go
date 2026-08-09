package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// RecordPostViewUseCase はユーザーの投稿閲覧を記録する（ユニーク閲覧のみカウント）。
type RecordPostViewUseCase struct {
	views repository.PostViewRepository
}

// NewRecordPostViewUseCase は RecordPostViewUseCase を生成する。
func NewRecordPostViewUseCase(views repository.PostViewRepository) *RecordPostViewUseCase {
	return &RecordPostViewUseCase{views: views}
}

// Execute は閲覧を記録する。既に閲覧済みの場合は何もしない。
func (uc *RecordPostViewUseCase) Execute(ctx context.Context, userID, postID uint) error {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return err
	}
	if err := domain.ValidateRequiredID(postID, "postID"); err != nil {
		return err
	}

	viewed, err := uc.views.HasViewed(ctx, userID, postID)
	if err != nil {
		return err
	}
	if viewed {
		return nil
	}
	return uc.views.RecordView(ctx, &model.PostView{UserID: userID, PostID: postID})
}

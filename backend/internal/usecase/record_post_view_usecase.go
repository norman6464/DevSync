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

// Execute は閲覧を記録する。既に閲覧済みの場合は何もしない（記録の判定と加算は原子的に行う）。
func (uc *RecordPostViewUseCase) Execute(ctx context.Context, userID, postID uint) error {
	if err := domain.ValidateRequiredID(userID, "userID"); err != nil {
		return err
	}
	if err := domain.ValidateRequiredID(postID, "postID"); err != nil {
		return err
	}

	_, err := uc.views.RecordViewIfAbsent(ctx, &model.PostView{UserID: userID, PostID: postID})
	return err
}

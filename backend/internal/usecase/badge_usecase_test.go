package usecase_test

import (
	"context"
	"testing"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestNotifyBadgeEarnedUseCase_Execute は、クライアントから届くbadge_idが
// EvaluateBadgesの既知バッジID集合に含まれない場合に、通知を作らずエラーを返すことを
// 検証する（DEVSYNC-159で塞いだ「クライアント由来badge_idの無検証保存」の回帰確認）。
func TestNotifyBadgeEarnedUseCase_Execute(t *testing.T) {
	t.Run("既知のバッジIDなら通知を作る", func(t *testing.T) {
		notifications := new(mockNotificationCreatorPort)
		notifications.On("Create", mock.Anything, mock.MatchedBy(func(n *model.Notification) bool {
			return n.UserID == 1 && n.ActorID == 1 && n.Type == model.NotificationTypeBadge && *n.BadgeID == "first-commit"
		})).Return(nil)
		uc := usecase.NewNotifyBadgeEarnedUseCase(notifications)

		err := uc.Execute(context.Background(), 1, "first-commit")

		assert.NoError(t, err)
		notifications.AssertExpectations(t)
	})

	t.Run("未知のバッジIDは通知を作らずエラーを返す", func(t *testing.T) {
		notifications := new(mockNotificationCreatorPort)
		uc := usecase.NewNotifyBadgeEarnedUseCase(notifications)

		err := uc.Execute(context.Background(), 1, "totally-made-up-badge")

		domainErr := domain.GetDomainError(err)
		if assert.NotNil(t, domainErr, "DomainError であること") {
			assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
		}
		notifications.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})
}

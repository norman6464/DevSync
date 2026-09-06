package usecase

import (
	"context"

	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase/repository"
)

// knownNotificationVerbs はnotifications.typeが取りうる既知の値の一覧。
// notification_verbsへの起動時シードに使う（notifications.typeのFK制約化、DEVSYNC-159）。
var knownNotificationVerbs = []string{
	string(model.NotificationTypePost),
	string(model.NotificationTypeMessage),
	string(model.NotificationTypeLike),
	string(model.NotificationTypeComment),
	string(model.NotificationTypeFollow),
	string(model.NotificationTypeAnswer),
	string(model.NotificationTypeBadge),
	string(model.NotificationTypeLevelUp),
	string(model.NotificationTypeMention),
}

// SeedNotificationVerbsUseCase はnotification_verbsへ既知の通知種別コードを登録する。
// notifications.typeがFK参照するため、通知作成より前に完了している必要がある。
type SeedNotificationVerbsUseCase struct {
	verbs repository.NotificationVerbRepository
}

// NewSeedNotificationVerbsUseCase は SeedNotificationVerbsUseCase を生成する。
func NewSeedNotificationVerbsUseCase(verbs repository.NotificationVerbRepository) *SeedNotificationVerbsUseCase {
	return &SeedNotificationVerbsUseCase{verbs: verbs}
}

// Execute は既知の通知種別コードをまとめて登録する（冪等）。
func (uc *SeedNotificationVerbsUseCase) Execute(ctx context.Context) error {
	return uc.verbs.SeedKnownVerbs(ctx, knownNotificationVerbs)
}

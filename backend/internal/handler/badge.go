package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/dto"
	"github.com/norman6464/devsync/backend/internal/usecase"
)

// BadgeHandler はバッジ関連のHTTPハンドラ。
// ユーザーバッジの取得・バッジ獲得通知の作成を処理する。
type BadgeHandler struct {
	getBadges *usecase.GetUserBadgesUseCase
	notify    *usecase.NotifyBadgeEarnedUseCase
}

// NewBadgeHandler は新しいBadgeHandlerインスタンスを生成する。
func NewBadgeHandler(
	getBadges *usecase.GetUserBadgesUseCase,
	notify *usecase.NotifyBadgeEarnedUseCase,
) *BadgeHandler {
	return &BadgeHandler{getBadges: getBadges, notify: notify}
}

// GetUserBadges は指定ユーザーの全バッジを獲得状況付きで返す。
func (h *BadgeHandler) GetUserBadges(c *gin.Context) {
	userID, ok := parseID(c, "userId")
	if !ok {
		return
	}

	badges, err := h.getBadges.Execute(c.Request.Context(), userID)
	if err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, dto.BadgesResponse{Badges: badges})
}

// NotifyBadgeEarned は新しく獲得したバッジの通知を作成する。
func (h *BadgeHandler) NotifyBadgeEarned(c *gin.Context) {
	userID := c.GetUint("userID")

	input := bindJSON[dto.NotifyBadgeEarnedRequest](c)
	if input == nil {
		return
	}

	if err := h.notify.Execute(c.Request.Context(), userID, input.BadgeID); err != nil {
		respondError(c, err)
		return
	}

	respondOK(c, domain.NewMessageResponse("badge notification created"))
}

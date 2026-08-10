package usecase

import "github.com/norman6464/devsync/backend/internal/domain"

// maxLanguageNameLength はランキングで受け付ける言語名の最大長。
const maxLanguageNameLength = 50

// validRankingPeriods はランキングで許可される期間パラメータ。
var validRankingPeriods = map[string]bool{
	"weekly": true, "monthly": true,
}

// validateRankingPeriod は期間パラメータを検証する。
func validateRankingPeriod(period string) error {
	if !validRankingPeriods[period] {
		return domain.NewError(domain.ErrCodeBadRequest, "periodはweekly/monthlyのいずれかを指定してください", nil)
	}
	return nil
}

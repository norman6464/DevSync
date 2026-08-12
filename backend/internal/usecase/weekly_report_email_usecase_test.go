package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/norman6464/devsync/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockEmailSender は usecase/repository.EmailSender のモック。
type mockEmailSender struct{ mock.Mock }

func (m *mockEmailSender) Send(ctx context.Context, to, subject, htmlBody string) error {
	return m.Called(ctx, to, subject, htmlBody).Error(0)
}

// mockWeeklyReportReader は usecase/repository.WeeklyActivityReportReader のモック。
type mockWeeklyReportReader struct{ mock.Mock }

func (m *mockWeeklyReportReader) GetWeeklyReport(ctx context.Context, userID uint) (*model.ActivityReport, error) {
	args := m.Called(ctx, userID)
	r, _ := args.Get(0).(*model.ActivityReport)
	return r, args.Error(1)
}

// weeklyTestReport はテスト用のレポートデータを返す。
func weeklyTestReport() *model.ActivityReport {
	return &model.ActivityReport{
		Period:             model.ReportPeriodWeekly,
		StartDate:          time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
		EndDate:            time.Date(2025, 1, 12, 0, 0, 0, 0, time.UTC),
		UserID:             1,
		TotalContributions: 42,
		PostsCreated:       3,
		CommentsCreated:    7,
		LikesReceived:      15,
		GoalsCompleted:     1,
		GoalsProgress:      65,
		NewFollowers:       2,
		MessagesExchanged:  10,
	}
}

// ============================================================
// 1 ユーザーへの送信
// ============================================================

func TestSendWeeklyReportUseCase(t *testing.T) {
	ctx := context.Background()

	t.Run("メールを送信する", func(t *testing.T) {
		sender := new(mockEmailSender)
		sender.On("Send", mock.Anything, "test@example.com", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

		user := &model.User{Name: "テストユーザー", Email: "test@example.com", EmailLanguage: "ja"}
		require.NoError(t, usecase.NewSendWeeklyReportUseCase(sender, "https://example.com").
			Execute(ctx, user, weeklyTestReport()))
		sender.AssertExpectations(t)
	})

	t.Run("件名は日本語のテキストを使う", func(t *testing.T) {
		sender := new(mockEmailSender)
		sender.On("Send", mock.Anything, mock.Anything, mock.MatchedBy(func(subject string) bool {
			return strings.Contains(subject, "[DevSync]") && strings.Contains(subject, "ウィークリーアクティビティレポート")
		}), mock.Anything).Return(nil)

		user := &model.User{Name: "u", Email: "u@example.com", EmailLanguage: "ja"}
		require.NoError(t, usecase.NewSendWeeklyReportUseCase(sender, "").Execute(ctx, user, weeklyTestReport()))
		sender.AssertExpectations(t)
	})

	t.Run("言語未設定なら日本語で送る", func(t *testing.T) {
		sender := new(mockEmailSender)
		sender.On("Send", mock.Anything, mock.Anything, mock.MatchedBy(func(subject string) bool {
			return strings.Contains(subject, "ウィークリーアクティビティレポート")
		}), mock.Anything).Return(nil)

		user := &model.User{Name: "u", Email: "u@example.com"}
		require.NoError(t, usecase.NewSendWeeklyReportUseCase(sender, "").Execute(ctx, user, weeklyTestReport()))
		sender.AssertExpectations(t)
	})

	t.Run("メールアドレスが空なら 400", func(t *testing.T) {
		sender := new(mockEmailSender)

		err := usecase.NewSendWeeklyReportUseCase(sender, "").
			Execute(ctx, &model.User{Name: "u"}, weeklyTestReport())
		var domainErr *domain.DomainError
		require.ErrorAs(t, err, &domainErr)
		assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
		sender.AssertNotCalled(t, "Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("レポートが nil なら 400", func(t *testing.T) {
		sender := new(mockEmailSender)

		err := usecase.NewSendWeeklyReportUseCase(sender, "").
			Execute(ctx, &model.User{Name: "u", Email: "u@example.com"}, nil)
		var domainErr *domain.DomainError
		require.ErrorAs(t, err, &domainErr)
		assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
	})

	t.Run("送信エラーはそのまま返す", func(t *testing.T) {
		sender := new(mockEmailSender)
		sendErr := errors.New("smtp error")
		sender.On("Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(sendErr)

		err := usecase.NewSendWeeklyReportUseCase(sender, "").
			Execute(ctx, &model.User{Name: "u", Email: "u@example.com"}, weeklyTestReport())
		assert.ErrorIs(t, err, sendErr)
	})
}

// ============================================================
// 一括送信
// ============================================================

func TestSendAllWeeklyReportsUseCase(t *testing.T) {
	ctx := context.Background()

	newUseCase := func(users *mockUserRepo, reports *mockWeeklyReportReader, sender *mockEmailSender) *usecase.SendAllWeeklyReportsUseCase {
		return usecase.NewSendAllWeeklyReportsUseCase(users, reports,
			usecase.NewSendWeeklyReportUseCase(sender, "https://example.com"))
	}

	t.Run("配信を無効にしているユーザーはスキップする", func(t *testing.T) {
		users := new(mockUserRepo)
		reports := new(mockWeeklyReportReader)
		sender := new(mockEmailSender)

		users.On("FindAll", mock.Anything).Return([]model.User{
			{ID: 1, Name: "有効", Email: "a@example.com", EmailWeeklyReport: true},
			{ID: 2, Name: "無効", Email: "b@example.com", EmailWeeklyReport: false},
		}, nil)
		reports.On("GetWeeklyReport", mock.Anything, uint(1)).Return(weeklyTestReport(), nil)
		sender.On("Send", mock.Anything, "a@example.com", mock.Anything, mock.Anything).Return(nil)

		require.NoError(t, newUseCase(users, reports, sender).Execute(ctx))
		reports.AssertNotCalled(t, "GetWeeklyReport", mock.Anything, uint(2))
		sender.AssertNumberOfCalls(t, "Send", 1)
	})

	t.Run("1 ユーザーの送信失敗で全体を止めない", func(t *testing.T) {
		users := new(mockUserRepo)
		reports := new(mockWeeklyReportReader)
		sender := new(mockEmailSender)

		users.On("FindAll", mock.Anything).Return([]model.User{
			{ID: 1, Name: "失敗", Email: "a@example.com", EmailWeeklyReport: true},
			{ID: 2, Name: "成功", Email: "b@example.com", EmailWeeklyReport: true},
		}, nil)
		reports.On("GetWeeklyReport", mock.Anything, uint(1)).Return(weeklyTestReport(), nil)
		reports.On("GetWeeklyReport", mock.Anything, uint(2)).Return(weeklyTestReport(), nil)
		sender.On("Send", mock.Anything, "a@example.com", mock.Anything, mock.Anything).Return(errors.New("smtp error"))
		sender.On("Send", mock.Anything, "b@example.com", mock.Anything, mock.Anything).Return(nil)

		require.NoError(t, newUseCase(users, reports, sender).Execute(ctx))
		sender.AssertNumberOfCalls(t, "Send", 2)
	})

	t.Run("レポート取得に失敗したユーザーはスキップする", func(t *testing.T) {
		users := new(mockUserRepo)
		reports := new(mockWeeklyReportReader)
		sender := new(mockEmailSender)

		users.On("FindAll", mock.Anything).Return([]model.User{
			{ID: 1, Name: "u", Email: "a@example.com", EmailWeeklyReport: true},
		}, nil)
		reports.On("GetWeeklyReport", mock.Anything, uint(1)).Return(nil, errors.New("db error"))

		require.NoError(t, newUseCase(users, reports, sender).Execute(ctx))
		sender.AssertNotCalled(t, "Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("ユーザー一覧の取得に失敗したらエラーを返す", func(t *testing.T) {
		users := new(mockUserRepo)
		reports := new(mockWeeklyReportReader)
		sender := new(mockEmailSender)

		users.On("FindAll", mock.Anything).Return(nil, errors.New("db error"))

		err := newUseCase(users, reports, sender).Execute(ctx)
		var domainErr *domain.DomainError
		require.ErrorAs(t, err, &domainErr)
		assert.Equal(t, domain.ErrCodeDatabase, domainErr.Code)
	})

	t.Run("ユーザーがいなければ何もしない", func(t *testing.T) {
		users := new(mockUserRepo)
		reports := new(mockWeeklyReportReader)
		sender := new(mockEmailSender)

		users.On("FindAll", mock.Anything).Return([]model.User{}, nil)

		require.NoError(t, newUseCase(users, reports, sender).Execute(ctx))
		sender.AssertNotCalled(t, "Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
}

// ============================================================
// HTML レンダリング
// ============================================================

func TestRenderWeeklyReportHTML(t *testing.T) {
	user := &model.User{Name: "テストユーザー", Email: "test@example.com"}
	report := weeklyTestReport()

	t.Run("日本語のテキストと数値が埋め込まれる", func(t *testing.T) {
		html, err := usecase.RenderWeeklyReportHTML(user, report, "ja", "https://example.com")
		require.NoError(t, err)
		assert.Contains(t, html, "テストユーザー")
		assert.Contains(t, html, "42")
		assert.Contains(t, html, "コントリビューション")
		assert.Contains(t, html, "https://example.com")
	})

	t.Run("英語のテキストが使われる", func(t *testing.T) {
		html, err := usecase.RenderWeeklyReportHTML(user, report, "en", "")
		require.NoError(t, err)
		assert.Contains(t, html, "Contributions")
	})

	t.Run("未対応の言語は日本語にフォールバックする", func(t *testing.T) {
		html, err := usecase.RenderWeeklyReportHTML(user, report, "xx", "")
		require.NoError(t, err)
		assert.Contains(t, html, "コントリビューション")
	})
}

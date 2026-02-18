package service

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/norman6464/devsync/backend/internal/domain"
	"github.com/norman6464/devsync/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// newTestWeeklyReportEmailService はWeeklyReportEmailServiceのテスト用インスタンスを生成するヘルパー。
func newTestWeeklyReportEmailService() (*WeeklyReportEmailService, *MockEmailSender, *MockUserRepository, *MockActivityReportRepository) {
	sender := new(MockEmailSender)
	userRepo := new(MockUserRepository)
	reportRepo := new(MockActivityReportRepository)
	reportService := NewActivityReportService(reportRepo)
	svc := NewWeeklyReportEmailService(sender, reportService, userRepo)
	return svc, sender, userRepo, reportRepo
}

// テスト用のレポートデータ生成ヘルパー
func testReport() *model.ActivityReport {
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
// ウィークリーレポートメール送信テスト
// ============================================================

func TestSendWeeklyReport_Success(t *testing.T) {
	svc, sender, _, _ := newTestWeeklyReportEmailService()

	user := &model.User{Name: "テストユーザー", Email: "test@example.com", EmailLanguage: "ja"}
	report := testReport()

	sender.On("Send", "test@example.com", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	err := svc.SendWeeklyReport(user, report)
	assert.NoError(t, err)
	sender.AssertExpectations(t)
}

func TestSendWeeklyReport_EmptyEmail(t *testing.T) {
	svc, _, _, _ := newTestWeeklyReportEmailService()

	user := &model.User{Name: "テストユーザー", Email: ""}
	report := testReport()

	err := svc.SendWeeklyReport(user, report)
	assert.Error(t, err)
	var domainErr *domain.DomainError
	assert.ErrorAs(t, err, &domainErr)
	assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
}

func TestSendWeeklyReport_NilReport(t *testing.T) {
	svc, _, _, _ := newTestWeeklyReportEmailService()

	user := &model.User{Name: "テストユーザー", Email: "test@example.com"}

	err := svc.SendWeeklyReport(user, nil)
	assert.Error(t, err)
	var domainErr *domain.DomainError
	assert.ErrorAs(t, err, &domainErr)
	assert.Equal(t, domain.ErrCodeBadRequest, domainErr.Code)
}

// ============================================================
// 一括送信テスト
// ============================================================

func TestSendAllWeeklyReports_SkipDisabledUsers(t *testing.T) {
	svc, sender, userRepo, reportRepo := newTestWeeklyReportEmailService()

	users := []model.User{
		{Name: "有効ユーザー", Email: "enabled@example.com", EmailWeeklyReport: true, EmailLanguage: "ja"},
		{Name: "無効ユーザー", Email: "disabled@example.com", EmailWeeklyReport: false, EmailLanguage: "ja"},
	}
	users[0].ID = 1
	users[1].ID = 2

	userRepo.On("FindAll").Return(users, nil)
	reportRepo.On("GetWeeklyReport", uint(1)).Return(testReport(), nil)
	// ユーザー2のレポートは呼ばれないはず
	sender.On("Send", "enabled@example.com", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	err := svc.SendAllWeeklyReports()
	assert.NoError(t, err)
	sender.AssertExpectations(t)
	// 無効ユーザーにはメール送信されないことを確認
	sender.AssertNumberOfCalls(t, "Send", 1)
}

func TestSendAllWeeklyReports_ContinueOnError(t *testing.T) {
	svc, sender, userRepo, reportRepo := newTestWeeklyReportEmailService()

	users := []model.User{
		{Name: "ユーザー1", Email: "user1@example.com", EmailWeeklyReport: true, EmailLanguage: "ja"},
		{Name: "ユーザー2", Email: "user2@example.com", EmailWeeklyReport: true, EmailLanguage: "ja"},
	}
	users[0].ID = 1
	users[1].ID = 2

	userRepo.On("FindAll").Return(users, nil)
	reportRepo.On("GetWeeklyReport", uint(1)).Return(testReport(), nil)
	reportRepo.On("GetWeeklyReport", uint(2)).Return(testReport(), nil)

	// ユーザー1のメール送信でエラー発生
	sender.On("Send", "user1@example.com", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(assert.AnError)
	// ユーザー2のメール送信は成功するはず
	sender.On("Send", "user2@example.com", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	err := svc.SendAllWeeklyReports()
	// 一部エラーでもnil返却（ログ出力のみ）
	assert.NoError(t, err)
	// 両方のユーザーに対してSendが呼ばれることを確認（1件エラーでも2件目は継続）
	sender.AssertNumberOfCalls(t, "Send", 2)
}

// ============================================================
// HTMLレンダリングテスト
// ============================================================

func TestRenderWeeklyReportHTML_Success(t *testing.T) {
	svc, _, _, _ := newTestWeeklyReportEmailService()

	user := &model.User{Name: "テストユーザー", Email: "test@example.com"}
	report := testReport()

	html, err := svc.RenderHTML(user, report, "ja")
	assert.NoError(t, err)
	assert.NotEmpty(t, html)
	// HTML内にユーザー名が含まれることを確認
	assert.True(t, strings.Contains(html, "テストユーザー"))
	// HTML内にレポートデータが含まれることを確認
	assert.True(t, strings.Contains(html, "42"))  // TotalContributions
	assert.True(t, strings.Contains(html, "3"))   // PostsCreated
}

// ============================================================
// エッジケーステスト
// ============================================================

func TestSendWeeklyReport_DefaultLanguage(t *testing.T) {
	svc, sender, _, _ := newTestWeeklyReportEmailService()

	// EmailLanguageが空の場合、デフォルトで"ja"が使われる
	user := &model.User{Name: "テストユーザー", Email: "test@example.com", EmailLanguage: ""}
	report := testReport()

	sender.On("Send", "test@example.com", mock.MatchedBy(func(subject string) bool {
		return strings.Contains(subject, "ウィークリー")
	}), mock.AnythingOfType("string")).Return(nil)

	err := svc.SendWeeklyReport(user, report)
	assert.NoError(t, err)
	sender.AssertExpectations(t)
}

func TestSendWeeklyReport_SenderError(t *testing.T) {
	svc, sender, _, _ := newTestWeeklyReportEmailService()

	user := &model.User{Name: "テストユーザー", Email: "test@example.com", EmailLanguage: "ja"}
	report := testReport()

	sender.On("Send", "test@example.com", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(errors.New("smtp error"))

	err := svc.SendWeeklyReport(user, report)
	assert.Error(t, err)
	sender.AssertExpectations(t)
}

func TestSendAllWeeklyReports_FindAllError(t *testing.T) {
	svc, _, userRepo, _ := newTestWeeklyReportEmailService()

	userRepo.On("FindAll").Return([]model.User(nil), errors.New("db error"))

	err := svc.SendAllWeeklyReports()
	assert.Error(t, err)
	userRepo.AssertExpectations(t)
}

func TestSendAllWeeklyReports_EmptyUserList(t *testing.T) {
	svc, _, userRepo, _ := newTestWeeklyReportEmailService()

	userRepo.On("FindAll").Return([]model.User{}, nil)

	err := svc.SendAllWeeklyReports()
	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestRenderWeeklyReportHTML_EnglishLanguage(t *testing.T) {
	svc, _, _, _ := newTestWeeklyReportEmailService()

	user := &model.User{Name: "TestUser", Email: "test@example.com"}
	report := testReport()

	html, err := svc.RenderHTML(user, report, "en")
	assert.NoError(t, err)
	assert.NotEmpty(t, html)
	assert.True(t, strings.Contains(html, "TestUser"))
	// 英語テキストが含まれることを確認
	assert.True(t, strings.Contains(html, "This Week"))
}

func TestRenderWeeklyReportHTML_UnsupportedLanguageFallback(t *testing.T) {
	svc, _, _, _ := newTestWeeklyReportEmailService()

	user := &model.User{Name: "テストユーザー", Email: "test@example.com"}
	report := testReport()

	// 未対応言語は日本語にフォールバック
	html, err := svc.RenderHTML(user, report, "xx")
	assert.NoError(t, err)
	assert.NotEmpty(t, html)
	// 日本語テキストが含まれることを確認
	assert.True(t, strings.Contains(html, "コントリビューション"))
}

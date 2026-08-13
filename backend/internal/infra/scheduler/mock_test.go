package scheduler

import (
	"context"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/mock"
)

// MockCronScheduler は CronScheduler のテスト用モック実装。
type MockCronScheduler struct {
	mock.Mock
}

var _ CronScheduler = (*MockCronScheduler)(nil)

func (m *MockCronScheduler) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	args := m.Called(spec, cmd)
	return args.Get(0).(cron.EntryID), args.Error(1)
}

func (m *MockCronScheduler) Start() {
	m.Called()
}

func (m *MockCronScheduler) Stop() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

// MockWeeklyReportSender は WeeklyReportSender のテスト用モック実装。
type MockWeeklyReportSender struct {
	mock.Mock
}

var _ WeeklyReportSender = (*MockWeeklyReportSender)(nil)

func (m *MockWeeklyReportSender) Execute(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

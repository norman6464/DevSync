package service

import (
	"context"
	"errors"
	"testing"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestScheduler() (*Scheduler, *MockCronScheduler, *MockWeeklyReportSender) {
	cronMock := new(MockCronScheduler)
	senderMock := new(MockWeeklyReportSender)
	s := &Scheduler{
		cron:     cronMock,
		emailSvc: senderMock,
	}
	return s, cronMock, senderMock
}

func TestNewScheduler(t *testing.T) {
	s := NewScheduler(nil)
	assert.NotNil(t, s)
	assert.NotNil(t, s.cron)
}

func TestStart_Success(t *testing.T) {
	s, cronMock, _ := newTestScheduler()

	cronMock.On("AddFunc", "0 9 * * 1", mock.AnythingOfType("func()")).
		Return(cron.EntryID(1), nil)
	cronMock.On("Start").Return()

	s.Start()

	cronMock.AssertExpectations(t)
}

func TestStart_AddFuncError(t *testing.T) {
	s, cronMock, _ := newTestScheduler()

	cronMock.On("AddFunc", "0 9 * * 1", mock.AnythingOfType("func()")).
		Return(cron.EntryID(0), errors.New("invalid cron spec"))

	s.Start()

	cronMock.AssertExpectations(t)
	cronMock.AssertNotCalled(t, "Start")
}

func TestStart_JobExecutes_SendSuccess(t *testing.T) {
	s, cronMock, senderMock := newTestScheduler()

	var capturedJob func()
	cronMock.On("AddFunc", "0 9 * * 1", mock.AnythingOfType("func()")).
		Run(func(args mock.Arguments) {
			capturedJob = args.Get(1).(func())
		}).
		Return(cron.EntryID(1), nil)
	cronMock.On("Start").Return()

	senderMock.On("Execute", mock.Anything).Return(nil)

	s.Start()
	assert.NotNil(t, capturedJob)

	// ジョブを手動実行
	capturedJob()

	senderMock.AssertExpectations(t)
}

func TestStart_JobExecutes_SendError(t *testing.T) {
	s, cronMock, senderMock := newTestScheduler()

	var capturedJob func()
	cronMock.On("AddFunc", "0 9 * * 1", mock.AnythingOfType("func()")).
		Run(func(args mock.Arguments) {
			capturedJob = args.Get(1).(func())
		}).
		Return(cron.EntryID(1), nil)
	cronMock.On("Start").Return()

	senderMock.On("Execute", mock.Anything).Return(errors.New("smtp error"))

	s.Start()
	assert.NotNil(t, capturedJob)

	// ジョブを手動実行（エラーでもpanicしない）
	capturedJob()

	senderMock.AssertExpectations(t)
}

func TestStop(t *testing.T) {
	s, cronMock, _ := newTestScheduler()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 即座にDoneになるcontextを作成
	cronMock.On("Stop").Return(ctx)

	s.Stop()

	cronMock.AssertExpectations(t)
}

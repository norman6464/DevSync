package scheduler

import (
	"context"
	"errors"
	"testing"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestScheduler() (*Scheduler, *MockCronScheduler, *MockWeeklyReportSender, *MockMetricsReconciler) {
	cronMock := new(MockCronScheduler)
	senderMock := new(MockWeeklyReportSender)
	reconcilerMock := new(MockMetricsReconciler)
	s := &Scheduler{
		cron:              cronMock,
		emailSvc:          senderMock,
		metricsReconciler: reconcilerMock,
	}
	return s, cronMock, senderMock, reconcilerMock
}

func TestNew(t *testing.T) {
	s := New(nil, nil)
	assert.NotNil(t, s)
	assert.NotNil(t, s.cron)
}

func TestStart_Success(t *testing.T) {
	s, cronMock, _, _ := newTestScheduler()

	cronMock.On("AddFunc", "0 9 * * 1", mock.AnythingOfType("func()")).
		Return(cron.EntryID(1), nil)
	cronMock.On("AddFunc", "0 3 * * *", mock.AnythingOfType("func()")).
		Return(cron.EntryID(2), nil)
	cronMock.On("Start").Return()

	s.Start()

	cronMock.AssertExpectations(t)
}

// AddFuncが片方失敗しても、もう片方のジョブ登録・cron起動自体は続行する
// （1ジョブの登録失敗が他の正常なジョブまで道連れにしない）。
func TestStart_AddFuncError_StillStartsOtherJob(t *testing.T) {
	s, cronMock, _, _ := newTestScheduler()

	cronMock.On("AddFunc", "0 9 * * 1", mock.AnythingOfType("func()")).
		Return(cron.EntryID(0), errors.New("invalid cron spec"))
	cronMock.On("AddFunc", "0 3 * * *", mock.AnythingOfType("func()")).
		Return(cron.EntryID(2), nil)
	cronMock.On("Start").Return()

	s.Start()

	cronMock.AssertExpectations(t)
}

func TestStart_NilWeeklyReportSender_SkipsThatJob(t *testing.T) {
	s, cronMock, _, _ := newTestScheduler()
	s.emailSvc = nil

	cronMock.On("AddFunc", "0 3 * * *", mock.AnythingOfType("func()")).
		Return(cron.EntryID(2), nil)
	cronMock.On("Start").Return()

	s.Start()

	cronMock.AssertExpectations(t)
	cronMock.AssertNotCalled(t, "AddFunc", "0 9 * * 1", mock.Anything)
}

func TestStart_NilMetricsReconciler_SkipsThatJob(t *testing.T) {
	s, cronMock, _, _ := newTestScheduler()
	s.metricsReconciler = nil

	cronMock.On("AddFunc", "0 9 * * 1", mock.AnythingOfType("func()")).
		Return(cron.EntryID(1), nil)
	cronMock.On("Start").Return()

	s.Start()

	cronMock.AssertExpectations(t)
	cronMock.AssertNotCalled(t, "AddFunc", "0 3 * * *", mock.Anything)
}

func TestStart_WeeklyReportJobExecutes_SendSuccess(t *testing.T) {
	s, cronMock, senderMock, _ := newTestScheduler()

	var capturedJob func()
	cronMock.On("AddFunc", "0 9 * * 1", mock.AnythingOfType("func()")).
		Run(func(args mock.Arguments) {
			capturedJob = args.Get(1).(func())
		}).
		Return(cron.EntryID(1), nil)
	cronMock.On("AddFunc", "0 3 * * *", mock.AnythingOfType("func()")).
		Return(cron.EntryID(2), nil)
	cronMock.On("Start").Return()

	senderMock.On("Execute", mock.Anything).Return(nil)

	s.Start()
	assert.NotNil(t, capturedJob)

	// ジョブを手動実行
	capturedJob()

	senderMock.AssertExpectations(t)
}

func TestStart_WeeklyReportJobExecutes_SendError(t *testing.T) {
	s, cronMock, senderMock, _ := newTestScheduler()

	var capturedJob func()
	cronMock.On("AddFunc", "0 9 * * 1", mock.AnythingOfType("func()")).
		Run(func(args mock.Arguments) {
			capturedJob = args.Get(1).(func())
		}).
		Return(cron.EntryID(1), nil)
	cronMock.On("AddFunc", "0 3 * * *", mock.AnythingOfType("func()")).
		Return(cron.EntryID(2), nil)
	cronMock.On("Start").Return()

	senderMock.On("Execute", mock.Anything).Return(errors.New("smtp error"))

	s.Start()
	assert.NotNil(t, capturedJob)

	// ジョブを手動実行（エラーでもpanicしない）
	capturedJob()

	senderMock.AssertExpectations(t)
}

func TestStart_ReconcileJobExecutes_Success(t *testing.T) {
	s, cronMock, _, reconcilerMock := newTestScheduler()

	var capturedJob func()
	cronMock.On("AddFunc", "0 9 * * 1", mock.AnythingOfType("func()")).
		Return(cron.EntryID(1), nil)
	cronMock.On("AddFunc", "0 3 * * *", mock.AnythingOfType("func()")).
		Run(func(args mock.Arguments) {
			capturedJob = args.Get(1).(func())
		}).
		Return(cron.EntryID(2), nil)
	cronMock.On("Start").Return()

	reconcilerMock.On("Execute", mock.Anything).Return(nil)

	s.Start()
	assert.NotNil(t, capturedJob)

	capturedJob()

	reconcilerMock.AssertExpectations(t)
}

func TestStart_ReconcileJobExecutes_Error(t *testing.T) {
	s, cronMock, _, reconcilerMock := newTestScheduler()

	var capturedJob func()
	cronMock.On("AddFunc", "0 9 * * 1", mock.AnythingOfType("func()")).
		Return(cron.EntryID(1), nil)
	cronMock.On("AddFunc", "0 3 * * *", mock.AnythingOfType("func()")).
		Run(func(args mock.Arguments) {
			capturedJob = args.Get(1).(func())
		}).
		Return(cron.EntryID(2), nil)
	cronMock.On("Start").Return()

	reconcilerMock.On("Execute", mock.Anything).Return(errors.New("db error"))

	s.Start()
	assert.NotNil(t, capturedJob)

	// ジョブを手動実行（エラーでもpanicしない）
	capturedJob()

	reconcilerMock.AssertExpectations(t)
}

func TestStop(t *testing.T) {
	s, cronMock, _, _ := newTestScheduler()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 即座にDoneになるcontextを作成
	cronMock.On("Stop").Return(ctx)

	s.Stop()

	cronMock.AssertExpectations(t)
}

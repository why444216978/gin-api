package cron

import (
	"log"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	lockmock "github.com/why444216978/gin-api/library/lock/mock"
	zapLogger "github.com/why444216978/gin-api/library/logger/zap"
)

func JobFunc() {
	log.Println("JobFunc handle")
}

func TestCron_AddJob(t *testing.T) {
	logger, err := zapLogger.NewLogger()
	assert.Equal(t, err, nil)

	ctl := gomock.NewController(t)
	locker := lockmock.NewMockLocker(ctl)
	locker.EXPECT().Lock(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	locker.EXPECT().Unlock(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(nil)

	cron, err := NewCron("JobFunc", logger, WithLocker(locker))
	assert.Equal(t, err, nil)

	entryID, err := cron.AddJob("*/3 * * * * *", JobFunc)
	assert.Equal(t, err, nil)
	assert.Equal(t, entryID > 0, true)

	cron.Start()
	time.Sleep(time.Second * 9)
	cron.Stop()
}
